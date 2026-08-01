package bootstrap

import (
	"fmt"
	"time"

	"github.com/m87/ctx/core"
	"github.com/m87/nod"
)

type databaseMigration struct {
	version string
	run     func(*nod.Repository) (int64, error)
}

type intervalTimeValue struct {
	NodeID    string     `gorm:"column:node_id"`
	Key       string     `gorm:"column:key"`
	ValueTime *time.Time `gorm:"column:value_time"`
}

type timeColumn struct {
	table  string
	column string
}

type timeColumnValue struct {
	RowID int64      `gorm:"column:row_id"`
	Value *time.Time `gorm:"column:value"`
}

var storedTimeColumns = []timeColumn{
	{table: "node_cores", column: "created_at"},
	{table: "node_cores", column: "updated_at"},
	{table: "node_contents", column: "created_at"},
	{table: "node_contents", column: "updated_at"},
	{table: "tags", column: "created_at"},
	{table: "edge_cores", column: "created_at"},
	{table: "edge_cores", column: "updated_at"},
	{table: "edge_contents", column: "created_at"},
	{table: "edge_contents", column: "updated_at"},
	{table: "edge_kvs", column: "value_time"},
	{table: "nod_properties", column: "updated_at"},
}

func applicationMigrations() []databaseMigration {
	return []databaseMigration{
		{
			version: "0.5.0",
			run: func(repository *nod.Repository) (int64, error) {
				count, err := newContextManager(repository).EnsureDefaultWorkspaceWithResult()
				return int64(count), err
			},
		},
		{
			version: "0.6.0",
			run:     migrateTimesToUTC,
		},
	}
}

func runApplicationMigrations(repository *nod.Repository, currentVersion string) (int64, error) {
	var updated int64
	for _, migration := range applicationMigrations() {
		needsMigration, err := core.DatabaseVersionNeedsMigration(currentVersion, migration.version)
		if err != nil {
			return updated, err
		}
		if !needsMigration {
			continue
		}

		count, err := migration.run(repository)
		if err != nil {
			return updated, fmt.Errorf("migrate database to %s: %w", migration.version, err)
		}
		updated += count
		currentVersion = migration.version
	}
	return updated, nil
}

// migrateTimesToUTC performs a preflight read of every application timestamp
// before mutating the database. The surrounding bootstrap transaction
// guarantees that an invalid value rolls back the whole migration.
func migrateTimesToUTC(repository *nod.Repository) (int64, error) {
	db := repository.DB()

	var values []intervalTimeValue
	if err := db.Table("node_kvs AS kv").
		Select("kv.node_id, kv.key, kv.value_time").
		Joins("JOIN node_cores AS core ON core.id = kv.node_id").
		Where("core.kind = ?", core.IntervalType).
		Where("kv.key IN ?", []string{"start", "end"}).
		Order("kv.node_id, kv.key").
		Scan(&values).Error; err != nil {
		return 0, fmt.Errorf("preflight interval timestamps: %w", err)
	}

	columnValues := make(map[timeColumn][]timeColumnValue, len(storedTimeColumns))
	for _, storedColumn := range storedTimeColumns {
		var storedValues []timeColumnValue
		selectExpression := fmt.Sprintf("rowid AS row_id, %s AS value", storedColumn.column)
		whereExpression := fmt.Sprintf("%s IS NOT NULL", storedColumn.column)
		if err := db.Table(storedColumn.table).
			Select(selectExpression).
			Where(whereExpression).
			Scan(&storedValues).Error; err != nil {
			return 0, fmt.Errorf("preflight %s.%s timestamps: %w", storedColumn.table, storedColumn.column, err)
		}
		columnValues[storedColumn] = storedValues
	}

	var updated int64
	for _, storedColumn := range storedTimeColumns {
		for _, storedValue := range columnValues[storedColumn] {
			if storedValue.Value == nil || storedValue.Value.IsZero() {
				continue
			}
			instant := storedValue.Value.UTC()
			result := db.Table(storedColumn.table).
				Where("rowid = ?", storedValue.RowID).
				Update(storedColumn.column, instant)
			if result.Error != nil {
				return updated, fmt.Errorf("normalize %s.%s timestamp: %w", storedColumn.table, storedColumn.column, result.Error)
			}
			updated += result.RowsAffected
		}
	}

	for _, value := range values {
		if value.ValueTime == nil {
			continue
		}

		if value.ValueTime.IsZero() {
			result := db.Table("node_kvs").
				Where("node_id = ? AND key = ?", value.NodeID, value.Key).
				Delete(nil)
			if result.Error != nil {
				return updated, fmt.Errorf("remove zero %s timestamp from interval %s: %w", value.Key, value.NodeID, result.Error)
			}
			updated += result.RowsAffected
			continue
		}

		instant := value.ValueTime.UTC()
		result := db.Table("node_kvs").
			Where("node_id = ? AND key = ?", value.NodeID, value.Key).
			Update("value_time", instant)
		if result.Error != nil {
			return updated, fmt.Errorf("normalize %s timestamp for interval %s: %w", value.Key, value.NodeID, result.Error)
		}
		updated += result.RowsAffected
	}

	result := db.Exec(`
		DELETE FROM node_kvs
		WHERE key IN ('start_timezone', 'end_timezone')
		  AND node_id IN (SELECT id FROM node_cores WHERE kind = ?)
	`, core.IntervalType)
	if result.Error != nil {
		return updated, fmt.Errorf("remove legacy interval timezones: %w", result.Error)
	}
	updated += result.RowsAffected

	result = db.Exec(`
		INSERT OR IGNORE INTO node_kvs (node_id, key, value_text)
		SELECT id, 'client.general.timeZone', 'browser'
		FROM node_cores
		WHERE kind = ?
		ORDER BY id
		LIMIT 1
	`, core.SettingsType)
	if result.Error != nil {
		return updated, fmt.Errorf("add default client timezone: %w", result.Error)
	}
	updated += result.RowsAffected

	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_node_kvs_key_value_time
		ON node_kvs (key, value_time, node_id)
	`).Error; err != nil {
		return updated, fmt.Errorf("create interval time index: %w", err)
	}

	var legacyTimezoneCount int64
	if err := db.Table("node_kvs AS kv").
		Joins("JOIN node_cores AS core ON core.id = kv.node_id").
		Where("core.kind = ?", core.IntervalType).
		Where("kv.key IN ?", []string{"start_timezone", "end_timezone"}).
		Count(&legacyTimezoneCount).Error; err != nil {
		return updated, fmt.Errorf("verify legacy interval timezones: %w", err)
	}
	if legacyTimezoneCount != 0 {
		return updated, fmt.Errorf("verify legacy interval timezones: %d values remain", legacyTimezoneCount)
	}

	return updated, nil
}

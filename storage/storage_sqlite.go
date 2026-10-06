package storage

import (
	"strings"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

func NewSqliteStorage(path string) (*Storage, error) {
	db, err := initSqliteStorage(path)
	if err != nil {
		return nil, err
	}

	return &Storage{
		DB: db,
	}, nil
}

func initSqliteStorage(path string) (*gorm.DB, error) {
	if strings.TrimSpace(path) == "" {
		return nil, NewErrInvalidStoragePath()
	}

	db, err := gorm.Open(sqlite.New(sqlite.Config{
		DSN:        sqliteDSN(path),
		DriverName: "sqlite",
	}), &gorm.Config{NowFunc: func() time.Time { return time.Now().UTC() }})

	if err != nil {
		return nil, err
	}

	if err := configureSqlite(db, path); err != nil {
		return nil, err
	}

	if err := Migrate(db); err != nil {
		return nil, err
	}

	return db, nil
}

func sqliteDSN(path string) string {
	separator := "?"
	if strings.Contains(path, "?") {
		separator = "&"
	}
	return path + separator + "_pragma=busy_timeout%285000%29&_pragma=foreign_keys%281%29"
}

func Migrate(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		hadDashboardActive := tx.Migrator().HasColumn(&DashboardEntity{}, "Active")
		if err := tx.SetupJoinTable(&ContextEntity{}, "Tags", &ContextTagEntity{}); err != nil {
			return err
		}

		if err := tx.AutoMigrate(
			&TagEntity{},
			&ContextTagEntity{},
			&WorkspaceEntity{},
			&WorkspaceLinkRuleEntity{},
			&Properties{},
			&ClientProperties{},
			&ProjectEntity{},
			&ContextEntity{},
			&IntervalEntity{},
			&SavedQueryEntity{},
			&DashboardEntity{},
		); err != nil {
			return err
		}

		return migrateDashboardActivation(tx, hadDashboardActive)
	})
}

func migrateDashboardActivation(db *gorm.DB, hadActive bool) error {
	if db.Migrator().HasIndex(&DashboardEntity{}, "idx_dashboard_target") {
		if err := db.Migrator().DropIndex(&DashboardEntity{}, "idx_dashboard_target"); err != nil {
			return err
		}
	}
	if !hadActive {
		if err := db.Exec(`UPDATE dashboards SET active = 1
			WHERE type = 'custom' OR id IN (
				SELECT MIN(id) FROM dashboards WHERE type <> 'custom' GROUP BY workspace_id, type
			)`).Error; err != nil {
			return err
		}
	}
	return db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_dashboard_active_type
		ON dashboards (workspace_id, type) WHERE active = 1 AND type <> 'custom'`).Error
}

func configureSqlite(db *gorm.DB, path string) error {
	sqliteDB, err := db.DB()
	if err != nil {
		return err
	}

	if path == ":memory:" || strings.Contains(path, "mode=memory") {
		sqliteDB.SetMaxOpenConns(1)
		sqliteDB.SetMaxIdleConns(1)
	}

	var fkEnabled int
	if err := db.Raw("PRAGMA foreign_keys").Scan(&fkEnabled).Error; err != nil {
		return err
	}

	if fkEnabled != 1 {
		return NewErrForeignKeyDisabled()
	}

	return nil
}

package bootstrap

import (
	"database/sql"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/m87/ctx/core"
	"github.com/stretchr/testify/require"
)

func TestLegacyDatabaseGeneratorMigratesToUTCFormat(t *testing.T) {
	databasePath := generateTestDatabase(t, "generate-legacy-test-db.py")

	legacyDB := openTestDatabase(t, databasePath)
	require.Equal(t, "0.5.0", querySetting(t, legacyDB, "kvs", "database_version"))
	require.Greater(t, queryKVCount(t, legacyDB, "kvs", "start_timezone", "end_timezone"), int64(0))
	forceLegacyOffsetStart(t, legacyDB)
	offsetCoreID := forceLegacyCoreCreatedAtOffset(t, legacyDB)
	zeroEndKey := forceLegacyZeroEnd(t, legacyDB)
	before := queryIntervalInstants(t, legacyDB, "kvs")
	coreCreatedBefore := queryCoreCreatedInstant(t, legacyDB, offsetCoreID)
	require.NotEmpty(t, before)
	delete(before, zeroEndKey)
	require.NoError(t, legacyDB.Close())

	t.Setenv("DATABASE_PATH", databasePath)
	_, err := CreateManager()
	require.NoError(t, err)

	migratedDB := openTestDatabase(t, databasePath)
	require.Equal(t, core.CurrentDatabaseVersion, querySetting(t, migratedDB, "node_kvs", "database_version"))
	require.Equal(t, "browser", querySetting(t, migratedDB, "node_kvs", "client.general.timeZone"))
	require.Zero(t, queryKVCount(t, migratedDB, "node_kvs", "start_timezone", "end_timezone"))
	require.Equal(t, before, queryIntervalInstants(t, migratedDB, "node_kvs"))
	require.Equal(t, coreCreatedBefore, queryCoreCreatedInstant(t, migratedDB, offsetCoreID))
	requireStoredCoreCreatedAtUTC(t, migratedDB, offsetCoreID)
	require.True(t, indexExists(t, migratedDB, "idx_node_kvs_key_value_time"))
	require.NoError(t, migratedDB.Close())

	_, err = CreateManager()
	require.NoError(t, err)

	idempotentDB := openTestDatabase(t, databasePath)
	require.Equal(t, before, queryIntervalInstants(t, idempotentDB, "node_kvs"))
	require.Zero(t, queryKVCount(t, idempotentDB, "node_kvs", "start_timezone", "end_timezone"))
	require.NoError(t, idempotentDB.Close())
}

func TestCurrentDatabaseGeneratorUsesUTCFormat(t *testing.T) {
	databasePath := generateTestDatabase(t, "generate-test-db.py")

	db := openTestDatabase(t, databasePath)
	require.Equal(t, core.CurrentDatabaseVersion, querySetting(t, db, "kvs", "database_version"))
	require.Equal(t, "browser", querySetting(t, db, "kvs", "client.general.timeZone"))
	require.Zero(t, queryKVCount(t, db, "kvs", "start_timezone", "end_timezone"))
	require.True(t, indexExists(t, db, "idx_node_kvs_key_value_time"))

	instants := queryIntervalInstants(t, db, "kvs")
	require.NotEmpty(t, instants)
	for _, instant := range instants {
		require.Equal(t, int64(0), instant%int64(time.Second))
	}
	require.Zero(t, queryNonCanonicalUTCTimeCount(t, db, "kvs"))
	require.NoError(t, db.Close())

	t.Setenv("DATABASE_PATH", databasePath)
	_, err := CreateManager()
	require.NoError(t, err)

	openedDB := openTestDatabase(t, databasePath)
	require.Equal(t, core.CurrentDatabaseVersion, querySetting(t, openedDB, "node_kvs", "database_version"))
	require.Zero(t, queryKVCount(t, openedDB, "node_kvs", "start_timezone", "end_timezone"))
	require.NoError(t, openedDB.Close())
}

func TestMigrationRollsBackApplicationChangesWhenTimestampIsInvalid(t *testing.T) {
	databasePath := generateTestDatabase(t, "generate-legacy-test-db.py")

	db := openTestDatabase(t, databasePath)
	_, err := db.Exec(`
		UPDATE kvs
		SET value_time = 'not-a-time'
		WHERE node_id = (
			SELECT core.id
			FROM node_cores AS core
			JOIN kvs AS kv ON kv.node_id = core.id AND kv.key = 'start'
			WHERE core.kind = 'interval'
			LIMIT 1
		) AND key = 'start'
	`)
	require.NoError(t, err)
	require.NoError(t, db.Close())

	t.Setenv("DATABASE_PATH", databasePath)
	_, err = CreateManager()
	require.ErrorContains(t, err, "preflight interval timestamps")

	rolledBackDB := openTestDatabase(t, databasePath)
	require.Equal(t, "0.5.0", querySetting(t, rolledBackDB, "node_kvs", "database_version"))
	require.Greater(t, queryKVCount(t, rolledBackDB, "node_kvs", "start_timezone", "end_timezone"), int64(0))
	require.Equal(t, "", queryOptionalSetting(t, rolledBackDB, "node_kvs", "client.general.timeZone"))
	require.NoError(t, rolledBackDB.Close())
}

func generateTestDatabase(t *testing.T, scriptName string) string {
	t.Helper()
	python, err := exec.LookPath("python3")
	if err != nil {
		t.Skip("python3 is required to test database generators")
	}

	databasePath := filepath.Join(t.TempDir(), strings.TrimSuffix(scriptName, ".py")+".db")
	scriptPath := filepath.Join("..", "scripts", scriptName)
	command := exec.Command(python, scriptPath, "--micro-contexts", "2", "--force", "--output", databasePath)
	output, err := command.CombinedOutput()
	require.NoError(t, err, "generator output: %s", output)
	return databasePath
}

func openTestDatabase(t *testing.T, databasePath string) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", databasePath)
	require.NoError(t, err)
	require.NoError(t, db.Ping())
	return db
}

func querySetting(t *testing.T, db *sql.DB, table string, key string) string {
	t.Helper()
	var value string
	query := fmt.Sprintf("SELECT value_text FROM %s WHERE key = ? LIMIT 1", table)
	require.NoError(t, db.QueryRow(query, key).Scan(&value))
	return value
}

func queryOptionalSetting(t *testing.T, db *sql.DB, table string, key string) string {
	t.Helper()
	var value string
	query := fmt.Sprintf("SELECT value_text FROM %s WHERE key = ? LIMIT 1", table)
	err := db.QueryRow(query, key).Scan(&value)
	if err == sql.ErrNoRows {
		return ""
	}
	require.NoError(t, err)
	return value
}

func queryKVCount(t *testing.T, db *sql.DB, table string, keys ...string) int64 {
	t.Helper()
	placeholders := strings.TrimRight(strings.Repeat("?,", len(keys)), ",")
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE key IN (%s)", table, placeholders)
	args := make([]any, len(keys))
	for index, key := range keys {
		args[index] = key
	}
	var count int64
	require.NoError(t, db.QueryRow(query, args...).Scan(&count))
	return count
}

func queryIntervalInstants(t *testing.T, db *sql.DB, table string) map[string]int64 {
	t.Helper()
	query := fmt.Sprintf(`
		SELECT kv.node_id, kv.key, kv.value_time
		FROM %s AS kv
		JOIN node_cores AS core ON core.id = kv.node_id
		WHERE core.kind = 'interval' AND kv.key IN ('start', 'end')
		ORDER BY kv.node_id, kv.key
	`, table)
	rows, err := db.Query(query)
	require.NoError(t, err)
	defer rows.Close()

	instants := make(map[string]int64)
	for rows.Next() {
		var nodeID string
		var key string
		var raw string
		require.NoError(t, rows.Scan(&nodeID, &key, &raw))
		parsed, err := parsePersistedTime(raw)
		require.NoError(t, err, "parse %s/%s value %q", nodeID, key, raw)
		require.Equal(t, time.UTC, parsed.Location())
		instants[nodeID+":"+key] = parsed.UnixNano()
	}
	require.NoError(t, rows.Err())
	return instants
}

func forceLegacyZeroEnd(t *testing.T, db *sql.DB) string {
	t.Helper()
	var nodeID string
	require.NoError(t, db.QueryRow(`
		SELECT core.id
		FROM node_cores AS core
		JOIN kvs AS kv ON kv.node_id = core.id AND kv.key = 'end'
		WHERE core.kind = 'interval'
		ORDER BY core.id
		LIMIT 1
	`).Scan(&nodeID))
	_, err := db.Exec(
		"UPDATE kvs SET value_time = ? WHERE node_id = ? AND key = 'end'",
		"0001-01-01 00:00:00.000000000+00:00",
		nodeID,
	)
	require.NoError(t, err)
	return nodeID + ":end"
}

func forceLegacyOffsetStart(t *testing.T, db *sql.DB) {
	t.Helper()
	_, err := db.Exec(`
		UPDATE kvs
		SET value_time = '2026-08-02 03:00:00.000000000+09:00'
		WHERE node_id = (
			SELECT core.id
			FROM node_cores AS core
			JOIN kvs AS kv ON kv.node_id = core.id AND kv.key = 'start'
			WHERE core.kind = 'interval'
			ORDER BY core.id
			LIMIT 1
		) AND key = 'start'
	`)
	require.NoError(t, err)
}

func forceLegacyCoreCreatedAtOffset(t *testing.T, db *sql.DB) string {
	t.Helper()
	var id string
	require.NoError(t, db.QueryRow("SELECT id FROM node_cores ORDER BY id LIMIT 1").Scan(&id))
	_, err := db.Exec(`
		UPDATE node_cores
		SET created_at = '2026-08-02 03:00:00.000000000+09:00'
		WHERE id = ?
	`, id)
	require.NoError(t, err)
	return id
}

func queryCoreCreatedInstant(t *testing.T, db *sql.DB, id string) int64 {
	t.Helper()
	var raw string
	require.NoError(t, db.QueryRow("SELECT created_at FROM node_cores WHERE id = ?", id).Scan(&raw))
	parsed, err := parsePersistedTime(raw)
	require.NoError(t, err, "parse node core %s created_at value %q", id, raw)
	return parsed.UnixNano()
}

func requireStoredCoreCreatedAtUTC(t *testing.T, db *sql.DB, id string) {
	t.Helper()
	var raw string
	require.NoError(t, db.QueryRow("SELECT created_at FROM node_cores WHERE id = ?", id).Scan(&raw))
	parsed, err := parsePersistedTimeRaw(raw)
	require.NoError(t, err)
	_, offset := parsed.Zone()
	require.Zero(t, offset)
}

func queryNonCanonicalUTCTimeCount(t *testing.T, db *sql.DB, table string) int64 {
	t.Helper()
	query := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM %s AS kv
		JOIN node_cores AS core ON core.id = kv.node_id
		WHERE core.kind = 'interval'
		  AND kv.key IN ('start', 'end')
		  AND kv.value_time NOT LIKE '%%+00:00'
	`, table)
	var count int64
	require.NoError(t, db.QueryRow(query).Scan(&count))
	return count
}

func parsePersistedTime(raw string) (time.Time, error) {
	parsed, err := parsePersistedTimeRaw(raw)
	if err != nil {
		return time.Time{}, err
	}
	return parsed.UTC(), nil
}

func parsePersistedTimeRaw(raw string) (time.Time, error) {
	value := strings.TrimSpace(raw)
	if monotonicIndex := strings.Index(value, " m=+"); monotonicIndex >= 0 {
		value = value[:monotonicIndex]
	}
	for _, layout := range []string{
		"2006-01-02 15:04:05.999999999-07:00",
		"2006-01-02 15:04:05.999999999 -0700 MST",
		time.RFC3339Nano,
	} {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, fmt.Errorf("unsupported persisted time %q", raw)
}

func indexExists(t *testing.T, db *sql.DB, indexName string) bool {
	t.Helper()
	var count int
	require.NoError(t, db.QueryRow(
		"SELECT COUNT(*) FROM sqlite_master WHERE type = 'index' AND name = ?",
		indexName,
	).Scan(&count))
	return count == 1
}

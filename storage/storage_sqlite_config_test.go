package storage

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSQLitePragmasApplyToEveryConnection(t *testing.T) {
	for _, dsn := range []string{
		filepath.Join(t.TempDir(), "ctx.db"),
		"file:" + filepath.Join(t.TempDir(), "ctx.db") + "?cache=shared",
	} {
		t.Run(dsn, func(t *testing.T) {
			database, err := NewSqliteStorage(dsn)
			require.NoError(t, err)
			sqlDB, err := database.DB.DB()
			require.NoError(t, err)
			defer sqlDB.Close()
			first, err := sqlDB.Conn(context.Background())
			require.NoError(t, err)
			defer first.Close()
			second, err := sqlDB.Conn(context.Background())
			require.NoError(t, err)
			defer second.Close()
			for _, connection := range []*sql.Conn{first, second} {
				var foreignKeys, busyTimeout int
				require.NoError(t, connection.QueryRowContext(context.Background(), "PRAGMA foreign_keys").Scan(&foreignKeys))
				require.NoError(t, connection.QueryRowContext(context.Background(), "PRAGMA busy_timeout").Scan(&busyTimeout))
				require.Equal(t, 1, foreignKeys)
				require.Equal(t, 5000, busyTimeout)
			}
		})
	}
}

func TestSQLiteBusyTimeoutAllowsSecondWriterToWait(t *testing.T) {
	database, err := NewSqliteStorage(filepath.Join(t.TempDir(), "ctx.db"))
	require.NoError(t, err)
	sqlDB, err := database.DB.DB()
	require.NoError(t, err)
	defer sqlDB.Close()
	_, err = sqlDB.Exec("CREATE TABLE busy_test (id INTEGER PRIMARY KEY)")
	require.NoError(t, err)
	first, err := sqlDB.Conn(context.Background())
	require.NoError(t, err)
	defer first.Close()
	second, err := sqlDB.Conn(context.Background())
	require.NoError(t, err)
	defer second.Close()
	_, err = first.ExecContext(context.Background(), "BEGIN IMMEDIATE")
	require.NoError(t, err)
	_, err = first.ExecContext(context.Background(), "INSERT INTO busy_test (id) VALUES (1)")
	require.NoError(t, err)
	result := make(chan error, 1)
	go func() {
		_, writeErr := second.ExecContext(context.Background(), "INSERT INTO busy_test (id) VALUES (2)")
		result <- writeErr
	}()
	time.Sleep(50 * time.Millisecond)
	_, err = first.ExecContext(context.Background(), "COMMIT")
	require.NoError(t, err)
	select {
	case err := <-result:
		require.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("second writer did not complete after releasing lock")
	}
}

func TestSQLiteReopensExistingDatabaseWithOrphanedProjectReference(t *testing.T) {
	path := filepath.Join(t.TempDir(), "existing.db")
	database, err := NewSqliteStorage(path)
	require.NoError(t, err)
	sqlDB, err := database.DB.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())

	legacy, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	legacy.SetMaxOpenConns(1)
	_, err = legacy.Exec("PRAGMA foreign_keys = OFF")
	require.NoError(t, err)
	_, err = legacy.Exec("INSERT INTO contexts (id, name, workspace_id, status, archived, project_id) VALUES (?, ?, ?, ?, ?, ?)", "context-1", "Existing", "workspace-1", "inactive", false, "missing-project")
	require.NoError(t, err)
	require.NoError(t, legacy.Close())

	reopened, err := NewSqliteStorage(path)
	require.NoError(t, err)
	defer func() {
		openedDB, dbErr := reopened.DB.DB()
		if dbErr == nil {
			openedDB.Close()
		}
	}()
	stored, err := NewContextRepository(reopened.DB).GetById("context-1")
	require.NoError(t, err)
	require.NotNil(t, stored)
	require.NotNil(t, stored.ProjectId)
	require.Equal(t, "missing-project", *stored.ProjectId)
}

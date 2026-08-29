package bootstrap

import (
	"path/filepath"
	"testing"

	"github.com/m87/ctx/core"
	"github.com/m87/ctx/storage"
	"github.com/stretchr/testify/require"
)

func TestCreateManagerRunsNodDataMigration(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "legacy.db")
	storageDB, err := storage.NewSqliteStorage(databasePath)
	require.NoError(t, err)
	require.NoError(t, storageDB.DB.Exec(`
		CREATE TABLE node_cores (
			id TEXT PRIMARY KEY,
			namespace_id TEXT,
			parent_id TEXT,
			kind TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT '',
			name TEXT NOT NULL
		)
	`).Error)
	require.NoError(t, storageDB.DB.Exec(`
		INSERT INTO node_cores (id, kind, name)
		VALUES ('legacy-workspace', 'workspace', 'Legacy workspace')
	`).Error)
	sqlDB, err := storageDB.DB.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())

	t.Setenv("DATABASE_PATH", databasePath)
	manager, err := CreateManager()
	require.NoError(t, err)

	workspace, err := manager.WorkspaceRepository.GetById("legacy-workspace")
	require.NoError(t, err)
	require.Equal(t, "Legacy workspace", workspace.Name)

	reopened, err := storage.NewSqliteStorage(databasePath)
	require.NoError(t, err)
	properties, err := storage.NewPropertiesRepository(reopened.DB).Load()
	require.NoError(t, err)
	require.Equal(t, core.CurrentDatabaseVersion, properties.DatabaseVersion)
	require.False(t, reopened.DB.Migrator().HasTable("node_cores"))
}

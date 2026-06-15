package bootstrap

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/m87/ctx/core"
	ctxlog "github.com/m87/ctx/log"
	"github.com/m87/nod/sqlite"
	"github.com/stretchr/testify/require"
)

func TestCreateManagerRunsWorkspaceMigrationOnlyOnce(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "ctx.db")
	t.Setenv("DATABASE_PATH", databasePath)

	manager, err := CreateManager()
	require.NoError(t, err)

	contextId, err := manager.ContextRepository.Save(&core.Context{Name: "orphan"})
	require.NoError(t, err)

	_, err = CreateManager()
	require.NoError(t, err)

	repository, err := sqlite.NewRepository(databasePath, ctxlog.Logger, NewMapperRegistry())
	require.NoError(t, err)
	context, err := NewContextRepository(repository).GetById(contextId)
	require.NoError(t, err)
	require.Empty(t, context.WorkspaceId)

	systemInfo, err := NewSystemInfoRepository(repository).Load()
	require.NoError(t, err)
	require.Equal(t, core.CurrentDatabaseVersion, systemInfo.DatabaseVersion)
}

func TestCreateManagerRunInTransactionRollsBackOnError(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "ctx.db")
	t.Setenv("DATABASE_PATH", databasePath)
	manager, err := CreateManager()
	require.NoError(t, err)
	wantErr := errors.New("rollback")

	err = manager.RunInTransaction(func(txManager *core.ContextManager) error {
		_, saveErr := txManager.WorkspaceRepository.Save(&core.Workspace{Name: "rollback-me"})
		require.NoError(t, saveErr)
		return wantErr
	})

	require.ErrorIs(t, err, wantErr)
	repository, err := sqlite.NewRepository(databasePath, ctxlog.Logger, NewMapperRegistry())
	require.NoError(t, err)
	workspaces, err := NewWorkspaceRepository(repository).List()
	require.NoError(t, err)
	for _, workspace := range workspaces {
		require.NotEqual(t, "rollback-me", workspace.Name)
	}
}

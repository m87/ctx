package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestArchiveContext(t *testing.T) {
	tm := newTestManager()
	workspaceID, err := tm.Workspaces.Save(&Workspace{Name: "archive-test"})
	require.NoError(t, err)
	contextID, err := tm.Manager.CreateContext(&Context{
		Name:        "context-to-archive",
		Description: "keep description",
		WorkspaceId: workspaceID,
		Tags:        []string{"keep-tag"},
	})
	require.NoError(t, err)

	err = tm.Manager.ArchiveContext(contextID)
	require.NoError(t, err)

	context, err := tm.Contexts.GetById(contextID)
	require.NoError(t, err)
	require.True(t, context.Archived)
	require.Equal(t, "archived", context.Status)
	require.Equal(t, "keep description", context.Description)
	require.Equal(t, []string{"keep-tag"}, context.Tags)
}

func TestRestoreContext(t *testing.T) {
	tm := newTestManager()
	workspaceID, err := tm.Workspaces.Save(&Workspace{Name: "restore-test"})
	require.NoError(t, err)
	contextID, err := tm.Manager.CreateContext(&Context{
		Name:        "context-to-restore",
		Description: "keep description",
		WorkspaceId: workspaceID,
		Tags:        []string{"keep-tag"},
	})
	require.NoError(t, err)

	err = tm.Manager.ArchiveContext(contextID)
	require.NoError(t, err)

	err = tm.Manager.RestoreContext(contextID)
	require.NoError(t, err)

	context, err := tm.Contexts.GetById(contextID)
	require.NoError(t, err)
	require.False(t, context.Archived)
	require.Equal(t, "inactive", context.Status)
	require.Equal(t, "keep description", context.Description)
	require.Equal(t, []string{"keep-tag"}, context.Tags)
}

func TestArchiveActiveContext(t *testing.T) {
	tm := newTestManager()
	workspaceID, err := tm.Workspaces.Save(&Workspace{Name: "archive-active-test"})
	require.NoError(t, err)
	contextID, err := tm.Manager.CreateContext(&Context{
		Name:        "active-context",
		Description: "keep description",
		WorkspaceId: workspaceID,
		Tags:        []string{"keep-tag"},
	})
	require.NoError(t, err)

	err = tm.Manager.SwitchContext(tm.Contexts.items[contextID])
	require.NoError(t, err)

	err = tm.Manager.ArchiveContext(contextID)
	require.Error(t, err)
	require.IsType(t, &ArchiveContextActiveError{}, err)
}

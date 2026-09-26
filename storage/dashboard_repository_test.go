package storage

import (
	"encoding/json"
	"testing"

	"github.com/m87/ctx/core"
	"github.com/stretchr/testify/require"
)

func TestDashboardRepositorySavesGetsAndListsByWorkspace(t *testing.T) {
	storage, err := CreateTestInMemoryStorage()
	require.NoError(t, err)
	workspaceRepository := NewWorkspaceRepository(storage.DB)
	repository := NewDashboardRepository(storage.DB)
	_, err = workspaceRepository.Save(&core.Workspace{Id: "workspace-1", Name: "Workspace"})
	require.NoError(t, err)
	_, err = workspaceRepository.Save(&core.Workspace{Id: "workspace-2", Name: "Other"})
	require.NoError(t, err)

	second := &core.Dashboard{
		WorkspaceId: "workspace-1",
		Type:        core.DashboardTypeCustom,
		Name:        "Second",
		Definition:  json.RawMessage(`{"widgets":[]}`),
	}
	secondID, err := repository.Save(second)
	require.NoError(t, err)
	require.NotEmpty(t, secondID)
	first := &core.Dashboard{
		WorkspaceId: "workspace-1",
		Type:        core.DashboardTypeWorkspace,
		TargetId:    "workspace-1",
		Name:        "First",
		Definition:  json.RawMessage(`{"columns":16}`),
	}
	firstID, err := repository.Save(first)
	require.NoError(t, err)
	_, err = repository.Save(&core.Dashboard{
		WorkspaceId: "workspace-2",
		Type:        core.DashboardTypeCustom,
		Name:        "Other",
		Definition:  json.RawMessage(`{}`),
	})
	require.NoError(t, err)

	retrieved, err := repository.GetById(firstID)
	require.NoError(t, err)
	require.Equal(t, first, retrieved)

	byTarget, err := repository.GetByTarget(
		"workspace-1",
		core.DashboardTypeWorkspace,
		"workspace-1",
	)
	require.NoError(t, err)
	require.Equal(t, first, byTarget)

	first.Type = core.DashboardTypeCustom
	first.TargetId = ""
	first.Name = "Edited"
	_, err = repository.Save(first)
	require.NoError(t, err)
	retrieved, err = repository.GetById(firstID)
	require.NoError(t, err)
	require.Equal(t, first, retrieved)
	byTarget, err = repository.GetByTarget(
		"workspace-1",
		core.DashboardTypeWorkspace,
		"workspace-1",
	)
	require.NoError(t, err)
	require.Nil(t, byTarget)

	dashboards, err := repository.ListByWorkspace("workspace-1")
	require.NoError(t, err)
	require.Equal(t, []*core.Dashboard{first, second}, dashboards)

	require.NoError(t, repository.Delete(firstID))
	deleted, err := repository.GetById(firstID)
	require.NoError(t, err)
	require.Nil(t, deleted)
}

func TestDashboardRepositoryCascadesWorkspaceDeletion(t *testing.T) {
	storage, err := CreateTestInMemoryStorage()
	require.NoError(t, err)
	workspaceRepository := NewWorkspaceRepository(storage.DB)
	repository := NewDashboardRepository(storage.DB)
	_, err = workspaceRepository.Save(&core.Workspace{Id: "workspace-1", Name: "Workspace"})
	require.NoError(t, err)
	dashboard := &core.Dashboard{
		WorkspaceId: "workspace-1",
		Type:        core.DashboardTypeCustom,
		Name:        "Dashboard",
		Definition:  json.RawMessage(`{}`),
	}
	id, err := repository.Save(dashboard)
	require.NoError(t, err)

	require.NoError(t, workspaceRepository.Delete("workspace-1"))
	retrieved, err := repository.GetById(id)
	require.NoError(t, err)
	require.Nil(t, retrieved)
}

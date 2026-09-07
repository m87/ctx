package storage

import (
	"testing"

	"github.com/m87/ctx/core"
	"github.com/stretchr/testify/require"
)

func TestSavedQueryRepositorySavesGetsAndListsByWorkspace(t *testing.T) {
	storage, err := CreateTestInMemoryStorage()
	require.NoError(t, err)
	workspaceRepository := NewWorkspaceRepository(storage.DB)
	repository := NewSavedQueryRepository(storage.DB)
	_, err = workspaceRepository.Save(&core.Workspace{Id: "workspace-1", Name: "Workspace"})
	require.NoError(t, err)
	_, err = workspaceRepository.Save(&core.Workspace{Id: "workspace-2", Name: "Other"})
	require.NoError(t, err)

	second := &core.SavedQuery{
		WorkspaceId: "workspace-1",
		Name:        "Second",
		Query:       "tag = second",
	}
	secondID, err := repository.Save(second)
	require.NoError(t, err)
	require.NotEmpty(t, secondID)
	first := &core.SavedQuery{
		WorkspaceId: "workspace-1",
		Name:        "First",
		Query:       "tag = first",
	}
	firstID, err := repository.Save(first)
	require.NoError(t, err)
	_, err = repository.Save(&core.SavedQuery{
		WorkspaceId: "workspace-2",
		Name:        "Other",
		Query:       "tag = other",
	})
	require.NoError(t, err)

	retrieved, err := repository.GetById(firstID)
	require.NoError(t, err)
	require.Equal(t, first, retrieved)

	queries, err := repository.ListByWorkspace("workspace-1")
	require.NoError(t, err)
	require.Equal(t, []*core.SavedQuery{first, second}, queries)

	missing, err := repository.GetById("missing")
	require.NoError(t, err)
	require.Nil(t, missing)

	require.NoError(t, repository.Delete(firstID))
	deleted, err := repository.GetById(firstID)
	require.NoError(t, err)
	require.Nil(t, deleted)
}

func TestSavedQueryRepositoryCascadesWorkspaceDeletion(t *testing.T) {
	storage, err := CreateTestInMemoryStorage()
	require.NoError(t, err)
	workspaceRepository := NewWorkspaceRepository(storage.DB)
	repository := NewSavedQueryRepository(storage.DB)
	_, err = workspaceRepository.Save(&core.Workspace{Id: "workspace-1", Name: "Workspace"})
	require.NoError(t, err)
	query := &core.SavedQuery{
		WorkspaceId: "workspace-1",
		Name:        "Query",
		Query:       "tag = query",
	}
	id, err := repository.Save(query)
	require.NoError(t, err)

	require.NoError(t, workspaceRepository.Delete("workspace-1"))
	retrieved, err := repository.GetById(id)
	require.NoError(t, err)
	require.Nil(t, retrieved)
}

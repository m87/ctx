package storage

import (
	"testing"

	"github.com/m87/ctx/core"
	"github.com/stretchr/testify/require"
)

func TestWorkspaceRepository(t *testing.T) {
	t.Run("GetById", func(t *testing.T) {
		repository := newTestWorkspaceRepository(t)
		workspace := &core.Workspace{
			Id:          "workspace-1",
			Name:        "Workspace",
			Description: "Description",
		}
		_, err := repository.Save(workspace)
		require.NoError(t, err)

		stored, err := repository.GetById(workspace.Id)
		require.NoError(t, err)
		require.Equal(t, workspace, stored)

		missing, err := repository.GetById("missing")
		require.NoError(t, err)
		require.Nil(t, missing)
	})

	t.Run("Save", func(t *testing.T) {
		repository := newTestWorkspaceRepository(t)
		workspace := &core.Workspace{
			Name:        "Workspace",
			Description: "Description",
		}

		id, err := repository.Save(workspace)
		require.NoError(t, err)
		require.NotEmpty(t, id)
		require.Equal(t, id, workspace.Id)

		workspace.Name = "Updated workspace"
		workspace.Description = "Updated description"
		updatedId, err := repository.Save(workspace)
		require.NoError(t, err)
		require.Equal(t, id, updatedId)

		stored, err := repository.GetById(id)
		require.NoError(t, err)
		require.Equal(t, workspace, stored)

		workspaces, err := repository.List()
		require.NoError(t, err)
		require.Len(t, workspaces, 1)

		_, err = repository.Save(nil)
		require.EqualError(t, err, "workspace is required")
	})

	t.Run("Delete", func(t *testing.T) {
		repository := newTestWorkspaceRepository(t)
		first := &core.Workspace{Id: "workspace-1", Name: "First"}
		second := &core.Workspace{Id: "workspace-2", Name: "Second"}
		_, err := repository.Save(first)
		require.NoError(t, err)
		_, err = repository.Save(second)
		require.NoError(t, err)

		require.NoError(t, repository.Delete(first.Id))

		deleted, err := repository.GetById(first.Id)
		require.NoError(t, err)
		require.Nil(t, deleted)

		remaining, err := repository.GetById(second.Id)
		require.NoError(t, err)
		require.Equal(t, second, remaining)
		require.NoError(t, repository.Delete("missing"))
	})

	t.Run("List", func(t *testing.T) {
		repository := newTestWorkspaceRepository(t)

		empty, err := repository.List()
		require.NoError(t, err)
		require.Empty(t, empty)

		first := &core.Workspace{Id: "workspace-1", Name: "First", Description: "First description"}
		second := &core.Workspace{Id: "workspace-2", Name: "Second", Description: "Second description"}
		_, err = repository.Save(first)
		require.NoError(t, err)
		_, err = repository.Save(second)
		require.NoError(t, err)

		workspaces, err := repository.List()
		require.NoError(t, err)
		require.ElementsMatch(t, []*core.Workspace{first, second}, workspaces)
	})

	t.Run("Save and replace link rules", func(t *testing.T) {
		repository := newTestWorkspaceRepository(t)
		workspace := &core.Workspace{
			Id:   "workspace-1",
			Name: "Workspace",
			Properties: &core.WorkspaceSettings{LinkRules: []core.LinkRule{
				{Regexp: `CTX-(\d+)`, Link: "https://example.test/$1"},
				{Regexp: `BUG-(\d+)`, Link: "https://bugs.test/$1"},
			}},
		}

		_, err := repository.Save(workspace)
		require.NoError(t, err)
		stored, err := repository.GetById(workspace.Id)
		require.NoError(t, err)
		require.Equal(t, workspace.Properties, stored.Properties)

		workspace.Properties.LinkRules = workspace.Properties.LinkRules[:1]
		_, err = repository.Save(workspace)
		require.NoError(t, err)
		stored, err = repository.GetById(workspace.Id)
		require.NoError(t, err)
		require.Equal(t, workspace.Properties, stored.Properties)

		var count int64
		require.NoError(t, repository.db.Model(&WorkspaceLinkRuleEntity{}).Count(&count).Error)
		require.Equal(t, int64(1), count)
	})

	t.Run("SaveAll returns transaction errors and rolls back", func(t *testing.T) {
		repository := newTestWorkspaceRepository(t)
		_, err := repository.SaveAll([]*core.Workspace{{Name: "Valid"}, nil})
		require.EqualError(t, err, "workspace is required")

		workspaces, listErr := repository.List()
		require.NoError(t, listErr)
		require.Empty(t, workspaces)
	})
}

func newTestWorkspaceRepository(t *testing.T) *WorkspaceRepository {
	t.Helper()
	storage, _ := CreateTestInMemoryStorage()
	return NewWorkspaceRepository(storage.DB)
}

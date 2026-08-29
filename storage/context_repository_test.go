package storage

import (
	"testing"

	"github.com/m87/ctx/core"
	"github.com/stretchr/testify/require"
)

func TestContextRepository(t *testing.T) {
	t.Helper()

	t.Run("Save and retrieve context", func(t *testing.T) {
		storage, err := CreateTestInMemoryStorage()
		require.NoError(t, err)
		repo := NewContextRepository(storage.DB)

		context := &core.Context{
			Name:        "Test Context",
			Description: "Description",
			WorkspaceId: "workspace-1",
			Tags:        []*core.Tag{{Name: "tag1"}},
		}

		id, err := repo.Save(context)
		require.NoError(t, err)
		require.NotEmpty(t, id)

		retrievedContext, err := repo.GetById(id)
		require.NoError(t, err)
		require.NotNil(t, retrievedContext)
		require.Equal(t, context.Name, retrievedContext.Name)
		require.Equal(t, context.Description, retrievedContext.Description)
		require.Equal(t, context.WorkspaceId, retrievedContext.WorkspaceId)
	})

	t.Run("Delete context", func(t *testing.T) {
		storage, err := CreateTestInMemoryStorage()
		require.NoError(t, err)
		repo := NewContextRepository(storage.DB)

		context := &core.Context{
			Name:        "Test Context",
			Description: "Description",
			WorkspaceId: "workspace-1",
			Tags:        []*core.Tag{{Name: "tag1"}},
		}

		id, err := repo.Save(context)
		require.NoError(t, err)
		require.NotEmpty(t, id)

		err = repo.Delete(id)
		require.NoError(t, err)

		deletedContext, err := repo.GetById(id)
		require.Error(t, err)
		require.Nil(t, deletedContext)
	})

	t.Run("Get context with tags and projectMetadata", func(t *testing.T) {
		storage, err := CreateTestInMemoryStorage()
		require.NoError(t, err)
		repo := NewContextRepository(storage.DB)
		projectRepo := NewProjectRepository(storage.DB)

		project := &core.Project{Id: "project-1", Name: "Project 1", WorkspaceId: "workspace-1"}
		_, err = projectRepo.Save(project)
		require.NoError(t, err)

		context := &core.Context{
			Name:        "Test Context",
			Description: "Description",
			WorkspaceId: "workspace-1",
			Tags:        []*core.Tag{{Name: "tag1"}, {Name: "tag2"}},
			Project:     &core.ProjectMetadata{Id: project.Id, Name: "Stale project name"},
		}

		id, err := repo.Save(context)
		require.NoError(t, err)
		require.NotEmpty(t, id)

		retrievedContext, err := repo.GetById(id)
		require.NoError(t, err)
		require.NotNil(t, retrievedContext)
		require.Equal(t, context.Name, retrievedContext.Name)
		require.Equal(t, context.Description, retrievedContext.Description)
		require.Equal(t, context.WorkspaceId, retrievedContext.WorkspaceId)
		require.Len(t, retrievedContext.Tags, len(context.Tags))
		require.NotNil(t, retrievedContext.Project)
		require.Equal(t, project.Id, retrievedContext.Project.Id)
		require.Equal(t, project.Name, retrievedContext.Project.Name)
		require.ElementsMatch(t, context.Tags, retrievedContext.Tags)
	})

	t.Run("Update context and replace tags", func(t *testing.T) {
		storage, err := CreateTestInMemoryStorage()
		require.NoError(t, err)
		repo := NewContextRepository(storage.DB)
		projectRepo := NewProjectRepository(storage.DB)

		project := &core.Project{Id: "project-1", Name: "Project 1", WorkspaceId: "workspace-1"}
		_, err = projectRepo.Save(project)
		require.NoError(t, err)

		context := &core.Context{
			Name:        "Original Context",
			WorkspaceId: "workspace-1",
			Tags:        []*core.Tag{{Name: "old-tag"}},
			Project:     &core.ProjectMetadata{Id: project.Id, Name: project.Name},
		}

		id, err := repo.Save(context)
		require.NoError(t, err)

		context.Name = "Updated Context"
		context.Tags = []*core.Tag{{Name: "new-tag"}}
		context.Project = nil
		context.ProjectId = nil

		updatedId, err := repo.Save(context)
		require.NoError(t, err)
		require.Equal(t, id, updatedId)

		retrievedContext, err := repo.GetById(id)
		require.NoError(t, err)
		require.Equal(t, "Updated Context", retrievedContext.Name)
		require.Equal(t, context.Tags, retrievedContext.Tags)
		require.Nil(t, retrievedContext.Project)

	})
}

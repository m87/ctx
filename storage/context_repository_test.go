package storage

import (
	"testing"

	"github.com/m87/ctx/core"
	"github.com/stretchr/testify/assert"
)


func TestContextRepository(t *testing.T) {
	t.Helper()

	t.Run("Save and retrieve context", func(t *testing.T) {
		storage, _ := CreateTestInMemoryStorage()
		repo := NewContextRepository(storage.DB)

		context := &core.Context{
			Name:        "Test Context",
			Description: "Description",
			WorkspaceId: "workspace-1",
		}

		id, err := repo.Save(context)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)

		retrievedContext, err := repo.GetById(id)
		assert.NoError(t, err)
		assert.NotNil(t, retrievedContext)
		assert.Equal(t, context.Name, retrievedContext.Name)
		assert.Equal(t, context.Description, retrievedContext.Description)
		assert.Equal(t, context.WorkspaceId, retrievedContext.WorkspaceId)
	})

	t.Run("Delete context", func(t *testing.T) {
		storage, _ := CreateTestInMemoryStorage()
		repo := NewContextRepository(storage.DB)

		context := &core.Context{
			Name:        "Test Context",
			Description: "Description",
			WorkspaceId: "workspace-1",
		}

		id, err := repo.Save(context)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)

		err = repo.Delete(id)
		assert.NoError(t, err)

		deletedContext, err := repo.GetById(id)
		assert.Error(t, err)
		assert.Nil(t, deletedContext)
	})

	t.Run("Get context with tags and projectMetadata", func(t *testing.T) {
		storage, _ := CreateTestInMemoryStorage()
		repo := NewContextRepository(storage.DB)

		context := &core.Context{
			Name:        "Test Context",
			Description: "Description",
			WorkspaceId: "workspace-1",
			Tags:        []*core.Tag{&core.Tag{Name: "tag1"}, &core.Tag{Name: "tag2"}},
			Project: &core.ProjectMetadata{
				Id: "project-1",
				Name: 		"Project 1",
			},
		}

		id, err := repo.Save(context)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)

		retrievedContext, err := repo.GetById(id)
		assert.NoError(t, err)
		assert.NotNil(t, retrievedContext)
		assert.Equal(t, context.Name, retrievedContext.Name)
		assert.Equal(t, context.Description, retrievedContext.Description)
		assert.Equal(t, context.WorkspaceId, retrievedContext.WorkspaceId)
		assert.Equal(t, len(context.Tags), len(retrievedContext.Tags))
		assert.Equal(t, context.Project.Id, retrievedContext.Project.Id)
		assert.Equal(t, context.Project.Name, retrievedContext.Project.Name)

	})
}	

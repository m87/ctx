package storage

import (
	"testing"

	"github.com/m87/ctx/core"
	"github.com/stretchr/testify/assert"
)


func TestProjectRepository(t *testing.T) {
	t.Helper()

  t.Run("Save and Get Project", func(t *testing.T) {
		storage, _ := CreateTestInMemoryStorage()
		projectRepo := NewProjectRepository(storage.DB)

		project := &core.Project{
			Name:        "Test Project",
			WorkspaceId: "workspace-1",
		}

		id, err := projectRepo.Save(project)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)

		retrievedProject, err := projectRepo.GetById(id)
		assert.NoError(t, err)
		assert.Equal(t, project.Name, retrievedProject.Name)
		assert.Equal(t, project.WorkspaceId, retrievedProject.WorkspaceId)
	})

	t.Run("Delete Project", func(t *testing.T) {
		storage, _ := CreateTestInMemoryStorage()
		projectRepo := NewProjectRepository(storage.DB)

		project := &core.Project{
			Name:        "Test Project",
			WorkspaceId: "workspace-1",
		}

		id, err := projectRepo.Save(project)
		assert.NoError(t, err)

		err = projectRepo.Delete(id)
		assert.NoError(t, err)

		retrievedProject, err := projectRepo.GetById(id)
		assert.Error(t, err)
		assert.Nil(t, retrievedProject)
	})


}

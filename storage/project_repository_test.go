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

	t.Run("List returns only root projects", func(t *testing.T) {
		storage, err := CreateTestInMemoryStorage()
		assert.NoError(t, err)
		repo := NewProjectRepository(storage.DB)
		_, err = repo.SaveAll([]*core.Project{
			{Id: "root", Name: "Root", WorkspaceId: "workspace-1"},
			{Id: "child", Name: "Child", ParentId: "root", WorkspaceId: "workspace-1"},
		})
		assert.NoError(t, err)

		projects, err := repo.List("workspace-1")
		assert.NoError(t, err)
		assert.Equal(t, []*core.Project{{Id: "root", Name: "Root", WorkspaceId: "workspace-1"}}, projects)
	})

	t.Run("SaveAll returns transaction errors and rolls back", func(t *testing.T) {
		storage, err := CreateTestInMemoryStorage()
		assert.NoError(t, err)
		repo := NewProjectRepository(storage.DB)
		_, err = repo.SaveAll([]*core.Project{{Name: "Valid", WorkspaceId: "workspace-1"}, nil})
		assert.EqualError(t, err, "project is required")

		var count int64
		assert.NoError(t, storage.DB.Model(&ProjectEntity{}).Count(&count).Error)
		assert.Zero(t, count)
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

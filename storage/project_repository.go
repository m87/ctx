package storage

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/m87/ctx/core"
	"gorm.io/gorm"
)

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) GetById(id string) (*core.Project, error) {
	var entity ProjectEntity
	if err := r.db.First(&entity, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return entity.ToModel(), nil
}

func (r *ProjectRepository) Save(project *core.Project) (string, error) {
	return saveProject(r.db, project)
}

func saveProject(db *gorm.DB, project *core.Project) (string, error) {
	if project == nil {
		return "", fmt.Errorf("project is required")
	}

	if project.Id == "" {
		project.Id = uuid.NewString()
	}
	entity := NewProjectEntityFromModel(project)
	if err := db.Save(entity).Error; err != nil {
		return "", err
	}

	return entity.Id, nil
}

func (r *ProjectRepository) Delete(id string) error {
	if err := r.db.Delete(&ProjectEntity{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *ProjectRepository) ListChildren(parentId string) ([]*core.Project, error) {
	if parentId == "" {
		return nil, nil
	}

	var entities []ProjectEntity
	if err := r.db.Where("parent_id = ?", parentId).Find(&entities).Error; err != nil {
		return nil, err
	}

	projects := make([]*core.Project, len(entities))
	for i, entity := range entities {
		projects[i] = entity.ToModel()
	}

	return projects, nil
}

func (r *ProjectRepository) SaveAll(projects []*core.Project) ([]string, error) {
	ids := make([]string, len(projects))
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		for i, project := range projects {
			id, err := saveProject(tx, project)
			if err != nil {
				return err
			}
			ids[i] = id
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return ids, nil

}

func (r *ProjectRepository) List(workspaceId string) ([]*core.Project, error) {
	var entities []ProjectEntity
	if err := r.db.Where("workspace_id = ? AND (parent_id = '' OR parent_id IS NULL)", workspaceId).Find(&entities).Error; err != nil {
		return nil, err
	}

	projects := make([]*core.Project, len(entities))
	for i, entity := range entities {
		projects[i] = entity.ToModel()
	}

	return projects, nil
}

func (r *ProjectRepository) ListIncludingArchived(workspaceId string) ([]*core.Project, error) {
	var entities []ProjectEntity
	if err := r.db.Unscoped().Where("workspace_id = ?", workspaceId).Find(&entities).Error; err != nil {
		return nil, err
	}

	projects := make([]*core.Project, len(entities))
	for i, entity := range entities {
		projects[i] = entity.ToModel()
	}

	return projects, nil
}

func (r *ProjectRepository) ListToSync(limit int) ([]*core.Project, error) {
	var entities []ProjectEntity
	if err := r.db.Limit(limit).Find(&entities).Error; err != nil {
		return nil, err
	}

	projects := make([]*core.Project, len(entities))
	for i, entity := range entities {
		projects[i] = entity.ToModel()
	}

	return projects, nil
}

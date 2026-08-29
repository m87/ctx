package storage

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/m87/ctx/core"
	"gorm.io/gorm"
)

type WorkspaceRepository struct {
	db *gorm.DB
}

func NewWorkspaceRepository(db *gorm.DB) *WorkspaceRepository {
	return &WorkspaceRepository{db: db}
}

func (r *WorkspaceRepository) GetById(id string) (*core.Workspace, error) {
	var entity WorkspaceEntity
	if err := r.db.First(&entity, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return entity.ToModel(), nil
}

func (r *WorkspaceRepository) Save(workspace *core.Workspace) (string, error) {
	return saveWorkspace(r.db, workspace)
}

func (r *WorkspaceRepository) Delete(id string) error {
	return r.db.Delete(&WorkspaceEntity{}, "id = ?", id).Error
}

func (r *WorkspaceRepository) List() ([]*core.Workspace, error) {
	var entities []WorkspaceEntity
	if err := r.db.Find(&entities).Error; err != nil {
		return nil, err
	}

	workspaces := make([]*core.Workspace, len(entities))
	for i := range entities {
		workspaces[i] = entities[i].ToModel()
	}

	return workspaces, nil
}

func (r *WorkspaceRepository) ListToSync(limit int) ([]*core.Workspace, error) {
	return r.List()
}

func (r *WorkspaceRepository) SaveAll(workspaces []*core.Workspace) ([]string, error) {
	ids := make([]string, len(workspaces))

	r.db.Transaction(func(tx *gorm.DB) error {
		for i, workspace := range workspaces {
			id, err := saveWorkspace(tx, workspace)
			if err != nil {
				return err
			}

			ids[i] = id
		}
		return nil
	})

	return ids, nil
}

func saveWorkspace(db *gorm.DB, workspace *core.Workspace) (string, error) {
	if workspace == nil {
		return "", fmt.Errorf("workspace is required")
	}
	if workspace.Id == "" {
		workspace.Id = uuid.NewString()
	}

	entity := NewWorkspaceEntityFromModel(workspace)
	if err := db.Save(entity).Error; err != nil {
		return "", err
	}

	return entity.Id, nil
}

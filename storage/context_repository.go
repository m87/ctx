package storage

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/m87/ctx/core"
	"gorm.io/gorm"
)

type ContextRepository struct {
	db *gorm.DB
}

func NewContextRepository(db *gorm.DB) *ContextRepository {
	return &ContextRepository{db: db}
}

func (r *ContextRepository) GetById(id string) (*core.Context, error) {
	var context ContextEntity
	if err := r.db.Preload("Tags").Preload("ProjectMetadata").First(&context, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return context.toModel(), nil
}

func (r *ContextRepository) ListByProject(projectId string) ([]*core.Context, error) {
	if projectId == "" {
		return nil, nil
	}

	entities := []*ContextEntity{}
	if err := r.db.Preload("Tags").Preload("ProjectMetadata").Find(&entities, "project_id = ?", projectId).Error; err != nil {
		return nil, err
	}

	contexts := []*core.Context{}
	for _, e := range entities {
		contexts = append(contexts, e.toModel())
	}

	return contexts, nil
}

func (r *ContextRepository) Save(context *core.Context) (string, error) {
	return save(r.db, context)
}

func save(tx *gorm.DB, context *core.Context) (string, error) {
	if context == nil {
		return "", fmt.Errorf("context is required")
	}

	isNew := context.Id == ""
	if isNew {
		context.Id = uuid.NewString()
	}

	entity := NewContextEntity(context)
	if err := tx.Transaction(func(db *gorm.DB) error {
		query := db.Omit("Tags", "ProjectMetadata")
		if isNew {
			if err := query.Create(entity).Error; err != nil {
				return err
			}
		} else if err := query.Save(entity).Error; err != nil {
			return err
		}

		for _, tag := range entity.Tags {
			if err := db.Save(tag).Error; err != nil {
				return err
			}
		}

		return db.Model(entity).Association("Tags").Replace(entity.Tags)
	}); err != nil {
		return "", err
	}

	return context.Id, nil
}

func (r *ContextRepository) Delete(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("context_id = ?", id).Delete(&ContextTagEntity{}).Error; err != nil {
			return err
		}

		return tx.Delete(&ContextEntity{}, "id = ?", id).Error
	})
}

func (r *ContextRepository) List() ([]*core.Context, error) {
	entities := []*ContextEntity{}
	if err := r.db.Preload("Tags").Preload("ProjectMetadata").Find(&entities).Error; err != nil {
		return nil, err
	}

	contexts := []*core.Context{}
	for _, e := range entities {
		contexts = append(contexts, e.toModel())
	}

	return contexts, nil
}

func (r *ContextRepository) GetActive() (*core.Context, error) {
	entity := &ContextEntity{}
	result := r.db.Preload("Tags").Preload("ProjectMetadata").Where("status = ?", "active").Limit(1).Find(entity)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return entity.toModel(), nil
}

func (r *ContextRepository) ListByWorkspace(workspaceId string) ([]*core.Context, error) {
	entities := []*ContextEntity{}
	if err := r.db.Preload("Tags").Preload("ProjectMetadata").Find(&entities, "workspace_id = ? AND archived = ?", workspaceId, false).Error; err != nil {
		return nil, err
	}

	contexts := []*core.Context{}
	for _, e := range entities {
		contexts = append(contexts, e.toModel())
	}

	return contexts, nil
}

func (r *ContextRepository) ListByWorkspaceIncludingArchived(workspaceId string) ([]*core.Context, error) {
	entities := []*ContextEntity{}
	if err := r.db.Preload("Tags").Preload("ProjectMetadata").Find(&entities, "workspace_id = ?", workspaceId).Error; err != nil {
		return nil, err
	}

	contexts := []*core.Context{}
	for _, e := range entities {
		contexts = append(contexts, e.toModel())
	}

	return contexts, nil
}

func (r *ContextRepository) ListToSync(limit int) ([]*core.Context, error) {
	return r.List()
}

func (r *ContextRepository) SaveAll(contexts []*core.Context) ([]string, error) {
	ids := make([]string, len(contexts))

	if err := r.db.Transaction(func(tx *gorm.DB) error {
		for i, context := range contexts {
			id, err := save(tx, context)
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

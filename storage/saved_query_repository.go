package storage

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/m87/ctx/core"
	"gorm.io/gorm"
)

type SavedQueryRepository struct {
	db *gorm.DB
}

func NewSavedQueryRepository(db *gorm.DB) *SavedQueryRepository {
	return &SavedQueryRepository{db: db}
}

func (r *SavedQueryRepository) GetById(id string) (*core.SavedQuery, error) {
	var entity SavedQueryEntity
	result := r.db.Where("id = ?", id).Limit(1).Find(&entity)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return entity.ToModel(), nil
}

func (r *SavedQueryRepository) Save(query *core.SavedQuery) (string, error) {
	if query == nil {
		return "", fmt.Errorf("saved query is required")
	}
	if query.Id == "" {
		query.Id = uuid.NewString()
	}

	entity := NewSavedQueryEntity(query)
	if err := r.db.Omit("Workspace").Save(entity).Error; err != nil {
		return "", err
	}
	return entity.Id, nil
}

func (r *SavedQueryRepository) ListByWorkspace(workspaceId string) ([]*core.SavedQuery, error) {
	var entities []*SavedQueryEntity
	if err := r.db.Where("workspace_id = ?", workspaceId).Order("name ASC, id ASC").Find(&entities).Error; err != nil {
		return nil, err
	}

	queries := make([]*core.SavedQuery, len(entities))
	for i, entity := range entities {
		queries[i] = entity.ToModel()
	}
	return queries, nil
}

func (r *SavedQueryRepository) Delete(id string) error {
	return r.db.Delete(&SavedQueryEntity{}, "id = ?", id).Error
}

package storage

import (
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
	if err := r.db.First(&context, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return context.toModel(), nil
}


func (r *ContextRepository) ListByProject(projectId string) ([]*core.Context, error) {
	if projectId == "" {
		return nil, nil
	}

	entities := []*ContextEntity{}
	if err := r.db.Find(&entities, "projectId = ?", projectId).Error; err != nil {
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
	if context.Id == "" {
		context.Id = uuid.NewString()
	}
	
	if err := tx.Save(NewContextEntity(context)).Error; err != nil {
		return "", err
	}
	return context.Id, nil
}


func (r *ContextRepository) Delete(id string) error {
	return r.db.Delete(&ContextEntity{}, "id = ?", id).Error
}

func (r *ContextRepository) List() ([]*core.Context, error) {
	entities := []*ContextEntity{}
	if err := r.db.Find(&entities).Error; err != nil {
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
	if err := r.db.Where("status = ?", "active").First(entity).Error; err != nil {
		return nil, err
	}

	return entity.toModel(), nil
}

func (r *ContextRepository) ListByWorkspace(workspaceId string) ([]*core.Context, error) {
	entities := []*ContextEntity{}
	if err := r.db.Find(&entities, "workspaceId = ? AND archived = ?", workspaceId, false).Error; err != nil {
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
	if err := r.db.Find(&entities, "workspaceId = ?", workspaceId).Error; err != nil {
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
	
	r.db.Transaction(func(tx *gorm.DB) error {
		for i, context := range contexts {
			id, err := save(tx, context)
			if err != nil {
				return err
			}

			ids[i] = id
		}
		return nil
	})

	return ids, nil
}


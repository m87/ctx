package storage

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/m87/ctx/core"
	"gorm.io/gorm"
)

type DashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

func (r *DashboardRepository) GetById(id string) (*core.Dashboard, error) {
	var entity DashboardEntity
	result := r.db.Where("id = ?", id).Limit(1).Find(&entity)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return entity.ToModel(), nil
}

func (r *DashboardRepository) GetByTarget(
	workspaceId string,
	dashboardType core.DashboardType,
	targetId string,
) (*core.Dashboard, error) {
	var entity DashboardEntity
	result := r.db.Where(
		"workspace_id = ? AND type = ? AND target_id = ?",
		workspaceId,
		dashboardType,
		targetId,
	).Limit(1).Find(&entity)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return entity.ToModel(), nil
}

func (r *DashboardRepository) Save(dashboard *core.Dashboard) (string, error) {
	if dashboard == nil {
		return "", fmt.Errorf("dashboard is required")
	}
	if dashboard.Id == "" {
		dashboard.Id = uuid.NewString()
	}

	entity := NewDashboardEntity(dashboard)
	if err := r.db.Omit("Workspace").Save(entity).Error; err != nil {
		return "", err
	}
	return entity.Id, nil
}

func (r *DashboardRepository) ListByWorkspace(workspaceId string) ([]*core.Dashboard, error) {
	var entities []*DashboardEntity
	if err := r.db.Where("workspace_id = ?", workspaceId).Order("name ASC, id ASC").Find(&entities).Error; err != nil {
		return nil, err
	}

	dashboards := make([]*core.Dashboard, len(entities))
	for i, entity := range entities {
		dashboards[i] = entity.ToModel()
	}
	return dashboards, nil
}

func (r *DashboardRepository) Delete(id string) error {
	return r.db.Delete(&DashboardEntity{}, "id = ?", id).Error
}

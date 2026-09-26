package storage

import (
	"encoding/json"

	"github.com/m87/ctx/core"
)

type DashboardEntity struct {
	Id          string             `gorm:"primaryKey"`
	WorkspaceId string             `gorm:"not null;index;uniqueIndex:idx_dashboard_target,priority:1"`
	Type        core.DashboardType `gorm:"not null;uniqueIndex:idx_dashboard_target,priority:2"`
	TargetId    *string            `gorm:"uniqueIndex:idx_dashboard_target,priority:3"`
	Name        string             `gorm:"not null"`
	Definition  string             `gorm:"not null;type:text"`
	Workspace   WorkspaceEntity    `gorm:"foreignKey:WorkspaceId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (DashboardEntity) TableName() string {
	return "dashboards"
}

func NewDashboardEntity(dashboard *core.Dashboard) *DashboardEntity {
	var targetId *string
	if dashboard.TargetId != "" {
		target := dashboard.TargetId
		targetId = &target
	}
	return &DashboardEntity{
		Id:          dashboard.Id,
		WorkspaceId: dashboard.WorkspaceId,
		Type:        dashboard.Type,
		TargetId:    targetId,
		Name:        dashboard.Name,
		Definition:  string(dashboard.Definition),
	}
}

func (e *DashboardEntity) ToModel() *core.Dashboard {
	targetId := ""
	if e.TargetId != nil {
		targetId = *e.TargetId
	}
	definition := json.RawMessage(e.Definition)
	return &core.Dashboard{
		Id:          e.Id,
		WorkspaceId: e.WorkspaceId,
		Type:        e.Type,
		TargetId:    targetId,
		Name:        e.Name,
		Definition:  definition,
	}
}

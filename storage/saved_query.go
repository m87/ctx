package storage

import "github.com/m87/ctx/core"

type SavedQueryEntity struct {
	Id          string          `gorm:"primaryKey"`
	WorkspaceId string          `gorm:"not null;index"`
	Name        string          `gorm:"not null"`
	Query       string          `gorm:"not null"`
	Workspace   WorkspaceEntity `gorm:"foreignKey:WorkspaceId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (SavedQueryEntity) TableName() string {
	return "saved_queries"
}

func NewSavedQueryEntity(query *core.SavedQuery) *SavedQueryEntity {
	return &SavedQueryEntity{
		Id:          query.Id,
		WorkspaceId: query.WorkspaceId,
		Name:        query.Name,
		Query:       query.Query,
	}
}

func (e *SavedQueryEntity) ToModel() *core.SavedQuery {
	return &core.SavedQuery{
		Id:          e.Id,
		WorkspaceId: e.WorkspaceId,
		Name:        e.Name,
		Query:       e.Query,
	}
}

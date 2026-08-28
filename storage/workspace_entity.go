package storage

import "github.com/m87/ctx/core"

type WorkspaceEntity struct {
	Id         string `gorm:"primaryKey"`
	Name       string
	Description string
}

func (WorkspaceEntity) TableName() string {
	return "workspaces"
}

func (e *WorkspaceEntity) ToModel() *core.Workspace {
	return &core.Workspace{
		Id:          e.Id,
		Name:        e.Name,
		Description: e.Description,
	}
}

func NewWorkspaceEntityFromModel(workspace *core.Workspace) *WorkspaceEntity {
	return &WorkspaceEntity{
		Id:          workspace.Id,
		Name:        workspace.Name,
		Description: workspace.Description,
	}
}	


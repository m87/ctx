package storage

import "github.com/m87/ctx/core"

type ProjectEntity struct {
	Id          string `gorm:"primaryKey"`
	Name        string
	ParentId    string `gorm:"index"`
	WorkspaceId string `gorm:"not null;index"`
}

func (ProjectEntity) TableName() string {
	return "projects"
}

func (e *ProjectEntity) ToModel() *core.Project {
	return &core.Project{
		Id:          e.Id,
		Name:        e.Name,
		WorkspaceId: e.WorkspaceId,
		ParentId:    e.ParentId,
	}
}

func NewProjectEntityFromModel(project *core.Project) *ProjectEntity {
	return &ProjectEntity{
		Id:          project.Id,
		Name:        project.Name,
		WorkspaceId: project.WorkspaceId,
		ParentId:    project.ParentId,
	}
}

package storage

import "github.com/m87/ctx/core"

type ContextEntity struct {
	Id string `gorm:"primaryKey"`
	Name string `gorm:"not null"`
	WorkspaceId string `gorm:"not null;index"`
	Status string `gorm:"not null"`
	Archived bool `gorm:"not null"`
	Description string
	ProjectId *string `gorm:"index"`
	ProjectMetadata ProjectMetadata `gorm:"foreignKey:ProjectId"`
	Tags []*TagEntity `gorm:"many2many:context_tags"`
}

type ProjectMetadata struct {
	Id string
	Name string
}


func (ContextEntity) TableName() string {
	return "contexts"
}

func NewContextEntity(context *core.Context) *ContextEntity {

	tags := []*TagEntity{}
	if context.Tags != nil {
		for _, tag := range context.Tags {
			tags = append(tags, NewTagEntity(tag))
		}
	}

	projectPtr := context.ProjectId
	if projectPtr == nil || *projectPtr == "" {
		projectPtr = nil
	}

	return &ContextEntity {
		Id: context.Id,
		Name: context.Name,
		ProjectId: projectPtr,
		WorkspaceId: context.WorkspaceId,
		Tags: tags,
		Status: context.Status,
		Archived: context.Archived,
		Description: context.Description,
	}
}

func (c *ContextEntity) toModel() *core.Context {
	tags := []*core.Tag{}

	if c.Tags != nil {
		for _, tag := range c.Tags {
			tags = append(tags, tag.ToModel())
		}
	}

	projectMetadata := &core.ProjectMetadata{
		Id: c.ProjectMetadata.Id,
		Name: c.ProjectMetadata.Name,
	}

	return &core.Context{
		Id: c.Id,
		Name: c.Id,
		WorkspaceId: c.WorkspaceId,
		ProjectId: c.ProjectId,
		Tags: tags,
		Project: projectMetadata,
		Status: c.Status,
		Archived: c.Archived,
		Description: c.Description,
	}
}

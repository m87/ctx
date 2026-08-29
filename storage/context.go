package storage

import "github.com/m87/ctx/core"

type ContextEntity struct {
	Id              string `gorm:"primaryKey"`
	Name            string `gorm:"not null"`
	WorkspaceId     string `gorm:"not null;index"`
	Status          string `gorm:"not null"`
	Archived        bool   `gorm:"not null"`
	Description     string
	ProjectId       *string          `gorm:"index"`
	ProjectMetadata *ProjectMetadata `gorm:"foreignKey:ProjectId;references:Id"`
	Tags            []*TagEntity     `gorm:"many2many:context_tags;joinForeignKey:ContextId;joinReferences:TagId"`
}

type ProjectMetadata struct {
	Id   string `gorm:"primaryKey"`
	Name string
}

func (ProjectMetadata) TableName() string {
	return "projects"
}

func (ContextEntity) TableName() string {
	return "contexts"
}

func NewContextEntity(context *core.Context) *ContextEntity {
	tags := []*TagEntity{}
	for _, tag := range context.Tags {
		if tag != nil {
			tags = append(tags, NewTagEntity(tag))
		}
	}

	var projectId *string
	if context.Project != nil && context.Project.Id != "" {
		id := context.Project.Id
		projectId = &id
	} else if context.ProjectId != nil && *context.ProjectId != "" {
		id := *context.ProjectId
		projectId = &id
	}

	return &ContextEntity{
		Id:          context.Id,
		Name:        context.Name,
		ProjectId:   projectId,
		WorkspaceId: context.WorkspaceId,
		Tags:        tags,
		Status:      context.Status,
		Archived:    context.Archived,
		Description: context.Description,
	}
}

func (c *ContextEntity) toModel() *core.Context {
	tags := []*core.Tag{}
	for _, tag := range c.Tags {
		if tag != nil {
			tags = append(tags, tag.ToModel())
		}
	}

	var project *core.ProjectMetadata
	if c.ProjectMetadata != nil {
		project = &core.ProjectMetadata{
			Id:   c.ProjectMetadata.Id,
			Name: c.ProjectMetadata.Name,
		}
	}

	return &core.Context{
		Id:          c.Id,
		Name:        c.Name,
		WorkspaceId: c.WorkspaceId,
		ProjectId:   c.ProjectId,
		Tags:        tags,
		Project:     project,
		Status:      c.Status,
		Archived:    c.Archived,
		Description: c.Description,
	}
}

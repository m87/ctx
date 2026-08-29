package storage

import (
	"github.com/google/uuid"
	"github.com/m87/ctx/core"
)

type TagEntity struct {
	Id   string `gorm:"primaryKey"`
	Name string `gorm:"not null"`
}

type ContextTagEntity struct {
	ContextId string `gorm:"primaryKey"`
	TagId     string `gorm:"primaryKey"`
}

func (TagEntity) TableName() string {
	return "tag"
}

func (ContextTagEntity) TableName() string {
	return "context_tags"
}

func NewTagEntity(tag *core.Tag) *TagEntity {
	if tag.Id == "" {
		tag.Id = uuid.NewString()
	}
	return &TagEntity{
		Id:   tag.Id,
		Name: tag.Name,
	}
}

func (e *TagEntity) ToModel() *core.Tag {
	return &core.Tag{
		Id:   e.Id,
		Name: e.Name,
	}
}

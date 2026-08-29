package storage

import (
	"time"

	"github.com/m87/ctx/core"
)

type IntervalEntity struct {
	Id string `gorm:"primaryKey"`
	ContextId string `gorm:"not null;index"`
	Start *time.Time `gorm:"not null"`
	End *time.Time
	Duration time.Duration `gorm:"not null"`
	Status string `gorm:"not null;index"`
	WorkspaceId string `gorm:"not null;index"`
}

func (IntervalEntity) TableName() string {
	return "intervals"
}

func (e *IntervalEntity) ToModel() *core.Interval {
	return &core.Interval{
		Id:        e.Id,
		ContextId: e.ContextId,
		Start: e.Start,
		End: e.End,
		Duration: e.Duration,
		Status: e.Status,
		WorkspaceId: e.WorkspaceId,
	}
}

func NewIntervalEntityFromModel(interval *core.Interval) *IntervalEntity {
	return &IntervalEntity{
		Id:        interval.Id,
		ContextId: interval.ContextId,
		Start: interval.Start,
		End: interval.End,
		Duration: interval.Duration,
		Status: interval.Status,
		WorkspaceId: interval.WorkspaceId,
	}
}

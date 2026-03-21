package core

import "time"

type IntervalRepository interface {
	GetById(id string) (*Interval, error)
	Save(interval *Interval) (string, error)
	Delete(id string) error
	ListByContextId(contextId string) ([]*Interval, error)
	GetActiveIntervalByContextId(contextId string) (*Interval, error)
	ListByDay(date time.Time) ([]*Interval, error)
}

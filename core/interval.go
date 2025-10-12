package core

import (
	"time"

	ctxtime "github.com/m87/ctx/time"
)

type Interval struct {
	Id       string            `json:"id"`
	Start    ctxtime.ZonedTime `json:"start"`
	End      ctxtime.ZonedTime `json:"end"`
	Duration time.Duration     `json:"duration"`
	Labels   []string          `json:"labels"`
}

func (interval *Interval) IsActive() bool {
	return interval.End.Time.IsZero()
}

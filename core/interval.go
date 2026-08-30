package core

import (
	"time"
)

const IntervalStatusActive = "active"

type Interval struct {
	Id          string        `json:"id"`
	ContextId   string        `json:"contextId"`
	Start       *time.Time    `json:"start"`
	End         *time.Time    `json:"end"`
	Duration    time.Duration `json:"duration"`
	Status      string        `json:"status"`
	WorkspaceId string        `json:"workspaceId"`
	Synced      bool          `json:"synced"`
}

const IntervalType = "interval"

func timeIsSet(value *time.Time) bool {
	return value != nil && !value.IsZero()
}

func utcTimePointer(value *time.Time) *time.Time {
	if !timeIsSet(value) {
		return nil
	}
	utc := value.UTC()
	return &utc
}

func durationBetween(start, end *time.Time) time.Duration {
	if !timeIsSet(start) || !timeIsSet(end) || !end.After(*start) {
		return 0
	}
	return end.Sub(*start)
}

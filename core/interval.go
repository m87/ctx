package core

import "time"

type Interval struct {
	Id        string    `json:"id"`
	ContextId string    `json:"contextId"`
	Start     time.Time `json:"start"`
	End       time.Time `json:"end"`
}

type IntervalMapper struct {
}

func NewIntervalMapper() *IntervalMapper {
	return &IntervalMapper{}
}

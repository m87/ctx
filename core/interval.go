package core

import (
	"time"

	"github.com/m87/nod"
)

type Interval struct {
	Id        string        `json:"id"`
	ContextId string        `json:"contextId"`
	Start     ZonedTime     `json:"start"`
	End       ZonedTime     `json:"end"`
	Duration  time.Duration `json:"duration"`
	Status    string        `json:"status"`
}

type IntervalMapper struct {
}

const IntervalType = "interval"

func NewIntervalMapper() *IntervalMapper {
	return &IntervalMapper{}
}

func (m *IntervalMapper) ToNode(interval *Interval) (*nod.Node, error) {
	durationNanos := interval.Duration.Nanoseconds()
	node := &nod.Node{
		Core: nod.NodeCore{
			Id:        interval.Id,
			Name:      interval.Id,
			Kind:      IntervalType,
			ParentId:  &interval.ContextId,
			Status:    interval.Status,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		KV: map[string]*nod.KV{
			"start":          &nod.KV{Key: "start", ValueTime: &interval.Start.Time},
			"start_timezone": &nod.KV{Key: "start_timezone", ValueText: &interval.Start.Timezone},
			"end":            &nod.KV{Key: "end", ValueTime: &interval.End.Time},
			"end_timezone":   &nod.KV{Key: "end_timezone", ValueText: &interval.End.Timezone},
			"duration":       &nod.KV{Key: "duration", ValueInt64: &durationNanos},
		},
	}
	return node, nil
}

func (m *IntervalMapper) FromNode(node *nod.Node) (*Interval, error) {
	return &Interval{
		Id:        node.Core.Id,
		ContextId: *node.Core.ParentId,
		Start: ZonedTime{
			Time:     nod.SafeTime(node.KV, "start"),
			Timezone: nod.SafeString(node.KV, "start_timezone"),
		},
		End: ZonedTime{
			Time:     nod.SafeTime(node.KV, "end"),
			Timezone: nod.SafeString(node.KV, "end_timezone"),
		},
		Duration: time.Duration(nod.SafeInt64(node.KV, "duration")),
		Status:   node.Core.Status,
	}, nil
}

func (m *IntervalMapper) IsApplicable(node *nod.Node) bool {
	return node.Core.Kind == IntervalType
}

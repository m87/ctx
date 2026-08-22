package core

import (
	"time"

	"github.com/m87/nod"
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

type IntervalMapper struct {
}

const IntervalType = "interval"

func NewIntervalMapper() *IntervalMapper {
	return &IntervalMapper{}
}

func (m *IntervalMapper) ToNode(interval *Interval) (*nod.Node, error) {
	durationNanos := interval.Duration.Nanoseconds()
	kv := map[string]*nod.NodeKV{
		"duration": {Key: "duration", ValueInt64: &durationNanos},
		"synced":   {Key: "synced", ValueBool: &interval.Synced},
	}
	if timeIsSet(interval.Start) {
		start := interval.Start.UTC()
		kv["start"] = &nod.NodeKV{Key: "start", ValueTime: &start}
	}
	if timeIsSet(interval.End) {
		end := interval.End.UTC()
		kv["end"] = &nod.NodeKV{Key: "end", ValueTime: &end}
	}

	node := &nod.Node{
		Core: nod.NodeCore{
			Id:          interval.Id,
			Name:        interval.Id,
			Kind:        IntervalType,
			ParentId:    stringPointerIfNotEmpty(interval.ContextId),
			NamespaceId: stringPointerIfNotEmpty(interval.WorkspaceId),
			Status:      interval.Status,
		},
		KV: kv,
	}
	return node, nil
}

func (m *IntervalMapper) FromNode(node *nod.Node) (*Interval, error) {
	contextId := ""
	if node.Core.ParentId != nil {
		contextId = *node.Core.ParentId
	}

	workspaceId := ""
	if node.Core.NamespaceId != nil {
		workspaceId = *node.Core.NamespaceId
	}
	start := nodTime(node.KV, "start").UTC()
	end := nodTime(node.KV, "end").UTC()
	var startPointer *time.Time
	if !start.IsZero() {
		startPointer = &start
	}
	var endPointer *time.Time
	if !end.IsZero() {
		endPointer = &end
	}
	return &Interval{
		Id:          node.Core.Id,
		ContextId:   contextId,
		Start:       startPointer,
		End:         endPointer,
		Duration:    time.Duration(nodInt64(node.KV, "duration")),
		Status:      node.Core.Status,
		WorkspaceId: workspaceId,
		Synced:      nodBool(node.KV, "synced"),
	}, nil
}

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

func (m *IntervalMapper) IsApplicable(node *nod.Node) bool {
	return node.Core.Kind == IntervalType
}

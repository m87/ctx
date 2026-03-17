package core

import (
	"time"

	"github.com/m87/nod"
)

type Interval struct {
	Id        string    `json:"id"`
	ContextId string    `json:"contextId"`
	Start     time.Time `json:"start"`
	End       time.Time `json:"end"`
}

type IntervalMapper struct {
}

const IntervalType = "interval"

func NewIntervalMapper() *IntervalMapper {
	return &IntervalMapper{}
}

func (m *IntervalMapper) ToNode(interval *Interval) (*nod.Node, error) {
	node := &nod.Node{
		Core: nod.NodeCore{
			Id:   interval.Id,
			Name: interval.Id,
			Kind: IntervalType,
		},
		KV: map[string]*nod.KV{
			"contextId": &nod.KV{Key: "contextId", ValueText: &interval.ContextId},
			"start":     &nod.KV{Key: "start", ValueTime: &interval.Start},
			"end":       &nod.KV{Key: "end", ValueTime: &interval.End},
		},
	}
	return node, nil
}

func (m *IntervalMapper) FromNode(node *nod.Node) (*Interval, error) {
	return &Interval{
		Id:        node.Core.Id,
		ContextId: *node.KV["contextId"].ValueText,
		Start:     *node.KV["start"].ValueTime,
		End:       *node.KV["end"].ValueTime,
	}, nil
}

func (m *IntervalMapper) IsApplicable(node *nod.Node) bool {
	return node.Core.Kind == IntervalType
}

package core

import (
	"time"

	"github.com/m87/nod"
)

type Workspace struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type WorkspaceContextStats struct {
	ContextId     string        `json:"contextId"`
	Duration      time.Duration `json:"duration"`
	Percentage    float64       `json:"percentage"`
	IntervalCount int           `json:"intervalCount"`
}

type WorkspaceStats struct {
	WorkspaceId   string                   `json:"workspaceId"`
	Contexts      []*Context               `json:"contexts"`
	ContextStats  []*WorkspaceContextStats `json:"contextStats"`
	TotalDuration time.Duration            `json:"totalDuration"`
	TotalSessions int                      `json:"totalSessions"`
}

type WorkspaceMapper struct {
}

func NewWorkspaceMapper() *WorkspaceMapper {
	return &WorkspaceMapper{}
}

func (m *WorkspaceMapper) ToNode(workspace *Workspace) (*nod.Node, error) {
	node := &nod.Node{
		Core: nod.NodeCore{
			Id:   workspace.Id,
			Name: workspace.Name,
			Kind: WorkspaceType,
		},
	}

	node.Content = ConvertToNodContent(map[string]string{
		"description": workspace.Description,
	})

	return node, nil
}

func (m *WorkspaceMapper) FromNode(node *nod.Node) (*Workspace, error) {
	return &Workspace{
		Id:          node.Core.Id,
		Name:        node.Core.Name,
		Description: ConvertFromNodContent(node.Content)["description"],
	}, nil
}

func (m *WorkspaceMapper) IsApplicable(node *nod.Node) bool {
	return node.Core.Kind == WorkspaceType
}

const WorkspaceType = "workspace"

package core

import "github.com/m87/nod"

type Workspace struct {
	Id   string `json:"id"`
	Name string `json:"name"`
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
	return node, nil
}

func (m *WorkspaceMapper) FromNode(node *nod.Node) (*Workspace, error) {
	return &Workspace{
		Id:   node.Core.Id,
		Name: node.Core.Name,
	}, nil
}

func (m *WorkspaceMapper) IsApplicable(node *nod.Node) bool {
	return node.Core.Kind == WorkspaceType
}

const WorkspaceType = "workspace"

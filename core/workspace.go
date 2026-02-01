package core

import "github.com/m87/nod"

const WorkspaceType = "workspace"

type Workspace struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (workspace *Workspace) Type() string {
	return WorkspaceType
}

func (workspace *Workspace) Kind() string {
	return ""
}

func NewWorkspace(name string) *Workspace {
	return &Workspace{
		Name: name,
	}
}

type WorkspaceMapper struct{}

func (wm *WorkspaceMapper) ToNode(node nod.NodeModel) (*nod.Node, error) {
	workspace, ok := node.(*Workspace)
	if !ok {
		return nil, nil
	}

	output := &nod.Node{
		Core: nod.NodeCore{
			Id:   workspace.Id,
			Type: workspace.Type(),
			Name: workspace.Name,
			Kind: workspace.Kind(),
		},
	}
	return output, nil
}

func (wm *WorkspaceMapper) FromNode(node *nod.Node) (nod.NodeModel, error) {
	workspace := &Workspace{
		Id:   node.Core.Id,
		Name: node.Core.Name,
	}
	return workspace, nil
}

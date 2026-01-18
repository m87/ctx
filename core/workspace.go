package core

import "github.com/m87/nod"

const WorkspaceType = "workspace"

type Workspace struct {
	Id          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Tags        []string          `json:"tags"`
	Properties  map[string]string `json:"properties"`
}

func (workspace *Workspace) Type() string {
	return WorkspaceType
}

func (workspace *Workspace) Kind() string {
	return ""
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
		Tags: ConvertToNodTags(workspace.Tags),
		KV:   ConvertToNodKV(workspace.Properties),
		Content: ConvertToNodContent(map[string]string{
			"description": workspace.Description,
		}),
	}
	return output, nil
}

func (wm *WorkspaceMapper) FromNode(node *nod.Node) (nod.NodeModel, error) {
	workspace := &Workspace{
		Id:   node.Core.Id,
		Name: node.Core.Name,
	}

	if desc, ok := node.Content["description"]; ok {
		workspace.Description = *desc.Value
	}

	workspace.Tags = ConvertFromNodTags(node.Tags)
	workspace.Properties = ConvertFromNodKV(node.KV)

	return workspace, nil
}



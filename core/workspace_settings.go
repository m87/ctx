package core

import "github.com/m87/nod"

type WorkspaceSettings struct {
	Id string
	WorkspaceId string
	LinkRules []LinkRule
}



func (m *WorkspaceSettings) ToNode() *nod.Node {
	node := &nod.Node{
		Core: nod.NodeCore{
			Id:   m.Id,
			ParentId: nod.Ptr(m.WorkspaceId),
			Kind: "workspace-settings",
		},
	}


	return node
}

func (m *WorkspaceSettings) FromNode(node *nod.Node) {
	m.Id = node.Core.Id
	if node.Core.ParentId != nil {
		m.WorkspaceId = *node.Core.ParentId
	}


}

func (m *WorkspaceSettings) IsApplicable(node *nod.Node) bool {
	return node.Core.Kind == "workspace-settings"
}

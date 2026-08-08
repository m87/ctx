package core

import "github.com/m87/nod"

const (
	ProjectType = "project"
)

type Project struct {
	Id string
	Name string
	ParentId string
	WorkspaceId string
}

func (p *Project) ToNode() (*nod.Node, error) {
	return &nod.Node{
		Core: nod.NodeCore{
			Kind: ProjectType,
			Id: p.Id,
			Name: p.Name,
			ParentId: stringPointerIfNotEmpty(p.ParentId),
			NamespaceId: stringPointerIfNotEmpty(p.WorkspaceId),
		},
	}, nil
}


func (p *Project) FromNode(node *nod.Node) error {
	p.Id = node.Core.Id
	p.Name = node.Core.Name

	if node.Core.ParentId != nil {
		p.ParentId = *node.Core.ParentId
	}

	if node.Core.NamespaceId != nil {
		p.WorkspaceId = *node.Core.NamespaceId
	}



	return nil
}

func (p *Project) IsApplicablet(node *nod.Node) bool {
	return node.Core.Kind == ProjectType
}

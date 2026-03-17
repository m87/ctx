package core

import "github.com/m87/nod"

type Context struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	ParentId string `json:"parentId"`
}

const ContextType = "context"

type ContextMapper struct {
}

func NewContextMapper() *ContextMapper {
	return &ContextMapper{}
}

func (m *ContextMapper) ToNode(context *Context) (*nod.Node, error) {
	node := &nod.Node{
		Core: nod.NodeCore{
			Id:       context.Id,
			Name:     context.Name,
			ParentId: &context.ParentId,
			Kind:     ContextType,
		},
	}
	return node, nil
}

func (m *ContextMapper) FromNode(node *nod.Node) (*Context, error) {
	return &Context{
		Id:       node.Core.Id,
		Name:     node.Core.Name,
		ParentId: *node.Core.ParentId,
	}, nil
}

func (m *ContextMapper) IsApplicable(node *nod.Node) bool {
	return node.Core.Kind == ContextType
}

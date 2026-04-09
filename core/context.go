package core

import "github.com/m87/nod"

type Context struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	ParentId    string   `json:"parentId"`
	Status      string   `json:"status"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

const ContextType = "context"

type ContextMapper struct {
}

func NewContextMapper() *ContextMapper {
	return &ContextMapper{}
}

func (m *ContextMapper) ToNode(context *Context) (*nod.Node, error) {
	var parentId *string
	if context.ParentId != "" {
		parentId = &context.ParentId
	}

	node := &nod.Node{
		Core: nod.NodeCore{
			Id:       context.Id,
			Name:     context.Name,
			ParentId: parentId,
			Kind:     ContextType,
			Status:   context.Status,
		},
	}

	node.Content = nod.ConvertStringMapToContent(map[string]string{
		"description": context.Description,
	})

	node.Tags = nod.ConvertStringSliceToTags(context.Tags)

	return node, nil
}

func (m *ContextMapper) FromNode(node *nod.Node) (*Context, error) {
	return &Context{
		Id:          node.Core.Id,
		Name:        node.Core.Name,
		ParentId:    *node.Core.ParentId,
		Status:      node.Core.Status,
		Description: nod.ConvertContentToStringMap(node.Content)["description"],
		Tags:        nod.ConvertTagsToStringSlice(node.Tags),
	}, nil
}

func (m *ContextMapper) IsApplicable(node *nod.Node) bool {
	return node.Core.Kind == ContextType
}

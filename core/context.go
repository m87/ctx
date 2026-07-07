package core

import "github.com/m87/nod"

type Context struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	ParentId    string   `json:"parentId"`
	WorkspaceId string   `json:"workspaceId"`
	Status      string   `json:"status"`
	Archived    bool     `json:"archived"`
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
	node := &nod.Node{
		Core: nod.NodeCore{
			Id:          context.Id,
			Name:        context.Name,
			ParentId:    stringPointerIfNotEmpty(context.ParentId),
			NamespaceId: stringPointerIfNotEmpty(context.WorkspaceId),
			Kind:        ContextType,
			Status:      context.Status,
		},
	}

	node.Content = nod.ConvertStringMapToContent(map[string]string{
		"description": context.Description,
	})
	node.KV = map[string]*nod.KV{
		"archived": {Key: "archived", ValueBool: &context.Archived},
	}

	node.Tags = nod.ConvertStringSliceToTags(context.Tags)

	return node, nil
}

func (m *ContextMapper) FromNode(node *nod.Node) (*Context, error) {
	parentId := ""
	if node.Core.ParentId != nil {
		parentId = *node.Core.ParentId
	}

	workspaceId := ""
	if node.Core.NamespaceId != nil {
		workspaceId = *node.Core.NamespaceId
	}

	return &Context{
		Id:          node.Core.Id,
		Name:        node.Core.Name,
		ParentId:    parentId,
		WorkspaceId: workspaceId,
		Status:      node.Core.Status,
		Archived:    nod.SafeBool(node.KV, "archived"),
		Description: nod.ConvertContentToStringMap(node.Content)["description"],
		Tags:        nod.ConvertTagsToStringSlice(node.Tags),
	}, nil
}

func (m *ContextMapper) IsApplicable(node *nod.Node) bool {
	return node.Core.Kind == ContextType
}

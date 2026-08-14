package core

import (
	"github.com/m87/nod"
)

type ProjectMetadata struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Context struct {
	Id          string           `json:"id"`
	Name        string           `json:"name"`
	ParentId    string           `json:"parentId"`
	WorkspaceId string           `json:"workspaceId"`
	Status      string           `json:"status"`
	Archived    bool             `json:"archived"`
	Description string           `json:"description,omitempty"`
	Tags        []string         `json:"tags,omitempty"`
	Synced      bool             `json:"synced"`
	Project     *ProjectMetadata `json:"project"`
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

	node.Content = ConvertToNodContent(map[string]string{
		"description": context.Description,
	})
	node.KV = map[string]*nod.NodeKV{
		"archived": {Key: "archived", ValueBool: &context.Archived},
		"synced":   {Key: "synced", ValueBool: &context.Synced},
	}

	if context.Project != nil {
		if context.Project.Id != "" {
			node.KV["projectId"] = &nod.NodeKV{Key: "projectId", ValueText: &context.Project.Id}
		}
		if context.Project.Name != "" {
			node.KV["projectName"] = &nod.NodeKV{Key: "projectName", ValueText: &context.Project.Name}
		}
	}

	node.Tags = ConvertToNodTags(context.Tags)

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

	var project *ProjectMetadata
	if projectIdKV, ok := node.KV["projectId"]; ok && projectIdKV.ValueText != nil {
		project = &ProjectMetadata{}
		project.Id = *projectIdKV.ValueText
		if projectNameKV, ok := node.KV["projectName"]; ok && projectNameKV.ValueText != nil {
			project.Name = *projectNameKV.ValueText
		}
	}

	return &Context{
		Id:          node.Core.Id,
		Name:        node.Core.Name,
		ParentId:    parentId,
		WorkspaceId: workspaceId,
		Status:      node.Core.Status,
		Archived:    nodBool(node.KV, "archived"),
		Description: ConvertFromNodContent(node.Content)["description"],
		Tags:        ConvertFromNodTags(node.Tags),
		Synced:      nodBool(node.KV, "synced"),
		Project:     project,
	}, nil
}

func (m *ContextMapper) IsApplicable(node *nod.Node) bool {
	return node.Core.Kind == ContextType
}

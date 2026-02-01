package core

import (
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/assert"
)

func TestConvertWorkspaceToFromNode(t *testing.T) {
	workspaceNode := &nod.Node{
		Core: nod.NodeCore{
			Id:   "workspace-123",
			Type: WorkspaceType,
			Name: "Test Workspace",
			Kind: "",
		},
	}
	mapper := &WorkspaceMapper{}

	workspaceModel, err := mapper.FromNode(workspaceNode)

	assert.NoError(t, err)
	workspace, ok := workspaceModel.(*Workspace)
	assert.True(t, ok)
	assert.Equal(t, "workspace-123", workspace.Id)
	assert.Equal(t, "Test Workspace", workspace.Name)
}

func TestConvertWorkspaceToNode(t *testing.T) {
	workspace := &Workspace{
		Id:   "workspace-456",
		Name: "Another Workspace",
	}
	mapper := &WorkspaceMapper{}

	workspaceNode, err := mapper.ToNode(workspace)

	assert.NoError(t, err)
	assert.Equal(t, "workspace-456", workspaceNode.Core.Id)
	assert.Equal(t, WorkspaceType, workspaceNode.Core.Type)
	assert.Equal(t, "Another Workspace", workspaceNode.Core.Name)
}

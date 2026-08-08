package core

import (
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/assert"
)


func testProjectToNode(t *testing.T) {
	project := &Project{
		Id:          "project1",
		Name:        "Project 1",
		ParentId:    "parent1",
		WorkspaceId: "workspace1",
	}

	node, err := project.ToNode()
	assert.NoError(t, err)
	assert.Equal(t, "project1", node.Core.Id)
	assert.Equal(t, "Project 1", node.Core.Name)
	assert.Equal(t, "parent1", *node.Core.ParentId)
	assert.Equal(t, "workspace1", *node.Core.NamespaceId)
	assert.Equal(t, ProjectType, node.Core.Kind)
}

func testProjectFromNode(t *testing.T) {
	node := &nod.Node{
		Core: nod.NodeCore{
			Id:          "project1",
			Name:        "Project 1",
			ParentId:    stringPointerIfNotEmpty("parent1"),
			NamespaceId: stringPointerIfNotEmpty("workspace1"),
			Kind:        ProjectType,
		},
	}
	
	project := &Project{}
	err := project.FromNode(node)
	assert.NoError(t, err)	
	assert.Equal(t, "project1", project.Id)
	assert.Equal(t, "Project 1", project.Name)
	assert.Equal(t, "parent1", project.ParentId)
	assert.Equal(t, "workspace1", project.WorkspaceId)
}

func testProjectIsApplicable(t *testing.T) {
	project := &Project{}
	node := &nod.Node{
		Core: nod.NodeCore{
			Kind: ProjectType,
		},
	}
	assert.True(t, project.IsApplicable(node))

	node.Core.Kind = "other"
	assert.False(t, project.IsApplicable(node))
}

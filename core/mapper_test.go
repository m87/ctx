package core

import (
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

func TestContextMapperToNodeOmitsEmptyParent(t *testing.T) {
	mapper := NewContextMapper()

	node, err := mapper.ToNode(&Context{
		Id:          "context-1",
		Name:        "Context",
		WorkspaceId: "workspace-1",
	})

	require.NoError(t, err)
	require.Nil(t, node.Core.ParentId)
	require.NotNil(t, node.Core.NamespaceId)
	require.Equal(t, "workspace-1", *node.Core.NamespaceId)
}

func TestContextMapperToNodeKeepsParent(t *testing.T) {
	mapper := NewContextMapper()

	node, err := mapper.ToNode(&Context{
		Id:          "context-1",
		Name:        "Context",
		ParentId:    "parent-1",
		WorkspaceId: "workspace-1",
	})

	require.NoError(t, err)
	require.NotNil(t, node.Core.ParentId)
	require.Equal(t, "parent-1", *node.Core.ParentId)
}

func TestIntervalMapperFromNodeHandlesMissingParent(t *testing.T) {
	mapper := NewIntervalMapper()

	got, err := mapper.FromNode(&nod.Node{
		Core: nod.NodeCore{
			Id:     "interval-1",
			Name:   "interval-1",
			Kind:   IntervalType,
			Status: "completed",
		},
		KV: map[string]*nod.NodeKV{},
	})

	require.NoError(t, err)
	require.Equal(t, "interval-1", got.Id)
	require.Empty(t, got.ContextId)
	require.Empty(t, got.WorkspaceId)
}

func TestIntervalMapperToNodeOmitsEmptyRelations(t *testing.T) {
	mapper := NewIntervalMapper()

	node, err := mapper.ToNode(&Interval{
		Id:     "interval-1",
		Status: "completed",
	})

	require.NoError(t, err)
	require.Nil(t, node.Core.ParentId)
	require.Nil(t, node.Core.NamespaceId)
}

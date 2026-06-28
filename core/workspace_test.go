package core

import (
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

func TestWorkspaceMapperRoundTripsContextLinkRules(t *testing.T) {
	mapper := NewWorkspaceMapper()
	workspace := &Workspace{
		Id:          "workspace-1",
		Name:        "Workspace",
		Description: "Description",
		ContextLinkRules: []ContextLinkRule{
			{Regex: "[A-Z]+-\\d+", LinkTemplate: "https://jira.example.com/browse/$0"},
			{Regex: "GH-(\\d+)", LinkTemplate: "https://github.com/example/repo/issues/$1"},
		},
	}

	node, err := mapper.ToNode(workspace)
	require.NoError(t, err)

	got, err := mapper.FromNode(node)
	require.NoError(t, err)
	require.Equal(t, workspace, got)
}

func TestWorkspaceMapperReadsLegacyContextLinkRule(t *testing.T) {
	mapper := NewWorkspaceMapper()
	node := &nod.Node{
		Core: nod.NodeCore{Id: "workspace-1", Name: "Workspace", Kind: WorkspaceType},
		Content: nod.ConvertStringMapToContent(map[string]string{
			"description":         "Description",
			"contextLinkRegex":    "[A-Z]+-\\d+",
			"contextLinkTemplate": "https://jira.example.com/browse/$0",
		}),
	}

	got, err := mapper.FromNode(node)
	require.NoError(t, err)
	require.Equal(t, []ContextLinkRule{
		{Regex: "[A-Z]+-\\d+", LinkTemplate: "https://jira.example.com/browse/$0"},
	}, got.ContextLinkRules)
}

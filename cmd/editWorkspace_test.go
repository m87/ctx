package cmd

import (
	"bytes"
	"testing"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEditWorkspaceDescription(t *testing.T) {
	manager := newIsolatedTestManager(t)
	id, err := manager.WorkspaceRepository.Save(&core.Workspace{
		Name:        "demo",
		Description: "old description",
	})
	require.NoError(t, err)

	oldOutputFormat := OutputFormat
	oldRemoteAddr := RemoteAddr
	OutputFormat = "text"
	RemoteAddr = ""
	t.Cleanup(func() {
		OutputFormat = oldOutputFormat
		RemoteAddr = oldRemoteAddr
	})

	command := NewEditWorkspaceCmd()
	var output bytes.Buffer
	command.SetOut(&output)
	command.SetErr(&output)
	command.SetArgs([]string{"--id", id, "--description", "new description"})
	require.NoError(t, command.Execute())

	workspace, err := bootstrap.CreateManager().WorkspaceRepository.GetById(id)
	require.NoError(t, err)
	require.NotNil(t, workspace)
	assert.Equal(t, "new description", workspace.Description)

	command = NewEditWorkspaceCmd()
	command.SetOut(&output)
	command.SetErr(&output)
	command.SetArgs([]string{"--id", id, "--description", ""})
	require.NoError(t, command.Execute())

	workspace, err = bootstrap.CreateManager().WorkspaceRepository.GetById(id)
	require.NoError(t, err)
	require.NotNil(t, workspace)
	assert.Empty(t, workspace.Description)
}

package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func executeCreateWorkspaceCommand(t *testing.T, commandArgs ...string) (string, error) {
	t.Helper()

	command := NewCreateWorkspaceCmd()

	var output bytes.Buffer
	command.SetOut(&output)
	command.SetErr(&output)
	command.SetArgs(commandArgs)

	err := command.Execute()
	return strings.TrimSpace(output.String()), err
}

func TestCreateWorkspaceTextOutputAndPersistence(t *testing.T) {
	newIsolatedTestManager(t)

	oldOutputFormat := OutputFormat
	oldRemoteAddr := RemoteAddr
	OutputFormat = "text"
	RemoteAddr = ""
	t.Cleanup(func() {
		OutputFormat = oldOutputFormat
		RemoteAddr = oldRemoteAddr
	})

	output, err := executeCreateWorkspaceCommand(
		t,
		"--name", "demo",
		"--description", "Workspace for demo contexts",
	)
	require.NoError(t, err)
	assert.Equal(t, "Workspace created successfully", output)

	workspaces, err := bootstrap.CreateManager().WorkspaceRepository.List()
	require.NoError(t, err)
	var created *core.Workspace
	for _, workspace := range workspaces {
		if workspace != nil && workspace.Name == "demo" {
			created = workspace
			break
		}
	}
	require.NotNil(t, created)
	assert.Equal(t, "Workspace for demo contexts", created.Description)
}

func TestCreateWorkspaceJsonOutput(t *testing.T) {
	newIsolatedTestManager(t)

	oldOutputFormat := OutputFormat
	oldRemoteAddr := RemoteAddr
	OutputFormat = "json"
	RemoteAddr = ""
	t.Cleanup(func() {
		OutputFormat = oldOutputFormat
		RemoteAddr = oldRemoteAddr
	})

	output, err := executeCreateWorkspaceCommand(
		t,
		"--name", "demo",
		"--description", "Workspace for demo contexts",
	)
	require.NoError(t, err)

	var workspace core.Workspace
	require.NoError(t, json.Unmarshal([]byte(output), &workspace), "output: %q", output)
	assert.Equal(t, "demo", workspace.Name)
	assert.Equal(t, "Workspace for demo contexts", workspace.Description)
	assert.NotEmpty(t, workspace.Id)
}

func TestCreateWorkspaceYamlAndShellOutput(t *testing.T) {
	newIsolatedTestManager(t)

	oldOutputFormat := OutputFormat
	oldRemoteAddr := RemoteAddr
	RemoteAddr = ""
	t.Cleanup(func() {
		OutputFormat = oldOutputFormat
		RemoteAddr = oldRemoteAddr
	})

	OutputFormat = "yaml"
	yamlOutput, err := executeCreateWorkspaceCommand(t, "--name", "demo-yaml")
	require.NoError(t, err)
	assert.Contains(t, yamlOutput, "name: demo-yaml")

	OutputFormat = "shell"
	shellOutput, err := executeCreateWorkspaceCommand(t, "--name", "demo-shell")
	require.NoError(t, err)
	assert.Contains(t, shellOutput, "RESULT_NAME=\"demo-shell\"")
	assert.Contains(t, shellOutput, "RESULT_ID=")
}

func TestCreateWorkspaceRequiresNameFlag(t *testing.T) {
	newIsolatedTestManager(t)

	oldOutputFormat := OutputFormat
	oldRemoteAddr := RemoteAddr
	OutputFormat = "text"
	RemoteAddr = ""
	t.Cleanup(func() {
		OutputFormat = oldOutputFormat
		RemoteAddr = oldRemoteAddr
	})

	_, err := executeCreateWorkspaceCommand(t)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "required flag(s) \"name\" not set")
}

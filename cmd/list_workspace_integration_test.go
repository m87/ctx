package cmd

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/stretchr/testify/require"
)

func newIsolatedTestManager(t *testing.T) *core.ContextManager {
	t.Helper()

	t.Setenv("DATABASE_PATH", filepath.Join(t.TempDir(), "ctx-test.db"))
	manager, err := bootstrap.CreateManager()
	require.NoError(t, err)
	return manager
}

func executeWorkspaceListCommand(t *testing.T, commandArgs ...string) (string, error) {
	t.Helper()

	newIsolatedTestManager(t)
	command := NewListWorkspaceCmd()

	var output bytes.Buffer
	command.SetOut(&output)
	command.SetErr(&output)
	command.SetArgs(commandArgs)

	err := command.Execute()
	return strings.TrimSpace(output.String()), err
}

func TestListWorkspaceTextOutputWithResults(t *testing.T) {
	manager := newIsolatedTestManager(t)
	_, err := manager.WorkspaceRepository.Save(&core.Workspace{Name: "demo"})
	if err != nil {
		t.Fatalf("save workspace failed: %v", err)
	}

	oldOutputFormat := OutputFormat
	oldRemoteAddr := RemoteAddr
	OutputFormat = "text"
	RemoteAddr = ""
	t.Cleanup(func() {
		OutputFormat = oldOutputFormat
		RemoteAddr = oldRemoteAddr
	})

	command := NewListWorkspaceCmd()
	var output bytes.Buffer
	command.SetOut(&output)
	command.SetErr(&output)

	err = command.Execute()
	if err != nil {
		t.Fatalf("command execute failed: %v", err)
	}

	result := strings.TrimSpace(output.String())
	if !strings.Contains(result, "Name: demo") {
		t.Fatalf("expected text output to include workspace name, got: %q", result)
	}
	if !strings.Contains(result, "- ID: ") {
		t.Fatalf("expected text output to include workspace id, got: %q", result)
	}
}

func TestListWorkspaceTextOutputIncludesDefault(t *testing.T) {
	oldOutputFormat := OutputFormat
	oldRemoteAddr := RemoteAddr
	OutputFormat = "text"
	RemoteAddr = ""
	t.Cleanup(func() {
		OutputFormat = oldOutputFormat
		RemoteAddr = oldRemoteAddr
	})

	output, err := executeWorkspaceListCommand(t)
	if err != nil {
		t.Fatalf("command execute failed: %v", err)
	}

	if !strings.Contains(output, "Name: Default") {
		t.Fatalf("expected default workspace, got: %q", output)
	}
}

func TestListWorkspaceJsonOutput(t *testing.T) {
	manager := newIsolatedTestManager(t)
	_, err := manager.WorkspaceRepository.Save(&core.Workspace{Name: "demo"})
	if err != nil {
		t.Fatalf("save workspace failed: %v", err)
	}

	oldOutputFormat := OutputFormat
	oldRemoteAddr := RemoteAddr
	OutputFormat = "json"
	RemoteAddr = ""
	t.Cleanup(func() {
		OutputFormat = oldOutputFormat
		RemoteAddr = oldRemoteAddr
	})

	command := NewListWorkspaceCmd()
	var output bytes.Buffer
	command.SetOut(&output)
	command.SetErr(&output)

	err = command.Execute()
	if err != nil {
		t.Fatalf("command execute failed: %v", err)
	}

	var workspaces []core.Workspace
	if unmarshalErr := json.Unmarshal(output.Bytes(), &workspaces); unmarshalErr != nil {
		t.Fatalf("json unmarshal failed: %v, output: %q", unmarshalErr, output.String())
	}

	if len(workspaces) != 2 {
		t.Fatalf("expected default and demo workspaces, got: %d", len(workspaces))
	}
	foundDemo := false
	for _, workspace := range workspaces {
		if workspace.Name == "demo" {
			foundDemo = true
			break
		}
	}
	if !foundDemo {
		t.Fatalf("expected workspace name demo, got: %#v", workspaces)
	}
}

func TestListWorkspaceYamlAndShellOutput(t *testing.T) {
	manager := newIsolatedTestManager(t)
	_, err := manager.WorkspaceRepository.Save(&core.Workspace{Name: "demo"})
	if err != nil {
		t.Fatalf("save workspace failed: %v", err)
	}

	oldOutputFormat := OutputFormat
	oldRemoteAddr := RemoteAddr
	oldVerbose := Verbose
	RemoteAddr = ""
	t.Cleanup(func() {
		OutputFormat = oldOutputFormat
		RemoteAddr = oldRemoteAddr
		Verbose = oldVerbose
	})

	command := NewListWorkspaceCmd()
	var yamlOutput bytes.Buffer
	command.SetOut(&yamlOutput)
	command.SetErr(&yamlOutput)

	OutputFormat = "yaml"
	err = command.Execute()
	if err != nil {
		t.Fatalf("yaml command execute failed: %v", err)
	}
	if !strings.Contains(yamlOutput.String(), "name: demo") {
		t.Fatalf("expected yaml output to include workspace name, got: %q", yamlOutput.String())
	}

	command = NewListWorkspaceCmd()
	var shellOutput bytes.Buffer
	command.SetOut(&shellOutput)
	command.SetErr(&shellOutput)

	OutputFormat = "shell"
	err = command.Execute()
	if err != nil {
		t.Fatalf("shell command execute failed: %v", err)
	}
	if !strings.Contains(shellOutput.String(), "NAME=\"demo\"") {
		t.Fatalf("expected shell output to include flattened workspace name, got: %q", shellOutput.String())
	}
}

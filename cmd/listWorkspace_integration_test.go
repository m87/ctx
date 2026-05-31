package cmd

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
)

func newIsolatedTestManager(t *testing.T) *core.ContextManager {
	t.Helper()

	t.Setenv("DATABASE_PATH", filepath.Join(t.TempDir(), "ctx-test.db"))
	return bootstrap.CreateManager()
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

func TestListWorkspaceTextOutputEmpty(t *testing.T) {
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

	if output != "No workspaces found" {
		t.Fatalf("expected empty text message, got: %q", output)
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

	if len(workspaces) != 1 {
		t.Fatalf("expected 1 workspace in json output, got: %d", len(workspaces))
	}
	if workspaces[0].Name != "demo" {
		t.Fatalf("expected workspace name demo, got: %q", workspaces[0].Name)
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
	if !strings.Contains(shellOutput.String(), "RESULT_0_NAME=\"demo\"") {
		t.Fatalf("expected shell output to include flattened workspace name, got: %q", shellOutput.String())
	}
}

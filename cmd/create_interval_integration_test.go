package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/m87/ctx/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func executeCreateIntervalCommand(t *testing.T, commandArgs ...string) (string, error) {
	t.Helper()
	command := NewCreateIntervalCmd()
	var output bytes.Buffer
	command.SetOut(&output)
	command.SetErr(&output)
	command.SetArgs(commandArgs)

	err := command.Execute()
	return strings.TrimSpace(output.String()), err
}

func createIntervalTestContext(t *testing.T) string {
	t.Helper()
	manager := newIsolatedTestManager(t)
	workspaces, err := manager.WorkspaceRepository.List()
	require.NoError(t, err)
	require.NotEmpty(t, workspaces)
	contextID, err := manager.CreateContext(&core.Context{
		Name:        "Interval test context",
		WorkspaceId: workspaces[0].Id,
	})
	require.NoError(t, err)
	return contextID
}

func TestCreateIntervalCommandCreatesActiveIntervalWithoutEnd(t *testing.T) {
	contextID := createIntervalTestContext(t)
	oldOutputFormat := OutputFormat
	oldRemoteAddr := RemoteAddr
	OutputFormat = "json"
	RemoteAddr = ""
	t.Cleanup(func() {
		OutputFormat = oldOutputFormat
		RemoteAddr = oldRemoteAddr
	})

	output, err := executeCreateIntervalCommand(
		t,
		"--context-id", contextID,
		"--start", "2026-08-14T08:00:00Z",
	)

	require.NoError(t, err)
	var interval core.Interval
	require.NoError(t, json.Unmarshal([]byte(output), &interval))
	assert.Equal(t, core.IntervalStatusActive, interval.Status)
	assert.Nil(t, interval.End)
	assert.Zero(t, interval.Duration)
}

func TestCreateIntervalCommandRejectsZeroDuration(t *testing.T) {
	contextID := createIntervalTestContext(t)
	oldOutputFormat := OutputFormat
	oldRemoteAddr := RemoteAddr
	OutputFormat = "text"
	RemoteAddr = ""
	t.Cleanup(func() {
		OutputFormat = oldOutputFormat
		RemoteAddr = oldRemoteAddr
	})
	instant := time.Date(2026, 8, 14, 8, 0, 0, 0, time.UTC).Format(time.RFC3339)

	_, err := executeCreateIntervalCommand(
		t,
		"--context-id", contextID,
		"--start", instant,
		"--end", instant,
	)

	require.Error(t, err)
	assert.Equal(t, "cannot create interval: end time must be after start time", err.Error())
}

package core

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type WorkspaceRepositoryMock struct {
	WorkspaceRepository
	workspaces []*Workspace
	called     bool
}

func (m *WorkspaceRepositoryMock) List() ([]*Workspace, error) {
	m.called = true
	if m.workspaces == nil {
		return nil, errors.New("WorkspaceRepository.List error")
	}
	return m.workspaces, nil
}

type ContextRepositoryMock struct {
	ContextRepository
	contexts []*Context
	called   bool
}

func (m *ContextRepositoryMock) List() ([]*Context, error) {
	m.called = true
	if m.contexts == nil {
		return nil, errors.New("ContextRepository.List error")
	}
	return m.contexts, nil
}

type IntervalRepositoryMock struct {
	IntervalRepository
	intervals []*Interval
	called    bool
}

func (m *IntervalRepositoryMock) List() ([]*Interval, error) {
	m.called = true
	if m.intervals == nil {
		return nil, errors.New("IntervalRepository.List error")
	}
	return m.intervals, nil
}

func setupManagerCorrectData() *ContextManager {
	workspaceRepo := &WorkspaceRepositoryMock{
		workspaces: []*Workspace{
			{Id: "workspace1"},
			{Id: "workspace2"},
		},
	}
	contextRepo := &ContextRepositoryMock{
		contexts: []*Context{
			{Id: "context1", Name: "Context 1", WorkspaceId: "workspace1"},
			{Id: "context2", Name: "Context 2", WorkspaceId: "workspace2"},
		},
	}
	intervalRepo := &IntervalRepositoryMock{
		intervals: []*Interval{
			{Id: "interval1", ContextId: "context1", WorkspaceId: "workspace1"},
			{Id: "interval2", ContextId: "context2", WorkspaceId: "workspace2"},
		},
	}

	manager := &ContextManager{
		WorkspaceRepository: workspaceRepo,
		ContextRepository:   contextRepo,
		IntervalRepository:  intervalRepo,
	}

	return manager
}

func TestPassIntegrityCheckTests(t *testing.T) {
	manager := setupManagerCorrectData()

	report, err := manager.CheckIntegrity()
	require.NoError(t, err)
	require.True(t, report.Healthy)
	require.Empty(t, report.Issues)
	require.Equal(t, 2, report.WorkspaceCount)
	require.Equal(t, 2, report.ContextCount)
	require.Equal(t, 2, report.IntervalCount)
}

func TestPassIntegrityCheckWithEmptyRepositories(t *testing.T) {
	manager := &ContextManager{
		WorkspaceRepository: &WorkspaceRepositoryMock{workspaces: []*Workspace{}},
		ContextRepository:   &ContextRepositoryMock{contexts: []*Context{}},
		IntervalRepository:  &IntervalRepositoryMock{intervals: []*Interval{}},
	}

	report, err := manager.CheckIntegrity()
	require.NoError(t, err)
	require.True(t, report.Healthy)
	require.Empty(t, report.Issues)
	require.Equal(t, 0, report.WorkspaceCount)
	require.Equal(t, 0, report.ContextCount)
	require.Equal(t, 0, report.IntervalCount)
}

func TestFailIntegrityCheckWithContextWithoutWorkspace(t *testing.T) {
	manager := setupManagerCorrectData()
	manager.ContextRepository.(*ContextRepositoryMock).contexts[0].WorkspaceId = ""
	manager.IntervalRepository.(*IntervalRepositoryMock).intervals = []*Interval{}

	report, err := manager.CheckIntegrity()
	require.NoError(t, err)
	require.False(t, report.Healthy)
	require.Len(t, report.Issues, 1)
	issue := report.Issues[0]
	require.Equal(t, "context", issue.EntityType)
	require.Equal(t, "context1", issue.EntityId)
	require.Equal(t, "CONTEXT_MISSING_WORKSPACE", issue.Code)
	require.True(t, issue.Repairable)
	require.Equal(t, "Context 1", issue.Details.Name)
}

func TestFailIntegrityCheckWithContextWithNonexistentWorkspace(t *testing.T) {
	manager := setupManagerCorrectData()
	manager.ContextRepository.(*ContextRepositoryMock).contexts[0].WorkspaceId = "nonexistent"
	manager.IntervalRepository.(*IntervalRepositoryMock).intervals = []*Interval{}

	report, err := manager.CheckIntegrity()
	require.NoError(t, err)
	require.False(t, report.Healthy)
	require.Len(t, report.Issues, 1)
	issue := report.Issues[0]
	require.Equal(t, "context", issue.EntityType)
	require.Equal(t, "context1", issue.EntityId)
	require.Equal(t, "CONTEXT_WORKSPACE_NOT_FOUND", issue.Code)
}

func TestFailIntegrityCheckWithIntervalWithoutContext(t *testing.T) {
	manager := setupManagerCorrectData()
	manager.IntervalRepository.(*IntervalRepositoryMock).intervals[0].ContextId = ""

	report, err := manager.CheckIntegrity()
	require.NoError(t, err)
	require.False(t, report.Healthy)
	require.Len(t, report.Issues, 1)
	issue := report.Issues[0]
	require.Equal(t, "interval", issue.EntityType)
	require.Equal(t, "interval1", issue.EntityId)
	require.Equal(t, "INTERVAL_MISSING_CONTEXT", issue.Code)
}

func TestFailIntegrityCheckWithIntervalWithNonexistentContext(t *testing.T) {
	manager := setupManagerCorrectData()
	manager.IntervalRepository.(*IntervalRepositoryMock).intervals[0].ContextId = "nonexistent"

	report, err := manager.CheckIntegrity()
	require.NoError(t, err)
	require.False(t, report.Healthy)
	require.Len(t, report.Issues, 1)
	issue := report.Issues[0]
	require.Equal(t, "interval", issue.EntityType)
	require.Equal(t, "interval1", issue.EntityId)
	require.Equal(t, "INTERVAL_CONTEXT_NOT_FOUND", issue.Code)
	require.False(t, issue.Repairable)
	require.Equal(t, "nonexistent", issue.Details.ContextId)
	require.Equal(t, "workspace1", issue.Details.WorkspaceId)
	require.NotNil(t, issue.Details.Start)
	require.NotNil(t, issue.Details.End)
}

func TestFailIntegrityCheckWithIntervalWithoutWorkspace(t *testing.T) {
	manager := setupManagerCorrectData()
	manager.IntervalRepository.(*IntervalRepositoryMock).intervals[0].WorkspaceId = ""

	report, err := manager.CheckIntegrity()
	require.NoError(t, err)
	require.False(t, report.Healthy)
	require.Len(t, report.Issues, 1)
	issue := report.Issues[0]
	require.Equal(t, "interval", issue.EntityType)
	require.Equal(t, "interval1", issue.EntityId)
	require.Equal(t, "INTERVAL_MISSING_WORKSPACE", issue.Code)
	require.True(t, issue.Repairable)
}

func TestFailIntegrityCheckWithIntervalWithNonexistentWorkspace(t *testing.T) {
	manager := setupManagerCorrectData()
	manager.IntervalRepository.(*IntervalRepositoryMock).intervals[0].WorkspaceId = "nonexistent"

	report, err := manager.CheckIntegrity()
	require.NoError(t, err)
	require.False(t, report.Healthy)
	require.Len(t, report.Issues, 1)
	issue := report.Issues[0]
	require.Equal(t, "interval", issue.EntityType)
	require.Equal(t, "interval1", issue.EntityId)
	require.Equal(t, "INTERVAL_WORKSPACE_NOT_FOUND", issue.Code)
}

func TestFailIntegrityCheckWithIntervalWorkspaceMismatch(t *testing.T) {
	manager := setupManagerCorrectData()
	manager.IntervalRepository.(*IntervalRepositoryMock).intervals[0].WorkspaceId = "workspace2"

	report, err := manager.CheckIntegrity()
	require.NoError(t, err)
	require.False(t, report.Healthy)
	require.Len(t, report.Issues, 1)
	issue := report.Issues[0]
	require.Equal(t, "interval", issue.EntityType)
	require.Equal(t, "interval1", issue.EntityId)
	require.Equal(t, "INTERVAL_WORKSPACE_MISMATCH", issue.Code)
	require.True(t, issue.Repairable)
}

func TestFailIntegrityCheckWithMultipleIssues(t *testing.T) {
	manager := setupManagerCorrectData()
	manager.ContextRepository.(*ContextRepositoryMock).contexts[0].WorkspaceId = ""
	manager.IntervalRepository.(*IntervalRepositoryMock).intervals[0].ContextId = "nonexistent"
	manager.IntervalRepository.(*IntervalRepositoryMock).intervals[0].WorkspaceId = "nonexistent"

	report, err := manager.CheckIntegrity()
	require.NoError(t, err)
	require.False(t, report.Healthy)
	require.Len(t, report.Issues, 3)
	require.Equal(t, "CONTEXT_MISSING_WORKSPACE", report.Issues[0].Code)
	require.Equal(t, "INTERVAL_CONTEXT_NOT_FOUND", report.Issues[1].Code)
	require.Equal(t, "INTERVAL_WORKSPACE_NOT_FOUND", report.Issues[2].Code)
}

func TestFailIntegrityCheckWithAllIssues(t *testing.T) {
	manager := setupManagerCorrectData()
	manager.ContextRepository.(*ContextRepositoryMock).contexts[0].WorkspaceId = ""
	manager.IntervalRepository.(*IntervalRepositoryMock).intervals[0].ContextId = ""
	manager.IntervalRepository.(*IntervalRepositoryMock).intervals[0].WorkspaceId = "nonexistent"

	report, err := manager.CheckIntegrity()
	require.NoError(t, err)
	require.False(t, report.Healthy)
	require.Len(t, report.Issues, 3)
	require.Equal(t, "CONTEXT_MISSING_WORKSPACE", report.Issues[0].Code)
	require.Equal(t, "INTERVAL_MISSING_CONTEXT", report.Issues[1].Code)
	require.Equal(t, "INTERVAL_WORKSPACE_NOT_FOUND", report.Issues[2].Code)
}

func TestIntegrityCheckOnRepositoryFail(t *testing.T) {
	manager := &ContextManager{
		WorkspaceRepository: &WorkspaceRepositoryMock{workspaces: nil},
		ContextRepository:   &ContextRepositoryMock{contexts: nil},
		IntervalRepository:  &IntervalRepositoryMock{intervals: nil},
	}

	report, err := manager.CheckIntegrity()
	require.Nil(t, report)
	require.Error(t, err)
	require.Equal(t, "WorkspaceRepository.List error", err.Error())
	require.False(t, manager.ContextRepository.(*ContextRepositoryMock).called)
	require.False(t, manager.IntervalRepository.(*IntervalRepositoryMock).called)

	manager = &ContextManager{
		WorkspaceRepository: &WorkspaceRepositoryMock{workspaces: []*Workspace{}},
		ContextRepository:   &ContextRepositoryMock{contexts: nil},
		IntervalRepository:  &IntervalRepositoryMock{intervals: nil},
	}

	report, err = manager.CheckIntegrity()
	require.Nil(t, report)
	require.Error(t, err)
	require.Equal(t, "ContextRepository.List error", err.Error())
	require.True(t, manager.ContextRepository.(*ContextRepositoryMock).called)
	require.False(t, manager.IntervalRepository.(*IntervalRepositoryMock).called)

	manager = &ContextManager{
		WorkspaceRepository: &WorkspaceRepositoryMock{workspaces: []*Workspace{}},
		ContextRepository:   &ContextRepositoryMock{contexts: []*Context{}},
		IntervalRepository:  &IntervalRepositoryMock{intervals: nil},
	}

	report, err = manager.CheckIntegrity()
	require.Nil(t, report)
	require.Error(t, err)
	require.Equal(t, "IntervalRepository.List error", err.Error())
	require.True(t, manager.ContextRepository.(*ContextRepositoryMock).called)
	require.True(t, manager.IntervalRepository.(*IntervalRepositoryMock).called)
}

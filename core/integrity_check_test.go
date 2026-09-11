package core

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var integrityTestBaseTime = time.Date(2026, 1, 2, 10, 0, 0, 0, time.UTC)

func integrityTestTime(offset time.Duration) *time.Time {
	return testTime(integrityTestBaseTime.Add(offset))
}

func newIntegrityCheckTestManager() *TestContextManager {
	test := NewEmptyTestContextManager()
	test.Workspaces.Seed(
		&Workspace{Id: "workspace1"},
		&Workspace{Id: "workspace2"},
	)
	test.Contexts.Seed(
		&Context{Id: "context1", Name: "Context 1", WorkspaceId: "workspace1"},
		&Context{Id: "context2", Name: "Context 2", WorkspaceId: "workspace2"},
	)
	test.Intervals.Seed(
		&Interval{Id: "interval1", ContextId: "context1", WorkspaceId: "workspace1", Status: "completed", Start: integrityTestTime(0), End: integrityTestTime(time.Hour)},
		&Interval{Id: "interval2", ContextId: "context2", WorkspaceId: "workspace2", Status: "completed", Start: integrityTestTime(2 * time.Hour), End: integrityTestTime(3 * time.Hour)},
	)
	return test
}

func TestPassIntegrityCheckTests(t *testing.T) {
	test := newIntegrityCheckTestManager()

	report, err := test.Manager.CheckIntegrity()
	require.NoError(t, err)
	require.True(t, report.Healthy)
	require.Empty(t, report.Issues)
	require.Equal(t, 2, report.WorkspaceCount)
	require.Equal(t, 2, report.ContextCount)
	require.Equal(t, 2, report.IntervalCount)
}

func TestPassIntegrityCheckWithEmptyRepositories(t *testing.T) {
	test := NewEmptyTestContextManager()

	report, err := test.Manager.CheckIntegrity()
	require.NoError(t, err)
	require.True(t, report.Healthy)
	require.Empty(t, report.Issues)
	require.Equal(t, 0, report.WorkspaceCount)
	require.Equal(t, 0, report.ContextCount)
	require.Equal(t, 0, report.IntervalCount)
}

func TestFailIntegrityCheckWithContextWithoutWorkspace(t *testing.T) {
	test := newIntegrityCheckTestManager()
	test.Contexts.Get("context1").WorkspaceId = ""
	test.Intervals.Seed()

	report, err := test.Manager.CheckIntegrity()
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
	test := newIntegrityCheckTestManager()
	test.Contexts.Get("context1").WorkspaceId = "nonexistent"
	test.Intervals.Seed()

	report, err := test.Manager.CheckIntegrity()
	require.NoError(t, err)
	require.False(t, report.Healthy)
	require.Len(t, report.Issues, 1)
	issue := report.Issues[0]
	require.Equal(t, "context", issue.EntityType)
	require.Equal(t, "context1", issue.EntityId)
	require.Equal(t, "CONTEXT_WORKSPACE_NOT_FOUND", issue.Code)
}

func TestFailIntegrityCheckWithIntervalWithoutContext(t *testing.T) {
	test := newIntegrityCheckTestManager()
	test.Intervals.Get("interval1").ContextId = ""

	report, err := test.Manager.CheckIntegrity()
	require.NoError(t, err)
	require.False(t, report.Healthy)
	require.Len(t, report.Issues, 1)
	issue := report.Issues[0]
	require.Equal(t, "interval", issue.EntityType)
	require.Equal(t, "interval1", issue.EntityId)
	require.Equal(t, "INTERVAL_MISSING_CONTEXT", issue.Code)
}

func TestFailIntegrityCheckWithIntervalWithNonexistentContext(t *testing.T) {
	test := newIntegrityCheckTestManager()
	test.Intervals.Get("interval1").ContextId = "nonexistent"

	report, err := test.Manager.CheckIntegrity()
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
	test := newIntegrityCheckTestManager()
	test.Intervals.Get("interval1").WorkspaceId = ""

	report, err := test.Manager.CheckIntegrity()
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
	test := newIntegrityCheckTestManager()
	test.Intervals.Get("interval1").WorkspaceId = "nonexistent"

	report, err := test.Manager.CheckIntegrity()
	require.NoError(t, err)
	require.False(t, report.Healthy)
	require.Len(t, report.Issues, 1)
	issue := report.Issues[0]
	require.Equal(t, "interval", issue.EntityType)
	require.Equal(t, "interval1", issue.EntityId)
	require.Equal(t, "INTERVAL_WORKSPACE_NOT_FOUND", issue.Code)
}

func TestFailIntegrityCheckWithIntervalWorkspaceMismatch(t *testing.T) {
	test := newIntegrityCheckTestManager()
	test.Intervals.Get("interval1").WorkspaceId = "workspace2"

	report, err := test.Manager.CheckIntegrity()
	require.NoError(t, err)
	require.False(t, report.Healthy)
	require.Len(t, report.Issues, 1)
	issue := report.Issues[0]
	require.Equal(t, "interval", issue.EntityType)
	require.Equal(t, "interval1", issue.EntityId)
	require.Equal(t, "INTERVAL_WORKSPACE_MISMATCH", issue.Code)
	require.True(t, issue.Repairable)
}

func TestFailIntegrityCheckWithInactiveIntervalMissingTime(t *testing.T) {
	test := newIntegrityCheckTestManager()
	interval := test.Intervals.Get("interval1")
	interval.Status = "completed"
	interval.Start = nil

	report, err := test.Manager.CheckIntegrity()
	require.NoError(t, err)
	require.False(t, report.Healthy)
	require.Len(t, report.Issues, 1)
	issue := report.Issues[0]
	require.Equal(t, "interval", issue.EntityType)
	require.Equal(t, "interval1", issue.EntityId)
	require.Equal(t, "INACTIVE_INTERVAL_MISSING_TIME", issue.Code)
	require.False(t, issue.Repairable)
}

func TestFailIntegrityCheckWithActiveIntervalWithEnd(t *testing.T) {
	test := newIntegrityCheckTestManager()
	interval := test.Intervals.Get("interval1")
	interval.Status = "active"
	interval.Start = integrityTestTime(0)
	interval.End = integrityTestTime(time.Hour)

	report, err := test.Manager.CheckIntegrity()
	require.NoError(t, err)
	require.False(t, report.Healthy)
	require.Len(t, report.Issues, 1)
	issue := report.Issues[0]
	require.Equal(t, "interval", issue.EntityType)
	require.Equal(t, "interval1", issue.EntityId)
	require.Equal(t, "ACTIVE_INTERVAL_HAS_END", issue.Code)
	require.True(t, issue.Repairable)
}

func TestFailIntegrityCheckWithMultipleActiveContexts(t *testing.T) {
	test := newIntegrityCheckTestManager()
	firstContext := test.Contexts.Get("context1")
	secondContext := test.Contexts.Get("context2")
	firstInterval := test.Intervals.Get("interval1")
	secondInterval := test.Intervals.Get("interval2")
	firstContext.Status = "active"
	secondContext.Status = "active"
	firstInterval.Status = "active"
	firstInterval.Start = integrityTestTime(0)
	firstInterval.End = nil
	secondInterval.Status = "active"
	secondInterval.Start = integrityTestTime(time.Hour)
	secondInterval.End = nil

	report, err := test.Manager.CheckIntegrity()
	require.NoError(t, err)
	require.False(t, report.Healthy)
	require.Len(t, report.Issues, 2)
	require.Equal(t, "MULTIPLE_ACTIVE_CONTEXTS", report.Issues[0].Code)
	require.Equal(t, "MULTIPLE_ACTIVE_CONTEXTS", report.Issues[1].Code)
	require.True(t, report.Issues[0].Repairable)
	require.True(t, report.Issues[1].Repairable)
}

func TestFailIntegrityCheckWithActiveContextWithoutOpenInterval(t *testing.T) {
	test := newIntegrityCheckTestManager()
	context := test.Contexts.Get("context1")
	interval := test.Intervals.Get("interval1")
	context.Status = "active"
	interval.Status = "completed"
	interval.Start = integrityTestTime(0)
	interval.End = integrityTestTime(time.Hour)

	report, err := test.Manager.CheckIntegrity()
	require.NoError(t, err)
	require.False(t, report.Healthy)
	require.Len(t, report.Issues, 1)
	issue := report.Issues[0]
	require.Equal(t, "context", issue.EntityType)
	require.Equal(t, "context1", issue.EntityId)
	require.Equal(t, "ACTIVE_CONTEXT_WITHOUT_OPEN_INTERVAL", issue.Code)
	require.True(t, issue.Repairable)
}

func TestFailIntegrityCheckWithMultipleIssues(t *testing.T) {
	test := newIntegrityCheckTestManager()
	test.Contexts.Get("context1").WorkspaceId = ""
	test.Intervals.Get("interval1").ContextId = "nonexistent"
	test.Intervals.Get("interval1").WorkspaceId = "nonexistent"

	report, err := test.Manager.CheckIntegrity()
	require.NoError(t, err)
	require.False(t, report.Healthy)
	require.Len(t, report.Issues, 3)
	require.Equal(t, "CONTEXT_MISSING_WORKSPACE", report.Issues[0].Code)
	require.Equal(t, "INTERVAL_CONTEXT_NOT_FOUND", report.Issues[1].Code)
	require.Equal(t, "INTERVAL_WORKSPACE_NOT_FOUND", report.Issues[2].Code)
}

func TestFailIntegrityCheckWithAllIssues(t *testing.T) {
	test := newIntegrityCheckTestManager()
	test.Contexts.Get("context1").WorkspaceId = ""
	test.Intervals.Get("interval1").ContextId = ""
	test.Intervals.Get("interval1").WorkspaceId = "nonexistent"

	report, err := test.Manager.CheckIntegrity()
	require.NoError(t, err)
	require.False(t, report.Healthy)
	require.Len(t, report.Issues, 3)
	require.Equal(t, "CONTEXT_MISSING_WORKSPACE", report.Issues[0].Code)
	require.Equal(t, "INTERVAL_MISSING_CONTEXT", report.Issues[1].Code)
	require.Equal(t, "INTERVAL_WORKSPACE_NOT_FOUND", report.Issues[2].Code)
}

func TestIntegrityCheckOnRepositoryFail(t *testing.T) {
	test := NewEmptyTestContextManager()
	test.Workspaces.listError = errors.New("WorkspaceRepository.List error")

	report, err := test.Manager.CheckIntegrity()
	require.Nil(t, report)
	require.EqualError(t, err, "WorkspaceRepository.List error")
	require.False(t, test.Contexts.listCalled)
	require.False(t, test.Intervals.listCalled)

	test = NewEmptyTestContextManager()
	test.Contexts.listError = errors.New("ContextRepository.List error")

	report, err = test.Manager.CheckIntegrity()
	require.Nil(t, report)
	require.EqualError(t, err, "ContextRepository.List error")
	require.True(t, test.Contexts.listCalled)
	require.False(t, test.Intervals.listCalled)

	test = NewEmptyTestContextManager()
	test.Intervals.listError = errors.New("IntervalRepository.List error")

	report, err = test.Manager.CheckIntegrity()
	require.Nil(t, report)
	require.EqualError(t, err, "IntervalRepository.List error")
	require.True(t, test.Contexts.listCalled)
	require.True(t, test.Intervals.listCalled)
}

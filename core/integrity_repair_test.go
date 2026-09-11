package core

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func newIntegrityRepairTestManager() *TestContextManager {
	test := NewEmptyTestContextManager()
	test.Workspaces.Seed(
		&Workspace{Id: "workspace1", Name: "Default"},
		&Workspace{Id: "workspace2"},
	)
	test.Contexts.Seed(
		&Context{Id: "context1", WorkspaceId: "workspace1"},
		&Context{Id: "context2", WorkspaceId: "workspace2"},
	)
	test.Intervals.Seed(
		&Interval{Id: "interval1", ContextId: "context1", WorkspaceId: "workspace1", Status: "completed", Start: integrityTestTime(0), End: integrityTestTime(time.Hour)},
		&Interval{Id: "interval2", ContextId: "context2", WorkspaceId: "workspace2", Status: "completed", Start: integrityTestTime(2 * time.Hour), End: integrityTestTime(3 * time.Hour)},
	)
	test.Projects.Seed(
		&Project{Id: "project1", WorkspaceId: "workspace1"},
		&Project{Id: "project2", WorkspaceId: "workspace2"},
	)
	return test
}

func TestPassIntegrityRepairWithCorrectData(t *testing.T) {
	test := newIntegrityRepairTestManager()

	result, err := test.Manager.RepairIntegrity()
	require.NoError(t, err)
	require.Equal(t, 0, result.RepairedCount)
	require.True(t, result.Report.Healthy)
	require.Empty(t, result.Report.Issues)
	require.Empty(t, test.Workspaces.saved)
	require.Empty(t, test.Contexts.saved)
	require.Empty(t, test.Intervals.saved)
}

func TestIntegrityRepairCreatesDefaultWorkspace(t *testing.T) {
	test := newIntegrityRepairTestManager()
	workspaceRepo := test.Workspaces
	workspaceRepo.Seed()
	test.Contexts.Seed()
	test.Intervals.Seed()

	result, err := test.Manager.RepairIntegrity()
	require.NoError(t, err)
	require.Equal(t, 1, result.RepairedCount)
	require.True(t, result.Report.Healthy)
	require.Len(t, workspaceRepo.saved, 1)
	require.NotEmpty(t, workspaceRepo.saved[0].Id)
	require.Equal(t, "Default", workspaceRepo.saved[0].Name)
}

func TestIntegrityRepairCreatesDefaultWorkspaceForUnassignedContext(t *testing.T) {
	test := newIntegrityRepairTestManager()
	workspaceRepo := test.Workspaces
	workspaceRepo.Get("workspace1").Name = "Workspace 1"
	contextRepo := test.Contexts
	intervalRepo := test.Intervals
	contextRepo.Get("context1").WorkspaceId = ""
	intervalRepo.Get("interval1").WorkspaceId = ""

	result, err := test.Manager.RepairIntegrity()
	require.NoError(t, err)
	require.Equal(t, 3, result.RepairedCount)
	require.True(t, result.Report.Healthy)
	require.Len(t, workspaceRepo.saved, 1)
	require.Equal(t, workspaceRepo.saved[0].Id, contextRepo.Get("context1").WorkspaceId)
	require.Equal(t, workspaceRepo.saved[0].Id, intervalRepo.Get("interval1").WorkspaceId)
}

func TestIntegrityRepairDoesNotCreateDefaultWorkspaceWhenDataIsValid(t *testing.T) {
	test := newIntegrityRepairTestManager()
	workspaceRepo := test.Workspaces
	workspaceRepo.Get("workspace1").Name = "Workspace 1"

	result, err := test.Manager.RepairIntegrity()
	require.NoError(t, err)
	require.Equal(t, 0, result.RepairedCount)
	require.True(t, result.Report.Healthy)
	require.Empty(t, workspaceRepo.saved)
}

func TestIntegrityRepairContextWithoutWorkspace(t *testing.T) {
	test := newIntegrityRepairTestManager()
	contextRepo := test.Contexts
	contextRepo.Get("context1").WorkspaceId = ""
	test.Intervals.Seed()

	result, err := test.Manager.RepairIntegrity()
	require.NoError(t, err)
	require.Equal(t, 1, result.RepairedCount)
	require.True(t, result.Report.Healthy)
	require.Len(t, contextRepo.saved, 1)
	require.Equal(t, "context1", contextRepo.saved[0].Id)
	require.Equal(t, "workspace1", contextRepo.saved[0].WorkspaceId)
}

func TestIntegrityRepairContextAndItsIntervalWithoutWorkspace(t *testing.T) {
	test := newIntegrityRepairTestManager()
	contextRepo := test.Contexts
	intervalRepo := test.Intervals
	contextRepo.Get("context1").WorkspaceId = ""
	intervalRepo.Get("interval1").WorkspaceId = ""

	report, err := test.Manager.CheckIntegrity()
	require.NoError(t, err)
	require.Len(t, report.Issues, 2)
	require.True(t, report.Issues[0].Repairable)
	require.True(t, report.Issues[1].Repairable)

	result, err := test.Manager.RepairIntegrity()
	require.NoError(t, err)
	require.Equal(t, 2, result.RepairedCount)
	require.True(t, result.Report.Healthy)
	require.Equal(t, "workspace1", contextRepo.Get("context1").WorkspaceId)
	require.Equal(t, "workspace1", intervalRepo.Get("interval1").WorkspaceId)
}

func TestIntegrityRepairContextWithNonexistentWorkspace(t *testing.T) {
	test := newIntegrityRepairTestManager()
	contextRepo := test.Contexts
	contextRepo.Get("context1").WorkspaceId = "nonexistent"
	test.Intervals.Seed()

	result, err := test.Manager.RepairIntegrity()
	require.NoError(t, err)
	require.Equal(t, 1, result.RepairedCount)
	require.True(t, result.Report.Healthy)
	require.Len(t, contextRepo.saved, 1)
	require.Equal(t, "context1", contextRepo.saved[0].Id)
	require.Equal(t, "workspace1", contextRepo.saved[0].WorkspaceId)
}

func TestIntegrityRepairIntervalWithoutWorkspace(t *testing.T) {
	test := newIntegrityRepairTestManager()
	intervalRepo := test.Intervals
	intervalRepo.Get("interval1").WorkspaceId = ""

	result, err := test.Manager.RepairIntegrity()
	require.NoError(t, err)
	require.Equal(t, 1, result.RepairedCount)
	require.True(t, result.Report.Healthy)
	require.Len(t, intervalRepo.saved, 1)
	require.Equal(t, "interval1", intervalRepo.saved[0].Id)
	require.Equal(t, "workspace1", intervalRepo.saved[0].WorkspaceId)
}

func TestIntegrityRepairIntervalWithNonexistentWorkspace(t *testing.T) {
	test := newIntegrityRepairTestManager()
	intervalRepo := test.Intervals
	intervalRepo.Get("interval1").WorkspaceId = "nonexistent"

	result, err := test.Manager.RepairIntegrity()
	require.NoError(t, err)
	require.Equal(t, 1, result.RepairedCount)
	require.True(t, result.Report.Healthy)
	require.Len(t, intervalRepo.saved, 1)
	require.Equal(t, "interval1", intervalRepo.saved[0].Id)
	require.Equal(t, "workspace1", intervalRepo.saved[0].WorkspaceId)
}

func TestIntegrityRepairIntervalWorkspaceMismatch(t *testing.T) {
	test := newIntegrityRepairTestManager()
	intervalRepo := test.Intervals
	intervalRepo.Get("interval1").WorkspaceId = "workspace2"

	result, err := test.Manager.RepairIntegrity()
	require.NoError(t, err)
	require.Equal(t, 1, result.RepairedCount)
	require.True(t, result.Report.Healthy)
	require.Len(t, intervalRepo.saved, 1)
	require.Equal(t, "interval1", intervalRepo.saved[0].Id)
	require.Equal(t, "workspace1", intervalRepo.saved[0].WorkspaceId)
}

func TestIntegrityRepairCompletesActiveIntervalWithEndAndStopsContext(t *testing.T) {
	test := newIntegrityRepairTestManager()
	contextRepo := test.Contexts
	intervalRepo := test.Intervals
	contextRepo.Get("context1").Status = "active"
	intervalRepo.Get("interval1").Status = "active"
	intervalRepo.Get("interval1").Start = integrityTestTime(0)
	intervalRepo.Get("interval1").End = integrityTestTime(time.Hour)

	result, err := test.Manager.RepairIntegrity()
	require.NoError(t, err)
	require.Equal(t, 2, result.RepairedCount)
	require.True(t, result.Report.Healthy)
	require.Len(t, intervalRepo.saved, 1)
	require.Equal(t, "completed", intervalRepo.Get("interval1").Status)
	require.Equal(t, time.Hour, intervalRepo.Get("interval1").Duration)
	require.Len(t, contextRepo.saved, 1)
	require.Equal(t, "inactive", contextRepo.Get("context1").Status)
}

func TestIntegrityRepairStopsActiveContextWithCompletedIntervalsOnly(t *testing.T) {
	test := newIntegrityRepairTestManager()
	contextRepo := test.Contexts
	intervalRepo := test.Intervals
	contextRepo.Get("context1").Status = "active"
	intervalRepo.Get("interval1").Status = "completed"
	intervalRepo.Get("interval1").Start = integrityTestTime(0)
	intervalRepo.Get("interval1").End = integrityTestTime(time.Hour)

	result, err := test.Manager.RepairIntegrity()
	require.NoError(t, err)
	require.Equal(t, 1, result.RepairedCount)
	require.True(t, result.Report.Healthy)
	require.Len(t, contextRepo.saved, 1)
	require.Equal(t, "inactive", contextRepo.Get("context1").Status)
	require.Empty(t, intervalRepo.saved)
}

func TestIntegrityRepairLeavesOnlyNewestActiveContextRunning(t *testing.T) {
	test := newIntegrityRepairTestManager()
	repairTime := integrityTestBaseTime.Add(3 * time.Hour)
	test.TimeProvider.Set(repairTime)
	contextRepo := test.Contexts
	intervalRepo := test.Intervals
	contextRepo.Get("context1").Status = "active"
	contextRepo.Get("context2").Status = "active"
	intervalRepo.Get("interval1").Status = "active"
	intervalRepo.Get("interval1").Start = integrityTestTime(time.Hour)
	intervalRepo.Get("interval1").End = nil
	intervalRepo.Get("interval2").Status = "active"
	intervalRepo.Get("interval2").Start = integrityTestTime(2 * time.Hour)
	intervalRepo.Get("interval2").End = nil

	result, err := test.Manager.RepairIntegrity()
	require.NoError(t, err)
	require.Equal(t, 2, result.RepairedCount)
	require.True(t, result.Report.Healthy)
	require.Equal(t, "inactive", contextRepo.Get("context1").Status)
	require.Equal(t, "active", contextRepo.Get("context2").Status)
	require.Len(t, contextRepo.saved, 1)
	require.Equal(t, "context1", contextRepo.saved[0].Id)
	require.Len(t, intervalRepo.saved, 1)
	require.Equal(t, "interval1", intervalRepo.saved[0].Id)
	require.Equal(t, "completed", intervalRepo.Get("interval1").Status)
	require.Equal(t, repairTime, *intervalRepo.Get("interval1").End)
	require.Equal(t, 2*time.Hour, intervalRepo.Get("interval1").Duration)
	require.Equal(t, "active", intervalRepo.Get("interval2").Status)
	require.False(t, timeIsSet(intervalRepo.Get("interval2").End))
}

func TestIntegrityRepairLeavesIntervalWithNonexistentContext(t *testing.T) {
	test := newIntegrityRepairTestManager()
	intervalRepo := test.Intervals
	intervalRepo.Get("interval1").ContextId = "nonexistent"

	result, err := test.Manager.RepairIntegrity()
	require.NoError(t, err)
	require.Equal(t, 0, result.RepairedCount)
	require.False(t, result.Report.Healthy)
	require.Len(t, result.Report.Issues, 1)
	require.Equal(t, "INTERVAL_CONTEXT_NOT_FOUND", result.Report.Issues[0].Code)
	require.Empty(t, intervalRepo.saved)
}

func TestIntegrityRepairWithMultipleIssues(t *testing.T) {
	test := newIntegrityRepairTestManager()
	contextRepo := test.Contexts
	intervalRepo := test.Intervals
	contextRepo.Get("context1").WorkspaceId = ""
	intervalRepo.Get("interval1").WorkspaceId = "nonexistent"

	result, err := test.Manager.RepairIntegrity()
	require.NoError(t, err)
	require.Equal(t, 2, result.RepairedCount)
	require.True(t, result.Report.Healthy)
	require.Len(t, contextRepo.saved, 1)
	require.Len(t, intervalRepo.saved, 1)
	require.Equal(t, "workspace1", contextRepo.saved[0].WorkspaceId)
	require.Equal(t, "workspace1", intervalRepo.saved[0].WorkspaceId)
}

func TestIntegrityRepairOnRepositoryFail(t *testing.T) {
	listError := errors.New("WorkspaceRepository.List error")
	test := newIntegrityRepairTestManager()
	test.Workspaces.listError = listError

	result, err := test.Manager.RepairIntegrity()
	require.Nil(t, result)
	require.ErrorIs(t, err, listError)
	require.False(t, test.Contexts.listCalled)
	require.False(t, test.Intervals.listCalled)

	listError = errors.New("ContextRepository.List error")
	test = newIntegrityRepairTestManager()
	test.Contexts.listError = listError

	result, err = test.Manager.RepairIntegrity()
	require.Nil(t, result)
	require.ErrorIs(t, err, listError)
	require.False(t, test.Intervals.listCalled)

	listError = errors.New("IntervalRepository.List error")
	test = newIntegrityRepairTestManager()
	test.Intervals.listError = listError

	result, err = test.Manager.RepairIntegrity()
	require.Nil(t, result)
	require.ErrorIs(t, err, listError)
}

func TestIntegrityRepairOnSaveFail(t *testing.T) {
	saveError := errors.New("WorkspaceRepository.Save error")
	test := newIntegrityRepairTestManager()
	test.Workspaces.Seed()
	test.Workspaces.saveError = saveError

	result, err := test.Manager.RepairIntegrity()
	require.Nil(t, result)
	require.ErrorIs(t, err, saveError)
	require.True(t, test.Contexts.listCalled)
	require.False(t, test.Intervals.listCalled)

	saveError = errors.New("ContextRepository.Save error")
	test = newIntegrityRepairTestManager()
	test.Contexts.Get("context1").WorkspaceId = ""
	test.Contexts.saveError = saveError

	result, err = test.Manager.RepairIntegrity()
	require.Nil(t, result)
	require.ErrorIs(t, err, saveError)
	require.False(t, test.Intervals.listCalled)

	saveError = errors.New("IntervalRepository.Save error")
	test = newIntegrityRepairTestManager()
	test.Intervals.Get("interval1").WorkspaceId = ""
	test.Intervals.saveError = saveError

	result, err = test.Manager.RepairIntegrity()
	require.Nil(t, result)
	require.ErrorIs(t, err, saveError)
}

package core

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type WorkspaceRepositoryRepairMock struct {
	WorkspaceRepository
	workspaces []*Workspace
	saved      []*Workspace
	listError  error
	saveError  error
	called     bool
}

func (m *WorkspaceRepositoryRepairMock) List() ([]*Workspace, error) {
	m.called = true
	return m.workspaces, m.listError
}

func (m *WorkspaceRepositoryRepairMock) Save(workspace *Workspace) (string, error) {
	if m.saveError != nil {
		return "", m.saveError
	}
	if workspace.Id == "" {
		workspace.Id = "default-workspace"
	}
	m.saved = append(m.saved, workspace)
	m.workspaces = append(m.workspaces, workspace)
	return workspace.Id, nil
}

type ContextRepositoryRepairMock struct {
	ContextRepository
	contexts  []*Context
	saved     []*Context
	listError error
	saveError error
	called    bool
}

func (m *ContextRepositoryRepairMock) List() ([]*Context, error) {
	m.called = true
	return m.contexts, m.listError
}

func (m *ContextRepositoryRepairMock) Save(context *Context) (string, error) {
	if m.saveError != nil {
		return "", m.saveError
	}
	m.saved = append(m.saved, context)
	return context.Id, nil
}

type IntervalRepositoryRepairMock struct {
	IntervalRepository
	intervals []*Interval
	saved     []*Interval
	listError error
	saveError error
	called    bool
}

func (m *IntervalRepositoryRepairMock) List() ([]*Interval, error) {
	m.called = true
	return m.intervals, m.listError
}

func (m *IntervalRepositoryRepairMock) Save(interval *Interval) (string, error) {
	if m.saveError != nil {
		return "", m.saveError
	}
	m.saved = append(m.saved, interval)
	return interval.Id, nil
}

func setupManagerCorrectDataForRepair() *ContextManager {
	workspaceRepo := &WorkspaceRepositoryRepairMock{
		workspaces: []*Workspace{
			{Id: "workspace1", Name: "Default"},
			{Id: "workspace2"},
		},
	}
	contextRepo := &ContextRepositoryRepairMock{
		contexts: []*Context{
			{Id: "context1", WorkspaceId: "workspace1"},
			{Id: "context2", WorkspaceId: "workspace2"},
		},
	}
	intervalRepo := &IntervalRepositoryRepairMock{
		intervals: []*Interval{
			{Id: "interval1", ContextId: "context1", WorkspaceId: "workspace1"},
			{Id: "interval2", ContextId: "context2", WorkspaceId: "workspace2"},
		},
	}

	return NewContextManager(nil, contextRepo, intervalRepo, workspaceRepo)
}

func TestPassIntegrityRepairWithCorrectData(t *testing.T) {
	manager := setupManagerCorrectDataForRepair()

	result, err := manager.RepairIntegrity()
	require.NoError(t, err)
	require.Equal(t, 0, result.RepairedCount)
	require.True(t, result.Report.Healthy)
	require.Empty(t, result.Report.Issues)
	require.Empty(t, manager.WorkspaceRepository.(*WorkspaceRepositoryRepairMock).saved)
	require.Empty(t, manager.ContextRepository.(*ContextRepositoryRepairMock).saved)
	require.Empty(t, manager.IntervalRepository.(*IntervalRepositoryRepairMock).saved)
}

func TestIntegrityRepairCreatesDefaultWorkspace(t *testing.T) {
	manager := setupManagerCorrectDataForRepair()
	workspaceRepo := manager.WorkspaceRepository.(*WorkspaceRepositoryRepairMock)
	workspaceRepo.workspaces = []*Workspace{}
	manager.ContextRepository.(*ContextRepositoryRepairMock).contexts = []*Context{}
	manager.IntervalRepository.(*IntervalRepositoryRepairMock).intervals = []*Interval{}

	result, err := manager.RepairIntegrity()
	require.NoError(t, err)
	require.Equal(t, 1, result.RepairedCount)
	require.True(t, result.Report.Healthy)
	require.Len(t, workspaceRepo.saved, 1)
	require.Equal(t, "default-workspace", workspaceRepo.saved[0].Id)
	require.Equal(t, "Default", workspaceRepo.saved[0].Name)
}

func TestIntegrityRepairCreatesDefaultWorkspaceForUnassignedContext(t *testing.T) {
	manager := setupManagerCorrectDataForRepair()
	workspaceRepo := manager.WorkspaceRepository.(*WorkspaceRepositoryRepairMock)
	workspaceRepo.workspaces[0].Name = "Workspace 1"
	contextRepo := manager.ContextRepository.(*ContextRepositoryRepairMock)
	intervalRepo := manager.IntervalRepository.(*IntervalRepositoryRepairMock)
	contextRepo.contexts[0].WorkspaceId = ""
	intervalRepo.intervals[0].WorkspaceId = ""

	result, err := manager.RepairIntegrity()
	require.NoError(t, err)
	require.Equal(t, 3, result.RepairedCount)
	require.True(t, result.Report.Healthy)
	require.Len(t, workspaceRepo.saved, 1)
	require.Equal(t, "default-workspace", contextRepo.contexts[0].WorkspaceId)
	require.Equal(t, "default-workspace", intervalRepo.intervals[0].WorkspaceId)
}

func TestIntegrityRepairDoesNotCreateDefaultWorkspaceWhenDataIsValid(t *testing.T) {
	manager := setupManagerCorrectDataForRepair()
	workspaceRepo := manager.WorkspaceRepository.(*WorkspaceRepositoryRepairMock)
	workspaceRepo.workspaces[0].Name = "Workspace 1"

	result, err := manager.RepairIntegrity()
	require.NoError(t, err)
	require.Equal(t, 0, result.RepairedCount)
	require.True(t, result.Report.Healthy)
	require.Empty(t, workspaceRepo.saved)
}

func TestIntegrityRepairContextWithoutWorkspace(t *testing.T) {
	manager := setupManagerCorrectDataForRepair()
	contextRepo := manager.ContextRepository.(*ContextRepositoryRepairMock)
	contextRepo.contexts[0].WorkspaceId = ""
	manager.IntervalRepository.(*IntervalRepositoryRepairMock).intervals = []*Interval{}

	result, err := manager.RepairIntegrity()
	require.NoError(t, err)
	require.Equal(t, 1, result.RepairedCount)
	require.True(t, result.Report.Healthy)
	require.Len(t, contextRepo.saved, 1)
	require.Equal(t, "context1", contextRepo.saved[0].Id)
	require.Equal(t, "workspace1", contextRepo.saved[0].WorkspaceId)
}

func TestIntegrityRepairContextAndItsIntervalWithoutWorkspace(t *testing.T) {
	manager := setupManagerCorrectDataForRepair()
	contextRepo := manager.ContextRepository.(*ContextRepositoryRepairMock)
	intervalRepo := manager.IntervalRepository.(*IntervalRepositoryRepairMock)
	contextRepo.contexts[0].WorkspaceId = ""
	intervalRepo.intervals[0].WorkspaceId = ""

	report, err := manager.CheckIntegrity()
	require.NoError(t, err)
	require.Len(t, report.Issues, 2)
	require.True(t, report.Issues[0].Repairable)
	require.True(t, report.Issues[1].Repairable)

	result, err := manager.RepairIntegrity()
	require.NoError(t, err)
	require.Equal(t, 2, result.RepairedCount)
	require.True(t, result.Report.Healthy)
	require.Equal(t, "workspace1", contextRepo.contexts[0].WorkspaceId)
	require.Equal(t, "workspace1", intervalRepo.intervals[0].WorkspaceId)
}

func TestIntegrityRepairContextWithNonexistentWorkspace(t *testing.T) {
	manager := setupManagerCorrectDataForRepair()
	contextRepo := manager.ContextRepository.(*ContextRepositoryRepairMock)
	contextRepo.contexts[0].WorkspaceId = "nonexistent"
	manager.IntervalRepository.(*IntervalRepositoryRepairMock).intervals = []*Interval{}

	result, err := manager.RepairIntegrity()
	require.NoError(t, err)
	require.Equal(t, 1, result.RepairedCount)
	require.True(t, result.Report.Healthy)
	require.Len(t, contextRepo.saved, 1)
	require.Equal(t, "context1", contextRepo.saved[0].Id)
	require.Equal(t, "workspace1", contextRepo.saved[0].WorkspaceId)
}

func TestIntegrityRepairIntervalWithoutWorkspace(t *testing.T) {
	manager := setupManagerCorrectDataForRepair()
	intervalRepo := manager.IntervalRepository.(*IntervalRepositoryRepairMock)
	intervalRepo.intervals[0].WorkspaceId = ""

	result, err := manager.RepairIntegrity()
	require.NoError(t, err)
	require.Equal(t, 1, result.RepairedCount)
	require.True(t, result.Report.Healthy)
	require.Len(t, intervalRepo.saved, 1)
	require.Equal(t, "interval1", intervalRepo.saved[0].Id)
	require.Equal(t, "workspace1", intervalRepo.saved[0].WorkspaceId)
}

func TestIntegrityRepairIntervalWithNonexistentWorkspace(t *testing.T) {
	manager := setupManagerCorrectDataForRepair()
	intervalRepo := manager.IntervalRepository.(*IntervalRepositoryRepairMock)
	intervalRepo.intervals[0].WorkspaceId = "nonexistent"

	result, err := manager.RepairIntegrity()
	require.NoError(t, err)
	require.Equal(t, 1, result.RepairedCount)
	require.True(t, result.Report.Healthy)
	require.Len(t, intervalRepo.saved, 1)
	require.Equal(t, "interval1", intervalRepo.saved[0].Id)
	require.Equal(t, "workspace1", intervalRepo.saved[0].WorkspaceId)
}

func TestIntegrityRepairIntervalWorkspaceMismatch(t *testing.T) {
	manager := setupManagerCorrectDataForRepair()
	intervalRepo := manager.IntervalRepository.(*IntervalRepositoryRepairMock)
	intervalRepo.intervals[0].WorkspaceId = "workspace2"

	result, err := manager.RepairIntegrity()
	require.NoError(t, err)
	require.Equal(t, 1, result.RepairedCount)
	require.True(t, result.Report.Healthy)
	require.Len(t, intervalRepo.saved, 1)
	require.Equal(t, "interval1", intervalRepo.saved[0].Id)
	require.Equal(t, "workspace1", intervalRepo.saved[0].WorkspaceId)
}

func TestIntegrityRepairLeavesIntervalWithNonexistentContext(t *testing.T) {
	manager := setupManagerCorrectDataForRepair()
	intervalRepo := manager.IntervalRepository.(*IntervalRepositoryRepairMock)
	intervalRepo.intervals[0].ContextId = "nonexistent"

	result, err := manager.RepairIntegrity()
	require.NoError(t, err)
	require.Equal(t, 0, result.RepairedCount)
	require.False(t, result.Report.Healthy)
	require.Len(t, result.Report.Issues, 1)
	require.Equal(t, "INTERVAL_CONTEXT_NOT_FOUND", result.Report.Issues[0].Code)
	require.Empty(t, intervalRepo.saved)
}

func TestIntegrityRepairWithMultipleIssues(t *testing.T) {
	manager := setupManagerCorrectDataForRepair()
	contextRepo := manager.ContextRepository.(*ContextRepositoryRepairMock)
	intervalRepo := manager.IntervalRepository.(*IntervalRepositoryRepairMock)
	contextRepo.contexts[0].WorkspaceId = ""
	intervalRepo.intervals[0].WorkspaceId = "nonexistent"

	result, err := manager.RepairIntegrity()
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
	manager := setupManagerCorrectDataForRepair()
	manager.WorkspaceRepository.(*WorkspaceRepositoryRepairMock).listError = listError

	result, err := manager.RepairIntegrity()
	require.Nil(t, result)
	require.ErrorIs(t, err, listError)
	require.False(t, manager.ContextRepository.(*ContextRepositoryRepairMock).called)
	require.False(t, manager.IntervalRepository.(*IntervalRepositoryRepairMock).called)

	listError = errors.New("ContextRepository.List error")
	manager = setupManagerCorrectDataForRepair()
	manager.ContextRepository.(*ContextRepositoryRepairMock).listError = listError

	result, err = manager.RepairIntegrity()
	require.Nil(t, result)
	require.ErrorIs(t, err, listError)
	require.False(t, manager.IntervalRepository.(*IntervalRepositoryRepairMock).called)

	listError = errors.New("IntervalRepository.List error")
	manager = setupManagerCorrectDataForRepair()
	manager.IntervalRepository.(*IntervalRepositoryRepairMock).listError = listError

	result, err = manager.RepairIntegrity()
	require.Nil(t, result)
	require.ErrorIs(t, err, listError)
}

func TestIntegrityRepairOnSaveFail(t *testing.T) {
	saveError := errors.New("WorkspaceRepository.Save error")
	manager := setupManagerCorrectDataForRepair()
	manager.WorkspaceRepository.(*WorkspaceRepositoryRepairMock).workspaces = []*Workspace{}
	manager.WorkspaceRepository.(*WorkspaceRepositoryRepairMock).saveError = saveError

	result, err := manager.RepairIntegrity()
	require.Nil(t, result)
	require.ErrorIs(t, err, saveError)
	require.True(t, manager.ContextRepository.(*ContextRepositoryRepairMock).called)
	require.False(t, manager.IntervalRepository.(*IntervalRepositoryRepairMock).called)

	saveError = errors.New("ContextRepository.Save error")
	manager = setupManagerCorrectDataForRepair()
	manager.ContextRepository.(*ContextRepositoryRepairMock).contexts[0].WorkspaceId = ""
	manager.ContextRepository.(*ContextRepositoryRepairMock).saveError = saveError

	result, err = manager.RepairIntegrity()
	require.Nil(t, result)
	require.ErrorIs(t, err, saveError)
	require.False(t, manager.IntervalRepository.(*IntervalRepositoryRepairMock).called)

	saveError = errors.New("IntervalRepository.Save error")
	manager = setupManagerCorrectDataForRepair()
	manager.IntervalRepository.(*IntervalRepositoryRepairMock).intervals[0].WorkspaceId = ""
	manager.IntervalRepository.(*IntervalRepositoryRepairMock).saveError = saveError

	result, err = manager.RepairIntegrity()
	require.Nil(t, result)
	require.ErrorIs(t, err, saveError)
}

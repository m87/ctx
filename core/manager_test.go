package core

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type fixedTimeProvider struct {
	now time.Time
}

func (p fixedTimeProvider) Now() ZonedTime {
	return ZonedTime{Time: p.now, Timezone: "UTC"}
}

type statsIntervalRepository struct {
	intervalsByContext map[string][]*Interval
	intervals          []*Interval
	savedIntervals     []*Interval
	deletedContextID   string
	deleteByContextErr error
}

func (r *statsIntervalRepository) GetById(string) (*Interval, error) { return nil, nil }
func (r *statsIntervalRepository) Save(interval *Interval) (string, error) {
	r.savedIntervals = append(r.savedIntervals, interval)
	return interval.Id, nil
}
func (r *statsIntervalRepository) Delete(string) error { return nil }
func (r *statsIntervalRepository) DeleteByContextId(contextID string) error {
	r.deletedContextID = contextID
	return r.deleteByContextErr
}
func (r *statsIntervalRepository) ListByContextId(contextId string) ([]*Interval, error) {
	return r.intervalsByContext[contextId], nil
}
func (r *statsIntervalRepository) GetActiveIntervalByContextId(string) (*Interval, error) {
	return nil, nil
}
func (r *statsIntervalRepository) ListByDay(time.Time, string) ([]*Interval, error) {
	return nil, nil
}
func (r *statsIntervalRepository) List() ([]*Interval, error) { return r.intervals, nil }

type mockContextRepository struct {
	contexts             []*Context
	contextsByID         map[string]*Context
	savedContexts        []*Context
	deletedContextID     string
	deleteErr            error
	listByWorkspaceErr   error
	listByWorkspaceCalls int
	listedWorkspaceID    string
}

func (r *mockContextRepository) GetById(id string) (*Context, error) {
	return r.contextsByID[id], nil
}

func (r *mockContextRepository) Save(context *Context) (string, error) {
	r.savedContexts = append(r.savedContexts, context)
	return context.Id, nil
}

func (r *mockContextRepository) Delete(contextID string) error {
	r.deletedContextID = contextID
	return r.deleteErr
}

func (r *mockContextRepository) List() ([]*Context, error) {
	return r.contexts, nil
}

func (r *mockContextRepository) ListByWorkspace(workspaceID string) ([]*Context, error) {
	r.listByWorkspaceCalls++
	r.listedWorkspaceID = workspaceID
	return r.contexts, r.listByWorkspaceErr
}

func (r *mockContextRepository) GetActive() (*Context, error) {
	return nil, nil
}

type mockWorkspaceRepository struct {
	deleteErr          error
	deleteCalls        int
	deletedWorkspaceID string
	workspacesByID     map[string]*Workspace
	workspaces         []*Workspace
}

func (r *mockWorkspaceRepository) GetById(id string) (*Workspace, error) {
	return r.workspacesByID[id], nil
}

func (r *mockWorkspaceRepository) Save(*Workspace) (string, error) {
	return "", nil
}

func (r *mockWorkspaceRepository) Delete(workspaceID string) error {
	r.deleteCalls++
	r.deletedWorkspaceID = workspaceID
	return r.deleteErr
}

func (r *mockWorkspaceRepository) List() ([]*Workspace, error) {
	return r.workspaces, nil
}

func TestContextManagerEnsureDefaultWorkspaceFillsOnlyMissingAssignments(t *testing.T) {
	unassignedContext := &Context{Id: "context-1"}
	assignedContext := &Context{Id: "context-2", WorkspaceId: "workspace-2"}
	unassignedInterval := &Interval{Id: "interval-1"}
	assignedInterval := &Interval{Id: "interval-2", WorkspaceId: "workspace-2"}
	contextRepo := &mockContextRepository{contexts: []*Context{
		unassignedContext,
		assignedContext,
	}}
	intervalRepo := &statsIntervalRepository{intervals: []*Interval{
		unassignedInterval,
		assignedInterval,
	}}
	workspaceRepo := &mockWorkspaceRepository{workspaces: []*Workspace{
		{Id: "default-workspace", Name: "Default"},
		{Id: "workspace-2", Name: "Second"},
	}}
	manager := NewContextManager(nil, contextRepo, intervalRepo, workspaceRepo)

	err := manager.EnsureDefaultWorkspace()

	require.NoError(t, err)
	require.Equal(t, "default-workspace", unassignedContext.WorkspaceId)
	require.Equal(t, "workspace-2", assignedContext.WorkspaceId)
	require.Equal(t, "default-workspace", unassignedInterval.WorkspaceId)
	require.Equal(t, "workspace-2", assignedInterval.WorkspaceId)
	require.Equal(t, []*Context{unassignedContext}, contextRepo.savedContexts)
	require.Equal(t, []*Interval{unassignedInterval}, intervalRepo.savedIntervals)
}

func TestContextManagerCheckIntegrityReportsOrphansAndWorkspaceMismatch(t *testing.T) {
	now := time.Date(2026, 1, 2, 10, 0, 0, 0, time.UTC)
	contextRepo := &mockContextRepository{contexts: []*Context{
		{Id: "context-without-workspace"},
		{Id: "context-1", WorkspaceId: "workspace-1"},
	}}
	intervalRepo := &statsIntervalRepository{intervals: []*Interval{
		{Id: "missing-context", ContextId: "does-not-exist", WorkspaceId: "workspace-1", Status: "completed", Start: ZonedTime{Time: now, Timezone: "UTC"}, End: ZonedTime{Time: now.Add(time.Hour), Timezone: "UTC"}},
		{Id: "workspace-mismatch", ContextId: "context-1", WorkspaceId: "workspace-2", Status: "completed", Start: ZonedTime{Time: now.Add(2 * time.Hour), Timezone: "UTC"}, End: ZonedTime{Time: now.Add(3 * time.Hour), Timezone: "UTC"}},
	}}
	workspaceRepo := &mockWorkspaceRepository{workspaces: []*Workspace{
		{Id: "workspace-1", Name: "First"},
		{Id: "workspace-2", Name: "Second"},
	}}
	manager := NewContextManager(nil, contextRepo, intervalRepo, workspaceRepo)

	report, err := manager.CheckIntegrity()

	require.NoError(t, err)
	require.False(t, report.Healthy)
	require.Equal(t, 2, report.WorkspaceCount)
	require.Equal(t, 2, report.ContextCount)
	require.Equal(t, 2, report.IntervalCount)
	require.Equal(t, []string{
		"CONTEXT_MISSING_WORKSPACE",
		"INTERVAL_CONTEXT_NOT_FOUND",
		"INTERVAL_WORKSPACE_MISMATCH",
	}, integrityIssueCodes(report.Issues))
}

func TestContextManagerRepairIntegrityRepairsWorkspaceAssignments(t *testing.T) {
	now := time.Date(2026, 1, 2, 10, 0, 0, 0, time.UTC)
	context := &Context{Id: "context-1", WorkspaceId: "missing-workspace"}
	matchingInterval := &Interval{Id: "interval-1", ContextId: context.Id, WorkspaceId: "other-workspace", Status: "completed", Start: ZonedTime{Time: now, Timezone: "UTC"}, End: ZonedTime{Time: now.Add(time.Hour), Timezone: "UTC"}}
	orphanInterval := &Interval{Id: "interval-2", ContextId: "missing-context", WorkspaceId: "default-workspace", Status: "completed", Start: ZonedTime{Time: now.Add(2 * time.Hour), Timezone: "UTC"}, End: ZonedTime{Time: now.Add(3 * time.Hour), Timezone: "UTC"}}
	contextRepo := &mockContextRepository{contexts: []*Context{context}}
	intervalRepo := &statsIntervalRepository{intervals: []*Interval{matchingInterval, orphanInterval}}
	workspaceRepo := &mockWorkspaceRepository{workspaces: []*Workspace{
		{Id: "default-workspace", Name: "Default"},
	}}
	manager := NewContextManager(nil, contextRepo, intervalRepo, workspaceRepo)

	result, err := manager.RepairIntegrity()

	require.NoError(t, err)
	require.Equal(t, 2, result.RepairedCount)
	require.Equal(t, "default-workspace", context.WorkspaceId)
	require.Equal(t, "default-workspace", matchingInterval.WorkspaceId)
	require.False(t, result.Report.Healthy)
	require.Equal(t, []string{"INTERVAL_CONTEXT_NOT_FOUND"}, integrityIssueCodes(result.Report.Issues))
}

func integrityIssueCodes(issues []*IntegrityIssue) []string {
	codes := make([]string, 0, len(issues))
	for _, issue := range issues {
		codes = append(codes, issue.Code)
	}
	return codes
}

func TestContextManagerCreateContextAssignsWorkspace(t *testing.T) {
	contextRepo := &mockContextRepository{}
	workspaceRepo := &mockWorkspaceRepository{workspacesByID: map[string]*Workspace{
		"workspace-1": {Id: "workspace-1", Name: "First"},
	}}
	manager := NewContextManager(nil, contextRepo, nil, workspaceRepo)
	context := &Context{Name: "Context", WorkspaceId: "workspace-1"}

	_, err := manager.CreateContext(context)

	require.NoError(t, err)
}

func TestContextManagerCreateContextRequiresExistingWorkspace(t *testing.T) {
	manager := NewContextManager(
		nil,
		&mockContextRepository{},
		nil,
		&mockWorkspaceRepository{},
	)

	_, err := manager.CreateContext(&Context{Name: "Context", WorkspaceId: "missing"})

	var workspaceNotFoundErr *WorkspaceNotFoundError
	require.ErrorAs(t, err, &workspaceNotFoundErr)
	require.Equal(t, "missing", workspaceNotFoundErr.WorkspaceId)
}

func TestContextManagerUpdateContextPreservesWorkspaceWhenPayloadOmitsIt(t *testing.T) {
	contextRepo := &mockContextRepository{contextsByID: map[string]*Context{
		"context-1": {Id: "context-1", Name: "Old", WorkspaceId: "workspace-1"},
	}}
	manager := NewContextManager(nil, contextRepo, nil, nil)
	updated := &Context{Id: "context-1", Name: "New"}

	err := manager.UpdateContext(updated)

	require.NoError(t, err)
	require.Equal(t, "workspace-1", updated.WorkspaceId)
	require.Equal(t, []*Context{updated}, contextRepo.savedContexts)
}

func TestContextManagerUpdateContextRejectsWorkspaceMove(t *testing.T) {
	contextRepo := &mockContextRepository{contextsByID: map[string]*Context{
		"context-1": {Id: "context-1", WorkspaceId: "workspace-1"},
	}}
	manager := NewContextManager(nil, contextRepo, nil, nil)

	err := manager.UpdateContext(&Context{Id: "context-1", WorkspaceId: "workspace-2"})

	var moveErr *ContextWorkspaceMoveNotAllowedError
	require.ErrorAs(t, err, &moveErr)
	require.Equal(t, "workspace-1", moveErr.FromWorkspaceId)
	require.Equal(t, "workspace-2", moveErr.ToWorkspaceId)
	require.Empty(t, contextRepo.savedContexts)
}

func TestContextManagerSaveIntervalUsesContextWorkspace(t *testing.T) {
	contextRepo := &mockContextRepository{contextsByID: map[string]*Context{
		"context-2": {Id: "context-2", WorkspaceId: "workspace-2"},
	}}
	intervalRepo := &statsIntervalRepository{}
	manager := NewContextManager(nil, contextRepo, intervalRepo, nil)
	interval := &Interval{
		Id:          "interval-1",
		ContextId:   "context-2",
		WorkspaceId: "workspace-1",
	}

	_, err := manager.SaveInterval(interval)

	require.NoError(t, err)
	require.Equal(t, "workspace-2", interval.WorkspaceId)
	require.Equal(t, []*Interval{interval}, intervalRepo.savedIntervals)
}

func TestContextManagerSaveIntervalRejectsMissingContext(t *testing.T) {
	manager := NewContextManager(
		nil,
		&mockContextRepository{},
		&statsIntervalRepository{},
		nil,
	)

	_, err := manager.SaveInterval(&Interval{ContextId: "missing"})

	var contextNotFoundErr *ContextNotFoundError
	require.ErrorAs(t, err, &contextNotFoundErr)
	require.Equal(t, "missing", contextNotFoundErr.ContextId)
}

func TestContextManagerDeleteContextDeletesIntervals(t *testing.T) {
	contextRepo := &mockContextRepository{}
	intervalRepo := &statsIntervalRepository{}
	manager := NewContextManager(nil, contextRepo, intervalRepo, nil)

	err := manager.DeleteContext("context-1")

	require.NoError(t, err)
	require.Equal(t, "context-1", intervalRepo.deletedContextID)
	require.Equal(t, "context-1", contextRepo.deletedContextID)
}

func TestContextManagerDeleteContextStopsWhenIntervalDeleteFails(t *testing.T) {
	wantErr := errors.New("delete intervals failed")
	contextRepo := &mockContextRepository{}
	intervalRepo := &statsIntervalRepository{deleteByContextErr: wantErr}
	manager := NewContextManager(nil, contextRepo, intervalRepo, nil)

	err := manager.DeleteContext("context-1")

	require.ErrorIs(t, err, wantErr)
	require.Equal(t, "context-1", intervalRepo.deletedContextID)
	require.Empty(t, contextRepo.deletedContextID)
}

func TestContextManagerDeleteWorkspaceDeletesUnusedWorkspace(t *testing.T) {
	contextRepo := &mockContextRepository{}
	workspaceRepo := &mockWorkspaceRepository{}
	manager := NewContextManager(nil, contextRepo, nil, workspaceRepo)

	err := manager.DeleteWorkspace("workspace-1")

	require.NoError(t, err)
	require.Equal(t, 1, contextRepo.listByWorkspaceCalls)
	require.Equal(t, "workspace-1", contextRepo.listedWorkspaceID)
	require.Equal(t, 1, workspaceRepo.deleteCalls)
	require.Equal(t, "workspace-1", workspaceRepo.deletedWorkspaceID)
}

func TestContextManagerDeleteWorkspaceReturnsWorkspaceInUseError(t *testing.T) {
	contextRepo := &mockContextRepository{
		contexts: []*Context{{Id: "context-1", WorkspaceId: "workspace-1"}},
	}
	workspaceRepo := &mockWorkspaceRepository{}
	manager := NewContextManager(nil, contextRepo, nil, workspaceRepo)

	err := manager.DeleteWorkspace("workspace-1")

	var workspaceInUseErr *WorkspaceInUseError
	require.ErrorAs(t, err, &workspaceInUseErr)
	require.Equal(t, "workspace-1", workspaceInUseErr.WorkspaceId)
	require.Equal(t, 0, workspaceRepo.deleteCalls)
}

func TestContextManagerDeleteWorkspaceReturnsContextRepositoryError(t *testing.T) {
	wantErr := errors.New("list contexts failed")
	contextRepo := &mockContextRepository{listByWorkspaceErr: wantErr}
	workspaceRepo := &mockWorkspaceRepository{}
	manager := NewContextManager(nil, contextRepo, nil, workspaceRepo)

	err := manager.DeleteWorkspace("workspace-1")

	require.ErrorIs(t, err, wantErr)
	require.Equal(t, 0, workspaceRepo.deleteCalls)
}

func TestContextManagerDeleteWorkspaceReturnsWorkspaceRepositoryError(t *testing.T) {
	wantErr := errors.New("delete workspace failed")
	contextRepo := &mockContextRepository{}
	workspaceRepo := &mockWorkspaceRepository{deleteErr: wantErr}
	manager := NewContextManager(nil, contextRepo, nil, workspaceRepo)

	err := manager.DeleteWorkspace("workspace-1")

	require.ErrorIs(t, err, wantErr)
	require.Equal(t, 1, workspaceRepo.deleteCalls)
	require.Equal(t, "workspace-1", workspaceRepo.deletedWorkspaceID)
}

func TestContextManagerGetWorkspaceStatsUsesAllIntervals(t *testing.T) {
	now := time.Date(2026, time.June, 14, 12, 0, 0, 0, time.UTC)
	contextRepo := &mockContextRepository{contexts: []*Context{
		{Id: "context-1", Name: "First", WorkspaceId: "workspace-1"},
		{Id: "context-2", Name: "Second", WorkspaceId: "workspace-1"},
	}}
	intervalRepo := &statsIntervalRepository{intervalsByContext: map[string][]*Interval{
		"context-1": {
			{
				Start:  ZonedTime{Time: now.Add(-3 * time.Hour)},
				End:    ZonedTime{Time: now.Add(-2 * time.Hour)},
				Status: "completed",
			},
			{
				Start:  ZonedTime{Time: now.Add(-30 * time.Minute)},
				Status: "active",
			},
		},
		"context-2": {
			{Duration: 30 * time.Minute, Status: "completed"},
		},
	}}
	manager := NewContextManager(
		fixedTimeProvider{now: now},
		contextRepo,
		intervalRepo,
		&mockWorkspaceRepository{},
	)

	stats, err := manager.GetWorkspaceStats("workspace-1")

	require.NoError(t, err)
	require.Equal(t, 2*time.Hour, stats.TotalDuration)
	require.Equal(t, 3, stats.TotalSessions)
	require.Len(t, stats.Contexts, 2)
	require.Len(t, stats.ContextStats, 2)
	require.Equal(t, "context-1", stats.ContextStats[0].ContextId)
	require.Equal(t, 90*time.Minute, stats.ContextStats[0].Duration)
	require.InDelta(t, 75, stats.ContextStats[0].Percentage, 0.001)
}

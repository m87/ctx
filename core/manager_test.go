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
	savedIntervals     []*Interval
}

func (r *statsIntervalRepository) GetById(string) (*Interval, error) { return nil, nil }
func (r *statsIntervalRepository) Save(interval *Interval) (string, error) {
	r.savedIntervals = append(r.savedIntervals, interval)
	return interval.Id, nil
}
func (r *statsIntervalRepository) Delete(string) error { return nil }
func (r *statsIntervalRepository) ListByContextId(contextId string) ([]*Interval, error) {
	return r.intervalsByContext[contextId], nil
}
func (r *statsIntervalRepository) GetActiveIntervalByContextId(string) (*Interval, error) {
	return nil, nil
}
func (r *statsIntervalRepository) ListByDay(time.Time, string) ([]*Interval, error) {
	return nil, nil
}
func (r *statsIntervalRepository) List() ([]*Interval, error) { return nil, nil }

type mockContextRepository struct {
	contexts             []*Context
	contextsByID         map[string]*Context
	listByWorkspaceErr   error
	listByWorkspaceCalls int
	listedWorkspaceID    string
}

func (r *mockContextRepository) GetById(id string) (*Context, error) {
	return r.contextsByID[id], nil
}

func (r *mockContextRepository) Save(*Context) (string, error) {
	return "", nil
}

func (r *mockContextRepository) Delete(string) error {
	return nil
}

func (r *mockContextRepository) List() ([]*Context, error) {
	return nil, nil
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
}

func (r *mockWorkspaceRepository) GetById(string) (*Workspace, error) {
	return nil, nil
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
	return nil, nil
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

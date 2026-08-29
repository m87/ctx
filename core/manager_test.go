package core

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSaveIntervalNormalizesIncomingInstantsToUTC(t *testing.T) {
	tokyo, err := time.LoadLocation("Asia/Tokyo")
	require.NoError(t, err)
	start := time.Date(2026, 8, 2, 3, 0, 0, 0, tokyo)
	end := start.Add(30 * time.Minute)
	intervalRepo := &statsIntervalRepository{}
	manager := NewContextManager(
		fixedTimeProvider{now: end},
		&mockContextRepository{contextsByID: map[string]*Context{
			"context-1": {Id: "context-1", WorkspaceId: "workspace-1"},
		}},
		intervalRepo,
		&mockWorkspaceRepository{},
		&mockProjectRepository{},
	)
	interval := &Interval{ContextId: "context-1", Start: &start, End: &end}

	_, err = manager.SaveInterval(interval)

	require.NoError(t, err)
	require.Equal(t, time.UTC, interval.Start.Location())
	require.Equal(t, "2026-08-01T18:00:00Z", interval.Start.Format(time.RFC3339))
	require.Equal(t, time.UTC, interval.End.Location())
	encoded, err := json.Marshal(interval)
	require.NoError(t, err)
	require.Contains(t, string(encoded), `"start":"2026-08-01T18:00:00Z"`)
	require.NotContains(t, string(encoded), "Asia/Tokyo")
}

func TestClipIntervalRangeToDayUsesSelectedTimeZone(t *testing.T) {
	tokyo, err := time.LoadLocation("Asia/Tokyo")
	require.NoError(t, err)
	date := time.Date(2026, 8, 2, 0, 0, 0, 0, tokyo)
	start := time.Date(2026, 8, 1, 14, 30, 0, 0, time.UTC)
	end := time.Date(2026, 8, 1, 15, 30, 0, 0, time.UTC)

	rng, ok := ClipIntervalRangeToDay(
		&Interval{Start: &start, End: &end, Status: "completed"},
		date,
		end,
	)

	require.True(t, ok)
	require.Equal(t, "2026-08-01T15:00:00Z", rng.Start.Format(time.RFC3339))
	require.Equal(t, "2026-08-01T15:30:00Z", rng.End.Format(time.RFC3339))
}

func TestClipIntervalRangeToDayHandlesDSTDayLength(t *testing.T) {
	newYork, err := time.LoadLocation("America/New_York")
	require.NoError(t, err)
	date := time.Date(2026, 3, 8, 0, 0, 0, 0, newYork)
	start := time.Date(2026, 3, 8, 4, 0, 0, 0, time.UTC)
	end := time.Date(2026, 3, 9, 5, 0, 0, 0, time.UTC)

	rng, ok := ClipIntervalRangeToDay(
		&Interval{Start: &start, End: &end, Status: "completed"},
		date,
		end,
	)

	require.True(t, ok)
	require.Equal(t, "2026-03-08T05:00:00Z", rng.Start.Format(time.RFC3339))
	require.Equal(t, "2026-03-09T04:00:00Z", rng.End.Format(time.RFC3339))
	require.Equal(t, 23*time.Hour, rng.End.Sub(rng.Start))
}

type fixedTimeProvider struct {
	now time.Time
}

func (p fixedTimeProvider) Now() time.Time {
	return p.now.UTC()
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

func (r *statsIntervalRepository) SaveAll(intervals []*Interval) ([]string, error) {
	ids := make([]string, 0, len(intervals))
	for _, interval := range intervals {
		id, err := r.Save(interval)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
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

func (r *statsIntervalRepository) ListToSync(limit int) ([]*Interval, error) {
	result := make([]*Interval, 0, len(r.intervals))
	for _, interval := range r.intervals {
		if interval == nil || interval.Synced {
			continue
		}
		result = append(result, interval)
		if limit > 0 && len(result) == limit {
			break
		}
	}
	return result, nil
}

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

func (r *mockContextRepository) SaveAll(contexts []*Context) ([]string, error) {
	ids := make([]string, 0, len(contexts))
	for _, context := range contexts {
		id, err := r.Save(context)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *mockContextRepository) Delete(contextID string) error {
	r.deletedContextID = contextID
	return r.deleteErr
}

func (r *mockContextRepository) List() ([]*Context, error) {
	return r.contexts, nil
}

func (r *mockContextRepository) ListToSync(limit int) ([]*Context, error) {
	result := make([]*Context, 0, len(r.contexts))
	for _, context := range r.contexts {
		if context == nil {
			continue
		}
		result = append(result, context)
		if limit > 0 && len(result) == limit {
			break
		}
	}
	return result, nil
}

func (r *mockContextRepository) ListByWorkspace(workspaceID string) ([]*Context, error) {
	r.listByWorkspaceCalls++
	r.listedWorkspaceID = workspaceID
	return r.contexts, r.listByWorkspaceErr
}

func (r *mockContextRepository) ListByProject(projectID string) ([]*Context, error) {
	result := make([]*Context, 0)
	for _, context := range r.contexts {
		if context != nil && context.Project != nil && context.Project.Id == projectID {
			result = append(result, context)
		}
	}
	return result, nil
}

func (r *mockContextRepository) ListByWorkspaceIncludingArchived(workspaceID string) ([]*Context, error) {
	return r.ListByWorkspace(workspaceID)
}

func (r *mockContextRepository) GetActive() (*Context, error) {
	return nil, nil
}

type mockProjectRepository struct {
	projectsByID     map[string]*Project
	projects         []*Project
	savedProjects    []*Project
	deletedProjectID string
}

func (r *mockProjectRepository) GetById(id string) (*Project, error) {
	return r.projectsByID[id], nil
}

func (r *mockProjectRepository) Save(project *Project) (string, error) {
	r.savedProjects = append(r.savedProjects, project)
	return project.Id, nil
}

func (r *mockProjectRepository) SaveAll(projects []*Project) ([]string, error) {
	ids := make([]string, 0, len(projects))
	for _, project := range projects {
		id, err := r.Save(project)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *mockProjectRepository) Delete(projectID string) error {
	r.deletedProjectID = projectID
	return nil
}

func (r *mockProjectRepository) List(workspaceID string) ([]*Project, error) {
	return r.projects, nil
}

func (r *mockProjectRepository) ListIncludingArchived(workspaceID string) ([]*Project, error) {
	return r.projects, nil
}

func (r *mockProjectRepository) ListChildren(parentID string) ([]*Project, error) {
	result := make([]*Project, 0)
	for _, project := range r.projects {
		if project != nil && project.ParentId == parentID {
			result = append(result, project)
		}
	}
	return result, nil
}

func (r *mockProjectRepository) ListToSync(limit int) ([]*Project, error) {
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

func (r *mockWorkspaceRepository) SaveAll(workspaces []*Workspace) ([]string, error) {
	ids := make([]string, 0, len(workspaces))
	for _, workspace := range workspaces {
		id, err := r.Save(workspace)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *mockWorkspaceRepository) Delete(workspaceID string) error {
	r.deleteCalls++
	r.deletedWorkspaceID = workspaceID
	return r.deleteErr
}

func (r *mockWorkspaceRepository) List() ([]*Workspace, error) {
	return r.workspaces, nil
}

func (r *mockWorkspaceRepository) ListToSync(limit int) ([]*Workspace, error) {
	result := make([]*Workspace, 0, len(r.workspaces))
	for _, workspace := range r.workspaces {
		if workspace == nil || workspace.Synced {
			continue
		}
		result = append(result, workspace)
		if limit > 0 && len(result) == limit {
			break
		}
	}
	return result, nil
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
	projectRepo := &mockProjectRepository{projects: []*Project{
		{Id: "project-1", WorkspaceId: "default-workspace"},
		{Id: "project-2", WorkspaceId: "workspace-2"},
	}}
	manager := NewContextManager(nil, contextRepo, intervalRepo, workspaceRepo, projectRepo)

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
		{Id: "missing-context", ContextId: "does-not-exist", WorkspaceId: "workspace-1", Status: "completed", Start: testTime(now), End: testTime(now.Add(time.Hour))},
		{Id: "workspace-mismatch", ContextId: "context-1", WorkspaceId: "workspace-2", Status: "completed", Start: testTime(now.Add(2 * time.Hour)), End: testTime(now.Add(3 * time.Hour))},
	}}
	workspaceRepo := &mockWorkspaceRepository{workspaces: []*Workspace{
		{Id: "workspace-1", Name: "First"},
		{Id: "workspace-2", Name: "Second"},
	}}
	projectRepo := &mockProjectRepository{projects: []*Project{
		{Id: "project-1", WorkspaceId: "workspace-1"},
		{Id: "project-2", WorkspaceId: "workspace-2"},
	}}
	manager := NewContextManager(nil, contextRepo, intervalRepo, workspaceRepo, projectRepo)

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
	matchingInterval := &Interval{Id: "interval-1", ContextId: context.Id, WorkspaceId: "other-workspace", Status: "completed", Start: testTime(now), End: testTime(now.Add(time.Hour))}
	orphanInterval := &Interval{Id: "interval-2", ContextId: "missing-context", WorkspaceId: "default-workspace", Status: "completed", Start: testTime(now.Add(2 * time.Hour)), End: testTime(now.Add(3 * time.Hour))}
	contextRepo := &mockContextRepository{contexts: []*Context{context}}
	intervalRepo := &statsIntervalRepository{intervals: []*Interval{matchingInterval, orphanInterval}}
	workspaceRepo := &mockWorkspaceRepository{workspaces: []*Workspace{
		{Id: "default-workspace", Name: "Default"},
	}}
	projectRepo := &mockProjectRepository{projects: []*Project{
		{Id: "project-1", WorkspaceId: "default-workspace"},
	}}
	manager := NewContextManager(nil, contextRepo, intervalRepo, workspaceRepo, projectRepo)

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
	projectRepo := &mockProjectRepository{projectsByID: map[string]*Project{
		"project-1": {Id: "project-1", WorkspaceId: "workspace-1"},
	}}
	manager := NewContextManager(nil, contextRepo, nil, workspaceRepo, projectRepo)
	context := &Context{Name: "Context", WorkspaceId: "workspace-1"}

	_, err := manager.CreateContext(context)

	require.NoError(t, err)
}

func TestContextManagerCreateContextCanonicalizesAssignedProject(t *testing.T) {
	contextRepo := &mockContextRepository{}
	workspaceRepo := &mockWorkspaceRepository{workspacesByID: map[string]*Workspace{
		"workspace-1": {Id: "workspace-1", Name: "First"},
	}}
	projectRepo := &mockProjectRepository{projectsByID: map[string]*Project{
		"project-1": {Id: "project-1", Name: "Canonical project", WorkspaceId: "workspace-1"},
	}}
	manager := NewContextManager(nil, contextRepo, nil, workspaceRepo, projectRepo)
	context := &Context{
		Name:        "Context",
		WorkspaceId: "workspace-1",
		Project:     &ProjectMetadata{Id: "project-1", Name: "Stale name"},
	}

	_, err := manager.CreateContext(context)

	require.NoError(t, err)
	require.Equal(t, &ProjectMetadata{Id: "project-1", Name: "Canonical project"}, context.Project)
}

func TestContextManagerCreateContextRejectsProjectFromAnotherWorkspace(t *testing.T) {
	manager := NewContextManager(
		nil,
		&mockContextRepository{},
		nil,
		&mockWorkspaceRepository{workspacesByID: map[string]*Workspace{
			"workspace-1": {Id: "workspace-1"},
		}},
		&mockProjectRepository{projectsByID: map[string]*Project{
			"project-2": {Id: "project-2", WorkspaceId: "workspace-2"},
		}},
	)

	_, err := manager.CreateContext(&Context{
		Name:        "Context",
		WorkspaceId: "workspace-1",
		Project:     &ProjectMetadata{Id: "project-2"},
	})

	var mismatchErr *ProjectWorkspaceMismatchError
	require.ErrorAs(t, err, &mismatchErr)
}

func TestContextManagerCreateContextRequiresExistingWorkspace(t *testing.T) {
	manager := NewContextManager(
		nil,
		&mockContextRepository{},
		nil,
		&mockWorkspaceRepository{},
		&mockProjectRepository{},
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
	manager := NewContextManager(nil, contextRepo, nil, nil, nil)
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
	manager := NewContextManager(nil, contextRepo, nil, nil, nil)

	err := manager.UpdateContext(&Context{Id: "context-1", WorkspaceId: "workspace-2"})

	var moveErr *ContextWorkspaceMoveNotAllowedError
	require.ErrorAs(t, err, &moveErr)
	require.Equal(t, "workspace-1", moveErr.FromWorkspaceId)
	require.Equal(t, "workspace-2", moveErr.ToWorkspaceId)
	require.Empty(t, contextRepo.savedContexts)
}

func TestContextManagerUpdateContextCanUnassignProject(t *testing.T) {
	existing := &Context{
		Id:          "context-1",
		WorkspaceId: "workspace-1",
		Project:     &ProjectMetadata{Id: "project-1", Name: "Project"},
	}
	contextRepo := &mockContextRepository{contextsByID: map[string]*Context{
		"context-1": existing,
	}}
	manager := NewContextManager(nil, contextRepo, nil, nil, &mockProjectRepository{})
	updated := &Context{Id: "context-1", Name: "Context"}

	err := manager.UpdateContext(updated)

	require.NoError(t, err)
	require.Nil(t, updated.Project)
	require.Equal(t, "workspace-1", updated.WorkspaceId)
}

func TestContextManagerUpdateProjectRejectsDescendantAsParent(t *testing.T) {
	projectRepo := &mockProjectRepository{projectsByID: map[string]*Project{
		"project-1": {Id: "project-1", WorkspaceId: "workspace-1"},
		"project-2": {Id: "project-2", ParentId: "project-1", WorkspaceId: "workspace-1"},
	}}
	manager := NewContextManager(nil, &mockContextRepository{}, nil, nil, projectRepo)

	err := manager.UpdateProject(&Project{
		Id:       "project-1",
		Name:     "Parent",
		ParentId: "project-2",
	})

	var cycleErr *ProjectHierarchyCycleError
	require.ErrorAs(t, err, &cycleErr)
	require.Empty(t, projectRepo.savedProjects)
}

func TestContextManagerUpdateProjectRefreshesContextMetadata(t *testing.T) {
	context := &Context{
		Id:          "context-1",
		WorkspaceId: "workspace-1",
		Project:     &ProjectMetadata{Id: "project-1", Name: "Old name"},
	}
	contextRepo := &mockContextRepository{contexts: []*Context{context}}
	projectRepo := &mockProjectRepository{projectsByID: map[string]*Project{
		"project-1": {Id: "project-1", Name: "Old name", WorkspaceId: "workspace-1"},
	}}
	manager := NewContextManager(nil, contextRepo, nil, nil, projectRepo)

	err := manager.UpdateProject(&Project{Id: "project-1", Name: "New name"})

	require.NoError(t, err)
	require.Equal(t, &ProjectMetadata{Id: "project-1", Name: "New name"}, context.Project)
	require.Equal(t, []*Context{context}, contextRepo.savedContexts)
}

func TestContextManagerDeleteProjectMovesContentsToParent(t *testing.T) {
	context := &Context{
		Id:          "context-1",
		WorkspaceId: "workspace-1",
		Project:     &ProjectMetadata{Id: "project-2", Name: "Child"},
	}
	childProject := &Project{
		Id:          "project-3",
		ParentId:    "project-2",
		WorkspaceId: "workspace-1",
	}
	contextRepo := &mockContextRepository{contexts: []*Context{context}}
	projectRepo := &mockProjectRepository{
		projectsByID: map[string]*Project{
			"project-1": {Id: "project-1", Name: "Parent", WorkspaceId: "workspace-1"},
			"project-2": {Id: "project-2", Name: "Deleted", ParentId: "project-1", WorkspaceId: "workspace-1"},
		},
		projects: []*Project{childProject},
	}
	manager := NewContextManager(nil, contextRepo, nil, nil, projectRepo)

	err := manager.DeleteProject("project-2")

	require.NoError(t, err)
	require.Equal(t, "project-1", childProject.ParentId)
	require.Equal(t, &ProjectMetadata{Id: "project-1", Name: "Parent"}, context.Project)
	require.Equal(t, "project-2", projectRepo.deletedProjectID)
}

func TestContextManagerSaveIntervalUsesContextWorkspace(t *testing.T) {
	contextRepo := &mockContextRepository{contextsByID: map[string]*Context{
		"context-2": {Id: "context-2", WorkspaceId: "workspace-2"},
	}}
	intervalRepo := &statsIntervalRepository{}
	manager := NewContextManager(nil, contextRepo, intervalRepo, nil, nil)
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
	manager := NewContextManager(nil, contextRepo, intervalRepo, nil, nil)

	err := manager.DeleteContext("context-1")

	require.NoError(t, err)
	require.Equal(t, "context-1", intervalRepo.deletedContextID)
	require.Equal(t, "context-1", contextRepo.deletedContextID)
}

func TestContextManagerDeleteContextStopsWhenIntervalDeleteFails(t *testing.T) {
	wantErr := errors.New("delete intervals failed")
	contextRepo := &mockContextRepository{}
	intervalRepo := &statsIntervalRepository{deleteByContextErr: wantErr}
	manager := NewContextManager(nil, contextRepo, intervalRepo, nil, nil)

	err := manager.DeleteContext("context-1")

	require.ErrorIs(t, err, wantErr)
	require.Equal(t, "context-1", intervalRepo.deletedContextID)
	require.Empty(t, contextRepo.deletedContextID)
}

func TestContextManagerDeleteWorkspaceDeletesUnusedWorkspace(t *testing.T) {
	contextRepo := &mockContextRepository{}
	workspaceRepo := &mockWorkspaceRepository{}
	manager := NewContextManager(nil, contextRepo, nil, workspaceRepo, nil)

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
	manager := NewContextManager(nil, contextRepo, nil, workspaceRepo, nil)

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
	manager := NewContextManager(nil, contextRepo, nil, workspaceRepo, nil)

	err := manager.DeleteWorkspace("workspace-1")

	require.ErrorIs(t, err, wantErr)
	require.Equal(t, 0, workspaceRepo.deleteCalls)
}

func TestContextManagerDeleteWorkspaceReturnsWorkspaceRepositoryError(t *testing.T) {
	wantErr := errors.New("delete workspace failed")
	contextRepo := &mockContextRepository{}
	workspaceRepo := &mockWorkspaceRepository{deleteErr: wantErr}
	manager := NewContextManager(nil, contextRepo, nil, workspaceRepo, nil)

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
				Start:  testTime(now.Add(-3 * time.Hour)),
				End:    testTime(now.Add(-2 * time.Hour)),
				Status: "completed",
			},
			{
				Start:  testTime(now.Add(-30 * time.Minute)),
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
		&mockProjectRepository{},
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

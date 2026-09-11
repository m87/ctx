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
	test := NewEmptyTestContextManager()
	test.TimeProvider.Set(end)
	test.Contexts.Seed(&Context{Id: "context-1", WorkspaceId: "workspace-1"})
	interval := &Interval{ContextId: "context-1", Start: &start, End: &end}

	_, err = test.Manager.SaveInterval(interval)

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

func TestContextManagerEnsureDefaultWorkspaceFillsOnlyMissingAssignments(t *testing.T) {
	test := NewEmptyTestContextManager()
	unassignedContext := &Context{Id: "context-1"}
	assignedContext := &Context{Id: "context-2", WorkspaceId: "workspace-2"}
	unassignedInterval := &Interval{Id: "interval-1"}
	assignedInterval := &Interval{Id: "interval-2", WorkspaceId: "workspace-2"}
	test.Contexts.Seed(
		unassignedContext,
		assignedContext,
	)
	test.Intervals.Seed(
		unassignedInterval,
		assignedInterval,
	)
	test.Workspaces.Seed(
		&Workspace{Id: "default-workspace", Name: "Default"},
		&Workspace{Id: "workspace-2", Name: "Second"},
	)
	test.Projects.Seed(
		&Project{Id: "project-1", WorkspaceId: "default-workspace"},
		&Project{Id: "project-2", WorkspaceId: "workspace-2"},
	)

	err := test.Manager.EnsureDefaultWorkspace()

	require.NoError(t, err)
	require.Equal(t, "default-workspace", unassignedContext.WorkspaceId)
	require.Equal(t, "workspace-2", assignedContext.WorkspaceId)
	require.Equal(t, "default-workspace", unassignedInterval.WorkspaceId)
	require.Equal(t, "workspace-2", assignedInterval.WorkspaceId)
	require.Equal(t, []*Context{unassignedContext}, test.Contexts.saved)
	require.Equal(t, []*Interval{unassignedInterval}, test.Intervals.saved)
}

func TestContextManagerCheckIntegrityReportsOrphansAndWorkspaceMismatch(t *testing.T) {
	test := NewEmptyTestContextManager()
	now := time.Date(2026, 1, 2, 10, 0, 0, 0, time.UTC)
	test.Contexts.Seed(
		&Context{Id: "context-without-workspace"},
		&Context{Id: "context-1", WorkspaceId: "workspace-1"},
	)
	test.Intervals.Seed(
		&Interval{Id: "missing-context", ContextId: "does-not-exist", WorkspaceId: "workspace-1", Status: "completed", Start: testTime(now), End: testTime(now.Add(time.Hour))},
		&Interval{Id: "workspace-mismatch", ContextId: "context-1", WorkspaceId: "workspace-2", Status: "completed", Start: testTime(now.Add(2 * time.Hour)), End: testTime(now.Add(3 * time.Hour))},
	)
	test.Workspaces.Seed(
		&Workspace{Id: "workspace-1", Name: "First"},
		&Workspace{Id: "workspace-2", Name: "Second"},
	)
	test.Projects.Seed(
		&Project{Id: "project-1", WorkspaceId: "workspace-1"},
		&Project{Id: "project-2", WorkspaceId: "workspace-2"},
	)

	report, err := test.Manager.CheckIntegrity()

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
	test := NewEmptyTestContextManager()
	now := time.Date(2026, 1, 2, 10, 0, 0, 0, time.UTC)
	context := &Context{Id: "context-1", WorkspaceId: "missing-workspace"}
	matchingInterval := &Interval{Id: "interval-1", ContextId: context.Id, WorkspaceId: "other-workspace", Status: "completed", Start: testTime(now), End: testTime(now.Add(time.Hour))}
	orphanInterval := &Interval{Id: "interval-2", ContextId: "missing-context", WorkspaceId: "default-workspace", Status: "completed", Start: testTime(now.Add(2 * time.Hour)), End: testTime(now.Add(3 * time.Hour))}
	test.Contexts.Seed(context)
	test.Intervals.Seed(matchingInterval, orphanInterval)
	test.Workspaces.Seed(
		&Workspace{Id: "default-workspace", Name: "Default"},
	)
	test.Projects.Seed(
		&Project{Id: "project-1", WorkspaceId: "default-workspace"},
	)

	result, err := test.Manager.RepairIntegrity()

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
	test := NewTestContextManager()
	context := &Context{Name: "Context", WorkspaceId: TestWorkspaceID}

	_, err := test.Manager.CreateContext(context)

	require.NoError(t, err)
}

func TestContextManagerCreateContextCanonicalizesAssignedProject(t *testing.T) {
	test := NewEmptyTestContextManager()
	test.Workspaces.Seed(&Workspace{Id: "workspace-1", Name: "First"})
	test.Projects.Seed(&Project{Id: "project-1", Name: "Canonical project", WorkspaceId: "workspace-1"})
	context := &Context{
		Name:        "Context",
		WorkspaceId: "workspace-1",
		Project:     &ProjectMetadata{Id: "project-1", Name: "Stale name"},
	}

	_, err := test.Manager.CreateContext(context)

	require.NoError(t, err)
	require.Equal(t, &ProjectMetadata{Id: "project-1", Name: "Canonical project"}, context.Project)
}

func TestContextManagerCreateContextRejectsProjectFromAnotherWorkspace(t *testing.T) {
	test := NewEmptyTestContextManager()
	test.Workspaces.Seed(&Workspace{Id: "workspace-1"})
	test.Projects.Seed(&Project{Id: "project-2", WorkspaceId: "workspace-2"})

	_, err := test.Manager.CreateContext(&Context{
		Name:        "Context",
		WorkspaceId: "workspace-1",
		Project:     &ProjectMetadata{Id: "project-2"},
	})

	var mismatchErr *ProjectWorkspaceMismatchError
	require.ErrorAs(t, err, &mismatchErr)
}

func TestContextManagerCreateContextRequiresExistingWorkspace(t *testing.T) {
	test := NewTestContextManager()

	_, err := test.Manager.CreateContext(&Context{Name: "Context", WorkspaceId: "missing"})

	var workspaceNotFoundErr *WorkspaceNotFoundError
	require.ErrorAs(t, err, &workspaceNotFoundErr)
	require.Equal(t, "missing", workspaceNotFoundErr.WorkspaceId)
}

func TestContextManagerUpdateContextPreservesWorkspaceWhenPayloadOmitsIt(t *testing.T) {
	test := NewEmptyTestContextManager()
	test.Contexts.Seed(&Context{Id: "context-1", Name: "Old", WorkspaceId: "workspace-1"})
	updated := &Context{Id: "context-1", Name: "New"}

	err := test.Manager.UpdateContext(updated)

	require.NoError(t, err)
	require.Equal(t, "workspace-1", updated.WorkspaceId)
	require.Equal(t, []*Context{updated}, test.Contexts.saved)
}

func TestContextManagerUpdateContextRejectsWorkspaceMove(t *testing.T) {
	test := NewEmptyTestContextManager()
	test.Contexts.Seed(&Context{Id: "context-1", WorkspaceId: "workspace-1"})

	err := test.Manager.UpdateContext(&Context{Id: "context-1", WorkspaceId: "workspace-2"})

	var moveErr *ContextWorkspaceMoveNotAllowedError
	require.ErrorAs(t, err, &moveErr)
	require.Equal(t, "workspace-1", moveErr.FromWorkspaceId)
	require.Equal(t, "workspace-2", moveErr.ToWorkspaceId)
	require.Empty(t, test.Contexts.saved)
}

func TestContextManagerUpdateContextCanUnassignProject(t *testing.T) {
	existing := &Context{
		Id:          "context-1",
		WorkspaceId: "workspace-1",
		Project:     &ProjectMetadata{Id: "project-1", Name: "Project"},
	}
	test := NewEmptyTestContextManager()
	test.Contexts.Seed(existing)
	updated := &Context{Id: "context-1", Name: "Context"}

	err := test.Manager.UpdateContext(updated)

	require.NoError(t, err)
	require.Nil(t, updated.Project)
	require.Equal(t, "workspace-1", updated.WorkspaceId)
}

func TestContextManagerUpdateProjectRejectsDescendantAsParent(t *testing.T) {
	test := NewEmptyTestContextManager()
	test.Projects.Seed(
		&Project{Id: "project-1", WorkspaceId: "workspace-1"},
		&Project{Id: "project-2", ParentId: "project-1", WorkspaceId: "workspace-1"},
	)

	err := test.Manager.UpdateProject(&Project{
		Id:       "project-1",
		Name:     "Parent",
		ParentId: "project-2",
	})

	var cycleErr *ProjectHierarchyCycleError
	require.ErrorAs(t, err, &cycleErr)
	require.Empty(t, test.Projects.saved)
}

func TestContextManagerUpdateProjectRefreshesContextMetadata(t *testing.T) {
	context := &Context{
		Id:          "context-1",
		WorkspaceId: "workspace-1",
		Project:     &ProjectMetadata{Id: "project-1", Name: "Old name"},
	}
	test := NewEmptyTestContextManager()
	test.Contexts.Seed(context)
	test.Projects.Seed(&Project{Id: "project-1", Name: "Old name", WorkspaceId: "workspace-1"})

	err := test.Manager.UpdateProject(&Project{Id: "project-1", Name: "New name"})

	require.NoError(t, err)
	require.Equal(t, &ProjectMetadata{Id: "project-1", Name: "New name"}, context.Project)
	require.Equal(t, []*Context{context}, test.Contexts.saved)
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
	test := NewEmptyTestContextManager()
	test.Contexts.Seed(context)
	test.Projects.Seed(
		&Project{Id: "project-1", Name: "Parent", WorkspaceId: "workspace-1"},
		&Project{Id: "project-2", Name: "Deleted", ParentId: "project-1", WorkspaceId: "workspace-1"},
		childProject,
	)

	err := test.Manager.DeleteProject("project-2")

	require.NoError(t, err)
	require.Equal(t, "project-1", childProject.ParentId)
	require.Equal(t, &ProjectMetadata{Id: "project-1", Name: "Parent"}, context.Project)
	require.Equal(t, []string{"project-2"}, test.Projects.deletedIDs)
}

func TestContextManagerSaveIntervalUsesContextWorkspace(t *testing.T) {
	test := NewEmptyTestContextManager()
	test.Contexts.Seed(&Context{Id: "context-2", WorkspaceId: "workspace-2"})
	interval := &Interval{
		Id:          "interval-1",
		ContextId:   "context-2",
		WorkspaceId: "workspace-1",
	}

	_, err := test.Manager.SaveInterval(interval)

	require.NoError(t, err)
	require.Equal(t, "workspace-2", interval.WorkspaceId)
	require.Equal(t, []*Interval{interval}, test.Intervals.saved)
}

func TestContextManagerSaveIntervalRejectsMissingContext(t *testing.T) {
	test := NewEmptyTestContextManager()

	_, err := test.Manager.SaveInterval(&Interval{ContextId: "missing"})

	var contextNotFoundErr *ContextNotFoundError
	require.ErrorAs(t, err, &contextNotFoundErr)
	require.Equal(t, "missing", contextNotFoundErr.ContextId)
}

func TestContextManagerDeleteContextDeletesIntervals(t *testing.T) {
	test := NewEmptyTestContextManager()

	err := test.Manager.DeleteContext("context-1")

	require.NoError(t, err)
	require.Equal(t, []string{"context-1"}, test.Intervals.deletedContextIDs)
	require.Equal(t, []string{"context-1"}, test.Contexts.deletedIDs)
}

func TestContextManagerDeleteContextStopsWhenIntervalDeleteFails(t *testing.T) {
	wantErr := errors.New("delete intervals failed")
	test := NewEmptyTestContextManager()
	test.Intervals.deleteByContextError = wantErr

	err := test.Manager.DeleteContext("context-1")

	require.ErrorIs(t, err, wantErr)
	require.Equal(t, []string{"context-1"}, test.Intervals.deletedContextIDs)
	require.Empty(t, test.Contexts.deletedIDs)
}

func TestContextManagerDeleteWorkspaceDeletesUnusedWorkspace(t *testing.T) {
	test := NewEmptyTestContextManager()

	err := test.Manager.DeleteWorkspace("workspace-1")

	require.NoError(t, err)
	require.Equal(t, 1, test.Contexts.listByWorkspaceCalls)
	require.Equal(t, "workspace-1", test.Contexts.listedWorkspaceID)
	require.Equal(t, []string{"workspace-1"}, test.Workspaces.deletedIDs)
}

func TestContextManagerDeleteWorkspaceReturnsWorkspaceInUseError(t *testing.T) {
	test := NewEmptyTestContextManager()
	test.Contexts.Seed(&Context{Id: "context-1", WorkspaceId: "workspace-1"})

	err := test.Manager.DeleteWorkspace("workspace-1")

	var workspaceInUseErr *WorkspaceInUseError
	require.ErrorAs(t, err, &workspaceInUseErr)
	require.Equal(t, "workspace-1", workspaceInUseErr.WorkspaceId)
	require.Empty(t, test.Workspaces.deletedIDs)
}

func TestContextManagerDeleteWorkspaceReturnsContextRepositoryError(t *testing.T) {
	wantErr := errors.New("list contexts failed")
	test := NewEmptyTestContextManager()
	test.Contexts.listByWorkspaceError = wantErr

	err := test.Manager.DeleteWorkspace("workspace-1")

	require.ErrorIs(t, err, wantErr)
	require.Empty(t, test.Workspaces.deletedIDs)
}

func TestContextManagerDeleteWorkspaceReturnsWorkspaceRepositoryError(t *testing.T) {
	wantErr := errors.New("delete workspace failed")
	test := NewEmptyTestContextManager()
	test.Workspaces.deleteError = wantErr

	err := test.Manager.DeleteWorkspace("workspace-1")

	require.ErrorIs(t, err, wantErr)
	require.Equal(t, []string{"workspace-1"}, test.Workspaces.deletedIDs)
}

func TestContextManagerGetWorkspaceStatsUsesAllIntervals(t *testing.T) {
	test := NewEmptyTestContextManager()
	now := time.Date(2026, time.June, 14, 12, 0, 0, 0, time.UTC)
	test.TimeProvider.Set(now)
	test.Contexts.Seed(
		&Context{Id: "context-1", Name: "First", WorkspaceId: "workspace-1"},
		&Context{Id: "context-2", Name: "Second", WorkspaceId: "workspace-1"},
	)
	test.Intervals.Seed(
		&Interval{
			ContextId: "context-1",
			Start:     testTime(now.Add(-3 * time.Hour)),
			End:       testTime(now.Add(-2 * time.Hour)),
			Status:    "completed",
		},
		&Interval{
			ContextId: "context-1",
			Start:     testTime(now.Add(-30 * time.Minute)),
			Status:    "active",
		},
		&Interval{ContextId: "context-2", Duration: 30 * time.Minute, Status: "completed"},
	)

	stats, err := test.Manager.GetWorkspaceStats("workspace-1")

	require.NoError(t, err)
	require.Equal(t, 2*time.Hour, stats.TotalDuration)
	require.Equal(t, 3, stats.TotalSessions)
	require.Len(t, stats.Contexts, 2)
	require.Len(t, stats.ContextStats, 2)
	require.Equal(t, "context-1", stats.ContextStats[0].ContextId)
	require.Equal(t, 90*time.Minute, stats.ContextStats[0].Duration)
	require.InDelta(t, 75, stats.ContextStats[0].Percentage, 0.001)
}

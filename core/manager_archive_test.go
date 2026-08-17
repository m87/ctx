package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestListArchiveCandidatesUsesLatestIntervalAndWorkspace(t *testing.T) {
	now := time.Date(2026, time.August, 18, 12, 0, 0, 0, time.UTC)
	tm := newTestManager()
	tm.Manager.TimeProvider = fixedTimeProvider{now: now}
	workspaceID, err := tm.Workspaces.Save(&Workspace{Name: "archive-candidates"})
	require.NoError(t, err)
	otherWorkspaceID, err := tm.Workspaces.Save(&Workspace{Name: "other-workspace"})
	require.NoError(t, err)

	staleID := saveArchiveTestContext(t, tm, workspaceID, "stale", "inactive", false)
	freshID := saveArchiveTestContext(t, tm, workspaceID, "fresh", "inactive", false)
	boundaryID := saveArchiveTestContext(t, tm, workspaceID, "boundary", "inactive", false)
	archivedID := saveArchiveTestContext(t, tm, workspaceID, "already archived", "archived", true)
	activeID := saveArchiveTestContext(t, tm, workspaceID, "active", "active", false)
	incompleteIntervalID := saveArchiveTestContext(t, tm, workspaceID, "incomplete interval", "inactive", false)
	saveArchiveTestContext(t, tm, workspaceID, "without intervals", "inactive", false)
	otherWorkspaceContextID := saveArchiveTestContext(t, tm, otherWorkspaceID, "other", "inactive", false)

	saveArchiveTestInterval(t, tm, staleID, now.AddDate(0, 0, -100), now.AddDate(0, 0, -90))
	latestStaleEnd := now.AddDate(0, 0, -31)
	saveArchiveTestInterval(t, tm, staleID, now.AddDate(0, 0, -32), latestStaleEnd)
	saveArchiveTestInterval(t, tm, freshID, now.AddDate(0, 0, -11), now.AddDate(0, 0, -10))
	saveArchiveTestInterval(t, tm, boundaryID, now.AddDate(0, 0, -31), now.AddDate(0, 0, -30))
	saveArchiveTestInterval(t, tm, archivedID, now.AddDate(0, 0, -60), now.AddDate(0, 0, -59))
	saveArchiveTestInterval(t, tm, activeID, now.AddDate(0, 0, -60), now.AddDate(0, 0, -59))
	incompleteStart := now.AddDate(0, 0, -60)
	_, err = tm.Intervals.Save(&Interval{
		ContextId: incompleteIntervalID,
		Start:     &incompleteStart,
		Status:    "active",
	})
	require.NoError(t, err)
	saveArchiveTestInterval(t, tm, otherWorkspaceContextID, now.AddDate(0, 0, -60), now.AddDate(0, 0, -59))

	preview, err := tm.Manager.ListArchiveCandidates(workspaceID, 30, time.UTC)

	require.NoError(t, err)
	require.Equal(t, time.Date(2026, time.July, 19, 0, 0, 0, 0, time.UTC), preview.Cutoff)
	require.Len(t, preview.Contexts, 1)
	require.Equal(t, staleID, preview.Contexts[0].Id)
	require.Equal(t, latestStaleEnd, preview.Contexts[0].LastIntervalAt)
}

func TestArchiveStaleContextsArchivesPreviewedContexts(t *testing.T) {
	now := time.Date(2026, time.August, 18, 12, 0, 0, 0, time.UTC)
	tm := newTestManager()
	tm.Manager.TimeProvider = fixedTimeProvider{now: now}
	workspaceID, err := tm.Workspaces.Save(&Workspace{Name: "bulk-archive"})
	require.NoError(t, err)
	staleID := saveArchiveTestContext(t, tm, workspaceID, "stale", "inactive", false)
	freshID := saveArchiveTestContext(t, tm, workspaceID, "fresh", "inactive", false)
	saveArchiveTestInterval(t, tm, staleID, now.AddDate(0, 0, -50), now.AddDate(0, 0, -45))
	saveArchiveTestInterval(t, tm, freshID, now.AddDate(0, 0, -5), now.AddDate(0, 0, -4))

	result, err := tm.Manager.ArchiveStaleContexts(workspaceID, 30, time.UTC)

	require.NoError(t, err)
	require.Equal(t, 1, result.ArchivedCount)
	require.Len(t, result.Contexts, 1)
	require.Equal(t, staleID, result.Contexts[0].Id)
	require.True(t, tm.Contexts.items[staleID].Archived)
	require.Equal(t, "archived", tm.Contexts.items[staleID].Status)
	require.False(t, tm.Contexts.items[freshID].Archived)
}

func TestArchiveStaleContextsReevaluatesCandidatesInTransaction(t *testing.T) {
	now := time.Date(2026, time.August, 18, 12, 0, 0, 0, time.UTC)
	tm := newTestManager()
	tm.Manager.TimeProvider = fixedTimeProvider{now: now}
	workspaceID, err := tm.Workspaces.Save(&Workspace{Name: "bulk-archive-race"})
	require.NoError(t, err)
	contextID := saveArchiveTestContext(t, tm, workspaceID, "recently used", "inactive", false)
	saveArchiveTestInterval(t, tm, contextID, now.AddDate(0, 0, -50), now.AddDate(0, 0, -45))

	preview, err := tm.Manager.ListArchiveCandidates(workspaceID, 30, time.UTC)
	require.NoError(t, err)
	require.Len(t, preview.Contexts, 1)

	tm.Manager.RunInTransaction = func(fn func(*ContextManager) error) error {
		saveArchiveTestInterval(t, tm, contextID, now.Add(-2*time.Hour), now.Add(-time.Hour))
		return fn(tm.Manager)
	}
	result, err := tm.Manager.ArchiveStaleContexts(workspaceID, 30, time.UTC)

	require.NoError(t, err)
	require.Zero(t, result.ArchivedCount)
	require.False(t, tm.Contexts.items[contextID].Archived)
}

func TestArchiveThresholdMustBePositive(t *testing.T) {
	tm := newTestManager()

	_, err := tm.Manager.ListArchiveCandidates("workspace", 0, time.UTC)

	require.ErrorAs(t, err, new(*InvalidArchiveThresholdError))
}

func TestArchiveCutoffUsesStartOfDayInRequestedTimeZone(t *testing.T) {
	now := time.Date(2026, time.August, 18, 22, 30, 0, 0, time.UTC)
	tm := newTestManager()
	tm.Manager.TimeProvider = fixedTimeProvider{now: now}
	location, err := time.LoadLocation("Europe/Warsaw")
	require.NoError(t, err)

	cutoff, err := tm.Manager.archiveCutoff(30, location)

	require.NoError(t, err)
	require.Equal(t, "2026-07-19T22:00:00Z", cutoff.Format(time.RFC3339))
}

func saveArchiveTestContext(t *testing.T, tm *testManager, workspaceID, name, status string, archived bool) string {
	t.Helper()
	id, err := tm.Contexts.Save(&Context{
		Name:        name,
		WorkspaceId: workspaceID,
		Status:      status,
		Archived:    archived,
	})
	require.NoError(t, err)
	return id
}

func saveArchiveTestInterval(t *testing.T, tm *testManager, contextID string, start, end time.Time) {
	t.Helper()
	_, err := tm.Intervals.Save(&Interval{
		ContextId: contextID,
		Start:     &start,
		End:       &end,
		Status:    "completed",
	})
	require.NoError(t, err)
}

func TestArchiveContext(t *testing.T) {
	tm := newTestManager()
	workspaceID, err := tm.Workspaces.Save(&Workspace{Name: "archive-test"})
	require.NoError(t, err)
	contextID, err := tm.Manager.CreateContext(&Context{
		Name:        "context-to-archive",
		Description: "keep description",
		WorkspaceId: workspaceID,
		Tags:        []string{"keep-tag"},
	})
	require.NoError(t, err)

	err = tm.Manager.ArchiveContext(contextID)
	require.NoError(t, err)

	context, err := tm.Contexts.GetById(contextID)
	require.NoError(t, err)
	require.True(t, context.Archived)
	require.Equal(t, "archived", context.Status)
	require.Equal(t, "keep description", context.Description)
	require.Equal(t, []string{"keep-tag"}, context.Tags)
}

func TestRestoreContext(t *testing.T) {
	tm := newTestManager()
	workspaceID, err := tm.Workspaces.Save(&Workspace{Name: "restore-test"})
	require.NoError(t, err)
	contextID, err := tm.Manager.CreateContext(&Context{
		Name:        "context-to-restore",
		Description: "keep description",
		WorkspaceId: workspaceID,
		Tags:        []string{"keep-tag"},
	})
	require.NoError(t, err)

	err = tm.Manager.ArchiveContext(contextID)
	require.NoError(t, err)

	err = tm.Manager.RestoreContext(contextID)
	require.NoError(t, err)

	context, err := tm.Contexts.GetById(contextID)
	require.NoError(t, err)
	require.False(t, context.Archived)
	require.Equal(t, "inactive", context.Status)
	require.Equal(t, "keep description", context.Description)
	require.Equal(t, []string{"keep-tag"}, context.Tags)
}

func TestArchiveActiveContext(t *testing.T) {
	tm := newTestManager()
	workspaceID, err := tm.Workspaces.Save(&Workspace{Name: "archive-active-test"})
	require.NoError(t, err)
	contextID, err := tm.Manager.CreateContext(&Context{
		Name:        "active-context",
		Description: "keep description",
		WorkspaceId: workspaceID,
		Tags:        []string{"keep-tag"},
	})
	require.NoError(t, err)

	err = tm.Manager.SwitchContext(tm.Contexts.items[contextID])
	require.NoError(t, err)

	err = tm.Manager.ArchiveContext(contextID)
	require.Error(t, err)
	require.IsType(t, &ArchiveContextActiveError{}, err)
}

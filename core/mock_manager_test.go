package core

import (
	"time"
)

var DefaultTestTime = time.Date(2026, time.January, 2, 12, 0, 0, 0, time.UTC)

const (
	TestWorkspaceID = "workspace-1"
	TestProjectID   = "project-1"
	TestContextID   = "context-1"
	TestIntervalID  = "interval-1"
)

type TestTimeProvider struct {
	now time.Time
}

func NewTestTimeProvider() *TestTimeProvider {
	return &TestTimeProvider{now: DefaultTestTime}
}

func (p *TestTimeProvider) Now() time.Time {
	return p.now.UTC()
}

func (p *TestTimeProvider) Set(now time.Time) {
	p.now = now.UTC()
}

func (p *TestTimeProvider) Advance(duration time.Duration) {
	p.now = p.now.Add(duration)
}

type TestContextManager struct {
	Manager      *ContextManager
	TimeProvider *TestTimeProvider
	Contexts     *ContextRepositoryMock
	Intervals    *IntervalRepositoryMock
	Projects     *ProjectRepositoryMock
	SavedQueries *SavedQueryRepositoryMock
	Workspaces   *WorkspaceRepositoryMock
}

func NewTestContextManager() *TestContextManager {
	test := NewEmptyTestContextManager()
	intervalStart := DefaultTestTime.Add(-time.Hour)
	intervalEnd := DefaultTestTime

	test.Workspaces.Seed(&Workspace{Id: TestWorkspaceID, Name: "Test workspace"})
	test.Projects.Seed(&Project{
		Id:          TestProjectID,
		Name:        "Test project",
		WorkspaceId: TestWorkspaceID,
	})
	test.Contexts.Seed(&Context{
		Id:          TestContextID,
		Name:        "Test context",
		WorkspaceId: TestWorkspaceID,
		Status:      "inactive",
		Project:     &ProjectMetadata{Id: TestProjectID, Name: "Test project"},
	})
	test.Intervals.Seed(&Interval{
		Id:          TestIntervalID,
		ContextId:   TestContextID,
		WorkspaceId: TestWorkspaceID,
		Start:       &intervalStart,
		End:         &intervalEnd,
		Duration:    time.Hour,
		Status:      "completed",
	})

	return test
}

func NewEmptyTestContextManager() *TestContextManager {
	timeProvider := NewTestTimeProvider()
	contexts := NewContextRepositoryMock()
	intervals := NewIntervalRepositoryMock()
	projects := NewProjectRepositoryMock()
	savedQueries := NewSavedQueryRepositoryMock()
	workspaces := NewWorkspaceRepositoryMock()
	manager := NewContextManager(timeProvider, contexts, intervals, workspaces, projects)
	manager.SavedQueryRepository = savedQueries

	return &TestContextManager{
		Manager:      manager,
		TimeProvider: timeProvider,
		Contexts:     contexts,
		Intervals:    intervals,
		Projects:     projects,
		SavedQueries: savedQueries,
		Workspaces:   workspaces,
	}
}

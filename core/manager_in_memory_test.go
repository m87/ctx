package core

import (
	"fmt"
	"sort"
	"time"
)

// testManager provides a ContextManager backed only by local, in-memory data.
// Tests can inspect the repositories directly when they need to assert persisted state.
type testManager struct {
	Manager    *ContextManager
	Contexts   *memoryContextRepository
	Intervals  *memoryIntervalRepository
	Workspaces *memoryWorkspaceRepository
}

func newTestManager() *testManager {
	contexts := &memoryContextRepository{items: make(map[string]*Context)}
	intervals := &memoryIntervalRepository{items: make(map[string]*Interval)}
	workspaces := &memoryWorkspaceRepository{items: make(map[string]*Workspace)}

	return &testManager{
		Manager:    NewContextManager(fixedTimeProvider{now: time.Now().UTC()}, contexts, intervals, workspaces),
		Contexts:   contexts,
		Intervals:  intervals,
		Workspaces: workspaces,
	}
}

type memoryContextRepository struct {
	items  map[string]*Context
	nextID int
}

func (r *memoryContextRepository) GetById(id string) (*Context, error) { return r.items[id], nil }

func (r *memoryContextRepository) Save(context *Context) (string, error) {
	if context == nil {
		return "", fmt.Errorf("context is required")
	}
	if context.Id == "" {
		r.nextID++
		context.Id = fmt.Sprintf("context-%d", r.nextID)
	}
	r.items[context.Id] = context
	return context.Id, nil
}

func (r *memoryContextRepository) Delete(id string) error {
	delete(r.items, id)
	return nil
}

func (r *memoryContextRepository) List() ([]*Context, error) {
	result := make([]*Context, 0, len(r.items))
	for _, context := range r.items {
		result = append(result, context)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Id < result[j].Id })
	return result, nil
}

func (r *memoryContextRepository) ListByWorkspace(workspaceID string) ([]*Context, error) {
	return r.listByWorkspace(workspaceID, false), nil
}

func (r *memoryContextRepository) ListByWorkspaceIncludingArchived(workspaceID string) ([]*Context, error) {
	return r.listByWorkspace(workspaceID, true), nil
}

func (r *memoryContextRepository) listByWorkspace(workspaceID string, includeArchived bool) []*Context {
	result := make([]*Context, 0)
	for _, context := range r.items {
		if context.WorkspaceId == workspaceID && (includeArchived || !context.Archived) {
			result = append(result, context)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Id < result[j].Id })
	return result
}

func (r *memoryContextRepository) GetActive() (*Context, error) {
	for _, context := range r.items {
		if context.Status == "active" {
			return context, nil
		}
	}
	return nil, nil
}

type memoryWorkspaceRepository struct {
	items  map[string]*Workspace
	nextID int
}

func (r *memoryWorkspaceRepository) GetById(id string) (*Workspace, error) { return r.items[id], nil }

func (r *memoryWorkspaceRepository) Save(workspace *Workspace) (string, error) {
	if workspace == nil {
		return "", fmt.Errorf("workspace is required")
	}
	if workspace.Id == "" {
		r.nextID++
		workspace.Id = fmt.Sprintf("workspace-%d", r.nextID)
	}
	r.items[workspace.Id] = workspace
	return workspace.Id, nil
}

func (r *memoryWorkspaceRepository) Delete(id string) error {
	delete(r.items, id)
	return nil
}

func (r *memoryWorkspaceRepository) List() ([]*Workspace, error) {
	result := make([]*Workspace, 0, len(r.items))
	for _, workspace := range r.items {
		result = append(result, workspace)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Id < result[j].Id })
	return result, nil
}

type memoryIntervalRepository struct {
	items  map[string]*Interval
	nextID int
}

func (r *memoryIntervalRepository) GetById(id string) (*Interval, error) { return r.items[id], nil }

func (r *memoryIntervalRepository) Save(interval *Interval) (string, error) {
	if interval == nil {
		return "", fmt.Errorf("interval is required")
	}
	if interval.Id == "" {
		r.nextID++
		interval.Id = fmt.Sprintf("interval-%d", r.nextID)
	}
	r.items[interval.Id] = interval
	return interval.Id, nil
}

func (r *memoryIntervalRepository) Delete(id string) error {
	delete(r.items, id)
	return nil
}

func (r *memoryIntervalRepository) DeleteByContextId(contextID string) error {
	for id, interval := range r.items {
		if interval.ContextId == contextID {
			delete(r.items, id)
		}
	}
	return nil
}

func (r *memoryIntervalRepository) ListByContextId(contextID string) ([]*Interval, error) {
	result := make([]*Interval, 0)
	for _, interval := range r.items {
		if interval.ContextId == contextID {
			result = append(result, interval)
		}
	}
	return result, nil
}

func (r *memoryIntervalRepository) GetActiveIntervalByContextId(contextID string) (*Interval, error) {
	for _, interval := range r.items {
		if interval.ContextId == contextID && interval.Status == "active" {
			return interval, nil
		}
	}
	return nil, nil
}

func (r *memoryIntervalRepository) ListByDay(date time.Time, workspaceID string) ([]*Interval, error) {
	result := make([]*Interval, 0)
	for _, interval := range r.items {
		if interval.WorkspaceId == workspaceID && interval.Start.Time.Year() == date.Year() && interval.Start.Time.YearDay() == date.YearDay() {
			result = append(result, interval)
		}
	}
	return result, nil
}

func (r *memoryIntervalRepository) List() ([]*Interval, error) {
	result := make([]*Interval, 0, len(r.items))
	for _, interval := range r.items {
		result = append(result, interval)
	}
	return result, nil
}

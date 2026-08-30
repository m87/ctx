package core

import (
	"fmt"
	"sort"
	"time"
)

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
	projects := &memoryProjectRepository{items: make(map[string]*Project)}

	return &testManager{
		Manager:    NewContextManager(fixedTimeProvider{now: time.Now().UTC()}, contexts, intervals, workspaces, projects),
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

func (r *memoryContextRepository) SaveAll(contexts []*Context) ([]string, error) {
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

func (r *memoryContextRepository) ListToSync(limit int) ([]*Context, error) {
	contexts, err := r.List()
	if err != nil {
		return nil, err
	}

	result := make([]*Context, 0, len(contexts))
	for _, context := range contexts {
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

func (r *memoryContextRepository) ListByWorkspace(workspaceID string) ([]*Context, error) {
	return r.listByWorkspace(workspaceID, false), nil
}

func (r *memoryContextRepository) ListByProject(projectID string) ([]*Context, error) {
	result := make([]*Context, 0)
	for _, context := range r.items {
		if context.Project != nil && context.Project.Id == projectID {
			result = append(result, context)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Id < result[j].Id })
	return result, nil
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

type memoryProjectRepository struct {
	items  map[string]*Project
	nextID int
}

func (r *memoryProjectRepository) GetById(id string) (*Project, error) { return r.items[id], nil }

func (r *memoryProjectRepository) Save(project *Project) (string, error) {
	if project == nil {
		return "", fmt.Errorf("project is required")
	}
	if project.Id == "" {
		r.nextID++
		project.Id = fmt.Sprintf("project-%d", r.nextID)
	}
	r.items[project.Id] = project
	return project.Id, nil
}

func (r *memoryProjectRepository) SaveAll(projects []*Project) ([]string, error) {
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

func (r *memoryProjectRepository) Delete(id string) error {
	delete(r.items, id)
	return nil
}

func (r *memoryProjectRepository) List(workspaceID string) ([]*Project, error) {
	result := make([]*Project, 0)
	for _, project := range r.items {
		result = append(result, project)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Id < result[j].Id })
	return result, nil
}

func (r *memoryProjectRepository) ListIncludingArchived(workspaceID string) ([]*Project, error) {
	result := make([]*Project, 0)
	for _, project := range r.items {
		if project.WorkspaceId == workspaceID {
			result = append(result, project)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Id < result[j].Id })
	return result, nil
}

func (r *memoryProjectRepository) ListChildren(parentID string) ([]*Project, error) {
	result := make([]*Project, 0)
	for _, project := range r.items {
		if project.ParentId == parentID {
			result = append(result, project)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Id < result[j].Id })
	return result, nil
}

func (r *memoryProjectRepository) ListToSync(limit int) ([]*Project, error) {
	projects := make([]*Project, 0, len(r.items))
	for _, project := range r.items {
		projects = append(projects, project)
		if limit > 0 && len(projects) == limit {
			break
		}
	}
	return projects, nil
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

func (r *memoryWorkspaceRepository) SaveAll(workspaces []*Workspace) ([]string, error) {
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

func (r *memoryWorkspaceRepository) ListToSync(limit int) ([]*Workspace, error) {
	workspaces, err := r.List()
	if err != nil {
		return nil, err
	}

	result := make([]*Workspace, 0, len(workspaces))
	for _, workspace := range workspaces {
		if workspace == nil {
			continue
		}
		result = append(result, workspace)
		if limit > 0 && len(result) == limit {
			break
		}
	}
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

func (r *memoryIntervalRepository) SaveAll(intervals []*Interval) ([]string, error) {
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
		if interval.WorkspaceId == workspaceID && interval.Start != nil && interval.Start.Year() == date.Year() && interval.Start.YearDay() == date.YearDay() {
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

func (r *memoryIntervalRepository) ListToSync(limit int) ([]*Interval, error) {
	intervals, err := r.List()
	if err != nil {
		return nil, err
	}

	result := make([]*Interval, 0, len(intervals))
	for _, interval := range intervals {
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

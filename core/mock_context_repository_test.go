package core

import (
	"errors"
	"fmt"
	"sort"
)

var _ ContextRepository = (*ContextRepositoryMock)(nil)

type ContextRepositoryMock struct {
	contexts             []*Context
	nextID               int
	listCalled           bool
	lastQuery            *ContextSQLQuery
	saved                []*Context
	deletedIDs           []string
	listError            error
	saveError            error
	deleteError          error
	listByWorkspaceError error
	listByWorkspaceCalls int
	listedWorkspaceID    string
}

func NewContextRepositoryMock(contexts ...*Context) *ContextRepositoryMock {
	mock := &ContextRepositoryMock{}
	mock.Seed(contexts...)
	return mock
}

func (m *ContextRepositoryMock) Seed(contexts ...*Context) {
	m.contexts = copyPointerSlice(contexts)
	m.nextID = 0
	m.listCalled = false
	m.lastQuery = nil
	m.saved = nil
	m.deletedIDs = nil
	m.listByWorkspaceCalls = 0
	m.listedWorkspaceID = ""
}

func (m *ContextRepositoryMock) Get(id string) *Context {
	context, _ := m.GetById(id)
	return context
}

func (m *ContextRepositoryMock) GetById(id string) (*Context, error) {
	for _, context := range m.contexts {
		if context != nil && context.Id == id {
			return context, nil
		}
	}
	return nil, nil
}

func (m *ContextRepositoryMock) Save(context *Context) (string, error) {
	if m.saveError != nil {
		return "", m.saveError
	}
	if context == nil {
		return "", errors.New("context is required")
	}
	if context.Id == "" {
		context.Id = m.nextContextID()
	}
	for i, existing := range m.contexts {
		if existing != nil && existing.Id == context.Id {
			m.contexts[i] = context
			m.saved = append(m.saved, context)
			return context.Id, nil
		}
	}
	m.contexts = append(m.contexts, context)
	m.saved = append(m.saved, context)
	return context.Id, nil
}

func (m *ContextRepositoryMock) SaveAll(contexts []*Context) ([]string, error) {
	ids := make([]string, 0, len(contexts))
	for _, context := range contexts {
		id, err := m.Save(context)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (m *ContextRepositoryMock) Delete(id string) error {
	m.deletedIDs = append(m.deletedIDs, id)
	if m.deleteError != nil {
		return m.deleteError
	}
	for i, context := range m.contexts {
		if context != nil && context.Id == id {
			m.contexts = append(m.contexts[:i], m.contexts[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *ContextRepositoryMock) List() ([]*Context, error) {
	m.listCalled = true
	if m.listError != nil {
		return nil, m.listError
	}
	return copyPointerSlice(m.contexts), nil
}

func (m *ContextRepositoryMock) Query(query *ContextSQLQuery) ([]*Context, error) {
	m.lastQuery = query
	if query == nil || query.WorkspaceId == "" {
		return copyPointerSlice(m.contexts), nil
	}
	return m.listByWorkspace(query.WorkspaceId, true), nil
}

func (m *ContextRepositoryMock) ListByWorkspace(workspaceID string) ([]*Context, error) {
	m.listByWorkspaceCalls++
	m.listedWorkspaceID = workspaceID
	if m.listByWorkspaceError != nil {
		return nil, m.listByWorkspaceError
	}
	return m.listByWorkspace(workspaceID, false), nil
}

func (m *ContextRepositoryMock) ListByProject(projectID string) ([]*Context, error) {
	contexts := make([]*Context, 0)
	for _, context := range m.contexts {
		if context == nil {
			continue
		}
		if context.Project != nil && context.Project.Id == projectID {
			contexts = append(contexts, context)
			continue
		}
		if context.ProjectId != nil && *context.ProjectId == projectID {
			contexts = append(contexts, context)
		}
	}
	sort.Slice(contexts, func(i, j int) bool { return contexts[i].Id < contexts[j].Id })
	return contexts, nil
}

func (m *ContextRepositoryMock) ListByWorkspaceIncludingArchived(workspaceID string) ([]*Context, error) {
	m.listByWorkspaceCalls++
	m.listedWorkspaceID = workspaceID
	if m.listByWorkspaceError != nil {
		return nil, m.listByWorkspaceError
	}
	return m.listByWorkspace(workspaceID, true), nil
}

func (m *ContextRepositoryMock) GetActive() (*Context, error) {
	for _, context := range m.contexts {
		if context != nil && context.Status == "active" {
			return context, nil
		}
	}
	return nil, nil
}

func (m *ContextRepositoryMock) ListToSync(limit int) ([]*Context, error) {
	return limited(copyPointerSlice(m.contexts), limit), nil
}

func (m *ContextRepositoryMock) listByWorkspace(workspaceID string, includeArchived bool) []*Context {
	contexts := make([]*Context, 0)
	for _, context := range m.contexts {
		if context != nil && context.WorkspaceId == workspaceID && (includeArchived || !context.Archived) {
			contexts = append(contexts, context)
		}
	}
	sort.Slice(contexts, func(i, j int) bool { return contexts[i].Id < contexts[j].Id })
	return contexts
}

func (m *ContextRepositoryMock) nextContextID() string {
	for {
		m.nextID++
		id := fmt.Sprintf("context-%d", m.nextID)
		context, _ := m.GetById(id)
		if context == nil {
			return id
		}
	}
}

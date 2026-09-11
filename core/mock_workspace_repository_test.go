package core

import (
	"errors"
	"fmt"
)

var _ WorkspaceRepository = (*WorkspaceRepositoryMock)(nil)

type WorkspaceRepositoryMock struct {
	workspaces  []*Workspace
	nextID      int
	listCalled  bool
	saved       []*Workspace
	deletedIDs  []string
	listError   error
	saveError   error
	deleteError error
}

func NewWorkspaceRepositoryMock(workspaces ...*Workspace) *WorkspaceRepositoryMock {
	mock := &WorkspaceRepositoryMock{}
	mock.Seed(workspaces...)
	return mock
}

func (m *WorkspaceRepositoryMock) Seed(workspaces ...*Workspace) {
	m.workspaces = copyPointerSlice(workspaces)
	m.nextID = 0
	m.listCalled = false
	m.saved = nil
	m.deletedIDs = nil
}

func (m *WorkspaceRepositoryMock) Get(id string) *Workspace {
	workspace, _ := m.GetById(id)
	return workspace
}

func (m *WorkspaceRepositoryMock) GetById(id string) (*Workspace, error) {
	for _, workspace := range m.workspaces {
		if workspace != nil && workspace.Id == id {
			return workspace, nil
		}
	}
	return nil, nil
}

func (m *WorkspaceRepositoryMock) Save(workspace *Workspace) (string, error) {
	if m.saveError != nil {
		return "", m.saveError
	}
	if workspace == nil {
		return "", errors.New("workspace is required")
	}
	if workspace.Id == "" {
		workspace.Id = m.nextWorkspaceID()
	}
	for i, existing := range m.workspaces {
		if existing != nil && existing.Id == workspace.Id {
			m.workspaces[i] = workspace
			m.saved = append(m.saved, workspace)
			return workspace.Id, nil
		}
	}
	m.workspaces = append(m.workspaces, workspace)
	m.saved = append(m.saved, workspace)
	return workspace.Id, nil
}

func (m *WorkspaceRepositoryMock) SaveAll(workspaces []*Workspace) ([]string, error) {
	ids := make([]string, 0, len(workspaces))
	for _, workspace := range workspaces {
		id, err := m.Save(workspace)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (m *WorkspaceRepositoryMock) Delete(id string) error {
	m.deletedIDs = append(m.deletedIDs, id)
	if m.deleteError != nil {
		return m.deleteError
	}
	for i, workspace := range m.workspaces {
		if workspace != nil && workspace.Id == id {
			m.workspaces = append(m.workspaces[:i], m.workspaces[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *WorkspaceRepositoryMock) List() ([]*Workspace, error) {
	m.listCalled = true
	if m.listError != nil {
		return nil, m.listError
	}
	return copyPointerSlice(m.workspaces), nil
}

func (m *WorkspaceRepositoryMock) ListToSync(limit int) ([]*Workspace, error) {
	return limited(copyPointerSlice(m.workspaces), limit), nil
}

func (m *WorkspaceRepositoryMock) nextWorkspaceID() string {
	for {
		m.nextID++
		id := fmt.Sprintf("workspace-%d", m.nextID)
		workspace, _ := m.GetById(id)
		if workspace == nil {
			return id
		}
	}
}

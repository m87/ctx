package core

import (
	"errors"
	"fmt"
	"sort"
)

var _ ProjectRepository = (*ProjectRepositoryMock)(nil)

type ProjectRepositoryMock struct {
	projects    []*Project
	nextID      int
	listCalled  bool
	saved       []*Project
	deletedIDs  []string
	listError   error
	saveError   error
	deleteError error
}

func NewProjectRepositoryMock(projects ...*Project) *ProjectRepositoryMock {
	mock := &ProjectRepositoryMock{}
	mock.Seed(projects...)
	return mock
}

func (m *ProjectRepositoryMock) Seed(projects ...*Project) {
	m.projects = copyPointerSlice(projects)
	m.nextID = 0
	m.listCalled = false
	m.saved = nil
	m.deletedIDs = nil
}

func (m *ProjectRepositoryMock) Get(id string) *Project {
	project, _ := m.GetById(id)
	return project
}

func (m *ProjectRepositoryMock) GetById(id string) (*Project, error) {
	for _, project := range m.projects {
		if project != nil && project.Id == id {
			return project, nil
		}
	}
	return nil, nil
}

func (m *ProjectRepositoryMock) Save(project *Project) (string, error) {
	if m.saveError != nil {
		return "", m.saveError
	}
	if project == nil {
		return "", errors.New("project is required")
	}
	if project.Id == "" {
		project.Id = m.nextProjectID()
	}
	for i, existing := range m.projects {
		if existing != nil && existing.Id == project.Id {
			m.projects[i] = project
			m.saved = append(m.saved, project)
			return project.Id, nil
		}
	}
	m.projects = append(m.projects, project)
	m.saved = append(m.saved, project)
	return project.Id, nil
}

func (m *ProjectRepositoryMock) SaveAll(projects []*Project) ([]string, error) {
	ids := make([]string, 0, len(projects))
	for _, project := range projects {
		id, err := m.Save(project)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (m *ProjectRepositoryMock) List(workspaceID string) ([]*Project, error) {
	m.listCalled = true
	if m.listError != nil {
		return nil, m.listError
	}
	projects := make([]*Project, 0)
	for _, project := range m.projects {
		if project != nil && project.WorkspaceId == workspaceID && project.ParentId == "" {
			projects = append(projects, project)
		}
	}
	sort.Slice(projects, func(i, j int) bool { return projects[i].Id < projects[j].Id })
	return projects, nil
}

func (m *ProjectRepositoryMock) ListToSync(limit int) ([]*Project, error) {
	if m.listError != nil {
		return nil, m.listError
	}
	return limited(copyPointerSlice(m.projects), limit), nil
}

func (m *ProjectRepositoryMock) ListIncludingArchived(workspaceID string) ([]*Project, error) {
	m.listCalled = true
	if m.listError != nil {
		return nil, m.listError
	}
	projects := make([]*Project, 0)
	for _, project := range m.projects {
		if project != nil && project.WorkspaceId == workspaceID {
			projects = append(projects, project)
		}
	}
	sort.Slice(projects, func(i, j int) bool { return projects[i].Id < projects[j].Id })
	return projects, nil
}

func (m *ProjectRepositoryMock) ListChildren(parentID string) ([]*Project, error) {
	if parentID == "" {
		return nil, nil
	}
	projects := make([]*Project, 0)
	for _, project := range m.projects {
		if project != nil && project.ParentId == parentID {
			projects = append(projects, project)
		}
	}
	sort.Slice(projects, func(i, j int) bool { return projects[i].Id < projects[j].Id })
	return projects, nil
}

func (m *ProjectRepositoryMock) Delete(id string) error {
	m.deletedIDs = append(m.deletedIDs, id)
	if m.deleteError != nil {
		return m.deleteError
	}
	for i, project := range m.projects {
		if project != nil && project.Id == id {
			m.projects = append(m.projects[:i], m.projects[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *ProjectRepositoryMock) nextProjectID() string {
	for {
		m.nextID++
		id := fmt.Sprintf("project-%d", m.nextID)
		project, _ := m.GetById(id)
		if project == nil {
			return id
		}
	}
}

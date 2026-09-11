package core

import (
	"errors"
	"fmt"
	"sort"
)

var _ SavedQueryRepository = (*SavedQueryRepositoryMock)(nil)

type SavedQueryRepositoryMock struct {
	queries    []*SavedQuery
	nextID     int
	saved      []*SavedQuery
	deletedIDs []string
}

func NewSavedQueryRepositoryMock(queries ...*SavedQuery) *SavedQueryRepositoryMock {
	mock := &SavedQueryRepositoryMock{}
	mock.Seed(queries...)
	return mock
}

func (m *SavedQueryRepositoryMock) Seed(queries ...*SavedQuery) {
	m.queries = copyPointerSlice(queries)
	m.nextID = 0
	m.saved = nil
	m.deletedIDs = nil
}

func (m *SavedQueryRepositoryMock) Get(id string) *SavedQuery {
	query, _ := m.GetById(id)
	return query
}

func (m *SavedQueryRepositoryMock) GetById(id string) (*SavedQuery, error) {
	for _, query := range m.queries {
		if query != nil && query.Id == id {
			return query, nil
		}
	}
	return nil, nil
}

func (m *SavedQueryRepositoryMock) Save(query *SavedQuery) (string, error) {
	if query == nil {
		return "", errors.New("saved query is required")
	}
	if query.Id == "" {
		query.Id = m.nextSavedQueryID()
	}
	for i, existing := range m.queries {
		if existing != nil && existing.Id == query.Id {
			m.queries[i] = query
			m.saved = append(m.saved, query)
			return query.Id, nil
		}
	}
	m.queries = append(m.queries, query)
	m.saved = append(m.saved, query)
	return query.Id, nil
}

func (m *SavedQueryRepositoryMock) ListByWorkspace(workspaceID string) ([]*SavedQuery, error) {
	queries := make([]*SavedQuery, 0)
	for _, query := range m.queries {
		if query != nil && query.WorkspaceId == workspaceID {
			queries = append(queries, query)
		}
	}
	sort.Slice(queries, func(i, j int) bool {
		if queries[i].Name == queries[j].Name {
			return queries[i].Id < queries[j].Id
		}
		return queries[i].Name < queries[j].Name
	})
	return queries, nil
}

func (m *SavedQueryRepositoryMock) Delete(id string) error {
	m.deletedIDs = append(m.deletedIDs, id)
	for i, query := range m.queries {
		if query != nil && query.Id == id {
			m.queries = append(m.queries[:i], m.queries[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *SavedQueryRepositoryMock) nextSavedQueryID() string {
	for {
		m.nextID++
		id := fmt.Sprintf("saved-query-%d", m.nextID)
		query, _ := m.GetById(id)
		if query == nil {
			return id
		}
	}
}

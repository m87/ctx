package core

import (
	"fmt"
	"strings"
)

type SavedQuery struct {
	Id          string `json:"id"`
	WorkspaceId string `json:"workspaceId"`
	Name        string `json:"name"`
	Query       string `json:"query"`
}

type SavedQueryRepository interface {
	GetById(id string) (*SavedQuery, error)
	Save(query *SavedQuery) (string, error)
	ListByWorkspace(workspaceId string) ([]*SavedQuery, error)
	Delete(id string) error
}

type SavedQueryNotFoundError struct {
	SavedQueryId string
}

func (e *SavedQueryNotFoundError) Error() string {
	return fmt.Sprintf("saved query %q not found", e.SavedQueryId)
}

func (m *ContextManager) CreateSavedQuery(query *SavedQuery) (string, error) {
	if query == nil {
		return "", fmt.Errorf("saved query is required")
	}

	query.WorkspaceId = strings.TrimSpace(query.WorkspaceId)
	query.Name = strings.TrimSpace(query.Name)
	query.Query = strings.TrimSpace(query.Query)
	if query.WorkspaceId == "" {
		return "", &WorkspaceNotFoundError{}
	}
	if query.Name == "" {
		return "", fmt.Errorf("saved query name is required")
	}
	if query.Query == "" {
		return "", fmt.Errorf("saved query text is required")
	}
	if m.SavedQueryRepository == nil {
		return "", fmt.Errorf("saved query repository is required")
	}

	workspace, err := m.WorkspaceRepository.GetById(query.WorkspaceId)
	if err != nil {
		return "", err
	}
	if workspace == nil {
		return "", &WorkspaceNotFoundError{WorkspaceId: query.WorkspaceId}
	}

	query.Id = ""
	return m.SavedQueryRepository.Save(query)
}

func (m *ContextManager) ListSavedQueries(workspaceId string) ([]*SavedQuery, error) {
	workspaceId = strings.TrimSpace(workspaceId)
	if workspaceId == "" {
		return nil, &WorkspaceNotFoundError{}
	}
	if m.SavedQueryRepository == nil {
		return nil, fmt.Errorf("saved query repository is required")
	}

	workspace, err := m.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, &WorkspaceNotFoundError{WorkspaceId: workspaceId}
	}
	return m.SavedQueryRepository.ListByWorkspace(workspaceId)
}

func (m *ContextManager) GetSavedQuery(id string) (*SavedQuery, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, &SavedQueryNotFoundError{}
	}
	if m.SavedQueryRepository == nil {
		return nil, fmt.Errorf("saved query repository is required")
	}

	query, err := m.SavedQueryRepository.GetById(id)
	if err != nil {
		return nil, err
	}
	if query == nil {
		return nil, &SavedQueryNotFoundError{SavedQueryId: id}
	}
	return query, nil
}

func (m *ContextManager) DeleteSavedQuery(id string) error {
	id = strings.TrimSpace(id)
	if _, err := m.GetSavedQuery(id); err != nil {
		return err
	}
	return m.SavedQueryRepository.Delete(id)
}

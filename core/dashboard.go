package core

import (
	"encoding/json"
	"fmt"
	"strings"
)

type DashboardType string

const (
	DashboardTypeWorkspace DashboardType = "workspace"
	DashboardTypeProject   DashboardType = "project"
	DashboardTypeContext   DashboardType = "context"
	DashboardTypeDaily     DashboardType = "daily"
	DashboardTypeCustom    DashboardType = "custom"
)

type Dashboard struct {
	Id          string          `json:"id"`
	WorkspaceId string          `json:"workspaceId"`
	Type        DashboardType   `json:"type"`
	TargetId    string          `json:"targetId,omitempty"`
	Name        string          `json:"name"`
	Active      bool            `json:"active"`
	Definition  json.RawMessage `json:"definition"`
}

type DashboardRepository interface {
	GetById(id string) (*Dashboard, error)
	GetByTarget(workspaceId string, dashboardType DashboardType, targetId string) (*Dashboard, error)
	Save(dashboard *Dashboard) (string, error)
	SaveActive(dashboard *Dashboard) (string, error)
	ListByWorkspace(workspaceId string) ([]*Dashboard, error)
	Delete(id string) error
}

type DashboardNotFoundError struct {
	DashboardId string
}

func (e *DashboardNotFoundError) Error() string {
	return fmt.Sprintf("dashboard %q not found", e.DashboardId)
}

type InvalidDashboardError struct {
	Description string
}

func (e *InvalidDashboardError) Error() string {
	return e.Description
}

func (m *ContextManager) CreateDashboard(dashboard *Dashboard) (string, error) {
	if dashboard == nil {
		return "", &InvalidDashboardError{Description: "dashboard is required"}
	}
	if m.DashboardRepository == nil {
		return "", fmt.Errorf("dashboard repository is required")
	}

	dashboard.Id = ""
	dashboard.Active = true
	if err := m.validateDashboard(dashboard); err != nil {
		return "", err
	}
	return m.saveDashboard(dashboard)
}

func (m *ContextManager) UpdateDashboard(dashboard *Dashboard) error {
	if dashboard == nil {
		return &InvalidDashboardError{Description: "dashboard is required"}
	}
	if m.DashboardRepository == nil {
		return fmt.Errorf("dashboard repository is required")
	}

	dashboard.Id = strings.TrimSpace(dashboard.Id)
	if dashboard.Id == "" {
		return &DashboardNotFoundError{}
	}
	existing, err := m.DashboardRepository.GetById(dashboard.Id)
	if err != nil {
		return err
	}
	if existing == nil {
		return &DashboardNotFoundError{DashboardId: dashboard.Id}
	}
	if err := m.validateDashboard(dashboard); err != nil {
		return err
	}
	_, err = m.saveDashboard(dashboard)
	return err
}

func (m *ContextManager) GetDashboard(id string) (*Dashboard, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, &DashboardNotFoundError{}
	}
	if m.DashboardRepository == nil {
		return nil, fmt.Errorf("dashboard repository is required")
	}

	dashboard, err := m.DashboardRepository.GetById(id)
	if err != nil {
		return nil, err
	}
	if dashboard == nil {
		return nil, &DashboardNotFoundError{DashboardId: id}
	}
	return dashboard, nil
}

func (m *ContextManager) ListDashboards(workspaceId string) ([]*Dashboard, error) {
	workspaceId = strings.TrimSpace(workspaceId)
	if workspaceId == "" {
		return nil, &WorkspaceNotFoundError{}
	}
	if m.DashboardRepository == nil {
		return nil, fmt.Errorf("dashboard repository is required")
	}

	workspace, err := m.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, &WorkspaceNotFoundError{WorkspaceId: workspaceId}
	}
	return m.DashboardRepository.ListByWorkspace(workspaceId)
}

func (m *ContextManager) DeleteDashboard(id string) error {
	dashboard, err := m.GetDashboard(id)
	if err != nil {
		return err
	}
	return m.DashboardRepository.Delete(dashboard.Id)
}

func (m *ContextManager) validateDashboard(dashboard *Dashboard) error {
	dashboard.WorkspaceId = strings.TrimSpace(dashboard.WorkspaceId)
	dashboard.TargetId = strings.TrimSpace(dashboard.TargetId)
	dashboard.Name = strings.TrimSpace(dashboard.Name)
	if dashboard.WorkspaceId == "" {
		return &WorkspaceNotFoundError{}
	}
	if dashboard.Name == "" {
		return &InvalidDashboardError{Description: "dashboard name is required"}
	}
	if len(dashboard.Definition) == 0 {
		dashboard.Definition = json.RawMessage(`{}`)
	}
	if !json.Valid(dashboard.Definition) {
		return &InvalidDashboardError{Description: "dashboard definition must be valid JSON"}
	}
	if err := validateDashboardDefinition(dashboard.Definition); err != nil {
		return err
	}

	workspace, err := m.WorkspaceRepository.GetById(dashboard.WorkspaceId)
	if err != nil {
		return err
	}
	if workspace == nil {
		return &WorkspaceNotFoundError{WorkspaceId: dashboard.WorkspaceId}
	}

	switch dashboard.Type {
	case DashboardTypeWorkspace, DashboardTypeDaily:
		if dashboard.TargetId == "" {
			dashboard.TargetId = dashboard.WorkspaceId
		}
		if dashboard.TargetId != dashboard.WorkspaceId {
			return &InvalidDashboardError{Description: "dashboard target must match its workspace"}
		}
	case DashboardTypeProject, DashboardTypeContext:
		dashboard.TargetId = ""
	case DashboardTypeCustom:
		if dashboard.TargetId != "" {
			return &InvalidDashboardError{Description: "custom dashboard cannot have a target"}
		}
	default:
		return &InvalidDashboardError{Description: "invalid dashboard type"}
	}

	return nil
}

func (m *ContextManager) saveDashboard(dashboard *Dashboard) (string, error) {
	if dashboard.Active && dashboard.Type != DashboardTypeCustom {
		return m.DashboardRepository.SaveActive(dashboard)
	}
	return m.DashboardRepository.Save(dashboard)
}

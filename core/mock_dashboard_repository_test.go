package core

import "fmt"

type DashboardRepositoryMock struct {
	dashboards []*Dashboard
	nextId     int
}

func NewDashboardRepositoryMock() *DashboardRepositoryMock {
	return &DashboardRepositoryMock{nextId: 1}
}

func (m *DashboardRepositoryMock) GetById(id string) (*Dashboard, error) {
	for _, dashboard := range m.dashboards {
		if dashboard.Id == id {
			return dashboard, nil
		}
	}
	return nil, nil
}

func (m *DashboardRepositoryMock) GetByTarget(
	workspaceId string,
	dashboardType DashboardType,
	targetId string,
) (*Dashboard, error) {
	for _, dashboard := range m.dashboards {
		if dashboard.WorkspaceId == workspaceId && dashboard.Type == dashboardType && dashboard.TargetId == targetId {
			return dashboard, nil
		}
	}
	return nil, nil
}

func (m *DashboardRepositoryMock) Save(dashboard *Dashboard) (string, error) {
	if dashboard.Id == "" {
		dashboard.Id = fmt.Sprintf("dashboard-%d", m.nextId)
		m.nextId++
	}
	for i, existing := range m.dashboards {
		if existing.Id == dashboard.Id {
			m.dashboards[i] = dashboard
			return dashboard.Id, nil
		}
	}
	m.dashboards = append(m.dashboards, dashboard)
	return dashboard.Id, nil
}

func (m *DashboardRepositoryMock) ListByWorkspace(workspaceId string) ([]*Dashboard, error) {
	dashboards := make([]*Dashboard, 0)
	for _, dashboard := range m.dashboards {
		if dashboard.WorkspaceId == workspaceId {
			dashboards = append(dashboards, dashboard)
		}
	}
	return dashboards, nil
}

func (m *DashboardRepositoryMock) Delete(id string) error {
	for i, dashboard := range m.dashboards {
		if dashboard.Id == id {
			m.dashboards = append(m.dashboards[:i], m.dashboards[i+1:]...)
			return nil
		}
	}
	return nil
}

package core

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDashboardLifecycle(t *testing.T) {
	test := NewTestContextManager()
	dashboard := &Dashboard{
		WorkspaceId: TestWorkspaceID,
		Type:        DashboardTypeProject,
		TargetId:    TestProjectID,
		Name:        "Project insight",
	}

	id, err := test.Manager.CreateDashboard(dashboard)
	require.NoError(t, err)
	require.Equal(t, "dashboard-1", id)
	require.JSONEq(t, `{}`, string(dashboard.Definition))

	dashboard.Name = "Updated insight"
	dashboard.Type = DashboardTypeCustom
	dashboard.TargetId = ""
	dashboard.Definition = json.RawMessage(`{"version":1}`)
	require.NoError(t, test.Manager.UpdateDashboard(dashboard))

	retrieved, err := test.Manager.GetDashboard(id)
	require.NoError(t, err)
	require.Equal(t, dashboard, retrieved)

	dashboards, err := test.Manager.ListDashboards(TestWorkspaceID)
	require.NoError(t, err)
	require.Equal(t, []*Dashboard{dashboard}, dashboards)

	require.NoError(t, test.Manager.DeleteDashboard(id))
	_, err = test.Manager.GetDashboard(id)
	var notFound *DashboardNotFoundError
	require.ErrorAs(t, err, &notFound)
}

func TestDashboardValidation(t *testing.T) {
	test := NewTestContextManager()

	testCases := []struct {
		name      string
		dashboard *Dashboard
	}{
		{name: "nil", dashboard: nil},
		{name: "missing workspace", dashboard: &Dashboard{Name: "Dashboard", Type: DashboardTypeCustom}},
		{name: "missing name", dashboard: &Dashboard{WorkspaceId: TestWorkspaceID, Type: DashboardTypeCustom}},
		{name: "invalid type", dashboard: &Dashboard{WorkspaceId: TestWorkspaceID, Name: "Dashboard", Type: "invalid"}},
		{name: "invalid definition", dashboard: &Dashboard{WorkspaceId: TestWorkspaceID, Name: "Dashboard", Type: DashboardTypeCustom, Definition: json.RawMessage(`{`)}},
		{name: "custom target", dashboard: &Dashboard{WorkspaceId: TestWorkspaceID, Name: "Dashboard", Type: DashboardTypeCustom, TargetId: TestProjectID}},
		{name: "missing project target", dashboard: &Dashboard{WorkspaceId: TestWorkspaceID, Name: "Dashboard", Type: DashboardTypeProject}},
		{name: "workspace target mismatch", dashboard: &Dashboard{WorkspaceId: TestWorkspaceID, Name: "Dashboard", Type: DashboardTypeWorkspace, TargetId: "other"}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := test.Manager.CreateDashboard(testCase.dashboard)
			require.Error(t, err)
		})
	}
}

func TestDashboardRejectsDuplicateSpecialTarget(t *testing.T) {
	test := NewTestContextManager()
	first := &Dashboard{
		WorkspaceId: TestWorkspaceID,
		Type:        DashboardTypeWorkspace,
		TargetId:    TestWorkspaceID,
		Name:        "Workspace insight",
	}
	_, err := test.Manager.CreateDashboard(first)
	require.NoError(t, err)

	_, err = test.Manager.CreateDashboard(&Dashboard{
		WorkspaceId: TestWorkspaceID,
		Type:        DashboardTypeWorkspace,
		TargetId:    TestWorkspaceID,
		Name:        "Duplicate",
	})
	var duplicate *DashboardAlreadyExistsError
	require.ErrorAs(t, err, &duplicate)

	_, err = test.Manager.CreateDashboard(&Dashboard{
		WorkspaceId: TestWorkspaceID,
		Type:        DashboardTypeCustom,
		Name:        "Custom one",
	})
	require.NoError(t, err)
	_, err = test.Manager.CreateDashboard(&Dashboard{
		WorkspaceId: TestWorkspaceID,
		Type:        DashboardTypeCustom,
		Name:        "Custom two",
	})
	require.NoError(t, err)
}

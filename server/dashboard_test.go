package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/stretchr/testify/require"
)

func TestDashboardRoutesSupportLifecycleAtCanonicalAndLegacyPaths(t *testing.T) {
	for _, basePath := range []string{"/api/dashboard", "/dashboard"} {
		t.Run(basePath, func(t *testing.T) {
			manager := bootstrap.NewTestContextManager(time.Now())
			_, err := manager.WorkspaceRepository.Save(&core.Workspace{
				Id:   "workspace-1",
				Name: "Workspace",
			})
			require.NoError(t, err)
			server := NewServer(manager, nil)

			createRequest := httptest.NewRequest(
				http.MethodPost,
				basePath,
				bytes.NewBufferString(
					`{"workspaceId":"workspace-1","type":"workspace","targetId":"workspace-1","name":"Workspace insight","definition":{}}`,
				),
			)
			createResponse := httptest.NewRecorder()
			server.Handler().ServeHTTP(createResponse, createRequest)
			require.Equal(t, http.StatusCreated, createResponse.Code)

			var created core.Dashboard
			require.NoError(t, json.NewDecoder(createResponse.Body).Decode(&created))
			require.NotEmpty(t, created.Id)
			require.Equal(t, core.DashboardTypeWorkspace, created.Type)
			require.JSONEq(t, `{}`, string(created.Definition))

			listRequest := httptest.NewRequest(
				http.MethodGet,
				basePath+"/?workspaceId=workspace-1",
				nil,
			)
			listResponse := httptest.NewRecorder()
			server.Handler().ServeHTTP(listResponse, listRequest)
			require.Equal(t, http.StatusOK, listResponse.Code)
			var dashboards []*core.Dashboard
			require.NoError(t, json.NewDecoder(listResponse.Body).Decode(&dashboards))
			require.Equal(t, []*core.Dashboard{&created}, dashboards)

			created.Name = "Updated insight"
			created.Type = core.DashboardTypeCustom
			created.TargetId = ""
			updateBody, err := json.Marshal(created)
			require.NoError(t, err)
			updateRequest := httptest.NewRequest(
				http.MethodPut,
				fmt.Sprintf("%s/%s", basePath, created.Id),
				bytes.NewReader(updateBody),
			)
			updateResponse := httptest.NewRecorder()
			server.Handler().ServeHTTP(updateResponse, updateRequest)
			require.Equal(t, http.StatusOK, updateResponse.Code)

			getRequest := httptest.NewRequest(
				http.MethodGet,
				fmt.Sprintf("%s/%s", basePath, created.Id),
				nil,
			)
			getResponse := httptest.NewRecorder()
			server.Handler().ServeHTTP(getResponse, getRequest)
			require.Equal(t, http.StatusOK, getResponse.Code)
			var retrieved core.Dashboard
			require.NoError(t, json.NewDecoder(getResponse.Body).Decode(&retrieved))
			require.Equal(t, "Updated insight", retrieved.Name)
			require.Equal(t, core.DashboardTypeCustom, retrieved.Type)
			require.Empty(t, retrieved.TargetId)

			deleteRequest := httptest.NewRequest(
				http.MethodDelete,
				fmt.Sprintf("%s/%s", basePath, created.Id),
				nil,
			)
			deleteResponse := httptest.NewRecorder()
			server.Handler().ServeHTTP(deleteResponse, deleteRequest)
			require.Equal(t, http.StatusNoContent, deleteResponse.Code)

			missingResponse := httptest.NewRecorder()
			server.Handler().ServeHTTP(missingResponse, getRequest)
			require.Equal(t, http.StatusNotFound, missingResponse.Code)
		})
	}
}

func TestDashboardHandlerRejectsDuplicateSpecialDashboard(t *testing.T) {
	manager := bootstrap.NewTestContextManager(time.Now())
	_, err := manager.WorkspaceRepository.Save(&core.Workspace{Id: "workspace-1", Name: "Workspace"})
	require.NoError(t, err)
	server := NewServer(manager, nil)
	body := `{"workspaceId":"workspace-1","type":"daily","targetId":"workspace-1","name":"Daily insight","definition":{}}`

	firstResponse := httptest.NewRecorder()
	server.Handler().ServeHTTP(
		firstResponse,
		httptest.NewRequest(http.MethodPost, "/api/dashboard/", bytes.NewBufferString(body)),
	)
	require.Equal(t, http.StatusCreated, firstResponse.Code)

	secondResponse := httptest.NewRecorder()
	server.Handler().ServeHTTP(
		secondResponse,
		httptest.NewRequest(http.MethodPost, "/api/dashboard/", bytes.NewBufferString(body)),
	)
	require.Equal(t, http.StatusConflict, secondResponse.Code)
	var response ErrorResponse
	require.NoError(t, json.NewDecoder(secondResponse.Body).Decode(&response))
	require.Equal(t, "DASHBOARD_ALREADY_EXISTS", response.Code)
}

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
			require.True(t, created.Active)
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
			created.Definition = json.RawMessage(`{"version":1,"query":"","widgets":[{"id":"text-1","type":"text","query":"","layout":{"x":0,"y":0,"width":4,"height":3},"properties":{"text":"Focus","horizontalAlign":"center","verticalAlign":"bottom","fontSize":32}}]}`)
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
			require.JSONEq(t, string(created.Definition), string(retrieved.Definition))

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

func TestDashboardHandlerReplacesActiveSpecialDashboard(t *testing.T) {
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
	require.Equal(t, http.StatusCreated, secondResponse.Code)
	var first, second core.Dashboard
	require.NoError(t, json.NewDecoder(firstResponse.Body).Decode(&first))
	require.NoError(t, json.NewDecoder(secondResponse.Body).Decode(&second))
	require.True(t, second.Active)
	stored, err := manager.GetDashboard(first.Id)
	require.NoError(t, err)
	require.False(t, stored.Active)
}

func TestDashboardEvaluationRoutesReturnEquivalentMetricResults(t *testing.T) {
	now := time.Date(2026, time.September, 27, 16, 0, 0, 0, time.UTC)
	manager := bootstrap.NewTestContextManager(now)
	_, err := manager.WorkspaceRepository.Save(&core.Workspace{Id: "workspace-1", Name: "Workspace"})
	require.NoError(t, err)
	contextID, err := manager.CreateContext(&core.Context{WorkspaceId: "workspace-1", Name: "Focus"})
	require.NoError(t, err)
	start := time.Date(2026, time.September, 27, 8, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)
	_, err = manager.SaveInterval(&core.Interval{ContextId: contextID, Start: &start, End: &end})
	require.NoError(t, err)
	definition := json.RawMessage(`{"version":1,"query":"","period":"today","widgets":[{"id":"kpi","type":"kpi","query":"","layout":{"x":0,"y":0,"width":4,"height":3},"properties":{"title":"Tracked time","description":"","metric":"duration.sum","period":"today","startDate":"","endDate":"","projectId":"","contextId":"","tag":"","archived":"include","includeRunning":false,"format":"duration","precision":0}},{"id":"gauge","type":"gauge","query":"","layout":{"x":4,"y":0,"width":6,"height":5},"properties":{"title":"Goal","description":"","metric":"duration.sum","period":"today","startDate":"","endDate":"","projectId":"","contextId":"","tag":"","archived":"include","includeRunning":false,"format":"duration","precision":0,"minimum":0,"maximum":144000000,"target":144000000,"goalMode":"at_least","thresholds":[],"showValue":true,"showTarget":true,"showPercentage":true}}]}`)
	var definitionObject map[string]interface{}
	require.NoError(t, json.Unmarshal(definition, &definitionObject))
	widgets := definitionObject["widgets"].([]interface{})
	widgets = append(widgets, map[string]interface{}{
		"id": "radar", "type": "radar", "query": "", "layout": map[string]int{"x": 10, "y": 0, "width": 6, "height": 6},
		"properties": map[string]interface{}{
			"title": "Week", "description": "", "metric": "duration.sum", "period": "today", "startDate": "", "endDate": "", "projectId": "", "contextId": "", "tag": "", "archived": "include", "includeRunning": false, "format": "duration", "precision": 0,
			"axesBy": "weekday", "axisKeys": []string{}, "scaleMaximum": nil, "targetValue": nil, "showValues": false, "showGrid": true, "showFill": true, "showSummary": true,
		},
	})
	definitionObject["widgets"] = widgets
	definition, err = json.Marshal(definitionObject)
	require.NoError(t, err)
	requestBody, err := json.Marshal(core.DashboardEvaluationRequest{
		WorkspaceID: "workspace-1", Definition: definition, TimeZone: "UTC", FirstDayOfWeek: 1,
	})
	require.NoError(t, err)
	server := NewServer(manager, nil)
	var results []*core.DashboardEvaluation
	for _, path := range []string{"/api/dashboard/evaluate", "/dashboard/evaluate"} {
		response := httptest.NewRecorder()
		server.Handler().ServeHTTP(response, httptest.NewRequest(http.MethodPost, path, bytes.NewReader(requestBody)))
		require.Equal(t, http.StatusOK, response.Code)
		var result core.DashboardEvaluation
		require.NoError(t, json.NewDecoder(response.Body).Decode(&result))
		results = append(results, &result)
		require.Equal(t, now, result.AsOf)
		require.Equal(t, float64(time.Hour/time.Millisecond), *result.Widgets["kpi"].Value)
		require.Equal(t, "ms", result.Widgets["gauge"].Unit)
		require.Equal(t, result.Widgets["kpi"].Start, result.Widgets["gauge"].Start)
		require.Equal(t, result.Widgets["kpi"].End, result.Widgets["gauge"].End)
		require.Equal(t, "radar", result.Widgets["radar"].Kind)
		require.Len(t, result.Widgets["radar"].Radar.Axes, 7)
		require.True(t, *result.Widgets["radar"].HasActivity)
	}
	require.Equal(t, results[0], results[1])
}

func TestDashboardEvaluationRejectsMalformedBodyAndInvalidTimeZone(t *testing.T) {
	manager := bootstrap.NewTestContextManager(time.Now())
	_, err := manager.WorkspaceRepository.Save(&core.Workspace{Id: "workspace-1", Name: "Workspace"})
	require.NoError(t, err)
	server := NewServer(manager, nil)
	for _, tc := range []struct {
		body string
		code int
	}{
		{body: `{`, code: http.StatusBadRequest},
		{body: `{"workspaceId":"workspace-1","definition":{},"timeZone":"Invalid/Zone"}`, code: http.StatusBadRequest},
	} {
		response := httptest.NewRecorder()
		server.Handler().ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/api/dashboard/evaluate", bytes.NewBufferString(tc.body)))
		require.Equal(t, tc.code, response.Code)
		var body map[string]string
		require.NoError(t, json.NewDecoder(response.Body).Decode(&body))
		require.NotEmpty(t, body["code"])
		require.NotEmpty(t, body["description"])
	}
}

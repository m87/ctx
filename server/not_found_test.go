package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/storage"
	"github.com/stretchr/testify/require"
)

func TestMissingResourcesReturnNotFoundOnBothRouteFamilies(t *testing.T) {
	manager := bootstrap.NewTestContextManager(time.Date(2026, 9, 29, 12, 0, 0, 0, time.UTC))
	handler := NewServer(manager, nil).Handler()
	for _, prefix := range []string{"", "/api"} {
		for _, testCase := range []struct {
			path string
			code string
		}{
			{path: "/context/missing", code: "CONTEXT_NOT_FOUND"},
			{path: "/project/missing", code: "PROJECT_NOT_FOUND"},
			{path: "/interval/missing/move/target", code: "INTERVAL_NOT_FOUND"},
		} {
			t.Run(prefix+testCase.path, func(t *testing.T) {
				method := http.MethodGet
				if testCase.path == "/interval/missing/move/target" {
					method = http.MethodPatch
				}
				request := httptest.NewRequest(method, prefix+testCase.path, nil)
				response := httptest.NewRecorder()
				handler.ServeHTTP(response, request)
				require.Equal(t, http.StatusNotFound, response.Code)
				var body struct {
					Code string `json:"code"`
				}
				require.NoError(t, json.NewDecoder(response.Body).Decode(&body))
				require.Equal(t, testCase.code, body.Code)
			})
		}
	}
}

func TestRepositoryFailureRemainsInternalError(t *testing.T) {
	database, err := storage.NewSqliteStorage(":memory:")
	require.NoError(t, err)
	sqlDB, err := database.DB.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())
	manager := core.NewContextManager(
		nil,
		storage.NewContextRepository(database.DB),
		storage.NewIntervalRepository(database.DB),
		storage.NewWorkspaceRepository(database.DB),
		storage.NewProjectRepository(database.DB),
	)
	handler := NewServer(manager, nil).Handler()
	for _, testCase := range []struct {
		method string
		path   string
	}{
		{method: http.MethodGet, path: "/api/context/record"},
		{method: http.MethodGet, path: "/api/project/record"},
		{method: http.MethodPatch, path: "/api/interval/record/move/target"},
	} {
		t.Run(testCase.path, func(t *testing.T) {
			request := httptest.NewRequest(testCase.method, testCase.path, nil)
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			require.Equal(t, http.StatusInternalServerError, response.Code)
		})
	}
}

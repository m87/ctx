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

func TestSavedQueryRoutesCreateListAndGetAtCanonicalAndLegacyPaths(t *testing.T) {
	for _, basePath := range []string{"/api/query/saved", "/query/saved"} {
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
				basePath+"/",
				bytes.NewBufferString(
					`{"workspaceId":"workspace-1","name":"Focus","query":"tag = focus"}`,
				),
			)
			createResponse := httptest.NewRecorder()
			server.Handler().ServeHTTP(createResponse, createRequest)
			require.Equal(t, http.StatusCreated, createResponse.Code)

			var created core.SavedQuery
			require.NoError(t, json.NewDecoder(createResponse.Body).Decode(&created))
			require.NotEmpty(t, created.Id)
			require.Equal(t, "Focus", created.Name)

			listRequest := httptest.NewRequest(
				http.MethodGet,
				basePath+"/?workspaceId=workspace-1",
				nil,
			)
			listResponse := httptest.NewRecorder()
			server.Handler().ServeHTTP(listResponse, listRequest)
			require.Equal(t, http.StatusOK, listResponse.Code)
			var queries []*core.SavedQuery
			require.NoError(t, json.NewDecoder(listResponse.Body).Decode(&queries))
			require.Equal(t, []*core.SavedQuery{&created}, queries)

			getRequest := httptest.NewRequest(
				http.MethodGet,
				fmt.Sprintf("%s/%s", basePath, created.Id),
				nil,
			)
			getResponse := httptest.NewRecorder()
			server.Handler().ServeHTTP(getResponse, getRequest)
			require.Equal(t, http.StatusOK, getResponse.Code)
			var retrieved core.SavedQuery
			require.NoError(t, json.NewDecoder(getResponse.Body).Decode(&retrieved))
			require.Equal(t, created, retrieved)

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

func TestSavedQueryHandlerValidatesCreateRequest(t *testing.T) {
	manager := bootstrap.NewTestContextManager(time.Now())
	server := NewServer(manager, nil)
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/query/saved/",
		bytes.NewBufferString(`{"workspaceId":"workspace-1","query":"tag = focus"}`),
	)
	response := httptest.NewRecorder()

	server.Handler().ServeHTTP(response, request)

	require.Equal(t, http.StatusBadRequest, response.Code)
	var body ErrorResponse
	require.NoError(t, json.NewDecoder(response.Body).Decode(&body))
	require.Equal(t, "MISSING_SAVED_QUERY_NAME", body.Code)
}

func TestSavedQueryHandlerReturnsNotFound(t *testing.T) {
	manager := bootstrap.NewTestContextManager(time.Now())
	server := NewServer(manager, nil)
	request := httptest.NewRequest(http.MethodGet, "/api/query/saved/missing", nil)
	response := httptest.NewRecorder()

	server.Handler().ServeHTTP(response, request)

	require.Equal(t, http.StatusNotFound, response.Code)
	var body ErrorResponse
	require.NoError(t, json.NewDecoder(response.Body).Decode(&body))
	require.Equal(t, "SAVED_QUERY_NOT_FOUND", body.Code)
}

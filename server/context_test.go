package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/stretchr/testify/require"
)

func TestContextHandlerReturnsTagObjects(t *testing.T) {
	manager := bootstrap.NewTestContextManager(time.Now())
	_, err := manager.WorkspaceRepository.Save(&core.Workspace{Id: "workspace-1", Name: "Workspace"})
	require.NoError(t, err)

	mux := http.NewServeMux()
	registerContextHandler(mux, manager)
	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{
		"name":"Context",
		"workspaceId":"workspace-1",
		"tags":[
			{"id":"","name":"important"},
			{"id":"","name":"backend"}
		]
	}`))
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)

	require.Equal(t, http.StatusCreated, response.Code)
	var body core.Context
	require.NoError(t, json.NewDecoder(response.Body).Decode(&body))
	require.Len(t, body.Tags, 2)
	require.ElementsMatch(t, []string{"important", "backend"}, []string{body.Tags[0].Name, body.Tags[1].Name})
	for _, tag := range body.Tags {
		require.NotEmpty(t, tag.Id)
	}
}

func TestContextHandlerTreatsMissingActiveContextAsNotFoundAndFreeAsNoOp(t *testing.T) {
	manager := bootstrap.NewTestContextManager(time.Now())
	mux := http.NewServeMux()
	registerContextHandler(mux, manager)

	activeRequest := httptest.NewRequest(http.MethodGet, "/active", nil)
	activeResponse := httptest.NewRecorder()
	mux.ServeHTTP(activeResponse, activeRequest)
	require.Equal(t, http.StatusNotFound, activeResponse.Code)

	freeRequest := httptest.NewRequest(http.MethodPost, "/free", bytes.NewBufferString(`{}`))
	freeResponse := httptest.NewRecorder()
	mux.ServeHTTP(freeResponse, freeRequest)
	require.Equal(t, http.StatusNoContent, freeResponse.Code)
}

package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/m87/ctx/bootstrap"
	"github.com/m87/ctx/core"
	"github.com/stretchr/testify/require"
)

func TestQueryHandlerReturnsAllContexts(t *testing.T) {
	manager := newQueryTestManager(t)
	mux := http.NewServeMux()
	registerQueryHandler(mux, manager)

	request := httptest.NewRequest(
		http.MethodPost,
		"/",
		bytes.NewBufferString(`{"query":"ignored query"}`),
	)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)

	require.Equal(t, http.StatusOK, response.Code)
	var contexts []*core.Context
	require.NoError(t, json.NewDecoder(response.Body).Decode(&contexts))
	require.Len(t, contexts, 3)
	require.ElementsMatch(
		t,
		[]string{"Active", "Archived", "Other workspace"},
		[]string{contexts[0].Name, contexts[1].Name, contexts[2].Name},
	)
}

func TestQueryHandlerValidatesRequest(t *testing.T) {
	manager := newQueryTestManager(t)
	mux := http.NewServeMux()
	registerQueryHandler(mux, manager)

	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{`))
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)

	require.Equal(t, http.StatusBadRequest, response.Code)
	var body ErrorResponse
	require.NoError(t, json.NewDecoder(response.Body).Decode(&body))
	require.Equal(t, "INVALID_REQUEST_BODY", body.Code)
}

func TestQueryHandlerReturnsInvalidQueryError(t *testing.T) {
	manager := newQueryTestManager(t)
	manager.QueryInterpreter = &invalidQueryInterpreter{}
	mux := http.NewServeMux()
	registerQueryHandler(mux, manager)

	request := httptest.NewRequest(
		http.MethodPost,
		"/",
		bytes.NewBufferString(`{"query":"broken"}`),
	)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)

	require.Equal(t, http.StatusBadRequest, response.Code)
	var body ErrorResponse
	require.NoError(t, json.NewDecoder(response.Body).Decode(&body))
	require.Equal(t, "INVALID_QUERY", body.Code)
	require.Equal(t, "invalid context query: unexpected token", body.Description)
}

func TestQueryRoutesAreMountedAtCanonicalAndLegacyPaths(t *testing.T) {
	manager := newQueryTestManager(t)
	server := NewServer(manager, nil)

	for _, path := range []string{"/api/query", "/api/query/", "/query", "/query/"} {
		t.Run(path, func(t *testing.T) {
			request := httptest.NewRequest(
				http.MethodPost,
				path,
				bytes.NewBufferString(`{"query":""}`),
			)
			response := httptest.NewRecorder()
			server.Handler().ServeHTTP(response, request)

			require.Equal(t, http.StatusOK, response.Code)
		})
	}
}

func newQueryTestManager(t *testing.T) *core.ContextManager {
	t.Helper()
	manager := bootstrap.NewTestContextManager(time.Now())
	_, err := manager.WorkspaceRepository.Save(&core.Workspace{Id: "workspace-1", Name: "Workspace"})
	require.NoError(t, err)
	_, err = manager.WorkspaceRepository.Save(&core.Workspace{Id: "workspace-2", Name: "Other"})
	require.NoError(t, err)
	_, err = manager.ContextRepository.Save(&core.Context{
		Id:          "context-1",
		Name:        "Active",
		WorkspaceId: "workspace-1",
	})
	require.NoError(t, err)
	_, err = manager.ContextRepository.Save(&core.Context{
		Id:          "context-2",
		Name:        "Archived",
		WorkspaceId: "workspace-1",
		Archived:    true,
	})
	require.NoError(t, err)
	_, err = manager.ContextRepository.Save(&core.Context{
		Id:          "context-3",
		Name:        "Other workspace",
		WorkspaceId: "workspace-2",
	})
	require.NoError(t, err)
	return manager
}

type invalidQueryInterpreter struct{}

func (i *invalidQueryInterpreter) Interpret(_ string) (*core.ContextSQLQuery, error) {
	return nil, &core.InvalidContextQueryError{Err: errors.New("unexpected token")}
}

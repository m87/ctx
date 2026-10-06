package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/m87/ctx/bootstrap"
	"github.com/stretchr/testify/require"
)

func TestSPARoutingKeepsJSONAndLegacyRequestsOnAPI(t *testing.T) {
	manager := bootstrap.NewTestContextManager(time.Date(2026, 9, 29, 12, 0, 0, 0, time.UTC))
	server := NewServer(manager, nil)
	server.spa = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Route", "SPA")
		w.WriteHeader(http.StatusOK)
	})
	handler := server.Handler()

	for _, testCase := range []struct {
		name   string
		method string
		path   string
		accept string
		spa    bool
	}{
		{name: "browser context", method: http.MethodGet, path: "/context/missing", accept: "text/html,application/xhtml+xml", spa: true},
		{name: "browser project", method: http.MethodGet, path: "/project/missing", accept: "text/html", spa: true},
		{name: "browser workspace", method: http.MethodGet, path: "/workspace", accept: "text/html", spa: true},
		{name: "browser dashboard", method: http.MethodGet, path: "/dashboard/missing", accept: "text/html", spa: true},
		{name: "browser query", method: http.MethodGet, path: "/query/missing", accept: "text/html", spa: true},
		{name: "browser head", method: http.MethodHead, path: "/context/missing", accept: "text/html", spa: true},
		{name: "json context", method: http.MethodGet, path: "/context/missing", accept: "application/json"},
		{name: "json dashboard", method: http.MethodGet, path: "/dashboard/missing", accept: "application/json"},
		{name: "disabled html", method: http.MethodGet, path: "/context/missing", accept: "text/html;q=0,application/json"},
		{name: "saved query api", method: http.MethodGet, path: "/query/saved/missing", accept: "text/html"},
		{name: "dashboard id matching action name", method: http.MethodGet, path: "/dashboard/evaluate", accept: "text/html", spa: true},
		{name: "api always api", method: http.MethodGet, path: "/api/context/missing", accept: "text/html"},
		{name: "api prefix boundary", method: http.MethodGet, path: "/apiary", accept: "text/html", spa: true},
		{name: "legacy prefix boundary", method: http.MethodGet, path: "/contextual", accept: "application/json", spa: true},
		{name: "mutation always api", method: http.MethodPost, path: "/context/missing", accept: "text/html"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			request := httptest.NewRequest(testCase.method, testCase.path, nil)
			request.Header.Set("Accept", testCase.accept)
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			require.Equal(t, testCase.spa, response.Header().Get("X-Route") == "SPA")
		})
	}
}

func TestRequestBodyLimitReturnsJSONError(t *testing.T) {
	handler := withRequestBodyLimit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}), 4)
	request := httptest.NewRequest(http.MethodPost, "/api/sync/intervals", bytes.NewBufferString("12345"))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	require.Equal(t, http.StatusRequestEntityTooLarge, response.Code)
	var body struct {
		Code string `json:"code"`
	}
	require.NoError(t, json.NewDecoder(response.Body).Decode(&body))
	require.Equal(t, "REQUEST_TOO_LARGE", body.Code)
}

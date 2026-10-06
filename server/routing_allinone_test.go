//go:build allinone

package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAllInOneServesSPAForBrowserNavigation(t *testing.T) {
	server := NewServer(nil, nil)
	handler := server.Handler()
	for _, method := range []string{http.MethodGet, http.MethodHead} {
		request := httptest.NewRequest(method, "/dashboard/example", nil)
		request.Header.Set("Accept", "text/html")
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		require.Equal(t, http.StatusOK, response.Code)
		require.Equal(t, "text/html; charset=utf-8", response.Header().Get("Content-Type"))
		require.Equal(t, "no-cache", response.Header().Get("Cache-Control"))
		if method == http.MethodHead {
			require.Empty(t, response.Body.String())
		} else {
			require.Contains(t, response.Body.String(), "<html")
		}
	}
}

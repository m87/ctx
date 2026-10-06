package core

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSyncRejectsInvalidDownloadedContextBatch(t *testing.T) {
	remote := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/version":
			w.Write([]byte(`{"dbVersion":"` + CurrentDatabaseVersion + `"}`))
		case "/sync/workspaces":
			w.Write([]byte(`[{"id":"workspace-1","name":"Remote"}]`))
		case "/sync/contexts":
			w.Write([]byte(`[{"id":"valid","name":"Valid","workspaceId":"workspace-1","status":"inactive"},{"id":"orphan","name":"Orphan","workspaceId":"missing"}]`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer remote.Close()

	test := NewEmptyTestContextManager()
	err := test.Manager.Sync(remote.URL)
	var validationError *SyncValidationError
	require.ErrorAs(t, err, &validationError)
	require.Equal(t, "context", validationError.Resource)
	contexts, err := test.Contexts.List()
	require.NoError(t, err)
	require.Empty(t, contexts)
	workspaces, err := test.Workspaces.List()
	require.NoError(t, err)
	require.Len(t, workspaces, 1)
}

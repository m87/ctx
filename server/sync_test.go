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

func TestSyncUploadRejectsInvalidBatchesWithoutSavingValidItems(t *testing.T) {
	cases := []struct {
		name     string
		resource string
		payload  string
		list     func(*core.ContextManager) (int, error)
	}{
		{
			name:     "null workspace",
			resource: "workspaces",
			payload:  `[{"id":"new-workspace","name":"New"},null]`,
			list: func(manager *core.ContextManager) (int, error) {
				items, err := manager.WorkspaceRepository.List()
				return len(items), err
			},
		},
		{
			name:     "orphan context",
			resource: "contexts",
			payload:  `[{"id":"new-context","name":"New","workspaceId":"workspace-1","status":"inactive"},{"id":"orphan","name":"Orphan","workspaceId":"missing"}]`,
			list: func(manager *core.ContextManager) (int, error) {
				items, err := manager.ContextRepository.List()
				return len(items), err
			},
		},
		{
			name:     "interval workspace mismatch",
			resource: "intervals",
			payload:  `[{"id":"new-interval","contextId":"context-1","workspaceId":"workspace-1","start":"2026-09-29T08:00:00Z","end":"2026-09-29T09:00:00Z","status":"completed"},{"id":"mismatch","contextId":"context-1","workspaceId":"missing","start":"2026-09-29T08:00:00Z","end":"2026-09-29T09:00:00Z","status":"completed"}]`,
			list: func(manager *core.ContextManager) (int, error) {
				items, err := manager.IntervalRepository.List()
				return len(items), err
			},
		},
	}

	for _, prefix := range []string{"/api/sync/", "/sync/"} {
		for _, testCase := range cases {
			t.Run(prefix+testCase.name, func(t *testing.T) {
				manager := bootstrap.NewTestContextManager(time.Date(2026, 9, 29, 12, 0, 0, 0, time.UTC))
				_, err := manager.WorkspaceRepository.Save(&core.Workspace{Id: "workspace-1", Name: "Existing"})
				require.NoError(t, err)
				_, err = manager.ContextRepository.Save(&core.Context{Id: "context-1", Name: "Existing", WorkspaceId: "workspace-1", Status: "inactive"})
				require.NoError(t, err)

				before, err := testCase.list(manager)
				require.NoError(t, err)
				request := httptest.NewRequest(http.MethodPost, prefix+testCase.resource, bytes.NewBufferString(testCase.payload))
				response := httptest.NewRecorder()
				NewServer(manager, nil).Handler().ServeHTTP(response, request)

				require.Equal(t, http.StatusBadRequest, response.Code)
				var body struct {
					Code        string `json:"code"`
					Description string `json:"description"`
				}
				require.NoError(t, json.NewDecoder(response.Body).Decode(&body))
				require.Equal(t, "INVALID_SYNC_ITEM", body.Code)
				require.NotEmpty(t, body.Description)
				after, err := testCase.list(manager)
				require.NoError(t, err)
				require.Equal(t, before, after)
			})
		}
	}
}

func TestSyncUploadNormalizesIntervalDurationAndAcceptsRepeatedBatch(t *testing.T) {
	manager := bootstrap.NewTestContextManager(time.Date(2026, 9, 29, 12, 0, 0, 0, time.UTC))
	_, err := manager.WorkspaceRepository.Save(&core.Workspace{Id: "workspace-1", Name: "Existing"})
	require.NoError(t, err)
	_, err = manager.ContextRepository.Save(&core.Context{Id: "context-1", Name: "Existing", WorkspaceId: "workspace-1", Status: "inactive"})
	require.NoError(t, err)

	payload := `[{"id":"interval-1","contextId":"context-1","workspaceId":"workspace-1","start":"2026-09-29T08:00:00+02:00","end":"2026-09-29T09:00:00+02:00","duration":999,"status":"completed"}]`
	for _, prefix := range []string{"/api/sync/", "/sync/"} {
		request := httptest.NewRequest(http.MethodPost, prefix+"intervals", bytes.NewBufferString(payload))
		response := httptest.NewRecorder()
		NewServer(manager, nil).Handler().ServeHTTP(response, request)
		require.Equal(t, http.StatusOK, response.Code)
	}
	items, err := manager.IntervalRepository.List()
	require.NoError(t, err)
	require.Len(t, items, 1)
	require.Equal(t, time.Hour, items[0].Duration)
	require.Equal(t, time.UTC, items[0].Start.Location())
}

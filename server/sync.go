package server

import (
	"net/http"

	"github.com/m87/ctx/core"
)

type SyncHandler struct {
	manager *core.ContextManager
}

func NewSyncHandler(manager *core.ContextManager) *SyncHandler {
	return &SyncHandler{manager: manager}
}

func registerSyncHandler(mux *http.ServeMux, manager *core.ContextManager) {
	handler := NewSyncHandler(manager)
	mux.HandleFunc("/workspaces", handler.GetWorkspacesToSync)
}

func (h *SyncHandler) GetWorkspacesToSync(w http.ResponseWriter, r *http.Request) {
	workspaces, err := h.manager.WorkspaceRepository.ListToSync(0)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "FAILED_TO_LIST_WORKSPACES", "Failed to list workspaces")
		return
	}
	writeJson(w, http.StatusOK, workspaces)
}

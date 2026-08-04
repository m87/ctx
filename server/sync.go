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
	mux.HandleFunc("/contexts", handler.GetContextsToSync)
	mux.HandleFunc("/intervals", handler.GetIntervalsToSync)
}

func (h *SyncHandler) GetWorkspacesToSync(w http.ResponseWriter, r *http.Request) {
	workspaces, err := h.manager.WorkspaceRepository.ListToSync(0)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "FAILED_TO_LIST_WORKSPACES", "Failed to list workspaces")
		return
	}
	writeJson(w, http.StatusOK, workspaces)
}

func (h *SyncHandler) GetContextsToSync(w http.ResponseWriter, r *http.Request) {
	contexts, err := h.manager.ContextRepository.ListToSync(0)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "FAILED_TO_LIST_CONTEXTS", "Failed to list contexts")
		return
	}
	writeJson(w, http.StatusOK, contexts)
}

func (h *SyncHandler) GetIntervalsToSync(w http.ResponseWriter, r *http.Request) {
	intervals, err := h.manager.IntervalRepository.ListToSync(0)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "FAILED_TO_LIST_INTERVALS", "Failed to list intervals")
		return
	}
	writeJson(w, http.StatusOK, intervals)
}

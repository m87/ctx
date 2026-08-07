package server

import (
	"encoding/json"
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
	mux.HandleFunc("GET /workspaces", handler.GetWorkspacesToSync)
	mux.HandleFunc("GET /contexts", handler.GetContextsToSync)
	mux.HandleFunc("GET /intervals", handler.GetIntervalsToSync)
	mux.HandleFunc("POST /workspaces", handler.UploadWorkspaces)
	mux.HandleFunc("POST /contexts", handler.UploadContexts)
	mux.HandleFunc("POST /intervals", handler.UploadIntervals)
}

func (h *SyncHandler) UploadWorkspaces(w http.ResponseWriter, r *http.Request) {
	var workspaces []*core.Workspace
	if err := json.NewDecoder(r.Body).Decode(&workspaces); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON payload")
		return
	}

	if _, err := h.manager.WorkspaceRepository.SaveAll(workspaces); err != nil {
		writeError(w, http.StatusInternalServerError, "FAILED_TO_SAVE_WORKSPACES", "Failed to save workspaces")
		return
	}
}

func (h *SyncHandler) UploadContexts(w http.ResponseWriter, r *http.Request) {
	var contexts []*core.Context
	if err := json.NewDecoder(r.Body).Decode(&contexts); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON payload")
		return
	}

	if _, err := h.manager.ContextRepository.SaveAll(contexts); err != nil {
		writeError(w, http.StatusInternalServerError, "FAILED_TO_SAVE_CONTEXTS", "Failed to save contexts")
		return
	}
}

func (h *SyncHandler) UploadIntervals(w http.ResponseWriter, r *http.Request) {
	var intervals []*core.Interval
	if err := json.NewDecoder(r.Body).Decode(&intervals); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON payload")
		return
	}

	if _, err := h.manager.IntervalRepository.SaveAll(intervals); err != nil {
		writeError(w, http.StatusInternalServerError, "FAILED_TO_SAVE_INTERVALS", "Failed to save intervals")
		return
	}
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

package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/m87/ctx/core"
)

type ContextHandler struct {
	manager *core.ContextManager
}

func registerContextHandler(mux *http.ServeMux, manager *core.ContextManager) {
	handler := &ContextHandler{manager: manager}
	mux.HandleFunc("GET /", handler.listContexts)
	mux.HandleFunc("POST /", handler.createContext)
	mux.HandleFunc("DELETE /{id}", handler.deleteContext)
	mux.HandleFunc("GET /{id}", handler.getContext)
	mux.HandleFunc("PUT /{id}", handler.updateContext)
	mux.HandleFunc("POST /switch", handler.switchContext)
	mux.HandleFunc("POST /free", handler.freeContext)
	mux.HandleFunc("GET /active", handler.getActiveContext)
	mux.HandleFunc("GET /{id}/intervals", handler.listIntervals)
	mux.HandleFunc("GET /{id}/stats/{date}", handler.getStats)
}

func (h *ContextHandler) getStats(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeError(w, http.StatusBadRequest, "MISSING_CONTEXT_ID", "Missing context ID")
		return
	}
	dateStr := r.PathValue("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_DATE_FORMAT", "Invalid date format, expected YYYY-MM-DD")
		return
	}
	stats, err := h.manager.GetStats(id, date)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "FAILED_TO_GET_CONTEXT_STATS", "Failed to get context stats")
		return
	}

	writeJson(w, http.StatusOK, stats)
}

func (h *ContextHandler) listIntervals(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeError(w, http.StatusBadRequest, "MISSING_CONTEXT_ID", "Missing context ID")
		return
	}
	intervals, err := h.manager.IntervalRepository.ListByContextId(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "FAILED_TO_LIST_INTERVALS", "Failed to list intervals")
		return
	}

	writeJson(w, http.StatusOK, intervals)
}

func (h *ContextHandler) getActiveContext(w http.ResponseWriter, r *http.Request) {
	activeContext, err := h.manager.ContextRepository.GetActive()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "FAILED_TO_GET_ACTIVE_CONTEXT", "Failed to get active context")
		return
	}
	if activeContext == nil {
		writeError(w, http.StatusNotFound, "ACTIVE_CONTEXT_NOT_FOUND", "No active context found")
		return
	}
	writeJson(w, http.StatusOK, activeContext)
}

func (h *ContextHandler) switchContext(w http.ResponseWriter, r *http.Request) {

	var req *core.Context
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", "Invalid request body")
		return
	}
	if err := h.manager.SwitchContext(req); err != nil {
		if _, ok := err.(*core.WorkspaceNotFoundError); ok {
			writeError(w, http.StatusBadRequest, "WORKSPACE_NOT_FOUND", err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "FAILED_TO_SWITCH_CONTEXT", "Failed to switch context")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ContextHandler) freeContext(w http.ResponseWriter, r *http.Request) {
	if err := h.manager.FreeActiveContext(); err != nil {
		writeError(w, http.StatusInternalServerError, "FAILED_TO_FREE_ACTIVE_CONTEXT", "Failed to free active context")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ContextHandler) listContexts(w http.ResponseWriter, r *http.Request) {
	workspaceId := r.URL.Query().Get("workspaceId")
	if workspaceId == "" {
		writeError(w, http.StatusBadRequest, "MISSING_WORKSPACE_ID", "Missing workspace ID")
		return
	}
	contexts, err := h.manager.ContextRepository.ListByWorkspace(workspaceId)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "FAILED_TO_LIST_CONTEXTS", "Failed to list contexts")
		return
	}
	writeJson(w, http.StatusOK, contexts)
}

func (h *ContextHandler) createContext(w http.ResponseWriter, r *http.Request) {
	var context core.Context
	if err := json.NewDecoder(r.Body).Decode(&context); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", "Invalid request body")
		return
	}
	id, err := h.manager.CreateContext(&context)
	if err != nil {
		if _, ok := err.(*core.WorkspaceNotFoundError); ok {
			writeError(w, http.StatusBadRequest, "WORKSPACE_NOT_FOUND", err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "FAILED_TO_CREATE_CONTEXT", "Failed to create context")
		return
	}
	context.Id = id
	writeJson(w, http.StatusCreated, &context)
}

func (h *ContextHandler) deleteContext(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeError(w, http.StatusBadRequest, "MISSING_CONTEXT_ID", "Missing context ID")
		return
	}
	if err := h.manager.DeleteContext(id); err != nil {
		writeError(w, http.StatusInternalServerError, "FAILED_TO_DELETE_CONTEXT", "Failed to delete context")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ContextHandler) getContext(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeError(w, http.StatusBadRequest, "MISSING_CONTEXT_ID", "Missing context ID")
		return
	}
	context, err := h.manager.ContextRepository.GetById(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "FAILED_TO_GET_CONTEXT", "Failed to get context")
		return
	}
	if context == nil {
		writeError(w, http.StatusNotFound, "CONTEXT_NOT_FOUND", "Context not found")
		return
	}
	writeJson(w, http.StatusOK, context)
}

func (h *ContextHandler) updateContext(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeError(w, http.StatusBadRequest, "MISSING_CONTEXT_ID", "Missing context ID")
		return
	}
	var context core.Context
	if err := json.NewDecoder(r.Body).Decode(&context); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", "Invalid request body")
		return
	}
	context.Id = id
	if err := h.manager.UpdateContext(&context); err != nil {
		if _, ok := err.(*core.ContextNotFoundError); ok {
			writeError(w, http.StatusNotFound, "CONTEXT_NOT_FOUND", "Context not found")
			return
		}
		if _, ok := err.(*core.ContextWorkspaceMoveNotAllowedError); ok {
			writeError(w, http.StatusBadRequest, "CONTEXT_WORKSPACE_MOVE_NOT_ALLOWED", err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "FAILED_TO_UPDATE_CONTEXT", "Failed to update context")
		return
	}
	writeJson(w, http.StatusOK, &context)
}

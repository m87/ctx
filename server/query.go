package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/m87/ctx/core"
)

type queryRequest struct {
	WorkspaceId string `json:"workspaceId"`
	Query       string `json:"query"`
}

type QueryHandler struct {
	manager *core.ContextManager
}

func registerQueryHandler(mux *http.ServeMux, manager *core.ContextManager) {
	handler := &QueryHandler{manager: manager}
	mux.HandleFunc("POST /", handler.queryContexts)
	mux.HandleFunc("POST /result", handler.queryContextResult)
	registerSavedQueryHandler(mux, manager)
}

func (h *QueryHandler) queryContextResult(w http.ResponseWriter, r *http.Request) {
	var request queryRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", "Invalid request body")
		return
	}
	if request.WorkspaceId == "" {
		writeError(w, http.StatusBadRequest, "MISSING_WORKSPACE_ID", "Missing workspace ID")
		return
	}

	result, err := h.manager.QueryWorkspaceContexts(request.WorkspaceId, request.Query)
	if err != nil {
		var invalidQuery *core.InvalidContextQueryError
		if errors.As(err, &invalidQuery) {
			writeError(w, http.StatusBadRequest, "INVALID_QUERY", invalidQuery.Error())
			return
		}
		var workspaceNotFound *core.WorkspaceNotFoundError
		if errors.As(err, &workspaceNotFound) {
			writeError(w, http.StatusNotFound, "WORKSPACE_NOT_FOUND", "Workspace not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "FAILED_TO_QUERY_CONTEXTS", "Failed to query contexts")
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *QueryHandler) queryContexts(w http.ResponseWriter, r *http.Request) {
	var request queryRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", "Invalid request body")
		return
	}

	contexts, err := h.manager.QueryContexts(request.Query)
	if err != nil {
		var invalidQuery *core.InvalidContextQueryError
		if errors.As(err, &invalidQuery) {
			writeError(w, http.StatusBadRequest, "INVALID_QUERY", invalidQuery.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "FAILED_TO_QUERY_CONTEXTS", "Failed to query contexts")
		return
	}

	writeJSON(w, http.StatusOK, contexts)
}

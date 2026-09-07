package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/m87/ctx/core"
)

type SavedQueryHandler struct {
	manager *core.ContextManager
}

func registerSavedQueryHandler(mux *http.ServeMux, manager *core.ContextManager) {
	handler := &SavedQueryHandler{manager: manager}
	mux.HandleFunc("GET /saved", handler.listSavedQueries)
	mux.HandleFunc("GET /saved/", handler.listSavedQueries)
	mux.HandleFunc("POST /saved", handler.createSavedQuery)
	mux.HandleFunc("POST /saved/", handler.createSavedQuery)
	mux.HandleFunc("GET /saved/{id}", handler.getSavedQuery)
	mux.HandleFunc("DELETE /saved/{id}", handler.deleteSavedQuery)
}

func (h *SavedQueryHandler) listSavedQueries(w http.ResponseWriter, r *http.Request) {
	workspaceId := strings.TrimSpace(r.URL.Query().Get("workspaceId"))
	if workspaceId == "" {
		writeError(w, http.StatusBadRequest, "MISSING_WORKSPACE_ID", "Missing workspace ID")
		return
	}

	queries, err := h.manager.ListSavedQueries(workspaceId)
	if err != nil {
		var workspaceNotFound *core.WorkspaceNotFoundError
		if errors.As(err, &workspaceNotFound) {
			writeError(w, http.StatusNotFound, "WORKSPACE_NOT_FOUND", "Workspace not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "FAILED_TO_LIST_SAVED_QUERIES", "Failed to list saved queries")
		return
	}
	writeJSON(w, http.StatusOK, queries)
}

func (h *SavedQueryHandler) createSavedQuery(w http.ResponseWriter, r *http.Request) {
	var query core.SavedQuery
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", "Invalid request body")
		return
	}
	if strings.TrimSpace(query.WorkspaceId) == "" {
		writeError(w, http.StatusBadRequest, "MISSING_WORKSPACE_ID", "Missing workspace ID")
		return
	}
	if strings.TrimSpace(query.Name) == "" {
		writeError(w, http.StatusBadRequest, "MISSING_SAVED_QUERY_NAME", "Missing saved query name")
		return
	}
	if strings.TrimSpace(query.Query) == "" {
		writeError(w, http.StatusBadRequest, "MISSING_QUERY", "Missing query")
		return
	}

	id, err := h.manager.CreateSavedQuery(&query)
	if err != nil {
		var workspaceNotFound *core.WorkspaceNotFoundError
		if errors.As(err, &workspaceNotFound) {
			writeError(w, http.StatusNotFound, "WORKSPACE_NOT_FOUND", "Workspace not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "FAILED_TO_CREATE_SAVED_QUERY", "Failed to create saved query")
		return
	}
	query.Id = id
	writeJSON(w, http.StatusCreated, &query)
}

func (h *SavedQueryHandler) getSavedQuery(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeError(w, http.StatusBadRequest, "MISSING_SAVED_QUERY_ID", "Missing saved query ID")
		return
	}

	query, err := h.manager.GetSavedQuery(id)
	if err != nil {
		var notFound *core.SavedQueryNotFoundError
		if errors.As(err, &notFound) {
			writeError(w, http.StatusNotFound, "SAVED_QUERY_NOT_FOUND", "Saved query not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "FAILED_TO_GET_SAVED_QUERY", "Failed to get saved query")
		return
	}
	writeJSON(w, http.StatusOK, query)
}

func (h *SavedQueryHandler) deleteSavedQuery(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeError(w, http.StatusBadRequest, "MISSING_SAVED_QUERY_ID", "Missing saved query ID")
		return
	}

	if err := h.manager.DeleteSavedQuery(id); err != nil {
		var notFound *core.SavedQueryNotFoundError
		if errors.As(err, &notFound) {
			writeError(w, http.StatusNotFound, "SAVED_QUERY_NOT_FOUND", "Saved query not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "FAILED_TO_DELETE_SAVED_QUERY", "Failed to delete saved query")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

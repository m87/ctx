package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/m87/ctx/core"
)

type queryRequest struct {
	Query string `json:"query"`
}

type QueryHandler struct {
	manager *core.ContextManager
}

func registerQueryHandler(mux *http.ServeMux, manager *core.ContextManager) {
	handler := &QueryHandler{manager: manager}
	mux.HandleFunc("POST /", handler.queryContexts)
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

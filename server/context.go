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
		http.Error(w, "Missing context ID", http.StatusBadRequest)
		return
	}
	dateStr := r.PathValue("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Invalid date format, expected YYYY-MM-DD", http.StatusBadRequest)
		return
	}
	stats, err := h.manager.GetStats(id, date)
	if err != nil {
		http.Error(w, "Failed to get stats", http.StatusInternalServerError)
		return
	}

	writeJson(w, http.StatusOK, stats)
}

func (h *ContextHandler) listIntervals(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		http.Error(w, "Missing context ID", http.StatusBadRequest)
		return
	}
	intervals, err := h.manager.IntervalRepository.ListByContextId(id)
	if err != nil {
		http.Error(w, "Failed to list intervals", http.StatusInternalServerError)
		return
	}

	writeJson(w, http.StatusOK, intervals)
}

func (h *ContextHandler) getActiveContext(w http.ResponseWriter, r *http.Request) {
	activeContext, err := h.manager.ContextRepository.GetActive()
	if err != nil {
		http.Error(w, "Failed to get active context", http.StatusInternalServerError)
		return
	}
	if activeContext == nil {
		http.Error(w, "No active context found", http.StatusNotFound)
		return
	}
	writeJson(w, http.StatusOK, activeContext)
}

func (h *ContextHandler) switchContext(w http.ResponseWriter, r *http.Request) {

	var req *core.Context
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.manager.SwitchContext(req); err != nil {
		http.Error(w, "Failed to switch context", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ContextHandler) freeContext(w http.ResponseWriter, r *http.Request) {
	if err := h.manager.FreeActiveContext(); err != nil {
		http.Error(w, "Failed to free active context", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ContextHandler) listContexts(w http.ResponseWriter, r *http.Request) {
	contexts, err := h.manager.ContextRepository.List()
	if err != nil {
		http.Error(w, "Failed to list contexts", http.StatusInternalServerError)
		return
	}
	writeJson(w, http.StatusOK, contexts)
}

func (h *ContextHandler) createContext(w http.ResponseWriter, r *http.Request) {
	var context core.Context
	if err := json.NewDecoder(r.Body).Decode(&context); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	context.Id = ""

	id, err := h.manager.ContextRepository.Save(&context)
	if err != nil {
		http.Error(w, "Failed to create context", http.StatusInternalServerError)
		return
	}
	context.Id = id
	writeJson(w, http.StatusCreated, &context)
}

func (h *ContextHandler) deleteContext(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		http.Error(w, "Missing context ID", http.StatusBadRequest)
		return
	}
	if err := h.manager.ContextRepository.Delete(id); err != nil {
		http.Error(w, "Failed to delete context", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ContextHandler) getContext(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		http.Error(w, "Missing context ID", http.StatusBadRequest)
		return
	}
	context, err := h.manager.ContextRepository.GetById(id)
	if err != nil {
		http.Error(w, "Failed to get context", http.StatusInternalServerError)
		return
	}
	if context == nil {
		http.Error(w, "Context not found", http.StatusNotFound)
		return
	}
	writeJson(w, http.StatusOK, context)
}

func (h *ContextHandler) updateContext(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		http.Error(w, "Missing context ID", http.StatusBadRequest)
		return
	}
	var context core.Context
	if err := json.NewDecoder(r.Body).Decode(&context); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	context.Id = id
	if _, err := h.manager.ContextRepository.Save(&context); err != nil {
		http.Error(w, "Failed to update context", http.StatusInternalServerError)
		return
	}
	writeJson(w, http.StatusOK, &context)
}

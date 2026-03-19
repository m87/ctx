package server

import (
	"encoding/json"
	"net/http"
	"strings"

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

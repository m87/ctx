package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/m87/ctx/core"
)

type WorkspaceHandler struct {
	manager *core.ContextManager
}

func registerWorkspaceHandler(mux *http.ServeMux, manager *core.ContextManager) {
	handler := &WorkspaceHandler{manager: manager}
	mux.HandleFunc("GET /", handler.listWorkspaces)
	mux.HandleFunc("POST /", handler.createWorkspace)
	mux.HandleFunc("DELETE /{id}", handler.deleteWorkspace)
	mux.HandleFunc("GET /{id}", handler.getWorkspace)
	mux.HandleFunc("PUT /{id}", handler.updateWorkspace)
}

func (h *WorkspaceHandler) listWorkspaces(w http.ResponseWriter, r *http.Request) {
	workspaces, err := h.manager.WorkspaceRepository.List()
	if err != nil {
		http.Error(w, "Failed to list workspaces", http.StatusInternalServerError)
		return
	}
	writeJson(w, http.StatusOK, workspaces)
}

func (h *WorkspaceHandler) createWorkspace(w http.ResponseWriter, r *http.Request) {
	var workspace core.Workspace
	if err := json.NewDecoder(r.Body).Decode(&workspace); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	workspace.Id = ""

	id, err := h.manager.WorkspaceRepository.Save(&workspace)
	if err != nil {
		http.Error(w, "Failed to create workspace", http.StatusInternalServerError)
		return
	}
	workspace.Id = id
	writeJson(w, http.StatusCreated, &workspace)
}

func (h *WorkspaceHandler) deleteWorkspace(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		http.Error(w, "Missing workspace ID", http.StatusBadRequest)
		return
	}

	if err := h.manager.DeleteWorkspace(id); err != nil {
		if _, ok := err.(*core.WorkspaceInUseError); ok {
			http.Error(w, "Cannot delete workspace because it is in use by one or more contexts", http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to delete workspace", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *WorkspaceHandler) getWorkspace(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		http.Error(w, "Missing workspace ID", http.StatusBadRequest)
		return
	}
	workspace, err := h.manager.WorkspaceRepository.GetById(id)
	if err != nil {
		http.Error(w, "Failed to get workspace", http.StatusInternalServerError)
		return
	}
	if workspace == nil {
		http.Error(w, "Workspace not found", http.StatusNotFound)
		return
	}
	writeJson(w, http.StatusOK, workspace)
}

func (h *WorkspaceHandler) updateWorkspace(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		http.Error(w, "Missing workspace ID", http.StatusBadRequest)
		return
	}
	var workspace core.Workspace
	if err := json.NewDecoder(r.Body).Decode(&workspace); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	workspace.Id = id
	if _, err := h.manager.WorkspaceRepository.Save(&workspace); err != nil {
		http.Error(w, "Failed to update workspace", http.StatusInternalServerError)
		return
	}
	writeJson(w, http.StatusOK, &workspace)
}

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
		writeError(w, http.StatusInternalServerError, "FAILED_TO_LIST_WORKSPACES", "Failed to list workspaces")
		return
	}
	writeJson(w, http.StatusOK, workspaces)
}

func (h *WorkspaceHandler) createWorkspace(w http.ResponseWriter, r *http.Request) {
	var workspace core.Workspace
	if err := json.NewDecoder(r.Body).Decode(&workspace); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", "Invalid request body")
		return
	}
	workspace.Id = ""

	id, err := h.manager.WorkspaceRepository.Save(&workspace)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "FAILED_TO_CREATE_WORKSPACE", "Failed to create workspace")
		return
	}
	workspace.Id = id
	writeJson(w, http.StatusCreated, &workspace)
}

func (h *WorkspaceHandler) deleteWorkspace(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeError(w, http.StatusBadRequest, "MISSING_WORKSPACE_ID", "Missing workspace ID")
		return
	}

	if err := h.manager.DeleteWorkspace(id); err != nil {
		if _, ok := err.(*core.WorkspaceInUseError); ok {
			writeError(w, http.StatusBadRequest, "WORKSPACE_IN_USE", "Cannot delete workspace because it is in use by one or more contexts")
			return
		}
		writeError(w, http.StatusInternalServerError, "FAILED_TO_DELETE_WORKSPACE", "Failed to delete workspace")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *WorkspaceHandler) getWorkspace(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeError(w, http.StatusBadRequest, "MISSING_WORKSPACE_ID", "Missing workspace ID")
		return
	}
	workspace, err := h.manager.WorkspaceRepository.GetById(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "FAILED_TO_GET_WORKSPACE", "Failed to get workspace")
		return
	}
	if workspace == nil {
		writeError(w, http.StatusNotFound, "WORKSPACE_NOT_FOUND", "Workspace not found")
		return
	}
	writeJson(w, http.StatusOK, workspace)
}

func (h *WorkspaceHandler) updateWorkspace(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeError(w, http.StatusBadRequest, "MISSING_WORKSPACE_ID", "Missing workspace ID")
		return
	}
	var workspace core.Workspace
	if err := json.NewDecoder(r.Body).Decode(&workspace); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", "Invalid request body")
		return
	}
	workspace.Id = id
	if _, err := h.manager.WorkspaceRepository.Save(&workspace); err != nil {
		writeError(w, http.StatusInternalServerError, "FAILED_TO_UPDATE_WORKSPACE", "Failed to update workspace")
		return
	}
	writeJson(w, http.StatusOK, &workspace)
}

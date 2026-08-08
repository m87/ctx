package server

import (
	"encoding/json"
	"net/http"

	"github.com/m87/ctx/core"
)

type ProjectHandler struct {
	manager *core.ContextManager
}

func registerProjectHandler(mux *http.ServeMux, manager *core.ContextManager) {
	handlder := &ProjectHandler{manager: manager}
	mux.HandleFunc("GET /", handlder.listProjects)
	mux.HandleFunc("POST /", handlder.createProject)
}

func (h *ProjectHandler) createProject(w http.ResponseWriter, r *http.Request) {
	var project core.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON body")
		return
	}

	if project.WorkspaceId == "" {
		writeError(w, http.StatusBadRequest, "MISSING_WORKSPACE_ID", "Missing workspace ID")
		return
	}

	if project.Name == "" {
		writeError(w, http.StatusBadRequest, "MISSING_PROJECT_NAME", "Missing project name")
		return
	}

	id, err := h.manager.CreateProject(&project)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "FAILED_TO_CREATE_PROJECT", "Failed to create project")
		return
	}

	project.Id = id
	writeJson(w, http.StatusCreated, project)	
}

func (h *ProjectHandler) listProjects(w http.ResponseWriter, r *http.Request) {
	workspaceId := r.URL.Query().Get("workspaceId")
	if workspaceId == "" {
		writeError(w, http.StatusBadRequest, "MISSING_WORKSPACE_ID", "Missing workspace ID")
		return
	}

	projects, err := h.manager.ProjectRepository.List(workspaceId)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "FAILED_TO_LIST_PROJECTS", "Failed to list projects")
		return
	}

	writeJson(w, http.StatusOK, projects)
}

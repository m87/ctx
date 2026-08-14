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
	mux.HandleFunc("GET /all", handlder.listAllProjects)
	mux.HandleFunc("GET /{id}", handlder.getProject)
	mux.HandleFunc("PUT /{id}", handlder.updateProject)
	mux.HandleFunc("POST /", handlder.createProject)
	mux.HandleFunc("DELETE /{id}", handlder.deleteProject)
	mux.HandleFunc("GET /projects", handlder.listRootProjects)
	mux.HandleFunc("GET /{id}/projects", handlder.listChildren)
	mux.HandleFunc("GET /{id}/contexts", handlder.listContexts)
}

func (h *ProjectHandler) listRootProjects(w http.ResponseWriter, r *http.Request) {
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

func (h *ProjectHandler) listAllProjects(w http.ResponseWriter, r *http.Request) {
	workspaceId := r.URL.Query().Get("workspaceId")
	if workspaceId == "" {
		writeError(w, http.StatusBadRequest, "MISSING_WORKSPACE_ID", "Missing workspace ID")
		return
	}

	projects, err := h.manager.ProjectRepository.ListIncludingArchived(workspaceId)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "FAILED_TO_LIST_PROJECTS", "Failed to list projects")
		return
	}

	writeJson(w, http.StatusOK, projects)
}

func (h *ProjectHandler) listContexts(w http.ResponseWriter, r *http.Request) {
	projectId := r.PathValue("id")
	if projectId == "" {
		writeError(w, http.StatusBadRequest, "MISSING_PROJECT_ID", "Missing project ID")
		return
	}

	contexts, err := h.manager.ContextRepository.ListByProject(projectId)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "FAILED_TO_LIST_CONTEXTS", "Failed to list contexts")
		return
	}

	writeJson(w, http.StatusOK, contexts)
}

func (h *ProjectHandler) listChildren(w http.ResponseWriter, r *http.Request) {
	projectId := r.PathValue("id")
	if projectId == "" {
		writeError(w, http.StatusBadRequest, "MISSING_PROJECT_ID", "Missing project ID")
		return
	}

	projects, err := h.manager.ProjectRepository.ListChildren(projectId)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "FAILED_TO_LIST_PROJECTS", "Failed to list projects")
		return
	}

	writeJson(w, http.StatusOK, projects)
}

func (h *ProjectHandler) deleteProject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "MISSING_PROJECT_ID", "Missing project ID")
		return
	}

	err := h.manager.DeleteProject(id)
	if err != nil {
		if _, ok := err.(*core.ProjectNotFoundError); ok {
			writeError(w, http.StatusNotFound, "PROJECT_NOT_FOUND", err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "FAILED_TO_DELETE_PROJECT", "Failed to delete project")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProjectHandler) getProject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "MISSING_PROJECT_ID", "Missing project ID")
		return
	}

	project, err := h.manager.ProjectRepository.GetById(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "FAILED_TO_GET_PROJECT", "Failed to get project")
		return
	}
	if project == nil {
		writeError(w, http.StatusNotFound, "PROJECT_NOT_FOUND", "Project not found")
		return
	}

	writeJson(w, http.StatusOK, project)
}

func (h *ProjectHandler) updateProject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "MISSING_PROJECT_ID", "Missing project ID")
		return
	}

	var project core.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON body")
		return
	}

	project.Id = id

	if err := h.manager.UpdateProject(&project); err != nil {
		if _, ok := err.(*core.ProjectNotFoundError); ok {
			writeError(w, http.StatusBadRequest, "PROJECT_NOT_FOUND", err.Error())
			return
		}
		if _, ok := err.(*core.ProjectWorkspaceMismatchError); ok {
			writeError(w, http.StatusBadRequest, "PROJECT_WORKSPACE_MISMATCH", err.Error())
			return
		}
		if _, ok := err.(*core.ProjectHierarchyCycleError); ok {
			writeError(w, http.StatusBadRequest, "PROJECT_HIERARCHY_CYCLE", err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "FAILED_TO_UPDATE_PROJECT", "Failed to update project")
		return
	}

	writeJson(w, http.StatusOK, project)
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
		if _, ok := err.(*core.WorkspaceNotFoundError); ok {
			writeError(w, http.StatusBadRequest, "WORKSPACE_NOT_FOUND", err.Error())
			return
		}
		if _, ok := err.(*core.ProjectNotFoundError); ok {
			writeError(w, http.StatusBadRequest, "PROJECT_NOT_FOUND", err.Error())
			return
		}
		if _, ok := err.(*core.ProjectWorkspaceMismatchError); ok {
			writeError(w, http.StatusBadRequest, "PROJECT_WORKSPACE_MISMATCH", err.Error())
			return
		}
		if _, ok := err.(*core.ProjectHierarchyCycleError); ok {
			writeError(w, http.StatusBadRequest, "PROJECT_HIERARCHY_CYCLE", err.Error())
			return
		}
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

package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/m87/ctx/core"
)

type DashboardHandler struct {
	manager *core.ContextManager
}

func registerDashboardHandler(mux *http.ServeMux, manager *core.ContextManager) {
	handler := &DashboardHandler{manager: manager}
	mux.HandleFunc("GET /", handler.listDashboards)
	mux.HandleFunc("POST /", handler.createDashboard)
	mux.HandleFunc("GET /{id}", handler.getDashboard)
	mux.HandleFunc("PUT /{id}", handler.updateDashboard)
	mux.HandleFunc("DELETE /{id}", handler.deleteDashboard)
}

func (h *DashboardHandler) listDashboards(w http.ResponseWriter, r *http.Request) {
	workspaceId := strings.TrimSpace(r.URL.Query().Get("workspaceId"))
	if workspaceId == "" {
		writeError(w, http.StatusBadRequest, "MISSING_WORKSPACE_ID", "Missing workspace ID")
		return
	}

	dashboards, err := h.manager.ListDashboards(workspaceId)
	if err != nil {
		writeDashboardError(w, err, "FAILED_TO_LIST_DASHBOARDS", "Failed to list dashboards")
		return
	}
	writeJSON(w, http.StatusOK, dashboards)
}

func (h *DashboardHandler) createDashboard(w http.ResponseWriter, r *http.Request) {
	var dashboard core.Dashboard
	if err := json.NewDecoder(r.Body).Decode(&dashboard); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", "Invalid request body")
		return
	}

	id, err := h.manager.CreateDashboard(&dashboard)
	if err != nil {
		writeDashboardError(w, err, "FAILED_TO_CREATE_DASHBOARD", "Failed to create dashboard")
		return
	}
	dashboard.Id = id
	writeJSON(w, http.StatusCreated, &dashboard)
}

func (h *DashboardHandler) getDashboard(w http.ResponseWriter, r *http.Request) {
	dashboard, err := h.manager.GetDashboard(r.PathValue("id"))
	if err != nil {
		writeDashboardError(w, err, "FAILED_TO_GET_DASHBOARD", "Failed to get dashboard")
		return
	}
	writeJSON(w, http.StatusOK, dashboard)
}

func (h *DashboardHandler) updateDashboard(w http.ResponseWriter, r *http.Request) {
	var dashboard core.Dashboard
	if err := json.NewDecoder(r.Body).Decode(&dashboard); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", "Invalid request body")
		return
	}
	dashboard.Id = strings.TrimSpace(r.PathValue("id"))
	if err := h.manager.UpdateDashboard(&dashboard); err != nil {
		writeDashboardError(w, err, "FAILED_TO_UPDATE_DASHBOARD", "Failed to update dashboard")
		return
	}
	writeJSON(w, http.StatusOK, &dashboard)
}

func (h *DashboardHandler) deleteDashboard(w http.ResponseWriter, r *http.Request) {
	if err := h.manager.DeleteDashboard(r.PathValue("id")); err != nil {
		writeDashboardError(w, err, "FAILED_TO_DELETE_DASHBOARD", "Failed to delete dashboard")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeDashboardError(w http.ResponseWriter, err error, fallbackCode string, fallbackDescription string) {
	var invalid *core.InvalidDashboardError
	if errors.As(err, &invalid) {
		writeError(w, http.StatusBadRequest, "INVALID_DASHBOARD", invalid.Error())
		return
	}
	var duplicate *core.DashboardAlreadyExistsError
	if errors.As(err, &duplicate) {
		writeError(w, http.StatusConflict, "DASHBOARD_ALREADY_EXISTS", duplicate.Error())
		return
	}
	var notFound *core.DashboardNotFoundError
	if errors.As(err, &notFound) {
		writeError(w, http.StatusNotFound, "DASHBOARD_NOT_FOUND", "Dashboard not found")
		return
	}
	var workspaceNotFound *core.WorkspaceNotFoundError
	if errors.As(err, &workspaceNotFound) {
		writeError(w, http.StatusNotFound, "WORKSPACE_NOT_FOUND", "Workspace not found")
		return
	}
	var projectNotFound *core.ProjectNotFoundError
	if errors.As(err, &projectNotFound) {
		writeError(w, http.StatusNotFound, "PROJECT_NOT_FOUND", "Project not found")
		return
	}
	var projectWorkspaceMismatch *core.ProjectWorkspaceMismatchError
	if errors.As(err, &projectWorkspaceMismatch) {
		writeError(w, http.StatusBadRequest, "PROJECT_WORKSPACE_MISMATCH", projectWorkspaceMismatch.Error())
		return
	}
	var contextNotFound *core.ContextNotFoundError
	if errors.As(err, &contextNotFound) {
		writeError(w, http.StatusNotFound, "CONTEXT_NOT_FOUND", "Context not found")
		return
	}
	writeError(w, http.StatusInternalServerError, fallbackCode, fallbackDescription)
}

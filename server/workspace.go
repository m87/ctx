package server

import (
	"net/http"

	"github.com/m87/ctx/core"
	"github.com/m87/nod"
)

type WorkspaceHandler struct {
	manager *core.ContextManager
}

func registerWorksapceHandler(mux *http.ServeMux, manager *core.ContextManager) {
	handler := &WorkspaceHandler{manager: manager}

	mux.HandleFunc("GET /", handler.listWorkspaces)
}

func (h *WorkspaceHandler) listWorkspaces(w http.ResponseWriter, r *http.Request) {
	h.manager.Execute(func(repository *nod.Repository) error {
		workspaces, err := repository.Query().TypeEquals(core.WorkspaceType).List()
		if err != nil {
			http.Error(w, "Failed to list workspaces", http.StatusInternalServerError)
			return err
		}

		var wsList []*core.Workspace
		for _, node := range workspaces {
			ws := node.(*core.Workspace)
			wsList = append(wsList, ws)
		}

		writeJson(w, http.StatusOK, wsList)
		return nil
	})
}

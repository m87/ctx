package server

import (
	"net/http"

	"github.com/m87/ctx/core"
)

type IntegrityHandler struct {
	manager *core.ContextManager
}

func registerIntegrityHandler(mux *http.ServeMux, manager *core.ContextManager) {
	handler := &IntegrityHandler{manager: manager}
	mux.HandleFunc("GET /", handler.checkIntegrity)
	mux.HandleFunc("GET /contexts", handler.listContexts)
	mux.HandleFunc("POST /repair", handler.repairIntegrity)
}

func (h *IntegrityHandler) listContexts(w http.ResponseWriter, _ *http.Request) {
	contexts, err := h.manager.ListIntegrityContextOptions()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "FAILED_TO_LIST_INTEGRITY_CONTEXTS", "Failed to list contexts for data integrity repair")
		return
	}
	writeJson(w, http.StatusOK, contexts)
}

func (h *IntegrityHandler) repairIntegrity(w http.ResponseWriter, _ *http.Request) {
	result, err := h.manager.RepairIntegrity()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "FAILED_TO_REPAIR_INTEGRITY", "Failed to repair data integrity")
		return
	}
	writeJson(w, http.StatusOK, result)
}

func (h *IntegrityHandler) checkIntegrity(w http.ResponseWriter, _ *http.Request) {
	report, err := h.manager.CheckIntegrity()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "FAILED_TO_CHECK_INTEGRITY", "Failed to check data integrity")
		return
	}
	writeJson(w, http.StatusOK, report)
}

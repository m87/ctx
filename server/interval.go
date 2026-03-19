package server

import (
	"net/http"

	"github.com/m87/ctx/core"
)

type IntervalHandler struct {
	manager *core.ContextManager
}

func registerIntervalHandler(mux *http.ServeMux, manager *core.ContextManager) {
	handler := &IntervalHandler{manager: manager}
	mux.HandleFunc("GET /", handler.listIntervals)
}

func (h *IntervalHandler) listIntervals(w http.ResponseWriter, r *http.Request) {

}

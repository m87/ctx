package server

import (
	"net/http"

	"github.com/m87/ctx/core"
)

func (s *Server) version(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, VersionResponse{Version: core.Version})
}

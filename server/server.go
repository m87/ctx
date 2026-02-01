package server

import (
	"net/http"

	"github.com/m87/ctx/core"
	ctxlog "github.com/m87/ctx/log"
)

type Server struct {
	Manager *core.ContextManager
	mux     *http.ServeMux
}

func NewServer(manager *core.ContextManager) *Server {
	s := &Server{
		Manager: manager,
		mux:     http.NewServeMux(),
	}

	workspaceMux := http.NewServeMux()
	registerWorksapceHandler(workspaceMux, manager)
	s.mux.Handle("/workspaces/", http.StripPrefix("/workspaces", workspaceMux))

	return s

}

func (s *Server) Handler() http.Handler {
	var h http.Handler = s.mux
	h = withLogging(h)
	return h
}

func (s *Server) Listen(addr string) error {
	ctxlog.Logger.Info("Starting server on ", addr)
	return http.ListenAndServe(addr, s.Handler())
}

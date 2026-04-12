package server

import (
	"net/http"
	"strings"

	"github.com/m87/ctx/core"
	ctxlog "github.com/m87/ctx/log"
)

type Server struct {
	Manager *core.ContextManager
	mux     *http.ServeMux
	spa     http.Handler
}

func NewServer(manager *core.ContextManager) *Server {
	s := &Server{
		Manager: manager,
		mux:     http.NewServeMux(),
	}

	s.spa = registerSpaHandler()
	registerApiRoutes(s.mux, manager)
	registerLegacyRoutes(s.mux, manager)

	return s

}

func registerApiRoutes(mux *http.ServeMux, manager *core.ContextManager) {
	apiMux := http.NewServeMux()

	contextMux := http.NewServeMux()
	registerContextHandler(contextMux, manager)
	apiMux.Handle("/context/", http.StripPrefix("/context", contextMux))
	apiMux.Handle("/context", http.StripPrefix("/context", contextMux))

	intervalMux := http.NewServeMux()
	registerIntervalHandler(intervalMux, manager)
	apiMux.Handle("/interval/", http.StripPrefix("/interval", intervalMux))
	apiMux.Handle("/interval", http.StripPrefix("/interval", intervalMux))

	mux.Handle("/api/", http.StripPrefix("/api", apiMux))
	mux.Handle("/api", http.StripPrefix("/api", apiMux))
}

func registerLegacyRoutes(mux *http.ServeMux, manager *core.ContextManager) {
	contextMux := http.NewServeMux()
	registerContextHandler(contextMux, manager)
	mux.Handle("/context/", http.StripPrefix("/context", contextMux))
	mux.Handle("/context", http.StripPrefix("/context", contextMux))

	intervalMux := http.NewServeMux()
	registerIntervalHandler(intervalMux, manager)
	mux.Handle("/interval/", http.StripPrefix("/interval", intervalMux))
	mux.Handle("/interval", http.StripPrefix("/interval", intervalMux))
}

func (s *Server) Handler() http.Handler {
	var base http.Handler = s.mux

	if s.spa == nil {
		return withLogging(base)
	}

	top := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			base.ServeHTTP(w, r)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/api") || strings.HasPrefix(r.URL.Path, "/context") || strings.HasPrefix(r.URL.Path, "/interval") {
			base.ServeHTTP(w, r)
			return
		}

		s.spa.ServeHTTP(w, r)
	})

	return withLogging(top)
}

func (s *Server) Listen(addr string) error {
	ctxlog.Logger.Info("Starting server on " + addr)
	return http.ListenAndServe(addr, s.Handler())
}

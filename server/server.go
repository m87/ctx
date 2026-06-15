package server

import (
	"net/http"
	"strings"

	"github.com/m87/ctx/core"
	ctxlog "github.com/m87/ctx/log"
)

type Server struct {
	Manager         *core.ContextManager
	SettingsManager *core.SettingsManager
	mux             *http.ServeMux
	spa             http.Handler
}

func NewServer(manager *core.ContextManager, settingsManager *core.SettingsManager) *Server {
	s := &Server{
		Manager:         manager,
		SettingsManager: settingsManager,
		mux:             http.NewServeMux(),
	}

	s.spa = registerSpaHandler()
	registerApiRoutes(s.mux, manager, settingsManager)
	registerLegacyRoutes(s.mux, manager, settingsManager)

	return s

}

func registerApiRoutes(mux *http.ServeMux, manager *core.ContextManager, settingsManager *core.SettingsManager) {
	apiMux := http.NewServeMux()

	versionMux := http.NewServeMux()
	registerVersionHandler(versionMux)
	apiMux.Handle("/version/", stripPrefixOrRoot("/version", versionMux))
	apiMux.Handle("/version", stripPrefixOrRoot("/version", versionMux))

	settingsMux := http.NewServeMux()
	registerSettingsHandler(settingsMux, settingsManager)
	apiMux.Handle("/settings/", stripPrefixOrRoot("/settings", settingsMux))
	apiMux.Handle("/settings", stripPrefixOrRoot("/settings", settingsMux))

	integrityMux := http.NewServeMux()
	registerIntegrityHandler(integrityMux, manager)
	apiMux.Handle("/integrity/", stripPrefixOrRoot("/integrity", integrityMux))
	apiMux.Handle("/integrity", stripPrefixOrRoot("/integrity", integrityMux))

	contextMux := http.NewServeMux()
	registerContextHandler(contextMux, manager)
	apiMux.Handle("/context/", http.StripPrefix("/context", contextMux))
	apiMux.Handle("/context", http.StripPrefix("/context", contextMux))

	workspaceMux := http.NewServeMux()
	registerWorkspaceHandler(workspaceMux, manager)
	apiMux.Handle("/workspace/", http.StripPrefix("/workspace", workspaceMux))
	apiMux.Handle("/workspace", http.StripPrefix("/workspace", workspaceMux))

	intervalMux := http.NewServeMux()
	registerIntervalHandler(intervalMux, manager)
	apiMux.Handle("/interval/", http.StripPrefix("/interval", intervalMux))
	apiMux.Handle("/interval", http.StripPrefix("/interval", intervalMux))

	mux.Handle("/api/", http.StripPrefix("/api", apiMux))
	mux.Handle("/api", http.StripPrefix("/api", apiMux))
}

func registerLegacyRoutes(mux *http.ServeMux, manager *core.ContextManager, settingsManager *core.SettingsManager) {
	versionMux := http.NewServeMux()
	registerVersionHandler(versionMux)
	mux.Handle("/version/", stripPrefixOrRoot("/version", versionMux))
	mux.Handle("/version", stripPrefixOrRoot("/version", versionMux))

	settingsMux := http.NewServeMux()
	registerSettingsHandler(settingsMux, settingsManager)
	mux.Handle("/settings/", stripPrefixOrRoot("/settings", settingsMux))
	mux.Handle("/settings", stripPrefixOrRoot("/settings", settingsMux))

	integrityMux := http.NewServeMux()
	registerIntegrityHandler(integrityMux, manager)
	mux.Handle("/integrity/", stripPrefixOrRoot("/integrity", integrityMux))
	mux.Handle("/integrity", stripPrefixOrRoot("/integrity", integrityMux))

	contextMux := http.NewServeMux()
	registerContextHandler(contextMux, manager)
	mux.Handle("/context/", http.StripPrefix("/context", contextMux))
	mux.Handle("/context", http.StripPrefix("/context", contextMux))

	workspaceMux := http.NewServeMux()
	registerWorkspaceHandler(workspaceMux, manager)
	mux.Handle("/workspace/", http.StripPrefix("/workspace", workspaceMux))
	mux.Handle("/workspace", http.StripPrefix("/workspace", workspaceMux))

	intervalMux := http.NewServeMux()
	registerIntervalHandler(intervalMux, manager)
	mux.Handle("/interval/", http.StripPrefix("/interval", intervalMux))
	mux.Handle("/interval", http.StripPrefix("/interval", intervalMux))
}

func stripPrefixOrRoot(prefix string, handler http.Handler) http.Handler {
	return http.StripPrefix(prefix, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "" {
			r.URL.Path = "/"
		}
		handler.ServeHTTP(w, r)
	}))
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
		if strings.HasPrefix(r.URL.Path, "/api") || strings.HasPrefix(r.URL.Path, "/context") || strings.HasPrefix(r.URL.Path, "/workspace") || strings.HasPrefix(r.URL.Path, "/interval") || strings.HasPrefix(r.URL.Path, "/version") || strings.HasPrefix(r.URL.Path, "/settings") || strings.HasPrefix(r.URL.Path, "/integrity") {
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

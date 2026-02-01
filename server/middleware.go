package server

import (
	"net/http"
	"time"

	ctxlog "github.com/m87/ctx/log"
)

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t0 := time.Now()
		next.ServeHTTP(w, r)
		ctxlog.Logger.Info(r.Method, r.URL.Path, time.Since(t0))
	})
}

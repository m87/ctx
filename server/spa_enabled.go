//go:build allinone

package server

import (
	"embed"
	"io/fs"
	"mime"
	"net/http"
	"path"
	"strings"

	ctxlog "github.com/m87/ctx/log"
)

//go:embed all:spa/dist
var spaAssets embed.FS

func registerSpaHandler() http.Handler {
	distFS, err := fs.Sub(spaAssets, "spa/dist")
	if err != nil {
		ctxlog.Logger.Warn("Failed to mount SPA assets", "error", err)
		return nil
	}

	if _, err := fs.Stat(distFS, "index.html"); err != nil {
		ctxlog.Logger.Warn("SPA assets missing index.html; skipping SPA handler")
		return nil
	}

	return spaHandler(distFS)
}

func spaHandler(distFS fs.FS) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assetPath := path.Clean(strings.TrimPrefix(r.URL.Path, "/"))

		if assetPath == "." || assetPath == "/" || assetPath == "index.html" {
			serveSpaIndex(distFS, w)
			return
		}

		if strings.HasPrefix(assetPath, "../") || assetPath == ".." {
			http.NotFound(w, r)
			return
		}

		if info, err := fs.Stat(distFS, assetPath); err == nil && !info.IsDir() {
			serveAsset(distFS, assetPath, w)
			return
		}

		serveSpaIndex(distFS, w)
	})
}

func serveSpaIndex(distFS fs.FS, w http.ResponseWriter) {
	content, err := fs.ReadFile(distFS, "index.html")
	if err != nil {
		http.Error(w, "Failed to load UI", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(content)
}

func serveAsset(distFS fs.FS, assetPath string, w http.ResponseWriter) {
	content, err := fs.ReadFile(distFS, assetPath)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if contentType := mime.TypeByExtension(path.Ext(assetPath)); contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(content)
}

//go:build allinone

package server

import (
	"embed"
	"io/fs"
	"mime"
	"net/http"
	"path"
	"strconv"
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
			serveSpaIndex(distFS, w, r.Method == http.MethodHead)
			return
		}

		if strings.HasPrefix(assetPath, "../") || assetPath == ".." {
			writeError(w, http.StatusNotFound, "ASSET_NOT_FOUND", "Asset not found")
			return
		}

		if info, err := fs.Stat(distFS, assetPath); err == nil && !info.IsDir() {
			serveAsset(distFS, assetPath, w, r.Method == http.MethodHead)
			return
		}

		serveSpaIndex(distFS, w, r.Method == http.MethodHead)
	})
}

func serveSpaIndex(distFS fs.FS, w http.ResponseWriter, head bool) {
	content, err := fs.ReadFile(distFS, "index.html")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "FAILED_TO_LOAD_UI", "Failed to load UI")
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))
	w.WriteHeader(http.StatusOK)
	if !head {
		_, _ = w.Write(content)
	}
}

func serveAsset(distFS fs.FS, assetPath string, w http.ResponseWriter, head bool) {
	content, err := fs.ReadFile(distFS, assetPath)
	if err != nil {
		writeError(w, http.StatusNotFound, "ASSET_NOT_FOUND", "Asset not found")
		return
	}
	if contentType := mime.TypeByExtension(path.Ext(assetPath)); contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))
	w.WriteHeader(http.StatusOK)
	if !head {
		_, _ = w.Write(content)
	}
}

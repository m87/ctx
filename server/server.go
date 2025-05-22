package server

import (
	"bytes"
	"embed"
	"encoding/json"
	"io"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/m87/ctx/ctx"
)

//go:embed ui/ctx-dashboard/dist/ctx-dashboard/assets/*
//go:embed ui/ctx-dashboard/dist/ctx-dashboard/*.html
//go:embed ui/ctx-dashboard/dist/ctx-dashboard/*.ico
var staticFiles embed.FS

func spaHandler(content fs.FS, fsHandler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request URL:", r.URL.Path)
		_, err := fs.Stat(content, r.URL.Path[1:])
		if err != nil {
			f, err := content.Open("index.html")
			if err != nil {
				http.Error(w, "index.html not found", http.StatusInternalServerError)
				return
			}
			defer f.Close()

			data, err := io.ReadAll(f)
			if err != nil {
				http.Error(w, "error reading index.html", http.StatusInternalServerError)
				return
			}

			http.ServeContent(w, r, "index.html", time.Time{}, bytes.NewReader(data))
			return
		}

		fsHandler.ServeHTTP(w, r)
	}
}

func contextList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	mgr := ctx.CreateManager()

	json.NewEncoder(w).Encode(mgr.ListJson2())
}

func Serve() {
	content, err := fs.Sub(staticFiles, "ui/ctx-dashboard/dist/ctx-dashboard")
	if err != nil {
		log.Fatal(err)
	}
	fs := http.FileServer(http.FS(content))

	http.Handle("/", spaHandler(content, fs))
	http.HandleFunc("/api/context/list", contextList)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

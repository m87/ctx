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
	"github.com/m87/ctx/ctx_model"
	"github.com/m87/ctx/util"
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

func currentContext(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	mgr := ctx.CreateManager()

	mgr.ContextStore.Read(func(s *ctx_model.State) error {
		if s.CurrentId != "" {
			json.NewEncoder(w).Encode(s.Contexts[s.CurrentId])
		} else {
			json.NewEncoder(w).Encode(nil)
		}
		return nil
	})

}

func createAndSwitchContext(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	mgr := ctx.CreateManager()

	var p createAndSwitchRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	mgr.CreateIfNotExistsAndSwitch(util.GenerateId(p.Description), p.Description)
}

func Serve() {
	content, err := fs.Sub(staticFiles, "ui/ctx-dashboard/dist/ctx-dashboard")
	if err != nil {
		log.Fatal(err)
	}
	fs := http.FileServer(http.FS(content))

	http.Handle("/", spaHandler(content, fs))
	http.HandleFunc("/api/context/list", contextList)
	http.HandleFunc("/api/context/current", currentContext)
	http.HandleFunc("/api/context/free", freeContext)
	http.HandleFunc("/api/context/switch", switchContext)
	http.HandleFunc("/api/context/createAndSwitch", createAndSwitchContext)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type SwitchRequest struct {
	Id string `json:"id"`
}

type createAndSwitchRequest struct {
	Description string `json:"description"`
}

func switchContext(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var p SwitchRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	mgr := ctx.CreateManager()
	mgr.Switch(p.Id)
}

func freeContext(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	mgr := ctx.CreateManager()
	mgr.Free()
}

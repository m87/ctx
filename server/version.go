package server

import "net/http"

var Release = "dev"
var Commit = ""
var Date = ""

type VersionInfo struct {
	Version string `json:"version"`
	Release string `json:"release"`
	Commit  string `json:"commit,omitempty"`
	Date    string `json:"date,omitempty"`
}

func CurrentVersion() VersionInfo {
	return VersionInfo{
		Version: Release,
		Release: Release,
		Commit:  Commit,
		Date:    Date,
	}
}

func registerVersionHandler(mux *http.ServeMux) {
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		writeJson(w, http.StatusOK, CurrentVersion())
	})
}

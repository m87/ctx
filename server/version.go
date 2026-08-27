package server

import (
	"net/http"

	"github.com/m87/ctx/core"
)

var Release = "dev"
var Commit = ""
var Date = ""

type VersionInfo struct {
	Version   string `json:"version"`
	Release   string `json:"release"`
	Commit    string `json:"commit,omitempty"`
	Date      string `json:"date,omitempty"`
	DBVersion string `json:"dbVersion,omitempty"`
}

func CurrentVersion() VersionInfo {
	return VersionInfo{
		Version:   Release,
		Release:   Release,
		Commit:    Commit,
		Date:      Date,
		DBVersion: core.CurrentDatabaseVersion,
	}
}

func registerVersionHandler(mux *http.ServeMux) {
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, CurrentVersion())
	})
}

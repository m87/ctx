package server

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

func writeJson(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, code, description string) {
	writeJson(w, status, ErrorResponse{
		Code:        code,
		Description: description,
	})
}

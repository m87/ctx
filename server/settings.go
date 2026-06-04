package server

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/m87/ctx/core"
)

func registerSettingsHandler(mux *http.ServeMux, manager *core.SettingsManager) {
	handler := &SettingsHandler{manager: manager}
	mux.HandleFunc("GET /key/{key}", handler.getClientSetting)
	mux.HandleFunc("GET /", handler.getClientSettings)
	mux.HandleFunc("PATCH /", handler.saveClientSettings)
}

type SettingsHandler struct {
	manager *core.SettingsManager
}

func (h *SettingsHandler) getClientSetting(w http.ResponseWriter, r *http.Request) {
	raw := r.PathValue("key")
	key, err := url.QueryUnescape(raw)
	if err != nil {
		http.Error(w, "Invalid setting key: "+err.Error(), http.StatusBadRequest)
		return
	}
	if key == "" {
		http.Error(w, "Missing setting key", http.StatusBadRequest)
		return
	}

	if !strings.HasPrefix(key, "client.") {
		http.Error(w, "Invalid setting key", http.StatusBadRequest)
		return
	}

	value, err := h.manager.GetClientKey(key)
	if err != nil {
		http.Error(w, "Error retrieving setting: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(value))
}

func (h *SettingsHandler) getClientSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := h.manager.GetClient()
	if err != nil {
		http.Error(w, "Error retrieving client settings: "+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(settings)
}

func (h *SettingsHandler) saveClientSettings(w http.ResponseWriter, r *http.Request) {
	var settings map[string]string
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.manager.SaveClient(settings); err != nil {
		http.Error(w, "Error saving settings: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

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
		writeError(w, http.StatusBadRequest, "INVALID_SETTING_KEY", "Invalid setting key")
		return
	}
	if key == "" {
		writeError(w, http.StatusBadRequest, "MISSING_SETTING_KEY", "Missing setting key")
		return
	}

	if !strings.HasPrefix(key, "client.") {
		writeError(w, http.StatusBadRequest, "INVALID_SETTING_KEY", "Invalid setting key")
		return
	}

	value, err := h.manager.GetClientKey(key)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "FAILED_TO_GET_SETTING", "Failed to get setting")
		return
	}

	w.Write([]byte(value))
}

func (h *SettingsHandler) getClientSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := h.manager.GetClient()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "FAILED_TO_GET_SETTINGS", "Failed to get settings")
		return
	}
	json.NewEncoder(w).Encode(settings)
}

func (h *SettingsHandler) saveClientSettings(w http.ResponseWriter, r *http.Request) {
	var settings map[string]string
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", "Invalid request body")
		return
	}

	if err := h.manager.SaveClient(settings); err != nil {
		writeError(w, http.StatusInternalServerError, "FAILED_TO_SAVE_SETTINGS", "Failed to save settings")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

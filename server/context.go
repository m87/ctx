package server

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"github.com/m87/ctx/core"
	"github.com/m87/ctx/util"
)

type contextHandlers struct{ mgr *core.ContextManager }

func registerContext(mux *http.ServeMux, mgr *core.ContextManager) {
	h := &contextHandlers{mgr: mgr}
	mux.HandleFunc("GET /list", h.list)
	mux.HandleFunc("GET /current", h.current)
	mux.HandleFunc("POST /free", h.free)
	mux.HandleFunc("POST /switch", h.switchContext)
	mux.HandleFunc("POST /createAndSwitch", h.createAndSwitch)
	mux.HandleFunc("PUT /interval", h.updateInterval)
	mux.HandleFunc("POST /rename", h.rename)
	mux.HandleFunc("DELETE /{ctxId}", h.delete)
	mux.HandleFunc("POST /labels", h.editLabels)
	mux.HandleFunc("POST /{ctxId}/comment", h.saveContextComment)
	mux.HandleFunc("DELETE /{ctxId}/comment/{commentId}", h.deleteContextComment)
}

func (h *contextHandlers) editLabels(w http.ResponseWriter, r *http.Request) {
	var p EditContextLabelsRequest
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	w.WriteHeader(http.StatusOK)
	h.mgr.WithSession(func(s core.Session) error {
		s.UpdateContextLabels(p.Id, p.Labels)
		return nil
	})
}

func (h *contextHandlers) delete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if ctxId := strings.TrimSpace(r.PathValue("ctxId")); ctxId != "" {
		h.mgr.WithSession(func(s core.Session) error { return s.Delete(ctxId) })
	}
}

func (h *contextHandlers) list(w http.ResponseWriter, r *http.Request) {
	h.mgr.WithSession(func(s core.Session) error {
		out := make([]core.Context, 0, len(s.State.Contexts))
		for _, c := range s.State.Contexts {
			out = append(out, c)
		}

		sort.Slice(out, func(i, j int) bool { return out[i].Description < out[j].Description })
		writeJSON(w, http.StatusOK, out)
		return nil
	})
}

func (h *contextHandlers) current(w http.ResponseWriter, r *http.Request) {
	res := CurrentContextResponse{CurrentDuration: 0}
	h.mgr.WithSession(func(s core.Session) error {
		if s.State.CurrentId == "" {
			writeJSON(w, http.StatusOK, nil)
			return nil
		}
		cur := s.MustGetCtx(s.State.CurrentId)
		res.Context = cur
		for _, iv := range cur.Intervals {
			if iv.End.Time.IsZero() {
				res.CurrentDuration = h.mgr.TimeProvider.Now().Time.Sub(iv.Start.Time)
			}
		}
		writeJSON(w, http.StatusOK, res)
		return nil
	})
}

func (h *contextHandlers) createAndSwitch(w http.ResponseWriter, r *http.Request) {
	var p createAndSwitchRequest
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	w.WriteHeader(http.StatusOK)
	h.mgr.WithArchiveSession(func(s core.Session) error {
		return s.CreateIfNotExistsAndSwitch(util.GenerateId(p.Description), p.Description)
	})
}

func (h *contextHandlers) rename(w http.ResponseWriter, r *http.Request) {
	var p RenameRequest
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	w.WriteHeader(http.StatusOK)
	h.mgr.WithSession(func(s core.Session) error {
		s.RenameContext(p.CtxId, util.GenerateId(p.Name), p.Name)
		return nil
	})
}

func (h *contextHandlers) updateInterval(w http.ResponseWriter, r *http.Request) {
	var p EditIntervalRequest
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	w.WriteHeader(http.StatusOK)
	h.mgr.WithSession(func(s core.Session) error {
		s.EditContextInterval(p.Id, p.IntervalId, p.Start, p.End)
		return nil
	})
}

func (h *contextHandlers) switchContext(w http.ResponseWriter, r *http.Request) {
	var p SwitchRequest
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	w.WriteHeader(http.StatusOK)
	h.mgr.WithSession(func(s core.Session) error { return s.Switch(p.Id) })
}

func (h *contextHandlers) free(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	h.mgr.WithSession(func(s core.Session) error { return s.Free() })
}

func (h *contextHandlers) saveContextComment(w http.ResponseWriter, r *http.Request) {
	var p core.Comment
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if ctxId := strings.TrimSpace(r.PathValue("ctxId")); ctxId != "" {
		w.WriteHeader(http.StatusOK)
		h.mgr.WithSession(func(s core.Session) error {
			return s.SaveContextComment(ctxId, p)
		})
	}
}

func (h *contextHandlers) deleteContextComment(w http.ResponseWriter, r *http.Request) {
	if ctxId := strings.TrimSpace(r.PathValue("ctxId")); ctxId != "" {
		if commentId := strings.TrimSpace(r.PathValue("commentId")); commentId != "" {
			w.WriteHeader(http.StatusOK)
			h.mgr.WithSession(func(s core.Session) error {
				return s.DeleteContextComment(ctxId, commentId)
			})
		}
	}
}

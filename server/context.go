package server

import (
	"encoding/json"
	"net/http"
	"sort"

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
	mux.HandleFunc("POST /interval", h.updateInterval)
	mux.HandleFunc("POST /rename", h.rename)
}

func (h *contextHandlers) list(w http.ResponseWriter, r *http.Request) {
	h.mgr.WithSession(func(s core.Session) error {
		out := make([]core.Context, 0, len(s.State.Contexts))
		for _, c := range s.State.Contexts {
			out = append(out, c)
		}
		// stabilna kolejność
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

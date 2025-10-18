package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/m87/ctx/core"
	ctxtime "github.com/m87/ctx/time"
)

type intervalsHandlers struct{ mgr *core.ContextManager }

func registerIntervals(mux *http.ServeMux, mgr *core.ContextManager) {
	h := &intervalsHandlers{mgr: mgr}
	mux.HandleFunc("GET /{date}", h.byDate)
	mux.HandleFunc("GET /", h.all)
	mux.HandleFunc("GET /recent/{n}", h.recent)
	mux.HandleFunc("POST /move", h.move)
	mux.HandleFunc("DELETE /{ctxId}/{id}", h.deleteOne)
	mux.HandleFunc("POST /{ctxId}/{id}/split", h.split) // POST, bo body
}

func (h *intervalsHandlers) all(w http.ResponseWriter, r *http.Request) {
	loc := getLoc()
	date := h.mgr.TimeProvider.Now().Time.In(loc)
	h.respondIntervals(w, date)
}

func (h *intervalsHandlers) byDate(w http.ResponseWriter, r *http.Request) {
	loc := getLoc()
	date := h.mgr.TimeProvider.Now().Time.In(loc)
	if raw := strings.TrimSpace(r.PathValue("date")); raw != "" {
		date, _ = time.ParseInLocation(time.DateOnly, raw, loc)
	}
	h.respondIntervals(w, date)
}

func (h *intervalsHandlers) recent(w http.ResponseWriter, r *http.Request) {
	loc := getLoc()
	date := h.mgr.TimeProvider.Now().Time.In(loc)

	n := 10
	if raw := strings.TrimSpace(r.PathValue("n")); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil {
			n = v
		}
	}
	resp := IntervalsResponse{}
	for i := 0; i < n; i++ {
		d := date.AddDate(0, 0, -i)
		entry, err := h.intervalsByDate(d)
		if err != nil {
			http.Error(w, "error fetching intervals for "+d.Format(time.DateOnly), http.StatusInternalServerError)
			return
		}
		resp.Days = append(resp.Days, entry)
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *intervalsHandlers) move(w http.ResponseWriter, r *http.Request) {
	var p MoveIntervalRequest
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	w.WriteHeader(http.StatusOK)
	h.mgr.WithArchiveSession(func(s core.Session) error {
		return s.MoveIntervalById(p.Src, p.Target, p.Id)
	})
}

func (h *intervalsHandlers) deleteOne(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	ctxID := strings.TrimSpace(r.PathValue("ctxId"))
	id := strings.TrimSpace(r.PathValue("id"))
	w.WriteHeader(http.StatusOK)
	h.mgr.WithSession(func(s core.Session) error { return s.DeleteInterval(ctxID, id) })
}

func (h *intervalsHandlers) split(w http.ResponseWriter, r *http.Request) {
	var p SplitRequest
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	ctxID := strings.TrimSpace(r.PathValue("ctxId"))
	id := strings.TrimSpace(r.PathValue("id"))
	w.WriteHeader(http.StatusOK)
	h.mgr.WithSession(func(s core.Session) error {
		return s.SplitContextIntervalById(ctxID, id, p.Split.H, p.Split.M, p.Split.S)
	})
}

func (h *intervalsHandlers) respondIntervals(w http.ResponseWriter, date time.Time) {
	entry, _ := h.intervalsByDate(date)
	writeJSON(w, http.StatusOK, IntervalsResponse{Days: []IntervalsResponseEntry{entry}})
}

func (h *intervalsHandlers) intervalsByDate(date time.Time) (IntervalsResponseEntry, error) {
	loc := getLoc()
	resp := IntervalsResponseEntry{Date: date.Format(time.DateOnly)}
	h.mgr.WithSession(func(s core.Session) error {
		for ctxID := range s.State.Contexts {
			ints := s.GetIntervalsByDate(ctxID, ctxtime.ZonedTime{Time: date, Timezone: loc.String()})
			for _, iv := range ints {
				resp.Intervals = append(resp.Intervals, IntervalEntry{
					Id:          iv.Id,
					CtxId:       ctxID,
					Description: s.State.Contexts[ctxID].Description,
					Interval:    iv,
				})
			}
		}
		return nil
	})
	return resp, nil
}

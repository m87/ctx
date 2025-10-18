package server

import (
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/m87/ctx/core"
	ctxtime "github.com/m87/ctx/time"
	"github.com/m87/ctx/util"
)

type summaryHandlers struct{ mgr *core.ContextManager }

func registerSummary(mux *http.ServeMux, mgr *core.ContextManager) {
	h := &summaryHandlers{mgr: mgr}
	mux.HandleFunc("GET /day/{date}", h.dayByDate)
	mux.HandleFunc("GET /day", h.dayToday)
	mux.HandleFunc("GET /day/list", h.dayList)
}

func (h *summaryHandlers) dayToday(w http.ResponseWriter, r *http.Request) {
	loc := getLoc()
	date := h.mgr.TimeProvider.Now().Time.In(loc)
	h.respondDay(w, date, r.URL.Query().Get("showAllContexts") == "true")
}

func (h *summaryHandlers) dayByDate(w http.ResponseWriter, r *http.Request) {
	loc := getLoc()
	date := h.mgr.TimeProvider.Now().Time.In(loc)
	if raw := strings.TrimSpace(r.PathValue("date")); raw != "" {
		date, _ = time.ParseInLocation(time.DateOnly, raw, loc)
	}
	h.respondDay(w, date, r.URL.Query().Get("showAllContexts") == "true")
}

func (h *summaryHandlers) dayList(w http.ResponseWriter, r *http.Request) {
	resp := map[string]int{}
	h.mgr.WithSession(func(s core.Session) error {
		resp = s.GetContextCountByDateMap()
		return nil
	})
	writeJSON(w, http.StatusOK, resp)
}

func (h *summaryHandlers) respondDay(w http.ResponseWriter, date time.Time, showAll bool) {
	resp := DaySummaryResponse{}
	loc := getLoc()
	durations := map[string]time.Duration{}
	total := time.Duration(0)

	h.mgr.WithSession(func(s core.Session) error {
		for ctxID := range s.State.Contexts {
			d, err := s.GetIntervalDurationsByDate(ctxID, ctxtime.ZonedTime{Time: date, Timezone: loc.String()})
			util.Checkm(err, "Unable to get interval durations for context "+ctxID)
			durations[ctxID] = roundDuration(d, "nanosecond")
		}
		ids := make([]string, 0, len(durations))
		for k := range durations {
			ids = append(ids, k)
		}
		sort.Strings(ids)

		for _, id := range ids {
			d := durations[id]
			if d <= 0 {
				continue
			}
			ctx := s.MustGetCtx(id)
			total += d
			ints := s.GetIntervalsByDate(id, ctxtime.ZonedTime{Time: date, Timezone: loc.String()})
			m := make(map[string]core.Interval)
			for _, it := range ints {
				m[it.Id] = it
			}
			resp.Contexts = append(resp.Contexts, core.Context{
				Id: id, Description: ctx.Description, Intervals: m, Duration: d,
			})
		}
		if showAll {
			resp.OtherContexts = []core.Context{}
			for k, v := range s.State.Contexts {
				if !contextInList(resp.Contexts, k) {
					resp.OtherContexts = append(resp.OtherContexts, v)
				}
			}
		}
		resp.Duration = total
		return nil
	})
	writeJSON(w, http.StatusOK, resp)
}

func contextInList(list []core.Context, id string) bool {
	for _, c := range list {
		if c.Id == id {
			return true
		}
	}
	return false
}

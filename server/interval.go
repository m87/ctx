package server

import (
	"net/http"
	"time"

	"github.com/m87/ctx/core"
)

type IntervalHandler struct {
	manager *core.ContextManager
}

type DayReport struct {
	Contexts  []*core.Context  `json:"contexts"`
	Intervals []*core.Interval `json:"intervals"`
}

func registerIntervalHandler(mux *http.ServeMux, manager *core.ContextManager) {
	handler := &IntervalHandler{manager: manager}
	mux.HandleFunc("GET /", handler.listIntervals)
	mux.HandleFunc("GET /day/{date}", handler.listByDay)
}

func (h *IntervalHandler) listIntervals(w http.ResponseWriter, r *http.Request) {

}

func (h *IntervalHandler) listByDay(w http.ResponseWriter, r *http.Request) {
	dateStr := r.PathValue("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Invalid date format, expected YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	intervals, err := h.manager.IntervalRepository.ListByDay(date)
	if err != nil {
		http.Error(w, "Failed to list intervals", http.StatusInternalServerError)
		return
	}

	seen := make(map[string]struct{})
	for _, interval := range intervals {
		seen[interval.ContextId] = struct{}{}
	}

	contexts := make([]*core.Context, 0, len(seen))
	for contextId := range seen {
		ctx, err := h.manager.ContextRepository.GetById(contextId)
		if err != nil || ctx == nil {
			continue
		}
		contexts = append(contexts, ctx)
	}

	writeJson(w, http.StatusOK, &DayReport{
		Contexts:  contexts,
		Intervals: intervals,
	})
}

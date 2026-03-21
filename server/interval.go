package server

import (
	"net/http"
	"sort"
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

type ContextStats struct {
	ContextId     string  `json:"contextId"`
	Duration      int64   `json:"duration"`
	Percentage    float64 `json:"percentage"`
	IntervalCount int     `json:"intervalCount"`
}

type DayStats struct {
	Date         string                      `json:"date"`
	ContextStats []*ContextStats             `json:"contextStats"`
	Contexts     []*core.Context             `json:"contexts"`
	Intervals    map[string][]*core.Interval `json:"intervals"`
	Distribution map[string]float64          `json:"distribution"`
}

func registerIntervalHandler(mux *http.ServeMux, manager *core.ContextManager) {
	handler := &IntervalHandler{manager: manager}
	mux.HandleFunc("GET /", handler.listIntervals)
	mux.HandleFunc("GET /day/{date}", handler.listByDay)
	mux.HandleFunc("GET /day/{date}/stats", handler.statsByDay)
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

func (h *IntervalHandler) statsByDay(w http.ResponseWriter, r *http.Request) {
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

	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	dayEnd := dayStart.Add(24 * time.Hour)
	now := h.manager.TimeProvider.Now().Time.UTC()

	contextsById := make(map[string]*core.Context)
	intervalsByContext := make(map[string][]*core.Interval)
	durationByContext := make(map[string]time.Duration)

	for _, interval := range intervals {
		if interval == nil || interval.ContextId == "" {
			continue
		}

		intervalsByContext[interval.ContextId] = append(intervalsByContext[interval.ContextId], interval)

		if _, exists := contextsById[interval.ContextId]; !exists {
			ctx, err := h.manager.ContextRepository.GetById(interval.ContextId)
			if err == nil && ctx != nil {
				contextsById[interval.ContextId] = ctx
			}
		}

		start := interval.Start.Time.UTC()
		end := interval.End.Time.UTC()
		if end.IsZero() {
			end = now
		}

		if end.Before(dayStart) || !start.Before(dayEnd) {
			continue
		}

		if start.Before(dayStart) {
			start = dayStart
		}
		if end.After(dayEnd) {
			end = dayEnd
		}
		if !end.After(start) {
			continue
		}

		durationByContext[interval.ContextId] += end.Sub(start)
	}

	contexts := make([]*core.Context, 0, len(contextsById))
	for _, ctx := range contextsById {
		contexts = append(contexts, ctx)
	}
	sort.Slice(contexts, func(i, j int) bool { return contexts[i].Name < contexts[j].Name })

	dayDuration := 24 * time.Hour
	distribution := make(map[string]float64, len(durationByContext))
	contextStats := make([]*ContextStats, 0, len(durationByContext))
	for contextId, duration := range durationByContext {
		percentage := (float64(duration) / float64(dayDuration)) * 100
		distribution[contextId] = percentage
		contextStats = append(contextStats, &ContextStats{
			ContextId:     contextId,
			Duration:      int64(duration),
			Percentage:    percentage,
			IntervalCount: len(intervalsByContext[contextId]),
		})
	}
	sort.Slice(contextStats, func(i, j int) bool { return contextStats[i].Duration > contextStats[j].Duration })

	writeJson(w, http.StatusOK, &DayStats{
		Date:         date.Format("2006-01-02"),
		ContextStats: contextStats,
		Contexts:     contexts,
		Intervals:    intervalsByContext,
		Distribution: distribution,
	})
}

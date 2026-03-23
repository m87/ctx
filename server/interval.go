package server

import (
	"encoding/json"
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
	mux.HandleFunc("GET /day/{date}", handler.listByDay)
	mux.HandleFunc("GET /day/{date}/stats", handler.statsByDay)
	mux.HandleFunc("DELETE /{id}", handler.deleteInterval)
	mux.HandleFunc("PUT /{id}", handler.updateInterval)
	mux.HandleFunc("POST /", handler.createInterval)
	mux.HandleFunc("PATCH /{id}/move/{targetId}", handler.moveInterval)
}

func (h *IntervalHandler) moveInterval(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	targetId := r.PathValue("targetId")
	if id == "" || targetId == "" {
		http.Error(w, "Missing interval ID or target ID", http.StatusBadRequest)
		return
	}

	interval, err := h.manager.IntervalRepository.GetById(id)
	if err != nil || interval == nil {
		http.Error(w, "Interval not found", http.StatusNotFound)
		return
	}

	context, err := h.manager.ContextRepository.GetById(targetId)
	if err != nil || context == nil {
		http.Error(w, "Target context not found", http.StatusNotFound)
		return
	}

	interval.ContextId = targetId
	_, err = h.manager.IntervalRepository.Save(interval)
	if err != nil {
		http.Error(w, "Failed to move interval", http.StatusInternalServerError)
		return
	}

	writeJson(w, http.StatusOK, interval)
}

func (h *IntervalHandler) createInterval(w http.ResponseWriter, r *http.Request) {
	var interval core.Interval
	err := json.NewDecoder(r.Body).Decode(&interval)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.manager.IntervalRepository.Save(&interval)
	if err != nil {
		http.Error(w, "Failed to save interval", http.StatusInternalServerError)
		return
	}
	interval.Id = id

	writeJson(w, http.StatusOK, &interval)
}

func (h *IntervalHandler) updateInterval(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing interval ID", http.StatusBadRequest)
		return
	}

	var interval core.Interval
	err := json.NewDecoder(r.Body).Decode(&interval)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	interval.Id = id
	_, err = h.manager.IntervalRepository.Save(&interval)
	if err != nil {
		http.Error(w, "Failed to save interval", http.StatusInternalServerError)
		return
	}

	writeJson(w, http.StatusOK, &interval)
}

func (h *IntervalHandler) deleteInterval(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing interval ID", http.StatusBadRequest)
		return
	}
	err := h.manager.IntervalRepository.Delete(id)
	if err != nil {
		http.Error(w, "Failed to delete interval", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
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
	now := h.manager.TimeProvider.Now().Time.UTC()
	clippedIntervals := make([]*core.Interval, 0, len(intervals))

	seen := make(map[string]struct{})
	for _, interval := range intervals {
		rng, ok := core.ClipIntervalRangeToDay(interval, date, now)
		if !ok {
			continue
		}

		clipped := *interval
		clipped.Start = core.ZonedTime{Time: rng.Start, Timezone: "UTC"}
		clipped.End = core.ZonedTime{Time: rng.End, Timezone: "UTC"}
		clipped.Duration = rng.End.Sub(rng.Start)

		clippedIntervals = append(clippedIntervals, &clipped)
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
		Intervals: clippedIntervals,
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

	now := h.manager.TimeProvider.Now().Time.UTC()

	contextsById := make(map[string]*core.Context)
	intervalsByContext := make(map[string][]*core.Interval)
	rangesByContext := make(map[string][]core.TimeRange)
	intervalCountByContext := make(map[string]int)

	for _, interval := range intervals {
		if interval == nil || interval.ContextId == "" {
			continue
		}

		rng, ok := core.ClipIntervalRangeToDay(interval, date, now)
		if !ok {
			continue
		}

		clipped := *interval
		clipped.Start = core.ZonedTime{Time: rng.Start, Timezone: "UTC"}
		clipped.End = core.ZonedTime{Time: rng.End, Timezone: "UTC"}
		clipped.Duration = rng.End.Sub(rng.Start)

		intervalsByContext[interval.ContextId] = append(intervalsByContext[interval.ContextId], &clipped)
		rangesByContext[interval.ContextId] = append(rangesByContext[interval.ContextId], rng)
		intervalCountByContext[interval.ContextId]++

		if _, exists := contextsById[interval.ContextId]; !exists {
			ctx, err := h.manager.ContextRepository.GetById(interval.ContextId)
			if err == nil && ctx != nil {
				contextsById[interval.ContextId] = ctx
			}
		}
	}

	contexts := make([]*core.Context, 0, len(contextsById))
	for _, ctx := range contextsById {
		contexts = append(contexts, ctx)
	}
	sort.Slice(contexts, func(i, j int) bool { return contexts[i].Name < contexts[j].Name })

	durationByContext := make(map[string]time.Duration, len(rangesByContext))
	var totalTrackedDuration time.Duration
	for contextId, ranges := range rangesByContext {
		duration := core.SumMergedRangesDuration(ranges)
		durationByContext[contextId] = duration
		totalTrackedDuration += duration
	}

	distribution := make(map[string]float64, len(durationByContext))
	contextStats := make([]*ContextStats, 0, len(durationByContext))
	for contextId, duration := range durationByContext {
		percentage := 0.0
		if totalTrackedDuration > 0 {
			percentage = (float64(duration) / float64(totalTrackedDuration)) * 100
		}
		distribution[contextId] = percentage
		contextStats = append(contextStats, &ContextStats{
			ContextId:     contextId,
			Duration:      int64(duration),
			Percentage:    percentage,
			IntervalCount: intervalCountByContext[contextId],
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

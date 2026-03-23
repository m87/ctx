package core

import (
	"sort"
	"time"
)

type TimeRange struct {
	Start time.Time
	End   time.Time
}

type ContextManager struct {
	TimeProvider       TimeProvider
	ContextRepository  ContextRepository
	IntervalRepository IntervalRepository
}

func NewContextManager(tp TimeProvider, contextRepo ContextRepository, intervalRepo IntervalRepository) *ContextManager {
	return &ContextManager{
		TimeProvider:       tp,
		ContextRepository:  contextRepo,
		IntervalRepository: intervalRepo,
	}
}

func (m *ContextManager) SwitchContext(context *Context) error {
	activeContext, _ := m.ContextRepository.GetActive()
	endTime := m.TimeProvider.Now()
	startTime := m.TimeProvider.Now()

	if activeContext != nil && activeContext.Id == context.Id {
		return nil
	}

	if activeContext != nil {
		activeContext.Status = "inactive"
		m.ContextRepository.Save(activeContext)

		activeInterval, _ := m.IntervalRepository.GetActiveIntervalByContextId(activeContext.Id)

		if activeInterval != nil {
			activeInterval.Duration = endTime.Time.Sub(activeInterval.Start.Time)
			activeInterval.End = endTime
			activeInterval.Status = "completed"
			m.IntervalRepository.Save(activeInterval)
		}
	}

	if context.Id == "" {
		id, _ := m.ContextRepository.Save(context)
		context.Id = id
	}

	context, err := m.ContextRepository.GetById(context.Id)
	if err != nil {
		return err
	}

	context.Status = "active"
	m.ContextRepository.Save(context)

	newInterval := &Interval{
		ContextId: context.Id,
		Start:     startTime,
		Status:    "active",
	}
	m.IntervalRepository.Save(newInterval)

	return nil
}

func (m *ContextManager) FreeActiveContext() error {
	activeContext, err := m.ContextRepository.GetActive()
	if err != nil {
		return err
	}
	if activeContext == nil {
		return nil
	}

	endTime := m.TimeProvider.Now()

	activeContext.Status = "inactive"
	if _, err := m.ContextRepository.Save(activeContext); err != nil {
		return err
	}

	activeInterval, err := m.IntervalRepository.GetActiveIntervalByContextId(activeContext.Id)
	if err != nil {
		return err
	}

	if activeInterval != nil {
		activeInterval.Duration = endTime.Time.Sub(activeInterval.Start.Time)
		activeInterval.End = endTime
		activeInterval.Status = "completed"
		if _, err := m.IntervalRepository.Save(activeInterval); err != nil {
			return err
		}
	}

	return nil
}

func (m *ContextManager) GetStats(contextId string, date time.Time) (*ContextStats, error) {
	allIntervalsByContext, err := m.IntervalRepository.ListByContextId(contextId)
	if err != nil {
		return nil, err
	}

	now := m.TimeProvider.Now().Time.UTC()

	var totalDuration time.Duration
	var totalSessions int
	var sessions int
	dayRanges := make([]TimeRange, 0, len(allIntervalsByContext))

	for _, interval := range allIntervalsByContext {
		intervalDuration := interval.Duration
		if intervalDuration <= 0 {
			start := interval.Start.Time.UTC()
			if start.IsZero() {
				intervalDuration = 0
			} else {
				end := interval.End.Time.UTC()
				if end.IsZero() {
					if interval.Status == "active" {
						end = now
					}
				}
				if end.After(start) {
					intervalDuration = end.Sub(start)
				}
			}
		}

		totalDuration += intervalDuration
		if intervalDuration > 0 {
			totalSessions++
		}

		if dayRange, ok := ClipIntervalRangeToDay(interval, date, now); ok {
			dayRanges = append(dayRanges, dayRange)
			sessions++
		}
	}

	duration := SumMergedRangesDuration(dayRanges)

	return &ContextStats{
		ContextId:     contextId,
		Date:          date,
		Duration:      duration,
		Sessions:      sessions,
		TotalDuration: totalDuration,
		TotalSessions: totalSessions,
	}, nil
}

func ClipIntervalRangeToDay(interval *Interval, date time.Time, now time.Time) (TimeRange, bool) {
	if interval == nil {
		return TimeRange{}, false
	}

	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	dayEnd := dayStart.Add(24 * time.Hour)

	start := interval.Start.Time.UTC()
	if start.IsZero() {
		return TimeRange{}, false
	}

	end := interval.End.Time.UTC()
	if end.IsZero() {
		if interval.Status != "active" {
			return TimeRange{}, false
		}
		end = now
	}

	if end.Before(dayStart) || !start.Before(dayEnd) {
		return TimeRange{}, false
	}

	if start.Before(dayStart) {
		start = dayStart
	}
	if end.After(dayEnd) {
		end = dayEnd
	}

	if !end.After(start) {
		return TimeRange{}, false
	}

	return TimeRange{Start: start, End: end}, true
}

func ClipIntervalDurationToDay(interval *Interval, date time.Time, now time.Time) time.Duration {
	rng, ok := ClipIntervalRangeToDay(interval, date, now)
	if !ok {
		return 0
	}
	return rng.End.Sub(rng.Start)
}

func SumMergedRangesDuration(ranges []TimeRange) time.Duration {
	if len(ranges) == 0 {
		return 0
	}

	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].Start.Before(ranges[j].Start)
	})

	mergedStart := ranges[0].Start
	mergedEnd := ranges[0].End
	var total time.Duration

	for _, rng := range ranges[1:] {
		if !rng.Start.After(mergedEnd) {
			if rng.End.After(mergedEnd) {
				mergedEnd = rng.End
			}
			continue
		}

		total += mergedEnd.Sub(mergedStart)
		mergedStart = rng.Start
		mergedEnd = rng.End
	}

	total += mergedEnd.Sub(mergedStart)
	return total
}

type ContextStats struct {
	ContextId     string        `json:"contextId"`
	Date          time.Time     `json:"date"`
	Duration      time.Duration `json:"duration"`
	Sessions      int           `json:"sessions"`
	TotalDuration time.Duration `json:"totalDuration"`
	TotalSessions int           `json:"totalSessions"`
}

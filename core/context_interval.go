package core

import (
	"errors"
	"time"

	ctxtime "github.com/m87/ctx/time"
)

func (manager *ContextManager) GetIntervalDurationsByDate(s *State, id string, date ctxtime.ZonedTime) (time.Duration, error) {
	var duration time.Duration = 0
	loc, err := time.LoadLocation(ctxtime.DetectTimezoneName())
	if err != nil {
		loc = time.UTC
	}
	startOfDay := time.Date(date.Time.Year(), date.Time.Month(), date.Time.Day(), 0, 0, 0, 0, loc)
	if ctx, ok := s.Contexts[id]; ok {
		for _, interval := range ctx.Intervals {
			if interval.Start.Time.Day() == startOfDay.Day() && interval.Start.Time.Month() == startOfDay.Month() && interval.Start.Time.Year() == startOfDay.Year() && interval.End.Time.Day() == startOfDay.Day() && interval.End.Time.Month() == startOfDay.Month() && interval.End.Time.Year() == startOfDay.Year() {
				duration += interval.Duration
			} else if interval.Start.Time.Before(startOfDay) && interval.End.Time.Day() == startOfDay.Day() && interval.End.Time.Month() == startOfDay.Month() && interval.End.Time.Year() == startOfDay.Year() {
				duration += interval.End.Time.Sub(startOfDay)
			} else if interval.Start.Time.Day() == startOfDay.Day() && interval.Start.Time.Month() == startOfDay.Month() && interval.Start.Time.Year() == startOfDay.Year() && interval.End.Time.After(startOfDay) {
				duration += 24*time.Hour - interval.Start.Time.Sub(startOfDay)
			} else if interval.Start.Time.Before(startOfDay) && interval.End.Time.After(startOfDay) {
				duration += 24 * time.Hour
			}
		}
	} else {
		return 0, errors.New("context does not exist")
	}
	return duration, nil
}

func (manager *ContextManager) GetIntervalsByDate(s *State, id string, date ctxtime.ZonedTime) []Interval {
	intervals := []Interval{}
	loc, err := time.LoadLocation(ctxtime.DetectTimezoneName())
	if err != nil {
		loc = time.UTC
	}
	startOfDay := time.Date(date.Time.Year(), date.Time.Month(), date.Time.Day(), 0, 0, 0, 0, loc)
	if ctx, ok := s.Contexts[id]; ok {
		for _, interval := range ctx.Intervals {
			if interval.Start.Time.Day() == startOfDay.Day() && interval.Start.Time.Month() == startOfDay.Month() && interval.Start.Time.Year() == startOfDay.Year() && interval.End.Time.Day() == startOfDay.Day() && interval.End.Time.Month() == startOfDay.Month() && interval.End.Time.Year() == startOfDay.Year() {
				intervals = append(intervals, Interval(interval))
			} else if interval.Start.Time.Before(startOfDay) && interval.End.Time.Day() == startOfDay.Day() && interval.End.Time.Month() == startOfDay.Month() && interval.End.Time.Year() == startOfDay.Year() {
				intervals = append(intervals, Interval(interval))
			} else if interval.Start.Time.Day() == startOfDay.Day() && interval.Start.Time.Month() == startOfDay.Month() && interval.Start.Time.Year() == startOfDay.Year() && interval.End.Time.After(startOfDay) {
				intervals = append(intervals, Interval(interval))
			} else if interval.Start.Time.Before(startOfDay) && interval.End.Time.After(startOfDay) {
				intervals = append(intervals, Interval(interval))
			}
		}
	}
	return intervals
}

func (state *State) DeleteInterval(ctxId string, id string) error {
	ctx, _ := state.Contexts[ctxId]
	for i, interval := range ctx.Intervals {
		if interval.Id == id {
			return state.DeleteIntervalByIndex(ctxId, i)
		}
	}
	return nil
}

func (state *State) DeleteIntervalByIndex(id string, index int) error {

	if state.CurrentId == id {
		return errors.New("context is active")
	}

	if _, ok := state.Contexts[id]; ok {
		if index < 0 || index >= len(state.Contexts[id].Intervals) {
			return errors.New("index out of range")
		}
		ctx := state.Contexts[id]
		interval := ctx.Intervals[index]
		ctx.Intervals = append(ctx.Intervals[:index], ctx.Intervals[index+1:]...)
		ctx.Duration = ctx.Duration - interval.Duration
		state.Contexts[id] = ctx
	} else {
		return errors.New("context does not exists")
	}
	return nil

}

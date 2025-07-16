package core

import (
	"errors"
	"time"

	ctxtime "github.com/m87/ctx/time"
)

func (session *Session) DeleteInterval(ctxId string, id string) error {
	if err := session.IsValidContext(ctxId); err != nil {
		return err
	}

	err := session.ValidateIntervalExists(ctxId, id)
	if err != nil {
		return err
	}

	ctx := session.State.Contexts[ctxId]
	ctx.Duration -= ctx.Intervals[id].Duration
	delete(ctx.Intervals, id)
	session.State.Contexts[ctxId] = ctx

	return nil
}

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

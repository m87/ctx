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


func (session *Session) DeleteInterval(ctxId string, id string) error {
	ctx, _ := session.State.Contexts[ctxId]
	for i, interval := range ctx.Intervals {
		if interval.Id == id {
			return session.DeleteIntervalByIndex(ctxId, i)
		}
	}
	return nil
}


func (session *Session) DeleteIntervalByIndex(id string, index int) error {
	if session.State.CurrentId == id {
		return errors.New("context is active")
	}

	if _, ok := session.State.Contexts[id]; ok {
		if index < 0 || index >= len(session.State.Contexts[id].Intervals) {
			return errors.New("index out of range")
		}
		ctx := session.State.Contexts[id]
		interval := ctx.Intervals[index]
		ctx.Intervals = append(ctx.Intervals[:index], ctx.Intervals[index+1:]...)
		ctx.Duration = ctx.Duration - interval.Duration
		session.State.Contexts[id] = ctx
	} else {
		return errors.New("context does not exists")
	}
	return nil

}

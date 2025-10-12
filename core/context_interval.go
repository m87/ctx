package core

import (
	"time"

	"github.com/google/uuid"
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

func (session *Session) GetActiveIntervals(ctxId string) ([]string, error) {
	intervals := []string{}
	if err := session.ValidateContextExists(ctxId); err != nil {
		return intervals, err
	}

	for _, interval := range session.State.Contexts[ctxId].Intervals {
		if interval.End.Time.IsZero() {
			intervals = append(intervals, interval.Id)
		}
	}

	return intervals, nil
}

func (session *Session) endInterval(ctxId string, now ctxtime.ZonedTime) error {
	if err := session.ValidateContextExists(ctxId); err != nil {
		return err
	}
	intervals, err := session.GetActiveIntervals(ctxId)
	if err != nil {
		return err
	}

	state := session.State

	for _, intervalId := range intervals {
		ctx := state.Contexts[ctxId]
		interval := ctx.Intervals[intervalId]
		interval.End = now
		interval.Duration = interval.End.Time.Sub(interval.Start.Time)
		ctx.Intervals[interval.Id] = interval
		ctx.Duration += interval.Duration
		session.SetCtx(ctx)
	}
	return nil
}

func (session *Session) SplitContextIntervalById(ctxId string, id string, h int, m int, s int) error {
	state := session.State

	if err := session.ValidateContextExists(ctxId); err != nil {
		return err
	}
	context := session.MustGetCtx(ctxId)

	interval := context.Intervals[id]

	newId := uuid.NewString()
	newStartTime := time.Date(interval.Start.Time.Year(), interval.Start.Time.Month(), interval.Start.Time.Day(), h, m, s, interval.Start.Time.Nanosecond(), interval.Start.Time.Location())

	context.Intervals[newId] = Interval{
		Id: newId,
		Start: ctxtime.ZonedTime{
			Time:     newStartTime,
			Timezone: interval.Start.Timezone,
		},
		End: ctxtime.ZonedTime{
			Time:     interval.End.Time,
			Timezone: interval.End.Timezone,
		},
		Duration: interval.End.Time.Sub(newStartTime),
	}

	newEndTime := time.Date(interval.End.Time.Year(), interval.End.Time.Month(), interval.End.Time.Day(), h, m, s, interval.End.Time.Nanosecond(), interval.End.Time.Location())
	interval.Duration = newEndTime.Sub(interval.Start.Time)
	interval.End.Time = newEndTime
	context.Intervals[id] = interval

	state.Contexts[ctxId] = context

	return nil

}

func (session *Session) EditContextInterval(ctxId string, intervalId string, start ctxtime.ZonedTime, end ctxtime.ZonedTime) error {

	if err := session.IsValidContext(ctxId); err != nil {
		return err
	}
	context := session.MustGetCtx(ctxId)

	for _, interval := range context.Intervals {
		if interval.Id == intervalId {
			session.EditContextIntervalById(ctxId, intervalId, start, end)
			return nil
		}
	}
	return nil
}

func (session *Session) MoveIntervalById(idSrc string, idTarget string, intervalId string) error {
	state := session.State
	if err := session.ValidateContextsExist(idSrc, idTarget); err != nil {
		return err
	}

	if err := session.ValidateActiveInterval(idSrc, intervalId); err != nil {
		return err
	}

	ctxSrc := state.Contexts[idSrc]
	ctxTarget := state.Contexts[idTarget]

	ctxTarget.Intervals[intervalId] = ctxSrc.Intervals[intervalId]
	delete(ctxSrc.Intervals, intervalId)

	ctxTarget.Duration += ctxTarget.Intervals[intervalId].Duration
	ctxSrc.Duration -= ctxTarget.Intervals[intervalId].Duration

	state.Contexts[idSrc] = ctxSrc
	state.Contexts[idTarget] = ctxTarget

	return nil
}

func (session *Session) EditContextIntervalById(ctxId string, intervalId string, start ctxtime.ZonedTime, end ctxtime.ZonedTime) error {
	state := session.State
	if err := session.ValidateActiveContext(ctxId); err != nil {
		return err
	}
	oldDuration := state.Contexts[ctxId].Intervals[intervalId].Duration

	ctx := session.MustGetCtx(ctxId)

	interval := ctx.Intervals[intervalId]

	interval.Start = start
	interval.End = end
	interval.Duration = end.Time.Sub(start.Time)
	ctx.Intervals[intervalId] = interval

	durationDiff := interval.Duration - oldDuration

	ctx.Duration = ctx.Duration + durationDiff

	session.SetCtx(ctx)
	return nil

}

func (session *Session) GetIntervalDurationsByDate(ctxId string, date ctxtime.ZonedTime) (time.Duration, error) {
	if err := session.ValidateContextExists(ctxId); err != nil {
		return 0, err
	}

	var duration time.Duration = 0
	loc, err := time.LoadLocation(ctxtime.DetectTimezoneName())
	if err != nil {
		loc = time.UTC
	}
	startOfDay := time.Date(date.Time.Year(), date.Time.Month(), date.Time.Day(), 0, 0, 0, 0, loc)
	ctx := session.MustGetCtx(ctxId)
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
	return duration, nil
}

func (session *Session) GetIntervalsByDate(ctxId string, date ctxtime.ZonedTime) []Interval {
	if err := session.ValidateContextExists(ctxId); err != nil {
		return []Interval{}
	}

	intervals := []Interval{}
	loc, err := time.LoadLocation(date.Timezone)
	if err != nil {
		loc = time.UTC
	}

	startOfDay := time.Date(date.Time.Year(), date.Time.Month(), date.Time.Day(), 0, 0, 0, 0, loc)
	ctx := session.MustGetCtx(ctxId)
	for _, interval := range ctx.Intervals {
		if interval.Start.Time.Day() == startOfDay.Day() && interval.Start.Time.Month() == startOfDay.Month() && interval.Start.Time.Year() == startOfDay.Year() && interval.End.Time.Day() == startOfDay.Day() && interval.End.Time.Month() == startOfDay.Month() && interval.End.Time.Year() == startOfDay.Year() {
			intervals = append(intervals, Interval(interval))
		} else if interval.Start.Time.Before(startOfDay) && interval.End.Time.Day() == startOfDay.Day() && interval.End.Time.Month() == startOfDay.Month() && interval.End.Time.Year() == startOfDay.Year() {
			interval := Interval(interval)
			interval.Start.Time = startOfDay
			interval.Duration = interval.End.Time.Sub(interval.Start.Time)
			intervals = append(intervals, interval)
		} else if interval.Start.Time.Day() == startOfDay.Day() && interval.Start.Time.Month() == startOfDay.Month() && interval.Start.Time.Year() == startOfDay.Year() && interval.End.Time.After(startOfDay) {
			interval := Interval(interval)
			interval.End.Time = startOfDay.Add(24*time.Hour - time.Second)
			interval.Duration = interval.End.Time.Sub(interval.Start.Time)
			intervals = append(intervals, interval)
		} else if interval.Start.Time.Before(startOfDay) && interval.End.Time.After(startOfDay) {
			interval := Interval(interval)
			interval.Start.Time = startOfDay
			interval.End.Time = startOfDay.Add(24*time.Hour - time.Second)
			interval.Duration = interval.End.Time.Sub(interval.Start.Time)
			intervals = append(intervals, interval)
		}
	}
	return intervals
}

func GetDateCountsFromInterval(interval Interval) []string {
	output := []string{}
	for d := interval.Start.Time; !d.After(interval.End.Time); d = d.AddDate(0, 0, 1) {
		output = append(output, d.Format(time.DateOnly))
	}

	return output
}

func (session *Session) GetDateCounts(ctxId string) map[string]int {
	ctx := session.MustGetCtx(ctxId)
	output := make(map[string]int)

	for _, interval := range ctx.Intervals {
		part := GetDateCountsFromInterval(interval)

		for _, d := range part {
			output[d]++
		}
	}

	return output
}

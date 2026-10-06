package core

import (
	"fmt"
	"time"
)

type radarBucketAggregate struct {
	duration float64
	sessions map[string]float64
	contexts map[string]bool
	days     map[string]bool
}

func evaluateDashboardRadar(intervals []*Interval, contexts []*Context, projects []*Project, selected, base map[string]bool, period dashboardPeriodRange, asOf time.Time, location *time.Location, firstDay int, p dashboardMetricProperties) (*DashboardRadarPayload, int, error) {
	unit := dashboardMetricUnit(p.Metric)
	contextByID := make(map[string]*Context, len(contexts))
	for _, context := range contexts {
		if context != nil {
			contextByID[context.Id] = context
		}
	}
	projectByID := make(map[string]*Project, len(projects))
	for _, project := range projects {
		if project != nil {
			projectByID[project.Id] = project
		}
	}
	keys := append([]string(nil), p.AxisKeys...)
	if p.AxesBy == "weekday" {
		keys = make([]string, 7)
		for i := range keys {
			keys[i] = fmt.Sprintf("weekday:%d", (firstDay+i)%7)
		}
	}
	observed := make(map[string]*int, len(keys))
	labels := make(map[string]string, len(keys))
	entities := make(map[string]*DashboardRadarEntity, len(keys))
	for _, key := range keys {
		if p.AxesBy == "weekday" {
			index := 0
			fmt.Sscanf(key, "weekday:%d", &index)
			labels[key] = time.Weekday(index).String()
		} else if p.AxesBy == "project" {
			if key == "project:none" {
				labels[key] = "No project"
			} else {
				id := key[len("project:"):]
				project := projectByID[id]
				if project == nil {
					return nil, 0, fmt.Errorf("A selected project is no longer available.")
				}
				labels[key], entities[key] = project.Name, &DashboardRadarEntity{Type: "project", ID: id}
			}
		} else {
			id := key[len("context:"):]
			context := contextByID[id]
			if context == nil {
				return nil, 0, fmt.Errorf("A selected context is no longer available.")
			}
			labels[key], entities[key] = context.Name, &DashboardRadarEntity{Type: "context", ID: id}
		}
	}
	partialDay := make(map[string]bool)
	if p.AxesBy == "weekday" {
		observedEnd := period.end
		if asOf.Before(observedEnd) {
			observedEnd = asOf.In(location)
		}
		for day := period.start.In(location); day.Before(period.end); day = day.AddDate(0, 0, 1) {
			key := fmt.Sprintf("weekday:%d", day.Weekday())
			if !day.Before(observedEnd) {
				continue
			}
			count := 0
			if observed[key] != nil {
				count = *observed[key]
			}
			count++
			observed[key] = &count
			if day.Before(period.start) || day.AddDate(0, 0, 1).After(period.end) || day.AddDate(0, 0, 1).After(observedEnd) {
				partialDay[key] = true
			}
		}
	}
	aggs := make(map[string]*radarBucketAggregate, len(keys))
	for _, key := range keys {
		aggs[key] = &radarBucketAggregate{sessions: map[string]float64{}, contexts: map[string]bool{}, days: map[string]bool{}}
	}
	axisContextIDs := make(map[string]bool, len(selected))
	if p.AxesBy != "weekday" {
		for _, context := range contexts {
			if context == nil || !selected[context.Id] {
				continue
			}
			key := "context:" + context.Id
			if p.AxesBy == "project" {
				key = "project:none"
				if context.ProjectId != nil {
					key = "project:" + *context.ProjectId
				}
			}
			for _, axisKey := range keys {
				if axisKey == key {
					axisContextIDs[context.Id] = true
					break
				}
			}
		}
	}
	denominators := make(map[string]float64, len(keys))
	for _, key := range keys {
		denominators[key] = 0
	}
	ignoredIDs := map[string]bool{}
	for _, interval := range intervals {
		if interval == nil || !base[interval.ContextId] {
			continue
		}
		if !timeIsSet(interval.Start) {
			if selected[interval.ContextId] {
				ignoredIDs[interval.Id] = true
			}
			continue
		}
		end := asOf
		if timeIsSet(interval.End) {
			end = *interval.End
		} else if interval.Status != IntervalStatusActive {
			if selected[interval.ContextId] {
				ignoredIDs[interval.Id] = true
			}
			continue
		} else if !p.IncludeRunning {
			continue
		}
		start := *interval.Start
		if !start.Before(asOf) {
			continue
		}
		if end.After(asOf) {
			end = asOf
		}
		if !end.After(start) {
			if selected[interval.ContextId] {
				ignoredIDs[interval.Id] = true
			}
			continue
		}
		if start.Before(period.start) {
			start = period.start
		}
		if end.After(period.end) {
			end = period.end
		}
		if !end.After(start) {
			continue
		}
		context := contextByID[interval.ContextId]
		if context == nil {
			continue
		}
		if p.AxesBy != "weekday" {
			key := ""
			switch p.AxesBy {
			case "context":
				key = "context:" + context.Id
			case "project":
				if context.ProjectId == nil {
					key = "project:none"
				} else {
					key = "project:" + *context.ProjectId
				}
			}
			if selected[interval.ContextId] {
				aggregateIntervalToBucket(aggs[key], interval.Id, interval.ContextId, start, end, location)
			}
			denominators[key] += end.Sub(start).Seconds() * 1000
			continue
		}
		for cursor := start.In(location); cursor.Before(end); {
			dayEnd := time.Date(cursor.Year(), cursor.Month(), cursor.Day()+1, 0, 0, 0, 0, location)
			segmentEnd := dayEnd
			if end.Before(segmentEnd) {
				segmentEnd = end
			}
			key := fmt.Sprintf("weekday:%d", cursor.Weekday())
			if selected[interval.ContextId] {
				aggregateIntervalToBucket(aggs[key], interval.Id, interval.ContextId, cursor, segmentEnd, location)
			}
			denominators[key] += segmentEnd.Sub(cursor).Seconds() * 1000
			cursor = segmentEnd.In(location)
		}
	}
	axes := make([]DashboardRadarAxis, 0, len(keys))
	for _, key := range keys {
		a := aggs[key]
		value := radarAggregateValue(a.duration, a.sessions, a.contexts, a.days, p.Metric)
		if p.Metric == "time_share.ratio" {
			if denominators[key] > 0 {
				share := a.duration / denominators[key] * 100
				value = &share
			} else {
				value = nil
			}
		}
		status := "complete"
		if p.AxesBy == "weekday" {
			if observed[key] == nil {
				status, value = "unobserved", nil
			} else if partialDay[key] {
				status = "partial"
			}
		} else if value == nil && p.Metric != "duration.mean" && p.Metric != "duration.max" {
			zero := 0.0
			value = &zero
		}
		axes = append(axes, DashboardRadarAxis{Key: key, Label: labels[key], Value: value, Status: status, ObservedDays: observed[key], Entity: entities[key]})
	}
	totalIDs := selected
	if p.AxesBy != "weekday" {
		totalIDs = axisContextIDs
	}
	total := aggregateDashboardIntervals(intervals, totalIDs, period, asOf, p.IncludeRunning).value(p.Metric)
	if p.Metric == "time_share.ratio" {
		denominator := aggregateDashboardIntervals(intervals, base, period, asOf, p.IncludeRunning)
		if denominator.duration > 0 {
			value := aggregateDashboardIntervals(intervals, totalIDs, period, asOf, p.IncludeRunning).duration / denominator.duration * 100
			total = &value
		}
	}
	return &DashboardRadarPayload{Unit: unit, Axes: axes, TotalValue: total}, len(ignoredIDs), nil
}

func aggregateIntervalToBucket(a *radarBucketAggregate, intervalID, contextID string, start, end time.Time, location *time.Location) {
	if a == nil || !end.After(start) {
		return
	}
	duration := end.Sub(start).Seconds() * 1000
	a.duration += duration
	a.sessions[intervalID] += duration
	a.contexts[contextID] = true
	local := start.In(location)
	day := time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, location)
	for day.Before(end) {
		a.days[day.Format(time.DateOnly)] = true
		day = day.AddDate(0, 0, 1)
	}
}

func radarAggregateValue(duration float64, sessions map[string]float64, contexts, days map[string]bool, metric string) *float64 {
	var value float64
	switch metric {
	case "duration.sum":
		value = duration
	case "duration.mean":
		if len(sessions) == 0 {
			return nil
		}
		for _, v := range sessions {
			value += v
		}
		value /= float64(len(sessions))
	case "duration.max":
		if len(sessions) == 0 {
			return nil
		}
		for _, v := range sessions {
			value = max(value, v)
		}
	case "sessions.count":
		value = float64(len(sessions))
	case "contexts.distinct_count":
		value = float64(len(contexts))
	case "active_days.distinct_count":
		value = float64(len(days))
	default:
		return nil
	}
	return &value
}

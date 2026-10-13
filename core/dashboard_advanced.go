package core

import (
	"sort"
	"strings"
	"time"
)

func initializeDashboardAdvancedResult(result *DashboardMetricResult, widgetType string) {
	active := false
	result.HasActivity = &active
	switch widgetType {
	case "bar":
		result.Kind = "categories"
		result.Categories = &DashboardCategoriesPayload{Unit: result.Unit, Categories: []DashboardCategory{}}
	case "heatmap":
		result.Kind = "heatmap"
		result.Heatmap = &DashboardHeatmapPayload{Unit: result.Unit, Days: []DashboardHeatmapDay{}}
	case "line":
		result.Kind = "series"
		result.Series = &DashboardLinePayload{Unit: result.Unit, Buckets: []DashboardLineBucket{}, Series: []DashboardLineSeries{}}
	case "table":
		result.Kind = "table"
		result.Table = &DashboardTablePayload{Columns: []DashboardTablePayloadColumn{}, Rows: []DashboardTableRow{}, Summary: map[string]*float64{}, ColumnMaxima: map[string]*float64{}}
	}
}

type dashboardBarGroup struct {
	key, label string
	ids        map[string]bool
	entity     *DashboardRadarEntity
	value      *float64
}

func evaluateDashboardBar(intervals []*Interval, contexts []*Context, selected, base map[string]bool, period dashboardPeriodRange, asOf time.Time, firstDayOfWeek int, p dashboardMetricProperties) (*DashboardCategoriesPayload, int, error) {
	payload := &DashboardCategoriesPayload{Unit: dashboardMetricUnit(p.Metric), Categories: []DashboardCategory{}}
	whole := aggregateDashboardIntervals(intervals, selected, period, asOf, p.IncludeRunning)
	payload.TotalValue = dashboardMetricValue(intervals, selected, base, period, asOf, p, whole)
	groups := map[string]*dashboardBarGroup{}
	if p.GroupBy == "time" {
		buckets := dashboardBuckets(period, p.Bucket, firstDayOfWeek)
		for _, bucketRange := range buckets {
			cursor, end := bucketRange.start, bucketRange.end
			key := "bucket:" + cursor.UTC().Format(time.RFC3339Nano)
			value := dashboardMetricValue(intervals, selected, base, bucketRange, asOf, p, aggregateDashboardIntervals(intervals, selected, bucketRange, asOf, p.IncludeRunning))
			status := "complete"
			if p.Bucket == "week" && (int(cursor.Weekday()) != firstDayOfWeek || int(end.Weekday()) != firstDayOfWeek) || p.Bucket == "month" && (cursor.Day() != 1 || end.Day() != 1) {
				status = "partial"
			}
			if cursor.Before(asOf) && end.After(asOf) {
				status = "partial"
			}
			if !cursor.Before(asOf) {
				status = "future"
				value = nil
			}
			payload.Categories = append(payload.Categories, DashboardCategory{Key: key, Label: cursor.Format("2006-01-02"), Value: value, Synthetic: false, MemberCount: 1, Start: timePtr(cursor.UTC()), End: timePtr(end.UTC()), Status: status})
			if len(payload.Categories) > 50 {
				return nil, whole.ignored, invalidDefinition("bar charts support at most 50 time buckets")
			}
		}
	} else {
		for _, context := range contexts {
			if context == nil || !selected[context.Id] {
				continue
			}
			key, label := "", ""
			var entity *DashboardRadarEntity
			if p.GroupBy == "context" {
				key, label = "context:"+context.Id, context.Name
				entity = &DashboardRadarEntity{Type: "context", ID: context.Id}
			} else if context.ProjectId == nil {
				key, label = "project:none", "No project"
			} else {
				key, label = "project:"+*context.ProjectId, *context.ProjectId
				if context.Project != nil {
					label = context.Project.Name
				}
				entity = &DashboardRadarEntity{Type: "project", ID: *context.ProjectId}
			}
			group := groups[key]
			if group == nil {
				group = &dashboardBarGroup{key: key, label: label, ids: map[string]bool{}, entity: entity}
				groups[key] = group
			}
			group.ids[context.Id] = true
		}
		for _, group := range groups {
			aggregate := aggregateDashboardIntervals(intervals, group.ids, period, asOf, p.IncludeRunning)
			if aggregate.duration <= 0 {
				continue
			}
			group.value = dashboardMetricValue(intervals, group.ids, base, period, asOf, p, aggregate)
			payload.Categories = append(payload.Categories, DashboardCategory{Key: group.key, Label: group.label, Value: group.value, Entity: group.entity, MemberCount: 1, Status: "complete"})
		}
		sortKey, direction := "value", "desc"
		if p.Sort != nil {
			if value, ok := p.Sort["by"].(string); ok {
				sortKey = value
			}
			if value, ok := p.Sort["direction"].(string); ok {
				direction = value
			}
		}
		sort.SliceStable(payload.Categories, func(i, j int) bool {
			a, b := payload.Categories[i], payload.Categories[j]
			if a.Value == nil || b.Value == nil {
				if a.Value == nil && b.Value != nil {
					return false
				}
				if b.Value == nil && a.Value != nil {
					return true
				}
			}
			less := false
			if sortKey == "label" {
				less = a.Label < b.Label
			} else if a.Value != nil && b.Value != nil {
				less = *a.Value < *b.Value
			}
			if direction == "desc" {
				less = (a.Value != nil && b.Value != nil && *a.Value > *b.Value)
				if sortKey == "label" {
					less = a.Label > b.Label
				}
			}
			if a.Value != nil && b.Value != nil && *a.Value == *b.Value || sortKey == "label" && a.Label == b.Label {
				return a.Key < b.Key
			}
			return less
		})
		limit := len(payload.Categories)
		if p.Limit != nil && *p.Limit < limit {
			limit = *p.Limit
		}
		payload.TotalCategories = len(payload.Categories)
		if limit < len(payload.Categories) {
			if p.ShowOther {
				visible := max(0, limit-1)
				otherValue := 0.0
				members := len(payload.Categories) - visible
				for _, category := range payload.Categories[visible:] {
					if category.Value != nil {
						otherValue += *category.Value
					}
				}
				other := DashboardCategory{Key: "aggregate:other", Label: "Other", Value: &otherValue, Synthetic: true, MemberCount: members, Status: "complete"}
				payload.Categories = append(payload.Categories[:visible], other)
			} else {
				payload.OmittedCategories = len(payload.Categories) - limit
				payload.Categories = payload.Categories[:limit]
			}
		}
	}
	if p.GroupBy == "time" {
		payload.TotalCategories = len(payload.Categories)
	} else if payload.TotalCategories == 0 {
		payload.TotalCategories = len(payload.Categories)
	}
	return payload, whole.ignored, nil
}

func evaluateDashboardHeatmap(intervals []*Interval, selected, base map[string]bool, period dashboardPeriodRange, asOf time.Time, p dashboardMetricProperties) (*DashboardHeatmapPayload, int, error) {
	payload := &DashboardHeatmapPayload{Unit: dashboardMetricUnit(p.Metric), Days: []DashboardHeatmapDay{}}
	whole := aggregateDashboardIntervals(intervals, selected, period, asOf, p.IncludeRunning)
	payload.TotalValue = dashboardMetricValue(intervals, selected, base, period, asOf, p, whole)
	location := period.start.Location()
	for day := period.start; day.Before(period.end); day = day.AddDate(0, 0, 1) {
		end := day.AddDate(0, 0, 1)
		if end.After(period.end) {
			end = period.end
		}
		rangeForDay := dashboardPeriodRange{day, end}
		aggregate := aggregateDashboardIntervals(intervals, selected, rangeForDay, asOf, p.IncludeRunning)
		value := dashboardMetricValue(intervals, selected, base, rangeForDay, asOf, p, aggregate)
		localAsOf := asOf.In(location)
		status := "complete"
		if !day.Before(time.Date(localAsOf.Year(), localAsOf.Month(), localAsOf.Day(), 0, 0, 0, 0, location).AddDate(0, 0, 1)) {
			status, value = "future", nil
		} else if end.After(localAsOf) {
			status = "partial"
		}
		payload.Days = append(payload.Days, DashboardHeatmapDay{Date: day.Format(time.DateOnly), Start: day.UTC(), End: end.UTC(), Value: value, Status: status})
		if aggregate.duration > 0 {
			payload.ActiveDays++
		}
		if len(payload.Days) > 400 {
			return nil, whole.ignored, invalidDefinition("heatmaps support at most 400 days")
		}
	}
	return payload, whole.ignored, nil
}

func dashboardBuckets(period dashboardPeriodRange, bucket string, firstDayOfWeek int) []dashboardPeriodRange {
	buckets := []dashboardPeriodRange{}
	for cursor := period.start; cursor.Before(period.end); {
		end := cursor.AddDate(0, 0, 1)
		if bucket == "week" {
			offset := (int(cursor.Weekday()) - firstDayOfWeek + 7) % 7
			end = cursor.AddDate(0, 0, 7-offset)
		}
		if bucket == "month" {
			end = time.Date(cursor.Year(), cursor.Month()+1, 1, 0, 0, 0, 0, cursor.Location())
		}
		if end.After(period.end) {
			end = period.end
		}
		buckets = append(buckets, dashboardPeriodRange{cursor, end})
		cursor = end
	}
	return buckets
}

func evaluateDashboardLine(intervals []*Interval, contexts []*Context, selected, base map[string]bool, period dashboardPeriodRange, asOf time.Time, firstDayOfWeek int, p dashboardMetricProperties) (*DashboardLinePayload, int, error) {
	buckets := dashboardBuckets(period, p.Bucket, firstDayOfWeek)
	if len(buckets) > 400 {
		return nil, 0, invalidDefinition("line charts support at most 400 buckets")
	}
	total := aggregateDashboardIntervals(intervals, selected, period, asOf, p.IncludeRunning)
	payload := &DashboardLinePayload{Unit: dashboardMetricUnit(p.Metric), TotalValue: dashboardMetricValue(intervals, selected, base, period, asOf, p, total), Buckets: []DashboardLineBucket{}, Series: []DashboardLineSeries{}}
	for _, bucket := range buckets {
		status := "complete"
		if p.Bucket == "week" && (int(bucket.start.Weekday()) != firstDayOfWeek || int(bucket.end.Weekday()) != firstDayOfWeek) || p.Bucket == "month" && (bucket.start.Day() != 1 || bucket.end.Day() != 1) {
			status = "partial"
		}
		if !bucket.end.Before(asOf) && bucket.start.Before(asOf) {
			status = "partial"
		}
		if !bucket.start.Before(asOf) {
			status = "future"
		}
		payload.Buckets = append(payload.Buckets, DashboardLineBucket{Key: bucket.start.UTC().Format(time.RFC3339Nano), Start: bucket.start.UTC(), End: bucket.end.UTC(), Status: status})
	}
	type group struct {
		key, label string
		ids        map[string]bool
		total      float64
	}
	groups := map[string]*group{}
	if p.SeriesBy == "none" {
		groups["all"] = &group{key: "all", label: "Total", ids: selected}
	} else {
		for _, context := range contexts {
			if context == nil || !selected[context.Id] {
				continue
			}
			key, label := "", ""
			if p.SeriesBy == "context" {
				key = "context:" + context.Id
				label = context.Name
			} else if context.ProjectId == nil {
				key = "project:none"
				label = "No project"
			} else {
				key = "project:" + *context.ProjectId
				label = *context.ProjectId
				if context.Project != nil {
					label = context.Project.Name
				}
			}
			if groups[key] == nil {
				groups[key] = &group{key: key, label: label, ids: map[string]bool{}}
			}
			groups[key].ids[context.Id] = true
		}
	}
	ordered := make([]*group, 0, len(groups))
	for _, g := range groups {
		aggregate := aggregateDashboardIntervals(intervals, g.ids, period, asOf, p.IncludeRunning)
		if aggregate.duration == 0 {
			continue
		}
		metricValue := dashboardMetricValue(intervals, g.ids, base, period, asOf, p, aggregate)
		if metricValue != nil {
			g.total = *metricValue
		}
		ordered = append(ordered, g)
	}
	sort.Slice(ordered, func(i, j int) bool {
		if ordered[i].total == ordered[j].total {
			return ordered[i].key < ordered[j].key
		}
		return ordered[i].total > ordered[j].total
	})
	payload.TotalSeries = len(ordered)
	if len(ordered) > p.SeriesLimit {
		payload.OmittedSeries = len(ordered) - p.SeriesLimit
		ordered = ordered[:p.SeriesLimit]
	}
	for _, g := range ordered {
		s := DashboardLineSeries{Key: g.key, Label: g.label, Values: make([]*float64, len(buckets))}
		running := 0.0
		for i, b := range buckets {
			agg := aggregateDashboardIntervals(intervals, g.ids, b, asOf, p.IncludeRunning)
			value := dashboardMetricValue(intervals, g.ids, base, b, asOf, p, agg)
			if p.ValueMode == "cumulative" && value != nil {
				running += *value
				v := running
				value = &v
			}
			if payload.Buckets[i].Status == "future" {
				value = nil
			}
			s.Values[i] = value
		}
		payload.Series = append(payload.Series, s)
	}
	return payload, total.ignored, nil
}

func evaluateDashboardTable(intervals []*Interval, contexts []*Context, selected, base map[string]bool, period dashboardPeriodRange, asOf time.Time, p dashboardMetricProperties) (*DashboardTablePayload, int, error) {
	type row struct {
		value    DashboardTableRow
		activity float64
	}
	rowsByKey := map[string]*row{}
	payload := &DashboardTablePayload{Columns: []DashboardTablePayloadColumn{}, Rows: []DashboardTableRow{}, Summary: map[string]*float64{}, ColumnMaxima: map[string]*float64{}}
	for _, column := range p.Columns {
		payload.Columns = append(payload.Columns, DashboardTablePayloadColumn{ID: column.ID, Metric: column.Metric, Unit: dashboardMetricUnit(column.Metric)})
		payload.ColumnMaxima[column.ID] = nil
	}
	for _, context := range contexts {
		if context == nil || !selected[context.Id] {
			continue
		}
		key, label := "", ""
		var entity *DashboardRadarEntity
		if p.GroupBy == "context" {
			key = "context:" + context.Id
			label = context.Name
			entity = &DashboardRadarEntity{Type: "context", ID: context.Id}
		} else if context.ProjectId == nil {
			key = "project:none"
			label = "No project"
		} else {
			key = "project:" + *context.ProjectId
			label = *context.ProjectId
			if context.Project != nil {
				label = context.Project.Name
			}
			entity = &DashboardRadarEntity{Type: "project", ID: *context.ProjectId}
		}
		entry := rowsByKey[key]
		if entry == nil {
			entry = &row{value: DashboardTableRow{Key: key, Label: label, Entity: entity, Values: map[string]*float64{}}, activity: 0}
			rowsByKey[key] = entry
		}
		entry.activity += aggregateDashboardIntervals(intervals, map[string]bool{context.Id: true}, period, asOf, p.IncludeRunning).duration
	}
	for _, entry := range rowsByKey {
		if entry.activity == 0 {
			continue
		}
		for _, column := range p.Columns {
			metric := p
			metric.Metric = column.Metric
			agg := aggregateDashboardIntervals(intervals, selectedIDsForTable(contexts, selected, entry.value.Key, p.GroupBy), period, asOf, p.IncludeRunning)
			value := dashboardMetricValue(intervals, selectedIDsForTable(contexts, selected, entry.value.Key, p.GroupBy), base, period, asOf, metric, agg)
			entry.value.Values[column.ID] = value
			if value != nil && (payload.ColumnMaxima[column.ID] == nil || *value > *payload.ColumnMaxima[column.ID]) {
				v := *value
				payload.ColumnMaxima[column.ID] = &v
			}
		}
		payload.Rows = append(payload.Rows, entry.value)
	}
	payload.TotalRows = len(payload.Rows)
	for _, column := range p.Columns {
		metric := p
		metric.Metric = column.Metric
		value := dashboardMetricValue(intervals, selected, base, period, asOf, metric, aggregateDashboardIntervals(intervals, selected, period, asOf, p.IncludeRunning))
		payload.Summary[column.ID] = value
	}
	sortColumn, direction := "time", "desc"
	if p.Sort != nil {
		if columnID, ok := p.Sort["columnId"].(string); ok && columnID != "" {
			sortColumn = columnID
		}
		if value, ok := p.Sort["direction"].(string); ok {
			direction = value
		}
	}
	sort.SliceStable(payload.Rows, func(i, j int) bool {
		a, b := payload.Rows[i], payload.Rows[j]
		less := false
		if sortColumn == "group" {
			less = strings.ToLower(a.Label) < strings.ToLower(b.Label)
		} else {
			av, bv := a.Values[sortColumn], b.Values[sortColumn]
			if av == nil || bv == nil {
				if av == nil && bv != nil {
					return false
				}
				if bv == nil && av != nil {
					return true
				}
			} else if *av != *bv {
				less = *av < *bv
			}
		}
		if direction == "desc" {
			if sortColumn == "group" {
				less = strings.ToLower(a.Label) > strings.ToLower(b.Label)
			} else if a.Values[sortColumn] != nil && b.Values[sortColumn] != nil {
				less = *a.Values[sortColumn] > *b.Values[sortColumn]
			}
		}
		if a.Label == b.Label || sortColumn != "group" && a.Values[sortColumn] != nil && b.Values[sortColumn] != nil && *a.Values[sortColumn] == *b.Values[sortColumn] {
			return a.Key < b.Key
		}
		return less
	})
	return payload, aggregateDashboardIntervals(intervals, selected, period, asOf, p.IncludeRunning).ignored, nil
}

func selectedIDsForTable(contexts []*Context, selected map[string]bool, key, groupBy string) map[string]bool {
	result := map[string]bool{}
	for _, c := range contexts {
		if c == nil || !selected[c.Id] {
			continue
		}
		matches := groupBy == "context" && key == "context:"+c.Id || groupBy == "project" && (key == "project:none" && c.ProjectId == nil || c.ProjectId != nil && key == "project:"+*c.ProjectId)
		if matches {
			result[c.Id] = true
		}
	}
	return result
}

func dashboardMetricValue(intervals []*Interval, selected, base map[string]bool, period dashboardPeriodRange, asOf time.Time, p dashboardMetricProperties, aggregate dashboardAggregate) *float64 {
	value := aggregate.value(p.Metric)
	if p.Metric == "time_share.ratio" {
		denominator := aggregateDashboardIntervals(intervals, base, period, asOf, p.IncludeRunning)
		if denominator.duration > 0 {
			share := aggregate.duration / denominator.duration * 100
			value = &share
		}
	}
	return value
}

func timePtr(value time.Time) *time.Time { return &value }

package core

import (
	"encoding/json"
	"math"
	"slices"
	"strings"
	"time"
)

type dashboardMetricProperties struct {
	Title          string                    `json:"title"`
	Description    string                    `json:"description"`
	Metric         string                    `json:"metric"`
	Period         string                    `json:"period"`
	StartDate      string                    `json:"startDate"`
	EndDate        string                    `json:"endDate"`
	ProjectID      string                    `json:"projectId"`
	ContextID      string                    `json:"contextId"`
	Tag            string                    `json:"tag"`
	Archived       string                    `json:"archived"`
	IncludeRunning bool                      `json:"includeRunning"`
	Format         string                    `json:"format"`
	Precision      int                       `json:"precision"`
	Minimum        float64                   `json:"minimum"`
	Maximum        float64                   `json:"maximum"`
	Target         *float64                  `json:"target"`
	GoalMode       string                    `json:"goalMode"`
	Thresholds     []dashboardGaugeThreshold `json:"thresholds"`
	ShowValue      bool                      `json:"showValue"`
	ShowTarget     bool                      `json:"showTarget"`
	ShowPercentage bool                      `json:"showPercentage"`
	AxesBy         string                    `json:"axesBy"`
	AxisKeys       []string                  `json:"axisKeys"`
	ScaleMaximum   *float64                  `json:"scaleMaximum"`
	TargetValue    *float64                  `json:"targetValue"`
	ShowValues     bool                      `json:"showValues"`
	ShowGrid       bool                      `json:"showGrid"`
	ShowFill       bool                      `json:"showFill"`
	ShowSummary    bool                      `json:"showSummary"`
	GroupBy        string                    `json:"groupBy"`
	Bucket         string                    `json:"bucket"`
	Orientation    string                    `json:"orientation"`
	Sort           map[string]interface{}    `json:"sort"`
	Limit          *int                      `json:"limit"`
	ShowOther      bool                      `json:"showOther"`
	ScaleMode      string                    `json:"scaleMode"`
	MonthLabels    bool                      `json:"showMonthLabels"`
	DayLabels      bool                      `json:"showDayLabels"`
	Legend         bool                      `json:"showLegend"`
	ThresholdsHeat []float64                 `json:"-"`
	SeriesBy       string                    `json:"seriesBy"`
	SeriesLimit    int                       `json:"seriesLimit"`
	ValueMode      string                    `json:"valueMode"`
	RenderMode     string                    `json:"renderMode"`
	ShowPoints     bool                      `json:"showPoints"`
	TargetLine     *float64                  `json:"targetValue"`
	Columns        []dashboardTableColumn    `json:"columns"`
	PageSize       int                       `json:"pageSize"`
}

type dashboardTableColumn struct {
	ID        string `json:"id"`
	Metric    string `json:"metric"`
	Label     string `json:"label"`
	Visible   bool   `json:"visible"`
	Format    string `json:"format"`
	Precision int    `json:"precision"`
	ShowBar   bool   `json:"showBar"`
}

type dashboardGaugeThreshold struct {
	Value float64 `json:"value"`
	Tone  string  `json:"tone"`
}

var dashboardPeriods = []string{"today", "yesterday", "this_week", "last_week", "this_month", "last_month", "last_7_days", "last_30_days"}

func dashboardMetricUnit(metric string) string {
	switch metric {
	case "duration.sum", "duration.mean", "duration.max":
		return "ms"
	case "sessions.count", "contexts.distinct_count", "active_days.distinct_count":
		return "count"
	case "time_share.ratio":
		return "percent"
	default:
		return ""
	}
}

func readDashboardMetric(widget dashboardWidget) (dashboardMetricProperties, error) {
	var p dashboardMetricProperties
	raw, err := json.Marshal(widget.Properties)
	if err != nil {
		return p, invalidDefinition("invalid metric settings")
	}
	if err := json.Unmarshal(raw, &p); err != nil {
		return p, invalidDefinition("invalid metric settings")
	}
	return p, nil
}

func validateDashboardMetric(widget dashboardWidget) error {
	p, err := readDashboardMetric(widget)
	if err != nil {
		return err
	}
	if widget.Type == "table" {
		return validateDashboardTableMetric(widget, p)
	}
	for _, key := range []string{"title", "description", "startDate", "endDate", "projectId", "contextId", "tag"} {
		text, ok := widget.Properties[key].(string)
		if !ok || len(text) > 500 {
			return invalidDefinition("invalid metric " + key)
		}
	}
	unit := dashboardMetricUnit(p.Metric)
	if unit == "" {
		return invalidDefinition("unknown metric")
	}
	localHeatmapPeriod := widget.Type == "heatmap" && slices.Contains([]string{"last_3_months", "last_6_months", "last_12_months"}, p.Period)
	if p.Period != "inherit" && p.Period != "custom" && !slices.Contains(dashboardPeriods, p.Period) && !localHeatmapPeriod {
		return invalidDefinition("invalid metric period")
	}
	if p.Period == "custom" {
		start, e1 := time.Parse(time.DateOnly, p.StartDate)
		end, e2 := time.Parse(time.DateOnly, p.EndDate)
		if e1 != nil || e2 != nil || end.Before(start) {
			return invalidDefinition("invalid metric date range")
		}
	}
	if !slices.Contains([]string{"include", "exclude", "only"}, p.Archived) {
		return invalidDefinition("invalid archived filter")
	}
	if _, ok := widget.Properties["includeRunning"].(bool); !ok {
		return invalidDefinition("includeRunning must be a boolean")
	}
	formats := map[string][]string{"ms": {"duration", "hours"}, "count": {"number"}, "percent": {"percent"}}
	if !slices.Contains(formats[unit], p.Format) {
		return invalidDefinition("metric format does not match its unit")
	}
	precision, ok := widget.Properties["precision"].(float64)
	if !ok || precision != float64(p.Precision) || p.Precision < 0 || p.Precision > 2 || unit == "count" && p.Precision != 0 {
		return invalidDefinition("invalid metric precision")
	}
	if widget.Type == "radar" {
		return validateDashboardRadar(widget, p, unit)
	}
	if widget.Type == "bar" || widget.Type == "heatmap" || widget.Type == "line" || widget.Type == "table" {
		return validateDashboardAdvancedMetric(widget, p, unit)
	}
	if widget.Type != "gauge" {
		return nil
	}
	min, minOK := widget.Properties["minimum"].(float64)
	max, maxOK := widget.Properties["maximum"].(float64)
	if !minOK || !maxOK || min < 0 || max <= min {
		return invalidDefinition("gauge maximum must exceed its minimum")
	}
	if !slices.Contains([]string{"at_least", "at_most", "informational"}, p.GoalMode) {
		return invalidDefinition("invalid gauge goal mode")
	}
	if _, ok := widget.Properties["target"]; !ok {
		return invalidDefinition("gauge target is required")
	}
	if p.Target != nil && (*p.Target <= 0 || *p.Target < min || *p.Target > max) {
		return invalidDefinition("gauge target must be positive and within its scale")
	}
	if p.GoalMode != "informational" && p.Target == nil {
		return invalidDefinition("gauge goal requires a target")
	}
	if unit == "count" && (min != math.Trunc(min) || max != math.Trunc(max) || p.Target != nil && *p.Target != math.Trunc(*p.Target)) {
		return invalidDefinition("count metrics require whole numbers")
	}
	for _, key := range []string{"showValue", "showTarget", "showPercentage"} {
		if _, ok := widget.Properties[key].(bool); !ok {
			return invalidDefinition(key + " must be a boolean")
		}
	}
	thresholds, ok := widget.Properties["thresholds"].([]interface{})
	if !ok || len(thresholds) > 8 {
		return invalidDefinition("gauge supports up to eight thresholds")
	}
	previous := math.Inf(-1)
	for i, threshold := range p.Thresholds {
		entry, ok := thresholds[i].(map[string]interface{})
		if !ok {
			return invalidDefinition("invalid gauge threshold")
		}
		if _, ok := entry["value"].(float64); !ok {
			return invalidDefinition("threshold value must be numeric")
		}
		if threshold.Value < min || threshold.Value > max || threshold.Value <= previous || !slices.Contains([]string{"muted", "info", "success", "warning", "destructive"}, threshold.Tone) || unit == "count" && threshold.Value != math.Trunc(threshold.Value) {
			return invalidDefinition("thresholds must increase within the gauge scale")
		}
		previous = threshold.Value
	}
	return nil
}

func validateDashboardTableMetric(widget dashboardWidget, p dashboardMetricProperties) error {
	for _, key := range []string{"metric", "format", "precision"} {
		if _, exists := widget.Properties[key]; exists {
			return invalidDefinition("table metrics are configured per column")
		}
	}
	for _, key := range []string{"title", "description", "startDate", "endDate", "projectId", "contextId", "tag"} {
		text, ok := widget.Properties[key].(string)
		if !ok || len(text) > 500 {
			return invalidDefinition("invalid table " + key)
		}
	}
	if p.Period != "inherit" && p.Period != "custom" && !slices.Contains(dashboardPeriods, p.Period) {
		return invalidDefinition("invalid table period")
	}
	if p.Period == "custom" {
		start, e1 := time.Parse(time.DateOnly, p.StartDate)
		end, e2 := time.Parse(time.DateOnly, p.EndDate)
		if e1 != nil || e2 != nil || end.Before(start) {
			return invalidDefinition("invalid table date range")
		}
	}
	if !slices.Contains([]string{"include", "exclude", "only"}, p.Archived) {
		return invalidDefinition("invalid archived filter")
	}
	if _, ok := widget.Properties["includeRunning"].(bool); !ok {
		return invalidDefinition("includeRunning must be a boolean")
	}
	return validateDashboardAdvancedMetric(widget, p, "")
}

func validateDashboardRadar(widget dashboardWidget, p dashboardMetricProperties, unit string) error {
	if p.AxesBy != "weekday" && p.AxesBy != "project" && p.AxesBy != "context" {
		return invalidDefinition("invalid radar axis grouping")
	}
	keys := p.AxisKeys
	if p.AxesBy == "weekday" {
		if len(keys) != 0 {
			return invalidDefinition("weekday radar axes are calendar-defined")
		}
	} else if len(keys) < 3 || len(keys) > 8 {
		return invalidDefinition("choose between three and eight radar axes")
	}
	seen := map[string]bool{}
	for _, key := range keys {
		valid := false
		if p.AxesBy == "project" {
			valid = key == "project:none" || strings.HasPrefix(key, "project:") && len(key) > len("project:")
		}
		if p.AxesBy == "context" {
			valid = strings.HasPrefix(key, "context:") && len(key) > len("context:")
		}
		if !valid || seen[key] {
			return invalidDefinition("radar axis keys must be unique and match the selected grouping")
		}
		seen[key] = true
	}
	for key, value := range map[string]*float64{"scaleMaximum": p.ScaleMaximum, "targetValue": p.TargetValue} {
		if value != nil && (*value < 0 || !isFinite(*value) || (unit == "count" && *value != math.Trunc(*value)) || (unit == "percent" && *value > 100)) {
			return invalidDefinition("invalid radar " + key)
		}
	}
	if p.ScaleMaximum != nil && *p.ScaleMaximum == 0 {
		return invalidDefinition("radar scale maximum must be positive")
	}
	for _, key := range []string{"showValues", "showGrid", "showFill", "showSummary"} {
		if _, ok := widget.Properties[key].(bool); !ok {
			return invalidDefinition("radar " + key + " must be a boolean")
		}
	}
	return nil
}

func validateDashboardAdvancedMetric(widget dashboardWidget, p dashboardMetricProperties, unit string) error {
	switch widget.Type {
	case "bar":
		group, _ := widget.Properties["groupBy"].(string)
		bucket, _ := widget.Properties["bucket"].(string)
		if !slices.Contains([]string{"project", "context", "time"}, group) {
			return invalidDefinition("invalid bar grouping")
		}
		if group == "time" {
			if !slices.Contains([]string{"day", "week", "month"}, bucket) {
				return invalidDefinition("time bars require a day, week, or month bucket")
			}
			sortBy, direction := "", ""
			if p.Sort != nil {
				sortBy, _ = p.Sort["by"].(string)
				direction, _ = p.Sort["direction"].(string)
			}
			if sortBy != "time" || direction != "asc" || p.Limit != nil || p.ShowOther {
				return invalidDefinition("time bars require chronological order without a limit")
			}
		} else if bucket != "" {
			return invalidDefinition("entity bars cannot use a time bucket")
		} else {
			sortBy, direction := "value", "desc"
			if p.Sort != nil {
				if v, ok := p.Sort["by"].(string); ok {
					sortBy = v
				}
				if v, ok := p.Sort["direction"].(string); ok {
					direction = v
				}
			}
			if !slices.Contains([]string{"value", "label"}, sortBy) || !slices.Contains([]string{"asc", "desc"}, direction) || p.Limit == nil || *p.Limit < 2 || *p.Limit > 50 {
				return invalidDefinition("invalid bar ranking settings")
			}
			if p.ShowOther && p.Metric != "duration.sum" {
				return invalidDefinition("other is available only for tracked time")
			}
		}
		for _, key := range []string{"showOther", "showValues", "showGrid", "showSummary"} {
			if _, ok := widget.Properties[key].(bool); !ok {
				return invalidDefinition("invalid bar display settings")
			}
		}
	case "heatmap":
		if values, ok := widget.Properties["thresholds"].([]interface{}); ok {
			for _, value := range values {
				number, ok := value.(float64)
				if !ok {
					return invalidDefinition("invalid heatmap threshold")
				}
				p.ThresholdsHeat = append(p.ThresholdsHeat, number)
			}
		}
		if !slices.Contains([]string{"auto", "fixed"}, p.ScaleMode) {
			return invalidDefinition("invalid heatmap scale mode")
		}
		if p.ScaleMode == "auto" && len(p.ThresholdsHeat) != 0 || p.ScaleMode == "fixed" && len(p.ThresholdsHeat) != 3 {
			return invalidDefinition("heatmap fixed scale requires three thresholds")
		}
		previous := 0.0
		for _, threshold := range p.ThresholdsHeat {
			if !isFinite(threshold) || threshold <= previous || unit == "percent" && threshold >= 100 || unit == "count" && threshold != math.Trunc(threshold) {
				return invalidDefinition("heatmap thresholds must be positive and increasing")
			}
			previous = threshold
		}
		for _, key := range []string{"showMonthLabels", "showDayLabels", "showLegend", "showSummary"} {
			if _, ok := widget.Properties[key].(bool); !ok {
				return invalidDefinition("invalid heatmap display settings")
			}
		}
	case "line":
		if !slices.Contains([]string{"day", "week", "month"}, p.Bucket) || !slices.Contains([]string{"none", "project", "context"}, p.SeriesBy) || p.SeriesLimit < 1 || p.SeriesLimit > 5 || !slices.Contains([]string{"periodic", "cumulative"}, p.ValueMode) || !slices.Contains([]string{"points", "linear", "smooth"}, p.RenderMode) {
			return invalidDefinition("invalid line chart settings")
		}
		if p.ValueMode == "cumulative" && p.Metric != "duration.sum" {
			return invalidDefinition("cumulative mode requires tracked time")
		}
		if p.TargetLine != nil && (!isFinite(*p.TargetLine) || *p.TargetLine < 0 || unit == "percent" && *p.TargetLine > 100 || unit == "count" && *p.TargetLine != math.Trunc(*p.TargetLine)) {
			return invalidDefinition("invalid line target")
		}
		for _, key := range []string{"showPoints", "showLegend", "showGrid", "showSummary"} {
			if _, ok := widget.Properties[key].(bool); !ok {
				return invalidDefinition("invalid line display settings")
			}
		}
	case "table":
		group, _ := widget.Properties["groupBy"].(string)
		if group != "project" && group != "context" {
			return invalidDefinition("invalid table grouping")
		}
		if len(p.Columns) < 1 || len(p.Columns) > 7 {
			return invalidDefinition("tables require between one and seven columns")
		}
		seen, seenMetrics, visible := map[string]bool{}, map[string]bool{}, false
		for _, column := range p.Columns {
			columnUnit := dashboardMetricUnit(column.Metric)
			if column.ID == "" || len(column.ID) > 64 || column.ID == "group" || seen[column.ID] || seenMetrics[column.Metric] || columnUnit == "" {
				return invalidDefinition("invalid table column")
			}
			seen[column.ID] = true
			seenMetrics[column.Metric] = true
			visible = visible || column.Visible
			formats := map[string][]string{"ms": {"duration", "hours"}, "count": {"number"}, "percent": {"percent"}}
			if !slices.Contains(formats[columnUnit], column.Format) || column.Precision < 0 || column.Precision > 2 || columnUnit == "count" && column.Precision != 0 || column.ShowBar && column.Metric != "duration.sum" && column.Metric != "time_share.ratio" {
				return invalidDefinition("invalid table column settings")
			}
		}
		if !visible || !slices.Contains([]int{10, 25, 50, 100}, p.PageSize) {
			return invalidDefinition("table requires a visible column and valid page size")
		}
		columnID, direction := "", "desc"
		if p.Sort != nil {
			columnID, _ = p.Sort["columnId"].(string)
			if v, ok := p.Sort["direction"].(string); ok {
				direction = v
			}
		}
		validSort := columnID == "group"
		for _, column := range p.Columns {
			if column.ID == columnID && column.Visible {
				validSort = true
			}
		}
		if !validSort || !slices.Contains([]string{"asc", "desc"}, direction) {
			return invalidDefinition("invalid table default sort")
		}
		if _, ok := widget.Properties["showSummary"].(bool); !ok {
			return invalidDefinition("invalid table display settings")
		}
	}
	return nil
}

func isFinite(value float64) bool { return !math.IsNaN(value) && !math.IsInf(value, 0) }

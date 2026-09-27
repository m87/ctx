package core

import (
	"encoding/json"
	"math"
	"slices"
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
	if p.Period != "inherit" && p.Period != "custom" && !slices.Contains(dashboardPeriods, p.Period) {
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

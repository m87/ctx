package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateDashboardMetricRejectsInvalidKPIAndGaugeSettings(t *testing.T) {
	base := map[string]interface{}{
		"title": "Tracked time", "description": "", "metric": "duration.sum", "period": "custom",
		"startDate": "2026-09-01", "endDate": "2026-09-30", "projectId": "", "contextId": "", "tag": "",
		"archived": "include", "includeRunning": false, "format": "duration", "precision": float64(0),
		"minimum": float64(0), "maximum": float64(40), "target": float64(30), "goalMode": "at_least",
		"thresholds": []interface{}{}, "showValue": true, "showTarget": true, "showPercentage": true,
	}
	cases := []struct {
		name   string
		mutate func(map[string]interface{})
	}{
		{"unknown metric", func(p map[string]interface{}) { p["metric"] = "missing" }},
		{"wrong format", func(p map[string]interface{}) { p["format"] = "number" }},
		{"invalid custom date", func(p map[string]interface{}) { p["endDate"] = "not-a-date" }},
		{"reversed custom dates", func(p map[string]interface{}) { p["endDate"] = "2026-08-31" }},
		{"invalid gauge scale", func(p map[string]interface{}) { p["maximum"] = float64(0) }},
		{"negative minimum", func(p map[string]interface{}) { p["minimum"] = float64(-1) }},
		{"target outside scale", func(p map[string]interface{}) { p["target"] = float64(41) }},
		{"missing required target", func(p map[string]interface{}) { p["target"] = nil }},
		{"unordered thresholds", func(p map[string]interface{}) {
			p["thresholds"] = []interface{}{map[string]interface{}{"value": float64(20), "tone": "info"}, map[string]interface{}{"value": float64(10), "tone": "success"}}
		}},
		{"too many thresholds", func(p map[string]interface{}) {
			values := make([]interface{}, 9)
			for i := range values {
				values[i] = map[string]interface{}{"value": float64(i + 1), "tone": "info"}
			}
			p["thresholds"] = values
		}},
		{"fractional count", func(p map[string]interface{}) {
			p["metric"], p["format"], p["minimum"] = "sessions.count", "number", float64(0.5)
		}},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			properties := make(map[string]interface{}, len(base))
			for key, value := range base {
				properties[key] = value
			}
			testCase.mutate(properties)
			err := validateDashboardMetric(dashboardWidget{Type: "gauge", Properties: properties})
			require.Error(t, err)
		})
	}
}

func TestValidateDashboardMetricAcceptsValidKPIMetricSettings(t *testing.T) {
	properties := map[string]interface{}{
		"title": "Sessions", "description": "", "metric": "sessions.count", "period": "today",
		"startDate": "", "endDate": "", "projectId": "", "contextId": "", "tag": "",
		"archived": "include", "includeRunning": true, "format": "number", "precision": float64(0),
	}
	err := validateDashboardMetric(dashboardWidget{Type: "kpi", Properties: properties})
	require.NoError(t, err)
}

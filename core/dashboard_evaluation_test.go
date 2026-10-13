package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestEvaluateDashboardTimeShareUsesDashboardScopeAsDenominator(t *testing.T) {
	test := NewTestContextManager()
	now := time.Date(2026, time.September, 27, 16, 0, 0, 0, time.UTC)
	test.TimeProvider.Set(now)
	archivedStart := time.Date(2026, time.September, 27, 8, 0, 0, 0, time.UTC)
	archivedEnd := archivedStart.Add(time.Hour)
	activeStart := archivedStart
	activeEnd := activeStart.Add(3 * time.Hour)
	test.Contexts.Seed(
		&Context{Id: "active", WorkspaceId: TestWorkspaceID},
		&Context{Id: "archived", WorkspaceId: TestWorkspaceID, Archived: true},
	)
	test.Intervals.Seed(
		&Interval{Id: "active-session", ContextId: "active", WorkspaceId: TestWorkspaceID, Start: &activeStart, End: &activeEnd},
		&Interval{Id: "archived-session", ContextId: "archived", WorkspaceId: TestWorkspaceID, Start: &archivedStart, End: &archivedEnd},
	)
	for _, tc := range []struct {
		archived string
		want     float64
	}{{"only", 25}, {"exclude", 75}, {"include", 100}} {
		t.Run(tc.archived, func(t *testing.T) {
			definition := metricDefinition("time_share.ratio", tc.archived, "")
			result, err := test.Manager.EvaluateDashboard(DashboardEvaluationRequest{
				WorkspaceID: TestWorkspaceID,
				Definition:  definition,
				TimeZone:    "UTC",
			})
			require.NoError(t, err)
			require.NotNil(t, result.Widgets["metric"].Value)
			require.Equal(t, tc.want, *result.Widgets["metric"].Value)
		})
	}
}

func TestEvaluateDashboardRadarGroupsWeekdayValuesOnSharedAxes(t *testing.T) {
	test := NewTestContextManager()
	endOfPeriod := time.Date(2026, time.September, 28, 0, 0, 0, 0, time.UTC)
	test.TimeProvider.Set(time.Date(2026, time.September, 27, 16, 0, 0, 0, time.UTC))
	day := time.Date(2026, time.September, 21, 8, 0, 0, 0, time.UTC)
	dayEnd := day.Add(2 * time.Hour)
	test.Intervals.Seed(&Interval{Id: "monday", ContextId: TestContextID, WorkspaceId: TestWorkspaceID, Start: &day, End: &dayEnd})
	definition, err := json.Marshal(map[string]interface{}{
		"version": 1, "query": "", "widgets": []interface{}{map[string]interface{}{
			"id": "radar", "type": "radar", "layout": map[string]int{"x": 0, "y": 0, "width": 6, "height": 6}, "query": "",
			"properties": map[string]interface{}{
				"title": "Week", "description": "", "metric": "duration.sum", "period": "last_7_days", "startDate": "", "endDate": "",
				"projectId": "", "contextId": "", "tag": "", "archived": "include", "includeRunning": false, "format": "duration", "precision": 1,
				"axesBy": "weekday", "axisKeys": []string{}, "scaleMaximum": nil, "targetValue": nil, "showValues": false, "showGrid": true, "showFill": true, "showSummary": true,
			},
		}},
	})
	require.NoError(t, err)
	result, err := test.Manager.EvaluateDashboard(DashboardEvaluationRequest{WorkspaceID: TestWorkspaceID, Definition: definition, TimeZone: "UTC", FirstDayOfWeek: 1})
	require.NoError(t, err)
	widget := result.Widgets["radar"]
	require.Equal(t, "radar", widget.Kind)
	require.NotNil(t, widget.Radar)
	require.Len(t, widget.Radar.Axes, 7)
	require.Equal(t, "Monday", widget.Radar.Axes[0].Label)
	require.Equal(t, float64(2*time.Hour/time.Millisecond), *widget.Radar.Axes[0].Value)
	require.Equal(t, float64(0), *widget.Radar.Axes[1].Value)
	require.Equal(t, float64(2*time.Hour/time.Millisecond), *widget.Radar.TotalValue)
	require.True(t, *widget.HasActivity)
	require.True(t, endOfPeriod.Equal(widget.End))
}

func TestEvaluateDashboardPassthroughQueryDoesNotFilterOrReject(t *testing.T) {
	test := NewTestContextManager()
	now := time.Date(2026, time.September, 27, 16, 0, 0, 0, time.UTC)
	test.TimeProvider.Set(now)
	start := time.Date(2026, time.September, 27, 8, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)
	test.Intervals.Seed(&Interval{Id: "session", ContextId: TestContextID, WorkspaceId: TestWorkspaceID, Start: &start, End: &end})
	for _, query := range []string{"", "arbitrary future query text"} {
		definition := metricDefinition("duration.sum", "include", query)
		result, err := test.Manager.EvaluateDashboard(DashboardEvaluationRequest{
			WorkspaceID: TestWorkspaceID,
			Definition:  definition,
			TimeZone:    "UTC",
		})
		require.NoError(t, err)
		require.NotNil(t, result.Widgets["metric"].Value)
		require.Equal(t, float64(time.Hour/time.Millisecond), *result.Widgets["metric"].Value)
	}
}

func TestEvaluateDashboardCountsOverlappingIntervalsIndependently(t *testing.T) {
	test := NewTestContextManager()
	test.TimeProvider.Set(time.Date(2026, time.September, 27, 16, 0, 0, 0, time.UTC))
	first := time.Date(2026, time.September, 27, 9, 0, 0, 0, time.UTC)
	firstEnd := first.Add(time.Hour)
	second := first.Add(30 * time.Minute)
	secondEnd := second.Add(time.Hour)
	test.Intervals.Seed(
		&Interval{Id: "one", ContextId: TestContextID, WorkspaceId: TestWorkspaceID, Start: &first, End: &firstEnd},
		&Interval{Id: "two", ContextId: TestContextID, WorkspaceId: TestWorkspaceID, Start: &second, End: &secondEnd},
	)
	definition := metricDefinition("duration.sum", "include", "")
	result, err := test.Manager.EvaluateDashboard(DashboardEvaluationRequest{WorkspaceID: TestWorkspaceID, Definition: definition, TimeZone: "UTC"})
	require.NoError(t, err)
	require.Equal(t, float64(2*time.Hour/time.Millisecond), *result.Widgets["metric"].Value)
}

func TestEvaluateDashboardBarAndHeatmapPayloads(t *testing.T) {
	test := NewTestContextManager()
	now := time.Date(2026, time.September, 27, 16, 0, 0, 0, time.UTC)
	test.TimeProvider.Set(now)
	start := time.Date(2026, time.September, 26, 23, 30, 0, 0, time.UTC)
	end := start.Add(time.Hour)
	test.Intervals.Seed(&Interval{Id: "overnight", ContextId: TestContextID, WorkspaceId: TestWorkspaceID, Start: &start, End: &end})
	common := map[string]interface{}{"title": "", "description": "", "metric": "duration.sum", "startDate": "", "endDate": "", "projectId": "", "contextId": "", "tag": "", "archived": "include", "includeRunning": false, "format": "duration", "precision": 0}
	bar := map[string]interface{}{"id": "bar", "type": "bar", "layout": map[string]int{"x": 0, "y": 0, "width": 8, "height": 5}, "query": "", "properties": map[string]interface{}{}}
	barProperties := bar["properties"].(map[string]interface{})
	for k, v := range common {
		barProperties[k] = v
	}
	for k, v := range map[string]interface{}{"period": "today", "groupBy": "context", "bucket": nil, "orientation": "horizontal", "sort": map[string]interface{}{"by": "value", "direction": "desc"}, "limit": 10, "showOther": false, "showValues": true, "showGrid": true, "showSummary": true} {
		barProperties[k] = v
	}
	heat := map[string]interface{}{"id": "heat", "type": "heatmap", "layout": map[string]int{"x": 8, "y": 0, "width": 8, "height": 5}, "query": "", "properties": map[string]interface{}{}}
	heatProperties := heat["properties"].(map[string]interface{})
	for k, v := range common {
		heatProperties[k] = v
	}
	for k, v := range map[string]interface{}{"period": "last_3_months", "scaleMode": "auto", "thresholds": []interface{}{}, "showMonthLabels": true, "showDayLabels": true, "showLegend": true, "showSummary": true} {
		heatProperties[k] = v
	}
	definition, _ := json.Marshal(map[string]interface{}{"version": 1, "query": "", "widgets": []interface{}{bar, heat}})
	result, err := test.Manager.EvaluateDashboard(DashboardEvaluationRequest{WorkspaceID: TestWorkspaceID, Definition: definition, TimeZone: "UTC", FirstDayOfWeek: 1})
	require.NoError(t, err)
	barResult := result.Widgets["bar"]
	require.Equal(t, "categories", barResult.Kind)
	require.Len(t, barResult.Categories.Categories, 1)
	require.Equal(t, float64(30*time.Minute/time.Millisecond), *barResult.Categories.Categories[0].Value)
	heatResult := result.Widgets["heat"]
	require.Equal(t, "heatmap", heatResult.Kind)
	require.Len(t, heatResult.Heatmap.Days, 92)
	require.Equal(t, "future", heatResult.Heatmap.Days[len(heatResult.Heatmap.Days)-1].Status)
	require.Equal(t, 2, heatResult.Heatmap.ActiveDays)
}

func tableDashboardDefinition(t *testing.T, groupBy string, sort map[string]interface{}) json.RawMessage {
	t.Helper()
	column := func(id, metric, label, format string, precision float64, showBar bool) map[string]interface{} {
		return map[string]interface{}{"id": id, "metric": metric, "label": label, "visible": true, "format": format, "precision": precision, "showBar": showBar}
	}
	properties := map[string]interface{}{
		"title": "", "description": "", "period": "today", "startDate": "", "endDate": "", "projectId": "", "contextId": "", "tag": "", "archived": "include", "includeRunning": false,
		"groupBy": groupBy,
		"columns": []interface{}{
			column("time", "duration.sum", "Time", "duration", 1, true),
			column("average", "duration.mean", "Average", "duration", 1, false),
			column("sessions", "sessions.count", "Sessions", "number", 0, false),
		},
		"sort": sort, "pageSize": 10, "showSummary": true,
	}
	definition, err := json.Marshal(map[string]interface{}{"version": 1, "query": "", "widgets": []interface{}{map[string]interface{}{"id": "table", "type": "table", "layout": map[string]int{"x": 0, "y": 0, "width": 10, "height": 6}, "query": "", "properties": properties}}})
	require.NoError(t, err)
	return definition
}

func seedTableContexts(test *TestContextManager, extraEmpty int) {
	day := time.Date(2026, time.September, 27, 8, 0, 0, 0, time.UTC)
	names := map[string]string{"ctx-a": "alpha", "ctx-b": "Beta", "ctx-c": "Charlie", "ctx-d": "Delta", "ctx-empty": "Echo"}
	minutes := map[string]int{"ctx-a": 30, "ctx-b": 90, "ctx-c": 90, "ctx-d": 45, "ctx-empty": 0}
	contexts := []*Context{}
	intervals := []*Interval{}
	for _, id := range []string{"ctx-a", "ctx-b", "ctx-c", "ctx-d", "ctx-empty"} {
		contexts = append(contexts, &Context{Id: id, Name: names[id], WorkspaceId: TestWorkspaceID, Status: "inactive"})
		if minutes[id] > 0 {
			start := day
			end := start.Add(time.Duration(minutes[id]) * time.Minute)
			intervals = append(intervals, &Interval{Id: "interval-" + id, ContextId: id, WorkspaceId: TestWorkspaceID, Start: &start, End: &end})
		}
	}
	for index := 0; index < extraEmpty; index++ {
		id := fmt.Sprintf("ctx-extra-%02d", index)
		start := day.Add(time.Hour + time.Duration(index)*15*time.Minute)
		end := start.Add(time.Duration(index+1) * time.Minute)
		contexts = append(contexts, &Context{Id: id, Name: fmt.Sprintf("Extra %02d", index), WorkspaceId: TestWorkspaceID, Status: "inactive"})
		intervals = append(intervals, &Interval{Id: "interval-" + id, ContextId: id, WorkspaceId: TestWorkspaceID, Start: &start, End: &end})
	}
	test.Contexts.Seed(contexts...)
	test.Intervals.Seed(intervals...)
}

func evaluateTable(t *testing.T, test *TestContextManager, definition json.RawMessage) *DashboardTablePayload {
	t.Helper()
	result, err := test.Manager.EvaluateDashboard(DashboardEvaluationRequest{WorkspaceID: TestWorkspaceID, Definition: definition, TimeZone: "UTC"})
	require.NoError(t, err)
	widget := result.Widgets["table"]
	require.Empty(t, widget.Error)
	require.Equal(t, "table", widget.Kind)
	return widget.Table
}

func tableRowLabels(table *DashboardTablePayload) []string {
	labels := make([]string, 0, len(table.Rows))
	for _, row := range table.Rows {
		labels = append(labels, row.Label)
	}
	return labels
}

func TestEvaluateDashboardTableUsesColumnMetrics(t *testing.T) {
	test := NewTestContextManager()
	test.TimeProvider.Set(time.Date(2026, time.September, 27, 16, 0, 0, 0, time.UTC))
	start := time.Date(2026, time.September, 27, 8, 0, 0, 0, time.UTC)
	end := start.Add(90 * time.Minute)
	test.Intervals.Seed(&Interval{Id: "session", ContextId: TestContextID, WorkspaceId: TestWorkspaceID, Start: &start, End: &end})
	table := evaluateTable(t, test, tableDashboardDefinition(t, "context", map[string]interface{}{"columnId": "time", "direction": "desc"}))
	require.Len(t, table.Columns, 3)
	require.Len(t, table.Rows, 1)
	require.Equal(t, 1, table.TotalRows)
	require.Equal(t, float64(90*time.Minute/time.Millisecond), *table.Summary["time"])
	require.Equal(t, float64(90*time.Minute/time.Millisecond), *table.ColumnMaxima["time"])
}

func TestEvaluateDashboardTableReturnsEveryRowSortedByTheSavedSort(t *testing.T) {
	newTest := func(extra int) *TestContextManager {
		test := NewTestContextManager()
		test.TimeProvider.Set(time.Date(2026, time.September, 27, 16, 0, 0, 0, time.UTC))
		seedTableContexts(test, extra)
		return test
	}
	sorted := func(column, direction string) map[string]interface{} {
		return map[string]interface{}{"columnId": column, "direction": direction}
	}

	t.Run("duration descending keeps ties ordered by key and returns more rows than the page size", func(t *testing.T) {
		table := evaluateTable(t, newTest(8), tableDashboardDefinition(t, "context", sorted("time", "desc")))
		require.Len(t, table.Rows, 12)
		require.Equal(t, 12, table.TotalRows)
		require.Equal(t, []string{"Beta", "Charlie", "Delta", "alpha"}, tableRowLabels(table)[:4])
	})

	t.Run("duration ascending keeps ties ordered by key too", func(t *testing.T) {
		table := evaluateTable(t, newTest(0), tableDashboardDefinition(t, "context", sorted("time", "asc")))
		require.Equal(t, []string{"alpha", "Delta", "Beta", "Charlie"}, tableRowLabels(table))
	})

	t.Run("group sort ignores case", func(t *testing.T) {
		table := evaluateTable(t, newTest(0), tableDashboardDefinition(t, "context", sorted("group", "asc")))
		require.Equal(t, []string{"alpha", "Beta", "Charlie", "Delta"}, tableRowLabels(table))
		table = evaluateTable(t, newTest(0), tableDashboardDefinition(t, "context", sorted("group", "desc")))
		require.Equal(t, []string{"Delta", "Charlie", "Beta", "alpha"}, tableRowLabels(table))
	})

	t.Run("the summary and column maxima cover all rows", func(t *testing.T) {
		table := evaluateTable(t, newTest(0), tableDashboardDefinition(t, "context", sorted("time", "desc")))
		require.Equal(t, float64(255*time.Minute/time.Millisecond), *table.Summary["time"])
		require.Equal(t, float64(90*time.Minute/time.Millisecond), *table.ColumnMaxima["time"])
	})
}

func TestEvaluateDashboardTablePayloadHasNoServerPaging(t *testing.T) {
	test := NewTestContextManager()
	test.TimeProvider.Set(time.Date(2026, time.September, 27, 16, 0, 0, 0, time.UTC))
	seedTableContexts(test, 0)
	result, err := test.Manager.EvaluateDashboard(DashboardEvaluationRequest{WorkspaceID: TestWorkspaceID, Definition: tableDashboardDefinition(t, "context", map[string]interface{}{"columnId": "time", "direction": "desc"}), TimeZone: "UTC"})
	require.NoError(t, err)
	encoded, err := json.Marshal(result.Widgets["table"].Table)
	require.NoError(t, err)
	var payload map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(encoded, &payload))
	for _, key := range []string{"columns", "rows", "summary", "columnMaxima", "totalRows"} {
		require.Contains(t, payload, key)
	}
	for _, key := range []string{"page", "pageSize", "pageCount", "sort"} {
		require.NotContains(t, payload, key)
	}
}

func TestEvaluateDashboardTableIgnoresLegacyViewState(t *testing.T) {
	var request DashboardEvaluationRequest
	require.NoError(t, json.Unmarshal([]byte(`{"workspaceId":"w","definition":{},"timeZone":"UTC","viewState":{"table":{"page":3,"sort":{"columnId":"time","direction":"asc"}}}}`), &request))
	require.Equal(t, "w", request.WorkspaceID)
}

func metricDefinition(metric, archived, query string) json.RawMessage {
	definition, _ := json.Marshal(map[string]interface{}{
		"version": 1,
		"query":   "",
		"widgets": []interface{}{map[string]interface{}{
			"id": "metric", "type": "kpi", "layout": map[string]int{"x": 0, "y": 0, "width": 4, "height": 3}, "query": query,
			"properties": map[string]interface{}{
				"title": "Time", "description": "", "metric": metric, "period": "today", "startDate": "", "endDate": "",
				"projectId": "", "contextId": "", "tag": "", "archived": archived, "includeRunning": false,
				"format": map[string]string{"duration.sum": "duration", "time_share.ratio": "percent"}[metric], "precision": 0,
			},
		}},
	})
	return definition
}

func TestAggregateDashboardIntervalsClipsHalfOpenRangesAndCountsUniqueSessions(t *testing.T) {
	period := dashboardPeriodRange{
		start: time.Date(2026, time.September, 27, 0, 0, 0, 0, time.UTC),
		end:   time.Date(2026, time.September, 28, 0, 0, 0, 0, time.UTC),
	}
	boundaryEnd := period.start
	boundaryStart := period.end
	start := period.start.Add(-30 * time.Minute)
	end := period.start.Add(30 * time.Minute)
	invalidStart := period.start.Add(time.Hour)
	interval := &Interval{Id: "duplicate", ContextId: "context", Start: &start, End: &end}
	aggregate := aggregateDashboardIntervals([]*Interval{
		{Id: "ends-at-start", ContextId: "context", Start: &start, End: &boundaryEnd},
		{Id: "starts-at-end", ContextId: "context", Start: &boundaryStart, End: &boundaryStart},
		interval,
		interval,
		{Id: "invalid", ContextId: "context", Start: &invalidStart, End: &invalidStart},
		{Id: "missing-start", ContextId: "context", End: &end},
	}, map[string]bool{"context": true}, period, period.end, false)
	require.Equal(t, float64(30*time.Minute/time.Millisecond), aggregate.duration)
	require.Equal(t, 1, aggregate.sessions)
	require.Equal(t, 3, aggregate.ignored)
}

func TestAggregateDashboardIntervalsUsesAsOfAndPeriodEndForRunningSessions(t *testing.T) {
	period := dashboardPeriodRange{
		start: time.Date(2026, time.September, 27, 0, 0, 0, 0, time.UTC),
		end:   time.Date(2026, time.September, 28, 0, 0, 0, 0, time.UTC),
	}
	start := period.start.Add(23 * time.Hour)
	interval := &Interval{Id: "running", ContextId: "context", Start: &start, Status: IntervalStatusActive}
	atAsOf := aggregateDashboardIntervals([]*Interval{interval}, map[string]bool{"context": true}, period, start.Add(30*time.Minute), true)
	require.Equal(t, float64(30*time.Minute/time.Millisecond), atAsOf.duration)
	afterEnd := aggregateDashboardIntervals([]*Interval{interval}, map[string]bool{"context": true}, period, period.end.Add(time.Hour), true)
	require.Equal(t, float64(time.Hour/time.Millisecond), afterEnd.duration)
	excluded := aggregateDashboardIntervals([]*Interval{interval}, map[string]bool{"context": true}, period, period.end, false)
	require.Zero(t, excluded.duration)
}

func TestDashboardAggregateEmptyValuesPreserveZeroAndNull(t *testing.T) {
	aggregate := dashboardAggregate{}
	require.Equal(t, float64(0), *aggregate.value("duration.sum"))
	require.Equal(t, float64(0), *aggregate.value("sessions.count"))
	require.Nil(t, aggregate.value("duration.mean"))
	require.Nil(t, aggregate.value("duration.max"))
}

func TestDashboardMetricPeriodUsesConfiguredWeekAndInclusiveCustomEndDate(t *testing.T) {
	location := time.UTC
	now := time.Date(2026, time.September, 27, 16, 0, 0, 0, location)
	for _, tc := range []struct {
		name        string
		period      string
		firstDay    int
		start, end  time.Time
		customStart string
		customEnd   string
	}{
		{name: "Monday week", period: "this_week", firstDay: 1, start: time.Date(2026, 9, 21, 0, 0, 0, 0, location), end: time.Date(2026, 9, 28, 0, 0, 0, 0, location)},
		{name: "Sunday week", period: "this_week", firstDay: 0, start: time.Date(2026, 9, 27, 0, 0, 0, 0, location), end: time.Date(2026, 10, 4, 0, 0, 0, 0, location)},
		{name: "inclusive custom end", period: "custom", firstDay: 1, customStart: "2026-09-25", customEnd: "2026-09-27", start: time.Date(2026, 9, 25, 0, 0, 0, 0, location), end: time.Date(2026, 9, 28, 0, 0, 0, 0, location)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			request := DashboardEvaluationRequest{TimeZone: "UTC", FirstDayOfWeek: tc.firstDay}
			period, err := dashboardMetricPeriod(dashboardMetricProperties{Period: tc.period, StartDate: tc.customStart, EndDate: tc.customEnd}, "", request, now)
			require.NoError(t, err)
			require.Equal(t, tc.start, period.start)
			require.Equal(t, tc.end, period.end)
		})
	}
}

func TestDashboardMetricPeriodUsesDailyScopeAndWarsawDSTBoundaries(t *testing.T) {
	warsaw, err := time.LoadLocation("Europe/Warsaw")
	require.NoError(t, err)
	for _, tc := range []struct {
		date     string
		duration time.Duration
	}{
		{date: "2026-03-29", duration: 23 * time.Hour},
		{date: "2026-10-25", duration: 25 * time.Hour},
	} {
		now := time.Date(2026, 11, 2, 12, 0, 0, 0, warsaw)
		request := DashboardEvaluationRequest{TimeZone: "Europe/Warsaw", Scope: DashboardScope{Date: tc.date}}
		period, err := dashboardMetricPeriod(dashboardMetricProperties{Period: "last_month"}, "today", request, now)
		require.NoError(t, err)
		require.Equal(t, tc.duration, period.end.Sub(period.start))
		require.Equal(t, tc.date, period.start.Format(time.DateOnly))
	}
}

type failingDashboardQueryInterpreter struct{ err error }

func (i failingDashboardQueryInterpreter) Interpret(string) (*ContextSQLQuery, error) {
	return nil, i.err
}

func TestEvaluateDashboardPreservesRepositoryAndInterpreterErrors(t *testing.T) {
	t.Run("context repository", func(t *testing.T) {
		test := NewTestContextManager()
		test.Contexts.listByWorkspaceError = errors.New("context storage unavailable")
		_, err := test.Manager.EvaluateDashboard(DashboardEvaluationRequest{
			WorkspaceID: TestWorkspaceID, Definition: metricDefinition("duration.sum", "include", ""), TimeZone: "UTC",
		})
		require.EqualError(t, err, "context storage unavailable")
	})
	t.Run("interval repository", func(t *testing.T) {
		test := NewTestContextManager()
		test.Intervals.listError = errors.New("interval storage unavailable")
		_, err := test.Manager.EvaluateDashboard(DashboardEvaluationRequest{
			WorkspaceID: TestWorkspaceID, Definition: metricDefinition("duration.sum", "include", ""), TimeZone: "UTC",
		})
		require.EqualError(t, err, "interval storage unavailable")
	})
	t.Run("query interpreter", func(t *testing.T) {
		test := NewTestContextManager()
		test.Manager.QueryInterpreter = failingDashboardQueryInterpreter{err: errors.New("interpreter unavailable")}
		result, err := test.Manager.EvaluateDashboard(DashboardEvaluationRequest{
			WorkspaceID: TestWorkspaceID, Definition: metricDefinition("duration.sum", "include", "saved query"), TimeZone: "UTC",
		})
		require.NoError(t, err)
		require.Contains(t, result.Widgets["metric"].Error, "interpreter unavailable")
		require.Nil(t, result.Widgets["metric"].Value)
	})
}

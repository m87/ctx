package core

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type DashboardScope struct {
	ProjectID string `json:"projectId,omitempty"`
	ContextID string `json:"contextId,omitempty"`
	Date      string `json:"date,omitempty"`
}

type DashboardEvaluationRequest struct {
	WorkspaceID    string                              `json:"workspaceId"`
	Definition     json.RawMessage                     `json:"definition"`
	TimeZone       string                              `json:"timeZone"`
	FirstDayOfWeek int                                 `json:"firstDayOfWeek"`
	Scope          DashboardScope                      `json:"scope"`
	ViewState      map[string]DashboardWidgetViewState `json:"viewState,omitempty"`
}

type DashboardWidgetViewState struct {
	Page int                 `json:"page"`
	Sort *DashboardTableSort `json:"sort,omitempty"`
}

type DashboardTableSort struct {
	ColumnID  string `json:"columnId"`
	Direction string `json:"direction"`
}

type DashboardMetricResult struct {
	Kind             string                      `json:"kind,omitempty"`
	Value            *float64                    `json:"value"`
	Unit             string                      `json:"unit,omitempty"`
	TimeZone         string                      `json:"timeZone"`
	Start            time.Time                   `json:"start"`
	End              time.Time                   `json:"end"`
	Error            string                      `json:"error,omitempty"`
	IgnoredIntervals int                         `json:"ignoredIntervals"`
	Radar            *DashboardRadarPayload      `json:"radar,omitempty"`
	HasActivity      *bool                       `json:"hasActivity,omitempty"`
	Categories       *DashboardCategoriesPayload `json:"categories,omitempty"`
	Heatmap          *DashboardHeatmapPayload    `json:"heatmap,omitempty"`
	Series           *DashboardLinePayload       `json:"series,omitempty"`
	Table            *DashboardTablePayload      `json:"table,omitempty"`
}

type DashboardCategoriesPayload struct {
	Unit              string              `json:"unit"`
	Categories        []DashboardCategory `json:"categories"`
	TotalValue        *float64            `json:"totalValue"`
	TotalCategories   int                 `json:"totalCategories"`
	OmittedCategories int                 `json:"omittedCategories"`
}
type DashboardCategory struct {
	Key         string                `json:"key"`
	Label       string                `json:"label"`
	Value       *float64              `json:"value"`
	Entity      *DashboardRadarEntity `json:"entity"`
	Synthetic   bool                  `json:"synthetic"`
	MemberCount int                   `json:"memberCount"`
	Start       *time.Time            `json:"start"`
	End         *time.Time            `json:"end"`
	Status      string                `json:"status"`
}
type DashboardHeatmapPayload struct {
	Unit       string                `json:"unit"`
	Days       []DashboardHeatmapDay `json:"days"`
	TotalValue *float64              `json:"totalValue"`
	ActiveDays int                   `json:"activeDays"`
}
type DashboardHeatmapDay struct {
	Date   string    `json:"date"`
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
	Value  *float64  `json:"value"`
	Status string    `json:"status"`
}
type DashboardLinePayload struct {
	Unit          string                `json:"unit"`
	Buckets       []DashboardLineBucket `json:"buckets"`
	Series        []DashboardLineSeries `json:"series"`
	TotalValue    *float64              `json:"totalValue"`
	TotalSeries   int                   `json:"totalSeries"`
	OmittedSeries int                   `json:"omittedSeries"`
}
type DashboardLineBucket struct {
	Key    string    `json:"key"`
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
	Status string    `json:"status"`
}
type DashboardLineSeries struct {
	Key    string     `json:"key"`
	Label  string     `json:"label"`
	Values []*float64 `json:"values"`
}
type DashboardTablePayload struct {
	Columns      []DashboardTablePayloadColumn `json:"columns"`
	Rows         []DashboardTableRow           `json:"rows"`
	Summary      map[string]*float64           `json:"summary"`
	ColumnMaxima map[string]*float64           `json:"columnMaxima"`
	TotalRows    int                           `json:"totalRows"`
	Page         int                           `json:"page"`
	PageSize     int                           `json:"pageSize"`
	PageCount    int                           `json:"pageCount"`
	Sort         DashboardTableSort            `json:"sort"`
}
type DashboardTablePayloadColumn struct {
	ID     string `json:"id"`
	Metric string `json:"metric"`
	Unit   string `json:"unit"`
}
type DashboardTableRow struct {
	Key    string                `json:"key"`
	Label  string                `json:"label"`
	Entity *DashboardRadarEntity `json:"entity"`
	Values map[string]*float64   `json:"values"`
}

type DashboardRadarPayload struct {
	Unit       string               `json:"unit"`
	Axes       []DashboardRadarAxis `json:"axes"`
	TotalValue *float64             `json:"totalValue"`
}

type DashboardRadarAxis struct {
	Key          string                `json:"key"`
	Label        string                `json:"label"`
	Value        *float64              `json:"value"`
	Status       string                `json:"status"`
	ObservedDays *int                  `json:"observedDays"`
	Entity       *DashboardRadarEntity `json:"entity"`
}

type DashboardRadarEntity struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type DashboardEvaluation struct {
	AsOf    time.Time                        `json:"asOf"`
	Widgets map[string]DashboardMetricResult `json:"widgets"`
}

type dashboardPeriodRange struct{ start, end time.Time }

func (m *ContextManager) EvaluateDashboard(request DashboardEvaluationRequest) (*DashboardEvaluation, error) {
	if err := validateDashboardDefinition(request.Definition); err != nil {
		return nil, err
	}
	var definition dashboardDefinition
	if err := json.Unmarshal(request.Definition, &definition); err != nil {
		return nil, invalidDefinition("invalid dashboard definition")
	}
	widgetTypes := make(map[string]string, len(definition.Widgets))
	for _, widget := range definition.Widgets {
		widgetTypes[widget.ID] = widget.Type
	}
	for id, state := range request.ViewState {
		if widgetTypes[id] != "table" || state.Page < 1 {
			return nil, invalidDefinition("invalid table view state")
		}
		if state.Sort != nil && (state.Sort.Direction != "asc" && state.Sort.Direction != "desc" || state.Sort.ColumnID == "") {
			return nil, invalidDefinition("invalid table view sort")
		}
	}
	workspace, err := m.WorkspaceRepository.GetById(request.WorkspaceID)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, &WorkspaceNotFoundError{WorkspaceId: request.WorkspaceID}
	}
	location, err := time.LoadLocation(request.TimeZone)
	if err != nil || request.TimeZone == "" {
		return nil, invalidDefinition("invalid time zone")
	}
	if request.FirstDayOfWeek != 0 && request.FirstDayOfWeek != 1 {
		return nil, invalidDefinition("first day of week must be Sunday or Monday")
	}
	contexts, err := m.ContextRepository.ListByWorkspaceIncludingArchived(request.WorkspaceID)
	if err != nil {
		return nil, err
	}
	if err := m.validateDashboardScope(request.WorkspaceID, request.Scope); err != nil {
		return nil, err
	}
	asOf := m.TimeProvider.Now().UTC()
	result := &DashboardEvaluation{AsOf: asOf, Widgets: make(map[string]DashboardMetricResult)}
	ranges := make(map[string]dashboardPeriodRange)
	var earliest, latest time.Time
	for _, widget := range definition.Widgets {
		if widget.Type == "text" {
			continue
		}
		p, err := readDashboardMetric(widget)
		if err != nil {
			return nil, err
		}
		period, err := dashboardMetricPeriod(p, definition.Period, request, asOf.In(location))
		if err != nil {
			return nil, err
		}
		ranges[widget.ID] = period
		if earliest.IsZero() || period.start.Before(earliest) {
			earliest = period.start
		}
		if latest.IsZero() || period.end.After(latest) {
			latest = period.end
		}
	}
	if len(ranges) == 0 {
		return result, nil
	}
	intervals, err := m.IntervalRepository.ListOverlapping(request.WorkspaceID, earliest.UTC(), latest.UTC())
	if err != nil {
		return nil, err
	}
	base, queryErr := m.dashboardQueryContexts(request.WorkspaceID, contexts, definition.Query)
	for _, widget := range definition.Widgets {
		if widget.Type == "text" {
			continue
		}
		p, err := readDashboardMetric(widget)
		if err != nil {
			return nil, err
		}
		period := ranges[widget.ID]
		metric := DashboardMetricResult{Unit: dashboardMetricUnit(p.Metric), TimeZone: request.TimeZone, Start: period.start.UTC(), End: period.end.UTC()}
		if widget.Type == "radar" {
			metric.Kind = "radar"
			metric.Radar = &DashboardRadarPayload{Unit: metric.Unit, Axes: []DashboardRadarAxis{}}
			hasActivity := false
			metric.HasActivity = &hasActivity
		}
		initializeDashboardAdvancedResult(&metric, widget.Type)
		if queryErr != nil {
			metric.Error = queryErr.Error()
			result.Widgets[widget.ID] = metric
			continue
		}
		selected, err := m.dashboardQueryContexts(request.WorkspaceID, base, widget.Query)
		if err != nil {
			metric.Error = err.Error()
			result.Widgets[widget.ID] = metric
			continue
		}
		if err := m.validateDashboardScope(request.WorkspaceID, DashboardScope{ProjectID: p.ProjectID, ContextID: p.ContextID}); err != nil {
			if _, invalid := err.(*InvalidDashboardError); invalid {
				metric.Error = "The selected project or context is no longer available in this workspace."
			} else {
				metric.Error = err.Error()
			}
			result.Widgets[widget.ID] = metric
			continue
		}
		baseIDs := dashboardContextIDs(base, request.Scope, dashboardMetricProperties{Archived: "include"})
		selectedIDs := dashboardContextIDs(selected, request.Scope, p)
		if widget.Type == "radar" {
			var projects []*Project
			if p.AxesBy == "project" {
				projects, err = m.ProjectRepository.ListIncludingArchived(request.WorkspaceID)
				if err != nil {
					return nil, err
				}
			}
			payload, ignored, radarErr := evaluateDashboardRadar(intervals, contexts, projects, selectedIDs, baseIDs, period, asOf, location, request.FirstDayOfWeek, p)
			if radarErr != nil {
				metric.Error = radarErr.Error()
			} else {
				metric.Radar, metric.IgnoredIntervals = payload, ignored
				hasActivity := aggregateDashboardIntervals(intervals, selectedIDs, period, asOf, p.IncludeRunning).duration > 0
				metric.HasActivity = &hasActivity
			}
			result.Widgets[widget.ID] = metric
			continue
		}
		if widget.Type == "bar" {
			payload, ignored, barErr := evaluateDashboardBar(intervals, contexts, selectedIDs, baseIDs, period, asOf, request.FirstDayOfWeek, p)
			if barErr != nil {
				metric.Error = barErr.Error()
			} else {
				metric.Categories, metric.IgnoredIntervals = payload, ignored
				active := aggregateDashboardIntervals(intervals, selectedIDs, period, asOf, p.IncludeRunning).duration > 0
				metric.HasActivity = &active
			}
			result.Widgets[widget.ID] = metric
			continue
		}
		if widget.Type == "heatmap" {
			payload, ignored, heatmapErr := evaluateDashboardHeatmap(intervals, selectedIDs, baseIDs, period, asOf, p)
			if heatmapErr != nil {
				metric.Error = heatmapErr.Error()
			} else {
				metric.Heatmap, metric.IgnoredIntervals = payload, ignored
				active := payload.ActiveDays > 0
				metric.HasActivity = &active
			}
			result.Widgets[widget.ID] = metric
			continue
		}
		if widget.Type == "line" {
			payload, ignored, lineErr := evaluateDashboardLine(intervals, contexts, selectedIDs, baseIDs, period, asOf, request.FirstDayOfWeek, p)
			if lineErr != nil {
				metric.Error = lineErr.Error()
			} else {
				metric.Series, metric.IgnoredIntervals = payload, ignored
				active := aggregateDashboardIntervals(intervals, selectedIDs, period, asOf, p.IncludeRunning).duration > 0
				metric.HasActivity = &active
			}
			result.Widgets[widget.ID] = metric
			continue
		}
		if widget.Type == "table" {
			payload, ignored, tableErr := evaluateDashboardTable(intervals, contexts, selectedIDs, baseIDs, period, asOf, p, request.ViewState[widget.ID])
			if tableErr != nil {
				metric.Error = tableErr.Error()
			} else {
				metric.Table, metric.IgnoredIntervals = payload, ignored
				active := aggregateDashboardIntervals(intervals, selectedIDs, period, asOf, p.IncludeRunning).duration > 0
				metric.HasActivity = &active
			}
			result.Widgets[widget.ID] = metric
			continue
		}
		aggregate := aggregateDashboardIntervals(intervals, selectedIDs, period, asOf, p.IncludeRunning)
		value := aggregate.value(p.Metric)
		if p.Metric == "time_share.ratio" {
			denominator := aggregateDashboardIntervals(intervals, baseIDs, period, asOf, p.IncludeRunning)
			if denominator.duration > 0 {
				share := aggregate.duration / denominator.duration * 100
				value = &share
			}
		}
		metric.Value = value
		metric.IgnoredIntervals = aggregate.ignored
		result.Widgets[widget.ID] = metric
	}
	return result, nil
}

func (m *ContextManager) validateDashboardScope(workspaceID string, scope DashboardScope) error {
	if scope.ProjectID != "" {
		project, err := m.ProjectRepository.GetById(scope.ProjectID)
		if err != nil {
			return err
		}
		if project == nil || project.WorkspaceId != workspaceID {
			return invalidDefinition("project is not available in this workspace")
		}
	}
	if scope.ContextID != "" {
		context, err := m.ContextRepository.GetById(scope.ContextID)
		if err != nil {
			return err
		}
		if context == nil || context.WorkspaceId != workspaceID {
			return invalidDefinition("context is not available in this workspace")
		}
	}
	if scope.Date != "" {
		if _, err := time.Parse(time.DateOnly, scope.Date); err != nil {
			return invalidDefinition("invalid daily summary date")
		}
	}
	return nil
}

func (m *ContextManager) dashboardQueryContexts(workspaceID string, contexts []*Context, query string) ([]*Context, error) {
	sqlQuery, err := m.interpretContextQuery(query)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(sqlQuery.WhereClause) == "" {
		return contexts, nil
	}
	sqlQuery.WorkspaceId = workspaceID
	matched, err := m.ContextRepository.Query(sqlQuery)
	if err != nil {
		return nil, fmt.Errorf("Unable to evaluate the query.")
	}
	ids := make(map[string]bool, len(matched))
	for _, context := range matched {
		if context != nil {
			ids[context.Id] = true
		}
	}
	selected := make([]*Context, 0, len(contexts))
	for _, context := range contexts {
		if context != nil && ids[context.Id] {
			selected = append(selected, context)
		}
	}
	return selected, nil
}

func dashboardContextIDs(contexts []*Context, scope DashboardScope, p dashboardMetricProperties) map[string]bool {
	ids := make(map[string]bool, len(contexts))
	for _, context := range contexts {
		if context == nil {
			continue
		}
		if scope.ContextID != "" && context.Id != scope.ContextID || p.ContextID != "" && context.Id != p.ContextID {
			continue
		}
		if scope.ProjectID != "" && (context.ProjectId == nil || *context.ProjectId != scope.ProjectID) || p.ProjectID != "" && (context.ProjectId == nil || *context.ProjectId != p.ProjectID) {
			continue
		}
		if p.Archived == "exclude" && context.Archived || p.Archived == "only" && !context.Archived {
			continue
		}
		if p.Tag != "" {
			found := false
			for _, tag := range context.Tags {
				if tag != nil && tag.Name == p.Tag {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		ids[context.Id] = true
	}
	return ids
}

func dashboardMetricPeriod(p dashboardMetricProperties, inherited string, request DashboardEvaluationRequest, now time.Time) (dashboardPeriodRange, error) {
	location := now.Location()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
	if request.Scope.Date != "" {
		day, err := time.ParseInLocation(time.DateOnly, request.Scope.Date, location)
		if err != nil {
			return dashboardPeriodRange{}, invalidDefinition("invalid daily summary date")
		}
		return dashboardPeriodRange{day, day.AddDate(0, 0, 1)}, nil
	}
	period := p.Period
	if period == "inherit" {
		period = inherited
		if period == "" {
			period = "this_week"
		}
	}
	start, end := today, today.AddDate(0, 0, 1)
	switch period {
	case "last_3_months", "last_6_months", "last_12_months":
		months := map[string]int{"last_3_months": 3, "last_6_months": 6, "last_12_months": 12}[period]
		start = time.Date(today.Year(), today.Month()-time.Month(months-1), 1, 0, 0, 0, 0, location)
		end = time.Date(today.Year(), today.Month()+1, 1, 0, 0, 0, 0, location)
	case "today":
	case "yesterday":
		start, end = today.AddDate(0, 0, -1), today
	case "this_week", "last_week":
		offset := (int(today.Weekday()) - request.FirstDayOfWeek + 7) % 7
		start = today.AddDate(0, 0, -offset)
		if period == "last_week" {
			start = start.AddDate(0, 0, -7)
		}
		end = start.AddDate(0, 0, 7)
	case "this_month", "last_month":
		start = time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, location)
		if period == "last_month" {
			start = start.AddDate(0, -1, 0)
		}
		end = start.AddDate(0, 1, 0)
	case "last_7_days":
		start = today.AddDate(0, 0, -6)
	case "last_30_days":
		start = today.AddDate(0, 0, -29)
	case "custom":
		var err error
		start, err = time.ParseInLocation(time.DateOnly, p.StartDate, location)
		if err != nil {
			return dashboardPeriodRange{}, invalidDefinition("invalid start date")
		}
		end, err = time.ParseInLocation(time.DateOnly, p.EndDate, location)
		if err != nil || end.Before(start) {
			return dashboardPeriodRange{}, invalidDefinition("invalid end date")
		}
		end = end.AddDate(0, 0, 1)
	default:
		return dashboardPeriodRange{}, invalidDefinition("invalid metric period")
	}
	return dashboardPeriodRange{start, end}, nil
}

type dashboardAggregate struct {
	duration, maximum        float64
	sessions, contexts, days int
	ignored                  int
}

func aggregateDashboardIntervals(intervals []*Interval, ids map[string]bool, period dashboardPeriodRange, asOf time.Time, running bool) dashboardAggregate {
	var result dashboardAggregate
	sessions, contexts, days := map[string]bool{}, map[string]bool{}, map[string]bool{}
	for _, interval := range intervals {
		if interval == nil || !ids[interval.ContextId] || sessions[interval.Id] {
			continue
		}
		if !timeIsSet(interval.Start) {
			result.ignored++
			continue
		}
		end := asOf
		if timeIsSet(interval.End) {
			end = *interval.End
		} else if interval.Status != IntervalStatusActive {
			result.ignored++
			continue
		} else if !running {
			continue
		}
		start := *interval.Start
		if end.After(asOf) {
			end = asOf
		}
		if !end.After(start) {
			result.ignored++
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
		duration := end.Sub(start).Seconds() * 1000
		result.duration += duration
		result.maximum = max(result.maximum, duration)
		sessions[interval.Id], contexts[interval.ContextId] = true, true
		local := start.In(period.start.Location())
		day := time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, local.Location())
		for day.Before(end) {
			days[day.Format(time.DateOnly)] = true
			day = day.AddDate(0, 0, 1)
		}
	}
	result.sessions, result.contexts, result.days = len(sessions), len(contexts), len(days)
	return result
}

func (a dashboardAggregate) value(metric string) *float64 {
	var value float64
	switch metric {
	case "duration.sum":
		value = a.duration
	case "duration.mean":
		if a.sessions == 0 {
			return nil
		}
		value = a.duration / float64(a.sessions)
	case "duration.max":
		if a.sessions == 0 {
			return nil
		}
		value = a.maximum
	case "sessions.count":
		value = float64(a.sessions)
	case "contexts.distinct_count":
		value = float64(a.contexts)
	case "active_days.distinct_count":
		value = float64(a.days)
	default:
		return nil
	}
	return &value
}

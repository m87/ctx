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
	WorkspaceID    string          `json:"workspaceId"`
	Definition     json.RawMessage `json:"definition"`
	TimeZone       string          `json:"timeZone"`
	FirstDayOfWeek int             `json:"firstDayOfWeek"`
	Scope          DashboardScope  `json:"scope"`
}

type DashboardMetricResult struct {
	Value            *float64  `json:"value"`
	Unit             string    `json:"unit"`
	TimeZone         string    `json:"timeZone"`
	Start            time.Time `json:"start"`
	End              time.Time `json:"end"`
	Error            string    `json:"error,omitempty"`
	IgnoredIntervals int       `json:"ignoredIntervals"`
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

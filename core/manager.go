package core

import (
	"fmt"
	"sort"
	"time"
)

type TimeRange struct {
	Start time.Time
	End   time.Time
}

type ContextManager struct {
	TimeProvider        TimeProvider
	ContextRepository   ContextRepository
	IntervalRepository  IntervalRepository
	WorkspaceRepository WorkspaceRepository
}

func NewContextManager(
	tp TimeProvider,
	contextRepo ContextRepository,
	intervalRepo IntervalRepository,
	workspaceRepo WorkspaceRepository,
) *ContextManager {
	return &ContextManager{
		TimeProvider:        tp,
		ContextRepository:   contextRepo,
		IntervalRepository:  intervalRepo,
		WorkspaceRepository: workspaceRepo,
	}
}

func (m *ContextManager) SaveInterval(interval *Interval) (string, error) {
	if interval == nil {
		return "", fmt.Errorf("interval is required")
	}

	context, err := m.ContextRepository.GetById(interval.ContextId)
	if err != nil {
		return "", err
	}
	if context == nil {
		return "", &ContextNotFoundError{ContextId: interval.ContextId}
	}

	interval.WorkspaceId = context.WorkspaceId
	return m.IntervalRepository.Save(interval)
}

type ContextNotFoundError struct {
	ContextId string
}

func (e *ContextNotFoundError) Error() string {
	return fmt.Sprintf("context %q not found", e.ContextId)
}

func (m *ContextManager) CreateContext(context *Context) (string, error) {
	if context == nil {
		return "", fmt.Errorf("context is required")
	}
	if context.WorkspaceId == "" {
		return "", &WorkspaceNotFoundError{}
	}

	workspace, err := m.WorkspaceRepository.GetById(context.WorkspaceId)
	if err != nil {
		return "", err
	}
	if workspace == nil {
		return "", &WorkspaceNotFoundError{WorkspaceId: context.WorkspaceId}
	}

	context.Id = ""
	return m.ContextRepository.Save(context)
}

type WorkspaceNotFoundError struct {
	WorkspaceId string
}

func (e *WorkspaceNotFoundError) Error() string {
	if e.WorkspaceId == "" {
		return "workspace is required"
	}
	return fmt.Sprintf("workspace %q not found", e.WorkspaceId)
}

func (m *ContextManager) SwitchContext(context *Context) error {
	activeContext, _ := m.ContextRepository.GetActive()
	endTime := m.TimeProvider.Now()
	startTime := m.TimeProvider.Now()

	if activeContext != nil && activeContext.Id == context.Id {
		return nil
	}

	if activeContext != nil {
		activeContext.Status = "inactive"
		m.ContextRepository.Save(activeContext)

		activeInterval, _ := m.IntervalRepository.GetActiveIntervalByContextId(activeContext.Id)

		if activeInterval != nil {
			activeInterval.Duration = endTime.Time.Sub(activeInterval.Start.Time)
			activeInterval.End = endTime
			activeInterval.Status = "completed"
			m.SaveInterval(activeInterval)
		}
	}

	if context.Id == "" {
		id, err := m.CreateContext(context)
		if err != nil {
			return err
		}
		context.Id = id
	}

	context, err := m.ContextRepository.GetById(context.Id)
	if err != nil {
		return err
	}

	context.Status = "active"
	m.ContextRepository.Save(context)

	newInterval := &Interval{
		ContextId:   context.Id,
		Start:       startTime,
		Status:      "active",
		WorkspaceId: context.WorkspaceId,
	}
	m.SaveInterval(newInterval)

	return nil
}

func (m *ContextManager) FreeActiveContext() error {
	activeContext, err := m.ContextRepository.GetActive()
	if err != nil {
		return err
	}
	if activeContext == nil {
		return nil
	}

	endTime := m.TimeProvider.Now()

	activeContext.Status = "inactive"
	if _, err := m.ContextRepository.Save(activeContext); err != nil {
		return err
	}

	activeInterval, err := m.IntervalRepository.GetActiveIntervalByContextId(activeContext.Id)
	if err != nil {
		return err
	}

	if activeInterval != nil {
		activeInterval.Duration = endTime.Time.Sub(activeInterval.Start.Time)
		activeInterval.End = endTime
		activeInterval.Status = "completed"
		if _, err := m.SaveInterval(activeInterval); err != nil {
			return err
		}
	}

	return nil
}

func (m *ContextManager) GetStats(contextId string, date time.Time) (*ContextStats, error) {
	allIntervalsByContext, err := m.IntervalRepository.ListByContextId(contextId)
	if err != nil {
		return nil, err
	}

	now := m.TimeProvider.Now().Time.UTC()

	var totalDuration time.Duration
	var totalSessions int
	var sessions int
	dayRanges := make([]TimeRange, 0, len(allIntervalsByContext))

	for _, interval := range allIntervalsByContext {
		intervalDuration := interval.Duration
		if intervalDuration <= 0 {
			start := interval.Start.Time.UTC()
			if start.IsZero() {
				intervalDuration = 0
			} else {
				end := interval.End.Time.UTC()
				if end.IsZero() {
					if interval.Status == "active" {
						end = now
					}
				}
				if end.After(start) {
					intervalDuration = end.Sub(start)
				}
			}
		}

		totalDuration += intervalDuration
		if intervalDuration > 0 {
			totalSessions++
		}

		if dayRange, ok := ClipIntervalRangeToDay(interval, date, now); ok {
			dayRanges = append(dayRanges, dayRange)
			sessions++
		}
	}

	duration := SumMergedRangesDuration(dayRanges)

	return &ContextStats{
		ContextId:     contextId,
		Date:          date,
		Duration:      duration,
		Sessions:      sessions,
		TotalDuration: totalDuration,
		TotalSessions: totalSessions,
	}, nil
}

func ClipIntervalRangeToDay(interval *Interval, date time.Time, now time.Time) (TimeRange, bool) {
	if interval == nil {
		return TimeRange{}, false
	}

	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	dayEnd := dayStart.Add(24 * time.Hour)

	start := interval.Start.Time.UTC()
	if start.IsZero() {
		return TimeRange{}, false
	}

	end := interval.End.Time.UTC()
	if end.IsZero() {
		if interval.Status != "active" {
			return TimeRange{}, false
		}
		end = now
	}

	if end.Before(dayStart) || !start.Before(dayEnd) {
		return TimeRange{}, false
	}

	if start.Before(dayStart) {
		start = dayStart
	}
	if end.After(dayEnd) {
		end = dayEnd
	}

	if !end.After(start) {
		return TimeRange{}, false
	}

	return TimeRange{Start: start, End: end}, true
}

func ClipIntervalDurationToDay(interval *Interval, date time.Time, now time.Time) time.Duration {
	rng, ok := ClipIntervalRangeToDay(interval, date, now)
	if !ok {
		return 0
	}
	return rng.End.Sub(rng.Start)
}

func SumMergedRangesDuration(ranges []TimeRange) time.Duration {
	if len(ranges) == 0 {
		return 0
	}

	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].Start.Before(ranges[j].Start)
	})

	mergedStart := ranges[0].Start
	mergedEnd := ranges[0].End
	var total time.Duration

	for _, rng := range ranges[1:] {
		if !rng.Start.After(mergedEnd) {
			if rng.End.After(mergedEnd) {
				mergedEnd = rng.End
			}
			continue
		}

		total += mergedEnd.Sub(mergedStart)
		mergedStart = rng.Start
		mergedEnd = rng.End
	}

	total += mergedEnd.Sub(mergedStart)
	return total
}

type ContextStats struct {
	ContextId     string        `json:"contextId"`
	Date          time.Time     `json:"date"`
	Duration      time.Duration `json:"duration"`
	Sessions      int           `json:"sessions"`
	TotalDuration time.Duration `json:"totalDuration"`
	TotalSessions int           `json:"totalSessions"`
}

func (m *ContextManager) GetWorkspaceStats(workspaceId string) (*WorkspaceStats, error) {
	contexts, err := m.ContextRepository.ListByWorkspace(workspaceId)
	if err != nil {
		return nil, err
	}

	now := m.TimeProvider.Now().Time.UTC()
	contextStats := make([]*WorkspaceContextStats, 0, len(contexts))
	var totalDuration time.Duration
	var totalSessions int

	for _, context := range contexts {
		if context == nil {
			continue
		}

		intervals, err := m.IntervalRepository.ListByContextId(context.Id)
		if err != nil {
			return nil, err
		}

		stats := &WorkspaceContextStats{ContextId: context.Id}
		for _, interval := range intervals {
			duration := intervalDurationAt(interval, now)
			if duration <= 0 {
				continue
			}
			stats.Duration += duration
			stats.IntervalCount++
		}

		totalDuration += stats.Duration
		totalSessions += stats.IntervalCount
		contextStats = append(contextStats, stats)
	}

	for _, stats := range contextStats {
		if totalDuration > 0 {
			stats.Percentage = float64(stats.Duration) / float64(totalDuration) * 100
		}
	}
	sort.Slice(contextStats, func(i, j int) bool {
		return contextStats[i].Duration > contextStats[j].Duration
	})

	return &WorkspaceStats{
		WorkspaceId:   workspaceId,
		Contexts:      contexts,
		ContextStats:  contextStats,
		TotalDuration: totalDuration,
		TotalSessions: totalSessions,
	}, nil
}

func intervalDurationAt(interval *Interval, now time.Time) time.Duration {
	if interval == nil {
		return 0
	}

	start := interval.Start.Time.UTC()
	end := interval.End.Time.UTC()
	if !start.IsZero() {
		if end.IsZero() && interval.Status == "active" {
			end = now
		}
		if end.After(start) {
			return end.Sub(start)
		}
	}

	if interval.Duration > 0 {
		return interval.Duration
	}
	return 0
}

func (m *ContextManager) DeleteWorkspace(workspaceId string) error {
	contexts, err := m.ContextRepository.ListByWorkspace(workspaceId)
	if err != nil {
		return err
	}

	if len(contexts) > 0 {
		return &WorkspaceInUseError{WorkspaceId: workspaceId}
	}

	return m.WorkspaceRepository.Delete(workspaceId)
}

type WorkspaceInUseError struct {
	WorkspaceId string
}

func (e *WorkspaceInUseError) Error() string {
	return "Cannot delete workspace because it is in use by one or more contexts"
}

func (m *ContextManager) EnsureDefaultWorkspace() error {
	workspaces, err := m.WorkspaceRepository.List()
	if err != nil {
		return err
	}
	var defaultWorkspaceId string
	if len(workspaces) == 0 {
		defaultWorkspace := &Workspace{
			Name: "Default",
		}
		id, err := m.WorkspaceRepository.Save(defaultWorkspace)
		if err != nil {
			return err
		}
		defaultWorkspaceId = id
	} else {
		for _, workspace := range workspaces {
			if workspace != nil && workspace.Name == "Default" {
				defaultWorkspaceId = workspace.Id
				break
			}
		}
	}
	if defaultWorkspaceId == "" {
		return nil
	}
	return m.setDefaultWorkspaceIfNotSet(defaultWorkspaceId)
}

func (m *ContextManager) setDefaultWorkspaceIfNotSet(defaultWorkspaceId string) error {
	contexts, err := m.ContextRepository.List()
	if err != nil {
		return err
	}
	for _, ctx := range contexts {
		if ctx == nil || ctx.WorkspaceId != "" {
			continue
		}
		ctx.WorkspaceId = defaultWorkspaceId
		if _, err := m.ContextRepository.Save(ctx); err != nil {
			return err
		}
	}

	intervals, err := m.IntervalRepository.List()
	if err != nil {
		return err
	}
	for _, interval := range intervals {
		if interval == nil || interval.WorkspaceId != "" {
			continue
		}
		interval.WorkspaceId = defaultWorkspaceId
		if _, err := m.IntervalRepository.Save(interval); err != nil {
			return err
		}
	}

	return nil
}

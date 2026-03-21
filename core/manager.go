package core

import "time"

type ContextManager struct {
	TimeProvider       TimeProvider
	ContextRepository  ContextRepository
	IntervalRepository IntervalRepository
}

func NewContextManager(tp TimeProvider, contextRepo ContextRepository, intervalRepo IntervalRepository) *ContextManager {
	return &ContextManager{
		TimeProvider:       tp,
		ContextRepository:  contextRepo,
		IntervalRepository: intervalRepo,
	}
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
			m.IntervalRepository.Save(activeInterval)
		}
	}

	if context.Id == "" {
		id, _ := m.ContextRepository.Save(context)
		context.Id = id
	}

	context, err := m.ContextRepository.GetById(context.Id)
	if err != nil {
		return err
	}

	context.Status = "active"
	m.ContextRepository.Save(context)

	newInterval := &Interval{
		ContextId: context.Id,
		Start:     startTime,
		Status:    "active",
	}
	m.IntervalRepository.Save(newInterval)

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
		if _, err := m.IntervalRepository.Save(activeInterval); err != nil {
			return err
		}
	}

	return nil
}

func (m *ContextManager) GetStats(contextId string, date time.Time) (*ContextStats, error) {
	intervalsByDay, err := m.IntervalRepository.ListByDay(date)
	if err != nil {
		return nil, err
	}
	allIntervalsByContext, err := m.IntervalRepository.ListByContextId(contextId)
	if err != nil {
		return nil, err
	}

	var totalDuration time.Duration
	var totalSessions int
	var duration time.Duration
	var sessions int

	for _, interval := range allIntervalsByContext {
		totalDuration += interval.Duration
		if interval.Duration > 0 {
			totalSessions++
		}
	}

	for _, interval := range intervalsByDay {
		if interval.ContextId == contextId {
			duration += interval.Duration
			if interval.Duration > 0 {
				sessions++
			}
		}
	}

	return &ContextStats{
		ContextId:     contextId,
		Date:          date,
		Duration:      duration,
		Sessions:      sessions,
		TotalDuration: totalDuration,
		TotalSessions: totalSessions,
	}, nil
}

type ContextStats struct {
	ContextId     string        `json:"contextId"`
	Date          time.Time     `json:"date"`
	Duration      time.Duration `json:"duration"`
	Sessions      int           `json:"sessions"`
	TotalDuration time.Duration `json:"totalDuration"`
	TotalSessions int           `json:"totalSessions"`
}

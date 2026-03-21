package core

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

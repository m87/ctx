package core

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
			activeInterval.Duration = durationBetween(activeInterval.Start, &endTime)
			activeInterval.End = &endTime
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
	if context == nil {
		return &ContextNotFoundError{}
	}
	if context.Archived {
		return &ContextArchivedError{ContextId: context.Id}
	}

	context.Status = "active"
	m.ContextRepository.Save(context)

	newInterval := &Interval{
		ContextId:   context.Id,
		Start:       &startTime,
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
		activeInterval.Duration = durationBetween(activeInterval.Start, &endTime)
		activeInterval.End = &endTime
		activeInterval.Status = "completed"
		if _, err := m.SaveInterval(activeInterval); err != nil {
			return err
		}
	}

	return nil
}

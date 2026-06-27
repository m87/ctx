package core

import "time"

func (m *ContextManager) RepairIntegrity() (*IntegrityRepairResult, error) {
	repairedCount := 0
	err := m.RunInTransaction(func(txManager *ContextManager) error {
		count, err := txManager.repairIntegrity()
		repairedCount = count
		return err
	})
	if err != nil {
		return nil, err
	}

	report, err := m.CheckIntegrity()
	if err != nil {
		return nil, err
	}
	return &IntegrityRepairResult{RepairedCount: repairedCount, Report: report}, nil
}

func (m *ContextManager) repairIntegrity() (int, error) {
	workspaces, err := m.WorkspaceRepository.List()
	if err != nil {
		return 0, err
	}

	workspaceIds := map[string]struct{}{}
	defaultWorkspaceId := ""
	for _, workspace := range workspaces {
		if workspace == nil || workspace.Id == "" {
			continue
		}
		workspaceIds[workspace.Id] = struct{}{}
		if workspace.Name == "Default" {
			defaultWorkspaceId = workspace.Id
		}
	}

	contexts, err := m.ContextRepository.List()
	if err != nil {
		return 0, err
	}

	needsDefaultWorkspace := len(workspaceIds) == 0
	if defaultWorkspaceId == "" && !needsDefaultWorkspace {
		for _, context := range contexts {
			if context == nil {
				continue
			}
			if _, ok := workspaceIds[context.WorkspaceId]; !ok {
				needsDefaultWorkspace = true
				break
			}
		}
	}

	repaired := 0
	if defaultWorkspaceId == "" && needsDefaultWorkspace {
		defaultWorkspaceId, err = m.WorkspaceRepository.Save(&Workspace{Name: "Default"})
		if err != nil {
			return 0, err
		}
		workspaceIds[defaultWorkspaceId] = struct{}{}
		repaired++
	}

	contextsById := make(map[string]*Context, len(contexts))
	activeContexts := make([]*Context, 0)
	for _, context := range contexts {
		if context == nil {
			continue
		}
		if _, ok := workspaceIds[context.WorkspaceId]; !ok {
			context.WorkspaceId = defaultWorkspaceId
			if _, err := m.ContextRepository.Save(context); err != nil {
				return repaired, err
			}
			repaired++
		}
		if context.Status == "active" {
			activeContexts = append(activeContexts, context)
		}
		contextsById[context.Id] = context
	}

	intervals, err := m.IntervalRepository.List()
	if err != nil {
		return repaired, err
	}
	repairTime := m.integrityRepairTime()
	for _, interval := range intervals {
		if interval == nil {
			continue
		}
		context := contextsById[interval.ContextId]
		if context != nil && interval.WorkspaceId != context.WorkspaceId {
			interval.WorkspaceId = context.WorkspaceId
			if _, err := m.IntervalRepository.Save(interval); err != nil {
				return repaired, err
			}
			repaired++
		}

		if interval.Status == "active" && zonedTimeIsSet(interval.End) {
			completeIntervalAtEnd(interval)
			if _, err := m.IntervalRepository.Save(interval); err != nil {
				return repaired, err
			}
			repaired++
		}
	}

	repairedContexts, err := m.repairActiveContexts(activeContexts, intervals, repairTime)
	if err != nil {
		return repaired, err
	}
	repaired += repairedContexts

	return repaired, nil
}

func (m *ContextManager) integrityRepairTime() ZonedTime {
	if m.TimeProvider != nil {
		return m.TimeProvider.Now()
	}
	now := time.Now().UTC()
	return ZonedTime{Time: now, Timezone: "UTC"}
}

func completeIntervalAtEnd(interval *Interval) {
	if zonedTimeIsSet(interval.Start) && zonedTimeIsSet(interval.End) && interval.End.Time.After(interval.Start.Time) {
		interval.Duration = interval.End.Time.Sub(interval.Start.Time)
	}
	interval.Status = "completed"
}

func (m *ContextManager) repairActiveContexts(activeContexts []*Context, intervals []*Interval, endTime ZonedTime) (int, error) {
	if len(activeContexts) == 0 {
		return 0, nil
	}

	openIntervalsByContext := map[string][]*Interval{}
	for _, interval := range intervals {
		if intervalIsOpenActive(interval) {
			openIntervalsByContext[interval.ContextId] = append(openIntervalsByContext[interval.ContextId], interval)
		}
	}

	contextToKeep := ""
	if len(activeContexts) > 1 {
		contextToKeep = newestOpenActiveContextId(activeContexts, openIntervalsByContext)
	}

	repaired := 0
	for _, context := range activeContexts {
		if context == nil {
			continue
		}

		openIntervals := openIntervalsByContext[context.Id]
		shouldStop := len(openIntervals) == 0 || (contextToKeep != "" && context.Id != contextToKeep)
		if !shouldStop {
			continue
		}

		context.Status = "inactive"
		if _, err := m.ContextRepository.Save(context); err != nil {
			return repaired, err
		}
		repaired++

		for _, interval := range openIntervals {
			interval.End = endTime
			completeIntervalAtEnd(interval)
			if _, err := m.IntervalRepository.Save(interval); err != nil {
				return repaired, err
			}
			repaired++
		}
	}

	return repaired, nil
}

func newestOpenActiveContextId(activeContexts []*Context, openIntervalsByContext map[string][]*Interval) string {
	newestContextId := ""
	var newestStart time.Time

	for _, context := range activeContexts {
		if context == nil {
			continue
		}
		for _, interval := range openIntervalsByContext[context.Id] {
			if !zonedTimeIsSet(interval.Start) {
				continue
			}
			if newestContextId == "" || interval.Start.Time.After(newestStart) {
				newestContextId = context.Id
				newestStart = interval.Start.Time
			}
		}
	}

	return newestContextId
}

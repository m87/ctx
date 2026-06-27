package core

func (m *ContextManager) CheckIntegrity() (*IntegrityReport, error) {
	workspaces, err := m.WorkspaceRepository.List()
	if err != nil {
		return nil, err
	}
	contexts, err := m.ContextRepository.List()
	if err != nil {
		return nil, err
	}
	intervals, err := m.IntervalRepository.List()
	if err != nil {
		return nil, err
	}

	workspaceIds := map[string]struct{}{}
	for _, workspace := range workspaces {
		workspaceIds[workspace.Id] = struct{}{}
	}

	contextsById := map[string]*Context{}
	activeContexts := []*Context{}
	issues := []*IntegrityIssue{}
	for _, context := range contexts {
		contextsById[context.Id] = context
		if context.Status == "active" {
			activeContexts = append(activeContexts, context)
		}
		if context.WorkspaceId == "" {
			issues = append(issues, contextIntegrityIssue(context, "CONTEXT_MISSING_WORKSPACE", "Context has no workspace assigned"))
			continue
		}
		if _, ok := workspaceIds[context.WorkspaceId]; !ok {
			issues = append(issues, contextIntegrityIssue(context, "CONTEXT_WORKSPACE_NOT_FOUND", "Context references a workspace that does not exist"))
		}
	}

	intervalsByContextId := map[string][]*Interval{}
	for _, interval := range intervals {
		intervalsByContextId[interval.ContextId] = append(intervalsByContextId[interval.ContextId], interval)
		context := contextsById[interval.ContextId]
		contextIsValid := interval.ContextId != "" && context != nil
		if interval.ContextId == "" {
			issues = append(issues, intervalIntegrityIssue(interval, "INTERVAL_MISSING_CONTEXT", "Interval has no context assigned", false))
		} else if context == nil {
			issues = append(issues, intervalIntegrityIssue(interval, "INTERVAL_CONTEXT_NOT_FOUND", "Interval references a context that does not exist", false))
		}

		if interval.WorkspaceId == "" {
			issues = append(issues, intervalIntegrityIssue(interval, "INTERVAL_MISSING_WORKSPACE", "Interval has no workspace assigned", contextIsValid))
		} else if _, ok := workspaceIds[interval.WorkspaceId]; !ok {
			issues = append(issues, intervalIntegrityIssue(interval, "INTERVAL_WORKSPACE_NOT_FOUND", "Interval references a workspace that does not exist", contextIsValid))
		} else if contextIsValid && interval.WorkspaceId != context.WorkspaceId {
			issues = append(issues, intervalIntegrityIssue(interval, "INTERVAL_WORKSPACE_MISMATCH", "Interval workspace differs from its context workspace", true))
		}

		if intervalIsInactive(interval) && (!zonedTimeIsSet(interval.Start) || !zonedTimeIsSet(interval.End)) {
			issues = append(issues, intervalIntegrityIssue(interval, "INACTIVE_INTERVAL_MISSING_TIME", "Inactive interval must have both start and end time set; set the missing time manually or delete the interval", false))
		}
		if interval.Status == "active" && zonedTimeIsSet(interval.End) {
			issues = append(issues, intervalIntegrityIssue(interval, "ACTIVE_INTERVAL_HAS_END", "Active interval has an end time set and can be completed automatically", true))
		}
	}

	if len(activeContexts) > 1 {
		for _, context := range activeContexts {
			issues = append(issues, contextIntegrityIssue(context, "MULTIPLE_ACTIVE_CONTEXTS", "More than one context is active; automatic repair keeps the newest open context and stops the others"))
		}
	}

	for _, context := range activeContexts {
		hasOpenActiveInterval := false
		for _, interval := range intervalsByContextId[context.Id] {
			if intervalIsOpenActive(interval) {
				hasOpenActiveInterval = true
				break
			}
		}
		if !hasOpenActiveInterval {
			issues = append(issues, contextIntegrityIssue(context, "ACTIVE_CONTEXT_WITHOUT_OPEN_INTERVAL", "Active context has no open active interval and can be stopped automatically"))
		}
	}

	return &IntegrityReport{
		Healthy:        len(issues) == 0,
		WorkspaceCount: len(workspaces),
		ContextCount:   len(contexts),
		IntervalCount:  len(intervals),
		Issues:         issues,
	}, nil
}

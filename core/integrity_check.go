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
	issues := []*IntegrityIssue{}
	for _, context := range contexts {
		contextsById[context.Id] = context
		if context.WorkspaceId == "" {
			issues = append(issues, integrityIssue("context", context.Id, "CONTEXT_MISSING_WORKSPACE", "Context has no workspace assigned"))
			continue
		}
		if _, ok := workspaceIds[context.WorkspaceId]; !ok {
			issues = append(issues, integrityIssue("context", context.Id, "CONTEXT_WORKSPACE_NOT_FOUND", "Context references a workspace that does not exist"))
		}
	}

	for _, interval := range intervals {
		context := contextsById[interval.ContextId]
		if interval.ContextId == "" {
			issues = append(issues, integrityIssue("interval", interval.Id, "INTERVAL_MISSING_CONTEXT", "Interval has no context assigned"))
		} else if context == nil {
			issues = append(issues, integrityIssue("interval", interval.Id, "INTERVAL_CONTEXT_NOT_FOUND", "Interval references a context that does not exist"))
		}

		if interval.WorkspaceId == "" {
			issues = append(issues, integrityIssue("interval", interval.Id, "INTERVAL_MISSING_WORKSPACE", "Interval has no workspace assigned"))
		} else if _, ok := workspaceIds[interval.WorkspaceId]; !ok {
			issues = append(issues, integrityIssue("interval", interval.Id, "INTERVAL_WORKSPACE_NOT_FOUND", "Interval references a workspace that does not exist"))
		}

		if context != nil && interval.WorkspaceId != "" && interval.WorkspaceId != context.WorkspaceId {
			issues = append(issues, integrityIssue("interval", interval.Id, "INTERVAL_WORKSPACE_MISMATCH", "Interval workspace differs from its context workspace"))
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

package core

type IntegrityIssue struct {
	EntityType  string `json:"entityType"`
	EntityId    string `json:"entityId"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

type IntegrityReport struct {
	Healthy        bool              `json:"healthy"`
	WorkspaceCount int               `json:"workspaceCount"`
	ContextCount   int               `json:"contextCount"`
	IntervalCount  int               `json:"intervalCount"`
	Issues         []*IntegrityIssue `json:"issues"`
}

type IntegrityRepairResult struct {
	RepairedCount int              `json:"repairedCount"`
	Report        *IntegrityReport `json:"report"`
}

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

	workspaceIds := make(map[string]struct{}, len(workspaces)+1)
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

	repaired := 0
	if defaultWorkspaceId == "" {
		defaultWorkspaceId, err = m.WorkspaceRepository.Save(&Workspace{Name: "Default"})
		if err != nil {
			return 0, err
		}
		workspaceIds[defaultWorkspaceId] = struct{}{}
		repaired++
	}

	contexts, err := m.ContextRepository.List()
	if err != nil {
		return repaired, err
	}
	contextsById := make(map[string]*Context, len(contexts))
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
		contextsById[context.Id] = context
	}

	intervals, err := m.IntervalRepository.List()
	if err != nil {
		return repaired, err
	}
	for _, interval := range intervals {
		if interval == nil {
			continue
		}
		context := contextsById[interval.ContextId]
		if context == nil || interval.WorkspaceId == context.WorkspaceId {
			continue
		}
		interval.WorkspaceId = context.WorkspaceId
		if _, err := m.IntervalRepository.Save(interval); err != nil {
			return repaired, err
		}
		repaired++
	}

	return repaired, nil
}

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

	workspaceIds := make(map[string]struct{}, len(workspaces))
	for _, workspace := range workspaces {
		if workspace != nil && workspace.Id != "" {
			workspaceIds[workspace.Id] = struct{}{}
		}
	}

	contextsById := make(map[string]*Context, len(contexts))
	issues := make([]*IntegrityIssue, 0)
	for _, context := range contexts {
		if context == nil {
			continue
		}
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
		if interval == nil {
			continue
		}

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

func integrityIssue(entityType, entityId, code, description string) *IntegrityIssue {
	return &IntegrityIssue{
		EntityType:  entityType,
		EntityId:    entityId,
		Code:        code,
		Description: description,
	}
}

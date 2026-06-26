package core

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

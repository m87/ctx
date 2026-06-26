package core

func (m *ContextManager) ListIntegrityContextOptions() ([]*IntegrityContextOption, error) {
	workspaces, err := m.WorkspaceRepository.List()
	if err != nil {
		return nil, err
	}

	workspaceNames := make(map[string]string, len(workspaces))
	for _, workspace := range workspaces {
		if workspace == nil || workspace.Id == "" {
			continue
		}
		workspaceNames[workspace.Id] = workspace.Name
	}

	contexts, err := m.ContextRepository.List()
	if err != nil {
		return nil, err
	}

	options := make([]*IntegrityContextOption, 0, len(contexts))
	for _, context := range contexts {
		if context == nil || context.Id == "" {
			continue
		}
		workspaceName, ok := workspaceNames[context.WorkspaceId]
		if !ok {
			continue
		}
		options = append(options, &IntegrityContextOption{
			Id:            context.Id,
			Name:          context.Name,
			WorkspaceId:   context.WorkspaceId,
			WorkspaceName: workspaceName,
		})
	}

	return options, nil
}

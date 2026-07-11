package core

type ArchiveContextActiveError struct {
	ContextId string
}

func (e *ArchiveContextActiveError) Error() string {
	return "cannot archive active context: " + e.ContextId
}

func (m *ContextManager) ArchiveContext(contextId string) error {
	context, err := m.ContextRepository.GetById(contextId)
	if err != nil {
		return err
	}

	if context.Status == "active" {
		return &ArchiveContextActiveError{ContextId: contextId}
	}

	context.Archived = true
	context.Status = "archived"

	_, err = m.ContextRepository.Save(context)
	if err != nil {
		return err
	}

	return nil
}

func (m *ContextManager) RestoreContext(contextId string) error {
	context, err := m.ContextRepository.GetById(contextId)
	if err != nil {
		return err
	}

	context.Archived = false
	context.Status = "inactive"

	_, err = m.ContextRepository.Save(context)
	if err != nil {
		return err
	}

	return nil
}

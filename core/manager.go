package core

type ContextManager struct {
	TimeProvider        TimeProvider
	ContextRepository   *ContextRepository
	WorkspaceRepository *WorkspaceRepository
	ProjectRepository   *ProjectRepository
	IntervalRepository  *IntervalRepository
}

func NewContextManager(tp TimeProvider, contextRepo *ContextRepository, workspaceRepo *WorkspaceRepository, projectRepo *ProjectRepository, intervalRepo *IntervalRepository) *ContextManager {
	return &ContextManager{
		TimeProvider:        tp,
		ContextRepository:   contextRepo,
		WorkspaceRepository: workspaceRepo,
		ProjectRepository:   projectRepo,
		IntervalRepository:  intervalRepo,
	}
}

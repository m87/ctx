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

package core

import (
	"github.com/m87/ctx/utils"
	"github.com/m87/nod"
)

type ContextManager struct {
	TimeProvider TimeProvider
	repository   *nod.Repository
}

func NewContextManager(tp TimeProvider, repo *nod.Repository) *ContextManager {
	return &ContextManager{
		TimeProvider: tp,
		repository:   repo,
	}
}

func (cm *ContextManager) ExecuteUnchecked(fn func(repository *nod.Repository) error) error {
	return cm.repository.Transaction(fn)
}

func (cm *ContextManager) ExecuteChecked(fn func(repository *nod.Repository) error) {
	utils.Check(cm.ExecuteUnchecked(fn))
}

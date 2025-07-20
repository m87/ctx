package core

import (
	"errors"

	ctxtime "github.com/m87/ctx/time"
)

type ContextManager struct {
	TimeProvider ctxtime.TimeProvider
	StateStore   TransactionalStore[State]
}

type Session struct {
	State        *State
	TimeProvider ctxtime.TimeProvider
}

func NewContextManager(timeProvider ctxtime.TimeProvider, stateStore TransactionalStore[State]) *ContextManager {
	return &ContextManager{
		TimeProvider: timeProvider,
		StateStore:   stateStore,
	}
}

func (manager *ContextManager) WithSession(fn func(session Session) error) error {
	stateTx, state, stateErr := manager.StateStore.BeginAndGet()
	if stateErr != nil {
		return errors.Join(stateErr)
	}

	if err := fn(Session{
		State:        state,
		TimeProvider: manager.TimeProvider,
	}); err != nil {
		return err
	}

	stateErr = stateTx.Commit()

	if stateErr != nil {
		stateRollbackEer := stateTx.Rollback()

		if stateRollbackEer != nil {
			panic(errors.Join(stateRollbackEer))
		}
	}

	return nil
}

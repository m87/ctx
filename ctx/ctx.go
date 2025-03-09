package ctx

import (
	"time"

	"github.com/m87/ctx/ctx_model"
)

type TimeProvider interface {
	Now() time.Time
}

type ContextStore interface {
	Apply(fn ctx_model.StatePatch)
}

type EventsRegistryStore interface {
}

type ArchiveStore interface {
}

type ContextManager struct {
	contextStore ContextStore
	timeProvider TimeProvider
}

func New(contextStore ContextStore, timeProvider TimeProvider) *ContextManager {
	return &ContextManager{
		contextStore: contextStore,
		timeProvider: timeProvider,
	}
}

func (manager *ContextManager) CreateContext(id string, description string) {
	manager.contextStore.Apply(
		func(state *ctx_model.State) {
			state.Contexts[id] = ctx_model.Context{
				Id:          id,
				Description: description,
				State:       ctx_model.ACTIVE,
				Intervals:   []ctx_model.Interval{},
			}
		},
	)
}

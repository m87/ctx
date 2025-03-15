package ctx

import (
	"time"

	"github.com/m87/ctx/ctx_model"
	"github.com/m87/ctx/ctx_store"
	"github.com/spf13/viper"
)

type RealTimeProvider struct{}

func (provider *RealTimeProvider) Now() time.Time {
	return time.Now().Local()
}

func NewTimer() *RealTimeProvider {
	return &RealTimeProvider{}
}

func CreateManager() *ContextManager {
	return New(ctx_store.New(viper.GetString("path")), NewTimer())
}

type ContextManager struct {
	ContextStore ctx_model.ContextStore
	TimeProvider ctx_model.TimeProvider
}

func New(contextStore ctx_model.ContextStore, timeProvider ctx_model.TimeProvider) *ContextManager {
	return &ContextManager{
		ContextStore: contextStore,
		TimeProvider: timeProvider,
	}
}

func (manager *ContextManager) CreateContext(id string, description string) {
	manager.ContextStore.Apply(
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

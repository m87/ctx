package ctx

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/m87/ctx/ctx_model"
	"github.com/m87/ctx/ctx_store"
	"github.com/spf13/viper"
)

type ContextState int

const (
	ACTIVE ContextState = iota
	FINISHED
)

type Interval struct {
	Start    time.Time     `json:"start"`
	End      time.Time     `json:"end"`
	Duration time.Duration `json:"duration"`
}

type Context struct {
	Id          string        `json:"id"`
	Description string        `json:"description"`
	Comments    []string      `json:"comments"`
	State       ContextState  `json:"state"`
	Duration    time.Duration `json:"duration"`
	Intervals   []Interval    `json:"intervals"`
}

type State struct {
	Contexts  map[string]Context `json:"contexts"`
	CurrentId string             `json:"currentId"`
}

type StatePatch func(*State) error

type TimeProvider interface {
	Now() time.Time
}

type ContextStore interface {
	Apply(fn StatePatch) error
}

type EventsRegistryStore interface {
}

type ArchiveStore interface {
}

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

func (manager *ContextManager) CreateContext(id string, description string) error {
	return manager.ContextStore.Apply(
		func(state *ctx_model.State) error {
			if _, ok := state.Contexts[id]; ok {
				return errors.New("Context already exists")
			} else {
				state.Contexts[id] = ctx_model.Context{
					Id:          id,
					Description: description,
					State:       ctx_model.ACTIVE,
					Intervals:   []ctx_model.Interval{},
				}
			}
			return nil
		},
	)
}

func (manager *ContextManager) List() {
	manager.ContextStore.Read(
		func(state *ctx_model.State) error {
			for _, v := range state.Contexts {
				fmt.Printf("- %s\n", v.Description)
			}
			return nil
		},
	)
}

func (manager *ContextManager) ListFull() {
	manager.ContextStore.Read(
		func(state *ctx_model.State) error {
			for _, v := range state.Contexts {
				fmt.Printf("- [%s] %s\n", v.Id, v.Description)
				for _, interval := range v.Intervals {
					fmt.Printf("\t- %s - %s\n", interval.Start.Local().Format(time.DateTime), interval.End.Local().Format(time.DateTime))
				}
			}
			return nil
		},
	)
}

func (manager *ContextManager) ListJson() {
	manager.ContextStore.Read(
		func(state *ctx_model.State) error {
			v := make([]ctx_model.Context, 0, len(state.Contexts))
			for _, c := range state.Contexts {
				v = append(v, c)
			}
			s, _ := json.Marshal(v)

			fmt.Printf("%s", string(s))
			return nil
		},
	)
}

func (manager *ContextManager) Switch(id string) error {
	return manager.ContextStore.Apply(
		func(state *ctx_model.State) error {
			if _, ok := state.Contexts[id]; ok {

			}
			return nil

		})
}

func (manager *ContextManager) CreateIfNotExistsAndSwitch(id string) error {
	return manager.ContextStore.Apply(
		func(s *ctx_model.State) error {
			return nil
		})
}

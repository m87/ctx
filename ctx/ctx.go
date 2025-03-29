package ctx

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/m87/ctx/ctx_model"
	"github.com/m87/ctx/localstore"
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

type RealTimeProvider struct{}

func (provider *RealTimeProvider) Now() time.Time {
	return time.Now().Local()
}

func NewTimer() *RealTimeProvider {
	return &RealTimeProvider{}
}

func CreateManager() *ContextManager {
	return New(localstore.NewContextStore(viper.GetString("storePath")), localstore.NewEventsStore(viper.GetString("storePath")), localstore.NewArchiveStore(viper.GetString("storePath")), NewTimer())
}

type ContextManager struct {
	ContextStore ctx_model.ContextStore
	EventsStore  ctx_model.EventsStore
	ArchiveStore ctx_model.ArchiveStore
	TimeProvider ctx_model.TimeProvider
}

func New(contextStore ctx_model.ContextStore, eventsStore ctx_model.EventsStore, archiveStore ctx_model.ArchiveStore, timeProvider ctx_model.TimeProvider) *ContextManager {
	return &ContextManager{
		ContextStore: contextStore,
		EventsStore:  eventsStore,
		ArchiveStore: archiveStore,
		TimeProvider: timeProvider,
	}
}

func (manager *ContextManager) createContetxtInternal(state *ctx_model.State, id string, description string) error {
	if len(strings.TrimSpace(id)) == 0 {
		return errors.New("empty id")
	}
	if len(strings.TrimSpace(description)) == 0 {
		return errors.New("empty description")
	}

	if _, ok := state.Contexts[id]; ok {
		return errors.New("context already exists")
	} else {
		state.Contexts[id] = ctx_model.Context{
			Id:          id,
			Description: description,
			State:       ctx_model.ACTIVE,
			Intervals:   []ctx_model.Interval{},
		}
		manager.PublishContextEvent(state.Contexts[id], manager.TimeProvider.Now(), ctx_model.CREATE_CTX, nil)
	}
	return nil
}

func (manager *ContextManager) CreateContext(id string, description string) error {
	return manager.ContextStore.Apply(
		func(state *ctx_model.State) error {
			return manager.createContetxtInternal(state, id, description)
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

func (manager *ContextManager) endInterval(state *ctx_model.State, id string, now time.Time) {
	prev := state.Contexts[state.CurrentId]
	interval := prev.Intervals[len(prev.Intervals)-1]
	interval.End = now
	interval.Duration = interval.End.Sub(interval.Start)
	state.Contexts[state.CurrentId].Intervals[len(prev.Intervals)-1] = interval
	prev.Duration = prev.Duration + interval.Duration
	state.Contexts[state.CurrentId] = prev
	manager.PublishContextEvent(state.Contexts[id], now, ctx_model.END_INTERVAL, map[string]string{
		"duration": interval.Duration.String(),
	})
}

func (manager *ContextManager) switchInternal(state *ctx_model.State, id string) error {
	if len(strings.TrimSpace(id)) == 0 {
		return errors.New("empty id")
	}

	if state.CurrentId == id {
		return errors.New("context already active")
	}

	if _, ok := state.Contexts[id]; !ok {
		return errors.New("context does not exist")
	}

	now := manager.TimeProvider.Now()
	prevId := state.CurrentId
	if state.CurrentId != "" {
		manager.endInterval(state, id, now)
	}

	if ctx, ok := state.Contexts[id]; ok {
		state.CurrentId = ctx.Id
		manager.PublishContextEvent(state.Contexts[id], now, ctx_model.SWITCH_CTX, map[string]string{
			"from": prevId,
		})
		ctx.Intervals = append(state.Contexts[id].Intervals, ctx_model.Interval{Start: now})
		manager.PublishContextEvent(state.Contexts[id], now, ctx_model.START_INTERVAL, nil)
		state.Contexts[id] = ctx
	}
	return nil
}

func (manager *ContextManager) Switch(id string) error {
	return manager.ContextStore.Apply(
		func(state *ctx_model.State) error {
			if _, ok := state.Contexts[id]; ok {
				return manager.switchInternal(state, id)
			} else {
				return errors.New("context does not exists")
			}
		})
}

func (manager *ContextManager) CreateIfNotExistsAndSwitch(id string, description string) error {
	return manager.ContextStore.Apply(
		func(state *ctx_model.State) error {
			if _, ok := state.Contexts[id]; !ok {
				err := manager.createContetxtInternal(state, id, description)
				if err != nil {
					return err
				}
			}
			return manager.switchInternal(state, id)
		})
}

func (manager *ContextManager) Ctx(id string) (ctx_model.Context, error) {
  ctx := ctx_model.Context{}

  manager.ContextStore.Read(func(s *ctx_model.State) error {
    ctx = s.Contexts[id]
    return nil
  })

  return ctx ,nil
} 

func (manager *ContextManager) PublishEvent(event ctx_model.Event) error {
	return manager.EventsStore.Apply(func(er *ctx_model.EventRegistry) error {
		event.UUID = uuid.NewString()
		er.Events = append(er.Events, event)
		return nil
	})
}

func (manager *ContextManager) PublishContextEvent(context ctx_model.Context, dateTime time.Time, eventType ctx_model.EventType, data map[string]string) error {
	return manager.PublishEvent(ctx_model.Event{
		DateTime:    dateTime,
		Type:        eventType,
		CtxId:       context.Id,
		Description: context.Description,
		Data:        data,
	})
}

func (manager *ContextManager) FilterEvents(filter ctx_model.EventsFilter) []ctx_model.Event {
  evs := []ctx_model.Event{}
  
  manager.EventsStore.Read(func(er *ctx_model.EventRegistry) error {
    
    evs = manager.filterEvents(er, filter)

    return nil
  },
  )

  return evs
}

func (manager *ContextManager) filterEvents(er *ctx_model.EventRegistry, filter ctx_model.EventsFilter) []ctx_model.Event {
	evs := er.Events
	tmpEvs := []ctx_model.Event{}

	if filter.CtxId != "" {
		for _, v := range evs {
			if v.CtxId == filter.CtxId {
				tmpEvs = append(tmpEvs, v)
			}
		}
		evs = tmpEvs
	}

	if filter.Date != "" {
		for _, v := range evs {
			if v.DateTime.Local().Format(time.DateOnly) == filter.Date {
				tmpEvs = append(tmpEvs, v)
			}
		}
		evs = tmpEvs
	}

	tmpEvs = []ctx_model.Event{}

	if len(filter.Types) > 0 {

		for _, v := range evs {
			for _, t := range filter.Types {
				if v.Type == ctx_model.StringAsEvent(t) {
					tmpEvs = append(tmpEvs, v)
				}
			}

		}
		evs = tmpEvs
	}

	return evs
}

func (manager *ContextManager) ListEvents(filter ctx_model.EventsFilter) {
	manager.EventsStore.Read(func(er *ctx_model.EventRegistry) error {
		evs := manager.filterEvents(er, filter)

		for _, v := range evs {
			fmt.Printf("[%s] [%s] %s\n", v.DateTime.Local().Format(time.DateTime), ctx_model.EventAsString(v.Type), v.Description)
		}
		return nil
	})
}
func (manager *ContextManager) ListEventsJson(filter ctx_model.EventsFilter) {
	manager.EventsStore.Read(func(er *ctx_model.EventRegistry) error {
		evs := manager.filterEvents(er, filter)

		s, _ := json.Marshal(evs)

		fmt.Printf("%s", string(s))

		return nil
	})
}

func (manager *ContextManager) formatEventData(event ctx_model.Event) string {
	if len(event.Data) == 0 {
		return ""
	}
	data, err := json.Marshal(event.Data)
	if err != nil {
		return ""
	}
	return string(data)
}

func (manager *ContextManager) ListEventsFull(filter ctx_model.EventsFilter) {
	manager.EventsStore.Read(func(er *ctx_model.EventRegistry) error {
		evs := manager.filterEvents(er, filter)

		for _, v := range evs {
			fmt.Printf("[%s] [%s] %s [%s]\n", v.DateTime.Local().Format(time.DateTime), ctx_model.EventAsString(v.Type), v.Description, manager.formatEventData(v))
		}
		return nil
	})
}

func (manager *ContextManager) Free() error {
	return manager.ContextStore.Apply(
		func(state *ctx_model.State) error {
			if state.CurrentId == "" {
				return errors.New("no active context")
			}

			now := manager.TimeProvider.Now()
			manager.endInterval(state, state.CurrentId, now)
			state.CurrentId = ""
			return nil
		})
}

func (manager *ContextManager) deleteEvents(id string) error {
	return manager.EventsStore.Apply(
		func(er *ctx_model.EventRegistry) error {
			tmpEvs := []ctx_model.Event{}
			for _, v := range er.Events {
				if v.CtxId != id {
					tmpEvs = append(tmpEvs, v)
				}
			}
			er.Events = tmpEvs
			return nil
		})
}

func (manager *ContextManager) Delete(id string) error {
	return manager.ContextStore.Apply(
		func(state *ctx_model.State) error {
			if state.CurrentId == id {
				return errors.New("context is active")
			}

			if _, ok := state.Contexts[id]; ok {
				delete(state.Contexts, id)
				manager.deleteEvents(id)
				return nil
			} else {
				return errors.New("context does not exists")
			}
		})
}

func (manager *ContextManager) groupEventsByDate(events []ctx_model.Event) map[string][]ctx_model.Event {
	eventsByDate := make(map[string][]ctx_model.Event)
	for _, event := range events {
		date := event.DateTime.Format(time.DateOnly)
		eventsByDate[date] = append(eventsByDate[date], event)
	}

	return eventsByDate
}

func (manager *ContextManager) upsertArchive(entry *ctx_model.ContextArchive) error {
	return manager.ArchiveStore.Apply(entry.Context.Id, func(entry2Update *ctx_model.ContextArchive) error {
		if entry2Update.Context.Id != entry.Context.Id {
			return errors.New("contexts mismatch, entry to update: " + entry2Update.Context.Id + ", entry to archive: " + entry.Context.Id)
		}

		entry2Update.Context.Id = entry.Context.Id
		entry2Update.Context.Description = entry.Context.Description
		entry2Update.Context.Duration = entry2Update.Context.Duration + entry.Context.Duration
		entry2Update.Context.Comments = append(entry2Update.Context.Comments, entry.Context.Comments...)
		entry2Update.Context.Intervals = append(entry2Update.Context.Intervals, entry.Context.Intervals...)
		entry2Update.Context.State = entry.Context.State
		entry2Update.Events = append(entry2Update.Events, entry.Events...)
		return nil
	})
}

func (manager *ContextManager) upsertEventsArchive(eventsByDate map[string][]ctx_model.Event) error {
	for k, v := range eventsByDate {
		return manager.ArchiveStore.ApplyEvents(k, func(entry *ctx_model.EventsArchive) error {
			entry.Events = append(entry.Events, v...)
			return nil
		})
	}
	return nil
}

func (manager *ContextManager) Archive(id string) error {
	if err := manager.ContextStore.Read(
		func(state *ctx_model.State) error {
			if state.CurrentId == id {
				return errors.New("context is active")
			}

			if _, ok := state.Contexts[id]; ok {
				return manager.EventsStore.Read(func(er *ctx_model.EventRegistry) error {
					archiveEntry := ctx_model.ContextArchive{
						Context: state.Contexts[id],
						Events:  manager.filterEvents(er, ctx_model.EventsFilter{CtxId: id}),
					}

					if err := manager.upsertArchive(&archiveEntry); err != nil {
						return err
					}

					return manager.upsertEventsArchive(manager.groupEventsByDate(manager.filterEvents(er, ctx_model.EventsFilter{CtxId: id})))
				})
			} else {
				return errors.New("context does not exists")
			}
		}); err != nil {
		return err
	}

	return manager.Delete(id)
}

func (manager *ContextManager) ArchiveAll() error {
	return manager.ContextStore.Read(
		func(state *ctx_model.State) error {
			for _, ctx := range state.Contexts {
				if err := manager.Archive(ctx.Id); err != nil {
					return err
				}
			}
			return nil
		})
}


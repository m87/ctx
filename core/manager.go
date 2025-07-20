package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	ctxtime "github.com/m87/ctx/time"
)

type ContextManager struct {
	ContextStore ContextStore
	EventsStore  EventsStore
	ArchiveStore ArchiveStore
	TimeProvider ctxtime.TimeProvider
	StateStore   TransactionalStore[State]
	EventStore   TransactionalStore[EventRegistry]
}

type Session struct {
	State          *State
	EventsRegistry *EventRegistry
	TimeProvider   ctxtime.TimeProvider
}

func NewContextManager(contextStore ContextStore, eventsStore EventsStore, archiveStore ArchiveStore, timeProvider ctxtime.TimeProvider, stateStore TransactionalStore[State], eventStore TransactionalStore[EventRegistry]) *ContextManager {
	return &ContextManager{
		ContextStore: contextStore,
		EventsStore:  eventsStore,
		ArchiveStore: archiveStore,
		TimeProvider: timeProvider,
		StateStore:   stateStore,
		EventStore:   eventStore,
	}
}

func (manager *ContextManager) WithSession(fn func(session Session) error) error {
	stateTx, state, stateErr := manager.StateStore.BeginAndGet()
	erTx, er, erErr := manager.EventStore.BeginAndGet()
	if stateErr != nil || erErr != nil {
		return errors.Join(stateErr, erErr)
	}

	if err := fn(Session{
		State:          state,
		EventsRegistry: er,
		TimeProvider:   manager.TimeProvider,
	}); err != nil {
		return err
	}

	erErr = erTx.Commit()
	stateErr = stateTx.Commit()

	if erErr != nil || stateErr != nil {
		erRollbackErr := erTx.Rollback()
		stateRollbackEer := stateTx.Rollback()

		if erRollbackErr != nil || stateRollbackEer != nil {
			panic(errors.Join(erRollbackErr, stateRollbackEer))
		}
	}

	return nil
}

func (manager *ContextManager) PublishEvent(event Event) error {
	return manager.EventsStore.Apply(func(er *EventRegistry) error {
		event.UUID = uuid.NewString()
		er.Events = append(er.Events, event)
		return nil
	})
}

func (manager *ContextManager) PublishContextEvent(context Context, dateTime ctxtime.ZonedTime, eventType EventType, data map[string]string) error {
	return manager.PublishEvent(Event{
		DateTime:    dateTime,
		Type:        eventType,
		CtxId:       context.Id,
		Description: context.Description,
		Data:        data,
	})
}

func (manager *ContextManager) FilterEvents(filter EventsFilter) []Event {
	evs := []Event{}

	manager.EventsStore.Read(func(er *EventRegistry) error {

		evs = manager.filterEvents(er, filter)

		return nil
	},
	)

	return evs
}

func (manager *ContextManager) filterEvents(er *EventRegistry, filter EventsFilter) []Event {
	evs := er.Events
	tmpEvs := []Event{}

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
			if v.DateTime.Time.Format(time.DateOnly) == filter.Date {
				tmpEvs = append(tmpEvs, v)
			}
		}
		evs = tmpEvs
	}

	tmpEvs = []Event{}

	if len(filter.Types) > 0 {

		for _, v := range evs {
			for _, t := range filter.Types {
				if v.Type == StringAsEvent(t) {
					tmpEvs = append(tmpEvs, v)
				}
			}

		}
		evs = tmpEvs
	}

	return evs
}

func (manager *ContextManager) ListEvents(filter EventsFilter) {
	manager.EventsStore.Read(func(er *EventRegistry) error {
		evs := manager.filterEvents(er, filter)

		for _, v := range evs {
			fmt.Printf("[%s] [%s] %s\n", v.DateTime.Time.Format(time.DateTime), EventAsString(v.Type), v.Description)
		}
		return nil
	})
}
func (manager *ContextManager) ListEventsJson(filter EventsFilter) {
	manager.EventsStore.Read(func(er *EventRegistry) error {
		evs := manager.filterEvents(er, filter)

		s, _ := json.Marshal(evs)

		fmt.Printf("%s", string(s))

		return nil
	})
}

func (manager *ContextManager) formatEventData(event Event) string {
	if len(event.Data) == 0 {
		return ""
	}
	data, err := json.Marshal(event.Data)
	if err != nil {
		return ""
	}
	return string(data)
}

func (manager *ContextManager) ListEventsFull(filter EventsFilter) {
	manager.EventsStore.Read(func(er *EventRegistry) error {
		evs := manager.filterEvents(er, filter)

		for _, v := range evs {
			fmt.Printf("[%s] [%s] %s [%s]\n", v.DateTime.Time.Format(time.DateTime), EventAsString(v.Type), v.Description, manager.formatEventData(v))
		}
		return nil
	})
}

func (manager *ContextManager) groupEventsByDate(events []Event) map[string][]Event {
	eventsByDate := make(map[string][]Event)
	for _, event := range events {
		date := event.DateTime.Time.Format(time.DateOnly)
		eventsByDate[date] = append(eventsByDate[date], event)
	}

	return eventsByDate
}

func (manager *ContextManager) upsertArchive(entry *ContextArchive) error {
	return manager.ArchiveStore.Apply(entry.Context.Id, func(entry2Update *ContextArchive) error {
		if entry2Update.Context.Id != entry.Context.Id {
			return errors.New("contexts mismatch, entry to update: " + entry2Update.Context.Id + ", entry to archive: " + entry.Context.Id)
		}

		entry2Update.Context.Id = entry.Context.Id
		entry2Update.Context.Description = entry.Context.Description
		entry2Update.Context.Duration = entry2Update.Context.Duration + entry.Context.Duration
		entry2Update.Context.Comments = append(entry2Update.Context.Comments, entry.Context.Comments...)

		for _, interval := range entry.Context.Intervals {
			if _, ok := entry2Update.Context.Intervals[interval.Id]; !ok {
				entry2Update.Context.Intervals[interval.Id] = interval
			}
		}

		entry2Update.Context.State = entry.Context.State
		return nil
	})
}

func (manager *ContextManager) upsertEventsArchive(eventsByDate map[string][]Event) error {
	for k, v := range eventsByDate {
		return manager.ArchiveStore.ApplyEvents(k, func(entry *EventsArchive) error {
			entry.Events = append(entry.Events, v...)
			return nil
		})
	}
	return nil
}

func (manager *ContextManager) Archive(id string) error {
	if err := manager.ContextStore.Read(
		func(state *State) error {
			if state.CurrentId == id {
				return errors.New("context is active")
			}

			if _, ok := state.Contexts[id]; ok {
				archiveEntry := ContextArchive{
					Context: state.Contexts[id],
				}

				if err := manager.upsertArchive(&archiveEntry); err != nil {
					return err
				}

			} else {
				return errors.New("context does not exists")
			}
			return nil
		}); err != nil {
		return err
	}
	return manager.WithSession(func(session Session) error {
		return session.deleteInternal(id)
	})
}

func (manager *ContextManager) ArchiveAllEvents() error {
	return manager.EventsStore.Apply(func(er *EventRegistry) error {
		eventsByDate := manager.groupEventsByDate(er.Events)
		err := manager.upsertEventsArchive(eventsByDate)
		if err != nil {
			return err
		}
		er.Events = []Event{}
		return nil
	})
}

func (manager *ContextManager) ArchiveAll() error {
	err := manager.ContextStore.Read(
		func(state *State) error {
			for _, ctx := range state.Contexts {
				if err := manager.Archive(ctx.Id); err != nil {
					return err
				}
			}
			return nil
		})

	if err != nil {
		return err
	}

	return manager.ArchiveAllEvents()
}

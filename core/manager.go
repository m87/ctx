package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
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

func (manager *ContextManager) GetIntervalDurationsByDate(s *State, id string, date ctxtime.ZonedTime) (time.Duration, error) {
	var duration time.Duration = 0
	loc, err := time.LoadLocation(ctxtime.DetectTimezoneName())
	if err != nil {
		loc = time.UTC
	}
	startOfDay := time.Date(date.Time.Year(), date.Time.Month(), date.Time.Day(), 0, 0, 0, 0, loc)
	if ctx, ok := s.Contexts[id]; ok {
		for _, interval := range ctx.Intervals {
			if interval.Start.Time.Day() == startOfDay.Day() && interval.Start.Time.Month() == startOfDay.Month() && interval.Start.Time.Year() == startOfDay.Year() && interval.End.Time.Day() == startOfDay.Day() && interval.End.Time.Month() == startOfDay.Month() && interval.End.Time.Year() == startOfDay.Year() {
				duration += interval.Duration
			} else if interval.Start.Time.Before(startOfDay) && interval.End.Time.Day() == startOfDay.Day() && interval.End.Time.Month() == startOfDay.Month() && interval.End.Time.Year() == startOfDay.Year() {
				duration += interval.End.Time.Sub(startOfDay)
			} else if interval.Start.Time.Day() == startOfDay.Day() && interval.Start.Time.Month() == startOfDay.Month() && interval.Start.Time.Year() == startOfDay.Year() && interval.End.Time.After(startOfDay) {
				duration += 24*time.Hour - interval.Start.Time.Sub(startOfDay)
			} else if interval.Start.Time.Before(startOfDay) && interval.End.Time.After(startOfDay) {
				duration += 24 * time.Hour
			}
		}
	} else {
		return 0, errors.New("context does not exist")
	}
	return duration, nil
}

func (manager *ContextManager) GetIntervalsByDate(s *State, id string, date ctxtime.ZonedTime) []Interval {
	intervals := []Interval{}
	loc, err := time.LoadLocation(ctxtime.DetectTimezoneName())
	if err != nil {
		loc = time.UTC
	}
	startOfDay := time.Date(date.Time.Year(), date.Time.Month(), date.Time.Day(), 0, 0, 0, 0, loc)
	if ctx, ok := s.Contexts[id]; ok {
		for _, interval := range ctx.Intervals {
			if interval.Start.Time.Day() == startOfDay.Day() && interval.Start.Time.Month() == startOfDay.Month() && interval.Start.Time.Year() == startOfDay.Year() && interval.End.Time.Day() == startOfDay.Day() && interval.End.Time.Month() == startOfDay.Month() && interval.End.Time.Year() == startOfDay.Year() {
				intervals = append(intervals, Interval(interval))
			} else if interval.Start.Time.Before(startOfDay) && interval.End.Time.Day() == startOfDay.Day() && interval.End.Time.Month() == startOfDay.Month() && interval.End.Time.Year() == startOfDay.Year() {
				intervals = append(intervals, Interval(interval))
			} else if interval.Start.Time.Day() == startOfDay.Day() && interval.Start.Time.Month() == startOfDay.Month() && interval.Start.Time.Year() == startOfDay.Year() && interval.End.Time.After(startOfDay) {
				intervals = append(intervals, Interval(interval))
			} else if interval.Start.Time.Before(startOfDay) && interval.End.Time.After(startOfDay) {
				intervals = append(intervals, Interval(interval))
			}
		}
	}
	return intervals
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

func (manager *ContextManager) SplitContextIntervalById(ctxId string, id string, split time.Time) error {
	manager.ContextStore.Apply(func(s *State) error {
		context, ok := s.Contexts[id]
		if !ok {
			return errors.New("context does not exists")
		}

		interval := context.Intervals[id]
		interval.End.Time = split
		interval.Duration = split.Sub(interval.Start.Time)
		newId := uuid.NewString()
		context.Intervals[newId] = Interval{
			Id: newId,
			Start: ctxtime.ZonedTime{
				Time:     split,
				Timezone: interval.Start.Timezone,
			},
			End: ctxtime.ZonedTime{
				Time:     interval.End.Time,
				Timezone: interval.End.Timezone,
			},
			Duration: interval.End.Time.Sub(split),
		}

		s.Contexts[id] = context

		return nil

	})

	return nil
}

func (manager *ContextManager) EditContextInterval(id string, intervalId string, start ctxtime.ZonedTime, end ctxtime.ZonedTime) error {
	manager.ContextStore.Read(func(s *State) error {
		context, ok := s.Contexts[id]
		if !ok {
			return errors.New("context does not exist")
		}
		for _, interval := range context.Intervals {
			if interval.Id == intervalId {
				manager.EditContextIntervalById(id, intervalId, start, end)
				return nil
			}
		}
		return nil
	})
	return nil
}

func (manager *ContextManager) MoveIntervalById(idSrc string, idTarget string, intervalId string) error {
	return manager.ContextStore.Apply(func(state *State) error {
		if state.CurrentId == idTarget {
			return errors.New("context is active")
		}

		ctxSrc := state.Contexts[idSrc]
		ctxTarget := state.Contexts[idTarget]

		ctxTarget.Intervals[intervalId] = ctxSrc.Intervals[intervalId]
		delete(ctxSrc.Intervals, intervalId)

		ctxTarget.Duration += ctxTarget.Intervals[intervalId].Duration
		ctxSrc.Duration -= ctxTarget.Intervals[intervalId].Duration

		state.Contexts[idSrc] = ctxSrc
		state.Contexts[idTarget] = ctxTarget

		return nil
	})
}

func (manager *ContextManager) EditContextIntervalById(id string, intervalId string, start ctxtime.ZonedTime, end ctxtime.ZonedTime) error {
	return manager.ContextStore.Apply(func(s *State) error {
		if s.CurrentId == id {
			return errors.New("context is active")
		}

		oldDuration := s.Contexts[id].Intervals[intervalId].Duration
		oldStart := s.Contexts[id].Intervals[intervalId].Start.Time.Format(time.RFC3339)
		oldEnd := s.Contexts[id].Intervals[intervalId].End.Time.Format(time.RFC3339)

		ctx := s.Contexts[id]

		interval := ctx.Intervals[intervalId]

		interval.Start = start
		interval.End = end
		interval.Duration = end.Time.Sub(start.Time)
		ctx.Intervals[intervalId] = interval

		durationDiff := interval.Duration - oldDuration

		ctx.Duration = ctx.Duration + durationDiff

		s.Contexts[id] = ctx
		manager.PublishContextEvent(ctx, manager.TimeProvider.Now(), EDIT_CTX_INTERVAL, map[string]string{
			"old.start": oldStart,
			"old.end":   oldEnd,
			"new.start": ctx.Intervals[intervalId].Start.Time.Format(time.RFC3339),
			"new.end":   ctx.Intervals[intervalId].End.Time.Format(time.RFC3339),
		})
		return nil
	})

}

func (manager *ContextManager) Search(regex string) ([]Context, error) {
	ctxs := []Context{}
	re := regexp.MustCompile(regex)
	err := manager.ContextStore.Read(func(s *State) error {
		for _, ctx := range s.Contexts {
			if re.MatchString(ctx.Description) {
				ctxs = append(ctxs, ctx)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return ctxs, nil
}

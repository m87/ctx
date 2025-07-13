package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	ctxtime "github.com/m87/ctx/time"
	"github.com/m87/ctx/util"
)

type ContextManager struct {
	ContextStore ContextStore
	EventsStore  EventsStore
	ArchiveStore ArchiveStore
	TimeProvider ctxtime.TimeProvider
	StateStore   TransactionalStore[State]
}

func NewContextManager(contextStore ContextStore, eventsStore EventsStore, archiveStore ArchiveStore, timeProvider ctxtime.TimeProvider, stateStore TransactionalStore[State]) *ContextManager {
	return &ContextManager{
		ContextStore: contextStore,
		EventsStore:  eventsStore,
		ArchiveStore: archiveStore,
		TimeProvider: timeProvider,
		StateStore:   stateStore,
	}
}

func (manager *ContextManager) DeleteInterval(ctxId string, id string) error {
	tx, state, err := manager.StateStore.BeginAndGet()
	if err != nil {
		return err
	}

	if err := state.DeleteInterval(ctxId, id); err != nil {
		return err
	}

	return tx.Commit()

	// return manager.StateStore.WithTx(func(state *State) error {
	// 	return state.DeleteInterval(ctxId, id)
	// 	// manager.PublishContextEvent(ctx, manager.TimeProvider.Now(), DELETE_CTX_INTERVAL, map[string]string{
	// 	// 	"start":    interval.Start.Time.Format(time.RFC3339),
	// 	// 	"end":      interval.End.Time.Format(time.RFC3339),
	// 	// 	"duration": interval.Duration.String(),
	// 	// })
	// })
}

func (manager *ContextManager) DeleteIntervalByIndex(ctxId string, index int) error {
	return manager.StateStore.WithTx(func(state *State) error {
		return state.DeleteIntervalByIndex(ctxId, index)
	})
}

func (manager *ContextManager) createContetxtInternal(state *State, id string, description string) error {
	if len(strings.TrimSpace(id)) == 0 {
		return errors.New("empty id")
	}
	if len(strings.TrimSpace(description)) == 0 {
		return errors.New("empty description")
	}

	if _, ok := state.Contexts[id]; ok {
		return errors.New("context already exists")
	} else {
		state.Contexts[id] = Context{
			Id:          id,
			Description: description,
			State:       ACTIVE,
			Intervals:   []Interval{},
		}
		manager.PublishContextEvent(state.Contexts[id], manager.TimeProvider.Now(), CREATE_CTX, nil)
	}
	return nil
}

func (manager *ContextManager) CreateContext(id string, description string) error {
	return manager.ContextStore.Apply(
		func(state *State) error {
			return manager.createContetxtInternal(state, id, description)
		},
	)
}

func (manager *ContextManager) List() {
	manager.ContextStore.Read(
		func(state *State) error {
			ids := manager.getSortedContextIds(state)
			for _, id := range ids {
				v := state.Contexts[id]
				fmt.Printf("- %s\n", v.Description)
			}
			return nil
		},
	)
}

func (manager *ContextManager) ListFull() {
	manager.ContextStore.Read(
		func(state *State) error {
			ids := manager.getSortedContextIds(state)
			for _, id := range ids {
				v := state.Contexts[id]
				fmt.Printf("- [%s] %s\n", id, v.Description)
				for i, interval := range v.Intervals {
					fmt.Printf("\t[%d] %s - %s\n", i, interval.Start.Time.Format(time.DateTime), interval.End.Time.Format(time.DateTime))
				}
			}
			return nil
		},
	)
}

func (manager *ContextManager) ListJson() {
	manager.ContextStore.Read(
		func(state *State) error {
			v := make([]Context, 0, len(state.Contexts))
			for _, c := range state.Contexts {
				v = append(v, c)
			}
			s, _ := json.Marshal(v)

			fmt.Printf("%s", string(s))
			return nil
		},
	)
}

func (manager *ContextManager) ListJson2() []Context {
	output := []Context{}
	manager.ContextStore.Read(
		func(state *State) error {
			for _, c := range state.Contexts {
				output = append(output, c)
			}
			return nil
		},
	)
	return output
}

func (manager *ContextManager) getSortedContextIds(state *State) []string {
	ids := make([]string, 0, len(state.Contexts))
	for k := range state.Contexts {
		ids = append(ids, k)
	}
	sort.Strings(ids)
	return ids
}

func (manager *ContextManager) endInterval(state *State, id string, now ctxtime.ZonedTime) {
	prev := state.Contexts[state.CurrentId]
	interval := prev.Intervals[len(prev.Intervals)-1]
	interval.End = now
	interval.Duration = interval.End.Time.Sub(interval.Start.Time)
	state.Contexts[state.CurrentId].Intervals[len(prev.Intervals)-1] = interval
	prev.Duration = prev.Duration + interval.Duration
	state.Contexts[state.CurrentId] = prev
	manager.PublishContextEvent(state.Contexts[id], now, END_INTERVAL, map[string]string{
		"duration": interval.Duration.String(),
	})
}

func (manager *ContextManager) switchInternal(state *State, id string) error {
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
		manager.endInterval(state, state.CurrentId, now)
	}

	if ctx, ok := state.Contexts[id]; ok {
		state.CurrentId = ctx.Id
		manager.PublishContextEvent(state.Contexts[id], now, SWITCH_CTX, map[string]string{
			"from": prevId,
		})
		ctx.Intervals = append(state.Contexts[id].Intervals, Interval{Id: uuid.NewString(), Start: now})
		manager.PublishContextEvent(state.Contexts[id], now, START_INTERVAL, nil)
		state.Contexts[id] = ctx
	}
	return nil
}

func (manager *ContextManager) Switch(id string) error {
	return manager.ContextStore.Apply(
		func(state *State) error {
			if _, ok := state.Contexts[id]; ok {
				return manager.switchInternal(state, id)
			} else {
				return errors.New("context does not exists")
			}
		})
}

func (manager *ContextManager) CreateIfNotExistsAndSwitch(id string, description string) error {
	return manager.ContextStore.Apply(
		func(state *State) error {
			if _, ok := state.Contexts[id]; !ok {
				err := manager.createContetxtInternal(state, id, description)
				if err != nil {
					return err
				}
			}
			return manager.switchInternal(state, id)
		})
}

func (manager *ContextManager) Ctx(id string) (Context, error) {
	ctx := Context{}

	manager.ContextStore.Read(func(s *State) error {
		ctx = s.Contexts[id]
		return nil
	})

	return ctx, nil
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

func (manager *ContextManager) Free() error {
	return manager.ContextStore.Apply(
		func(state *State) error {
			if state.CurrentId == "" {
				return errors.New("no active context")
			}

			now := manager.TimeProvider.Now()
			manager.endInterval(state, state.CurrentId, now)
			state.CurrentId = ""
			return nil
		})
}

func (manager *ContextManager) deleteInternal(state *State, id string) error {
	if state.CurrentId == id {
		return errors.New("context is active")
	}

	if _, ok := state.Contexts[id]; ok {
		context := state.Contexts[id]
		delete(state.Contexts, id)
		manager.PublishContextEvent(context, manager.TimeProvider.Now(), DELETE_CTX, nil)
		return nil
	} else {
		return errors.New("context does not exists")
	}
}

func (manager *ContextManager) Delete(id string) error {
	return manager.ContextStore.Apply(
		func(state *State) error {
			return manager.deleteInternal(state, id)
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
		entry2Update.Context.Intervals = append(entry2Update.Context.Intervals, entry.Context.Intervals...)
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

	return manager.Delete(id)
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

func (manager *ContextManager) MergeContext(from string, to string) error {
	return manager.ContextStore.Apply(func(state *State) error {
		if from == to {
			return errors.New("contexts are the same")
		}

		if from == state.CurrentId {
			return errors.New("from context is active")
		}

		if _, ok := state.Contexts[from]; !ok {
			return errors.New("context does not exists: " + from)
		}
		if _, ok := state.Contexts[to]; !ok {
			return errors.New("context does not exists: " + to)
		}

		fromCtx := state.Contexts[from]
		toCtx := state.Contexts[to]

		toCtx.Comments = append(toCtx.Comments, fromCtx.Comments...)
		toCtx.Duration = toCtx.Duration + fromCtx.Duration
		toCtx.Intervals = append(toCtx.Intervals, fromCtx.Intervals...)

		state.Contexts[to] = toCtx
		manager.deleteInternal(state, from)

		manager.PublishContextEvent(state.Contexts[to], manager.TimeProvider.Now(), MERGE_CTX, map[string]string{
			"from": from,
			"to":   to,
		})

		return nil
	},
	)
}

func (manager *ContextManager) SplitContextIntervalById(ctxId string, id string, split time.Time) error {
	ctx, _ := manager.Ctx(ctxId)
	index := -1
	for i, interval := range ctx.Intervals {
		if interval.Id == id {
			index = i
		}
	}

	return manager.SplitContextIntervalByIndex(ctxId, index, split)
}

func (manager *ContextManager) SplitContextIntervalByIndex(id string, intervalIndex int, split time.Time) error {
	manager.ContextStore.Apply(func(s *State) error {
		context, ok := s.Contexts[id]
		if !ok {
			return errors.New("context does not exists")
		}

		interval := context.Intervals[intervalIndex]
		context.Intervals[intervalIndex].End.Time = split
		context.Intervals[intervalIndex].Duration = split.Sub(interval.Start.Time)
		context.Intervals = append(context.Intervals, Interval{
			Id: uuid.NewString(),
			Start: ctxtime.ZonedTime{
				Time:     split,
				Timezone: interval.Start.Timezone,
			},
			End: ctxtime.ZonedTime{
				Time:     interval.End.Time,
				Timezone: interval.End.Timezone,
			},
			Duration: interval.End.Time.Sub(split),
		})

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
		for i, interval := range context.Intervals {
			if interval.Id == intervalId {
				manager.EditContextIntervalByIndex(id, i, start, end)
				return nil
			}
		}
		return nil
	})
	return nil
}

func (manager *ContextManager) MoveIntervalById(idSrc string, idTarget string, id string) error {
	ctx, _ := manager.Ctx(idSrc)
	for i, interval := range ctx.Intervals {
		if interval.Id == id {
			return manager.MoveIntervalByIndex(idSrc, idTarget, i)
		}
	}
	return nil
}

func (manager *ContextManager) MoveIntervalByIndex(idSrc string, idTarget string, intervalIndex int) error {
	return manager.ContextStore.Apply(func(state *State) error {
		if state.CurrentId == idTarget {
			return errors.New("context is active")
		}

		ctxSrc := state.Contexts[idSrc]
		ctxTarget := state.Contexts[idTarget]

		ctxTarget.Intervals = append(ctxTarget.Intervals, ctxSrc.Intervals[intervalIndex])
		ctxSrc.Intervals = append(ctxSrc.Intervals[:intervalIndex], ctxSrc.Intervals[intervalIndex+1:]...)

		state.Contexts[idSrc] = ctxSrc
		state.Contexts[idTarget] = ctxTarget

		return nil
	})
}

func (manager *ContextManager) EditContextIntervalByIndex(id string, intervalIndex int, start ctxtime.ZonedTime, end ctxtime.ZonedTime) error {
	return manager.ContextStore.Apply(func(s *State) error {
		if s.CurrentId == id {
			return errors.New("context is active")
		}

		oldDuration := s.Contexts[id].Intervals[intervalIndex].Duration
		oldStart := s.Contexts[id].Intervals[intervalIndex].Start.Time.Format(time.RFC3339)
		oldEnd := s.Contexts[id].Intervals[intervalIndex].End.Time.Format(time.RFC3339)

		ctx := s.Contexts[id]

		ctx.Intervals[intervalIndex].Start = start
		ctx.Intervals[intervalIndex].End = end
		ctx.Intervals[intervalIndex].Duration = ctx.Intervals[intervalIndex].End.Time.Sub(ctx.Intervals[intervalIndex].Start.Time)

		durationDiff := ctx.Intervals[intervalIndex].Duration - oldDuration

		ctx.Duration = ctx.Duration + durationDiff

		s.Contexts[id] = ctx
		manager.PublishContextEvent(ctx, manager.TimeProvider.Now(), EDIT_CTX_INTERVAL, map[string]string{
			"old.start": oldStart,
			"old.end":   oldEnd,
			"new.start": ctx.Intervals[intervalIndex].Start.Time.Format(time.RFC3339),
			"new.end":   ctx.Intervals[intervalIndex].End.Time.Format(time.RFC3339),
		})
		return nil
	})

}

func (manager *ContextManager) RenameContext(srcId string, targetId string, name string) error {
	return manager.ContextStore.Apply(func(s *State) error {
		s.Contexts[targetId] = Context{
			Id:          targetId,
			Description: name,
			Intervals:   append([]Interval{}, s.Contexts[srcId].Intervals...),
			State:       s.Contexts[srcId].State,
			Duration:    s.Contexts[srcId].Duration,
			Comments:    append([]string{}, s.Contexts[srcId].Comments...),
		}

		ctx := s.Contexts[srcId]
		delete(s.Contexts, srcId)
		manager.PublishContextEvent(ctx, manager.TimeProvider.Now(), RENAME_CTX, map[string]string{
			"src.id":             ctx.Id,
			"src.description":    ctx.Description,
			"target.id":          targetId,
			"target:description": name,
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

func (manager *ContextManager) LabelContext(id string, label string) error {
	return manager.ContextStore.Apply(func(s *State) error {
		if _, ok := s.Contexts[id]; !ok {
			return errors.New("context does not exist")
		}
		ctx := s.Contexts[id]
		if !util.Contains(s.Contexts[id].Labels, label) {
			ctx.Labels = append(s.Contexts[id].Labels, label)
			manager.PublishContextEvent(s.Contexts[id], manager.TimeProvider.Now(), LABEL_CTX, map[string]string{
				"label": label,
			})
			s.Contexts[id] = ctx
		}
		return nil
	})
}

func (manager *ContextManager) DeleteLabelContext(id string, label string) error {
	return manager.ContextStore.Apply(func(s *State) error {
		if _, ok := s.Contexts[id]; !ok {
			return errors.New("context does not exist")
		}

		ctx := s.Contexts[id]
		if util.Contains(s.Contexts[id].Labels, label) {
			ctx.Labels = util.Remove(s.Contexts[id].Labels, label)
			manager.PublishContextEvent(s.Contexts[id], manager.TimeProvider.Now(), DELETE_CTX_LABEL, map[string]string{
				"label": label,
			})
			s.Contexts[id] = ctx
		}
		return nil
	})
}

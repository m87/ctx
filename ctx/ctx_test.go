package ctx

import (
	"encoding/json"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/m87/ctx/ctx_model"
	"github.com/m87/ctx/test"
	"github.com/stretchr/testify/assert"
)

type TestTimeProvider struct {
	Current ctx_model.ZonedTime
}

func NewTestTimerProvider(dateTime string) *TestTimeProvider {
	loc, err := time.LoadLocation(ctx_model.DetectTimezoneName())
	if err != nil {
		loc = time.UTC
	}
	dt, _ := time.ParseInLocation(time.DateTime, dateTime, loc)
	return &TestTimeProvider{
		Current: ctx_model.ZonedTime{Time: dt, Timezone: ctx_model.DetectTimezoneName()},
	}
}

func (provider *TestTimeProvider) Now() ctx_model.ZonedTime {
	return provider.Current
}

type TestContextStore struct {
	store []byte
}

type TestEventsStore struct {
	store []byte
}

type TestArchiveStore struct {
	store       []byte
	storeEvents []byte
}

func (cs *TestContextStore) Load() ctx_model.State {
	state := ctx_model.State{}
	err := json.Unmarshal(cs.store, &state)
	if err != nil {
		log.Fatal("Unable to parse state store")
	}

	return state
}

func (cs *TestContextStore) Save(state *ctx_model.State) {
	data, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	cs.store = data
}

func NewTestContextStore() *TestContextStore {
	tcs := TestContextStore{}
	tcs.Save(&ctx_model.State{
		Contexts:  map[string]ctx_model.Context{},
		CurrentId: "",
	})
	return &tcs
}

func (store *TestContextStore) Apply(fn ctx_model.StatePatch) error {
	state := store.Load()
	err := fn(&state)
	if err != nil {
		return err
	} else {
		store.Save(&state)
		return nil
	}
}

func (store *TestContextStore) Read(fn ctx_model.StatePatch) error {
	state := store.Load()
	return fn(&state)
}

func NewTestEventsStore() *TestEventsStore {
	tes := TestEventsStore{}
	tes.Save(&ctx_model.EventRegistry{
		Events: []ctx_model.Event{},
	})
	return &tes
}

func (es *TestEventsStore) Apply(fn ctx_model.EventsPatch) error {
	registry := es.Load()
	err := fn(&registry)
	if err != nil {
		return err
	} else {
		es.Save(&registry)
		return nil
	}
}

func (es *TestEventsStore) Read(fn ctx_model.EventsPatch) error {
	registry := es.Load()
	return fn(&registry)
}

func NewTestArchiveStore() *TestArchiveStore {
	tas := TestArchiveStore{}
	tas.Save(map[string]ctx_model.ContextArchive{})
	tas.SaveEvents(&ctx_model.EventsArchive{})
	return &tas
}

func (as *TestArchiveStore) Load() map[string]ctx_model.ContextArchive {
	state := map[string]ctx_model.ContextArchive{}
	err := json.Unmarshal(as.store, &state)
	if err != nil {
		log.Fatal("Unable to parse archvie store")
	}

	return state
}

func (as *TestArchiveStore) Save(archive map[string]ctx_model.ContextArchive) {
	data, err := json.Marshal(archive)
	if err != nil {
		panic(err)
	}
	as.store = data
}

func (as *TestArchiveStore) LoadEvents() *ctx_model.EventsArchive {
	events := ctx_model.EventsArchive{}
	err := json.Unmarshal(as.storeEvents, &events)
	if err != nil {
		log.Fatal("Unable to parse events archvie store")
	}

	return &events
}

func (as *TestArchiveStore) SaveEvents(entry *ctx_model.EventsArchive) {
	data, err := json.Marshal(entry)
	if err != nil {
		panic(err)
	}
	as.storeEvents = data
}

func (store *TestArchiveStore) Apply(id string, fn ctx_model.ArchivePatch) error {
	archive := store.Load()
	context := archive[id]
	if _, ok := archive[id]; !ok {
		context = ctx_model.ContextArchive{
			Context: ctx_model.Context{
				Id: id,
			},
		}
	}

	err := fn(&context)
	archive[id] = context
	if err != nil {
		return err
	} else {
		store.Save(archive)
		return nil
	}
}

func (store *TestArchiveStore) ApplyEvents(date string, fn ctx_model.ArchiveEventsPatch) error {
	archive := store.LoadEvents()
	err := fn(archive)
	if err != nil {
		return err
	} else {
		store.SaveEvents(archive)
		return nil
	}
}

func (es *TestEventsStore) Save(eventsRegistry *ctx_model.EventRegistry) {
	data, err := json.Marshal(eventsRegistry)
	if err != nil {
		panic(err)
	}
	es.store = data
}

func (es *TestEventsStore) Load() ctx_model.EventRegistry {
	eventsRegistry := ctx_model.EventRegistry{}
	err := json.Unmarshal(es.store, &eventsRegistry)
	if err != nil {
		log.Fatal("Unable to parse state store")
	}

	return eventsRegistry
}

func TestCreateContext(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTestEventsStore(), NewTestArchiveStore(), NewTimer())
	cm.CreateContext(test.TestId, test.TestDescription)

	state := cs.Load()
	assert.Len(t, state.Contexts, 1)
	createdContext := state.Contexts[test.TestId]
	assert.Equal(t, createdContext.Id, test.TestId)
	assert.Equal(t, createdContext.Description, test.TestDescription)
	assert.Equal(t, createdContext.State, ctx_model.ACTIVE)
	assert.Equal(t, createdContext.Duration, time.Duration(0))
	assert.Len(t, createdContext.Intervals, 0)
	assert.Len(t, createdContext.Comments, 0)
}

func TestCreateExistingId(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTestEventsStore(), NewTestArchiveStore(), NewTimer())
	cm.CreateContext(test.TestId, test.TestDescription)
	err := cm.CreateContext(test.TestId, test.TestDescription)

	assert.Error(t, err, errors.New("context already exists"))
}

func TestDontCreateContextWithEmptyDescription(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTestEventsStore(), NewTestArchiveStore(), NewTimer())
	err := cm.CreateContext(test.TestId, "  \t")

	assert.Error(t, err, errors.New("empty description"))
}

func TestDontCreateContextWithEmptyId(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTestEventsStore(), NewTestArchiveStore(), NewTimer())
	err := cm.CreateContext(" \t", test.TestDescription)

	assert.Error(t, err, errors.New("empty id"))
}

func TestEmitCreateEvent(t *testing.T) {
	loc, err := time.LoadLocation(ctx_model.DetectTimezoneName())
	if err != nil {
		loc = time.UTC
	}
	dt1, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:00:00", loc)
	es := NewTestEventsStore()
	cm := New(NewTestContextStore(), es, NewTestArchiveStore(), NewTestTimerProvider("2025-03-13 13:00:00"))
	cm.CreateContext(test.TestId, test.TestDescription)

	registry := es.Load()
	assert.Len(t, registry.Events, 1)
	assert.Equal(t, registry.Events[0].Type, ctx_model.CREATE_CTX)
	assert.Equal(t, registry.Events[0].CtxId, test.TestId)
	assert.Equal(t, registry.Events[0].DateTime, ctx_model.ZonedTime{Time: dt1, Timezone: ctx_model.DetectTimezoneName()})
}

func TestSwitchContext(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTestEventsStore(), NewTestArchiveStore(), NewTimer())
	cm.CreateContext(test.TestId, test.TestDescription)

	err := cm.Switch(test.TestId)

	state := cs.Load()
	assert.NoError(t, err)
	assert.Equal(t, test.TestId, state.CurrentId)

}

func TestDontSwitchWithEmptyId(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTestEventsStore(), NewTestArchiveStore(), NewTimer())
	cm.CreateContext(test.TestId, test.TestDescription)

	err := cm.Switch("\t")

	state := cs.Load()
	assert.Error(t, err, errors.New("empty id"))
	assert.Equal(t, "", state.CurrentId)
}

func TestDontSwitchIfDoesNotExists(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTestEventsStore(), NewTestArchiveStore(), NewTimer())
	cm.CreateContext(test.TestId, test.TestDescription)
	cm.Switch(test.TestId)
	err := cm.Switch("test")

	state := cs.Load()
	assert.Error(t, err, errors.New("context does not exist"))
	assert.Equal(t, test.TestId, state.CurrentId)
}

func TestSwitchNotExistingContext(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTestEventsStore(), NewTestArchiveStore(), NewTimer())
	cm.CreateContext(test.PrevTestId, test.TestDescription)

	cm.Switch(test.PrevTestId)
	err := cm.Switch(test.TestId)

	state := cs.Load()
	assert.Error(t, err, errors.New("context does not exist"))
	assert.Equal(t, test.PrevTestId, state.CurrentId)

}
func TestSwitchCreateIfNotExists(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTestEventsStore(), NewTestArchiveStore(), NewTimer())

	err := cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)

	state := cs.Load()
	assert.NoError(t, err)
	assert.Equal(t, test.TestId, state.CurrentId)
	assert.NotNil(t, state.Contexts[test.TestId])
	assert.Len(t, state.Contexts, 1)
	createdContext := state.Contexts[test.TestId]
	assert.Equal(t, createdContext.Id, test.TestId)
	assert.Equal(t, createdContext.Description, test.TestDescription)
	assert.Equal(t, createdContext.State, ctx_model.ACTIVE)
	assert.Equal(t, createdContext.Duration, time.Duration(0))
	assert.Len(t, createdContext.Intervals, 1)
	assert.Len(t, createdContext.Comments, 0)

}

func TestDontSwitchOrCreateWithEmptyId(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTestEventsStore(), NewTestArchiveStore(), NewTimer())

	err := cm.CreateIfNotExistsAndSwitch("\t", test.TestDescription)

	state := cs.Load()
	assert.Error(t, err, errors.New("empty id"))
	assert.Equal(t, "", state.CurrentId)
}

func TestDontSwitchOrCreateWithEmptyDescription(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTestEventsStore(), NewTestArchiveStore(), NewTimer())

	err := cm.CreateIfNotExistsAndSwitch(test.TestId, " \t ")

	state := cs.Load()
	assert.Error(t, err, errors.New("empty id"))
	assert.Equal(t, "", state.CurrentId)
}

func TestSwitchCreateIfNotExistsOnExistingContext(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTestEventsStore(), NewTestArchiveStore(), NewTimer())
	cm.CreateContext(test.TestId, test.TestDescription)

	assert.Equal(t, cs.Load().CurrentId, "")
	err := cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)

	assert.Equal(t, cs.Load().CurrentId, test.TestId)
	assert.NoError(t, err)
}

func TestSwitchAlreadyActiveContext(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTestEventsStore(), NewTestArchiveStore(), NewTimer())
	cm.CreateContext(test.TestId, test.TestDescription)

	err := cm.Switch(test.TestId)
	assert.NoError(t, err)

	err = cm.Switch(test.TestId)
	state := cs.Load()
	assert.Error(t, err, errors.New("context already active"))
	assert.Len(t, state.Contexts[test.TestId].Intervals, 1)

	err = cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
	state = cs.Load()
	assert.Error(t, err, errors.New("context already active"))
	assert.Len(t, state.Contexts[test.TestId].Intervals, 1)

}

func TestIntervals(t *testing.T) {
	cs := NewTestContextStore()
	loc, err := time.LoadLocation(ctx_model.DetectTimezoneName())
	if err != nil {
		loc = time.UTC
	}
	dt1, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:00:00", loc)
	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:05:00", loc)
	dt3, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:10:00", loc)

	tp := NewTestTimerProvider("2025-03-13 13:00:00")
	cm := New(cs, NewTestEventsStore(), NewTestArchiveStore(), tp)
	cm.CreateContext(test.TestId, test.TestDescription)

	tp.Current = ctx_model.ZonedTime{Time: dt1, Timezone: ctx_model.DetectTimezoneName()}
	cm.Switch(test.TestId)
	state := cs.Load()
	assert.Equal(t, test.TestId, state.CurrentId)
	prevCtx := state.Contexts[state.CurrentId]
	assert.Len(t, prevCtx.Intervals, 1)
	assert.Equal(t, prevCtx.Intervals[0].Start, tp.Current)
	assert.True(t, prevCtx.Intervals[0].End.Time.IsZero())

	tp.Current = ctx_model.ZonedTime{Time: dt2, Timezone: ctx_model.DetectTimezoneName()}
	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
	state = cs.Load()
	prevCtx = state.Contexts[test.TestId]
	assert.Equal(t, prevCtx.Intervals[0].Start, ctx_model.ZonedTime{Time: dt1, Timezone: ctx_model.DetectTimezoneName()})
	assert.Equal(t, prevCtx.Intervals[0].End, ctx_model.ZonedTime{Time: dt2, Timezone: ctx_model.DetectTimezoneName()})
	assert.Equal(t, test.PrevTestId, state.CurrentId)
	nextCtx := state.Contexts[state.CurrentId]
	assert.Len(t, nextCtx.Intervals, 1)
	assert.Equal(t, nextCtx.Intervals[0].Start, ctx_model.ZonedTime{Time: dt2, Timezone: ctx_model.DetectTimezoneName()})
	assert.True(t, nextCtx.Intervals[0].End.Time.IsZero())

	tp.Current = ctx_model.ZonedTime{Time: dt3, Timezone: ctx_model.DetectTimezoneName()}
	cm.Switch(test.TestId)
	state = cs.Load()
	nextCtx = state.Contexts[test.PrevTestId]
	assert.Equal(t, nextCtx.Intervals[0].Start, ctx_model.ZonedTime{Time: dt2, Timezone: ctx_model.DetectTimezoneName()})
	assert.Equal(t, nextCtx.Intervals[0].End, ctx_model.ZonedTime{Time: dt3, Timezone: ctx_model.DetectTimezoneName()})
	assert.Equal(t, test.TestId, state.CurrentId)
	prevCtx = state.Contexts[state.CurrentId]
	assert.Len(t, prevCtx.Intervals, 2)
	assert.Equal(t, prevCtx.Intervals[1].Start, ctx_model.ZonedTime{Time: dt3, Timezone: ctx_model.DetectTimezoneName()})
	assert.True(t, prevCtx.Intervals[1].End.Time.IsZero())

}

func TestEventsFlow(t *testing.T) {
	es := NewTestEventsStore()
	loc, err := time.LoadLocation(ctx_model.DetectTimezoneName())
	if err != nil {
		loc = time.UTC // fallback
	}
	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:05:00", loc)
	dt3, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:10:00", loc)

	tp := NewTestTimerProvider("2025-03-13 13:00:00")
	cm := New(NewTestContextStore(), es, NewTestArchiveStore(), tp)
	cm.CreateIfNotExistsAndSwitch(test.TestId, test.PrevDescription)
	tp.Current = ctx_model.ZonedTime{Time: dt2, Timezone: ctx_model.DetectTimezoneName()}
	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
	tp.Current = ctx_model.ZonedTime{Time: dt3, Timezone: ctx_model.DetectTimezoneName()}
	cm.Switch(test.TestId)

	registry := es.Load()
	assert.Len(t, registry.Events, 10)
	assert.Equal(t, registry.Events[0].Type, ctx_model.CREATE_CTX)
	assert.Equal(t, registry.Events[1].Type, ctx_model.SWITCH_CTX)
	assert.Equal(t, registry.Events[2].Type, ctx_model.START_INTERVAL)
	assert.Equal(t, registry.Events[3].Type, ctx_model.CREATE_CTX)
	assert.Equal(t, registry.Events[4].Type, ctx_model.END_INTERVAL)
	assert.Equal(t, registry.Events[5].Type, ctx_model.SWITCH_CTX)
	assert.Equal(t, registry.Events[6].Type, ctx_model.START_INTERVAL)
	assert.Equal(t, registry.Events[7].Type, ctx_model.END_INTERVAL)
	assert.Equal(t, registry.Events[8].Type, ctx_model.SWITCH_CTX)
	assert.Equal(t, registry.Events[9].Type, ctx_model.START_INTERVAL)

}

func TestFree(t *testing.T) {
	cs := NewTestContextStore()
	loc, err := time.LoadLocation(ctx_model.DetectTimezoneName())
	if err != nil {
		loc = time.UTC
	}
	cm := New(cs, NewTestEventsStore(), NewTestArchiveStore(), NewTestTimerProvider("2025-03-13 13:00:00"))
	dt, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:00:00", loc)

	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
	assert.Equal(t, test.TestId, cs.Load().CurrentId)

	cm.Free()
	state := cs.Load()
	assert.Equal(t, "", state.CurrentId)
	assert.Equal(t, ctx_model.ZonedTime{Time: dt, Timezone: ctx_model.DetectTimezoneName()}, state.Contexts[test.TestId].Intervals[0].Start)
	assert.Equal(t, ctx_model.ZonedTime{Time: dt, Timezone: ctx_model.DetectTimezoneName()}, state.Contexts[test.TestId].Intervals[0].End)

}

func TestFreeWithNowCurrentContext(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTestEventsStore(), NewTestArchiveStore(), NewTimer())
	cm.CreateContext(test.TestId, test.TestDescription)
	assert.Equal(t, "", cs.Load().CurrentId)
	err := cm.Free()
	assert.Error(t, err, errors.New("no active context"))
}

func TestDeleteContext(t *testing.T) {
	cs := NewTestContextStore()
	es := NewTestEventsStore()
	cm := New(cs, es, NewTestArchiveStore(), NewTimer())
	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
	cm.Free()
	assert.Len(t, cs.Load().Contexts, 1)
	assert.Len(t, es.Load().Events, 4)
	err := cm.Delete(test.TestId)
	assert.NoError(t, err)
	assert.Len(t, es.Load().Events, 5)
	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Type, ctx_model.DELETE_CTX)
	assert.Len(t, cs.Load().Contexts, 0)
}

func TestDeleteActiveContext(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTestEventsStore(), NewTestArchiveStore(), NewTimer())
	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
	assert.Len(t, cs.Load().Contexts, 1)
	err := cm.Delete(test.TestId)
	assert.Error(t, err, errors.New("context is active"))
	assert.Len(t, cs.Load().Contexts, 1)
}

func TestDeleteNotExistingContext(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTestEventsStore(), NewTestArchiveStore(), NewTimer())
	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
	err := cm.Delete(test.PrevTestId)
	assert.Error(t, err, errors.New("context does not exist"))
	assert.Len(t, cs.Load().Contexts, 1)
}

func TestEventFilter(t *testing.T) {
	es := NewTestEventsStore()
	cm := New(NewTestContextStore(), es, NewTestArchiveStore(), NewTimer())
	loc, err := time.LoadLocation(ctx_model.DetectTimezoneName())
	if err != nil {
		loc = time.UTC
	}
	dt1, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:00:00", loc)
	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-14 13:00:00", loc)
	cm.EventsStore.Apply(func(er *ctx_model.EventRegistry) error {
		er.Events = append(er.Events, ctx_model.Event{
			DateTime: ctx_model.ZonedTime{Time: dt1}, Description: "test1", Type: ctx_model.CREATE_CTX,
		})
		er.Events = append(er.Events, ctx_model.Event{
			DateTime: ctx_model.ZonedTime{Time: dt2}, Description: "test2", Type: ctx_model.SWITCH_CTX,
		})
		er.Events = append(er.Events, ctx_model.Event{
			DateTime: ctx_model.ZonedTime{Time: dt1}, Description: "test3", Type: ctx_model.SWITCH_CTX,
		})
		er.Events = append(er.Events, ctx_model.Event{
			DateTime: ctx_model.ZonedTime{Time: dt2}, Description: "test4", Type: ctx_model.START_INTERVAL,
		})
		return nil
	})

	er := es.Load()
	assert.Len(t, er.Events, 4)
	events := cm.filterEvents(&er, ctx_model.EventsFilter{
		Date: "2025-03-14",
	})

	assert.Len(t, events, 2)
	assert.Equal(t, ctx_model.ZonedTime{Time: dt2, Timezone: ctx_model.DetectTimezoneName()}, events[0].DateTime)
	assert.Equal(t, "test2", events[0].Description)
	assert.Equal(t, ctx_model.ZonedTime{Time: dt2, Timezone: ctx_model.DetectTimezoneName()}, events[1].DateTime)
	assert.Equal(t, "test4", events[1].Description)

	events = cm.filterEvents(&er, ctx_model.EventsFilter{
		Date:  "2025-03-13",
		Types: []string{"CREATE"},
	})

	assert.Len(t, events, 1)
	assert.Equal(t, ctx_model.ZonedTime{Time: dt1, Timezone: ctx_model.DetectTimezoneName()}, events[0].DateTime)
	assert.Equal(t, "test1", events[0].Description)

	events = cm.filterEvents(&er, ctx_model.EventsFilter{
		Types: []string{"SWITCH"},
	})

	assert.Len(t, events, 2)
	assert.Equal(t, ctx_model.ZonedTime{Time: dt2, Timezone: ctx_model.DetectTimezoneName()}, events[0].DateTime)
	assert.Equal(t, "test2", events[0].Description)
	assert.Equal(t, ctx_model.ZonedTime{Time: dt1, Timezone: ctx_model.DetectTimezoneName()}, events[1].DateTime)
	assert.Equal(t, "test3", events[1].Description)

	events = cm.filterEvents(&er, ctx_model.EventsFilter{
		Types: []string{"CREATE", "START_INTERVAL"},
	})

	assert.Len(t, events, 2)
	assert.Equal(t, ctx_model.ZonedTime{Time: dt1, Timezone: ctx_model.DetectTimezoneName()}, events[0].DateTime)
	assert.Equal(t, "test1", events[0].Description)
	assert.Equal(t, ctx_model.ZonedTime{Time: dt2, Timezone: ctx_model.DetectTimezoneName()}, events[1].DateTime)
	assert.Equal(t, "test4", events[1].Description)
}

func TestArchiveContext(t *testing.T) {
	as := NewTestArchiveStore()
	cs := NewTestContextStore()
	tp := NewTestTimerProvider("2025-03-13 13:00:00")
	es := NewTestEventsStore()
	cm := New(cs, es, as, tp)
	loc, err := time.LoadLocation(ctx_model.DetectTimezoneName())
	if err != nil {
		loc = time.UTC
	}
	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:05:00", loc)
	dt3, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:10:00", loc)

	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
	tp.Current = ctx_model.ZonedTime{Time: dt2, Timezone: ctx_model.DetectTimezoneName()}
	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
	tp.Current = ctx_model.ZonedTime{Time: dt3, Timezone: ctx_model.DetectTimezoneName()}
	cm.Switch(test.TestId)
	cm.Switch(test.PrevTestId)
	state := cs.Load()
	assert.Len(t, state.Contexts, 2)
	assert.Len(t, es.Load().Events, 13)
	err = cm.Archive(test.TestId)

	assert.NoError(t, err)
	archive := as.Load()[test.TestId]
	state = cs.Load()
	assert.Equal(t, archive.Context.Id, test.TestId)
	assert.Equal(t, archive.Context.Description, test.TestDescription)
	assert.Len(t, state.Contexts, 1)
	assert.Len(t, es.Load().Events, 14)
}

func TestDontArchiveActiveContext(t *testing.T) {
	as := NewTestArchiveStore()
	cs := NewTestContextStore()
	tp := NewTestTimerProvider("2025-03-13 13:00:00")
	es := NewTestEventsStore()
	cm := New(cs, es, as, tp)
	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
	err := cm.Archive(test.TestId)

	assert.Error(t, err, errors.New("context is active"))
}

func TestArchiveAll(t *testing.T) {
	as := NewTestArchiveStore()
	cs := NewTestContextStore()
	tp := NewTestTimerProvider("2025-03-13 13:00:00")
	es := NewTestEventsStore()
	cm := New(cs, es, as, tp)
	loc, err := time.LoadLocation(ctx_model.DetectTimezoneName())
	if err != nil {
		loc = time.UTC
	}
	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:05:00", loc)
	dt3, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:10:00", loc)

	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
	tp.Current = ctx_model.ZonedTime{Time: dt2, Timezone: ctx_model.DetectTimezoneName()}
	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
	tp.Current = ctx_model.ZonedTime{Time: dt3, Timezone: ctx_model.DetectTimezoneName()}
	cm.Switch(test.TestId)
	cm.Switch(test.PrevTestId)
	cm.Free()
	assert.Len(t, es.Load().Events, 14)
	assert.Len(t, cs.Load().Contexts, 2)
	err = cm.ArchiveAll()

	assert.NoError(t, err)
	events := as.LoadEvents().Events
	assert.Len(t, events, 16)
	assert.Len(t, cs.Load().Contexts, 0)
}

func TestMergeContexts(t *testing.T) {
	as := NewTestArchiveStore()
	cs := NewTestContextStore()
	tp := NewTestTimerProvider("2025-03-13 13:00:00")
	es := NewTestEventsStore()
	cm := New(cs, es, as, tp)
	loc, err := time.LoadLocation(ctx_model.DetectTimezoneName())
	if err != nil {
		loc = time.UTC
	}
	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:05:00", loc)
	dt3, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:10:00", loc)

	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
	cm.CreateIfNotExistsAndSwitch(test.TestIdExtra, test.PrevDescription)
	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
	tp.Current = ctx_model.ZonedTime{Time: dt2, Timezone: ctx_model.DetectTimezoneName()}
	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
	tp.Current = ctx_model.ZonedTime{Time: dt3, Timezone: ctx_model.DetectTimezoneName()}
	cm.Switch(test.TestId)
	cm.Switch(test.PrevTestId)
	cm.Free()
	assert.Len(t, es.Load().Events, 21)
	assert.Len(t, cs.Load().Contexts, 3)
	assert.Equal(t, cs.Load().Contexts[test.PrevTestId].Duration, dt3.Sub(dt2))

	err = cm.MergeContext(test.TestId, test.PrevTestId)

	assert.NoError(t, err)
	assert.Equal(t, cs.Load().Contexts[test.PrevTestId].Duration, dt3.Sub(dt2)*2)
	events := es.Load().Events
	assert.Len(t, events, 23)
	assert.Len(t, cs.Load().Contexts, 2)
	for _, event := range events {
		if event.Type == ctx_model.SWITCH_CTX {
			if v, ok := event.Data["from"]; ok && v != "" && event.CtxId == test.TestIdExtra {
				assert.Equal(t, v, test.TestId)
			}
		}
	}
}

func TestArchiveAllEvents(t *testing.T) {
	as := NewTestArchiveStore()
	cs := NewTestContextStore()
	tp := NewTestTimerProvider("2025-03-13 13:00:00")
	es := NewTestEventsStore()
	cm := New(cs, es, as, tp)
	loc, err := time.LoadLocation(ctx_model.DetectTimezoneName())
	if err != nil {
		loc = time.UTC
	}
	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:05:00", loc)
	dt3, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:10:00", loc)

	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
	tp.Current = ctx_model.ZonedTime{Time: dt2, Timezone: ctx_model.DetectTimezoneName()}
	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
	tp.Current = ctx_model.ZonedTime{Time: dt3, Timezone: ctx_model.DetectTimezoneName()}
	cm.Switch(test.TestId)
	cm.Switch(test.PrevTestId)
	cm.Free()
	assert.Len(t, es.Load().Events, 14)
	assert.Len(t, cs.Load().Contexts, 2)
	err = cm.ArchiveAllEvents()

	assert.NoError(t, err)
	events := as.LoadEvents().Events
	assert.Len(t, events, 14)
	assert.Len(t, es.Load().Events, 0)
}

func TestErrorOnEditCurrentContextInterval(t *testing.T) {
	as := NewTestArchiveStore()
	cs := NewTestContextStore()
	tp := NewTestTimerProvider("2025-03-13 13:00:00")
	es := NewTestEventsStore()
	cm := New(cs, es, as, tp)

	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)

	err := cm.EditContextInterval(test.TestId, 0, ctx_model.ZonedTime{Time: time.Now().Local()}, ctx_model.ZonedTime{Time: time.Now().Local()})

	assert.Error(t, err, errors.New("context is active"))

}

func TestEditContextInterval(t *testing.T) {
	as := NewTestArchiveStore()
	cs := NewTestContextStore()
	tp := NewTestTimerProvider("2025-03-13 13:00:00")
	es := NewTestEventsStore()
	cm := New(cs, es, as, tp)
	loc, err := time.LoadLocation(ctx_model.DetectTimezoneName())
	if err != nil {
		loc = time.UTC
	}
	dt1, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:00:00", loc)
	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:05:00", loc)
	dt3, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:10:00", loc)

	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
	tp.Current = ctx_model.ZonedTime{Time: dt2, Timezone: ctx_model.DetectTimezoneName()}
	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.TestDescription)
	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
	tp.Current = ctx_model.ZonedTime{Time: dt3, Timezone: ctx_model.DetectTimezoneName()}
	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.TestDescription)

	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[0].Start, ctx_model.ZonedTime{Time: dt1, Timezone: ctx_model.DetectTimezoneName()})
	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[0].End, ctx_model.ZonedTime{Time: dt2, Timezone: ctx_model.DetectTimezoneName()})
	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[0].Duration, dt2.Sub(dt1))
	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[1].Start, ctx_model.ZonedTime{Time: dt2, Timezone: ctx_model.DetectTimezoneName()})
	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[1].End, ctx_model.ZonedTime{Time: dt3, Timezone: ctx_model.DetectTimezoneName()})
	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[1].Duration, dt3.Sub(dt2))
	assert.Equal(t, cs.Load().Contexts[test.TestId].Duration, cs.Load().Contexts[test.TestId].Intervals[0].Duration+cs.Load().Contexts[test.TestId].Intervals[1].Duration)

	err = cm.EditContextInterval(test.TestId, 0, ctx_model.ZonedTime{Time: dt1}, ctx_model.ZonedTime{Time: dt3, Timezone: ctx_model.DetectTimezoneName()})
	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[0].Start, ctx_model.ZonedTime{Time: dt1, Timezone: ctx_model.DetectTimezoneName()})
	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[0].End, ctx_model.ZonedTime{Time: dt3, Timezone: ctx_model.DetectTimezoneName()})
	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[0].Duration, dt3.Sub(dt1))
	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[1].Start, ctx_model.ZonedTime{Time: dt2, Timezone: ctx_model.DetectTimezoneName()})
	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[1].End, ctx_model.ZonedTime{Time: dt3, Timezone: ctx_model.DetectTimezoneName()})
	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[1].Duration, dt3.Sub(dt2))
	assert.Equal(t, cs.Load().Contexts[test.TestId].Duration, cs.Load().Contexts[test.TestId].Intervals[0].Duration+cs.Load().Contexts[test.TestId].Intervals[1].Duration)
	assert.NoError(t, err, errors.New("context is active"))
	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Type, ctx_model.EDIT_CTX_INTERVAL)
	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Data["old.start"], dt1.Format(time.RFC3339))
	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Data["old.end"], dt2.Format(time.RFC3339))
	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Data["new.start"], dt1.Format(time.RFC3339))
	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Data["new.end"], dt3.Format(time.RFC3339))

	err = cm.EditContextInterval(test.TestId, 0, ctx_model.ZonedTime{Time: dt2, Timezone: ctx_model.DetectTimezoneName()}, ctx_model.ZonedTime{Time: dt3, Timezone: ctx_model.DetectTimezoneName()})
	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[0].Start, ctx_model.ZonedTime{Time: dt2, Timezone: ctx_model.DetectTimezoneName()})
	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[0].End, ctx_model.ZonedTime{Time: dt3, Timezone: ctx_model.DetectTimezoneName()})
	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[0].Duration, dt3.Sub(dt2))
	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[1].Start, ctx_model.ZonedTime{Time: dt2, Timezone: ctx_model.DetectTimezoneName()})
	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[1].End, ctx_model.ZonedTime{Time: dt3, Timezone: ctx_model.DetectTimezoneName()})
	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[1].Duration, dt3.Sub(dt2))
	assert.Equal(t, cs.Load().Contexts[test.TestId].Duration, cs.Load().Contexts[test.TestId].Intervals[0].Duration+cs.Load().Contexts[test.TestId].Intervals[1].Duration)
	assert.NoError(t, err, errors.New("context is active"))
	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Type, ctx_model.EDIT_CTX_INTERVAL)
	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Data["old.start"], dt1.Format(time.RFC3339))
	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Data["old.end"], dt3.Format(time.RFC3339))
	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Data["new.start"], dt2.Format(time.RFC3339))
	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Data["new.end"], dt3.Format(time.RFC3339))
}

func TestRename(t *testing.T) {
	as := NewTestArchiveStore()
	cs := NewTestContextStore()
	tp := NewTestTimerProvider("2025-03-13 13:00:00")
	es := NewTestEventsStore()
	cm := New(cs, es, as, tp)

	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)

	cm.RenameContext(test.TestId, test.PrevTestId, test.PrevDescription)

	state := cs.Load()
	assert.Contains(t, state.Contexts, test.PrevTestId)
	assert.NotContains(t, state.Contexts, test.TestId)
	assert.Len(t, state.Contexts[test.PrevTestId].Intervals, 1)
	assert.Equal(t, state.Contexts[test.PrevTestId].Description, test.PrevDescription)

}

func TestGetIntervalDurationForDateInBetween(t *testing.T) {
	cs := NewTestContextStore()
	tp := NewTestTimerProvider("2025-03-13 13:00:00")
	es := NewTestEventsStore()
	cm := New(cs, es, NewTestArchiveStore(), tp)
	loc, err := time.LoadLocation(ctx_model.DetectTimezoneName())
	if err != nil {
		loc = time.UTC
	}
	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-15 13:05:00", loc)
	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
	tp.Current = ctx_model.ZonedTime{Time: dt2, Timezone: ctx_model.DetectTimezoneName()}
	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
	state := cs.Load()
	date, _ := time.ParseInLocation(time.DateOnly, "2025-03-14", loc)
	duration, err := cm.GetIntervalDurationsByDate(&state, test.TestId, ctx_model.ZonedTime{Time: date, Timezone: ctx_model.DetectTimezoneName()})
	assert.NoError(t, err)
	assert.Equal(t, 24*time.Hour, duration)
}

func TestGetIntervalDurationForDateOutOfBounds(t *testing.T) {
	cs := NewTestContextStore()
	tp := NewTestTimerProvider("2025-03-13 13:00:00")
	es := NewTestEventsStore()
	loc, err := time.LoadLocation(ctx_model.DetectTimezoneName())
	if err != nil {
		loc = time.UTC
	}
	cm := New(cs, es, NewTestArchiveStore(), tp)
	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-15 13:05:00", loc)
	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
	tp.Current = ctx_model.ZonedTime{Time: dt2, Timezone: ctx_model.DetectTimezoneName()}
	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
	state := cs.Load()
	date, _ := time.ParseInLocation(time.DateOnly, "2025-03-16", loc)
	duration, err := cm.GetIntervalDurationsByDate(&state, test.TestId, ctx_model.ZonedTime{Time: date, Timezone: ctx_model.DetectTimezoneName()})
	assert.NoError(t, err)
	assert.Equal(t, time.Duration(0), duration)
	date, _ = time.ParseInLocation(time.DateOnly, "2025-03-12", loc)
	duration, err = cm.GetIntervalDurationsByDate(&state, test.TestId, ctx_model.ZonedTime{Time: date, Timezone: ctx_model.DetectTimezoneName()})
	assert.NoError(t, err)
	assert.Equal(t, time.Duration(0), duration)
}

func TestGetIntervalDurationForDateBefore(t *testing.T) {
	cs := NewTestContextStore()
	tp := NewTestTimerProvider("2025-03-13 10:00:00")
	es := NewTestEventsStore()
	loc, err := time.LoadLocation(ctx_model.DetectTimezoneName())
	if err != nil {
		loc = time.UTC
	}
	cm := New(cs, es, NewTestArchiveStore(), tp)
	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-15 13:00:00", loc)
	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
	tp.Current = ctx_model.ZonedTime{Time: dt2, Timezone: ctx_model.DetectTimezoneName()}
	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
	state := cs.Load()
	date, _ := time.ParseInLocation(time.DateOnly, "2025-03-15", loc)
	duration, err := cm.GetIntervalDurationsByDate(&state, test.TestId, ctx_model.ZonedTime{Time: date, Timezone: ctx_model.DetectTimezoneName()})
	assert.NoError(t, err)
	assert.Equal(t, 13*time.Hour, duration)
}

func TestGetIntervalDurationForDateAfter(t *testing.T) {
	cs := NewTestContextStore()
	tp := NewTestTimerProvider("2025-03-13 10:00:00")
	es := NewTestEventsStore()
	loc, err := time.LoadLocation(ctx_model.DetectTimezoneName())
	if err != nil {
		loc = time.UTC
	}
	cm := New(cs, es, NewTestArchiveStore(), tp)
	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-15 13:00:00", loc)
	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
	tp.Current = ctx_model.ZonedTime{Time: dt2, Timezone: ctx_model.DetectTimezoneName()}
	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
	state := cs.Load()
	date, _ := time.ParseInLocation(time.DateOnly, "2025-03-13", loc)
	duration, err := cm.GetIntervalDurationsByDate(&state, test.TestId, ctx_model.ZonedTime{Time: date, Timezone: ctx_model.DetectTimezoneName()})
	assert.NoError(t, err)
	assert.Equal(t, 14*time.Hour, duration)
}

func TestGetIntervalDurationForDateEqual(t *testing.T) {
	cs := NewTestContextStore()
	tp := NewTestTimerProvider("2025-03-13 10:00:00")
	es := NewTestEventsStore()
	loc, err := time.LoadLocation(ctx_model.DetectTimezoneName())
	if err != nil {
		loc = time.UTC
	}
	cm := New(cs, es, NewTestArchiveStore(), tp)
	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:00:00", loc)
	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
	tp.Current = ctx_model.ZonedTime{Time: dt2, Timezone: ctx_model.DetectTimezoneName()}
	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
	state := cs.Load()
	date, _ := time.ParseInLocation(time.DateOnly, "2025-03-13", loc)
	duration, err := cm.GetIntervalDurationsByDate(&state, test.TestId, ctx_model.ZonedTime{Time: date, Timezone: ctx_model.DetectTimezoneName()})
	assert.NoError(t, err)
	assert.Equal(t, 3*time.Hour, duration)
}

func TestDeleteINtervalActiveContext(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTestEventsStore(), NewTestArchiveStore(), NewTimer())
	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
	err := cm.DeleteInterval(test.TestId, 0)

	assert.Error(t, err, errors.New("context is active"))
}

func TestDeleteIntervalNotExistingContext(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTestEventsStore(), NewTestArchiveStore(), NewTimer())
	err := cm.DeleteInterval(test.TestId, 0)

	assert.Error(t, err, errors.New("context does not exist"))
}

func TestDeleteIntervalOutOfBounds(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTestEventsStore(), NewTestArchiveStore(), NewTimer())
	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
	err := cm.DeleteInterval(test.TestId, 0)

	assert.Error(t, err, errors.New("interval out of bounds"))
}

func TestDeleteInterval(t *testing.T) {
	cs := NewTestContextStore()
	tp := NewTestTimerProvider("2025-03-13 13:00:00")
	es := NewTestEventsStore()
	loc, err := time.LoadLocation(ctx_model.DetectTimezoneName())
	if err != nil {
		loc = time.UTC
	}
	cm := New(cs, es, NewTestArchiveStore(), tp)
	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:05:00", loc)
	dt3, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:10:00", loc)
	dt4, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:15:00", loc)

	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
	tp.Current = ctx_model.ZonedTime{Time: dt2, Timezone: ctx_model.DetectTimezoneName()}
	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
	tp.Current = ctx_model.ZonedTime{Time: dt3, Timezone: ctx_model.DetectTimezoneName()}
	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
	tp.Current = ctx_model.ZonedTime{Time: dt4, Timezone: ctx_model.DetectTimezoneName()}
	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)

	assert.Len(t, cs.Load().Contexts[test.TestId].Intervals, 2)
	assert.Equal(t, cs.Load().Contexts[test.TestId].Duration.Seconds(), time.Duration(600000000000).Seconds())

	err = cm.DeleteInterval(test.TestId, 0)

	assert.NoError(t, err)
	assert.Len(t, cs.Load().Contexts[test.TestId].Intervals, 1)
	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[0].Start, ctx_model.ZonedTime{Time: dt3, Timezone: ctx_model.DetectTimezoneName()})
	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[0].End, ctx_model.ZonedTime{Time: dt4, Timezone: ctx_model.DetectTimezoneName()})
	assert.Equal(t, cs.Load().Contexts[test.TestId].Duration.Seconds(), time.Duration(300000000000).Seconds())
	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Type, ctx_model.DELETE_CTX_INTERVAL)
}

func TestSearchContextWithRegex(t *testing.T) {
	cs := NewTestContextStore()
	tp := NewTestTimerProvider("2025-03-13 13:00:00")
	es := NewTestEventsStore()
	loc, err := time.LoadLocation(ctx_model.DetectTimezoneName())
	if err != nil {
		loc = time.UTC
	}
	cm := New(cs, es, NewTestArchiveStore(), tp)
	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:05:00", loc)
	dt3, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:10:00", loc)

	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
	tp.Current = ctx_model.ZonedTime{Time: dt2, Timezone: ctx_model.DetectTimezoneName()}
	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
	tp.Current = ctx_model.ZonedTime{Time: dt3, Timezone: ctx_model.DetectTimezoneName()}
	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)

	assert.Len(t, cs.Load().Contexts[test.TestId].Intervals, 2)

	contexts, err := cm.Search("p.*test.*")
	assert.NoError(t, err)
	assert.Len(t, contexts, 1)
	assert.Equal(t, contexts[0].Description, test.PrevDescription)
}

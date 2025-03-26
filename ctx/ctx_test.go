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
	Current time.Time
}

func NewTestTimerProvider(dateTime string) *TestTimeProvider {
	dt, _ := time.Parse(time.DateTime, dateTime)
	return &TestTimeProvider{
		Current: dt,
	}
}

func (provider *TestTimeProvider) Now() time.Time {
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

func NewTestArchiveStore(id string) *TestArchiveStore {
	tas := TestArchiveStore{}
	tas.Save(&ctx_model.ArchiveEntry{
		Context: ctx_model.Context{
			Id: id,
		},
	})
	tas.SaveEvents([]ctx_model.Event{})
	return &tas
}

func (as *TestArchiveStore) Load(id string) ctx_model.ArchiveEntry {
	state := ctx_model.ArchiveEntry{}
	err := json.Unmarshal(as.store, &state)
	if err != nil {
		log.Fatal("Unable to parse archvie store")
	}

	return state
}

func (as *TestArchiveStore) Save(archive *ctx_model.ArchiveEntry) {
	data, err := json.Marshal(archive)
	if err != nil {
		panic(err)
	}
	as.store = data
}

func (as *TestArchiveStore) LoadEvents() []ctx_model.Event {
	events := []ctx_model.Event{}
	err := json.Unmarshal(as.storeEvents, &events)
	if err != nil {
		log.Fatal("Unable to parse events archvie store")
	}

	return events
}

func (as *TestArchiveStore) SaveEvents(events []ctx_model.Event) {
	data, err := json.Marshal(events)
	if err != nil {
		panic(err)
	}
	as.storeEvents = data
}

func (store *TestArchiveStore) Apply(id string, fn ctx_model.ArchivePatch) error {
	archive := store.Load(id)
	err := fn(&archive)
	if err != nil {
		return err
	} else {
		store.Save(&archive)
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
	cm := New(cs, NewTestEventsStore(), NewTestArchiveStore(test.TestId), NewTimer())
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
	cm := New(cs, NewTestEventsStore(), NewTestArchiveStore(test.TestId), NewTimer())
	cm.CreateContext(test.TestId, test.TestDescription)
	err := cm.CreateContext(test.TestId, test.TestDescription)

	assert.Error(t, err, errors.New("context already exists"))
}

func TestDontCreateContextWithEmptyDescription(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTestEventsStore(), NewTestArchiveStore(test.TestId), NewTimer())
	err := cm.CreateContext(test.TestId, "  \t")

	assert.Error(t, err, errors.New("empty description"))
}

func TestDontCreateContextWithEmptyId(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTestEventsStore(), NewTestArchiveStore(test.TestId), NewTimer())
	err := cm.CreateContext(" \t", test.TestDescription)

	assert.Error(t, err, errors.New("empty id"))
}

func TestEmitCreateEvent(t *testing.T) {
	dt1, _ := time.Parse(time.DateTime, "2025-03-13 13:00:00")
	es := NewTestEventsStore()
	cm := New(NewTestContextStore(), es, NewTestArchiveStore(test.TestId), NewTestTimerProvider("2025-03-13 13:00:00"))
	cm.CreateContext(test.TestId, test.TestDescription)

	registry := es.Load()
	assert.Len(t, registry.Events, 1)
	assert.Equal(t, registry.Events[0].Type, ctx_model.CREATE_CTX)
	assert.Equal(t, registry.Events[0].CtxId, test.TestId)
	assert.Equal(t, registry.Events[0].DateTime, dt1)
}

func TestSwitchContext(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTestEventsStore(), NewTestArchiveStore(test.TestId), NewTimer())
	cm.CreateContext(test.TestId, test.TestDescription)

	err := cm.Switch(test.TestId)

	state := cs.Load()
	assert.NoError(t, err)
	assert.Equal(t, test.TestId, state.CurrentId)

}

func TestDontSwitchWithEmptyId(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTestEventsStore(), NewTestArchiveStore(test.TestId), NewTimer())
	cm.CreateContext(test.TestId, test.TestDescription)

	err := cm.Switch("\t")

	state := cs.Load()
	assert.Error(t, err, errors.New("empty id"))
	assert.Equal(t, "", state.CurrentId)
}

func TestDontSwitchIfDoesNotExists(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTestEventsStore(), NewTestArchiveStore(test.TestId), NewTimer())
	cm.CreateContext(test.TestId, test.TestDescription)
	cm.Switch(test.TestId)
	err := cm.Switch("test")

	state := cs.Load()
	assert.Error(t, err, errors.New("context does not exist"))
	assert.Equal(t, test.TestId, state.CurrentId)
}

func TestSwitchNotExistingContext(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTestEventsStore(), NewTestArchiveStore(test.TestId), NewTimer())
	cm.CreateContext(test.PrevTestId, test.TestDescription)

	cm.Switch(test.PrevTestId)
	err := cm.Switch(test.TestId)

	state := cs.Load()
	assert.Error(t, err, errors.New("context does not exist"))
	assert.Equal(t, test.PrevTestId, state.CurrentId)

}
func TestSwitchCreateIfNotExists(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTestEventsStore(), NewTestArchiveStore(test.TestId), NewTimer())

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
	cm := New(cs, NewTestEventsStore(), NewTestArchiveStore(test.TestId), NewTimer())

	err := cm.CreateIfNotExistsAndSwitch("\t", test.TestDescription)

	state := cs.Load()
	assert.Error(t, err, errors.New("empty id"))
	assert.Equal(t, "", state.CurrentId)
}

func TestDontSwitchOrCreateWithEmptyDescription(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTestEventsStore(), NewTestArchiveStore(test.TestId), NewTimer())

	err := cm.CreateIfNotExistsAndSwitch(test.TestId, " \t ")

	state := cs.Load()
	assert.Error(t, err, errors.New("empty id"))
	assert.Equal(t, "", state.CurrentId)
}

func TestSwitchCreateIfNotExistsOnExistingContext(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTestEventsStore(), NewTestArchiveStore(test.TestId), NewTimer())
	cm.CreateContext(test.TestId, test.TestDescription)

	assert.Equal(t, cs.Load().CurrentId, "")
	err := cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)

	assert.Equal(t, cs.Load().CurrentId, test.TestId)
	assert.NoError(t, err)
}

func TestSwitchAlreadyActiveContext(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTestEventsStore(), NewTestArchiveStore(test.TestId), NewTimer())
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
	dt1, _ := time.Parse(time.DateTime, "2025-03-13 13:00:00")
	dt2, _ := time.Parse(time.DateTime, "2025-03-13 13:05:00")
	dt3, _ := time.Parse(time.DateTime, "2025-03-13 13:10:00")

	tp := NewTestTimerProvider("2025-03-13 13:00:00")
	cm := New(cs, NewTestEventsStore(), NewTestArchiveStore(test.TestId), tp)
	cm.CreateContext(test.TestId, test.TestDescription)

	tp.Current = dt1
	cm.Switch(test.TestId)
	state := cs.Load()
	assert.Equal(t, test.TestId, state.CurrentId)
	prevCtx := state.Contexts[state.CurrentId]
	assert.Len(t, prevCtx.Intervals, 1)
	assert.Equal(t, prevCtx.Intervals[0].Start, tp.Current)
	assert.True(t, prevCtx.Intervals[0].End.IsZero())

	tp.Current = dt2
	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
	state = cs.Load()
	prevCtx = state.Contexts[test.TestId]
	assert.Equal(t, prevCtx.Intervals[0].Start, dt1)
	assert.Equal(t, prevCtx.Intervals[0].End, dt2)
	assert.Equal(t, test.PrevTestId, state.CurrentId)
	nextCtx := state.Contexts[state.CurrentId]
	assert.Len(t, nextCtx.Intervals, 1)
	assert.Equal(t, nextCtx.Intervals[0].Start, dt2)
	assert.True(t, nextCtx.Intervals[0].End.IsZero())

	tp.Current = dt3
	cm.Switch(test.TestId)
	state = cs.Load()
	nextCtx = state.Contexts[test.PrevTestId]
	assert.Equal(t, nextCtx.Intervals[0].Start, dt2)
	assert.Equal(t, nextCtx.Intervals[0].End, dt3)
	assert.Equal(t, test.TestId, state.CurrentId)
	prevCtx = state.Contexts[state.CurrentId]
	assert.Len(t, prevCtx.Intervals, 2)
	assert.Equal(t, prevCtx.Intervals[1].Start, dt3)
	assert.True(t, prevCtx.Intervals[1].End.IsZero())

}

func TestEventsFlow(t *testing.T) {
	es := NewTestEventsStore()
	dt2, _ := time.Parse(time.DateTime, "2025-03-13 13:05:00")
	dt3, _ := time.Parse(time.DateTime, "2025-03-13 13:10:00")

	tp := NewTestTimerProvider("2025-03-13 13:00:00")
	cm := New(NewTestContextStore(), es, NewTestArchiveStore(test.TestId), tp)
	cm.CreateIfNotExistsAndSwitch(test.TestId, test.PrevDescription)
	tp.Current = dt2
	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
	tp.Current = dt3
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
	cm := New(cs, NewTestEventsStore(), NewTestArchiveStore(test.TestId), NewTestTimerProvider("2025-03-13 13:00:00"))
	dt, _ := time.Parse(time.DateTime, "2025-03-13 13:00:00")

	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
	assert.Equal(t, test.TestId, cs.Load().CurrentId)

	cm.Free()
	state := cs.Load()
	assert.Equal(t, "", state.CurrentId)
	assert.Equal(t, dt, state.Contexts[test.TestId].Intervals[0].Start)
	assert.Equal(t, dt, state.Contexts[test.TestId].Intervals[0].End)

}

func TestFreeWithNowCurrentContext(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTestEventsStore(), NewTestArchiveStore(test.TestId), NewTimer())
	cm.CreateContext(test.TestId, test.TestDescription)
	assert.Equal(t, "", cs.Load().CurrentId)
	err := cm.Free()
	assert.Error(t, err, errors.New("no active context"))
}

func TestEventFilter(t *testing.T) {
	es := NewTestEventsStore()
	cm := New(NewTestContextStore(), es, NewTestArchiveStore(test.TestId), NewTimer())
	dt1, _ := time.Parse(time.DateTime, "2025-03-13 13:00:00")
	dt2, _ := time.Parse(time.DateTime, "2025-03-14 13:00:00")
	cm.EventsStore.Apply(func(er *ctx_model.EventRegistry) error {
		er.Events = append(er.Events, ctx_model.Event{
			DateTime: dt1, Description: "test1", Type: ctx_model.CREATE_CTX,
		})
		er.Events = append(er.Events, ctx_model.Event{
			DateTime: dt2, Description: "test2", Type: ctx_model.SWITCH_CTX,
		})
		er.Events = append(er.Events, ctx_model.Event{
			DateTime: dt1, Description: "test3", Type: ctx_model.SWITCH_CTX,
		})
		er.Events = append(er.Events, ctx_model.Event{
			DateTime: dt2, Description: "test4", Type: ctx_model.START_INTERVAL,
		})
		return nil
	})

	er := es.Load()
	assert.Len(t, er.Events, 4)
	events := cm.filterEvents(&er, ctx_model.EventsFilter{
		Date: "2025-03-14",
	})

	assert.Len(t, events, 2)
	assert.Equal(t, dt2, events[0].DateTime)
	assert.Equal(t, "test2", events[0].Description)
	assert.Equal(t, dt2, events[1].DateTime)
	assert.Equal(t, "test4", events[1].Description)

	events = cm.filterEvents(&er, ctx_model.EventsFilter{
		Date:  "2025-03-13",
		Types: []string{"CREATE"},
	})

	assert.Len(t, events, 1)
	assert.Equal(t, dt1, events[0].DateTime)
	assert.Equal(t, "test1", events[0].Description)

	events = cm.filterEvents(&er, ctx_model.EventsFilter{
		Types: []string{"SWITCH"},
	})

	assert.Len(t, events, 2)
	assert.Equal(t, dt2, events[0].DateTime)
	assert.Equal(t, "test2", events[0].Description)
	assert.Equal(t, dt1, events[1].DateTime)
	assert.Equal(t, "test3", events[1].Description)

	events = cm.filterEvents(&er, ctx_model.EventsFilter{
		Types: []string{"CREATE", "START_INTERVAL"},
	})

	assert.Len(t, events, 2)
	assert.Equal(t, dt1, events[0].DateTime)
	assert.Equal(t, "test1", events[0].Description)
	assert.Equal(t, dt2, events[1].DateTime)
	assert.Equal(t, "test4", events[1].Description)
}

func TestArchiveContext(t *testing.T) {
	as := NewTestArchiveStore(test.TestId)
	cs := NewTestContextStore()
	tp := NewTestTimerProvider("2025-03-13 13:00:00")
	cm := New(cs, NewTestEventsStore(), as, tp)
	dt1, _ := time.Parse(time.DateTime, "2025-03-13 13:00:00")
	dt2, _ := time.Parse(time.DateTime, "2025-03-13 13:05:00")
	dt3, _ := time.Parse(time.DateTime, "2025-03-13 13:10:00")

	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
	tp.Current = dt2
	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
	tp.Current = dt3
	cm.Switch(test.TestId)
	cm.Switch(test.PrevTestId)
	cm.Archive(test.TestId)

	archive := as.Load(test.TestId)
	assert.Len(t, archive.Events, 6)
	assert.Equal(t, archive.Context.Id, test.TestId)
	assert.Equal(t, archive.Context.Description, test.TestDescription)
	assert.Equal(t, archive.Events[0].DateTime, dt1)
	assert.Equal(t, archive.Events[0].Type, ctx_model.CREATE_CTX)
	assert.Equal(t, archive.Events[1].DateTime, dt1)
	assert.Equal(t, archive.Events[1].Type, ctx_model.SWITCH_CTX)
	assert.Equal(t, archive.Events[2].DateTime, dt1)
	assert.Equal(t, archive.Events[2].Type, ctx_model.START_INTERVAL)
	assert.Equal(t, archive.Events[3].DateTime, dt3)
	assert.Equal(t, archive.Events[3].Type, ctx_model.END_INTERVAL)
	assert.Equal(t, archive.Events[4].DateTime, dt3)
	assert.Equal(t, archive.Events[4].Type, ctx_model.SWITCH_CTX)
	assert.Equal(t, archive.Events[5].DateTime, dt3)
	assert.Equal(t, archive.Events[5].Type, ctx_model.START_INTERVAL)
}

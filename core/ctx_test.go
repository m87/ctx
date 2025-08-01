package core_test

// import (
// 	"encoding/json"
// 	"errors"
// 	"log"
// 	"testing"
// 	"time"

// 	"github.com/m87/ctx/core"
// 	"github.com/m87/ctx/test"
// 	ctxtime "github.com/m87/ctx/time"
// 	"github.com/stretchr/testify/assert"
// )

// type TestTimeProvider struct {
// 	Current ctxtime.ZonedTime
// }

// func NewTestTimerProvider(dateTime string) *TestTimeProvider {
// 	dt, _ := time.ParseInLocation(time.DateTime, dateTime, time.UTC)
// 	return &TestTimeProvider{
// 		Current: ctxtime.ZonedTime{Time: dt, Timezone: time.UTC.String()},
// 	}
// }

// func (provider *TestTimeProvider) Now() ctxtime.ZonedTime {
// 	return provider.Current
// }

// type TestContextStore struct {
// 	store []byte
// }

// type TestTransactionalStore struct {
// 	store []byte
// }

// type TestEventsStore struct {
// 	store []byte
// }

// type TestArchiveStore struct {
// 	store       []byte
// 	storeEvents []byte
// }

// func (cs *TestContextStore) Load() core.State {
// 	state := core.State{}
// 	err := json.Unmarshal(cs.store, &state)
// 	if err != nil {
// 		log.Fatal("Unable to parse state store")
// 	}

// 	return state
// }

// func (cs *TestContextStore) Save(state *core.State) {
// 	data, err := json.Marshal(state)
// 	if err != nil {
// 		panic(err)
// 	}
// 	cs.store = data
// }

// func NewTestContextStore() *TestContextStore {
// 	tcs := TestContextStore{}
// 	tcs.Save(&core.State{
// 		Contexts:  map[string]core.Context{},
// 		CurrentId: "",
// 	})
// 	return &tcs
// }

// func (store *TestContextStore) Apply(fn core.StatePatch) error {
// 	state := store.Load()
// 	err := fn(&state)
// 	if err != nil {
// 		return err
// 	} else {
// 		store.Save(&state)
// 		return nil
// 	}
// }

// func (store *TestContextStore) Read(fn core.StatePatch) error {
// 	state := store.Load()
// 	return fn(&state)
// }

// func NewTestEventsStore() *TestEventsStore {
// 	tes := TestEventsStore{}
// 	tes.Save(&core.EventRegistry{
// 		Events: []core.Event{},
// 	})
// 	return &tes
// }

// func (es *TestEventsStore) Apply(fn core.EventsPatch) error {
// 	registry := es.Load()
// 	err := fn(&registry)
// 	if err != nil {
// 		return err
// 	} else {
// 		es.Save(&registry)
// 		return nil
// 	}
// }

// func (es *TestEventsStore) Read(fn core.EventsPatch) error {
// 	registry := es.Load()
// 	return fn(&registry)
// }

// func NewTestArchiveStore() *TestArchiveStore {
// 	tas := TestArchiveStore{}
// 	tas.Save(map[string]core.ContextArchive{})
// 	tas.SaveEvents(&core.EventsArchive{})
// 	return &tas
// }

// func (as *TestArchiveStore) Load() map[string]core.ContextArchive {
// 	state := map[string]core.ContextArchive{}
// 	err := json.Unmarshal(as.store, &state)
// 	if err != nil {
// 		log.Fatal("Unable to parse archvie store")
// 	}

// 	return state
// }

// func (as *TestArchiveStore) Save(archive map[string]core.ContextArchive) {
// 	data, err := json.Marshal(archive)
// 	if err != nil {
// 		panic(err)
// 	}
// 	as.store = data
// }

// func (as *TestArchiveStore) LoadEvents() *core.EventsArchive {
// 	events := core.EventsArchive{}
// 	err := json.Unmarshal(as.storeEvents, &events)
// 	if err != nil {
// 		log.Fatal("Unable to parse events archvie store")
// 	}

// 	return &events
// }

// func (as *TestArchiveStore) SaveEvents(entry *core.EventsArchive) {
// 	data, err := json.Marshal(entry)
// 	if err != nil {
// 		panic(err)
// 	}
// 	as.storeEvents = data
// }

// func (store *TestArchiveStore) Apply(id string, fn core.ArchivePatch) error {
// 	archive := store.Load()
// 	context := archive[id]
// 	if _, ok := archive[id]; !ok {
// 		context = core.ContextArchive{
// 			Context: core.Context{
// 				Id: id,
// 			},
// 		}
// 	}

// 	err := fn(&context)
// 	archive[id] = context
// 	if err != nil {
// 		return err
// 	} else {
// 		store.Save(archive)
// 		return nil
// 	}
// }

// func (store *TestArchiveStore) ApplyEvents(date string, fn core.ArchiveEventsPatch) error {
// 	archive := store.LoadEvents()
// 	err := fn(archive)
// 	if err != nil {
// 		return err
// 	} else {
// 		store.SaveEvents(archive)
// 		return nil
// 	}
// }

// func (es *TestEventsStore) Save(eventsRegistry *core.EventRegistry) {
// 	data, err := json.Marshal(eventsRegistry)
// 	if err != nil {
// 		panic(err)
// 	}
// 	es.store = data
// }

// func (es *TestEventsStore) Load() core.EventRegistry {
// 	eventsRegistry := core.EventRegistry{}
// 	err := json.Unmarshal(es.store, &eventsRegistry)
// 	if err != nil {
// 		log.Fatal("Unable to parse state store")
// 	}

// 	return eventsRegistry
// }

// func TestCreateContext(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())
// 	cm.CreateContext(test.TestId, test.TestDescription)

// 	state := cs.Load()
// 	assert.Len(t, state.Contexts, 1)
// 	createdContext := state.Contexts[test.TestId]
// 	assert.Equal(t, createdContext.Id, test.TestId)
// 	assert.Equal(t, createdContext.Description, test.TestDescription)
// 	assert.Equal(t, createdContext.State, core.ACTIVE)
// 	assert.Equal(t, createdContext.Duration, time.Duration(0))
// 	assert.Len(t, createdContext.Intervals, 0)
// 	assert.Len(t, createdContext.Comments, 0)
// }

// func TestCreateExistingId(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())
// 	cm.CreateContext(test.TestId, test.TestDescription)
// 	err := cm.CreateContext(test.TestId, test.TestDescription)

// 	assert.Error(t, err, errors.New("context already exists"))
// }

// func TestDontCreateContextWithEmptyDescription(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())
// 	err := cm.CreateContext(test.TestId, "  \t")

// 	assert.Error(t, err, errors.New("empty description"))
// }

// func TestDontCreateContextWithEmptyId(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())
// 	err := cm.CreateContext(" \t", test.TestDescription)

// 	assert.Error(t, err, errors.New("empty id"))
// }

// func TestEmitCreateEvent(t *testing.T) {
// 	dt1, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:00:00", time.UTC)
// 	es := NewTestEventsStore()
// 	cm := core.NewContextManager(NewTestContextStore(), es, NewTestArchiveStore(), NewTestTimerProvider("2025-03-13 13:00:00"))
// 	cm.CreateContext(test.TestId, test.TestDescription)

// 	registry := es.Load()
// 	assert.Len(t, registry.Events, 1)
// 	assert.Equal(t, registry.Events[0].Type, core.CREATE_CTX)
// 	assert.Equal(t, registry.Events[0].CtxId, test.TestId)
// 	assert.Equal(t, registry.Events[0].DateTime, ctxtime.ZonedTime{Time: dt1, Timezone: time.UTC.String()})
// }

// func TestSwitchContext(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())
// 	cm.CreateContext(test.TestId, test.TestDescription)

// 	err := cm.Switch(test.TestId)

// 	state := cs.Load()
// 	assert.NoError(t, err)
// 	assert.Equal(t, test.TestId, state.CurrentId)

// }

// func TestDontSwitchWithEmptyId(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())
// 	cm.CreateContext(test.TestId, test.TestDescription)

// 	err := cm.Switch("\t")

// 	state := cs.Load()
// 	assert.Error(t, err, errors.New("empty id"))
// 	assert.Equal(t, "", state.CurrentId)
// }

// func TestDontSwitchIfDoesNotExists(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())
// 	cm.CreateContext(test.TestId, test.TestDescription)
// 	cm.Switch(test.TestId)
// 	err := cm.Switch("test")

// 	state := cs.Load()
// 	assert.Error(t, err, errors.New("context does not exist"))
// 	assert.Equal(t, test.TestId, state.CurrentId)
// }

// func TestSwitchNotExistingContext(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())
// 	cm.CreateContext(test.PrevTestId, test.TestDescription)

// 	cm.Switch(test.PrevTestId)
// 	err := cm.Switch(test.TestId)

// 	state := cs.Load()
// 	assert.Error(t, err, errors.New("context does not exist"))
// 	assert.Equal(t, test.PrevTestId, state.CurrentId)

// }
// func TestSwitchCreateIfNotExists(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())

// 	err := cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)

// 	state := cs.Load()
// 	assert.NoError(t, err)
// 	assert.Equal(t, test.TestId, state.CurrentId)
// 	assert.NotNil(t, state.Contexts[test.TestId])
// 	assert.Len(t, state.Contexts, 1)
// 	createdContext := state.Contexts[test.TestId]
// 	assert.Equal(t, createdContext.Id, test.TestId)
// 	assert.Equal(t, createdContext.Description, test.TestDescription)
// 	assert.Equal(t, createdContext.State, core.ACTIVE)
// 	assert.Equal(t, createdContext.Duration, time.Duration(0))
// 	assert.Len(t, createdContext.Intervals, 1)
// 	assert.Len(t, createdContext.Comments, 0)

// }

// func TestDontSwitchOrCreateWithEmptyId(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())

// 	err := cm.CreateIfNotExistsAndSwitch("\t", test.TestDescription)

// 	state := cs.Load()
// 	assert.Error(t, err, errors.New("empty id"))
// 	assert.Equal(t, "", state.CurrentId)
// }

// func TestDontSwitchOrCreateWithEmptyDescription(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())

// 	err := cm.CreateIfNotExistsAndSwitch(test.TestId, " \t ")

// 	state := cs.Load()
// 	assert.Error(t, err, errors.New("empty id"))
// 	assert.Equal(t, "", state.CurrentId)
// }

// func TestSwitchCreateIfNotExistsOnExistingContext(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())
// 	cm.CreateContext(test.TestId, test.TestDescription)

// 	assert.Equal(t, cs.Load().CurrentId, "")
// 	err := cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)

// 	assert.Equal(t, cs.Load().CurrentId, test.TestId)
// 	assert.NoError(t, err)
// }

// func TestSwitchAlreadyActiveContext(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())
// 	cm.CreateContext(test.TestId, test.TestDescription)

// 	err := cm.Switch(test.TestId)
// 	assert.NoError(t, err)

// 	err = cm.Switch(test.TestId)
// 	state := cs.Load()
// 	assert.Error(t, err, errors.New("context already active"))
// 	assert.Len(t, state.Contexts[test.TestId].Intervals, 1)

// 	err = cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	state = cs.Load()
// 	assert.Error(t, err, errors.New("context already active"))
// 	assert.Len(t, state.Contexts[test.TestId].Intervals, 1)

// }

// func TestIntervals(t *testing.T) {
// 	cs := NewTestContextStore()
// 	dt1, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:00:00", time.UTC)
// 	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:05:00", time.UTC)
// 	dt3, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:10:00", time.UTC)

// 	tp := NewTestTimerProvider("2025-03-13 13:00:00")
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), tp)
// 	cm.CreateContext(test.TestId, test.TestDescription)

// 	tp.Current = ctxtime.ZonedTime{Time: dt1, Timezone: time.UTC.String()}
// 	cm.Switch(test.TestId)
// 	state := cs.Load()
// 	assert.Equal(t, test.TestId, state.CurrentId)
// 	prevCtx := state.Contexts[state.CurrentId]
// 	assert.Len(t, prevCtx.Intervals, 1)
// 	assert.Equal(t, prevCtx.Intervals[0].Start, tp.Current)
// 	assert.True(t, prevCtx.Intervals[0].End.Time.IsZero())

// 	tp.Current = ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()}
// 	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
// 	state = cs.Load()
// 	prevCtx = state.Contexts[test.TestId]
// 	assert.Equal(t, prevCtx.Intervals[0].Start, ctxtime.ZonedTime{Time: dt1, Timezone: time.UTC.String()})
// 	assert.Equal(t, prevCtx.Intervals[0].End, ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()})
// 	assert.Equal(t, test.PrevTestId, state.CurrentId)
// 	nextCtx := state.Contexts[state.CurrentId]
// 	assert.Len(t, nextCtx.Intervals, 1)
// 	assert.Equal(t, nextCtx.Intervals[0].Start, ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()})
// 	assert.True(t, nextCtx.Intervals[0].End.Time.IsZero())

// 	tp.Current = ctxtime.ZonedTime{Time: dt3, Timezone: time.UTC.String()}
// 	cm.Switch(test.TestId)
// 	state = cs.Load()
// 	nextCtx = state.Contexts[test.PrevTestId]
// 	assert.Equal(t, nextCtx.Intervals[0].Start, ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()})
// 	assert.Equal(t, nextCtx.Intervals[0].End, ctxtime.ZonedTime{Time: dt3, Timezone: time.UTC.String()})
// 	assert.Equal(t, test.TestId, state.CurrentId)
// 	prevCtx = state.Contexts[state.CurrentId]
// 	assert.Len(t, prevCtx.Intervals, 2)
// 	assert.Equal(t, prevCtx.Intervals[1].Start, ctxtime.ZonedTime{Time: dt3, Timezone: time.UTC.String()})
// 	assert.True(t, prevCtx.Intervals[1].End.Time.IsZero())

// }

// func TestEventsFlow(t *testing.T) {
// 	es := NewTestEventsStore()
// 	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:05:00", time.UTC)
// 	dt3, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:10:00", time.UTC)

// 	tp := NewTestTimerProvider("2025-03-13 13:00:00")
// 	cm := core.NewContextManager(NewTestContextStore(), es, NewTestArchiveStore(), tp)
// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.PrevDescription)
// 	tp.Current = ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()}
// 	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
// 	tp.Current = ctxtime.ZonedTime{Time: dt3, Timezone: time.UTC.String()}
// 	cm.Switch(test.TestId)

// 	registry := es.Load()
// 	assert.Len(t, registry.Events, 10)
// 	assert.Equal(t, registry.Events[0].Type, core.CREATE_CTX)
// 	assert.Equal(t, registry.Events[1].Type, core.SWITCH_CTX)
// 	assert.Equal(t, registry.Events[2].Type, core.START_INTERVAL)
// 	assert.Equal(t, registry.Events[3].Type, core.CREATE_CTX)
// 	assert.Equal(t, registry.Events[4].Type, core.END_INTERVAL)
// 	assert.Equal(t, registry.Events[5].Type, core.SWITCH_CTX)
// 	assert.Equal(t, registry.Events[6].Type, core.START_INTERVAL)
// 	assert.Equal(t, registry.Events[7].Type, core.END_INTERVAL)
// 	assert.Equal(t, registry.Events[8].Type, core.SWITCH_CTX)
// 	assert.Equal(t, registry.Events[9].Type, core.START_INTERVAL)

// }

// func TestFree(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), NewTestTimerProvider("2025-03-13 13:00:00"))
// 	dt, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:00:00", time.UTC)

// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	assert.Equal(t, test.TestId, cs.Load().CurrentId)

// 	cm.Free()
// 	state := cs.Load()
// 	assert.Equal(t, "", state.CurrentId)
// 	assert.Equal(t, ctxtime.ZonedTime{Time: dt, Timezone: time.UTC.String()}, state.Contexts[test.TestId].Intervals[0].Start)
// 	assert.Equal(t, ctxtime.ZonedTime{Time: dt, Timezone: time.UTC.String()}, state.Contexts[test.TestId].Intervals[0].End)

// }

// func TestFreeWithNowCurrentContext(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())
// 	cm.CreateContext(test.TestId, test.TestDescription)
// 	assert.Equal(t, "", cs.Load().CurrentId)
// 	err := cm.Free()
// 	assert.Error(t, err, errors.New("no active context"))
// }

// func TestDeleteContext(t *testing.T) {
// 	cs := NewTestContextStore()
// 	es := NewTestEventsStore()
// 	cm := core.NewContextManager(cs, es, NewTestArchiveStore(), ctxtime.NewTimer())
// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	cm.Free()
// 	assert.Len(t, cs.Load().Contexts, 1)
// 	assert.Len(t, es.Load().Events, 4)
// 	err := cm.Delete(test.TestId)
// 	assert.NoError(t, err)
// 	assert.Len(t, es.Load().Events, 5)
// 	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Type, core.DELETE_CTX)
// 	assert.Len(t, cs.Load().Contexts, 0)
// }

// func TestDeleteActiveContext(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())
// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	assert.Len(t, cs.Load().Contexts, 1)
// 	err := cm.Delete(test.TestId)
// 	assert.Error(t, err, errors.New("context is active"))
// 	assert.Len(t, cs.Load().Contexts, 1)
// }

// func TestDeleteNotExistingContext(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())
// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	err := cm.Delete(test.PrevTestId)
// 	assert.Error(t, err, errors.New("context does not exist"))
// 	assert.Len(t, cs.Load().Contexts, 1)
// }

// // func TestEventFilter(t *testing.T) {
// // 	es := NewTestEventsStore()
// // 	cm := core.New(NewTestContextStore(), es, NewTestArchiveStore(), ctxtime.NewTimer())
// // 	dt1, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:00:00", time.UTC)
// // 	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-14 13:00:00", time.UTC)
// // 	cm.EventsStore.Apply(func(er *core.EventRegistry) error {
// // 		er.Events = append(er.Events, core.Event{
// // 			DateTime: ctxtime.ZonedTime{Time: dt1}, Description: "test1", Type: core.CREATE_CTX,
// // 		})
// // 		er.Events = append(er.Events, core.Event{
// // 			DateTime: ctxtime.ZonedTime{Time: dt2}, Description: "test2", Type: core.SWITCH_CTX,
// // 		})
// // 		er.Events = append(er.Events, core.Event{
// // 			DateTime: ctxtime.ZonedTime{Time: dt1}, Description: "test3", Type: core.SWITCH_CTX,
// // 		})
// // 		er.Events = append(er.Events, core.Event{
// // 			DateTime: ctxtime.ZonedTime{Time: dt2}, Description: "test4", Type: core.START_INTERVAL,
// // 		})
// // 		return nil
// // 	})

// // 	er := es.Load()
// // 	assert.Len(t, er.Events, 4)
// // 	events := cm.filterEvents(&er, core.EventsFilter{
// // 		Date: "2025-03-14",
// // 	})

// // 	assert.Len(t, events, 2)
// // 	assert.Equal(t, ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()}, events[0].DateTime)
// // 	assert.Equal(t, "test2", events[0].Description)
// // 	assert.Equal(t, ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()}, events[1].DateTime)
// // 	assert.Equal(t, "test4", events[1].Description)

// // 	events = cm.filterEvents(&er, core.EventsFilter{
// // 		Date:  "2025-03-13",
// // 		Types: []string{"CREATE"},
// // 	})

// // 	assert.Len(t, events, 1)
// // 	assert.Equal(t, ctxtime.ZonedTime{Time: dt1, Timezone: time.UTC.String()}, events[0].DateTime)
// // 	assert.Equal(t, "test1", events[0].Description)

// // 	events = cm.filterEvents(&er, core.EventsFilter{
// // 		Types: []string{"SWITCH"},
// // 	})

// // 	assert.Len(t, events, 2)
// // 	assert.Equal(t, ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()}, events[0].DateTime)
// // 	assert.Equal(t, "test2", events[0].Description)
// // 	assert.Equal(t, ctxtime.ZonedTime{Time: dt1, Timezone: time.UTC.String()}, events[1].DateTime)
// // 	assert.Equal(t, "test3", events[1].Description)

// // 	events = cm.filterEvents(&er, core.EventsFilter{
// // 		Types: []string{"CREATE", "START_INTERVAL"},
// // 	})

// // 	assert.Len(t, events, 2)
// // 	assert.Equal(t, ctxtime.ZonedTime{Time: dt1, Timezone: time.UTC.String()}, events[0].DateTime)
// // 	assert.Equal(t, "test1", events[0].Description)
// // 	assert.Equal(t, ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()}, events[1].DateTime)
// // 	assert.Equal(t, "test4", events[1].Description)
// // }

// func TestArchiveContext(t *testing.T) {
// 	as := NewTestArchiveStore()
// 	cs := NewTestContextStore()
// 	tp := NewTestTimerProvider("2025-03-13 13:00:00")
// 	es := NewTestEventsStore()
// 	cm := core.NewContextManager(cs, es, as, tp)
// 	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:05:00", time.UTC)
// 	dt3, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:10:00", time.UTC)

// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	tp.Current = ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()}
// 	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
// 	tp.Current = ctxtime.ZonedTime{Time: dt3, Timezone: time.UTC.String()}
// 	cm.Switch(test.TestId)
// 	cm.Switch(test.PrevTestId)
// 	state := cs.Load()
// 	assert.Len(t, state.Contexts, 2)
// 	assert.Len(t, es.Load().Events, 13)
// 	err := cm.Archive(test.TestId)

// 	assert.NoError(t, err)
// 	archive := as.Load()[test.TestId]
// 	state = cs.Load()
// 	assert.Equal(t, archive.Context.Id, test.TestId)
// 	assert.Equal(t, archive.Context.Description, test.TestDescription)
// 	assert.Len(t, state.Contexts, 1)
// 	assert.Len(t, es.Load().Events, 14)
// }

// func TestDontArchiveActiveContext(t *testing.T) {
// 	as := NewTestArchiveStore()
// 	cs := NewTestContextStore()
// 	tp := NewTestTimerProvider("2025-03-13 13:00:00")
// 	es := NewTestEventsStore()
// 	cm := core.NewContextManager(cs, es, as, tp)
// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	err := cm.Archive(test.TestId)

// 	assert.Error(t, err, errors.New("context is active"))
// }

// func TestArchiveAll(t *testing.T) {
// 	as := NewTestArchiveStore()
// 	cs := NewTestContextStore()
// 	tp := NewTestTimerProvider("2025-03-13 13:00:00")
// 	es := NewTestEventsStore()
// 	cm := core.NewContextManager(cs, es, as, tp)
// 	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:05:00", time.UTC)
// 	dt3, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:10:00", time.UTC)

// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	tp.Current = ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()}
// 	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
// 	tp.Current = ctxtime.ZonedTime{Time: dt3, Timezone: time.UTC.String()}
// 	cm.Switch(test.TestId)
// 	cm.Switch(test.PrevTestId)
// 	cm.Free()
// 	assert.Len(t, es.Load().Events, 14)
// 	assert.Len(t, cs.Load().Contexts, 2)
// 	err := cm.ArchiveAll()

// 	assert.NoError(t, err)
// 	events := as.LoadEvents().Events
// 	assert.Len(t, events, 16)
// 	assert.Len(t, cs.Load().Contexts, 0)
// }

// func TestMergeContexts(t *testing.T) {
// 	as := NewTestArchiveStore()
// 	cs := NewTestContextStore()
// 	tp := NewTestTimerProvider("2025-03-13 13:00:00")
// 	es := NewTestEventsStore()
// 	cm := core.NewContextManager(cs, es, as, tp)
// 	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:05:00", time.UTC)
// 	dt3, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:10:00", time.UTC)

// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	cm.CreateIfNotExistsAndSwitch(test.TestIdExtra, test.PrevDescription)
// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	tp.Current = ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()}
// 	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
// 	tp.Current = ctxtime.ZonedTime{Time: dt3, Timezone: time.UTC.String()}
// 	cm.Switch(test.TestId)
// 	cm.Switch(test.PrevTestId)
// 	cm.Free()
// 	assert.Len(t, es.Load().Events, 21)
// 	assert.Len(t, cs.Load().Contexts, 3)
// 	assert.Equal(t, cs.Load().Contexts[test.PrevTestId].Duration, dt3.Sub(dt2))

// 	err := cm.MergeContext(test.TestId, test.PrevTestId)

// 	assert.NoError(t, err)
// 	assert.Equal(t, cs.Load().Contexts[test.PrevTestId].Duration, dt3.Sub(dt2)*2)
// 	events := es.Load().Events
// 	assert.Len(t, events, 23)
// 	assert.Len(t, cs.Load().Contexts, 2)
// 	for _, event := range events {
// 		if event.Type == core.SWITCH_CTX {
// 			if v, ok := event.Data["from"]; ok && v != "" && event.CtxId == test.TestIdExtra {
// 				assert.Equal(t, v, test.TestId)
// 			}
// 		}
// 	}
// }

// func TestArchiveAllEvents(t *testing.T) {
// 	as := NewTestArchiveStore()
// 	cs := NewTestContextStore()
// 	tp := NewTestTimerProvider("2025-03-13 13:00:00")
// 	es := NewTestEventsStore()
// 	cm := core.NewContextManager(cs, es, as, tp)
// 	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:05:00", time.UTC)
// 	dt3, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:10:00", time.UTC)

// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	tp.Current = ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()}
// 	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
// 	tp.Current = ctxtime.ZonedTime{Time: dt3, Timezone: time.UTC.String()}
// 	cm.Switch(test.TestId)
// 	cm.Switch(test.PrevTestId)
// 	cm.Free()
// 	assert.Len(t, es.Load().Events, 14)
// 	assert.Len(t, cs.Load().Contexts, 2)
// 	err := cm.ArchiveAllEvents()

// 	assert.NoError(t, err)
// 	events := as.LoadEvents().Events
// 	assert.Len(t, events, 14)
// 	assert.Len(t, es.Load().Events, 0)
// }

// func TestErrorOnEditCurrentContextInterval(t *testing.T) {
// 	as := NewTestArchiveStore()
// 	cs := NewTestContextStore()
// 	tp := NewTestTimerProvider("2025-03-13 13:00:00")
// 	es := NewTestEventsStore()
// 	cm := core.NewContextManager(cs, es, as, tp)

// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)

// 	id := cs.Load().Contexts[test.TestId].Intervals[0].Id
// 	err := cm.EditContextInterval(test.TestId, id, ctxtime.ZonedTime{Time: time.Now().Local()}, ctxtime.ZonedTime{Time: time.Now().Local()})

// 	assert.Error(t, errors.New("context is active"), err)

// }

// func TestEditContextInterval(t *testing.T) {
// 	as := NewTestArchiveStore()
// 	cs := NewTestContextStore()
// 	tp := NewTestTimerProvider("2025-03-13 13:00:00")
// 	es := NewTestEventsStore()
// 	cm := core.NewContextManager(cs, es, as, tp)
// 	dt1, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:00:00", time.UTC)
// 	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:05:00", time.UTC)
// 	dt3, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:10:00", time.UTC)

// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	tp.Current = ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()}
// 	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.TestDescription)
// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	tp.Current = ctxtime.ZonedTime{Time: dt3, Timezone: time.UTC.String()}
// 	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.TestDescription)

// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[0].Start, ctxtime.ZonedTime{Time: dt1, Timezone: time.UTC.String()})
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[0].End, ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()})
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[0].Duration, dt2.Sub(dt1))
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[1].Start, ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()})
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[1].End, ctxtime.ZonedTime{Time: dt3, Timezone: time.UTC.String()})
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[1].Duration, dt3.Sub(dt2))
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Duration, cs.Load().Contexts[test.TestId].Intervals[0].Duration+cs.Load().Contexts[test.TestId].Intervals[1].Duration)

// 	id := cs.Load().Contexts[test.TestId].Intervals[0].Id
// 	err := cm.EditContextInterval(test.TestId, id, ctxtime.ZonedTime{Time: dt1}, ctxtime.ZonedTime{Time: dt3, Timezone: time.UTC.String()})
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[0].Start, ctxtime.ZonedTime{Time: dt1, Timezone: time.UTC.String()})
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[0].End, ctxtime.ZonedTime{Time: dt3, Timezone: time.UTC.String()})
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[0].Duration, dt3.Sub(dt1))
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[1].Start, ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()})
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[1].End, ctxtime.ZonedTime{Time: dt3, Timezone: time.UTC.String()})
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[1].Duration, dt3.Sub(dt2))
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Duration, cs.Load().Contexts[test.TestId].Intervals[0].Duration+cs.Load().Contexts[test.TestId].Intervals[1].Duration)
// 	assert.NoError(t, err, errors.New("context is active"))
// 	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Type, core.EDIT_CTX_INTERVAL)
// 	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Data["old.start"], dt1.Format(time.RFC3339))
// 	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Data["old.end"], dt2.Format(time.RFC3339))
// 	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Data["new.start"], dt1.Format(time.RFC3339))
// 	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Data["new.end"], dt3.Format(time.RFC3339))

// 	id = cs.Load().Contexts[test.TestId].Intervals[0].Id
// 	err = cm.EditContextInterval(test.TestId, id, ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()}, ctxtime.ZonedTime{Time: dt3, Timezone: time.UTC.String()})
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[0].Start, ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()})
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[0].End, ctxtime.ZonedTime{Time: dt3, Timezone: time.UTC.String()})
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[0].Duration, dt3.Sub(dt2))
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[1].Start, ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()})
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[1].End, ctxtime.ZonedTime{Time: dt3, Timezone: time.UTC.String()})
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[1].Duration, dt3.Sub(dt2))
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Duration, cs.Load().Contexts[test.TestId].Intervals[0].Duration+cs.Load().Contexts[test.TestId].Intervals[1].Duration)
// 	assert.NoError(t, err, errors.New("context is active"))
// 	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Type, core.EDIT_CTX_INTERVAL)
// 	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Data["old.start"], dt1.Format(time.RFC3339))
// 	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Data["old.end"], dt3.Format(time.RFC3339))
// 	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Data["new.start"], dt2.Format(time.RFC3339))
// 	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Data["new.end"], dt3.Format(time.RFC3339))
// }

// func TestRename(t *testing.T) {
// 	as := NewTestArchiveStore()
// 	cs := NewTestContextStore()
// 	tp := NewTestTimerProvider("2025-03-13 13:00:00")
// 	es := NewTestEventsStore()
// 	cm := core.NewContextManager(cs, es, as, tp)

// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)

// 	cm.RenameContext(test.TestId, test.PrevTestId, test.PrevDescription)

// 	state := cs.Load()
// 	assert.Contains(t, state.Contexts, test.PrevTestId)
// 	assert.NotContains(t, state.Contexts, test.TestId)
// 	assert.Len(t, state.Contexts[test.PrevTestId].Intervals, 1)
// 	assert.Equal(t, state.Contexts[test.PrevTestId].Description, test.PrevDescription)

// }

// func TestGetIntervalDurationForDateInBetween(t *testing.T) {
// 	cs := NewTestContextStore()
// 	tp := NewTestTimerProvider("2025-03-13 13:00:00")
// 	es := NewTestEventsStore()
// 	cm := core.NewContextManager(cs, es, NewTestArchiveStore(), tp)
// 	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-15 13:05:00", time.UTC)
// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	tp.Current = ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()}
// 	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
// 	state := cs.Load()
// 	date, _ := time.ParseInLocation(time.DateOnly, "2025-03-14", time.UTC)
// 	duration, err := cm.GetIntervalDurationsByDate(&state, test.TestId, ctxtime.ZonedTime{Time: date, Timezone: time.UTC.String()})
// 	assert.NoError(t, err)
// 	assert.Equal(t, 24*time.Hour, duration)
// }

// func TestGetIntervalDurationForDateOutOfBounds(t *testing.T) {
// 	cs := NewTestContextStore()
// 	tp := NewTestTimerProvider("2025-03-13 13:00:00")
// 	es := NewTestEventsStore()
// 	cm := core.NewContextManager(cs, es, NewTestArchiveStore(), tp)
// 	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-15 13:05:00", time.UTC)
// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	tp.Current = ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()}
// 	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
// 	state := cs.Load()
// 	date, _ := time.ParseInLocation(time.DateOnly, "2025-03-16", time.UTC)
// 	duration, err := cm.GetIntervalDurationsByDate(&state, test.TestId, ctxtime.ZonedTime{Time: date, Timezone: time.UTC.String()})
// 	assert.NoError(t, err)
// 	assert.Equal(t, time.Duration(0), duration)
// 	date, _ = time.ParseInLocation(time.DateOnly, "2025-03-12", time.UTC)
// 	duration, err = cm.GetIntervalDurationsByDate(&state, test.TestId, ctxtime.ZonedTime{Time: date, Timezone: time.UTC.String()})
// 	assert.NoError(t, err)
// 	assert.Equal(t, time.Duration(0), duration)
// }

// func TestGetIntervalDurationForDateBefore(t *testing.T) {
// 	cs := NewTestContextStore()
// 	tp := NewTestTimerProvider("2025-03-13 10:00:00")
// 	es := NewTestEventsStore()
// 	cm := core.NewContextManager(cs, es, NewTestArchiveStore(), tp)
// 	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-15 13:00:00", time.UTC)
// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	tp.Current = ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()}
// 	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
// 	state := cs.Load()
// 	date, _ := time.ParseInLocation(time.DateOnly, "2025-03-15", time.UTC)
// 	duration, err := cm.GetIntervalDurationsByDate(&state, test.TestId, ctxtime.ZonedTime{Time: date, Timezone: time.UTC.String()})
// 	assert.NoError(t, err)
// 	assert.Equal(t, 13*time.Hour, duration)
// }

// func TestGetIntervalDurationForDateAfter(t *testing.T) {
// 	cs := NewTestContextStore()
// 	tp := NewTestTimerProvider("2025-03-13 10:00:00")
// 	es := NewTestEventsStore()
// 	cm := core.NewContextManager(cs, es, NewTestArchiveStore(), tp)
// 	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-15 13:00:00", time.UTC)
// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	tp.Current = ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()}
// 	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
// 	state := cs.Load()
// 	date, _ := time.ParseInLocation(time.DateOnly, "2025-03-13", time.UTC)
// 	duration, err := cm.GetIntervalDurationsByDate(&state, test.TestId, ctxtime.ZonedTime{Time: date, Timezone: time.UTC.String()})
// 	assert.NoError(t, err)
// 	assert.Equal(t, 14*time.Hour, duration)
// }

// func TestGetIntervalDurationForDateEqual(t *testing.T) {
// 	cs := NewTestContextStore()
// 	tp := NewTestTimerProvider("2025-03-13 10:00:00")
// 	es := NewTestEventsStore()
// 	cm := core.NewContextManager(cs, es, NewTestArchiveStore(), tp)
// 	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:00:00", time.UTC)
// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	tp.Current = ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()}
// 	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
// 	state := cs.Load()
// 	date, _ := time.ParseInLocation(time.DateOnly, "2025-03-13", time.UTC)
// 	duration, err := cm.GetIntervalDurationsByDate(&state, test.TestId, ctxtime.ZonedTime{Time: date, Timezone: time.UTC.String()})
// 	assert.NoError(t, err)
// 	assert.Equal(t, 3*time.Hour, duration)
// }

// func TestDeleteINtervalActiveContext(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())
// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	err := cm.DeleteIntervalByIndex(test.TestId, 0)

// 	assert.Error(t, err, errors.New("context is active"))
// }

// func TestDeleteIntervalNotExistingContext(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())
// 	err := cm.DeleteIntervalByIndex(test.TestId, 0)

// 	assert.Error(t, err, errors.New("context does not exist"))
// }

// func TestDeleteIntervalOutOfBounds(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())
// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	err := cm.DeleteIntervalByIndex(test.TestId, 0)

// 	assert.Error(t, err, errors.New("interval out of bounds"))
// }

// func TestDeleteInterval(t *testing.T) {
// 	cs := NewTestContextStore()
// 	tp := NewTestTimerProvider("2025-03-13 13:00:00")
// 	es := NewTestEventsStore()
// 	cm := core.NewContextManager(cs, es, NewTestArchiveStore(), tp)
// 	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:05:00", time.UTC)
// 	dt3, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:10:00", time.UTC)
// 	dt4, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:15:00", time.UTC)

// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	tp.Current = ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()}
// 	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
// 	tp.Current = ctxtime.ZonedTime{Time: dt3, Timezone: time.UTC.String()}
// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	tp.Current = ctxtime.ZonedTime{Time: dt4, Timezone: time.UTC.String()}
// 	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)

// 	assert.Len(t, cs.Load().Contexts[test.TestId].Intervals, 2)
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Duration.Seconds(), time.Duration(600000000000).Seconds())

// 	err := cm.DeleteIntervalByIndex(test.TestId, 0)

// 	assert.NoError(t, err)
// 	assert.Len(t, cs.Load().Contexts[test.TestId].Intervals, 1)
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[0].Start, ctxtime.ZonedTime{Time: dt3, Timezone: time.UTC.String()})
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[0].End, ctxtime.ZonedTime{Time: dt4, Timezone: time.UTC.String()})
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Duration.Seconds(), time.Duration(300000000000).Seconds())
// 	// assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Type, core.DELETE_CTX_INTERVAL)
// }

// func TestSearchContextWithRegex(t *testing.T) {
// 	cs := NewTestContextStore()
// 	tp := NewTestTimerProvider("2025-03-13 13:00:00")
// 	es := NewTestEventsStore()
// 	cm := core.NewContextManager(cs, es, NewTestArchiveStore(), tp)
// 	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:05:00", time.UTC)
// 	dt3, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:10:00", time.UTC)

// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	tp.Current = ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()}
// 	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
// 	tp.Current = ctxtime.ZonedTime{Time: dt3, Timezone: time.UTC.String()}
// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)

// 	assert.Len(t, cs.Load().Contexts[test.TestId].Intervals, 2)

// 	contexts, err := cm.Search("p.*test.*")
// 	assert.NoError(t, err)
// 	assert.Len(t, contexts, 1)
// 	assert.Equal(t, contexts[0].Description, test.PrevDescription)
// }

// func TestAddLabelToContext(t *testing.T) {
// 	cs := NewTestContextStore()
// 	tp := NewTestTimerProvider("2025-03-13 13:00:00")
// 	es := NewTestEventsStore()
// 	cm := core.NewContextManager(cs, es, NewTestArchiveStore(), tp)

// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	err := cm.LabelContext(test.TestId, "test-label")

// 	assert.NoError(t, err)
// 	assert.Contains(t, cs.Load().Contexts[test.TestId].Labels, "test-label")
// 	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Type, core.LABEL_CTX)
// }

// func TestRemoveLabelFromContext(t *testing.T) {
// 	cs := NewTestContextStore()
// 	tp := NewTestTimerProvider("2025-03-13 13:00:00")
// 	es := NewTestEventsStore()
// 	cm := core.NewContextManager(cs, es, NewTestArchiveStore(), tp)

// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	err := cm.LabelContext(test.TestId, "test-label")
// 	assert.NoError(t, err)

// 	err = cm.DeleteLabelContext(test.TestId, "test-label")
// 	assert.NoError(t, err)
// 	assert.NotContains(t, cs.Load().Contexts[test.TestId].Labels, "test-label")
// 	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Type, core.DELETE_CTX_LABEL)
// }

// func TestMoveInterval(t *testing.T) {
// 	cs := NewTestContextStore()
// 	tp := NewTestTimerProvider("2025-03-13 13:00:00")
// 	es := NewTestEventsStore()
// 	cm := core.NewContextManager(cs, es, NewTestArchiveStore(), tp)

// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
// 	cm.Free()

// 	assert.Len(t, cs.Load().Contexts[test.TestId].Intervals, 1)
// 	assert.Len(t, cs.Load().Contexts[test.PrevTestId].Intervals, 1)

// 	err := cm.MoveIntervalByIndex(test.TestId, test.PrevTestId, 0)
// 	assert.NoError(t, err)

// 	assert.Len(t, cs.Load().Contexts[test.TestId].Intervals, 0)
// 	assert.Len(t, cs.Load().Contexts[test.PrevTestId].Intervals, 2)
// }

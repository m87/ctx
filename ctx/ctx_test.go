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

func (store *TestContextStore) Read(fn ctx_model.StatePatch) {
	state := store.Load()
	fn(&state)
}

func TestCreateContext(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTimer())
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
	cm := New(cs, NewTimer())
	cm.CreateContext(test.TestId, test.TestDescription)
	err := cm.CreateContext(test.TestId, test.TestDescription)

	assert.Error(t, errors.New("Context already exists"), err)
}

func TestSwitchContext(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTimer())
	cm.CreateContext(test.TestId, test.TestDescription)

	err := cm.Switch(test.TestId)

	state := cs.Load()
	assert.NoError(t, err)
	assert.Equal(t, test.TestId, state.CurrentId)

}

func TestSwitchNotExistingContext(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTimer())
	cm.CreateContext(test.PrevTestId, test.TestDescription)

	cm.Switch(test.PrevTestId)
	err := cm.Switch(test.TestId)

	state := cs.Load()
	assert.Error(t, errors.New("Context does not exist"), err)
	assert.Equal(t, test.PrevTestId, state.CurrentId)

}
func TestSwitchCreateIfNotExists(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTimer())

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
func TestSwitchCreateIfNotExistsOnExistingContext(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTimer())
	cm.CreateContext(test.TestId, test.TestDescription)

	err := cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)

	assert.Error(t, errors.New("Context already exists"), err)

}

func TestSwitchAlreadyActiveContext(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTimer())
	cm.CreateContext(test.TestId, test.TestDescription)

	err := cm.Switch(test.TestId)
	assert.NoError(t, err)

	err = cm.Switch(test.TestId)
	state := cs.Load()
	assert.Error(t, errors.New("Context already active"), err)
	assert.Len(t, state.Contexts[test.TestId].Intervals, 1)

	err = cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
	state = cs.Load()
	assert.Error(t, errors.New("Context already active"), err)
	assert.Len(t, state.Contexts[test.TestId].Intervals, 1)

}

func TestIntervals(t *testing.T) {
	cs := NewTestContextStore()
	dt1, _ := time.Parse(time.DateTime, "2025-03-13 13:00:00")
	dt2, _ := time.Parse(time.DateTime, "2025-03-13 13:05:00")
	dt3, _ := time.Parse(time.DateTime, "2025-03-13 13:10:00")

	tp := NewTestTimerProvider("2025-03-13 13:00:00")
	cm := New(cs, tp)
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
	assert.Equal(t, nextCtx.Intervals[0].Start, dt2)
	assert.Equal(t, nextCtx.Intervals[0].End, dt3)
	assert.Equal(t, test.TestId, state.CurrentId)
	prevCtx = state.Contexts[state.CurrentId]
	assert.Len(t, prevCtx.Intervals, 2)
	assert.Equal(t, prevCtx.Intervals[1].Start, dt3)
	assert.True(t, prevCtx.Intervals[1].End.IsZero())

}

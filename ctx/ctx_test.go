package ctx

import (
	"errors"
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
	state ctx_model.State
}

func NewTestContextStore() *TestContextStore {
	return &TestContextStore{
		state: ctx_model.State{
			Contexts:  map[string]ctx_model.Context{},
			CurrentId: "",
		},
	}
}

func (store *TestContextStore) Apply(fn ctx_model.StatePatch) error {
	return fn(&store.state)
}

func (store *TestContextStore) Read(fn ctx_model.StatePatch) {
	fn(&store.state)
}

func TestCreateContext(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTimer())
	cm.CreateContext(test.TestId, test.TestDescription)

	assert.Len(t, cs.state.Contexts, 1)
	createdContext := cs.state.Contexts[test.TestId]
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

	assert.NoError(t, err)
	assert.Equal(t, test.TestId, cs.state.CurrentId)

}

func TestSwitchNotExistingContext(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTimer())
	cm.CreateContext(test.PrevTestId, test.TestDescription)

	cm.Switch(test.PrevTestId)
	err := cm.Switch(test.TestId)

	assert.Error(t, errors.New("Context does not exist"), err)
	assert.Equal(t, test.PrevTestId, cs.state.CurrentId)

}
func TestSwitctCreateIfNotExists(t *testing.T) {
	cs := NewTestContextStore()
	cm := New(cs, NewTimer())

	err := cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)

	assert.NoError(t, err)
	assert.Equal(t, test.TestId, cs.state.CurrentId)
	assert.NotNil(t, cs.state.Contexts[test.TestId])
	assert.Len(t, cs.state.Contexts, 1)
	createdContext := cs.state.Contexts[test.TestId]
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
	assert.Error(t, errors.New("Context already active"), err)
	assert.Len(t, cs.state.Contexts[test.TestId].Intervals, 1)

	err = cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
	assert.Error(t, errors.New("Context already active"), err)
	assert.Len(t, cs.state.Contexts[test.TestId].Intervals, 1)

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
	assert.Equal(t, test.TestId, cs.state.CurrentId)
	prevCtx := cs.state.Contexts[cs.state.CurrentId]
	assert.Len(t, prevCtx.Intervals, 1)
	assert.Equal(t, prevCtx.Intervals[0].Start, tp.Current)
	assert.True(t, prevCtx.Intervals[0].End.IsZero())

	tp.Current = dt2
	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
	assert.Equal(t, prevCtx.Intervals[0].Start, dt1)
	assert.Equal(t, prevCtx.Intervals[0].End, dt2)
	assert.Equal(t, test.PrevTestId, cs.state.CurrentId)
	nextCtx := cs.state.Contexts[cs.state.CurrentId]
	assert.Len(t, nextCtx.Intervals, 1)
	assert.Equal(t, nextCtx.Intervals[0].Start, dt2)
	assert.True(t, nextCtx.Intervals[0].End.IsZero())

	tp.Current = dt3
	cm.Switch(test.TestId)
	assert.Equal(t, nextCtx.Intervals[0].Start, dt2)
	assert.Equal(t, nextCtx.Intervals[0].End, dt3)
	assert.Equal(t, test.TestId, cs.state.CurrentId)
	prevCtx = cs.state.Contexts[cs.state.CurrentId]
	assert.Len(t, prevCtx.Intervals, 2)
	assert.Equal(t, prevCtx.Intervals[1].Start, dt3)
	assert.True(t, prevCtx.Intervals[1].End.IsZero())

}

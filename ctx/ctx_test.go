package ctx

import (
	"errors"
	"testing"
	"time"

	"github.com/m87/ctx/ctx_model"
	"github.com/m87/ctx/test"
	"github.com/stretchr/testify/assert"
)

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

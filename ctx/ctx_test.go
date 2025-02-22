package ctx

import (
	"testing"

	"github.com/m87/ctx/assert"
	"github.com/m87/ctx/ctx_model"
)

func TestCreateContext(t *testing.T) {

	state := ctx_model.State{
		Contexts: map[string]ctx_model.Context{},
	}
	Create(&state, "123", "test")
	assert.Equal(t, len(state.Contexts), 1)
	createdContext := state.Contexts["123"]
	assert.Equal(t, createdContext.Id, "123")
	assert.Equal(t, createdContext.Description, "test")
	assert.Equal(t, createdContext.State, ctx_model.ACTIVE)
	assert.Equal(t, createdContext.Duration, 0)
	assert.Size(t, createdContext.Intervals, 0)
	assert.Size(t, createdContext.Comments, 0)
}

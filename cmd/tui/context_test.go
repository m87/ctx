package tui

import (
	"testing"

	"github.com/m87/ctx/core"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	session := core.CreateTestSession()
	output := List(*session)
	expected := "- Test Context\n- Test2 Context\n"
	assert.NotEmpty(t, output)
	assert.Equal(t, expected, output)
}

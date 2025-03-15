package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateShaFromDescription(t *testing.T) {
	expected := "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
	description := "test"

	id, err := Id(description, false)
	assert.NoError(t, err)
	assert.Equal(t, id, expected)

}

func TestReturnErrIfEmptyDescription(t *testing.T) {
	_, err := Id("", false)
	assert.Error(t, err)

	_, err = Id(" \t", false)
	assert.Error(t, err)
}

func TestReturnIdAsIsIfIsRawFlagSet(t *testing.T) {
	expected := "test"

	id, err := Id(expected, true)
	assert.NoError(t, err)
	assert.Equal(t, id, expected)
}

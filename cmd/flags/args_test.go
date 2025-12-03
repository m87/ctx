package flags

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMissingArgument(t *testing.T) {
	args := []string{
		"test1",
	}

	_, err := GetStringArg(args, 1, "name")

	assert.Error(t, err)
	assert.ErrorContains(t, err, "missing name")

}

func TestGetStringArgument(t *testing.T) {
	args := []string{
		"21",
	}

	val, err := GetStringArg(args, 0, "name")

	assert.NoError(t, err)
	assert.Equal(t, "21", val)
}

func TestResolveArgAsId(t *testing.T) {
	args := []string{
		"21",
	}

	val, err := ResolveArgumentAsContextId(args, 0, "name")

	assert.NoError(t, err)
	assert.Equal(t, "6f4b6612125fb3a0daecd2799dfd6c9c299424fd920f9b308110a2c1fbd8f443", val)
}

func TestConditionalIndexProvider(t *testing.T) {
	providerWithCtxId := ConditionalIndexProvider(true)
	providerWithoutCtxId := ConditionalIndexProvider(false)
	assert.Equal(t, 0, providerWithCtxId(1))
	assert.Equal(t, 1, providerWithoutCtxId(1))
}

package flags

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestResolveIntervalIdWithIdFlag(t *testing.T) {
	cmd := cobra.Command{}
	AddIntervalFlag(&cmd)

	args := []string{
		"--interval=test",
	}
	cmd.SetArgs(args)
	cmd.ParseFlags(args)

	id, err := ResolveIntervalId(&cmd)

	assert.NoError(t, err)
	assert.Equal(t, "test", id)
}

func TestResolveIntervalIdNoFlags(t *testing.T) {
	cmd := cobra.Command{}
	AddIntervalFlag(&cmd)

	args := []string{}
	cmd.SetArgs(args)
	cmd.ParseFlags(args)

	_, err := ResolveIntervalId(&cmd)

	assert.ErrorContains(t, err, "--interval must be provided")
}

func TestResolveIntervalIdWithShortIdFlag(t *testing.T) {
	cmd := cobra.Command{}
	AddIntervalFlag(&cmd)

	args := []string{
		"-i=test2",
	}
	cmd.SetArgs(args)
	cmd.ParseFlags(args)

	id, err := ResolveIntervalId(&cmd)

	assert.NoError(t, err)
	assert.Equal(t, "test2", id)
}

func TestResolveIntervalIdWithEmptyIdFlag(t *testing.T) {
	cmd := cobra.Command{}
	AddIntervalFlag(&cmd)

	args := []string{
		"-i",
	}
	cmd.SetArgs(args)
	cmd.ParseFlags(args)

	_, err := ResolveIntervalId(&cmd)

	assert.Error(t, err)
	assert.ErrorContains(t, err, "--interval must be provided")

}

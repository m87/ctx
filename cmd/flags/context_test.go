package cmd

import (
	"testing"

	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestResolveContextId_WithCtxAndCtxIdFlag(t *testing.T) {
	cmd := cobra.Command{}
	AddContxtFlag(&cmd)

	args := []string{
		"--ctx=test2",
		"--ctx-id=test",
	}
	cmd.SetArgs(args)
	cmd.ParseFlags(args)

  _, err := ResolveId(&cmd)

	assert.Error(t, err)
	assert.ErrorContains(t, err, "both --ctx and -ctx-id provided")
}

func TestResolveContextId_WithCtxFlag(t *testing.T) {
	cmd := cobra.Command{}
	AddContxtFlag(&cmd)
	args := []string{
		"-c=test2",
	}
	cmd.SetArgs(args)
	cmd.ParseFlags(args)

	id, err := ResolveId(&cmd)

	assert.NoError(t, err)
	assert.Equal(t, util.GenerateId("test2"), id)
}

func TestResolveContextId_NoFlags(t *testing.T) {
	cmd := cobra.Command{}
	AddContxtFlag(&cmd)

	args := []string{
	}
	cmd.SetArgs(args)
	cmd.ParseFlags(args)

	_, err := ResolveId(&cmd)

	assert.ErrorContains(t, err, "either --ctx or --ctx-id must be provided")
}

func TestResolveContextId_WithEmptyCtxAndCtxIdFlag(t *testing.T) {
	cmd := cobra.Command{}
	AddContxtFlag(&cmd)

	args := []string{
		"-C=test2",
	}
	cmd.SetArgs(args)
	cmd.ParseFlags(args)

	id, err := ResolveId(&cmd)

	assert.NoError(t, err)
	assert.Equal(t, "test2", id)
}

func TestResolveContextId_WithEmptyCtxIdFlag(t *testing.T) {
	cmd := cobra.Command{}
	AddContxtFlag(&cmd)

	args := []string{
		"-ctx=test",
		"-C=",
	}
	cmd.SetArgs(args)
	cmd.ParseFlags(args)

	_, err := ResolveId(&cmd)

	assert.Error(t, err)
	assert.ErrorContains(t, err, "both --ctx and -ctx-id provided")

}

func TestResolveContextId_WithEmptyCtxFlag(t *testing.T) {
	cmd := cobra.Command{}
	AddContxtFlag(&cmd)
	
	args := []string{
		"-c=",
		"--ctx-id=test",
	}
	cmd.SetArgs(args)
	cmd.ParseFlags(args)

	_, err := ResolveId(&cmd)

	assert.Error(t, err)
	assert.ErrorContains(t, err, "both --ctx and -ctx-id provided")

}

package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSwitchContext(t *testing.T) {
	test := NewTestContextManager()
	context := test.Contexts.Get(TestContextID)

	err := test.Manager.SwitchContext(context)
	require.NoError(t, err)

	c, err := test.Contexts.GetActive()
	require.NoError(t, err)
	require.NotNil(t, c)
	require.Equal(t, TestContextID, c.Id)
	require.Equal(t, "active", c.Status)

	interval, err := test.Intervals.GetActiveIntervalByContextId(c.Id)
	require.NoError(t, err)
	require.NotNil(t, interval)
	require.Equal(t, test.TimeProvider.Now(), *interval.Start)
	require.Equal(t, c.WorkspaceId, interval.WorkspaceId)
}

func TestSwitchContextToSameContext(t *testing.T) {
	test := NewTestContextManager()
	context := test.Contexts.Get(TestContextID)

	err := test.Manager.SwitchContext(context)
	require.NoError(t, err)

	err = test.Manager.SwitchContext(context)
	require.NoError(t, err)

	c, err := test.Contexts.GetActive()
	require.NoError(t, err)
	require.NotNil(t, c)
	require.Equal(t, TestContextID, c.Id)
	require.Equal(t, "active", c.Status)

	intervals, err := test.Intervals.ListByContextId(c.Id)
	require.NoError(t, err)
	require.Len(t, intervals, 2)
}

func TestSwitchContextToArchivedContext(t *testing.T) {
	test := NewTestContextManager()
	context := test.Contexts.Get(TestContextID)
	context.Archived = true
	test.Contexts.Save(context)

	err := test.Manager.SwitchContext(context)
	require.Error(t, err)
	require.IsType(t, &ContextArchivedError{}, err)

	c, err := test.Contexts.GetActive()
	require.NoError(t, err)
	require.Nil(t, c)

	intervals, err := test.Intervals.ListByContextId(context.Id)
	require.NoError(t, err)
	require.Len(t, intervals, 1)
}

func TestSwitchContextToNewContext(t *testing.T) {
	test := NewTestContextManager()
	newContext := &Context{
		Name:        "New Context",
		WorkspaceId: TestWorkspaceID,
	}

	err := test.Manager.SwitchContext(newContext)
	require.NoError(t, err)

	c, err := test.Contexts.GetActive()
	require.NoError(t, err)
	require.NotNil(t, c)
	require.Equal(t, "New Context", c.Name)
	require.Equal(t, "active", c.Status)

	interval, err := test.Intervals.GetActiveIntervalByContextId(c.Id)
	require.NoError(t, err)
	require.NotNil(t, interval)
	require.Equal(t, test.TimeProvider.Now(), *interval.Start)
	require.Equal(t, c.WorkspaceId, interval.WorkspaceId)
}

func TestFreeActiveContext(t *testing.T) {
	test := NewTestContextManager()
	test.Intervals.Seed()
	context := test.Contexts.Get(TestContextID)

	err := test.Manager.SwitchContext(context)
	require.NoError(t, err)

	err = test.Manager.FreeActiveContext()
	require.NoError(t, err)

	c, err := test.Contexts.GetActive()
	require.NoError(t, err)
	require.Nil(t, c)

	intervals, err := test.Intervals.ListByContextId(context.Id)
	require.NoError(t, err)
	require.Len(t, intervals, 1)
	require.Equal(t, "completed", intervals[0].Status)
}

func TestFreeActiveContextWhenNoActiveContext(t *testing.T) {
	test := NewTestContextManager()

	err := test.Manager.FreeActiveContext()
	require.NoError(t, err)

	c, err := test.Contexts.GetActive()
	require.NoError(t, err)
	require.Nil(t, c)
}

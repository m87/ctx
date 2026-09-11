package core

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestQueryContextsReturnsAllContextsWithPassthroughInterpreter(t *testing.T) {
	test := NewEmptyTestContextManager()
	_, err := test.Contexts.Save(&Context{
		Id:          "context-1",
		Name:        "Active",
		WorkspaceId: "workspace-1",
	})
	require.NoError(t, err)
	_, err = test.Contexts.Save(&Context{
		Id:          "context-2",
		Name:        "Archived",
		WorkspaceId: "workspace-1",
		Archived:    true,
	})
	require.NoError(t, err)
	_, err = test.Contexts.Save(&Context{
		Id:          "context-3",
		Name:        "Other workspace",
		WorkspaceId: "workspace-2",
	})
	require.NoError(t, err)

	contexts, err := test.Manager.QueryContexts("ignored query")

	require.NoError(t, err)
	require.Equal(t, &ContextSQLQuery{}, test.Contexts.lastQuery)
	require.Equal(t, []*Context{
		test.Contexts.Get("context-1"),
		test.Contexts.Get("context-2"),
		test.Contexts.Get("context-3"),
	}, contexts)
}

func TestQueryContextsDelegatesToConfiguredInterpreter(t *testing.T) {
	test := NewEmptyTestContextManager()
	_, err := test.Contexts.Save(&Context{
		Id:          "context-1",
		Name:        "Context",
		WorkspaceId: "workspace-1",
	})
	require.NoError(t, err)

	interpreter := &recordingContextQueryInterpreter{
		err: &InvalidContextQueryError{Err: errors.New("unexpected token")},
	}
	test.Manager.QueryInterpreter = interpreter

	contexts, err := test.Manager.QueryContexts("broken query")

	require.Nil(t, contexts)
	require.EqualError(t, err, "invalid context query: unexpected token")
	require.Equal(t, "broken query", interpreter.query)
	require.Nil(t, test.Contexts.lastQuery)
}

func TestQueryContextsExecutesTheInterpretedSQLPlan(t *testing.T) {
	test := NewEmptyTestContextManager()
	plan := &ContextSQLQuery{
		WhereClause: "name = ?",
		Arguments:   []any{"Context"},
	}
	interpreter := &recordingContextQueryInterpreter{result: plan}
	test.Manager.QueryInterpreter = interpreter

	_, err := test.Manager.QueryContexts("name = Context")

	require.NoError(t, err)
	require.Equal(t, "name = Context", interpreter.query)
	require.Equal(t, plan, test.Contexts.lastQuery)
}

func TestQueryWorkspaceContextsScopesPassthroughResultsAndBuildsSummary(t *testing.T) {
	test := NewEmptyTestContextManager()
	_, err := test.Workspaces.Save(&Workspace{Id: "workspace-1", Name: "Workspace"})
	require.NoError(t, err)
	_, err = test.Contexts.Save(&Context{
		Id:          "context-1",
		Name:        "Tracked",
		WorkspaceId: "workspace-1",
	})
	require.NoError(t, err)
	_, err = test.Contexts.Save(&Context{
		Id:          "context-2",
		Name:        "Untracked",
		WorkspaceId: "workspace-1",
		Archived:    true,
	})
	require.NoError(t, err)
	_, err = test.Contexts.Save(&Context{
		Id:          "context-3",
		Name:        "Other workspace",
		WorkspaceId: "workspace-2",
	})
	require.NoError(t, err)
	_, err = test.Intervals.Save(&Interval{
		Id:        "interval-1",
		ContextId: "context-1",
		Duration:  45 * time.Minute,
		Status:    "completed",
	})
	require.NoError(t, err)

	result, err := test.Manager.QueryWorkspaceContexts(" workspace-1 ", "anything")

	require.NoError(t, err)
	require.Equal(t, "workspace-1", result.WorkspaceId)
	require.Equal(t, "anything", result.Query)
	require.Equal(t, "workspace-1", test.Contexts.lastQuery.WorkspaceId)
	require.ElementsMatch(t, []string{"context-1", "context-2"}, []string{
		result.Contexts[0].Id,
		result.Contexts[1].Id,
	})
	require.Equal(t, int64(45*time.Minute), result.TotalDuration)
	require.Equal(t, 1, result.TotalSessions)
	require.Len(t, result.ContextStats, 2)
	require.Equal(t, "context-1", result.ContextStats[0].ContextId)
	require.Equal(t, 45*time.Minute, result.ContextStats[0].Duration)
	require.Equal(t, float64(100), result.ContextStats[0].Percentage)
	require.Equal(t, 1, result.ContextStats[0].IntervalCount)
}

func TestQueryWorkspaceContextsDoesNotMutateInterpreterPlan(t *testing.T) {
	test := NewEmptyTestContextManager()
	_, err := test.Workspaces.Save(&Workspace{Id: "workspace-1", Name: "Workspace"})
	require.NoError(t, err)
	plan := &ContextSQLQuery{WhereClause: "name = ?", Arguments: []any{"Context"}}
	test.Manager.QueryInterpreter = &recordingContextQueryInterpreter{result: plan}

	_, err = test.Manager.QueryWorkspaceContexts("workspace-1", "name = Context")

	require.NoError(t, err)
	require.Empty(t, plan.WorkspaceId)
	require.Equal(t, "workspace-1", test.Contexts.lastQuery.WorkspaceId)
}

func TestQueryWorkspaceContextsRequiresExistingWorkspace(t *testing.T) {
	test := NewEmptyTestContextManager()

	result, err := test.Manager.QueryWorkspaceContexts("missing", "anything")

	require.Nil(t, result)
	var notFound *WorkspaceNotFoundError
	require.ErrorAs(t, err, &notFound)
}

type recordingContextQueryInterpreter struct {
	query  string
	result *ContextSQLQuery
	err    error
}

func (i *recordingContextQueryInterpreter) Interpret(query string) (*ContextSQLQuery, error) {
	i.query = query
	return i.result, i.err
}

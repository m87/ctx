package core

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueryContextsReturnsAllContextsWithPassthroughInterpreter(t *testing.T) {
	test := newTestManager()
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
		test.Contexts.items["context-1"],
		test.Contexts.items["context-2"],
		test.Contexts.items["context-3"],
	}, contexts)
}

func TestQueryContextsDelegatesToConfiguredInterpreter(t *testing.T) {
	test := newTestManager()
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
	test := newTestManager()
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

type recordingContextQueryInterpreter struct {
	query  string
	result *ContextSQLQuery
	err    error
}

func (i *recordingContextQueryInterpreter) Interpret(query string) (*ContextSQLQuery, error) {
	i.query = query
	return i.result, i.err
}

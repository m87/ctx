package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContextManagerCreatesAndListsSavedQueries(t *testing.T) {
	test := NewTestContextManager()

	query := &SavedQuery{
		Id:          "client-id",
		WorkspaceId: " " + TestWorkspaceID + " ",
		Name:        " Focus work ",
		Query:       " project = ctx ",
	}
	id, err := test.Manager.CreateSavedQuery(query)

	require.NoError(t, err)
	require.Equal(t, "saved-query-1", id)
	require.Equal(t, TestWorkspaceID, query.WorkspaceId)
	require.Equal(t, "Focus work", query.Name)
	require.Equal(t, "project = ctx", query.Query)
	require.Equal(t, "saved-query-1", query.Id)

	queries, err := test.Manager.ListSavedQueries(TestWorkspaceID)
	require.NoError(t, err)
	require.Equal(t, []*SavedQuery{query}, queries)

	retrieved, err := test.Manager.GetSavedQuery(id)
	require.NoError(t, err)
	require.Equal(t, query, retrieved)
}

func TestContextManagerRejectsInvalidSavedQueries(t *testing.T) {
	test := NewTestContextManager()

	tests := []struct {
		name  string
		query *SavedQuery
	}{
		{name: "nil", query: nil},
		{name: "missing workspace", query: &SavedQuery{Name: "Name", Query: "query"}},
		{name: "unknown workspace", query: &SavedQuery{WorkspaceId: "missing", Name: "Name", Query: "query"}},
		{name: "missing name", query: &SavedQuery{WorkspaceId: TestWorkspaceID, Query: "query"}},
		{name: "missing query", query: &SavedQuery{WorkspaceId: TestWorkspaceID, Name: "Name"}},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			id, err := test.Manager.CreateSavedQuery(testCase.query)
			require.Error(t, err)
			require.Empty(t, id)
		})
	}
}

func TestContextManagerReturnsSavedQueryNotFound(t *testing.T) {
	test := NewEmptyTestContextManager()

	query, err := test.Manager.GetSavedQuery("missing")

	require.Nil(t, query)
	var notFound *SavedQueryNotFoundError
	require.ErrorAs(t, err, &notFound)
}

func TestContextManagerDeletesSavedQuery(t *testing.T) {
	test := NewEmptyTestContextManager()
	repository := test.SavedQueries
	repository.Seed(&SavedQuery{Id: "query-1", Name: "Focus"})

	require.NoError(t, test.Manager.DeleteSavedQuery(" query-1 "))
	require.Nil(t, repository.Get("query-1"))

	err := test.Manager.DeleteSavedQuery("query-1")
	var notFound *SavedQueryNotFoundError
	require.ErrorAs(t, err, &notFound)
}

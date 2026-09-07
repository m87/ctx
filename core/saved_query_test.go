package core

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContextManagerCreatesAndListsSavedQueries(t *testing.T) {
	test := newTestManager()
	repository := &memorySavedQueryRepository{items: make(map[string]*SavedQuery)}
	test.Manager.SavedQueryRepository = repository
	_, err := test.Workspaces.Save(&Workspace{Id: "workspace-1", Name: "Workspace"})
	require.NoError(t, err)

	query := &SavedQuery{
		Id:          "client-id",
		WorkspaceId: " workspace-1 ",
		Name:        " Focus work ",
		Query:       " project = ctx ",
	}
	id, err := test.Manager.CreateSavedQuery(query)

	require.NoError(t, err)
	require.Equal(t, "saved-query-1", id)
	require.Equal(t, "workspace-1", query.WorkspaceId)
	require.Equal(t, "Focus work", query.Name)
	require.Equal(t, "project = ctx", query.Query)
	require.Equal(t, "saved-query-1", query.Id)

	queries, err := test.Manager.ListSavedQueries("workspace-1")
	require.NoError(t, err)
	require.Equal(t, []*SavedQuery{query}, queries)

	retrieved, err := test.Manager.GetSavedQuery(id)
	require.NoError(t, err)
	require.Equal(t, query, retrieved)
}

func TestContextManagerRejectsInvalidSavedQueries(t *testing.T) {
	test := newTestManager()
	test.Manager.SavedQueryRepository = &memorySavedQueryRepository{items: make(map[string]*SavedQuery)}
	_, err := test.Workspaces.Save(&Workspace{Id: "workspace-1", Name: "Workspace"})
	require.NoError(t, err)

	tests := []struct {
		name  string
		query *SavedQuery
	}{
		{name: "nil", query: nil},
		{name: "missing workspace", query: &SavedQuery{Name: "Name", Query: "query"}},
		{name: "unknown workspace", query: &SavedQuery{WorkspaceId: "missing", Name: "Name", Query: "query"}},
		{name: "missing name", query: &SavedQuery{WorkspaceId: "workspace-1", Query: "query"}},
		{name: "missing query", query: &SavedQuery{WorkspaceId: "workspace-1", Name: "Name"}},
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
	test := newTestManager()
	test.Manager.SavedQueryRepository = &memorySavedQueryRepository{items: make(map[string]*SavedQuery)}

	query, err := test.Manager.GetSavedQuery("missing")

	require.Nil(t, query)
	var notFound *SavedQueryNotFoundError
	require.ErrorAs(t, err, &notFound)
}

func TestContextManagerDeletesSavedQuery(t *testing.T) {
	test := newTestManager()
	repository := &memorySavedQueryRepository{items: make(map[string]*SavedQuery)}
	test.Manager.SavedQueryRepository = repository
	repository.items["query-1"] = &SavedQuery{Id: "query-1", Name: "Focus"}

	require.NoError(t, test.Manager.DeleteSavedQuery(" query-1 "))
	require.NotContains(t, repository.items, "query-1")

	err := test.Manager.DeleteSavedQuery("query-1")
	var notFound *SavedQueryNotFoundError
	require.ErrorAs(t, err, &notFound)
}

type memorySavedQueryRepository struct {
	items  map[string]*SavedQuery
	nextID int
}

func (r *memorySavedQueryRepository) GetById(id string) (*SavedQuery, error) {
	return r.items[id], nil
}

func (r *memorySavedQueryRepository) Save(query *SavedQuery) (string, error) {
	if query == nil {
		return "", fmt.Errorf("saved query is required")
	}
	if query.Id == "" {
		r.nextID++
		query.Id = fmt.Sprintf("saved-query-%d", r.nextID)
	}
	r.items[query.Id] = query
	return query.Id, nil
}

func (r *memorySavedQueryRepository) ListByWorkspace(workspaceId string) ([]*SavedQuery, error) {
	queries := make([]*SavedQuery, 0)
	for _, query := range r.items {
		if query.WorkspaceId == workspaceId {
			queries = append(queries, query)
		}
	}
	sort.Slice(queries, func(i, j int) bool { return queries[i].Name < queries[j].Name })
	return queries, nil
}

func (r *memorySavedQueryRepository) Delete(id string) error {
	delete(r.items, id)
	return nil
}

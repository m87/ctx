package core

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTagJSONUsesObjectFormat(t *testing.T) {
	encoded, err := json.Marshal(&Context{Tags: []*Tag{{Id: "tag-1", Name: "important"}}})
	require.NoError(t, err)
	require.JSONEq(t, `{"projectId":null,"id":"","name":"","workspaceId":"","status":"","archived":false,"tags":[{"id":"tag-1","name":"important"}],"project":null}`, string(encoded))

	var context Context
	require.NoError(t, json.Unmarshal([]byte(`{"tags":[{"id":"tag-1","name":"important"}]}`), &context))
	require.Equal(t, []*Tag{{Id: "tag-1", Name: "important"}}, context.Tags)
	require.Error(t, json.Unmarshal([]byte(`{"tags":["important"]}`), &context))
}

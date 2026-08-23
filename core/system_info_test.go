package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSystemInfoMapperUsesClientIdAsNodeId(t *testing.T) {
	info := &SystemInfo{
		DatabaseVersion: CurrentDatabaseVersion,
		ClientId:        "7aa13ad4-7af8-42db-bd6a-927ac9573d8f",
	}

	node, err := NewSystemInfoMapper().ToNode(info)
	require.NoError(t, err)
	require.Equal(t, info.ClientId, node.Core.Id)
	require.Equal(t, SystemInfoName, node.Core.Name)
	require.NotContains(t, node.KV, "client_id")

	restored, err := NewSystemInfoMapper().FromNode(node)
	require.NoError(t, err)
	require.Equal(t, info, restored)
}

func TestSystemInfoMapperRequiresClientId(t *testing.T) {
	_, err := NewSystemInfoMapper().ToNode(&SystemInfo{DatabaseVersion: CurrentDatabaseVersion})
	require.ErrorContains(t, err, "client ID is required")
}

func TestDatabaseVersionNeedsMigration(t *testing.T) {
	tests := []struct {
		current string
		want    bool
	}{
		{current: "", want: true},
		{current: "0.4.9", want: true},
		{current: "0.5.0", want: false},
		{current: "0.5.1", want: false},
		{current: "1.0.0", want: false},
	}

	for _, test := range tests {
		got, err := DatabaseVersionNeedsMigration(test.current, "0.5.0")
		require.NoError(t, err)
		require.Equal(t, test.want, got, test.current)
	}
}

func TestDatabaseVersionNeedsMigrationRejectsInvalidVersion(t *testing.T) {
	_, err := DatabaseVersionNeedsMigration("invalid", "0.5.0")
	require.Error(t, err)
}

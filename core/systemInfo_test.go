package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
		got, err := DatabaseVersionNeedsMigration(test.current, CurrentDatabaseVersion)
		require.NoError(t, err)
		require.Equal(t, test.want, got, test.current)
	}
}

func TestDatabaseVersionNeedsMigrationRejectsInvalidVersion(t *testing.T) {
	_, err := DatabaseVersionNeedsMigration("invalid", CurrentDatabaseVersion)
	require.Error(t, err)
}

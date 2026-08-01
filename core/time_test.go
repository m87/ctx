package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func testTime(value time.Time) *time.Time {
	return &value
}

func TestRealTimeProviderReturnsUTC(t *testing.T) {
	now := (&RealTimeProvider{}).Now()

	require.Equal(t, time.UTC, now.Location())
	require.False(t, now.IsZero())
}

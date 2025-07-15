package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDeleteInterval(t *testing.T) {
	session := CreateTestSession()

	assert.Len(t, session.State.Contexts[TEST_ID].Intervals, 2)
	assert.Equal(t, session.State.Contexts[TEST_ID].Duration, 2*time.Hour)

	err := session.DeleteInterval(TEST_ID, TEST_INTERVAL_ID)
	assert.NoError(t, err)
	assert.Len(t, session.State.Contexts[TEST_ID].Intervals, 1)
	assert.Equal(t, session.State.Contexts[TEST_ID].Duration, 1*time.Hour)

}

func TestDeleteIntervalNonExistentContext(t *testing.T) {
	session := CreateTestSession()

	err := session.DeleteInterval("non-existent-context", TEST_INTERVAL_ID)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "context does not exist")
}
func TestDeleteIntervalNonExistentInterval(t *testing.T) {
	session := CreateTestSession()

	err := session.DeleteInterval(TEST_ID, "non-existent-interval")
	assert.Error(t, err)
	assert.ErrorContains(t, err, "interval does not exist")
}

func TestDeleteINtervalActiveContext(t *testing.T) {
	session := CreateTestSession()

	session.State.CurrentId = TEST_ID

	err := session.DeleteInterval(TEST_ID, TEST_INTERVAL_ID)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "context is active")
}

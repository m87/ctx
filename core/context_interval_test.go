package core

import (
	"testing"
	"time"

	ctxtime "github.com/m87/ctx/time"
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

func TestGetActiveIntervals(t *testing.T) {
	session := CreateTestSession()
	session.State.Contexts[TEST_ID].Intervals["active1"] = Interval{
		Id: "active1",
		Start: ctxtime.ZonedTime{
			Time: time.Now(),
		},
		End: ctxtime.ZonedTime{},
	}
	session.State.Contexts[TEST_ID].Intervals["active2"] = Interval{
		Id: "active2",
		Start: ctxtime.ZonedTime{
			Time: time.Now(),
		},
		End: ctxtime.ZonedTime{},
	}

	active, err := session.GetActiveIntervals(TEST_ID)

	assert.NoError(t, err)
	assert.Len(t, session.State.Contexts[TEST_ID].Intervals, 4)
	assert.Len(t, active, 2)
	assert.Contains(t, active, "active1")
	assert.Contains(t, active, "active2")

}

func TestGetActiveIntervalsNonExistentContext(t *testing.T) {
	session := CreateTestSession()

	_, err := session.GetActiveIntervals("non-existent-context")
	assert.Error(t, err)
	assert.ErrorContains(t, err, "context does not exist")
}

func TestGetActiveIntervalsNoActiveIntervals(t *testing.T) {
	session := CreateTestSession()

	active, err := session.GetActiveIntervals(TEST_ID)
	assert.NoError(t, err)
	assert.Len(t, active, 0)
}

func TestGetActiveIntervl(t *testing.T) {
	session := CreateTestSession()

	interval, err := session.GetActiveIntervals(TEST_ID)
	assert.NoError(t, err)
	assert.NotNil(t, interval)

	assert.Len(t, session.State.Contexts[TEST_ID].Intervals, 2)
	assert.Contains(t, session.State.Contexts[TEST_ID].Intervals, TEST_INTERVAL_ID)
	assert.Contains(t, session.State.Contexts[TEST_ID].Intervals, TEST_INTERVAL_2_ID)
}

func TestEndInterval(t *testing.T) {
	session := CreateTestSession()

	assert.Len(t, session.State.Contexts[TEST_ID].Intervals, 2)
	assert.Equal(t, session.State.Contexts[TEST_ID].Duration, 2*time.Hour)

	err := session.endInterval(TEST_ID, session.TimeProvider.Now())
	assert.NoError(t, err)

	interval := session.State.Contexts[TEST_ID].Intervals[TEST_INTERVAL_2_ID]
	assert.NotEqual(t, interval.End.Time, session.TimeProvider.Now().Time)
}

func TestEndIntervalNonExistentContext(t *testing.T) {
	session := CreateTestSession()

	err := session.endInterval("non-existent-context", session.TimeProvider.Now())
	assert.Error(t, err)
	assert.ErrorContains(t, err, "context does not exist")
}

func TestEndIntervalNoActiveIntervals(t *testing.T) {
	session := CreateTestSession()
	ctx := session.MustGetCtx(TEST_ID)
	ctx.Intervals = map[string]Interval{}
	session.State.Contexts[TEST_ID] = ctx

	err := session.endInterval(TEST_ID, session.TimeProvider.Now())
	assert.NoError(t, err)

	assert.Len(t, session.State.Contexts[TEST_ID].Intervals, 0)
}

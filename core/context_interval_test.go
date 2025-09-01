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

func TestDeleteIntervalActiveContext(t *testing.T) {
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

	session.State.Contexts[TEST_ID].Intervals[TEST_INTERVAL_2_ID] = Interval{
		Id:       TEST_INTERVAL_2_ID,
		Start:    session.TimeProvider.Now(),
		End:      ctxtime.ZonedTime{},
		Duration: 0,
	}

	dt, _ := time.Parse(time.DateTime, "2025-02-02T15:15:15Z")
	session.TimeProvider = &TestTimeProvider{
		currentTime: dt,
	}

	err := session.endInterval(TEST_ID, session.TimeProvider.Now())
	assert.NoError(t, err)

	interval := session.State.Contexts[TEST_ID].Intervals[TEST_INTERVAL_2_ID]
	assert.Equal(t, interval.End.Time, session.TimeProvider.Now().Time)
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

func TestGetIntervalsByData(t *testing.T) {
	session := CreateTestSession()

	intervals := session.GetIntervalsByDate(TEST_ID, ctxtime.ZonedTime{Time: session.TimeProvider.Now().Time})
	assert.Len(t, intervals, 2)
}

func TestGetIntervalsByDateCropEndToCurrentDay(t *testing.T) {
	session := CreateTestSession()
	ctx := session.MustGetCtx(TEST_ID)
	interval := ctx.Intervals[TEST_INTERVAL_2_ID]
	interval.End = ctxtime.ZonedTime{Time: session.TimeProvider.Now().Time.Add(48 * time.Hour), Timezone: "UTC"}
	ctx.Intervals[TEST_INTERVAL_2_ID] = interval
	session.SetCtx(ctx)

	intervals := session.GetIntervalsByDate(TEST_ID, ctxtime.ZonedTime{Time: session.TimeProvider.Now().Time})
	assert.Len(t, intervals, 2)
	for _, interval := range intervals {
		if interval.Id == TEST_INTERVAL_2_ID {
			assert.Equal(t, interval.Duration, 10*time.Hour + 47*time.Minute + 47*time.Second)
	    loc, _ := time.LoadLocation("UTC")
			dt, _ := time.ParseInLocation(time.DateTime, "2025-02-02 23:59:59", loc)
			assert.Equal(t, interval.End.Time, dt)
		}
	}
}

func TestMoveInterval(t *testing.T) {
	session := CreateTestSession()

	assert.Len(t, session.MustGetCtx(TEST_ID).Intervals, 2)
	assert.Len(t, session.MustGetCtx(TEST_ID_2).Intervals, 2)
	session.MoveIntervalById(TEST_ID, TEST_ID_2, TEST_INTERVAL_ID)
	
	assert.Len(t, session.MustGetCtx(TEST_ID).Intervals, 1)
	assert.Len(t, session.MustGetCtx(TEST_ID_2).Intervals, 3)
}
 


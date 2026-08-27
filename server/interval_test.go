package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/m87/ctx/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSplitIntervalResponse(t *testing.T) {
	start := time.Date(2026, 8, 14, 8, 0, 0, 0, time.UTC)
	end := time.Date(2026, 8, 14, 9, 0, 0, 0, time.UTC)
	origin := &core.Interval{
		Id:          "interval-1",
		ContextId:   "context-1",
		Start:       &start,
		End:         &end,
		Duration:    time.Hour,
		Status:      "inactive",
		WorkspaceId: "workspace-1",
		Synced:      true,
	}
	repository := &splitIntervalRepository{
		intervals: map[string]*core.Interval{origin.Id: origin},
	}
	manager := core.NewContextManager(nil, nil, repository, nil, nil)
	mux := http.NewServeMux()
	registerIntervalHandler(mux, manager)

	request := httptest.NewRequest(
		http.MethodPost,
		"/interval-1/split?timeZone=UTC",
		bytes.NewBufferString(`{"splitTime":"2026-08-14T08:30:00"}`),
	)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)

	require.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "application/json", response.Header().Get("Content-Type"))

	var result core.IntervalSplitResult
	require.NoError(t, json.NewDecoder(response.Body).Decode(&result))
	assert.Equal(t, origin, result.Origin)
	assert.Equal(t, "split-1", result.SplitResult[0].Id)
	assert.Equal(t, "split-2", result.SplitResult[1].Id)
	assert.Equal(t, start, *result.SplitResult[0].Start)
	assert.Equal(t, time.Date(2026, 8, 14, 8, 30, 0, 0, time.UTC), *result.SplitResult[0].End)
	assert.Equal(t, *result.SplitResult[0].End, *result.SplitResult[1].Start)
	assert.Equal(t, end, *result.SplitResult[1].End)
}

func TestUndoSplitIntervalRestoresOriginAfterPartialClientUndo(t *testing.T) {
	result := intervalSplitResultFixture()
	repository := &splitIntervalRepository{
		intervals: map[string]*core.Interval{
			result.SplitResult[1].Id: result.SplitResult[1],
		},
	}
	manager := core.NewContextManager(
		nil,
		&splitContextRepository{context: &core.Context{Id: result.Origin.ContextId, WorkspaceId: result.Origin.WorkspaceId}},
		repository,
		nil,
		nil,
	)
	transactionCalled := false
	manager.RunInTransaction = func(fn func(*core.ContextManager) error) error {
		transactionCalled = true
		return fn(manager)
	}
	mux := http.NewServeMux()
	registerIntervalHandler(mux, manager)
	body, err := json.Marshal(result)
	require.NoError(t, err)

	request := httptest.NewRequest(http.MethodPost, "/interval-1/split/undo", bytes.NewReader(body))
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)

	require.Equal(t, http.StatusOK, response.Code)
	assert.True(t, transactionCalled)
	assert.Len(t, repository.intervals, 1)
	assert.Equal(t, result.Origin, repository.intervals[result.Origin.Id])
}

func TestUndoSplitIntervalRollsBackWhenDeleteFails(t *testing.T) {
	result := intervalSplitResultFixture()
	repository := &splitIntervalRepository{
		intervals: map[string]*core.Interval{
			result.SplitResult[0].Id: result.SplitResult[0],
			result.SplitResult[1].Id: result.SplitResult[1],
		},
		failDeleteId: result.SplitResult[1].Id,
	}
	manager := core.NewContextManager(
		nil,
		&splitContextRepository{context: &core.Context{Id: result.Origin.ContextId, WorkspaceId: result.Origin.WorkspaceId}},
		repository,
		nil,
		nil,
	)
	manager.RunInTransaction = func(fn func(*core.ContextManager) error) error {
		before := maps.Clone(repository.intervals)
		err := fn(manager)
		if err != nil {
			repository.intervals = before
		}
		return err
	}
	mux := http.NewServeMux()
	registerIntervalHandler(mux, manager)
	body, err := json.Marshal(result)
	require.NoError(t, err)

	request := httptest.NewRequest(http.MethodPost, "/interval-1/split/undo", bytes.NewReader(body))
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)

	require.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Contains(t, repository.intervals, result.SplitResult[0].Id)
	assert.Contains(t, repository.intervals, result.SplitResult[1].Id)
	assert.NotContains(t, repository.intervals, result.Origin.Id)
}

func intervalSplitResultFixture() *core.IntervalSplitResult {
	start := time.Date(2026, 8, 14, 8, 0, 0, 0, time.UTC)
	splitTime := time.Date(2026, 8, 14, 8, 30, 0, 0, time.UTC)
	end := time.Date(2026, 8, 14, 9, 0, 0, 0, time.UTC)
	origin := &core.Interval{
		Id:          "interval-1",
		ContextId:   "context-1",
		Start:       &start,
		End:         &end,
		Duration:    time.Hour,
		Status:      "inactive",
		WorkspaceId: "workspace-1",
	}
	return &core.IntervalSplitResult{
		Origin: origin,
		SplitResult: [2]*core.Interval{
			{Id: "split-1", ContextId: origin.ContextId, Start: &start, End: &splitTime, Duration: 30 * time.Minute, WorkspaceId: origin.WorkspaceId},
			{Id: "split-2", ContextId: origin.ContextId, Start: &splitTime, End: &end, Duration: 30 * time.Minute, WorkspaceId: origin.WorkspaceId},
		},
	}
}

type splitIntervalRepository struct {
	core.IntervalRepository
	intervals    map[string]*core.Interval
	nextId       int
	failDeleteId string
}

func (r *splitIntervalRepository) GetById(id string) (*core.Interval, error) {
	interval, ok := r.intervals[id]
	if !ok {
		return nil, errors.New("interval not found")
	}
	return interval, nil
}

func (r *splitIntervalRepository) Save(interval *core.Interval) (string, error) {
	if interval.Id == "" {
		r.nextId++
		interval.Id = fmt.Sprintf("split-%d", r.nextId)
	}
	r.intervals[interval.Id] = interval
	return interval.Id, nil
}

func (r *splitIntervalRepository) Delete(id string) error {
	if id == r.failDeleteId {
		return errors.New("failed to delete interval")
	}
	delete(r.intervals, id)
	return nil
}

type splitContextRepository struct {
	core.ContextRepository
	context *core.Context
}

func (r *splitContextRepository) GetById(id string) (*core.Context, error) {
	if r.context == nil || r.context.Id != id {
		return nil, errors.New("context not found")
	}
	return r.context, nil
}

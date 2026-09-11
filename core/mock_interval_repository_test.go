package core

import (
	"errors"
	"fmt"
	"time"
)

var _ IntervalRepository = (*IntervalRepositoryMock)(nil)

type IntervalRepositoryMock struct {
	intervals            []*Interval
	nextID               int
	listCalled           bool
	saved                []*Interval
	deletedIDs           []string
	deletedContextIDs    []string
	listError            error
	saveError            error
	deleteError          error
	deleteByContextError error
}

func NewIntervalRepositoryMock(intervals ...*Interval) *IntervalRepositoryMock {
	mock := &IntervalRepositoryMock{}
	mock.Seed(intervals...)
	return mock
}

func (m *IntervalRepositoryMock) Seed(intervals ...*Interval) {
	m.intervals = copyPointerSlice(intervals)
	m.nextID = 0
	m.listCalled = false
	m.saved = nil
	m.deletedIDs = nil
	m.deletedContextIDs = nil
}

func (m *IntervalRepositoryMock) Get(id string) *Interval {
	interval, _ := m.GetById(id)
	return interval
}

func (m *IntervalRepositoryMock) GetById(id string) (*Interval, error) {
	for _, interval := range m.intervals {
		if interval != nil && interval.Id == id {
			return interval, nil
		}
	}
	return nil, nil
}

func (m *IntervalRepositoryMock) Save(interval *Interval) (string, error) {
	if m.saveError != nil {
		return "", m.saveError
	}
	if interval == nil {
		return "", errors.New("interval is required")
	}
	if interval.Id == "" {
		interval.Id = m.nextIntervalID()
	}
	for i, existing := range m.intervals {
		if existing != nil && existing.Id == interval.Id {
			m.intervals[i] = interval
			m.saved = append(m.saved, interval)
			return interval.Id, nil
		}
	}
	m.intervals = append(m.intervals, interval)
	m.saved = append(m.saved, interval)
	return interval.Id, nil
}

func (m *IntervalRepositoryMock) SaveAll(intervals []*Interval) ([]string, error) {
	ids := make([]string, 0, len(intervals))
	for _, interval := range intervals {
		id, err := m.Save(interval)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (m *IntervalRepositoryMock) Delete(id string) error {
	m.deletedIDs = append(m.deletedIDs, id)
	if m.deleteError != nil {
		return m.deleteError
	}
	for i, interval := range m.intervals {
		if interval != nil && interval.Id == id {
			m.intervals = append(m.intervals[:i], m.intervals[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *IntervalRepositoryMock) DeleteByContextId(contextID string) error {
	m.deletedContextIDs = append(m.deletedContextIDs, contextID)
	if m.deleteByContextError != nil {
		return m.deleteByContextError
	}
	intervals := m.intervals[:0]
	for _, interval := range m.intervals {
		if interval == nil || interval.ContextId != contextID {
			intervals = append(intervals, interval)
		}
	}
	m.intervals = intervals
	return nil
}

func (m *IntervalRepositoryMock) ListByContextId(contextID string) ([]*Interval, error) {
	intervals := make([]*Interval, 0)
	for _, interval := range m.intervals {
		if interval != nil && interval.ContextId == contextID {
			intervals = append(intervals, interval)
		}
	}
	return intervals, nil
}

func (m *IntervalRepositoryMock) GetActiveIntervalByContextId(contextID string) (*Interval, error) {
	for _, interval := range m.intervals {
		if interval != nil && interval.ContextId == contextID && interval.Status == IntervalStatusActive {
			return interval, nil
		}
	}
	return nil, nil
}

func (m *IntervalRepositoryMock) ListByDay(date time.Time, workspaceID string) ([]*Interval, error) {
	location := date.Location()
	localDayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, location)
	dayStart := localDayStart.UTC()
	dayEnd := localDayStart.AddDate(0, 0, 1).UTC()
	intervals := make([]*Interval, 0)
	for _, interval := range m.intervals {
		if interval == nil || interval.WorkspaceId != workspaceID || !timeIsSet(interval.Start) {
			continue
		}
		if interval.Start.After(dayEnd) {
			continue
		}
		if timeIsSet(interval.End) && interval.End.Before(dayStart) {
			continue
		}
		intervals = append(intervals, interval)
	}
	return intervals, nil
}

func (m *IntervalRepositoryMock) List() ([]*Interval, error) {
	m.listCalled = true
	if m.listError != nil {
		return nil, m.listError
	}
	return copyPointerSlice(m.intervals), nil
}

func (m *IntervalRepositoryMock) ListToSync(limit int) ([]*Interval, error) {
	intervals := make([]*Interval, 0)
	for _, interval := range m.intervals {
		if interval != nil && !interval.Synced {
			intervals = append(intervals, interval)
		}
	}
	return limited(intervals, limit), nil
}

func (m *IntervalRepositoryMock) nextIntervalID() string {
	for {
		m.nextID++
		id := fmt.Sprintf("interval-%d", m.nextID)
		interval, _ := m.GetById(id)
		if interval == nil {
			return id
		}
	}
}

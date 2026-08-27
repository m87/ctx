package core

import (
	"errors"
	"github.com/google/uuid"
)

func NewIntervalRepositoryMock(intervals []*Interval) *IntervalRepositoryMock {
	return &IntervalRepositoryMock{
		intervals: intervals,
	}
}

type IntervalRepositoryMock struct {
	IntervalRepository
	intervals []*Interval
	called    bool
}

func (m *IntervalRepositoryMock) List() ([]*Interval, error) {
	m.called = true
	if m.intervals == nil {
		return nil, errors.New("IntervalRepository.List error")
	}
	return m.intervals, nil
}

func (m *IntervalRepositoryMock) GetById(id string) (*Interval, error) {
	m.called = true
	for _, interval := range m.intervals {
		if interval.Id == id {
			return interval, nil
		}
	}
	return nil, errors.New("IntervalRepository.GetById error")
}

func (m *IntervalRepositoryMock) Save(interval *Interval) (string, error) {
	m.called = true
	if interval.Id == "" {
		interval.Id = uuid.NewString()
	}

	for i, existingInterval := range m.intervals {
		if existingInterval.Id == interval.Id {
			m.intervals[i] = interval
			return interval.Id, nil
		}
	}
	m.intervals = append(m.intervals, interval)
	return interval.Id, nil
}

func (m *IntervalRepositoryMock) Delete(id string) error {
	m.called = true
	for i, interval := range m.intervals {
		if interval.Id == id {
			m.intervals = append(m.intervals[:i], m.intervals[i+1:]...)
			return nil
		}
	}
	return errors.New("IntervalRepository.Delete error")
}

func (m *IntervalRepositoryMock) ListByContextId(contextId string) ([]*Interval, error) {
	m.called = true
	var result []*Interval
	for _, interval := range m.intervals {
		if interval.ContextId == contextId {
			result = append(result, interval)
		}
	}
	return result, nil
}

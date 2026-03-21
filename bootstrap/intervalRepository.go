package bootstrap

import (
	"time"

	"github.com/m87/ctx/core"
	"github.com/m87/nod"
)

type IntervalRepository struct {
	repository nod.TypedRepository[core.Interval]
}

func NewIntervalRepository(repository *nod.Repository) *IntervalRepository {
	return &IntervalRepository{
		repository: nod.As[core.Interval](repository),
	}
}

func (r *IntervalRepository) GetById(id string) (*core.Interval, error) {
	return r.repository.Query().NodeId(id).First()
}

func (r *IntervalRepository) Save(interval *core.Interval) (string, error) {
	return r.repository.Save(interval)
}

func (r *IntervalRepository) Delete(id string) error {
	return r.repository.Query().NodeId(id).Delete()
}

func (r *IntervalRepository) ListByContextId(contextId string) ([]*core.Interval, error) {
	return r.repository.Query().KindEquals(core.IntervalType).ParentId(contextId).KV().List()
}

func (r *IntervalRepository) GetActiveIntervalByContextId(contextId string) (*core.Interval, error) {
	return r.repository.Query().KindEquals(core.IntervalType).ParentId(contextId).StatusEquals("active").KV().First()
}

func (r *IntervalRepository) ListByDay(date time.Time) ([]*core.Interval, error) {
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	dayEnd := dayStart.Add(24 * time.Hour)

	key := "start"
	all, err := r.repository.Query().
		KindEquals(core.IntervalType).
		KVFilter(&nod.KVFilter{Key: &key, TimeTo: &dayEnd}).
		KV().
		List()
	if err != nil {
		return nil, err
	}

	result := make([]*core.Interval, 0, len(all))
	for _, interval := range all {
		if interval.End.Time.IsZero() {
			result = append(result, interval)
			continue
		}
		if !interval.End.Time.Before(dayStart) {
			result = append(result, interval)
		}
	}
	return result, nil
}

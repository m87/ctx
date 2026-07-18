package bootstrap

import (
	"sort"
	"time"

	"github.com/m87/ctx/core"
	"github.com/m87/nod"
)

type IntervalRepository struct {
	scope *nod.NodeScope[core.Interval]
}

func NewIntervalRepository(repository *nod.Repository) *IntervalRepository {
	return &IntervalRepository{scope: nod.Nodes[core.Interval](repository)}
}

func (r *IntervalRepository) GetById(id string) (*core.Interval, error) {
	return r.scope.Query().
		Where(nod.NodeFields.Id.Equals(id)).
		WithKV().
		FindFirst()
}

func (r *IntervalRepository) Save(interval *core.Interval) (string, error) {
	return r.scope.SaveNode(interval)
}

func (r *IntervalRepository) Delete(id string) error {
	return r.scope.Query().
		Where(nod.NodeFields.Id.Equals(id)).
		DeleteAll()
}

func (r *IntervalRepository) DeleteByContextId(contextId string) error {
	return r.scope.Query().
		Where(nod.NodeFields.Kind.Equals(core.IntervalType)).
		Where(nod.NodeFields.ParentId.Equals(contextId)).
		DeleteAll()
}

func (r *IntervalRepository) ListByContextId(contextId string) ([]*core.Interval, error) {
	intervals, err := r.scope.Query().
		Where(nod.NodeFields.Kind.Equals(core.IntervalType)).
		Where(nod.NodeFields.ParentId.Equals(contextId)).
		WithKV().
		FindAll()
	if err != nil {
		return nil, err
	}

	sort.SliceStable(intervals, func(i, j int) bool {
		if intervals[i].Start.Time.Equal(intervals[j].Start.Time) {
			return intervals[i].End.Time.After(intervals[j].End.Time)
		}
		return intervals[i].Start.Time.After(intervals[j].Start.Time)
	})

	return intervals, nil
}

func (r *IntervalRepository) GetActiveIntervalByContextId(contextId string) (*core.Interval, error) {
	return r.scope.Query().
		Where(nod.NodeFields.Kind.Equals(core.IntervalType)).
		Where(nod.NodeFields.ParentId.Equals(contextId)).
		Where(nod.NodeFields.Status.Equals("active")).
		WithKV().
		FindFirst()
}

func (r *IntervalRepository) ListByDay(date time.Time, workspaceId string) ([]*core.Interval, error) {
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	dayEnd := dayStart.Add(24 * time.Hour)

	all, err := r.scope.Query().
		Where(nod.NodeFields.Kind.Equals(core.IntervalType)).
		Where(nod.NodeFields.NamespaceId.Equals(workspaceId)).
		Where(nod.KvTime("start").LessThanOrEqual(dayEnd)).
		WithKV().
		FindAll()
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

func (r *IntervalRepository) List() ([]*core.Interval, error) {
	return r.scope.Query().
		Where(nod.NodeFields.Kind.Equals(core.IntervalType)).
		WithKV().
		FindAll()
}

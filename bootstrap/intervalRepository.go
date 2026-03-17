package bootstrap

import (
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

package bootstrap

import (
	"github.com/m87/ctx/core"
	"github.com/m87/nod"
)

type ContextRepository struct {
	repository nod.TypedRepository[core.Context]
}

func NewContextRepository(repository *nod.Repository) *ContextRepository {
	return &ContextRepository{
		repository: nod.As[core.Context](repository),
	}
}

func (r *ContextRepository) GetById(id string) (*core.Context, error) {
	return r.repository.Query().NodeId(id).KV().Tags().Content().First()
}

func (r *ContextRepository) Save(context *core.Context) (string, error) {
	return r.repository.Save(context)
}

func (r *ContextRepository) Delete(id string) error {
	return r.repository.Query().NodeId(id).Delete()
}

func (r *ContextRepository) List() ([]*core.Context, error) {
	return r.repository.Query().KindEquals(core.ContextType).List()
}

func (r *ContextRepository) GetActive() (*core.Context, error) {
	return r.repository.Query().KindEquals(core.ContextType).StatusEquals("active").KV().First()
}

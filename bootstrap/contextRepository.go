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
	return r.repository.Query().KindEquals(core.ContextType).KV().List()
}

func (r *ContextRepository) GetActive() (*core.Context, error) {
	return r.repository.Query().KindEquals(core.ContextType).StatusEquals("active").KV().First()
}

func (r *ContextRepository) ListByWorkspace(workspaceId string) ([]*core.Context, error) {
	contexts, err := r.ListByWorkspaceIncludingArchived(workspaceId)
	if err != nil {
		return nil, err
	}

	activeContexts := make([]*core.Context, 0, len(contexts))
	for _, context := range contexts {
		if context != nil && !context.Archived {
			activeContexts = append(activeContexts, context)
		}
	}

	return activeContexts, nil
}

func (r *ContextRepository) ListByWorkspaceIncludingArchived(workspaceId string) ([]*core.Context, error) {
	return r.repository.Query().KindEquals(core.ContextType).NamespaceId(workspaceId).KV().Content().Tags().List()
}

package bootstrap

import (
	"github.com/m87/ctx/core"
	"github.com/m87/nod"
)

type ContextRepository struct {
	scope *nod.NodeScope[core.Context]
}

func NewContextRepository(repository *nod.Repository) *ContextRepository {
	return &ContextRepository{scope: nod.Nodes[core.Context](repository)}
}

func (r *ContextRepository) GetById(id string) (*core.Context, error) {
	return r.scope.Query().
		Where(nod.NodeFields.Id.Equals(id)).
		WithKV().
		WithTags().
		WithContent().
		FindFirst()
}

func (r *ContextRepository) Save(context *core.Context) (string, error) {
	return r.scope.SaveNode(context)
}

func (r *ContextRepository) Delete(id string) error {
	return r.scope.Query().
		Where(nod.NodeFields.Id.Equals(id)).
		DeleteAll()
}

func (r *ContextRepository) List() ([]*core.Context, error) {
	return r.scope.Query().
		Where(nod.NodeFields.Kind.Equals(core.ContextType)).
		WithKV().
		FindAll()
}

func (r *ContextRepository) GetActive() (*core.Context, error) {
	return r.scope.Query().
		Where(nod.NodeFields.Kind.Equals(core.ContextType)).
		Where(nod.NodeFields.Status.Equals("active")).
		WithKV().
		FindFirst()
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
	return r.scope.Query().
		Where(nod.NodeFields.Kind.Equals(core.ContextType)).
		Where(nod.NodeFields.NamespaceId.Equals(workspaceId)).
		WithKV().
		WithContent().
		WithTags().
		FindAll()
}

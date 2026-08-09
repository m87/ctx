package bootstrap

import (
	"github.com/m87/ctx/core"
	"github.com/m87/nod"
	"gorm.io/gorm"
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

func (r *ContextRepository) ListByProject(projectId string) ([]*core.Context, error) {
	return r.scope.Query().
		Where(nod.NodeFields.Kind.Equals(core.ContextType)).
		Where(nod.NodeFields.ParentId.Equals(projectId)).
		WithKV().
		WithContent().
		WithTags().
		FindAll()
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
	context, err := r.scope.Query().
		Where(nod.NodeFields.Kind.Equals(core.ContextType)).
		Where(nod.NodeFields.Status.Equals("active")).
		WithKV().
		FindFirst()

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return context, err
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

func (r *ContextRepository) ListToSync(limit int) ([]*core.Context, error) {
	return r.scope.Query().
		Where(nod.NodeFields.Kind.Equals(core.ContextType)).
		Where(nod.Or(nod.Kv("synced").NotExists(), nod.KvBool("synced").Equals(false))).
		WithKV().
		WithTags().
		WithContent().
		FindAll()
}

func (r *ContextRepository) SaveAll(contexts []*core.Context) ([]string, error) {
	return r.scope.SaveNodes(contexts)
}

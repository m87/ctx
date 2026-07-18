package bootstrap

import (
	"github.com/m87/ctx/core"
	"github.com/m87/nod"
)

type WorkspaceRepository struct {
	scope *nod.NodeScope[core.Workspace]
}

func NewWorkspaceRepository(repository *nod.Repository) *WorkspaceRepository {
	return &WorkspaceRepository{scope: nod.Nodes[core.Workspace](repository)}
}

func (r *WorkspaceRepository) GetById(id string) (*core.Workspace, error) {
	return r.scope.Query().
		Where(nod.NodeFields.Id.Equals(id)).
		WithContent().
		FindFirst()
}

func (r *WorkspaceRepository) Save(workspace *core.Workspace) (string, error) {
	return r.scope.SaveNode(workspace)
}

func (r *WorkspaceRepository) Delete(id string) error {
	return r.scope.Query().
		Where(nod.NodeFields.Id.Equals(id)).
		DeleteAll()
}

func (r *WorkspaceRepository) List() ([]*core.Workspace, error) {
	return r.scope.Query().
		Where(nod.NodeFields.Kind.Equals(core.WorkspaceType)).
		WithContent().
		FindAll()
}

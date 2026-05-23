package bootstrap

import (
	"github.com/m87/ctx/core"
	"github.com/m87/nod"
)

type WorkspaceRepository struct {
	repository nod.TypedRepository[core.Workspace]
}

func NewWorkspaceRepository(repository *nod.Repository) *WorkspaceRepository {
	return &WorkspaceRepository{
		repository: nod.As[core.Workspace](repository),
	}
}

func (r *WorkspaceRepository) GetById(id string) (*core.Workspace, error) {
	return r.repository.Query().NodeId(id).First()
}

func (r *WorkspaceRepository) Save(workspace *core.Workspace) (string, error) {
	return r.repository.Save(workspace)
}

func (r *WorkspaceRepository) Delete(id string) error {
	return r.repository.Query().NodeId(id).Delete()
}

func (r *WorkspaceRepository) List() ([]*core.Workspace, error) {
	return r.repository.Query().KindEquals(core.WorkspaceType).List()
}

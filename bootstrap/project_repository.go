package bootstrap

import (
	"github.com/m87/ctx/core"
	"github.com/m87/nod"
)

type ProjectRepository struct {
	scope *nod.NodeScope[core.Project]
}

func NewProjectRepository(repository *nod.Repository) *ProjectRepository {
	return &ProjectRepository{scope: nod.Nodes[core.Project](repository)}
}

func (r *ProjectRepository) GetById(id string) (*core.Project, error) {
	return r.scope.Query().
		Where(nod.NodeFields.Id.Equals(id)).
		FindFirst()
}

func (r *ProjectRepository) ListChildren(parentId string) ([]*core.Project, error) {
	return r.scope.Query().
		Where(nod.NodeFields.Kind.Equals(core.ProjectType)).
		Where(nod.NodeFields.ParentId.Equals(parentId)).
		FindAll()
}

func (r *ProjectRepository) Save(project *core.Project) (string, error) {
	return r.scope.SaveNode(project)
}

func (r *ProjectRepository) SaveAll(projects []*core.Project) ([]string, error) {
	return r.scope.SaveNodes(projects)
}

func (r *ProjectRepository) Delete(id string) error {
	return r.scope.Query().
		Where(nod.NodeFields.Id.Equals(id)).
		DeleteAll()
}

func (r *ProjectRepository) List(workspaceId string) ([]*core.Project, error) {
	return r.scope.Query().
		Where(nod.NodeFields.Kind.Equals(core.ProjectType)).
		Where(nod.NodeFields.NamespaceId.Equals(workspaceId)).
		Where(nod.NodeFields.ParentId.IsNil()).
		FindAll()
}

func (r *ProjectRepository) ListIncludingArchived(workspaceId string) ([]*core.Project, error) {
	return r.scope.Query().
		Where(nod.NodeFields.Kind.Equals(core.ProjectType)).
		Where(nod.NodeFields.NamespaceId.Equals(workspaceId)).
		FindAll()
}

func (r *ProjectRepository) ListToSync(limit int) ([]*core.Project, error) {
	return r.scope.Query().
		Where(nod.NodeFields.Kind.Equals(core.ProjectType)).
		FindAll()
}

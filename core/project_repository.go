package core

type ProjectRepository interface {
	Save(project *Project) (string, error)
	SaveAll(projects []*Project) ([]string, error)
	List(workspaceId string) ([]*Project, error)
	ListToSync(limit int) ([]*Project, error)
	ListIncludingArchived(workspaceId string) ([]*Project, error)
	ListChildren(parentId string) ([]*Project, error)
	Delete(id string) error
	GetById(id string) (*Project, error)
}

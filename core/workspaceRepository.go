package core

type WorkspaceRepository interface {
	GetById(id string) (*Workspace, error)
	Save(workspace *Workspace) (string, error)
	Delete(id string) error
	List() ([]*Workspace, error)
}

package core

type WorkspaceRepository interface {
	GetById(id string) (*Workspace, error)
	Save(workspace *Workspace) error
	Delete(id string) error
}

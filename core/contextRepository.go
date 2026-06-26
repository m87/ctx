package core

type ContextRepository interface {
	GetById(id string) (*Context, error)
	Save(context *Context) (string, error)
	Delete(id string) error
	List() ([]*Context, error)
	ListByWorkspace(workspaceId string) ([]*Context, error)
	GetActive() (*Context, error)
}

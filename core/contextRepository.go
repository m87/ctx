package core

type ContextRepository interface {
	GetById(id string) (*Context, error)
	Save(context *Context) error
	Delete(id string) error
}

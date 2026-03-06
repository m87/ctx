package core

type ProjectRepository interface {
	GetById(id string) (*Project, error)
	Save(project *Project) error
	Delete(id string) error
}

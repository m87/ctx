package core

type Archiver[T any] interface {
	Archive(data []T, session Session) error
	Update(data []T, session Session) error
}

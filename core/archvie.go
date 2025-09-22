package core

type Archiver[T any] interface {
	Archvie(data []T) error
}

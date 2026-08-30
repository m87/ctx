package storage

type ErrInvalidStoragePath struct{}

func (e ErrInvalidStoragePath) Error() string {
	return "invalid storage path"
}

func NewErrInvalidStoragePath() error {
	return ErrInvalidStoragePath{}
}

type ErrForeignKeyDisabled struct{}

func (e ErrForeignKeyDisabled) Error() string {
	return "foreign key support is disabled"
}

func NewErrForeignKeyDisabled() error {
	return ErrForeignKeyDisabled{}
}

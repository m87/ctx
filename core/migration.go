package core

type Migrator interface {
	Id() string
	Migrate() error
	MigrateArchive(archiver Archiver[Context]) error
}

type MigrationManager interface {
	CallMigrationChain(fromVersion Version, toVersion Version, registry MigrationRegistry) error
}

type MigrationRegistry interface {
	RegisterMigration(migrator Migrator) error
	MigrationExecuted(migrator Migrator) bool
	Save() error
}

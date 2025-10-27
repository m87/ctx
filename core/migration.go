package core

type Migrator interface {
	Migrate(session Session) error
	MigrateArchive(archiver Archiver[Context]) error
}

type MigrationManager interface {
	CreateMigrationChain(fromVersion Version, toVersion Version) []Migrator
}

package core

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Migrator interface {
	Id() string
	Migrate() error
	MigrateArchive(archiver Archiver[Context]) error
}

type MigrationManager interface {
	CreateMigrationMap(fromVersion Version, toVersion Version) map[Version]Migrator
}

type MigrationRegistry interface {
	RegisterMigration(migrator Migrator) error
	MigrationExecuted(migrator Migrator) bool
	Save() error
}

func ParseVersion(versionStr string) Version {
	if Release == "dev" {
		return Version{Major: 99999, Minor: 99999, Patch: 99999}
	} else {
		var major, minor, patch int
		fmt.Sscanf(versionStr, "%d.%d.%d", &major, &minor, &patch)
		return Version{Major: major, Minor: minor, Patch: patch}
	}

}

func Migrate(manager MigrationManager, registry MigrationRegistry) error {
	log.Println("Migration process started...")

	localVersion := viper.GetString("version")
	err := callMigrationChain(ParseVersion(localVersion), ParseVersion(Release), registry, manager)

	if err != nil {
		return err
	}

	log.Println("Migration completed")
	return nil
}

func callMigrationChain(fromVersion Version, toVersion Version, registry MigrationRegistry, manager MigrationManager) error {

	migrations := manager.CreateMigrationMap(fromVersion, toVersion)

	chain := make([]Version, 0, len(migrations))
	for v := range migrations {
		chain = append(chain, v)
	}
	Sort(chain)

	migrationsToApply := []Migrator{}
	for _, version := range chain {
		if version.Major < fromVersion.Major ||
			(version.Major == fromVersion.Major && version.Minor < fromVersion.Minor) ||
			(version.Major == fromVersion.Major && version.Minor == fromVersion.Minor && version.Patch < fromVersion.Patch) {
			continue
		}
		if version.Major > toVersion.Major ||
			(version.Major == toVersion.Major && version.Minor > toVersion.Minor) ||
			(version.Major == toVersion.Major && version.Minor == toVersion.Minor && version.Patch > toVersion.Patch) {
			break
		}
		migrationsToApply = append(migrationsToApply, migrations[version])
	}

	if len(migrationsToApply) == 0 {
		log.Println("No migrations to apply")
		return nil
	}

	for _, migrator := range migrationsToApply {
		if registry.MigrationExecuted(migrator) {
			log.Println("Skipping already executed migration:", migrator.Id())
			continue
		}
		err := migrator.Migrate()
		if err != nil {
			return err
		}
		registry.RegisterMigration(migrator)
	}

	err := registry.Save()
	if err != nil {
		return err
	}

	return nil
}

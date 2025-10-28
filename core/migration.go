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
	CallMigrationChain(fromVersion Version, toVersion Version, registry MigrationRegistry) error
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
	err := manager.CallMigrationChain(ParseVersion(localVersion), ParseVersion(Release), registry)

	if err != nil {
		return err
	}

	log.Println("Migration completed")
	return nil
}

package localstorage

import (
	"log"

	"github.com/m87/ctx/core"
	"golang.org/x/exp/maps"
)

type LocalStorageMigrationManager struct {
	statePath   string
	archivePath string
}

type LocalStorageMigratorV2_1_0_V_2_2_0 struct {
	statePath   string
	archivePath string
}

func (manager *LocalStorageMigrationManager) CreateMigrationChain(fromVersion core.Version, toVersion core.Version) []core.Migrator {
	migrations := map[core.Version]core.Migrator{
		core.Version{Major: 2, Minor: 1, Patch: 0}: &LocalStorageMigratorV2_1_0_V_2_2_0{statePath: manager.statePath, archivePath: manager.archivePath},
	}

	chain := maps.Keys(migrations)
	core.Sort(chain)

	migrationsToApply := []core.Migrator{}
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

	return migrationsToApply

}

func (migrator *LocalStorageMigratorV2_1_0_V_2_2_0) Migrate(session core.Session) error {
	log.Println(`Migration v2.1.0 -> v2.2.0
	Migration plan:
	- Convert context comments to objects with ids
	`)

	return nil
}

func (migrator *LocalStorageMigratorV2_1_0_V_2_2_0) MigrateArchive(archiver core.Archiver[core.Context]) error {
	log.Println(`Archive migration v2.1.0 -> v2.2.0
	Migration plan:
	- Convert context comments to objects with ids
	`)

	return nil
}

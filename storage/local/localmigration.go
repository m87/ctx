package localstorage

import (
	"github.com/m87/ctx/core"
	"golang.org/x/exp/maps"
)

type LocalStorageMigrationManager struct {
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

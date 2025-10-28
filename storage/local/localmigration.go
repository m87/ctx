package localstorage

import (
	"encoding/json"
	"log"
	"os"

	"github.com/m87/ctx/core"
	"golang.org/x/exp/maps"
)

type LocalMigrationRegistry struct {
	path               string
	executedMigrations map[string]bool
}

type LocalStorageMigrationManager struct {
	statePath   string
	archivePath string
}

func (manager *LocalStorageMigrationManager) CallMigrationChain(fromVersion core.Version, toVersion core.Version, registry core.MigrationRegistry) error {

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

func LoadLocalMigrationRegistry(path string) (*LocalMigrationRegistry, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NewLocalMigrationRegistry(path), nil
		}
		return nil, err
	}
	defer f.Close()

	var executedMigrations map[string]bool
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&executedMigrations); err != nil {
		return nil, err
	}
	return &LocalMigrationRegistry{
		path:               path,
		executedMigrations: executedMigrations,
	}, nil
}

func (registry *LocalMigrationRegistry) Save() error {
	f, err := os.Create(registry.path)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	return encoder.Encode(registry.executedMigrations)
}

func NewLocalMigrationRegistry(path string) *LocalMigrationRegistry {
	return &LocalMigrationRegistry{
		path:               path,
		executedMigrations: make(map[string]bool),
	}
}

func (registry *LocalMigrationRegistry) RegisterMigration(migrator core.Migrator) error {
	registry.executedMigrations[migrator.Id()] = true
	return nil
}

func (registry *LocalMigrationRegistry) MigrationExecuted(migrator core.Migrator) bool {
	_, executed := registry.executedMigrations[migrator.Id()]
	return executed
}

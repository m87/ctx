package localstorage

import (
	"encoding/json"
	"os"

	"github.com/m87/ctx/core"
)

type LocalMigrationRegistry struct {
	path               string
	executedMigrations map[string]bool
}

type LocalStorageMigrationManager struct {
	statePath   string
	archivePath string
}

func (manager *LocalStorageMigrationManager) CreateMigrationMap(fromVersion core.Version, toVersion core.Version) map[core.Version]core.Migrator {
	return map[core.Version]core.Migrator{
		core.Version{Major: 2, Minor: 1, Patch: 0}: &LocalStorageMigratorV2_1_0_V_2_2_0{statePath: manager.statePath, archivePath: manager.archivePath},
	}
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

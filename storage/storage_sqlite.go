package storage

import (
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

func NewSqliteStorage(path string) (*Storage, error) {
	db, err := initSqliteStorage(path)
	if err != nil {
		return nil, err
	}

	return &Storage{
		DB: db,
	}, nil
}

func initSqliteStorage(path string) (*gorm.DB, error) {
	if strings.TrimSpace(path) == "" {
		return nil, NewErrInvalidStoragePath()
	}

	db, err := gorm.Open(sqlite.New(sqlite.Config{
		DSN:                  path,
		DriverName: 				 "sqlite",
	}), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	if err := configureSqlite(db, path); err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&TagEntity{},
		&ContextTagEntity{},
		&WorkspaceEntity{},
		&Properties{},
		&ProjectEntity{},
		&ContextEntity{},
		&IntervalEntity{},
	); err != nil {
		return nil, err
	}

	return db, nil
}

func configureSqlite(db *gorm.DB, path string) error {
	sqliteDB, err := db.DB()
	if err != nil {
		return err
	}

	if path == ":memory:" || strings.Contains(path, "mode=memory") {
		sqliteDB.SetMaxOpenConns(1)
		sqliteDB.SetMaxIdleConns(1)
	}

	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		return err
	}

	var fkEnabeld int
	if err := db.Raw("PRAGMA foreign_keys").Scan(&fkEnabeld).Error; err != nil {
		return err
	}
	
	if fkEnabeld != 1 {
		return NewErrForeignKeyDisabled()
	}

	return nil
}

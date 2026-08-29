package storage

import (
	"strings"
	"time"

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
		DSN:        path,
		DriverName: "sqlite",
	}), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	if err := configureSqlite(db, path); err != nil {
		return nil, err
	}

	if err := Migrate(db); err != nil {
		return nil, err
	}

	db.Config.NowFunc = func() time.Time {
		return time.Now().UTC()
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.SetupJoinTable(&ContextEntity{}, "Tags", &ContextTagEntity{}); err != nil {
			return err
		}

		if err := tx.AutoMigrate(
			&TagEntity{},
			&ContextTagEntity{},
			&WorkspaceEntity{},
			&Properties{},
			&ClientProperties{},
			&ProjectEntity{},
			&ContextEntity{},
			&IntervalEntity{},
		); err != nil {
			return err
		}

		return nil
	})
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

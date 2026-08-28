package storage



func CreateTestInMemoryStorage() (*Storage, error) {
	db, err := initSqliteStorage(":memory:")
	if err != nil {
		return nil, err
	}

	return &Storage{
		DB: db,
	}, nil
}

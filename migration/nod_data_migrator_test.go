package migration

import (
	"bytes"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/m87/ctx/core"
	ctxlog "github.com/m87/ctx/log"
	"github.com/m87/ctx/storage"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestNodDataMigratorMigratesLegacyTables(t *testing.T) {
	var logOutput bytes.Buffer
	originalLogger := ctxlog.Logger
	ctxlog.Logger = slog.New(slog.NewTextHandler(&logOutput, nil))
	t.Cleanup(func() { ctxlog.Logger = originalLogger })

	storageDB, err := storage.NewSqliteStorage(":memory:")
	require.NoError(t, err)
	createLegacyTables(t, storageDB.DB)
	seedLegacyData(t, storageDB.DB)

	migrator := NewNodDataMigrator(storageDB.DB)
	require.NoError(t, migrator.MigrateIfNeeded())

	properties, err := storage.NewPropertiesRepository(storageDB.DB).Load()
	require.NoError(t, err)
	require.Equal(t, "0.7.1", properties.DatabaseVersion)
	require.NotEmpty(t, properties.ClientId)

	workspace, err := storage.NewWorkspaceRepository(storageDB.DB).GetById("workspace-1")
	require.NoError(t, err)
	require.Equal(t, &core.Workspace{
		Id:          "workspace-1",
		Name:        "Workspace",
		Description: "Workspace description",
	}, workspace)

	project, err := storage.NewProjectRepository(storageDB.DB).GetById("project-1")
	require.NoError(t, err)
	require.Equal(t, &core.Project{
		Id:          "project-1",
		Name:        "Project",
		WorkspaceId: "workspace-1",
	}, project)

	context, err := storage.NewContextRepository(storageDB.DB).GetById("context-1")
	require.NoError(t, err)
	require.Equal(t, "Context", context.Name)
	require.Equal(t, "workspace-1", context.WorkspaceId)
	require.Equal(t, "inactive", context.Status)
	require.True(t, context.Archived)
	require.Equal(t, "Context description", context.Description)
	require.NotNil(t, context.ProjectId)
	require.Equal(t, "project-1", *context.ProjectId)
	require.Equal(t, &core.ProjectMetadata{Id: "project-1", Name: "Project"}, context.Project)
	require.Equal(t, []*core.Tag{{Id: "tag-1", Name: "important"}}, context.Tags)

	interval, err := storage.NewIntervalRepository(storageDB.DB).GetById("interval-1")
	require.NoError(t, err)
	require.Equal(t, "context-1", interval.ContextId)
	require.Equal(t, "workspace-1", interval.WorkspaceId)
	require.Equal(t, "completed", interval.Status)
	require.Equal(t, 45*time.Minute, interval.Duration)
	require.Equal(t, time.Date(2026, time.August, 29, 18, 0, 0, 0, time.UTC), *interval.Start)
	require.Equal(t, time.Date(2026, time.August, 29, 18, 45, 0, 0, time.UTC), *interval.End)

	settings, err := storage.NewClientPropertiesRepository(storageDB.DB).Load()
	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"client.general.theme":    "dark",
		"client.general.firstDay": "Sunday",
		"client.general.timeZone": "Europe/Warsaw",
	}, settings.Values())

	for _, table := range []string{"node_cores", "node_kvs", "node_contents", "tags", "node_tags"} {
		require.False(t, storageDB.DB.Migrator().HasTable(table), "legacy table %s must be removed", table)
	}
	require.Contains(t, logOutput.String(), "Starting legacy nod data migration")
	require.Contains(t, logOutput.String(), "Legacy nod data migration completed")
	require.Contains(t, logOutput.String(), "Legacy nod tables removed")

	require.NoError(t, migrator.MigrateIfNeeded())
	migratedAgain, err := storage.NewWorkspaceRepository(storageDB.DB).GetById("workspace-1")
	require.NoError(t, err)
	require.Equal(t, "Workspace", migratedAgain.Name)
}

func TestNodDataMigratorUsesEmptyNewVersionAsMigrationFlag(t *testing.T) {
	t.Run("missing properties row initializes a clean database", func(t *testing.T) {
		storageDB, err := storage.NewSqliteStorage(":memory:")
		require.NoError(t, err)

		require.NoError(t, NewNodDataMigrator(storageDB.DB).MigrateIfNeeded())

		properties, err := storage.NewPropertiesRepository(storageDB.DB).Load()
		require.NoError(t, err)
		require.Equal(t, core.CurrentDatabaseVersion, properties.DatabaseVersion)
		require.NotEmpty(t, properties.ClientId)
	})

	t.Run("empty new version imports data and preserves client ID", func(t *testing.T) {
		storageDB, err := storage.NewSqliteStorage(":memory:")
		require.NoError(t, err)
		createLegacyTables(t, storageDB.DB)
		insertLegacyNode(t, storageDB.DB, legacyNodeCore{
			Id:   "legacy-workspace",
			Kind: "workspace",
			Name: "Legacy workspace",
		})
		require.NoError(t, storage.NewPropertiesRepository(storageDB.DB).Save(&core.SystemInfo{
			ClientId: "existing-client",
		}))

		require.NoError(t, NewNodDataMigrator(storageDB.DB).MigrateIfNeeded())

		workspace, err := storage.NewWorkspaceRepository(storageDB.DB).GetById("legacy-workspace")
		require.NoError(t, err)
		require.Equal(t, "Legacy workspace", workspace.Name)
		properties, err := storage.NewPropertiesRepository(storageDB.DB).Load()
		require.NoError(t, err)
		require.Equal(t, core.CurrentDatabaseVersion, properties.DatabaseVersion)
		require.Equal(t, "existing-client", properties.ClientId)
	})

	t.Run("non-empty new version skips import but removes legacy tables", func(t *testing.T) {
		storageDB, err := storage.NewSqliteStorage(":memory:")
		require.NoError(t, err)
		createLegacyTables(t, storageDB.DB)
		insertLegacyNode(t, storageDB.DB, legacyNodeCore{
			Id:   "legacy-workspace",
			Kind: "workspace",
			Name: "Legacy workspace",
		})
		require.NoError(t, storage.NewPropertiesRepository(storageDB.DB).Save(&core.SystemInfo{
			DatabaseVersion: "0.6.0",
			ClientId:        "client-1",
		}))

		require.NoError(t, NewNodDataMigrator(storageDB.DB).MigrateIfNeeded())

		workspace, err := storage.NewWorkspaceRepository(storageDB.DB).GetById("legacy-workspace")
		require.NoError(t, err)
		require.Nil(t, workspace)
		properties, err := storage.NewPropertiesRepository(storageDB.DB).Load()
		require.NoError(t, err)
		require.Equal(t, "0.6.0", properties.DatabaseVersion)
		require.False(t, storageDB.DB.Migrator().HasTable("node_cores"))
	})
}

func TestNodDataMigratorRollsBackDataAndVersionOnFailure(t *testing.T) {
	storageDB, err := storage.NewSqliteStorage(":memory:")
	require.NoError(t, err)
	createLegacyTables(t, storageDB.DB)
	insertLegacyNode(t, storageDB.DB,
		legacyNodeCore{Id: "workspace-1", Kind: "workspace", Name: "Workspace"},
		legacyNodeCore{
			Id:          "interval-without-start",
			Kind:        "interval",
			Name:        "interval-without-start",
			Status:      "completed",
			NamespaceId: stringPointer("workspace-1"),
			ParentId:    stringPointer("context-1"),
		},
	)

	err = NewNodDataMigrator(storageDB.DB).MigrateIfNeeded()
	require.Error(t, err)

	var workspaceCount int64
	require.NoError(t, storageDB.DB.Model(&storage.WorkspaceEntity{}).Count(&workspaceCount).Error)
	require.Zero(t, workspaceCount)
	_, err = storage.NewPropertiesRepository(storageDB.DB).Load()
	require.True(t, errors.Is(err, gorm.ErrRecordNotFound))
	require.True(t, storageDB.DB.Migrator().HasTable("node_cores"))
}

func createLegacyTables(t *testing.T, db *gorm.DB) {
	t.Helper()
	statements := []string{
		`CREATE TABLE node_cores (
			id TEXT PRIMARY KEY,
			namespace_id TEXT,
			parent_id TEXT,
			kind TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT '',
			name TEXT NOT NULL
		)`,
		`CREATE TABLE node_kvs (
			node_id TEXT NOT NULL,
			key TEXT NOT NULL,
			value_text TEXT,
			value_int64 INTEGER,
			value_bool BOOLEAN,
			value_time DATETIME,
			PRIMARY KEY (node_id, key)
		)`,
		`CREATE TABLE node_contents (
			node_id TEXT NOT NULL,
			key TEXT NOT NULL,
			value TEXT,
			PRIMARY KEY (node_id, key)
		)`,
		`CREATE TABLE tags (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL
		)`,
		`CREATE TABLE node_tags (
			node_id TEXT NOT NULL,
			tag_id TEXT NOT NULL,
			PRIMARY KEY (node_id, tag_id)
		)`,
	}
	for _, statement := range statements {
		require.NoError(t, db.Exec(statement).Error)
	}
}

func seedLegacyData(t *testing.T, db *gorm.DB) {
	t.Helper()
	insertLegacyNode(t, db,
		legacyNodeCore{Id: "workspace-1", Kind: "workspace", Name: "Workspace"},
		legacyNodeCore{
			Id:          "project-1",
			Kind:        "project",
			Name:        "Project",
			NamespaceId: stringPointer("workspace-1"),
		},
		legacyNodeCore{
			Id:          "context-1",
			Kind:        "context",
			Name:        "Context",
			Status:      "inactive",
			NamespaceId: stringPointer("workspace-1"),
		},
		legacyNodeCore{
			Id:          "interval-1",
			Kind:        "interval",
			Name:        "interval-1",
			Status:      "completed",
			NamespaceId: stringPointer("workspace-1"),
			ParentId:    stringPointer("context-1"),
		},
		legacyNodeCore{Id: "settingsV1", Kind: "settings", Name: "settingsV1"},
	)

	start := time.Date(2026, time.August, 29, 20, 0, 0, 0, time.FixedZone("CEST", 2*60*60))
	end := start.Add(45 * time.Minute)
	duration := int64(45 * time.Minute)
	archived := true
	projectId := "project-1"
	theme := "dark"
	firstDay := "Sunday"
	timeZone := "Europe/Warsaw"
	require.NoError(t, db.Table("node_kvs").Create(&[]legacyNodeKV{
		{NodeId: "context-1", Key: "archived", ValueBool: &archived},
		{NodeId: "context-1", Key: "projectId", ValueText: &projectId},
		{NodeId: "interval-1", Key: "start", ValueTime: &start},
		{NodeId: "interval-1", Key: "end", ValueTime: &end},
		{NodeId: "interval-1", Key: "duration", ValueInt64: &duration},
		{NodeId: "settingsV1", Key: "client.general.theme", ValueText: &theme},
		{NodeId: "settingsV1", Key: "client.general.firstDay", ValueText: &firstDay},
		{NodeId: "settingsV1", Key: "client.general.timeZone", ValueText: &timeZone},
	}).Error)

	workspaceDescription := "Workspace description"
	contextDescription := "Context description"
	require.NoError(t, db.Table("node_contents").Create(&[]legacyNodeContent{
		{NodeId: "workspace-1", Key: "description", Value: &workspaceDescription},
		{NodeId: "context-1", Key: "description", Value: &contextDescription},
	}).Error)

	require.NoError(t, db.Table("tags").Create(&legacyTag{Id: "tag-1", Name: "important"}).Error)
	require.NoError(t, db.Table("node_tags").Create(&legacyNodeTag{NodeId: "context-1", TagId: "tag-1"}).Error)
}

func insertLegacyNode(t *testing.T, db *gorm.DB, nodes ...legacyNodeCore) {
	t.Helper()
	require.NoError(t, db.Table("node_cores").Create(&nodes).Error)
}

func stringPointer(value string) *string {
	return &value
}

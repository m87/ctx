package migration

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/m87/ctx/core"
	ctxlog "github.com/m87/ctx/log"
	"github.com/m87/ctx/storage"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	legacyNodeCoreTable = "node_cores"
	legacyTagTable      = "tags"
	legacyNodeTagTable  = "node_tags"
)

var (
	legacyKVTables      = []string{"node_kvs", "kvs", "kv", "node_kv"}
	legacyContentTables = []string{"node_contents", "contents", "content", "node_content"}
	legacyNodTables     = []string{
		"edge_tags",
		"edge_kvs",
		"edge_contents",
		"edge_cores",
		legacyNodeTagTable,
		"node_kvs",
		"kvs",
		"kv",
		"node_kv",
		"node_contents",
		"contents",
		"content",
		"node_content",
		legacyTagTable,
		legacyNodeCoreTable,
		"nod_properties",
	}
)

type NodDataMigrator struct {
	db *gorm.DB
}

type legacyNodeCore struct {
	Id          string
	NamespaceId *string
	ParentId    *string
	Kind        string
	Status      string
	Name        string
}

type legacyNodeKV struct {
	NodeId     string
	Key        string
	ValueText  *string
	ValueInt64 *int64
	ValueBool  *bool
	ValueTime  *time.Time
}

type legacyNodeContent struct {
	NodeId string
	Key    string
	Value  *string
}

type legacyTag struct {
	Id   string
	Name string
}

type legacyNodeTag struct {
	NodeId string
	TagId  string
}

type legacyData struct {
	nodes      []legacyNodeCore
	kv         map[string]map[string]legacyNodeKV
	content    map[string]map[string]legacyNodeContent
	tags       map[string]legacyTag
	contextTag []legacyNodeTag
}

func NewNodDataMigrator(db *gorm.DB) *NodDataMigrator {
	return &NodDataMigrator{db: db}
}

func (m *NodDataMigrator) MigrateIfNeeded() error {
	if m == nil || m.db == nil {
		return fmt.Errorf("database is required")
	}

	var migrationPerformed bool
	var removedTables []string
	err := m.db.Transaction(func(tx *gorm.DB) error {
		propertiesRepository := storage.NewPropertiesRepository(tx)
		systemInfo, err := propertiesRepository.Load()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			systemInfo = &core.SystemInfo{}
		} else if err != nil {
			return fmt.Errorf("load database version: %w", err)
		}

		legacyTables := existingLegacyTables(tx)
		if strings.TrimSpace(systemInfo.DatabaseVersion) != "" {
			if err := dropLegacyTables(tx, legacyTables); err != nil {
				return err
			}
			removedTables = legacyTables
			return nil
		}

		if len(legacyTables) > 0 {
			migrationPerformed = true
			ctxlog.Logger.Info("Starting legacy nod data migration",
				"target_database_version", core.CurrentDatabaseVersion,
			)
		}

		if err := (&NodDataMigrator{db: tx}).migrateLegacyData(); err != nil {
			return err
		}
		if err := dropLegacyTables(tx, legacyTables); err != nil {
			return err
		}
		removedTables = legacyTables

		if strings.TrimSpace(systemInfo.ClientId) == "" {
			systemInfo.ClientId = uuid.NewString()
		}
		systemInfo.DatabaseVersion = core.CurrentDatabaseVersion
		if err := propertiesRepository.Save(systemInfo); err != nil {
			return fmt.Errorf("save database version: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	if migrationPerformed {
		ctxlog.Logger.Info("Legacy nod data migration completed",
			"database_version", core.CurrentDatabaseVersion,
		)
	}
	if len(removedTables) > 0 {
		ctxlog.Logger.Info("Legacy nod tables removed",
			"tables", strings.Join(removedTables, ","),
		)
	}

	return nil
}

func (m *NodDataMigrator) migrateLegacyData() error {
	if !m.db.Migrator().HasTable(legacyNodeCoreTable) {
		return nil
	}

	data, err := m.loadLegacyData()
	if err != nil {
		return err
	}

	if err := m.migrateWorkspaces(data); err != nil {
		return err
	}
	if err := m.migrateProjects(data); err != nil {
		return err
	}
	if err := m.migrateTags(data); err != nil {
		return err
	}
	if err := m.migrateContexts(data); err != nil {
		return err
	}
	if err := m.migrateContextTags(data); err != nil {
		return err
	}
	if err := m.migrateIntervals(data); err != nil {
		return err
	}
	if err := m.migrateClientProperties(data); err != nil {
		return err
	}

	return nil
}

func (m *NodDataMigrator) loadLegacyData() (*legacyData, error) {
	data := &legacyData{
		kv:      make(map[string]map[string]legacyNodeKV),
		content: make(map[string]map[string]legacyNodeContent),
		tags:    make(map[string]legacyTag),
	}

	if err := m.db.Table(legacyNodeCoreTable).
		Where("kind IN ?", []string{"workspace", "project", "context", "interval", "settings"}).
		Order("id").
		Find(&data.nodes).Error; err != nil {
		return nil, fmt.Errorf("load legacy nodes: %w", err)
	}

	if table := firstExistingTable(m.db, legacyKVTables); table != "" {
		var values []legacyNodeKV
		if err := m.db.Table(table).Find(&values).Error; err != nil {
			return nil, fmt.Errorf("load legacy key-values: %w", err)
		}
		for _, value := range values {
			if data.kv[value.NodeId] == nil {
				data.kv[value.NodeId] = make(map[string]legacyNodeKV)
			}
			data.kv[value.NodeId][value.Key] = value
		}
	}

	if table := firstExistingTable(m.db, legacyContentTables); table != "" {
		var values []legacyNodeContent
		if err := m.db.Table(table).Find(&values).Error; err != nil {
			return nil, fmt.Errorf("load legacy contents: %w", err)
		}
		for _, value := range values {
			if data.content[value.NodeId] == nil {
				data.content[value.NodeId] = make(map[string]legacyNodeContent)
			}
			data.content[value.NodeId][value.Key] = value
		}
	}

	if m.db.Migrator().HasTable(legacyTagTable) {
		var tags []legacyTag
		if err := m.db.Table(legacyTagTable).Find(&tags).Error; err != nil {
			return nil, fmt.Errorf("load legacy tags: %w", err)
		}
		for _, tag := range tags {
			data.tags[tag.Id] = tag
		}
	}

	if m.db.Migrator().HasTable(legacyNodeTagTable) {
		if err := m.db.Table(legacyNodeTagTable).Find(&data.contextTag).Error; err != nil {
			return nil, fmt.Errorf("load legacy node tags: %w", err)
		}
	}

	return data, nil
}

func (m *NodDataMigrator) migrateWorkspaces(data *legacyData) error {
	entities := make([]storage.WorkspaceEntity, 0)
	for _, node := range data.nodes {
		if node.Kind != "workspace" {
			continue
		}
		entities = append(entities, storage.WorkspaceEntity{
			Id:          node.Id,
			Name:        node.Name,
			Description: legacyContentValue(data.content, node.Id, "description"),
		})
	}

	if err := createByIdIfMissing(m.db, entities); err != nil {
		return fmt.Errorf("migrate workspaces: %w", err)
	}
	return nil
}

func (m *NodDataMigrator) migrateProjects(data *legacyData) error {
	entities := make([]storage.ProjectEntity, 0)
	for _, node := range data.nodes {
		if node.Kind != "project" {
			continue
		}
		entities = append(entities, storage.ProjectEntity{
			Id:          node.Id,
			Name:        node.Name,
			ParentId:    stringValue(node.ParentId),
			WorkspaceId: stringValue(node.NamespaceId),
		})
	}

	if err := createByIdIfMissing(m.db, entities); err != nil {
		return fmt.Errorf("migrate projects: %w", err)
	}
	return nil
}

func (m *NodDataMigrator) migrateTags(data *legacyData) error {
	entities := make([]storage.TagEntity, 0, len(data.tags))
	for _, tag := range data.tags {
		entities = append(entities, storage.TagEntity{Id: tag.Id, Name: tag.Name})
	}

	if err := createByIdIfMissing(m.db, entities); err != nil {
		return fmt.Errorf("migrate tags: %w", err)
	}
	return nil
}

func (m *NodDataMigrator) migrateContexts(data *legacyData) error {
	entities := make([]storage.ContextEntity, 0)
	for _, node := range data.nodes {
		if node.Kind != "context" {
			continue
		}

		var projectId *string
		if value := legacyTextValue(data.kv, node.Id, "projectId"); value != "" {
			projectId = &value
		}

		entities = append(entities, storage.ContextEntity{
			Id:          node.Id,
			Name:        node.Name,
			WorkspaceId: stringValue(node.NamespaceId),
			Status:      node.Status,
			Archived:    legacyBoolValue(data.kv, node.Id, "archived"),
			Description: legacyContentValue(data.content, node.Id, "description"),
			ProjectId:   projectId,
		})
	}

	if len(entities) == 0 {
		return nil
	}
	if err := m.db.Omit("Tags", "ProjectMetadata").
		Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "id"}}, DoNothing: true}).
		Create(&entities).Error; err != nil {
		return fmt.Errorf("migrate contexts: %w", err)
	}
	return nil
}

func (m *NodDataMigrator) migrateContextTags(data *legacyData) error {
	contextIds := make(map[string]struct{})
	for _, node := range data.nodes {
		if node.Kind == "context" {
			contextIds[node.Id] = struct{}{}
		}
	}

	entities := make([]storage.ContextTagEntity, 0)
	for _, relation := range data.contextTag {
		if _, isContext := contextIds[relation.NodeId]; !isContext {
			continue
		}
		if _, tagExists := data.tags[relation.TagId]; !tagExists {
			continue
		}
		entities = append(entities, storage.ContextTagEntity{
			ContextId: relation.NodeId,
			TagId:     relation.TagId,
		})
	}

	if len(entities) == 0 {
		return nil
	}
	if err := m.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&entities).Error; err != nil {
		return fmt.Errorf("migrate context tags: %w", err)
	}
	return nil
}

func (m *NodDataMigrator) migrateIntervals(data *legacyData) error {
	entities := make([]storage.IntervalEntity, 0)
	for _, node := range data.nodes {
		if node.Kind != "interval" {
			continue
		}
		entities = append(entities, storage.IntervalEntity{
			Id:          node.Id,
			ContextId:   stringValue(node.ParentId),
			Start:       legacyTimeValue(data.kv, node.Id, "start"),
			End:         legacyTimeValue(data.kv, node.Id, "end"),
			Duration:    time.Duration(legacyInt64Value(data.kv, node.Id, "duration")),
			Status:      node.Status,
			WorkspaceId: stringValue(node.NamespaceId),
		})
	}

	if err := createByIdIfMissing(m.db, entities); err != nil {
		return fmt.Errorf("migrate intervals: %w", err)
	}
	return nil
}

func (m *NodDataMigrator) migrateClientProperties(data *legacyData) error {
	for _, node := range data.nodes {
		if node.Kind != "settings" {
			continue
		}

		settings := map[string]string{
			"client.general.theme":    "light",
			"client.general.firstDay": "Monday",
			"client.general.timeZone": "browser",
		}
		for key := range settings {
			if value, ok := legacyTextValueIfPresent(data.kv, node.Id, key); ok {
				settings[key] = value
			}
		}

		if err := storage.NewClientPropertiesRepository(m.db).Save(core.NewSettings(settings)); err != nil {
			return fmt.Errorf("migrate client properties: %w", err)
		}
		return nil
	}

	return nil
}

func createByIdIfMissing[T any](db *gorm.DB, entities []T) error {
	if len(entities) == 0 {
		return nil
	}
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&entities).Error
}

func firstExistingTable(db *gorm.DB, tables []string) string {
	for _, table := range tables {
		if db.Migrator().HasTable(table) {
			return table
		}
	}
	return ""
}

func existingLegacyTables(db *gorm.DB) []string {
	tables := make([]string, 0, len(legacyNodTables))
	for _, table := range legacyNodTables {
		if db.Migrator().HasTable(table) {
			tables = append(tables, table)
		}
	}
	return tables
}

func dropLegacyTables(db *gorm.DB, tables []string) error {
	for _, table := range tables {
		if err := db.Migrator().DropTable(table); err != nil {
			return fmt.Errorf("drop legacy nod table %s: %w", table, err)
		}
	}
	return nil
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func legacyTextValue(values map[string]map[string]legacyNodeKV, nodeId, key string) string {
	value, _ := legacyTextValueIfPresent(values, nodeId, key)
	return value
}

func legacyTextValueIfPresent(values map[string]map[string]legacyNodeKV, nodeId, key string) (string, bool) {
	value, exists := values[nodeId][key]
	if !exists || value.ValueText == nil {
		return "", false
	}
	return *value.ValueText, true
}

func legacyBoolValue(values map[string]map[string]legacyNodeKV, nodeId, key string) bool {
	value, exists := values[nodeId][key]
	return exists && value.ValueBool != nil && *value.ValueBool
}

func legacyInt64Value(values map[string]map[string]legacyNodeKV, nodeId, key string) int64 {
	value, exists := values[nodeId][key]
	if !exists || value.ValueInt64 == nil {
		return 0
	}
	return *value.ValueInt64
}

func legacyTimeValue(values map[string]map[string]legacyNodeKV, nodeId, key string) *time.Time {
	value, exists := values[nodeId][key]
	if !exists || value.ValueTime == nil || value.ValueTime.IsZero() {
		return nil
	}
	instant := value.ValueTime.UTC()
	return &instant
}

func legacyContentValue(values map[string]map[string]legacyNodeContent, nodeId, key string) string {
	value, exists := values[nodeId][key]
	if !exists || value.Value == nil {
		return ""
	}
	return *value.Value
}

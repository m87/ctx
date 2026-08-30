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
		databaseVersion := strings.TrimSpace(systemInfo.DatabaseVersion)
		if databaseVersion != "" {
			if databaseVersion != core.CurrentDatabaseVersion {
				return fmt.Errorf("unsupported database version %q; expected %q", databaseVersion, core.CurrentDatabaseVersion)
			}
			if len(legacyTables) > 0 {
				return fmt.Errorf("legacy nod tables remain in an initialized %s database: %s", databaseVersion, strings.Join(legacyTables, ","))
			}
			return nil
		}
		if len(legacyTables) > 0 && !tx.Migrator().HasTable(legacyNodeCoreTable) {
			return fmt.Errorf("legacy nod schema is incomplete: missing %s", legacyNodeCoreTable)
		}

		if len(legacyTables) > 0 {
			migrationPerformed = true
			ctxlog.Logger.Info("Starting legacy nod data migration",
				"target_database_version", core.CurrentDatabaseVersion,
			)
		}

		legacyData, err := (&NodDataMigrator{db: tx}).migrateLegacyData()
		if err != nil {
			return err
		}
		if err := (&NodDataMigrator{db: tx}).verifyLegacyMigration(legacyData); err != nil {
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

func (m *NodDataMigrator) migrateLegacyData() (*legacyData, error) {
	if !m.db.Migrator().HasTable(legacyNodeCoreTable) {
		return nil, nil
	}

	data, err := m.loadLegacyData()
	if err != nil {
		return nil, err
	}

	if err := m.migrateWorkspaces(data); err != nil {
		return nil, err
	}
	if err := m.migrateProjects(data); err != nil {
		return nil, err
	}
	if err := m.migrateTags(data); err != nil {
		return nil, err
	}
	if err := m.migrateContexts(data); err != nil {
		return nil, err
	}
	if err := m.migrateContextTags(data); err != nil {
		return nil, err
	}
	if err := m.migrateIntervals(data); err != nil {
		return nil, err
	}
	if err := m.migrateClientProperties(data); err != nil {
		return nil, err
	}

	return data, nil
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

	if len(entities) > 0 {
		if err := m.db.Omit("LinkRules").Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(&entities).Error; err != nil {
			return fmt.Errorf("migrate workspaces: %w", err)
		}
	}

	for _, entity := range entities {
		if err := m.db.Delete(&storage.WorkspaceLinkRuleEntity{}, "workspace_id = ?", entity.Id).Error; err != nil {
			return fmt.Errorf("replace workspace link rules: %w", err)
		}
		rules := legacyWorkspaceLinkRules(data, entity.Id)
		if len(rules) > 0 {
			if err := m.db.Create(&rules).Error; err != nil {
				return fmt.Errorf("migrate workspace link rules: %w", err)
			}
		}
	}

	return nil
}

func legacyWorkspaceLinkRules(data *legacyData, workspaceId string) []storage.WorkspaceLinkRuleEntity {
	rules := make([]storage.WorkspaceLinkRuleEntity, 0)
	for position := 0; ; position++ {
		prefix := fmt.Sprintf("linkRule.%d", position)
		regexp, hasRegexp := legacyTextValueIfPresent(data.kv, workspaceId, prefix+".regexp")
		link, hasLink := legacyTextValueIfPresent(data.kv, workspaceId, prefix+".link")
		if !hasRegexp || !hasLink {
			break
		}
		rules = append(rules, storage.WorkspaceLinkRuleEntity{
			WorkspaceId: workspaceId,
			Position:    position,
			Regexp:      regexp,
			Link:        link,
		})
	}
	return rules
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

	if err := upsertById(m.db, entities); err != nil {
		return fmt.Errorf("migrate projects: %w", err)
	}
	return nil
}

func (m *NodDataMigrator) migrateTags(data *legacyData) error {
	entities := make([]storage.TagEntity, 0, len(data.tags))
	for _, tag := range data.tags {
		entities = append(entities, storage.TagEntity{Id: tag.Id, Name: tag.Name})
	}

	if err := upsertById(m.db, entities); err != nil {
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
		Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "id"}}, UpdateAll: true}).
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
	for contextId := range contextIds {
		if err := m.db.Delete(&storage.ContextTagEntity{}, "context_id = ?", contextId).Error; err != nil {
			return fmt.Errorf("replace context tags: %w", err)
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

	if err := upsertById(m.db, entities); err != nil {
		return fmt.Errorf("migrate intervals: %w", err)
	}
	return nil
}

func (m *NodDataMigrator) migrateClientProperties(data *legacyData) error {
	settings, exists := legacyClientSettings(data)
	if exists {
		if err := storage.NewClientPropertiesRepository(m.db).Save(core.NewSettings(settings)); err != nil {
			return fmt.Errorf("migrate client properties: %w", err)
		}
	}
	return nil
}

func legacyClientSettings(data *legacyData) (map[string]string, bool) {
	for _, node := range data.nodes {
		if node.Kind != "settings" {
			continue
		}

		settings := map[string]string{
			"client.general.theme":    "light",
			"client.general.firstDay": "Monday",
			"client.general.timeZone": "browser",
		}
		for key, value := range data.kv[node.Id] {
			if strings.HasPrefix(key, "client.") && value.ValueText != nil {
				settings[key] = *value.ValueText
			}
		}
		return settings, true
	}
	return nil, false
}

func (m *NodDataMigrator) verifyLegacyMigration(data *legacyData) error {
	if data == nil {
		return nil
	}

	idsByKind := make(map[string][]string)
	workspaceIds := make(map[string]struct{})
	contextIds := make(map[string]struct{})
	for _, node := range data.nodes {
		idsByKind[node.Kind] = append(idsByKind[node.Kind], node.Id)
		if node.Kind == "workspace" {
			workspaceIds[node.Id] = struct{}{}
		}
		if node.Kind == "context" {
			contextIds[node.Id] = struct{}{}
		}
	}

	for _, check := range []struct {
		kind  string
		table string
	}{
		{kind: "workspace", table: "workspaces"},
		{kind: "project", table: "projects"},
		{kind: "context", table: "contexts"},
		{kind: "interval", table: "intervals"},
	} {
		if err := verifyMigratedIds(m.db, check.table, idsByKind[check.kind]); err != nil {
			return fmt.Errorf("verify migrated %s records: %w", check.kind, err)
		}
	}

	tagIds := make([]string, 0, len(data.tags))
	for id := range data.tags {
		tagIds = append(tagIds, id)
	}
	if err := verifyMigratedIds(m.db, "tag", tagIds); err != nil {
		return fmt.Errorf("verify migrated tags: %w", err)
	}

	expectedRules := make(map[string]storage.WorkspaceLinkRuleEntity)
	for workspaceId := range workspaceIds {
		for _, rule := range legacyWorkspaceLinkRules(data, workspaceId) {
			expectedRules[workspaceRuleKey(rule.WorkspaceId, rule.Position)] = rule
		}
	}
	var storedRules []storage.WorkspaceLinkRuleEntity
	if err := m.db.Find(&storedRules).Error; err != nil {
		return fmt.Errorf("verify workspace link rules: %w", err)
	}
	actualRuleCount := 0
	for _, rule := range storedRules {
		if _, migratedWorkspace := workspaceIds[rule.WorkspaceId]; !migratedWorkspace {
			continue
		}
		actualRuleCount++
		expected, exists := expectedRules[workspaceRuleKey(rule.WorkspaceId, rule.Position)]
		if !exists || expected.Regexp != rule.Regexp || expected.Link != rule.Link {
			return fmt.Errorf("verify workspace link rules: unexpected rule for workspace %s at position %d", rule.WorkspaceId, rule.Position)
		}
	}
	if actualRuleCount != len(expectedRules) {
		return fmt.Errorf("verify workspace link rules: expected %d, found %d", len(expectedRules), actualRuleCount)
	}

	expectedRelations := make(map[string]struct{})
	for _, relation := range data.contextTag {
		if _, isContext := contextIds[relation.NodeId]; !isContext {
			continue
		}
		if _, tagExists := data.tags[relation.TagId]; !tagExists {
			continue
		}
		expectedRelations[contextTagKey(relation.NodeId, relation.TagId)] = struct{}{}
	}
	var storedRelations []storage.ContextTagEntity
	if err := m.db.Find(&storedRelations).Error; err != nil {
		return fmt.Errorf("verify context tags: %w", err)
	}
	actualRelationCount := 0
	for _, relation := range storedRelations {
		if _, migratedContext := contextIds[relation.ContextId]; !migratedContext {
			continue
		}
		actualRelationCount++
		if _, exists := expectedRelations[contextTagKey(relation.ContextId, relation.TagId)]; !exists {
			return fmt.Errorf("verify context tags: unexpected relation %s/%s", relation.ContextId, relation.TagId)
		}
	}
	if actualRelationCount != len(expectedRelations) {
		return fmt.Errorf("verify context tags: expected %d, found %d", len(expectedRelations), actualRelationCount)
	}

	if expectedSettings, exists := legacyClientSettings(data); exists {
		storedSettings, err := storage.NewClientPropertiesRepository(m.db).Load()
		if err != nil {
			return fmt.Errorf("verify client properties: %w", err)
		}
		actualSettings := storedSettings.Values()
		if len(actualSettings) != len(expectedSettings) {
			return fmt.Errorf("verify client properties: expected %d, found %d", len(expectedSettings), len(actualSettings))
		}
		for key, expected := range expectedSettings {
			if actualSettings[key] != expected {
				return fmt.Errorf("verify client properties: unexpected value for %s", key)
			}
		}
	}

	return nil
}

func verifyMigratedIds(db *gorm.DB, table string, ids []string) error {
	const batchSize = 500
	var found int64
	for start := 0; start < len(ids); start += batchSize {
		end := min(start+batchSize, len(ids))
		var count int64
		if err := db.Table(table).Where("id IN ?", ids[start:end]).Count(&count).Error; err != nil {
			return err
		}
		found += count
	}
	if found != int64(len(ids)) {
		return fmt.Errorf("expected %d source IDs, found %d", len(ids), found)
	}
	return nil
}

func workspaceRuleKey(workspaceId string, position int) string {
	return fmt.Sprintf("%s\x00%d", workspaceId, position)
}

func contextTagKey(contextId, tagId string) string {
	return contextId + "\x00" + tagId
}

func upsertById[T any](db *gorm.DB, entities []T) error {
	if len(entities) == 0 {
		return nil
	}
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
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

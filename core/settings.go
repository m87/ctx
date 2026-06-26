package core

import (
	"errors"
	"strings"
	"time"

	"github.com/m87/nod"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type GeneralSettings struct {
	theme    string
	firstDay string
}

type Settings struct {
	raw     map[string]string
	general GeneralSettings
}

type SettingsMapper struct{}

const SettingsType = "settings"

type SettingsManager struct {
	SettingsRepository   SettingsRepository
	SystemInfoRepository SystemInfoRepository
	cache                *Settings
}

var defaultClientSettings = map[string]string{
	"client.general.theme":    "light",
	"client.general.firstDay": "Monday",
}

func NewSettingsManager(settingsRepo SettingsRepository, systemInfoRepo SystemInfoRepository) *SettingsManager {
	return &SettingsManager{
		SettingsRepository:   settingsRepo,
		SystemInfoRepository: systemInfoRepo,
	}
}

func (m *SettingsManager) InitSettingsIfNotExists() error {
	settings, err := m.SettingsRepository.Load()
	if err == nil {
		m.cache = settings
		return nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	defaultSettings := NewSettings(defaultClientSettings)
	if err := m.SettingsRepository.Save(defaultSettings); err != nil {
		return err
	}
	m.cache = defaultSettings
	return nil
}

func (m *SettingsManager) GetClientKey(key string) (string, error) {
	if strings.HasPrefix(key, "client.") {
		return m.GetKey(key)
	}
	return "", nil
}

func (m *SettingsManager) GetClient() (map[string]string, error) {
	if m.cache == nil {
		settings, err := m.SettingsRepository.Load()
		if err != nil {
			return nil, err
		}
		m.cache = settings
	}

	clientSettings := m.filterClientSettings(m.cache.raw)
	return clientSettings, nil
}

func (m *SettingsManager) sanitizeSettings(settings map[string]string) map[string]string {
	sanitized := make(map[string]string)
	for key, value := range settings {
		if strings.HasPrefix(key, "client.") {
			sanitized[key] = value
		}
	}
	return sanitized
}

func (m *SettingsManager) SaveClient(settings map[string]string) error {
	clientSettings := m.sanitizeSettings(settings)
	if m.cache == nil {
		current, err := m.SettingsRepository.Load()
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		m.cache = current
	}

	mergedSettings := make(map[string]string)
	if m.cache != nil {
		for key, value := range m.cache.raw {
			mergedSettings[key] = value
		}
	}
	for key, value := range clientSettings {
		mergedSettings[key] = value
	}

	s := NewSettings(mergedSettings)
	err := m.SettingsRepository.Save(s)
	if err != nil {
		return err
	}
	m.cache = s
	return nil
}

func (m *SettingsManager) filterClientSettings(settings map[string]string) map[string]string {
	clientSettings := make(map[string]string)
	for key, value := range settings {
		if strings.HasPrefix(key, "client.") {
			clientSettings[key] = value
		}
	}
	return clientSettings
}

func (m *SettingsManager) GetKey(key string) (string, error) {
	if m.cache == nil {
		settings, err := m.SettingsRepository.Load()
		if err != nil {
			return "", err
		}
		m.cache = settings
	}

	if value, ok := m.cache.raw[key]; ok {
		return value, nil
	}

	if viper.InConfig(key) {
		return viper.GetString(key), nil
	}
	return "", nil
}

func (m *SettingsManager) Save(settings map[string]string) error {
	s := NewSettings(settings)
	err := m.SettingsRepository.Save(s)
	if err != nil {
		return err
	}
	m.cache = s
	return nil
}

func NewSettingsMapper() *SettingsMapper {
	return &SettingsMapper{}
}

func NewSettings(raw map[string]string) *Settings {
	settings := copySettings(raw)
	return &Settings{
		raw: settings,
		general: GeneralSettings{
			theme:    settings["client.general.theme"],
			firstDay: settings["client.general.firstDay"],
		},
	}
}

func copySettings(settings map[string]string) map[string]string {
	copied := make(map[string]string, len(settings))
	for key, value := range settings {
		copied[key] = value
	}
	return copied
}

func (m *SettingsMapper) ToNode(settings *Settings) (*nod.Node, error) {
	raw := copySettings(settings.raw)
	if raw["client.general.theme"] == "" {
		raw["client.general.theme"] = settings.general.theme
	}
	if raw["client.general.firstDay"] == "" {
		raw["client.general.firstDay"] = settings.general.firstDay
	}

	kv := make(map[string]*nod.KV, len(raw))
	for key, value := range raw {
		valueText := value
		kv[key] = &nod.KV{Key: key, ValueText: &valueText}
	}

	node := &nod.Node{
		Core: nod.NodeCore{
			Id:        "settingsV1",
			Name:      "settingsV1",
			Kind:      SettingsType,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		KV: kv,
	}
	return node, nil
}

func (m *SettingsMapper) FromNode(node *nod.Node) (*Settings, error) {
	raw := make(map[string]string, len(node.KV))
	for key, value := range node.KV {
		if value.ValueText != nil {
			raw[key] = nod.SafeString(node.KV, key)
		}
	}
	return &Settings{
		raw: raw,
		general: GeneralSettings{
			theme:    nod.SafeString(node.KV, "client.general.theme"),
			firstDay: nod.SafeString(node.KV, "client.general.firstDay"),
		},
	}, nil
}

func (m *SettingsMapper) IsApplicable(node *nod.Node) bool {
	return node.Core.Kind == SettingsType
}

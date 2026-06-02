package core

import (
	"strings"
	"time"

	"github.com/m87/nod"
	"github.com/spf13/viper"
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
	SettingsRepository SettingsRepository
	cache              *Settings
}

func NewSettingsManager(settingsRepo SettingsRepository) *SettingsManager {
	return &SettingsManager{
		SettingsRepository: settingsRepo,
	}
}

func (m *SettingsManager) InitSettingsIfNotExists() error {
	_, err := m.SettingsRepository.Load()
	if err != nil {
		defaultSettings := &Settings{
			raw: map[string]string{
				"client.general.theme":    "light",
				"client.general.firstDay": "Monday",
			},
			general: GeneralSettings{
				theme:    "light",
				firstDay: "Monday",
			},
		}
		m.cache = defaultSettings
		return m.SettingsRepository.Save(defaultSettings)
	}
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

func (m *SettingsManager) SaveClient(settings map[string]string) error {
	s := &Settings{
		raw: settings,
		general: GeneralSettings{
			theme:    settings["client.general.theme"],
			firstDay: settings["client.general.firstDay"],
		},
	}
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
	s := &Settings{
		raw: settings,
		general: GeneralSettings{
			theme:    settings["client.general.theme"],
			firstDay: settings["client.general.firstDay"],
		},
	}
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

func (m *SettingsMapper) ToNode(settings *Settings) (*nod.Node, error) {
	node := &nod.Node{
		Core: nod.NodeCore{
			Id:        "settingsV1",
			Name:      "settingsV1",
			Kind:      SettingsType,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		KV: map[string]*nod.KV{
			"client.general.theme":    &nod.KV{Key: "client.general.theme", ValueText: &settings.general.theme},
			"client.general.firstDay": &nod.KV{Key: "client.general.firstDay", ValueText: &settings.general.firstDay},
		},
	}
	return node, nil
}

func (m *SettingsMapper) FromNode(node *nod.Node) (*Settings, error) {
	return &Settings{
		general: GeneralSettings{
			theme:    *node.KV["client.general.theme"].ValueText,
			firstDay: *node.KV["client.general.firstDay"].ValueText,
		},
	}, nil
}

func (m *SettingsMapper) IsApplicable(node *nod.Node) bool {
	return node.Core.Kind == SettingsType
}

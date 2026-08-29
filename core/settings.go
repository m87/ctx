package core

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Settings struct {
	raw map[string]string
}

type SettingsManager struct {
	ClientPropertiesRepository SettingsRepository
	cache                      *Settings
}

var defaultClientSettings = map[string]string{
	"client.general.theme":    "light",
	"client.general.firstDay": "Monday",
	"client.general.timeZone": "browser",
}

func NewSettingsManager(clientPropertiesRepo SettingsRepository) *SettingsManager {
	return &SettingsManager{
		ClientPropertiesRepository: clientPropertiesRepo,
	}
}

func (m *SettingsManager) InitSettingsIfNotExists() error {
	settings, err := m.ClientPropertiesRepository.Load()
	if err == nil {
		m.cache = settings
		return nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	defaultSettings := NewSettings(defaultClientSettings)
	if err := m.ClientPropertiesRepository.Save(defaultSettings); err != nil {
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
		settings, err := m.ClientPropertiesRepository.Load()
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
	if err := validateClientTimeZone(clientSettings); err != nil {
		return err
	}
	if m.cache == nil {
		current, err := m.ClientPropertiesRepository.Load()
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
	err := m.ClientPropertiesRepository.Save(s)
	if err != nil {
		return err
	}
	m.cache = s
	return nil
}

type InvalidTimeZoneError struct {
	TimeZone string
}

func (e *InvalidTimeZoneError) Error() string {
	return fmt.Sprintf("invalid time zone %q", e.TimeZone)
}

func validateClientTimeZone(settings map[string]string) error {
	zone, exists := settings["client.general.timeZone"]
	if !exists {
		return nil
	}
	zone = strings.TrimSpace(zone)
	if zone == "browser" {
		return nil
	}
	if zone == "" {
		return &InvalidTimeZoneError{TimeZone: zone}
	}
	if _, err := time.LoadLocation(zone); err != nil {
		return &InvalidTimeZoneError{TimeZone: zone}
	}
	settings["client.general.timeZone"] = zone
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
		settings, err := m.ClientPropertiesRepository.Load()
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
	return m.SaveClient(settings)
}

func NewSettings(raw map[string]string) *Settings {
	return &Settings{raw: copySettings(raw)}
}

func (s *Settings) Values() map[string]string {
	if s == nil {
		return nil
	}
	return copySettings(s.raw)
}

func copySettings(settings map[string]string) map[string]string {
	copied := make(map[string]string, len(settings))
	for key, value := range settings {
		copied[key] = value
	}
	return copied
}

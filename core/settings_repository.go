package core

type SettingsRepository interface {
	Save(settings *Settings) error
	Load() (*Settings, error)
}

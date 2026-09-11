package core

var _ SettingsRepository = (*SettingsRepositoryMock)(nil)

type SettingsRepositoryMock struct {
	settings  *Settings
	loadError error
	saveError error
	loadCalls int
	saveCalls int
	saved     *Settings
}

func (m *SettingsRepositoryMock) Save(settings *Settings) error {
	m.saveCalls++
	m.saved = settings
	if m.saveError != nil {
		return m.saveError
	}
	m.settings = settings
	return nil
}

func (m *SettingsRepositoryMock) Load() (*Settings, error) {
	m.loadCalls++
	if m.loadError != nil {
		return nil, m.loadError
	}
	return m.settings, nil
}

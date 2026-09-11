package core

var _ SystemInfoRepository = (*SystemInfoRepositoryMock)(nil)

type SystemInfoRepositoryMock struct {
	info      *SystemInfo
	loadError error
	saveError error
}

func (m *SystemInfoRepositoryMock) Load() (*SystemInfo, error) {
	return m.info, m.loadError
}

func (m *SystemInfoRepositoryMock) Save(info *SystemInfo) error {
	if m.saveError != nil {
		return m.saveError
	}
	m.info = info
	return nil
}

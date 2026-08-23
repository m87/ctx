package core

type SystemInfoRepository interface {
	Load() (*SystemInfo, error)
	Save(info *SystemInfo) error
}

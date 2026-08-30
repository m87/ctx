package storage

const (
	systemPropertiesId = "system"
	clientPropertiesId = "client"
)

type Properties struct {
	Id              string `gorm:"primaryKey"`
	DatabaseVersion string
	ClientId        string
}

func (Properties) TableName() string {
	return "properties"
}

type ClientProperties struct {
	Id       string `gorm:"primaryKey"`
	Theme    string
	FirstDay string
	Timezone string
	Values   string
}

func (ClientProperties) TableName() string {
	return "client_properties"
}

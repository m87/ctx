package core

const CurrentDatabaseVersion = "0.7.1"

type SystemInfo struct {
	DatabaseVersion string `json:"dbVersion,omitempty"`
	ClientId        string `json:"clientId,omitempty"`
}

package core

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/m87/nod"
)

const (
	SystemInfoType         = "system_info"
	SystemInfoName         = "systemInfoV1"
	CurrentDatabaseVersion = "0.6.0"
)

type SystemInfo struct {
	DatabaseVersion string `json:"dbVersion,omitempty"`
	ClientId        string
}

type SystemInfoMapper struct{}

func NewSystemInfoMapper() *SystemInfoMapper {
	return &SystemInfoMapper{}
}

func (m *SystemInfoMapper) ToNode(info *SystemInfo) (*nod.Node, error) {
	clientId := strings.TrimSpace(info.ClientId)
	if clientId == "" {
		return nil, fmt.Errorf("system info client ID is required")
	}
	databaseVersion := info.DatabaseVersion
	return &nod.Node{
		Core: nod.NodeCore{
			Id:   clientId,
			Name: SystemInfoName,
			Kind: SystemInfoType,
		},
		KV: map[string]*nod.NodeKV{
			"database_version": {Key: "database_version", ValueText: &databaseVersion},
		},
	}, nil
}

func (m *SystemInfoMapper) FromNode(node *nod.Node) (*SystemInfo, error) {
	return &SystemInfo{
		DatabaseVersion: nodString(node.KV, "database_version"),
		ClientId:        node.Core.Id,
	}, nil
}

func (m *SystemInfoMapper) IsApplicable(node *nod.Node) bool {
	return node.Core.Kind == SystemInfoType
}

func DatabaseVersionNeedsMigration(current, target string) (bool, error) {
	if strings.TrimSpace(current) == "" {
		return true, nil
	}
	currentParts, err := parseDatabaseVersion(current)
	if err != nil {
		return false, err
	}
	targetParts, err := parseDatabaseVersion(target)
	if err != nil {
		return false, err
	}
	for i := range currentParts {
		if currentParts[i] < targetParts[i] {
			return true, nil
		}
		if currentParts[i] > targetParts[i] {
			return false, nil
		}
	}
	return false, nil
}

func parseDatabaseVersion(version string) ([3]int, error) {
	var parsed [3]int
	parts := strings.Split(strings.TrimSpace(version), ".")
	if len(parts) != len(parsed) {
		return parsed, fmt.Errorf("invalid database version %q", version)
	}
	for i, part := range parts {
		value, err := strconv.Atoi(part)
		if err != nil || value < 0 {
			return parsed, fmt.Errorf("invalid database version %q", version)
		}
		parsed[i] = value
	}
	return parsed, nil
}

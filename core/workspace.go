package core

import (
	"time"
)

type Workspace struct {
	Id          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description,omitempty"`
	Properties  *WorkspaceSettings `json:"properties,omitempty"`
}

type WorkspaceContextStats struct {
	ContextId     string        `json:"contextId"`
	Duration      time.Duration `json:"duration"`
	Percentage    float64       `json:"percentage"`
	IntervalCount int           `json:"intervalCount"`
}

type WorkspaceStats struct {
	WorkspaceId   string                   `json:"workspaceId"`
	Contexts      []*Context               `json:"contexts"`
	ContextStats  []*WorkspaceContextStats `json:"contextStats"`
	TotalDuration time.Duration            `json:"totalDuration"`
	TotalSessions int                      `json:"totalSessions"`
}

package core

import (
	"fmt"
	"strings"
)

type SyncValidationError struct {
	Resource string
	Index    int
	Reason   string
}

func (e *SyncValidationError) Error() string {
	return fmt.Sprintf("invalid %s at index %d: %s", e.Resource, e.Index, e.Reason)
}

func invalidSyncItem(resource string, index int, reason string) error {
	return &SyncValidationError{Resource: resource, Index: index, Reason: reason}
}

func (m *ContextManager) ImportSyncWorkspaces(workspaces []*Workspace) error {
	return m.RunInTransaction(func(tx *ContextManager) error {
		seen := make(map[string]struct{}, len(workspaces))
		for i, workspace := range workspaces {
			if workspace == nil {
				return invalidSyncItem("workspace", i, "item is null")
			}
			if strings.TrimSpace(workspace.Id) == "" || strings.TrimSpace(workspace.Name) == "" {
				return invalidSyncItem("workspace", i, "id and name are required")
			}
			if _, exists := seen[workspace.Id]; exists {
				return invalidSyncItem("workspace", i, "duplicate id")
			}
			seen[workspace.Id] = struct{}{}
		}
		_, err := tx.WorkspaceRepository.SaveAll(workspaces)
		return err
	})
}

func (m *ContextManager) ImportSyncContexts(contexts []*Context) error {
	return m.RunInTransaction(func(tx *ContextManager) error {
		existing, err := tx.ContextRepository.List()
		if err != nil {
			return err
		}
		existingByID := make(map[string]*Context, len(existing))
		for _, context := range existing {
			if context != nil {
				existingByID[context.Id] = context
			}
		}
		seen := make(map[string]struct{}, len(contexts))
		for i, context := range contexts {
			if context == nil {
				return invalidSyncItem("context", i, "item is null")
			}
			if strings.TrimSpace(context.Id) == "" || strings.TrimSpace(context.Name) == "" || strings.TrimSpace(context.WorkspaceId) == "" {
				return invalidSyncItem("context", i, "id, name and workspace id are required")
			}
			if context.Status != "" && context.Status != "active" && context.Status != "inactive" {
				return invalidSyncItem("context", i, "invalid status")
			}
			if context.Archived && context.Status == "active" {
				return invalidSyncItem("context", i, "archived context cannot be active")
			}
			if _, exists := seen[context.Id]; exists {
				return invalidSyncItem("context", i, "duplicate id")
			}
			seen[context.Id] = struct{}{}
			if previous := existingByID[context.Id]; previous != nil && previous.WorkspaceId != context.WorkspaceId {
				return invalidSyncItem("context", i, "workspace cannot change")
			}
			workspace, err := tx.WorkspaceRepository.GetById(context.WorkspaceId)
			if err != nil {
				return err
			}
			if workspace == nil {
				return invalidSyncItem("context", i, "workspace does not exist")
			}
			projectID := ""
			if context.ProjectId != nil {
				projectID = *context.ProjectId
			}
			if context.Project != nil && context.Project.Id != "" {
				if projectID != "" && projectID != context.Project.Id {
					return invalidSyncItem("context", i, "project ids disagree")
				}
				projectID = context.Project.Id
			}
			if projectID != "" {
				projects, err := tx.ProjectRepository.ListIncludingArchived(context.WorkspaceId)
				if err != nil {
					return err
				}
				found := false
				for _, project := range projects {
					if project != nil && project.Id == projectID {
						found = true
						context.Project = &ProjectMetadata{Id: project.Id, Name: project.Name}
						break
					}
				}
				if !found {
					return invalidSyncItem("context", i, "project does not exist in workspace")
				}
			}
			for _, tag := range context.Tags {
				if tag == nil || strings.TrimSpace(tag.Id) == "" || strings.TrimSpace(tag.Name) == "" {
					return invalidSyncItem("context", i, "tag id and name are required")
				}
			}
		}
		_, err = tx.ContextRepository.SaveAll(contexts)
		return err
	})
}

func (m *ContextManager) ImportSyncIntervals(intervals []*Interval) error {
	return m.RunInTransaction(func(tx *ContextManager) error {
		existing, err := tx.IntervalRepository.List()
		if err != nil {
			return err
		}
		existingByID := make(map[string]*Interval, len(existing))
		for _, interval := range existing {
			if interval != nil {
				existingByID[interval.Id] = interval
			}
		}
		contexts, err := tx.ContextRepository.List()
		if err != nil {
			return err
		}
		contextsByID := make(map[string]*Context, len(contexts))
		for _, context := range contexts {
			if context != nil {
				contextsByID[context.Id] = context
			}
		}
		seen := make(map[string]struct{}, len(intervals))
		for i, interval := range intervals {
			if interval == nil {
				return invalidSyncItem("interval", i, "item is null")
			}
			if strings.TrimSpace(interval.Id) == "" || strings.TrimSpace(interval.ContextId) == "" || strings.TrimSpace(interval.WorkspaceId) == "" {
				return invalidSyncItem("interval", i, "id, context id and workspace id are required")
			}
			if _, exists := seen[interval.Id]; exists {
				return invalidSyncItem("interval", i, "duplicate id")
			}
			seen[interval.Id] = struct{}{}
			if previous := existingByID[interval.Id]; previous != nil && (previous.ContextId != interval.ContextId || previous.WorkspaceId != interval.WorkspaceId) {
				return invalidSyncItem("interval", i, "context and workspace cannot change")
			}
			context := contextsByID[interval.ContextId]
			if context == nil || context.WorkspaceId != interval.WorkspaceId {
				return invalidSyncItem("interval", i, "context does not exist in workspace")
			}
			if !timeIsSet(interval.Start) {
				return invalidSyncItem("interval", i, "start time is required")
			}
			if interval.Status != "active" && interval.Status != "completed" && interval.Status != "inactive" && interval.Status != "" {
				return invalidSyncItem("interval", i, "invalid status")
			}
			if interval.Status == "active" {
				if timeIsSet(interval.End) {
					return invalidSyncItem("interval", i, "active interval cannot have an end time")
				}
			} else if !timeIsSet(interval.End) || !interval.End.After(*interval.Start) {
				return invalidSyncItem("interval", i, "end time must be after start time")
			}
			interval.Start = utcTimePointer(interval.Start)
			interval.End = utcTimePointer(interval.End)
			interval.Duration = durationBetween(interval.Start, interval.End)
		}
		_, err = tx.IntervalRepository.SaveAll(intervals)
		return err
	})
}

package core

import (
	"fmt"
	"sort"
	"time"
)

const MaxArchiveThresholdDays = 365000

type ArchiveCandidate struct {
	Id             string           `json:"id"`
	Name           string           `json:"name"`
	LastIntervalAt time.Time        `json:"lastIntervalAt"`
	Project        *ProjectMetadata `json:"project,omitempty"`
}

type ArchiveCandidatesPreview struct {
	Cutoff   time.Time           `json:"cutoff"`
	Contexts []*ArchiveCandidate `json:"contexts"`
}

type BulkArchiveResult struct {
	Cutoff        time.Time           `json:"cutoff"`
	ArchivedCount int                 `json:"archivedCount"`
	Contexts      []*ArchiveCandidate `json:"contexts"`
}

type InvalidArchiveThresholdError struct {
	Days int
}

func (e *InvalidArchiveThresholdError) Error() string {
	return fmt.Sprintf("archive threshold must be between 1 and %d days", MaxArchiveThresholdDays)
}

type ArchiveContextActiveError struct {
	ContextId string
}

func (e *ArchiveContextActiveError) Error() string {
	return "cannot archive active context: " + e.ContextId
}

func (m *ContextManager) ArchiveContext(contextId string) error {
	context, err := m.ContextRepository.GetById(contextId)
	if err != nil {
		return err
	}

	if context.Status == "active" {
		return &ArchiveContextActiveError{ContextId: contextId}
	}

	context.Archived = true
	context.Status = "archived"

	_, err = m.ContextRepository.Save(context)
	if err != nil {
		return err
	}

	return nil
}

func (m *ContextManager) RestoreContext(contextId string) error {
	context, err := m.ContextRepository.GetById(contextId)
	if err != nil {
		return err
	}

	context.Archived = false
	context.Status = "inactive"

	_, err = m.ContextRepository.Save(context)
	if err != nil {
		return err
	}

	return nil
}

func (m *ContextManager) ListArchiveCandidates(workspaceId string, olderThanDays int, location *time.Location) (*ArchiveCandidatesPreview, error) {
	cutoff, err := m.archiveCutoff(olderThanDays, location)
	if err != nil {
		return nil, err
	}

	preview, _, err := m.listArchiveCandidatesBefore(workspaceId, cutoff)
	return preview, err
}

func (m *ContextManager) ArchiveStaleContexts(workspaceId string, olderThanDays int, location *time.Location) (*BulkArchiveResult, error) {
	cutoff, err := m.archiveCutoff(olderThanDays, location)
	if err != nil {
		return nil, err
	}

	result := &BulkArchiveResult{Cutoff: cutoff, Contexts: []*ArchiveCandidate{}}
	err = m.RunInTransaction(func(txManager *ContextManager) error {
		preview, contextsById, err := txManager.listArchiveCandidatesBefore(workspaceId, cutoff)
		if err != nil {
			return err
		}

		contexts := make([]*Context, 0, len(preview.Contexts))
		for _, candidate := range preview.Contexts {
			context := contextsById[candidate.Id]
			if context == nil {
				continue
			}
			context.Archived = true
			context.Status = "archived"
			contexts = append(contexts, context)
		}

		if len(contexts) > 0 {
			if _, err := txManager.ContextRepository.SaveAll(contexts); err != nil {
				return err
			}
		}

		result.Contexts = preview.Contexts
		result.ArchivedCount = len(contexts)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *ContextManager) archiveCutoff(olderThanDays int, location *time.Location) (time.Time, error) {
	if olderThanDays < 1 || olderThanDays > MaxArchiveThresholdDays {
		return time.Time{}, &InvalidArchiveThresholdError{Days: olderThanDays}
	}
	if m.TimeProvider == nil {
		return time.Time{}, fmt.Errorf("time provider is required")
	}

	if location == nil {
		location = time.UTC
	}
	now := m.TimeProvider.Now().In(location)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
	return today.AddDate(0, 0, -olderThanDays).UTC(), nil
}

func (m *ContextManager) listArchiveCandidatesBefore(workspaceId string, cutoff time.Time) (*ArchiveCandidatesPreview, map[string]*Context, error) {
	contexts, err := m.ContextRepository.ListByWorkspace(workspaceId)
	if err != nil {
		return nil, nil, err
	}

	contextsById := make(map[string]*Context, len(contexts))
	for _, context := range contexts {
		if context == nil || context.Archived || context.Status == "active" {
			continue
		}
		contextsById[context.Id] = context
	}

	intervals, err := m.IntervalRepository.List()
	if err != nil {
		return nil, nil, err
	}

	latestByContextId := make(map[string]time.Time, len(contextsById))
	ineligibleContextIds := make(map[string]struct{})
	for _, interval := range intervals {
		if interval == nil {
			continue
		}
		if _, exists := contextsById[interval.ContextId]; !exists {
			continue
		}
		if interval.Status == "active" || !timeIsSet(interval.End) {
			ineligibleContextIds[interval.ContextId] = struct{}{}
			continue
		}

		latest, ok := latestIntervalTime(interval)
		if !ok {
			continue
		}
		if previous, exists := latestByContextId[interval.ContextId]; !exists || latest.After(previous) {
			latestByContextId[interval.ContextId] = latest
		}
	}

	candidates := make([]*ArchiveCandidate, 0)
	for contextId, context := range contextsById {
		if _, ineligible := ineligibleContextIds[contextId]; ineligible {
			continue
		}
		latest, exists := latestByContextId[contextId]
		if !exists || !latest.Before(cutoff) {
			continue
		}
		candidates = append(candidates, &ArchiveCandidate{
			Id:             context.Id,
			Name:           context.Name,
			LastIntervalAt: latest,
			Project:        context.Project,
		})
	}

	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].LastIntervalAt.Equal(candidates[j].LastIntervalAt) {
			return candidates[i].Name < candidates[j].Name
		}
		return candidates[i].LastIntervalAt.Before(candidates[j].LastIntervalAt)
	})

	return &ArchiveCandidatesPreview{
		Cutoff:   cutoff,
		Contexts: candidates,
	}, contextsById, nil
}

func latestIntervalTime(interval *Interval) (time.Time, bool) {
	var latest time.Time
	if timeIsSet(interval.Start) {
		latest = interval.Start.UTC()
	}
	if timeIsSet(interval.End) && interval.End.UTC().After(latest) {
		latest = interval.End.UTC()
	}
	return latest, !latest.IsZero()
}

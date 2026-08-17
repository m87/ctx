import { inject, Injectable } from '@angular/core';
import { lastValueFrom } from 'rxjs';
import { ContextService } from './context.service';

export const contextQueryKeys = {
  all: ['contexts'] as const,
  intervals: () => [...contextQueryKeys.all, 'intervals'] as const,
  intervalsFor: (contextId: string) => [...contextQueryKeys.intervals(), contextId] as const,
  lists: () => [...contextQueryKeys.all, 'list'] as const,
  list: (workspaceId: string | null, includeArchived: boolean) =>
    [...contextQueryKeys.lists(), workspaceId, includeArchived] as const,
  details: () => [...contextQueryKeys.all, 'get'] as const,
  detail: (contextId: string) => [...contextQueryKeys.details(), contextId] as const,
  active: () => [...contextQueryKeys.all, 'active'] as const,
  stats: () => [...contextQueryKeys.all, 'stats'] as const,
  statsFor: (contextId: string) => [...contextQueryKeys.stats(), contextId] as const,
  stat: (contextId: string, date: string, timeZone: string) =>
    [...contextQueryKeys.statsFor(contextId), date, timeZone] as const,
  archivization: () => [...contextQueryKeys.all, 'archivization'] as const,
  archiveCandidates: (workspaceId: string | null, olderThanDays: number, timeZone: string) =>
    [...contextQueryKeys.archivization(), workspaceId, olderThanDays, timeZone] as const,
};

@Injectable({
  providedIn: 'root',
})
export class ContextQueries {
  private readonly contextService = inject(ContextService);

  intervals(contextId: string) {
    return {
      queryKey: contextQueryKeys.intervalsFor(contextId),
      queryFn: () => lastValueFrom(this.contextService.getIntervals(contextId)),
      enabled: contextId.length > 0,
    };
  }

  list(workspaceId: string | null, includeArchived = false) {
    return {
      queryKey: contextQueryKeys.list(workspaceId, includeArchived),
      queryFn: () => lastValueFrom(this.contextService.getContexts(workspaceId!, includeArchived)),
      enabled: workspaceId !== null && workspaceId.length > 0,
    };
  }

  get(id: string) {
    return {
      queryKey: contextQueryKeys.detail(id),
      queryFn: () => lastValueFrom(this.contextService.getContext(id)),
      enabled: id.length > 0,
    };
  }

  active() {
    return {
      queryKey: contextQueryKeys.active(),
      queryFn: () => lastValueFrom(this.contextService.getActiveContext()),
    };
  }

  stats(contextId: string, date: string, timeZone: string) {
    return {
      queryKey: contextQueryKeys.stat(contextId, date, timeZone),
      queryFn: () => lastValueFrom(this.contextService.getStats(contextId, date, timeZone)),
      enabled: contextId.length > 0 && date.length > 0,
    };
  }

  archiveCandidates(workspaceId: string | null, olderThanDays: number, timeZone: string) {
    return {
      queryKey: contextQueryKeys.archiveCandidates(workspaceId, olderThanDays, timeZone),
      queryFn: () =>
        lastValueFrom(
          this.contextService.getArchiveCandidates(workspaceId!, olderThanDays, timeZone),
        ),
      enabled:
        workspaceId !== null && workspaceId.length > 0 && olderThanDays > 0 && timeZone.length > 0,
    };
  }
}

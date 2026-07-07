import { inject, Inject, Injectable } from '@angular/core';
import { ContextService } from './context.service';
import { lastValueFrom } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class ContextQueries {
  static readonly key = ['contexts'];
  contextService = inject(ContextService);

  intervals(contextId: string) {
    return {
      queryKey: [...ContextQueries.key, 'intervals', contextId],
      queryFn: () => lastValueFrom(this.contextService.getIntervals(contextId)),
      enabled: contextId.length > 0,
    };
  }

  list(workspaceId: string | null) {
    return {
      queryKey: [ContextQueries.key, 'list', workspaceId],
      queryFn: () => lastValueFrom(this.contextService.getContexts(workspaceId!)),
      enabled: workspaceId !== null && workspaceId.length > 0,
    };
  }

  get(id: string) {
    return {
      queryKey: [ContextQueries.key, 'get', id],
      queryFn: () => lastValueFrom(this.contextService.getContext(id)),
    };
  }

  active() {
    return {
      queryKey: [ContextQueries.key, 'active'],
      queryFn: () => lastValueFrom(this.contextService.getActiveContext()),
    };
  }

  stats(contextId: string, date: string) {
    return {
      queryKey: [ContextQueries.key, 'stats', contextId, date],
      queryFn: () => lastValueFrom(this.contextService.getStats(contextId, date)),
    };
  }

  dayStats(workspaceId: string | null, date: string) {
    return {
      queryKey: [ContextQueries.key, 'day-stats', workspaceId, date],
      queryFn: () => lastValueFrom(this.contextService.getDayStats(workspaceId!, date)),
      enabled: workspaceId !== null && workspaceId.length > 0,
    };
  }
}

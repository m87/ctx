import { inject, Injectable } from '@angular/core';
import { lastValueFrom } from 'rxjs';
import { WorkspaceService } from './workspace.service';

export const workspaceQueryKeys = {
  all: ['workspaces'] as const,
  list: () => [...workspaceQueryKeys.all, 'list'] as const,
  details: () => [...workspaceQueryKeys.all, 'get'] as const,
  detail: (workspaceId: string) => [...workspaceQueryKeys.details(), workspaceId] as const,
  stats: () => [...workspaceQueryKeys.all, 'stats'] as const,
  statsFor: (workspaceId: string) => [...workspaceQueryKeys.stats(), workspaceId] as const,
};

@Injectable({
  providedIn: 'root',
})
export class WorkspaceQueries {
  private readonly workspaceService = inject(WorkspaceService);

  list() {
    return {
      queryKey: workspaceQueryKeys.list(),
      queryFn: () => lastValueFrom(this.workspaceService.listWorkspaces()),
    };
  }

  get(id: string) {
    return {
      queryKey: workspaceQueryKeys.detail(id),
      queryFn: () => lastValueFrom(this.workspaceService.getWorkspace(id)),
    };
  }

  stats(id: string) {
    return {
      queryKey: workspaceQueryKeys.statsFor(id),
      queryFn: () => lastValueFrom(this.workspaceService.getWorkspaceStats(id)),
      enabled: id.length > 0,
    };
  }
}

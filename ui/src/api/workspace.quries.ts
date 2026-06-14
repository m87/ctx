import { inject, Injectable } from '@angular/core';
import { lastValueFrom } from 'rxjs';
import { WorkspaceService } from './workspace.service';

@Injectable({
  providedIn: 'root',
})
export class WorkspaceQueries {
  static readonly key = ['workspaces'];
  private workspaceService = inject(WorkspaceService);

  list() {
    return {
      queryKey: [WorkspaceQueries.key, 'list'],
      queryFn: () => lastValueFrom(this.workspaceService.listWorkspaces()),
    };
  }

  get(id: string) {
    return {
      queryKey: [WorkspaceQueries.key, 'get', id],
      queryFn: () => lastValueFrom(this.workspaceService.getWorkspace(id)),
    };
  }

  stats(id: string) {
    return {
      queryKey: [WorkspaceQueries.key, 'stats', id],
      queryFn: () => lastValueFrom(this.workspaceService.getWorkspaceStats(id)),
      enabled: id.length > 0,
    };
  }
}

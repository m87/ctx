import { inject, Injectable } from '@angular/core';
import { mutationOptions, QueryClient } from '@tanstack/angular-query-experimental';
import { lastValueFrom } from 'rxjs';
import { Workspace, WorkspaceService } from './workspace.service';
import { WorkspaceQueries } from './workspace.quries';
import { Router } from '@angular/router';
import { SelectWorkspace } from '../app/sidebar/workspace.state';
import { Store } from '@ngxs/store';
import { toastError } from './error';

@Injectable({
  providedIn: 'root',
})
export class WorkspaceMutations {
  private workspaceService = inject(WorkspaceService);
  private queryClient = inject(QueryClient);
  private router = inject(Router);
  private readonly store = inject(Store);

  create() {
    return mutationOptions({
      mutationFn: (name: string) => lastValueFrom(this.workspaceService.createWorkspace(name)),
      onSuccess: (data) => {
        this.queryClient.invalidateQueries({ queryKey: [WorkspaceQueries.key, 'list'] });
        this.router.navigate(['/workspace', data.id]);
      },
      onError(error) {
        toastError(error);
      },
    });
  }

  update() {
    return mutationOptions({
      mutationFn: (workspace: Workspace) =>
        lastValueFrom(this.workspaceService.updateWorkspace(workspace)),
      onSuccess: (data) => {
        this.queryClient.invalidateQueries({ queryKey: [WorkspaceQueries.key, 'list'] });
        this.queryClient.invalidateQueries({ queryKey: [WorkspaceQueries.key, 'stats', data.id] });
        this.queryClient.invalidateQueries({ queryKey: [WorkspaceQueries.key, 'get', data.id] });
      },
      onError(error) {
        toastError(error);
      },
    });
  }

  delete() {
    return mutationOptions({
      mutationFn: (id: string) => lastValueFrom(this.workspaceService.deleteWorkspace(id)),
      onSuccess: (_, id) => {
        this.queryClient.invalidateQueries({ queryKey: [WorkspaceQueries.key, 'list'] });
        this.queryClient.removeQueries({ queryKey: [WorkspaceQueries.key, 'get', id] });
        this.router.navigate(['/day']);
        this.store.dispatch(new SelectWorkspace(null));
      },
      onError(error) {
        toastError(error);
      },
    });
  }
}

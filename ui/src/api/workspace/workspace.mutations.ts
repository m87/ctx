import { inject, Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { Store } from '@ngxs/store';
import { mutationOptions } from '@tanstack/angular-query-experimental';
import { lastValueFrom } from 'rxjs';
import { SelectWorkspace } from '../../app/sidebar/workspace.state';
import { CacheService } from '../cache/cache.service';
import { Workspace, WorkspaceService } from './workspace.service';

@Injectable({
  providedIn: 'root',
})
export class WorkspaceMutations {
  private readonly workspaceService = inject(WorkspaceService);
  private readonly cache = inject(CacheService);
  private readonly router = inject(Router);
  private readonly store = inject(Store);

  create() {
    return mutationOptions({
      mutationFn: (name: string) => lastValueFrom(this.workspaceService.createWorkspace(name)),
      onSuccess: async (data) => {
        await this.cache.afterWorkspaceListChange();
        await this.router.navigate(['/workspace', data.id]);
      },
    });
  }

  update() {
    return mutationOptions({
      mutationFn: (workspace: Workspace) =>
        lastValueFrom(this.workspaceService.updateWorkspace(workspace)),
      onSuccess: (data) => this.cache.afterWorkspaceUpdate(data.id),
    });
  }

  delete() {
    return mutationOptions({
      mutationFn: (id: string) => lastValueFrom(this.workspaceService.deleteWorkspace(id)),
      onSuccess: async (_, id) => {
        await this.cache.afterWorkspaceDelete(id);
        await this.router.navigate(['/day']);
        this.store.dispatch(new SelectWorkspace(null));
      },
    });
  }
}

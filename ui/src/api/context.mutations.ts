import { inject, Injectable } from '@angular/core';
import { Context, ContextService } from './context.service';
import { mutationOptions, QueryClient } from '@tanstack/angular-query-experimental';
import { Router } from '@angular/router';
import { lastValueFrom } from 'rxjs';
import { ContextQueries } from './context.quries';
import { toastError } from './error';
import { WorkspaceQueries } from './workspace.quries';

@Injectable({
  providedIn: 'root',
})
export class ContextMutations {
  private contextService = inject(ContextService);
  private queryClient = inject(QueryClient);
  private router = inject(Router);

  create() {
    return mutationOptions({
      mutationFn: (context: Context) => lastValueFrom(this.contextService.createContext(context)),
      onSuccess: (data) => {
        this.queryClient.invalidateQueries({ queryKey: [ContextQueries.key, 'list'] });
        this.queryClient.invalidateQueries({ queryKey: [WorkspaceQueries.key, 'stats'] });
        this.router.navigate(['/contexts', data.id]);
      },
      onError(error) {
        toastError(error);
      },
    });
  }

  switch() {
    return mutationOptions({
      mutationFn: (context: Context) => lastValueFrom(this.contextService.switchContext(context)),
      onSuccess: () => {
        this.queryClient.invalidateQueries({ queryKey: [ContextQueries.key, 'list'] });
        this.queryClient.invalidateQueries({ queryKey: [ContextQueries.key, 'active'] });
        this.queryClient.invalidateQueries({ queryKey: ['interval', 'day'] });
        this.queryClient.invalidateQueries({ queryKey: [WorkspaceQueries.key, 'stats'] });
      },
      onError(error) {
        toastError(error);
      },
    });
  }

  free() {
    return mutationOptions({
      mutationFn: () => lastValueFrom(this.contextService.freeContext()),
      onSuccess: () => {
        this.queryClient.invalidateQueries({ queryKey: [ContextQueries.key, 'list'] });
        this.queryClient.invalidateQueries({ queryKey: [ContextQueries.key, 'active'] });
        this.queryClient.invalidateQueries({ queryKey: ['interval', 'day'] });
        this.queryClient.invalidateQueries({ queryKey: [WorkspaceQueries.key, 'stats'] });
      },
      onError(error) {
        toastError(error);
      },
    });
  }

  update() {
    return mutationOptions({
      mutationFn: ({ id, context }: { id: string; context: Context }) =>
        lastValueFrom(this.contextService.updateContext(id, context)),
      onSuccess: (data) => {
        this.queryClient.invalidateQueries({ queryKey: [ContextQueries.key, 'list'] });
        this.queryClient.invalidateQueries({ queryKey: [WorkspaceQueries.key, 'stats'] });
        this.queryClient.invalidateQueries({ queryKey: [ContextQueries.key, 'active'] });
        this.queryClient.invalidateQueries({ queryKey: [ContextQueries.key, 'get', data.id] });
      },
      onError(error) {
        toastError(error);
      },
    });
  }

  delete() {
    return mutationOptions({
      mutationFn: (id: string) => lastValueFrom(this.contextService.deleteContext(id)),
      onSuccess: () => {
        this.queryClient.invalidateQueries({ queryKey: [ContextQueries.key, 'list'] });
        this.queryClient.invalidateQueries({ queryKey: [WorkspaceQueries.key, 'stats'] });
        this.queryClient.invalidateQueries({ queryKey: [ContextQueries.key, 'active'] });
        this.queryClient.invalidateQueries({ queryKey: ['interval', 'day'] });
      },
      onError(error) {
        toastError(error);
      },
    });
  }
}

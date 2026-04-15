import { inject, Injectable } from '@angular/core';
import { Context, ContextService } from './context.service';
import { mutationOptions, QueryClient } from '@tanstack/angular-query-experimental';
import { Router } from '@angular/router';
import { lastValueFrom } from 'rxjs';
import { ContextQueries } from './context.quries';

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
        this.router.navigate(['/contexts', data.id]);
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
      },
    });
  }

  update() {
    return mutationOptions({
      mutationFn: ({ id, context }: { id: string; context: Context }) =>
        lastValueFrom(this.contextService.updateContext(id, context)),
      onSuccess: (data) => {
        this.queryClient.invalidateQueries({ queryKey: [ContextQueries.key, 'list'] });
        this.queryClient.invalidateQueries({ queryKey: [ContextQueries.key, 'active'] });
        this.queryClient.invalidateQueries({ queryKey: [ContextQueries.key, 'get', data.id] });
      },
    });
  }

  delete() {
    return mutationOptions({
      mutationFn: (id: string) => lastValueFrom(this.contextService.deleteContext(id)),
      onSuccess: () => {
        this.queryClient.invalidateQueries({ queryKey: [ContextQueries.key, 'list'] });
        this.queryClient.invalidateQueries({ queryKey: [ContextQueries.key, 'active'] });
        this.queryClient.invalidateQueries({ queryKey: ['interval', 'day'] });
      },
    });
  }
}

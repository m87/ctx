import { Router } from '@angular/router';
import { inject, Injectable } from '@angular/core';
import { mutationOptions } from '@tanstack/angular-query-experimental';
import { lastValueFrom } from 'rxjs';
import { CacheService } from '../cache/cache.service';
import { Context, ContextService, CreateContextInput, SwitchContextInput } from './context.service';

@Injectable({
  providedIn: 'root',
})
export class ContextMutations {
  private readonly contextService = inject(ContextService);
  private readonly cache = inject(CacheService);
  private readonly router = inject(Router);

  create() {
    return mutationOptions({
      mutationFn: (context: CreateContextInput) =>
        lastValueFrom(this.contextService.createContext(context)),
      onSuccess: async (data) => {
        await this.cache.afterContextCreate();
        await this.router.navigate(['/context', data.id]);
      },
    });
  }

  switch() {
    return mutationOptions({
      mutationFn: (context: SwitchContextInput) =>
        lastValueFrom(this.contextService.switchContext(context)),
      onSuccess: () => this.cache.afterActiveIntervalChange(),
    });
  }

  free() {
    return mutationOptions({
      mutationFn: () => lastValueFrom(this.contextService.freeContext()),
      onSuccess: () => this.cache.afterActiveIntervalChange(),
    });
  }

  update() {
    return mutationOptions({
      mutationFn: ({ id, context }: { id: string; context: Context }) =>
        lastValueFrom(this.contextService.updateContext(id, context)),
      onSuccess: (data) => this.cache.afterContextMetadataChange(data.id, true),
    });
  }

  delete() {
    return mutationOptions({
      mutationFn: (id: string) => lastValueFrom(this.contextService.deleteContext(id)),
      onSuccess: (_data, contextId) => this.cache.afterContextDelete(contextId),
    });
  }

  archive() {
    return mutationOptions({
      mutationFn: (contextId: string) =>
        lastValueFrom(this.contextService.archiveContext(contextId)),
      onSuccess: (_data, contextId) => this.cache.afterContextMetadataChange(contextId),
    });
  }

  restore() {
    return mutationOptions({
      mutationFn: (contextId: string) =>
        lastValueFrom(this.contextService.restoreContext(contextId)),
      onSuccess: (_data, contextId) => this.cache.afterContextMetadataChange(contextId),
    });
  }
}

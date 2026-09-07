import { inject, Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { mutationOptions } from '@tanstack/angular-query-experimental';
import { lastValueFrom } from 'rxjs';
import { CacheService } from '../cache/cache.service';
import { CreateSavedQueryInput, SavedQuery, SavedQueryService } from './saved-query.service';

@Injectable({ providedIn: 'root' })
export class SavedQueryMutations {
  private readonly service = inject(SavedQueryService);
  private readonly cache = inject(CacheService);
  private readonly router = inject(Router);

  create() {
    return mutationOptions({
      mutationFn: (input: CreateSavedQueryInput) => lastValueFrom(this.service.create(input)),
      onSuccess: async (query) => {
        await this.cache.afterSavedQueryCreate(query.workspaceId);
        await this.router.navigate(['/query', query.id]);
      },
    });
  }

  delete() {
    return mutationOptions({
      mutationFn: (query: SavedQuery) => lastValueFrom(this.service.delete(query.id)),
      onSuccess: async (_data, query) => {
        await this.cache.afterSavedQueryDelete(query.id, query.workspaceId);
        if (this.router.url.split('?')[0] === `/query/${query.id}`) {
          await this.router.navigate(['/query']);
        }
      },
    });
  }
}

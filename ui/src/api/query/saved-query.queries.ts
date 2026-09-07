import { inject, Injectable } from '@angular/core';
import { lastValueFrom } from 'rxjs';
import { SavedQueryService } from './saved-query.service';

export const savedQueryKeys = {
  all: ['saved-queries'] as const,
  lists: () => [...savedQueryKeys.all, 'list'] as const,
  list: (workspaceId: string | null) => [...savedQueryKeys.lists(), workspaceId] as const,
  details: () => [...savedQueryKeys.all, 'detail'] as const,
  detail: (id: string) => [...savedQueryKeys.details(), id] as const,
};

@Injectable({ providedIn: 'root' })
export class SavedQueryQueries {
  private readonly service = inject(SavedQueryService);

  list(workspaceId: string | null) {
    return {
      queryKey: savedQueryKeys.list(workspaceId),
      queryFn: () => lastValueFrom(this.service.list(workspaceId!)),
      enabled: workspaceId !== null && workspaceId.length > 0,
    };
  }

  get(id: string) {
    return {
      queryKey: savedQueryKeys.detail(id),
      queryFn: () => lastValueFrom(this.service.get(id)),
      enabled: id.length > 0,
    };
  }
}

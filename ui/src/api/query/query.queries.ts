import { inject, Injectable } from '@angular/core';
import { lastValueFrom } from 'rxjs';
import { QueryService } from './query.service';

export const queryKeys = {
  all: ['query'] as const,
  results: () => [...queryKeys.all, 'result'] as const,
  result: (workspaceId: string, query: string) =>
    [...queryKeys.results(), workspaceId, query] as const,
};

@Injectable({ providedIn: 'root' })
export class QueryQueries {
  private readonly queryService = inject(QueryService);

  result(workspaceId: string, query: string, enabled = true) {
    return {
      queryKey: queryKeys.result(workspaceId, query),
      queryFn: () => lastValueFrom(this.queryService.execute({ workspaceId, query })),
      enabled: enabled && workspaceId.length > 0,
    };
  }
}

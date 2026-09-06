import { inject, Injectable } from '@angular/core';
import { lastValueFrom } from 'rxjs';
import { QueryService } from './query.service';

export const queryKeys = {
  all: ['query'] as const,
  results: () => [...queryKeys.all, 'result'] as const,
  result: (query: string) => [...queryKeys.results(), query] as const,
};

@Injectable({ providedIn: 'root' })
export class QueryQueries {
  private readonly queryService = inject(QueryService);

  contexts(query: string) {
    return {
      queryKey: queryKeys.result(query),
      queryFn: () => lastValueFrom(this.queryService.execute({ query })),
    };
  }
}

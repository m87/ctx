import { inject, Injectable } from '@angular/core';
import { mutationOptions } from '@tanstack/angular-query-experimental';
import { lastValueFrom } from 'rxjs';
import { CacheService } from './cache.service';
import { Interval, IntervalService } from './interval.service';

@Injectable({ providedIn: 'root' })
export class IntervalMutations {
  private readonly intervalService = inject(IntervalService);
  private readonly cache = inject(CacheService);

  create() {
    return mutationOptions({
      mutationFn: (interval: Interval) =>
        lastValueFrom(this.intervalService.createInterval(interval)),
      onSuccess: (data) => this.cache.afterIntervalChange(data.contextId ?? ''),
    });
  }

  update() {
    return mutationOptions({
      mutationFn: ({ id, interval }: { id: string; interval: Interval }) =>
        lastValueFrom(this.intervalService.updateInterval(id, interval)),
      onSuccess: (data) => this.cache.afterIntervalChange(data.contextId ?? ''),
    });
  }

  delete() {
    return mutationOptions({
      mutationFn: ({ id, contextId }: { id: string; contextId: string }) =>
        lastValueFrom(this.intervalService.deleteInterval(id)),
      onSuccess: (_, variables) =>
        this.cache.afterIntervalDelete(variables.id, variables.contextId),
    });
  }

  move() {
    return mutationOptions({
      mutationFn: ({ id, targetContextId }: { id: string; targetContextId: string }) =>
        lastValueFrom(this.intervalService.moveInterval(id, targetContextId)),
      onSuccess: (data, variables) =>
        this.cache.afterIntervalChange(data.contextId ?? '', variables.targetContextId),
    });
  }
}

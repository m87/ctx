import { inject, Injectable } from '@angular/core';
import { mutationOptions, QueryClient } from '@tanstack/angular-query-experimental';
import { lastValueFrom } from 'rxjs';
import { ContextQueries } from './context.quries';
import { Interval, IntervalService } from './interval.service';

@Injectable({ providedIn: 'root' })
export class IntervalMutations {
  private intervalService = inject(IntervalService);
  private queryClient = inject(QueryClient);

  create() {
    return mutationOptions({
      mutationFn: (interval: Interval) =>
        lastValueFrom(this.intervalService.createInterval(interval)),
      onSuccess: (data) => {
        this.invalidateAfterIntervalChange(data.contextId ?? data.context_id ?? '');
      },
    });
  }

  update() {
    return mutationOptions({
      mutationFn: ({ id, interval }: { id: string; interval: Interval }) =>
        lastValueFrom(this.intervalService.updateInterval(id, interval)),
      onSuccess: (data) => {
        this.invalidateAfterIntervalChange(data.contextId ?? data.context_id ?? '');
      },
    });
  }

  delete() {
    return mutationOptions({
      mutationFn: ({ id, contextId }: { id: string; contextId: string }) =>
        lastValueFrom(this.intervalService.deleteInterval(id)),
      onSuccess: (_, variables) => {
        this.invalidateAfterIntervalChange(variables.contextId);
      },
    });
  }

  move() {
    return mutationOptions({
      mutationFn: ({ id, targetContextId }: { id: string; targetContextId: string }) =>
        lastValueFrom(this.intervalService.moveInterval(id, targetContextId)),
      onSuccess: (data, variables) => {
        const sourceContextId = data.contextId ?? data.context_id ?? '';
        this.invalidateAfterIntervalChange(sourceContextId);
        this.invalidateAfterIntervalChange(variables.targetContextId);
      },
    });
  }

  private invalidateAfterIntervalChange(contextId: string) {
    this.queryClient.invalidateQueries({ queryKey: ['interval', 'day'] });
    this.queryClient.invalidateQueries({ queryKey: [ContextQueries.key, 'day-stats'] });

    if (contextId) {
      this.queryClient.invalidateQueries({
        queryKey: [ContextQueries.key, 'intervals', contextId],
      });
      this.queryClient.invalidateQueries({ queryKey: [ContextQueries.key, 'stats', contextId] });
      this.queryClient.invalidateQueries({ queryKey: [ContextQueries.key, 'get', contextId] });
    }
  }
}

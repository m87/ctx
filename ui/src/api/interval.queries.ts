import { inject, Injectable } from '@angular/core';
import { IntervalService } from './interval.service';
import { lastValueFrom } from 'rxjs';

export const intervalQueryKeys = {
  all: ['interval'] as const,
  detail: (intervalId: string) => [...intervalQueryKeys.all, intervalId] as const,
  days: () => [...intervalQueryKeys.all, 'day'] as const,
  day: (workspaceId: string | null, date: string) =>
    [...intervalQueryKeys.days(), workspaceId, date] as const,
  dayStats: () => [...intervalQueryKeys.all, 'day-stats'] as const,
  dayStatsFor: (workspaceId: string | null, date: string) =>
    [...intervalQueryKeys.dayStats(), workspaceId, date] as const,
};

@Injectable({ providedIn: 'root' })
export class IntervalQueries {
  private readonly intervalService = inject(IntervalService);

  get(id: string) {
    return {
      queryKey: intervalQueryKeys.detail(id),
      queryFn: () => lastValueFrom(this.intervalService.getInterval(id)),
    };
  }

  day(workspaceId: string | null, day: string) {
    return {
      queryKey: intervalQueryKeys.day(workspaceId, day),
      queryFn: () => lastValueFrom(this.intervalService.getDayIntervals(workspaceId!, day)),
      enabled: workspaceId !== null && workspaceId.length > 0,
    };
  }

  dayStats(workspaceId: string | null, date: string) {
    return {
      queryKey: intervalQueryKeys.dayStatsFor(workspaceId, date),
      queryFn: () => lastValueFrom(this.intervalService.getDayStats(workspaceId!, date)),
      enabled: workspaceId !== null && workspaceId.length > 0,
    };
  }
}

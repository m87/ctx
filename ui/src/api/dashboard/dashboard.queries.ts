import { inject, Injectable } from '@angular/core';
import { lastValueFrom } from 'rxjs';
import { DashboardService, DashboardEvaluationRequest } from './dashboard.service';

export const dashboardQueryKeys = {
  all: ['dashboards'] as const,
  evaluations: () => [...dashboardQueryKeys.all, 'evaluation'] as const,
  lists: () => [...dashboardQueryKeys.all, 'list'] as const,
  list: (workspaceId: string | null) => [...dashboardQueryKeys.lists(), workspaceId] as const,
  details: () => [...dashboardQueryKeys.all, 'detail'] as const,
  detail: (id: string) => [...dashboardQueryKeys.details(), id] as const,
};

@Injectable({ providedIn: 'root' })
export class DashboardQueries {
  private readonly service = inject(DashboardService);

  list(workspaceId: string | null) {
    return {
      queryKey: dashboardQueryKeys.list(workspaceId),
      queryFn: () => lastValueFrom(this.service.list(workspaceId ?? '')),
      enabled: workspaceId !== null && workspaceId.length > 0,
    };
  }

  evaluate(
    request: DashboardEvaluationRequest,
    dataKey: unknown,
    enabled: boolean,
    refetchInterval: number | false = 30_000,
  ) {
    return {
      queryKey: [
        ...dashboardQueryKeys.evaluations(),
        request.workspaceId,
        request.scope,
        request.timeZone,
        request.firstDayOfWeek,
        dataKey,
      ],
      queryFn: () => lastValueFrom(this.service.evaluate(request)),
      enabled: enabled && request.workspaceId.length > 0,
      staleTime: 15_000,
      gcTime: 60_000,
      refetchInterval,
    };
  }

  get(id: string) {
    return {
      queryKey: dashboardQueryKeys.detail(id),
      queryFn: () => lastValueFrom(this.service.get(id)),
      enabled: id.length > 0,
    };
  }
}

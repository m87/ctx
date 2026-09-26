import { inject, Injectable } from '@angular/core';
import { lastValueFrom } from 'rxjs';
import { DashboardService } from './dashboard.service';

export const dashboardQueryKeys = {
  all: ['dashboards'] as const,
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
      queryFn: () => lastValueFrom(this.service.list(workspaceId!)),
      enabled: workspaceId !== null && workspaceId.length > 0,
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

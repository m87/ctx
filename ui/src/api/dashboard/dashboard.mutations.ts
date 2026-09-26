import { inject, Injectable } from '@angular/core';
import { mutationOptions } from '@tanstack/angular-query-experimental';
import { lastValueFrom } from 'rxjs';
import { CacheService } from '../cache/cache.service';
import { CreateDashboardInput, Dashboard, DashboardService } from './dashboard.service';

@Injectable({ providedIn: 'root' })
export class DashboardMutations {
  private readonly service = inject(DashboardService);
  private readonly cache = inject(CacheService);

  create() {
    return mutationOptions({
      mutationFn: (input: CreateDashboardInput) => lastValueFrom(this.service.create(input)),
      onSuccess: (dashboard) => this.cache.afterDashboardCreate(dashboard.workspaceId),
    });
  }

  update() {
    return mutationOptions({
      mutationFn: (dashboard: Dashboard) => lastValueFrom(this.service.update(dashboard)),
      onSuccess: (dashboard) => this.cache.afterDashboardUpdate(dashboard),
    });
  }

  delete() {
    return mutationOptions({
      mutationFn: (dashboard: Dashboard) => lastValueFrom(this.service.delete(dashboard.id)),
      onSuccess: (_data, dashboard) => this.cache.afterDashboardDelete(dashboard),
    });
  }
}

import { Component, inject } from '@angular/core';
import { toSignal } from '@angular/core/rxjs-interop';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { map } from 'rxjs';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { HlmButtonImports } from '@spartan-ng/helm/button';
import { HlmSkeletonImports } from '@spartan-ng/helm/skeleton';
import { DashboardQueries } from '../../api/dashboard/dashboard.queries';
import { QueryErrorStateComponent } from '../shared/query-error-state.component';
import { DashboardViewComponent } from './dashboard-view.component';

@Component({
  selector: 'ctx-dashboard-page',
  imports: [
    DashboardViewComponent,
    QueryErrorStateComponent,
    RouterLink,
    HlmButtonImports,
    HlmSkeletonImports,
  ],
  host: { class: 'block h-full min-h-0 w-full' },
  template: `
    <div class="flex h-full min-h-0 flex-col">
      @if (dashboardQuery.isLoading()) {
        <div class="flex flex-col gap-4 p-4 md:p-6" role="status" aria-label="Loading dashboard">
          <hlm-skeleton class="h-6 w-48"></hlm-skeleton>
          <hlm-skeleton class="h-32 w-full"></hlm-skeleton>
        </div>
      } @else if (dashboardQuery.isError() || dashboardQuery.isPaused()) {
        <ctx-query-error-state
          [error]="dashboardQuery.error()"
          [paused]="dashboardQuery.isPaused()"
          [retrying]="dashboardQuery.isFetching()"
          resourceName="dashboard"
          (retry)="dashboardQuery.refetch()"
        ></ctx-query-error-state>
      } @else if (dashboardQuery.data(); as dashboard) {
        <header class="ui-page-header shrink-0 border-b border-border/65 p-4 md:p-6">
          <div class="min-w-0">
            <div class="ui-page-kicker">Lens</div>
            <h1 class="ui-page-title break-words">{{ dashboard.name }}</h1>
          </div>
          <a
            hlmBtn
            variant="outline"
            size="sm"
            routerLink="/editor"
            [queryParams]="{ dashboard: dashboard.id }"
            >Edit</a
          >
        </header>
        <div class="min-h-0 flex-1 overflow-auto bg-muted dark:bg-background">
          <ctx-dashboard-view [dashboard]="dashboard"></ctx-dashboard-view>
        </div>
      }
    </div>
  `,
})
export class DashboardPageComponent {
  private readonly route = inject(ActivatedRoute);
  private readonly queries = inject(DashboardQueries);
  readonly dashboardId = toSignal(
    this.route.paramMap.pipe(map((params) => params.get('id') ?? '')),
    { initialValue: '' },
  );
  readonly dashboardQuery = injectQuery(() => this.queries.get(this.dashboardId()));
}

import { Component, computed, inject, input } from '@angular/core';
import { RouterLink } from '@angular/router';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { HlmButtonImports } from '@spartan-ng/helm/button';
import { HlmSkeletonImports } from '@spartan-ng/helm/skeleton';
import { DashboardQueries } from '../../api/dashboard/dashboard.queries';
import { DashboardType } from '../../api/dashboard/dashboard.service';
import { QueryErrorStateComponent } from '../shared/query-error-state.component';
import { DashboardViewComponent } from './dashboard-view.component';

@Component({
  selector: 'ctx-dashboard-slot',
  imports: [
    DashboardViewComponent,
    QueryErrorStateComponent,
    HlmSkeletonImports,
    HlmButtonImports,
    RouterLink,
  ],
  host: {
    class:
      'block h-full min-h-0 w-full flex-1 overflow-auto bg-muted text-foreground dark:bg-background',
  },
  template: `
    @if (dashboardsQuery.isLoading()) {
      <div class="flex flex-col gap-3 p-4" role="status" aria-label="Loading dashboard">
        <hlm-skeleton class="h-6 w-48"></hlm-skeleton>
        <hlm-skeleton class="h-32 w-full"></hlm-skeleton>
      </div>
    } @else if (dashboardsQuery.isError() || dashboardsQuery.isPaused()) {
      <ctx-query-error-state
        [error]="dashboardsQuery.error()"
        [paused]="dashboardsQuery.isPaused()"
        [retrying]="dashboardsQuery.isFetching()"
        resourceName="dashboard"
        (retry)="dashboardsQuery.refetch()"
      ></ctx-query-error-state>
    } @else if (dashboard(); as dashboard) {
      <ctx-dashboard-view [dashboard]="dashboard"></ctx-dashboard-view>
    } @else {
      <div class="ui-empty-state flex flex-col items-center gap-3 py-10">
        <p>No active dashboard for this view.</p>
        <a hlmBtn variant="outline" size="sm" routerLink="/editor">Open dashboard editor</a>
      </div>
    }
  `,
})
export class DashboardSlotComponent {
  readonly workspaceId = input.required<string | null>();
  readonly type = input.required<Exclude<DashboardType, 'custom'>>();
  private readonly queries = inject(DashboardQueries);
  readonly dashboardsQuery = injectQuery(() => this.queries.list(this.workspaceId()));
  readonly dashboard = computed(() =>
    (this.dashboardsQuery.data() ?? []).find(
      (dashboard) => dashboard.active && dashboard.type === this.type(),
    ),
  );
}

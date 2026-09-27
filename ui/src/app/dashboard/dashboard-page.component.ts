import { Component, inject } from '@angular/core';
import { toSignal } from '@angular/core/rxjs-interop';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { map } from 'rxjs';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucidePencil } from '@ng-icons/lucide';
import { HlmButtonImports } from '@spartan-ng/helm/button';
import { HlmSkeletonImports } from '@spartan-ng/helm/skeleton';
import { DashboardQueries } from '../../api/dashboard/dashboard.queries';
import { NameComponent } from '../shared/name.component';
import { QueryErrorStateComponent } from '../shared/query-error-state.component';
import { DashboardViewComponent } from './dashboard-view.component';

@Component({
  selector: 'ctx-dashboard-page',
  imports: [
    DashboardViewComponent,
    NameComponent,
    QueryErrorStateComponent,
    HlmButtonImports,
    HlmSkeletonImports,
    NgIcon,
    RouterLink,
  ],
  providers: [provideIcons({ lucidePencil })],
  host: { class: 'block h-full min-h-0 w-full min-w-0 max-w-[1000px]' },
  template: `
    <div class="flex h-full min-h-0 w-full flex-col overflow-hidden p-4 md:p-6">
      @if (dashboardQuery.isLoading()) {
        <div class="flex flex-col gap-4" role="status" aria-label="Loading dashboard">
          <hlm-skeleton class="h-6 w-48"></hlm-skeleton>
          <hlm-skeleton class="h-32 w-full"></hlm-skeleton>
        </div>
      } @else if (dashboardQuery.isError() || dashboardQuery.isPaused()) {
        <ctx-query-error-state
          class="min-h-0 flex-1"
          [error]="dashboardQuery.error()"
          [paused]="dashboardQuery.isPaused()"
          [retrying]="dashboardQuery.isFetching()"
          resourceName="dashboard"
          (retry)="dashboardQuery.refetch()"
        ></ctx-query-error-state>
      } @else if (dashboardQuery.data(); as dashboard) {
        <div class="mb-5 flex shrink-0 items-start justify-between gap-4">
          <ctx-name
            class="w-0 min-w-0 flex-1"
            label="Lens"
            [name]="dashboard.name"
            [showDescription]="false"
            [readonly]="true"
          ></ctx-name>
          <a
            hlmBtn
            variant="outline"
            size="icon"
            routerLink="/editor"
            [queryParams]="{ dashboard: dashboard.id }"
            class="mt-5 shrink-0 text-muted-foreground hover:text-foreground"
            aria-label="Edit dashboard"
            title="Edit"
          >
            <ng-icon name="lucidePencil"></ng-icon>
          </a>
        </div>
        <div class="-mx-2.5 min-h-0 flex-1 overflow-auto bg-muted dark:bg-background">
          <ctx-dashboard-view class="px-0 sm:px-0" [dashboard]="dashboard"></ctx-dashboard-view>
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

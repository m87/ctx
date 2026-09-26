import { Component, computed, inject, signal } from '@angular/core';
import { RouterLink, RouterLinkActive } from '@angular/router';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideChevronDown, lucidePanelsTopLeft } from '@ng-icons/lucide';
import { Store } from '@ngxs/store';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { HlmSkeletonImports } from '@spartan-ng/helm/skeleton';
import { DashboardQueries } from '../../api/dashboard/dashboard.queries';
import { SidebarStore } from './sidebar.store';
import { WorkspaceState } from './workspace.state';

@Component({
  selector: 'ctx-sidebar-lenses',
  imports: [RouterLink, RouterLinkActive, NgIcon, HlmSkeletonImports],
  providers: [provideIcons({ lucideChevronDown, lucidePanelsTopLeft })],
  host: { class: 'block min-w-0' },
  template: `
    @if (
      dashboardsQuery.isLoading() ||
      dashboards().length ||
      dashboardsQuery.isError() ||
      dashboardsQuery.isPaused()
    ) {
      <section aria-labelledby="lenses-title">
        <button
          type="button"
          class="sticky top-0 z-10 flex min-w-0 w-full items-center justify-between gap-2 bg-sidebar px-4 py-2 text-xs font-semibold text-muted-foreground transition-colors hover:bg-muted hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring sidebar-collapsed:justify-center sidebar-collapsed:px-2"
          [attr.aria-expanded]="expanded()"
          aria-controls="lenses-list"
          (click)="expanded.set(!expanded())"
        >
          <span
            id="lenses-title"
            class="min-w-0 truncate uppercase tracking-label sidebar-collapsed:hidden"
            >Lenses</span
          >
          <ng-icon
            name="lucideChevronDown"
            class="shrink-0 text-xs transition-transform"
            [class.-rotate-90]="!expanded()"
          ></ng-icon>
        </button>
        @if (expanded()) {
          <div id="lenses-list" class="flex flex-col gap-0.5 px-2 pb-2">
            @if (dashboardsQuery.isLoading()) {
              <div class="px-2 py-1" role="status" aria-label="Loading lenses">
                <hlm-skeleton class="h-6 w-full"></hlm-skeleton>
              </div>
            } @else if (dashboardsQuery.isError() || dashboardsQuery.isPaused()) {
              <button
                type="button"
                class="rounded-md px-2 py-1.5 text-left text-xs text-muted-foreground hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                [disabled]="dashboardsQuery.isFetching()"
                (click)="dashboardsQuery.refetch()"
              >
                Couldn't load lenses. Retry
              </button>
            } @else {
              @for (dashboard of dashboards(); track dashboard.id) {
                <a
                  [routerLink]="['/dashboard', dashboard.id]"
                  routerLinkActive="bg-muted text-foreground"
                  class="flex min-w-0 items-center gap-2 rounded-md px-2 py-1.5 text-dense font-medium text-muted-foreground transition-colors hover:bg-muted/60 hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                  [title]="dashboard.name"
                  (click)="sidebar.closeMobile()"
                >
                  <ng-icon name="lucidePanelsTopLeft" class="shrink-0 text-xs"></ng-icon>
                  <span class="truncate">{{ dashboard.name }}</span>
                </a>
              }
            }
          </div>
        }
      </section>
    }
  `,
})
export class SidebarLensesComponent {
  private readonly store = inject(Store);
  private readonly queries = inject(DashboardQueries);
  readonly sidebar = inject(SidebarStore);
  readonly expanded = signal(true);
  readonly workspaceId = this.store.selectSignal(WorkspaceState.selectedWorkspaceId);
  readonly dashboardsQuery = injectQuery(() => this.queries.list(this.workspaceId()));
  readonly dashboards = computed(() =>
    (this.dashboardsQuery.data() ?? []).filter(
      (dashboard) => dashboard.type === 'custom' && dashboard.active,
    ),
  );
}

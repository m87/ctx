import { Component, computed, inject, signal } from '@angular/core';
import { RouterLink, RouterLinkActive } from '@angular/router';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideChevronDown, lucideSearch } from '@ng-icons/lucide';
import { Store } from '@ngxs/store';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { SavedQueryQueries } from '../../api/query/saved-query.queries';
import { HlmSkeletonImports } from '@spartan-ng/helm/skeleton';
import { SidebarStore } from './sidebar.store';
import { WorkspaceState } from './workspace.state';

export function isSavedQuerySectionVisible(loading: boolean, queryCount: number): boolean {
  return loading || queryCount > 0;
}

@Component({
  selector: 'ctx-sidebar-saved-query-list',
  imports: [HlmSkeletonImports, NgIcon, RouterLink, RouterLinkActive],
  providers: [provideIcons({ lucideChevronDown, lucideSearch })],
  host: {
    class: 'block min-w-0',
  },
  template: `
    @if (visible()) {
      <section aria-labelledby="saved-queries-title">
        <button
          type="button"
          class="sticky top-0 z-10 flex min-w-0 w-full items-center justify-between gap-2 bg-sidebar px-4 py-2 text-xs font-semibold text-muted-foreground transition-colors hover:bg-muted hover:text-foreground sidebar-collapsed:justify-center sidebar-collapsed:px-2"
          [attr.aria-expanded]="expanded()"
          aria-controls="saved-queries-list"
          (click)="expanded.set(!expanded())"
        >
          <span
            id="saved-queries-title"
            class="min-w-0 truncate uppercase tracking-[0.08em] sidebar-collapsed:hidden"
            >Queries</span
          >
          <ng-icon
            name="lucideChevronDown"
            class="shrink-0 text-xs transition-transform"
            [class.-rotate-90]="!expanded()"
          ></ng-icon>
        </button>

        @if (expanded()) {
          <div id="saved-queries-list" class="flex flex-col gap-0.5 px-2 pb-2">
            @if (savedQueriesQuery.isLoading()) {
              <div class="flex flex-col gap-2 px-2 py-1" role="status">
                <span class="sr-only">Loading saved queries</span>
                @for (item of skeletonItems; track item) {
                  <hlm-skeleton class="h-6 w-full" [class.max-w-36]="item === 1"></hlm-skeleton>
                }
              </div>
            } @else {
              @for (query of queries(); track query.id) {
                <a
                  [routerLink]="['/query', query.id]"
                  routerLinkActive="bg-muted text-foreground"
                  class="flex min-w-0 items-center gap-2 rounded-md px-2 py-1.5 text-[13px] font-medium text-muted-foreground transition-colors hover:bg-muted/60 hover:text-foreground"
                  [title]="query.name"
                  (click)="sidebar.closeMobile()"
                >
                  <ng-icon name="lucideSearch" class="shrink-0 text-xs"></ng-icon>
                  <span class="truncate">{{ query.name }}</span>
                </a>
              }
            }
          </div>
        }
      </section>
    }
  `,
})
export class SidebarSavedQueryListComponent {
  readonly skeletonItems = [0, 1];
  readonly expanded = signal(false);
  private readonly store = inject(Store);
  private readonly savedQueryQueries = inject(SavedQueryQueries);
  readonly sidebar = inject(SidebarStore);

  readonly workspaceId = this.store.selectSignal(WorkspaceState.selectedWorkspaceId);
  readonly savedQueriesQuery = injectQuery(() => this.savedQueryQueries.list(this.workspaceId()));
  readonly queries = computed(() => this.savedQueriesQuery.data() ?? []);
  readonly visible = computed(() =>
    isSavedQuerySectionVisible(this.savedQueriesQuery.isLoading(), this.queries().length),
  );
}

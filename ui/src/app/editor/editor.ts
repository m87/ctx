import { Component, computed, effect, inject, signal } from '@angular/core';
import { RouterLink } from '@angular/router';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideArrowLeft, lucideEye, lucideGanttChart } from '@ng-icons/lucide';
import { BrnDialogImports } from '@spartan-ng/brain/dialog';
import { HlmButtonImports } from '@spartan-ng/helm/button';
import { HlmDialogImports } from '@spartan-ng/helm/dialog';
import { HlmSkeletonImports } from '@spartan-ng/helm/skeleton';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { Store } from '@ngxs/store';
import { ContextQueries } from '../../api/context/context.queries';
import { ProjectQueries } from '../../api/project/project.queries';
import { QueryQueries } from '../../api/query/query.queries';
import { ContextQueryResult } from '../../api/query/query.service';
import { SavedQueryQueries } from '../../api/query/saved-query.queries';
import { SavedQuery } from '../../api/query/saved-query.service';
import { ContextListComponent } from '../context/context-list.component';
import { contextQueryResultAsListItems } from '../query/query-summary.component';
import { QueryErrorStateComponent } from '../shared/query-error-state.component';
import { SearchDropdownSelectComponent } from '../shared/search-dropdown-select.component';
import { SearchSelectComponent, SearchSelectOption } from '../shared/search-select.component';
import { SidebarWorkspaceSelectComponent } from '../sidebar/sidebar-workspace-select.component';
import { WorkspaceState } from '../sidebar/workspace.state';
import { colorHash, durationAsHM } from '../utils';

export type QueryScope = 'workspace' | 'project' | 'context' | 'daily';

export const EDITOR_QUERY_SCOPE_OPTIONS: readonly SearchSelectOption[] = [
  { value: 'workspace', label: 'Workspace' },
  { value: 'project', label: 'Project' },
  { value: 'context', label: 'Context' },
  { value: 'daily', label: 'Daily' },
];

export function queryScopeHasEntity(scope: QueryScope): boolean {
  return scope === 'project' || scope === 'context';
}

export function savedQueriesAsOptions(queries: readonly SavedQuery[]): SearchSelectOption[] {
  return queries.map((query) => ({
    value: query.id,
    label: query.name,
    description: query.query,
    keywords: [query.query],
  }));
}

export function selectedSavedQueryText(
  queries: readonly SavedQuery[],
  selectedQueryId: string,
): string {
  return queries.find((query) => query.id === selectedQueryId)?.query ?? '';
}

export function queryPreviewSummary(result: ContextQueryResult | undefined): string {
  if (!result) {
    return '';
  }

  const contextCount = result.contexts.length;
  const contextLabel = contextCount === 1 ? 'context' : 'contexts';
  const duration = durationAsHM(result.totalDuration).trim() || '0m';
  return `${contextCount} ${contextLabel} · ${duration} · ${result.totalSessions} sessions`;
}

@Component({
  selector: 'ctx-editor',
  imports: [
    BrnDialogImports,
    HlmButtonImports,
    HlmDialogImports,
    HlmSkeletonImports,
    NgIcon,
    ContextListComponent,
    QueryErrorStateComponent,
    RouterLink,
    SearchDropdownSelectComponent,
    SearchSelectComponent,
    SidebarWorkspaceSelectComponent,
  ],
  providers: [provideIcons({ lucideArrowLeft, lucideEye, lucideGanttChart })],
  template: `
    <div class="fixed inset-0 z-50 flex h-dvh min-h-0 flex-col bg-background text-foreground">
      <header class="flex h-12 shrink-0 items-center border-b bg-card/70 px-3">
        <div class="flex min-w-0 items-center gap-2">
          <a
            routerLink="/"
            class="flex size-8 shrink-0 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-muted hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
            aria-label="Back to application"
            title="Back to application"
          >
            <ng-icon name="lucideArrowLeft"></ng-icon>
          </a>

          <div class="h-5 w-px bg-border"></div>

          <div class="flex min-w-0 items-center gap-2 pl-1">
            <ng-icon name="lucideGanttChart" class="shrink-0 text-primary"></ng-icon>
            <span class="font-semibold tracking-tight text-primary">Ctx</span>
            <span class="text-muted-foreground">/</span>
            <span class="truncate text-sm font-medium">Editor</span>
          </div>

          <div class="ml-1 h-5 w-px bg-border"></div>

          <ctx-sidebar-workspace-select class="ml-1 w-48"></ctx-sidebar-workspace-select>
        </div>

        <div class="ml-auto flex shrink-0 items-center gap-2 pl-4">
          <ctx-search-dropdown-select
            class="w-32"
            inputId="editor-query-scope"
            ariaLabel="Query scope"
            actionLabel="Add custom"
            panelWidth="100%"
            align="end"
            [searchable]="false"
            [options]="queryScopeOptions"
            [value]="queryScope()"
            (selectionChange)="setQueryScope($event)"
          ></ctx-search-dropdown-select>

          @if (showEntitySelect()) {
            <ctx-search-dropdown-select
              class="w-52"
              inputId="editor-entity-select"
              align="end"
              [ariaLabel]="entitySelectAriaLabel()"
              [placeholder]="entitySelectPlaceholder()"
              [searchPlaceholder]="entitySearchPlaceholder()"
              [emptyText]="entityEmptyText()"
              [options]="entityOptions()"
              [value]="selectedEntityId()"
              (selectionChange)="selectEntity($event)"
            ></ctx-search-dropdown-select>
          }

          <label
            class="flex h-8 items-center rounded-md border border-border/60 bg-muted/30 transition-[background-color,border-color,box-shadow] hover:border-border hover:bg-muted/50 focus-within:border-ring/70 focus-within:ring-2 focus-within:ring-ring/30"
          >
            <span
              class="pl-2.5 text-caption font-medium uppercase tracking-label text-muted-foreground"
            >
              Days
            </span>
            <input
              type="number"
              min="1"
              step="1"
              class="days-input h-full w-14 bg-transparent px-2 text-right text-xs outline-none"
              aria-label="Number of days"
              [value]="days()"
              (input)="setDays($event)"
            />
          </label>
        </div>
      </header>

      <div class="flex min-h-0 flex-1">
        <aside
          class="hidden w-60 shrink-0 flex-col border-r bg-sidebar md:flex"
          aria-label="Editor tools"
        >
          <section class="flex min-h-0 flex-1 flex-col border-b" aria-label="Query">
            <div class="flex items-center justify-between gap-2 border-b px-3 py-2">
              <div class="text-meta font-semibold uppercase tracking-label text-muted-foreground">
                Query
              </div>

              <hlm-dialog>
                <button
                  hlmBtn
                  hlmDialogTrigger
                  type="button"
                  variant="ghost"
                  size="icon-sm"
                  class="text-muted-foreground"
                  aria-label="Preview query"
                  title="Preview query"
                  [disabled]="!selectedWorkspaceId()"
                  (click)="loadQueryPreview()"
                >
                  <ng-icon name="lucideEye"></ng-icon>
                </button>

                <hlm-dialog-content
                  class="editor-preview-dialog flex h-[calc(100dvh-2rem)] max-h-[calc(100dvh-2rem)] flex-col gap-0 overflow-hidden p-0 sm:h-[min(44rem,calc(100dvh-4rem))]"
                  *brnDialogContent
                >
                  <div hlmDialogHeader class="shrink-0 border-b p-5 pr-12">
                    <h2 hlmDialogTitle>Query preview</h2>
                    <p hlmDialogDescription>
                      {{ previewSummary() || 'Contexts returned for the selected query.' }}
                    </p>
                  </div>
                  <div class="min-h-0 flex-1 overflow-y-auto p-5">
                    @if (previewShowError()) {
                      <ctx-query-error-state
                        [error]="queryPreview.error()"
                        [paused]="queryPreview.isPaused()"
                        resourceName="query preview"
                        [retrying]="queryPreview.isFetching()"
                        (retry)="loadQueryPreview()"
                      ></ctx-query-error-state>
                    } @else if (queryPreview.isLoading()) {
                      <div class="flex flex-col gap-2" role="status" aria-label="Loading contexts">
                        <span class="sr-only">Loading contexts</span>
                        @for (item of previewSkeletonItems; track item) {
                          <div class="rounded-lg border bg-card p-3">
                            <div class="mb-3 flex items-center gap-2">
                              <hlm-skeleton class="size-2 shrink-0"></hlm-skeleton>
                              <hlm-skeleton class="h-3.5 w-2/5"></hlm-skeleton>
                              <hlm-skeleton class="ml-auto h-3 w-12"></hlm-skeleton>
                            </div>
                            <hlm-skeleton class="mb-2 h-1.5 w-full"></hlm-skeleton>
                            <hlm-skeleton class="h-2.5 w-24"></hlm-skeleton>
                          </div>
                        }
                      </div>
                    } @else {
                      <ctx-context-list
                        [items]="previewContexts()"
                        emptyMessage="No contexts match this query."
                      ></ctx-context-list>
                    }
                  </div>
                </hlm-dialog-content>
              </hlm-dialog>
            </div>

            <div class="flex flex-1 items-start gap-3 overflow-y-auto bg-background/40 p-3">
              <div class="font-mono text-xs text-muted-foreground/50" aria-hidden="true">1</div>
              <div
                class="min-w-0 whitespace-pre-wrap break-words font-mono text-xs text-foreground"
              >
                {{ queryText() }}
              </div>
            </div>
          </section>

          <section class="flex min-h-0 flex-1 flex-col" aria-labelledby="saved-queries-title">
            <div class="border-b px-3 py-3">
              <div id="saved-queries-title" class="ui-section-label">Saved queries</div>
            </div>

            <ctx-search-select
              class="min-h-0"
              inputId="editor-saved-query"
              ariaLabel="Saved queries"
              searchPlaceholder="Search saved queries…"
              [emptyText]="savedQueryEmptyText()"
              [disabled]="savedQueriesQuery.isLoading()"
              [options]="savedQueryOptions()"
              [value]="selectedSavedQueryId()"
              [embedded]="true"
              (selectionChange)="selectSavedQuery($event)"
            ></ctx-search-select>
          </section>
        </aside>

        <main
          class="editor-canvas flex min-w-0 flex-1 items-center justify-center overflow-auto p-6"
        >
          <div
            class="aspect-[16/10] w-full max-w-[1000px] rounded-lg border bg-card shadow-sm"
            aria-label="Empty editor canvas"
          ></div>
        </main>

        <aside
          class="hidden w-72 shrink-0 flex-col border-l bg-sidebar lg:flex"
          aria-label="Widget tools"
        >
          <section class="flex min-h-0 flex-1 flex-col border-b" aria-label="Widgets">
            <div class="border-b px-3 py-3">
              <div class="text-meta font-semibold uppercase tracking-label text-muted-foreground">
                Widgets
              </div>
            </div>

            <div class="flex flex-1 items-center justify-center p-5">
              <div class="max-w-44 text-center">
                <div class="text-xs font-medium text-foreground/80">No widgets yet</div>
                <div class="mt-1 text-meta leading-relaxed text-muted-foreground">
                  Dashboard widgets will appear here.
                </div>
              </div>
            </div>
          </section>

          <section class="flex min-h-0 flex-1 flex-col" aria-label="Properties">
            <div class="border-b px-3 py-3">
              <div class="text-meta font-semibold uppercase tracking-label text-muted-foreground">
                Properties
              </div>
            </div>

            <div class="flex flex-1 items-center justify-center p-5">
              <div class="max-w-44 text-center">
                <div class="text-xs font-medium text-foreground/80">Nothing selected</div>
                <div class="mt-1 text-meta leading-relaxed text-muted-foreground">
                  Select an element to view its properties.
                </div>
              </div>
            </div>
          </section>
        </aside>
      </div>
    </div>
  `,
  styles: `
    :host {
      display: block;
    }

    .editor-canvas {
      background-color: color-mix(in oklab, var(--muted) 55%, var(--background));
      background-image: radial-gradient(
        circle,
        color-mix(in oklab, var(--muted-foreground) 22%, transparent) 1px,
        transparent 1px
      );
      background-size: 20px 20px;
    }

    .editor-preview-dialog {
      width: calc(100vw - 2rem);
      max-width: 72rem;
    }

    .days-input {
      appearance: textfield;
      -moz-appearance: textfield;
    }

    .days-input::-webkit-inner-spin-button,
    .days-input::-webkit-outer-spin-button {
      margin: 0;
      appearance: none;
      -webkit-appearance: none;
    }
  `,
})
export class EditorComponent {
  readonly previewSkeletonItems = [0, 1, 2];
  private readonly projectQueries = inject(ProjectQueries);
  private readonly contextQueries = inject(ContextQueries);
  private readonly queryQueries = inject(QueryQueries);
  private readonly savedQueryQueries = inject(SavedQueryQueries);
  private readonly store = inject(Store);

  readonly queryScopeOptions = EDITOR_QUERY_SCOPE_OPTIONS;
  readonly queryScope = signal<QueryScope>('workspace');
  readonly selectedEntityId = signal('');
  readonly selectedSavedQueryId = signal('');
  readonly days = signal(30);
  readonly selectedWorkspaceId = this.store.selectSignal(WorkspaceState.selectedWorkspaceId);

  readonly projectsQuery = injectQuery(() =>
    this.projectQueries.all(this.selectedWorkspaceId() ?? ''),
  );
  readonly contextsQuery = injectQuery(() => this.contextQueries.list(this.selectedWorkspaceId()));
  readonly savedQueriesQuery = injectQuery(() =>
    this.savedQueryQueries.list(this.selectedWorkspaceId()),
  );

  private readonly resetEntityOnWorkspaceChange = effect(() => {
    this.selectedWorkspaceId();
    this.selectedEntityId.set('');
    this.selectedSavedQueryId.set('');
  });

  readonly projectOptions = computed<SearchSelectOption[]>(() =>
    (this.projectsQuery.data() ?? []).map((project) => ({
      value: project.id,
      label: project.name,
      color: colorHash(`project:${project.id}`),
      description: 'Project',
    })),
  );
  readonly contextOptions = computed<SearchSelectOption[]>(() =>
    (this.contextsQuery.data() ?? []).map((context) => ({
      value: context.id,
      label: context.name,
      color: colorHash(context.id),
      badge: context.project?.name,
      keywords: context.project?.name ? [context.project.name] : [],
    })),
  );
  readonly savedQueryOptions = computed<SearchSelectOption[]>(() =>
    savedQueriesAsOptions(this.savedQueriesQuery.data() ?? []),
  );
  readonly queryText = computed(() =>
    selectedSavedQueryText(this.savedQueriesQuery.data() ?? [], this.selectedSavedQueryId()),
  );
  readonly queryPreview = injectQuery(() =>
    this.queryQueries.result(this.selectedWorkspaceId() ?? '', this.queryText(), false),
  );
  readonly previewSummary = computed(() => queryPreviewSummary(this.queryPreview.data()));
  readonly previewContexts = computed(() =>
    contextQueryResultAsListItems(this.queryPreview.data()),
  );
  readonly previewShowError = computed(
    () =>
      this.queryPreview.data() === undefined &&
      (this.queryPreview.isError() || this.queryPreview.isPaused()),
  );
  readonly savedQueryEmptyText = computed(() => {
    if (this.savedQueriesQuery.isLoading()) {
      return 'Loading saved queries…';
    }
    if (this.savedQueriesQuery.isError()) {
      return 'Unable to load saved queries';
    }
    return 'No saved queries';
  });
  readonly showEntitySelect = computed(() => queryScopeHasEntity(this.queryScope()));
  readonly entityOptions = computed<SearchSelectOption[]>(() => {
    switch (this.queryScope()) {
      case 'workspace':
        return [];
      case 'project':
        return this.projectOptions();
      case 'context':
        return this.contextOptions();
      case 'daily':
        return [];
    }
  });
  readonly entityOptionsLoading = computed(() => {
    switch (this.queryScope()) {
      case 'workspace':
        return false;
      case 'project':
        return this.projectsQuery.isLoading();
      case 'context':
        return this.contextsQuery.isLoading();
      case 'daily':
        return false;
    }
  });

  readonly entityName = computed(() => this.queryScope());
  readonly entityPluralName = computed(() => `${this.entityName()}s`);
  readonly entitySelectAriaLabel = computed(() => `Select ${this.entityName()}`);
  readonly entitySelectPlaceholder = computed(() => `Select ${this.entityName()}…`);
  readonly entitySearchPlaceholder = computed(() => `Search ${this.entityPluralName()}…`);
  readonly entityEmptyText = computed(() =>
    this.entityOptionsLoading()
      ? `Loading ${this.entityPluralName()}…`
      : `No matching ${this.entityPluralName()}`,
  );
  setQueryScope(scope: string): void {
    this.queryScope.set(scope as QueryScope);
    this.selectedEntityId.set('');
  }

  selectEntity(entityId: string): void {
    this.selectedEntityId.set(entityId);
  }

  selectSavedQuery(queryId: string): void {
    this.selectedSavedQueryId.set(queryId);
  }

  loadQueryPreview(): void {
    if (this.selectedWorkspaceId()) {
      void this.queryPreview.refetch();
    }
  }

  setDays(event: Event): void {
    const value = (event.target as HTMLInputElement).valueAsNumber;
    this.days.set(Number.isFinite(value) && value > 0 ? Math.floor(value) : 30);
  }
}

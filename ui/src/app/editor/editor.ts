import { Component, computed, effect, inject, signal } from '@angular/core';
import { RouterLink } from '@angular/router';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideArrowLeft, lucideGanttChart } from '@ng-icons/lucide';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { Store } from '@ngxs/store';
import { ContextQueries } from '../../api/context/context.queries';
import { ProjectQueries } from '../../api/project/project.queries';
import { SearchDropdownSelectComponent } from '../shared/search-dropdown-select.component';
import { SearchSelectOption } from '../shared/search-select.component';
import { SidebarWorkspaceSelectComponent } from '../sidebar/sidebar-workspace-select.component';
import { WorkspaceState } from '../sidebar/workspace.state';
import { colorHash } from '../utils';

type QueryScope = 'workspace' | 'project' | 'context' | 'daily';

@Component({
  selector: 'ctx-editor',
  imports: [NgIcon, RouterLink, SearchDropdownSelectComponent, SidebarWorkspaceSelectComponent],
  providers: [provideIcons({ lucideArrowLeft, lucideGanttChart })],
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

          @if (queryScope() === 'project' || queryScope() === 'context') {
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
              class="pl-2.5 text-[10px] font-medium uppercase tracking-[0.06em] text-muted-foreground"
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
            <div class="border-b px-3 py-3">
              <div
                class="text-[11px] font-semibold uppercase tracking-[0.08em] text-muted-foreground"
              >
                Query
              </div>
            </div>

            <div class="flex flex-1 items-start bg-background/40 p-3">
              <div class="font-mono text-xs text-muted-foreground/50" aria-hidden="true">1</div>
            </div>
          </section>

          <section class="flex min-h-0 flex-1 flex-col" aria-label="Query preview">
            <div class="border-b px-3 py-3">
              <div
                class="text-[11px] font-semibold uppercase tracking-[0.08em] text-muted-foreground"
              >
                Preview
              </div>
            </div>

            <div class="flex flex-1 items-center justify-center p-5">
              <div class="max-w-40 text-center">
                <div class="text-xs font-medium text-foreground/80">No preview yet</div>
                <div class="mt-1 text-[11px] leading-relaxed text-muted-foreground">
                  Results matching the query will appear here.
                </div>
              </div>
            </div>
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
              <div
                class="text-[11px] font-semibold uppercase tracking-[0.08em] text-muted-foreground"
              >
                Widgets
              </div>
            </div>

            <div class="flex flex-1 items-center justify-center p-5">
              <div class="max-w-44 text-center">
                <div class="text-xs font-medium text-foreground/80">No widgets yet</div>
                <div class="mt-1 text-[11px] leading-relaxed text-muted-foreground">
                  Dashboard widgets will appear here.
                </div>
              </div>
            </div>
          </section>

          <section class="flex min-h-0 flex-1 flex-col" aria-label="Properties">
            <div class="border-b px-3 py-3">
              <div
                class="text-[11px] font-semibold uppercase tracking-[0.08em] text-muted-foreground"
              >
                Properties
              </div>
            </div>

            <div class="flex flex-1 items-center justify-center p-5">
              <div class="max-w-44 text-center">
                <div class="text-xs font-medium text-foreground/80">Nothing selected</div>
                <div class="mt-1 text-[11px] leading-relaxed text-muted-foreground">
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
  private readonly projectQueries = inject(ProjectQueries);
  private readonly contextQueries = inject(ContextQueries);
  private readonly store = inject(Store);

  readonly queryScopeOptions: readonly SearchSelectOption[] = [
    { value: 'workspace', label: 'Workspace' },
    { value: 'project', label: 'Project' },
    { value: 'context', label: 'Context' },
    { value: 'daily', label: 'Daily' },
  ];
  readonly queryScope = signal<QueryScope>('workspace');
  readonly selectedEntityId = signal('');
  readonly days = signal(30);
  readonly selectedWorkspaceId = this.store.selectSignal(WorkspaceState.selectedWorkspaceId);

  readonly projectsQuery = injectQuery(() =>
    this.projectQueries.all(this.selectedWorkspaceId() ?? ''),
  );
  readonly contextsQuery = injectQuery(() => this.contextQueries.list(this.selectedWorkspaceId()));

  private readonly resetEntityOnWorkspaceChange = effect(() => {
    this.selectedWorkspaceId();
    this.selectedEntityId.set('');
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

  readonly entitySelectAriaLabel = computed(() => `Select ${this.queryScope()}`);
  readonly entitySelectPlaceholder = computed(() => `Select ${this.queryScope()}…`);
  readonly entitySearchPlaceholder = computed(() => `Search ${this.queryScope()}s…`);
  readonly entityEmptyText = computed(() =>
    this.entityOptionsLoading()
      ? `Loading ${this.queryScope()}s…`
      : `No matching ${this.queryScope()}s`,
  );

  setQueryScope(scope: string): void {
    this.queryScope.set(scope as QueryScope);
    this.selectedEntityId.set('');
  }

  selectEntity(entityId: string): void {
    this.selectedEntityId.set(entityId);
  }

  setDays(event: Event): void {
    const value = (event.target as HTMLInputElement).valueAsNumber;
    this.days.set(Number.isFinite(value) && value > 0 ? Math.floor(value) : 30);
  }
}

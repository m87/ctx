import { Component, computed, inject, signal } from '@angular/core';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { NgIcon, provideIcons } from '@ng-icons/core';
import {
  lucideArchive,
  lucideArchiveRestore,
  lucideChevronRight,
  lucideFolder,
  lucidePause,
  lucidePencil,
  lucidePlay,
  lucideTrash2,
  lucideX,
} from '@ng-icons/lucide';
import { HlmButtonImports } from '@spartan-ng/helm/button';
import { HlmCardImports } from '@spartan-ng/helm/card';
import { map } from 'rxjs';
import { toSignal } from '@angular/core/rxjs-interop';
import { ContextQueries } from '../../api/context/context.queries';
import { injectMutation, injectQuery } from '@tanstack/angular-query-experimental';
import { ContextMutations } from '../../api/context/context.mutations';
import { colorHash, durationAsH, durationAsM } from '../utils';
import { Store } from '@ngxs/store';
import { SelectProject, WorkspaceState } from '../sidebar/workspace.state';
import { NameComponent, NameSaveValue } from '../shared/name.component';
import { QueryErrorStateComponent } from '../shared/query-error-state.component';
import { ContextIntervalListComponent } from './context-interval-list.component';
import { TimeZoneService } from '../shared/time-zone.service';
import { ProjectQueries } from '../../api/project/project.queries';
import { HlmSkeletonImports } from '@spartan-ng/helm/skeleton';
import { InsightsEmptyStateComponent } from '../shared/insights-empty-state.component';
import { SearchSelectComponent, SearchSelectOption } from '../shared/search-select.component';

const WORKSPACE_ROOT_VALUE = '__workspace_root__';
type DetailView = 'overview' | 'insights';

@Component({
  imports: [
    NameComponent,
    ContextIntervalListComponent,
    NgIcon,
    HlmButtonImports,
    HlmCardImports,
    QueryErrorStateComponent,
    RouterLink,
    InsightsEmptyStateComponent,
    HlmSkeletonImports,
    SearchSelectComponent,
  ],
  providers: [
    provideIcons({
      lucideArchive,
      lucideArchiveRestore,
      lucideChevronRight,
      lucideFolder,
      lucidePause,
      lucidePencil,
      lucidePlay,
      lucideTrash2,
      lucideX,
    }),
  ],
  selector: 'ctx-context',
  template: `
    <div
      class="w-full h-full overflow-y-auto overflow-x-hidden md:overflow-hidden flex flex-col items-start justify-start p-4 md:p-6 gap-5 relative"
    >
      @if (showContextError()) {
        <ctx-query-error-state
          class="flex-1 min-h-0"
          [error]="contextQuery.error()"
          [paused]="contextQuery.isPaused()"
          resourceName="context"
          [retrying]="contextQuery.isFetching()"
          (retry)="retryContext()"
        ></ctx-query-error-state>
      } @else if (contextQuery.isLoading()) {
        <div class="w-full flex-1 min-h-0" role="status" aria-label="Loading context">
          <span class="sr-only">Loading context</span>
          <div class="flex flex-col md:flex-row justify-between gap-4 mb-5">
            <div class="flex-1">
              <hlm-skeleton class="h-2.5 w-16 mb-2"></hlm-skeleton>
              <hlm-skeleton class="h-7 w-52 mb-2"></hlm-skeleton>
              <hlm-skeleton class="h-3 w-80 max-w-full"></hlm-skeleton>
            </div>
            <div class="flex gap-2 md:pt-5">
              @for (item of actionSkeletonItems; track item) {
                <hlm-skeleton class="h-9 w-20"></hlm-skeleton>
              }
            </div>
          </div>
          <hlm-skeleton class="h-14 w-full mb-5"></hlm-skeleton>
          <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-5">
            @for (item of statsSkeletonItems; track item) {
              <div class="rounded-lg border bg-card p-3">
                <hlm-skeleton class="h-2.5 w-20 mb-2"></hlm-skeleton>
                <hlm-skeleton class="h-5 w-14"></hlm-skeleton>
              </div>
            }
          </div>
          <hlm-skeleton class="h-2.5 w-16 mb-4"></hlm-skeleton>
          <hlm-skeleton class="h-24 w-full mb-3"></hlm-skeleton>
          <hlm-skeleton class="h-16 w-full"></hlm-skeleton>
        </div>
      } @else if (context(); as currentContext) {
        <div class="w-full flex flex-col md:flex-row justify-between items-start gap-4">
          <ctx-name
            class="w-full min-w-0"
            label="Context"
            accentColor="#d97706"
            [name]="currentContext.name"
            [description]="currentContext.description ?? ''"
            [tags]="currentContext.tags ?? []"
            [showTags]="true"
            [readonly]="currentContext.archived ?? false"
            namePlaceholder="Context name"
            descriptionPlaceholder="What this context is for"
            tagsPlaceholder="Comma separated"
            [savePending]="updateContextMutation.isPending()"
            (save)="saveContextName($event)"
          ></ctx-name>

          <div class="flex items-center gap-2 w-full md:w-auto flex-nowrap md:pt-5">
            @if (currentContext.archived) {
              <span
                class="h-9 inline-flex items-center rounded-md border px-3 text-xs text-muted-foreground"
              >
                Archived
              </span>
            }
            @if (currentContext.archived) {
              <button
                hlmBtn
                variant="outline"
                class="h-9 px-3 text-xs bg-blue-200/70 text-blue-600"
                [disabled]="restoreContextMutation.isPending()"
                (click)="restoreContext()"
              >
                <ng-icon name="lucideArchiveRestore"></ng-icon>
                <span>Restore</span>
              </button>
            } @else {
              <button
                hlmBtn
                variant="outline"
                class="h-9 px-3 text-xs"
                [disabled]="archiveContextMutation.isPending()"
                (click)="archiveContext()"
              >
                <ng-icon name="lucideArchive"></ng-icon>
                <span>Archive</span>
              </button>
            }
            <button
              hlmBtn
              variant="outline"
              class="size-9 p-0 text-xs bg-red-100/70 text-red-700"
              [disabled]="deleteContextMutation.isPending()"
              (click)="deleteContext()"
            >
              <ng-icon name="lucideTrash2"></ng-icon>
            </button>
            @if (isActiveContext()) {
              <button
                hlmBtn
                variant="outline"
                class="h-9 px-3 text-xs bg-amber-100/70 text-amber-700"
                [disabled]="contextOperationPending()"
                [attr.aria-busy]="contextOperationPending()"
                (click)="pauseContext()"
              >
                @if (freeContextMutation.isPending()) {
                  <span
                    class="size-3.5 shrink-0 rounded-full border-2 border-current border-t-transparent animate-spin"
                    aria-hidden="true"
                  ></span>
                  <span class="font-semibold">Pausing...</span>
                } @else {
                  <ng-icon name="lucidePause"></ng-icon>
                  <span class="font-semibold">Pause</span>
                }
              </button>
            } @else {
              <button
                hlmBtn
                variant="outline"
                class="h-9 px-3 text-xs bg-blue-200/70 text-blue-600"
                [disabled]="currentContext.archived || contextOperationPending()"
                [attr.aria-busy]="contextOperationPending()"
                (click)="startContext()"
              >
                @if (switchContextMutation.isPending()) {
                  <span
                    class="size-3.5 shrink-0 rounded-full border-2 border-current border-t-transparent animate-spin"
                    aria-hidden="true"
                  ></span>
                  <span class="font-semibold text-blue-600">Starting...</span>
                } @else {
                  <ng-icon name="lucidePlay"></ng-icon>
                  <span class="font-semibold text-blue-600">Start</span>
                }
              </button>
            }
          </div>
        </div>

        <div
          class="inline-flex self-start rounded-lg bg-muted p-1 shrink-0"
          role="tablist"
          aria-label="Context view"
        >
          <button
            type="button"
            id="context-overview-tab"
            class="rounded-md px-3 py-1.5 text-xs font-medium transition-colors"
            [class.bg-background]="detailView() === 'overview'"
            [class.shadow-sm]="detailView() === 'overview'"
            [class.text-foreground]="detailView() === 'overview'"
            [class.text-muted-foreground]="detailView() !== 'overview'"
            role="tab"
            aria-controls="context-view-panel"
            [attr.aria-selected]="detailView() === 'overview'"
            (click)="detailView.set('overview')"
          >
            Overview
          </button>
          <button
            type="button"
            id="context-insights-tab"
            class="rounded-md px-3 py-1.5 text-xs font-medium transition-colors"
            [class.bg-background]="detailView() === 'insights'"
            [class.shadow-sm]="detailView() === 'insights'"
            [class.text-foreground]="detailView() === 'insights'"
            [class.text-muted-foreground]="detailView() !== 'insights'"
            role="tab"
            aria-controls="context-view-panel"
            [attr.aria-selected]="detailView() === 'insights'"
            (click)="detailView.set('insights')"
          >
            Insights
          </button>
        </div>

        <div
          id="context-view-panel"
          class="w-full flex-1 min-h-0 flex flex-col gap-5"
          role="tabpanel"
          [attr.aria-labelledby]="
            detailView() === 'overview' ? 'context-overview-tab' : 'context-insights-tab'
          "
        >
          @if (detailView() === 'overview') {
            @if (currentContext.project; as project) {
              <div
                class="w-full flex items-center gap-3 rounded-lg border bg-card p-3 cursor-pointer hover:bg-muted/30 transition-colors"
                [routerLink]="['/project', project.id]"
                role="link"
                tabindex="0"
                (click)="selectProject(project.id)"
              >
                <ctx-name
                  class="min-w-0 flex-1"
                  label="Project"
                  [name]="project.name"
                  [showDescription]="false"
                  [readonly]="true"
                  [compact]="true"
                  [accentColor]="projectColor(project.id)"
                ></ctx-name>
                <button
                  type="button"
                  class="size-8 shrink-0 rounded-md border text-muted-foreground hover:text-foreground hover:bg-muted/60 flex items-center justify-center"
                  aria-label="Change project"
                  title="Change project"
                  (click)="startProjectAssignmentEdit($event)"
                >
                  <ng-icon name="lucidePencil"></ng-icon>
                </button>
                <ng-icon
                  name="lucideChevronRight"
                  class="text-sm shrink-0 text-muted-foreground/70"
                ></ng-icon>
              </div>
            } @else {
              <button
                type="button"
                class="h-9 px-3 rounded-md border border-dashed text-xs text-muted-foreground hover:text-foreground hover:bg-muted/40 inline-flex items-center gap-2"
                (click)="startProjectAssignmentEdit($event)"
              >
                <ng-icon name="lucideFolder"></ng-icon>
                <span>Assign to project</span>
              </button>
            }

            @if (showContextStatsError()) {
              <ctx-query-error-state
                class="w-full min-h-40"
                [error]="contextStatsQuery.error()"
                [paused]="contextStatsQuery.isPaused()"
                resourceName="context statistics"
                [retrying]="contextStatsQuery.isFetching()"
                (retry)="retryContextStats()"
              ></ctx-query-error-state>
            } @else if (contextStatsQuery.isLoading()) {
              <div class="grid grid-cols-2 md:grid-cols-4 gap-4 w-full" role="status">
                <span class="sr-only">Loading context statistics</span>
                @for (item of statsSkeletonItems; track item) {
                  <div class="rounded-lg border bg-card p-3">
                    <hlm-skeleton class="h-2.5 w-20 mb-2"></hlm-skeleton>
                    <hlm-skeleton class="h-5 w-14"></hlm-skeleton>
                  </div>
                }
              </div>
            } @else {
              <div class="flex w-full">
                <div class="w-full flex items-center justify-center gap-4">
                  <div hlmCard class="w-full p-3 rounded-lg border">
                    <h3
                      class="text-[11px] uppercase tracking-[0.08em] text-muted-foreground font-semibold"
                      hlmCardTitle
                    >
                      Total time
                    </h3>
                    <div class="text-lg font-semibold">
                      {{ parseDuration(contextStats()?.totalDuration) }}
                    </div>
                  </div>
                  <div hlmCard class="w-full p-3 rounded-lg border">
                    <h3
                      class="text-[11px] uppercase tracking-[0.08em] text-muted-foreground font-semibold"
                      hlmCardTitle
                    >
                      Today
                    </h3>
                    <div class="text-lg font-semibold">
                      {{ parseDuration(contextStats()?.duration) }}
                    </div>
                  </div>
                  <div hlmCard class="w-full p-3 rounded-lg border">
                    <h3
                      class="text-[11px] uppercase tracking-[0.08em] text-muted-foreground font-semibold"
                      hlmCardTitle
                    >
                      Sessions
                    </h3>
                    <div class="text-lg font-semibold">
                      {{ contextStats()?.totalSessions }}
                    </div>
                  </div>
                  <div hlmCard class="w-full p-3 rounded-lg border">
                    <h3
                      class="text-[11px] uppercase tracking-[0.08em] text-muted-foreground font-semibold"
                      hlmCardTitle
                    >
                      Today sessions
                    </h3>
                    <div class="text-lg font-semibold">
                      {{ contextStats()?.sessions }}
                    </div>
                  </div>
                </div>
              </div>
            }
            <ctx-context-interval-list
              [contextId]="contextId()"
              [activeWorkspaceId]="activeWorkspaceId()"
              [contexts]="contexts()"
              [readonly]="currentContext.archived ?? false"
            ></ctx-context-interval-list>
          } @else {
            <ctx-insights-empty-state />
          }
        </div>

        @if (editingProjectAssignment()) {
          <div
            class="fixed inset-0 z-50 flex items-end sm:items-center sm:justify-center"
            (click)="closeProjectAssignmentDialog()"
          >
            <div class="absolute inset-0 bg-background/65 backdrop-blur-sm"></div>
            <section
              class="relative max-h-[92dvh] w-full overflow-hidden rounded-t-3xl border border-border/80 bg-popover text-popover-foreground shadow-2xl sm:w-[min(92vw,30rem)] sm:rounded-2xl"
              role="dialog"
              aria-modal="true"
              aria-labelledby="project-assignment-dialog-title"
              (click)="$event.stopPropagation()"
            >
              <header class="flex items-start gap-3 border-b border-border/70 px-5 py-4 sm:px-6">
                <div
                  class="flex size-10 shrink-0 items-center justify-center rounded-xl bg-muted text-foreground"
                >
                  <ng-icon name="lucideFolder" class="text-lg"></ng-icon>
                </div>
                <div class="min-w-0 flex-1 pt-0.5">
                  <h2
                    id="project-assignment-dialog-title"
                    class="text-base font-semibold tracking-tight"
                  >
                    Assign to project
                  </h2>
                  <p class="mt-0.5 truncate text-xs leading-5 text-muted-foreground">
                    Choose a project for {{ currentContext.name }}.
                  </p>
                </div>
                <button
                  type="button"
                  class="flex size-9 shrink-0 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:bg-muted hover:text-foreground focus-visible:ring-2 focus-visible:ring-ring focus-visible:outline-none"
                  aria-label="Close project assignment dialog"
                  [disabled]="updateContextMutation.isPending()"
                  (click)="closeProjectAssignmentDialog()"
                >
                  <ng-icon name="lucideX" class="text-base"></ng-icon>
                </button>
              </header>

              <div class="px-5 py-5 sm:px-6">
                <label
                  for="context-project-search"
                  class="mb-2 block text-xs font-medium text-foreground"
                >
                  Project
                </label>
                <ctx-search-select
                  inputId="context-project-search"
                  ariaLabel="Project"
                  searchPlaceholder="Search projects…"
                  emptyText="No matching projects"
                  [options]="projectSelectOptions()"
                  [value]="projectAssignmentTargetId()"
                  [disabled]="updateContextMutation.isPending()"
                  (selectionChange)="projectAssignmentTargetId.set($event)"
                ></ctx-search-select>
              </div>

              <footer
                class="flex flex-col-reverse gap-2 border-t border-border/70 bg-muted/20 px-5 py-4 sm:flex-row sm:justify-end sm:px-6"
              >
                <button
                  hlmBtn
                  variant="ghost"
                  class="h-10 px-4 text-xs"
                  [disabled]="updateContextMutation.isPending()"
                  (click)="closeProjectAssignmentDialog()"
                >
                  Cancel
                </button>
                <button
                  hlmBtn
                  variant="outline"
                  class="h-10 px-5 text-xs shadow-sm"
                  [disabled]="!projectAssignmentChanged() || updateContextMutation.isPending()"
                  (click)="confirmProjectAssignment()"
                >
                  <ng-icon name="lucideFolder"></ng-icon>
                  {{ updateContextMutation.isPending() ? 'Saving…' : 'Save assignment' }}
                </button>
              </footer>
            </section>
          </div>
        }
      }
    </div>
  `,
  styles: `
    :host {
      display: block;
      width: 100%;
      max-width: 1000px;
      height: 100%;
      min-height: 0;
    }
  `,
})
export class ContextComponent {
  readonly actionSkeletonItems = [0, 1, 2];
  readonly statsSkeletonItems = [0, 1, 2, 3];
  readonly detailView = signal<DetailView>('overview');
  private contextQueries = inject(ContextQueries);
  private contextMutations = inject(ContextMutations);
  private router = inject(Router);
  private store = inject(Store);
  private route = inject(ActivatedRoute);
  private timeZone = inject(TimeZoneService);
  private projectQueries = inject(ProjectQueries);
  readonly activeWorkspaceId = this.store.selectSignal(WorkspaceState.selectedWorkspaceId);
  readonly today = computed(() => this.timeZone.today());
  readonly contextId = toSignal(this.route.paramMap.pipe(map((pm) => pm.get('id') ?? '')), {
    initialValue: '',
  });

  switchContextMutation = injectMutation(() => this.contextMutations.switch());
  freeContextMutation = injectMutation(() => this.contextMutations.free());
  updateContextMutation = injectMutation(() => this.contextMutations.update());
  deleteContextMutation = injectMutation(() => this.contextMutations.delete());
  archiveContextMutation = injectMutation(() => this.contextMutations.archive());
  restoreContextMutation = injectMutation(() => this.contextMutations.restore());
  contextQuery = injectQuery(() => this.contextQueries.get(this.contextId()));
  activeContextQuery = injectQuery(() => this.contextQueries.active());
  contextsQuery = injectQuery(() => this.contextQueries.list(this.activeWorkspaceId()));
  context = computed(() => this.contextQuery.data() ?? null);
  readonly contextWorkspaceId = computed(
    () => this.context()?.workspaceId ?? this.activeWorkspaceId() ?? '',
  );
  projectsQuery = injectQuery(() => this.projectQueries.all(this.contextWorkspaceId()));
  readonly availableProjects = computed(() => this.projectsQuery.data() ?? []);
  readonly projectSelectOptions = computed<SearchSelectOption[]>(() => [
    {
      value: WORKSPACE_ROOT_VALUE,
      label: 'Workspace root',
      color: 'var(--muted-foreground)',
      description: 'No project',
    },
    ...this.availableProjects().map((project) => ({
      value: project.id,
      label: project.name,
      color: colorHash(project.id),
    })),
  ]);
  readonly editingProjectAssignment = signal(false);
  readonly workspaceRootValue = WORKSPACE_ROOT_VALUE;
  readonly projectAssignmentTargetId = signal(WORKSPACE_ROOT_VALUE);
  readonly projectAssignmentChanged = computed(
    () =>
      this.projectAssignmentTargetId() !== (this.context()?.project?.id ?? WORKSPACE_ROOT_VALUE),
  );
  isActiveContext = computed(() => this.activeContextQuery.data()?.id === this.contextId());
  readonly contextOperationPending = computed(
    () => this.switchContextMutation.isPending() || this.freeContextMutation.isPending(),
  );
  contextStatsQuery = injectQuery(() =>
    this.contextQueries.stats(this.contextId(), this.today(), this.timeZone.effectiveTimeZone()),
  );
  contextStats = computed(() => this.contextStatsQuery.data());
  readonly showContextError = computed(
    () =>
      this.contextQuery.data() === undefined &&
      (this.contextQuery.isError() || this.contextQuery.isPaused()),
  );
  readonly showContextStatsError = computed(
    () =>
      this.contextStatsQuery.data() === undefined &&
      (this.contextStatsQuery.isError() || this.contextStatsQuery.isPaused()),
  );
  contexts = computed(() => this.contextsQuery.data() ?? []);

  startContext() {
    const context = this.context();
    if (!context || context.archived) {
      return;
    }
    this.switchContextMutation.mutate(context);
  }

  pauseContext(): void {
    if (!this.isActiveContext()) {
      return;
    }
    this.freeContextMutation.mutate();
  }

  deleteContext() {
    const context = this.context();

    if (!context?.id) {
      return;
    }

    if (!window.confirm(`Delete context "${context.name}"?`)) {
      return;
    }

    this.deleteContextMutation.mutate(context.id, {
      onSuccess: () => {
        this.router.navigate(['/day', this.today()]);
      },
    });
  }

  saveContextName(value: NameSaveValue): void {
    const context = this.context();
    if (!context || context.archived) {
      return;
    }

    this.updateContextMutation.mutate({
      id: context.id,
      context: {
        ...context,
        name: value.name,
        description: value.description,
        tags: value.tags ?? [],
      },
    });
  }

  archiveContext(): void {
    const context = this.context();
    if (!context || context.archived) {
      return;
    }
    if (!window.confirm(`Archive context "${context.name}"?`)) {
      return;
    }

    this.archiveContextMutation.mutate(context.id);
  }

  restoreContext(): void {
    const context = this.context();
    if (!context || !context.archived) {
      return;
    }

    this.restoreContextMutation.mutate(context.id);
  }

  startProjectAssignmentEdit(event: MouseEvent): void {
    event.preventDefault();
    event.stopPropagation();
    this.projectAssignmentTargetId.set(this.context()?.project?.id ?? WORKSPACE_ROOT_VALUE);
    this.editingProjectAssignment.set(true);
  }

  closeProjectAssignmentDialog(): void {
    this.editingProjectAssignment.set(false);
    this.projectAssignmentTargetId.set(WORKSPACE_ROOT_VALUE);
  }

  confirmProjectAssignment(): void {
    if (!this.projectAssignmentChanged() || this.updateContextMutation.isPending()) {
      return;
    }

    this.assignContextToProject(this.projectAssignmentTargetId());
  }

  assignContextToProject(selectedValue: string): void {
    const context = this.context();
    if (!context) {
      return;
    }

    const projectId = selectedValue === WORKSPACE_ROOT_VALUE ? '' : selectedValue;
    const project = this.availableProjects().find((candidate) => candidate.id === projectId);
    this.updateContextMutation.mutate(
      {
        id: context.id,
        context: {
          ...context,
          project: project ? { id: project.id, name: project.name } : undefined,
        },
      },
      {
        onSuccess: () => this.closeProjectAssignmentDialog(),
      },
    );
  }

  selectProject(projectId: string): void {
    this.store.dispatch(new SelectProject(projectId));
  }

  projectColor(projectId: string): string {
    return colorHash(projectId);
  }

  parseDuration(duration: number | undefined): string {
    if (duration === undefined) {
      return '0h 0m';
    }
    return `${durationAsH(duration)}h ${durationAsM(duration)}m`;
  }

  retryContext(): void {
    void this.contextQuery.refetch();
  }

  retryContextStats(): void {
    void this.contextStatsQuery.refetch();
  }
}

import { Component, computed, inject, signal } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { map } from 'rxjs';
import { Store } from '@ngxs/store';
import { injectMutation, injectQuery } from '@tanstack/angular-query-experimental';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideTrash2 } from '@ng-icons/lucide';
import { WorkspaceQueries } from '../../api/workspace/workspace.queries';
import { WorkspaceMutations } from '../../api/workspace/workspace.mutations';
import { WorkspaceState } from '../sidebar/workspace.state';
import { WorkspaceStats } from '../../api/workspace/workspace.service';
import { ContextListComponent } from '../context/context-list.component';
import { ContextListGroup } from '../context/context-list-group.component';
import { ContextListItem } from '../context/context-list-item.component';
import { DistributionComponent, DistributionItem } from '../shared/distribution.component';
import { NameComponent, NameSaveValue } from '../shared/name.component';
import { QueryErrorStateComponent } from '../shared/query-error-state.component';
import { colorHash, durationAsHM } from '../utils';
import { HlmSkeletonImports } from '@spartan-ng/helm/skeleton';
import {
  ProjectTimeListComponent,
  ProjectTimeListItem,
} from '../shared/project-time-list.component';
import { summarizeContextsByProject, UNASSIGNED_PROJECT_ID } from '../shared/project-time-summary';
import { InsightsEmptyStateComponent } from '../shared/insights-empty-state.component';

type SummaryView = 'contexts' | 'projects' | 'insights';

const GROUPED_CONTEXT_ID = '__contexts_below_1_percent__';
const GROUPED_CONTEXT_THRESHOLD = 1;

const EMPTY_WORKSPACE_STATS: WorkspaceStats = {
  workspaceId: '',
  contexts: [],
  contextStats: [],
  totalDuration: 0,
  totalSessions: 0,
};

@Component({
  selector: 'ctx-workspace',
  imports: [
    ContextListComponent,
    DistributionComponent,
    NameComponent,
    NgIcon,
    QueryErrorStateComponent,
    ProjectTimeListComponent,
    InsightsEmptyStateComponent,
    HlmSkeletonImports,
  ],
  providers: [provideIcons({ lucideTrash2 })],
  template: `
    <div class="w-full h-full overflow-hidden flex flex-col p-4 md:p-6">
      @if (showWorkspaceError()) {
        <ctx-query-error-state
          class="flex-1 min-h-0"
          [error]="workspaceError()"
          [paused]="workspaceErrorPaused()"
          [resourceName]="workspaceErrorResourceName()"
          [retrying]="workspaceErrorRetrying()"
          (retry)="retryWorkspaceData()"
        ></ctx-query-error-state>
      } @else if (isWorkspaceLoading()) {
        <div class="w-full flex-1 min-h-0" role="status" aria-label="Loading workspace">
          <span class="sr-only">Loading workspace</span>
          <hlm-skeleton class="h-2.5 w-20 mb-2"></hlm-skeleton>
          <hlm-skeleton class="h-7 w-56 mb-2"></hlm-skeleton>
          <hlm-skeleton class="h-3 w-80 max-w-full mb-12"></hlm-skeleton>
          <hlm-skeleton class="h-2.5 w-28 mb-2"></hlm-skeleton>
          <div class="grid grid-cols-2 md:grid-cols-4 gap-2.5 mb-6">
            @for (item of summarySkeletonItems; track item) {
              <div class="rounded-lg border bg-card px-3 py-2.5">
                <hlm-skeleton class="h-2.5 w-20 mb-2"></hlm-skeleton>
                <hlm-skeleton class="h-5 w-14"></hlm-skeleton>
              </div>
            }
          </div>
          <hlm-skeleton class="h-8 w-44 mb-4"></hlm-skeleton>
          <hlm-skeleton class="h-2.5 w-20 mb-2"></hlm-skeleton>
          <hlm-skeleton class="h-2 w-full mb-6"></hlm-skeleton>
          <hlm-skeleton class="h-2.5 w-16 mb-2"></hlm-skeleton>
          <div class="flex flex-col gap-2">
            @for (item of contextSkeletonItems; track item) {
              <div class="rounded-lg border bg-card p-3">
                <div class="flex items-center gap-2 mb-3">
                  <hlm-skeleton class="size-2"></hlm-skeleton>
                  <hlm-skeleton class="h-3.5 w-2/5"></hlm-skeleton>
                  <hlm-skeleton class="h-3 w-12 ml-auto"></hlm-skeleton>
                </div>
                <hlm-skeleton class="h-1.5 w-full"></hlm-skeleton>
              </div>
            }
          </div>
        </div>
      } @else {
        <div class="mb-5">
          @if (workspace()) {
            <div class="flex items-start justify-between gap-4">
              <ctx-name
                class="min-w-0 flex-1"
                label="Workspace"
                [name]="workspace()?.name ?? ''"
                [description]="workspace()?.description ?? ''"
                namePlaceholder="Workspace name"
                descriptionPlaceholder="What this workspace is for"
                [savePending]="updateWorkspaceMutation.isPending()"
                (save)="saveWorkspaceName($event)"
              ></ctx-name>

              <button
                type="button"
                class="size-9 rounded-md border border-destructive/30 text-destructive hover:bg-destructive/10 flex items-center justify-center shrink-0 mt-5"
                aria-label="Delete workspace"
                title="Delete"
                (click)="deleteWorkspace()"
              >
                <ng-icon name="lucideTrash2"></ng-icon>
              </button>
            </div>
          } @else {
            <div class="text-[11px] uppercase tracking-widest text-muted-foreground font-semibold">
              Workspace
            </div>
            <h1 class="text-2xl font-semibold tracking-tight mt-1">Default workspace</h1>
          }
        </div>

        @if (workspace()) {
          <div class="mt-6 flex-1 min-h-0 overflow-auto pr-1 pb-2">
            <div
              class="text-[11px] uppercase tracking-[0.08em] text-muted-foreground font-semibold mb-2"
            >
              Workspace summary
            </div>

            <div class="grid grid-cols-2 md:grid-cols-4 gap-2.5 mb-6">
              <div class="rounded-lg border bg-card px-3 py-2.5">
                <div class="text-[10px] uppercase tracking-[0.08em] text-muted-foreground">
                  Total tracked
                </div>
                <div class="text-base font-semibold mt-1">{{ totalTracked() }}</div>
              </div>
              <div class="rounded-lg border bg-card px-3 py-2.5">
                <div class="text-[10px] uppercase tracking-[0.08em] text-muted-foreground">
                  Contexts
                </div>
                <div class="text-base font-semibold mt-1">
                  {{ workspaceStats().contexts.length }}
                </div>
              </div>
              <div class="rounded-lg border bg-card px-3 py-2.5">
                <div class="text-[10px] uppercase tracking-[0.08em] text-muted-foreground">
                  Sessions
                </div>
                <div class="text-base font-semibold mt-1">
                  {{ workspaceStats().totalSessions }}
                </div>
              </div>
              <div class="rounded-lg border bg-card px-3 py-2.5">
                <div class="text-[10px] uppercase tracking-[0.08em] text-muted-foreground">
                  Top context
                </div>
                <div class="text-sm font-medium mt-1 truncate">{{ topContext() }}</div>
              </div>
            </div>

            <div
              class="inline-flex rounded-lg bg-muted p-1 mb-4"
              role="tablist"
              aria-label="Workspace summary view"
            >
              <button
                type="button"
                id="workspace-contexts-tab"
                class="rounded-md px-3 py-1.5 text-xs font-medium transition-colors"
                [class.bg-background]="summaryView() === 'contexts'"
                [class.shadow-sm]="summaryView() === 'contexts'"
                [class.text-foreground]="summaryView() === 'contexts'"
                [class.text-muted-foreground]="summaryView() !== 'contexts'"
                role="tab"
                aria-controls="workspace-summary-panel"
                [attr.aria-selected]="summaryView() === 'contexts'"
                (click)="selectSummaryView('contexts')"
              >
                Contexts
              </button>
              <button
                type="button"
                id="workspace-projects-tab"
                class="rounded-md px-3 py-1.5 text-xs font-medium transition-colors"
                [class.bg-background]="summaryView() === 'projects'"
                [class.shadow-sm]="summaryView() === 'projects'"
                [class.text-foreground]="summaryView() === 'projects'"
                [class.text-muted-foreground]="summaryView() !== 'projects'"
                role="tab"
                aria-controls="workspace-summary-panel"
                [attr.aria-selected]="summaryView() === 'projects'"
                (click)="selectSummaryView('projects')"
              >
                Projects
              </button>
              <button
                type="button"
                id="workspace-insights-tab"
                class="rounded-md px-3 py-1.5 text-xs font-medium transition-colors"
                [class.bg-background]="summaryView() === 'insights'"
                [class.shadow-sm]="summaryView() === 'insights'"
                [class.text-foreground]="summaryView() === 'insights'"
                [class.text-muted-foreground]="summaryView() !== 'insights'"
                role="tab"
                aria-controls="workspace-summary-panel"
                [attr.aria-selected]="summaryView() === 'insights'"
                (click)="selectSummaryView('insights')"
              >
                Insights
              </button>
            </div>

            @if (summaryView() !== 'insights') {
              <ctx-distribution
                class="block mb-6"
                [label]="distributionLabel()"
                [items]="activeDistribution()"
                emptyMessage="No tracked time in this workspace."
              ></ctx-distribution>
            }

            <div
              id="workspace-summary-panel"
              role="tabpanel"
              [attr.aria-labelledby]="summaryViewTabId()"
            >
              @if (summaryView() === 'contexts') {
                <ctx-context-list
                  [items]="largeSummaryContexts()"
                  [group]="groupedSummaryContext()"
                  emptyMessage="No contexts tracked in this workspace."
                ></ctx-context-list>
              } @else if (summaryView() === 'projects') {
                <ctx-project-time-list
                  [items]="projectSummaries()"
                  emptyMessage="No projects tracked in this workspace."
                ></ctx-project-time-list>
              } @else {
                <ctx-insights-empty-state />
              }
            </div>
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
export class WorkspaceComponent {
  readonly summarySkeletonItems = [0, 1, 2, 3];
  readonly contextSkeletonItems = [0, 1, 2];
  readonly summaryView = signal<SummaryView>('contexts');
  private readonly route = inject(ActivatedRoute);
  private readonly store = inject(Store);
  private readonly workspaceQueries = inject(WorkspaceQueries);
  private readonly workspaceMutations = inject(WorkspaceMutations);

  private readonly routeWorkspaceId = toSignal(
    this.route.paramMap.pipe(map((params) => params.get('id'))),
    { initialValue: null },
  );
  private readonly selectedWorkspaceId = this.store.selectSignal(
    WorkspaceState.selectedWorkspaceId,
  );

  listWorkspacesQuery = injectQuery(() => this.workspaceQueries.list());
  updateWorkspaceMutation = injectMutation(() => this.workspaceMutations.update());
  deleteWorkspaceMutation = injectMutation(() => this.workspaceMutations.delete());

  readonly activeWorkspaceId = computed(
    () => this.routeWorkspaceId() ?? this.selectedWorkspaceId(),
  );
  workspaceStatsQuery = injectQuery(() =>
    this.workspaceQueries.stats(this.activeWorkspaceId() ?? ''),
  );
  readonly showWorkspaceListError = computed(
    () =>
      this.listWorkspacesQuery.data() === undefined &&
      (this.listWorkspacesQuery.isError() || this.listWorkspacesQuery.isPaused()),
  );
  readonly showWorkspaceStatsError = computed(
    () =>
      this.activeWorkspaceId() !== null &&
      this.workspaceStatsQuery.data() === undefined &&
      (this.workspaceStatsQuery.isError() || this.workspaceStatsQuery.isPaused()),
  );
  readonly showWorkspaceError = computed(
    () => this.showWorkspaceListError() || this.showWorkspaceStatsError(),
  );
  readonly isWorkspaceLoading = computed(() => {
    if (this.listWorkspacesQuery.isLoading()) {
      return true;
    }

    const workspaceId = this.activeWorkspaceId();
    return !!workspaceId && this.workspaceStatsQuery.isLoading();
  });
  readonly workspaceError = computed(() =>
    this.showWorkspaceListError()
      ? this.listWorkspacesQuery.error()
      : this.workspaceStatsQuery.error(),
  );
  readonly workspaceErrorPaused = computed(() =>
    this.showWorkspaceListError()
      ? this.listWorkspacesQuery.isPaused()
      : this.workspaceStatsQuery.isPaused(),
  );
  readonly workspaceErrorResourceName = computed(() =>
    this.showWorkspaceListError() ? 'workspaces' : 'workspace summary',
  );
  readonly workspaceErrorRetrying = computed(() =>
    this.showWorkspaceListError()
      ? this.listWorkspacesQuery.isFetching()
      : this.workspaceStatsQuery.isFetching(),
  );
  readonly workspace = computed(() => {
    const id = this.activeWorkspaceId();
    return this.listWorkspacesQuery.data()?.find((workspace) => workspace.id === id) ?? null;
  });
  readonly workspaceStats = computed(
    () => this.workspaceStatsQuery.data() ?? EMPTY_WORKSPACE_STATS,
  );
  readonly allSummaryContexts = computed<ContextListItem[]>(() => {
    const contextsById = new Map(
      this.workspaceStats().contexts.map((context) => [context.id, context]),
    );

    const contexts = this.workspaceStats()
      .contextStats.filter((stats) => stats.duration > 0)
      .map((stats) => {
        const context = contextsById.get(stats.contextId);

        return {
          id: stats.contextId,
          name: context?.name ?? stats.contextId,
          duration: durationAsHM(stats.duration).trim() || '0m',
          durationValue: stats.duration,
          sessions: stats.intervalCount,
          percentage: stats.percentage,
          color: colorHash(stats.contextId),
          archived: context?.archived ?? false,
          project: context?.project,
        };
      });

    return contexts;
  });

  readonly smallSummaryContexts = computed(() =>
    this.allSummaryContexts().filter((context) => context.percentage < GROUPED_CONTEXT_THRESHOLD),
  );
  readonly largeSummaryContexts = computed(() =>
    this.allSummaryContexts().filter((context) => context.percentage >= GROUPED_CONTEXT_THRESHOLD),
  );
  readonly groupedSummaryContext = computed<ContextListGroup | null>(() => {
    const groupedContexts = this.smallSummaryContexts();
    if (groupedContexts.length === 0) {
      return null;
    }

    const durationValue = groupedContexts.reduce(
      (duration, context) => duration + (context.durationValue ?? 0),
      0,
    );

    return {
      id: GROUPED_CONTEXT_ID,
      name: 'Other contexts (<1% each)',
      duration: durationAsHM(durationValue).trim() || '0m',
      durationValue,
      sessions: groupedContexts.reduce(
        (sessions, context) => sessions + (context.sessions ?? 0),
        0,
      ),
      percentage: groupedContexts.reduce(
        (percentage, context) => percentage + context.percentage,
        0,
      ),
      color: '#94a3b8',
      groupedCount: groupedContexts.length,
      items: groupedContexts,
    };
  });
  readonly distributionContexts = computed<DistributionItem[]>(() => {
    const groupedContext = this.groupedSummaryContext();
    return groupedContext
      ? [...this.largeSummaryContexts(), groupedContext]
      : this.allSummaryContexts();
  });
  readonly projectSummaries = computed<ProjectTimeListItem[]>(() =>
    summarizeContextsByProject(this.workspaceStats()).map((summary) => ({
      ...summary,
      duration: durationAsHM(summary.duration).trim() || '0m',
      color: summary.id === UNASSIGNED_PROJECT_ID ? '#94a3b8' : colorHash(`project:${summary.id}`),
    })),
  );
  readonly projectDistribution = computed<DistributionItem[]>(() =>
    this.projectSummaries().map((project) => ({
      id: project.id,
      name: project.name,
      duration: project.duration,
      percentage: project.percentage,
      color: project.color,
    })),
  );
  readonly activeDistribution = computed(() =>
    this.summaryView() === 'contexts' ? this.distributionContexts() : this.projectDistribution(),
  );
  readonly distributionLabel = computed(() =>
    this.summaryView() === 'contexts' ? 'Context distribution' : 'Project distribution',
  );
  readonly summaryViewTabId = computed(() => `workspace-${this.summaryView()}-tab`);
  readonly totalTracked = computed(
    () => durationAsHM(this.workspaceStats().totalDuration).trim() || '0m',
  );
  readonly topContext = computed(() => this.allSummaryContexts()[0]?.name ?? '-');

  selectSummaryView(view: SummaryView): void {
    this.summaryView.set(view);
  }

  saveWorkspaceName(value: NameSaveValue): void {
    const workspace = this.workspace();
    if (!workspace) {
      return;
    }

    this.updateWorkspaceMutation.mutate({
      ...workspace,
      name: value.name,
      description: value.description,
    });
  }

  deleteWorkspace(): void {
    const workspace = this.workspace();
    if (!workspace) {
      return;
    }

    this.deleteWorkspaceMutation.mutate(workspace.id);
  }

  retryWorkspaceData(): void {
    if (this.showWorkspaceListError()) {
      void this.listWorkspacesQuery.refetch();
      return;
    }

    void this.workspaceStatsQuery.refetch();
  }
}

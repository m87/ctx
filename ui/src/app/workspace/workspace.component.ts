import { Component, computed, inject } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { map } from 'rxjs';
import { Store } from '@ngxs/store';
import { injectMutation, injectQuery } from '@tanstack/angular-query-experimental';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideTrash2 } from '@ng-icons/lucide';
import { WorkspaceQueries } from '../../api/workspace.quries';
import { WorkspaceMutations } from '../../api/workspace.mutations';
import { WorkspaceState } from '../sidebar/workspace.state';
import { WorkspaceStats } from '../../api/workspace.service';
import { ContextListComponent } from '../context/context-list.component';
import { ContextListGroup } from '../context/context-list-group.component';
import { ContextListItem } from '../context/context-list-item.component';
import { DistributionComponent, DistributionItem } from '../shared/distribution.component';
import { NameComponent, NameSaveValue } from '../shared/name.component';
import { colorHash, durationAsHM } from '../utils';

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
  imports: [ContextListComponent, DistributionComponent, NameComponent, NgIcon],
  providers: [provideIcons({ lucideTrash2 })],
  template: `
    <div class="w-full h-full overflow-hidden flex flex-col p-4 md:p-6">
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
              class="h-8 w-8 rounded-md border border-destructive/30 text-destructive hover:bg-destructive/10 flex items-center justify-center shrink-0 mt-5"
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
              <div class="text-base font-semibold mt-1">{{ workspaceStats().contexts.length }}</div>
            </div>
            <div class="rounded-lg border bg-card px-3 py-2.5">
              <div class="text-[10px] uppercase tracking-[0.08em] text-muted-foreground">
                Sessions
              </div>
              <div class="text-base font-semibold mt-1">{{ workspaceStats().totalSessions }}</div>
            </div>
            <div class="rounded-lg border bg-card px-3 py-2.5">
              <div class="text-[10px] uppercase tracking-[0.08em] text-muted-foreground">
                Top context
              </div>
              <div class="text-sm font-medium mt-1 truncate">{{ topContext() }}</div>
            </div>
          </div>

          <ctx-distribution
            class="block mb-6"
            [items]="distributionContexts()"
            emptyMessage="No tracked time in this workspace."
          ></ctx-distribution>

          @if (largeSummaryContexts().length > 0 || groupedSummaryContext()) {
            <div
              class="text-[11px] uppercase tracking-[0.08em] text-muted-foreground font-semibold mb-2"
            >
              Contexts
            </div>
            <ctx-context-list
              [items]="largeSummaryContexts()"
              [group]="groupedSummaryContext()"
            ></ctx-context-list>
          }
        </div>
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
  readonly totalTracked = computed(
    () => durationAsHM(this.workspaceStats().totalDuration).trim() || '0m',
  );
  readonly topContext = computed(() => this.allSummaryContexts()[0]?.name ?? '-');

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
}

import { Component, computed, inject, signal } from '@angular/core';
import { NgTemplateOutlet } from '@angular/common';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { map } from 'rxjs';
import { Store } from '@ngxs/store';
import { injectMutation, injectQuery } from '@tanstack/angular-query-experimental';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideCheck, lucidePencil, lucideTrash2, lucideX } from '@ng-icons/lucide';
import { WorkspaceQueries } from '../../api/workspace.quries';
import { WorkspaceMutations } from '../../api/workspace.mutations';
import { SelectWorkspace, WorkspaceState } from '../sidebar/workspace.state';
import { WorkspaceStats } from '../../api/workspace.service';
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

type SummaryContext = {
  id: string;
  name: string;
  duration: string;
  durationValue: number;
  sessions: number;
  percentage: number;
  color: string;
  grouped?: boolean;
  groupedChild?: boolean;
  groupedCount?: number;
};

@Component({
  selector: 'app-workspace',
  imports: [NgIcon, NgTemplateOutlet, RouterLink],
  providers: [provideIcons({ lucideCheck, lucidePencil, lucideTrash2, lucideX })],
  template: `
    <div class="w-full h-full overflow-hidden flex flex-col p-4 md:p-6">
      <div class="mb-5">
        <div class="text-[11px] uppercase tracking-widest text-muted-foreground font-semibold">
          Workspace
        </div>
        <div class="mt-1 flex items-start justify-between gap-4">
          <div class="min-w-0 flex-1">
            @if (workspace()) {
              @if (isEditing()) {
                <div class="grid max-w-2xl gap-3">
                  <label class="flex flex-col gap-1">
                    <span
                      class="text-[11px] uppercase tracking-[0.08em] text-muted-foreground font-semibold"
                      >Name</span
                    >
                    <input
                      type="text"
                      class="h-9 w-full rounded-md border bg-background px-2.5 text-sm outline-none focus:ring-1 focus:ring-ring"
                      placeholder="Workspace name"
                      [value]="editName()"
                      (input)="editName.set(getInputValue($event))"
                      (keydown.escape)="cancelEdit()"
                    />
                  </label>

                  <label class="flex flex-col gap-1">
                    <span
                      class="text-[11px] uppercase tracking-[0.08em] text-muted-foreground font-semibold"
                      >Description</span
                    >
                    <textarea
                      class="min-h-24 w-full rounded-md border bg-background px-2.5 py-2 text-sm outline-none focus:ring-1 focus:ring-ring"
                      placeholder="What this workspace is for"
                      [value]="editDescription()"
                      (input)="editDescription.set(getInputValue($event))"
                      (keydown.escape)="cancelEdit()"
                    ></textarea>
                  </label>
                </div>
              } @else {
                <h1 class="text-2xl font-semibold tracking-tight truncate">
                  {{ workspace()?.name }}
                </h1>
                @if (workspace()?.description) {
                  <p
                    class="mt-1 whitespace-pre-wrap text-sm text-muted-foreground text-ellipsis overflow-hidden"
                  >
                    {{ workspace()?.description }}
                  </p>
                } @else {
                  <p class="mt-1 text-sm text-muted-foreground">No description</p>
                }
              }
            } @else {
              <h1 class="text-2xl font-semibold tracking-tight">Default workspace</h1>
            }
          </div>

          @if (workspace()) {
            <div class="flex shrink-0 items-center gap-1 pt-0.5">
              @if (isEditing()) {
                <button
                  type="button"
                  class="h-8 w-8 rounded-md bg-primary text-primary-foreground hover:bg-primary/90 flex items-center justify-center"
                  aria-label="Save workspace"
                  title="Save"
                  [disabled]="updateWorkspaceMutation.isPending()"
                  (click)="saveEdit()"
                >
                  <ng-icon name="lucideCheck"></ng-icon>
                </button>
                <button
                  type="button"
                  class="h-8 w-8 rounded-md border hover:bg-muted/60 flex items-center justify-center"
                  aria-label="Cancel workspace edit"
                  title="Cancel"
                  (click)="cancelEdit()"
                >
                  <ng-icon name="lucideX"></ng-icon>
                </button>
              } @else {
                <button
                  type="button"
                  class="h-8 w-8 rounded-md border text-muted-foreground hover:text-foreground hover:bg-muted/60 flex items-center justify-center"
                  aria-label="Edit workspace"
                  title="Edit"
                  (click)="startEdit()"
                >
                  <ng-icon name="lucidePencil"></ng-icon>
                </button>
                <button
                  type="button"
                  class="h-8 w-8 rounded-md border border-destructive/30 text-destructive hover:bg-destructive/10 flex items-center justify-center"
                  aria-label="Delete workspace"
                  title="Delete"
                  (click)="deleteWorkspace()"
                >
                  <ng-icon name="lucideTrash2"></ng-icon>
                </button>
              }
            </div>
          }
        </div>
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

          <div class="mb-6">
            <div
              class="text-[11px] uppercase tracking-[0.08em] text-muted-foreground font-semibold mb-2"
            >
              Distribution
            </div>
            @if (distributionContexts().length > 0) {
              <div class="flex h-2 rounded-md overflow-hidden gap-px bg-muted/40">
                @for (context of distributionContexts(); track context.id) {
                  <div
                    [style.width.%]="context.percentage"
                    [style.background-color]="context.color"
                    [title]="context.name + ': ' + context.duration"
                  ></div>
                }
              </div>
            } @else {
              <div class="h-2 rounded-md bg-muted/40"></div>
              <p class="mt-2 text-xs text-muted-foreground">No tracked time in this workspace.</p>
            }
          </div>

          @if (summaryContexts().length > 0) {
            <div
              class="text-[11px] uppercase tracking-[0.08em] text-muted-foreground font-semibold mb-2"
            >
              Contexts
            </div>
            <div class="flex flex-col gap-2">
              @for (context of summaryContexts(); track context.id) {
                @if (context.grouped) {
                  <div
                    class="rounded-lg border bg-card p-3 hover:bg-muted/30 transition-colors cursor-pointer"
                    role="button"
                    tabindex="0"
                    [attr.aria-expanded]="showGroupedContexts()"
                    (click)="toggleGroupedContexts()"
                    (keydown.enter)="toggleGroupedContexts()"
                    (keydown.space)="$event.preventDefault(); toggleGroupedContexts()"
                  >
                    <ng-container
                      *ngTemplateOutlet="summaryContextItem; context: { $implicit: context }"
                    ></ng-container>
                    <div class="mt-2 text-[10px] text-muted-foreground">
                      {{
                        showGroupedContexts()
                          ? 'Hide smaller contexts'
                          : 'Show ' +
                            context.groupedCount +
                            ' smaller ' +
                            (context.groupedCount === 1 ? 'context' : 'contexts')
                      }}
                    </div>
                  </div>
                } @else {
                  <a
                    class="rounded-lg border bg-card p-3 hover:bg-muted/30 transition-colors"
                    [class.ml-4]="context.groupedChild"
                    [class.border-dashed]="context.groupedChild"
                    [routerLink]="['/context', context.id]"
                  >
                    <ng-container
                      *ngTemplateOutlet="summaryContextItem; context: { $implicit: context }"
                    ></ng-container>
                  </a>
                }
              }
            </div>
          }
        </div>
      }

      <ng-template #summaryContextItem let-context>
        <div class="flex items-center gap-2 mb-2">
          <span
            class="w-2 h-2 rounded-sm shrink-0"
            [style.background-color]="context.color"
          ></span>
          <span class="text-sm font-medium flex-1 truncate">{{ context.name }}</span>
          <span class="text-xs text-muted-foreground">{{ context.duration }}</span>
        </div>
        <div class="h-1.5 rounded bg-muted/40 overflow-hidden">
          <div
            class="h-full rounded"
            [style.width.%]="context.percentage"
            [style.background-color]="context.color"
          ></div>
        </div>
        <div class="mt-2 text-[10px] text-muted-foreground">
          {{ context.sessions }} {{ context.sessions === 1 ? 'session' : 'sessions' }} ·
          {{ context.percentage.toFixed(1) }}%
        </div>
      </ng-template>
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

  readonly isEditing = signal(false);
  readonly editName = signal('');
  readonly editDescription = signal('');
  readonly showGroupedContexts = signal(false);
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
  readonly allSummaryContexts = computed<SummaryContext[]>(() => {
    const contextsById = new Map(
      this.workspaceStats().contexts.map((context) => [context.id, context]),
    );

    const contexts = this.workspaceStats()
      .contextStats.filter((stats) => stats.duration > 0)
      .map((stats) => ({
        id: stats.contextId,
        name: contextsById.get(stats.contextId)?.name ?? stats.contextId,
        duration: durationAsHM(stats.duration).trim() || '0m',
        durationValue: stats.duration,
        sessions: stats.intervalCount,
        percentage: stats.percentage,
        color: colorHash(stats.contextId),
      }));

    return contexts;
  });

  readonly smallSummaryContexts = computed(() =>
    this.allSummaryContexts().filter(
      (context) => context.percentage < GROUPED_CONTEXT_THRESHOLD,
    ),
  );
  readonly largeSummaryContexts = computed(() =>
    this.allSummaryContexts().filter(
      (context) => context.percentage >= GROUPED_CONTEXT_THRESHOLD,
    ),
  );
  readonly groupedSummaryContext = computed<SummaryContext | null>(() => {
    const groupedContexts = this.smallSummaryContexts();
    if (groupedContexts.length === 0) {
      return null;
    }

    const durationValue = groupedContexts.reduce(
      (duration, context) => duration + context.durationValue,
      0,
    );

    return {
      id: GROUPED_CONTEXT_ID,
      name: 'Other contexts (<1% each)',
      duration: durationAsHM(durationValue).trim() || '0m',
      durationValue,
      sessions: groupedContexts.reduce((sessions, context) => sessions + context.sessions, 0),
      percentage: groupedContexts.reduce(
        (percentage, context) => percentage + context.percentage,
        0,
      ),
      color: '#94a3b8',
      grouped: true,
      groupedCount: groupedContexts.length,
    };
  });
  readonly distributionContexts = computed(() => {
    const groupedContext = this.groupedSummaryContext();
    return groupedContext
      ? [...this.largeSummaryContexts(), groupedContext]
      : this.allSummaryContexts();
  });
  readonly summaryContexts = computed(() => {
    const groupedContext = this.groupedSummaryContext();
    if (!groupedContext) {
      return this.allSummaryContexts();
    }

    if (!this.showGroupedContexts()) {
      return [...this.largeSummaryContexts(), groupedContext];
    }

    return [
      ...this.largeSummaryContexts(),
      groupedContext,
      ...this.smallSummaryContexts().map((context) => ({ ...context, groupedChild: true })),
    ];
  });
  readonly totalTracked = computed(
    () => durationAsHM(this.workspaceStats().totalDuration).trim() || '0m',
  );
  readonly topContext = computed(() => this.allSummaryContexts()[0]?.name ?? '-');

  toggleGroupedContexts(): void {
    this.showGroupedContexts.update((expanded) => !expanded);
  }

  startEdit(): void {
    const workspace = this.workspace();
    if (!workspace) {
      return;
    }

    this.editName.set(workspace.name);
    this.editDescription.set(workspace.description ?? '');
    this.isEditing.set(true);
  }

  cancelEdit(): void {
    this.isEditing.set(false);
    this.editName.set('');
    this.editDescription.set('');
  }

  getInputValue(event: Event): string {
    return (event.target as HTMLInputElement | HTMLTextAreaElement).value;
  }

  saveEdit(): void {
    const workspace = this.workspace();
    const name = this.editName().trim();
    if (!workspace || !name) {
      return;
    }

    this.updateWorkspaceMutation.mutate(
      {
        ...workspace,
        name,
        description: this.editDescription().trim(),
      },
      {
        onSuccess: () => this.cancelEdit(),
      },
    );
  }

  deleteWorkspace(): void {
    const workspace = this.workspace();
    if (!workspace) {
      return;
    }

    this.deleteWorkspaceMutation.mutate(workspace.id);
  }
}

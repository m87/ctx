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

const WORKSPACE_ROOT_VALUE = '__workspace_root__';

@Component({
  imports: [
    NameComponent,
    ContextIntervalListComponent,
    NgIcon,
    HlmButtonImports,
    HlmCardImports,
    QueryErrorStateComponent,
    RouterLink,
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
                [disabled]="freeContextMutation.isPending()"
                (click)="pauseContext()"
              >
                <ng-icon name="lucidePause"></ng-icon>
                <span class="font-semibold">Pause</span>
              </button>
            } @else {
              <button
                hlmBtn
                variant="outline"
                class="h-9 px-3 text-xs bg-blue-200/70 text-blue-600"
                [disabled]="currentContext.archived || switchContextMutation.isPending()"
                (click)="startContext()"
              >
                <ng-icon name="lucidePlay"></ng-icon>
                <span class="font-semibold text-blue-600">Start</span>
              </button>
            }
          </div>
        </div>

        @if (editingProjectAssignment()) {
          <div class="w-full rounded-lg border bg-card p-3">
            <div
              class="text-[11px] font-semibold text-muted-foreground flex items-center gap-2 uppercase"
            >
              Project
            </div>
            <div class="mt-2 flex items-center gap-2">
              <select
                class="h-9 min-w-0 flex-1 rounded-md border border-border bg-background px-3 text-sm outline-none focus-visible:ring-2 focus-visible:ring-ring/50"
                [disabled]="updateContextMutation.isPending()"
                (change)="assignContextToProject($event)"
              >
                <option [value]="workspaceRootValue" [selected]="!currentContext.project?.id">
                  Workspace root
                </option>
                @for (project of availableProjects(); track project.id) {
                  <option
                    [value]="project.id"
                    [selected]="currentContext.project?.id === project.id"
                  >
                    {{ project.name }}
                  </option>
                }
              </select>
              <button
                type="button"
                class="size-9 shrink-0 rounded-md border text-muted-foreground hover:text-foreground hover:bg-muted/60 flex items-center justify-center"
                aria-label="Cancel project assignment"
                title="Cancel"
                [disabled]="updateContextMutation.isPending()"
                (click)="editingProjectAssignment.set(false)"
              >
                <ng-icon name="lucideX"></ng-icon>
              </button>
            </div>
          </div>
        } @else if (currentContext.project; as project) {
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
            (click)="editingProjectAssignment.set(true)"
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
                <div class="text-lg font-semibold" hlmCardContet>
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
                <div class="text-lg font-semibold" hlmCardContet>
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
                <div class="text-lg font-semibold" hlmCardContet>
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
                <div class="text-lg font-semibold" hlmCardContet>
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
  readonly editingProjectAssignment = signal(false);
  readonly workspaceRootValue = WORKSPACE_ROOT_VALUE;
  isActiveContext = computed(() => this.activeContextQuery.data()?.id === this.contextId());
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
    this.editingProjectAssignment.set(true);
  }

  assignContextToProject(event: Event): void {
    const context = this.context();
    if (!context) {
      return;
    }

    const selectedValue = (event.target as HTMLSelectElement).value;
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
        onSuccess: () => this.editingProjectAssignment.set(false),
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

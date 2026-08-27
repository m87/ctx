import { Component, computed, inject, signal } from '@angular/core';
import { toSignal } from '@angular/core/rxjs-interop';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { NgIcon, provideIcons } from '@ng-icons/core';
import {
  lucideChevronRight,
  lucideFolder,
  lucidePencil,
  lucideTrash2,
  lucideX,
} from '@ng-icons/lucide';
import { HlmButtonImports } from '@spartan-ng/helm/button';
import { injectMutation, injectQuery } from '@tanstack/angular-query-experimental';
import { map } from 'rxjs';
import { Store } from '@ngxs/store';
import { ProjectMutations } from '../../api/project/project.mutations';
import { ProjectQueries } from '../../api/project/project.queries';
import { NameComponent, NameSaveValue } from '../shared/name.component';
import { QueryErrorStateComponent } from '../shared/query-error-state.component';
import { SelectProject } from '../sidebar/workspace.state';
import { LinkifiedTextComponent } from '../shared/linkified-text.component';
import { colorHash } from '../utils';
import { HlmSkeletonImports } from '@spartan-ng/helm/skeleton';
import { InsightsEmptyStateComponent } from '../shared/insights-empty-state.component';
import { SearchSelectComponent, SearchSelectOption } from '../shared/search-select.component';

const WORKSPACE_ROOT_VALUE = '__workspace_root__';
type DetailView = 'overview' | 'insights';

@Component({
  selector: 'ctx-project',
  imports: [
    NameComponent,
    NgIcon,
    HlmButtonImports,
    QueryErrorStateComponent,
    RouterLink,
    LinkifiedTextComponent,
    InsightsEmptyStateComponent,
    HlmSkeletonImports,
    SearchSelectComponent,
  ],
  providers: [
    provideIcons({
      lucideTrash2,
      lucideFolder,
      lucideChevronRight,
      lucidePencil,
      lucideX,
    }),
  ],
  template: `
    <div
      class="w-full h-full overflow-hidden flex flex-col items-start justify-start p-4 md:p-6 gap-5 relative"
    >
      @if (showProjectError()) {
        <ctx-query-error-state
          class="flex-1 min-h-0"
          [error]="projectDetailsQuery.error()"
          [paused]="projectDetailsQuery.isPaused()"
          resourceName="project"
          [retrying]="projectDetailsQuery.isFetching()"
          (retry)="retryProject()"
        ></ctx-query-error-state>
      } @else if (projectDetailsQuery.isLoading()) {
        <div class="w-full flex-1 min-h-0" role="status" aria-label="Loading project">
          <span class="sr-only">Loading project</span>
          <div class="flex justify-between gap-4 mb-5">
            <div class="flex-1">
              <hlm-skeleton class="h-2.5 w-14 mb-2"></hlm-skeleton>
              <hlm-skeleton class="h-7 w-52"></hlm-skeleton>
            </div>
            <hlm-skeleton class="size-9 mt-5"></hlm-skeleton>
          </div>
          <hlm-skeleton class="h-14 w-full mb-6"></hlm-skeleton>
          <hlm-skeleton class="h-2.5 w-20 mb-2"></hlm-skeleton>
          <div class="flex flex-col gap-2 mb-6">
            @for (item of listSkeletonItems; track item) {
              <hlm-skeleton class="h-12 w-full"></hlm-skeleton>
            }
          </div>
          <hlm-skeleton class="h-2.5 w-16 mb-2"></hlm-skeleton>
          <div class="flex flex-col gap-2">
            @for (item of listSkeletonItems; track item) {
              <hlm-skeleton class="h-12 w-full"></hlm-skeleton>
            }
          </div>
        </div>
      } @else if (project(); as currentProject) {
        <div class="w-full flex flex-col md:flex-row justify-between items-start gap-4">
          <ctx-name
            class="w-full min-w-0"
            label="Project"
            [name]="currentProject.name"
            [showDescription]="false"
            namePlaceholder="Project name"
            [savePending]="updateProjectMutation.isPending()"
            (save)="saveProjectName($event)"
          ></ctx-name>

          <div class="flex items-center gap-2 w-full md:w-auto flex-nowrap md:pt-5">
            <button
              hlmBtn
              variant="outline"
              class="size-9 p-0 text-xs bg-red-100/70 text-red-700"
              aria-label="Delete project"
              title="Delete"
              [disabled]="deleteProjectMutation.isPending()"
              (click)="deleteProject()"
            >
              <ng-icon name="lucideTrash2"></ng-icon>
            </button>
          </div>
        </div>

        <div
          class="inline-flex self-start rounded-lg bg-muted p-1 shrink-0"
          role="tablist"
          aria-label="Project view"
        >
          <button
            type="button"
            id="project-overview-tab"
            class="rounded-md px-3 py-1.5 text-xs font-medium transition-colors"
            [class.bg-background]="detailView() === 'overview'"
            [class.shadow-sm]="detailView() === 'overview'"
            [class.text-foreground]="detailView() === 'overview'"
            [class.text-muted-foreground]="detailView() !== 'overview'"
            role="tab"
            aria-controls="project-view-panel"
            [attr.aria-selected]="detailView() === 'overview'"
            (click)="detailView.set('overview')"
          >
            Overview
          </button>
          <button
            type="button"
            id="project-insights-tab"
            class="rounded-md px-3 py-1.5 text-xs font-medium transition-colors"
            [class.bg-background]="detailView() === 'insights'"
            [class.shadow-sm]="detailView() === 'insights'"
            [class.text-foreground]="detailView() === 'insights'"
            [class.text-muted-foreground]="detailView() !== 'insights'"
            role="tab"
            aria-controls="project-view-panel"
            [attr.aria-selected]="detailView() === 'insights'"
            (click)="detailView.set('insights')"
          >
            Insights
          </button>
        </div>

        <div
          id="project-view-panel"
          class="w-full flex-1 min-h-0 flex flex-col gap-5"
          role="tabpanel"
          [attr.aria-labelledby]="
            detailView() === 'overview' ? 'project-overview-tab' : 'project-insights-tab'
          "
        >
          @if (detailView() === 'overview') {
            @if (editingParentAssignment()) {
              <div class="w-full rounded-lg border bg-card p-3">
                <label
                  for="parent-project-search"
                  class="text-[11px] font-semibold text-muted-foreground flex items-center gap-2 uppercase"
                >
                  Parent project
                </label>
                <div class="mt-2 flex items-start gap-2">
                  <ctx-search-select
                    class="min-w-0 flex-1"
                    inputId="parent-project-search"
                    ariaLabel="Parent project"
                    searchPlaceholder="Search projects…"
                    emptyText="No matching projects"
                    [options]="parentProjectSelectOptions()"
                    [value]="currentProject.parentId ?? workspaceRootValue"
                    [disabled]="updateProjectMutation.isPending()"
                    (selectionChange)="assignParentProject($event)"
                  ></ctx-search-select>
                  <button
                    type="button"
                    class="flex size-9 shrink-0 items-center justify-center rounded-md border text-muted-foreground hover:bg-muted/60 hover:text-foreground"
                    aria-label="Cancel parent project assignment"
                    title="Cancel"
                    [disabled]="updateProjectMutation.isPending()"
                    (click)="editingParentAssignment.set(false)"
                  >
                    <ng-icon name="lucideX"></ng-icon>
                  </button>
                </div>
              </div>
            } @else if (currentProject.parentId && allProjectsQuery.isLoading()) {
              <hlm-skeleton class="h-16 w-full"></hlm-skeleton>
            } @else if (parentProject(); as parent) {
              <div
                class="w-full flex items-center gap-3 rounded-lg border bg-card p-3 cursor-pointer hover:bg-muted/30 transition-colors"
                [routerLink]="['/project', parent.id]"
                role="link"
                tabindex="0"
                (click)="selectProject(parent.id)"
              >
                <ctx-name
                  class="min-w-0 flex-1"
                  label="Parent project"
                  [name]="parent.name"
                  [showDescription]="false"
                  [readonly]="true"
                  [compact]="true"
                  [accentColor]="itemColor(parent.id)"
                ></ctx-name>
                <button
                  type="button"
                  class="size-8 shrink-0 rounded-md border text-muted-foreground hover:text-foreground hover:bg-muted/60 flex items-center justify-center"
                  aria-label="Change parent project"
                  title="Change parent project"
                  (click)="startParentAssignmentEdit($event)"
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
                (click)="editingParentAssignment.set(true)"
              >
                <ng-icon name="lucideFolder"></ng-icon>
                <span>Assign parent project</span>
              </button>
            }

            <div class="w-full flex-1 min-h-0 overflow-auto pr-1 pb-2">
              <section>
                <div
                  class="text-[11px] uppercase tracking-[0.08em] text-muted-foreground font-semibold mb-2"
                >
                  Subprojects
                </div>
                @if (projectSubprojectsQuery.isLoading()) {
                  <div class="flex flex-col gap-2" role="status">
                    <span class="sr-only">Loading subprojects</span>
                    @for (item of listSkeletonItems; track item) {
                      <hlm-skeleton class="h-12 w-full"></hlm-skeleton>
                    }
                  </div>
                } @else if ((projectSubprojectsQuery.data() ?? []).length > 0) {
                  <div class="flex flex-col gap-2">
                    @for (subproject of projectSubprojectsQuery.data() ?? []; track subproject.id) {
                      <div
                        class="block cursor-pointer rounded-lg border bg-card p-3 hover:bg-muted/30 transition-colors"
                        [routerLink]="['/project', subproject.id]"
                        role="link"
                        tabindex="0"
                        (click)="selectProject(subproject.id)"
                      >
                        <div class="flex items-center gap-2">
                          <ng-icon
                            name="lucideFolder"
                            class="text-sm shrink-0"
                            [style.color]="itemColor(subproject.id)"
                          ></ng-icon>
                          <span class="text-sm font-medium min-w-0 flex-1 truncate">
                            <ctx-linkified-text [text]="subproject.name" />
                          </span>
                          <ng-icon
                            name="lucideChevronRight"
                            class="text-sm shrink-0 text-muted-foreground/70"
                          ></ng-icon>
                        </div>
                      </div>
                    }
                  </div>
                } @else {
                  <p class="text-xs text-muted-foreground">No subprojects.</p>
                }
              </section>

              <section class="mt-6">
                <div
                  class="text-[11px] uppercase tracking-[0.08em] text-muted-foreground font-semibold mb-2"
                >
                  Contexts
                </div>
                @if (projectContextsQuery.isLoading()) {
                  <div class="flex flex-col gap-2" role="status">
                    <span class="sr-only">Loading project contexts</span>
                    @for (item of listSkeletonItems; track item) {
                      <hlm-skeleton class="h-12 w-full"></hlm-skeleton>
                    }
                  </div>
                } @else if ((projectContextsQuery.data() ?? []).length > 0) {
                  <div class="flex flex-col gap-2">
                    @for (context of projectContextsQuery.data() ?? []; track context.id) {
                      <div
                        class="block cursor-pointer rounded-lg border bg-card p-3 hover:bg-muted/30 transition-colors"
                        [routerLink]="['/context', context.id]"
                        role="link"
                        tabindex="0"
                      >
                        <div class="flex items-center gap-2">
                          <span
                            class="w-2 h-2 rounded-sm shrink-0"
                            [style.background-color]="itemColor(context.id)"
                          ></span>
                          <span class="text-sm font-medium min-w-0 flex-1 truncate">
                            <ctx-linkified-text [text]="context.name" />
                          </span>
                          @if (context.archived) {
                            <span
                              class="text-[10px] font-medium rounded border px-1.5 py-0.5 text-muted-foreground"
                            >
                              Archived
                            </span>
                          }
                          <ng-icon
                            name="lucideChevronRight"
                            class="text-sm shrink-0 text-muted-foreground/70"
                          ></ng-icon>
                        </div>
                      </div>
                    }
                  </div>
                } @else {
                  <p class="text-xs text-muted-foreground">No contexts.</p>
                }
              </section>
            </div>
          } @else {
            <ctx-insights-empty-state />
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
export class ProjectComponent {
  readonly listSkeletonItems = [0, 1, 2];
  readonly detailView = signal<DetailView>('overview');
  private readonly projectQueries = inject(ProjectQueries);
  private readonly projectMutations = inject(ProjectMutations);
  private readonly activeRoute = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly store = inject(Store);

  readonly projectId = toSignal(
    this.activeRoute.paramMap.pipe(map((params) => params.get('id') ?? '')),
    { initialValue: '' },
  );
  readonly selectedWorkspaceId = this.store.selectSignal(
    (state) => state.workspace.selectedWorkspaceId,
  );
  readonly projectDetailsQuery = injectQuery(() => this.projectQueries.get(this.projectId()));
  readonly projectWorkspaceId = computed(
    () => this.projectDetailsQuery.data()?.workspaceId ?? this.selectedWorkspaceId() ?? '',
  );
  readonly allProjectsQuery = injectQuery(() => this.projectQueries.all(this.projectWorkspaceId()));
  readonly updateProjectMutation = injectMutation(() => this.projectMutations.update());
  readonly deleteProjectMutation = injectMutation(() => this.projectMutations.delete());
  readonly projectContextsQuery = injectQuery(() => this.projectQueries.contexts(this.projectId()));
  readonly projectSubprojectsQuery = injectQuery(() =>
    this.projectQueries.subprojects(this.projectId(), this.selectedWorkspaceId()),
  );
  readonly project = computed(() => this.projectDetailsQuery.data() ?? null);
  readonly parentProject = computed(() => {
    const parentId = this.project()?.parentId;
    return parentId
      ? ((this.allProjectsQuery.data() ?? []).find((project) => project.id === parentId) ?? null)
      : null;
  });
  readonly availableParentProjects = computed(() => {
    const currentProjectId = this.projectId();
    const projects = this.allProjectsQuery.data() ?? [];
    const projectsById = new Map(projects.map((project) => [project.id, project]));

    return projects.filter((candidate) => {
      if (candidate.id === currentProjectId) {
        return false;
      }

      const visited = new Set<string>();
      let ancestorId = candidate.parentId;
      while (ancestorId) {
        if (ancestorId === currentProjectId) {
          return false;
        }
        if (visited.has(ancestorId)) {
          return false;
        }
        visited.add(ancestorId);
        ancestorId = projectsById.get(ancestorId)?.parentId;
      }
      return true;
    });
  });
  readonly parentProjectSelectOptions = computed<SearchSelectOption[]>(() => [
    {
      value: WORKSPACE_ROOT_VALUE,
      label: 'Workspace root',
      color: 'var(--muted-foreground)',
      description: 'No parent project',
    },
    ...this.availableParentProjects().map((project) => ({
      value: project.id,
      label: project.name,
      color: colorHash(project.id),
    })),
  ]);
  readonly editingParentAssignment = signal(false);
  readonly workspaceRootValue = WORKSPACE_ROOT_VALUE;
  readonly showProjectError = computed(
    () =>
      this.projectDetailsQuery.data() === undefined &&
      (this.projectDetailsQuery.isError() || this.projectDetailsQuery.isPaused()),
  );

  saveProjectName(value: NameSaveValue): void {
    const project = this.project();
    if (!project) {
      return;
    }

    this.updateProjectMutation.mutate({
      ...project,
      name: value.name,
    });
  }

  startParentAssignmentEdit(event: MouseEvent): void {
    event.preventDefault();
    event.stopPropagation();
    this.editingParentAssignment.set(true);
  }

  assignParentProject(selectedValue: string): void {
    const project = this.project();
    if (!project) {
      return;
    }

    const parentId = selectedValue === WORKSPACE_ROOT_VALUE ? '' : selectedValue;
    this.updateProjectMutation.mutate(
      {
        ...project,
        parentId: parentId || undefined,
      },
      {
        onSuccess: () => this.editingParentAssignment.set(false),
      },
    );
  }

  deleteProject(): void {
    const project = this.project();
    if (!project) {
      return;
    }

    if (!window.confirm(`Delete project "${project.name}"?`)) {
      return;
    }

    this.deleteProjectMutation.mutate(project, {
      onSuccess: () => {
        const parentId = project.parentId ?? null;
        this.store.dispatch(new SelectProject(parentId));
        void this.router.navigate(parentId ? ['/project', parentId] : ['/day']);
      },
    });
  }

  selectProject(projectId: string): void {
    this.store.dispatch(new SelectProject(projectId));
  }

  itemColor(id: string): string {
    return colorHash(id);
  }

  retryProject(): void {
    void this.projectDetailsQuery.refetch();
  }
}

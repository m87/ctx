import { Component, computed, inject } from '@angular/core';
import { toSignal } from '@angular/core/rxjs-interop';
import { ActivatedRoute, Router } from '@angular/router';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideFolder, lucideTrash2 } from '@ng-icons/lucide';
import { HlmButtonImports } from '@spartan-ng/helm/button';
import { injectMutation, injectQuery } from '@tanstack/angular-query-experimental';
import { map } from 'rxjs';
import { Store } from '@ngxs/store';
import { ProjectMutations } from '../../api/project/project.mutations';
import { ProjectQueries } from '../../api/project/project.queries';
import { NameComponent, NameSaveValue } from '../shared/name.component';
import { QueryErrorStateComponent } from '../shared/query-error-state.component';
import { SelectProject } from '../sidebar/workspace.state';
import { HlmSeparatorImports } from '@spartan-ng/helm/separator';

@Component({
  selector: 'ctx-project',
  imports: [NameComponent, NgIcon, HlmButtonImports, QueryErrorStateComponent, HlmSeparatorImports],
  providers: [provideIcons({ lucideTrash2, lucideFolder })],
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

        <div class="w-full flex-1 overflow-hidden">
          <div class="flex flex-col gap-1.5 overflow-y-auto">
            @if ((projectSubprojectsQuery.data() ?? []).length > 0) {
              @for (subproject of projectSubprojectsQuery.data() ?? []; track subproject.id) {
                <div
                  class="flex items-center gap-2 text-[13px] pl-2 pr-1 py-1 cursor-pointer rounded-lg border bg-card p-3 hover:bg-muted/30 transition-colors"
                >
                  <ng-icon name="lucideFolder" class="text-[10px] text-muted-foreground"></ng-icon>
                  <span class="min-w-0 flex-1 truncate">{{ subproject.name }}</span>
                </div>
              }
            } @else {
              <div class="text-[13px] text-muted-foreground mb-1.5">No subprojects.</div>
            }
          </div>
          <hlm-separator class="my-2"></hlm-separator>
          <div class="flex flex-col gap-1.5 overflow-y-auto mt-4">
            @if ((projectContextsQuery.data() ?? []).length > 0) {
              @for (context of projectContextsQuery.data() ?? []; track context.id) {
                <div
                  class="flex items-center gap-2 text-[13px] pl-2 pr-1 py-1 font-medium cursor-pointer rounded-lg border bg-card p-3 hover:bg-muted/30 transition-colors"
                >
                  <span class="min-w-0 flex-1 truncate">{{ context.name }}</span>
                </div>
              }
            } @else {
              <div class="text-[13px] text-muted-foreground">No contexts.</div>
            }
          </div>
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
  private readonly projectQueries = inject(ProjectQueries);
  private readonly projectMutations = inject(ProjectMutations);
  private readonly activeRoute = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly store = inject(Store);

  readonly projectId = toSignal(
    this.activeRoute.paramMap.pipe(map((params) => params.get('id') ?? '')),
    { initialValue: '' },
  );
  readonly selecetedWorkspaceId = this.store.selectSignal(
    (state) => state.workspace.selectedWorkspaceId,
  );
  readonly projectDetailsQuery = injectQuery(() => this.projectQueries.get(this.projectId()));
  readonly updateProjectMutation = injectMutation(() => this.projectMutations.update());
  readonly deleteProjectMutation = injectMutation(() => this.projectMutations.delete());
  readonly projectContextsQuery = injectQuery(() => this.projectQueries.contexts(this.projectId()));
  readonly projectSubprojectsQuery = injectQuery(() =>
    this.projectQueries.subprojects(this.projectId(), this.selecetedWorkspaceId()),
  );
  readonly project = computed(() => this.projectDetailsQuery.data() ?? null);
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

  retryProject(): void {
    void this.projectDetailsQuery.refetch();
  }
}

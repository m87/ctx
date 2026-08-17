import { Component, inject, signal } from '@angular/core';
import { NgIcon, provideIcons } from '@ng-icons/core';
import {
  lucideArrowLeft,
  lucideCheck,
  lucideChevronRight,
  lucideFolder,
  lucidePlus,
  lucideX,
} from '@ng-icons/lucide';
import { Store } from '@ngxs/store';
import { HlmButtonImports } from '@spartan-ng/helm/button';
import { HlmIcon } from '@spartan-ng/helm/icon';
import { ProjectQueries } from '../../api/project/project.queries';
import { injectMutation, injectQuery } from '@tanstack/angular-query-experimental';
import { ProjectMutations } from '../../api/project/project.mutations';
import { SelectProject } from './workspace.state';
import { Router, RouterLink } from '@angular/router';
import { LinkifiedTextComponent } from '../shared/linkified-text.component';
import { HlmSkeletonImports } from '@spartan-ng/helm/skeleton';

@Component({
  selector: 'ctx-project-widget',
  imports: [
    NgIcon,
    HlmIcon,
    HlmButtonImports,
    RouterLink,
    LinkifiedTextComponent,
    HlmSkeletonImports,
  ],
  providers: [
    provideIcons({
      lucidePlus,
      lucideArrowLeft,
      lucideCheck,
      lucideX,
      lucideFolder,
      lucideChevronRight,
    }),
  ],
  template: `
    <div class="flex flex-col h-full">
      <div class="flex-1 overflow-auto">
        <div
          class="flex items-center justify-between w-full gap-1 px-4 py-2 text-sm font-semibold text-muted-foreground"
        >
          @if (selectedProjectId()) {
            <div class="flex min-w-0 items-center">
              <ng-icon
                hlm
                name="lucideArrowLeft"
                size="18px"
                class="shrink-0 cursor-pointer"
                (click)="selectProject(this.projectDetailsQuery.data()?.parentId ?? '')"
              ></ng-icon>
              @if (projectDetailsQuery.isLoading()) {
                <hlm-skeleton class="h-3.5 w-24 ml-1"></hlm-skeleton>
              } @else {
                <span
                  class="ml-1 min-w-0 truncate cursor-pointer"
                  [routerLink]="['/project', selectedProjectId()]"
                >
                  <ctx-linkified-text [text]="projectDetailsQuery.data()?.name ?? ''" />
                </span>
              }
            </div>
          } @else {
            <span class="ml-1 text-xs">PROJECTS</span>
          }
          <ng-icon
            hlm
            name="lucidePlus"
            size="20px"
            class="w-4 h-4 mr-1 justify-end cursor-pointer"
            (click)="newProject.set(!newProject())"
          ></ng-icon>
        </div>

        @if (newProject()) {
          <div class="flex border rounded-md bg-card border-dashed items-center mx-2 my-2">
            <input
              type="text"
              placeholder="New project name"
              class="w-full px-2 py-1 focus:outline-none focus:ring focus:border-blue-300 outline-none bg-transparent text-sm"
              [value]="newProjectName()"
              (input)="onNewProjectNameInput($event)"
              (keydown.enter)="createProject()"
              (keydown.escape)="cancelProject()"
            />
            <button
              hlmBtn
              variant="ghost"
              [disabled]="createProjectMutation.isPending()"
              (click)="createProject()"
            >
              <ng-icon hlm name="lucideCheck" size="15px"></ng-icon>
            </button>
            <button hlmBtn variant="ghost" (click)="cancelProject()">
              <ng-icon hlm name="lucideX" size="15px"></ng-icon>
            </button>
          </div>
        }

        @if (subprojectsQuery.isLoading()) {
          <div class="px-5 py-2 flex flex-col gap-4" role="status">
            <span class="sr-only">Loading projects</span>
            @for (item of projectSkeletonItems; track item) {
              <div class="flex items-center gap-2">
                <hlm-skeleton class="size-3.5 shrink-0"></hlm-skeleton>
                <hlm-skeleton class="h-3.5 w-full" [class.max-w-28]="item % 2 === 0"></hlm-skeleton>
              </div>
            }
          </div>
        } @else {
          @for (project of subprojectsQuery.data(); track project.id) {
            <div
              class="flex items-center justify-between w-full gap-1 px-5 py-2 text-sm font-semibold text-muted-foreground hover:bg-muted/50 cursor-pointer"
              (click)="selectProject(project.id)"
            >
              <div class="flex items-center">
                <ng-icon hlm name="lucideFolder" size="14px" class="w-4 h-4 mr-1"></ng-icon>
                <span class="ml-1">
                  <ctx-linkified-text [text]="project.name" />
                </span>
              </div>
              <div>
                <ng-icon
                  hlm
                  name="lucideChevronRight"
                  size="14px"
                  class="w-4 h-4 ml-auto"
                ></ng-icon>
              </div>
            </div>
          }
        }
      </div>
    </div>
  `,
  styles: [],
})
export class ProjectWidgetComponent {
  readonly projectSkeletonItems = [0, 1, 2];
  private readonly store = inject(Store);
  private readonly projectQueries = inject(ProjectQueries);
  private readonly projectMutations = inject(ProjectMutations);
  readonly router = inject(Router);
  readonly subprojectsQuery = injectQuery(() =>
    this.projectQueries.subprojects(this.selectedProjectId(), this.selectedWorkspaceId()),
  );
  readonly projectDetailsQuery = injectQuery(() =>
    this.projectQueries.get(this.selectedProjectId() ?? ''),
  );
  readonly createProjectMutation = injectMutation(() => this.projectMutations.create());
  readonly selectedProjectId = this.store.selectSignal(
    (state) => state.workspace.selectedProjectId,
  );
  readonly selectedWorkspaceId = this.store.selectSignal(
    (state) => state.workspace.selectedWorkspaceId,
  );
  newProject = signal(false);
  newProjectName = signal('');

  onNewProjectNameInput(event: Event): void {
    this.newProjectName.set((event.target as HTMLInputElement).value);
  }

  createProject(): void {
    const name = this.newProjectName().trim();
    const workspaceId = this.selectedWorkspaceId();
    if (!name || !workspaceId) {
      return;
    }

    this.createProjectMutation.mutate(
      {
        id: '',
        name,
        workspaceId,
        parentId: this.selectedProjectId() ?? undefined,
      },
      {
        onSuccess: (project) => {
          this.newProject.set(false);
          this.newProjectName.set('');
          this.selectProject(project.id);
        },
      },
    );
  }

  cancelProject(): void {
    this.newProject.set(false);
    this.newProjectName.set('');
  }

  selectProject(projectId: string): void {
    const normalizedProjectId = projectId || null;
    this.store.dispatch(new SelectProject(normalizedProjectId));
    void this.router.navigate(normalizedProjectId ? ['/project', normalizedProjectId] : ['/day']);
  }
}

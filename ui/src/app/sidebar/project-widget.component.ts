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
import { injectQuery } from '@tanstack/angular-query-experimental';
import { SelectProject } from './workspace.state';
import { Router, RouterLink } from '@angular/router';

@Component({
  selector: 'ctx-project-widget',
  imports: [NgIcon, HlmIcon, HlmButtonImports, RouterLink],
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
            <div
              class="flex items-center cursor-pointer"
              [routerLink]="['/project', selectedProjectId()]"
            >
              <ng-icon
                hlm
                name="lucideArrowLeft"
                size="18px"
                (click)="selectProject(this.projectDetailsQuery.data()?.parentId ?? '')"
              ></ng-icon>
              <span class="ml-1">{{ projectDetailsQuery.data()?.name }}</span>
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
          <div class="flex border rounded-md bg-black border-dashed items-center mx-2 my-2">
            <input
              type="text"
              placeholder="New project name"
              class="w-full px-2 py-1 focus:outline-none focus:ring focus:border-blue-300 outline-none bg-transparent text-sm"
            />
            <button hlmBtn variant="ghost" (click)="createProject()">
              <ng-icon hlm name="lucideCheck" size="15px"></ng-icon>
            </button>
            <button hlmBtn variant="ghost" (click)="cancelProject()">
              <ng-icon hlm name="lucideX" size="15px"></ng-icon>
            </button>
          </div>
        }

        @if (subprojectsQuery.isSuccess()) {
          @for (project of subprojectsQuery.data(); track project.id) {
            <div
              class="flex items-center justify-between w-full gap-1 px-5 py-2 text-sm font-semibold text-muted-foreground hover:bg-muted/50 cursor-pointer"
              (click)="selectProject(project.id)"
            >
              <div class="flex items-center">
                <ng-icon hlm name="lucideFolder" size="14px" class="w-4 h-4 mr-1"></ng-icon>
                <span class="ml-1">{{ project.name }}</span>
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
        } @else {
          <div class="px-4 py-2 text-sm text-muted-foreground">Loading projects...</div>
        }
      </div>
    </div>
  `,
  styles: [],
})
export class ProjectWidgetComponent {
  private readonly store = inject(Store);
  private readonly projectQueries = inject(ProjectQueries);
  readonly router = inject(Router);
  readonly subprojectsQuery = injectQuery(() =>
    this.projectQueries.subprojects(this.selectedProjectId(), this.selectedWorkspaceId()),
  );
  readonly projectDetailsQuery = injectQuery(() =>
    this.projectQueries.get(this.selectedProjectId()),
  );
  readonly selectedProjectId = this.store.selectSignal(
    (state) => state.workspace.selectedProjectId,
  );
  readonly selectedWorkspaceId = this.store.selectSignal(
    (state) => state.workspace.selectedWorkspaceId,
  );
  newProject = signal(false);

  createProject() {
    this.newProject.set(false);
  }

  cancelProject() {
    this.newProject.set(false);
  }

  selectProject(projectId: string) {
    this.store.dispatch(new SelectProject(projectId));
    this.router.navigate(['/project', projectId]);
  }
}

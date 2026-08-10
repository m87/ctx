import { Component, inject } from '@angular/core';
import { toSignal } from '@angular/core/rxjs-interop';
import { ActivatedRoute } from '@angular/router';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { map } from 'rxjs';
import { ProjectQueries } from '../../api/project.queries';

@Component({
  selector: "ctx-project",
  imports: [],
  template: `
    <div class="flex flex-col h-full">
      <div class="flex-1 overflow-auto">
        @if (projectDetailsQuery.isLoading()) {
          <div class="flex items-center justify-center h-full">
            <span class="text-sm text-muted-foreground">Loading project details...</span>
          </div>
        } @else if (projectDetailsQuery.isError()) {
          <div class="flex items-center justify-center h-full">
            <span class="text-sm text-destructive">Error loading project details.</span>
          </div>
        } @else {
          <div class="p-4">
            <h2 class="text-lg font-semibold">{{ projectDetailsQuery.data()?.name }}</h2>
            <p class="text-sm text-muted-foreground">Project ID: {{ projectDetailsQuery.data()?.id }}</p>
            <p class="text-sm text-muted-foreground">Workspace ID: {{ projectDetailsQuery.data()?.workspaceId }}</p>
            @if (projectDetailsQuery.data()?.parentId) {
              <p class="text-sm text-muted-foreground">Parent Project ID: {{ projectDetailsQuery.data()?.parentId }}</p>
            }
          </div>
        }
      </div>
    </div>
  `,
  styles: [],
})
export class ProjectComponent {
  private readonly projectQueries = inject(ProjectQueries);
  private readonly activeRoute = inject(ActivatedRoute);

  readonly projectId = toSignal(
    this.activeRoute.paramMap.pipe(map((params) => params.get('id') ?? '')),
    { initialValue: '' },
  );
  readonly projectDetailsQuery = injectQuery(() =>
    this.projectQueries.get(this.projectId()),
  );
}

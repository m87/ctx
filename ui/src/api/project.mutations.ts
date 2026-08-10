import { inject, Injectable } from '@angular/core';
import { Project, ProjectService } from './project.service';
import { CacheService } from './cache.service';
import { lastValueFrom } from 'rxjs';
import { mutationOptions } from '@tanstack/angular-query-experimental';

@Injectable({
  providedIn: 'root',
})
export class ProjectMutations {
  private readonly projectService = inject(ProjectService);
  private readonly cache = inject(CacheService);

  create() {
    return mutationOptions({
      mutationFn: (project: Project) => lastValueFrom(this.projectService.create(project)),
      onSuccess: (data) => {
        this.cache.afterProjectChange(data.workspaceId);
      },
    });
  }

  update() {
    return mutationOptions({
      mutationFn: (project: Project) =>
        lastValueFrom(this.projectService.update(project.id, project)),
      onSuccess: (data) => this.cache.afterProjectChange(data.workspaceId),
    });
  }

  delete() {
    return mutationOptions({
      mutationFn: (project: Project) => lastValueFrom(this.projectService.delete(project.id)),
      onSuccess: (_data, project) => this.cache.afterProjectDelete(project.id, project.workspaceId),
    });
  }
}

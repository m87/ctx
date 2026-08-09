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
}

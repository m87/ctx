import { inject, Injectable } from '@angular/core';
import { ProjectService } from './project.service';
import { lastValueFrom } from 'rxjs';

export const projectQueryKeys = {
  all: ['project'] as const,
  detail: (projectId: string) => [...projectQueryKeys.all, projectId] as const,
};

@Injectable({
  providedIn: 'root',
})
export class ProjectQueries {
  private readonly projectService = inject(ProjectService);
  get(projectId: string) {
    return {
      queryKey: projectQueryKeys.detail(projectId),
      queryFn: () => lastValueFrom(this.projectService.get(projectId)),
    };
  }
}

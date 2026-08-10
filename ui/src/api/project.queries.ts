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

  subprojects(projectId: string, workspaceId: string) {
    return {
      queryKey: [...projectQueryKeys.detail(projectId), 'subprojects'] as const,
      queryFn: () => lastValueFrom(this.projectService.subprojects(projectId, workspaceId)),
    };
  }

  contexts(projectId: string) {
    return {
      queryKey: [...projectQueryKeys.detail(projectId), 'contexts'] as const,
      queryFn: () => lastValueFrom(this.projectService.contexts(projectId)),
    };
  }
}

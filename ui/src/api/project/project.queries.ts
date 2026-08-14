import { inject, Injectable } from '@angular/core';
import { ProjectService } from './project.service';
import { lastValueFrom } from 'rxjs';

export const projectQueryKeys = {
  all: ['project'] as const,
  lists: () => [...projectQueryKeys.all, 'list'] as const,
  list: (workspaceId: string) => [...projectQueryKeys.lists(), workspaceId] as const,
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
      enabled: projectId.length > 0,
    };
  }

  all(workspaceId: string) {
    return {
      queryKey: projectQueryKeys.list(workspaceId),
      queryFn: () => lastValueFrom(this.projectService.all(workspaceId)),
      enabled: workspaceId.length > 0,
    };
  }

  subprojects(projectId: string, workspaceId: string) {
    return {
      queryKey: [...projectQueryKeys.detail(projectId), 'subprojects', workspaceId] as const,
      queryFn: () => lastValueFrom(this.projectService.subprojects(projectId, workspaceId)),
      enabled: workspaceId.length > 0,
    };
  }

  contexts(projectId: string) {
    return {
      queryKey: [...projectQueryKeys.detail(projectId), 'contexts'] as const,
      queryFn: () => lastValueFrom(this.projectService.contexts(projectId)),
      enabled: projectId.length > 0,
    };
  }
}

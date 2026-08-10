import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Context } from './context.service';

export interface Project {
  id: string;
  name: string;
  workspaceId: string;
  parentId?: string;
}

@Injectable({ providedIn: 'root' })
export class ProjectService {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = '/api/project';

  get(projectId: string): Observable<Project> {
    return this.http.get<Project>(this.url(projectId));
  }

  create(project: Project): Observable<Project> {
    return this.http.post<Project>(this.url(), project);
  }

  update(projectId: string, project: Project): Observable<Project> {
    return this.http.put<Project>(this.url(projectId), project);
  }

  delete(projectId: string): Observable<void> {
    return this.http.delete<void>(this.url(projectId));
  }

  subprojects(projectId: string, workspaceId: string): Observable<Project[]> {
    if (projectId) {
      return this.http.get<Project[]>(this.url(projectId, 'projects'), { params: { workspaceId } });
    } else {
      return this.http.get<Project[]>(this.url('projects'), { params: { workspaceId } });
    }
  }

  contexts(projectId: string): Observable<Context[]> {
    return this.http.get<Context[]>(this.url(projectId, 'contexts'));
  }

  private url(...segments: string[]): string {
    let url = [this.baseUrl, ...segments].join('/');
    if (segments.length === 0 && !url.endsWith('/')) {
      url += '/';
    }
    return url;
  }
}

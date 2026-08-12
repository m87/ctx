import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Observable } from 'rxjs/internal/Observable';
import { Context } from '../context/context.service';

export interface Workspace {
  id: string;
  name: string;
  description?: string;
  properties: WorkspaceProperties;
}

export interface WorkspaceProperties {
  linkRules?: LinkRule[];
}

export interface LinkRule {
  regexp: string;
  link: string;
}

export interface WorkspaceContextStats {
  contextId: string;
  duration: number;
  percentage: number;
  intervalCount: number;
}

export interface WorkspaceStats {
  workspaceId: string;
  contexts: Context[];
  contextStats: WorkspaceContextStats[];
  totalDuration: number;
  totalSessions: number;
}

@Injectable({
  providedIn: 'root',
})
export class WorkspaceService {
  http = inject(HttpClient);
  private readonly baseUrl = '/api/workspace';

  listWorkspaces(): Observable<Workspace[]> {
    return this.http.get<Workspace[]>(this.url());
  }

  createWorkspace(name: string): Observable<Workspace> {
    return this.http.post<Workspace>(this.url(), { name });
  }

  deleteWorkspace(id: string): Observable<void> {
    return this.http.delete<void>(this.url(id));
  }

  getWorkspace(id: string): Observable<Workspace> {
    return this.http.get<Workspace>(this.url(id));
  }

  getWorkspaceStats(id: string): Observable<WorkspaceStats> {
    return this.http.get<WorkspaceStats>(this.url(id, 'stats'));
  }

  updateWorkspace(workspace: Workspace): Observable<Workspace> {
    return this.http.put<Workspace>(this.url(workspace.id), workspace);
  }

  private url(...segments: string[]): string {
    let url = [this.baseUrl, ...segments].join('/');
    if(segments.length === 0 && !url.endsWith('/')) {
      url += '/';
    }
    return url;
  }
}

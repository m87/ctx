import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Observable } from 'rxjs/internal/Observable';
import { Context } from './context.service';

export interface Workspace {
  id: string;
  name: string;
  description?: string;
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

  listWorkspaces(): Observable<Workspace[]> {
    return this.http.get<Workspace[]>('/api/workspace/');
  }

  createWorkspace(name: string): Observable<Workspace> {
    return this.http.post<Workspace>('/api/workspace/', { name });
  }

  deleteWorkspace(id: string): Observable<void> {
    return this.http.delete<void>(`/api/workspace/${id}`);
  }

  getWorkspace(id: string): Observable<Workspace> {
    return this.http.get<Workspace>(`/api/workspace/${id}`);
  }

  getWorkspaceStats(id: string): Observable<WorkspaceStats> {
    return this.http.get<WorkspaceStats>(`/api/workspace/${id}/stats`);
  }

  updateWorkspace(workspace: Workspace): Observable<Workspace> {
    return this.http.put<Workspace>(`/api/workspace/${workspace.id}`, workspace);
  }
}

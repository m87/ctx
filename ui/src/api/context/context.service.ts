import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { catchError, Observable, of, throwError } from 'rxjs';
import { hasApiErrorCode } from '../error';
import { Interval } from '../interval/interval.service';
import { Store } from '@ngxs/store';
import { WorkspaceState } from '../../app/sidebar/workspace.state';

export interface ProjectMetadata {
  id: string;
  name: string;
}

export interface Context {
  id: string;
  name: string;
  description?: string;
  workspaceId: string;
  archived?: boolean;
  status?: string;
  tags?: string[];
  project?: ProjectMetadata;
}

export type CreateContextInput = Pick<Context, 'name'> &
  Partial<Pick<Context, 'description' | 'workspaceId' | 'tags'>>;

export type SwitchContextInput = CreateContextInput & Partial<Pick<Context, 'id'>>;

export interface ContextStats {
  contextId: string;
  date: string;
  totalDuration: number;
  sessions: number;
  totalSessions: number;
  duration: number;
}

@Injectable({
  providedIn: 'root',
})
export class ContextService {
  private readonly http = inject(HttpClient);
  private readonly store = inject(Store);
  private readonly baseUrl = '/api/context';

  getIntervals(contextId: string): Observable<Interval[]> {
    return this.http.get<Interval[]>(this.url(contextId, 'intervals'));
  }

  getActiveContext(): Observable<Context | null> {
    return this.http
      .get<Context>(this.url('active'))
      .pipe(
        catchError((error: unknown) =>
          hasApiErrorCode(error, 'ACTIVE_CONTEXT_NOT_FOUND') ? of(null) : throwError(() => error),
        ),
      );
  }

  getContexts(workspaceId: string, includeArchived = false): Observable<Context[]> {
    const params = new URLSearchParams({ workspaceId });
    if (includeArchived) {
      params.set('includeArchived', 'true');
    }
    return this.http.get<Context[]>(this.url(`?${params.toString()}`));
  }

  createContext(context: CreateContextInput): Observable<Context> {
    return this.http.post<Context>(this.url(), this.withWorkspace(context));
  }

  deleteContext(id: string): Observable<void> {
    return this.http.delete<void>(this.url(id));
  }

  updateContext(id: string, context: Context): Observable<Context> {
    return this.http.put<Context>(this.url(id), context);
  }

  getContext(id: string): Observable<Context> {
    return this.http.get<Context>(this.url(id));
  }

  switchContext(context: SwitchContextInput): Observable<void> {
    return this.http.post<void>(this.url('switch'), this.withWorkspace(context));
  }

  freeContext(): Observable<void> {
    return this.http.post<void>(this.url('free'), {});
  }

  getStats(contextId: string, date: string, timeZone: string): Observable<ContextStats> {
    const params = new URLSearchParams({ timeZone });
    const path = [this.baseUrl, contextId, 'stats', date].join('/');
    return this.http.get<ContextStats>(`${path}?${params.toString()}`);
  }

  archiveContext(contextId: string): Observable<void> {
    return this.http.post<void>(this.url(contextId, 'archive'), {});
  }

  restoreContext(contextId: string): Observable<void> {
    return this.http.post<void>(this.url(contextId, 'restore'), {});
  }

  private withWorkspace<T extends CreateContextInput>(context: T): T & { workspaceId: string } {
    return {
      ...context,
      workspaceId:
        context.workspaceId ?? this.store.selectSnapshot(WorkspaceState.selectedWorkspaceId)!,
    };
  }

  private url(...segments: string[]): string {
    let url = [this.baseUrl, ...segments].join('/');
    if (segments.length === 0 && !url.endsWith('/')) {
      url += '/';
    }
    return url;
  }
}

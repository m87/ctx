import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { catchError, map, Observable, of } from 'rxjs';
import { deserializeIntervals, Interval, RawInterval } from './interval.service';
import { Store } from '@ngxs/store';
import { WorkspaceState } from '../app/sidebar/workspace.state';

export interface Context {
  id: string;
  name: string;
  description: string;
  workspaceId: string;
  archived?: boolean;
  status?: string;
  tags?: string[];
}

export interface ContextStats {
  contextId: string;
  date: string;
  totalDuration: number;
  sessions: number;
  totalSessions: number;
  duration: number;
}

export interface DayContextStats {
  contextId: string;
  duration: number;
  percentage: number;
  intervalCount: number;
}

export interface DayStats {
  date: string;
  contextStats: DayContextStats[];
  contexts: Context[];
  intervals: { [key: string]: Interval[] };
  distribution: { [contextId: string]: number };
}

type RawDayStats = Omit<DayStats, 'intervals'> & {
  intervals: { [key: string]: RawInterval[] };
};

@Injectable({
  providedIn: 'root',
})
export class ContextService {
  http = inject(HttpClient);
  store = inject(Store);

  getIntervals(contextId: string): Observable<Interval[]> {
    return this.http
      .get<RawInterval[]>(`/api/context/${contextId}/intervals`)
      .pipe(map((intervals) => deserializeIntervals(intervals)));
  }

  getActiveContext(): Observable<Context> {
    return this.http
      .get<Context>('/api/context/active')
      .pipe(catchError(() => of<Context>({} as Context)));
  }

  getContexts(workspaceId: string): Observable<Context[]> {
    return this.http.get<Context[]>(`/api/context/?workspaceId=${workspaceId}`);
  }

  createContext(context: Context): Observable<Context> {
    if (context.workspaceId == null) {
      context.workspaceId = this.store.selectSnapshot(WorkspaceState.selectedWorkspaceId)!;
    }
    return this.http.post<Context>('/api/context/', context);
  }

  deleteContext(id: string): Observable<void> {
    return this.http.delete<void>(`/api/context/${id}`);
  }

  updateContext(id: string, context: Context): Observable<Context> {
    return this.http.put<Context>(`/api/context/${id}`, context);
  }

  getContext(id: string): Observable<Context> {
    return this.http.get<Context>(`/api/context/${id}`);
  }

  switchContext(context: Context): Observable<void> {
    if (context.workspaceId == null) {
      context.workspaceId = this.store.selectSnapshot(WorkspaceState.selectedWorkspaceId)!;
    }
    return this.http.post<void>(`/api/context/switch`, context);
  }

  freeContext(): Observable<void> {
    return this.http.post<void>('/api/context/free', {});
  }

  getStats(contextId: string, date: string): Observable<ContextStats> {
    return this.http.get<ContextStats>(`/api/context/${contextId}/stats/${date}`);
  }

  getDayStats(workspaceId: string, date: string): Observable<DayStats> {
    return this.http
      .get<RawDayStats>(`/api/interval/day/${date}/stats?workspaceId=${workspaceId}`)
      .pipe(
        map((response) => ({
          ...response,
          intervals: Object.fromEntries(
            Object.entries(response.intervals).map(([contextId, intervals]) => [
              contextId,
              deserializeIntervals(intervals),
            ]),
          ),
        })),
      );
  }
}

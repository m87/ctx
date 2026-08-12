import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Store } from '@ngxs/store';
import { Observable } from 'rxjs';
import { WorkspaceState } from '../../app/sidebar/workspace.state';
import type { Context } from '../context/context.service';

export interface Interval {
  id: string;
  contextId: string;
  start: string | null;
  end: string | null;
  duration: number;
  status?: string;
  workspaceId: string;
}

export interface DayIntervalsResponse {
  contexts: Context[];
  intervals: Interval[];
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

@Injectable({ providedIn: 'root' })
export class IntervalService {
  private readonly http = inject(HttpClient);
  private readonly store = inject(Store);
  private readonly baseUrl = '/api/interval';

  createInterval(interval: Interval): Observable<Interval> {
    if (interval.workspaceId == null) {
      interval.workspaceId = this.store.selectSnapshot(WorkspaceState.selectedWorkspaceId)!;
    }
    return this.http.post<Interval>(this.url('/'), interval);
  }

  deleteInterval(id: string): Observable<void> {
    return this.http.delete<void>(this.url(id));
  }

  updateInterval(id: string, interval: Interval): Observable<Interval> {
    return this.http.put<Interval>(this.url(id), interval);
  }

  moveInterval(id: string, targetContextId: string): Observable<Interval> {
    return this.http.patch<Interval>(this.url(id, 'move', targetContextId), {});
  }

  getInterval(id: string): Observable<Interval> {
    return this.http.get<Interval>(this.url(id));
  }

  getDayIntervals(
    workspaceId: string,
    day: string,
    timeZone: string,
  ): Observable<DayIntervalsResponse> {
    return this.http.get<DayIntervalsResponse>(
      this.urlWithParams({ workspaceId, timeZone }, 'day', day),
    );
  }

  getDayStats(workspaceId: string, date: string, timeZone: string): Observable<DayStats> {
    return this.http.get<DayStats>(
      this.urlWithParams({ workspaceId, timeZone }, 'day', date, 'stats'),
    );
  }

  private url(...segments: string[]): string {
    let url = [this.baseUrl, ...segments].join('/');
    if (segments.length === 0 && !url.endsWith('/')) {
      url += '/';
    }
    return url;
  }

  private urlWithParams(params: { [key: string]: string }, ...segments: string[]): string {
    let url = [this.baseUrl, ...segments].join('/');
    if (Object.keys(params).length > 0) {
      url += `?${new URLSearchParams(params).toString()}`;
    }
    return url;
  }
}

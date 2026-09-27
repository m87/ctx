import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Observable } from 'rxjs';

export type DashboardType = 'workspace' | 'project' | 'context' | 'daily' | 'custom';

export interface Dashboard {
  id: string;
  workspaceId: string;
  type: DashboardType;
  targetId?: string;
  name: string;
  active: boolean;
  definition: Record<string, unknown>;
}

export interface DashboardScope {
  projectId?: string;
  contextId?: string;
  date?: string;
}

export interface DashboardEvaluationRequest {
  workspaceId: string;
  definition: Record<string, unknown>;
  timeZone: string;
  firstDayOfWeek: number;
  scope: DashboardScope;
}

export interface DashboardMetricResult {
  value: number | null;
  unit: 'ms' | 'count' | 'percent';
  timeZone: string;
  start: string;
  end: string;
  error?: string;
  ignoredIntervals: number;
}

export interface DashboardEvaluation {
  asOf: string;
  widgets: Record<string, DashboardMetricResult>;
}

export type CreateDashboardInput = Omit<Dashboard, 'id' | 'active'>;

@Injectable({ providedIn: 'root' })
export class DashboardService {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = '/api/dashboard';

  list(workspaceId: string): Observable<Dashboard[]> {
    return this.http.get<Dashboard[]>(this.url(), { params: { workspaceId } });
  }

  get(id: string): Observable<Dashboard> {
    return this.http.get<Dashboard>(this.url(id));
  }

  create(input: CreateDashboardInput): Observable<Dashboard> {
    return this.http.post<Dashboard>(this.url(), input);
  }

  update(dashboard: Dashboard): Observable<Dashboard> {
    return this.http.put<Dashboard>(this.url(dashboard.id), dashboard);
  }

  evaluate(request: DashboardEvaluationRequest): Observable<DashboardEvaluation> {
    return this.http.post<DashboardEvaluation>(this.url('evaluate'), request);
  }

  delete(id: string): Observable<void> {
    return this.http.delete<void>(this.url(id));
  }

  private url(...segments: string[]): string {
    let url = [this.baseUrl, ...segments].join('/');
    if (segments.length === 0 && !url.endsWith('/')) {
      url += '/';
    }
    return url;
  }
}

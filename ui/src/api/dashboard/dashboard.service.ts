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

interface DashboardMetricResultMetadata {
  timeZone: string;
  start: string;
  end: string;
  error?: string;
  ignoredIntervals: number;
}

export interface DashboardScalarMetricResult extends DashboardMetricResultMetadata {
  value: number | null;
  unit: 'ms' | 'count' | 'percent';
  kind?: never;
}

export interface DashboardRadarMetricResult extends DashboardMetricResultMetadata {
  kind: 'radar';
  value: null;
  unit: 'ms' | 'count' | 'percent';
  radar: DashboardRadarPayload;
  hasActivity: boolean;
}

export interface DashboardCategoriesMetricResult extends DashboardMetricResultMetadata {
  kind: 'categories';
  value: null;
  unit: 'ms' | 'count' | 'percent';
  categories: DashboardCategoriesPayload;
  hasActivity: boolean;
}
export interface DashboardHeatmapMetricResult extends DashboardMetricResultMetadata {
  kind: 'heatmap';
  value: null;
  unit: 'ms' | 'count' | 'percent';
  heatmap: DashboardHeatmapPayload;
  hasActivity: boolean;
}
export interface DashboardSeriesMetricResult extends DashboardMetricResultMetadata {
  kind: 'series';
  value: null;
  unit: 'ms' | 'count' | 'percent';
  series: DashboardLinePayload;
  hasActivity: boolean;
}
export interface DashboardTableMetricResult extends DashboardMetricResultMetadata {
  kind: 'table';
  value: null;
  table: DashboardTablePayload;
  hasActivity: boolean;
}
export type DashboardMetricResult =
  | DashboardScalarMetricResult
  | DashboardRadarMetricResult
  | DashboardCategoriesMetricResult
  | DashboardHeatmapMetricResult
  | DashboardSeriesMetricResult
  | DashboardTableMetricResult;

export interface DashboardCategoriesPayload {
  unit: 'ms' | 'count' | 'percent';
  categories: {
    key: string;
    label: string;
    value: number | null;
    entity: { type: 'project' | 'context'; id: string } | null;
    synthetic: boolean;
    memberCount: number;
    start: string | null;
    end: string | null;
    status: 'complete' | 'partial' | 'future';
  }[];
  totalValue: number | null;
  totalCategories: number;
  omittedCategories: number;
}
export interface DashboardHeatmapPayload {
  unit: 'ms' | 'count' | 'percent';
  days: {
    date: string;
    start: string;
    end: string;
    value: number | null;
    status: 'complete' | 'partial' | 'future';
  }[];
  totalValue: number | null;
  activeDays: number;
}
export interface DashboardLinePayload {
  unit: 'ms' | 'count' | 'percent';
  buckets: { key: string; start: string; end: string; status: 'complete' | 'partial' | 'future' }[];
  series: { key: string; label: string; values: (number | null)[] }[];
  totalValue: number | null;
  totalSeries: number;
  omittedSeries: number;
}
export interface DashboardTablePayload {
  columns: { id: string; metric: string; unit: 'ms' | 'count' | 'percent' }[];
  rows: {
    key: string;
    label: string;
    entity: { type: 'project' | 'context'; id: string } | null;
    values: Record<string, number | null>;
  }[];
  summary: Record<string, number | null>;
  columnMaxima: Record<string, number | null>;
  totalRows: number;
}

export interface DashboardRadarPayload {
  unit: 'ms' | 'count' | 'percent';
  axes: DashboardRadarAxis[];
  totalValue: number | null;
}

export interface DashboardRadarAxis {
  key: string;
  label: string;
  value: number | null;
  status: 'complete' | 'partial' | 'unobserved';
  observedDays: number | null;
  entity: { type: 'project' | 'context'; id: string } | null;
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

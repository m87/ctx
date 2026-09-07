import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Context } from '../context/context.service';

export interface ContextQueryInput {
  workspaceId: string;
  query: string;
}

export interface ContextQueryStats {
  contextId: string;
  duration: number;
  percentage: number;
  intervalCount: number;
}

export interface ContextQueryResult extends ContextQueryInput {
  contexts: Context[];
  contextStats: ContextQueryStats[];
  totalDuration: number;
  totalSessions: number;
}

@Injectable({ providedIn: 'root' })
export class QueryService {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = '/api/query/';

  execute(input: ContextQueryInput): Observable<ContextQueryResult> {
    return this.http.post<ContextQueryResult>(`${this.baseUrl}result`, input);
  }
}

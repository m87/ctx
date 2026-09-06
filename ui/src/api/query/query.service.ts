import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Context } from '../context/context.service';

export interface ContextQueryInput {
  query: string;
}

@Injectable({ providedIn: 'root' })
export class QueryService {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = '/api/query/';

  execute(input: ContextQueryInput): Observable<Context[]> {
    return this.http.post<Context[]>(this.baseUrl, input);
  }
}

import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Observable } from 'rxjs';

export interface SavedQuery {
  id: string;
  workspaceId: string;
  name: string;
  query: string;
}

export type CreateSavedQueryInput = Omit<SavedQuery, 'id'>;

@Injectable({ providedIn: 'root' })
export class SavedQueryService {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = '/api/query/saved';

  list(workspaceId: string): Observable<SavedQuery[]> {
    const params = new URLSearchParams({ workspaceId });
    return this.http.get<SavedQuery[]>(`${this.url()}?${params.toString()}`);
  }

  get(id: string): Observable<SavedQuery> {
    return this.http.get<SavedQuery>(this.url(id));
  }

  create(input: CreateSavedQueryInput): Observable<SavedQuery> {
    return this.http.post<SavedQuery>(this.url(), input);
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

import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Observable } from 'rxjs';

export type Settings = { [key: string]: string };
@Injectable({
  providedIn: 'root',
})
export class SettingsService {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = '/api/settings';

  getSettings(): Observable<Settings> {
    return this.http.get<Settings>(this.url());
  }

  saveSettings(settings: Settings): Observable<void> {
    return this.http.patch<void>(this.url(), settings);
  }

  getSetting(key: string): Observable<string> {
    return this.http.get(this.url('key', encodeURIComponent(key)), { responseType: 'text' });
  }

  private url(...segments: string[]): string {
    let url = [this.baseUrl, ...segments].join('/');
    if(segments.length === 0 && !url.endsWith('/')) {
      url += '/';
    }
    return url;
  }
}

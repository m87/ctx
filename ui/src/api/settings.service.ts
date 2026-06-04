import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

export type Settings = { [key: string]: string };

@Injectable({
  providedIn: 'root',
})
export class SettingsService {
  constructor(private http: HttpClient) {}

  getSettings(): Observable<Settings> {
    return this.http.get<Settings>('/api/settings/');
  }

  saveSettings(settings: Settings): Observable<void> {
    return this.http.patch<void>('/api/settings/', settings);
  }

  getSetting(key: string): Observable<string> {
    return this.http.get(`/api/settings/key/${encodeURIComponent(key)}`, { responseType: 'text' });
  }
}

import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

export type Settings = { [key: string]: string };

export type IntegrityIssue = {
  entityType: 'context' | 'interval';
  entityId: string;
  code: string;
  description: string;
};

export type IntegrityReport = {
  healthy: boolean;
  workspaceCount: number;
  contextCount: number;
  intervalCount: number;
  issues: IntegrityIssue[];
};

export type IntegrityRepairResult = {
  repairedCount: number;
  report: IntegrityReport;
};

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

  checkIntegrity(): Observable<IntegrityReport> {
    return this.http.get<IntegrityReport>('/api/integrity/');
  }

  repairIntegrity(): Observable<IntegrityRepairResult> {
    return this.http.post<IntegrityRepairResult>('/api/integrity/repair', null);
  }
}

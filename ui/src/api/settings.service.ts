import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

export type Settings = { [key: string]: string };

export type IntegrityDateTime = {
  time: string | null;
  timezone: string | null;
  isZero: boolean | null;
};

export type IntegrityIssueDetails = {
  name?: string;
  contextId?: string;
  workspaceId?: string;
  start?: IntegrityDateTime;
  end?: IntegrityDateTime;
};

export type IntegrityIssue = {
  entityType: 'context' | 'interval';
  entityId: string;
  code: string;
  description: string;
  repairable: boolean;
  details?: IntegrityIssueDetails;
};

export type IntegrityReport = {
  healthy: boolean;
  workspaceCount: number;
  contextCount: number;
  intervalCount: number;
  issues: IntegrityIssue[];
};

export type IntegrityContextOption = {
  id: string;
  name: string;
  workspaceId: string;
  workspaceName: string;
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

  getIntegrityContexts(): Observable<IntegrityContextOption[]> {
    return this.http.get<IntegrityContextOption[]>('/api/integrity/contexts');
  }

  repairIntegrity(): Observable<IntegrityRepairResult> {
    return this.http.post<IntegrityRepairResult>('/api/integrity/repair', null);
  }
}

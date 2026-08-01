import { HttpClient } from "@angular/common/http";
import { inject, Injectable } from "@angular/core";
import { Observable } from "rxjs";


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



@Injectable({providedIn: 'root'})
export class IntegrityService {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = '/api/integrity';

  checkIntegrity(): Observable<IntegrityReport> {
    return this.http.get<IntegrityReport>(this.url());
  }

  getIntegrityContexts(): Observable<IntegrityContextOption[]> {
    return this.http.get<IntegrityContextOption[]>(this.url('contexts'));
  }

  repairIntegrity(): Observable<IntegrityRepairResult> {
    return this.http.post<IntegrityRepairResult>(this.url('repair'), null);
  }


  private url(...segments: string[]): string {
    let url = [this.baseUrl, ...segments].join('/');
    if(segments.length === 0 && !url.endsWith('/')) {
      url += '/';
    }
    return url;
  }
}

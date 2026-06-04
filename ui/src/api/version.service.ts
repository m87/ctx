import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Observable } from 'rxjs';

export interface VersionInfo {
  version: string;
}

@Injectable({
  providedIn: 'root',
})
export class VersionService {
  http = inject(HttpClient);

  getVersion(): Observable<VersionInfo> {
    return this.http.get<VersionInfo>('/api/version');
  }
}

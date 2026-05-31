import { HttpClient } from "@angular/common/http";
import { inject, Injectable } from "@angular/core";
import { Observable } from "rxjs/internal/Observable";


export interface Workspace {
    id: string;
    name: string;
}

@Injectable({
    providedIn: 'root',
})
export class WorkspaceService {
    http = inject(HttpClient);

    listWorkspaces(): Observable<Workspace[]> {
        return this.http.get<Workspace[]>('/api/workspace');
    }

    createWorkspace(name: string): Observable<Workspace> {
        return this.http.post<Workspace>('/api/workspace', { name });
    }

    deleteWorkspace(id: string): Observable<void> {
        return this.http.delete<void>(`/api/workspace/${id}`);
    }

    updateWorkspace(workspace: Workspace): Observable<Workspace> {
        return this.http.put<Workspace>(`/api/workspace/${workspace.id}`, workspace);
    }
} 
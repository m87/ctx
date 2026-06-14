import { Component, computed, inject, signal } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { map } from 'rxjs';
import { Store } from '@ngxs/store';
import { injectMutation, injectQuery } from '@tanstack/angular-query-experimental';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideCheck, lucidePencil, lucideTrash2, lucideX } from '@ng-icons/lucide';
import { WorkspaceQueries } from '../../api/workspace.quries';
import { WorkspaceMutations } from '../../api/workspace.mutations';
import { SelectWorkspace, WorkspaceState } from '../sidebar/workspace.state';

@Component({
  selector: 'app-workspace',
  imports: [NgIcon],
  providers: [provideIcons({ lucideCheck, lucidePencil, lucideTrash2, lucideX })],
  template: `
    <div class="w-full h-full overflow-hidden flex flex-col p-4 md:p-6">
      <div class="rounded-lg border bg-card px-3 py-2.5">
        <div class="flex items-center justify-between gap-3">
          <div class="min-w-0 flex-1">
            @if (workspace()) {
              @if (isEditingName()) {
                <input
                  type="text"
                  class="h-9 w-full rounded-md border bg-background px-2.5 text-2xl font-semibold tracking-tight outline-none focus:ring-1 focus:ring-ring"
                  [value]="workspaceName()"
                  (input)="onWorkspaceNameInput($event)"
                  (keydown.enter)="confirmRenameWorkspace()"
                  (keydown.escape)="cancelRenameWorkspace()"
                />
              } @else {
                <h1 class="text-2xl font-semibold tracking-tight truncate">
                  {{ workspace()?.name }}
                </h1>
              }
            } @else {
              <h1 class="text-2xl font-semibold tracking-tight">Default workspace</h1>
            }
          </div>

          @if (workspace()) {
            <div class="flex shrink-0 items-center gap-1">
              @if (isEditingName()) {
                <button
                  type="button"
                  class="h-8 w-8 rounded-md bg-primary text-primary-foreground hover:bg-primary/90 flex items-center justify-center"
                  aria-label="Save workspace name"
                  title="Save name"
                  (click)="confirmRenameWorkspace()"
                >
                  <ng-icon name="lucideCheck"></ng-icon>
                </button>
                <button
                  type="button"
                  class="h-8 w-8 rounded-md border hover:bg-muted/60 flex items-center justify-center"
                  aria-label="Cancel workspace rename"
                  title="Cancel"
                  (click)="cancelRenameWorkspace()"
                >
                  <ng-icon name="lucideX"></ng-icon>
                </button>
              } @else {
                <button
                  type="button"
                  class="h-8 w-8 rounded-md border text-muted-foreground hover:text-foreground hover:bg-muted/60 flex items-center justify-center"
                  aria-label="Edit workspace name"
                  title="Edit name"
                  (click)="startRenameWorkspace()"
                >
                  <ng-icon name="lucidePencil"></ng-icon>
                </button>
                <button
                  type="button"
                  class="h-8 w-8 rounded-md border border-destructive/30 text-destructive hover:bg-destructive/10 flex items-center justify-center"
                  aria-label="Delete workspace"
                  title="Delete"
                  (click)="deleteWorkspace()"
                >
                  <ng-icon name="lucideTrash2"></ng-icon>
                </button>
              }
            </div>
          }
        </div>
      </div>
    </div>
  `,
  styles: `
    :host {
      display: block;
      width: 100%;
      max-width: 1000px;
      height: 100%;
      min-height: 0;
    }
  `,
})
export class WorkspaceComponent {
  private readonly route = inject(ActivatedRoute);
  private readonly store = inject(Store);
  private readonly workspaceQueries = inject(WorkspaceQueries);
  private readonly workspaceMutations = inject(WorkspaceMutations);

  private readonly routeWorkspaceId = toSignal(
    this.route.paramMap.pipe(map((params) => params.get('id'))),
    { initialValue: null },
  );
  private readonly selectedWorkspaceId = this.store.selectSignal(WorkspaceState.selectedWorkspaceId);

  listWorkspacesQuery = injectQuery(() => this.workspaceQueries.list());
  updateWorkspaceMutation = injectMutation(() => this.workspaceMutations.update());
  deleteWorkspaceMutation = injectMutation(() => this.workspaceMutations.delete());

  readonly isEditingName = signal<boolean>(false);
  readonly workspaceName = signal<string>('');
  readonly activeWorkspaceId = computed(() => this.routeWorkspaceId() ?? this.selectedWorkspaceId());
  readonly workspace = computed(() => {
    const id = this.activeWorkspaceId();
    return this.listWorkspacesQuery.data()?.find((workspace) => workspace.id === id) ?? null;
  });

  startRenameWorkspace(): void {
    const workspace = this.workspace();
    if (!workspace) {
      return;
    }

    this.workspaceName.set(workspace.name);
    this.isEditingName.set(true);
  }

  cancelRenameWorkspace(): void {
    this.isEditingName.set(false);
    this.workspaceName.set('');
  }

  onWorkspaceNameInput(event: Event): void {
    const target = event.target as HTMLInputElement;
    this.workspaceName.set(target.value);
  }

  confirmRenameWorkspace(): void {
    const workspace = this.workspace();
    const name = this.workspaceName().trim();
    if (!workspace || !name) {
      this.cancelRenameWorkspace();
      return;
    }

    this.updateWorkspaceMutation.mutate({ ...workspace, name });
    this.cancelRenameWorkspace();
  }

  deleteWorkspace(): void {
    const workspace = this.workspace();
    if (!workspace) {
      return;
    }

    this.deleteWorkspaceMutation.mutate(workspace.id);

  }
}

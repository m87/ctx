import { Component, computed, effect, inject, signal } from '@angular/core';
import { Store } from '@ngxs/store';
import { injectMutation, injectQuery } from '@tanstack/angular-query-experimental';
import { WorkspaceMutations } from '../../api/workspace/workspace.mutations';
import { WorkspaceQueries } from '../../api/workspace/workspace.queries';
import { Workspace } from '../../api/workspace/workspace.service';
import { SearchDropdownSelectComponent } from '../shared/search-dropdown-select.component';
import { SearchSelectOption } from '../shared/search-select.component';
import { colorHash } from '../utils';
import { SelectWorkspace, WorkspaceState } from './workspace.state';

let nextWorkspaceSelectId = 0;

@Component({
  selector: 'ctx-sidebar-workspace-select',
  imports: [SearchDropdownSelectComponent],
  host: {
    class: 'block min-w-0',
  },
  template: `
    <div class="relative">
      <ctx-search-dropdown-select
        [inputId]="selectInputId"
        ariaLabel="Select workspace"
        placeholder="Select workspace…"
        searchPlaceholder="Search workspaces…"
        actionLabel="Add workspace"
        [emptyText]="workspaceEmptyText()"
        [options]="workspaceOptions()"
        [value]="activeWorkspaceId() ?? ''"
        [actionExpanded]="isAddingWorkspace()"
        [keepOpenOnAction]="true"
        [actionValue]="newWorkspaceName()"
        [actionPending]="createWorkspaceMutation.isPending()"
        actionInputPlaceholder="Workspace name"
        actionInputAriaLabel="New workspace name"
        actionConfirmAriaLabel="Create workspace"
        actionCancelAriaLabel="Cancel workspace creation"
        (selectionChange)="selectWorkspace($event)"
        (action)="startAddWorkspace()"
        (actionValueChange)="newWorkspaceName.set($event)"
        (actionConfirm)="confirmAddWorkspace()"
        (actionCancel)="cancelAddWorkspace()"
        (openChange)="onDropdownOpenChange($event)"
      ></ctx-search-dropdown-select>
    </div>
  `,
})
export class SidebarWorkspaceSelectComponent {
  readonly selectInputId = `workspace-select-${nextWorkspaceSelectId++}`;
  private readonly store = inject(Store);
  private readonly workspaceQueries = inject(WorkspaceQueries);
  private readonly workspaceMutations = inject(WorkspaceMutations);

  readonly listWorkspacesQuery = injectQuery(() => this.workspaceQueries.list());
  readonly createWorkspaceMutation = injectMutation(() => this.workspaceMutations.create());
  readonly isAddingWorkspace = signal(false);
  readonly newWorkspaceName = signal('');

  readonly activeWorkspaceId = this.store.selectSignal(WorkspaceState.selectedWorkspaceId);
  readonly workspaceStateInitialized = this.store.selectSignal(WorkspaceState.initialized);
  readonly workspaces = computed<readonly Workspace[]>(() => this.listWorkspacesQuery.data() ?? []);
  readonly workspaceOptions = computed<SearchSelectOption[]>(() =>
    this.workspaces().map((workspace) => ({
      value: workspace.id,
      label: workspace.name,
      description: workspace.description,
      color: colorHash(`workspace:${workspace.id}`),
    })),
  );
  readonly workspaceEmptyText = computed(() =>
    this.listWorkspacesQuery.isLoading() ? 'Loading workspaces…' : 'No matching workspaces',
  );

  private readonly selectFirstWorkspaceEffect = effect(() => {
    if (!this.workspaceStateInitialized() || this.listWorkspacesQuery.data() === undefined) {
      return;
    }

    const selectedWorkspaceId = this.activeWorkspaceId();
    const workspaces = this.workspaces();
    const firstWorkspace = workspaces[0];
    if (selectedWorkspaceId === null && firstWorkspace) {
      this.store.dispatch(new SelectWorkspace(firstWorkspace.id));
      return;
    }

    const selectedWorkspaceExists = workspaces.some(
      (workspace) => workspace.id === selectedWorkspaceId,
    );
    if (selectedWorkspaceId !== null && !selectedWorkspaceExists) {
      this.store.dispatch(new SelectWorkspace(firstWorkspace?.id ?? null));
    }
  });

  selectWorkspace(workspaceId: string): void {
    this.store.dispatch(new SelectWorkspace(workspaceId));
  }

  startAddWorkspace(): void {
    this.isAddingWorkspace.set(true);
    this.newWorkspaceName.set('');
  }

  onDropdownOpenChange(open: boolean): void {
    if (!open && this.isAddingWorkspace()) {
      this.cancelAddWorkspace();
    }
  }

  cancelAddWorkspace(): void {
    this.isAddingWorkspace.set(false);
    this.newWorkspaceName.set('');
  }

  confirmAddWorkspace(): void {
    const name = this.newWorkspaceName().trim();
    if (!name) {
      this.cancelAddWorkspace();
      return;
    }

    this.createWorkspaceMutation.mutate(name, {
      onSuccess: (workspace) => {
        this.store.dispatch(new SelectWorkspace(workspace.id));
        this.cancelAddWorkspace();
      },
    });
  }
}

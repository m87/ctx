import { Component, ElementRef, HostListener, computed, effect, inject, signal, viewChild } from '@angular/core';
import { NgIcon, provideIcons } from '@ng-icons/core';
import {
  lucideCheck,
  lucideChevronDown,
  lucideChevronUp,
  lucidePlus,
  lucideX,
} from '@ng-icons/lucide';
import { Store } from '@ngxs/store';
import { SelectWorkspace, WorkspaceState } from './workspace.state';
import { injectMutation, injectQuery } from '@tanstack/angular-query-experimental';
import { WorkspaceQueries } from '../../api/workspace.quries';
import { WorkspaceMutations } from '../../api/workspace.mutations';
import { Workspace } from '../../api/workspace.service';

@Component({
  selector: 'app-sidebar-workspace-select',
  imports: [NgIcon],
  providers: [provideIcons({ lucideCheck, lucideChevronDown, lucideChevronUp, lucidePlus, lucideX })],
  template: `
    <div class="relative">
      <button
        type="button"
        class="w-full px-2.5 py-1 flex items-center justify-between gap-2 rounded-md hover:bg-muted/50 transition-colors"
        (click)="toggleOpen()"
        aria-label="Select workspace"
      >
        <div class="min-w-0 text-left leading-tight">
          <div class="text-sm font-semibold truncate">{{ activeWorkspaceName() }}</div>
          <div class="text-[10px] uppercase tracking-[0.08em] text-muted-foreground/80">workspace</div>
        </div>
        <ng-icon [name]="isOpen() ? 'lucideChevronUp' : 'lucideChevronDown'" class="text-muted-foreground"></ng-icon>
      </button>

      @if (isOpen()) {
        <div
          class="absolute left-0 right-0 top-[calc(100%+0.35rem)] z-30 border rounded-md bg-popover text-popover-foreground shadow-sm p-2 flex flex-col gap-1.5 origin-top animate-in fade-in-0 zoom-in-95 slide-in-from-top-2 duration-200"
        >
          @if (isAddingWorkspace()) {
            <div class="p-1.5 rounded-lg bg-muted/40 border border-dashed flex items-center gap-1 min-w-0 overflow-hidden">
              <input
                type="text"
                class="h-8 min-w-0 flex-1 rounded-md border bg-background px-2 text-[13px] outline-none focus:ring-1 focus:ring-ring"
                placeholder="name"
                [value]="newWorkspaceName()"
                (input)="onNewWorkspaceInput($event)"
                (keydown.enter)="confirmAddWorkspace()"
                (keydown.escape)="cancelAddWorkspace()"
                #newWorkspaceInput
              />
              <button type="button" class="h-8 w-8 shrink-0 rounded-md hover:bg-background flex items-center justify-center" (click)="confirmAddWorkspace()">
                <ng-icon name="lucideCheck"></ng-icon>
              </button>
              <button type="button" class="h-8 w-8 shrink-0 rounded-md hover:bg-background flex items-center justify-center" (click)="cancelAddWorkspace()">
                <ng-icon name="lucideX"></ng-icon>
              </button>
            </div>
          } @else {
            <button
              type="button"
              class="h-8 px-2 rounded-lg border border-dashed text-[11px] uppercase tracking-[0.08em] text-muted-foreground hover:text-foreground hover:bg-muted/40 flex items-center justify-center gap-1.5"
              (click)="startAddWorkspace()"
            >
              <ng-icon name="lucidePlus"></ng-icon>
              add workspace
            </button>
          }

          <div class="max-h-64 overflow-y-auto pr-1">
            <button
              type="button"
              class="w-full rounded-md px-2.5 py-2 text-left hover:bg-muted/50 transition-colors"
              (click)="selectWorkspace(null)"
            >
              <div class="flex items-center gap-1.5">
                <span class="text-[13px] font-medium truncate">Default</span>
                @if (activeWorkspaceId() === null) {
                  <ng-icon name="lucideCheck" class="text-muted-foreground text-[12px]"></ng-icon>
                }
              </div>
              <div class="text-[11px] text-muted-foreground mt-0.5">show contexts without selecting a workspace</div>
            </button>

            @if (listWorkspacesQuery.isLoading()) {
              <div class="px-2.5 py-2 text-[12px] text-muted-foreground">Loading workspaces...</div>
            }

            @for (workspace of workspaces(); track workspace.id) {
              <div class="group rounded-md border border-transparent hover:bg-muted/50 transition-colors">
                <div class="px-2.5 py-2 flex items-start justify-between gap-1.5">
                  <button
                    type="button"
                    class="min-w-0 flex-1 text-left overflow-hidden"
                    (click)="selectWorkspace(workspace.id)"
                  >
                    <div class="flex items-center gap-1.5">
                      <span class="text-[13px] font-medium truncate">{{ workspace.name }}</span>
                      @if (activeWorkspaceId() === workspace.id) {
                        <ng-icon name="lucideCheck" class="text-muted-foreground text-[12px]"></ng-icon>
                      }
                    </div>
                    <div class="text-[11px] text-muted-foreground mt-0.5">workspace</div>
                  </button>
                </div>
              </div>
            } @empty {
              @if (!listWorkspacesQuery.isLoading()) {
                <div class="px-2.5 py-2 text-[12px] text-muted-foreground">No workspaces yet</div>
              }
            }
          </div>
        </div>
      }
    </div>
  `,
})
export class SidebarWorkspaceSelectComponent {
  private readonly host = inject(ElementRef<HTMLElement>);
  private readonly store = inject(Store);
  private readonly workspaceQueries = inject(WorkspaceQueries);
  private readonly workspaceMutations = inject(WorkspaceMutations);
  private readonly newWorkspaceInput = viewChild<ElementRef<HTMLInputElement>>('newWorkspaceInput');

  listWorkspacesQuery = injectQuery(() => this.workspaceQueries.list());
  createWorkspaceMutation = injectMutation(() => this.workspaceMutations.create());

  readonly isOpen = signal<boolean>(false);
  readonly isAddingWorkspace = signal<boolean>(false);
  readonly newWorkspaceName = signal<string>('');

  readonly activeWorkspaceId = this.store.selectSignal(WorkspaceState.selectedWorkspaceId);
  readonly workspaces = computed<readonly Workspace[]>(() => this.listWorkspacesQuery.data() ?? []);

  readonly activeWorkspaceName = computed<string>(() => {
    const selected = this.workspaces().find((workspace) => workspace.id === this.activeWorkspaceId());
    return selected?.name ?? 'Default';
  });

  private readonly focusInputEffect = effect(() => {
    if (this.isAddingWorkspace()) {
      this.newWorkspaceInput()?.nativeElement.focus();
    }
  });

  toggleOpen(): void {
    this.isOpen.update((open) => !open);
  }

  @HostListener('document:mousedown', ['$event'])
  onDocumentMouseDown(event: MouseEvent): void {
    if (!this.isOpen()) {
      return;
    }

    const target = event.target as Node | null;
    if (target && !this.host.nativeElement.contains(target)) {
      this.isOpen.set(false);
    }
  }

  selectWorkspace(workspaceId: string | null): void {
    this.store.dispatch(new SelectWorkspace(workspaceId));
    this.isOpen.set(false);
  }

  startAddWorkspace(): void {
    this.isAddingWorkspace.set(true);
    this.newWorkspaceName.set('');
  }

  cancelAddWorkspace(): void {
    this.isAddingWorkspace.set(false);
    this.newWorkspaceName.set('');
  }

  onNewWorkspaceInput(event: Event): void {
    const target = event.target as HTMLInputElement;
    this.newWorkspaceName.set(target.value);
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
      },
    });
    this.cancelAddWorkspace();
  }
}

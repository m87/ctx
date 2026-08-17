import { Component, ElementRef, computed, effect, inject, signal, viewChild } from '@angular/core';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucidePause, lucidePlay, lucidePlus } from '@ng-icons/lucide';
import { ContextMutations } from '../../api/context/context.mutations';
import { ContextQueries } from '../../api/context/context.queries';
import { injectMutation, injectQuery } from '@tanstack/angular-query-experimental';
import { Context } from '../../api/context/context.service';
import { ProjectQueries } from '../../api/project/project.queries';
import { RouterLink } from '@angular/router';
import { Store } from '@ngxs/store';
import { LinkifiedTextComponent } from '../shared/linkified-text.component';
import { WorkspaceState } from './workspace.state';
import { HlmSkeletonImports } from '@spartan-ng/helm/skeleton';

@Component({
  selector: 'ctx-sidebar-context-list',
  imports: [NgIcon, RouterLink, LinkifiedTextComponent, HlmSkeletonImports],
  providers: [provideIcons({ lucidePause, lucidePlay, lucidePlus })],
  template: ` <div class="group/list flex flex-col gap-1 p-2">
    @if (isAddingContext()) {
      <input
        #newContextInput
        type="text"
        class="h-7 px-2 rounded-md border bg-background text-[13px] outline-none focus:ring-1 focus:ring-ring"
        placeholder="Context name"
        [value]="newContextName()"
        (input)="onNewContextNameInput($event)"
        (keydown.enter)="confirmAddContext()"
        (keydown.escape)="cancelAddContext()"
      />
    } @else {
      <button
        class="h-7 px-2 rounded-md border border-dashed text-[11px] uppercase tracking-[0.08em] text-muted-foreground hover:text-foreground hover:bg-muted/40 flex items-center justify-center gap-1.5"
        (click)="startAddContext()"
        aria-label="Add new context"
      >
        <ng-icon name="lucidePlus"></ng-icon>
        add context
      </button>
    }

    @if (listContextsQuery.isLoading()) {
      <div class="flex flex-col gap-2 px-2 py-1" role="status">
        <span class="sr-only">Loading contexts</span>
        @for (item of contextSkeletonItems; track item) {
          <div class="h-6 flex items-center gap-2">
            <hlm-skeleton class="h-3 w-full" [class.max-w-36]="item % 2 === 0"></hlm-skeleton>
            <hlm-skeleton class="size-5 ml-auto shrink-0"></hlm-skeleton>
          </div>
        }
      </div>
    } @else {
      @for (context of contexts(); track context.id) {
        <div
          class="group flex items-center gap-2 text-[13px] pl-2 pr-1 py-1 font-medium hover:bg-muted/60 rounded-md cursor-pointer"
          [routerLink]="['/context', context.id]"
          role="link"
          tabindex="0"
        >
          <span class="min-w-0 flex-1 truncate">
            <ctx-linkified-text [text]="context.name" />
          </span>
          <button
            type="button"
            class="h-6 w-6 shrink-0 rounded flex items-center justify-center text-muted-foreground opacity-60 hover:opacity-100 hover:bg-muted md:opacity-0 md:group-hover:opacity-100 focus:opacity-100"
            [disabled]="switchContextMutation.isPending() || freeContextMutation.isPending()"
            [attr.aria-label]="(isActiveContext(context.id) ? 'Pause ' : 'Start ') + context.name"
            [title]="(isActiveContext(context.id) ? 'Pause ' : 'Start ') + context.name"
            (click)="toggleContext($event, context)"
          >
            <ng-icon
              [name]="isActiveContext(context.id) ? 'lucidePause' : 'lucidePlay'"
              class="text-[12px] pointer-events-none"
            ></ng-icon>
          </button>
        </div>
      }
    }
  </div>`,
})
export class SidebarContextListComponent {
  readonly contextSkeletonItems = [0, 1, 2, 3, 4];
  private contextQueries = inject(ContextQueries);
  private contextMutations = inject(ContextMutations);
  private projectQueries = inject(ProjectQueries);
  private store = inject(Store);
  private readonly newContextInput = viewChild<ElementRef<HTMLInputElement>>('newContextInput');

  readonly selectedWorkspaceId = this.store.selectSignal(WorkspaceState.selectedWorkspaceId);
  readonly selectedProject = this.store.selectSignal(WorkspaceState.selectedProjectId);
  readonly selectedProjectId = computed(() => this.selectedProject() ?? '');

  listContextsQuery = injectQuery(() => this.contextQueries.list(this.selectedWorkspaceId()));
  selectedProjectQuery = injectQuery(() => this.projectQueries.get(this.selectedProjectId()));
  activeContextQuery = injectQuery(() => this.contextQueries.active());
  createContextMutation = injectMutation(() => this.contextMutations.create());
  switchContextMutation = injectMutation(() => this.contextMutations.switch());
  freeContextMutation = injectMutation(() => this.contextMutations.free());

  readonly contexts = computed<readonly Context[]>(() =>
    this.filterBySelectedProject(this.listContextsQuery.data() ?? []),
  );
  readonly isAddingContext = signal<boolean>(false);
  readonly newContextName = signal<string>('');

  private readonly focusInputEffect = effect(() => {
    if (this.isAddingContext()) {
      this.newContextInput()?.nativeElement.focus();
    }
  });

  filterBySelectedProject(contexts: readonly Context[]): readonly Context[] {
    const selectedProjectId = this.selectedProject();
    if (!selectedProjectId) {
      return contexts.filter((context) => !context.project || !context.project.id);
    }
    return contexts.filter(
      (context) => context.project && context.project.id === selectedProjectId,
    );
  }

  startAddContext(): void {
    this.isAddingContext.set(true);
    this.newContextName.set('');
  }

  cancelAddContext(): void {
    this.isAddingContext.set(false);
    this.newContextName.set('');
  }

  onNewContextNameInput(event: Event): void {
    const target = event.target as HTMLInputElement;
    this.newContextName.set(target.value);
  }

  confirmAddContext(): void {
    const name = this.newContextName().trim();
    if (!name) {
      this.cancelAddContext();
      return;
    }

    const projectId = this.selectedProjectId();
    this.createContextMutation.mutate({
      name,
      workspaceId: this.selectedWorkspaceId() ?? '',
      project: projectId
        ? {
            id: projectId,
            name: this.selectedProjectQuery.data()?.name ?? '',
          }
        : undefined,
    });
    this.cancelAddContext();
  }

  isActiveContext(contextId: string): boolean {
    return this.activeContextQuery.data()?.id === contextId;
  }

  toggleContext(event: MouseEvent, context: Context): void {
    event.preventDefault();
    event.stopPropagation();
    if (this.isActiveContext(context.id)) {
      this.freeContextMutation.mutate();
      return;
    }
    this.switchContextMutation.mutate(context);
  }
}

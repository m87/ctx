import { Component, ElementRef, computed, effect, inject, signal, viewChild } from '@angular/core';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucidePlay, lucidePlus } from '@ng-icons/lucide';
import { ContextMutations } from '../../api/context.mutations';
import { ContextQueries } from '../../api/context.quries';
import { injectMutation, injectQuery } from '@tanstack/angular-query-experimental';
import { Context } from '../../api/context.service';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-sidebar-context-list',
  imports: [NgIcon, RouterLink],
  providers: [provideIcons({ lucidePlay, lucidePlus })],
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

    @for (context of contexts(); track context.id) {
      <div
        class="group flex justify-between items-center text-[13px] px-2 py-1.5 font-medium hover:bg-muted/60 rounded-md cursor-pointer"
      >
        <span [routerLink]="['/context', context.id]" class="truncate">{{ context.name }}</span>
        <span class="relative h-4 text-muted-foreground text-[13px] flex items-center justify-end">
          <ng-icon
            name="lucidePlay"
            class="absolute inset-0 flex items-center justify-end opacity-0 group-hover:opacity-100"
            (click)="switchContextMutation.mutate(context)"
          ></ng-icon>
        </span>
      </div>
    }
  </div>`,
})
export class SidebarContextListComponent {
  private contextQueries = inject(ContextQueries);
  private contextMutations = inject(ContextMutations);
  private readonly newContextInput = viewChild<ElementRef<HTMLInputElement>>('newContextInput');

  listContextsQuery = injectQuery(() => this.contextQueries.list());
  createContextMutation = injectMutation(() => this.contextMutations.create());
  switchContextMutation = injectMutation(() => this.contextMutations.switch());

  readonly contexts = computed<readonly Context[]>(() => this.listContextsQuery.data() ?? []);
  readonly isAddingContext = signal<boolean>(false);
  readonly newContextName = signal<string>('');

  private readonly focusInputEffect = effect(() => {
    if (this.isAddingContext()) {
      this.newContextInput()?.nativeElement.focus();
    }
  });

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

    this.createContextMutation.mutate({ name } as Context);
    this.cancelAddContext();
  }
}

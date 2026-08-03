import { Component, inject, input } from '@angular/core';
import { injectMutation, injectQuery } from '@tanstack/angular-query-experimental';
import { ContextMutations } from '../../api/context.mutations';
import { ContextQueries } from '../../api/context.queries';
import { ContextListGroup, ContextListGroupComponent } from './context-list-group.component';
import { ContextListItem, ContextListItemComponent } from './context-list-item.component';

@Component({
  selector: 'ctx-context-list',
  imports: [ContextListGroupComponent, ContextListItemComponent],
  template: `
    @if (items().length > 0 || group()) {
      <div class="flex flex-col gap-2">
        @for (item of items(); track item.id) {
          <ctx-context-list-item
            [item]="item"
            [active]="activeContextQuery.data()?.id === item.id"
            [startPending]="switchContextMutation.isPending() || freeContextMutation.isPending()"
            (start)="startContext($event)"
          ></ctx-context-list-item>
        }
        @if (group(); as groupedContexts) {
          <ctx-context-list-group
            [group]="groupedContexts"
            [activeContextId]="activeContextQuery.data()?.id ?? null"
            [startPending]="switchContextMutation.isPending() || freeContextMutation.isPending()"
            (start)="startContext($event)"
          ></ctx-context-list-group>
        }
      </div>
    } @else if (emptyMessage()) {
      <p class="text-xs text-muted-foreground">{{ emptyMessage() }}</p>
    }
  `,
})
export class ContextListComponent {
  private readonly contextMutations = inject(ContextMutations);
  private readonly contextQueries = inject(ContextQueries);

  readonly activeContextQuery = injectQuery(() => this.contextQueries.active());
  readonly switchContextMutation = injectMutation(() => this.contextMutations.switch());
  readonly freeContextMutation = injectMutation(() => this.contextMutations.free());
  readonly items = input<readonly ContextListItem[]>([]);
  readonly group = input<ContextListGroup | null>(null);
  readonly emptyMessage = input('');

  startContext(item: ContextListItem): void {
    if (item.archived) {
      return;
    }
    if (this.activeContextQuery.data()?.id === item.id) {
      this.freeContextMutation.mutate();
      return;
    }
    this.switchContextMutation.mutate({ id: item.id, name: item.name });
  }
}

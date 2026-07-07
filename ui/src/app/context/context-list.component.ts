import { Component, input } from '@angular/core';
import { ContextListGroup, ContextListGroupComponent } from './context-list-group.component';
import { ContextListItem, ContextListItemComponent } from './context-list-item.component';

@Component({
  selector: 'ctx-context-list',
  imports: [ContextListGroupComponent, ContextListItemComponent],
  template: `
    @if (items().length > 0 || group()) {
      <div class="flex flex-col gap-2">
        @for (item of items(); track item.id) {
          <ctx-context-list-item [item]="item"></ctx-context-list-item>
        }
        @if (group(); as groupedContexts) {
          <ctx-context-list-group [group]="groupedContexts"></ctx-context-list-group>
        }
      </div>
    } @else if (emptyMessage()) {
      <p class="text-xs text-muted-foreground">{{ emptyMessage() }}</p>
    }
  `,
})
export class ContextListComponent {
  readonly items = input<readonly ContextListItem[]>([]);
  readonly group = input<ContextListGroup | null>(null);
  readonly emptyMessage = input('');
}

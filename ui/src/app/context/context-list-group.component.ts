import { Component, input, signal } from '@angular/core';
import { ContextListGroupItemComponent } from './context-list-group-item.component';
import { ContextListItem } from './context-list-item.component';

export interface ContextListGroup {
  id: string;
  name: string;
  duration: string;
  percentage: number;
  color: string;
  sessions: number;
  groupedCount: number;
  items: readonly ContextListItem[];
}

@Component({
  selector: 'ctx-context-list-group',
  imports: [ContextListGroupItemComponent],
  template: `
    <div class="flex flex-col gap-2">
      <div
        class="rounded-lg border bg-card p-3 hover:bg-muted/30 transition-colors cursor-pointer"
        role="button"
        tabindex="0"
        [attr.aria-expanded]="expanded()"
        (click)="toggle()"
        (keydown.enter)="toggle()"
        (keydown.space)="$event.preventDefault(); toggle()"
      >
        <div class="flex items-center gap-2 mb-2">
          <span class="w-2 h-2 rounded-sm shrink-0" [style.background-color]="group().color"></span>
          <span class="text-sm font-medium flex-1 truncate">{{ group().name }}</span>
          <span class="text-xs text-muted-foreground">{{ group().duration }}</span>
        </div>
        <div class="h-1.5 rounded bg-muted/40 overflow-hidden">
          <div
            class="h-full rounded"
            [style.width.%]="boundedPercentage(group().percentage)"
            [style.background-color]="group().color"
          ></div>
        </div>
        <div class="mt-2 text-[10px] text-muted-foreground">
          {{ group().sessions }} {{ group().sessions === 1 ? 'session' : 'sessions' }} ·
          {{ boundedPercentage(group().percentage).toFixed(1) }}%
        </div>
        <div class="mt-2 text-[10px] text-muted-foreground">
          {{
            expanded()
              ? 'Hide smaller contexts'
              : 'Show ' +
                group().groupedCount +
                ' smaller ' +
                (group().groupedCount === 1 ? 'context' : 'contexts')
          }}
        </div>
      </div>

      @if (expanded()) {
        @for (item of group().items; track item.id) {
          <ctx-context-list-group-item [item]="item"></ctx-context-list-group-item>
        }
      }
    </div>
  `,
})
export class ContextListGroupComponent {
  readonly group = input.required<ContextListGroup>();
  readonly expanded = signal(false);

  toggle(): void {
    this.expanded.update((expanded) => !expanded);
  }

  boundedPercentage(value: number): number {
    if (!Number.isFinite(value)) {
      return 0;
    }
    return Math.min(100, Math.max(0, value));
  }
}

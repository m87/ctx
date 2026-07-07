import { Component, input } from '@angular/core';
import { RouterLink } from '@angular/router';

export interface ContextListItem {
  id: string;
  name: string;
  duration: string;
  durationValue?: number;
  percentage: number;
  color: string;
  sessions?: number;
  sessionRanges?: readonly string[];
  distributionPercentage?: number;
  archived?: boolean;
}

@Component({
  selector: 'ctx-context-list-item',
  imports: [RouterLink],
  template: `
    <a
      class="block rounded-lg border bg-card p-3 hover:bg-muted/30 transition-colors"
      [routerLink]="['/context', item().id]"
    >
      <div class="flex items-center gap-2 mb-2">
        <span class="w-2 h-2 rounded-sm shrink-0" [style.background-color]="item().color"></span>
        <span class="text-sm font-medium flex-1 truncate">{{ item().name }}</span>
        @if (item().archived) {
          <span class="text-[10px] font-medium rounded border px-1.5 py-0.5 text-muted-foreground">
            Archived
          </span>
        }
        <span class="text-xs text-muted-foreground">{{ item().duration }}</span>
      </div>
      <div class="h-1.5 rounded bg-muted/40 overflow-hidden">
        <div
          class="h-full rounded"
          [style.width.%]="boundedPercentage(item().percentage)"
          [style.background-color]="item().color"
        ></div>
      </div>
      <div class="mt-2 flex flex-wrap gap-x-2 gap-y-1 text-[10px] text-muted-foreground">
        @if ((item().sessionRanges ?? []).length > 0) {
          @for (session of item().sessionRanges ?? []; track session) {
            <span>{{ session }}</span>
          }
        } @else if (item().sessions !== undefined) {
          <span>
            {{ item().sessions }}
            {{ item().sessions === 1 ? 'session' : 'sessions' }} ·
            {{ boundedPercentage(item().percentage).toFixed(1) }}%
          </span>
        } @else {
          <span>{{ boundedPercentage(item().percentage).toFixed(1) }}%</span>
        }
      </div>
    </a>
  `,
})
export class ContextListItemComponent {
  readonly item = input.required<ContextListItem>();

  boundedPercentage(value: number): number {
    if (!Number.isFinite(value)) {
      return 0;
    }
    return Math.min(100, Math.max(0, value));
  }
}

import { Component, input } from '@angular/core';

export interface DistributionItem {
  id: string;
  name: string;
  percentage: number;
  color: string;
  duration?: string;
}

@Component({
  selector: 'ctx-distribution',
  template: `
    <div>
      <div class="text-[11px] uppercase tracking-[0.08em] text-muted-foreground font-semibold mb-2">
        {{ label() }}
      </div>
      @if (items().length > 0) {
        <div class="flex h-2 rounded-md overflow-hidden gap-px bg-muted/40">
          @for (item of items(); track item.id) {
            <div
              [style.width.%]="boundedPercentage(item.percentage)"
              [style.background-color]="item.color"
              [title]="itemTitle(item)"
            ></div>
          }
        </div>
      } @else {
        <div class="h-2 rounded-md bg-muted/40"></div>
        @if (emptyMessage()) {
          <p class="mt-2 text-xs text-muted-foreground">{{ emptyMessage() }}</p>
        }
      }
    </div>
  `,
})
export class DistributionComponent {
  readonly label = input('Distribution');
  readonly emptyMessage = input('');
  readonly items = input<readonly DistributionItem[]>([]);

  boundedPercentage(value: number): number {
    if (!Number.isFinite(value)) {
      return 0;
    }
    return Math.min(100, Math.max(0, value));
  }

  itemTitle(item: DistributionItem): string {
    return item.duration ? `${item.name}: ${item.duration}` : item.name;
  }
}

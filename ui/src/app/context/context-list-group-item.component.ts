import { Component, input, output } from '@angular/core';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucidePause, lucidePlay } from '@ng-icons/lucide';
import { RouterLink } from '@angular/router';
import { LinkifiedTextComponent } from '../shared/linkified-text.component';
import { ContextListItem } from './context-list-item.component';

@Component({
  selector: 'ctx-context-list-group-item',
  imports: [NgIcon, RouterLink, LinkifiedTextComponent],
  providers: [provideIcons({ lucidePause, lucidePlay })],
  template: `
    <div
      class="block ml-4 cursor-pointer rounded-lg border border-dashed bg-card p-3 hover:bg-muted/30 transition-colors"
      [routerLink]="['/context', item().id]"
      role="link"
      tabindex="0"
    >
      <div class="flex items-center gap-2 mb-2">
        <span class="w-2 h-2 rounded-sm shrink-0" [style.background-color]="item().color"></span>
        <span class="text-sm font-medium flex-1 truncate">
          <ctx-linkified-text [text]="item().name" />
        </span>
        @if (item().archived) {
          <span class="text-[10px] font-medium rounded border px-1.5 py-0.5 text-muted-foreground">
            Archived
          </span>
        }
        <span class="text-xs text-muted-foreground">{{ item().duration }}</span>
        @if (!item().archived) {
          <button
            type="button"
            class="h-7 w-7 -my-1 -mr-1 shrink-0 rounded-md text-muted-foreground/70 hover:text-foreground hover:bg-muted transition-colors flex items-center justify-center"
            [disabled]="startPending()"
            [attr.aria-label]="(active() ? 'Pause ' : 'Start ') + item().name"
            [title]="(active() ? 'Pause ' : 'Start ') + item().name"
            (click)="requestStart($event)"
          >
            <ng-icon
              [name]="active() ? 'lucidePause' : 'lucidePlay'"
              class="text-[13px] pointer-events-none"
            ></ng-icon>
          </button>
        }
      </div>
      <div class="h-1.5 rounded bg-muted/40 overflow-hidden">
        <div
          class="h-full rounded"
          [style.width.%]="boundedPercentage(item().percentage)"
          [style.background-color]="item().color"
        ></div>
      </div>
      <div class="mt-2 text-[10px] text-muted-foreground">
        {{ item().sessions ?? 0 }} {{ item().sessions === 1 ? 'session' : 'sessions' }} ·
        {{ boundedPercentage(item().percentage).toFixed(1) }}%
      </div>
    </div>
  `,
})
export class ContextListGroupItemComponent {
  readonly item = input.required<ContextListItem>();
  readonly active = input(false);
  readonly startPending = input(false);
  readonly start = output<ContextListItem>();

  requestStart(event: MouseEvent): void {
    event.preventDefault();
    event.stopPropagation();
    this.start.emit(this.item());
  }

  boundedPercentage(value: number): number {
    if (!Number.isFinite(value)) {
      return 0;
    }
    return Math.min(100, Math.max(0, value));
  }
}

import { Component, input, output } from '@angular/core';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideArrowRightLeft, lucidePencil, lucideTrash2 } from '@ng-icons/lucide';
import { HlmButtonImports } from '@spartan-ng/helm/button';
import { Interval } from '../../api/interval.service';
import { durationAsH, durationAsM } from '../utils';

@Component({
  selector: 'ctx-context-interval-item',
  imports: [NgIcon, HlmButtonImports],
  providers: [
    provideIcons({
      lucidePencil,
      lucideTrash2,
      lucideArrowRightLeft,
    }),
  ],
  template: `
    <div class="w-full flex justify-start px-3 py-2.5 bg-card border border-border rounded-lg">
      @if (isEditing()) {
        <div class="w-full flex flex-col gap-2">
          <div class="w-full flex flex-col md:flex-row items-stretch md:items-end gap-2">
            <label class="flex-1 text-xs text-muted-foreground">
              Start
              <input
                type="datetime-local"
                class="w-full h-9 rounded-md border border-border bg-background px-3 text-sm mt-1"
                [value]="editStartInput()"
                (input)="editStartInputChange.emit(getInputValue($event))"
              />
            </label>
            <label class="flex-1 text-xs text-muted-foreground">
              End
              <input
                type="datetime-local"
                class="w-full h-9 rounded-md border border-border bg-background px-3 text-sm mt-1"
                [value]="editEndInput()"
                (input)="editEndInputChange.emit(getInputValue($event))"
              />
            </label>
            <button hlmBtn variant="outline" class="h-9 px-3 text-xs" (click)="cancel.emit()">
              Cancel
            </button>
            <button
              hlmBtn
              variant="outline"
              class="h-9 px-3 text-xs bg-blue-200/70 text-blue-600"
              [disabled]="updatePending()"
              (click)="save.emit(interval())"
            >
              Save
            </button>
          </div>
        </div>
      } @else {
        <div class="flex gap-2 ml-4 w-full items-center">
          <div class="flex flex-col flex-1">
            <div class="text-sm font-medium">
              {{ interval().start.toTimeString() }} - {{ interval().end.toTimeString() }}
            </div>
            <div class="text-xs text-muted-foreground">
              {{
                intervalStartDateEqEndDate(interval())
                  ? interval().start.toDateString()
                  : interval().start.toDateString() + ' - ' + interval().end.toDateString()
              }}
            </div>
          </div>
          <div class="text-xs text-muted-foreground">
            {{ parseDuration(interval().duration) }}
          </div>
          <div class="flex items-center gap-1.5">
            @if (!readonly()) {
              <button
                hlmBtn
                variant="outline"
                class="h-7 px-2 text-xs"
                (click)="edit.emit(interval())"
              >
                <ng-icon name="lucidePencil"></ng-icon>
              </button>
              <button
                hlmBtn
                variant="outline"
                class="h-7 px-2 text-xs"
                [disabled]="!canMove()"
                (click)="move.emit(interval())"
              >
                <ng-icon name="lucideArrowRightLeft"></ng-icon>
              </button>
              <button
                hlmBtn
                variant="outline"
                class="h-7 px-2 text-xs text-red-700 bg-red-100/60"
                [disabled]="deletePending()"
                (click)="delete.emit(interval())"
              >
                <ng-icon name="lucideTrash2"></ng-icon>
              </button>
            }
          </div>
        </div>
      }
    </div>
  `,
})
export class ContextIntervalItemComponent {
  readonly interval = input.required<Interval>();
  readonly isEditing = input(false);
  readonly editStartInput = input('');
  readonly editEndInput = input('');
  readonly updatePending = input(false);
  readonly deletePending = input(false);
  readonly canMove = input(false);
  readonly readonly = input(false);

  readonly editStartInputChange = output<string>();
  readonly editEndInputChange = output<string>();
  readonly edit = output<Interval>();
  readonly save = output<Interval>();
  readonly cancel = output<void>();
  readonly move = output<Interval>();
  readonly delete = output<Interval>();

  getInputValue(event: Event): string {
    return (event.target as HTMLInputElement | HTMLTextAreaElement).value;
  }

  intervalStartDateEqEndDate(interval: Interval): boolean {
    return interval.start.toDateString() === interval.end.toDateString();
  }

  parseDuration(duration: number | undefined): string {
    if (duration === undefined) {
      return '0h 0m';
    }
    return `${durationAsH(duration)}h ${durationAsM(duration)}m`;
  }
}

import { Component, computed, inject, input, signal } from '@angular/core';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucidePlus } from '@ng-icons/lucide';
import { HlmButtonImports } from '@spartan-ng/helm/button';
import { injectMutation, injectQuery } from '@tanstack/angular-query-experimental';
import { DateTime } from 'luxon';
import { ContextQueries } from '../../api/context.quries';
import { Context } from '../../api/context.service';
import { IntervalMutations } from '../../api/interval.mutations';
import { Interval, ZonedDateTime } from '../../api/interval.service';
import { ContextIntervalItemComponent } from './context-interval-item.component';

@Component({
  selector: 'ctx-context-interval-list',
  imports: [ContextIntervalItemComponent, NgIcon, HlmButtonImports],
  providers: [
    provideIcons({
      lucidePlus,
    }),
  ],
  template: `
    <div class="w-full flex flex-col gap-4 flex-1 min-h-0">
      <div
        class="w-full flex flex-wrap items-center justify-between gap-2 text-[11px] uppercase tracking-[0.08em] text-muted-foreground font-semibold"
      >
        <span>Intervals</span>
      </div>

      <div class="w-full rounded-lg border bg-card p-3 flex flex-col gap-2">
        <div class="text-[11px] uppercase tracking-[0.08em] text-muted-foreground font-semibold">
          Add interval
        </div>
        <div class="w-full flex flex-col md:flex-row items-stretch md:items-end gap-2">
          <label class="flex-1 text-xs text-muted-foreground">
            Start
            <input
              type="datetime-local"
              class="w-full h-9 rounded-md border border-border bg-background px-3 text-sm mt-1"
              [value]="newIntervalStartInput()"
              (input)="newIntervalStartInput.set(getInputValue($event))"
            />
          </label>
          <label class="flex-1 text-xs text-muted-foreground">
            End
            <input
              type="datetime-local"
              class="w-full h-9 rounded-md border border-border bg-background px-3 text-sm mt-1"
              [value]="newIntervalEndInput()"
              (input)="newIntervalEndInput.set(getInputValue($event))"
            />
          </label>
          <button
            hlmBtn
            variant="outline"
            class="h-9 px-3 text-xs bg-blue-200/70 text-blue-600"
            [disabled]="createIntervalMutation.isPending()"
            (click)="addInterval()"
          >
            <ng-icon name="lucidePlus"></ng-icon>
            <span>Add</span>
          </button>
        </div>
        @if (intervalFormError()) {
          <div class="text-xs text-red-600">{{ intervalFormError() }}</div>
        }
      </div>

      <div class="w-full flex flex-col gap-2 flex-1 min-h-0 overflow-auto pr-1 pb-2">
        @for (interval of intervals(); track interval.id) {
          <ctx-context-interval-item
            [interval]="interval"
            [isEditing]="editingIntervalId() === interval.id"
            [editStartInput]="editIntervalStartInput()"
            [editEndInput]="editIntervalEndInput()"
            [updatePending]="updateIntervalMutation.isPending()"
            [deletePending]="deleteIntervalMutation.isPending()"
            [canMove]="movableContexts().length > 0"
            (editStartInputChange)="editIntervalStartInput.set($event)"
            (editEndInputChange)="editIntervalEndInput.set($event)"
            (edit)="startIntervalEdit($event)"
            (save)="saveIntervalEdit($event)"
            (cancel)="cancelIntervalEdit()"
            (move)="openMoveDialog($event)"
            (delete)="deleteInterval($event)"
          ></ctx-context-interval-item>
        }
      </div>

      @if (moveDialogIntervalId()) {
        <div class="fixed inset-0 bg-black/30 flex items-center justify-center z-50 p-4">
          <div class="w-full max-w-md rounded-lg border bg-card p-4 flex flex-col gap-3">
            <div class="text-sm font-semibold">Move interval</div>
            <div class="text-xs text-muted-foreground">Select target context</div>
            <select
              class="w-full h-9 rounded-md border border-border bg-background px-3 text-sm"
              [value]="moveTargetContextId()"
              (change)="moveTargetContextId.set(getSelectValue($event))"
            >
              @for (context of movableContexts(); track context.id) {
                <option [value]="context.id">{{ context.name }}</option>
              }
            </select>
            <div class="flex justify-end gap-2 pt-1">
              <button hlmBtn variant="outline" class="h-9 px-3 text-xs" (click)="closeMoveDialog()">
                Cancel
              </button>
              <button
                hlmBtn
                variant="outline"
                class="h-9 px-3 text-xs bg-blue-200/70 text-blue-600"
                [disabled]="moveIntervalMutation.isPending()"
                (click)="confirmMoveInterval()"
              >
                Move
              </button>
            </div>
          </div>
        </div>
      }
    </div>
  `,
  styles: `
    :host {
      display: flex;
      width: 100%;
      flex: 1 1 auto;
      min-height: 0;
    }
  `,
})
export class ContextIntervalListComponent {
  private contextQueries = inject(ContextQueries);
  private intervalMutations = inject(IntervalMutations);

  readonly contextId = input.required<string>();
  readonly activeWorkspaceId = input<string | null>(null);
  readonly contexts = input<readonly Context[]>([]);

  createIntervalMutation = injectMutation(() => this.intervalMutations.create());
  updateIntervalMutation = injectMutation(() => this.intervalMutations.update());
  deleteIntervalMutation = injectMutation(() => this.intervalMutations.delete());
  moveIntervalMutation = injectMutation(() => this.intervalMutations.move());
  contextIntervalsQuery = injectQuery(() => this.contextQueries.intervals(this.contextId()));

  readonly intervals = computed(() => this.contextIntervalsQuery.data() ?? []);
  readonly movableContexts = computed(() =>
    this.contexts().filter((context) => context.id && context.id !== this.contextId()),
  );

  readonly newIntervalStartInput = signal('');
  readonly newIntervalEndInput = signal('');
  readonly editingIntervalId = signal<string | null>(null);
  readonly editIntervalStartInput = signal('');
  readonly editIntervalEndInput = signal('');
  readonly moveDialogIntervalId = signal<string | null>(null);
  readonly moveTargetContextId = signal('');
  readonly intervalFormError = signal('');

  constructor() {
    this.resetNewIntervalForm();
  }

  addInterval() {
    this.intervalFormError.set('');
    const parsed = this.parseIntervalInput(
      this.newIntervalStartInput(),
      this.newIntervalEndInput(),
    );

    if (!parsed) {
      return;
    }

    this.createIntervalMutation.mutate(
      {
        id: '',
        contextId: this.contextId(),
        start: parsed.start,
        end: parsed.end,
        duration: 0,
        workspaceId: this.activeWorkspaceId() ?? '',
      },
      {
        onSuccess: () => {
          this.resetNewIntervalForm();
        },
      },
    );
  }

  startIntervalEdit(interval: Interval) {
    this.intervalFormError.set('');
    this.editingIntervalId.set(interval.id);
    this.editIntervalStartInput.set(interval.start.toInputValue());
    this.editIntervalEndInput.set(interval.end.toInputValue());
  }

  cancelIntervalEdit() {
    this.editingIntervalId.set(null);
    this.editIntervalStartInput.set('');
    this.editIntervalEndInput.set('');
    this.intervalFormError.set('');
  }

  saveIntervalEdit(interval: Interval) {
    this.intervalFormError.set('');
    const parsed = this.parseIntervalInput(
      this.editIntervalStartInput(),
      this.editIntervalEndInput(),
    );

    if (!parsed) {
      return;
    }

    this.updateIntervalMutation.mutate(
      {
        id: interval.id,
        interval: {
          ...interval,
          contextId: this.contextId(),
          start: parsed.start,
          end: parsed.end,
        },
      },
      {
        onSuccess: () => {
          this.cancelIntervalEdit();
        },
      },
    );
  }

  deleteInterval(interval: Interval) {
    if (!window.confirm('Delete this interval?')) {
      return;
    }

    this.deleteIntervalMutation.mutate({ id: interval.id, contextId: this.contextId() });
  }

  openMoveDialog(interval: Interval) {
    const contexts = this.movableContexts();
    if (contexts.length === 0) {
      return;
    }

    this.moveDialogIntervalId.set(interval.id);
    this.moveTargetContextId.set(contexts[0]?.id ?? '');
  }

  closeMoveDialog() {
    this.moveDialogIntervalId.set(null);
    this.moveTargetContextId.set('');
  }

  confirmMoveInterval() {
    const intervalId = this.moveDialogIntervalId();
    const targetContextId = this.moveTargetContextId();

    if (!intervalId || !targetContextId) {
      return;
    }

    this.moveIntervalMutation.mutate(
      { id: intervalId, targetContextId },
      {
        onSuccess: () => {
          this.closeMoveDialog();
        },
      },
    );
  }

  getInputValue(event: Event): string {
    return (event.target as HTMLInputElement | HTMLTextAreaElement).value;
  }

  getSelectValue(event: Event): string {
    return (event.target as HTMLSelectElement).value;
  }

  private parseIntervalInput(
    startInput: string,
    endInput: string,
  ): { start: ZonedDateTime; end: ZonedDateTime } | null {
    const startDateTime = DateTime.fromFormat(startInput, "yyyy-MM-dd'T'HH:mm");
    const endDateTime = DateTime.fromFormat(endInput, "yyyy-MM-dd'T'HH:mm");

    if (!startDateTime.isValid || !endDateTime.isValid) {
      this.intervalFormError.set('Invalid start or end date/time.');
      return null;
    }

    if (endDateTime <= startDateTime) {
      this.intervalFormError.set('End must be later than start.');
      return null;
    }

    return {
      start: ZonedDateTime.fromDateTime(startDateTime),
      end: ZonedDateTime.fromDateTime(endDateTime),
    };
  }

  private resetNewIntervalForm() {
    const end = DateTime.local().startOf('minute');
    const start = end.minus({ minutes: 30 });
    this.newIntervalStartInput.set(start.toFormat("yyyy-MM-dd'T'HH:mm"));
    this.newIntervalEndInput.set(end.toFormat("yyyy-MM-dd'T'HH:mm"));
    this.intervalFormError.set('');
  }
}

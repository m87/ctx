import { Component, computed, inject, input, signal } from '@angular/core';
import { NgIcon, provideIcons } from '@ng-icons/core';
import {
  lucideArrowRightLeft,
  lucideClock3,
  lucidePlus,
  lucideScissors,
  lucideX,
} from '@ng-icons/lucide';
import { HlmButtonImports } from '@spartan-ng/helm/button';
import { injectMutation, injectQuery } from '@tanstack/angular-query-experimental';
import { ContextQueries } from '../../api/context/context.queries';
import { Context } from '../../api/context/context.service';
import { IntervalMutations } from '../../api/interval/interval.mutations';
import { Interval } from '../../api/interval/interval.service';
import { QueryErrorStateComponent } from '../shared/query-error-state.component';
import { SearchSelectComponent, SearchSelectOption } from '../shared/search-select.component';
import { TimeZoneService } from '../shared/time-zone.service';
import { ContextIntervalItemComponent } from './context-interval-item.component';
import { colorHash } from '../utils';
import { HlmSkeletonImports } from '@spartan-ng/helm/skeleton';
import { HlmSliderImports } from '@spartan-ng/helm/slider';
import { toast } from 'ngx-sonner';

interface SplitProperties {
  start: number;
  end: number;
  split: number;
  step: number;
}

@Component({
  selector: 'ctx-context-interval-list',
  imports: [
    ContextIntervalItemComponent,
    NgIcon,
    HlmButtonImports,
    QueryErrorStateComponent,
    HlmSkeletonImports,
    HlmSliderImports,
    SearchSelectComponent,
  ],
  providers: [
    provideIcons({
      lucideArrowRightLeft,
      lucideClock3,
      lucidePlus,
      lucideScissors,
      lucideX,
    }),
  ],
  template: `
    <div class="w-full flex flex-col gap-4 md:flex-1 md:min-h-0">
      <div
        class="w-full flex flex-wrap items-center justify-between gap-2 text-[11px] uppercase tracking-[0.08em] text-muted-foreground font-semibold"
      >
        <span>Intervals</span>
      </div>

      @if (showIntervalsError()) {
        <ctx-query-error-state
          class="flex-1 min-h-0"
          [error]="contextIntervalsQuery.error()"
          [paused]="contextIntervalsQuery.isPaused()"
          resourceName="intervals"
          [retrying]="contextIntervalsQuery.isFetching()"
          (retry)="retryIntervals()"
        ></ctx-query-error-state>
      } @else {
        @if (!readonly()) {
          <div class="w-full rounded-lg border bg-card p-3 flex flex-col gap-2">
            <div
              class="text-[11px] uppercase tracking-[0.08em] text-muted-foreground font-semibold"
            >
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
                class="h-9 px-3 text-xs"
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
        }

        <div class="w-full flex flex-col gap-2 md:flex-1 md:min-h-0 md:overflow-auto pr-1 pb-2">
          @if (contextIntervalsQuery.isLoading()) {
            <div class="flex flex-col gap-2" role="status">
              <span class="sr-only">Loading intervals</span>
              @for (item of intervalSkeletonItems; track item) {
                <div class="rounded-lg border bg-card p-3 flex items-center gap-3">
                  <hlm-skeleton class="h-3.5 w-28"></hlm-skeleton>
                  <hlm-skeleton class="h-3.5 w-36"></hlm-skeleton>
                  <hlm-skeleton class="h-7 w-16 ml-auto"></hlm-skeleton>
                </div>
              }
            </div>
          } @else {
            @for (interval of intervals(); track interval.id) {
              <ctx-context-interval-item
                [interval]="interval"
                [isEditing]="editingIntervalId() === interval.id"
                [editStartInput]="editIntervalStartInput()"
                [editEndInput]="editIntervalEndInput()"
                [updatePending]="updateIntervalMutation.isPending()"
                [deletePending]="deleteIntervalMutation.isPending()"
                [readonly]="readonly()"
                [canMove]="!readonly() && movableContexts().length > 0"
                (editStartInputChange)="editIntervalStartInput.set($event)"
                (editEndInputChange)="editIntervalEndInput.set($event)"
                (edit)="startIntervalEdit($event)"
                (save)="saveIntervalEdit($event)"
                (cancel)="cancelIntervalEdit()"
                (move)="openMoveDialog($event)"
                (delete)="deleteInterval($event)"
                (split)="openSplitDialog($event)"
              ></ctx-context-interval-item>
            }
          }
        </div>
      }

      @if (splitDialogInterval(); as interval) {
        <div
          class="fixed inset-0 z-50 flex items-end sm:items-center sm:justify-center"
          (click)="closeSplitDialog()"
        >
          <div class="absolute inset-0 bg-background/65 backdrop-blur-sm"></div>
          <section
            class="relative w-full sm:w-[min(92vw,34rem)] max-h-[92dvh] overflow-hidden rounded-t-3xl sm:rounded-2xl border border-border/80 bg-popover text-popover-foreground shadow-2xl"
            role="dialog"
            aria-modal="true"
            aria-labelledby="split-dialog-title"
            (click)="$event.stopPropagation()"
          >
            <header class="flex items-start gap-3 border-b border-border/70 px-5 py-4 sm:px-6">
              <div
                class="flex size-10 shrink-0 items-center justify-center rounded-xl bg-muted text-foreground"
              >
                <ng-icon name="lucideScissors" class="text-lg"></ng-icon>
              </div>
              <div class="min-w-0 flex-1 pt-0.5">
                <h2 id="split-dialog-title" class="text-base font-semibold tracking-tight">
                  Split interval
                </h2>
                <p class="mt-0.5 text-xs leading-5 text-muted-foreground">
                  Choose where this interval should become two separate entries.
                </p>
              </div>
              <button
                type="button"
                class="flex size-9 shrink-0 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:bg-muted hover:text-foreground focus-visible:ring-2 focus-visible:ring-ring focus-visible:outline-none"
                aria-label="Close split dialog"
                (click)="closeSplitDialog()"
              >
                <ng-icon name="lucideX" class="text-base"></ng-icon>
              </button>
            </header>

            <div
              class="flex max-h-[calc(92dvh-9rem)] flex-col gap-5 overflow-y-auto px-5 py-5 sm:px-6"
            >
              <div class="rounded-2xl border border-border/80 bg-muted/25 p-4">
                <div class="flex items-start justify-between gap-4 text-xs">
                  <div class="min-w-0">
                    <div class="font-medium text-foreground">
                      {{ timeZone.formatTime(interval.start) }}
                    </div>
                    <div class="mt-0.5 truncate text-muted-foreground">
                      {{ timeZone.formatDate(interval.start) }}
                    </div>
                  </div>
                  <div class="min-w-0 text-right">
                    <div class="font-medium text-foreground">
                      {{ timeZone.formatTime(interval.end) }}
                    </div>
                    <div class="mt-0.5 truncate text-muted-foreground">
                      {{ timeZone.formatDate(interval.end) }}
                    </div>
                  </div>
                </div>

                <hlm-slider
                  class="my-4 py-1 [&_[brnSliderTrack]]:h-2 [&_[brnSliderTrack]]:bg-muted [&_[brnSliderRange]]:bg-primary [&_[brnSliderThumb]]:size-5 [&_[brnSliderThumb]]:border-2 [&_[brnSliderThumb]]:border-primary"
                  [min]="splitProperties.start"
                  [max]="splitProperties.end"
                  [step]="splitProperties.step"
                  [value]="[splitProperties.split]"
                  (valueChange)="splitProperties.split = $event[0]"
                ></hlm-slider>

                <div class="flex justify-center">
                  <div
                    class="inline-flex items-center gap-2 rounded-full border border-border bg-background px-3 py-1.5 text-foreground shadow-xs"
                  >
                    <ng-icon name="lucideClock3" class="text-sm"></ng-icon>
                    <span class="text-xs font-semibold tabular-nums">
                      {{ timeZone.formatDateTime(splitAsInstant(splitProperties.split)) }}
                    </span>
                  </div>
                </div>
              </div>

              <div>
                <div class="mb-2.5 flex items-center justify-between gap-3">
                  <div
                    class="text-xs font-semibold uppercase tracking-[0.08em] text-muted-foreground"
                  >
                    Result preview
                  </div>
                  <div class="text-[11px] text-muted-foreground">2 intervals</div>
                </div>
                <div class="grid gap-2.5 sm:grid-cols-2">
                  <div class="rounded-xl border border-border/80 bg-card p-3.5 shadow-xs">
                    <div class="mb-3 flex items-center justify-between gap-2">
                      <span class="text-xs font-semibold text-foreground">First interval</span>
                      <span
                        class="rounded-md bg-muted px-2 py-0.5 text-[10px] font-semibold text-muted-foreground"
                      >
                        {{ formatSplitDuration(splitProperties.start, splitProperties.split) }}
                      </span>
                    </div>
                    <div class="flex items-center gap-2 text-sm font-medium tabular-nums">
                      <span>{{ timeZone.formatTime(interval.start) }}</span>
                      <span class="h-px flex-1 bg-border"></span>
                      <span>{{ timeZone.formatTime(splitAsInstant(splitProperties.split)) }}</span>
                    </div>
                  </div>

                  <div class="rounded-xl border border-border/80 bg-card p-3.5 shadow-xs">
                    <div class="mb-3 flex items-center justify-between gap-2">
                      <span class="text-xs font-semibold text-foreground">Second interval</span>
                      <span
                        class="rounded-md bg-muted px-2 py-0.5 text-[10px] font-semibold text-muted-foreground"
                      >
                        {{ formatSplitDuration(splitProperties.split, splitProperties.end) }}
                      </span>
                    </div>
                    <div class="flex items-center gap-2 text-sm font-medium tabular-nums">
                      <span>{{ timeZone.formatTime(splitAsInstant(splitProperties.split)) }}</span>
                      <span class="h-px flex-1 bg-border"></span>
                      <span>{{ timeZone.formatTime(interval.end) }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <footer
              class="flex flex-col-reverse gap-2 border-t border-border/70 bg-muted/20 px-5 py-4 sm:flex-row sm:justify-end sm:px-6"
            >
              <button hlmBtn variant="ghost" class="h-10 px-4 text-xs" (click)="closeSplitDialog()">
                Cancel
              </button>
              <button
                hlmBtn
                variant="outline"
                class="h-10 px-5 text-xs shadow-sm"
                [disabled]="splitIntervalMutation.isPending()"
                (click)="splitInterval()"
              >
                <ng-icon name="lucideScissors"></ng-icon>
                {{ splitIntervalMutation.isPending() ? 'Splitting…' : 'Split interval' }}
              </button>
            </footer>
          </section>
        </div>
      }

      @if (moveDialogInterval(); as interval) {
        <div
          class="fixed inset-0 z-50 flex items-end sm:items-center sm:justify-center"
          (click)="closeMoveDialog()"
        >
          <div class="absolute inset-0 bg-background/65 backdrop-blur-sm"></div>
          <section
            class="relative w-full sm:w-[min(92vw,30rem)] max-h-[92dvh] overflow-hidden rounded-t-3xl sm:rounded-2xl border border-border/80 bg-popover text-popover-foreground shadow-2xl"
            role="dialog"
            aria-modal="true"
            aria-labelledby="move-dialog-title"
            (click)="$event.stopPropagation()"
          >
            <header class="flex items-start gap-3 border-b border-border/70 px-5 py-4 sm:px-6">
              <div
                class="flex size-10 shrink-0 items-center justify-center rounded-xl bg-muted text-foreground"
              >
                <ng-icon name="lucideArrowRightLeft" class="text-lg"></ng-icon>
              </div>
              <div class="min-w-0 flex-1 pt-0.5">
                <h2 id="move-dialog-title" class="text-base font-semibold tracking-tight">
                  Move interval
                </h2>
                <p class="mt-0.5 text-xs leading-5 text-muted-foreground">
                  {{ timeZone.formatDate(interval.start) }} ·
                  {{ timeZone.formatTime(interval.start) }}–{{ timeZone.formatTime(interval.end) }}
                </p>
              </div>
              <button
                type="button"
                class="flex size-9 shrink-0 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:bg-muted hover:text-foreground focus-visible:ring-2 focus-visible:ring-ring focus-visible:outline-none"
                aria-label="Close move dialog"
                (click)="closeMoveDialog()"
              >
                <ng-icon name="lucideX" class="text-base"></ng-icon>
              </button>
            </header>

            <div class="flex flex-col gap-3 px-5 py-5 sm:px-6">
              <div>
                <label
                  for="move-context-search"
                  class="mb-2 block text-xs font-medium text-foreground"
                >
                  Destination context
                </label>
                <ctx-search-select
                  inputId="move-context-search"
                  ariaLabel="Destination context"
                  searchPlaceholder="Search contexts…"
                  emptyText="No matching contexts"
                  [options]="moveContextOptions()"
                  [value]="moveTargetContextId()"
                  (selectionChange)="selectMoveContext($event)"
                ></ctx-search-select>
              </div>

              <p class="text-xs leading-5 text-muted-foreground">
                The interval keeps its start, end and duration. Only its context changes.
              </p>
            </div>

            <footer
              class="flex flex-col-reverse gap-2 border-t border-border/70 bg-muted/20 px-5 py-4 sm:flex-row sm:justify-end sm:px-6"
            >
              <button hlmBtn variant="ghost" class="h-10 px-4 text-xs" (click)="closeMoveDialog()">
                Cancel
              </button>
              <button
                hlmBtn
                variant="outline"
                class="h-10 px-5 text-xs shadow-sm"
                [disabled]="!moveTargetContextId() || moveIntervalMutation.isPending()"
                (click)="confirmMoveInterval()"
              >
                <ng-icon name="lucideArrowRightLeft"></ng-icon>
                {{ moveIntervalMutation.isPending() ? 'Moving…' : 'Move interval' }}
              </button>
            </footer>
          </section>
        </div>
      }
    </div>
  `,
  styles: `
    :host {
      display: flex;
      width: 100%;
      flex: 0 0 auto;
    }

    @media (min-width: 48rem) {
      :host {
        flex: 1 1 auto;
        min-height: 0;
      }
    }
  `,
})
export class ContextIntervalListComponent {
  readonly intervalSkeletonItems = [0, 1, 2];
  private contextQueries = inject(ContextQueries);
  private intervalMutations = inject(IntervalMutations);
  timeZone = inject(TimeZoneService);

  readonly contextId = input.required<string>();
  readonly activeWorkspaceId = input<string | null>(null);
  readonly contexts = input<readonly Context[]>([]);
  readonly readonly = input(false);
  splitProperties: SplitProperties = {
    start: 0,
    end: 100,
    split: 50,
    step: 1,
  };

  createIntervalMutation = injectMutation(() => this.intervalMutations.create());
  updateIntervalMutation = injectMutation(() => this.intervalMutations.update());
  deleteIntervalMutation = injectMutation(() => this.intervalMutations.delete());
  splitIntervalMutation = injectMutation(() => this.intervalMutations.split());
  undoSplitMutation = injectMutation(() => this.intervalMutations.undoSplit());
  moveIntervalMutation = injectMutation(() => this.intervalMutations.move());
  contextIntervalsQuery = injectQuery(() => this.contextQueries.intervals(this.contextId()));

  readonly intervals = computed(() => this.contextIntervalsQuery.data() ?? []);
  readonly showIntervalsError = computed(
    () =>
      this.contextIntervalsQuery.data() === undefined &&
      (this.contextIntervalsQuery.isError() || this.contextIntervalsQuery.isPaused()),
  );
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
  readonly splitDialogIntervalId = signal<string | null>(null);
  readonly splitDialogInterval = computed(() => {
    const intervalId = this.splitDialogIntervalId();
    return intervalId ? this.intervals().find((interval) => interval.id === intervalId) : undefined;
  });
  readonly moveDialogInterval = computed(() => {
    const intervalId = this.moveDialogIntervalId();
    return intervalId ? this.intervals().find((interval) => interval.id === intervalId) : undefined;
  });
  readonly moveContextOptions = computed<SearchSelectOption[]>(() =>
    this.movableContexts().map((context) => ({
      value: context.id,
      label: context.name,
      color: colorHash(context.id),
      description: context.project?.name,
    })),
  );

  constructor() {
    this.resetNewIntervalForm();
  }

  splitInterval() {
    if (this.readonly()) {
      return;
    }

    const interval = this.intervals().find((i) => i.id === this.splitDialogIntervalId());
    if (!interval) {
      return;
    }

    const splitTime = this.splitAsDateTime(this.splitProperties.split);

    this.splitIntervalMutation.mutate(
      {
        id: interval.id,
        splitTime: `${splitTime}:00`,
        timeZone: this.timeZone.effectiveTimeZone(),
      },
      {
        onSuccess: (result) => {
          this.closeSplitDialog();
          toast('Interval split', {
            description: 'Two intervals were created.',
            duration: 10_000,
            classes: {
              toast: 'split-undo-toast',
              actionButton: 'split-undo-action',
            },
            action: {
              label: 'Undo',
              onClick: () => {
                this.undoSplitMutation.mutate(result, {
                  onSuccess: () => toast.success('Split undone'),
                });
              },
            },
          });
        },
      },
    );
  }

  splitAsDateTime(split: number): string {
    return this.timeZone.toInputValue(this.splitAsInstant(split));
  }

  splitAsInstant(split: number): string {
    return new Date(split).toISOString();
  }

  formatSplitDuration(start: number, end: number): string {
    const totalMinutes = Math.max(0, Math.round((end - start) / 60_000));
    const hours = Math.floor(totalMinutes / 60);
    const minutes = totalMinutes % 60;

    if (hours === 0) {
      return `${minutes}m`;
    }
    return minutes === 0 ? `${hours}h` : `${hours}h ${minutes}m`;
  }

  calcSplitProperties(interval: Interval): SplitProperties {
    const start = Date.parse(interval.start!);
    const end = Date.parse(interval.end!);
    const split = start + (end - start) / 2;
    const step = 60 * 1000;
    return {
      start,
      end,
      split,
      step,
    };
  }

  openSplitDialog(interval: Interval) {
    if (this.readonly()) {
      return;
    }

    this.splitProperties = this.calcSplitProperties(interval);
    this.splitDialogIntervalId.set(interval.id);
  }

  closeSplitDialog() {
    this.splitDialogIntervalId.set(null);
  }

  addInterval() {
    if (this.readonly()) {
      return;
    }

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
    if (this.readonly()) {
      return;
    }

    this.intervalFormError.set('');
    this.editingIntervalId.set(interval.id);
    this.editIntervalStartInput.set(this.timeZone.toInputValue(interval.start));
    this.editIntervalEndInput.set(this.timeZone.toInputValue(interval.end));
  }

  cancelIntervalEdit() {
    this.editingIntervalId.set(null);
    this.editIntervalStartInput.set('');
    this.editIntervalEndInput.set('');
    this.intervalFormError.set('');
  }

  saveIntervalEdit(interval: Interval) {
    if (this.readonly()) {
      return;
    }

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
    if (this.readonly()) {
      return;
    }

    if (!window.confirm('Delete this interval?')) {
      return;
    }

    this.deleteIntervalMutation.mutate({ id: interval.id, contextId: this.contextId() });
  }

  openMoveDialog(interval: Interval) {
    if (this.readonly()) {
      return;
    }

    const contexts = this.movableContexts();
    if (contexts.length === 0) {
      return;
    }

    this.moveDialogIntervalId.set(interval.id);
    this.moveTargetContextId.set('');
  }

  closeMoveDialog() {
    this.moveDialogIntervalId.set(null);
    this.moveTargetContextId.set('');
  }

  selectMoveContext(contextId: string) {
    this.moveTargetContextId.set(contextId);
  }

  confirmMoveInterval() {
    if (this.readonly()) {
      return;
    }

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

  retryIntervals(): void {
    void this.contextIntervalsQuery.refetch();
  }

  private parseIntervalInput(
    startInput: string,
    endInput: string,
  ): { start: string; end: string } | null {
    const start = this.timeZone.inputToUTC(startInput);
    const end = this.timeZone.inputToUTC(endInput);

    if (!start || !end) {
      this.intervalFormError.set('Invalid start or end date/time.');
      return null;
    }

    if (Date.parse(end) <= Date.parse(start)) {
      this.intervalFormError.set('End must be later than start.');
      return null;
    }

    return {
      start,
      end,
    };
  }

  private resetNewIntervalForm() {
    const end = this.timeZone.now().startOf('minute');
    const start = end.minus({ minutes: 30 });
    this.newIntervalStartInput.set(start.toFormat("yyyy-MM-dd'T'HH:mm"));
    this.newIntervalEndInput.set(end.toFormat("yyyy-MM-dd'T'HH:mm"));
    this.intervalFormError.set('');
  }
}

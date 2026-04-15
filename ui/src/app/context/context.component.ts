import { Component, computed, effect, inject, signal } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { NgIcon, provideIcons } from '@ng-icons/core';
import {
  lucideArrowRightLeft,
  lucideDot,
  lucidePencil,
  lucidePlay,
  lucidePlus,
  lucideTrash2,
} from '@ng-icons/lucide';
import { HlmButtonImports } from '@spartan-ng/helm/button';
import { HlmCardImports } from '@spartan-ng/helm/card';
import { map } from 'rxjs';
import { toSignal } from '@angular/core/rxjs-interop';
import { ContextQueries } from '../../api/context.quries';
import { injectMutation, injectQuery } from '@tanstack/angular-query-experimental';
import { ContextMutations } from '../../api/context.mutations';
import { Context, EMPTY_CONTEXT } from '../../api/context.service';
import { Interval, ZonedDateTime } from '../../api/interval.service';
import { durationAsH, durationAsM } from '../utils';
import { DateTime } from 'luxon';
import { IntervalMutations } from '../../api/interval.mutations';

@Component({
  imports: [NgIcon, HlmButtonImports, HlmCardImports],
  providers: [
    provideIcons({
      lucidePlay,
      lucideDot,
      lucidePlus,
      lucideTrash2,
      lucidePencil,
      lucideArrowRightLeft,
    }),
  ],
  selector: 'app-context',
  template: `
    <div
      class="w-full h-full overflow-hidden flex flex-col items-start justify-start p-4 md:p-6 gap-5 relative"
    >
      <div
        class="w-full flex flex-col justify-between items-start gap-4"
        [class.md:flex-row]="!isEditing()"
        [class.md:items-center]="!isEditing()"
      >
        <div class="flex flex-col gap-1 w-full min-w-0">
          <div
            class="text-[11px] font-semibold text-muted-foreground flex items-center gap-2 uppercase"
          >
            <span class="bg-amber-600 w-2 h-2 rounded-full"></span><span>CONTEXT</span>
          </div>
          @if (!isEditing()) {
            <div class="text-2xl font-semibold tracking-tight">{{ context().name }}</div>
            <div class="text-sm text-muted-foreground/90">{{ context().description }}</div>
            <div class="flex flex-wrap items-center gap-1.5 md:gap-2 mt-2">
              @for (tag of contextTags(); track tag) {
                <span
                  class="text-[10px] md:text-[11px] font-medium text-blue-600 bg-blue-50/80 px-2 py-1 rounded-md"
                >
                  #{{ tag }}
                </span>
              }
            </div>
          } @else {
            <div class="w-full mt-2 rounded-lg border bg-card p-3 md:p-4">
              <div class="grid gap-3">
                <label class="flex flex-col gap-1">
                  <span
                    class="text-[11px] uppercase tracking-[0.08em] text-muted-foreground font-semibold"
                    >Name</span
                  >
                  <input
                    class="w-full h-9 rounded-md border border-border bg-background px-3 text-sm transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring/50"
                    [value]="editName()"
                    (input)="editName.set(getInputValue($event))"
                    placeholder="Context name"
                  />
                </label>

                <label class="flex flex-col gap-1">
                  <span
                    class="text-[11px] uppercase tracking-[0.08em] text-muted-foreground font-semibold"
                    >Description</span
                  >
                  <textarea
                    class="w-full min-h-24 rounded-md border border-border bg-background px-3 py-2 text-sm transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring/50"
                    [value]="editDescription()"
                    (input)="editDescription.set(getInputValue($event))"
                    placeholder="What this context is for"
                  ></textarea>
                </label>

                <label class="flex flex-col gap-1">
                  <span
                    class="text-[11px] uppercase tracking-[0.08em] text-muted-foreground font-semibold"
                    >Tags</span
                  >
                  <input
                    class="w-full h-9 rounded-md border border-border bg-background px-3 text-sm transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring/50"
                    [value]="editTagsInput()"
                    (input)="editTagsInput.set(getInputValue($event))"
                    placeholder="Comma separated"
                  />
                </label>

                <div class="flex flex-wrap items-center gap-1.5 min-h-5">
                  @if (editTagsPreview().length > 0) {
                    @for (tag of editTagsPreview(); track tag) {
                      <span
                        class="text-[10px] md:text-[11px] font-medium text-primary bg-primary/10 px-2 py-1 rounded-md"
                      >
                        #{{ tag }}
                      </span>
                    }
                  } @else {
                    <span class="text-xs text-muted-foreground">No tags yet</span>
                  }
                </div>
              </div>
            </div>
          }
        </div>
        <div
          class="flex items-center gap-2 w-full"
          [class.md:w-auto]="!isEditing()"
          [class.flex-wrap]="isEditing()"
          [class.flex-nowrap]="!isEditing()"
        >
          <button
            hlmBtn
            variant="outline"
            class="h-9 px-3 text-xs bg-red-100/70 text-red-700"
            [disabled]="deleteContextMutation.isPending()"
            (click)="deleteContext()"
          >
            <ng-icon name="lucideTrash2"></ng-icon>
            <span>Delete</span>
          </button>
          @if (!isEditing()) {
            <button hlmBtn variant="outline" class="h-9 px-3 text-xs" (click)="startEdit()">
              Edit
            </button>
          } @else {
            <button hlmBtn variant="outline" class="h-9 px-3 text-xs" (click)="cancelEdit()">
              Cancel
            </button>
            <button
              hlmBtn
              variant="outline"
              class="h-9 px-3 text-xs bg-primary/10 text-primary border-primary/25"
              [disabled]="updateContextMutation.isPending()"
              (click)="saveEdit()"
            >
              Save
            </button>
          }
          <button
            hlmBtn
            variant="outline"
            class="h-9 px-3 text-xs bg-blue-200/70 text-blue-600"
            (click)="startContext()"
          >
            <ng-icon name="lucidePlay"></ng-icon>
            <span class="font-semibold text-blue-600">Start</span>
          </button>
        </div>
      </div>

      <div class="flex w-full">
        <div class="w-full flex items-center justify-center gap-4">
          <div hlmCard class="w-full p-3 rounded-lg border">
            <h3
              class="text-[11px] uppercase tracking-[0.08em] text-muted-foreground font-semibold"
              hlmCardTitle
            >
              Total time
            </h3>
            <div class="text-lg font-semibold" hlmCardContet>
              {{ parseDuration(contextStats()?.totalDuration) }}
            </div>
          </div>
          <div hlmCard class="w-full p-3 rounded-lg border">
            <h3
              class="text-[11px] uppercase tracking-[0.08em] text-muted-foreground font-semibold"
              hlmCardTitle
            >
              Today
            </h3>
            <div class="text-lg font-semibold" hlmCardContet>
              {{ parseDuration(contextStats()?.duration) }}
            </div>
          </div>
          <div hlmCard class="w-full p-3 rounded-lg border">
            <h3
              class="text-[11px] uppercase tracking-[0.08em] text-muted-foreground font-semibold"
              hlmCardTitle
            >
              Sessions
            </h3>
            <div class="text-lg font-semibold" hlmCardContet>
              {{ contextStats()?.totalSessions }}
            </div>
          </div>
          <div hlmCard class="w-full p-3 rounded-lg border">
            <h3
              class="text-[11px] uppercase tracking-[0.08em] text-muted-foreground font-semibold"
              hlmCardTitle
            >
              Today sessions
            </h3>
            <div class="text-lg font-semibold" hlmCardContet>{{ contextStats()?.sessions }}</div>
          </div>
        </div>
      </div>
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
            <div
              class="w-full flex justify-start px-3 py-2.5 bg-card border border-border rounded-lg"
            >
              @if (editingIntervalId() === interval.id) {
                <div class="w-full flex flex-col gap-2">
                  <div class="w-full flex flex-col md:flex-row items-stretch md:items-end gap-2">
                    <label class="flex-1 text-xs text-muted-foreground">
                      Start
                      <input
                        type="datetime-local"
                        class="w-full h-9 rounded-md border border-border bg-background px-3 text-sm mt-1"
                        [value]="editIntervalStartInput()"
                        (input)="editIntervalStartInput.set(getInputValue($event))"
                      />
                    </label>
                    <label class="flex-1 text-xs text-muted-foreground">
                      End
                      <input
                        type="datetime-local"
                        class="w-full h-9 rounded-md border border-border bg-background px-3 text-sm mt-1"
                        [value]="editIntervalEndInput()"
                        (input)="editIntervalEndInput.set(getInputValue($event))"
                      />
                    </label>
                    <button
                      hlmBtn
                      variant="outline"
                      class="h-9 px-3 text-xs"
                      (click)="cancelIntervalEdit()"
                    >
                      Cancel
                    </button>
                    <button
                      hlmBtn
                      variant="outline"
                      class="h-9 px-3 text-xs bg-blue-200/70 text-blue-600"
                      [disabled]="updateIntervalMutation.isPending()"
                      (click)="saveIntervalEdit(interval)"
                    >
                      Save
                    </button>
                  </div>
                </div>
              } @else {
                <div class="flex gap-2 ml-4 w-full items-center">
                  <div class="flex flex-col flex-1">
                    <div class="text-sm font-medium">
                      {{ interval.start.toTimeString() }} - {{ interval.end.toTimeString() }}
                    </div>
                    <div class="text-xs text-muted-foreground">
                      {{
                        intervalStartDateEqEndDate(interval)
                          ? interval.start.toDateString()
                          : interval.start.toDateString() + ' - ' + interval.end.toDateString()
                      }}
                    </div>
                  </div>
                  <div class="text-xs text-muted-foreground">
                    {{ parseDuration(interval.duration) }}
                  </div>
                  <div class="flex items-center gap-1.5">
                    <button
                      hlmBtn
                      variant="outline"
                      class="h-7 px-2 text-xs"
                      (click)="startIntervalEdit(interval)"
                    >
                      <ng-icon name="lucidePencil"></ng-icon>
                    </button>
                    <button
                      hlmBtn
                      variant="outline"
                      class="h-7 px-2 text-xs"
                      [disabled]="movableContexts().length === 0"
                      (click)="openMoveDialog(interval)"
                    >
                      <ng-icon name="lucideArrowRightLeft"></ng-icon>
                    </button>
                    <button
                      hlmBtn
                      variant="outline"
                      class="h-7 px-2 text-xs text-red-700 bg-red-100/60"
                      [disabled]="deleteIntervalMutation.isPending()"
                      (click)="deleteInterval(interval)"
                    >
                      <ng-icon name="lucideTrash2"></ng-icon>
                    </button>
                  </div>
                </div>
              }
            </div>
          }
        </div>
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
      display: block;
      width: 100%;
      max-width: 1000px;
      height: 100%;
      min-height: 0;
    }
  `,
})
export class ContextComponent {
  private contextQueries = inject(ContextQueries);
  private contextMutations = inject(ContextMutations);
  private intervalMutations = inject(IntervalMutations);
  private router = inject(Router);

  switchContextMutation = injectMutation(() => this.contextMutations.switch());
  updateContextMutation = injectMutation(() => this.contextMutations.update());
  deleteContextMutation = injectMutation(() => this.contextMutations.delete());
  createIntervalMutation = injectMutation(() => this.intervalMutations.create());
  updateIntervalMutation = injectMutation(() => this.intervalMutations.update());
  deleteIntervalMutation = injectMutation(() => this.intervalMutations.delete());
  moveIntervalMutation = injectMutation(() => this.intervalMutations.move());
  contextQuery = injectQuery(() => this.contextQueries.get(this.contextId()));
  contextIntervalsQuery = injectQuery(() => this.contextQueries.intervals(this.contextId()));
  contextsQuery = injectQuery(() => this.contextQueries.list());
  context = computed(() => this.contextQuery.data() ?? EMPTY_CONTEXT);
  contextStatsQuery = injectQuery(() => this.contextQueries.stats(this.contextId(), this.today()));
  contextStats = computed(() => this.contextStatsQuery.data());
  today = signal(DateTime.local().toFormat('yyyy-MM-dd'));
  contextTags = computed(() => this.context().tags ?? []);
  editTagsPreview = computed(() =>
    this.editTagsInput()
      .split(',')
      .map((tag) => tag.trim())
      .filter((tag) => tag.length > 0),
  );
  intervals = computed(() => this.contextIntervalsQuery.data() ?? []);
  contexts = computed(() => this.contextsQuery.data() ?? []);
  movableContexts = computed(() =>
    this.contexts().filter((context) => context.id && context.id !== this.contextId()),
  );
  isEditing = signal(false);
  editName = signal('');
  editDescription = signal('');
  editTagsInput = signal('');
  newIntervalStartInput = signal('');
  newIntervalEndInput = signal('');
  editingIntervalId = signal<string | null>(null);
  editIntervalStartInput = signal('');
  editIntervalEndInput = signal('');
  moveDialogIntervalId = signal<string | null>(null);
  moveTargetContextId = signal('');
  intervalFormError = signal('');

  route = inject(ActivatedRoute);
  readonly contextId = toSignal(this.route.paramMap.pipe(map((pm) => pm.get('id') ?? '')), {
    initialValue: '',
  });

  constructor() {
    this.resetNewIntervalForm();

    effect(() => {
      const context = this.context();

      if (!this.isEditing()) {
        this.editName.set(context.name);
        this.editDescription.set(context.description);
        this.editTagsInput.set((context.tags ?? []).join(', '));
      }
    });
  }

  startContext() {
    this.switchContextMutation.mutate(this.context()!);
  }

  deleteContext() {
    const context = this.context();

    if (!context.id) {
      return;
    }

    if (!window.confirm(`Delete context "${context.name}"?`)) {
      return;
    }

    this.deleteContextMutation.mutate(context.id, {
      onSuccess: () => {
        this.router.navigate(['/day', this.today()]);
      },
    });
  }

  startEdit() {
    const context = this.context();
    this.editName.set(context.name);
    this.editDescription.set(context.description);
    this.editTagsInput.set((context.tags ?? []).join(', '));
    this.isEditing.set(true);
  }

  cancelEdit() {
    this.isEditing.set(false);
  }

  saveEdit() {
    const context = this.context();
    const tags = this.editTagsInput()
      .split(',')
      .map((tag) => tag.trim())
      .filter((tag) => tag.length > 0);

    this.updateContextMutation.mutate(
      {
        id: context.id,
        context: {
          ...context,
          name: this.editName(),
          description: this.editDescription(),
          tags,
        },
      },
      {
        onSuccess: () => {
          this.isEditing.set(false);
        },
      },
    );
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

  getInputValue(event: Event): string {
    return (event.target as HTMLInputElement | HTMLTextAreaElement).value;
  }

  getSelectValue(event: Event): string {
    return (event.target as HTMLSelectElement).value;
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

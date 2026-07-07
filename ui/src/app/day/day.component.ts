import { Component, computed, inject, signal } from '@angular/core';
import { toSignal } from '@angular/core/rxjs-interop';
import { ActivatedRoute } from '@angular/router';
import { map } from 'rxjs/operators';
import { DateTime } from 'luxon';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideFlag, lucidePlay } from '@ng-icons/lucide';
import { ContextQueries } from '../../api/context.quries';
import { DayStats } from '../../api/context.service';
import { ContextListComponent } from '../context/context-list.component';
import { ContextListItem } from '../context/context-list-item.component';
import { DistributionComponent, DistributionItem } from '../shared/distribution.component';
import { colorHash, durationAsHM } from '../utils';
import { Store } from '@ngxs/store';
import { WorkspaceState } from '../sidebar/workspace.state';

const EMPTY_DAY_STATS: DayStats = {
  date: '',
  contextStats: [],
  contexts: [],
  intervals: {},
  distribution: {},
};

@Component({
  selector: 'ctx-day',
  imports: [ContextListComponent, DistributionComponent, NgIcon],
  providers: [
    provideIcons({
      lucidePlay,
      lucideFlag,
    }),
  ],
  template: `
    <div class="w-full h-full overflow-hidden flex flex-col p-4 md:p-6">
      <div class="mb-5">
        <div class="text-[11px] uppercase tracking-widest text-muted-foreground font-semibold">
          Daily summary
        </div>
        <div class="flex justify-between items-center">
          <h1 class="text-2xl font-semibold tracking-tight mt-1">
            {{ formatDate(selectedDate()) }}
          </h1>
          <div>
            <div class="mt-2 flex flex-wrap items-center gap-3 text-xs text-muted-foreground">
              <div class="flex items-center justify-between gap-1">
                <span class="inline-flex items-center gap-1.5" [title]="'Start of first session'">
                  <ng-icon name="lucidePlay" class="text-[11px]"></ng-icon>
                </span>
                <span class="tabular-nums">{{ firstContextStart() }}</span>
              </div>
              <div class="flex items-center justify-between gap-1">
                <span class="inline-flex items-center gap-1.5" [title]="'End of last session'">
                  <ng-icon name="lucideFlag" class="text-[11px]"></ng-icon>
                </span>
                <span class="tabular-nums">{{ lastContextEnd() }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="grid grid-cols-2 md:grid-cols-3 xl:grid-cols-4 gap-2.5 mb-6">
        <div class="rounded-lg border bg-card px-3 py-2.5">
          <div class="text-[10px] uppercase tracking-[0.08em] text-muted-foreground">
            Total tracked
          </div>
          <div class="text-base font-semibold mt-1">{{ totalTracked() }}</div>
        </div>
        <div class="rounded-lg border bg-card px-3 py-2.5">
          <div class="text-[10px] uppercase tracking-[0.08em] text-muted-foreground">Contexts</div>
          <div class="text-base font-semibold mt-1">{{ contexts().length }}</div>
        </div>
        <div class="rounded-lg border bg-card px-3 py-2.5">
          <div class="text-[10px] uppercase tracking-[0.08em] text-muted-foreground">Sessions</div>
          <div class="text-base font-semibold mt-1">{{ totalSessions() }}</div>
        </div>
        <div class="rounded-lg border bg-card px-3 py-2.5">
          <div class="text-[10px] uppercase tracking-[0.08em] text-muted-foreground">
            Top context
          </div>
          <div class="text-sm font-medium mt-1 truncate">{{ topContext() }}</div>
        </div>
      </div>

      <ctx-distribution
        class="block mb-6"
        [items]="distributionContexts()"
        emptyMessage="No tracked time for this day."
      ></ctx-distribution>

      <div
        class="text-[11px] uppercase tracking-[0.08em] text-muted-foreground font-semibold mb-2 shrink-0"
      >
        Contexts
      </div>
      <div class="flex-1 min-h-0 overflow-auto pr-1 pb-2">
        <ctx-context-list
          [items]="contexts()"
          emptyMessage="No contexts tracked for this day."
        ></ctx-context-list>
      </div>
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
export class DayComponent {
  private contextQueries = inject(ContextQueries);
  private store = inject(Store);
  private activeWorkspaceId = this.store.selectSignal(WorkspaceState.selectedWorkspaceId);
  route = inject(ActivatedRoute);
  today = signal(DateTime.local().toFormat('yyyy-MM-dd'));
  readonly selectedDate = toSignal(
    this.route.paramMap.pipe(map((pm) => pm.get('date') ?? this.today())),
    {
      initialValue: this.today(),
    },
  );

  dayStatsQuery = injectQuery(() =>
    this.contextQueries.dayStats(this.activeWorkspaceId(), this.selectedDate()),
  );
  dayStats = computed(() => this.dayStatsQuery.data() ?? EMPTY_DAY_STATS);
  firstContextStart = computed(() => {
    const intervals = Object.values(this.dayStats().intervals).flat();
    if (intervals.length === 0) {
      return '-';
    }

    let firstStart = intervals[0].start;
    for (const interval of intervals) {
      if (interval.start < firstStart) {
        firstStart = interval.start;
      }
    }
    return this.formatTime(firstStart);
  });
  lastContextEnd = computed(() => {
    const intervals = Object.values(this.dayStats().intervals).flat();
    if (intervals.length === 0) {
      return '-';
    }

    let lastEnd = intervals[0].end;
    for (const interval of intervals) {
      if (interval.end > lastEnd) {
        lastEnd = interval.end;
      }
    }
    return this.formatTime(lastEnd);
  });

  contexts = computed<ContextListItem[]>(() => {
    const contextsById = new Map(this.dayStats().contexts.map((context) => [context.id, context]));

    const mappedContexts = this.dayStats()
      .contextStats.map((contextStats) => {
        const context = contextsById.get(contextStats.contextId);
        const distributionValue = this.dayStats().distribution[contextStats.contextId];
        const distributionPercent = distributionValue ?? contextStats.percentage;

        return {
          id: contextStats.contextId,
          name: context?.name ?? contextStats.contextId,
          duration: durationAsHM(contextStats.duration),
          percentage: distributionPercent,
          distributionPercentage: distributionPercent,
          color: colorHash(context?.id ?? contextStats.contextId),
          sessions: contextStats.intervalCount,
          sessionRanges: (this.dayStats().intervals[contextStats.contextId] ?? []).map(
            (interval) => `${interval.start.toTimeString()}–${interval.end.toTimeString()}`,
          ),
        };
      })
      .sort((left, right) => right.percentage - left.percentage);

    return mappedContexts;
  });

  distributionContexts = computed<DistributionItem[]>(() =>
    this.contexts().map((context) => ({
      id: context.id,
      name: context.name,
      duration: context.duration,
      percentage: context.distributionPercentage ?? context.percentage,
      color: context.color,
    })),
  );

  totalTracked = computed(() => {
    const duration = this.dayStats().contextStats.reduce(
      (sum, context) => sum + context.duration,
      0,
    );
    return durationAsHM(duration);
  });

  totalSessions = computed(() =>
    this.dayStats().contextStats.reduce((sum, context) => sum + context.intervalCount, 0),
  );

  topContext = computed(() => this.contexts()[0]?.name ?? '-');

  formatDate(date: string): string {
    return DateTime.fromFormat(date, 'yyyy-MM-dd').toFormat('dd.MM.yyyy');
  }

  formatTime(date: { toTimeString: () => string }): string {
    return date.toTimeString().slice(0, 5);
  }
}

import { Component, computed, inject, signal } from '@angular/core';
import { toSignal } from '@angular/core/rxjs-interop';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { map } from 'rxjs/operators';
import { DateTime } from 'luxon';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { ContextQueries } from '../../api/context.quries';
import { DayStats } from '../../api/context.service';
import { colorHash, durationAsHM } from '../utils';

const EMPTY_DAY_STATS: DayStats = {
  date: '',
  contextStats: [],
  contexts: [],
  intervals: {},
  distribution: {},
};

@Component({
  selector: 'app-day',
  imports: [RouterLink],
  template: `
    <div class="w-full h-full overflow-hidden flex flex-col p-4 md:p-6">
      <div class="mb-5">
        <div class="text-[11px] uppercase tracking-widest text-muted-foreground font-semibold">
          Daily summary
        </div>
        <h1 class="text-2xl font-semibold tracking-tight mt-1">{{ formatDate(selectedDate()) }}</h1>
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
        <div class="rounded-lg border bg-card px-3 py-2.5 col-span-2 md:col-span-1">
          <div class="text-[10px] uppercase tracking-[0.08em] text-muted-foreground">
            Top context
          </div>
          <div class="text-sm font-medium mt-1 truncate">{{ topContext() }}</div>
        </div>
      </div>

      <div class="mb-6">
        <div
          class="text-[11px] uppercase tracking-[0.08em] text-muted-foreground font-semibold mb-2"
        >
          Distribution
        </div>
        <div class="flex h-2 rounded-md overflow-hidden gap-px bg-muted/40">
          @for (context of contexts(); track context.id) {
            <div
              [style.width.%]="context.distributionPercent"
              [style.background-color]="context.color"
            ></div>
          }
        </div>
      </div>

      <div
        class="text-[11px] uppercase tracking-[0.08em] text-muted-foreground font-semibold mb-2 shrink-0"
      >
        Contexts
      </div>
      <div class="flex-1 min-h-0 overflow-auto pr-1 pb-2">
        <div class="flex flex-col gap-2">
          @for (context of contexts(); track context.id) {
            <div
              class="rounded-lg border bg-card p-3 hover:bg-muted/30 transition-colors cursor-pointer"
              [routerLink]="['/context', context.id]"
            >
              <div class="flex items-center gap-2 mb-2">
                <span
                  class="w-2 h-2 rounded-sm shrink-0"
                  [style.background-color]="context.color"
                ></span>
                <span class="text-sm font-medium flex-1 truncate">{{ context.name }}</span>
                <span class="text-xs text-muted-foreground">{{ context.duration }}</span>
              </div>
              <div class="h-1.5 rounded bg-muted/40 overflow-hidden">
                <div
                  class="h-full rounded"
                  [style.width.%]="context.percent"
                  [style.background-color]="context.color"
                ></div>
              </div>
              <div class="mt-2 flex flex-wrap gap-x-2 gap-y-1 text-[10px] text-muted-foreground">
                @for (session of context.sessions; track session) {
                  <span>{{ session }}</span>
                }
              </div>
            </div>
          }
        </div>
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
  route = inject(ActivatedRoute);
  today = signal(DateTime.local().toFormat('yyyy-MM-dd'));
  readonly selectedDate = toSignal(
    this.route.paramMap.pipe(map((pm) => pm.get('date') ?? this.today())),
    {
      initialValue: this.today(),
    },
  );

  dayStatsQuery = injectQuery(() => this.contextQueries.dayStats(this.selectedDate()));
  dayStats = computed(() => this.dayStatsQuery.data() ?? EMPTY_DAY_STATS);

  contexts = computed(() => {
    const contextsById = new Map(this.dayStats().contexts.map((context) => [context.id, context]));

    const mappedContexts = this.dayStats()
      .contextStats.map((contextStats) => {
        const context = contextsById.get(contextStats.contextId);
        const distributionValue = this.dayStats().distribution[contextStats.contextId];
        const distributionPercent =
          distributionValue === undefined
            ? undefined
            : distributionValue <= 1
              ? distributionValue * 100
              : distributionValue;

        return {
          id: contextStats.contextId,
          name: context?.name ?? contextStats.contextId,
          duration: durationAsHM(contextStats.duration),
          percent: distributionPercent ?? contextStats.percentage,
          distributionPercent: 0,
          color: colorHash(context?.id ?? contextStats.contextId),
          sessions: (this.dayStats().intervals[contextStats.contextId] ?? []).map(
            (interval) => `${interval.start.toTimeString()}–${interval.end.toTimeString()}`,
          ),
        };
      })
      .sort((left, right) => right.percent - left.percent);

    const totalPercent = mappedContexts.reduce(
      (sum, context) => sum + Math.max(context.percent, 0),
      0,
    );

    return mappedContexts.map((context) => ({
      ...context,
      distributionPercent:
        totalPercent > 0 ? (Math.max(context.percent, 0) / totalPercent) * 100 : 0,
    }));
  });

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
}

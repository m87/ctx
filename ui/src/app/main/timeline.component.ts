import { Component, computed, inject, signal } from '@angular/core';
import { toSignal } from '@angular/core/rxjs-interop';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { DateTime } from 'luxon';
import { IntervalQueries } from '../../api/interval.queries';
import { DayIntervalsResponse } from '../../api/interval.service';
import { colorHash } from '../utils';
import { NavigationEnd, Router, RouterLink } from '@angular/router';
import { filter, map, startWith } from 'rxjs/operators';

const EMPTY_DAY_INTERVALS: DayIntervalsResponse = {
  contexts: [],
  intervals: [],
};

@Component({
  imports: [RouterLink],
  selector: 'app-timeline',
  template: `
    <div class="w-full border-t bg-background px-4 py-2">
      <div class="text-[10px] text-muted-foreground mb-1.5 tracking-[0.08em] uppercase">
        Timeline — {{ formatDate(selectedDate()) }}
      </div>

      <div class="relative h-3.25 mb-1">
        @for (mark of hourMarks; track mark.hour) {
          <div
            class="absolute text-[9px] text-muted-foreground whitespace-nowrap -translate-x-1/2 leading-none"
            [style.left.%]="getHourPosition(mark.hour)"
          >
            {{ mark.label }}
          </div>
        }
      </div>

      <div class="relative h-3.5">
        @for (mark of hourMarks; track mark.hour) {
          <div
            class="absolute top-0 bottom-0 border-l border-border pointer-events-none"
            [style.left.%]="getHourPosition(mark.hour)"
          ></div>
        }

        <div class="absolute inset-0 bg-muted/30 rounded-lg"></div>

        @for (interval of intervals(); track interval.id) {
          @if (getWidth(interval.from, interval.to) > 0) {
            <div
              class="absolute top-0 h-full rounded-[3px] opacity-85 hover:opacity-100 hover:scale-y-110 transition-all duration-100 origin-center cursor-pointer"
              [style.background-color]="interval.color"
              [style.left.%]="getLeft(interval.from)"
              [style.width.%]="getWidth(interval.from, interval.to)"
              [attr.aria-label]="'Interwał od ' + interval.from + ' do ' + interval.to"
              (click)="selectLegendContext(interval.contextId)"
            ></div>
          }
        }
      </div>

      <div class="flex flex-wrap gap-x-3 gap-y-1.5 mt-2">
        @for (context of visibleLegendContexts(); track context.id) {
          <div
            class="flex items-center gap-1.5 text-[10px] text-muted-foreground hover:text-foreground cursor-default"
            [routerLink]="['/context', context.id]"
          >
            <span
              class="w-1.75 h-1.75 rounded-sm shrink-0"
              [style.background-color]="context.color"
            ></span>
            {{ context.name }}
          </div>
        }
      </div>
    </div>
  `,
})
export class TimelineComponent {
  private intervalQueries = inject(IntervalQueries);
  private router = inject(Router);
  private today = DateTime.local().toFormat('yyyy-MM-dd');

  selectedDay = toSignal(
    this.router.events.pipe(
      filter((event): event is NavigationEnd => event instanceof NavigationEnd),
      startWith(null),
      map(() => this.extractDayFromUrl(this.router.url) ?? this.today),
    ),
    {
      initialValue: this.extractDayFromUrl(this.router.url) ?? this.today,
    },
  );
  selectedDate = computed(() => DateTime.fromFormat(this.selectedDay(), 'yyyy-MM-dd').toJSDate());
  dayIntervalsQuery = injectQuery(() => this.intervalQueries.day(this.selectedDay()));
  dayIntervals = computed(() => this.dayIntervalsQuery.data() ?? EMPTY_DAY_INTERVALS);
  private selectedLegendContextId = signal<string | null>(null);

  intervals = computed(() => {
    const contextsById = new Map(
      this.dayIntervals().contexts.map((context) => [context.id, context]),
    );

    return this.dayIntervals()
      .intervals.map((interval) => {
        const contextId = interval.contextId ?? interval.context_id ?? '';
        const context = contextsById.get(contextId);
        const colorKey = context?.id || contextId || interval.id;
        const durationMinutes = Math.max(
          interval.duration > 0
            ? interval.duration / 60
            : interval.end.toDateTime().diff(interval.start.toDateTime(), 'minutes').minutes,
          0,
        );

        return {
          id: interval.id,
          contextId,
          from: interval.start.toTimeString(),
          to: interval.end.toTimeString(),
          durationMinutes,
          color: colorHash(colorKey),
        };
      })
      .filter((interval) => interval.contextId !== '' && interval.durationMinutes > 0);
  });

  todayContexts = computed(() => {
    const durationsByContextId = this.intervals().reduce((result, interval) => {
      result.set(
        interval.contextId,
        (result.get(interval.contextId) ?? 0) + interval.durationMinutes,
      );
      return result;
    }, new Map<string, number>());

    return this.dayIntervals()
      .contexts.filter((context) => durationsByContextId.has(context.id))
      .map((context) => ({
        id: context.id,
        name: context.name,
        color: colorHash(context.id),
        durationMinutes: durationsByContextId.get(context.id) ?? 0,
      }))
      .sort((left, right) => right.durationMinutes - left.durationMinutes);
  });

  visibleLegendContexts = computed(() => {
    const contexts = this.todayContexts();
    const defaultContexts = contexts.slice(0, 5);
    const selectedContextId = this.selectedLegendContextId();

    if (!selectedContextId || defaultContexts.some((context) => context.id === selectedContextId)) {
      return defaultContexts;
    }

    const selectedContext = contexts.find((context) => context.id === selectedContextId);
    return selectedContext ? [...defaultContexts, selectedContext] : defaultContexts;
  });

  dayStartHour = 0;
  dayEndHour = 24;

  hourMarks = [
    { hour: 0, label: '0:00' },
    { hour: 6, label: '6:00' },
    { hour: 12, label: '12:00' },
    { hour: 18, label: '18:00' },
    { hour: 24, label: '24:00' },
  ];

  private readonly totalDayMinutes = 24 * 60;

  private extractDayFromUrl(url: string): string | null {
    const normalizedUrl = url.split('?')[0].split('#')[0];
    const dayMatch = normalizedUrl.match(/\/day\/([^/]+)/);

    if (!dayMatch) {
      return null;
    }

    const parsedDate = DateTime.fromFormat(dayMatch[1], 'yyyy-MM-dd');
    return parsedDate.isValid ? dayMatch[1] : null;
  }

  formatDate(date: Date): string {
    return new Intl.DateTimeFormat('pl-PL', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
    }).format(date);
  }

  getHourPosition(hour: number): number {
    return ((hour - this.dayStartHour) / (this.dayEndHour - this.dayStartHour)) * 100;
  }

  private toMinutes(time: string): number {
    const [hour, minute] = time.split(':').map(Number);
    if (Number.isNaN(hour) || Number.isNaN(minute)) {
      return 0;
    }

    return hour * 60 + minute;
  }

  getLeft(from: string): number {
    return (this.toMinutes(from) / this.totalDayMinutes) * 100;
  }

  getWidth(from: string, to: string): number {
    const start = this.toMinutes(from);
    const end = this.toMinutes(to);
    return (Math.max(end - start, 1) / this.totalDayMinutes) * 100;
  }

  selectLegendContext(contextId: string): void {
    this.selectedLegendContextId.set(contextId);
  }
}

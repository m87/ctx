import { Component, computed, inject, signal } from '@angular/core';
import { toSignal } from '@angular/core/rxjs-interop';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { DateTime } from 'luxon';
import { IntervalQueries } from '../../api/interval/interval.queries';
import { DayIntervalsResponse } from '../../api/interval/interval.service';
import { colorHash } from '../utils';
import { NavigationEnd, Router, RouterLink } from '@angular/router';
import { filter, map, startWith } from 'rxjs/operators';
import { Store } from '@ngxs/store';
import { WorkspaceState } from '../sidebar/workspace.state';
import { TimeZoneService } from '../shared/time-zone.service';
import { HlmSkeletonImports } from '@spartan-ng/helm/skeleton';
import { SettingsQueries } from '../../api/settings/settings.queries';
import {
  createTimelineMarks,
  minuteWithinDay,
  normalizeTimelineRangeMode,
  resolveTimelineRange,
  timelinePosition,
  timelineRangeSettingKey,
  timelineWidth,
} from '../shared/timeline-range';

const EMPTY_DAY_INTERVALS: DayIntervalsResponse = {
  contexts: [],
  intervals: [],
};

@Component({
  imports: [RouterLink, HlmSkeletonImports],
  selector: 'ctx-timeline',
  template: `
    <div class="w-full border-t bg-background px-4 py-2">
      <div class="text-caption text-muted-foreground mb-1.5 tracking-label uppercase">
        Timeline — {{ formatDate(selectedDay()) }}
      </div>

      <div class="relative h-3.25 mb-1">
        @for (mark of hourMarks(); track mark.minute; let first = $first; let last = $last) {
          <div
            class="absolute text-label text-muted-foreground whitespace-nowrap leading-none"
            [class.-translate-x-1/2]="!first && !last"
            [class.-translate-x-full]="last"
            [style.left.%]="getHourPosition(mark.minute)"
          >
            {{ mark.label }}
          </div>
        }
      </div>

      <div class="relative h-3.5">
        @for (mark of hourMarks(); track mark.minute) {
          <div
            class="absolute top-0 bottom-0 border-l border-border pointer-events-none"
            [style.left.%]="getHourPosition(mark.minute)"
          ></div>
        }

        @if (dayIntervalsQuery.isLoading()) {
          <hlm-skeleton class="h-full w-full"></hlm-skeleton>
        } @else {
          <div class="absolute inset-0 bg-muted/30 rounded-lg"></div>

          @for (interval of intervals(); track interval.id) {
            @if (getWidth(interval.startMinutes, interval.endMinutes) > 0) {
              <button
                type="button"
                class="absolute top-0 h-full rounded-md opacity-85 hover:opacity-100 hover:scale-y-110 transition-all duration-100 origin-center cursor-pointer focus-visible:ring-2 focus-visible:ring-ring/40 focus-visible:outline-none"
                [style.background-color]="interval.color"
                [style.left.%]="getLeft(interval.startMinutes)"
                [style.width.%]="getWidth(interval.startMinutes, interval.endMinutes)"
                [attr.aria-label]="'Interval from ' + interval.from + ' to ' + interval.to"
                (click)="selectLegendContext(interval.contextId)"
              ></button>
            }
          }
        }
      </div>

      <div class="flex flex-wrap gap-x-3 gap-y-1.5 mt-2">
        @if (dayIntervalsQuery.isLoading()) {
          <span class="sr-only" role="status">Loading timeline</span>
          @for (item of legendSkeletonItems; track item) {
            <hlm-skeleton class="h-2.5 w-20"></hlm-skeleton>
          }
        } @else {
          @for (context of visibleLegendContexts(); track context.id) {
            <a
              class="flex items-center gap-1.5 text-caption text-muted-foreground hover:text-foreground cursor-pointer"
              [routerLink]="['/context', context.id]"
            >
              <span
                class="w-1.75 h-1.75 rounded-sm shrink-0"
                [style.background-color]="context.color"
              ></span>
              {{ context.name }}
            </a>
          }
        }
      </div>
    </div>
  `,
})
export class TimelineComponent {
  readonly legendSkeletonItems = [0, 1, 2];
  private intervalQueries = inject(IntervalQueries);
  private settingsQueries = inject(SettingsQueries);
  private router = inject(Router);
  private store = inject(Store);
  private timeZone = inject(TimeZoneService);
  private activeWorkspaceId = this.store.selectSignal(WorkspaceState.selectedWorkspaceId);

  private routedDay = toSignal(
    this.router.events.pipe(
      filter((event): event is NavigationEnd => event instanceof NavigationEnd),
      startWith(null),
      map(() => this.extractDayFromUrl(this.router.url)),
    ),
    {
      initialValue: this.extractDayFromUrl(this.router.url),
    },
  );
  selectedDay = computed(() => this.routedDay() ?? this.timeZone.today());
  dayIntervalsQuery = injectQuery(() =>
    this.intervalQueries.day(
      this.activeWorkspaceId(),
      this.selectedDay(),
      this.timeZone.effectiveTimeZone(),
    ),
  );
  dayIntervals = computed(() => this.dayIntervalsQuery.data() ?? EMPTY_DAY_INTERVALS);
  settingsQuery = injectQuery(() => this.settingsQueries.settings());
  private selectedLegendContextId = signal<string | null>(null);

  intervals = computed(() => {
    const contextsById = new Map(
      this.dayIntervals().contexts.map((context) => [context.id, context]),
    );

    return this.dayIntervals()
      .intervals.map((interval) => {
        const contextId = interval.contextId ?? '';
        const context = contextsById.get(contextId);
        const colorKey = context?.id || contextId || interval.id;
        const startMinutes = minuteWithinDay(
          interval.start,
          this.selectedDay(),
          this.timeZone.effectiveTimeZone(),
        );
        const endMinutes = minuteWithinDay(
          interval.end,
          this.selectedDay(),
          this.timeZone.effectiveTimeZone(),
        );
        if (startMinutes === null || endMinutes === null || endMinutes <= startMinutes) {
          return null;
        }
        const durationMinutes = Math.max(
          interval.duration > 0
            ? interval.duration / 60_000_000_000
            : this.timeZone
                .parseInstant(interval.end)
                .diff(this.timeZone.parseInstant(interval.start), 'minutes').minutes,
          0,
        );

        return {
          id: interval.id,
          contextId,
          from: this.timeZone.formatTime(interval.start),
          to: this.timeZone.formatTime(interval.end),
          startMinutes,
          endMinutes,
          durationMinutes,
          color: colorHash(colorKey),
        };
      })
      .filter(
        (interval): interval is NonNullable<typeof interval> =>
          interval !== null && interval.contextId !== '' && interval.durationMinutes > 0,
      );
  });

  timelineRangeMode = computed(() =>
    normalizeTimelineRangeMode(this.settingsQuery.data()?.[timelineRangeSettingKey]),
  );
  timelineRange = computed(() => resolveTimelineRange(this.intervals(), this.timelineRangeMode()));
  hourMarks = computed(() => createTimelineMarks(this.timelineRange()));

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

  private extractDayFromUrl(url: string): string | null {
    const normalizedUrl = url.split('?')[0].split('#')[0];
    const dayMatch = normalizedUrl.match(/\/day\/([^/]+)/);

    if (!dayMatch) {
      return null;
    }

    const parsedDate = DateTime.fromFormat(dayMatch[1], 'yyyy-MM-dd');
    return parsedDate.isValid ? dayMatch[1] : null;
  }

  formatDate(date: string): string {
    return DateTime.fromFormat(date, 'yyyy-MM-dd').toFormat('dd.MM.yyyy');
  }

  getHourPosition(minute: number): number {
    return timelinePosition(minute, this.timelineRange());
  }

  getLeft(startMinutes: number): number {
    return timelinePosition(startMinutes, this.timelineRange());
  }

  getWidth(startMinutes: number, endMinutes: number): number {
    return timelineWidth(startMinutes, endMinutes, this.timelineRange());
  }

  selectLegendContext(contextId: string): void {
    this.selectedLegendContextId.set(contextId);
  }
}

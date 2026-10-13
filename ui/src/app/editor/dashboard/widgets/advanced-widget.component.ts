import { RouterLink } from '@angular/router';
import { Component, computed, input, signal } from '@angular/core';
import { DashboardMetricResult } from '../../../../api/dashboard/dashboard.service';
import {
  BarWidgetDefinition,
  HeatmapWidgetDefinition,
  LineWidgetDefinition,
  TableWidgetDefinition,
} from '../dashboard-definition';
import { BarWidgetComponent } from './bar-widget.component';
import { LineWidgetComponent } from './line-widget.component';
import { TableWidgetComponent } from './table-widget.component';
import { defaultMetricProperties, formatMetric } from '../metric-widget-definition';

type AdvancedWidget =
  BarWidgetDefinition | HeatmapWidgetDefinition | LineWidgetDefinition | TableWidgetDefinition;

@Component({
  selector: 'ctx-advanced-widget',
  imports: [BarWidgetComponent, LineWidgetComponent, RouterLink, TableWidgetComponent],
  host: {
    class: 'block h-full min-h-0 min-w-0 overflow-auto',
    '(pointerdown)': 'stopPointer($event)',
    '(keydown)': 'stopChildKeyboard($event)',
  },
  template: `
    <div class="flex h-full min-h-0 flex-col gap-3 p-3">
      <header class="flex min-w-0 items-baseline justify-between gap-2">
        <h3 class="truncate text-sm font-medium">{{ title() }}</h3>
        <span class="shrink-0 text-xs text-muted-foreground">{{ period() }}</span>
      </header>
      @if (error()) {
        <p role="status" class="text-sm text-destructive">{{ error() }}</p>
      } @else if (result(); as data) {
        @if (data.error) {
          <p role="status" class="text-sm text-destructive">{{ data.error }}</p>
        } @else if (data.kind === 'categories') {
          @if (bar(); as config) {
            @defer {
              <ctx-bar-widget [widget]="config" [result]="data" [editing]="editing()" />
            } @placeholder {
              <p role="status" class="text-sm text-muted-foreground">Loading…</p>
            }
          }
        } @else if (data.kind === 'heatmap') {
          @if (heatmap(); as config) {
            @if (config.properties.showMonthLabels) {
              <div
                class="grid grid-flow-col grid-rows-1 gap-1 overflow-x-auto"
                [style.grid-auto-columns.px]="24"
                aria-hidden="true"
              >
                @for (week of heatmapWeeks(); let i = $index; track $index) {
                  <span class="text-meta text-muted-foreground">{{ weekMonthLabel(i) }}</span>
                }
              </div>
            }
            <div class="flex min-h-0 flex-1 gap-1 overflow-x-auto">
              @if (config.properties.showDayLabels) {
                <div
                  class="grid shrink-0 grid-rows-7 gap-1 text-meta text-muted-foreground"
                  aria-hidden="true"
                >
                  @for (day of weekDayLabels(); track $index) {
                    <span class="flex size-5 items-center justify-center sm:size-6">{{ day }}</span>
                  }
                </div>
              }
              <div
                class="grid min-h-0 flex-1 grid-flow-col grid-rows-7 content-start gap-1"
                role="grid"
                [attr.aria-label]="config.properties.title || 'Daily activity calendar'"
                [attr.aria-activedescendant]="activeHeatmapCellId()"
              >
                @for (week of heatmapWeeks(); track $index) {
                  @for (cell of week; track cell.date ?? 'blank') {
                    @if (cell.day; as day) {
                      <a
                        role="gridcell"
                        [routerLink]="['/day', day.date]"
                        class="size-5 rounded-sm border border-border focus-visible:outline focus-visible:outline-2 focus-visible:outline-ring sm:size-6"
                        [class.bg-muted]="day.value === 0 || day.value === null"
                        [class.bg-primary/25]="intensity(day.value, data.unit) === 1"
                        [class.bg-primary/50]="intensity(day.value, data.unit) === 2"
                        [class.bg-primary/75]="intensity(day.value, data.unit) === 3"
                        [class.bg-primary]="intensity(day.value, data.unit) === 4"
                        [class.border-dashed]="day.status === 'future' || day.value === null"
                        [class.ring-1]="day.date === today(data.timeZone)"
                        [class.ring-ring]="day.date === today(data.timeZone)"
                        [tabIndex]="activeHeatmapIndex() === cell.index ? 0 : -1"
                        [id]="heatmapCellId(cell.index)"
                        [attr.aria-label]="
                          heatmapLabel(
                            day,
                            data.unit,
                            intensity(day.value, data.unit),
                            config.properties.format,
                            config.properties.precision
                          )
                        "
                        [title]="
                          heatmapLabel(
                            day,
                            data.unit,
                            intensity(day.value, data.unit),
                            config.properties.format,
                            config.properties.precision
                          )
                        "
                        (focus)="activeHeatmapIndex.set(cell.index)"
                        (keydown)="onHeatmapKeydown($event, cell.index)"
                      ></a>
                    } @else {
                      <span class="size-5 sm:size-6" aria-hidden="true"></span>
                    }
                  }
                }
              </div>
            </div>
            @if (config.properties.showLegend) {
              <div
                class="flex items-center gap-1 text-meta text-muted-foreground"
                aria-label="Activity intensity legend"
              >
                <span>Less</span><span class="size-3 rounded-sm border bg-muted"></span
                ><span class="size-3 rounded-sm bg-primary/25"></span
                ><span class="size-3 rounded-sm bg-primary/50"></span
                ><span class="size-3 rounded-sm bg-primary/75"></span
                ><span class="size-3 rounded-sm bg-primary"></span><span>More</span>
              </div>
            }
            @if (config.properties.showSummary) {
              <p class="text-xs text-muted-foreground">
                {{ data.heatmap.activeDays }} active days ·
                {{
                  format(
                    data.heatmap.totalValue,
                    data.unit,
                    config.properties.format,
                    config.properties.precision
                  )
                }}
              </p>
            }
          }
        } @else if (data.kind === 'series') {
          @if (line(); as config) {
            @defer {
              <ctx-line-widget [widget]="config" [result]="data" [editing]="editing()" />
            } @placeholder {
              <p role="status" class="text-sm text-muted-foreground">Loading…</p>
            }
          }
        } @else if (data.kind === 'table') {
          @if (table(); as config) {
            <ctx-table-widget [widget]="config" [result]="data" />
          }
        }
      } @else {
        <p class="text-sm text-muted-foreground">Loading…</p>
      }
    </div>
  `,
})
export class AdvancedWidgetComponent {
  readonly widget = input.required<AdvancedWidget>();
  readonly result = input<DashboardMetricResult>();
  readonly error = input('');
  readonly period = input('');
  readonly firstDayOfWeek = input<0 | 1>(1);
  readonly editing = input(false);
  readonly bar = computed(() => {
    const widget = this.widget();
    return widget.type === 'bar' ? widget : null;
  });
  readonly heatmap = computed(() => {
    const widget = this.widget();
    return widget.type === 'heatmap' ? widget : null;
  });
  readonly line = computed(() => {
    const widget = this.widget();
    return widget.type === 'line' ? widget : null;
  });
  readonly table = computed(() => {
    const widget = this.widget();
    return widget.type === 'table' ? widget : null;
  });
  readonly title = computed(() => this.widget().properties.title || this.widget().type);
  readonly activeHeatmapIndex = signal(0);
  readonly activeHeatmapCellId = computed(() => this.heatmapCellId(this.activeHeatmapIndex()));
  readonly heatmapWeeks = computed(() => {
    const data = this.result();
    if (data?.kind !== 'heatmap') return [];
    const config = this.heatmap(),
      firstDay = this.firstDayOfWeek();
    const days = [...data.heatmap.days];
    if (!days.length) return [];
    const start = new Date(days[0].date + 'T00:00:00Z'),
      leading = (start.getUTCDay() - firstDay + 7) % 7;
    const cells: { date: string | null; day: (typeof days)[number] | null; index: number }[] =
      Array.from({ length: leading }, (_, i) => ({ date: null, day: null, index: -1 }));
    days.forEach((day, index) => cells.push({ date: day.date, day, index }));
    while (cells.length % 7) cells.push({ date: null, day: null, index: -1 });
    void config;
    return Array.from({ length: cells.length / 7 }, (_, week) =>
      cells.slice(week * 7, week * 7 + 7),
    );
  });

  format(
    value: number | null,
    unit: 'ms' | 'count' | 'percent',
    format?: 'duration' | 'hours' | 'number' | 'percent',
    precision?: number,
  ): string {
    if (value === null) return '—';
    const widget = this.widget();
    const properties = 'format' in widget.properties ? widget.properties : null;
    const selectedFormat =
      format ??
      properties?.format ??
      (unit === 'ms' ? 'duration' : unit === 'count' ? 'number' : 'percent');
    const selectedPrecision = precision ?? properties?.precision ?? (unit === 'count' ? 0 : 1);
    const metric =
      unit === 'ms' ? 'duration.sum' : unit === 'count' ? 'sessions.count' : 'time_share.ratio';
    return formatMetric(value, {
      ...defaultMetricProperties(),
      metric,
      format: selectedFormat,
      precision: selectedPrecision,
    });
  }
  intensity(value: number | null, unit: 'ms' | 'count' | 'percent'): number {
    if (value === null || value <= 0) return 0;
    const config = this.heatmap();
    if (!config) return 0;
    let thresholds = config.properties.thresholds;
    if (config.properties.scaleMode === 'auto') {
      const data = this.result();
      const maxValue =
        data?.kind === 'heatmap'
          ? Math.max(0, ...data.heatmap.days.map((day) => day.value ?? 0))
          : 0;
      thresholds = [maxValue / 4, maxValue / 2, (maxValue * 3) / 4];
    }
    void unit;
    return value <= thresholds[0] ? 1 : value <= thresholds[1] ? 2 : value <= thresholds[2] ? 3 : 4;
  }
  heatmapCellId(index: number): string {
    return `${this.widget().id}-day-${index}`;
  }
  heatmapLabel(
    day: { date: string; value: number | null; status: string },
    unit: 'ms' | 'count' | 'percent',
    level: number,
    format?: 'duration' | 'hours' | 'number' | 'percent',
    precision?: number,
  ): string {
    const date = new Date(`${day.date}T00:00:00Z`).toLocaleDateString(undefined, {
      weekday: 'long',
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      timeZone: 'UTC',
    });
    const value =
      day.status === 'future' ? 'future day' : this.format(day.value, unit, format, precision);
    return `${date}, ${value}, ${day.status === 'partial' ? 'day in progress' : day.status === 'future' ? '' : 'complete day'}${level ? `, level ${level} of 4` : ''}`;
  }
  weekMonthLabel(index: number): string {
    const weeks = this.heatmapWeeks(),
      week = weeks[index] ?? [];
    const first =
      week.find((cell) => cell.day?.date.endsWith('-01'))?.day ??
      (index === 0 ? week.find((cell) => cell.day)?.day : null);
    return first
      ? new Date(`${first.date}T00:00:00Z`).toLocaleDateString(undefined, {
          month: 'short',
          timeZone: 'UTC',
        })
      : '';
  }
  weekDayLabels(): string[] {
    return Array.from({ length: 7 }, (_, index) =>
      new Intl.DateTimeFormat(undefined, { weekday: 'narrow', timeZone: 'UTC' }).format(
        new Date(Date.UTC(2023, 0, 1, (this.firstDayOfWeek() + index) * 24)),
      ),
    );
  }
  today(timeZone: string): string {
    const parts = new Intl.DateTimeFormat('en-CA', {
      timeZone,
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
    }).formatToParts(new Date());
    return `${parts.find((part) => part.type === 'year')?.value}-${parts.find((part) => part.type === 'month')?.value}-${parts.find((part) => part.type === 'day')?.value}`;
  }
  onHeatmapKeydown(event: KeyboardEvent, index: number): void {
    let next = index;
    if (event.key === 'ArrowRight') next = index + 7;
    else if (event.key === 'ArrowLeft') next = index - 7;
    else if (event.key === 'ArrowDown') next = index + 1;
    else if (event.key === 'ArrowUp') next = index - 1;
    else if (event.key === 'Home') next = 0;
    else if (event.key === 'End') {
      const data = this.result();
      next = data?.kind === 'heatmap' ? data.heatmap.days.length - 1 : 0;
    } else return;
    event.preventDefault();
    const data = this.result();
    const length = data?.kind === 'heatmap' ? data.heatmap.days.length : 0;
    next = Math.max(0, Math.min(length - 1, next));
    this.activeHeatmapIndex.set(next);
    setTimeout(() => document.getElementById(this.heatmapCellId(next))?.focus());
  }
  stopPointer(event: PointerEvent): void {
    const target = event.target;
    if (
      target instanceof Element &&
      target.closest('button, a, input, select, textarea, svg, [data-no-drag]')
    )
      event.stopPropagation();
  }
  stopChildKeyboard(event: KeyboardEvent): void {
    if (event.target !== event.currentTarget) event.stopPropagation();
  }
}

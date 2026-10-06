import { DatePipe } from '@angular/common';
import { RouterLink } from '@angular/router';
import { Component, computed, input, output, signal } from '@angular/core';
import {
  DashboardMetricResult,
  DashboardWidgetViewState,
} from '../../../../api/dashboard/dashboard.service';
import {
  BarWidgetDefinition,
  HeatmapWidgetDefinition,
  LineWidgetDefinition,
  TableWidgetDefinition,
} from '../dashboard-definition';
import { LinePoint, lineCoordinates, linePath } from './line-chart-geometry';
import { defaultMetricProperties, formatMetric } from '../metric-widget-definition';

type AdvancedWidget =
  BarWidgetDefinition | HeatmapWidgetDefinition | LineWidgetDefinition | TableWidgetDefinition;

@Component({
  selector: 'ctx-advanced-widget',
  imports: [DatePipe, RouterLink],
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
            @if (config.properties.orientation === 'vertical') {
              <div
                class="flex min-h-0 flex-1 items-end gap-2 overflow-x-auto border-b border-border pb-1"
                role="list"
                aria-label="Categories"
              >
                @for (category of data.categories.categories; track category.key) {
                  <div
                    role="listitem"
                    class="flex h-full min-w-12 flex-col items-center justify-end gap-1 text-xs"
                    [attr.aria-label]="category.label + ': ' + format(category.value, data.unit)"
                  >
                    @if (config.properties.showValues) {
                      <span class="tabular-nums">{{ format(category.value, data.unit) }}</span>
                    }
                    <span class="flex min-h-1 w-6 flex-1 items-end rounded-sm bg-muted"
                      ><span
                        class="w-full rounded-sm bg-primary"
                        [style.height.%]="barWidth(category.value)"
                      ></span
                    ></span>
                    <a
                      class="max-w-20 truncate"
                      [title]="category.label"
                      [routerLink]="
                        category.entity
                          ? [
                              category.entity.type === 'project' ? '/project' : '/context',
                              category.entity.id,
                            ]
                          : null
                      "
                      >{{ category.label }}</a
                    >
                  </div>
                }
              </div>
            } @else {
              <div
                class="min-h-0 flex-1 space-y-2 overflow-auto"
                role="list"
                aria-label="Categories"
              >
                @for (category of data.categories.categories; track category.key) {
                  <div
                    role="listitem"
                    class="grid grid-cols-[minmax(5rem,1fr)_3fr_auto] items-center gap-2 text-xs"
                    [attr.aria-label]="category.label + ': ' + format(category.value, data.unit)"
                  >
                    <a
                      class="truncate"
                      [title]="category.label"
                      [routerLink]="
                        category.entity
                          ? [
                              category.entity.type === 'project' ? '/project' : '/context',
                              category.entity.id,
                            ]
                          : null
                      "
                      >{{ category.label }}</a
                    >
                    <span class="h-3 overflow-hidden rounded-sm bg-muted"
                      ><span
                        class="block h-full rounded-sm bg-primary"
                        [style.width.%]="barWidth(category.value)"
                      ></span
                    ></span>
                    @if (config.properties.showValues) {
                      <span class="tabular-nums">{{ format(category.value, data.unit) }}</span>
                    }
                  </div>
                }
              </div>
            }
            @if (config.properties.showSummary) {
              <p class="text-xs text-muted-foreground">
                {{ data.categories.totalCategories }} categories ·
                {{
                  format(
                    data.categories.totalValue,
                    data.unit,
                    config.properties.format,
                    config.properties.precision
                  )
                }}
                total
                @if (data.categories.omittedCategories) {
                  · {{ data.categories.omittedCategories }} omitted
                }
              </p>
            }
          }
          @if (data.categories.categories.length === 0) {
            <p role="status" class="text-sm text-muted-foreground">
              No activity in the selected period.
            </p>
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
            <svg
              class="min-h-0 min-w-0 flex-1 overflow-visible"
              viewBox="0 0 100 100"
              preserveAspectRatio="none"
              role="img"
              [attr.aria-label]="
                lineDescription(data.series.series.length, data.series.buckets.length)
              "
              tabindex="0"
              (keydown)="onLineKeydown($event, data.series.buckets.length)"
            >
              @if (config.properties.showGrid) {
                <path
                  d="M 0 25 H 100 M 0 50 H 100 M 0 75 H 100"
                  fill="none"
                  stroke="var(--border)"
                  stroke-dasharray="1 1"
                  vector-effect="non-scaling-stroke"
                />
              }
              @for (
                series of lineCoordinates(data.series.series, data.series.buckets.length);
                track series.key;
                let index = $index
              ) {
                @for (path of series.paths; track $index) {
                  <path
                    [attr.d]="path"
                    fill="none"
                    [attr.stroke]="seriesColor(index)"
                    stroke-width="1.5"
                    vector-effect="non-scaling-stroke"
                  />
                }
                @if (
                  config.properties.showPoints ||
                  config.properties.renderMode === 'points' ||
                  series.points.length === 1
                ) {
                  @for (point of series.points; track point.index) {
                    <circle
                      [attr.cx]="point.x"
                      [attr.cy]="point.y"
                      r="1.5"
                      [attr.fill]="seriesColor(index)"
                      [attr.aria-label]="
                        series.label +
                        ', ' +
                        data.series.buckets[point.index].start +
                        ': ' +
                        format(point.value, data.unit)
                      "
                    />
                  }
                }
              }
              @if (config.properties.targetValue !== null) {
                <line
                  x1="0"
                  x2="100"
                  [attr.y1]="targetY(config.properties.targetValue, data.series.series)"
                  [attr.y2]="targetY(config.properties.targetValue, data.series.series)"
                  stroke="var(--warning)"
                  stroke-dasharray="3 2"
                  vector-effect="non-scaling-stroke"
                />
              }
            </svg>
            @if (data.series.buckets.length) {
              <p class="text-xs text-muted-foreground" aria-live="polite">
                {{ selectedLinePoint() + 1 }} / {{ data.series.buckets.length }} ·
                {{
                  data.series.buckets[selectedLinePoint()]?.start
                    | date: 'mediumDate' : undefined : data.timeZone
                }}
                ·
                {{
                  linePointValue(
                    data,
                    selectedLinePoint(),
                    config.properties.format,
                    config.properties.precision
                  )
                }}
              </p>
            }
            @if (config.properties.showLegend) {
              <div class="flex flex-wrap gap-3 text-xs">
                @for (series of data.series.series; track series.key; let i = $index) {
                  <span class="inline-flex items-center gap-1"
                    ><span class="size-2 rounded-full" [style.background]="seriesColor(i)"></span
                    >{{ series.label }}</span
                  >
                }
              </div>
            }
            @if (config.properties.showSummary) {
              <p class="text-xs text-muted-foreground">
                {{
                  format(
                    data.series.totalValue,
                    data.unit,
                    config.properties.format,
                    config.properties.precision
                  )
                }}
                total
                @if (data.series.omittedSeries) {
                  · {{ data.series.omittedSeries }} series omitted
                }
              </p>
            }
            @if (data.series.series.length === 0 || data.series.buckets.length === 0) {
              <p role="status" class="text-sm text-muted-foreground">
                No activity in the selected period.
              </p>
            }
          }
        } @else if (data.kind === 'table') {
          @if (table(); as config) {
            <div class="min-h-0 flex-1 overflow-auto" [attr.aria-busy]="false">
              @if (data.table.totalRows === 0) {
                <p class="p-4 text-sm text-muted-foreground">No matching groups.</p>
              } @else {
                <table class="w-full border-collapse text-left text-xs">
                  <thead class="sticky top-0 bg-card">
                    <tr>
                      <th
                        scope="col"
                        aria-sort="{{
                          data.table.sort.columnId === 'group'
                            ? sortDirection(data.table.sort.direction)
                            : 'none'
                        }}"
                      >
                        <button type="button" class="font-medium" (click)="sortTable('group')">
                          {{ config.properties.groupBy === 'project' ? 'Project' : 'Context' }}
                        </button>
                      </th>
                      @for (column of visibleTableColumns(); track column.id) {
                        <th
                          scope="col"
                          [attr.aria-sort]="
                            column.id === data.table.sort.columnId
                              ? sortDirection(data.table.sort.direction)
                              : 'none'
                          "
                        >
                          <button type="button" class="font-medium" (click)="sortTable(column.id)">
                            {{ column.label || metricLabel(column.metric)
                            }}<span aria-hidden="true">{{
                              column.id === data.table.sort.columnId
                                ? data.table.sort.direction === 'asc'
                                  ? ' ↑'
                                  : ' ↓'
                                : ''
                            }}</span>
                          </button>
                        </th>
                      }
                    </tr>
                  </thead>
                  <tbody>
                    @for (row of data.table.rows; track row.key) {
                      <tr class="border-t border-border/60">
                        <th scope="row" class="py-1 pr-2 font-medium" [title]="row.label">
                          <a
                            [routerLink]="
                              row.entity
                                ? [
                                    row.entity.type === 'project' ? '/project' : '/context',
                                    row.entity.id,
                                  ]
                                : null
                            "
                            >{{ row.label }}</a
                          >
                        </th>
                        @for (column of visibleTableColumns(); track column.id) {
                          <td class="py-1 pl-2 text-right tabular-nums">
                            {{
                              format(
                                row.values[column.id] ?? null,
                                metricUnit(column.metric),
                                column.format,
                                column.precision
                              )
                            }}
                            @if (column.showBar && row.values[column.id] !== null) {
                              <span
                                class="ml-1 inline-block h-1 rounded-sm bg-primary align-middle"
                                [style.width.%]="
                                  tableBarWidth(row.values[column.id], column.metric, data)
                                "
                              ></span>
                            }
                          </td>
                        }
                      </tr>
                    }
                  </tbody>
                  @if (config.properties.showSummary) {
                    <tfoot>
                      <tr class="border-t border-border font-semibold">
                        <th scope="row" class="py-1 pr-2">Total</th>
                        @for (column of visibleTableColumns(); track column.id) {
                          <td class="py-1 pl-2 text-right tabular-nums">
                            {{
                              format(
                                data.table.summary[column.id] ?? null,
                                metricUnit(column.metric),
                                column.format,
                                column.precision
                              )
                            }}
                          </td>
                        }
                      </tr>
                    </tfoot>
                  }
                </table>
              }
            </div>
            <footer class="flex items-center justify-between gap-2 text-xs text-muted-foreground">
              <span aria-live="polite">
                @if (data.table.totalRows) {
                  {{ (data.table.page - 1) * data.table.pageSize + 1 }}–{{
                    pageEnd(data.table.page, data.table.pageSize, data.table.totalRows)
                  }}
                  of {{ data.table.totalRows }}
                }</span
              ><span class="flex gap-1"
                ><button
                  type="button"
                  class="rounded border border-border px-2 py-1 disabled:opacity-50"
                  [disabled]="data.table.page <= 1"
                  (click)="changePage(data, data.table.page - 1)"
                >
                  Previous
                </button>
                @if (data.table.pageCount) {
                  <span class="px-1 py-1">{{ data.table.page }} / {{ data.table.pageCount }}</span>
                }
                <button
                  type="button"
                  class="rounded border border-border px-2 py-1 disabled:opacity-50"
                  [disabled]="data.table.pageCount === 0 || data.table.page >= data.table.pageCount"
                  (click)="changePage(data, data.table.page + 1)"
                >
                  Next
                </button></span
              >
            </footer>
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
  readonly viewStateChange = output<DashboardWidgetViewState>();
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
  readonly selectedLinePoint = signal(0);
  readonly activeHeatmapCellId = computed(() => this.heatmapCellId(this.activeHeatmapIndex()));
  readonly maxCategory = computed(() => {
    const data = this.result();
    return data?.kind === 'categories'
      ? Math.max(0, ...data.categories.categories.map((item) => item.value ?? 0))
      : 0;
  });
  readonly maximumLineValue = computed(() => {
    const data = this.result();
    const target = this.line()?.properties.targetValue ?? 0;
    return data?.kind === 'series'
      ? Math.max(
          0,
          target,
          ...data.series.series.flatMap((item) =>
            item.values.filter((value): value is number => value !== null),
          ),
        )
      : 0;
  });
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

  barWidth(value: number | null): number {
    return value === null || this.maxCategory() <= 0 ? 0 : (value / this.maxCategory()) * 100;
  }
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
  lineCoordinates(
    series: { key: string; label: string; values: (number | null)[] }[],
    bucketCount: number,
  ): { key: string; label: string; points: LinePoint[]; paths: string[] }[] {
    const config = this.line();
    const data = this.result();
    if (!config || data?.kind !== 'series') return [];
    return series.map((item) => {
      const points = lineCoordinates(item.values, this.maximumLineValue() || 1);
      const future = data.series.buckets.map((bucket) => bucket.status === 'future');
      const visible = points.map((point, index) => (future[index] ? null : point));
      const paths =
        config.properties.renderMode === 'points'
          ? []
          : linePath(visible, config.properties.renderMode === 'smooth' ? 'smooth' : 'linear');
      return {
        key: item.key,
        label: item.label,
        points: visible.filter((point): point is LinePoint => point !== null),
        paths,
      };
    });
  }
  seriesColor(index: number): string {
    return [
      'var(--primary)',
      'var(--info)',
      'var(--accent)',
      'var(--foreground)',
      'var(--muted-foreground)',
    ][index % 5];
  }
  targetY(value: number, series: { values: (number | null)[] }[]): number {
    const maxValue = Math.max(this.maximumLineValue(), value, Number.EPSILON);
    return 100 - (value / maxValue) * 100;
  }
  lineDescription(series: number, buckets: number): string {
    return `${this.line()?.properties.renderMode ?? 'line'} chart with ${series} series and ${buckets} time buckets`;
  }
  onLineKeydown(event: KeyboardEvent, length: number): void {
    if (event.key === 'ArrowRight')
      this.selectedLinePoint.update((value) => Math.min(length - 1, value + 1));
    else if (event.key === 'ArrowLeft')
      this.selectedLinePoint.update((value) => Math.max(0, value - 1));
    else if (event.key === 'Home') this.selectedLinePoint.set(0);
    else if (event.key === 'End') this.selectedLinePoint.set(Math.max(0, length - 1));
    else return;
    event.preventDefault();
  }
  linePointValue(
    data: Extract<DashboardMetricResult, { kind: 'series' }>,
    index: number,
    format?: 'duration' | 'hours' | 'number' | 'percent',
    precision?: number,
  ): string {
    return data.series.series
      .map(
        (series) =>
          `${series.label}: ${this.format(series.values[index] ?? null, data.unit, format, precision)}`,
      )
      .join(' · ');
  }
  visibleTableColumns() {
    const config = this.table();
    return config?.properties.columns.filter((column) => column.visible) ?? [];
  }
  sortDirection(direction: 'asc' | 'desc'): 'ascending' | 'descending' {
    return direction === 'asc' ? 'ascending' : 'descending';
  }
  sortTable(columnId: string): void {
    const data = this.result();
    if (data?.kind !== 'table') return;
    const direction =
      data.table.sort.columnId === columnId && data.table.sort.direction === 'desc'
        ? 'asc'
        : 'desc';
    this.viewStateChange.emit({ page: 1, sort: { columnId, direction } });
  }
  changePage(data: Extract<DashboardMetricResult, { kind: 'table' }>, page: number): void {
    this.viewStateChange.emit({ page, sort: data.table.sort });
  }
  metricLabel(metric: string): string {
    return metric.replaceAll('.', ' ');
  }
  metricUnit(metric: string): 'ms' | 'count' | 'percent' {
    return metric.startsWith('duration.')
      ? 'ms'
      : metric === 'time_share.ratio'
        ? 'percent'
        : 'count';
  }
  tableBarWidth(
    value: number | null,
    metric: string,
    data: Extract<DashboardMetricResult, { kind: 'table' }>,
  ): number {
    if (value === null) return 0;
    const maxValue =
      metric === 'time_share.ratio'
        ? 100
        : (data.table.columnMaxima[
            this.table()?.properties.columns.find((column) => column.metric === metric)?.id ?? ''
          ] ?? 0);
    return maxValue <= 0 ? 0 : Math.min(100, (value / maxValue) * 100);
  }
  pageEnd(page: number, pageSize: number, total: number): number {
    return Math.min(page * pageSize, total);
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

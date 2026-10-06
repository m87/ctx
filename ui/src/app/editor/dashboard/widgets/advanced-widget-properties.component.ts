import { Component, computed, input, output } from '@angular/core';
import {
  DashboardWidgetDefinition,
  BarWidgetDefinition,
  HeatmapWidgetDefinition,
  LineWidgetDefinition,
  TableWidgetDefinition,
} from '../dashboard-definition';
import {
  DashboardPeriod,
  METRICS,
  PERIODS,
  MetricId,
  metricInfo,
} from '../metric-widget-definition';

type AdvancedWidget =
  BarWidgetDefinition | HeatmapWidgetDefinition | LineWidgetDefinition | TableWidgetDefinition;

@Component({
  selector: 'ctx-advanced-widget-properties',
  host: { class: 'block min-h-0 overflow-y-auto' },
  template: `
    <fieldset class="flex min-w-0 flex-col gap-4 p-3" [disabled]="disabled()">
      <legend class="sr-only">{{ widget().type }} settings</legend>
      <label class="ui-field"
        ><span class="ui-field-label">Title</span
        ><input
          class="h-9 rounded-md border border-border bg-background px-2 text-sm"
          [value]="widget().properties.title"
          (input)="setTitle($event)"
      /></label>
      @if (widget().type !== 'table') {
        <label class="ui-field"
          ><span class="ui-field-label">Metric</span
          ><select
            class="h-9 rounded-md border border-border bg-background px-2 text-sm"
            [value]="currentMetric()"
            (change)="setMetric($event)"
          >
            @for (metric of metrics; track metric.id) {
              <option [value]="metric.id">{{ metric.label }}</option>
            }
          </select></label
        >
      }
      @if (widget().type === 'heatmap') {
        <label class="ui-field"
          ><span class="ui-field-label">Period</span
          ><select
            class="h-9 rounded-md border border-border bg-background px-2 text-sm"
            [value]="widget().properties.period"
            (change)="setHeatmapPeriod($event)"
          >
            <option value="last_3_months">Last 3 months</option>
            <option value="last_6_months">Last 6 months</option>
            <option value="last_12_months">Last 12 months</option>
            @for (period of periods; track period.id) {
              <option [value]="period.id">{{ period.label }}</option>
            }
          </select></label
        >
      } @else if (widget().type === 'bar') {
        <label class="ui-field"
          ><span class="ui-field-label">Group by</span
          ><select
            class="h-9 rounded-md border border-border bg-background px-2 text-sm"
            [value]="currentGroup()"
            (change)="setBarGroup($event)"
          >
            <option value="project">Project</option>
            <option value="context">Context</option>
            <option value="time">Time</option>
          </select></label
        >
        @if (currentGroup() === 'time') {
          <label class="ui-field"
            ><span class="ui-field-label">Bucket</span
            ><select
              class="h-9 rounded-md border border-border bg-background px-2 text-sm"
              [value]="currentBucket()"
              (change)="setBarBucket($event)"
            >
              <option value="day">Day</option>
              <option value="week">Week</option>
              <option value="month">Month</option>
            </select></label
          >
        }
      } @else if (widget().type === 'line') {
        <label class="ui-field"
          ><span class="ui-field-label">Bucket</span
          ><select
            class="h-9 rounded-md border border-border bg-background px-2 text-sm"
            [value]="currentBucket()"
            (change)="setLineBucket($event)"
          >
            <option value="day">Day</option>
            <option value="week">Week</option>
            <option value="month">Month</option>
          </select></label
        >
        <label class="ui-field"
          ><span class="ui-field-label">Values</span
          ><select
            class="h-9 rounded-md border border-border bg-background px-2 text-sm"
            [value]="currentValueMode()"
            (change)="setLineMode($event)"
          >
            <option value="periodic">Per period</option>
            <option value="cumulative" [disabled]="currentMetric() !== 'duration.sum'">
              Cumulative
            </option>
          </select></label
        >
      } @else if (widget().type === 'table') {
        <label class="ui-field"
          ><span class="ui-field-label">Group by</span
          ><select
            class="h-9 rounded-md border border-border bg-background px-2 text-sm"
            [value]="currentGroup()"
            (change)="setTableGroup($event)"
          >
            <option value="project">Project</option>
            <option value="context">Context</option>
          </select></label
        >
      }
      @if (widget().type !== 'heatmap') {
        <label class="ui-field"
          ><span class="ui-field-label">Period</span
          ><select
            class="h-9 rounded-md border border-border bg-background px-2 text-sm"
            [value]="widget().properties.period"
            (change)="setPeriod($event)"
          >
            <option value="inherit">Dashboard period</option>
            @for (period of periods; track period.id) {
              <option [value]="period.id">{{ period.label }}</option>
            }
            <option value="custom">Custom</option>
          </select></label
        >
      }
      @if (bar(); as config) {
        @if (config.properties.groupBy !== 'time') {
          <label class="ui-field"
            ><span class="ui-field-label">Ranking</span
            ><select
              class="h-9 rounded-md border border-border bg-background px-2 text-sm"
              [value]="config.properties.sort.by + ':' + config.properties.sort.direction"
              (change)="setBarSort($event)"
            >
              <option value="value:desc">Value descending</option>
              <option value="value:asc">Value ascending</option>
              <option value="label:asc">Name ascending</option>
              <option value="label:desc">Name descending</option>
            </select></label
          >
          <label class="ui-field"
            ><span class="ui-field-label">Limit</span
            ><input
              type="number"
              min="2"
              max="50"
              class="h-9 rounded-md border border-border bg-background px-2 text-sm"
              [value]="config.properties.limit ?? 10"
              (change)="setBarLimit($event)"
          /></label>
          @if (config.properties.metric === 'duration.sum') {
            <label class="flex items-center gap-2 text-sm"
              ><input
                type="checkbox"
                [checked]="config.properties.showOther"
                (change)="setBarOther($event)"
              />Include Other</label
            >
          }
        }
        <label class="ui-field"
          ><span class="ui-field-label">Orientation</span
          ><select
            class="h-9 rounded-md border border-border bg-background px-2 text-sm"
            [value]="config.properties.orientation"
            (change)="setBarOrientation($event)"
          >
            <option value="horizontal">Horizontal</option>
            <option value="vertical">Vertical</option>
          </select></label
        >
        <label class="flex items-center gap-2 text-sm"
          ><input
            type="checkbox"
            [checked]="config.properties.showValues"
            (change)="setBarBoolean('showValues', $event)"
          />Show values</label
        >
        <label class="flex items-center gap-2 text-sm"
          ><input
            type="checkbox"
            [checked]="config.properties.showGrid"
            (change)="setBarGrid($event)"
          />Show grid</label
        >
        <label class="flex items-center gap-2 text-sm"
          ><input
            type="checkbox"
            [checked]="config.properties.showSummary"
            (change)="setBarBoolean('showSummary', $event)"
          />Show summary</label
        >
      }
      @if (heatmap(); as config) {
        <label class="ui-field"
          ><span class="ui-field-label">Scale</span
          ><select
            class="h-9 rounded-md border border-border bg-background px-2 text-sm"
            [value]="config.properties.scaleMode"
            (change)="setHeatmapScale($event)"
          >
            <option value="auto">Automatic</option>
            <option value="fixed">Fixed thresholds</option>
          </select></label
        >
        @if (config.properties.scaleMode === 'fixed') {
          @for (threshold of config.properties.thresholds; track $index) {
            <label class="ui-field"
              ><span class="ui-field-label">Threshold {{ $index + 1 }}</span
              ><input
                type="number"
                min="0"
                step="any"
                class="h-9 rounded-md border border-border bg-background px-2 text-sm"
                [value]="threshold"
                (change)="setHeatmapThreshold($index, $event)"
            /></label>
          }
        }
        <label class="flex items-center gap-2 text-sm"
          ><input
            type="checkbox"
            [checked]="config.properties.showMonthLabels"
            (change)="setHeatmapLabels('showMonthLabels', $event)"
          />Show month labels</label
        >
        <label class="flex items-center gap-2 text-sm"
          ><input
            type="checkbox"
            [checked]="config.properties.showDayLabels"
            (change)="setHeatmapLabels('showDayLabels', $event)"
          />Show day labels</label
        >
        <label class="flex items-center gap-2 text-sm"
          ><input
            type="checkbox"
            [checked]="config.properties.showLegend"
            (change)="setHeatmapBoolean('showLegend', $event)"
          />Show legend</label
        >
        <label class="flex items-center gap-2 text-sm"
          ><input
            type="checkbox"
            [checked]="config.properties.showSummary"
            (change)="setHeatmapBoolean('showSummary', $event)"
          />Show summary</label
        >
      }
      @if (line(); as config) {
        <label class="ui-field"
          ><span class="ui-field-label">Series</span
          ><select
            class="h-9 rounded-md border border-border bg-background px-2 text-sm"
            [value]="config.properties.seriesBy"
            (change)="setLineSeriesBy($event)"
          >
            <option value="none">Total</option>
            <option value="project">Project</option>
            <option value="context">Context</option>
          </select></label
        >
        @if (config.properties.seriesBy !== 'none') {
          <label class="ui-field"
            ><span class="ui-field-label">Series limit</span
            ><input
              type="number"
              min="1"
              max="5"
              class="h-9 rounded-md border border-border bg-background px-2 text-sm"
              [value]="config.properties.seriesLimit"
              (change)="setLineSeriesLimit($event)"
          /></label>
        }
        <label class="ui-field"
          ><span class="ui-field-label">Render style</span
          ><select
            class="h-9 rounded-md border border-border bg-background px-2 text-sm"
            [value]="config.properties.renderMode"
            (change)="setLineRenderMode($event)"
          >
            <option value="points">Points</option>
            <option value="linear">Straight lines</option>
            <option value="smooth">Smoothed lines</option>
          </select></label
        >
        <label class="ui-field"
          ><span class="ui-field-label">Target</span
          ><input
            type="number"
            min="0"
            step="any"
            class="h-9 rounded-md border border-border bg-background px-2 text-sm"
            [value]="config.properties.targetValue ?? ''"
            (change)="setLineTarget($event)"
        /></label>
        <label class="flex items-center gap-2 text-sm"
          ><input
            type="checkbox"
            [checked]="config.properties.showLegend"
            (change)="setLineBoolean('showLegend', $event)"
          />Show legend</label
        >
        <label class="flex items-center gap-2 text-sm"
          ><input
            type="checkbox"
            [checked]="config.properties.showSummary"
            (change)="setLineBoolean('showSummary', $event)"
          />Show summary</label
        >
      }
      @if (table(); as config) {
        <label class="ui-field"
          ><span class="ui-field-label">Sort direction</span
          ><select
            class="h-9 rounded-md border border-border bg-background px-2 text-sm"
            [value]="config.properties.sort.direction"
            (change)="setTableDirection($event)"
          >
            <option value="desc">Descending</option>
            <option value="asc">Ascending</option>
          </select></label
        >
        <label class="ui-field"
          ><span class="ui-field-label">Page size</span
          ><select
            class="h-9 rounded-md border border-border bg-background px-2 text-sm"
            [value]="config.properties.pageSize"
            (change)="setTablePageSize($event)"
          >
            <option value="10">10</option>
            <option value="25">25</option>
            <option value="50">50</option>
            <option value="100">100</option>
          </select></label
        >
        <section class="flex flex-col gap-2 border-t border-border/65 pt-3">
          <h4 class="text-sm font-medium">Columns</h4>
          @for (column of config.properties.columns; track column.id; let index = $index) {
            <div class="flex flex-col gap-2 rounded-md border border-border p-2">
              <select
                class="rounded-md border border-border bg-background p-1 text-xs"
                [value]="column.metric"
                (change)="setTableColumnMetric(index, $event)"
              >
                @for (metric of metrics; track metric.id) {
                  <option [value]="metric.id" [disabled]="tableMetricUsed(metric.id, index)">
                    {{ metric.label }}
                  </option>
                }
              </select>
              <input
                class="h-8 rounded-md border border-border bg-background px-2 text-xs"
                maxlength="100"
                aria-label="Column label"
                [value]="column.label"
                (input)="setTableColumnLabel(index, $event)"
              />
              <div class="flex flex-wrap gap-2">
                <button
                  type="button"
                  class="rounded border border-border px-2 py-1 text-xs"
                  [disabled]="index === 0"
                  (click)="moveTableColumn(index, -1)"
                >
                  Move up</button
                ><button
                  type="button"
                  class="rounded border border-border px-2 py-1 text-xs"
                  [disabled]="index === config.properties.columns.length - 1"
                  (click)="moveTableColumn(index, 1)"
                >
                  Move down</button
                ><button
                  type="button"
                  class="rounded border border-border px-2 py-1 text-xs"
                  [disabled]="config.properties.columns.length <= 1"
                  (click)="removeTableColumn(index)"
                >
                  Remove
                </button>
              </div>
              <label class="flex items-center gap-2 text-xs"
                ><input
                  type="checkbox"
                  [checked]="column.visible"
                  [disabled]="column.visible && visibleTableCount() <= 1"
                  (change)="setTableColumnVisible(index, $event)"
                />Visible</label
              >
              <label class="flex items-center gap-2 text-xs"
                ><input
                  type="checkbox"
                  [checked]="column.showBar"
                  [disabled]="
                    column.metric !== 'duration.sum' && column.metric !== 'time_share.ratio'
                  "
                  (change)="setTableColumnBar(index, $event)"
                />Show bar</label
              >
            </div>
          }
          <button
            type="button"
            class="rounded border border-border px-2 py-1 text-sm"
            [disabled]="config.properties.columns.length >= 7 || !availableTableMetric()"
            (click)="addTableColumn()"
          >
            Add column
          </button>
        </section>
      }
    </fieldset>
  `,
})
export class AdvancedWidgetPropertiesComponent {
  readonly widget = input.required<AdvancedWidget>();
  readonly disabled = input(false);
  readonly widgetChange = output<DashboardWidgetDefinition>();
  readonly metrics = METRICS;
  readonly periods = PERIODS;
  readonly currentMetric = computed(() => {
    const widget = this.widget();
    return 'metric' in widget.properties ? widget.properties.metric : '';
  });
  readonly currentGroup = computed(() => {
    const widget = this.widget();
    return 'groupBy' in widget.properties ? widget.properties.groupBy : '';
  });
  readonly currentBucket = computed(() => {
    const widget = this.widget();
    return widget.type === 'bar' || widget.type === 'line' ? (widget.properties.bucket ?? '') : '';
  });
  readonly currentValueMode = computed(() => {
    const widget = this.widget();
    return widget.type === 'line' ? widget.properties.valueMode : '';
  });
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

  setTitle(event: Event): void {
    const widget = this.widget();
    const title = (event.target as HTMLInputElement).value;
    if (widget.type === 'bar')
      this.widgetChange.emit({ ...widget, properties: { ...widget.properties, title } });
    else if (widget.type === 'heatmap')
      this.widgetChange.emit({ ...widget, properties: { ...widget.properties, title } });
    else if (widget.type === 'line')
      this.widgetChange.emit({ ...widget, properties: { ...widget.properties, title } });
    else this.widgetChange.emit({ ...widget, properties: { ...widget.properties, title } });
  }
  setMetric(event: Event): void {
    const id = (event.target as HTMLSelectElement).value;
    const metric = METRICS.find((item) => item.id === id)?.id;
    const widget = this.widget();
    if (!metric || widget.type === 'table') return;
    const info = metricInfo(metric);
    const format = info.unit === 'ms' ? 'duration' : info.unit === 'count' ? 'number' : 'percent';
    const precision = info.unit === 'count' ? 0 : widget.properties.precision;
    if (widget.type === 'bar')
      this.widgetChange.emit({
        ...widget,
        properties: {
          ...widget.properties,
          metric,
          format,
          precision,
          showOther: metric === 'duration.sum' ? widget.properties.showOther : false,
        },
      });
    else if (widget.type === 'heatmap')
      this.widgetChange.emit({
        ...widget,
        properties: {
          ...widget.properties,
          metric,
          format,
          precision,
          scaleMode: 'auto',
          thresholds: [],
        },
      });
    else
      this.widgetChange.emit({
        ...widget,
        properties: {
          ...widget.properties,
          metric,
          format,
          precision,
          targetValue: null,
          valueMode: metric === 'duration.sum' ? widget.properties.valueMode : 'periodic',
        },
      });
  }
  setHeatmapPeriod(event: Event): void {
    const value = (event.target as HTMLSelectElement).value;
    const widget = this.widget();
    if (
      widget.type === 'heatmap' &&
      [
        'last_3_months',
        'last_6_months',
        'last_12_months',
        ...this.periods.map((p) => p.id),
        'inherit',
        'custom',
      ].includes(value)
    )
      this.widgetChange.emit({
        ...widget,
        properties: { ...widget.properties, period: value as DashboardPeriod },
      });
  }
  setPeriod(event: Event): void {
    const value = (event.target as HTMLSelectElement).value;
    const widget = this.widget();
    if (!['inherit', 'custom', ...this.periods.map((p) => p.id)].includes(value)) return;
    if (widget.type === 'bar')
      this.widgetChange.emit({
        ...widget,
        properties: { ...widget.properties, period: value as DashboardPeriod },
      });
    else if (widget.type === 'line')
      this.widgetChange.emit({
        ...widget,
        properties: { ...widget.properties, period: value as DashboardPeriod },
      });
    else if (widget.type === 'table')
      this.widgetChange.emit({
        ...widget,
        properties: { ...widget.properties, period: value as DashboardPeriod },
      });
  }
  setBarGroup(event: Event): void {
    const value = (event.target as HTMLSelectElement).value;
    const widget = this.widget();
    if (widget.type === 'bar' && ['project', 'context', 'time'].includes(value))
      this.widgetChange.emit({
        ...widget,
        properties: {
          ...widget.properties,
          groupBy: value as 'project' | 'context' | 'time',
          bucket: value === 'time' ? (widget.properties.bucket ?? 'day') : null,
          sort:
            value === 'time'
              ? { by: 'time', direction: 'asc' }
              : widget.properties.sort.by === 'time'
                ? { by: 'value', direction: 'desc' }
                : widget.properties.sort,
          limit: value === 'time' ? null : (widget.properties.limit ?? 10),
          showOther: value === 'time' ? false : widget.properties.showOther,
        },
      });
  }
  setBarBucket(event: Event): void {
    const value = (event.target as HTMLSelectElement).value;
    const widget = this.widget();
    if (widget.type === 'bar' && ['day', 'week', 'month'].includes(value))
      this.widgetChange.emit({
        ...widget,
        properties: { ...widget.properties, bucket: value as 'day' | 'week' | 'month' },
      });
  }
  setLineBucket(event: Event): void {
    const value = (event.target as HTMLSelectElement).value;
    const widget = this.widget();
    if (widget.type === 'line' && ['day', 'week', 'month'].includes(value))
      this.widgetChange.emit({
        ...widget,
        properties: { ...widget.properties, bucket: value as 'day' | 'week' | 'month' },
      });
  }
  setLineMode(event: Event): void {
    const value = (event.target as HTMLSelectElement).value;
    const widget = this.widget();
    if (
      widget.type === 'line' &&
      (value === 'periodic' ||
        (value === 'cumulative' && widget.properties.metric === 'duration.sum'))
    )
      this.widgetChange.emit({ ...widget, properties: { ...widget.properties, valueMode: value } });
  }
  setTableGroup(event: Event): void {
    const value = (event.target as HTMLSelectElement).value;
    const widget = this.widget();
    if (widget.type === 'table' && (value === 'project' || value === 'context'))
      this.widgetChange.emit({ ...widget, properties: { ...widget.properties, groupBy: value } });
  }
  barSort(): string {
    const p = this.bar()?.properties.sort;
    return p ? `${p.by}:${p.direction}` : 'value:desc';
  }
  setBarSort(event: Event): void {
    const widget = this.widget(),
      value = (event.target as HTMLSelectElement).value;
    const [by, direction] = value.split(':');
    if (
      widget.type === 'bar' &&
      (by === 'value' || by === 'label') &&
      (direction === 'asc' || direction === 'desc')
    )
      this.widgetChange.emit({
        ...widget,
        properties: { ...widget.properties, sort: { by, direction } },
      });
  }
  setBarLimit(event: Event): void {
    const widget = this.widget(),
      value = Number((event.target as HTMLInputElement).value);
    if (widget.type === 'bar' && Number.isInteger(value) && value >= 2 && value <= 50)
      this.widgetChange.emit({ ...widget, properties: { ...widget.properties, limit: value } });
  }
  setBarOther(event: Event): void {
    const widget = this.widget();
    if (widget.type === 'bar')
      this.widgetChange.emit({
        ...widget,
        properties: { ...widget.properties, showOther: (event.target as HTMLInputElement).checked },
      });
  }
  setBarOrientation(event: Event): void {
    const widget = this.widget(),
      value = (event.target as HTMLSelectElement).value;
    if (widget.type === 'bar' && (value === 'horizontal' || value === 'vertical'))
      this.widgetChange.emit({
        ...widget,
        properties: { ...widget.properties, orientation: value },
      });
  }
  setBarBoolean(key: 'showValues' | 'showSummary', event: Event): void {
    const widget = this.widget();
    if (widget.type === 'bar')
      this.widgetChange.emit({
        ...widget,
        properties: { ...widget.properties, [key]: (event.target as HTMLInputElement).checked },
      });
  }
  setBarGrid(event: Event): void {
    const widget = this.widget();
    if (widget.type === 'bar')
      this.widgetChange.emit({
        ...widget,
        properties: { ...widget.properties, showGrid: (event.target as HTMLInputElement).checked },
      });
  }
  setHeatmapScale(event: Event): void {
    const widget = this.widget(),
      value = (event.target as HTMLSelectElement).value;
    if (widget.type === 'heatmap' && (value === 'auto' || value === 'fixed')) {
      const unit = metricInfo(widget.properties.metric).unit;
      const thresholds =
        unit === 'ms'
          ? [3600000, 7200000, 14400000]
          : unit === 'percent'
            ? [10, 25, 50]
            : [1, 2, 3];
      this.widgetChange.emit({
        ...widget,
        properties: {
          ...widget.properties,
          scaleMode: value,
          thresholds: value === 'fixed' ? thresholds : [],
        },
      });
    }
  }
  setHeatmapThreshold(index: number, event: Event): void {
    const widget = this.widget(),
      value = Number((event.target as HTMLInputElement).value);
    if (widget.type !== 'heatmap' || !Number.isFinite(value) || value <= 0) return;
    const thresholds = [...widget.properties.thresholds];
    thresholds[index] = value;
    if (thresholds.length === 3 && thresholds[0] < thresholds[1] && thresholds[1] < thresholds[2])
      this.widgetChange.emit({ ...widget, properties: { ...widget.properties, thresholds } });
  }
  setHeatmapBoolean(key: 'showLegend' | 'showSummary', event: Event): void {
    const widget = this.widget();
    if (widget.type === 'heatmap')
      this.widgetChange.emit({
        ...widget,
        properties: { ...widget.properties, [key]: (event.target as HTMLInputElement).checked },
      });
  }
  setHeatmapLabels(key: 'showMonthLabels' | 'showDayLabels', event: Event): void {
    const widget = this.widget();
    if (widget.type === 'heatmap')
      this.widgetChange.emit({
        ...widget,
        properties: { ...widget.properties, [key]: (event.target as HTMLInputElement).checked },
      });
  }
  setLineSeriesBy(event: Event): void {
    const widget = this.widget(),
      value = (event.target as HTMLSelectElement).value;
    if (widget.type === 'line' && (value === 'none' || value === 'project' || value === 'context'))
      this.widgetChange.emit({ ...widget, properties: { ...widget.properties, seriesBy: value } });
  }
  setLineSeriesLimit(event: Event): void {
    const widget = this.widget(),
      value = Number((event.target as HTMLInputElement).value);
    if (widget.type === 'line' && Number.isInteger(value) && value >= 1 && value <= 5)
      this.widgetChange.emit({
        ...widget,
        properties: { ...widget.properties, seriesLimit: value },
      });
  }
  setLineRenderMode(event: Event): void {
    const widget = this.widget(),
      value = (event.target as HTMLSelectElement).value;
    if (widget.type === 'line' && (value === 'points' || value === 'linear' || value === 'smooth'))
      this.widgetChange.emit({
        ...widget,
        properties: { ...widget.properties, renderMode: value },
      });
  }
  setLineTarget(event: Event): void {
    const widget = this.widget(),
      raw = (event.target as HTMLInputElement).value,
      value = raw === '' ? null : Number(raw);
    if (widget.type === 'line' && (value === null || (Number.isFinite(value) && value >= 0)))
      this.widgetChange.emit({
        ...widget,
        properties: { ...widget.properties, targetValue: value },
      });
  }
  setLineBoolean(key: 'showLegend' | 'showSummary', event: Event): void {
    const widget = this.widget();
    if (widget.type === 'line')
      this.widgetChange.emit({
        ...widget,
        properties: { ...widget.properties, [key]: (event.target as HTMLInputElement).checked },
      });
  }
  setTableDirection(event: Event): void {
    const widget = this.widget(),
      direction = (event.target as HTMLSelectElement).value;
    if (widget.type === 'table' && (direction === 'asc' || direction === 'desc'))
      this.widgetChange.emit({
        ...widget,
        properties: { ...widget.properties, sort: { ...widget.properties.sort, direction } },
      });
  }
  setTablePageSize(event: Event): void {
    const widget = this.widget(),
      pageSize = Number((event.target as HTMLSelectElement).value);
    if (
      widget.type === 'table' &&
      (pageSize === 10 || pageSize === 25 || pageSize === 50 || pageSize === 100)
    )
      this.widgetChange.emit({ ...widget, properties: { ...widget.properties, pageSize } });
  }
  tableMetricUsed(metric: MetricId, index: number): boolean {
    const columns = this.table()?.properties.columns ?? [];
    return columns.some((column, i) => i !== index && column.metric === metric);
  }
  visibleTableCount(): number {
    return this.table()?.properties.columns.filter((column) => column.visible).length ?? 0;
  }
  availableTableMetric(): MetricId | undefined {
    return METRICS.find(
      (metric) => !this.table()?.properties.columns.some((column) => column.metric === metric.id),
    )?.id;
  }
  setTableColumnMetric(index: number, event: Event): void {
    const widget = this.widget(),
      id = (event.target as HTMLSelectElement).value,
      metric = METRICS.find((item) => item.id === id)?.id;
    if (widget.type !== 'table' || !metric || this.tableMetricUsed(metric, index)) return;
    const info = metricInfo(metric),
      format: TableWidgetDefinition['properties']['columns'][number]['format'] =
        info.unit === 'ms' ? 'duration' : info.unit === 'count' ? 'number' : 'percent',
      columns = widget.properties.columns.map((column, i) =>
        i === index
          ? {
              ...column,
              metric,
              format,
              precision: info.unit === 'count' ? 0 : column.precision,
              showBar:
                column.showBar && (metric === 'duration.sum' || metric === 'time_share.ratio'),
            }
          : column,
      );
    this.widgetChange.emit({ ...widget, properties: { ...widget.properties, columns } });
  }
  setTableColumnLabel(index: number, event: Event): void {
    const widget = this.widget(),
      label = (event.target as HTMLInputElement).value;
    if (widget.type === 'table')
      this.widgetChange.emit({
        ...widget,
        properties: {
          ...widget.properties,
          columns: widget.properties.columns.map((column, i) =>
            i === index ? { ...column, label } : column,
          ),
        },
      });
  }
  setTableColumnVisible(index: number, event: Event): void {
    const widget = this.widget(),
      visible = (event.target as HTMLInputElement).checked;
    if (
      widget.type !== 'table' ||
      (visible === false &&
        widget.properties.columns[index].visible &&
        this.visibleTableCount() <= 1)
    )
      return;
    if (widget.type === 'table') {
      const columns = widget.properties.columns.map((column, i) =>
        i === index ? { ...column, visible } : column,
      );
      const sort = columns.some(
        (column) => column.id === widget.properties.sort.columnId && column.visible,
      )
        ? widget.properties.sort
        : {
            columnId: columns.find((column) => column.visible)?.id ?? 'group',
            direction: 'desc' as const,
          };
      this.widgetChange.emit({ ...widget, properties: { ...widget.properties, columns, sort } });
    }
  }
  setTableColumnBar(index: number, event: Event): void {
    const widget = this.widget(),
      showBar = (event.target as HTMLInputElement).checked;
    if (widget.type === 'table')
      this.widgetChange.emit({
        ...widget,
        properties: {
          ...widget.properties,
          columns: widget.properties.columns.map((column, i) =>
            i === index ? { ...column, showBar } : column,
          ),
        },
      });
  }
  moveTableColumn(index: number, delta: -1 | 1): void {
    const widget = this.widget();
    if (widget.type !== 'table') return;
    const target = index + delta;
    if (target < 0 || target >= widget.properties.columns.length) return;
    const columns = [...widget.properties.columns];
    [columns[index], columns[target]] = [columns[target], columns[index]];
    this.widgetChange.emit({ ...widget, properties: { ...widget.properties, columns } });
  }
  removeTableColumn(index: number): void {
    const widget = this.widget();
    if (widget.type !== 'table' || widget.properties.columns.length <= 1) return;
    const columns = widget.properties.columns.filter((_, i) => i !== index);
    const sort = columns.some(
      (column) => column.id === widget.properties.sort.columnId && column.visible,
    )
      ? widget.properties.sort
      : {
          columnId: columns.find((column) => column.visible)?.id ?? 'group',
          direction: 'desc' as const,
        };
    this.widgetChange.emit({ ...widget, properties: { ...widget.properties, columns, sort } });
  }
  addTableColumn(): void {
    const widget = this.widget(),
      metric = this.availableTableMetric();
    if (widget.type !== 'table' || !metric || widget.properties.columns.length >= 7) return;
    const info = metricInfo(metric),
      id = metric.replaceAll('.', '-');
    const column = {
      id,
      metric,
      label: info.label,
      visible: true,
      format:
        info.unit === 'ms'
          ? ('duration' as const)
          : info.unit === 'count'
            ? ('number' as const)
            : ('percent' as const),
      precision: info.unit === 'count' ? 0 : 1,
      showBar: false,
    };
    this.widgetChange.emit({
      ...widget,
      properties: { ...widget.properties, columns: [...widget.properties.columns, column] },
    });
  }
  metricLabel(metric: MetricId): string {
    return metricInfo(metric).label;
  }
}

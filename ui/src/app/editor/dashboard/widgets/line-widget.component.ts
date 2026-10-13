import { DatePipe } from '@angular/common';
import { Component, computed, input, signal } from '@angular/core';
import { ChartOptions } from '@tanstack/charts/angular';
import { ChartPoint } from '@tanstack/charts';
import { DashboardMetricResult } from '../../../../api/dashboard/dashboard.service';
import { ChartFrameComponent } from './chart-frame.component';
import { formatResultValue } from './format-result-value';
import { LineWidgetDefinition } from '../dashboard-definition';
import { lineChartRows } from './line-chart-data';
import { bucketLabelFormatter, lineChartDefinition, seriesColor } from './line-chart-definition';

type SeriesResult = Extract<DashboardMetricResult, { kind: 'series' }>;

@Component({
  selector: 'ctx-line-widget',
  imports: [ChartFrameComponent, DatePipe],
  host: { class: 'flex min-h-0 flex-1 flex-col gap-3' },
  template: `
    <ctx-chart-frame [chart]="options()" />
    <p class="min-h-4 text-xs text-muted-foreground" aria-live="polite">
      @if (focusedBucket(); as index) {
        {{ result().series.buckets[index.value]?.start | date: 'mediumDate' : result().timeZone }}
        · {{ pointValue(index.value) }}
      }
    </p>
    @if (widget().properties.showLegend) {
      <div class="flex flex-wrap gap-3 text-xs">
        @for (series of result().series.series; track series.key; let i = $index) {
          <span class="inline-flex items-center gap-1"
            ><span class="size-2 rounded-full" [style.background]="seriesColor(i)"></span
            >{{ series.label }}</span
          >
        }
      </div>
    }
    @if (widget().properties.showSummary) {
      <p class="text-xs text-muted-foreground">
        {{ format(result().series.totalValue) }}
        total
        @if (result().series.omittedSeries) {
          · {{ result().series.omittedSeries }} series omitted
        }
      </p>
    }
    @if (result().series.series.length === 0 || result().series.buckets.length === 0) {
      <p role="status" class="text-sm text-muted-foreground">No activity in the selected period.</p>
    }
  `,
})
export class LineWidgetComponent {
  readonly widget = input.required<LineWidgetDefinition>();
  readonly result = input.required<SeriesResult>();
  readonly editing = input(false);
  readonly focused = signal<ChartPoint | null>(null);
  readonly seriesColor = seriesColor;
  readonly focusedBucket = computed(() => {
    const point = this.focused();
    if (!point) return null;
    const index = this.result().series.buckets.findIndex((bucket) => bucket.start === point.xValue);
    return index < 0 ? null : { value: index };
  });
  readonly description = computed(() => {
    const data = this.result();
    return `${this.widget().properties.renderMode} chart with ${data.series.series.length} series and ${data.series.buckets.length} time buckets`;
  });
  readonly options = computed<ChartOptions>(() => {
    const data = this.result();
    const properties = this.widget().properties;
    return {
      definition: lineChartDefinition(lineChartRows(data, properties), properties, {
        value: (value) => this.format(value),
        bucket: bucketLabelFormatter(properties.bucket, data.timeZone),
      }),
      ariaLabel: this.description(),
      tabIndex: this.editing() ? -1 : 0,
      onFocusChange: (point) => this.focused.set(point),
    };
  });

  format(value: number | null): string {
    return formatResultValue(value, this.result().unit, this.widget().properties);
  }
  pointValue(index: number): string {
    return this.result()
      .series.series.map(
        (series) => `${series.label}: ${this.format(series.values[index] ?? null)}`,
      )
      .join(' · ');
  }
}

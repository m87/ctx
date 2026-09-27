import { Component, computed, input } from '@angular/core';
import { DashboardMetricResult } from '../../../../api/dashboard/dashboard.service';
import { MetricWidgetProperties, formatMetric, metricInfo } from '../metric-widget-definition';

@Component({
  selector: 'ctx-kpi-widget',
  host: { class: 'dashboard-metric-card' },
  template: `
    <div class="min-w-0">
      <div
        class="truncate text-meta font-semibold uppercase tracking-label text-muted-foreground"
        [title]="title()"
      >
        {{ title() }}
      </div>
      @if (properties().description) {
        <p class="mt-1 line-clamp-2 text-xs text-muted-foreground">
          {{ properties().description }}
        </p>
      }
    </div>
    <div class="flex min-h-0 flex-1 items-center">
      @if (error() || result()?.error) {
        <p class="text-xs text-destructive" role="status">{{ error() || result()?.error }}</p>
      } @else if (!result()) {
        <span class="text-xs text-muted-foreground" role="status">Loading…</span>
      } @else {
        <div
          class="dashboard-metric-value"
          [class.text-muted-foreground]="result()?.value === null"
        >
          {{ value() }}
        </div>
      }
    </div>
    <div class="mt-auto border-t border-border/65 pt-2">
      <div class="truncate text-meta text-muted-foreground" [title]="period()">{{ period() }}</div>
    </div>
    @if (result()?.ignoredIntervals) {
      <p class="text-caption text-warning">Some invalid sessions were excluded.</p>
    }
  `,
})
export class KpiWidgetComponent {
  readonly properties = input.required<MetricWidgetProperties>();
  readonly result = input<DashboardMetricResult>();
  readonly error = input('');
  readonly title = computed(
    () => this.properties().title || metricInfo(this.properties().metric).label,
  );
  readonly value = computed(() => formatMetric(this.result()?.value ?? null, this.properties()));
  readonly period = input('');
}

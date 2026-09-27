import { Component, computed, input } from '@angular/core';
import { DashboardMetricResult } from '../../../../api/dashboard/dashboard.service';
import {
  GaugeWidgetProperties,
  formatMetric,
  gaugeTone,
  metricInfo,
} from '../metric-widget-definition';

@Component({
  selector: 'ctx-gauge-widget',
  host: { class: 'dashboard-metric-card' },
  template: `
    <div class="min-w-0">
      <div class="truncate text-sm font-semibold" [title]="title()">{{ title() }}</div>
      @if (properties().description) {
        <p class="mt-1 truncate text-meta text-muted-foreground">{{ properties().description }}</p>
      }
      <p class="mt-1 truncate text-meta text-muted-foreground">{{ period() }}</p>
    </div>
    @if (error() || result()?.error) {
      <p class="flex flex-1 items-center text-xs text-destructive" role="status">
        {{ error() || result()?.error }}
      </p>
    } @else if (!result()) {
      <p
        class="flex flex-1 items-center justify-center text-sm text-muted-foreground"
        role="status"
      >
        Loading…
      </p>
    } @else {
      <div
        class="dashboard-gauge-chart"
        role="img"
        [attr.aria-label]="accessibleValue()"
        [attr.data-tone]="tone()"
      >
        <svg viewBox="0 0 240 140" aria-hidden="true" class="dashboard-gauge-svg">
          <path
            d="M 20 120 A 100 100 0 0 1 220 120"
            pathLength="100"
            class="dashboard-gauge-track"
          />
          @if (progress() > 0) {
            <path
              d="M 20 120 A 100 100 0 0 1 220 120"
              pathLength="100"
              class="dashboard-gauge-progress"
              [attr.stroke-dasharray]="progress() + ' 100'"
            />
          }
        </svg>
        <div class="dashboard-gauge-value">
          @if (properties().showValue || result()?.value === null) {
            <div
              class="dashboard-metric-value"
              [class.text-muted-foreground]="result()?.value === null"
            >
              {{ value() }}
            </div>
          }
          @if (properties().showPercentage && percentage() !== null) {
            <p class="mt-1 text-meta text-muted-foreground">
              {{ percentage() }}% of {{ properties().goalMode === 'at_most' ? 'limit' : 'target' }}
            </p>
          }
        </div>
      </div>
      <div class="flex items-end justify-between gap-2 text-caption text-muted-foreground">
        <span>{{ format(properties().minimum) }}</span>
        @if (properties().showTarget && properties().target !== null) {
          <span class="text-center"
            >{{ properties().goalMode === 'at_most' ? 'Limit' : 'Target' }}
            <strong class="font-medium text-foreground">{{
              format(properties().target)
            }}</strong></span
          >
        }
        <span>{{ format(properties().maximum) }}</span>
      </div>
      @if (goalStatus()) {
        <p class="text-meta text-muted-foreground" role="status">{{ goalStatus() }}</p>
      }
      @if (result()?.ignoredIntervals) {
        <p class="text-caption text-warning">Some invalid sessions were excluded.</p>
      }
    }
  `,
})
export class GaugeWidgetComponent {
  readonly properties = input.required<GaugeWidgetProperties>();
  readonly result = input<DashboardMetricResult>();
  readonly error = input('');
  readonly period = input('');
  readonly title = computed(
    () => this.properties().title || metricInfo(this.properties().metric).label,
  );
  readonly value = computed(() => this.format(this.result()?.value ?? null));
  readonly progress = computed(() => {
    const p = this.properties(),
      value = this.result()?.value;
    return value == null
      ? 0
      : Math.max(0, Math.min(100, ((value - p.minimum) / (p.maximum - p.minimum)) * 100));
  });
  readonly tone = computed(() => gaugeTone(this.result()?.value ?? null, this.properties()));
  readonly percentage = computed(() => {
    const value = this.result()?.value,
      target = this.properties().target;
    return value == null || !target
      ? null
      : new Intl.NumberFormat(undefined, { maximumFractionDigits: 1 }).format(
          (value / target) * 100,
        );
  });
  readonly accessibleValue = computed(
    () =>
      this.title() +
      ': ' +
      this.value() +
      (this.properties().target !== null
        ? ', ' +
          (this.properties().goalMode === 'at_most' ? 'limit ' : 'target ') +
          this.format(this.properties().target)
        : '') +
      ', ' +
      this.period(),
  );
  readonly goalStatus = computed(() => {
    const value = this.result()?.value;
    const target = this.properties().target;
    if (value == null || target === null || this.properties().goalMode === 'informational')
      return '';
    if (this.properties().goalMode === 'at_most')
      return value > target ? 'Limit exceeded' : 'Within limit';
    return value >= target ? 'Target reached' : 'Target not reached';
  });
  format(value: number | null): string {
    return formatMetric(value, this.properties());
  }
}

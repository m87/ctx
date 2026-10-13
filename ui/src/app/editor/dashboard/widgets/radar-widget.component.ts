import { Component, computed, input, signal } from '@angular/core';
import { ChartOptions } from '@tanstack/charts/angular';
import {
  DashboardMetricResult,
  DashboardRadarAxis,
} from '../../../../api/dashboard/dashboard.service';
import { RadarWidgetProperties, formatMetric, metricInfo } from '../metric-widget-definition';
import { ChartFrameComponent } from './chart-frame.component';
import { radarChartDefinition, radarEffectiveMaximum } from './radar-chart-definition';

@Component({
  selector: 'ctx-radar-widget',
  imports: [ChartFrameComponent],
  host: { class: 'dashboard-metric-card' },
  template: `
    <header class="min-w-0">
      <h3 class="truncate text-sm font-semibold" [title]="title()">{{ title() }}</h3>
      @if (properties().description) {
        <p class="mt-1 truncate text-meta text-muted-foreground">{{ properties().description }}</p>
      }
      <p class="mt-1 truncate text-meta text-muted-foreground">
        {{ period() }} · {{ metricLabel() }}
      </p>
    </header>
    @if (properties().axesBy !== 'weekday' && properties().axisKeys.length < 3) {
      <p
        class="flex flex-1 items-center justify-center text-xs text-muted-foreground"
        role="status"
      >
        Choose at least 3 axes to preview the radar chart.
      </p>
    } @else if (error() || result()?.error) {
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
    } @else if (noActivity()) {
      <p
        class="flex flex-1 items-center justify-center text-xs text-muted-foreground"
        role="status"
      >
        No activity in the selected period.
      </p>
    } @else if (!radar()?.axes?.length || allMissing()) {
      <p
        class="flex flex-1 items-center justify-center text-xs text-muted-foreground"
        role="status"
      >
        No data for this period.
      </p>
    } @else {
      <div class="flex min-h-0 flex-1 flex-col gap-1">
        <ctx-chart-frame class="max-h-96" [chart]="options()" />
        @if (activeAxis(); as axis) {
          <p class="text-center text-caption text-muted-foreground" aria-live="polite">
            <strong class="font-medium text-foreground">{{ axis.label }}</strong
            >: {{ format(axis.value) }}
            @if (axis.observedDays !== null) {
              · {{ axis.observedDays }} observed days
            }
            @if (axis.status === 'unobserved') {
              · not observed
            } @else if (axis.status === 'partial') {
              · partial period
            }
          </p>
        }
      </div>
      <div
        class="flex flex-wrap items-center justify-between gap-2 border-t border-border/65 pt-2 text-caption text-muted-foreground"
      >
        <span class="flex items-center gap-1.5"
          ><i class="dashboard-radar-legend"></i>{{ metricLabel() }}</span
        >
        @if (properties().targetValue !== null) {
          <span>Target: {{ format(properties().targetValue) }} per axis</span>
        }
        <span>Scale 0–{{ format(effectiveMaximum()) }}</span>
      </div>
      @if (properties().showSummary) {
        <p class="text-caption text-muted-foreground">
          {{ summaryLabel() }}:
          <span class="tabular-nums text-foreground">{{
            format(radar()?.totalValue ?? null)
          }}</span>
        </p>
      }
      @if (incompleteCount()) {
        <p class="text-caption text-warning">Incomplete data on {{ incompleteCount() }} axes.</p>
      }
      @if (extendedScale()) {
        <p class="text-caption text-warning">Scale extended to {{ format(effectiveMaximum()) }}.</p>
      }
      <details class="text-caption">
        @if (!editing()) {
          <summary
            class="cursor-pointer rounded-md text-muted-foreground focus-visible:ring-2 focus-visible:ring-ring"
          >
            Data table
          </summary>
          <table class="mt-2 w-full text-left">
            <caption class="sr-only">
              {{
                title()
              }}
              values by axis
            </caption>
            <thead>
              <tr>
                <th scope="col">Category</th>
                <th scope="col">Value</th>
                <th scope="col">Status</th>
              </tr>
            </thead>
            <tbody>
              @for (axis of axes(); track axis.key) {
                <tr>
                  <th scope="row" class="py-1 font-medium">{{ axis.label }}</th>
                  <td class="py-1 tabular-nums">{{ format(axis.value) }}</td>
                  <td class="py-1 text-muted-foreground">{{ axis.status }}</td>
                </tr>
              }
            </tbody>
          </table>
        }
      </details>
    }
  `,
})
export class RadarWidgetComponent {
  readonly properties = input.required<RadarWidgetProperties>();
  readonly result = input<DashboardMetricResult>();
  readonly error = input('');
  readonly period = input('');
  readonly editing = input(false);
  readonly activeIndex = signal(0);
  readonly radar = computed(() => {
    const result = this.result();
    return result?.kind === 'radar' ? result.radar : undefined;
  });
  readonly noActivity = computed(() => {
    const result = this.result();
    return result?.kind === 'radar' && !result.hasActivity;
  });
  readonly allMissing = computed(() =>
    (this.radar()?.axes ?? []).every((axis) => axis.value === null),
  );
  readonly axes = computed(() => {
    const axes = this.radar()?.axes ?? [];
    return this.properties().axesBy === 'weekday'
      ? axes
      : this.properties()
          .axisKeys.map((key) => axes.find((axis) => axis.key === key))
          .filter((axis): axis is DashboardRadarAxis => !!axis);
  });
  readonly title = computed(
    () => this.properties().title || metricInfo(this.properties().metric).label,
  );
  readonly metricLabel = computed(() => metricInfo(this.properties().metric).label);
  readonly effectiveMaximum = computed(() => {
    const p = this.properties();
    return radarEffectiveMaximum(
      this.axes().map((axis) => axis.value),
      p.scaleMaximum,
      this.radar()?.unit ?? metricInfo(p.metric).unit,
      p.targetValue,
    );
  });
  readonly extendedScale = computed(
    () =>
      this.properties().scaleMaximum !== null &&
      this.effectiveMaximum() > (this.properties().scaleMaximum ?? 0),
  );
  readonly incompleteCount = computed(
    () => this.axes().filter((axis) => axis.value === null).length,
  );
  readonly activeAxis = computed(() => this.axes()[this.activeIndex()] ?? null);
  readonly summaryLabel = computed(() =>
    this.properties().axesBy === 'weekday' ? 'Period total' : 'Selected categories',
  );
  readonly options = computed<ChartOptions>(() => ({
    definition: radarChartDefinition(this.axes(), this.properties(), this.effectiveMaximum(), {
      value: (value) => this.format(value),
      metric: this.metricLabel(),
    }),
    ariaLabel: `${this.title()}, ${this.scaleDescription()}`,
    ariaDescription: 'Use arrow keys to inspect axes.',
    tabIndex: this.editing() ? -1 : 0,
    onFocusChange: (point) => {
      const index = this.axes().findIndex((axis) => axis.key === point?.xValue);
      if (index >= 0) this.activeIndex.set(index);
    },
  }));
  readonly scaleDescription = computed(
    () =>
      `Shared scale from zero to ${this.format(this.effectiveMaximum())}. Values use radial distance.`,
  );

  format(value: number | null) {
    return value === null ? '—' : formatMetric(value, this.properties());
  }
}

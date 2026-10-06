import { Component, computed, input, signal } from '@angular/core';
import {
  DashboardMetricResult,
  DashboardRadarAxis,
} from '../../../../api/dashboard/dashboard.service';
import { RadarWidgetProperties, formatMetric, metricInfo } from '../metric-widget-definition';
import {
  radarClosedPolygon,
  radarEffectiveMaximum,
  radarPoint,
  radarSegments,
} from './radar-chart-geometry';

@Component({
  selector: 'ctx-radar-widget',
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
      <div class="flex min-h-0 flex-1 flex-col items-center justify-center gap-1">
        <svg
          viewBox="0 0 300 260"
          role="img"
          [attr.aria-label]="title() + ', ' + scaleDescription() + ', ' + activeAxisLabel()"
          [attr.tabindex]="editing() ? -1 : 0"
          class="dashboard-radar-svg"
          (focus)="selectAxis(0)"
          (keydown)="moveAxis($event)"
        >
          <desc>{{ title() }}. {{ scaleDescription() }}. Use arrow keys to inspect axes.</desc>
          @if (properties().showGrid) {
            @for (ring of rings; track ring) {
              <polygon [attr.points]="ringPoints(ring)" class="dashboard-radar-grid" />
            }
            @for (axis of axes(); track axis.key; let i = $index) {
              <path [attr.d]="axisLine(i)" class="dashboard-radar-grid" />
            }
          }
          @if (targetPath()) {
            <path [attr.d]="targetPath()" class="dashboard-radar-target" />
          }
          @if (fillPath()) {
            <path [attr.d]="fillPath()" class="dashboard-radar-fill" />
          }
          @for (segment of segments(); track $index) {
            <path [attr.d]="segment" class="dashboard-radar-series" />
          }
          @for (axis of axes(); track axis.key; let i = $index) {
            @if (axis.value !== null) {
              <circle
                [attr.cx]="point(i, axis.value).x"
                [attr.cy]="point(i, axis.value).y"
                r="3"
                class="dashboard-radar-dot"
                (mouseenter)="selectAxis(i)"
                (click)="selectAxis(i)"
              />
            }
            <text
              [attr.x]="labelPoint(i).x"
              [attr.y]="labelPoint(i).y"
              text-anchor="middle"
              class="dashboard-radar-label"
              (mouseenter)="selectAxis(i)"
              (click)="selectAxis(i)"
            >
              {{ shortLabel(axis.label) }}
              <title>{{ axis.label }}: {{ format(axis.value) }}</title>
            </text>
            @if (properties().showValues && axis.value !== null) {
              <text
                [attr.x]="labelPoint(i).x"
                [attr.y]="labelPoint(i).y + 12"
                text-anchor="middle"
                class="dashboard-radar-value-label"
              >
                {{ format(axis.value) }}
              </text>
            }
          }
          @if (activeAxis(); as axis) {
            <g class="dashboard-radar-focus" aria-hidden="true">
              <circle [attr.cx]="focusPoint().x" [attr.cy]="focusPoint().y" r="9" />
            </g>
          }
        </svg>
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
  readonly rings = [0.25, 0.5, 0.75, 1];
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
  readonly fillPath = computed(() =>
    this.properties().showFill
      ? radarClosedPolygon(
          this.axes().map((axis) => axis.value),
          this.effectiveMaximum(),
        )
      : null,
  );
  readonly targetPath = computed(() => {
    const target = this.properties().targetValue;
    return target === null
      ? null
      : radarClosedPolygon(
          this.axes().map(() => target),
          this.effectiveMaximum(),
        );
  });
  readonly segments = computed(() =>
    radarSegments(
      this.axes().map((axis) => axis.value),
      this.effectiveMaximum(),
    ),
  );
  readonly incompleteCount = computed(
    () => this.axes().filter((axis) => axis.value === null).length,
  );
  readonly activeAxis = computed(() => this.axes()[this.activeIndex()] ?? null);
  readonly focusPoint = computed(() => {
    const axis = this.activeAxis();
    return axis ? this.point(this.activeIndex(), axis.value ?? 0) : { x: 150, y: 130 };
  });
  readonly activeAxisLabel = computed(() =>
    this.activeAxis()
      ? `${this.activeAxis()?.label}: ${this.format(this.activeAxis()?.value ?? null)}`
      : 'No radar axes',
  );
  readonly summaryLabel = computed(() =>
    this.properties().axesBy === 'weekday' ? 'Period total' : 'Selected categories',
  );
  readonly scaleDescription = computed(
    () =>
      `Shared scale from zero to ${this.format(this.effectiveMaximum())}. Values use radial distance.`,
  );

  point(index: number, value: number | null) {
    return radarPoint(
      index,
      Math.max(1, this.axes().length),
      Math.max(0, (value ?? 0) / this.effectiveMaximum()),
    );
  }
  labelPoint(index: number) {
    return radarPoint(index, Math.max(1, this.axes().length), 1.18);
  }
  axisLine(index: number) {
    const point = radarPoint(index, this.axes().length, 1);
    return `M 150 130 L ${point.x} ${point.y}`;
  }
  ringPoints(ratio: number) {
    return this.axes()
      .map((_, i) => {
        const point = radarPoint(i, this.axes().length, ratio);
        return `${point.x},${point.y}`;
      })
      .join(' ');
  }
  shortLabel(label: string) {
    return this.axes().length === 7 && label.length > 2 ? label.slice(0, 2) : label;
  }
  format(value: number | null) {
    return value === null ? '—' : formatMetric(value, this.properties());
  }
  selectAxis(index: number) {
    this.activeIndex.set(index);
  }
  moveAxis(event: KeyboardEvent) {
    if (event.key === 'ArrowRight' || event.key === 'ArrowDown') {
      event.preventDefault();
      this.activeIndex.update((i) => (i + 1) % this.axes().length);
    } else if (event.key === 'ArrowLeft' || event.key === 'ArrowUp') {
      event.preventDefault();
      this.activeIndex.update((i) => (i - 1 + this.axes().length) % this.axes().length);
    } else if (event.key === 'Home') {
      event.preventDefault();
      this.activeIndex.set(0);
    } else if (event.key === 'End') {
      event.preventDefault();
      this.activeIndex.set(this.axes().length - 1);
    }
  }
}

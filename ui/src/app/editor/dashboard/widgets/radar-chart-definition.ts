import { ChartFocusStrategy, ChartPoint, defineChart } from '@tanstack/charts';
import {
  angleGrid,
  focusGroupAngle,
  polar,
  radialArea,
  radialGrid,
  radialLine,
  radialText,
} from '@tanstack/charts/polar';
import { scaleLinear } from '@tanstack/charts/scales/linear';
import { scalePoint } from '@tanstack/charts/scales/point';
import { tooltip } from '@tanstack/charts/tooltip';
import { curveLinear, curveLinearClosed } from 'd3-shape';
import { DashboardRadarAxis } from '../../../../api/dashboard/dashboard.service';
import { RadarWidgetProperties } from '../metric-widget-definition';

const GRID_RATIOS = [0.25, 0.5, 0.75, 1];
const WEEKDAY_AXES = 7;
const PROFILE_MARK_ID = 'profile';

export function radarEffectiveMaximum(
  values: readonly (number | null)[],
  requested: number | null,
  unit: 'ms' | 'count' | 'percent',
  target: number | null,
): number {
  const maximum = Math.max(0, ...values.map((value) => value ?? 0), target ?? 0);
  const automatic = unit === 'percent' ? 100 : roundedMaximum(maximum);
  return Math.max(
    requested ?? automatic,
    maximum > (requested ?? automatic) ? roundedMaximum(maximum) : 0,
    1,
  );
}

function roundedMaximum(value: number): number {
  if (value <= 0) return 1;
  const magnitude = 10 ** Math.floor(Math.log10(value));
  const normalized = value / magnitude;
  const step = normalized <= 1 ? 1 : normalized <= 2 ? 2 : normalized <= 5 ? 5 : 10;
  return step * magnitude;
}

export function radarAxisLabel(label: string, axisCount: number): string {
  return axisCount === WEEKDAY_AXES && label.length > 2 ? label.slice(0, 2) : label;
}

export function radarProfileRows(axes: readonly DashboardRadarAxis[]): DashboardRadarAxis[] {
  const firstGap = axes.findIndex((axis) => axis.value === null);
  return firstGap < 0 ? [...axes] : [...axes.slice(firstGap + 1), ...axes.slice(0, firstGap + 1)];
}

export interface RadarChartFormatters {
  value: (value: number | null) => string;
  metric: string;
}

type RadarPoint = ChartPoint<DashboardRadarAxis, string, number>;

const profilePoints = (points: readonly RadarPoint[]) =>
  points.filter((point) => point.markId === PROFILE_MARK_ID);

export const radarFocus: ChartFocusStrategy<DashboardRadarAxis, string, number> = {
  resolve: (points, context) =>
    focusGroupAngle.resolve(
      profilePoints(points) as readonly ChartPoint[],
      context,
    ) as RadarPoint[],
  group: (points, context) =>
    focusGroupAngle.group(profilePoints(points) as readonly ChartPoint[], {
      point: context.point as ChartPoint,
    }) as RadarPoint[],
  navigation: (points) =>
    focusGroupAngle.navigation(profilePoints(points) as readonly ChartPoint[]) as RadarPoint[],
};

export function radarChartTooltip(formatters: RadarChartFormatters) {
  return {
    use: tooltip,
    content: (points: readonly { datum: unknown }[]) => {
      const axis = points[0]?.datum as DashboardRadarAxis | undefined;
      return {
        title: axis?.label ?? '',
        rows: axis
          ? [
              { label: formatters.metric, value: formatters.value(axis.value) },
              ...(axis.status === 'complete'
                ? []
                : [
                    {
                      label: 'Status',
                      value: axis.status === 'partial' ? 'partial period' : 'not observed',
                    },
                  ]),
            ]
          : [],
      };
    },
  };
}

export function radarChartSpec(
  axes: readonly DashboardRadarAxis[],
  properties: RadarWidgetProperties,
  maximum: number,
  formatters: RadarChartFormatters,
) {
  const keys = axes.map((axis) => axis.key);
  const complete = axes.every((axis) => axis.value !== null);
  const rows = radarProfileRows(axes);
  const channels = { angle: 'key', radius: 'value', z: () => 'profile' } as const;
  const curve = complete ? curveLinearClosed : curveLinear;
  const target = properties.targetValue;
  return {
    marks: [
      polar({
        radiusRatio: 0.72,
        scales: {
          angle: { scale: scalePoint<string>().domain(keys), wrap: true },
          radius: { scale: scaleLinear().domain([0, maximum]) },
        },
        guides: [
          ...(properties.showGrid
            ? [
                radialGrid({
                  values: GRID_RATIOS.map((ratio) => ratio * maximum),
                  shape: 'polygon',
                  stroke: 'var(--border)',
                  strokeOpacity: 1,
                }),
              ]
            : []),
          ...(target === null
            ? []
            : [
                radialGrid({
                  values: [target],
                  shape: 'polygon',
                  stroke: 'var(--muted-foreground)',
                  strokeOpacity: 1,
                  strokeDasharray: '4 3',
                }),
              ]),
          angleGrid({
            labels: true,
            format: (key) =>
              radarAxisLabel(
                axes.find((axis) => axis.key === key)?.label ?? String(key),
                axes.length,
              ),
            stroke: 'var(--border)',
            strokeOpacity: properties.showGrid ? 1 : 0,
          }),
        ],
        marks: [
          ...(properties.showFill && complete
            ? [
                radialArea(rows, {
                  ...channels,
                  curve,
                  fill: 'var(--primary)',
                  fillOpacity: 0.12,
                }),
              ]
            : []),
          radialLine(rows, {
            id: PROFILE_MARK_ID,
            ...channels,
            curve,
            stroke: 'var(--primary)',
            strokeWidth: 2,
            points: true,
          }),
          ...(properties.showValues
            ? [
                radialText(
                  rows.filter((axis) => axis.value !== null),
                  {
                    ...channels,
                    text: (axis: DashboardRadarAxis) => formatters.value(axis.value),
                    fill: 'var(--foreground)',
                    fontSize: 10,
                    radiusOffset: 12,
                  },
                ),
              ]
            : []),
        ],
      }),
    ],
    scales: { x: null, y: null },
  };
}

export function radarChartDefinition(
  axes: readonly DashboardRadarAxis[],
  properties: RadarWidgetProperties,
  maximum: number,
  formatters: RadarChartFormatters,
) {
  return defineChart({
    focus: radarFocus,
    tooltip: radarChartTooltip(formatters),
    chart: () => radarChartSpec(axes, properties, maximum, formatters),
  });
}

import { defineChart, dot, lineY, ruleY } from '@tanstack/charts';
import { d3Curve } from '@tanstack/charts/d3/shape';
import { scaleLinear } from '@tanstack/charts/scales/linear';
import { scaleOrdinal } from '@tanstack/charts/scales/ordinal';
import { scalePoint } from '@tanstack/charts/scales/point';
import { tooltip } from '@tanstack/charts/tooltip';
import { curveMonotoneX } from 'd3-shape';
import { LineWidgetProperties } from '../metric-widget-definition';
import { LineChartData, LineChartRow } from './line-chart-data';

export const COMPACT_CHART_WIDTH = 240;
const TICK_SPACING = 80;

export interface LineChartFormatters {
  value: (value: number) => string;
  bucket: (bucket: string) => string;
}

export function bucketLabelFormatter(
  bucket: LineWidgetProperties['bucket'],
  timeZone: string,
  locale?: string,
): (bucket: string) => string {
  const formatter = new Intl.DateTimeFormat(
    locale,
    bucket === 'month'
      ? { month: 'short', timeZone }
      : { day: 'numeric', month: 'short', timeZone },
  );
  return (value) => formatter.format(new Date(value));
}

const SERIES_COLORS = [
  'var(--chart-1)',
  'var(--chart-2)',
  'var(--chart-3)',
  'var(--chart-4)',
  'var(--chart-5)',
  'var(--chart-6)',
];

export function seriesColor(index: number): string {
  return SERIES_COLORS[index % SERIES_COLORS.length];
}

export function evenlySpaced<T>(values: readonly T[], maximum: number): T[] {
  const step = Math.max(1, Math.ceil(values.length / maximum));
  return values.filter((_, index) => index % step === 0);
}

function seriesWithIsolatedPoint(rows: readonly LineChartRow[]): Set<string> {
  const filled = (row: LineChartRow | undefined, series: string) =>
    row?.series === series && row.value !== null;
  return new Set(
    rows
      .filter(
        (row, index) =>
          row.value !== null &&
          !filled(rows[index - 1], row.series) &&
          !filled(rows[index + 1], row.series),
      )
      .map((row) => row.series),
  );
}

function lineChartMarks(rows: readonly LineChartRow[], properties: LineWidgetProperties) {
  const channels = { x: 'bucket', y: 'value', z: 'series', color: 'series' } as const;
  if (properties.renderMode === 'points') return [dot(rows, channels)];
  const curve = properties.renderMode === 'smooth' ? d3Curve(curveMonotoneX) : undefined;
  const isolated = seriesWithIsolatedPoint(rows);
  const marked = (row: LineChartRow) => properties.showPoints || isolated.has(row.series);
  return [
    lineY(rows.filter(marked), { ...channels, curve, points: true }),
    lineY(
      rows.filter((row) => !marked(row)),
      { ...channels, curve },
    ),
  ];
}

export function lineChartSpec(
  data: LineChartData,
  properties: LineWidgetProperties,
  formatters: LineChartFormatters,
  width: number,
) {
  const compact = width < COMPACT_CHART_WIDTH;
  const keys = [...new Set(data.rows.map((row) => row.series))];
  const buckets = [...new Set(data.rows.map((row) => row.bucket))];
  const tickBuckets = evenlySpaced(buckets, Math.max(2, Math.floor(width / TICK_SPACING)));
  return {
    marks: [
      ...lineChartMarks(data.rows, properties),
      ...(properties.targetValue === null
        ? []
        : [
            ruleY([properties.targetValue], {
              stroke: 'var(--warning)',
              strokeOpacity: 1,
              strokeDasharray: '3 2',
            }),
          ]),
    ],
    scales: {
      x: {
        scale: () => scalePoint<string>().padding(0.3),
        axis: compact
          ? (false as const)
          : { ticks: { values: tickBuckets, format: formatters.bucket } },
      },
      y: {
        scale: scaleLinear().domain([0, data.maximum || 1]),
        nice: true,
        grid: properties.showGrid,
        axis: compact ? (false as const) : { ticks: { format: formatters.value } },
      },
    },
    color: {
      scale: scaleOrdinal(
        keys,
        keys.map((_, index) => seriesColor(index)),
      ),
    },
  };
}

export function lineChartTooltip(data: LineChartData, formatters: LineChartFormatters) {
  const labels = new Map(data.rows.map((row) => [row.series, row.label]));
  return {
    use: tooltip,
    sort: 'color-domain' as const,
    content: (
      points: readonly { groupLabel: string; xValue: unknown; yValue: unknown; color: string }[],
    ) => ({
      title: formatters.bucket(String(points[0]?.xValue ?? '')),
      rows: points.map((point) => ({
        label: labels.get(point.groupLabel) ?? point.groupLabel,
        value: formatters.value(Number(point.yValue)),
        color: point.color,
      })),
    }),
  };
}

export function lineChartDefinition(
  data: LineChartData,
  properties: LineWidgetProperties,
  formatters: LineChartFormatters,
) {
  return defineChart({
    focus: 'group-x',
    tooltip: lineChartTooltip(data, formatters),
    chart: ({ width }) => lineChartSpec(data, properties, formatters, width),
  });
}

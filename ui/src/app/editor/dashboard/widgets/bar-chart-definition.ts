import { barX, barY, defineChart, text } from '@tanstack/charts';
import { scaleBand } from '@tanstack/charts/scales/band';
import { scaleLinear } from '@tanstack/charts/scales/linear';
import { tooltip } from '@tanstack/charts/tooltip';
import { DashboardCategoriesPayload } from '../../../../api/dashboard/dashboard.service';
import { BarWidgetProperties } from '../metric-widget-definition';

export type BarCategory = DashboardCategoriesPayload['categories'][number];

export const BAR_LABEL_WIDTH = 96;
export const BAR_LABEL_HEIGHT = 28;
export const BAR_ROW_HEIGHT = 24;
export const BAR_COLUMN_WIDTH = 48;
const VALUE_LABEL_WIDTH = 72;
const VALUE_LABEL_HEIGHT = 18;
const MAXIMUM_BAR_THICKNESS = 16;
const VALUE_AXIS_ALLOWANCE = 80;
const CHARACTER_WIDTH = 6.4;

export interface BarChartRow {
  key: string;
  label: string;
  value: number | null;
  entity: BarCategory['entity'];
}

export function barChartRows(categories: readonly BarCategory[]): BarChartRow[] {
  return categories.map((category) => ({
    key: category.key,
    label: category.label,
    value: category.value,
    entity: category.entity,
  }));
}

export function barChartMaximum(rows: readonly BarChartRow[]): number {
  return Math.max(1, ...rows.map((row) => row.value ?? 0));
}

export function barChartMargin(properties: BarWidgetProperties) {
  return properties.orientation === 'vertical'
    ? {
        bottom: BAR_LABEL_HEIGHT,
        ...(properties.showValues ? { top: VALUE_LABEL_HEIGHT } : {}),
      }
    : {
        left: BAR_LABEL_WIDTH,
        ...(properties.showValues ? { right: VALUE_LABEL_WIDTH } : {}),
      };
}

export function barChartTooltip(
  rows: readonly BarChartRow[],
  format: (value: number | null) => string,
) {
  const labels = new Map(rows.map((row) => [row.key, row.label]));
  return {
    use: tooltip,
    content: (points: readonly { datum: unknown }[]) => {
      const row = points[0]?.datum as BarChartRow | undefined;
      return {
        title: row ? (labels.get(row.key) ?? row.label) : '',
        rows: row ? [{ label: 'Value', value: format(row.value) }] : [],
      };
    },
  };
}

export function valueLabelsFit(
  rows: readonly BarChartRow[],
  properties: BarWidgetProperties,
  format: (value: number | null) => string,
  width: number,
): boolean {
  if (properties.orientation !== 'vertical') return true;
  const longest = Math.max(0, ...rows.map((row) => format(row.value).length));
  return (width - VALUE_AXIS_ALLOWANCE) / Math.max(1, rows.length) >= longest * CHARACTER_WIDTH;
}

export function barChartSpec(
  rows: readonly BarChartRow[],
  properties: BarWidgetProperties,
  format: (value: number | null) => string,
  width: number,
) {
  const vertical = properties.orientation === 'vertical';
  const showValues = properties.showValues && valueLabelsFit(rows, properties, format, width);
  const valueScale = {
    scale: scaleLinear().domain([0, barChartMaximum(rows)]),
    nice: true,
    grid: properties.showGrid,
    axis: { ticks: { format: (value: number) => format(value) } },
  };
  const categoryScale = {
    scale: () => scaleBand<string>().paddingInner(0.3).paddingOuter(0.15),
    axis: false as const,
  };
  const labels = {
    text: (row: BarChartRow) => format(row.value),
    fill: 'var(--foreground)',
    fontSize: 11,
  };
  const withValue = rows.filter((row) => row.value !== null);
  return {
    margin: barChartMargin({ ...properties, showValues }),
    marks: vertical
      ? [
          barY(rows, {
            x: 'key',
            y: 'value',
            fill: 'var(--primary)',
            maxThickness: MAXIMUM_BAR_THICKNESS * 2,
          }),
          ...(showValues ? [text(withValue, { x: 'key', y: 'value', ...labels, dy: -8 })] : []),
        ]
      : [
          barX(rows, {
            x: 'value',
            y: 'key',
            fill: 'var(--primary)',
            maxThickness: MAXIMUM_BAR_THICKNESS,
          }),
          ...(showValues
            ? [text(withValue, { x: 'value', y: 'key', ...labels, anchor: 'start', dx: 6 })]
            : []),
        ],
    scales: vertical ? { x: categoryScale, y: valueScale } : { x: valueScale, y: categoryScale },
  };
}

export function barChartDefinition(
  rows: readonly BarChartRow[],
  properties: BarWidgetProperties,
  format: (value: number | null) => string,
) {
  return defineChart({
    focus: 'nearest',
    tooltip: barChartTooltip(rows, format),
    chart: ({ width }) => barChartSpec(rows, properties, format, width),
  });
}

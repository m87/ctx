import { DashboardMetricResult } from '../../../../api/dashboard/dashboard.service';
import { LineWidgetProperties } from '../metric-widget-definition';

export interface LineChartRow {
  bucket: string;
  series: string;
  label: string;
  value: number | null;
}

export interface LineChartData {
  rows: LineChartRow[];
  maximum: number;
}

export function lineChartRows(
  result: Extract<DashboardMetricResult, { kind: 'series' }>,
  properties: Pick<LineWidgetProperties, 'targetValue'>,
): LineChartData {
  const rows = result.series.series.flatMap((series) =>
    result.series.buckets.map((bucket, index) => ({
      bucket: bucket.start,
      series: series.key,
      label: series.label,
      value: bucket.status === 'future' ? null : (series.values[index] ?? null),
    })),
  );
  const maximum = Math.max(0, properties.targetValue ?? 0, ...rows.map((row) => row.value ?? 0));
  return { rows, maximum };
}

import { DashboardMetricResult } from '../../../../api/dashboard/dashboard.service';
import {
  MetricWidgetProperties,
  defaultMetricProperties,
  formatMetric,
} from '../metric-widget-definition';

type Unit = Extract<DashboardMetricResult, { unit: string }>['unit'];

export function formatResultValue(
  value: number | null,
  unit: Unit,
  properties: Pick<MetricWidgetProperties, 'format' | 'precision'>,
): string {
  if (value === null) return '—';
  return formatMetric(value, {
    ...defaultMetricProperties(),
    metric:
      unit === 'ms' ? 'duration.sum' : unit === 'count' ? 'sessions.count' : 'time_share.ratio',
    format: properties.format,
    precision: properties.precision,
  });
}

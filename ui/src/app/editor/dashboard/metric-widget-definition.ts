export const METRICS = [
  { id: 'duration.sum', label: 'Tracked time', unit: 'ms' },
  { id: 'duration.mean', label: 'Average session', unit: 'ms' },
  { id: 'duration.max', label: 'Longest session', unit: 'ms' },
  { id: 'sessions.count', label: 'Sessions', unit: 'count' },
  { id: 'contexts.distinct_count', label: 'Active contexts', unit: 'count' },
  { id: 'active_days.distinct_count', label: 'Active days', unit: 'count' },
  { id: 'time_share.ratio', label: 'Share of tracked time', unit: 'percent' },
] as const;

export const PERIODS = [
  { id: 'today', label: 'Today' },
  { id: 'yesterday', label: 'Yesterday' },
  { id: 'this_week', label: 'This week' },
  { id: 'last_week', label: 'Last week' },
  { id: 'this_month', label: 'This month' },
  { id: 'last_month', label: 'Last month' },
  { id: 'last_7_days', label: 'Last 7 days' },
  { id: 'last_30_days', label: 'Last 30 days' },
] as const;

export type MetricId = (typeof METRICS)[number]['id'];
export type DashboardPeriod = (typeof PERIODS)[number]['id'];
export type MetricTone = 'muted' | 'info' | 'success' | 'warning' | 'destructive';
export const METRIC_TONES: readonly MetricTone[] = [
  'muted',
  'info',
  'success',
  'warning',
  'destructive',
];

export interface MetricWidgetProperties {
  title: string;
  description: string;
  metric: MetricId;
  period: DashboardPeriod | 'inherit' | 'custom';
  startDate: string;
  endDate: string;
  projectId: string;
  contextId: string;
  tag: string;
  archived: 'include' | 'exclude' | 'only';
  includeRunning: boolean;
  format: 'duration' | 'hours' | 'number' | 'percent';
  precision: number;
}

export interface GaugeWidgetProperties extends MetricWidgetProperties {
  minimum: number;
  maximum: number;
  target: number | null;
  goalMode: 'at_least' | 'at_most' | 'informational';
  thresholds: { value: number; tone: MetricTone }[];
  showValue: boolean;
  showTarget: boolean;
  showPercentage: boolean;
}

export function metricInfo(id: MetricId) {
  return METRICS.find((metric) => metric.id === id) ?? METRICS[0];
}

export function defaultMetricProperties(): MetricWidgetProperties {
  return {
    title: '',
    description: '',
    metric: 'duration.sum',
    period: 'inherit',
    startDate: '',
    endDate: '',
    projectId: '',
    contextId: '',
    tag: '',
    archived: 'include',
    includeRunning: true,
    format: 'duration',
    precision: 1,
  };
}

export function defaultGaugeProperties(): GaugeWidgetProperties {
  return {
    ...defaultMetricProperties(),
    minimum: 0,
    maximum: 144000000,
    target: 144000000,
    goalMode: 'at_least',
    thresholds: [],
    showValue: true,
    showTarget: true,
    showPercentage: true,
  };
}

export function metricSettingsError(value: unknown, gauge = false): string {
  if (!value || typeof value !== 'object' || Array.isArray(value))
    return 'Metric settings are invalid.';
  const p = value as Record<string, unknown>;
  if (!METRICS.some((metric) => metric.id === p['metric'])) return 'Choose a metric.';
  for (const field of [
    'title',
    'description',
    'startDate',
    'endDate',
    'projectId',
    'contextId',
    'tag',
  ]) {
    if (typeof p[field] !== 'string' || p[field].length > 500)
      return 'Metric settings are invalid.';
  }
  if (!['inherit', 'custom', ...PERIODS.map((period) => period.id)].includes(String(p['period'])))
    return 'Choose a period.';
  if (
    p['period'] === 'custom' &&
    (!validDate(String(p['startDate'])) ||
      !validDate(String(p['endDate'])) ||
      String(p['endDate']) < String(p['startDate']))
  )
    return 'Choose a valid date range.';
  if (
    !['include', 'exclude', 'only'].includes(String(p['archived'])) ||
    typeof p['includeRunning'] !== 'boolean'
  )
    return 'Metric filters are invalid.';
  const unit = metricInfo(p['metric'] as MetricId).unit;
  const formats =
    unit === 'ms' ? ['duration', 'hours'] : unit === 'count' ? ['number'] : ['percent'];
  if (!formats.includes(String(p['format']))) return 'Choose a format matching the metric.';
  if (
    typeof p['precision'] !== 'number' ||
    !Number.isInteger(p['precision']) ||
    p['precision'] < 0 ||
    p['precision'] > 2 ||
    (unit === 'count' && p['precision'] !== 0)
  )
    return 'Precision must match the metric.';
  if (!gauge) return '';
  const min = p['minimum'],
    max = p['maximum'],
    target = p['target'];
  if (
    typeof min !== 'number' ||
    typeof max !== 'number' ||
    !Number.isFinite(min) ||
    !Number.isFinite(max) ||
    min < 0 ||
    max <= min
  )
    return 'Maximum must be greater than minimum.';
  if (!['at_least', 'at_most', 'informational'].includes(String(p['goalMode'])))
    return 'Choose a goal mode.';
  if (
    target !== null &&
    (typeof target !== 'number' ||
      !Number.isFinite(target) ||
      target <= 0 ||
      target < min ||
      target > max)
  )
    return 'The target must be positive and within the scale.';
  if (p['goalMode'] !== 'informational' && target === null) return 'Enter a target.';
  if (
    unit === 'count' &&
    (!Number.isInteger(min) ||
      !Number.isInteger(max) ||
      (target !== null && !Number.isInteger(target)))
  )
    return 'Counts require whole numbers.';
  if (['showValue', 'showTarget', 'showPercentage'].some((key) => typeof p[key] !== 'boolean'))
    return 'Gauge display settings are invalid.';
  if (!Array.isArray(p['thresholds']) || p['thresholds'].length > 8)
    return 'Use up to eight thresholds.';
  let previous = -Infinity;
  for (const entry of p['thresholds']) {
    if (!entry || typeof entry !== 'object') return 'Thresholds are invalid.';
    const threshold = entry as Record<string, unknown>;
    const value = threshold['value'];
    if (
      typeof value !== 'number' ||
      !Number.isFinite(value) ||
      value < min ||
      value > max ||
      value <= previous ||
      !METRIC_TONES.includes(threshold['tone'] as MetricTone) ||
      (unit === 'count' && !Number.isInteger(value))
    )
      return 'Thresholds must be increasing and within the scale.';
    previous = value;
  }
  return '';
}

function validDate(value: string): boolean {
  if (!/^\d{4}-\d{2}-\d{2}$/.test(value)) return false;
  const date = new Date(value + 'T00:00:00Z');
  return Number.isFinite(date.getTime()) && date.toISOString().slice(0, 10) === value;
}

export function formatMetric(value: number | null, p: MetricWidgetProperties): string {
  if (value === null || !Number.isFinite(value)) return 'No data';
  if (p.format === 'duration') {
    const minutes = Math.floor(value / 60000);
    if (minutes < 1 && value > 0) return '<1 min';
    const hours = Math.floor(minutes / 60);
    return hours
      ? `${hours} h${minutes % 60 ? ' ' + (minutes % 60) + ' min' : ''}`
      : `${minutes} min`;
  }
  const number = p.format === 'hours' ? value / 3600000 : value;
  const text = new Intl.NumberFormat(undefined, { maximumFractionDigits: p.precision }).format(
    number,
  );
  return text + (p.format === 'hours' ? ' h' : p.format === 'percent' ? '%' : '');
}

export function gaugeTone(value: number | null, p: GaugeWidgetProperties): MetricTone {
  if (value === null) return 'muted';
  if (p.thresholds.length)
    return p.thresholds.filter((threshold) => threshold.value <= value).at(-1)?.tone ?? 'muted';
  if (p.target === null || p.goalMode === 'informational') return 'info';
  if (p.goalMode === 'at_most')
    return value > p.target ? 'destructive' : value >= p.target * 0.8 ? 'warning' : 'success';
  return value >= p.target ? 'success' : 'info';
}

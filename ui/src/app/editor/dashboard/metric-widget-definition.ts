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
export type DashboardPeriod =
  (typeof PERIODS)[number]['id'] | 'last_3_months' | 'last_6_months' | 'last_12_months';
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

export interface RadarWidgetProperties extends MetricWidgetProperties {
  axesBy: 'weekday' | 'project' | 'context';
  axisKeys: string[];
  scaleMaximum: number | null;
  targetValue: number | null;
  showValues: boolean;
  showGrid: boolean;
  showFill: boolean;
  showSummary: boolean;
}

export interface BarWidgetProperties extends MetricWidgetProperties {
  groupBy: 'project' | 'context' | 'time';
  bucket: null | 'day' | 'week' | 'month';
  orientation: 'horizontal' | 'vertical';
  sort: { by: 'value' | 'label' | 'time'; direction: 'asc' | 'desc' };
  limit: number | null;
  showOther: boolean;
  showValues: boolean;
  showGrid: boolean;
  showSummary: boolean;
}
export interface HeatmapWidgetProperties extends MetricWidgetProperties {
  scaleMode: 'auto' | 'fixed';
  thresholds: number[];
  showMonthLabels: boolean;
  showDayLabels: boolean;
  showLegend: boolean;
  showSummary: boolean;
}
export interface LineWidgetProperties extends MetricWidgetProperties {
  bucket: 'day' | 'week' | 'month';
  seriesBy: 'none' | 'project' | 'context';
  seriesLimit: number;
  valueMode: 'periodic' | 'cumulative';
  renderMode: 'points' | 'linear' | 'smooth';
  showPoints: boolean;
  showLegend: boolean;
  showGrid: boolean;
  showSummary: boolean;
  targetValue: number | null;
}
export interface TableWidgetProperties extends Omit<
  MetricWidgetProperties,
  'metric' | 'format' | 'precision'
> {
  groupBy: 'project' | 'context';
  columns: {
    id: string;
    metric: MetricId;
    label: string;
    visible: boolean;
    format: MetricWidgetProperties['format'];
    precision: number;
    showBar: boolean;
  }[];
  sort: { columnId: string; direction: 'asc' | 'desc' };
  pageSize: 10 | 25 | 50 | 100;
  showSummary: boolean;
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

export function defaultRadarProperties(): RadarWidgetProperties {
  return {
    ...defaultMetricProperties(),
    axesBy: 'weekday',
    axisKeys: [],
    scaleMaximum: null,
    targetValue: null,
    showValues: false,
    showGrid: true,
    showFill: true,
    showSummary: true,
  };
}

export function defaultBarProperties(): BarWidgetProperties {
  return {
    ...defaultMetricProperties(),
    groupBy: 'project',
    bucket: null,
    orientation: 'horizontal',
    sort: { by: 'value', direction: 'desc' },
    limit: 10,
    showOther: true,
    showValues: true,
    showGrid: true,
    showSummary: true,
  };
}
export function defaultHeatmapProperties(): HeatmapWidgetProperties {
  return {
    ...defaultMetricProperties(),
    period: 'last_3_months',
    scaleMode: 'auto',
    thresholds: [],
    showMonthLabels: true,
    showDayLabels: true,
    showLegend: true,
    showSummary: true,
  };
}
export function defaultLineProperties(): LineWidgetProperties {
  return {
    ...defaultMetricProperties(),
    bucket: 'day',
    seriesBy: 'none',
    seriesLimit: 5,
    valueMode: 'periodic',
    renderMode: 'smooth',
    showPoints: false,
    showLegend: true,
    showGrid: true,
    showSummary: true,
    targetValue: null,
  };
}
export function defaultTableProperties(): TableWidgetProperties {
  const {
    metric: _metric,
    format: _format,
    precision: _precision,
    ...common
  } = defaultMetricProperties();
  return {
    ...common,
    groupBy: 'project',
    columns: [
      {
        id: 'time',
        metric: 'duration.sum',
        label: 'Time',
        visible: true,
        format: 'duration',
        precision: 1,
        showBar: false,
      },
      {
        id: 'sessions',
        metric: 'sessions.count',
        label: 'Sessions',
        visible: true,
        format: 'number',
        precision: 0,
        showBar: false,
      },
      {
        id: 'average',
        metric: 'duration.mean',
        label: 'Average session',
        visible: true,
        format: 'duration',
        precision: 1,
        showBar: false,
      },
      {
        id: 'share',
        metric: 'time_share.ratio',
        label: 'Share',
        visible: true,
        format: 'percent',
        precision: 1,
        showBar: false,
      },
    ],
    sort: { columnId: 'time', direction: 'desc' },
    pageSize: 10,
    showSummary: true,
  };
}

export function metricSettingsError(
  value: unknown,
  gauge = false,
  radar = false,
  type = '',
): string {
  if (!value || typeof value !== 'object' || Array.isArray(value))
    return 'Metric settings are invalid.';
  const p = value as Record<string, unknown>;
  if (type !== 'table' && !METRICS.some((metric) => metric.id === p['metric']))
    return 'Choose a metric.';
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
  const extraPeriods =
    type === 'heatmap' ? ['last_3_months', 'last_6_months', 'last_12_months'] : [];
  if (
    !['inherit', 'custom', ...PERIODS.map((period) => period.id), ...extraPeriods].includes(
      String(p['period']),
    )
  )
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
  if (type === 'table') return tableSettingsError(p);
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
  if (radar) return radarSettingsError(p, unit);
  if (type === 'heatmap') {
    if (
      !['auto', 'fixed'].includes(String(p['scaleMode'])) ||
      !Array.isArray(p['thresholds']) ||
      (p['scaleMode'] === 'auto' ? p['thresholds'].length !== 0 : p['thresholds'].length !== 3)
    )
      return 'Heatmap scale settings are invalid.';
    let previous = 0;
    for (const value of p['thresholds']) {
      if (
        typeof value !== 'number' ||
        !Number.isFinite(value) ||
        value <= previous ||
        (unit === 'count' && !Number.isInteger(value)) ||
        (unit === 'percent' && value >= 100)
      )
        return 'Heatmap thresholds must increase.';
      previous = value;
    }
    if (
      ['showMonthLabels', 'showDayLabels', 'showLegend', 'showSummary'].some(
        (key) => typeof p[key] !== 'boolean',
      )
    )
      return 'Heatmap display settings are invalid.';
    return '';
  }
  if (type === 'bar') {
    if (
      !['project', 'context', 'time'].includes(String(p['groupBy'])) ||
      !['horizontal', 'vertical'].includes(String(p['orientation']))
    )
      return 'Bar chart settings are invalid.';
    if (
      p['groupBy'] === 'time'
        ? !['day', 'week', 'month'].includes(String(p['bucket'])) ||
          (p['sort'] as { by?: string })?.by !== 'time' ||
          p['limit'] !== null ||
          p['showOther'] !== false
        : p['bucket'] !== null && p['bucket'] !== undefined
    )
      return 'Bar chart grouping settings are invalid.';
    const sort = p['sort'] as { by?: unknown; direction?: unknown } | undefined;
    if (
      !sort ||
      !['asc', 'desc'].includes(String(sort.direction)) ||
      (p['groupBy'] === 'time'
        ? sort.direction !== 'asc'
        : !['value', 'label'].includes(String(sort.by))) ||
      (p['groupBy'] !== 'time' &&
        (typeof p['limit'] !== 'number' ||
          !Number.isInteger(p['limit']) ||
          p['limit'] < 2 ||
          p['limit'] > 50)) ||
      ['showValues', 'showGrid', 'showSummary'].some((key) => typeof p[key] !== 'boolean')
    )
      return 'Bar ranking settings are invalid.';
    if (p['showOther'] === true && p['metric'] !== 'duration.sum')
      return 'Other is available only for tracked time.';
    return '';
  }
  if (type === 'line') {
    if (
      !['day', 'week', 'month'].includes(String(p['bucket'])) ||
      !['none', 'project', 'context'].includes(String(p['seriesBy'])) ||
      typeof p['seriesLimit'] !== 'number' ||
      p['seriesLimit'] < 1 ||
      p['seriesLimit'] > 5 ||
      !['periodic', 'cumulative'].includes(String(p['valueMode'])) ||
      (p['valueMode'] === 'cumulative' && p['metric'] !== 'duration.sum') ||
      !['points', 'linear', 'smooth'].includes(String(p['renderMode'])) ||
      (p['targetValue'] !== null &&
        (typeof p['targetValue'] !== 'number' ||
          !Number.isFinite(p['targetValue']) ||
          p['targetValue'] < 0 ||
          (unit === 'percent' && p['targetValue'] > 100) ||
          (unit === 'count' && !Number.isInteger(p['targetValue'])))) ||
      ['showPoints', 'showLegend', 'showGrid', 'showSummary'].some(
        (key) => typeof p[key] !== 'boolean',
      )
    )
      return 'Line chart settings are invalid.';
    return '';
  }
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

function radarSettingsError(p: Record<string, unknown>, unit: string): string {
  if (!['weekday', 'project', 'context'].includes(String(p['axesBy'])))
    return 'Choose an axis grouping.';
  const keys = p['axisKeys'];
  if (!Array.isArray(keys) || (p['axesBy'] === 'weekday' ? keys.length !== 0 : keys.length > 8))
    return 'Choose no more than eight axes.';
  const prefix = p['axesBy'] === 'project' ? 'project:' : 'context:';
  if (p['axesBy'] !== 'weekday') {
    const unique = new Set<string>();
    for (const key of keys) {
      if (
        typeof key !== 'string' ||
        (key !== 'project:none' && !key.startsWith(prefix)) ||
        unique.has(key)
      )
        return 'Axis selections are invalid.';
      unique.add(key);
    }
  }
  for (const key of ['scaleMaximum', 'targetValue'] as const) {
    const value = p[key];
    if (
      value !== null &&
      (typeof value !== 'number' ||
        !Number.isFinite(value) ||
        value < 0 ||
        (key === 'scaleMaximum' && value === 0) ||
        (unit === 'count' && !Number.isInteger(value)) ||
        (unit === 'percent' && value > 100))
    )
      return `Invalid radar ${key}.`;
  }
  if (
    ['showValues', 'showGrid', 'showFill', 'showSummary'].some((key) => typeof p[key] !== 'boolean')
  )
    return 'Radar appearance settings are invalid.';
  return '';
}

function tableSettingsError(p: Record<string, unknown>): string {
  if (
    !['project', 'context'].includes(String(p['groupBy'])) ||
    !Array.isArray(p['columns']) ||
    p['columns'].length < 1 ||
    p['columns'].length > 7 ||
    ![10, 25, 50, 100].includes(Number(p['pageSize']))
  )
    return 'Table settings are invalid.';
  const columns = p['columns'] as {
    id?: unknown;
    metric?: unknown;
    label?: unknown;
    visible?: unknown;
    format?: unknown;
    precision?: unknown;
    showBar?: unknown;
  }[];
  const ids = new Set<unknown>(),
    metrics = new Set<unknown>();
  for (const column of columns) {
    const metric = METRICS.find((entry) => entry.id === column.metric);
    if (
      !metric ||
      typeof column.id !== 'string' ||
      !column.id ||
      column.id.length > 64 ||
      column.id === 'group' ||
      ids.has(column.id) ||
      metrics.has(column.metric) ||
      typeof column.label !== 'string' ||
      column.label.length > 100 ||
      typeof column.visible !== 'boolean'
    )
      return 'Table columns are invalid.';
    const formats =
      metric.unit === 'ms'
        ? ['duration', 'hours']
        : metric.unit === 'count'
          ? ['number']
          : ['percent'];
    if (
      !formats.includes(String(column.format)) ||
      typeof column.precision !== 'number' ||
      !Number.isInteger(column.precision) ||
      column.precision < 0 ||
      column.precision > 2 ||
      (metric.unit === 'count' && column.precision !== 0) ||
      typeof column.showBar !== 'boolean' ||
      (column.showBar && metric.id !== 'duration.sum' && metric.id !== 'time_share.ratio')
    )
      return 'Table column settings are invalid.';
    ids.add(column.id);
    metrics.add(column.metric);
  }
  if (!columns.some((column) => column.visible)) return 'Show at least one table column.';
  if (typeof p['showSummary'] !== 'boolean') return 'Table display settings are invalid.';
  const sort = p['sort'] as { columnId?: unknown; direction?: unknown } | undefined;
  if (
    !sort ||
    !['asc', 'desc'].includes(String(sort.direction)) ||
    !(
      sort.columnId === 'group' ||
      columns.some((column) => column.id === sort.columnId && column.visible)
    )
  )
    return 'Choose a visible table sort column.';
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

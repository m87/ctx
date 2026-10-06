import {
  DashboardPeriod,
  MetricWidgetProperties,
  GaugeWidgetProperties,
  RadarWidgetProperties,
  BarWidgetProperties,
  HeatmapWidgetProperties,
  LineWidgetProperties,
  TableWidgetProperties,
  defaultMetricProperties,
  defaultGaugeProperties,
  defaultRadarProperties,
  defaultBarProperties,
  defaultHeatmapProperties,
  defaultLineProperties,
  defaultTableProperties,
  metricSettingsError,
  PERIODS,
} from './metric-widget-definition';

export interface DashboardWidgetLayout {
  x: number;
  y: number;
  width: number;
  height: number;
}

export interface TextWidgetProperties {
  text: string;
  horizontalAlign?: 'left' | 'center' | 'right';
  verticalAlign?: 'top' | 'center' | 'bottom';
  fontSize?: number;
}

export const TEXT_WIDGET_FONT_SIZE = { minimum: 8, maximum: 96, default: 14 };

export interface TextWidgetDefinition {
  id: string;
  type: 'text';
  layout: DashboardWidgetLayout;
  query: string;
  properties: TextWidgetProperties;
}

export interface KpiWidgetDefinition {
  id: string;
  type: 'kpi';
  layout: DashboardWidgetLayout;
  query: string;
  properties: MetricWidgetProperties;
}

export interface GaugeWidgetDefinition {
  id: string;
  type: 'gauge';
  layout: DashboardWidgetLayout;
  query: string;
  properties: GaugeWidgetProperties;
}

export interface RadarWidgetDefinition {
  id: string;
  type: 'radar';
  layout: DashboardWidgetLayout;
  query: string;
  properties: RadarWidgetProperties;
}

export interface BarWidgetDefinition {
  id: string;
  type: 'bar';
  layout: DashboardWidgetLayout;
  query: string;
  properties: BarWidgetProperties;
}
export interface HeatmapWidgetDefinition {
  id: string;
  type: 'heatmap';
  layout: DashboardWidgetLayout;
  query: string;
  properties: HeatmapWidgetProperties;
}
export interface LineWidgetDefinition {
  id: string;
  type: 'line';
  layout: DashboardWidgetLayout;
  query: string;
  properties: LineWidgetProperties;
}
export interface TableWidgetDefinition {
  id: string;
  type: 'table';
  layout: DashboardWidgetLayout;
  query: string;
  properties: TableWidgetProperties;
}

export type DashboardWidgetDefinition =
  | TextWidgetDefinition
  | KpiWidgetDefinition
  | GaugeWidgetDefinition
  | RadarWidgetDefinition
  | BarWidgetDefinition
  | HeatmapWidgetDefinition
  | LineWidgetDefinition
  | TableWidgetDefinition;
export type DashboardWidgetType = DashboardWidgetDefinition['type'];

export interface DashboardDefinition {
  [key: string]: unknown;
  version: 1;
  query: string;
  period?: DashboardPeriod;
  widgets: DashboardWidgetDefinition[];
}

export const TEXT_WIDGET_CATALOG_ENTRY = {
  type: 'text' as const,
  label: 'Text',
  description: 'Display a text note on your dashboard.',
  icon: 'lucideType',
  defaultLayout: { width: 4, height: 3 },
  minimumLayout: { width: 2, height: 2 },
  createProperties: (): TextWidgetProperties => ({ text: 'Text' }),
};

export const WIDGET_CATALOG = [
  TEXT_WIDGET_CATALOG_ENTRY,
  {
    type: 'kpi' as const,
    label: 'KPI',
    icon: 'lucideHash',
    defaultLayout: { width: 4, height: 3 },
    minimumLayout: { width: 3, height: 2 },
    createProperties: defaultMetricProperties,
  },
  {
    type: 'gauge' as const,
    label: 'Gauge',
    icon: 'lucideGauge',
    defaultLayout: { width: 6, height: 5 },
    minimumLayout: { width: 4, height: 4 },
    createProperties: defaultGaugeProperties,
  },
  {
    type: 'radar' as const,
    label: 'Radar chart',
    description: 'Compare one metric across categories.',
    icon: 'lucideRadar',
    defaultLayout: { width: 6, height: 6 },
    minimumLayout: { width: 5, height: 5 },
    createProperties: defaultRadarProperties,
  },
  {
    type: 'bar' as const,
    label: 'Bar chart',
    icon: 'lucideChartBar',
    defaultLayout: { width: 8, height: 5 },
    minimumLayout: { width: 5, height: 4 },
    createProperties: defaultBarProperties,
  },
  {
    type: 'heatmap' as const,
    label: 'Heatmap',
    icon: 'lucideCalendarDays',
    defaultLayout: { width: 12, height: 5 },
    minimumLayout: { width: 6, height: 4 },
    createProperties: defaultHeatmapProperties,
  },
  {
    type: 'line' as const,
    label: 'Line chart',
    icon: 'lucideChartNoAxesCombined',
    defaultLayout: { width: 8, height: 5 },
    minimumLayout: { width: 5, height: 4 },
    createProperties: defaultLineProperties,
  },
  {
    type: 'table' as const,
    label: 'Table',
    icon: 'lucideTable',
    defaultLayout: { width: 10, height: 6 },
    minimumLayout: { width: 6, height: 4 },
    createProperties: defaultTableProperties,
  },
];

export function widgetCatalogEntry(type: DashboardWidgetType) {
  return WIDGET_CATALOG.find((entry) => entry.type === type) ?? TEXT_WIDGET_CATALOG_ENTRY;
}

export function createWidget(
  type: DashboardWidgetType,
  id: string,
  layout: DashboardWidgetLayout,
): DashboardWidgetDefinition {
  switch (type) {
    case 'text':
      return {
        id,
        type,
        layout,
        query: '',
        properties: TEXT_WIDGET_CATALOG_ENTRY.createProperties(),
      };
    case 'kpi':
      return { id, type, layout, query: '', properties: defaultMetricProperties() };
    case 'gauge':
      return { id, type, layout, query: '', properties: defaultGaugeProperties() };
    case 'radar':
      return { id, type, layout, query: '', properties: defaultRadarProperties() };
    case 'bar':
      return { id, type, layout, query: '', properties: defaultBarProperties() };
    case 'heatmap':
      return { id, type, layout, query: '', properties: defaultHeatmapProperties() };
    case 'line':
      return { id, type, layout, query: '', properties: defaultLineProperties() };
    case 'table':
      return { id, type, layout, query: '', properties: defaultTableProperties() };
  }
}

export function normalizeDashboardDefinition(value: unknown): DashboardDefinition {
  if (!value || typeof value !== 'object' || Array.isArray(value)) {
    throw new Error('Dashboard definition must be an object.');
  }
  const definition = value as Record<string, unknown>;
  if (Object.keys(definition).length === 0) {
    return { version: 1, query: '', widgets: [] };
  }
  if (definition['version'] !== 1) {
    throw new Error('This dashboard uses an unsupported definition version.');
  }
  if (typeof definition['query'] !== 'string' || !Array.isArray(definition['widgets'])) {
    throw new Error('Dashboard definition is invalid.');
  }
  if (
    definition['period'] !== undefined &&
    !PERIODS.some((period) => period.id === definition['period'])
  )
    throw new Error('Dashboard period is invalid.');
  const widgets: unknown[] = definition['widgets'];
  for (const value of widgets) {
    if (!value || typeof value !== 'object' || Array.isArray(value)) {
      throw new Error('Dashboard widget is invalid.');
    }
    const widget = value as Record<string, unknown>;
    if (
      !['text', 'kpi', 'gauge', 'radar', 'bar', 'heatmap', 'line', 'table'].includes(
        String(widget['type']),
      )
    ) {
      throw new Error('This dashboard contains an unsupported widget type.');
    }
    const properties: unknown = widget['properties'];
    if (
      widget['type'] === 'text'
        ? !isTextWidgetProperties(properties)
        : !!metricSettingsError(
            properties,
            widget['type'] === 'gauge',
            widget['type'] === 'radar',
            String(widget['type']),
          )
    ) {
      throw new Error('Widget settings are invalid.');
    }
    if (
      widget['type'] === 'radar' &&
      (properties as { axesBy?: string; axisKeys?: unknown[] }).axesBy !== 'weekday' &&
      ((properties as { axisKeys?: unknown[] }).axisKeys?.length ?? 0) < 3
    )
      throw new Error('Choose at least three radar axes.');
  }
  return structuredClone(value) as DashboardDefinition;
}

export function isTextWidgetProperties(value: unknown): value is TextWidgetProperties {
  if (!value || typeof value !== 'object' || Array.isArray(value)) return false;
  const properties = value as Record<string, unknown>;
  return (
    typeof properties['text'] === 'string' &&
    properties['text'].length <= 20_000 &&
    (properties['horizontalAlign'] === undefined ||
      (typeof properties['horizontalAlign'] === 'string' &&
        ['left', 'center', 'right'].includes(properties['horizontalAlign']))) &&
    (properties['verticalAlign'] === undefined ||
      (typeof properties['verticalAlign'] === 'string' &&
        ['top', 'center', 'bottom'].includes(properties['verticalAlign']))) &&
    (properties['fontSize'] === undefined ||
      (typeof properties['fontSize'] === 'number' &&
        Number.isInteger(properties['fontSize']) &&
        properties['fontSize'] >= TEXT_WIDGET_FONT_SIZE.minimum &&
        properties['fontSize'] <= TEXT_WIDGET_FONT_SIZE.maximum))
  );
}

export function dashboardDefinitionWithWidgets(
  definition: DashboardDefinition,
  widgets: DashboardWidgetDefinition[],
): DashboardDefinition {
  return { ...definition, widgets };
}

export function firstFreeWidgetPosition(
  widgets: readonly DashboardWidgetDefinition[],
  width = TEXT_WIDGET_CATALOG_ENTRY.defaultLayout.width,
  height = TEXT_WIDGET_CATALOG_ENTRY.defaultLayout.height,
): { x: number; y: number } | null {
  for (let y = 0; ; y++) {
    for (let x = 0; x <= 16 - width; x++) {
      if (widgets.every((widget) => !layoutsOverlap({ x, y, width, height }, widget.layout))) {
        return { x, y };
      }
    }
  }
}

export function layoutsOverlap(a: DashboardWidgetLayout, b: DashboardWidgetLayout): boolean {
  return a.x < b.x + b.width && a.x + a.width > b.x && a.y < b.y + b.height && a.y + a.height > b.y;
}

export function pixelDeltaToCells(delta: number, cellSize: number): number {
  return cellSize > 0 ? Math.round(delta / cellSize) : 0;
}

export function dashboardGridRowCount(widgets: readonly DashboardWidgetDefinition[]): number {
  return widgets.reduce(
    (count, widget) => Math.max(count, widget.layout.y + widget.layout.height),
    32,
  );
}

export function addDashboardWidget(
  definition: DashboardDefinition,
  id: string,
  type: DashboardWidgetType = 'text',
  target?: Pick<DashboardWidgetLayout, 'x' | 'y'>,
): DashboardDefinition {
  if (definition.widgets.length >= 200 || definition.widgets.some((widget) => widget.id === id)) {
    return definition;
  }
  const entry = widgetCatalogEntry(type);
  const position =
    target ??
    firstFreeWidgetPosition(
      definition.widgets,
      entry.defaultLayout.width,
      entry.defaultLayout.height,
    );
  if (!position) {
    return definition;
  }
  const layout = { ...position, ...entry.defaultLayout };
  if (
    !Number.isSafeInteger(layout.x) ||
    !Number.isSafeInteger(layout.y) ||
    layout.x < 0 ||
    layout.y < 0 ||
    layout.x + layout.width > 16 ||
    definition.widgets.some((widget) => layoutsOverlap(layout, widget.layout))
  ) {
    return definition;
  }
  const widget = createWidget(type, id, layout);
  return { ...definition, widgets: [...definition.widgets, widget] };
}

export function updateTextWidget(
  definition: DashboardDefinition,
  widgetId: string,
  update: Partial<Pick<TextWidgetDefinition, 'query' | 'properties'>>,
): DashboardDefinition {
  return {
    ...definition,
    widgets: definition.widgets.map((widget) =>
      widget.id === widgetId && widget.type === 'text'
        ? { ...widget, ...update, properties: update.properties ?? widget.properties }
        : widget,
    ),
  };
}

export function moveTextWidget(
  definition: DashboardDefinition,
  widgetId: string,
  x: number,
  y: number,
): DashboardDefinition {
  const widget = definition.widgets.find((item) => item.id === widgetId);
  if (!widget) return definition;
  return updateTextWidgetLayout(definition, widgetId, { ...widget.layout, x, y });
}

export function updateTextWidgetLayout(
  definition: DashboardDefinition,
  widgetId: string,
  target: DashboardWidgetLayout,
): DashboardDefinition {
  if (
    !Object.values(target).every(Number.isSafeInteger) ||
    target.x < 0 ||
    target.y < 0 ||
    target.width <
      widgetCatalogEntry(
        definition.widgets.find((widget) => widget.id === widgetId)?.type ?? 'text',
      ).minimumLayout.width ||
    target.height <
      widgetCatalogEntry(
        definition.widgets.find((widget) => widget.id === widgetId)?.type ?? 'text',
      ).minimumLayout.height ||
    target.x + target.width > 16 ||
    definition.widgets.some(
      (other) => other.id !== widgetId && layoutsOverlap(target, other.layout),
    )
  ) {
    return definition;
  }
  return {
    ...definition,
    widgets: definition.widgets.map((item) =>
      item.id === widgetId ? { ...item, layout: target } : item,
    ),
  };
}

export function resizeTextWidget(
  definition: DashboardDefinition,
  widgetId: string,
  width: number,
  height: number,
): DashboardDefinition {
  const widget = definition.widgets.find((item) => item.id === widgetId);
  if (!widget) return definition;
  return updateTextWidgetLayout(definition, widgetId, { ...widget.layout, width, height });
}

export function deleteTextWidget(
  definition: DashboardDefinition,
  widgetId: string,
): DashboardDefinition {
  return { ...definition, widgets: definition.widgets.filter((widget) => widget.id !== widgetId) };
}

export function addTextWidget(
  definition: DashboardDefinition,
  id: string,
  target?: Pick<DashboardWidgetLayout, 'x' | 'y'>,
): DashboardDefinition {
  return addDashboardWidget(definition, id, 'text', target);
}

export function replaceDashboardWidget(
  definition: DashboardDefinition,
  widget: DashboardWidgetDefinition,
): DashboardDefinition {
  return {
    ...definition,
    widgets: definition.widgets.map((item) => (item.id === widget.id ? widget : item)),
  };
}

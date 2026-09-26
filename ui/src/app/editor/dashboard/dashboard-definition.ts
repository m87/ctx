export interface DashboardWidgetLayout {
  x: number;
  y: number;
  width: number;
  height: number;
}

export interface TextWidgetProperties {
  text: string;
}

export interface TextWidgetDefinition {
  id: string;
  type: 'text';
  layout: DashboardWidgetLayout;
  query: string;
  properties: TextWidgetProperties;
}

export interface DashboardDefinition {
  [key: string]: unknown;
  version: 1;
  query: string;
  widgets: TextWidgetDefinition[];
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
  return structuredClone(value) as DashboardDefinition;
}

export function dashboardDefinitionWithWidgets(
  definition: DashboardDefinition,
  widgets: TextWidgetDefinition[],
): DashboardDefinition {
  return { ...definition, widgets };
}

export function firstFreeWidgetPosition(
  widgets: readonly TextWidgetDefinition[],
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

export function dashboardGridRowCount(widgets: readonly TextWidgetDefinition[]): number {
  return widgets.reduce(
    (count, widget) => Math.max(count, widget.layout.y + widget.layout.height),
    32,
  );
}

export function addTextWidget(
  definition: DashboardDefinition,
  id: string,
  target?: Pick<DashboardWidgetLayout, 'x' | 'y'>,
): DashboardDefinition {
  if (definition.widgets.length >= 200 || definition.widgets.some((widget) => widget.id === id)) {
    return definition;
  }
  const position = target ?? firstFreeWidgetPosition(definition.widgets);
  if (!position) {
    return definition;
  }
  const layout = { ...position, ...TEXT_WIDGET_CATALOG_ENTRY.defaultLayout };
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
  const widget: TextWidgetDefinition = {
    id,
    type: 'text',
    layout,
    query: '',
    properties: TEXT_WIDGET_CATALOG_ENTRY.createProperties(),
  };
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
      widget.id === widgetId
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
    target.width < TEXT_WIDGET_CATALOG_ENTRY.minimumLayout.width ||
    target.height < TEXT_WIDGET_CATALOG_ENTRY.minimumLayout.height ||
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

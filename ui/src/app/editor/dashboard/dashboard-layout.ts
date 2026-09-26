import {
  DashboardWidgetLayout,
  layoutsOverlap,
  pixelDeltaToCells,
  TextWidgetDefinition,
  TEXT_WIDGET_CATALOG_ENTRY,
} from './dashboard-definition';

export const DASHBOARD_COLUMNS = 16;
export const DASHBOARD_WIDGET_INSET = 10;

export interface GridPoint {
  x: number;
  y: number;
}

export interface DashboardGridGeometry {
  left: number;
  top: number;
  width: number;
  height: number;
}

export type WidgetResizeEdge = 'n' | 'ne' | 'e' | 'se' | 's' | 'sw' | 'w' | 'nw';

export function widgetResizeLayout(
  layout: DashboardWidgetLayout,
  edge: WidgetResizeEdge,
  pixelDelta: GridPoint,
  cellSize: number,
): DashboardWidgetLayout {
  const dx = pixelDeltaToCells(pixelDelta.x, cellSize);
  const dy = pixelDeltaToCells(pixelDelta.y, cellSize);
  const minimum = TEXT_WIDGET_CATALOG_ENTRY.minimumLayout;
  let left = layout.x;
  let top = layout.y;
  let right = layout.x + layout.width;
  let bottom = layout.y + layout.height;
  if (edge.includes('w')) left = Math.max(0, Math.min(right - minimum.width, left + dx));
  if (edge.includes('e'))
    right = Math.max(left + minimum.width, Math.min(DASHBOARD_COLUMNS, right + dx));
  if (edge.includes('n')) top = Math.max(0, Math.min(bottom - minimum.height, top + dy));
  if (edge.includes('s')) bottom = Math.max(top + minimum.height, bottom + dy);
  return { x: left, y: top, width: right - left, height: bottom - top };
}

export function widgetDropLayout(
  pointer: GridPoint,
  geometry: DashboardGridGeometry,
  layout: DashboardWidgetLayout,
  pickup?: GridPoint,
): DashboardWidgetLayout {
  const cellSize = geometry.width / DASHBOARD_COLUMNS;
  if (cellSize <= 0) return layout;
  const offset = pickup ?? { x: 0, y: 0 };
  const inset = pickup ? DASHBOARD_WIDGET_INSET : 0;
  const snap = pickup ? Math.round : Math.floor;
  const x = snap((pointer.x - geometry.left - offset.x - inset) / cellSize);
  const y = snap((pointer.y - geometry.top - offset.y - inset) / cellSize);
  return {
    ...layout,
    x: Math.max(0, Math.min(DASHBOARD_COLUMNS - layout.width, x)),
    y: Math.max(0, y),
  };
}

export function isPointerOnGrid(pointer: GridPoint, geometry: DashboardGridGeometry): boolean {
  return (
    pointer.x >= geometry.left &&
    pointer.x < geometry.left + geometry.width &&
    pointer.y >= geometry.top &&
    pointer.y < geometry.top + geometry.height
  );
}

export function canPlaceWidget(
  widgets: readonly TextWidgetDefinition[],
  layout: DashboardWidgetLayout,
  excludedId?: string,
): boolean {
  const minimum = TEXT_WIDGET_CATALOG_ENTRY.minimumLayout;
  return (
    Object.values(layout).every(Number.isSafeInteger) &&
    layout.x >= 0 &&
    layout.y >= 0 &&
    layout.width >= minimum.width &&
    layout.height >= minimum.height &&
    layout.x + layout.width <= DASHBOARD_COLUMNS &&
    widgets.every((widget) => widget.id === excludedId || !layoutsOverlap(layout, widget.layout))
  );
}

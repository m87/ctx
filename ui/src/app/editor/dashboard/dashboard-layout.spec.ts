import { addTextWidget, normalizeDashboardDefinition } from './dashboard-definition';
import {
  canPlaceWidget,
  isPointerOnGrid,
  widgetDropLayout,
  widgetResizeLayout,
  WidgetResizeEdge,
} from './dashboard-layout';

describe('dashboard drag geometry', () => {
  const geometry = { left: 100, top: 60, width: 800, height: 1600 };
  const layout = { x: 0, y: 0, width: 4, height: 3 };

  it('snaps moves around the pickup offset with a ten pixel inset', () => {
    const pickup = { x: 30, y: 15 };
    expect(widgetDropLayout({ x: 140, y: 85 }, geometry, layout, pickup)).toEqual(layout);
    expect(widgetDropLayout({ x: 164, y: 109 }, geometry, layout, pickup)).toEqual(layout);
    expect(widgetDropLayout({ x: 166, y: 111 }, geometry, layout, pickup)).toEqual({
      ...layout,
      x: 1,
      y: 1,
    });
  });

  it('places a new palette widget at the cell under the cursor and clamps to grid edges', () => {
    expect(widgetDropLayout({ x: 275, y: 185 }, geometry, layout)).toEqual({
      ...layout,
      x: 3,
      y: 2,
    });
    expect(widgetDropLayout({ x: 899, y: 80 }, geometry, layout).x).toBe(12);
    expect(widgetDropLayout({ x: -100, y: -100 }, geometry, layout, { x: 10, y: 10 })).toEqual(
      layout,
    );
    expect(isPointerOnGrid({ x: 901, y: 100 }, geometry)).toBe(false);
  });

  it('uses current canvas geometry after scrolling or resizing', () => {
    const pickup = { x: 15, y: 15 };
    const before = widgetDropLayout({ x: 225, y: 235 }, geometry, layout, pickup);
    const afterScroll = widgetDropLayout(
      { x: 225, y: 135 },
      { ...geometry, top: -40 },
      layout,
      pickup,
    );
    expect(afterScroll).toEqual(before);
    expect(
      widgetDropLayout({ x: 175, y: 160 }, { ...geometry, width: 400 }, layout, pickup),
    ).toEqual(before);
  });

  it('rejects occupied cells without blocking the original position of a moved widget', () => {
    const { widgets } = addTextWidget(normalizeDashboardDefinition({}), 'one');
    expect(canPlaceWidget(widgets, layout)).toBe(false);
    expect(canPlaceWidget(widgets, layout, 'one')).toBe(true);
    expect(canPlaceWidget(widgets, { ...layout, x: 4 })).toBe(true);
    expect(canPlaceWidget(widgets, { ...layout, x: 13 })).toBe(false);
  });

  it.each<{
    edge: WidgetResizeEdge;
    expected: { x: number; y: number; width: number; height: number };
  }>([
    { edge: 'n', expected: { x: 4, y: 5, width: 4, height: 2 } },
    { edge: 'ne', expected: { x: 4, y: 5, width: 5, height: 2 } },
    { edge: 'e', expected: { x: 4, y: 4, width: 5, height: 3 } },
    { edge: 'se', expected: { x: 4, y: 4, width: 5, height: 4 } },
    { edge: 's', expected: { x: 4, y: 4, width: 4, height: 4 } },
    { edge: 'sw', expected: { x: 5, y: 4, width: 3, height: 4 } },
    { edge: 'w', expected: { x: 5, y: 4, width: 3, height: 3 } },
    { edge: 'nw', expected: { x: 5, y: 5, width: 3, height: 2 } },
  ])('resizes the $edge edge while keeping the opposite edges fixed', ({ edge, expected }) => {
    expect(
      widgetResizeLayout({ x: 4, y: 4, width: 4, height: 3 }, edge, { x: 50, y: 50 }, 50),
    ).toEqual(expected);
  });

  it('clamps resize to the grid and minimum size but allows expansion below row 32', () => {
    const initial = { x: 1, y: 1, width: 4, height: 3 };
    expect(widgetResizeLayout(initial, 'nw', { x: -500, y: -500 }, 50)).toEqual({
      x: 0,
      y: 0,
      width: 5,
      height: 4,
    });
    expect(widgetResizeLayout(initial, 'e', { x: 1000, y: 0 }, 50)).toEqual({
      ...initial,
      width: 15,
    });
    expect(widgetResizeLayout(initial, 'se', { x: -500, y: -500 }, 50)).toEqual({
      ...initial,
      width: 2,
      height: 2,
    });
    expect(widgetResizeLayout({ ...initial, y: 30 }, 's', { x: 0, y: 250 }, 50)).toEqual({
      ...initial,
      y: 30,
      height: 8,
    });
  });
});

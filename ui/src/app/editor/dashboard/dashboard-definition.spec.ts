import {
  addTextWidget,
  dashboardGridRowCount,
  deleteTextWidget,
  firstFreeWidgetPosition,
  layoutsOverlap,
  moveTextWidget,
  normalizeDashboardDefinition,
  pixelDeltaToCells,
  resizeTextWidget,
  updateTextWidget,
  updateTextWidgetLayout,
} from './dashboard-definition';

describe('dashboard definition', () => {
  it('normalizes legacy definitions and rejects unknown versions', () => {
    expect(normalizeDashboardDefinition({})).toEqual({ version: 1, query: '', widgets: [] });
    expect(() => normalizeDashboardDefinition({ version: 2 })).toThrow(/unsupported/);
  });

  it('adds independent widget instances in the first available grid cells', () => {
    const first = addTextWidget(normalizeDashboardDefinition({}), 'one');
    const second = addTextWidget(first, 'two');
    expect(second.widgets.map(({ id }) => id)).toEqual(['one', 'two']);
    expect(second.widgets[0].layout).toEqual({ x: 0, y: 0, width: 4, height: 3 });
    expect(second.widgets[1].layout).toEqual({ x: 4, y: 0, width: 4, height: 3 });
    expect(second.widgets[0].properties).not.toBe(second.widgets[1].properties);
  });

  it('converts pixel deltas, checks overlap, and expands rows', () => {
    expect(pixelDeltaToCells(31, 20)).toBe(2);
    expect(
      layoutsOverlap({ x: 0, y: 0, width: 3, height: 3 }, { x: 2, y: 2, width: 2, height: 2 }),
    ).toBe(true);
    const definition = addTextWidget(normalizeDashboardDefinition({}), 'one');
    expect(firstFreeWidgetPosition(definition.widgets)).toEqual({ x: 4, y: 0 });
    expect(
      dashboardGridRowCount([
        { ...definition.widgets[0], layout: { x: 0, y: 34, width: 4, height: 3 } },
      ]),
    ).toBe(37);
  });

  it('updates and deletes only the requested instance', () => {
    const initial = addTextWidget(addTextWidget(normalizeDashboardDefinition({}), 'one'), 'two');
    const updated = updateTextWidget(initial, 'two', {
      query: 'label = "important"',
      properties: { text: 'Updated' },
    });
    expect(updated.widgets[0]).toEqual(initial.widgets[0]);
    expect(updated.widgets[1]).toMatchObject({
      query: 'label = "important"',
      properties: { text: 'Updated' },
    });
    expect(deleteTextWidget(updated, 'one').widgets.map(({ id }) => id)).toEqual(['two']);
  });

  it('moves and resizes widgets by cells while enforcing boundaries and collisions', () => {
    const initial = addTextWidget(addTextWidget(normalizeDashboardDefinition({}), 'one'), 'two');
    const moved = moveTextWidget(initial, 'two', 4, 4);
    expect(moved.widgets[1].layout).toMatchObject({ x: 4, y: 4 });
    expect(moveTextWidget(initial, 'two', 2, 0)).toBe(initial);
    const resized = resizeTextWidget(moved, 'two', 2, 5);
    expect(resized.widgets[1].layout).toMatchObject({ width: 2, height: 5 });
    expect(resizeTextWidget(resized, 'two', 1, 5)).toBe(resized);
    expect(resizeTextWidget(moved, 'two', 13, 3)).toBe(moved);
  });

  it('adds at a requested drop position without modifying the draft on a rejected drop', () => {
    const empty = normalizeDashboardDefinition({});
    const added = addTextWidget(empty, 'one', { x: 8, y: 30 });
    expect(added.widgets[0].layout).toEqual({ x: 8, y: 30, width: 4, height: 3 });
    expect(empty.widgets).toEqual([]);
    expect(addTextWidget(added, 'two', { x: 8, y: 30 })).toBe(added);
    expect(addTextWidget(added, 'two', { x: 14, y: 0 })).toBe(added);
  });

  it('commits a left/top resize as one complete layout while preserving the cached definition', () => {
    const initial = addTextWidget(normalizeDashboardDefinition({}), 'one', { x: 4, y: 4 });
    const resized = updateTextWidgetLayout(initial, 'one', { x: 2, y: 3, width: 6, height: 4 });
    expect(resized.widgets[0].layout).toEqual({ x: 2, y: 3, width: 6, height: 4 });
    expect(initial.widgets[0].layout).toEqual({ x: 4, y: 4, width: 4, height: 3 });
    expect(resized.widgets[0].properties).toEqual(initial.widgets[0].properties);
  });
});

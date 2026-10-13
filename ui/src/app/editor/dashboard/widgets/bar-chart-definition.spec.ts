import { createChartRuntime } from '@tanstack/charts/runtime';
import { defaultBarProperties } from '../metric-widget-definition';
import {
  BAR_LABEL_HEIGHT,
  BAR_LABEL_WIDTH,
  BarChartRow,
  barChartDefinition,
  barChartMargin,
  barChartMaximum,
  barChartRows,
  barChartTooltip,
  valueLabelsFit,
} from './bar-chart-definition';

const format = (value: number | null) => (value === null ? '—' : `${value} h`);

function rows(values: (number | null)[]): BarChartRow[] {
  return values.map((value, index) => ({
    key: `k${index}`,
    label: `Category ${index}`,
    value,
    entity: index === 0 ? { type: 'project', id: 'p1' } : null,
  }));
}

function render(
  values: (number | null)[],
  overrides: Partial<ReturnType<typeof defaultBarProperties>> = {},
  width = 480,
) {
  const properties = { ...defaultBarProperties(), ...overrides };
  const definition = barChartDefinition(rows(values), properties, format);
  return createChartRuntime<BarChartRow, string | number, string | number>().render(definition, {
    width,
    height: 260,
  });
}

describe('barChartDefinition', () => {
  it('draws one bar per category with a value, in both orientations', () => {
    for (const orientation of ['horizontal', 'vertical'] as const) {
      const scene = render([3, 5, null, 2], { orientation, showValues: false });
      expect(scene.points.map((point) => point.datum.key)).toEqual(['k0', 'k1', 'k3']);
    }
  });

  it('keeps zero as a bar of zero length and starts the value scale at zero', () => {
    const scene = render([0, 4], { showValues: false });
    const value = scene.scales['x'];
    expect(value.domain[0]).toBe(0);
    expect(scene.points.map((point) => point.xValue)).toEqual([0, 4]);
  });

  it('shows value labels only when requested', () => {
    const plain = render([3, 5], { showValues: false });
    const labelled = render([3, 5], { showValues: true });
    expect(plain.points).toHaveLength(2);
    expect(labelled.points).toHaveLength(4);
  });

  it('draws the grid only when showGrid is set', () => {
    const grid = (scene: ReturnType<typeof render>) =>
      scene.nodes.some((node) => JSON.stringify(node).includes('grid'));
    expect(grid(render([3, 5], { showGrid: true }))).toBe(true);
    expect(grid(render([3, 5], { showGrid: false }))).toBe(false);
  });

  it('reserves room for the category labels on the label side only', () => {
    expect(barChartMargin({ ...defaultBarProperties(), orientation: 'horizontal' })).toMatchObject({
      left: BAR_LABEL_WIDTH,
    });
    expect(barChartMargin({ ...defaultBarProperties(), orientation: 'vertical' })).toMatchObject({
      bottom: BAR_LABEL_HEIGHT,
    });
    const scene = render([3, 5], { orientation: 'horizontal' });
    expect(scene.margin.left).toBe(BAR_LABEL_WIDTH);
  });

  it('computes a maximum of at least one and builds rows from categories', () => {
    expect(barChartMaximum(rows([]))).toBe(1);
    expect(barChartMaximum(rows([3, null, 7]))).toBe(7);
    expect(
      barChartRows([
        {
          key: 'a',
          label: 'A',
          value: 2,
          entity: null,
          synthetic: false,
          memberCount: 1,
          start: null,
          end: null,
          status: 'complete',
        },
      ]),
    ).toEqual([{ key: 'a', label: 'A', value: 2, entity: null }]);
  });

  it('titles the tooltip with the category label and formats its value', () => {
    const tooltip = barChartTooltip(rows([3]), format);
    expect(tooltip.content([{ datum: rows([3])[0] }])).toEqual({
      title: 'Category 0',
      rows: [{ label: 'Value', value: '3 h' }],
    });
  });

  it('centres every band so that label links can be placed by index', () => {
    const scene = render([3, 5, 2, 4], { showValues: false });
    const step = scene.chart.height / 4;
    scene.points.forEach((point, index) => {
      expect(point.y).toBeCloseTo(scene.chart.y + (index + 0.5) * step, 1);
    });
    const vertical = render([3, 5, 2, 4], { orientation: 'vertical', showValues: false });
    const columns = vertical.chart.width / 4;
    vertical.points.forEach((point, index) => {
      expect(point.x).toBeCloseTo(vertical.chart.x + (index + 0.5) * columns, 1);
    });
  });

  it('drops value labels of vertical bars that would overlap', () => {
    const many = Array.from({ length: 30 }, () => 3);
    const properties = { ...defaultBarProperties(), orientation: 'vertical' as const };
    expect(valueLabelsFit(rows([3, 4, 5]), properties, format, 480)).toBe(true);
    expect(valueLabelsFit(rows(many), properties, format, 480)).toBe(false);
    expect(
      valueLabelsFit(rows(many), { ...properties, orientation: 'horizontal' }, format, 480),
    ).toBe(true);
    const scene = render(many, { orientation: 'vertical', showValues: true });
    expect(scene.points).toHaveLength(30);
  });
});

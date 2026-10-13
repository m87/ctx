import { createChartRuntime } from '@tanstack/charts/runtime';
import { DashboardRadarAxis } from '../../../../api/dashboard/dashboard.service';
import { defaultRadarProperties } from '../metric-widget-definition';
import {
  radarAxisLabel,
  radarChartDefinition,
  radarChartTooltip,
  radarEffectiveMaximum,
  radarFocus,
  radarProfileRows,
} from './radar-chart-definition';

const formatters = {
  value: (value: number | null) => (value === null ? '—' : `${value} h`),
  metric: 'Tracked time',
};

function axes(values: (number | null)[]): DashboardRadarAxis[] {
  return values.map((value, index) => ({
    key: `axis-${index}`,
    label: `Axis ${index}`,
    value,
    status: value === null ? 'unobserved' : 'complete',
    observedDays: value === null ? null : 3,
    entity: null,
  }));
}

function render(
  values: (number | null)[],
  overrides: Partial<ReturnType<typeof defaultRadarProperties>> = {},
  maximum = 10,
) {
  const properties = { ...defaultRadarProperties(), ...overrides };
  const definition = radarChartDefinition(axes(values), properties, maximum, formatters);
  return createChartRuntime<DashboardRadarAxis, string, number>().render(definition, {
    width: 320,
    height: 280,
  });
}

describe('radarEffectiveMaximum', () => {
  it('rounds automatic and extended scales without clipping values', () => {
    expect(radarEffectiveMaximum([2, 4.1], null, 'ms', null)).toBe(5);
    expect(radarEffectiveMaximum([12], 10, 'count', null)).toBe(20);
    expect(radarEffectiveMaximum([35], null, 'percent', null)).toBe(100);
  });

  it('extends the scale to a target above the data', () => {
    expect(radarEffectiveMaximum([2, 3], null, 'count', 8)).toBe(10);
  });
});

describe('radar profile', () => {
  it('starts after the first gap so that runs of values stay contiguous', () => {
    const rows = radarProfileRows(axes([1, null, 3, 4]));
    expect(rows.map((row) => row.value)).toEqual([3, 4, 1, null]);
    expect(radarProfileRows(axes([1, 2, 3])).map((row) => row.value)).toEqual([1, 2, 3]);
  });

  it('shortens weekday labels only for seven axes', () => {
    expect(radarAxisLabel('Monday', 7)).toBe('Mo');
    expect(radarAxisLabel('Monday', 5)).toBe('Monday');
  });
});

describe('radarChartDefinition', () => {
  it('has one profile point per axis with a value, at the same angles as the axes', () => {
    const scene = render([1, 2, 3, 4, 5, 6, 7]);
    const profile = scene.points.filter((point) => point.markId === 'profile');
    expect(profile.map((point) => point.xValue)).toEqual(
      Array.from({ length: 7 }, (_, index) => `axis-${index}`),
    );
  });

  it('does not create profile points for axes without a value', () => {
    const scene = render([1, null, 3, 4]);
    const profile = scene.points.filter((point) => point.markId === 'profile');
    expect(profile.map((point) => point.xValue)).toEqual(['axis-2', 'axis-3', 'axis-0']);
  });

  it('lets only the profile line take part in focus and keyboard navigation', () => {
    const scene = render([1, 2, 3], { showFill: true, showValues: true });
    expect(new Set(scene.points.map((point) => point.markId)).size).toBeGreaterThan(1);
    const profile = scene.points.filter((point) => point.markId === 'profile');
    const resolved = radarFocus.resolve(scene.points, {
      x: profile[0].x,
      y: profile[0].y,
      maxDistance: 48,
    });
    expect(resolved.length).toBeGreaterThan(0);
    expect(resolved.every((point) => point.markId === 'profile')).toBe(true);
    const visited = radarFocus.navigation(scene.points);
    expect(visited.map((point) => point.xValue)).toEqual(['axis-0', 'axis-1', 'axis-2']);
    expect(visited.every((point) => point.markId === 'profile')).toBe(true);
  });

  it('adds fill and value labels only when requested', () => {
    const base = render([1, 2, 3], { showFill: false, showValues: false });
    const rich = render([1, 2, 3], { showFill: true, showValues: true });
    const marks = (scene: typeof base) => new Set(scene.points.map((point) => point.markId));
    expect(marks(base)).toEqual(new Set(['profile']));
    expect(marks(rich).size).toBe(3);
  });

  it('does not fill an incomplete profile', () => {
    const complete = render([1, 2, 3], { showFill: true });
    const incomplete = render([1, null, 3], { showFill: true });
    const areaPoints = (scene: typeof complete) =>
      scene.points.filter((point) => point.markId !== 'profile');
    expect(areaPoints(complete).length).toBeGreaterThan(0);
    expect(areaPoints(incomplete)).toHaveLength(0);
  });

  it('describes an axis in the tooltip, including a partial period', () => {
    const tooltip = radarChartTooltip(formatters);
    const [axis] = axes([4]);
    expect(tooltip.content([{ datum: { ...axis, status: 'partial' } }])).toEqual({
      title: 'Axis 0',
      rows: [
        { label: 'Tracked time', value: '4 h' },
        { label: 'Status', value: 'partial period' },
      ],
    });
  });
});

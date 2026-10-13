import { createChartRuntime } from '@tanstack/charts/runtime';
import { defaultLineProperties } from '../metric-widget-definition';
import { LineChartData } from './line-chart-data';
import {
  COMPACT_CHART_WIDTH,
  bucketLabelFormatter,
  evenlySpaced,
  lineChartDefinition,
  lineChartTooltip,
} from './line-chart-definition';

const formatters = { value: (value: number) => `${value} h`, bucket: (bucket: string) => bucket };

function data(values: Record<string, (number | null)[]>, maximum: number): LineChartData {
  const buckets = ['2026-10-10', '2026-10-11', '2026-10-12', '2026-10-13'];
  return {
    maximum,
    rows: Object.entries(values).flatMap(([series, list]) =>
      list.map((value, index) => ({
        bucket: buckets[index],
        series,
        label: series.toUpperCase(),
        value,
      })),
    ),
  };
}

function render(
  chart: LineChartData,
  overrides: Partial<ReturnType<typeof defaultLineProperties>> = {},
  width = 640,
) {
  const runtime = createChartRuntime();
  const properties = { ...defaultLineProperties(), ...overrides };
  return runtime.render(lineChartDefinition(chart, properties, formatters), { width, height: 320 });
}

describe('lineChartDefinition', () => {
  it('starts the value axis at zero and covers the maximum including the target', () => {
    const scene = render(data({ a: [1, 2, 3, 4] }, 10), { targetValue: 10 });
    expect(scene.scales['y'].domain[0]).toBe(0);
    expect(scene.scales['y'].domain[1]).toBeGreaterThanOrEqual(10);
  });

  it('keeps null gaps and omits points for future or missing buckets', () => {
    const scene = render(data({ a: [1, null, 3, null] }, 3), { renderMode: 'linear' });
    expect(scene.points.filter((point) => point.group === 'a')).toHaveLength(2);
  });

  it('keeps zero as an observation', () => {
    const scene = render(data({ a: [0, 0, null, null] }, 1));
    expect(scene.points.map((point) => point.yValue)).toEqual([0, 0]);
  });

  it('shows an isolated observation as a point even without showPoints', () => {
    const scene = render(data({ a: [null, 5, null, null] }, 5), {
      renderMode: 'linear',
      showPoints: false,
    });
    expect(scene.points).toHaveLength(1);
  });

  it('creates one focus point per valid value for several series', () => {
    const scene = render(data({ a: [1, 2, 3, 4], b: [4, null, 2, 1] }, 4), {
      renderMode: 'smooth',
    });
    expect(scene.points).toHaveLength(7);
  });

  it('draws only points in points mode', () => {
    const scene = render(data({ a: [1, 2, 3, 4] }, 4), { renderMode: 'points' });
    expect(scene.points).toHaveLength(4);
  });

  it('formats tooltip rows with the series label, the value formatter and the series color', () => {
    const chart = data({ a: [1, 2, 3, 4], b: [4, 3, 2, 1] }, 4);
    const tooltip = lineChartTooltip(chart, formatters);
    const scene = render(chart);
    const content = tooltip.content(scene.points.filter((point) => point.xValue === '2026-10-11'));
    expect(content.title).toBe('2026-10-11');
    expect(content.rows.map((row) => [row.label, row.value])).toEqual([
      ['A', '2 h'],
      ['B', '3 h'],
    ]);
    expect(content.rows.map((row) => row.color)).toEqual(['var(--chart-1)', 'var(--chart-2)']);
  });

  it('hides the axes below the compact width', () => {
    const wide = render(data({ a: [1, 2, 3, 4] }, 4), {}, COMPACT_CHART_WIDTH + 40);
    const narrow = render(data({ a: [1, 2, 3, 4] }, 4), {}, COMPACT_CHART_WIDTH - 40);
    expect(wide.chart.height).toBeLessThan(320 - 10);
    expect(narrow.chart.height).toBeGreaterThan(wide.chart.height);
  });

  it('labels buckets by day, week and month in the dashboard time zone', () => {
    expect(bucketLabelFormatter('day', 'UTC', 'en-US')('2026-10-12T00:00:00Z')).toBe('Oct 12');
    expect(bucketLabelFormatter('month', 'UTC', 'en-US')('2026-10-01T00:00:00Z')).toBe('Oct');
    expect(bucketLabelFormatter('day', 'Pacific/Auckland', 'en-US')('2026-10-12T20:00:00Z')).toBe(
      'Oct 13',
    );
  });

  it('picks evenly spaced tick buckets that always include the first one', () => {
    const values = Array.from({ length: 30 }, (_, index) => index);
    expect(evenlySpaced(values, 8)).toEqual([0, 4, 8, 12, 16, 20, 24, 28]);
    expect(evenlySpaced(values.slice(0, 5), 8)).toEqual([0, 1, 2, 3, 4]);
    expect(evenlySpaced([], 4)).toEqual([]);
  });

  it('uses the picked buckets as the x axis ticks, formatted by the bucket formatter', () => {
    const buckets = Array.from({ length: 30 }, (_, index) => `2026-09-${index + 1}`);
    const chart: LineChartData = {
      maximum: 4,
      rows: buckets.map((bucket) => ({ bucket, series: 'a', label: 'A', value: 1 })),
    };
    const scene = render(chart, {}, 400);
    const ticks = scene.scales['x'].ticks;
    expect(ticks.map((tick) => tick.value)).toEqual(evenlySpaced(buckets, 5));
    expect(ticks.map((tick) => tick.label)).toEqual(ticks.map((tick) => String(tick.value)));
  });
});

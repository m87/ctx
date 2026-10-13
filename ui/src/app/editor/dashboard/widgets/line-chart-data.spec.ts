import { DashboardMetricResult } from '../../../../api/dashboard/dashboard.service';
import { lineChartRows } from './line-chart-data';

type SeriesResult = Extract<DashboardMetricResult, { kind: 'series' }>;

function seriesResult(
  buckets: { start: string; status: 'complete' | 'partial' | 'future' }[],
  series: { key: string; label: string; values: (number | null)[] }[],
): SeriesResult {
  return { series: { buckets, series } } as unknown as SeriesResult;
}

const buckets = [
  { start: '2026-10-10', status: 'complete' as const },
  { start: '2026-10-11', status: 'partial' as const },
  { start: '2026-10-12', status: 'future' as const },
];

describe('lineChartRows', () => {
  it('emits one row per series and bucket and keeps null gaps and zero values', () => {
    const { rows } = lineChartRows(
      seriesResult(buckets, [{ key: 'a', label: 'A', values: [0, null, null] }]),
      { targetValue: null },
    );
    expect(rows).toEqual([
      { bucket: '2026-10-10', series: 'a', label: 'A', value: 0 },
      { bucket: '2026-10-11', series: 'a', label: 'A', value: null },
      { bucket: '2026-10-12', series: 'a', label: 'A', value: null },
    ]);
  });

  it('turns future buckets into null even when the API sent a value', () => {
    const { rows, maximum } = lineChartRows(
      seriesResult(buckets, [{ key: 'a', label: 'A', values: [1, 2, 99] }]),
      { targetValue: null },
    );
    expect(rows.map((row) => row.value)).toEqual([1, 2, null]);
    expect(maximum).toBe(2);
  });

  it('extends the maximum to a target above the data', () => {
    const { maximum } = lineChartRows(
      seriesResult(buckets, [{ key: 'a', label: 'A', values: [1, 2, null] }]),
      { targetValue: 10 },
    );
    expect(maximum).toBe(10);
  });

  it('returns no rows and a zero maximum without series', () => {
    expect(lineChartRows(seriesResult(buckets, []), { targetValue: null })).toEqual({
      rows: [],
      maximum: 0,
    });
  });
});

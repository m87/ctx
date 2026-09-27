import {
  defaultGaugeProperties,
  defaultMetricProperties,
  formatMetric,
  metricSettingsError,
} from './metric-widget-definition';

describe('metric widget definitions', () => {
  it('provides independent defaults and formats empty and zero results', () => {
    const first = defaultMetricProperties();
    const second = defaultMetricProperties();
    expect(first).toEqual(second);
    expect(first).not.toBe(second);
    expect(formatMetric(0, first)).toBe('0 min');
    expect(formatMetric(null, first)).toBe('No data');
    expect(formatMetric(1_800_000, first)).toBe('30 min');
    expect(formatMetric(3_600_000, { ...first, format: 'hours' })).toBe('1 h');
  });

  it('validates gauge scales, goals, thresholds and whole-number counts', () => {
    const gauge = defaultGaugeProperties();
    expect(metricSettingsError(gauge, true)).toBe('');
    expect(metricSettingsError({ ...gauge, metric: 'missing' }, true)).toMatch(/metric/i);
    expect(metricSettingsError({ ...gauge, maximum: 0 }, true)).toMatch(/maximum/i);
    expect(metricSettingsError({ ...gauge, target: 200_000_000 }, true)).toMatch(/target/i);
    expect(
      metricSettingsError(
        {
          ...gauge,
          metric: 'sessions.count',
          format: 'number',
          precision: 0,
          minimum: 0.5,
        },
        true,
      ),
    ).toMatch(/whole numbers/i);
    expect(
      metricSettingsError(
        {
          ...gauge,
          thresholds: [
            { value: 20, tone: 'info' },
            { value: 10, tone: 'success' },
          ],
        },
        true,
      ),
    ).toMatch(/thresholds/i);
  });
});

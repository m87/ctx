import {
  radarClosedPolygon,
  radarEffectiveMaximum,
  radarPoint,
  radarSegments,
} from './radar-chart-geometry';

describe('radar chart geometry', () => {
  it('uses a shared maximum and starts the first axis at twelve o’clock', () => {
    expect(radarPoint(0, 4, 1)).toEqual({ x: 150, y: 44 });
    expect(radarPoint(1, 4, 0.5).x).toBeCloseTo(193);
  });

  it('does not bridge missing values or fill an incomplete profile', () => {
    expect(radarSegments([2, null, 4, 0], 4)).toHaveLength(2);
    expect(radarClosedPolygon([2, null, 4], 4)).toBeNull();
  });

  it('closes a complete profile while retaining zero at the center', () => {
    const polygon = radarClosedPolygon([0, 2, 4], 4);
    expect(polygon).toContain('M 150 130');
    expect(polygon?.endsWith(' Z')).toBe(true);
  });

  it('rounds automatic and extended scales without clipping values', () => {
    expect(radarEffectiveMaximum([2, 4.1], null, 'ms', null)).toBe(5);
    expect(radarEffectiveMaximum([12], 10, 'count', null)).toBe(20);
    expect(radarEffectiveMaximum([35], null, 'percent', null)).toBe(100);
  });
});

import { lineCoordinates, linePath } from './line-chart-geometry';

describe('line chart geometry', () => {
  it('keeps null gaps separate and leaves zero as an observation', () => {
    const points = lineCoordinates([2, null, 0, 4], 4);
    expect(linePath(points, 'linear')).toHaveLength(1);
    expect(linePath(points, 'linear')[0]).toContain('M 66.66666666666666 100 L 100 0');
    expect(points[2]?.value).toBe(0);
  });

  it('builds finite straight and smooth paths through every point without overshoot', () => {
    for (const values of [
      [0, 10, 0],
      [1, 2, 3, 4],
      [4, 3, 2, 1],
      [2, 2, 2],
    ]) {
      const points = lineCoordinates(values, Math.max(...values));
      const paths = linePath(points, 'smooth');
      expect(paths.every((path) => !/NaN|Infinity/.test(path))).toBe(true);
      expect(paths[0]).toContain(` ${points.at(-1)?.x} ${points.at(-1)?.y}`);
      expect(linePath(points, 'linear')[0]).toContain(' L ');
    }
  });

  it('does not draw a path for isolated points or points mode', () => {
    expect(linePath(lineCoordinates([3, null, 5], 5), 'smooth')).toEqual([]);
    expect(linePath(lineCoordinates([3, 5], 5), 'smooth')).toHaveLength(1);
  });
});

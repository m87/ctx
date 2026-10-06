export interface LinePoint {
  index: number;
  value: number;
  x: number;
  y: number;
}

export function lineCoordinates(
  values: readonly (number | null)[],
  maximum: number,
  width = 100,
  height = 100,
): (LinePoint | null)[] {
  const divisor = Math.max(1, values.length - 1);
  return values.map((value, index) =>
    value === null
      ? null
      : {
          index,
          value,
          x: (index / divisor) * width,
          y: height - (value / Math.max(maximum, Number.EPSILON)) * height,
        },
  );
}

export function linePath(
  points: readonly (LinePoint | null)[],
  mode: 'linear' | 'smooth',
): string[] {
  const segments: LinePoint[][] = [];
  let current: LinePoint[] = [];
  for (const point of points) {
    if (point) current.push(point);
    else if (current.length) {
      segments.push(current);
      current = [];
    }
  }
  if (current.length) segments.push(current);
  return segments
    .filter((segment) => segment.length > 1)
    .map((segment) => {
      let path = `M ${segment[0].x} ${segment[0].y}`;
      if (mode === 'linear')
        return (
          path +
          segment
            .slice(1)
            .map((point) => ` L ${point.x} ${point.y}`)
            .join('')
        );
      const slopes = segment
        .slice(0, -1)
        .map((point, index) => (segment[index + 1].y - point.y) / (segment[index + 1].x - point.x));
      const tangents = segment.map((_, index) => {
        if (index === 0) return slopes[0];
        if (index === segment.length - 1) return slopes.at(-1) ?? 0;
        const before = slopes[index - 1],
          after = slopes[index];
        return before * after <= 0 ? 0 : (2 * before * after) / (before + after);
      });
      for (let index = 0; index < segment.length - 1; index++) {
        const start = segment[index],
          end = segment[index + 1],
          delta = (end.x - start.x) / 3;
        let firstY = start.y + tangents[index] * delta,
          secondY = end.y - tangents[index + 1] * delta;
        firstY = Math.max(Math.min(start.y, end.y), Math.min(Math.max(start.y, end.y), firstY));
        secondY = Math.max(Math.min(start.y, end.y), Math.min(Math.max(start.y, end.y), secondY));
        path += ` C ${start.x + delta} ${firstY} ${end.x - delta} ${secondY} ${end.x} ${end.y}`;
      }
      return path;
    });
}

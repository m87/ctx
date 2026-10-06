export interface RadarPoint {
  x: number;
  y: number;
}

export function radarEffectiveMaximum(
  values: readonly (number | null)[],
  requested: number | null,
  unit: 'ms' | 'count' | 'percent',
  target: number | null,
): number {
  const maximum = Math.max(0, ...values.map((value) => value ?? 0), target ?? 0);
  const automatic = unit === 'percent' ? 100 : roundedMaximum(maximum);
  return Math.max(
    requested ?? automatic,
    maximum > (requested ?? automatic) ? roundedMaximum(maximum) : 0,
    1,
  );
}

function roundedMaximum(value: number): number {
  if (value <= 0) return 1;
  const magnitude = 10 ** Math.floor(Math.log10(value));
  const normalized = value / magnitude;
  const step = normalized <= 1 ? 1 : normalized <= 2 ? 2 : normalized <= 5 ? 5 : 10;
  return step * magnitude;
}

export function radarPoint(
  index: number,
  count: number,
  ratio: number,
  cx = 150,
  cy = 130,
  radius = 86,
): RadarPoint {
  const angle = -Math.PI / 2 + (Math.PI * 2 * index) / count;
  return { x: cx + Math.cos(angle) * radius * ratio, y: cy + Math.sin(angle) * radius * ratio };
}

export function radarSegments(values: readonly (number | null)[], maximum: number): string[] {
  if (!values.length) return [];
  return values.flatMap((value, index) => {
    const next = values[(index + 1) % values.length];
    if (value === null || next === null) return [];
    const a = radarPoint(index, values.length, Math.max(0, value / maximum));
    const b = radarPoint((index + 1) % values.length, values.length, Math.max(0, next / maximum));
    return [`M ${a.x} ${a.y} L ${b.x} ${b.y}`];
  });
}

export function radarClosedPolygon(
  values: readonly (number | null)[],
  maximum: number,
): string | null {
  if (!values.length || values.some((value) => value === null)) return null;
  return (
    values
      .map((value, index) => {
        const point = radarPoint(index, values.length, Math.max(0, (value ?? 0) / maximum));
        return `${index === 0 ? 'M' : 'L'} ${point.x} ${point.y}`;
      })
      .join(' ') + ' Z'
  );
}

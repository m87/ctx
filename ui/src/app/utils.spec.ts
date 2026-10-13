import { durationAsDecimalHours } from './utils';

const SECOND = 1_000_000_000;
const MINUTE = 60 * SECOND;
const HOUR = 60 * MINUTE;

describe('durationAsDecimalHours', () => {
  it.each([
    ['0', 0, '0h'],
    ['45 s', 45 * SECOND, '0h'],
    ['1 min', MINUTE, '0.02h'],
    ['15 min', 15 * MINUTE, '0.25h'],
    ['20 min', 20 * MINUTE, '0.33h'],
    ['40 min', 40 * MINUTE, '0.67h'],
    ['1 h', HOUR, '1h'],
    ['1 h 20 min', HOUR + 20 * MINUTE, '1.33h'],
    ['1 h 30 min', HOUR + 30 * MINUTE, '1.5h'],
    ['1 h 29 min 50 s', HOUR + 29 * MINUTE + 50 * SECOND, '1.48h'],
    ['7 h 45 min', 7 * HOUR + 45 * MINUTE, '7.75h'],
    ['8 h 10 min', 8 * HOUR + 10 * MINUTE, '8.17h'],
    ['37 h 5 min', 37 * HOUR + 5 * MINUTE, '37.08h'],
  ])('converts %s', (_label, duration, expected) => {
    expect(durationAsDecimalHours(duration)).toBe(expected);
  });

  it('uses a dot separator regardless of locale', () => {
    const toLocaleString = vi.spyOn(Number.prototype, 'toLocaleString').mockReturnValue('1,5');
    expect(durationAsDecimalHours(HOUR + 30 * MINUTE)).toBe('1.5h');
    toLocaleString.mockRestore();
  });
});

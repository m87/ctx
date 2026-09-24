import {
  createTimelineMarks,
  minuteWithinDay,
  normalizeTimelineRangeMode,
  resolveTimelineRange,
  timelinePosition,
  timelineWidth,
} from './timeline-range';

describe('timeline range', () => {
  it('uses the full day by default', () => {
    expect(normalizeTimelineRangeMode(undefined)).toBe('full-day');
    expect(
      resolveTimelineRange([{ startMinutes: 9 * 60, endMinutes: 10 * 60 }], 'full-day'),
    ).toEqual({
      startMinutes: 0,
      endMinutes: 24 * 60,
    });
  });

  it('fits the range to the first start and last end', () => {
    const range = resolveTimelineRange(
      [
        { startMinutes: 13 * 60, endMinutes: 15 * 60 + 30 },
        { startMinutes: 8 * 60 + 15, endMinutes: 10 * 60 },
      ],
      'intervals',
    );

    expect(range).toEqual({ startMinutes: 8 * 60 + 15, endMinutes: 15 * 60 + 30 });
    expect(timelinePosition(8 * 60 + 15, range)).toBe(0);
    expect(timelinePosition(15 * 60 + 30, range)).toBe(100);
    expect(timelineWidth(13 * 60, 15 * 60 + 30, range)).toBeCloseTo(34.48, 2);
  });

  it('falls back to the full day when there are no intervals', () => {
    expect(resolveTimelineRange([], 'intervals')).toEqual({ startMinutes: 0, endMinutes: 1440 });
  });

  it('maps day boundaries in the selected time zone', () => {
    expect(minuteWithinDay('2026-09-24T06:15:00Z', '2026-09-24', 'Europe/Warsaw')).toBe(
      8 * 60 + 15,
    );
    expect(minuteWithinDay('2026-09-24T22:00:00Z', '2026-09-24', 'Europe/Warsaw')).toBe(1440);
  });

  it('keeps range endpoints in fitted timeline marks', () => {
    expect(createTimelineMarks({ startMinutes: 8 * 60 + 15, endMinutes: 17 * 60 + 45 })).toEqual([
      { minute: 8 * 60 + 15, label: '8:15' },
      { minute: 9 * 60, label: '9:00' },
      { minute: 12 * 60, label: '12:00' },
      { minute: 15 * 60, label: '15:00' },
      { minute: 17 * 60 + 45, label: '17:45' },
    ]);
  });
});

import { DateTime } from 'luxon';

export const timelineRangeSettingKey = 'client.general.timelineRange';

export type TimelineRangeMode = 'full-day' | 'intervals';

export interface TimelineRange {
  startMinutes: number;
  endMinutes: number;
}

export interface TimelineMark {
  minute: number;
  label: string;
}

const fullDayRange: TimelineRange = {
  startMinutes: 0,
  endMinutes: 24 * 60,
};

export function normalizeTimelineRangeMode(value: string | undefined): TimelineRangeMode {
  return value === 'intervals' ? 'intervals' : 'full-day';
}

export function minuteWithinDay(
  instant: string | null,
  day: string,
  timeZone: string,
): number | null {
  if (!instant) {
    return null;
  }

  const parsed = DateTime.fromISO(instant, { setZone: true }).setZone(timeZone);
  if (!parsed.isValid) {
    return null;
  }

  const instantDay = parsed.toFormat('yyyy-MM-dd');
  if (instantDay < day) {
    return fullDayRange.startMinutes;
  }
  if (instantDay > day) {
    return fullDayRange.endMinutes;
  }

  return parsed.hour * 60 + parsed.minute + parsed.second / 60 + parsed.millisecond / 60_000;
}

export function resolveTimelineRange(
  intervals: TimelineRange[],
  mode: TimelineRangeMode,
): TimelineRange {
  if (mode === 'full-day') {
    return fullDayRange;
  }

  const validIntervals = intervals.filter(
    (interval) =>
      Number.isFinite(interval.startMinutes) &&
      Number.isFinite(interval.endMinutes) &&
      interval.endMinutes > interval.startMinutes,
  );
  if (validIntervals.length === 0) {
    return fullDayRange;
  }

  const startMinutes = Math.max(
    fullDayRange.startMinutes,
    Math.min(...validIntervals.map((interval) => interval.startMinutes)),
  );
  const endMinutes = Math.min(
    fullDayRange.endMinutes,
    Math.max(...validIntervals.map((interval) => interval.endMinutes)),
  );

  return endMinutes > startMinutes ? { startMinutes, endMinutes } : fullDayRange;
}

export function createTimelineMarks(range: TimelineRange): TimelineMark[] {
  if (
    range.startMinutes === fullDayRange.startMinutes &&
    range.endMinutes === fullDayRange.endMinutes
  ) {
    return [0, 6 * 60, 12 * 60, 18 * 60, 24 * 60].map(createTimelineMark);
  }

  const largestStep = 720;
  const steps = [15, 30, 60, 120, 180, 240, 360, largestStep];
  const step =
    steps.find((candidate) => interiorMarks(range, candidate).length <= 3) ?? largestStep;
  const minutes = [range.startMinutes, ...interiorMarks(range, step), range.endMinutes].filter(
    (minute, index, values) => index === 0 || Math.abs(minute - values[index - 1]) > 0.01,
  );

  return minutes.map(createTimelineMark);
}

export function timelinePosition(minute: number, range: TimelineRange): number {
  const duration = range.endMinutes - range.startMinutes;
  if (duration <= 0) {
    return 0;
  }
  return ((minute - range.startMinutes) / duration) * 100;
}

export function timelineWidth(
  startMinutes: number,
  endMinutes: number,
  range: TimelineRange,
): number {
  const start = Math.max(startMinutes, range.startMinutes);
  const end = Math.min(endMinutes, range.endMinutes);
  return Math.max(timelinePosition(end, range) - timelinePosition(start, range), 0);
}

function interiorMarks(range: TimelineRange, step: number): number[] {
  const result: number[] = [];
  let minute = Math.ceil(range.startMinutes / step) * step;
  if (Math.abs(minute - range.startMinutes) <= 0.01) {
    minute += step;
  }

  while (minute < range.endMinutes - 0.01) {
    result.push(minute);
    minute += step;
  }
  return result;
}

function createTimelineMark(minute: number): TimelineMark {
  const roundedMinute = Math.min(Math.max(Math.floor(minute), 0), 24 * 60);
  const hour = Math.floor(roundedMinute / 60);
  const minutes = roundedMinute % 60;
  return {
    minute,
    label: `${hour}:${minutes.toString().padStart(2, '0')}`,
  };
}

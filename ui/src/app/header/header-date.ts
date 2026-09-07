import { DateTime } from 'luxon';

const isoDateFormat = 'yyyy-MM-dd';

export function formatHeaderDate(date: string): string {
  const parsed = parseHeaderDate(date);
  return parsed?.toFormat('d LLL yyyy') ?? date;
}

export function headerDateToLocalDate(date: string): Date | undefined {
  const parsed = parseHeaderDate(date);
  return parsed ? new Date(parsed.year, parsed.month - 1, parsed.day) : undefined;
}

function parseHeaderDate(date: string): DateTime | null {
  const parsed = DateTime.fromFormat(date, isoDateFormat, { zone: 'UTC' });
  return parsed.isValid && parsed.toFormat(isoDateFormat) === date ? parsed : null;
}

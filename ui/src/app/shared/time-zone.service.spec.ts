import { instantInTimeZone, inputValueToUTC, resolveTimeZone } from './time-zone.service';

describe('time zone conversion', () => {
  it('renders the same instant in the selected time zone', () => {
    const instant = '2026-08-01T18:00:00Z';

    expect(instantInTimeZone(instant, 'Asia/Tokyo').toFormat('yyyy-MM-dd HH:mm')).toBe(
      '2026-08-02 03:00',
    );
    expect(instantInTimeZone(instant, 'Europe/Warsaw').toFormat('yyyy-MM-dd HH:mm')).toBe(
      '2026-08-01 20:00',
    );
  });

  it('converts a datetime-local value in the selected zone to UTC', () => {
    expect(inputValueToUTC('2026-08-02T03:00', 'Asia/Tokyo')).toBe('2026-08-01T18:00:00.000Z');
  });

  it('uses a valid explicit IANA preference', () => {
    expect(resolveTimeZone('America/New_York')).toBe('America/New_York');
  });

  it('rejects a local time skipped by a DST transition', () => {
    expect(inputValueToUTC('2026-03-08T02:30', 'America/New_York')).toBeNull();
  });
});

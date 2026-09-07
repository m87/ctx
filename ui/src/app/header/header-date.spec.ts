import { formatHeaderDate, headerDateToLocalDate } from './header-date';

describe('header date navigation', () => {
  it('formats the selected date for the top bar', () => {
    expect(formatHeaderDate('2026-09-06')).toBe('6 Sep 2026');
  });

  it('creates a local calendar date without changing its day', () => {
    const date = headerDateToLocalDate('2026-09-06');

    expect(date?.getFullYear()).toBe(2026);
    expect(date?.getMonth()).toBe(8);
    expect(date?.getDate()).toBe(6);
  });
});

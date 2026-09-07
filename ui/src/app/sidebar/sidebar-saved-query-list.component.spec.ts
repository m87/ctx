import { isSavedQuerySectionVisible } from './sidebar-saved-query-list.component';

describe('saved query sidebar section', () => {
  it('is hidden after an empty list finishes loading', () => {
    expect(isSavedQuerySectionVisible(false, 0)).toBe(false);
  });

  it('is visible while loading or when saved queries exist', () => {
    expect(isSavedQuerySectionVisible(true, 0)).toBe(true);
    expect(isSavedQuerySectionVisible(false, 1)).toBe(true);
  });
});

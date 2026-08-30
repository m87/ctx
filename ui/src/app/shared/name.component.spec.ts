import { parseTagNames, resolveTags } from './name.component';

describe('name tags', () => {
  it('keeps existing tag IDs and creates objects for new names', () => {
    expect(resolveTags('existing, new, existing', [{ id: 'tag-1', name: 'existing' }])).toEqual([
      { id: 'tag-1', name: 'existing' },
      { id: '', name: 'new' },
    ]);
  });

  it('trims and deduplicates names', () => {
    expect(parseTagNames(' first, second, first, , second ')).toEqual(['first', 'second']);
  });
});

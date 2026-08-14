import { parseProjectPicker } from './search-project-picker';

describe('search project picker', () => {
  it('does not activate without a hash', () => {
    expect(parseProjectPicker('review pull request')).toBeNull();
  });

  it('shows every project for a trailing hash', () => {
    expect(parseProjectPicker('review pull request #')).toEqual({
      contextName: 'review pull request',
      projectQuery: '',
    });
  });

  it('extracts the context name and project query', () => {
    expect(parseProjectPicker('review pull request # client app ')).toEqual({
      contextName: 'review pull request',
      projectQuery: 'client app',
    });
  });

  it('uses a bare hash to request every project', () => {
    expect(parseProjectPicker('#')).toEqual({
      contextName: '',
      projectQuery: '',
    });
  });
});

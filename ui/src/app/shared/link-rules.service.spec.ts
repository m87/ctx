import { describe, expect, it } from 'vitest';
import { resolveLinkTemplate } from './link-rules.service';

function match(expression: RegExp, text: string): RegExpExecArray {
  const result = expression.exec(text);
  if (result === null) {
    throw new Error(`Expected ${expression} to match ${text}`);
  }
  return result;
}

describe('resolveLinkTemplate', () => {
  it('resolves capture groups and URL-encoded context values', () => {
    const result = match(/([A-Z]+-\d+)/, 'Work on CTX-42');

    expect(
      resolveLinkTemplate(
        'https://jira.test/$1?name=${name};time=${duration};startDate=${start}',
        result,
        {
          name: 'Work on CTX-42',
          duration: '1h 30m',
          start: '2026-09-08T08:15:00Z',
        },
      ),
    ).toBe(
      'https://jira.test/CTX-42?name=Work%20on;time=1h%2030m;startDate=2026-09-08T08%3A15%3A00Z',
    );
  });

  it('uses the context name without the matched issue key', () => {
    const result = match(/([A-Z]+-\d+)/, 'TAR-123 opis mojego ticketa');

    expect(
      resolveLinkTemplate('/$1?description=${name}', result, {
        name: 'TAR-123 opis mojego ticketa',
      }),
    ).toBe('/TAR-123?description=opis%20mojego%20ticketa');
  });

  it('supports numeric and boolean values', () => {
    const result = match(/CTX/, 'CTX');

    expect(
      resolveLinkTemplate('/${durationValue}/${archived}', result, {
        durationValue: 90_000_000_000,
        archived: false,
      }),
    ).toBe('/90000000000/false');
  });

  it('resolves nested context values', () => {
    const result = match(/CTX/, 'CTX');

    expect(
      resolveLinkTemplate('/projects/${project.id}?name=${project.name}', result, {
        project: { id: 'project-1', name: 'Internal tools' },
      }),
    ).toBe('/projects/project-1?name=Internal%20tools');
  });

  it('preserves a placeholder when its value is unavailable', () => {
    const result = match(/CTX/, 'CTX');

    expect(resolveLinkTemplate('/${start}/$&/$$', result, {})).toBe('/${start}/CTX/$');
  });
});

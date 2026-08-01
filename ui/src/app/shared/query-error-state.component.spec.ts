import { HttpErrorResponse } from '@angular/common/http';
import { getQueryErrorContent } from './query-error-state.component';

describe('query error state', () => {
  it('describes a server connection failure', () => {
    expect(getQueryErrorContent(httpError(0), 'daily summary')).toEqual({
      icon: 'lucideServerOff',
      title: "Can't reach the server",
      description: "We couldn't load daily summary. Check that ctx is running and try again.",
      retryable: true,
    });
  });

  it('treats a paused query as a connection failure', () => {
    expect(getQueryErrorContent(null, 'workspace', true).title).toBe("Can't reach the server");
  });

  it('recognizes an unavailable development proxy', () => {
    expect(getQueryErrorContent(httpError(500), 'workspaces').title).toBe("Can't reach the server");
  });

  it('does not offer a retry for a missing resource', () => {
    const content = getQueryErrorContent(
      httpError(404, { code: 'CONTEXT_NOT_FOUND', description: 'Context not found' }),
      'context',
    );

    expect(content.title).toBe('Context not found');
    expect(content.retryable).toBe(false);
  });
});

function httpError(status: number, error?: unknown): HttpErrorResponse {
  return new HttpErrorResponse({ error, status, statusText: 'Test error' });
}

import { provideHttpClient } from '@angular/common/http';
import { provideHttpClientTesting } from '@angular/common/http/testing';
import { TestBed } from '@angular/core/testing';
import { provideRouter } from '@angular/router';
import { provideStore } from '@ngxs/store';
import { provideTanStackQuery, QueryClient } from '@tanstack/angular-query-experimental';
import { ContextQueryResult } from '../../api/query/query.service';
import { QuerySummaryComponent } from './query-summary.component';

const MINUTE = 60 * 1_000_000_000;

beforeEach(() => {
  TestBed.configureTestingModule({
    providers: [
      provideRouter([]),
      provideHttpClient(),
      provideHttpClientTesting(),
      provideStore([]),
      provideTanStackQuery(new QueryClient()),
    ],
  });
});

function render(inputs: { result?: ContextQueryResult | null; loading?: boolean }) {
  const fixture = TestBed.createComponent(QuerySummaryComponent);
  fixture.componentRef.setInput('result', inputs.result ?? null);
  fixture.componentRef.setInput('loading', inputs.loading ?? false);
  fixture.detectChanges();
  return fixture.nativeElement as HTMLElement;
}

function resultWithTotal(totalDuration: number): ContextQueryResult {
  return { contexts: [], contextStats: [], totalDuration } as unknown as ContextQueryResult;
}

describe('QuerySummaryComponent total tracked tile', () => {
  it('shows decimal hours under the total', () => {
    const text = render({ result: resultWithTotal(90 * MINUTE) }).textContent ?? '';

    expect(text).toContain('1h 30m');
    expect(text).toContain('1.5h');
  });

  it('shows zero hours for a zero total', () => {
    expect(render({ result: resultWithTotal(0) }).textContent).toContain('0h');
  });

  it('shows zero hours without a result', () => {
    const text = render({ result: null }).textContent ?? '';

    expect(text).toContain('Total tracked');
    expect(text).toContain('0h');
  });

  it('shows neither tile nor decimal line while loading', () => {
    const text = render({ result: resultWithTotal(90 * MINUTE), loading: true }).textContent ?? '';

    expect(text).not.toContain('Total tracked');
    expect(text).not.toContain('1.5h');
  });
});

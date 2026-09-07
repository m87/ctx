import { provideHttpClient } from '@angular/common/http';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { TestBed } from '@angular/core/testing';
import { QueryService } from './query.service';

describe('QueryService', () => {
  let http: HttpTestingController;
  let service: QueryService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [provideHttpClient(), provideHttpClientTesting()],
    });
    http = TestBed.inject(HttpTestingController);
    service = TestBed.inject(QueryService);
  });

  afterEach(() => http.verify());

  it('executes a query in its workspace and requests a summary result', () => {
    const input = { workspaceId: 'workspace-1', query: 'tag = focus' };

    service.execute(input).subscribe();

    const request = http.expectOne('/api/query/result');
    expect(request.request.method).toBe('POST');
    expect(request.request.body).toEqual(input);
    request.flush({
      ...input,
      contexts: [],
      contextStats: [],
      totalDuration: 0,
      totalSessions: 0,
    });
  });
});

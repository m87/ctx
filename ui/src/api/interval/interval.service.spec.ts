import { provideHttpClient } from '@angular/common/http';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { TestBed } from '@angular/core/testing';
import { Store } from '@ngxs/store';
import { Interval, IntervalService } from './interval.service';

describe('IntervalService', () => {
  let http: HttpTestingController;
  let service: IntervalService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        provideHttpClient(),
        provideHttpClientTesting(),
        { provide: Store, useValue: { selectSnapshot: () => 'workspace-1' } },
      ],
    });

    http = TestBed.inject(HttpTestingController);
    service = TestBed.inject(IntervalService);
  });

  afterEach(() => http.verify());

  it('posts a new interval to the canonical collection URL', () => {
    const interval: Interval = {
      id: '',
      contextId: 'context-1',
      start: '2026-08-14T08:00:00.000Z',
      end: '2026-08-14T08:30:00.000Z',
      duration: 0,
      workspaceId: 'workspace-1',
    };

    service.createInterval(interval).subscribe();

    const request = http.expectOne('/api/interval/');
    expect(request.request.method).toBe('POST');
    expect(request.request.body).toEqual(interval);
    request.flush(interval);
  });
});

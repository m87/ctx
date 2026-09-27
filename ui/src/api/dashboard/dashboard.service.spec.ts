import { provideHttpClient } from '@angular/common/http';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { TestBed } from '@angular/core/testing';
import { DashboardService } from './dashboard.service';

describe('DashboardService', () => {
  let http: HttpTestingController;
  let service: DashboardService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [provideHttpClient(), provideHttpClientTesting()],
    });
    http = TestBed.inject(HttpTestingController);
    service = TestBed.inject(DashboardService);
  });

  afterEach(() => http.verify());

  it('lists dashboards for a workspace', () => {
    service.list('workspace 1').subscribe();

    const request = http.expectOne(
      (candidate) =>
        candidate.url === '/api/dashboard/' &&
        candidate.params.get('workspaceId') === 'workspace 1',
    );
    expect(request.request.method).toBe('GET');
    request.flush([]);
  });

  it('creates, updates, gets and deletes a dashboard', () => {
    const input = {
      workspaceId: 'workspace-1',
      type: 'custom' as const,
      name: 'Focus',
      definition: {},
    };
    const dashboard = { id: 'dashboard-1', active: true, ...input };

    service.create(input).subscribe();
    const createRequest = http.expectOne('/api/dashboard/');
    expect(createRequest.request.method).toBe('POST');
    expect(createRequest.request.body).toEqual(input);
    createRequest.flush(dashboard);

    service.update(dashboard).subscribe();
    const updateRequest = http.expectOne('/api/dashboard/dashboard-1');
    expect(updateRequest.request.method).toBe('PUT');
    expect(updateRequest.request.body).toEqual(dashboard);
    updateRequest.flush(dashboard);

    service.get('dashboard-1').subscribe();
    const getRequest = http.expectOne('/api/dashboard/dashboard-1');
    expect(getRequest.request.method).toBe('GET');
    getRequest.flush(dashboard);

    service.delete('dashboard-1').subscribe();
    const deleteRequest = http.expectOne('/api/dashboard/dashboard-1');
    expect(deleteRequest.request.method).toBe('DELETE');
    deleteRequest.flush(null);
  });

  it('evaluates dashboard metrics with the selected view scope', () => {
    const body = {
      workspaceId: 'workspace-1',
      definition: { version: 1, query: '', widgets: [] },
      timeZone: 'Europe/Warsaw',
      firstDayOfWeek: 1,
      scope: { projectId: 'project-1' },
    };
    service.evaluate(body).subscribe();
    const request = http.expectOne('/api/dashboard/evaluate');
    expect(request.request.method).toBe('POST');
    expect(request.request.body).toEqual(body);
    const response = { asOf: '2026-09-27T16:00:00Z', widgets: { kpi: { value: 0 } } };
    request.flush(response);
  });
});

import { provideHttpClient } from '@angular/common/http';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { TestBed } from '@angular/core/testing';
import { SavedQueryService } from './saved-query.service';

describe('SavedQueryService', () => {
  let http: HttpTestingController;
  let service: SavedQueryService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [provideHttpClient(), provideHttpClientTesting()],
    });
    http = TestBed.inject(HttpTestingController);
    service = TestBed.inject(SavedQueryService);
  });

  afterEach(() => http.verify());

  it('lists saved queries for a workspace', () => {
    service.list('workspace 1').subscribe();

    const request = http.expectOne('/api/query/saved/?workspaceId=workspace+1');
    expect(request.request.method).toBe('GET');
    request.flush([]);
  });

  it('creates and gets a saved query', () => {
    const input = { workspaceId: 'workspace-1', name: 'Focus', query: 'tag = focus' };
    const savedQuery = { id: 'query-1', ...input };

    service.create(input).subscribe();
    const createRequest = http.expectOne('/api/query/saved/');
    expect(createRequest.request.method).toBe('POST');
    expect(createRequest.request.body).toEqual(input);
    createRequest.flush(savedQuery);

    service.get('query-1').subscribe();
    const getRequest = http.expectOne('/api/query/saved/query-1');
    expect(getRequest.request.method).toBe('GET');
    getRequest.flush(savedQuery);
  });

  it('deletes a saved query', () => {
    service.delete('query-1').subscribe();

    const request = http.expectOne('/api/query/saved/query-1');
    expect(request.request.method).toBe('DELETE');
    request.flush(null);
  });
});

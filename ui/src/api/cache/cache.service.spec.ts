import { TestBed } from '@angular/core/testing';
import { QueryClient } from '@tanstack/angular-query-experimental';
import { dashboardQueryKeys } from '../dashboard/dashboard.queries';
import { CacheService } from './cache.service';

describe('CacheService dashboard metric invalidation', () => {
  it('invalidates evaluations after interval, context metadata and project changes', async () => {
    const queryClient = new QueryClient();
    const invalidate = vi.spyOn(queryClient, 'invalidateQueries').mockResolvedValue();
    TestBed.configureTestingModule({
      providers: [{ provide: QueryClient, useValue: queryClient }],
    });
    const cache = TestBed.inject(CacheService);

    await cache.afterIntervalChange('context-1');
    expect(invalidate).toHaveBeenCalledWith({ queryKey: dashboardQueryKeys.evaluations() });
    invalidate.mockClear();
    await cache.afterContextMetadataChange('context-1');
    expect(invalidate).toHaveBeenCalledWith({ queryKey: dashboardQueryKeys.evaluations() });
    invalidate.mockClear();
    await cache.afterProjectChange('workspace-1');
    expect(invalidate).toHaveBeenCalledWith({ queryKey: dashboardQueryKeys.evaluations() });
    queryClient.clear();
  });
});

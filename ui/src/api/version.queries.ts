import { inject, Injectable } from '@angular/core';
import { lastValueFrom } from 'rxjs';
import { VersionService } from './version.service';

export const versionQueryKeys = {
  all: ['version'] as const,
  current: () => [...versionQueryKeys.all] as const,
};

@Injectable({
  providedIn: 'root',
})
export class VersionQueries {
  private readonly versionService = inject(VersionService);

  version() {
    return {
      queryKey: versionQueryKeys.current(),
      queryFn: () => lastValueFrom(this.versionService.getVersion()),
    };
  }
}

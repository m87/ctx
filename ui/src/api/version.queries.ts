import { inject, Injectable } from '@angular/core';
import { lastValueFrom } from 'rxjs';
import { VersionService } from './version.service';

@Injectable({
  providedIn: 'root',
})
export class VersionQueries {
  static readonly key = ['version'];
  private versionService = inject(VersionService);

  version() {
    return {
      queryKey: VersionQueries.key,
      queryFn: () => lastValueFrom(this.versionService.getVersion()),
    };
  }
}

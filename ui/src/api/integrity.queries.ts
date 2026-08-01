import { inject, Injectable } from "@angular/core";
import { lastValueFrom } from "rxjs";
import { IntegrityService } from "./integrity.service";

export const integrityQueryKeys = {
  all: ['integrity'] as const,
  integrity: () => [...integrityQueryKeys.all, 'integrity'] as const,
  integrityContexts: () => [...integrityQueryKeys.all, 'integrity-contexts'] as const,
};


@Injectable({ providedIn: "root" })
export class IntegrityQueries {
 private readonly integrityService = inject(IntegrityService);

  integrity() {
    return {
      queryKey: integrityQueryKeys.integrity(),
      queryFn: () => lastValueFrom(this.integrityService.checkIntegrity()),
      enabled: false,
    };
  }

  integrityContexts(enabled: boolean) {
    return {
      queryKey: integrityQueryKeys.integrityContexts(),
      queryFn: () => lastValueFrom(this.integrityService.getIntegrityContexts()),
      enabled,
    };
  }
}

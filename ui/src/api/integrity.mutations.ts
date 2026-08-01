import { inject, Injectable } from '@angular/core';
import { mutationOptions, QueryClient } from '@tanstack/angular-query-experimental';
import { toast } from 'ngx-sonner';
import { lastValueFrom } from 'rxjs';
import { IntegrityService } from './integrity.service';
import { integrityQueryKeys } from './integrity.queries';

@Injectable({ providedIn: 'root' })
export class IntegrityMutations {
  private readonly integrityService = inject(IntegrityService);
  private readonly queryClient = inject(QueryClient);

  repairIntegrity() {
    return mutationOptions({
      mutationFn: () => lastValueFrom(this.integrityService.repairIntegrity()),
      onSuccess: (result) => {
        this.queryClient.setQueryData(integrityQueryKeys.integrity(), result.report);
        toast.success(`Integrity repair completed. Repaired ${result.repairedCount} records.`);
      },
    });
  }

  checkIntegrity() {
    return mutationOptions({
      mutationFn: () => lastValueFrom(this.integrityService.checkIntegrity()),
      onSuccess: (report) => {
        this.queryClient.setQueryData(integrityQueryKeys.integrity(), report);
        if (report.healthy) {
          toast.success('Data integrity check passed. No issues found.');
        } else {
          toast.warning(`Data integrity check found ${report.issues.length} issues.`);
        }
      },
    });
  }
}

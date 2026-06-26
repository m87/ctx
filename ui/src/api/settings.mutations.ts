import { inject, Injectable } from '@angular/core';
import { mutationOptions, QueryClient } from '@tanstack/angular-query-experimental';
import { lastValueFrom } from 'rxjs';
import { Settings, SettingsService } from './settings.service';
import { SettingsQueries } from './settings.queries';
import { toastError } from './error';
import { toast } from 'ngx-sonner';

@Injectable({ providedIn: 'root' })
export class SettingsMutations {
  private settingsService = inject(SettingsService);
  private queryClient = inject(QueryClient);

  save() {
    return mutationOptions({
      mutationFn: (settings: Settings) =>
        lastValueFrom(this.settingsService.saveSettings(settings)),
      onSuccess: (_, settings) => {
        this.queryClient.setQueryData([...SettingsQueries.key, 'settings'], settings);
        this.queryClient.invalidateQueries({ queryKey: [...SettingsQueries.key, 'settings'] });

        Object.keys(settings).forEach((key) => {
          this.queryClient.setQueryData([...SettingsQueries.key, 'setting', key], settings[key]);
          this.queryClient.invalidateQueries({
            queryKey: [...SettingsQueries.key, 'setting', key],
          });
        });
      },
      onError(error) {
        toastError(error);
      },
    });
  }

  repairIntegrity() {
    return mutationOptions({
      mutationFn: () => lastValueFrom(this.settingsService.repairIntegrity()),
      onSuccess: (result) => {
        this.queryClient.setQueryData([...SettingsQueries.key, 'integrity'], result.report);
        toast.success(`Integrity repair completed. Repaired ${result.repairedCount} records.`);
      },
      onError(error) {
        toastError(error);
      },
    });
  }

  checkIntegrity() {
    return mutationOptions({
      mutationFn: () => lastValueFrom(this.settingsService.checkIntegrity()),
      onSuccess: (report) => {
        this.queryClient.setQueryData([...SettingsQueries.key, 'integrity'], report);
        if (report.healthy) {
          toast.success('Data integrity check passed. No issues found.');
        } else {
          toast.warning(`Data integrity check found ${report.issues.length} issues.`);
        }
      },
      onError(error) {
        toastError(error);
      },
    });
  }
}

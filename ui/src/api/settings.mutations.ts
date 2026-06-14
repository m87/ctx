import { inject, Injectable } from '@angular/core';
import { mutationOptions, QueryClient } from '@tanstack/angular-query-experimental';
import { lastValueFrom } from 'rxjs';
import { Settings, SettingsService } from './settings.service';
import { SettingsQueries } from './settings.queries';
import { toastError } from './error';

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
}

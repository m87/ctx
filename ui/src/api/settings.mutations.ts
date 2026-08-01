import { inject, Injectable } from '@angular/core';
import { mutationOptions, QueryClient } from '@tanstack/angular-query-experimental';
import { lastValueFrom } from 'rxjs';
import { CacheService } from './cache.service';
import { settingsQueryKeys } from './settings.queries';
import { Settings, SettingsService } from './settings.service';

@Injectable({ providedIn: 'root' })
export class SettingsMutations {
  private readonly settingsService = inject(SettingsService);
  private readonly queryClient = inject(QueryClient);
  private readonly cache = inject(CacheService);

  save() {
    return mutationOptions({
      mutationFn: (settings: Settings) =>
        lastValueFrom(this.settingsService.saveSettings(settings)),
      onSuccess: (_, settings) => {
        this.queryClient.setQueryData(settingsQueryKeys.settings(), settings);

        for (const [key, value] of Object.entries(settings)) {
          this.queryClient.setQueryData(settingsQueryKeys.setting(key), value);
        }

        return this.cache.afterSettingsSave(Object.keys(settings));
      },
    });
  }
}

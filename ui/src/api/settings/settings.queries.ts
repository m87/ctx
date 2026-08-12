import { inject, Injectable } from '@angular/core';
import { lastValueFrom } from 'rxjs';
import { SettingsService } from './settings.service';

export const settingsQueryKeys = {
  all: ['settings'] as const,
  settings: () => [...settingsQueryKeys.all, 'settings'] as const,
  settingValues: () => [...settingsQueryKeys.all, 'setting'] as const,
  setting: (key: string) => [...settingsQueryKeys.settingValues(), key] as const,
};

@Injectable({ providedIn: 'root' })
export class SettingsQueries {
  private readonly settingsService = inject(SettingsService);

  settings() {
    return {
      queryKey: settingsQueryKeys.settings(),
      queryFn: () => lastValueFrom(this.settingsService.getSettings()),
    };
  }

  getSetting(key: string) {
    return {
      queryKey: settingsQueryKeys.setting(key),
      queryFn: () => lastValueFrom(this.settingsService.getSetting(key)),
    };
  }

}

import { inject, Injectable } from '@angular/core';
import { lastValueFrom } from 'rxjs';
import { SettingsService } from './settings.service';

@Injectable({ providedIn: 'root' })
export class SettingsQueries {
  static readonly key = ['settings'];
  private settingsService = inject(SettingsService);

  settings() {
    return {
      queryKey: [...SettingsQueries.key, 'settings'],
      queryFn: () => lastValueFrom(this.settingsService.getSettings()),
    };
  }

  getSetting(key: string) {
    return {
      queryKey: [...SettingsQueries.key, 'setting', key],
      queryFn: () => lastValueFrom(this.settingsService.getSetting(key)),
    };
  }

  integrity() {
    return {
      queryKey: [...SettingsQueries.key, 'integrity'],
      queryFn: () => lastValueFrom(this.settingsService.checkIntegrity()),
      enabled: false,
    };
  }
}

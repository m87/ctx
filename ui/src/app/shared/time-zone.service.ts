import { computed, inject, Injectable, signal } from '@angular/core';
import { DateTime, IANAZone } from 'luxon';
import { take } from 'rxjs';
import { SettingsService } from '../../api/settings/settings.service';

export const timeZoneSettingKey = 'client.general.timeZone';
export const browserTimeZonePreference = 'browser';
const dateTimeInputFormat = "yyyy-MM-dd'T'HH:mm";

export function browserTimeZone(): string {
  return Intl.DateTimeFormat().resolvedOptions().timeZone || 'UTC';
}

export function resolveTimeZone(preference: string | null | undefined): string {
  if (!preference || preference === browserTimeZonePreference) {
    return browserTimeZone();
  }
  return IANAZone.isValidZone(preference) ? preference : browserTimeZone();
}

export function instantInTimeZone(instant: string | null | undefined, zone: string): DateTime {
  if (!instant) {
    return DateTime.invalid('Missing timestamp');
  }
  return DateTime.fromISO(instant, { setZone: true }).setZone(zone);
}

export function inputValueToUTC(input: string, zone: string): string | null {
  const parsed = DateTime.fromFormat(input, dateTimeInputFormat, { zone, setZone: true });
  if (!parsed.isValid || parsed.toFormat(dateTimeInputFormat) !== input) {
    return null;
  }
  return parsed.toUTC().toISO();
}

export function supportedTimeZones(): string[] {
  const intl = Intl as typeof Intl & {
    supportedValuesOf?: (key: 'timeZone') => string[];
  };
  const zones = intl.supportedValuesOf?.('timeZone') ?? [];
  return ['UTC', ...zones.filter((zone) => zone !== 'UTC')];
}

@Injectable({ providedIn: 'root' })
export class TimeZoneService {
  private readonly settingsService = inject(SettingsService);
  private readonly preferenceState = signal(browserTimeZonePreference);

  readonly preference = this.preferenceState.asReadonly();
  readonly effectiveTimeZone = computed(() => resolveTimeZone(this.preference()));
  readonly options = supportedTimeZones();

  constructor() {
    this.settingsService
      .getSetting(timeZoneSettingKey)
      .pipe(take(1))
      .subscribe({
        next: (preference) => this.setPreference(preference),
        error: () => this.setPreference(browserTimeZonePreference),
      });
  }

  setPreference(preference: string | null | undefined): void {
    const normalized = preference?.trim() || browserTimeZonePreference;
    this.preferenceState.set(
      normalized === browserTimeZonePreference || IANAZone.isValidZone(normalized)
        ? normalized
        : browserTimeZonePreference,
    );
  }

  now(): DateTime {
    return DateTime.now().setZone(this.effectiveTimeZone());
  }

  today(): string {
    return this.now().toFormat('yyyy-MM-dd');
  }

  parseInstant(instant: string | null | undefined): DateTime {
    return instantInTimeZone(instant, this.effectiveTimeZone());
  }

  toInputValue(instant: string | null | undefined): string {
    const parsed = this.parseInstant(instant);
    return parsed.isValid ? parsed.toFormat(dateTimeInputFormat) : '';
  }

  inputToUTC(input: string): string | null {
    return inputValueToUTC(input, this.effectiveTimeZone());
  }

  formatDateTime(instant: string | null | undefined): string {
    const parsed = this.parseInstant(instant);
    return parsed.isValid ? parsed.toFormat('yyyy-MM-dd HH:mm') : '...';
  }

  formatTime(instant: string | null | undefined): string {
    const parsed = this.parseInstant(instant);
    return parsed.isValid ? parsed.toFormat('HH:mm') : '...';
  }

  formatDate(instant: string | null | undefined): string {
    const parsed = this.parseInstant(instant);
    return parsed.isValid ? parsed.toFormat('yyyy-MM-dd') : '...';
  }
}

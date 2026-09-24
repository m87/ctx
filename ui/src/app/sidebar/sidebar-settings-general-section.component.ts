import { Component, computed, effect, inject, signal } from '@angular/core';
import { injectMutation, injectQuery } from '@tanstack/angular-query-experimental';
import { SettingsMutations } from '../../api/settings/settings.mutations';
import { SettingsQueries } from '../../api/settings/settings.queries';
import { Settings } from '../../api/settings/settings.service';
import {
  browserTimeZonePreference,
  browserTimeZone,
  timeZoneSettingKey,
  TimeZoneService,
} from '../shared/time-zone.service';
import {
  normalizeTimelineRangeMode,
  timelineRangeSettingKey,
  TimelineRangeMode,
} from '../shared/timeline-range';

const themeKey = 'client.general.theme';
const firstDayKey = 'client.general.firstDay';

@Component({
  selector: 'ctx-sidebar-settings-general-section',
  template: `
    <div class="space-y-7">
      <div class="space-y-2">
        <div class="text-foreground font-medium text-lead">Theme mode</div>
        <div class="text-dense sm:text-sm">Choose your preferred app theme.</div>
        <div class="grid grid-cols-2 gap-2 sm:gap-3 pt-1">
          <button
            type="button"
            class="h-12 rounded-md border text-sm font-medium hover:bg-muted/50"
            [class.bg-muted]="colorMode() === 'light'"
            [class.text-foreground]="colorMode() === 'light'"
            [disabled]="saveSettingsMutation.isPending()"
            (click)="setColorMode('light')"
          >
            Light
          </button>
          <button
            type="button"
            class="h-12 rounded-md border text-sm font-medium hover:bg-muted/50"
            [class.bg-muted]="colorMode() === 'dark'"
            [class.text-foreground]="colorMode() === 'dark'"
            [disabled]="saveSettingsMutation.isPending()"
            (click)="setColorMode('dark')"
          >
            Dark
          </button>
        </div>
      </div>

      <div class="space-y-2">
        <label for="general-time-zone" class="text-foreground font-medium text-lead">
          Time zone
        </label>
        <div class="text-dense sm:text-sm">
          Display every recorded event in this time zone. Browser currently resolves to
          {{ browserZone }}.
        </div>
        <select
          id="general-time-zone"
          class="w-full h-10 rounded-md border border-border bg-background px-3 text-sm mt-1"
          [disabled]="saveSettingsMutation.isPending()"
          (change)="setTimeZone(getSelectValue($event))"
        >
          <option
            [value]="browserTimeZonePreference"
            [selected]="selectedTimeZone() === browserTimeZonePreference"
          >
            Browser ({{ browserZone }})
          </option>
          @for (zone of timeZoneOptions(); track zone) {
            <option [value]="zone" [selected]="selectedTimeZone() === zone">{{ zone }}</option>
          }
        </select>
      </div>

      <div class="space-y-2">
        <div class="text-foreground font-medium text-lead">First day of week</div>
        <div class="text-dense sm:text-sm">Choose which day starts the week.</div>
        <div class="grid grid-cols-2 gap-2 sm:gap-3 pt-1">
          <button
            type="button"
            class="h-12 rounded-md border text-sm font-medium hover:bg-muted/50"
            [class.bg-muted]="weekStart() === 'monday'"
            [class.text-foreground]="weekStart() === 'monday'"
            [disabled]="saveSettingsMutation.isPending()"
            (click)="setWeekStart('monday')"
          >
            Monday
          </button>
          <button
            type="button"
            class="h-12 rounded-md border text-sm font-medium hover:bg-muted/50"
            [class.bg-muted]="weekStart() === 'sunday'"
            [class.text-foreground]="weekStart() === 'sunday'"
            [disabled]="saveSettingsMutation.isPending()"
            (click)="setWeekStart('sunday')"
          >
            Sunday
          </button>
        </div>
      </div>

      <div class="space-y-2">
        <div class="text-foreground font-medium text-lead">Timeline range</div>
        <div class="text-dense sm:text-sm">
          Show the entire day or fit the timeline to the day's recorded intervals.
        </div>
        <div class="ui-tablist mt-1" role="radiogroup" aria-label="Timeline range">
          <button
            type="button"
            class="ui-tab disabled:pointer-events-none disabled:opacity-50"
            role="radio"
            [attr.aria-checked]="timelineRangeMode() === 'full-day'"
            [disabled]="saveSettingsMutation.isPending()"
            (click)="setTimelineRangeMode('full-day')"
          >
            Full day
          </button>
          <button
            type="button"
            class="ui-tab disabled:pointer-events-none disabled:opacity-50"
            role="radio"
            [attr.aria-checked]="timelineRangeMode() === 'intervals'"
            [disabled]="saveSettingsMutation.isPending()"
            (click)="setTimelineRangeMode('intervals')"
          >
            Fit to intervals
          </button>
        </div>
      </div>
    </div>
  `,
})
export class SidebarSettingsGeneralSectionComponent {
  private settingsQueries = inject(SettingsQueries);
  private settingsMutations = inject(SettingsMutations);
  private timeZone = inject(TimeZoneService);

  readonly colorMode = signal<'light' | 'dark'>('light');
  readonly weekStart = signal<'monday' | 'sunday'>('monday');
  readonly timelineRangeMode = signal<TimelineRangeMode>('full-day');
  readonly selectedTimeZone = this.timeZone.preference;
  readonly browserTimeZonePreference = browserTimeZonePreference;
  readonly browserZone = browserTimeZone();
  readonly timeZoneOptions = computed(() => {
    const selectedTimeZone = this.selectedTimeZone();
    const options = this.timeZone.options;
    return selectedTimeZone === browserTimeZonePreference || options.includes(selectedTimeZone)
      ? options
      : [selectedTimeZone, ...options];
  });

  settingsQuery = injectQuery(() => this.settingsQueries.settings());
  saveSettingsMutation = injectMutation(() => this.settingsMutations.save());

  private readonly settings = computed<Settings>(() => this.settingsQuery.data() ?? {});

  private readonly syncSettingsEffect = effect(() => {
    const settings = this.settings();
    const theme = settings[themeKey];
    const firstDay = settings[firstDayKey];
    const selectedTimeZone = settings[timeZoneSettingKey];
    const timelineRangeMode = settings[timelineRangeSettingKey];

    if (theme === 'light' || theme === 'dark') {
      this.colorMode.set(theme);
    }

    if (firstDay === 'Monday') {
      this.weekStart.set('monday');
    }

    if (firstDay === 'Sunday') {
      this.weekStart.set('sunday');
    }

    if (selectedTimeZone) {
      this.timeZone.setPreference(selectedTimeZone);
    }

    this.timelineRangeMode.set(normalizeTimelineRangeMode(timelineRangeMode));
  });

  setColorMode(mode: 'light' | 'dark'): void {
    this.colorMode.set(mode);
    this.saveSettings();
  }

  setWeekStart(day: 'monday' | 'sunday'): void {
    this.weekStart.set(day);
    this.saveSettings();
  }

  setTimeZone(timeZone: string): void {
    this.timeZone.setPreference(timeZone);
    this.saveSettings();
  }

  setTimelineRangeMode(mode: TimelineRangeMode): void {
    this.timelineRangeMode.set(mode);
    this.saveSettings();
  }

  getSelectValue(event: Event): string {
    return (event.target as HTMLSelectElement).value;
  }

  private saveSettings(): void {
    this.saveSettingsMutation.mutate({
      ...this.settings(),
      [themeKey]: this.colorMode(),
      [firstDayKey]: this.weekStart() === 'monday' ? 'Monday' : 'Sunday',
      [timeZoneSettingKey]: this.selectedTimeZone(),
      [timelineRangeSettingKey]: this.timelineRangeMode(),
    });
  }
}

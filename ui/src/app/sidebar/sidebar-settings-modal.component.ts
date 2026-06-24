import {
  Component,
  EventEmitter,
  Input,
  Output,
  computed,
  effect,
  inject,
  signal,
} from '@angular/core';
import { injectMutation, injectQuery } from '@tanstack/angular-query-experimental';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideX } from '@ng-icons/lucide';
import { SettingsMutations } from '../../api/settings.mutations';
import { SettingsQueries } from '../../api/settings.queries';
import { Settings } from '../../api/settings.service';

const themeKey = 'client.general.theme';
const firstDayKey = 'client.general.firstDay';

@Component({
  selector: 'app-sidebar-settings-modal',
  imports: [NgIcon],
  providers: [provideIcons({ lucideX })],
  template: `
    @if (open) {
      <div
        class="fixed inset-0 z-50 w-screen flex items-end sm:items-center sm:justify-center"
        (click)="requestClose()"
      >
        <div class="absolute inset-0 bg-background/60 backdrop-blur-[1px]"></div>
        <div
          class="relative w-full sm:w-[min(94vw,1040px)] h-[92vh] sm:h-auto max-h-[92vh] sm:max-h-[720px] bg-popover text-popover-foreground border rounded-t-2xl sm:rounded-xl shadow-xl flex flex-col sm:flex-row min-h-0"
          (click)="$event.stopPropagation()"
        >
          <div class="sm:w-60 border-b sm:border-b-0 sm:border-r p-3 sm:p-4 shrink-0">
            <div
              class="text-[11px] uppercase tracking-[0.08em] text-muted-foreground px-2 py-1 hidden sm:block"
            >
              Settings
            </div>
            <div class="flex items-center gap-2 sm:block">
              <div class="flex-1 flex sm:flex-col gap-1 overflow-x-auto sm:overflow-visible">
                @for (section of settingsSections; track section) {
                  <button
                    type="button"
                    class="text-left px-3 sm:px-2 py-2 sm:py-1.5 text-[14px] sm:text-[13px] rounded-md hover:bg-muted/60 whitespace-nowrap min-w-fit"
                    [class.bg-muted]="activeSettingsSection() === section"
                    [class.text-foreground]="activeSettingsSection() === section"
                    [class.text-muted-foreground]="activeSettingsSection() !== section"
                    (click)="activeSettingsSection.set(section)"
                  >
                    {{ section }}
                  </button>
                }
              </div>
              <button
                type="button"
                class="h-9 w-9 rounded-md text-muted-foreground hover:text-foreground hover:bg-muted/60 flex items-center justify-center sm:hidden"
                aria-label="Close settings"
                (click)="requestClose()"
              >
                <ng-icon name="lucideX" class="text-[14px]"></ng-icon>
              </button>
            </div>
          </div>
          <div class="flex-1 min-h-0 flex flex-col">
            <div
              class="h-14 sm:h-14 border-b px-5 sm:px-7 flex items-center justify-between hidden sm:flex"
            >
              <div class="font-semibold text-[15px]">{{ activeSettingsSection() }}</div>
              <button
                type="button"
                class="h-9 w-9 rounded-md text-muted-foreground hover:text-foreground hover:bg-muted/60 flex items-center justify-center"
                aria-label="Close settings"
                (click)="requestClose()"
              >
                <ng-icon name="lucideX" class="text-[14px]"></ng-icon>
              </button>
            </div>
            <div class="p-5 sm:p-7 overflow-auto text-[14px] text-muted-foreground">
              <div class="max-w-[760px] pb-5">
                @if (activeSettingsSection() === 'General') {
                  <div class="space-y-7">
                    <div class="space-y-2">
                      <div class="text-foreground font-medium text-[15px]">Theme mode</div>
                      <div class="text-[13px] sm:text-[14px]">Choose your preferred app theme.</div>
                      <div class="grid grid-cols-2 gap-2 sm:gap-3 pt-1">
                        <button
                          type="button"
                          class="h-12 rounded-md border text-[14px] font-medium hover:bg-muted/50"
                          [class.bg-muted]="colorMode() === 'light'"
                          [class.text-foreground]="colorMode() === 'light'"
                          [disabled]="saveSettingsMutation.isPending()"
                          (click)="setColorMode('light')"
                        >
                          Light
                        </button>
                        <button
                          type="button"
                          class="h-12 rounded-md border text-[14px] font-medium hover:bg-muted/50"
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
                      <div class="text-foreground font-medium text-[15px]">First day of week</div>
                      <div class="text-[13px] sm:text-[14px]">
                        Choose which day starts the week.
                      </div>
                      <div class="grid grid-cols-2 gap-2 sm:gap-3 pt-1">
                        <button
                          type="button"
                          class="h-12 rounded-md border text-[14px] font-medium hover:bg-muted/50"
                          [class.bg-muted]="weekStart() === 'monday'"
                          [class.text-foreground]="weekStart() === 'monday'"
                          [disabled]="saveSettingsMutation.isPending()"
                          (click)="setWeekStart('monday')"
                        >
                          Monday
                        </button>
                        <button
                          type="button"
                          class="h-12 rounded-md border text-[14px] font-medium hover:bg-muted/50"
                          [class.bg-muted]="weekStart() === 'sunday'"
                          [class.text-foreground]="weekStart() === 'sunday'"
                          [disabled]="saveSettingsMutation.isPending()"
                          (click)="setWeekStart('sunday')"
                        >
                          Sunday
                        </button>
                      </div>
                    </div>
                  </div>
                } @else {
                  <div class="space-y-5">
                    <div class="space-y-1.5">
                      <div class="text-foreground font-medium text-[15px]">Data integrity</div>
                      <div class="text-[13px] sm:text-[14px]">
                        Check workspace assignments and references after migration.
                      </div>
                    </div>

                    <button
                      type="button"
                      class="h-10 px-4 mr-4 rounded-md border text-foreground text-[14px] font-medium hover:bg-muted/50 disabled:opacity-50"
                      [disabled]="checkIntegrityMutation.isPending()"
                      (click)="checkIntegrity()"
                    >
                      {{
                        checkIntegrityMutation.isPending() ? 'Checking...' : 'Run integrity check'
                      }}
                    </button>

                    @if (integrityReport(); as report) {
                      @if (!report.healthy) {
                        <button
                          type="button"
                          class="h-10 px-4 rounded-md bg-primary text-primary-foreground text-[14px] font-medium hover:bg-primary/90 disabled:opacity-50"
                          [disabled]="repairIntegrityMutation.isPending()"
                          (click)="repairIntegrity()"
                        >
                          {{
                            repairIntegrityMutation.isPending()
                              ? 'Repairing...'
                              : 'Repair automatically'
                          }}
                        </button>
                      }
                    }

                    @if (repairIntegrityMutation.data(); as repairResult) {
                      <div class="rounded-md border bg-muted/30 p-3 text-[13px]">
                        Repaired {{ repairResult.repairedCount }} records. Issues that cannot be
                        repaired safely remain listed below.
                      </div>
                    }

                    @if (integrityReport(); as report) {
                      <div
                        class="rounded-lg border p-4"
                        [class.border-emerald-500]="report.healthy"
                        [class.border-destructive]="!report.healthy"
                      >
                        <div
                          class="font-medium"
                          [class.text-emerald-600]="report.healthy"
                          [class.text-destructive]="!report.healthy"
                        >
                          {{ report.healthy ? 'Integrity check passed' : 'Integrity issues found' }}
                        </div>
                        <div class="mt-2 text-[12px] text-muted-foreground">
                          {{ report.workspaceCount }} workspaces,
                          {{ report.contextCount }} contexts, {{ report.intervalCount }} intervals
                        </div>
                      </div>

                      @if (report.issues.length > 0) {
                        <div class="space-y-2">
                          @for (issue of report.issues; track issue.code + issue.entityId) {
                            <div class="rounded-md border p-3 text-[13px]">
                              <div class="flex flex-wrap items-center gap-2">
                                <span class="font-medium text-foreground">{{ issue.code }}</span>
                                <span class="text-[11px] uppercase text-muted-foreground">
                                  {{ issue.entityType }}
                                </span>
                              </div>
                              <div class="mt-1">{{ issue.description }}</div>
                              <div
                                class="mt-1 font-mono text-[11px] break-all text-muted-foreground"
                              >
                                {{ issue.entityId || '(missing id)' }}
                              </div>
                            </div>
                          }
                        </div>
                      }
                    }
                  </div>
                }
              </div>
            </div>
          </div>
        </div>
      </div>
    }
  `,
})
export class SidebarSettingsModalComponent {
  private settingsQueries = inject(SettingsQueries);
  private settingsMutations = inject(SettingsMutations);

  @Input() open = false;
  @Output() openChange = new EventEmitter<boolean>();

  readonly settingsSections = ['General', 'Data integrity'] as const;
  readonly activeSettingsSection = signal<(typeof this.settingsSections)[number]>('General');
  readonly colorMode = signal<'light' | 'dark'>('light');
  readonly weekStart = signal<'monday' | 'sunday'>('monday');

  settingsQuery = injectQuery(() => this.settingsQueries.settings());
  integrityQuery = injectQuery(() => this.settingsQueries.integrity());
  saveSettingsMutation = injectMutation(() => this.settingsMutations.save());
  checkIntegrityMutation = injectMutation(() => this.settingsMutations.checkIntegrity());
  repairIntegrityMutation = injectMutation(() => this.settingsMutations.repairIntegrity());

  readonly integrityReport = computed(
    () =>
      this.repairIntegrityMutation.data()?.report ??
      this.checkIntegrityMutation.data() ??
      this.integrityQuery.data(),
  );

  private readonly settings = computed<Settings>(() => this.settingsQuery.data() ?? {});

  private readonly syncSettingsEffect = effect(() => {
    const settings = this.settings();
    const theme = settings[themeKey];
    const firstDay = settings[firstDayKey];

    if (theme === 'light' || theme === 'dark') {
      this.colorMode.set(theme);
    }

    if (firstDay === 'Monday') {
      this.weekStart.set('monday');
    }

    if (firstDay === 'Sunday') {
      this.weekStart.set('sunday');
    }
  });

  requestClose(): void {
    this.openChange.emit(false);
  }

  checkIntegrity(): void {
    this.checkIntegrityMutation.mutate();
  }

  repairIntegrity(): void {
    this.repairIntegrityMutation.mutate();
  }

  setColorMode(mode: 'light' | 'dark'): void {
    this.colorMode.set(mode);
    this.saveSettings();
  }

  setWeekStart(day: 'monday' | 'sunday'): void {
    this.weekStart.set(day);
    this.saveSettings();
  }

  private saveSettings(): void {
    this.saveSettingsMutation.mutate({
      ...this.settings(),
      [themeKey]: this.colorMode(),
      [firstDayKey]: this.weekStart() === 'monday' ? 'Monday' : 'Sunday',
    });
  }
}

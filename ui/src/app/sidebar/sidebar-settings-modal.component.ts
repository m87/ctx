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
import { lucideTrash2, lucideX } from '@ng-icons/lucide';
import { DateTime } from 'luxon';
import { ContextMutations } from '../../api/context.mutations';
import { IntervalMutations } from '../../api/interval.mutations';
import { Interval, ZonedDateTime } from '../../api/interval.service';
import { SettingsMutations } from '../../api/settings.mutations';
import { SettingsQueries } from '../../api/settings.queries';
import {
  IntegrityDateTime,
  IntegrityIssue,
  IntegrityReport,
  Settings,
} from '../../api/settings.service';

const themeKey = 'client.general.theme';
const firstDayKey = 'client.general.firstDay';

type IntegrityIssueGroup = {
  key: string;
  issue: IntegrityIssue;
  hiddenIssueCount: number;
};

type IntegrityIntervalTimeField = 'start' | 'end';

type IntegrityIntervalTimeInputs = {
  start?: string;
  end?: string;
};

@Component({
  selector: 'ctx-sidebar-settings-modal',
  imports: [NgIcon],
  providers: [provideIcons({ lucideTrash2, lucideX })],
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
                      @if (hasRepairableIssues(report)) {
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

                      @if (visibleIntegrityIssueGroups().length > 0) {
                        <div class="space-y-2">
                          @for (group of visibleIntegrityIssueGroups(); track group.key) {
                            @let issue = group.issue;
                            <div class="rounded-md border p-3 text-[13px]">
                              <div class="flex items-start justify-between gap-3">
                                <div class="flex flex-wrap items-center gap-2">
                                  <span class="font-medium text-foreground">{{ issue.code }}</span>
                                  <span class="text-[11px] uppercase text-muted-foreground">
                                    {{ issue.entityType }}
                                  </span>
                                  <span
                                    class="rounded-full px-2 py-0.5 text-[11px] font-medium"
                                    [class.bg-emerald-500/10]="issue.repairable"
                                    [class.text-emerald-600]="issue.repairable"
                                    [class.bg-amber-500/10]="!issue.repairable"
                                    [class.text-amber-700]="!issue.repairable"
                                  >
                                    {{
                                      issue.repairable
                                        ? 'Auto-repairable'
                                        : 'Manual action required'
                                    }}
                                  </span>
                                </div>

                                @if (!issue.repairable && issue.entityId) {
                                  <button
                                    type="button"
                                    class="h-8 w-8 shrink-0 rounded-md text-destructive hover:bg-destructive/10 disabled:opacity-50 flex items-center justify-center"
                                    [disabled]="isDeletingIntegrityEntity()"
                                    [attr.aria-label]="'Delete problematic ' + issue.entityType"
                                    [title]="'Delete ' + issue.entityType"
                                    (click)="deleteIntegrityIssue(issue)"
                                  >
                                    <ng-icon name="lucideTrash2" class="text-[15px]"></ng-icon>
                                  </button>
                                }
                              </div>
                              <div class="mt-1">{{ issue.description }}</div>

                              @if (group.hiddenIssueCount > 0) {
                                <div
                                  class="mt-2 rounded-md bg-muted/40 px-2 py-1.5 text-[12px] text-muted-foreground"
                                >
                                  {{ group.hiddenIssueCount }} more integrity
                                  {{ group.hiddenIssueCount === 1 ? 'issue was' : 'issues were' }}
                                  detected for this {{ issue.entityType }}. Resolve this first step,
                                  then run the check again.
                                </div>
                              }

                              @if (issue.entityType === 'context' && issue.details?.name) {
                                <div class="mt-2 font-medium text-foreground">
                                  {{ issue.details?.name }}
                                </div>
                              }

                              @if (issue.entityType === 'interval') {
                                <div
                                  class="mt-2 grid grid-cols-[auto_1fr] gap-x-3 gap-y-1 text-[12px]"
                                >
                                  <span class="text-muted-foreground">Interval</span>
                                  <span class="text-foreground">
                                    {{ formatIntegrityTime(issue.details?.start) }} –
                                    {{ formatIntegrityTime(issue.details?.end) }}
                                  </span>
                                  <span class="text-muted-foreground">Context</span>
                                  <span class="font-mono break-all">
                                    {{ issue.details?.contextId || '(missing)' }}
                                  </span>
                                </div>
                              }

                              @if (isContextAssignmentIssue(issue)) {
                                <div class="mt-3 flex flex-col sm:flex-row gap-2">
                                  <select
                                    class="h-9 min-w-0 flex-1 rounded-md border border-border bg-background px-3 text-[12px]"
                                    [value]="selectedIntegrityContext(issue.entityId)"
                                    [disabled]="
                                      availableIntegrityContexts().length === 0 ||
                                      moveIntegrityIntervalMutation.isPending()
                                    "
                                    (change)="
                                      selectIntegrityContext(issue.entityId, getSelectValue($event))
                                    "
                                  >
                                    <option value="">
                                      {{
                                        integrityContextsQuery.isPending()
                                          ? 'Loading contexts...'
                                          : availableIntegrityContexts().length === 0
                                            ? 'No contexts available'
                                            : 'Select context...'
                                      }}
                                    </option>
                                    @for (
                                      context of availableIntegrityContexts();
                                      track context.id
                                    ) {
                                      <option [value]="context.id">
                                        {{ context.name || '(unnamed context)' }} ·
                                        {{ context.workspaceName || context.workspaceId }}
                                      </option>
                                    }
                                  </select>
                                  <button
                                    type="button"
                                    class="h-9 px-3 rounded-md bg-primary text-primary-foreground text-[12px] font-medium hover:bg-primary/90 disabled:opacity-50"
                                    [disabled]="
                                      !selectedIntegrityContext(issue.entityId) ||
                                      moveIntegrityIntervalMutation.isPending()
                                    "
                                    (click)="assignIntegrityContext(issue)"
                                  >
                                    {{
                                      moveIntegrityIntervalMutation.isPending()
                                        ? 'Assigning...'
                                        : 'Assign'
                                    }}
                                  </button>
                                </div>
                              }

                              @if (isIntervalTimeEditIssue(issue)) {
                                <div class="mt-3 rounded-md border bg-muted/20 p-3">
                                  <div
                                    class="mb-2 text-[11px] uppercase tracking-[0.08em] text-muted-foreground font-semibold"
                                  >
                                    Set interval time
                                  </div>
                                  <div class="flex flex-col sm:flex-row gap-2 sm:items-end">
                                    <label class="flex-1 text-[12px] text-muted-foreground">
                                      Start
                                      <input
                                        type="datetime-local"
                                        class="mt-1 h-9 w-full rounded-md border border-border bg-background px-3 text-[12px] text-foreground"
                                        [value]="integrityIntervalTimeInput(issue, 'start')"
                                        (input)="
                                          setIntegrityIntervalTimeInput(
                                            issue.entityId,
                                            'start',
                                            getInputValue($event)
                                          )
                                        "
                                      />
                                    </label>
                                    <label class="flex-1 text-[12px] text-muted-foreground">
                                      End
                                      <input
                                        type="datetime-local"
                                        class="mt-1 h-9 w-full rounded-md border border-border bg-background px-3 text-[12px] text-foreground"
                                        [value]="integrityIntervalTimeInput(issue, 'end')"
                                        (input)="
                                          setIntegrityIntervalTimeInput(
                                            issue.entityId,
                                            'end',
                                            getInputValue($event)
                                          )
                                        "
                                      />
                                    </label>
                                    <button
                                      type="button"
                                      class="h-9 px-3 rounded-md bg-primary text-primary-foreground text-[12px] font-medium hover:bg-primary/90 disabled:opacity-50"
                                      [disabled]="
                                        updateIntegrityIntervalMutation.isPending() ||
                                        !issue.details?.contextId
                                      "
                                      (click)="saveIntegrityIntervalTime(issue)"
                                    >
                                      {{
                                        updateIntegrityIntervalMutation.isPending()
                                          ? 'Saving...'
                                          : 'Save time'
                                      }}
                                    </button>
                                  </div>
                                  @if (integrityIntervalTimeError(issue.entityId)) {
                                    <div class="mt-2 text-[12px] text-destructive">
                                      {{ integrityIntervalTimeError(issue.entityId) }}
                                    </div>
                                  }
                                  @if (!issue.details?.contextId) {
                                    <div class="mt-2 text-[12px] text-muted-foreground">
                                      Assign a context first, then set the interval time.
                                    </div>
                                  }
                                </div>
                              }

                              @if (issue.details?.workspaceId && !isContextAssignmentIssue(issue)) {
                                <div class="mt-1 text-[12px]">
                                  <span class="text-muted-foreground">Workspace:</span>
                                  <span class="ml-2 font-mono break-all">
                                    {{ issue.details?.workspaceId }}
                                  </span>
                                </div>
                              }

                              <div
                                class="mt-2 font-mono text-[11px] break-all text-muted-foreground"
                              >
                                ID: {{ issue.entityId || '(missing id)' }}
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
  private contextMutations = inject(ContextMutations);
  private intervalMutations = inject(IntervalMutations);

  @Input() open = false;
  @Output() openChange = new EventEmitter<boolean>();

  readonly settingsSections = ['General', 'Data integrity'] as const;
  readonly activeSettingsSection = signal<(typeof this.settingsSections)[number]>('General');
  readonly colorMode = signal<'light' | 'dark'>('light');
  readonly weekStart = signal<'monday' | 'sunday'>('monday');

  settingsQuery = injectQuery(() => this.settingsQueries.settings());
  integrityQuery = injectQuery(() => this.settingsQueries.integrity());

  private readonly latestIntegrityReport = signal<IntegrityReport | undefined>(undefined);
  readonly integrityReport = computed(
    () => this.latestIntegrityReport() ?? this.integrityQuery.data(),
  );
  integrityContextsQuery = injectQuery(() =>
    this.settingsQueries.integrityContexts(
      this.integrityReport()?.issues.some((issue) => this.isContextAssignmentIssue(issue)) ?? false,
    ),
  );
  readonly availableIntegrityContexts = computed(() => this.integrityContextsQuery.data() ?? []);
  readonly visibleIntegrityIssueGroups = computed(() =>
    this.groupIntegrityIssues(this.integrityReport()?.issues ?? []),
  );

  saveSettingsMutation = injectMutation(() => this.settingsMutations.save());
  checkIntegrityMutation = injectMutation(() => this.settingsMutations.checkIntegrity());
  repairIntegrityMutation = injectMutation(() => this.settingsMutations.repairIntegrity());
  deleteIntegrityContextMutation = injectMutation(() => this.contextMutations.delete());
  deleteIntegrityIntervalMutation = injectMutation(() => this.intervalMutations.delete());
  moveIntegrityIntervalMutation = injectMutation(() => this.intervalMutations.move());
  updateIntegrityIntervalMutation = injectMutation(() => this.intervalMutations.update());

  private readonly integrityContextSelections = signal<Record<string, string>>({});
  private readonly integrityIntervalTimeInputs = signal<
    Record<string, IntegrityIntervalTimeInputs>
  >({});
  private readonly integrityIntervalTimeErrors = signal<Record<string, string>>({});

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
    this.checkIntegrityMutation.mutate(undefined, {
      onSuccess: (report) => this.updateIntegrityReport(report),
    });
  }

  repairIntegrity(): void {
    this.repairIntegrityMutation.mutate(undefined, {
      onSuccess: (result) => this.updateIntegrityReport(result.report),
    });
  }

  hasRepairableIssues(report: IntegrityReport): boolean {
    return report.issues.some((issue) => issue.repairable);
  }

  isDeletingIntegrityEntity(): boolean {
    return (
      this.deleteIntegrityContextMutation.isPending() ||
      this.deleteIntegrityIntervalMutation.isPending()
    );
  }

  isContextAssignmentIssue(issue: IntegrityIssue): boolean {
    return (
      issue.entityType === 'interval' &&
      (issue.code === 'INTERVAL_MISSING_CONTEXT' || issue.code === 'INTERVAL_CONTEXT_NOT_FOUND')
    );
  }

  isIntervalTimeEditIssue(issue: IntegrityIssue): boolean {
    return issue.entityType === 'interval' && issue.code === 'INACTIVE_INTERVAL_MISSING_TIME';
  }

  private groupIntegrityIssues(issues: IntegrityIssue[]): IntegrityIssueGroup[] {
    const groupedIssues = new Map<string, IntegrityIssue[]>();
    const keys: string[] = [];

    issues.forEach((issue, index) => {
      const key = issue.entityId
        ? `${issue.entityType}:${issue.entityId}`
        : `${issue.entityType}:missing-id:${index}`;

      if (!groupedIssues.has(key)) {
        groupedIssues.set(key, []);
        keys.push(key);
      }

      groupedIssues.get(key)?.push(issue);
    });

    return keys.map((key) => {
      const groupIssues = groupedIssues.get(key) ?? [];
      const issue =
        groupIssues.find((item) => this.isContextAssignmentIssue(item)) ?? groupIssues[0];

      return {
        key,
        issue,
        hiddenIssueCount: Math.max(0, groupIssues.length - 1),
      };
    });
  }

  selectedIntegrityContext(intervalId: string): string {
    return this.integrityContextSelections()[intervalId] ?? '';
  }

  selectIntegrityContext(intervalId: string, contextId: string): void {
    this.integrityContextSelections.update((selections) => ({
      ...selections,
      [intervalId]: contextId,
    }));
  }

  assignIntegrityContext(issue: IntegrityIssue): void {
    const targetContextId = this.selectedIntegrityContext(issue.entityId);
    if (!this.isContextAssignmentIssue(issue) || !issue.entityId || !targetContextId) {
      return;
    }

    this.moveIntegrityIntervalMutation.mutate(
      { id: issue.entityId, targetContextId },
      { onSuccess: () => this.refreshIntegrityReport() },
    );
  }

  integrityIntervalTimeInput(issue: IntegrityIssue, field: IntegrityIntervalTimeField): string {
    const savedValue = this.integrityIntervalTimeInputs()[issue.entityId]?.[field];
    if (savedValue !== undefined) {
      return savedValue;
    }

    return this.integrityDateTimeToInputValue(issue.details?.[field]);
  }

  setIntegrityIntervalTimeInput(
    intervalId: string,
    field: IntegrityIntervalTimeField,
    value: string,
  ): void {
    this.integrityIntervalTimeInputs.update((inputs) => ({
      ...inputs,
      [intervalId]: {
        ...(inputs[intervalId] ?? {}),
        [field]: value,
      },
    }));
    this.setIntegrityIntervalTimeError(intervalId, '');
  }

  integrityIntervalTimeError(intervalId: string): string {
    return this.integrityIntervalTimeErrors()[intervalId] ?? '';
  }

  saveIntegrityIntervalTime(issue: IntegrityIssue): void {
    if (!this.isIntervalTimeEditIssue(issue) || !issue.entityId) {
      return;
    }

    const contextId = issue.details?.contextId ?? '';
    if (!contextId) {
      this.setIntegrityIntervalTimeError(issue.entityId, 'Assign a context first.');
      return;
    }

    const parsed = this.parseIntegrityIntervalTimeInputs(
      issue.entityId,
      this.integrityIntervalTimeInput(issue, 'start'),
      this.integrityIntervalTimeInput(issue, 'end'),
    );
    if (!parsed) {
      return;
    }

    const interval: Interval = {
      id: issue.entityId,
      contextId,
      start: parsed.start,
      end: parsed.end,
      duration: 0,
      workspaceId: issue.details?.workspaceId ?? '',
    };

    this.updateIntegrityIntervalMutation.mutate(
      { id: issue.entityId, interval },
      {
        onSuccess: () => {
          this.clearIntegrityIntervalTimeState(issue.entityId);
          this.refreshIntegrityReport();
        },
      },
    );
  }

  getSelectValue(event: Event): string {
    return (event.target as HTMLSelectElement).value;
  }

  getInputValue(event: Event): string {
    return (event.target as HTMLInputElement | HTMLTextAreaElement).value;
  }

  deleteIntegrityIssue(issue: IntegrityIssue): void {
    if (issue.repairable || !issue.entityId) {
      return;
    }

    const entityLabel = issue.details?.name
      ? `${issue.entityType} "${issue.details.name}"`
      : `${issue.entityType} "${issue.entityId}"`;
    if (!window.confirm(`Delete ${entityLabel}? This action cannot be undone.`)) {
      return;
    }

    if (issue.entityType === 'context') {
      this.deleteIntegrityContextMutation.mutate(issue.entityId, {
        onSuccess: () => this.refreshIntegrityReport(),
      });
      return;
    }

    this.deleteIntegrityIntervalMutation.mutate(
      {
        id: issue.entityId,
        contextId: issue.details?.contextId ?? '',
      },
      { onSuccess: () => this.refreshIntegrityReport() },
    );
  }

  formatIntegrityTime(value: IntegrityDateTime | undefined): string {
    if (!value?.time || value.isZero) {
      return 'Not set';
    }

    return new Date(value.time).toLocaleString();
  }

  private integrityDateTimeToInputValue(value: IntegrityDateTime | undefined): string {
    if (!value?.time || value.isZero) {
      return '';
    }

    const dateTime = value.timezone
      ? DateTime.fromISO(value.time, { zone: value.timezone })
      : DateTime.fromISO(value.time);
    if (!dateTime.isValid) {
      return '';
    }

    return dateTime.toFormat("yyyy-MM-dd'T'HH:mm");
  }

  private parseIntegrityIntervalTimeInputs(
    intervalId: string,
    startInput: string,
    endInput: string,
  ): { start: ZonedDateTime; end: ZonedDateTime } | null {
    const startDateTime = DateTime.fromFormat(startInput, "yyyy-MM-dd'T'HH:mm");
    const endDateTime = DateTime.fromFormat(endInput, "yyyy-MM-dd'T'HH:mm");

    if (!startDateTime.isValid || !endDateTime.isValid) {
      this.setIntegrityIntervalTimeError(intervalId, 'Set both start and end date/time.');
      return null;
    }

    if (endDateTime <= startDateTime) {
      this.setIntegrityIntervalTimeError(intervalId, 'End must be later than start.');
      return null;
    }

    return {
      start: ZonedDateTime.fromDateTime(startDateTime),
      end: ZonedDateTime.fromDateTime(endDateTime),
    };
  }

  private setIntegrityIntervalTimeError(intervalId: string, message: string): void {
    this.integrityIntervalTimeErrors.update((errors) => ({
      ...errors,
      [intervalId]: message,
    }));
  }

  private clearIntegrityIntervalTimeState(intervalId: string): void {
    this.integrityIntervalTimeInputs.update((inputs) => {
      const next = { ...inputs };
      delete next[intervalId];
      return next;
    });
    this.integrityIntervalTimeErrors.update((errors) => {
      const next = { ...errors };
      delete next[intervalId];
      return next;
    });
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

  private refreshIntegrityReport(): void {
    void this.integrityQuery.refetch().then(({ data: report }) => {
      if (report) {
        this.updateIntegrityReport(report);
      }
    });
  }

  private updateIntegrityReport(report: IntegrityReport): void {
    this.latestIntegrityReport.set(report);
    if (report.issues.some((issue) => this.isContextAssignmentIssue(issue))) {
      void this.integrityContextsQuery.refetch();
    }
  }
}

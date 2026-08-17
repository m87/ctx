import { Component, computed, inject, signal } from '@angular/core';
import { injectMutation, injectQuery } from '@tanstack/angular-query-experimental';
import { HlmSkeletonImports } from '@spartan-ng/helm/skeleton';
import { ContextMutations } from '../../api/context/context.mutations';
import { ContextQueries } from '../../api/context/context.queries';
import { Store } from '@ngxs/store';
import { WorkspaceState } from './workspace.state';
import { QueryErrorStateComponent } from '../shared/query-error-state.component';
import { TimeZoneService } from '../shared/time-zone.service';

const DEFAULT_ARCHIVE_DAYS = '30';
const MAX_ARCHIVE_DAYS = 365000;

@Component({
  selector: 'ctx-sidebar-settings-archivization-section',
  imports: [HlmSkeletonImports, QueryErrorStateComponent],
  template: `
    <div class="space-y-6">
      <div class="space-y-1.5">
        <div class="flex items-center gap-2">
          <div class="text-foreground font-medium text-[15px]">Archive inactive contexts</div>
          <span
            class="rounded-full border border-border bg-muted/50 px-2 py-0.5 text-[10px] font-medium uppercase tracking-[0.06em] text-muted-foreground"
          >
            Workspace
          </span>
        </div>
        <div class="text-[13px] sm:text-[14px]">
          Archive every context whose latest interval is older than the selected threshold. Contexts
          without intervals and already archived contexts are ignored.
        </div>
      </div>

      <div class="space-y-2">
        <label for="archive-older-than-days" class="text-foreground font-medium text-[14px]">
          Inactive for at least
        </label>
        <div class="flex items-center gap-2">
          <input
            id="archive-older-than-days"
            type="text"
            inputmode="numeric"
            autocomplete="off"
            class="h-10 w-32 rounded-md border border-border bg-background px-3 text-sm text-foreground outline-none focus-visible:ring-2 focus-visible:ring-ring/50"
            [value]="daysInput()"
            [attr.aria-invalid]="!validDays()"
            (input)="setDaysInput($event)"
          />
          <span class="text-[13px]">days</span>
        </div>
        @if (!validDays()) {
          <div class="text-[12px] text-destructive">
            Enter a whole number between 1 and {{ maxArchiveDays }}.
          </div>
        }
      </div>

      @if (!activeWorkspaceId()) {
        <div class="rounded-lg border border-dashed p-5 text-center text-[13px]">
          Select a workspace to preview contexts.
        </div>
      } @else if (validDays()) {
        @if (showPreviewError()) {
          <ctx-query-error-state
            class="min-h-40"
            [error]="archiveCandidatesQuery.error()"
            [paused]="archiveCandidatesQuery.isPaused()"
            resourceName="archive preview"
            [retrying]="archiveCandidatesQuery.isFetching()"
            (retry)="retryPreview()"
          ></ctx-query-error-state>
        } @else if (archiveCandidatesQuery.isLoading()) {
          <div class="space-y-3" role="status">
            <span class="sr-only">Loading archive candidates</span>
            <hlm-skeleton class="h-4 w-64 max-w-full" />
            @for (item of candidateSkeletonItems; track item) {
              <div class="rounded-lg border p-3">
                <hlm-skeleton class="h-3.5 w-2/5 mb-2" />
                <hlm-skeleton class="h-3 w-48 max-w-full" />
              </div>
            }
          </div>
        } @else if (archiveCandidatesQuery.data(); as preview) {
          <div class="space-y-3">
            <div class="text-[13px]">
              Cutoff:
              <span class="font-medium text-foreground">{{ formatCutoff(preview.cutoff) }}</span>
            </div>

            <div class="flex items-center justify-between gap-3">
              <div class="text-foreground font-medium text-[14px]">Contexts to archive</div>
              <span class="rounded-full border bg-muted/40 px-2 py-0.5 text-[11px] font-medium">
                {{ preview.contexts.length }}
              </span>
            </div>

            @if (preview.contexts.length > 0) {
              <div class="space-y-2">
                @for (context of preview.contexts; track context.id) {
                  <div class="rounded-lg border bg-card p-3">
                    <div class="flex items-start justify-between gap-3">
                      <div class="min-w-0">
                        <div class="font-medium text-foreground truncate">{{ context.name }}</div>
                        @if (context.project; as project) {
                          <div class="mt-0.5 text-[11px] truncate">Project: {{ project.name }}</div>
                        }
                      </div>
                      <div class="shrink-0 text-right text-[11px]">
                        Latest interval<br />
                        <span class="text-foreground">{{
                          formatLastInterval(context.lastIntervalAt)
                        }}</span>
                      </div>
                    </div>
                  </div>
                }
              </div>

              <button
                type="button"
                class="h-10 rounded-md bg-primary px-4 text-[14px] font-medium text-primary-foreground hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-50"
                [disabled]="archiveMutation.isPending() || archiveCandidatesQuery.isFetching()"
                (click)="archiveContexts()"
              >
                {{
                  archiveMutation.isPending()
                    ? 'Archiving...'
                    : 'Archive ' +
                      preview.contexts.length +
                      (preview.contexts.length === 1 ? ' context' : ' contexts')
                }}
              </button>
            } @else {
              <div class="rounded-lg border border-dashed p-5 text-center text-[13px]">
                No contexts match this threshold.
              </div>
            }
          </div>
        }
      }

      @if (lastArchivedCount() !== null) {
        <div
          class="rounded-md border border-emerald-500/40 bg-emerald-500/5 p-3 text-[13px] text-emerald-700"
          aria-live="polite"
        >
          Archived {{ lastArchivedCount() }}
          {{ lastArchivedCount() === 1 ? 'context' : 'contexts' }}.
        </div>
      }
    </div>
  `,
})
export class SidebarSettingsArchivizationSectionComponent {
  private readonly store = inject(Store);
  private readonly contextQueries = inject(ContextQueries);
  private readonly contextMutations = inject(ContextMutations);
  private readonly timeZone = inject(TimeZoneService);

  readonly maxArchiveDays = MAX_ARCHIVE_DAYS;
  readonly candidateSkeletonItems = [0, 1, 2];
  readonly daysInput = signal(DEFAULT_ARCHIVE_DAYS);
  readonly lastArchivedCount = signal<number | null>(null);
  readonly activeWorkspaceId = this.store.selectSignal(WorkspaceState.selectedWorkspaceId);
  readonly olderThanDays = computed(() => this.parseDays(this.daysInput()));
  readonly validDays = computed(() => this.olderThanDays() !== null);

  readonly archiveCandidatesQuery = injectQuery(() =>
    this.contextQueries.archiveCandidates(
      this.activeWorkspaceId(),
      this.olderThanDays() ?? 0,
      this.timeZone.effectiveTimeZone(),
    ),
  );
  readonly archiveMutation = injectMutation(() => this.contextMutations.bulkArchive());
  readonly showPreviewError = computed(
    () =>
      this.archiveCandidatesQuery.data() === undefined &&
      (this.archiveCandidatesQuery.isError() || this.archiveCandidatesQuery.isPaused()),
  );

  setDaysInput(event: Event): void {
    this.daysInput.set((event.target as HTMLInputElement).value);
    this.lastArchivedCount.set(null);
  }

  archiveContexts(): void {
    const workspaceId = this.activeWorkspaceId();
    const olderThanDays = this.olderThanDays();
    const count = this.archiveCandidatesQuery.data()?.contexts.length ?? 0;
    if (!workspaceId || olderThanDays === null || count === 0) {
      return;
    }

    if (!window.confirm(`Archive ${count} ${count === 1 ? 'context' : 'contexts'}?`)) {
      return;
    }

    this.archiveMutation.mutate(
      { workspaceId, olderThanDays, timeZone: this.timeZone.effectiveTimeZone() },
      {
        onSuccess: (result) => this.lastArchivedCount.set(result.archivedCount),
      },
    );
  }

  retryPreview(): void {
    void this.archiveCandidatesQuery.refetch();
  }

  formatCutoff(value: string): string {
    return this.timeZone.formatDate(value);
  }

  formatLastInterval(value: string): string {
    return this.timeZone.formatDateTime(value);
  }

  private parseDays(value: string): number | null {
    const normalized = value.trim();
    if (!/^[1-9]\d*$/.test(normalized)) {
      return null;
    }

    const days = Number(normalized);
    return Number.isSafeInteger(days) && days <= MAX_ARCHIVE_DAYS ? days : null;
  }
}

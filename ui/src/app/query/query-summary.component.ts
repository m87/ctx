import { Component, computed, input, output, signal } from '@angular/core';
import { HlmSkeletonImports } from '@spartan-ng/helm/skeleton';
import { ContextQueryResult } from '../../api/query/query.service';
import { ContextListItem } from '../context/context-list-item.component';
import { ContextListComponent } from '../context/context-list.component';
import { DistributionComponent, DistributionItem } from '../shared/distribution.component';
import {
  ProjectTimeListComponent,
  ProjectTimeListItem,
} from '../shared/project-time-list.component';
import { summarizeContextsByProject, UNASSIGNED_PROJECT_ID } from '../shared/project-time-summary';
import { QueryErrorStateComponent } from '../shared/query-error-state.component';
import { colorHash, durationAsHM } from '../utils';

type SummaryView = 'contexts' | 'projects';

@Component({
  selector: 'ctx-query-summary',
  imports: [
    ContextListComponent,
    DistributionComponent,
    HlmSkeletonImports,
    ProjectTimeListComponent,
    QueryErrorStateComponent,
  ],
  template: `
    @if (showError()) {
      <ctx-query-error-state
        [error]="error()"
        [paused]="paused()"
        [resourceName]="resourceName()"
        [retrying]="retrying()"
        (retry)="retry.emit()"
      ></ctx-query-error-state>
    } @else if (loading()) {
      <div role="status" aria-label="Loading query summary">
        <span class="sr-only">Loading query summary</span>
        <div class="mb-6 grid grid-cols-2 gap-2.5 md:grid-cols-4">
          @for (item of skeletonItems; track item) {
            <div class="rounded-lg border bg-card px-3 py-2.5">
              <hlm-skeleton class="mb-2 h-2.5 w-20"></hlm-skeleton>
              <hlm-skeleton class="h-5 w-14"></hlm-skeleton>
            </div>
          }
        </div>
        <hlm-skeleton class="mb-4 h-8 w-40"></hlm-skeleton>
        <hlm-skeleton class="mb-6 h-2 w-full"></hlm-skeleton>
        <div class="flex flex-col gap-2">
          @for (item of contextSkeletonItems; track item) {
            <hlm-skeleton class="h-24 w-full rounded-lg"></hlm-skeleton>
          }
        </div>
      </div>
    } @else {
      <div class="mb-6 grid grid-cols-2 gap-2.5 md:grid-cols-4">
        <div class="rounded-lg border bg-card px-3 py-2.5">
          <div class="text-[10px] uppercase tracking-[0.08em] text-muted-foreground">
            Total tracked
          </div>
          <div class="mt-1 text-base font-semibold">{{ totalTracked() }}</div>
        </div>
        <div class="rounded-lg border bg-card px-3 py-2.5">
          <div class="text-[10px] uppercase tracking-[0.08em] text-muted-foreground">Contexts</div>
          <div class="mt-1 text-base font-semibold">{{ contexts().length }}</div>
        </div>
        <div class="rounded-lg border bg-card px-3 py-2.5">
          <div class="text-[10px] uppercase tracking-[0.08em] text-muted-foreground">Sessions</div>
          <div class="mt-1 text-base font-semibold">{{ result()?.totalSessions ?? 0 }}</div>
        </div>
        <div class="rounded-lg border bg-card px-3 py-2.5">
          <div class="text-[10px] uppercase tracking-[0.08em] text-muted-foreground">
            Top context
          </div>
          <div class="mt-1 truncate text-sm font-medium">{{ topContext() }}</div>
        </div>
      </div>

      <div
        class="mb-4 inline-flex rounded-lg bg-muted p-1"
        role="tablist"
        aria-label="Query summary view"
      >
        <button
          type="button"
          id="query-contexts-tab"
          class="rounded-md px-3 py-1.5 text-xs font-medium transition-colors"
          [class.bg-background]="summaryView() === 'contexts'"
          [class.shadow-sm]="summaryView() === 'contexts'"
          [class.text-foreground]="summaryView() === 'contexts'"
          [class.text-muted-foreground]="summaryView() !== 'contexts'"
          role="tab"
          aria-controls="query-summary-panel"
          [attr.aria-selected]="summaryView() === 'contexts'"
          (click)="summaryView.set('contexts')"
        >
          Contexts
        </button>
        <button
          type="button"
          id="query-projects-tab"
          class="rounded-md px-3 py-1.5 text-xs font-medium transition-colors"
          [class.bg-background]="summaryView() === 'projects'"
          [class.shadow-sm]="summaryView() === 'projects'"
          [class.text-foreground]="summaryView() === 'projects'"
          [class.text-muted-foreground]="summaryView() !== 'projects'"
          role="tab"
          aria-controls="query-summary-panel"
          [attr.aria-selected]="summaryView() === 'projects'"
          (click)="summaryView.set('projects')"
        >
          Projects
        </button>
      </div>

      <ctx-distribution
        class="mb-6 block"
        [label]="summaryView() === 'contexts' ? 'Context distribution' : 'Project distribution'"
        [items]="activeDistribution()"
        emptyMessage="No tracked time for this query."
      ></ctx-distribution>

      <div
        id="query-summary-panel"
        role="tabpanel"
        [attr.aria-labelledby]="
          summaryView() === 'contexts' ? 'query-contexts-tab' : 'query-projects-tab'
        "
      >
        @if (summaryView() === 'contexts') {
          <ctx-context-list
            [items]="contexts()"
            emptyMessage="No contexts match this query."
          ></ctx-context-list>
        } @else {
          <ctx-project-time-list
            [items]="projects()"
            emptyMessage="No tracked projects match this query."
          ></ctx-project-time-list>
        }
      </div>
    }
  `,
})
export class QuerySummaryComponent {
  readonly skeletonItems = [0, 1, 2, 3];
  readonly contextSkeletonItems = [0, 1, 2];
  readonly summaryView = signal<SummaryView>('contexts');
  readonly result = input<ContextQueryResult | null>(null);
  readonly loading = input(false);
  readonly showError = input(false);
  readonly error = input<unknown>(null);
  readonly paused = input(false);
  readonly retrying = input(false);
  readonly resourceName = input('query result');
  readonly retry = output<void>();

  readonly contexts = computed<ContextListItem[]>(() => {
    const result = this.result();
    const statsByContext = new Map(
      (result?.contextStats ?? []).map((stats) => [stats.contextId, stats]),
    );

    return (result?.contexts ?? [])
      .map((context) => {
        const stats = statsByContext.get(context.id);
        return {
          ...context,
          id: context.id,
          name: context.name,
          duration: durationAsHM(stats?.duration ?? 0).trim() || '0m',
          durationValue: stats?.duration ?? 0,
          percentage: stats?.percentage ?? 0,
          color: colorHash(context.id),
          sessions: stats?.intervalCount ?? 0,
          archived: context.archived ?? false,
          project: context.project,
        };
      })
      .sort(
        (left, right) =>
          (right.durationValue ?? 0) - (left.durationValue ?? 0) ||
          left.name.localeCompare(right.name),
      );
  });

  readonly contextDistribution = computed<DistributionItem[]>(() =>
    this.contexts()
      .filter((context) => (context.durationValue ?? 0) > 0)
      .map((context) => ({
        id: context.id,
        name: context.name,
        duration: context.duration,
        percentage: context.percentage,
        color: context.color,
      })),
  );

  readonly projects = computed<ProjectTimeListItem[]>(() =>
    summarizeContextsByProject(this.result() ?? { contexts: [], contextStats: [] }).map(
      (project) => ({
        ...project,
        duration: durationAsHM(project.duration).trim() || '0m',
        color:
          project.id === UNASSIGNED_PROJECT_ID ? '#94a3b8' : colorHash(`project:${project.id}`),
      }),
    ),
  );

  readonly projectDistribution = computed<DistributionItem[]>(() =>
    this.projects().map((project) => ({
      id: project.id,
      name: project.name,
      duration: project.duration,
      percentage: project.percentage,
      color: project.color,
    })),
  );

  readonly activeDistribution = computed(() =>
    this.summaryView() === 'contexts' ? this.contextDistribution() : this.projectDistribution(),
  );
  readonly totalTracked = computed(
    () => durationAsHM(this.result()?.totalDuration ?? 0).trim() || '0m',
  );
  readonly topContext = computed(() => {
    const top = this.contexts()[0];
    return top && (top.durationValue ?? 0) > 0 ? top.name : '-';
  });
}

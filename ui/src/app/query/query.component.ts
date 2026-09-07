import { Component, computed, inject, signal } from '@angular/core';
import { toSignal } from '@angular/core/rxjs-interop';
import { ActivatedRoute } from '@angular/router';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideTrash2 } from '@ng-icons/lucide';
import { Store } from '@ngxs/store';
import { injectMutation, injectQuery } from '@tanstack/angular-query-experimental';
import { map } from 'rxjs';
import { QueryQueries } from '../../api/query/query.queries';
import { SavedQueryMutations } from '../../api/query/saved-query.mutations';
import { SavedQueryQueries } from '../../api/query/saved-query.queries';
import { WorkspaceState } from '../sidebar/workspace.state';
import { QueryFormComponent } from './query-form.component';
import { QuerySummaryComponent } from './query-summary.component';

@Component({
  selector: 'ctx-query',
  imports: [NgIcon, QueryFormComponent, QuerySummaryComponent],
  providers: [provideIcons({ lucideTrash2 })],
  template: `
    <div class="h-full w-full overflow-y-auto p-4 md:p-6">
      <div class="mx-auto flex w-full max-w-5xl flex-col gap-6">
        @if (savedQueryId()) {
          <header>
            <div class="flex items-start justify-between gap-4">
              <div class="min-w-0">
                <div
                  class="text-[11px] font-semibold uppercase tracking-[0.08em] text-muted-foreground"
                >
                  Saved query
                </div>
                <h1 class="mt-1 truncate text-2xl font-semibold tracking-tight">
                  {{ savedQueryQuery.data()?.name ?? 'Loading query…' }}
                </h1>
              </div>
              @if (savedQueryQuery.data(); as savedQuery) {
                <button
                  type="button"
                  class="mt-1 flex size-9 shrink-0 items-center justify-center rounded-md border border-destructive/30 text-destructive transition-colors hover:bg-destructive/10 disabled:opacity-50"
                  [disabled]="deleteSavedQueryMutation.isPending()"
                  [attr.aria-label]="'Delete query ' + savedQuery.name"
                  title="Delete query"
                  (click)="deleteSavedQuery()"
                >
                  <ng-icon name="lucideTrash2"></ng-icon>
                </button>
              }
            </div>
            @if (savedQueryQuery.data(); as savedQuery) {
              <pre
                class="mt-3 overflow-x-auto rounded-lg border bg-muted/35 px-3 py-2 text-xs leading-5 text-muted-foreground"
              ><code>{{ savedQuery.query }}</code></pre>
            }
          </header>
        } @else {
          <header>
            <h1 class="text-2xl font-semibold tracking-tight">Query</h1>
            <p class="mt-1 max-w-2xl text-sm text-muted-foreground">
              Explore workspace contexts and save useful queries for later.
            </p>
          </header>

          <ctx-query-form
            [query]="queryInput()"
            [canSave]="canSave()"
            [saveExpanded]="saveExpanded()"
            [saveName]="saveName()"
            [savePending]="createSavedQueryMutation.isPending()"
            (queryChange)="queryInput.set($event)"
            (run)="runQuery()"
            (clear)="clearQuery()"
            (saveStart)="startSaving()"
            (saveNameChange)="saveName.set($event)"
            (saveConfirm)="saveQuery()"
            (saveCancel)="cancelSaving()"
          />

          @if (createSavedQueryMutation.isError()) {
            <p class="-mt-4 text-sm text-destructive" role="alert">
              Could not save this query. Please try again.
            </p>
          }
        }

        <ctx-query-summary
          [result]="queryResultQuery.data() ?? null"
          [loading]="isLoading()"
          [showError]="showError()"
          [error]="pageError()"
          [paused]="isPaused()"
          [retrying]="isRetrying()"
          [resourceName]="savedQueryId() ? 'saved query' : 'query result'"
          (retry)="retry()"
        />
      </div>
    </div>
  `,
  styles: `
    :host {
      display: block;
      width: 100%;
      height: 100%;
      min-height: 0;
    }
  `,
})
export class QueryComponent {
  private readonly route = inject(ActivatedRoute);
  private readonly store = inject(Store);
  private readonly queryQueries = inject(QueryQueries);
  private readonly savedQueryQueries = inject(SavedQueryQueries);
  private readonly savedQueryMutations = inject(SavedQueryMutations);

  readonly savedQueryId = toSignal(
    this.route.paramMap.pipe(map((params) => params.get('id') ?? '')),
    { initialValue: '' },
  );
  readonly activeWorkspaceId = this.store.selectSignal(WorkspaceState.selectedWorkspaceId);
  readonly queryInput = signal('');
  readonly submittedQuery = signal('');
  readonly saveExpanded = signal(false);
  readonly saveName = signal('');

  readonly savedQueryQuery = injectQuery(() => this.savedQueryQueries.get(this.savedQueryId()));
  readonly resultWorkspaceId = computed(
    () => this.savedQueryQuery.data()?.workspaceId ?? this.activeWorkspaceId() ?? '',
  );
  readonly resultQuery = computed(
    () => this.savedQueryQuery.data()?.query ?? this.submittedQuery(),
  );
  readonly resultEnabled = computed(
    () => this.savedQueryId().length === 0 || this.savedQueryQuery.data() !== undefined,
  );
  readonly queryResultQuery = injectQuery(() =>
    this.queryQueries.result(this.resultWorkspaceId(), this.resultQuery(), this.resultEnabled()),
  );
  readonly createSavedQueryMutation = injectMutation(() => this.savedQueryMutations.create());
  readonly deleteSavedQueryMutation = injectMutation(() => this.savedQueryMutations.delete());

  readonly canSave = computed(
    () => this.activeWorkspaceId() !== null && this.queryInput().trim().length > 0,
  );
  readonly showSavedQueryError = computed(
    () =>
      this.savedQueryId().length > 0 &&
      this.savedQueryQuery.data() === undefined &&
      (this.savedQueryQuery.isError() || this.savedQueryQuery.isPaused()),
  );
  readonly showResultError = computed(
    () =>
      this.resultEnabled() &&
      this.queryResultQuery.data() === undefined &&
      (this.queryResultQuery.isError() || this.queryResultQuery.isPaused()),
  );
  readonly showError = computed(() => this.showSavedQueryError() || this.showResultError());
  readonly isLoading = computed(
    () =>
      (this.savedQueryId().length > 0 && this.savedQueryQuery.isLoading()) ||
      (this.resultEnabled() && this.queryResultQuery.isLoading()),
  );
  readonly pageError = computed(() =>
    this.showSavedQueryError() ? this.savedQueryQuery.error() : this.queryResultQuery.error(),
  );
  readonly isPaused = computed(() =>
    this.showSavedQueryError() ? this.savedQueryQuery.isPaused() : this.queryResultQuery.isPaused(),
  );
  readonly isRetrying = computed(() =>
    this.showSavedQueryError()
      ? this.savedQueryQuery.isFetching()
      : this.queryResultQuery.isFetching(),
  );

  runQuery(): void {
    const query = this.queryInput().trim();
    if (query === this.submittedQuery()) {
      void this.queryResultQuery.refetch();
      return;
    }
    this.submittedQuery.set(query);
  }

  clearQuery(): void {
    this.queryInput.set('');
    this.submittedQuery.set('');
    this.cancelSaving();
  }

  startSaving(): void {
    if (!this.canSave()) {
      return;
    }
    this.createSavedQueryMutation.reset();
    this.saveName.set(this.queryInput().trim().slice(0, 60));
    this.saveExpanded.set(true);
  }

  cancelSaving(): void {
    this.saveExpanded.set(false);
    this.saveName.set('');
    this.createSavedQueryMutation.reset();
  }

  saveQuery(): void {
    const workspaceId = this.activeWorkspaceId();
    const name = this.saveName().trim();
    const query = this.queryInput().trim();
    if (!workspaceId || !name || !query || this.createSavedQueryMutation.isPending()) {
      return;
    }

    this.createSavedQueryMutation.mutate({ workspaceId, name, query });
  }

  deleteSavedQuery(): void {
    const query = this.savedQueryQuery.data();
    if (!query || this.deleteSavedQueryMutation.isPending()) {
      return;
    }
    if (!window.confirm(`Delete query "${query.name}"?`)) {
      return;
    }
    this.deleteSavedQueryMutation.mutate(query);
  }

  retry(): void {
    if (this.showSavedQueryError()) {
      void this.savedQueryQuery.refetch();
      return;
    }
    void this.queryResultQuery.refetch();
  }
}

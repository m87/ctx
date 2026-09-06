import { Component, computed, inject, signal } from '@angular/core';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { Context } from '../../api/context/context.service';
import { QueryQueries } from '../../api/query/query.queries';
import { QueryFormComponent } from './query-form.component';
import { QueryResultsComponent } from './query-results.component';

@Component({
  selector: 'ctx-query',
  imports: [QueryFormComponent, QueryResultsComponent],
  template: `
    <div class="h-full w-full overflow-y-auto p-4 md:p-6">
      <div class="mx-auto flex w-full max-w-5xl flex-col gap-6">
        <header>
          <h1 class="text-2xl font-semibold tracking-tight">Query</h1>
          <p class="mt-1 max-w-2xl text-sm text-muted-foreground">
            Explore all contexts with a query.
          </p>
        </header>

        <ctx-query-form
          [query]="queryInput()"
          (queryChange)="queryInput.set($event)"
          (run)="runQuery()"
          (clear)="clearQuery()"
        />

        <ctx-query-results
          [contexts]="results()"
          [summary]="resultSummary()"
          [loading]="queryResultQuery.isLoading()"
          [showError]="showContextsError()"
          [error]="queryResultQuery.error()"
          [paused]="queryResultQuery.isPaused()"
          [retrying]="queryResultQuery.isFetching()"
          (retry)="retryContexts()"
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
  private readonly queryQueries = inject(QueryQueries);

  readonly queryInput = signal('');
  readonly submittedQuery = signal('');
  readonly queryResultQuery = injectQuery(() => this.queryQueries.contexts(this.submittedQuery()));
  readonly results = computed<readonly Context[]>(() => this.queryResultQuery.data() ?? []);
  readonly showContextsError = computed(
    () =>
      this.queryResultQuery.data() === undefined &&
      (this.queryResultQuery.isError() || this.queryResultQuery.isPaused()),
  );
  readonly resultSummary = computed(() => {
    if (this.queryResultQuery.isLoading()) {
      return 'Loading contexts…';
    }
    const count = this.results().length;
    return `${count} ${count === 1 ? 'context' : 'contexts'}`;
  });

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
  }

  retryContexts(): void {
    void this.queryResultQuery.refetch();
  }
}

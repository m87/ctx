import { Component, input, output } from '@angular/core';
import { HlmSkeletonImports } from '@spartan-ng/helm/skeleton';
import { Context } from '../../api/context/context.service';
import { QueryErrorStateComponent } from '../shared/query-error-state.component';
import { QueryResultItemComponent } from './query-result-item.component';

@Component({
  selector: 'ctx-query-results',
  imports: [HlmSkeletonImports, QueryErrorStateComponent, QueryResultItemComponent],
  template: `
    <section class="min-h-0" aria-labelledby="query-results-title">
      <div class="mb-3 flex items-center justify-between gap-3">
        <div>
          <h2 id="query-results-title" class="text-sm font-semibold">Results</h2>
          <p class="text-xs text-muted-foreground">{{ summary() }}</p>
        </div>
      </div>

      @if (showError()) {
        <ctx-query-error-state
          [error]="error()"
          [paused]="paused()"
          resourceName="contexts"
          [retrying]="retrying()"
          (retry)="retry.emit()"
        ></ctx-query-error-state>
      } @else if (loading()) {
        <div class="grid gap-2" role="status" aria-label="Loading contexts">
          <span class="sr-only">Loading contexts</span>
          @for (item of contextSkeletonItems; track item) {
            <div class="rounded-lg border bg-card p-4">
              <div class="flex items-center gap-3">
                <hlm-skeleton class="size-2.5 shrink-0 rounded-sm"></hlm-skeleton>
                <div class="min-w-0 flex-1">
                  <hlm-skeleton class="mb-2 h-3 w-40 max-w-full"></hlm-skeleton>
                  <hlm-skeleton class="h-2.5 w-64 max-w-full"></hlm-skeleton>
                </div>
              </div>
            </div>
          }
        </div>
      } @else if (contexts().length === 0) {
        <div class="rounded-xl border border-dashed bg-muted/20 px-6 py-10 text-center">
          <p class="text-sm font-medium">No contexts</p>
          <p class="mt-1 text-xs text-muted-foreground">There are no contexts to display yet.</p>
        </div>
      } @else {
        <div class="grid gap-2" role="list">
          @for (context of contexts(); track context.id) {
            <ctx-query-result-item [context]="context" />
          }
        </div>
      }
    </section>
  `,
})
export class QueryResultsComponent {
  readonly contextSkeletonItems = [0, 1, 2, 3];
  readonly contexts = input<readonly Context[]>([]);
  readonly summary = input('');
  readonly loading = input(false);
  readonly showError = input(false);
  readonly error = input<unknown>(null);
  readonly paused = input(false);
  readonly retrying = input(false);
  readonly retry = output<void>();
}

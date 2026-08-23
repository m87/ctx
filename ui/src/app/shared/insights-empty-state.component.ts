import { Component } from '@angular/core';

@Component({
  selector: 'ctx-insights-empty-state',
  template: `
    <div
      class="min-h-40 w-full rounded-lg border border-dashed bg-muted/20 px-6 py-10 flex flex-col items-center justify-center text-center"
    >
      <p class="text-sm font-medium">No insights yet</p>
      <p class="mt-1 max-w-md text-xs text-muted-foreground">
        Insights will appear here when they become available.
      </p>
    </div>
  `,
})
export class InsightsEmptyStateComponent {}

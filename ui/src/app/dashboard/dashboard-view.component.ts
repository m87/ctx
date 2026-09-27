import { Component, computed, input } from '@angular/core';
import { Dashboard, DashboardScope } from '../../api/dashboard/dashboard.service';
import { DashboardCanvasComponent } from '../editor/dashboard/dashboard-canvas.component';
import { injectDashboardMetrics } from './dashboard-metrics';
import { normalizeDashboardDefinition } from '../editor/dashboard/dashboard-definition';

@Component({
  selector: 'ctx-dashboard-view',
  imports: [DashboardCanvasComponent],
  host: { class: 'block min-h-40 w-full bg-muted p-3 text-foreground dark:bg-background sm:p-6' },
  template: `
    @if (content(); as content) {
      @if (content.error) {
        <div class="ui-empty-state text-destructive" role="alert">{{ content.error }}</div>
      } @else if (content.definition?.widgets?.length) {
        <ctx-dashboard-canvas
          class="mx-auto"
          [widgets]="content.definition?.widgets ?? []"
          [results]="metrics.results()"
          [error]="metrics.error()"
        ></ctx-dashboard-canvas>
      } @else {
        <div class="ui-empty-state">No widgets in this dashboard yet.</div>
      }
    }
  `,
})
export class DashboardViewComponent {
  readonly dashboard = input.required<Dashboard>();
  readonly scope = input<DashboardScope>({});
  readonly metrics = injectDashboardMetrics(
    () => this.dashboard(),
    () => this.scope(),
  );
  readonly content = computed(() => {
    try {
      return { definition: normalizeDashboardDefinition(this.dashboard().definition), error: '' };
    } catch {
      return {
        definition: null,
        error: 'This dashboard cannot be displayed. Check its definition in the editor.',
      };
    }
  });
}

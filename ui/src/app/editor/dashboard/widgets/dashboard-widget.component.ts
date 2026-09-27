import { Component, computed, input } from '@angular/core';
import { DashboardMetricResult } from '../../../../api/dashboard/dashboard.service';
import { DashboardWidgetDefinition } from '../dashboard-definition';
import { PERIODS } from '../metric-widget-definition';
import { TextWidgetComponent } from './text-widget.component';
import { KpiWidgetComponent } from './kpi-widget.component';
import { GaugeWidgetComponent } from './gauge-widget.component';

@Component({
  selector: 'ctx-dashboard-widget',
  imports: [TextWidgetComponent, KpiWidgetComponent, GaugeWidgetComponent],
  host: { class: 'block h-full min-h-0 min-w-0' },
  template: `
    @let item = widget();
    @switch (item.type) {
      @case ('text') {
        <ctx-text-widget [properties]="item.properties" [editing]="editing()" />
      }
      @case ('kpi') {
        <ctx-kpi-widget
          [properties]="item.properties"
          [result]="result()"
          [error]="error()"
          [period]="period()"
        />
      }
      @case ('gauge') {
        <ctx-gauge-widget
          [properties]="item.properties"
          [result]="result()"
          [error]="error()"
          [period]="period()"
        />
      }
    }
  `,
})
export class DashboardWidgetComponent {
  readonly widget = input.required<DashboardWidgetDefinition>();
  readonly editing = input(false);
  readonly result = input<DashboardMetricResult>();
  readonly error = input('');
  readonly period = computed(() => {
    const result = this.result();
    if (result) {
      const options: Intl.DateTimeFormatOptions = {
        month: 'short',
        day: 'numeric',
        timeZone: result.timeZone,
      };
      const start = new Date(result.start).toLocaleDateString(undefined, options);
      const end = new Date(new Date(result.end).getTime() - 1).toLocaleDateString(
        undefined,
        options,
      );
      return start === end ? start : start + ' – ' + end;
    }
    const widget = this.widget();
    return widget.type === 'text'
      ? ''
      : (PERIODS.find((period) => period.id === widget.properties.period)?.label ??
          'Dashboard period');
  });
}

import { Component, input, output } from '@angular/core';
import { SearchSelectOption } from '../../../shared/search-select.component';
import { DashboardWidgetDefinition, TextWidgetProperties } from '../dashboard-definition';
import { TextWidgetPropertiesComponent } from './text-widget-properties.component';
import { MetricWidgetPropertiesComponent } from './metric-widget-properties.component';
import { AdvancedWidgetPropertiesComponent } from './advanced-widget-properties.component';

@Component({
  selector: 'ctx-widget-properties',
  imports: [
    TextWidgetPropertiesComponent,
    MetricWidgetPropertiesComponent,
    AdvancedWidgetPropertiesComponent,
  ],
  host: { class: 'block min-h-0 overflow-y-auto' },
  template: `
    @let item = widget();
    @if (item.type === 'text') {
      <ctx-text-widget-properties
        [widget]="item"
        [disabled]="disabled()"
        [idSuffix]="idSuffix()"
        (propertiesChange)="setTextProperties($event)"
        (queryChange)="setQuery($event)"
        (deleteRequested)="deleteRequested.emit()"
      />
    } @else if (item.type === 'kpi' || item.type === 'gauge' || item.type === 'radar') {
      <ctx-metric-widget-properties
        [widget]="item"
        [disabled]="disabled()"
        [idSuffix]="idSuffix()"
        [daily]="daily()"
        [today]="today()"
        [projects]="projects()"
        [contexts]="contexts()"
        [tags]="tags()"
        (widgetChange)="widgetChange.emit($event)"
        (deleteRequested)="deleteRequested.emit()"
      />
    } @else if (
      item.type === 'bar' ||
      item.type === 'heatmap' ||
      item.type === 'line' ||
      item.type === 'table'
    ) {
      <ctx-advanced-widget-properties
        [widget]="item"
        [disabled]="disabled()"
        (widgetChange)="widgetChange.emit($event)"
      />
    }
  `,
})
export class WidgetPropertiesComponent {
  readonly widget = input.required<DashboardWidgetDefinition>();
  readonly disabled = input(false);
  readonly daily = input(false);
  readonly today = input('');
  readonly idSuffix = input('');
  readonly projects = input<readonly SearchSelectOption[]>([]);
  readonly contexts = input<readonly SearchSelectOption[]>([]);
  readonly tags = input<readonly string[]>([]);
  readonly widgetChange = output<DashboardWidgetDefinition>();
  readonly deleteRequested = output<void>();
  setTextProperties(properties: TextWidgetProperties): void {
    const widget = this.widget();
    if (widget.type === 'text') this.widgetChange.emit({ ...widget, properties });
  }
  setQuery(query: string): void {
    this.widgetChange.emit({ ...this.widget(), query });
  }
}

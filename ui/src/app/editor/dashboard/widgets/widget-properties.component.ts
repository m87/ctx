import { Component, input, output } from '@angular/core';
import { SearchSelectOption } from '../../../shared/search-select.component';
import { DashboardWidgetDefinition, TextWidgetProperties } from '../dashboard-definition';
import { TextWidgetPropertiesComponent } from './text-widget-properties.component';
import { MetricWidgetPropertiesComponent } from './metric-widget-properties.component';

@Component({
  selector: 'ctx-widget-properties',
  imports: [TextWidgetPropertiesComponent, MetricWidgetPropertiesComponent],
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
    } @else {
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

import { Component, computed, input, output } from '@angular/core';
import { HlmButtonImports } from '@spartan-ng/helm/button';
import { HlmInputImports } from '@spartan-ng/helm/input';
import { HlmLabelImports } from '@spartan-ng/helm/label';
import { SearchSelectOption } from '../../../shared/search-select.component';
import { metricInfo, RadarWidgetProperties } from '../metric-widget-definition';

@Component({
  selector: 'ctx-radar-widget-properties',
  imports: [HlmButtonImports, HlmInputImports, HlmLabelImports],
  host: { class: 'contents' },
  template: `
    @let p = properties();
    <section class="flex flex-col gap-3 border-t border-border/65 pt-4">
      <h4 class="ui-section-label">Axes</h4>
      <div class="ui-field">
        <label hlmLabel [for]="fieldId('axes-by')">Group categories by</label>
        <select
          class="ui-select-trigger w-full"
          [id]="fieldId('axes-by')"
          [value]="p.axesBy"
          (change)="setAxesBy($event)"
        >
          <option value="weekday">Day of week</option>
          <option value="project">Projects</option>
          <option value="context">Contexts</option>
        </select>
      </div>
      @if (p.axesBy === 'weekday') {
        <p class="text-xs text-muted-foreground">
          Each axis sums all observed occurrences of that weekday in the selected period.
        </p>
      } @else {
        <div
          class="flex flex-col gap-1"
          role="group"
          [attr.aria-label]="p.axesBy === 'project' ? 'Radar project axes' : 'Radar context axes'"
        >
          @for (option of radarAxisOptions(); track option.value) {
            <div class="flex min-h-8 items-center gap-2 rounded-md px-2 text-xs hover:bg-muted/40">
              <label class="flex min-w-0 flex-1 items-center gap-2">
                <input
                  type="checkbox"
                  class="size-4 accent-primary focus-visible:ring-2 focus-visible:ring-ring"
                  [checked]="p.axisKeys.includes(option.value)"
                  [disabled]="!p.axisKeys.includes(option.value) && p.axisKeys.length >= 8"
                  (change)="toggleAxis(option.value, $event)"
                />
                <span class="truncate">{{ option.label }}</span>
              </label>
              @if (axisIndex(option.value) >= 0) {
                <span class="ml-auto flex gap-1">
                  <button
                    hlmBtn
                    variant="ghost"
                    size="icon-sm"
                    type="button"
                    [disabled]="axisIndex(option.value) === 0"
                    [attr.aria-label]="'Move ' + option.label + ' axis up'"
                    (click)="moveAxis(axisIndex(option.value), -1)"
                  >
                    ↑
                  </button>
                  <button
                    hlmBtn
                    variant="ghost"
                    size="icon-sm"
                    type="button"
                    [disabled]="axisIndex(option.value) === p.axisKeys.length - 1"
                    [attr.aria-label]="'Move ' + option.label + ' axis down'"
                    (click)="moveAxis(axisIndex(option.value), 1)"
                  >
                    ↓
                  </button>
                </span>
              }
            </div>
          }
          <p class="text-xs text-muted-foreground" role="status">
            Choose 3–8 axes. Selected: {{ p.axisKeys.length }}.
          </p>
        </div>
      }
    </section>
    <section class="flex flex-col gap-3 border-t border-border/65 pt-4">
      <h4 class="ui-section-label">Scale and target</h4>
      @if (unitResetNotice()) {
        <p class="text-xs text-muted-foreground" role="status">{{ unitResetNotice() }}</p>
      }
      <div class="ui-field">
        <label hlmLabel [for]="fieldId('scale-maximum')">Scale maximum ({{ unitLabel() }})</label
        ><input
          hlmInput
          class="dashboard-field-control h-9"
          type="number"
          min="0"
          [step]="metric().unit === 'count' ? 1 : 'any'"
          [id]="fieldId('scale-maximum')"
          [value]="displayNumber(p.scaleMaximum) ?? ''"
          (change)="setRadarNumber('scaleMaximum', $event)"
        />
      </div>
      <div class="ui-field">
        <label hlmLabel [for]="fieldId('target-value')">Target per axis ({{ unitLabel() }})</label
        ><input
          hlmInput
          class="dashboard-field-control h-9"
          type="number"
          min="0"
          [step]="metric().unit === 'count' ? 1 : 'any'"
          [id]="fieldId('target-value')"
          [value]="displayNumber(p.targetValue) ?? ''"
          (change)="setRadarNumber('targetValue', $event)"
        />
      </div>
      <p class="text-xs text-muted-foreground">
        Targets are reference values; the chart does not rate higher or lower values.
      </p>
    </section>
    <section class="flex flex-col gap-3 border-t border-border/65 pt-4">
      <h4 class="ui-section-label">Appearance</h4>
      @for (field of radarVisibilityFields; track field.key) {
        <label hlmLabel class="flex items-center gap-2"
          ><input
            type="checkbox"
            class="size-4 accent-primary focus-visible:ring-2 focus-visible:ring-ring"
            [checked]="p[field.key]"
            (change)="setRadarBoolean(field.key, $event)"
          />{{ field.label }}</label
        >
      }
    </section>
  `,
})
export class RadarWidgetPropertiesComponent {
  readonly properties = input.required<RadarWidgetProperties>();
  readonly widgetId = input.required<string>();
  readonly idSuffix = input('');
  readonly projects = input<readonly SearchSelectOption[]>([]);
  readonly contexts = input<readonly SearchSelectOption[]>([]);
  readonly unitResetNotice = input('');
  readonly propertiesChange = output<Partial<RadarWidgetProperties>>();
  readonly metric = computed(() => metricInfo(this.properties().metric));
  readonly unitLabel = computed(() =>
    this.metric().unit === 'ms' ? 'hours' : this.metric().unit === 'percent' ? '%' : 'count',
  );
  readonly radarVisibilityFields = [
    { key: 'showValues', label: 'Show values by axis' },
    { key: 'showGrid', label: 'Show grid' },
    { key: 'showFill', label: 'Fill complete profile' },
    { key: 'showSummary', label: 'Show period summary' },
  ] as const;
  readonly radarAxisOptions = computed(() => {
    const properties = this.properties();
    const options =
      properties.axesBy === 'project'
        ? [
            { value: 'project:none', label: 'No project' },
            ...this.projects().map((option) => ({
              value: 'project:' + option.value,
              label: option.label,
            })),
          ]
        : this.contexts().map((option) => ({
            value: 'context:' + option.value,
            label: option.label,
          }));
    const available = new Set(options.map((option) => option.value));
    return [
      ...options,
      ...properties.axisKeys
        .filter((key) => !available.has(key))
        .map((key) => ({ value: key, label: `Unavailable category (${key})` })),
    ];
  });

  fieldId(field: string): string {
    return 'metric-' + this.widgetId() + '-' + field + this.idSuffix();
  }
  displayNumber(value: number | null): number | null {
    return value === null ? null : value / (this.metric().unit === 'ms' ? 3600000 : 1);
  }
  setAxesBy(event: Event): void {
    const axesBy = (event.target as HTMLSelectElement).value as RadarWidgetProperties['axesBy'];
    this.propertiesChange.emit({ axesBy, axisKeys: [] });
  }
  axisIndex(key: string): number {
    return this.properties().axisKeys.indexOf(key);
  }
  toggleAxis(key: string, event: Event): void {
    const checked = (event.target as HTMLInputElement).checked;
    const axisKeys = checked
      ? [...this.properties().axisKeys, key]
      : this.properties().axisKeys.filter((item) => item !== key);
    this.propertiesChange.emit({ axisKeys });
  }
  moveAxis(index: number, delta: number): void {
    const axisKeys = [...this.properties().axisKeys];
    const target = index + delta;
    if (target < 0 || target >= axisKeys.length) return;
    [axisKeys[index], axisKeys[target]] = [axisKeys[target], axisKeys[index]];
    this.propertiesChange.emit({ axisKeys });
  }
  setRadarNumber(key: 'scaleMaximum' | 'targetValue', event: Event): void {
    const input = event.target as HTMLInputElement;
    const value =
      input.value === '' ? null : input.valueAsNumber * (this.metric().unit === 'ms' ? 3600000 : 1);
    if (value !== null && !Number.isFinite(value)) return;
    this.propertiesChange.emit({ [key]: value });
  }
  setRadarBoolean(key: 'showValues' | 'showGrid' | 'showFill' | 'showSummary', event: Event): void {
    this.propertiesChange.emit({ [key]: (event.target as HTMLInputElement).checked });
  }
}

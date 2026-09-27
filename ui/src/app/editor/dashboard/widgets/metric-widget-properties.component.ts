import { Component, computed, effect, input, output, signal } from '@angular/core';
import { HlmButtonImports } from '@spartan-ng/helm/button';
import { HlmInputImports } from '@spartan-ng/helm/input';
import { HlmLabelImports } from '@spartan-ng/helm/label';
import { HlmTextareaImports } from '@spartan-ng/helm/textarea';
import { SearchDropdownSelectComponent } from '../../../shared/search-dropdown-select.component';
import { SearchSelectOption } from '../../../shared/search-select.component';
import { GaugeWidgetDefinition, KpiWidgetDefinition } from '../dashboard-definition';
import {
  METRICS,
  PERIODS,
  METRIC_TONES,
  MetricId,
  MetricWidgetProperties,
  GaugeWidgetProperties,
  MetricTone,
  metricInfo,
  metricSettingsError,
} from '../metric-widget-definition';

type MetricWidget = KpiWidgetDefinition | GaugeWidgetDefinition;

@Component({
  selector: 'ctx-metric-widget-properties',
  imports: [
    HlmButtonImports,
    HlmInputImports,
    HlmLabelImports,
    HlmTextareaImports,
    SearchDropdownSelectComponent,
  ],
  host: { class: 'block min-h-0 overflow-y-auto' },
  template: `
    <fieldset class="flex min-w-0 flex-col gap-4 p-3" [disabled]="disabled()">
      <legend class="sr-only">{{ widget().type === 'kpi' ? 'KPI' : 'Gauge' }} settings</legend>
      @for (field of textFields; track field.key) {
        <div class="ui-field">
          <label hlmLabel [for]="fieldId(field.key)">{{ field.label }}</label>
          <input
            hlmInput
            class="dashboard-field-control h-9"
            [id]="fieldId(field.key)"
            [value]="widget().properties[field.key]"
            [placeholder]="field.key === 'title' ? metric().label : 'Optional description'"
            maxlength="500"
            (input)="setText(field.key, $event)"
          />
        </div>
      }
      <section class="flex flex-col gap-3 border-t border-border/65 pt-4">
        <h4 class="ui-section-label">Data</h4>
        <div class="ui-field">
          <label hlmLabel [for]="fieldId('metric')">Metric</label>
          <select
            class="ui-select-trigger w-full"
            [id]="fieldId('metric')"
            [value]="widget().properties.metric"
            (change)="setMetric($event)"
          >
            @for (metric of metrics; track metric.id) {
              <option [value]="metric.id">{{ metric.label }}</option>
            }
          </select>
        </div>
        <div class="ui-field">
          <label hlmLabel [for]="fieldId('period')">Period</label>
          <select
            class="ui-select-trigger w-full"
            [id]="fieldId('period')"
            [value]="widget().properties.period"
            [disabled]="disabled() || daily()"
            (change)="setPeriod($event)"
          >
            <option value="inherit">Dashboard period</option>
            @for (period of periods; track period.id) {
              <option [value]="period.id">{{ period.label }}</option>
            }
            <option value="custom">Custom dates</option>
          </select>
          @if (daily()) {
            <p class="text-xs text-muted-foreground">
              Daily Insights always uses the selected summary date.
            </p>
          }
        </div>
        @if (widget().properties.period === 'custom') {
          @for (field of dateFields; track field.key) {
            <div class="ui-field">
              <label hlmLabel [for]="fieldId(field.key)">{{ field.label }}</label>
              <input
                hlmInput
                class="dashboard-field-control h-9"
                type="date"
                [id]="fieldId(field.key)"
                [value]="widget().properties[field.key]"
                (change)="setDate(field.key, $event)"
              />
            </div>
          }
        }
        <label hlmLabel class="flex items-center gap-2">
          <input
            type="checkbox"
            class="size-4 accent-primary focus-visible:ring-2 focus-visible:ring-ring"
            [checked]="widget().properties.includeRunning"
            (change)="setBoolean('includeRunning', $event)"
          />
          Include running sessions
        </label>
        <details>
          <summary
            class="cursor-pointer rounded-md text-xs text-muted-foreground focus-visible:ring-2 focus-visible:ring-ring"
          >
            Filters
          </summary>
          <div class="mt-3 flex flex-col gap-3">
            <div class="ui-field">
              <label hlmLabel [for]="fieldId('project')">Project</label>
              <ctx-search-dropdown-select
                [inputId]="fieldId('project')"
                ariaLabel="Metric project"
                panelWidth="100%"
                [disabled]="disabled()"
                [options]="projectOptions()"
                [value]="widget().properties.projectId"
                (selectionChange)="update({ projectId: $event })"
              />
            </div>
            <div class="ui-field">
              <label hlmLabel [for]="fieldId('context')">Context</label>
              <ctx-search-dropdown-select
                [inputId]="fieldId('context')"
                ariaLabel="Metric context"
                panelWidth="100%"
                [disabled]="disabled()"
                [options]="contextOptions()"
                [value]="widget().properties.contextId"
                (selectionChange)="update({ contextId: $event })"
              />
            </div>
            <div class="ui-field">
              <label hlmLabel [for]="fieldId('tag')">Tag</label>
              <ctx-search-dropdown-select
                [inputId]="fieldId('tag')"
                ariaLabel="Metric tag"
                panelWidth="100%"
                [disabled]="disabled()"
                [options]="tagOptions()"
                [value]="widget().properties.tag"
                (selectionChange)="update({ tag: $event })"
              />
            </div>
            <div class="ui-field">
              <label hlmLabel [for]="fieldId('archived')">Archived contexts</label>
              <select
                class="ui-select-trigger w-full"
                [id]="fieldId('archived')"
                [value]="widget().properties.archived"
                (change)="setArchived($event)"
              >
                <option value="include">Include</option>
                <option value="exclude">Exclude</option>
                <option value="only">Only archived</option>
              </select>
            </div>
            <p class="text-xs text-muted-foreground">
              Filters narrow the dashboard scope. Tags filter contexts; they are not metric values.
            </p>
          </div>
        </details>
      </section>
      <section class="flex flex-col gap-3 border-t border-border/65 pt-4">
        <h4 class="ui-section-label">Appearance</h4>
        <div class="ui-field">
          <label hlmLabel [for]="fieldId('format')">Value format</label>
          <select
            class="ui-select-trigger w-full"
            [id]="fieldId('format')"
            [value]="widget().properties.format"
            (change)="setFormat($event)"
          >
            @if (metric().unit === 'ms') {
              <option value="duration">Hours and minutes</option>
              <option value="hours">Decimal hours</option>
            }
            @if (metric().unit === 'count') {
              <option value="number">Number</option>
            }
            @if (metric().unit === 'percent') {
              <option value="percent">Percentage</option>
            }
          </select>
        </div>
        @if (widget().properties.format !== 'duration' && metric().unit !== 'count') {
          <div class="ui-field">
            <label hlmLabel [for]="fieldId('precision')">Decimal places</label>
            <select
              class="ui-select-trigger w-full"
              [id]="fieldId('precision')"
              [value]="widget().properties.precision"
              (change)="setPrecision($event)"
            >
              <option value="0">0</option>
              <option value="1">1</option>
              <option value="2">2</option>
            </select>
          </div>
        }
      </section>
      @if (gauge(); as p) {
        <section class="flex flex-col gap-3 border-t border-border/65 pt-4">
          <h4 class="ui-section-label">Scale & goal</h4>
          @if (unitResetNotice()) {
            <p class="text-xs text-muted-foreground" role="status">{{ unitResetNotice() }}</p>
          }
          <p class="text-xs text-muted-foreground">Values in {{ unitLabel() }}.</p>
          <div class="ui-field">
            <label hlmLabel [for]="fieldId('goal-mode')">Goal</label>
            <select
              class="ui-select-trigger w-full"
              [id]="fieldId('goal-mode')"
              [value]="p.goalMode"
              (change)="setGoalMode($event)"
            >
              <option value="at_least">Reach at least</option>
              <option value="at_most">Stay at most</option>
              <option value="informational">Informational</option>
            </select>
          </div>
          @for (field of gaugeFields; track field.key) {
            <div class="ui-field">
              <label hlmLabel [for]="fieldId(field.key)"
                >{{ field.label }} ({{ unitLabel() }})</label
              >
              <input
                hlmInput
                class="dashboard-field-control h-9"
                type="number"
                min="0"
                [step]="metric().unit === 'count' ? 1 : 'any'"
                [id]="fieldId(field.key)"
                [value]="displayNumber(p[field.key])"
                [attr.aria-invalid]="fieldError(field.key) ? 'true' : null"
                [attr.aria-describedby]="
                  fieldError(field.key) ? fieldId(field.key) + '-error' : null
                "
                (change)="setGaugeNumber(field.key, $event)"
              />
              @if (fieldError(field.key)) {
                <p
                  [id]="fieldId(field.key) + '-error'"
                  class="text-xs text-destructive"
                  role="alert"
                >
                  {{ error() }}
                </p>
              }
            </div>
          }
          @for (field of visibilityFields; track field.key) {
            <label hlmLabel class="flex items-center gap-2">
              <input
                type="checkbox"
                class="size-4 accent-primary focus-visible:ring-2 focus-visible:ring-ring"
                [checked]="p[field.key]"
                (change)="setBoolean(field.key, $event)"
              />
              {{ field.label }}
            </label>
          }
          <details>
            <summary
              class="cursor-pointer rounded-md text-xs text-muted-foreground focus-visible:ring-2 focus-visible:ring-ring"
            >
              Color thresholds
            </summary>
            <p class="mt-3 text-xs text-muted-foreground">
              Without thresholds, the color follows the goal. Each threshold applies from its value
              upward.
            </p>
            @for (threshold of p.thresholds; track $index; let i = $index) {
              <div class="mt-2 flex items-center gap-2">
                <input
                  hlmInput
                  class="dashboard-field-control h-9 min-w-0 flex-1"
                  type="number"
                  [attr.aria-label]="'Threshold ' + (i + 1) + ' in ' + unitLabel()"
                  [step]="metric().unit === 'count' ? 1 : 'any'"
                  [value]="displayNumber(threshold.value)"
                  (change)="setThresholdValue(i, $event)"
                />
                <select
                  class="ui-select-trigger min-w-0 flex-1"
                  [attr.aria-label]="'Threshold ' + (i + 1) + ' color'"
                  [value]="threshold.tone"
                  (change)="setThresholdTone(i, $event)"
                >
                  @for (tone of tones; track tone) {
                    <option [value]="tone">{{ tone }}</option>
                  }
                </select>
                <button
                  hlmBtn
                  variant="ghost"
                  size="icon"
                  type="button"
                  [attr.aria-label]="'Remove threshold ' + (i + 1)"
                  (click)="removeThreshold(i)"
                >
                  ×
                </button>
              </div>
            }
            <button
              hlmBtn
              variant="outline"
              size="sm"
              type="button"
              class="mt-3"
              [disabled]="!canAddThreshold()"
              (click)="addThreshold()"
            >
              Add threshold
            </button>
          </details>
        </section>
      }
      @if (error()) {
        <p class="text-xs text-destructive" role="alert">{{ error() }}</p>
      }
      <details class="border-t border-border/65 pt-3">
        <summary
          class="cursor-pointer rounded-md text-xs text-muted-foreground focus-visible:ring-2 focus-visible:ring-ring"
        >
          Widget query
        </summary>
        <div class="ui-field mt-3">
          <label hlmLabel [for]="fieldId('query')">Query</label>
          <textarea
            hlmTextarea
            class="dashboard-field-control min-h-24"
            [id]="fieldId('query')"
            [value]="widget().query"
            maxlength="8000"
            (input)="setQuery($event)"
          ></textarea>
          <p class="text-xs text-muted-foreground">
            Query text is saved for future use. It does not filter data yet.
          </p>
        </div>
      </details>
      <button
        hlmBtn
        variant="ghost"
        type="button"
        class="justify-start text-destructive"
        (click)="deleteRequested.emit()"
      >
        Delete widget
      </button>
    </fieldset>
  `,
})
export class MetricWidgetPropertiesComponent {
  readonly widget = input.required<MetricWidget>();
  readonly disabled = input(false);
  readonly daily = input(false);
  readonly today = input('');
  readonly idSuffix = input('');
  readonly projects = input<readonly SearchSelectOption[]>([]);
  readonly contexts = input<readonly SearchSelectOption[]>([]);
  readonly tags = input<readonly string[]>([]);
  readonly widgetChange = output<MetricWidget>();
  readonly deleteRequested = output<void>();
  readonly metrics = METRICS;
  readonly periods = PERIODS;
  readonly tones = METRIC_TONES;
  readonly error = signal('');
  readonly unitResetNotice = signal('');
  private readonly resetFeedbackOnWidgetChange = effect(() => {
    this.widget().id;
    this.error.set('');
    this.unitResetNotice.set('');
  });
  readonly textFields = [
    { key: 'title', label: 'Title' },
    { key: 'description', label: 'Description' },
  ] as const;
  readonly dateFields = [
    { key: 'startDate', label: 'From' },
    { key: 'endDate', label: 'Through' },
  ] as const;
  readonly gaugeFields = [
    { key: 'minimum', label: 'Minimum' },
    { key: 'maximum', label: 'Maximum' },
    { key: 'target', label: 'Target' },
  ] as const;
  readonly visibilityFields = [
    { key: 'showValue', label: 'Show value' },
    { key: 'showTarget', label: 'Show target' },
    { key: 'showPercentage', label: 'Show percentage of target' },
  ] as const;
  readonly metric = computed(() => metricInfo(this.widget().properties.metric));
  readonly gauge = computed(() => {
    const widget = this.widget();
    return widget.type === 'gauge' ? widget.properties : null;
  });
  readonly unitLabel = computed(() =>
    this.metric().unit === 'ms' ? 'hours' : this.metric().unit === 'percent' ? '%' : 'count',
  );
  readonly projectOptions = computed(() => [
    { value: '', label: 'All projects' },
    ...this.projects(),
  ]);
  readonly contextOptions = computed(() => [
    { value: '', label: 'All contexts' },
    ...this.contexts(),
  ]);
  readonly tagOptions = computed(() => [
    { value: '', label: 'All tags' },
    ...this.tags().map((tag) => ({ value: tag, label: tag })),
  ]);
  readonly canAddThreshold = computed(() => {
    const p = this.gauge();
    return (
      !!p && p.thresholds.length < 8 && (p.thresholds.at(-1)?.value ?? p.minimum - 1) < p.maximum
    );
  });

  fieldId(field: string): string {
    return 'metric-' + this.widget().id + '-' + field + this.idSuffix();
  }
  displayNumber(value: number | null): number | null {
    return value === null ? null : value / (this.metric().unit === 'ms' ? 3600000 : 1);
  }
  fieldError(key: string): boolean {
    const message = this.error();
    if (message.includes('Maximum must')) return key === 'minimum' || key === 'maximum';
    if (message.includes('target'))
      return (
        key === 'target' || (message.includes('scale') && ['minimum', 'maximum'].includes(key))
      );
    if (message.includes('Counts require')) return ['minimum', 'maximum', 'target'].includes(key);
    return false;
  }
  setText(key: 'title' | 'description', event: Event): void {
    this.update({ [key]: (event.target as HTMLInputElement).value });
  }
  setQuery(event: Event): void {
    if (!this.disabled())
      this.widgetChange.emit({
        ...this.widget(),
        query: (event.target as HTMLTextAreaElement).value,
      });
  }
  setMetric(event: Event): void {
    const metric = (event.target as HTMLSelectElement).value as MetricId;
    if (!METRICS.some((item) => item.id === metric)) return;
    const unit = metricInfo(metric).unit;
    const update: Partial<GaugeWidgetProperties> = {
      metric,
      format: unit === 'ms' ? 'duration' : unit === 'count' ? 'number' : 'percent',
      precision: unit === 'count' ? 0 : 1,
    };
    if (this.gauge() && unit !== this.metric().unit) {
      const maximum = unit === 'ms' ? 144000000 : 100;
      Object.assign(update, { minimum: 0, maximum, target: maximum, thresholds: [] });
      this.unitResetNotice.set('Scale, goal, and thresholds reset for the new unit.');
    } else {
      this.unitResetNotice.set('');
    }
    this.update(update);
  }
  setPeriod(event: Event): void {
    const period = (event.target as HTMLSelectElement).value as MetricWidgetProperties['period'];
    const p = this.widget().properties;
    const today = this.today();
    this.update({ period, startDate: p.startDate || today, endDate: p.endDate || today });
  }
  setDate(key: 'startDate' | 'endDate', event: Event): void {
    const input = event.target as HTMLInputElement,
      value = input.value;
    const p = this.widget().properties;
    const update = {
      [key]: value,
      ...(key === 'startDate' && value > p.endDate ? { endDate: value } : {}),
      ...(key === 'endDate' && value < p.startDate ? { startDate: value } : {}),
    };
    if (!this.update(update)) input.value = p[key];
  }
  setArchived(event: Event): void {
    this.update({
      archived: (event.target as HTMLSelectElement).value as MetricWidgetProperties['archived'],
    });
  }
  setFormat(event: Event): void {
    this.update({
      format: (event.target as HTMLSelectElement).value as MetricWidgetProperties['format'],
    });
  }
  setPrecision(event: Event): void {
    this.update({ precision: Number((event.target as HTMLSelectElement).value) });
  }
  setBoolean(
    key: 'includeRunning' | 'showValue' | 'showTarget' | 'showPercentage',
    event: Event,
  ): void {
    this.update({ [key]: (event.target as HTMLInputElement).checked });
  }
  setGoalMode(event: Event): void {
    const goalMode = (event.target as HTMLSelectElement).value as GaugeWidgetProperties['goalMode'];
    const p = this.gauge();
    if (p)
      this.update({
        goalMode,
        target: goalMode !== 'informational' && p.target === null ? p.maximum : p.target,
      });
  }
  setGaugeNumber(key: 'minimum' | 'maximum' | 'target', event: Event): void {
    const p = this.gauge();
    if (!p) return;
    const input = event.target as HTMLInputElement;
    const value =
      input.value === '' && key === 'target'
        ? null
        : input.valueAsNumber * (this.metric().unit === 'ms' ? 3600000 : 1);
    if (typeof value === 'number' && !Number.isFinite(value)) {
      this.error.set('Enter a valid number.');
      input.setAttribute('aria-invalid', 'true');
      input.value = String(this.displayNumber(p[key]) ?? '');
      return;
    }
    const update: Partial<GaugeWidgetProperties> = { [key]: value };
    if (!this.update(update)) input.value = String(this.displayNumber(p[key]) ?? '');
  }
  setThresholdValue(index: number, event: Event): void {
    const p = this.gauge();
    if (!p) return;
    const input = event.target as HTMLInputElement;
    const value = input.valueAsNumber * (this.metric().unit === 'ms' ? 3600000 : 1);
    if (
      !this.update({ thresholds: p.thresholds.map((t, i) => (i === index ? { ...t, value } : t)) })
    )
      input.value = String(this.displayNumber(p.thresholds[index].value));
  }
  setThresholdTone(index: number, event: Event): void {
    const p = this.gauge();
    if (!p) return;
    const tone = (event.target as HTMLSelectElement).value as MetricTone;
    this.update({ thresholds: p.thresholds.map((t, i) => (i === index ? { ...t, tone } : t)) });
  }
  addThreshold(): void {
    const p = this.gauge();
    if (!p || !this.canAddThreshold()) return;
    const value = p.thresholds.length
      ? Math.min(
          p.maximum,
          (p.thresholds.at(-1)?.value ?? p.minimum) +
            Math.max(1, Math.round((p.maximum - p.minimum) / 4)),
        )
      : p.minimum;
    this.update({ thresholds: [...p.thresholds, { value, tone: 'info' }] });
  }
  removeThreshold(index: number): void {
    const p = this.gauge();
    if (p) this.update({ thresholds: p.thresholds.filter((_, i) => i !== index) });
  }
  update(patch: Partial<GaugeWidgetProperties>): boolean {
    if (this.disabled()) return false;
    const widget = this.widget();
    const properties = { ...widget.properties, ...patch };
    const error = metricSettingsError(properties, widget.type === 'gauge');
    this.error.set(error);
    if (error) return false;
    if (widget.type === 'gauge')
      this.widgetChange.emit({ ...widget, properties: { ...widget.properties, ...patch } });
    else this.widgetChange.emit({ ...widget, properties });
    return true;
  }
}

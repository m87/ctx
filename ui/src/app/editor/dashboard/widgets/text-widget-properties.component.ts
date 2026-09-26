import { Component, input, output, signal } from '@angular/core';
import { HlmButtonImports } from '@spartan-ng/helm/button';
import { HlmInputImports } from '@spartan-ng/helm/input';
import { HlmLabelImports } from '@spartan-ng/helm/label';
import { HlmTextareaImports } from '@spartan-ng/helm/textarea';
import {
  TextWidgetDefinition,
  TextWidgetProperties,
  TEXT_WIDGET_FONT_SIZE,
} from '../dashboard-definition';

@Component({
  selector: 'ctx-text-widget-properties',
  imports: [HlmButtonImports, HlmInputImports, HlmLabelImports, HlmTextareaImports],
  host: { class: 'block min-h-0 overflow-y-auto' },
  template: `
    <div class="flex flex-col gap-4 p-3">
      <div class="ui-field">
        <label hlmLabel [for]="fieldId('text')">Text</label>
        <textarea
          hlmTextarea
          class="dashboard-field-control min-h-28 resize-y"
          [id]="fieldId('text')"
          [value]="widget().properties.text"
          [disabled]="disabled()"
          maxlength="20000"
          placeholder="Write a note or heading…"
          (input)="updateText($event)"
        ></textarea>
      </div>
      <section
        class="flex flex-col gap-3 border-t border-border/65 pt-4"
        aria-label="Text appearance"
      >
        <h4 class="ui-section-label">Appearance</h4>
        <div class="grid grid-cols-2 gap-3">
          <div class="ui-field">
            <label hlmLabel [for]="fieldId('horizontal')">Horizontal</label>
            <select
              class="ui-select-trigger w-full"
              [id]="fieldId('horizontal')"
              [value]="widget().properties.horizontalAlign ?? 'left'"
              [disabled]="disabled()"
              (change)="updateHorizontal($event)"
            >
              <option value="left">Left</option>
              <option value="center">Center</option>
              <option value="right">Right</option>
            </select>
          </div>
          <div class="ui-field">
            <label hlmLabel [for]="fieldId('vertical')">Vertical</label>
            <select
              class="ui-select-trigger w-full"
              [id]="fieldId('vertical')"
              [value]="widget().properties.verticalAlign ?? 'top'"
              [disabled]="disabled()"
              (change)="updateVertical($event)"
            >
              <option value="top">Top</option>
              <option value="center">Center</option>
              <option value="bottom">Bottom</option>
            </select>
          </div>
        </div>
        <div class="ui-field">
          <label hlmLabel [for]="fieldId('font-size')">Font size (px)</label>
          <input
            hlmInput
            class="dashboard-field-control h-9"
            type="number"
            step="1"
            [id]="fieldId('font-size')"
            [min]="fontSize.minimum"
            [max]="fontSize.maximum"
            [value]="widget().properties.fontSize ?? fontSize.default"
            [disabled]="disabled()"
            [attr.aria-invalid]="fontSizeError() ? true : null"
            [attr.aria-describedby]="fontSizeError() ? fieldId('font-error') : null"
            (input)="updateFontSize($event)"
            (blur)="restoreFontSize($event)"
          />
          @if (fontSizeError()) {
            <p [id]="fieldId('font-error')" class="text-xs text-destructive" role="alert">
              Enter a whole number from {{ fontSize.minimum }} to {{ fontSize.maximum }}.
            </p>
          }
        </div>
      </section>
      <details class="border-t border-border/65 pt-3">
        <summary
          class="cursor-pointer rounded-md text-xs text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        >
          Widget query
        </summary>
        <div class="ui-field mt-3">
          <label hlmLabel [for]="fieldId('query')">Query</label>
          <textarea
            hlmTextarea
            class="dashboard-field-control min-h-24 resize-y"
            [id]="fieldId('query')"
            [value]="widget().query"
            [disabled]="disabled()"
            maxlength="8000"
            (input)="updateQuery($event)"
          ></textarea>
        </div>
      </details>
      <button
        hlmBtn
        type="button"
        variant="ghost"
        class="justify-start text-destructive"
        [disabled]="disabled()"
        (click)="deleteRequested.emit()"
      >
        Delete widget
      </button>
    </div>
  `,
})
export class TextWidgetPropertiesComponent {
  readonly widget = input.required<TextWidgetDefinition>();
  readonly disabled = input(false);
  readonly idSuffix = input('');
  readonly queryChange = output<string>();
  readonly propertiesChange = output<TextWidgetProperties>();
  readonly deleteRequested = output<void>();
  readonly fontSize = TEXT_WIDGET_FONT_SIZE;
  readonly fontSizeError = signal(false);

  fieldId(field: string): string {
    return `widget-${field}-${this.widget().id}${this.idSuffix()}`;
  }

  updateQuery(event: Event): void {
    if (!this.disabled()) this.queryChange.emit((event.target as HTMLTextAreaElement).value);
  }

  updateText(event: Event): void {
    this.update({ text: (event.target as HTMLTextAreaElement).value });
  }

  updateHorizontal(event: Event): void {
    const value = (event.target as HTMLSelectElement).value;
    if (value === 'left' || value === 'center' || value === 'right') {
      this.update({ horizontalAlign: value });
    }
  }

  updateVertical(event: Event): void {
    const value = (event.target as HTMLSelectElement).value;
    if (value === 'top' || value === 'center' || value === 'bottom') {
      this.update({ verticalAlign: value });
    }
  }

  updateFontSize(event: Event): void {
    const value = (event.target as HTMLInputElement).valueAsNumber;
    const valid =
      Number.isInteger(value) && value >= this.fontSize.minimum && value <= this.fontSize.maximum;
    this.fontSizeError.set(!valid);
    if (valid) this.update({ fontSize: value });
  }

  restoreFontSize(event: Event): void {
    if (this.fontSizeError()) {
      (event.target as HTMLInputElement).value = String(
        this.widget().properties.fontSize ?? this.fontSize.default,
      );
      this.fontSizeError.set(false);
    }
  }

  private update(properties: Partial<TextWidgetProperties>): void {
    if (!this.disabled())
      this.propertiesChange.emit({ ...this.widget().properties, ...properties });
  }
}

import type { BooleanInput, NumberInput } from '@angular/cdk/coercion';
import {
  booleanAttribute,
  ChangeDetectionStrategy,
  Component,
  computed,
  contentChild,
  forwardRef,
  input,
  linkedSignal,
  numberAttribute,
  output,
  signal,
  viewChild,
} from '@angular/core';
import { type ControlValueAccessor, NG_VALUE_ACCESSOR } from '@angular/forms';
import type { BrnDialogState } from '@spartan-ng/brain/dialog';
import type { Weekday } from '@spartan-ng/brain/calendar';
import { BrnFieldControl, provideBrnLabelable } from '@spartan-ng/brain/field';
import type { ChangeFn, TouchFn } from '@spartan-ng/brain/forms';
import { BrnPopover } from '@spartan-ng/brain/popover';
import { HlmCalendar } from '@spartan-ng/helm/calendar';
import { HlmPopoverImports } from '@spartan-ng/helm/popover';
import { HlmDatePickerTriggerToken } from './hlm-date-picker-trigger.token';
import {
  HlmDatePickerBase,
  injectHlmDatePickerConfig,
  provideHlmDatePicker,
} from './hlm-date-picker.token';

export const HLM_DATE_PICKER_VALUE_ACCESSOR = {
  provide: NG_VALUE_ACCESSOR,
  useExisting: forwardRef(() => HlmDatePicker),
  multi: true,
};

@Component({
  selector: 'hlm-date-picker',
  imports: [HlmPopoverImports, HlmCalendar],
  providers: [
    HLM_DATE_PICKER_VALUE_ACCESSOR,
    provideHlmDatePicker(HlmDatePicker),
    provideBrnLabelable(HlmDatePicker),
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
  hostDirectives: [BrnFieldControl],
  host: { class: 'block' },
  template: `
    <hlm-popover
      [align]="align()"
      sideOffset="5"
      autoFocus="first-heading"
      [state]="_popoverState()"
      (stateChanged)="_popoverState.set($event)"
      (closed)="_onTouched?.()"
    >
      <ng-content />

      <hlm-popover-content class="w-fit p-0" *hlmPopoverPortal="let ctx">
        <ng-content select="[hlmDatePickerHeader]" />
        <hlm-calendar
          calendarClass="rounded-none border-0"
          [captionLayout]="captionLayout()"
          [date]="_mutableDate()"
          [defaultFocusedDate]="_mutableDate() ?? defaultFocusedDate()"
          [min]="min()"
          [max]="max()"
          [disabled]="_disabled()"
          [weekStartsOn]="weekStartsOn()"
          (dateChange)="_handleChange($event)"
        />
        <ng-content select="[hlmDatePickerFooter]" />
      </hlm-popover-content>
    </hlm-popover>
  `,
})
export class HlmDatePicker<T> implements HlmDatePickerBase<T>, ControlValueAccessor {
  private readonly _config = injectHlmDatePickerConfig<T>();

  public readonly popover = viewChild.required(BrnPopover);

  private readonly _trigger = contentChild(HlmDatePickerTriggerToken);

  public readonly captionLayout = input<
    'dropdown' | 'label' | 'dropdown-months' | 'dropdown-years'
  >('label');

  public readonly min = input<T>();

  public readonly align = input<'start' | 'center' | 'end'>('start');

  public readonly weekStartsOn = input<Weekday, NumberInput>(undefined, {
    transform: (value: unknown) => numberAttribute(value) as Weekday,
  });

  public readonly max = input<T>();

  public readonly disabled = input<boolean, BooleanInput>(false, {
    transform: booleanAttribute,
  });

  public readonly date = input<T>();

  public readonly defaultFocusedDate = input<T>();

  protected readonly _mutableDate = linkedSignal(this.date);

  public readonly autoCloseOnSelect = input<boolean, BooleanInput>(this._config.autoCloseOnSelect, {
    transform: booleanAttribute,
  });

  public readonly formatDate = input<(date: T) => string>(this._config.formatDate);

  public readonly transformDate = input<(date: T) => T>(this._config.transformDate);

  protected readonly _popoverState = signal<BrnDialogState | null>(null);

  protected readonly _disabled = linkedSignal(this.disabled);

  public readonly disabledState = this._disabled.asReadonly();

  public readonly formattedDate = computed(() => {
    const date = this._mutableDate();
    return date ? this.formatDate()(date) : undefined;
  });

  public readonly dateChange = output<T>();

  public readonly labelableId = computed(() => this._trigger()?.triggerId());

  public readonly hasDate = computed(() => !!this._mutableDate());

  protected _onChange?: ChangeFn<T>;
  protected _onTouched?: TouchFn;

  protected _handleChange(value: T | undefined) {
    if (this._disabled()) return;
    this.updateDate(value);

    if (this.autoCloseOnSelect()) {
      this._popoverState.set('closed');
    }
  }

  public updateDate(value: T | undefined) {
    if (this._disabled()) return;
    const transformedDate = value !== undefined ? this.transformDate()(value) : undefined;

    this._mutableDate.set(transformedDate);
    this._onChange?.(transformedDate as T);
    this.dateChange.emit(transformedDate as T);
  }

  public writeValue(value: T | null): void {
    this._mutableDate.set(value ? this.transformDate()(value) : undefined);
  }

  public registerOnChange(fn: ChangeFn<T>): void {
    this._onChange = fn;
  }

  public registerOnTouched(fn: TouchFn): void {
    this._onTouched = fn;
  }

  public touched(): void {
    this._onTouched?.();
  }

  public setDisabledState(isDisabled: boolean): void {
    this._disabled.set(isDisabled);
  }

  public open() {
    this._popoverState.set('open');
  }

  public close() {
    this._popoverState.set('closed');
  }

  public reset() {
    this._mutableDate.set(undefined);
    this._onChange?.(undefined as T);
    this.dateChange.emit(undefined as T);
  }
}

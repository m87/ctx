import { BooleanInput } from '@angular/cdk/coercion';
import {
  booleanAttribute,
  ChangeDetectionStrategy,
  Component,
  computed,
  inject,
  input,
} from '@angular/core';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideChevronDown } from '@ng-icons/lucide';
import { BrnFieldControl, BrnFieldControlDescribedBy } from '@spartan-ng/brain/field';
import { buttonVariants, type ButtonVariants } from '@spartan-ng/helm/button';
import { HlmPopoverTrigger } from '@spartan-ng/helm/popover';
import { hlm } from '@spartan-ng/helm/utils';
import { ClassValue } from 'clsx';
import {
  HlmDatePickerTriggerBase,
  provideHlmDatePickerTrigger,
} from './hlm-date-picker-trigger.token';
import { injectHlmDatePicker } from './hlm-date-picker.token';

@Component({
  selector: 'hlm-date-picker-trigger',
  imports: [HlmPopoverTrigger, NgIcon, BrnFieldControlDescribedBy],
  providers: [
    provideIcons({ lucideChevronDown }),
    provideHlmDatePickerTrigger(HlmDatePickerTrigger),
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
  host: { 'data-slot': 'date-picker-trigger' },
  template: `
    <button
      [id]="buttonId()"
      type="button"
      [class]="_computedClass()"
      [disabled]="_disabled()"
      [attr.aria-invalid]="_ariaInvalid()"
      [attr.aria-label]="ariaLabel()"
      [attr.data-invalid]="_ariaInvalid()"
      [attr.data-touched]="_touched?.() ? 'true' : null"
      [attr.data-dirty]="_dirty?.() ? 'true' : null"
      [attr.data-matches-spartan-invalid]="_spartanInvalid() ? 'true' : null"
      hlmPopoverTrigger
      [hlmPopoverTriggerFor]="_popover()"
      brnFieldControlDescribedBy
      [attr.data-placeholder]="_isPlaceholder() ? '' : null"
    >
      <span class="truncate">
        @if (_formattedDate(); as formattedDate) {
          {{ formattedDate }}
        } @else {
          <ng-content />
        }
      </span>

      @if (showTrigger()) {
        <ng-icon name="lucideChevronDown" />
      }
    </button>
  `,
})
export class HlmDatePickerTrigger implements HlmDatePickerTriggerBase {
  private static _nextId = 0;

  private readonly _fieldControl = inject(BrnFieldControl, { optional: true });
  private readonly _datePicker = injectHlmDatePicker();

  private readonly _invalid = this._fieldControl?.invalid;
  protected readonly _spartanInvalid = computed(
    () => this.forceInvalid() || this._fieldControl?.spartanInvalid(),
  );
  protected readonly _dirty = this._fieldControl?.dirty;
  protected readonly _touched = this._fieldControl?.touched;

  protected readonly _ariaInvalid = computed(() => (this._invalid?.() ? 'true' : null));

  public readonly userClass = input<ClassValue>('', { alias: 'class' });
  protected readonly _computedClass = computed(() =>
    hlm(
      buttonVariants({ variant: this.variant(), size: 'default' }),
      'spartan-date-picker-trigger ring-offset-background border-input w-[280px] cursor-default justify-between px-3 text-left font-normal shadow-none',
      'focus-visible:border-input focus-visible:ring-ring focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:outline-none',
      this.userClass(),
    ),
  );

  protected readonly _isPlaceholder = computed(() => !this._datePicker.hasDate());

  public readonly buttonId = input<string>(`hlm-date-picker-${++HlmDatePickerTrigger._nextId}`);

  public readonly ariaLabel = input<string | null>(null, { alias: 'aria-label' });

  public readonly triggerId = this.buttonId;

  public readonly forceInvalid = input<boolean, BooleanInput>(false, {
    transform: booleanAttribute,
  });

  public readonly variant = input<ButtonVariants['variant']>('outline');

  public readonly showTrigger = input<boolean, BooleanInput>(true, { transform: booleanAttribute });

  protected readonly _popover = this._datePicker.popover;
  protected readonly _disabled = this._datePicker.disabledState;
  protected readonly _formattedDate = this._datePicker.formattedDate;
}

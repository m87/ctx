import {
  type ExistingProvider,
  inject,
  InjectionToken,
  type Signal,
  type Type,
  type ValueProvider,
} from '@angular/core';
import type { BrnPopover } from '@spartan-ng/brain/popover';

export interface HlmDatePickerBase<T> {
  popover: Signal<BrnPopover>;
  disabledState: Signal<boolean>;
  formattedDate: Signal<string | undefined>;
  hasDate: Signal<boolean>;
  updateDate?(value: T | undefined): void;
  touched?(): void;
}

export const HlmDatePickerToken = new InjectionToken<HlmDatePickerBase<unknown>>(
  'HlmDatePickerToken',
);

export function provideHlmDatePicker(instance: Type<HlmDatePickerBase<unknown>>): ExistingProvider {
  return { provide: HlmDatePickerToken, useExisting: instance };
}

export function injectHlmDatePicker<T>(): HlmDatePickerBase<T> {
  return inject(HlmDatePickerToken) as HlmDatePickerBase<T>;
}

export interface HlmDatePickerConfig<T> {
  autoCloseOnSelect: boolean;

  formatDate: (date: T) => string;

  transformDate: (date: T) => T;

  parseDate: (value: string) => T | undefined;
}

function getDefaultConfig<T>(): HlmDatePickerConfig<T> {
  return {
    formatDate: (date) => (date instanceof Date ? date.toDateString() : `${date}`),
    transformDate: (date) => date,
    parseDate: (value) => {
      const date = new Date(value);
      return isNaN(date.getTime()) ? undefined : (date as T);
    },
    autoCloseOnSelect: false,
  };
}

const HlmDatePickerConfigToken = new InjectionToken<HlmDatePickerConfig<unknown>>(
  'HlmDatePickerConfig',
);

export function provideHlmDatePickerConfig<T>(
  config: Partial<HlmDatePickerConfig<T>>,
): ValueProvider {
  return { provide: HlmDatePickerConfigToken, useValue: { ...getDefaultConfig(), ...config } };
}

export function injectHlmDatePickerConfig<T>(): HlmDatePickerConfig<T> {
  const injectedConfig = inject(HlmDatePickerConfigToken, { optional: true });
  return injectedConfig ? (injectedConfig as HlmDatePickerConfig<T>) : getDefaultConfig();
}

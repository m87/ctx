import type { BooleanInput } from '@angular/cdk/coercion';
import { NgTemplateOutlet } from '@angular/common';
import {
	booleanAttribute,
	ChangeDetectionStrategy,
	Component,
	computed,
	effect,
	ElementRef,
	forwardRef,
	inject,
	input,
	linkedSignal,
	model,
	output,
	type TemplateRef,
	viewChild,
} from '@angular/core';
import { type ControlValueAccessor, NG_VALUE_ACCESSOR } from '@angular/forms';
import { provideIcons } from '@ng-icons/core';
import { lucideChevronDown, lucideCircleX, lucideSearch } from '@ng-icons/lucide';
import { BrnAutocomplete, BrnAutocompleteImports } from '@spartan-ng/brain/autocomplete';
import { debouncedSignal } from '@spartan-ng/brain/core';
import type { ChangeFn, TouchFn } from '@spartan-ng/brain/forms';
import { BrnPopoverImports } from '@spartan-ng/brain/popover';
import { HlmIconImports } from '@spartan-ng/helm/icon';
import { HlmPopoverImports } from '@spartan-ng/helm/popover';
import { classes, hlm } from '@spartan-ng/helm/utils';
import type { ClassValue } from 'clsx';
import { HlmAutocompleteEmpty } from './hlm-autocomplete-empty';
import { HlmAutocompleteGroup } from './hlm-autocomplete-group';
import { HlmAutocompleteItem } from './hlm-autocomplete-item';
import { HlmAutocompleteList } from './hlm-autocomplete-list';

import { HlmInputGroupImports } from '@spartan-ng/helm/input-group';
import { HlmAutocompleteTrigger } from './hlm-autocomplete-trigger';
import { injectHlmAutocompleteConfig } from './hlm-autocomplete.token';

export const HLM_AUTOCOMPLETE_VALUE_ACCESSOR = {
	provide: NG_VALUE_ACCESSOR,
	useExisting: forwardRef(() => HlmAutocomplete),
	multi: true,
};

@Component({
	selector: 'hlm-autocomplete',
	imports: [
		NgTemplateOutlet,
		BrnAutocompleteImports,
		HlmAutocompleteEmpty,
		HlmAutocompleteGroup,
		HlmAutocompleteItem,
		HlmAutocompleteList,
		HlmAutocompleteTrigger,
		BrnPopoverImports,
		HlmPopoverImports,
		HlmIconImports,
		HlmInputGroupImports,
	],
	providers: [HLM_AUTOCOMPLETE_VALUE_ACCESSOR, provideIcons({ lucideSearch, lucideChevronDown, lucideCircleX })],
	changeDetection: ChangeDetectionStrategy.OnPush,
	template: `
		@let transformer = transformOptionToValue();
		<hlm-popover
			#popover
			align="start"
			autoFocus="first-heading"
			sideOffset="5"
			closeDelay="100"
			closeOnOutsidePointerEvents="true"
			(closed)="_closed()"
		>
			<div brnAutocomplete (selectionCleared)="_selectionCleared()">
				<div
					hlmInputGroup
					hlmAutocompleteTrigger
					[class]="_computedAutocompleteSearchClass()"
					[disabledTrigger]="!_search() || _disabled()"
				>
					<input
						#input
						brnAutocompleteSearchInput
						hlmInputGroupInput
						type="text"
						autocomplete="off"
						[id]="inputId()"
						[class]="_computedAutocompleteInputClass()"
						[placeholder]="searchPlaceholderText()"
						[disabled]="_disabled()"
						[value]="search()"
						(input)="_searchChanged($event)"
					/>
					<div hlmInputGroupAddon>
						<ng-icon name="lucideSearch" />
					</div>
					<div hlmInputGroupAddon align="inline-end">
						@if (showClearBtn() && value() !== undefined) {
							<button
								hlmInputGroupButton
								type="button"
								tabindex="-1"
								aria-label="clear"
								[disabled]="_disabled()"
								(click)="_selectionCleared()"
								size="icon-xs"
							>
								<ng-icon name="lucideCircleX" />
							</button>
						}

						<button
							hlmInputGroupButton
							type="button"
							tabindex="-1"
							[attr.aria-label]="ariaLabelToggleButton()"
							[disabled]="_disabled()"
							(click)="_toggleOptions()"
							size="icon-xs"
						>
							<ng-icon name="lucideChevronDown" />
						</button>
					</div>
				</div>

				<div
					*brnPopoverContent="let ctx"
					hlmPopoverContent
					class="p-0"
					[style.width.px]="_elementRef.nativeElement.offsetWidth"
				>
					<hlm-autocomplete-list
						[class]="_computedAutocompleteListClass()"
						[class.hidden]="filteredOptions().length === 0"
					>
						<hlm-autocomplete-group>
							@for (option of filteredOptions(); track option) {
								<button
									hlm-autocomplete-item
									[class]="_computedAutocompleteItemClass()"
									[value]="transformer ? transformer(option) : option"
									(selected)="_optionSelected(option)"
								>
									@if (optionTemplate(); as optionTemplate) {
										<ng-container *ngTemplateOutlet="optionTemplate; context: { $implicit: option }" />
									} @else {
										{{ transformOptionToString()(option) }}
									}
								</button>
							}
						</hlm-autocomplete-group>
					</hlm-autocomplete-list>

					<div *brnAutocompleteEmpty hlmAutocompleteEmpty [class]="_computedAutocompleteEmptyClass()">
						@if (loading()) {
							<ng-content select="[loading]">{{ loadingText() }}</ng-content>
						} @else {
							<ng-content select="[empty]">{{ emptyText() }}</ng-content>
						}
					</div>
				</div>
			</div>
		</hlm-popover>
	`,
})
export class HlmAutocomplete<T, V = T> implements ControlValueAccessor {
	private static _id = 0;
	private readonly _config = injectHlmAutocompleteConfig<T, V>();

	private readonly _brnAutocomplete = viewChild.required(BrnAutocomplete);

	private readonly _inputRef = viewChild.required('input', { read: ElementRef });

	protected readonly _elementRef = inject(ElementRef<HTMLElement>);

	public readonly autocompleteSearchClass = input<ClassValue>('');
	protected readonly _computedAutocompleteSearchClass = computed(() => hlm(this.autocompleteSearchClass()));

	public readonly autocompleteInputClass = input<ClassValue>('');
	protected readonly _computedAutocompleteInputClass = computed(() => hlm(this.autocompleteInputClass()));

	public readonly autocompleteListClass = input<ClassValue>('');
	protected readonly _computedAutocompleteListClass = computed(() => hlm(this.autocompleteListClass()));

	public readonly autocompleteItemClass = input<ClassValue>('');
	protected readonly _computedAutocompleteItemClass = computed(() => hlm(this.autocompleteItemClass()));

	public readonly autocompleteEmptyClass = input<ClassValue>('');
	protected readonly _computedAutocompleteEmptyClass = computed(() => hlm(this.autocompleteEmptyClass()));

	public readonly filteredOptions = input<T[]>([]);

	public readonly value = model<T | V>();

	public readonly debounceTime = input<number>(this._config.debounceTime);

	public readonly search = model<string>('');

	protected readonly _search = debouncedSignal(this.search, this.debounceTime());

	public readonly transformValueToSearch = input<(option: T) => string>(this._config.transformValueToSearch);

	public readonly requireSelection = input<boolean, BooleanInput>(this._config.requireSelection, {
		transform: booleanAttribute,
	});

	public readonly transformOptionToString = input<(option: T) => string>(this._config.transformOptionToString);

	public readonly transformOptionToValue = input<((option: T) => V) | undefined>(this._config.transformOptionToValue);

	public readonly displayWith = input<((value: V) => string) | undefined>(undefined);

	protected readonly _displaySearchValue = computed(() => {
		const displayWith = this.displayWith();
		if (displayWith) {
			return displayWith;
		} else {
			return this.transformValueToSearch();
		}
	});

	public readonly optionTemplate = input<TemplateRef<HlmAutocompleteOption<T>>>();

	public readonly loading = input<boolean, BooleanInput>(false, { transform: booleanAttribute });

	public readonly showClearBtn = input<boolean, BooleanInput>(this._config.showClearBtn, {
		transform: booleanAttribute,
	});

	public readonly searchPlaceholderText = input<string>('Select an option');

	public readonly loadingText = input<string>('Loading options...');

	public readonly emptyText = input<string>('No options found');

	public readonly ariaLabelToggleButton = input<string>('Toggle options');

	public readonly inputId = input<string>(`hlm-autocomplete-input-${++HlmAutocomplete._id}`);

	public readonly disabled = input<boolean, BooleanInput>(false, { transform: booleanAttribute });

	protected readonly _disabled = linkedSignal(() => this.disabled());

	public readonly valueChange = output<T | V | null>();

	public readonly searchChange = output<string>();

	protected _onChange?: ChangeFn<T | V | null>;
	protected _onTouched?: TouchFn;

	constructor() {
		classes(() => 'block w-full');
		effect(() => {
			const search = this._search();
			this.searchChange.emit(search);
		});
	}

	protected _searchChanged(event: Event) {
		const value = (event.target as HTMLInputElement).value;
		this.search.set(value ?? '');

		if (!this._brnAutocomplete().isExpanded() && value.length > 0) {
			this._brnAutocomplete().open();
		}
	}

	protected _toggleOptions() {
		if (this._search() || this.filteredOptions().length > 0) {
			this._brnAutocomplete().toggle();
		}

		this._inputRef().nativeElement.focus();
	}

	protected _selectionCleared() {
		this.value.set(undefined);
		this._onChange?.(null);
		this.valueChange.emit(null);
		this.search.set('');
	}

	protected _optionSelected(option: T) {
		const transformer = this.transformOptionToValue();

		const value = transformer ? transformer(option) : option;

		this.value.set(value);
		this._onChange?.(value);
		this.valueChange.emit(value);

		const searchValue = this._displaySearchValue()(value as any);
		this.search.set(searchValue ?? '');
		this._brnAutocomplete().close();
	}

	public writeValue(value: T | V | null): void {
		this.value.set(value ? value : undefined);

		const searchValue = value ? this._displaySearchValue()(value as any) : '';
		this.search.set(searchValue);
	}

	public registerOnChange(fn: ChangeFn<T | V | null>): void {
		this._onChange = fn;
	}

	public registerOnTouched(fn: TouchFn): void {
		this._onTouched = fn;
	}

	public setDisabledState(isDisabled: boolean): void {
		this._disabled.set(isDisabled);
	}

	protected _closed() {
		if (this.requireSelection()) {
			const value = this.value();
			const searchValue = value ? this._displaySearchValue()(value as any) : '';
			this.search.set(searchValue ?? '');
		}
	}
}

export interface HlmAutocompleteOption<T> {
	$implicit: T;
}

import { Component, computed, input, output, signal } from '@angular/core';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideCheck, lucideSearch } from '@ng-icons/lucide';

export interface SearchSelectOption {
  value: string;
  label: string;
  color?: string;
  description?: string;
  keywords?: readonly string[];
}

@Component({
  selector: 'ctx-search-select',
  imports: [NgIcon],
  providers: [provideIcons({ lucideCheck, lucideSearch })],
  host: {
    class: 'block min-w-0',
  },
  template: `
    <div
      class="overflow-hidden rounded-xl border border-input bg-background shadow-xs transition-[border-color,box-shadow] focus-within:border-ring focus-within:ring-[3px] focus-within:ring-ring/50"
      [class.pointer-events-none]="disabled()"
      [class.opacity-60]="disabled()"
      [attr.aria-busy]="disabled()"
    >
      <div class="flex h-11 items-center gap-2.5 border-b border-border/70 px-3">
        <ng-icon name="lucideSearch" class="shrink-0 text-sm text-muted-foreground"></ng-icon>
        <input
          [id]="inputId()"
          type="search"
          role="combobox"
          aria-autocomplete="list"
          aria-expanded="true"
          autocomplete="off"
          class="min-w-0 flex-1 bg-transparent text-sm outline-none placeholder:text-muted-foreground"
          [attr.aria-label]="ariaLabel()"
          [attr.aria-controls]="listboxId()"
          [attr.aria-activedescendant]="activeOptionId()"
          [placeholder]="searchPlaceholder()"
          [disabled]="disabled()"
          [value]="search()"
          (input)="updateSearch($event)"
          (keydown)="handleSearchKeydown($event)"
        />
      </div>

      <div
        [id]="listboxId()"
        class="max-h-56 overflow-y-auto p-1.5"
        role="listbox"
        [attr.aria-label]="ariaLabel()"
      >
        @for (option of filteredOptions(); track option.value; let index = $index) {
          <button
            [id]="optionId(index)"
            type="button"
            class="flex w-full items-center gap-3 rounded-lg px-2.5 py-2 text-left transition-colors hover:bg-muted/70 focus-visible:bg-muted focus-visible:outline-none"
            role="option"
            [disabled]="disabled()"
            [attr.aria-selected]="value() === option.value"
            [class.bg-muted]="value() === option.value || activeIndex() === index"
            (click)="selectOption(option)"
            (mouseenter)="activeIndex.set(index)"
          >
            <span
              class="size-2.5 shrink-0 rounded-full"
              [style.background-color]="option.color ?? 'var(--muted-foreground)'"
              aria-hidden="true"
            ></span>
            <span class="min-w-0 flex-1">
              <span class="block truncate text-sm font-medium text-foreground">
                {{ option.label }}
              </span>
              @if (option.description) {
                <span class="block truncate text-[11px] text-muted-foreground">
                  {{ option.description }}
                </span>
              }
            </span>
            @if (value() === option.value) {
              <ng-icon name="lucideCheck" class="shrink-0 text-sm text-foreground"></ng-icon>
            }
          </button>
        } @empty {
          <div class="px-4 py-8 text-center text-xs text-muted-foreground">
            {{ emptyText() }}
          </div>
        }
      </div>
    </div>
  `,
})
export class SearchSelectComponent {
  readonly inputId = input.required<string>();
  readonly options = input.required<readonly SearchSelectOption[]>();
  readonly value = input('');
  readonly ariaLabel = input('Options');
  readonly searchPlaceholder = input('Search…');
  readonly emptyText = input('No matching options');
  readonly disabled = input(false);
  readonly selectionChange = output<string>();

  readonly search = signal('');
  readonly activeIndex = signal(-1);
  readonly listboxId = computed(() => `${this.inputId()}-listbox`);
  readonly filteredOptions = computed(() => {
    const search = this.search().trim().toLocaleLowerCase();
    if (!search) {
      return this.options();
    }

    return this.options().filter((option) =>
      [option.label, option.description, ...(option.keywords ?? [])]
        .filter((value): value is string => Boolean(value))
        .some((value) => value.toLocaleLowerCase().includes(search)),
    );
  });
  readonly activeOptionId = computed(() => {
    const index = this.activeIndex();
    return index >= 0 && index < this.filteredOptions().length ? this.optionId(index) : null;
  });

  updateSearch(event: Event): void {
    this.search.set((event.target as HTMLInputElement).value);
    this.activeIndex.set(this.filteredOptions().length > 0 ? 0 : -1);
  }

  handleSearchKeydown(event: KeyboardEvent): void {
    const options = this.filteredOptions();
    if (options.length === 0) {
      return;
    }

    if (event.key === 'ArrowDown' || event.key === 'ArrowUp') {
      event.preventDefault();
      const direction = event.key === 'ArrowDown' ? 1 : -1;
      const current = this.activeIndex();
      const next = current < 0 ? (direction > 0 ? 0 : options.length - 1) : current + direction;
      this.activeIndex.set((next + options.length) % options.length);
      return;
    }

    if (event.key === 'Enter' && this.activeIndex() >= 0) {
      event.preventDefault();
      this.selectOption(options[this.activeIndex()]);
    }
  }

  selectOption(option: SearchSelectOption): void {
    if (!this.disabled()) {
      this.selectionChange.emit(option.value);
    }
  }

  optionId(index: number): string {
    return `${this.inputId()}-option-${index}`;
  }
}

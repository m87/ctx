import {
  Component,
  ElementRef,
  HostListener,
  computed,
  effect,
  inject,
  input,
  output,
  signal,
  viewChild,
} from '@angular/core';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideCheck, lucideChevronDown, lucidePlus, lucideX } from '@ng-icons/lucide';
import { SearchSelectComponent, SearchSelectOption } from './search-select.component';

@Component({
  selector: 'ctx-search-dropdown-select',
  imports: [NgIcon, SearchSelectComponent],
  providers: [provideIcons({ lucideCheck, lucideChevronDown, lucidePlus, lucideX })],
  host: {
    class: 'block min-w-0',
  },
  template: `
    <div class="relative">
      <button
        type="button"
        class="flex h-8 w-full items-center justify-between gap-2 rounded-md border border-border/60 bg-muted/30 px-2.5 text-left text-xs outline-none transition-[background-color,border-color,box-shadow] hover:border-border hover:bg-muted/50 focus-visible:border-ring/70 focus-visible:ring-2 focus-visible:ring-ring/30 disabled:pointer-events-none disabled:opacity-60"
        aria-haspopup="listbox"
        [attr.aria-expanded]="open()"
        [attr.aria-label]="ariaLabel()"
        [disabled]="disabled()"
        (click)="toggle()"
      >
        @if (selectedOption(); as selected) {
          <span class="flex min-w-0 items-center gap-1.5">
            <span class="min-w-0 truncate">{{ selected.label }}</span>
            @if (selected.badge) {
              <span
                class="max-w-20 shrink-0 truncate rounded-md bg-primary/10 px-1.5 py-0.5 text-[9px] font-medium text-primary"
              >
                {{ selected.badge }}
              </span>
            }
          </span>
        } @else {
          <span class="truncate text-muted-foreground">{{ placeholder() }}</span>
        }
        <ng-icon
          name="lucideChevronDown"
          class="shrink-0 text-xs text-muted-foreground transition-transform"
          [class.rotate-180]="open()"
        ></ng-icon>
      </button>

      @if (open()) {
        <div
          class="absolute top-[calc(100%+0.375rem)] z-50 origin-top overflow-hidden rounded-xl border border-border/70 bg-popover/95 text-popover-foreground shadow-md backdrop-blur-sm animate-in fade-in-0 zoom-in-95 slide-in-from-top-2 duration-150"
          [class.left-0]="align() === 'start'"
          [class.right-0]="align() === 'end'"
          [style.width]="panelWidth()"
          [style.min-width]="'100%'"
        >
          <ctx-search-select
            [inputId]="inputId()"
            [ariaLabel]="ariaLabel()"
            [searchPlaceholder]="searchPlaceholder()"
            [emptyText]="emptyText()"
            [options]="options()"
            [value]="value()"
            [showSearch]="searchable()"
            [embedded]="true"
            (selectionChange)="select($event)"
          ></ctx-search-select>

          @if (actionExpanded()) {
            <div class="border-t border-border/70 bg-muted/20 p-1.5">
              <div class="flex min-w-0 items-center gap-1">
                <input
                  type="text"
                  class="h-8 min-w-0 flex-1 rounded-md border border-border/70 bg-background/70 px-2 text-[13px] outline-none focus:border-ring/70 focus:ring-2 focus:ring-ring/30"
                  [placeholder]="actionInputPlaceholder()"
                  [attr.aria-label]="actionInputAriaLabel()"
                  [value]="actionValue()"
                  [disabled]="actionPending()"
                  (input)="updateActionValue($event)"
                  (keydown.enter)="actionConfirm.emit()"
                  (keydown.escape)="actionCancel.emit()"
                  #actionInput
                />
                <button
                  type="button"
                  class="flex size-8 shrink-0 items-center justify-center rounded-md hover:bg-muted"
                  [attr.aria-label]="actionConfirmAriaLabel()"
                  [disabled]="actionPending()"
                  (click)="actionConfirm.emit()"
                >
                  <ng-icon name="lucideCheck"></ng-icon>
                </button>
                <button
                  type="button"
                  class="flex size-8 shrink-0 items-center justify-center rounded-md hover:bg-muted"
                  [attr.aria-label]="actionCancelAriaLabel()"
                  [disabled]="actionPending()"
                  (click)="actionCancel.emit()"
                >
                  <ng-icon name="lucideX"></ng-icon>
                </button>
              </div>
            </div>
          } @else if (actionLabel()) {
            <button
              type="button"
              class="flex h-9 w-full items-center justify-center gap-1.5 border-t border-border/70 bg-muted/20 px-3 text-[11px] font-medium uppercase tracking-[0.06em] text-muted-foreground transition-colors hover:bg-muted/60 hover:text-foreground"
              (click)="triggerAction()"
            >
              <ng-icon name="lucidePlus" class="text-xs"></ng-icon>
              {{ actionLabel() }}
            </button>
          }
        </div>
      }
    </div>
  `,
})
export class SearchDropdownSelectComponent {
  private readonly host = inject(ElementRef<HTMLElement>);

  readonly inputId = input.required<string>();
  readonly options = input.required<readonly SearchSelectOption[]>();
  readonly value = input('');
  readonly ariaLabel = input('Options');
  readonly placeholder = input('Select…');
  readonly searchPlaceholder = input('Search…');
  readonly emptyText = input('No matching options');
  readonly searchable = input(true);
  readonly disabled = input(false);
  readonly align = input<'start' | 'end'>('start');
  readonly panelWidth = input('20rem');
  readonly actionLabel = input('');
  readonly actionExpanded = input(false);
  readonly keepOpenOnAction = input(false);
  readonly actionInputPlaceholder = input('Name');
  readonly actionInputAriaLabel = input('Name');
  readonly actionConfirmAriaLabel = input('Confirm');
  readonly actionCancelAriaLabel = input('Cancel');
  readonly actionValue = input('');
  readonly actionPending = input(false);
  readonly selectionChange = output<string>();
  readonly action = output<void>();
  readonly actionValueChange = output<string>();
  readonly actionConfirm = output<void>();
  readonly actionCancel = output<void>();
  readonly openChange = output<boolean>();

  readonly open = signal(false);
  readonly selectedOption = computed(() =>
    this.options().find((option) => option.value === this.value()),
  );
  private readonly actionInput = viewChild<ElementRef<HTMLInputElement>>('actionInput');
  private readonly focusActionInputEffect = effect(() => {
    if (this.actionExpanded()) {
      this.actionInput()?.nativeElement.focus();
    }
  });

  toggle(): void {
    if (!this.disabled()) {
      this.setOpen(!this.open());
    }
  }

  select(value: string): void {
    this.selectionChange.emit(value);
    this.setOpen(false);
  }

  triggerAction(): void {
    this.action.emit();
    if (!this.keepOpenOnAction()) {
      this.setOpen(false);
    }
  }

  updateActionValue(event: Event): void {
    this.actionValueChange.emit((event.target as HTMLInputElement).value);
  }

  @HostListener('document:mousedown', ['$event'])
  closeOnOutsideClick(event: MouseEvent): void {
    const target = event.target as Node | null;
    if (this.open() && target && !this.host.nativeElement.contains(target)) {
      this.setOpen(false);
    }
  }

  private setOpen(open: boolean): void {
    this.open.set(open);
    this.openChange.emit(open);
  }
}

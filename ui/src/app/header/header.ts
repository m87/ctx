import { Component, computed, effect, inject, signal } from '@angular/core';
import { NgIcon, provideIcons } from '@ng-icons/core';
import {
  lucideCalendar,
  lucideGanttChart,
  lucidePanelLeft,
  lucidePause,
  lucideSearch,
  lucideX,
} from '@ng-icons/lucide';
import { HlmBreadCrumbImports } from '@spartan-ng/helm/breadcrumb';
import { HlmInputImports } from '@spartan-ng/helm/input';
import { BreadcrumbService } from '../header/breadcrumbs';
import { SidebarStore } from '../sidebar/sidebar.store';
import { HlmButtonImports } from '@spartan-ng/helm/button';
import { injectMutation, injectQuery } from '@tanstack/angular-query-experimental';
import { Context } from '../../api/context.service';
import { ContextQueries } from '../../api/context.quries';
import { ContextMutations } from '../../api/context.mutations';
import { Router, RouterLink } from '@angular/router';
import { HlmDatePickerImports } from '@spartan-ng/helm/date-picker';
import { DateTime } from 'luxon';

@Component({
  selector: 'app-header',
  imports: [
    HlmBreadCrumbImports,
    NgIcon,
    HlmInputImports,
    HlmButtonImports,
    RouterLink,
    HlmDatePickerImports,
  ],
  providers: [
    provideIcons({
      lucideGanttChart,
      lucideSearch,
      lucideCalendar,
      lucidePanelLeft,
      lucidePause,
      lucideX,
    }),
  ],
  template: `
    <div class="w-full border-b bg-card/70">
      <div class="w-full h-12 flex items-center justify-between px-3 gap-2">
        <div class="flex items-center gap-3 min-w-0">
          <div class="gap-2 flex items-center shrink-0">
            <button
              hlmBtn
              variant="ghost"
              class="md:hidden h-8 w-8 px-0"
              (click)="sidebar.toggleMobile()"
              aria-label="Toggle sidebar"
            >
              <ng-icon name="lucidePanelLeft"></ng-icon>
            </button>
            <ng-icon name="lucideGanttChart" class="cursor-pointer text-primary"></ng-icon>
            <span class="font-semibold tracking-tight text-primary">Ctx</span>
          </div>

          <div class="hidden md:block flex-1 max-w-xl w-xl relative">
            <input
              hlmInput
              type="text"
              placeholder="Search or create new context"
              class="h-8 w-full text-xs"
              [value]="searchTerm()"
              (input)="onSearchInput($event)"
              (focus)="onSearchFocus()"
              (blur)="onSearchBlur()"
              (keydown)="onSearchKeydown($event)"
            />

            @if (showSuggestions()) {
              <div
                class="absolute top-9 left-0 right-0 z-30 border rounded-md bg-popover text-popover-foreground shadow-sm p-1 max-h-56 overflow-auto"
              >
                @for (context of filteredContexts(); track context.id; let index = $index) {
                  <button
                    type="button"
                    class="w-full text-left px-2 py-1.5 rounded-sm text-xs hover:bg-muted"
                    [class.bg-muted]="activeSuggestionIndex() === index"
                    [class.text-foreground]="activeSuggestionIndex() === index"
                    [class.text-muted-foreground]="activeSuggestionIndex() !== index"
                    (mouseenter)="setActiveSuggestionIndex(index)"
                    (mousedown)="selectContext(context)"
                  >
                    {{ context.name }}
                  </button>
                }
              </div>
            }
          </div>
        </div>

        <div class="flex items-center gap-2 shrink-0">
          <button
            hlmBtn
            variant="outline"
            class="h-8 w-8 px-0 md:hidden"
            (click)="openMobileSearch()"
            [class.hidden]="mobileSearchOpen()"
            aria-label="Open search"
          >
            <ng-icon name="lucideSearch"></ng-icon>
          </button>

          @if (activeContextName()) {
            <div class="flex items-center max-w-40">
              <div
                class="h-8 px-2 rounded-l-md border bg-muted/40 flex items-center gap-2 max-w-28"
              >
                <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 shrink-0"></span>
                <span class="text-xs font-medium truncate">{{ activeContextName() }}</span>
              </div>
              <button
                hlmBtn
                variant="outline"
                class="h-8 w-8 px-0 sm:px-2 rounded-l-none -ml-px"
                [disabled]="freeContextMutation.isPending()"
                (click)="stopContext()"
                aria-label="Stop active context"
              >
                <ng-icon name="lucidePause" class="text-xs"></ng-icon>
              </button>
            </div>
          } @else {
            <div
              class="h-8 px-2 rounded-md border bg-muted/30 flex items-center max-w-28 sm:max-w-none"
            >
              <span class="text-xs text-muted-foreground truncate">No context</span>
            </div>
          }

          <div class="flex items-center gap-2">
            <hlm-date-picker
              [buttonId]="'dupa'"
              align="end"
              class="w-auto"
              [autoCloseOnSelect]="true"
              (dateChange)="navigateToDate($event)"
            >
              <button
                id="dupa"
                class="h-8 px-2 sm:px-3 text-xs text-muted-foreground gap-2 flex items-center"
                aria-label="Select date"
              >
                <span class="hidden sm:inline">{{ today() }}</span>
                <ng-icon name="lucideCalendar" class="cursor-pointer"></ng-icon>
              </button>
            </hlm-date-picker>
          </div>
          <button
            hlmBtn
            variant="outline"
            class="hidden sm:inline-flex h-8 px-3 text-xs"
            [routerLink]="['/day', today()]"
          >
            Today
          </button>
        </div>
      </div>

      @if (mobileSearchOpen()) {
        <div class="md:hidden px-3 pb-2">
          <div class="relative">
            <input
              hlmInput
              type="text"
              placeholder="Search or create new context"
              class="h-8 w-full text-xs pr-9"
              [value]="searchTerm()"
              (input)="onSearchInput($event)"
              (focus)="onSearchFocus()"
              (blur)="onSearchBlur()"
              (keydown)="onSearchKeydown($event)"
            />
            <button
              hlmBtn
              variant="ghost"
              class="absolute right-1 top-1/2 -translate-y-1/2 h-6 w-6 px-0"
              (click)="closeMobileSearch()"
              aria-label="Close search"
            >
              <ng-icon name="lucideX"></ng-icon>
            </button>

            @if (showSuggestions()) {
              <div
                class="absolute top-9 left-0 right-0 z-30 border rounded-md bg-popover text-popover-foreground shadow-sm p-1 max-h-56 overflow-auto"
              >
                @for (context of filteredContexts(); track context.id; let index = $index) {
                  <button
                    type="button"
                    class="w-full text-left px-2 py-1.5 rounded-sm text-xs hover:bg-muted"
                    [class.bg-muted]="activeSuggestionIndex() === index"
                    [class.text-foreground]="activeSuggestionIndex() === index"
                    [class.text-muted-foreground]="activeSuggestionIndex() !== index"
                    (mouseenter)="setActiveSuggestionIndex(index)"
                    (mousedown)="selectContext(context)"
                  >
                    {{ context.name }}
                  </button>
                }
              </div>
            }
          </div>
        </div>
      }
    </div>
  `,
  styles: ``,
})
export class HeaderComponent {
  breadcrumbService = inject(BreadcrumbService);
  sidebar = inject(SidebarStore);
  private contextQueries = inject(ContextQueries);
  private contextMutations = inject(ContextMutations);
  private router = inject(Router);
  today = signal(DateTime.local().toFormat('yyyy-MM-dd'));

  listContextsQuery = injectQuery(() => this.contextQueries.list());
  switchContextMutation = injectMutation(() => this.contextMutations.switch());
  freeContextMutation = injectMutation(() => this.contextMutations.free());
  activeContextQuery = injectQuery(() => this.contextQueries.active());
  activeContextName = computed(() => this.activeContextQuery.data()?.name ?? '');

  readonly searchTerm = signal<string>('');
  readonly searchFocused = signal<boolean>(false);
  readonly mobileSearchOpen = signal<boolean>(false);
  readonly activeSuggestionIndex = signal<number>(-1);
  readonly contexts = computed<readonly Context[]>(() => this.listContextsQuery.data() ?? []);
  readonly filteredContexts = computed<readonly Context[]>(() => {
    const term = this.searchTerm().trim().toLowerCase();
    if (!term) {
      return this.contexts();
    }
    return this.contexts().filter((context) => context.name.toLowerCase().includes(term));
  });
  readonly showSuggestions = computed<boolean>(
    () =>
      this.searchFocused() &&
      this.searchTerm().trim().length > 0 &&
      this.filteredContexts().length > 0,
  );

  private readonly syncActiveSuggestionEffect = effect(() => {
    const visible = this.showSuggestions();
    const suggestionsLength = this.filteredContexts().length;
    const currentIndex = this.activeSuggestionIndex();

    if (!visible || suggestionsLength === 0) {
      if (currentIndex !== -1) {
        this.activeSuggestionIndex.set(-1);
      }
      return;
    }

    if (currentIndex < 0 || currentIndex >= suggestionsLength) {
      this.activeSuggestionIndex.set(0);
    }
  });

  onSearchInput(event: Event): void {
    const target = event.target as HTMLInputElement;
    this.searchTerm.set(target.value);
  }

  onSearchFocus(): void {
    this.searchFocused.set(true);
  }

  openMobileSearch(): void {
    this.mobileSearchOpen.set(true);
    this.searchFocused.set(true);
  }

  closeMobileSearch(): void {
    this.mobileSearchOpen.set(false);
    this.searchFocused.set(false);
    this.activeSuggestionIndex.set(-1);
  }

  onSearchBlur(): void {
    setTimeout(() => this.searchFocused.set(false), 100);
  }

  setActiveSuggestionIndex(index: number): void {
    this.activeSuggestionIndex.set(index);
  }

  selectContext(context: Context): void {
    this.searchTerm.set(context.name);
    this.searchFocused.set(false);
    this.activeSuggestionIndex.set(-1);
    this.mobileSearchOpen.set(false);
    this.switchContextMutation.mutate(context);
  }

  onSearchKeydown(event: KeyboardEvent): void {
    if (event.key === 'ArrowDown') {
      if (!this.showSuggestions()) {
        return;
      }
      event.preventDefault();
      const suggestionsLength = this.filteredContexts().length;
      const currentIndex = this.activeSuggestionIndex();
      this.activeSuggestionIndex.set((currentIndex + 1 + suggestionsLength) % suggestionsLength);
      return;
    }

    if (event.key === 'ArrowUp') {
      if (!this.showSuggestions()) {
        return;
      }
      event.preventDefault();
      const suggestionsLength = this.filteredContexts().length;
      const currentIndex = this.activeSuggestionIndex();
      this.activeSuggestionIndex.set((currentIndex - 1 + suggestionsLength) % suggestionsLength);
      return;
    }

    if (event.key === 'Escape') {
      this.searchFocused.set(false);
      this.activeSuggestionIndex.set(-1);
      this.mobileSearchOpen.set(false);
      return;
    }

    if (event.key === 'Enter') {
      event.preventDefault();
      this.onSearchEnter();
    }
  }

  onSearchEnter(): void {
    const term = this.searchTerm().trim();
    if (!term) {
      return;
    }

    if (this.showSuggestions()) {
      const selectedContext = this.filteredContexts()[this.activeSuggestionIndex()];
      if (selectedContext) {
        this.selectContext(selectedContext);
        return;
      }
    }

    const existingContext = this.contexts().find(
      (context) => context.name.toLowerCase() === term.toLowerCase(),
    );

    if (existingContext) {
      this.switchContextMutation.mutate(existingContext);
      this.searchTerm.set(existingContext.name);
      this.searchFocused.set(false);
      this.activeSuggestionIndex.set(-1);
      this.mobileSearchOpen.set(false);
      return;
    }

    this.switchContextMutation.mutate({ id: '', name: term } as Context);
    this.searchFocused.set(false);
    this.activeSuggestionIndex.set(-1);
    this.mobileSearchOpen.set(false);
  }

  navigateToDate(date: Date): void {
    this.router.navigate(['day', DateTime.fromJSDate(date).toFormat('yyyy-MM-dd')]);
  }

  stopContext(): void {
    this.freeContextMutation.mutate();
  }
}

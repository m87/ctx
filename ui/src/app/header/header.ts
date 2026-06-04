import { Component, computed, effect, inject, signal } from '@angular/core';
import { toObservable, toSignal } from '@angular/core/rxjs-interop';
import { NgIcon, provideIcons } from '@ng-icons/core';
import {
  lucideCalendar,
  lucideClock3,
  lucideGanttChart,
  lucideHistory,
  lucidePanelLeft,
  lucidePause,
  lucidePlus,
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
import { NavigationEnd, Router, RouterLink } from '@angular/router';
import { HlmDatePickerImports } from '@spartan-ng/helm/date-picker';
import { DateTime } from 'luxon';
import { catchError, filter, forkJoin, map, of, startWith, switchMap } from 'rxjs';
import { ContextService, ContextStats } from '../../api/context.service';
import { SettingsQueries } from '../../api/settings.queries';

const firstDayKey = 'client.general.firstDay';

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
      lucidePlus,
      lucideClock3,
      lucideHistory,
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
                class="absolute top-9 left-0 right-0 z-30 border rounded-md bg-popover text-popover-foreground shadow-sm p-1 max-h-72 overflow-auto origin-top animate-in fade-in-0 zoom-in-95 slide-in-from-top-2 duration-200"
              >
                <button
                  type="button"
                  class="w-full flex items-center gap-2 text-left px-2 py-2 rounded-sm text-xs hover:bg-muted border border-dashed border-border/80 mb-1"
                  [class.bg-muted]="activeSuggestionIndex() === 0"
                  [class.text-foreground]="activeSuggestionIndex() === 0"
                  [class.text-muted-foreground]="activeSuggestionIndex() !== 0"
                  (mouseenter)="setActiveSuggestionIndex(0)"
                  (mousedown)="createContextFromTerm(searchTerm().trim())"
                >
                  <ng-icon name="lucidePlus" class="text-xs"></ng-icon>
                  <span class="truncate font-medium">{{ searchTerm().trim() }}</span>
                </button>

                @if (dayMatchedContexts().length > 0) {
                  <div
                    class="px-2 pt-1 pb-1 text-[10px] uppercase tracking-[0.08em] text-muted-foreground"
                  >
                    {{ daySectionLabel() }}
                  </div>
                }
                @for (context of dayMatchedContexts(); track context.id) {
                  <button
                    type="button"
                    class="w-full flex items-center justify-between gap-2 text-left px-2 py-1.5 rounded-sm text-xs hover:bg-muted"
                    [class.bg-muted]="activeSuggestionIndex() === suggestionIndex(context.id)"
                    [class.text-foreground]="
                      activeSuggestionIndex() === suggestionIndex(context.id)
                    "
                    [class.text-muted-foreground]="
                      activeSuggestionIndex() !== suggestionIndex(context.id)
                    "
                    (mouseenter)="setActiveSuggestionIndex(suggestionIndex(context.id))"
                    (mousedown)="selectContext(context)"
                  >
                    <span class="truncate">{{ context.name }}</span>
                    <span
                      class="shrink-0 text-[10px] text-muted-foreground/80 flex items-center gap-2"
                    >
                      @if (contextTodayDuration(context.id); as todayDuration) {
                        <span class="inline-flex items-center gap-1">
                          <ng-icon name="lucideClock3" class="text-[10px]"></ng-icon>
                          {{ todayDuration }}
                        </span>
                      }
                      @if (contextTotalDuration(context.id); as totalDuration) {
                        <span class="inline-flex items-center gap-1">
                          <ng-icon name="lucideHistory" class="text-[10px]"></ng-icon>
                          {{ totalDuration }}
                        </span>
                      }
                    </span>
                  </button>
                }

                @if (otherMatchedContexts().length > 0) {
                  <div class="my-1 border-t border-border/70"></div>
                }
                @for (context of otherMatchedContexts(); track context.id) {
                  <button
                    type="button"
                    class="w-full flex items-center justify-between gap-2 text-left px-2 py-1.5 rounded-sm text-xs hover:bg-muted"
                    [class.bg-muted]="activeSuggestionIndex() === suggestionIndex(context.id)"
                    [class.text-foreground]="
                      activeSuggestionIndex() === suggestionIndex(context.id)
                    "
                    [class.text-muted-foreground]="
                      activeSuggestionIndex() !== suggestionIndex(context.id)
                    "
                    (mouseenter)="setActiveSuggestionIndex(suggestionIndex(context.id))"
                    (mousedown)="selectContext(context)"
                  >
                    <span class="truncate">{{ context.name }}</span>
                    <span
                      class="shrink-0 text-[10px] text-muted-foreground/80 flex items-center gap-2"
                    >
                      @if (contextTodayDuration(context.id); as todayDuration) {
                        <span class="inline-flex items-center gap-1">
                          <ng-icon name="lucideClock3" class="text-[10px]"></ng-icon>
                          {{ todayDuration }}
                        </span>
                      }
                      @if (contextTotalDuration(context.id); as totalDuration) {
                        <span class="inline-flex items-center gap-1">
                          <ng-icon name="lucideHistory" class="text-[10px]"></ng-icon>
                          {{ totalDuration }}
                        </span>
                      }
                    </span>
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
              [weekStartsOn]="weekStartsOn()"
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
                class="absolute top-9 left-0 right-0 z-30 border rounded-md bg-popover text-popover-foreground shadow-sm p-1 max-h-72 overflow-auto origin-top animate-in fade-in-0 zoom-in-95 slide-in-from-top-2 duration-200"
              >
                <button
                  type="button"
                  class="w-full flex items-center gap-2 text-left px-2 py-2 rounded-sm text-xs hover:bg-muted border border-dashed border-border/80 mb-1"
                  [class.bg-muted]="activeSuggestionIndex() === 0"
                  [class.text-foreground]="activeSuggestionIndex() === 0"
                  [class.text-muted-foreground]="activeSuggestionIndex() !== 0"
                  (mouseenter)="setActiveSuggestionIndex(0)"
                  (mousedown)="createContextFromTerm(searchTerm().trim())"
                >
                  <ng-icon name="lucidePlus" class="text-xs"></ng-icon>
                  <span class="truncate font-medium">{{ searchTerm().trim() }}</span>
                </button>

                @if (dayMatchedContexts().length > 0) {
                  <div
                    class="px-2 pt-1 pb-1 text-[10px] uppercase tracking-[0.08em] text-muted-foreground"
                  >
                    {{ daySectionLabel() }}
                  </div>
                }
                @for (context of dayMatchedContexts(); track context.id) {
                  <button
                    type="button"
                    class="w-full flex items-center justify-between gap-2 text-left px-2 py-1.5 rounded-sm text-xs hover:bg-muted"
                    [class.bg-muted]="activeSuggestionIndex() === suggestionIndex(context.id)"
                    [class.text-foreground]="
                      activeSuggestionIndex() === suggestionIndex(context.id)
                    "
                    [class.text-muted-foreground]="
                      activeSuggestionIndex() !== suggestionIndex(context.id)
                    "
                    (mouseenter)="setActiveSuggestionIndex(suggestionIndex(context.id))"
                    (mousedown)="selectContext(context)"
                  >
                    <span class="truncate">{{ context.name }}</span>
                    <span
                      class="shrink-0 text-[10px] text-muted-foreground/80 flex items-center gap-2"
                    >
                      @if (contextTodayDuration(context.id); as todayDuration) {
                        <span class="inline-flex items-center gap-1">
                          <ng-icon name="lucideClock3" class="text-[10px]"></ng-icon>
                          {{ todayDuration }}
                        </span>
                      }
                      @if (contextTotalDuration(context.id); as totalDuration) {
                        <span class="inline-flex items-center gap-1">
                          <ng-icon name="lucideHistory" class="text-[10px]"></ng-icon>
                          {{ totalDuration }}
                        </span>
                      }
                    </span>
                  </button>
                }

                @if (otherMatchedContexts().length > 0) {
                  <div class="my-1 border-t border-border/70"></div>
                }
                @for (context of otherMatchedContexts(); track context.id) {
                  <button
                    type="button"
                    class="w-full flex items-center justify-between gap-2 text-left px-2 py-1.5 rounded-sm text-xs hover:bg-muted"
                    [class.bg-muted]="activeSuggestionIndex() === suggestionIndex(context.id)"
                    [class.text-foreground]="
                      activeSuggestionIndex() === suggestionIndex(context.id)
                    "
                    [class.text-muted-foreground]="
                      activeSuggestionIndex() !== suggestionIndex(context.id)
                    "
                    (mouseenter)="setActiveSuggestionIndex(suggestionIndex(context.id))"
                    (mousedown)="selectContext(context)"
                  >
                    <span class="truncate">{{ context.name }}</span>
                    <span
                      class="shrink-0 text-[10px] text-muted-foreground/80 flex items-center gap-2"
                    >
                      @if (contextTodayDuration(context.id); as todayDuration) {
                        <span class="inline-flex items-center gap-1">
                          <ng-icon name="lucideClock3" class="text-[10px]"></ng-icon>
                          {{ todayDuration }}
                        </span>
                      }
                      @if (contextTotalDuration(context.id); as totalDuration) {
                        <span class="inline-flex items-center gap-1">
                          <ng-icon name="lucideHistory" class="text-[10px]"></ng-icon>
                          {{ totalDuration }}
                        </span>
                      }
                    </span>
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
  private contextService = inject(ContextService);
  private settingsQueries = inject(SettingsQueries);
  private router = inject(Router);
  today = signal(DateTime.local().toFormat('yyyy-MM-dd'));

  listContextsQuery = injectQuery(() => this.contextQueries.list());
  settingsQuery = injectQuery(() => this.settingsQueries.settings());
  switchContextMutation = injectMutation(() => this.contextMutations.switch());
  freeContextMutation = injectMutation(() => this.contextMutations.free());
  activeContextQuery = injectQuery(() => this.contextQueries.active());
  selectedDate = toSignal(
    this.router.events.pipe(
      filter((event): event is NavigationEnd => event instanceof NavigationEnd),
      startWith(null),
      map(() => this.resolveSelectedDate()),
    ),
    { initialValue: this.today() },
  );
  dayStatsQuery = injectQuery(() => this.contextQueries.dayStats(this.selectedDate()));
  activeContextName = computed(() => this.activeContextQuery.data()?.name ?? '');
  weekStartsOn = computed(() => (this.settingsQuery.data()?.[firstDayKey] === 'Sunday' ? 0 : 1));
  daySectionLabel = computed(() =>
    this.selectedDate() === DateTime.local().toFormat('yyyy-MM-dd') ? 'Today' : this.selectedDate(),
  );

  readonly searchTerm = signal<string>('');
  readonly searchFocused = signal<boolean>(false);
  readonly mobileSearchOpen = signal<boolean>(false);
  readonly activeSuggestionIndex = signal<number>(-1);
  readonly contexts = computed<readonly Context[]>(() => this.listContextsQuery.data() ?? []);
  readonly filteredContexts = computed<readonly Context[]>(() => {
    const term = this.searchTerm().trim().toLowerCase();
    if (!term) {
      return [];
    }
    return this.contexts().filter((context) => context.name.toLowerCase().includes(term));
  });
  readonly usedContextIdsForDay = computed(
    () => new Set(this.dayStatsQuery.data()?.contextStats.map((stats) => stats.contextId) ?? []),
  );
  readonly dayMatchedContexts = computed<readonly Context[]>(() =>
    this.filteredContexts().filter((context) => this.usedContextIdsForDay().has(context.id)),
  );
  readonly otherMatchedContexts = computed<readonly Context[]>(() =>
    this.filteredContexts().filter((context) => !this.usedContextIdsForDay().has(context.id)),
  );
  readonly suggestionContexts = computed<readonly Context[]>(() => [
    ...this.dayMatchedContexts(),
    ...this.otherMatchedContexts(),
  ]);
  readonly suggestionCount = computed<number>(() =>
    this.searchTerm().trim().length > 0 ? this.suggestionContexts().length + 1 : 0,
  );
  readonly statsByContextId = toSignal(
    toObservable(
      computed(() => ({
        contexts: this.filteredContexts(),
        date: this.selectedDate(),
      })),
    ).pipe(
      switchMap(({ contexts, date }) => {
        if (contexts.length === 0) {
          return of({} as Record<string, ContextStats>);
        }

        return forkJoin(
          contexts.map((context) =>
            this.contextService.getStats(context.id, date).pipe(
              map((stats) => [context.id, stats] as const),
              catchError(() => of([context.id, null] as const)),
            ),
          ),
        ).pipe(
          map((entries) =>
            Object.fromEntries(
              entries.filter(
                (entry): entry is readonly [string, ContextStats] => entry[1] !== null,
              ),
            ),
          ),
        );
      }),
    ),
    { initialValue: {} as Record<string, ContextStats> },
  );
  readonly showSuggestions = computed<boolean>(
    () => this.searchFocused() && this.searchTerm().trim().length > 0,
  );

  private readonly syncActiveSuggestionEffect = effect(() => {
    const visible = this.showSuggestions();
    const suggestionsLength = this.suggestionCount();
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
    this.resetSearchUi();
    this.switchContextMutation.mutate(context);
  }

  createContextFromTerm(term: string): void {
    const normalizedTerm = term.trim();
    if (!normalizedTerm) {
      return;
    }
    this.searchTerm.set(normalizedTerm);
    this.resetSearchUi();
    this.switchContextMutation.mutate({ id: '', name: normalizedTerm } as Context);
  }

  onSearchKeydown(event: KeyboardEvent): void {
    if (event.key === 'ArrowDown') {
      if (!this.showSuggestions()) {
        return;
      }
      event.preventDefault();
      const suggestionsLength = this.suggestionCount();
      const currentIndex = this.activeSuggestionIndex();
      this.activeSuggestionIndex.set((currentIndex + 1 + suggestionsLength) % suggestionsLength);
      return;
    }

    if (event.key === 'ArrowUp') {
      if (!this.showSuggestions()) {
        return;
      }
      event.preventDefault();
      const suggestionsLength = this.suggestionCount();
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

    if (this.showSuggestions() && this.activeSuggestionIndex() > 0) {
      const selectedContext = this.suggestionContexts()[this.activeSuggestionIndex() - 1];
      if (selectedContext) {
        this.selectContext(selectedContext);
        return;
      }
    }

    this.createContextFromTerm(term);
  }

  suggestionIndex(contextId: string): number {
    return this.suggestionContexts().findIndex((context) => context.id === contextId) + 1;
  }

  contextTodayDuration(contextId: string): string | null {
    const stats = this.statsByContextId()[contextId];
    return this.formatDuration(stats?.duration ?? 0);
  }

  contextTotalDuration(contextId: string): string | null {
    const stats = this.statsByContextId()[contextId];
    return this.formatDuration(stats?.totalDuration ?? 0);
  }

  private resolveSelectedDate(): string {
    const dayMatch = this.router.url.match(/\/day\/(\d{4}-\d{2}-\d{2})/);
    return dayMatch?.[1] ?? DateTime.local().toFormat('yyyy-MM-dd');
  }

  private formatDuration(duration: number): string | null {
    const totalMinutes = Math.max(0, Math.floor(duration / 60000000000));
    if (totalMinutes === 0) {
      return null;
    }

    const hours = Math.floor(totalMinutes / 60);
    const minutes = totalMinutes % 60;

    if (hours > 0 && minutes > 0) {
      return `${hours}h ${minutes}m`;
    }
    if (hours > 0) {
      return `${hours}h`;
    }
    return `${minutes}m`;
  }

  private resetSearchUi(): void {
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

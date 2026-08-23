import { Component, computed, effect, inject, signal } from '@angular/core';
import { toObservable, toSignal } from '@angular/core/rxjs-interop';
import { NgIcon, provideIcons } from '@ng-icons/core';
import {
  lucideCalendar,
  lucideClock3,
  lucideFolder,
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
import {
  Context,
  ContextService,
  ContextStats,
  ProjectMetadata,
} from '../../api/context/context.service';
import { ContextQueries } from '../../api/context/context.queries';
import { ContextMutations } from '../../api/context/context.mutations';
import { NavigationEnd, Router, RouterLink } from '@angular/router';
import { HlmDatePicker, HlmDatePickerTrigger } from '@spartan-ng/helm/date-picker';
import { DateTime } from 'luxon';
import { catchError, filter, forkJoin, map, of, startWith, switchMap } from 'rxjs';
import { ProjectQueries } from '../../api/project/project.queries';
import { Project } from '../../api/project/project.service';
import { SettingsQueries } from '../../api/settings/settings.queries';
import { Store } from '@ngxs/store';
import { WorkspaceState } from '../sidebar/workspace.state';
import { IntervalQueries } from '../../api/interval/interval.queries';
import { TimeZoneService } from '../shared/time-zone.service';
import { SearchProjectBadgeComponent } from './search-project-badge.component';
import { parseProjectPicker } from './search-project-picker';
import { HlmSkeletonImports } from '@spartan-ng/helm/skeleton';

const firstDayKey = 'client.general.firstDay';

@Component({
  selector: 'ctx-header',
  imports: [
    HlmBreadCrumbImports,
    NgIcon,
    HlmInputImports,
    HlmButtonImports,
    RouterLink,
    HlmDatePicker,
    HlmDatePickerTrigger,
    SearchProjectBadgeComponent,
    HlmSkeletonImports,
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
      lucideFolder,
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
              placeholder="Search or create new context; use #project"
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
                @if (projectCreationSuggestion(); as project) {
                  <button
                    type="button"
                    class="w-full flex items-center gap-2 text-left px-2 py-2 rounded-sm text-xs hover:bg-muted border border-dashed border-border/80 mb-1"
                    [class.bg-muted]="activeSuggestionIndex() === 0"
                    [class.text-foreground]="activeSuggestionIndex() === 0"
                    [class.text-muted-foreground]="activeSuggestionIndex() !== 0"
                    (mouseenter)="setActiveSuggestionIndex(0)"
                    (mousedown)="createContextFromTerm(searchTerm().trim(), project)"
                  >
                    <ng-icon name="lucidePlus" class="text-xs shrink-0"></ng-icon>
                    <span class="min-w-0 flex-1 truncate font-medium">{{
                      searchTerm().trim()
                    }}</span>
                    <span
                      class="max-w-40 shrink-0 inline-flex items-center gap-1 rounded-md bg-primary/10 px-1.5 py-0.5 text-[10px] text-primary"
                    >
                      <ng-icon name="lucideFolder" class="shrink-0 text-[10px]"></ng-icon>
                      <span class="truncate">{{ project.name }}</span>
                    </span>
                  </button>
                  <button
                    type="button"
                    class="w-full flex items-center gap-2 text-left px-2 py-2 rounded-sm text-xs hover:bg-muted border border-dashed border-border/80 mb-1"
                    [class.bg-muted]="activeSuggestionIndex() === 1"
                    [class.text-foreground]="activeSuggestionIndex() === 1"
                    [class.text-muted-foreground]="activeSuggestionIndex() !== 1"
                    (mouseenter)="setActiveSuggestionIndex(1)"
                    (mousedown)="createContextFromTerm(searchTerm().trim(), null)"
                  >
                    <ng-icon name="lucidePlus" class="text-xs shrink-0"></ng-icon>
                    <span class="min-w-0 flex-1 truncate font-medium">{{
                      searchTerm().trim()
                    }}</span>
                    <span class="shrink-0 text-[10px] text-muted-foreground">Workspace</span>
                  </button>
                } @else if (!projectPickerMode()) {
                  <button
                    type="button"
                    class="w-full flex items-center gap-2 text-left px-2 py-2 rounded-sm text-xs hover:bg-muted border border-dashed border-border/80 mb-1"
                    [class.bg-muted]="activeSuggestionIndex() === 0"
                    [class.text-foreground]="activeSuggestionIndex() === 0"
                    [class.text-muted-foreground]="activeSuggestionIndex() !== 0"
                    (mouseenter)="setActiveSuggestionIndex(0)"
                    (mousedown)="createContextFromTerm(searchTerm().trim(), null)"
                  >
                    <ng-icon name="lucidePlus" class="text-xs shrink-0"></ng-icon>
                    <span class="min-w-0 flex-1 truncate font-medium">{{
                      searchTerm().trim()
                    }}</span>
                    <span class="shrink-0 text-[10px] text-muted-foreground">Workspace</span>
                  </button>
                }

                @if (matchedProjects().length > 0) {
                  <div
                    class="px-2 pt-1 pb-1 text-[10px] uppercase tracking-[0.08em] text-muted-foreground"
                  >
                    Project for context
                  </div>
                }
                @for (project of matchedProjects(); track project.id) {
                  <button
                    type="button"
                    class="w-full flex items-center justify-between gap-2 text-left px-2 py-1.5 rounded-sm text-xs hover:bg-muted"
                    [class.bg-muted]="
                      activeSuggestionIndex() === projectSuggestionIndex(project.id)
                    "
                    [class.text-foreground]="
                      activeSuggestionIndex() === projectSuggestionIndex(project.id)
                    "
                    [class.text-muted-foreground]="
                      activeSuggestionIndex() !== projectSuggestionIndex(project.id)
                    "
                    (mouseenter)="setActiveSuggestionIndex(projectSuggestionIndex(project.id))"
                    (mousedown)="selectProjectSuggestion($event, project)"
                  >
                    <span class="min-w-0 flex items-center gap-1.5">
                      <ng-icon name="lucideFolder" class="shrink-0 text-[11px]"></ng-icon>
                      <span class="truncate">{{ project.name }}</span>
                    </span>
                    <span class="shrink-0 text-[10px] text-muted-foreground">Project</span>
                  </button>
                }

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
                      @if (context.project; as project) {
                        <ctx-search-project-badge [project]="project" />
                      }
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
                      @if (context.project; as project) {
                        <ctx-search-project-badge [project]="project" />
                      }
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

                @if (archivedMatchedContexts().length > 0) {
                  <div class="my-1 border-t border-border/70"></div>
                  <div
                    class="px-2 pt-1 pb-1 text-[10px] uppercase tracking-[0.08em] text-muted-foreground"
                  >
                    Archived
                  </div>
                }
                @for (context of archivedMatchedContexts(); track context.id) {
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
                    (mousedown)="openContext(context)"
                  >
                    <span class="truncate">{{ context.name }}</span>
                    <span
                      class="shrink-0 text-[10px] text-muted-foreground/80 flex items-center gap-2"
                    >
                      @if (context.project; as project) {
                        <ctx-search-project-badge [project]="project" />
                      }
                      <span class="rounded-sm border px-1 py-0.5 leading-none">Archived</span>
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

          @if (activeContextQuery.isLoading()) {
            <hlm-skeleton class="h-8 w-28"></hlm-skeleton>
          } @else if (activeContextName()) {
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
                [attr.aria-busy]="freeContextMutation.isPending()"
                (click)="stopContext()"
                aria-label="Stop active context"
              >
                @if (freeContextMutation.isPending()) {
                  <span
                    class="size-3.5 rounded-full border-2 border-current border-t-transparent animate-spin"
                    aria-hidden="true"
                  ></span>
                } @else {
                  <ng-icon name="lucidePause" class="text-xs"></ng-icon>
                }
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
              align="end"
              class="w-auto"
              [autoCloseOnSelect]="true"
              [weekStartsOn]="weekStartsOn()"
              (dateChange)="navigateToDate($event)"
            >
              <hlm-date-picker-trigger buttonId="dupa" class="w-auto" aria-label="Select date">
                <span
                  class="h-8 px-2 sm:px-3 text-xs text-muted-foreground gap-2 flex items-center"
                >
                  <span class="hidden sm:inline">{{ today() }}</span>
                  <ng-icon name="lucideCalendar" class="cursor-pointer"></ng-icon>
                </span>
              </hlm-date-picker-trigger>
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
              placeholder="Search or create new context; use #project"
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
                @if (projectCreationSuggestion(); as project) {
                  <button
                    type="button"
                    class="w-full flex items-center gap-2 text-left px-2 py-2 rounded-sm text-xs hover:bg-muted border border-dashed border-border/80 mb-1"
                    [class.bg-muted]="activeSuggestionIndex() === 0"
                    [class.text-foreground]="activeSuggestionIndex() === 0"
                    [class.text-muted-foreground]="activeSuggestionIndex() !== 0"
                    (mouseenter)="setActiveSuggestionIndex(0)"
                    (mousedown)="createContextFromTerm(searchTerm().trim(), project)"
                  >
                    <ng-icon name="lucidePlus" class="text-xs shrink-0"></ng-icon>
                    <span class="min-w-0 flex-1 truncate font-medium">{{
                      searchTerm().trim()
                    }}</span>
                    <span
                      class="max-w-32 shrink-0 inline-flex items-center gap-1 rounded-md bg-primary/10 px-1.5 py-0.5 text-[10px] text-primary"
                    >
                      <ng-icon name="lucideFolder" class="shrink-0 text-[10px]"></ng-icon>
                      <span class="truncate">{{ project.name }}</span>
                    </span>
                  </button>
                  <button
                    type="button"
                    class="w-full flex items-center gap-2 text-left px-2 py-2 rounded-sm text-xs hover:bg-muted border border-dashed border-border/80 mb-1"
                    [class.bg-muted]="activeSuggestionIndex() === 1"
                    [class.text-foreground]="activeSuggestionIndex() === 1"
                    [class.text-muted-foreground]="activeSuggestionIndex() !== 1"
                    (mouseenter)="setActiveSuggestionIndex(1)"
                    (mousedown)="createContextFromTerm(searchTerm().trim(), null)"
                  >
                    <ng-icon name="lucidePlus" class="text-xs shrink-0"></ng-icon>
                    <span class="min-w-0 flex-1 truncate font-medium">{{
                      searchTerm().trim()
                    }}</span>
                    <span class="shrink-0 text-[10px] text-muted-foreground">Workspace</span>
                  </button>
                } @else if (!projectPickerMode()) {
                  <button
                    type="button"
                    class="w-full flex items-center gap-2 text-left px-2 py-2 rounded-sm text-xs hover:bg-muted border border-dashed border-border/80 mb-1"
                    [class.bg-muted]="activeSuggestionIndex() === 0"
                    [class.text-foreground]="activeSuggestionIndex() === 0"
                    [class.text-muted-foreground]="activeSuggestionIndex() !== 0"
                    (mouseenter)="setActiveSuggestionIndex(0)"
                    (mousedown)="createContextFromTerm(searchTerm().trim(), null)"
                  >
                    <ng-icon name="lucidePlus" class="text-xs shrink-0"></ng-icon>
                    <span class="min-w-0 flex-1 truncate font-medium">{{
                      searchTerm().trim()
                    }}</span>
                    <span class="shrink-0 text-[10px] text-muted-foreground">Workspace</span>
                  </button>
                }

                @if (matchedProjects().length > 0) {
                  <div
                    class="px-2 pt-1 pb-1 text-[10px] uppercase tracking-[0.08em] text-muted-foreground"
                  >
                    Project for context
                  </div>
                }
                @for (project of matchedProjects(); track project.id) {
                  <button
                    type="button"
                    class="w-full flex items-center justify-between gap-2 text-left px-2 py-1.5 rounded-sm text-xs hover:bg-muted"
                    [class.bg-muted]="
                      activeSuggestionIndex() === projectSuggestionIndex(project.id)
                    "
                    [class.text-foreground]="
                      activeSuggestionIndex() === projectSuggestionIndex(project.id)
                    "
                    [class.text-muted-foreground]="
                      activeSuggestionIndex() !== projectSuggestionIndex(project.id)
                    "
                    (mouseenter)="setActiveSuggestionIndex(projectSuggestionIndex(project.id))"
                    (mousedown)="selectProjectSuggestion($event, project)"
                  >
                    <span class="min-w-0 flex items-center gap-1.5">
                      <ng-icon name="lucideFolder" class="shrink-0 text-[11px]"></ng-icon>
                      <span class="truncate">{{ project.name }}</span>
                    </span>
                    <span class="shrink-0 text-[10px] text-muted-foreground">Project</span>
                  </button>
                }

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
                      @if (context.project; as project) {
                        <ctx-search-project-badge [project]="project" />
                      }
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
                      @if (context.project; as project) {
                        <ctx-search-project-badge [project]="project" />
                      }
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

                @if (archivedMatchedContexts().length > 0) {
                  <div class="my-1 border-t border-border/70"></div>
                  <div
                    class="px-2 pt-1 pb-1 text-[10px] uppercase tracking-[0.08em] text-muted-foreground"
                  >
                    Archived
                  </div>
                }
                @for (context of archivedMatchedContexts(); track context.id) {
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
                    (mousedown)="openContext(context)"
                  >
                    <span class="truncate">{{ context.name }}</span>
                    <span
                      class="shrink-0 text-[10px] text-muted-foreground/80 flex items-center gap-2"
                    >
                      @if (context.project; as project) {
                        <ctx-search-project-badge [project]="project" />
                      }
                      <span class="rounded-sm border px-1 py-0.5 leading-none">Archived</span>
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
  private inteverlaQueries = inject(IntervalQueries);
  private contextMutations = inject(ContextMutations);
  private projectQueries = inject(ProjectQueries);
  private contextService = inject(ContextService);
  private settingsQueries = inject(SettingsQueries);
  private timeZone = inject(TimeZoneService);
  private readonly store = inject(Store);
  readonly activeWorkspaceId = this.store.selectSignal(WorkspaceState.selectedWorkspaceId);
  private readonly selectedProjectState = this.store.selectSignal(WorkspaceState.selectedProjectId);
  readonly selectedProjectId = computed(() => this.selectedProjectState() ?? '');
  private router = inject(Router);
  today = computed(() => this.timeZone.today());

  listContextsQuery = injectQuery(() => this.contextQueries.list(this.activeWorkspaceId(), true));
  allProjectsQuery = injectQuery(() => this.projectQueries.all(this.activeWorkspaceId() ?? ''));
  selectedProjectQuery = injectQuery(() => this.projectQueries.get(this.selectedProjectId()));
  settingsQuery = injectQuery(() => this.settingsQueries.settings());
  switchContextMutation = injectMutation(() => this.contextMutations.switch());
  freeContextMutation = injectMutation(() => this.contextMutations.free());
  activeContextQuery = injectQuery(() => this.contextQueries.active());
  private routedSelectedDate = toSignal(
    this.router.events.pipe(
      filter((event): event is NavigationEnd => event instanceof NavigationEnd),
      startWith(null),
      map(() => this.resolveSelectedDate()),
    ),
    { initialValue: this.today() },
  );
  selectedDate = computed(() => this.routedSelectedDate() ?? this.today());
  dayStatsQuery = injectQuery(() =>
    this.inteverlaQueries.dayStats(
      this.activeWorkspaceId(),
      this.selectedDate(),
      this.timeZone.effectiveTimeZone(),
    ),
  );
  activeContextName = computed(() => this.activeContextQuery.data()?.name ?? '');
  weekStartsOn = computed(() => (this.settingsQuery.data()?.[firstDayKey] === 'Sunday' ? 0 : 1));
  daySectionLabel = computed(() =>
    this.selectedDate() === this.today() ? 'Today' : this.selectedDate(),
  );

  readonly searchTerm = signal<string>('');
  readonly searchFocused = signal<boolean>(false);
  readonly mobileSearchOpen = signal<boolean>(false);
  readonly activeSuggestionIndex = signal<number>(-1);
  readonly contexts = computed<readonly Context[]>(() => this.listContextsQuery.data() ?? []);
  readonly projects = computed<readonly Project[]>(() => this.allProjectsQuery.data() ?? []);
  readonly selectedProject = computed<ProjectMetadata | null>(() => {
    const projectId = this.selectedProjectId();
    if (!projectId) {
      return null;
    }
    return {
      id: projectId,
      name: this.selectedProjectQuery.data()?.name ?? 'Current project',
    };
  });
  readonly projectPicker = computed(() => parseProjectPicker(this.searchTerm()));
  readonly projectPickerMode = computed(() => this.projectPicker() !== null);
  readonly projectCreationSuggestion = computed<ProjectMetadata | null>(() =>
    this.projectPickerMode() ? null : this.selectedProject(),
  );
  readonly filteredContexts = computed<readonly Context[]>(() => {
    if (this.projectPickerMode()) {
      return [];
    }
    const term = this.searchTerm().trim().toLowerCase();
    if (!term) {
      return [];
    }
    return this.contexts().filter((context) => context.name.toLowerCase().includes(term));
  });
  readonly matchedProjects = computed<readonly Project[]>(() => {
    const picker = this.projectPicker();
    if (!picker) {
      return [];
    }
    const term = picker.projectQuery.toLowerCase();
    return this.projects().filter((project) => project.name.toLowerCase().includes(term));
  });
  readonly activeFilteredContexts = computed<readonly Context[]>(() =>
    this.filteredContexts().filter((context) => !context.archived),
  );
  readonly archivedMatchedContexts = computed<readonly Context[]>(() =>
    this.filteredContexts().filter((context) => context.archived),
  );
  readonly usedContextIdsForDay = computed(
    () => new Set(this.dayStatsQuery.data()?.contextStats.map((stats) => stats.contextId) ?? []),
  );
  readonly dayMatchedContexts = computed<readonly Context[]>(() =>
    this.activeFilteredContexts().filter((context) => this.usedContextIdsForDay().has(context.id)),
  );
  readonly otherMatchedContexts = computed<readonly Context[]>(() =>
    this.activeFilteredContexts().filter((context) => !this.usedContextIdsForDay().has(context.id)),
  );
  readonly suggestionContexts = computed<readonly Context[]>(() => [
    ...this.dayMatchedContexts(),
    ...this.otherMatchedContexts(),
    ...this.archivedMatchedContexts(),
  ]);
  readonly creationSuggestionCount = computed(() => {
    if (this.projectPickerMode()) {
      return 0;
    }
    return this.selectedProject() ? 2 : 1;
  });
  readonly suggestionCount = computed<number>(() =>
    this.searchTerm().trim().length > 0
      ? this.suggestionContexts().length +
        this.matchedProjects().length +
        this.creationSuggestionCount()
      : 0,
  );
  readonly statsByContextId = toSignal(
    toObservable(
      computed(() => ({
        contexts: this.filteredContexts(),
        date: this.selectedDate(),
        timeZone: this.timeZone.effectiveTimeZone(),
      })),
    ).pipe(
      switchMap(({ contexts, date, timeZone }) => {
        if (contexts.length === 0) {
          return of({} as Record<string, ContextStats>);
        }

        return forkJoin(
          contexts.map((context) =>
            this.contextService.getStats(context.id, date, timeZone).pipe(
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
    if (context.archived) {
      this.openContext(context);
      return;
    }
    this.resetSearchUi();
    this.switchContextMutation.mutate(context);
  }

  openContext(context: Context): void {
    this.resetSearchUi();
    this.router.navigate(['/context', context.id]);
  }

  selectProjectSuggestion(event: MouseEvent, project: Project): void {
    event.preventDefault();
    const picker = this.projectPicker();
    if (!picker) {
      return;
    }
    this.createContextFromTerm(picker.contextName, { id: project.id, name: project.name });
  }

  createContextFromTerm(
    term: string,
    project: ProjectMetadata | null = this.selectedProject(),
  ): void {
    const normalizedTerm = term.trim();
    if (!normalizedTerm) {
      return;
    }
    this.resetSearchUi();
    this.switchContextMutation.mutate({
      id: '',
      name: normalizedTerm,
      project: project ?? undefined,
    });
  }

  onSearchKeydown(event: KeyboardEvent): void {
    if (event.key === 'ArrowDown') {
      if (!this.showSuggestions()) {
        return;
      }
      event.preventDefault();
      const suggestionsLength = this.suggestionCount();
      if (suggestionsLength === 0) {
        return;
      }
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
      if (suggestionsLength === 0) {
        return;
      }
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

    const activeIndex = Math.max(0, this.activeSuggestionIndex());
    if (this.projectPickerMode()) {
      const project = this.matchedProjects()[activeIndex];
      if (project) {
        const picker = this.projectPicker();
        if (picker) {
          this.createContextFromTerm(picker.contextName, {
            id: project.id,
            name: project.name,
          });
        }
      }
      return;
    }

    const selectedProject = this.selectedProject();
    if (this.showSuggestions() && activeIndex < this.creationSuggestionCount()) {
      this.createContextFromTerm(
        term,
        selectedProject && activeIndex === 0 ? selectedProject : null,
      );
      return;
    }

    if (this.showSuggestions()) {
      const resultIndex = activeIndex - this.creationSuggestionCount();
      const selectedContext = this.suggestionContexts()[resultIndex];
      if (selectedContext) {
        this.selectContext(selectedContext);
        return;
      }
    }

    this.createContextFromTerm(term);
  }

  suggestionIndex(contextId: string): number {
    return (
      this.suggestionContexts().findIndex((context) => context.id === contextId) +
      this.creationSuggestionCount() +
      this.matchedProjects().length
    );
  }

  projectSuggestionIndex(projectId: string): number {
    return (
      this.matchedProjects().findIndex((project) => project.id === projectId) +
      this.creationSuggestionCount()
    );
  }

  contextTodayDuration(contextId: string): string | null {
    const stats = this.statsByContextId()[contextId];
    return this.formatDuration(stats?.duration ?? 0);
  }

  contextTotalDuration(contextId: string): string | null {
    const stats = this.statsByContextId()[contextId];
    return this.formatDuration(stats?.totalDuration ?? 0);
  }

  private resolveSelectedDate(): string | null {
    const dayMatch = this.router.url.match(/\/day\/(\d{4}-\d{2}-\d{2})/);
    return dayMatch?.[1] ?? null;
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
    this.searchTerm.set('');
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

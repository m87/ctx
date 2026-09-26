import { Component, computed, effect, inject, signal } from '@angular/core';
import { RouterLink } from '@angular/router';
import { NgIcon, provideIcons } from '@ng-icons/core';
import {
  lucideArrowLeft,
  lucideCheck,
  lucideEye,
  lucideGanttChart,
  lucidePencil,
  lucidePlus,
  lucideSettings2,
  lucideTrash2,
  lucideX,
} from '@ng-icons/lucide';
import { BrnAlertDialogImports } from '@spartan-ng/brain/alert-dialog';
import { BrnDialogImports } from '@spartan-ng/brain/dialog';
import { HlmAlertDialogImports } from '@spartan-ng/helm/alert-dialog';
import { HlmButtonImports } from '@spartan-ng/helm/button';
import { HlmDialogImports } from '@spartan-ng/helm/dialog';
import { HlmInputImports } from '@spartan-ng/helm/input';
import { HlmLabelImports } from '@spartan-ng/helm/label';
import { HlmSkeletonImports } from '@spartan-ng/helm/skeleton';
import { HlmTextareaImports } from '@spartan-ng/helm/textarea';
import { injectMutation, injectQuery } from '@tanstack/angular-query-experimental';
import { Store } from '@ngxs/store';
import { ContextQueries } from '../../api/context/context.queries';
import { DashboardMutations } from '../../api/dashboard/dashboard.mutations';
import { DashboardQueries } from '../../api/dashboard/dashboard.queries';
import {
  CreateDashboardInput,
  Dashboard,
  DashboardType,
} from '../../api/dashboard/dashboard.service';
import { ProjectQueries } from '../../api/project/project.queries';
import { QueryQueries } from '../../api/query/query.queries';
import { ContextQueryResult } from '../../api/query/query.service';
import { SavedQueryQueries } from '../../api/query/saved-query.queries';
import { SavedQuery } from '../../api/query/saved-query.service';
import { ContextListComponent } from '../context/context-list.component';
import {
  addTextWidget,
  deleteTextWidget,
  normalizeDashboardDefinition,
  moveTextWidget,
  resizeTextWidget,
  TextWidgetDefinition,
  TextWidgetProperties,
  updateTextWidget,
  updateTextWidgetLayout,
  DashboardWidgetLayout,
} from './dashboard/dashboard-definition';
import { DashboardCanvasComponent } from './dashboard/dashboard-canvas.component';
import { GridPoint } from './dashboard/dashboard-layout';
import { WidgetPaletteComponent } from './dashboard/widget-palette.component';
import { TextWidgetPropertiesComponent } from './dashboard/widgets/text-widget-properties.component';
import { contextQueryResultAsListItems } from '../query/query-summary.component';
import { QueryErrorStateComponent } from '../shared/query-error-state.component';
import { SearchDropdownSelectComponent } from '../shared/search-dropdown-select.component';
import { SearchSelectComponent, SearchSelectOption } from '../shared/search-select.component';
import { SidebarWorkspaceSelectComponent } from '../sidebar/sidebar-workspace-select.component';
import { WorkspaceState } from '../sidebar/workspace.state';
import { colorHash, durationAsHM } from '../utils';

export interface DashboardWidgetPlacement {
  y: number;
  height: number;
}

type DashboardDialogMode = 'create' | 'properties';

export const DASHBOARD_COLUMN_COUNT = 16;
export const DASHBOARD_MINIMUM_ROW_COUNT = 32;

export const DASHBOARD_TYPE_OPTIONS: readonly SearchSelectOption[] = [
  { value: 'custom', label: 'Custom' },
  { value: 'workspace', label: 'Workspace insight' },
  { value: 'project', label: 'Project insight' },
  { value: 'context', label: 'Context insight' },
  { value: 'daily', label: 'Daily insight' },
];

export function dashboardTypeRequiresTarget(type: DashboardType): boolean {
  return type === 'project' || type === 'context';
}

export function isDashboardType(value: string): value is DashboardType {
  return (
    value === 'workspace' ||
    value === 'project' ||
    value === 'context' ||
    value === 'daily' ||
    value === 'custom'
  );
}

export function dashboardTargetId(
  type: DashboardType,
  workspaceId: string,
  selectedTargetId: string,
): string | undefined {
  if (type === 'workspace' || type === 'daily') {
    return workspaceId;
  }
  if (dashboardTypeRequiresTarget(type)) {
    return selectedTargetId || undefined;
  }
  return undefined;
}

export function dashboardsAsOptions(dashboards: readonly Dashboard[]): SearchSelectOption[] {
  return dashboards.map((dashboard) => ({
    value: dashboard.id,
    label: dashboard.name,
    description:
      DASHBOARD_TYPE_OPTIONS.find((option) => option.value === dashboard.type)?.label ??
      dashboard.type,
    keywords: [dashboard.type],
  }));
}

export function dashboardGridRowCount(widgets: readonly DashboardWidgetPlacement[]): number {
  return widgets.reduce(
    (rowCount, widget) => Math.max(rowCount, widget.y + widget.height),
    DASHBOARD_MINIMUM_ROW_COUNT,
  );
}

export function savedQueriesAsOptions(queries: readonly SavedQuery[]): SearchSelectOption[] {
  return queries.map((query) => ({
    value: query.id,
    label: query.name,
    description: query.query,
    keywords: [query.query],
  }));
}

export function selectedSavedQueryText(
  queries: readonly SavedQuery[],
  selectedQueryId: string,
): string {
  return queries.find((query) => query.id === selectedQueryId)?.query ?? '';
}

export function queryPreviewSummary(result: ContextQueryResult | undefined): string {
  if (!result) {
    return '';
  }

  const contextCount = result.contexts.length;
  const contextLabel = contextCount === 1 ? 'context' : 'contexts';
  const duration = durationAsHM(result.totalDuration).trim() || '0m';
  return `${contextCount} ${contextLabel} · ${duration} · ${result.totalSessions} sessions`;
}

@Component({
  selector: 'ctx-editor',
  imports: [
    BrnAlertDialogImports,
    BrnDialogImports,
    HlmAlertDialogImports,
    HlmButtonImports,
    HlmDialogImports,
    HlmInputImports,
    HlmLabelImports,
    HlmSkeletonImports,
    HlmTextareaImports,
    DashboardCanvasComponent,
    WidgetPaletteComponent,
    TextWidgetPropertiesComponent,
    NgIcon,
    ContextListComponent,
    QueryErrorStateComponent,
    RouterLink,
    SearchDropdownSelectComponent,
    SearchSelectComponent,
    SidebarWorkspaceSelectComponent,
  ],
  providers: [
    provideIcons({
      lucideArrowLeft,
      lucideCheck,
      lucideEye,
      lucideGanttChart,
      lucidePencil,
      lucidePlus,
      lucideSettings2,
      lucideTrash2,
      lucideX,
    }),
  ],
  template: `
    <div class="fixed inset-0 z-50 flex h-dvh min-h-0 flex-col bg-background text-foreground">
      <header
        class="flex min-h-12 shrink-0 flex-wrap items-center gap-y-2 border-b bg-card/70 px-3 py-2 sm:flex-nowrap"
      >
        <div class="flex min-w-0 items-center gap-2">
          <a
            routerLink="/"
            class="flex size-8 shrink-0 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-muted hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
            aria-label="Back to application"
            title="Back to application"
          >
            <ng-icon name="lucideArrowLeft"></ng-icon>
          </a>

          <div class="h-5 w-px bg-border"></div>

          <div class="flex min-w-0 items-center gap-2 pl-1">
            <ng-icon name="lucideGanttChart" class="shrink-0 text-primary"></ng-icon>
            <span class="font-semibold tracking-tight text-primary">Ctx</span>
            <span class="text-muted-foreground">/</span>
            <span class="truncate text-sm font-medium">Editor</span>
          </div>

          <div class="ml-1 h-5 w-px bg-border"></div>

          <ctx-sidebar-workspace-select
            class="ml-1 hidden w-48 md:block"
          ></ctx-sidebar-workspace-select>
        </div>

        <div class="ml-auto flex min-w-0 items-center gap-2 sm:shrink-0 sm:pl-4">
          <ctx-search-dropdown-select
            class="w-36 sm:w-56"
            inputId="editor-dashboard"
            ariaLabel="Dashboard"
            placeholder="Select dashboard…"
            searchPlaceholder="Search dashboards…"
            [emptyText]="dashboardListEmptyText()"
            [disabled]="dashboardEditing() || dashboardsQuery.isLoading()"
            [options]="dashboardOptions()"
            [value]="activeDashboardId()"
            (selectionChange)="selectDashboard($event)"
          ></ctx-search-dropdown-select>

          <button
            hlmBtn
            type="button"
            variant="outline"
            size="icon-sm"
            aria-label="Create dashboard"
            title="Create dashboard"
            [disabled]="!selectedWorkspaceId() || dashboardEditing()"
            (click)="openCreateDashboardDialog()"
          >
            <ng-icon name="lucidePlus"></ng-icon>
          </button>

          @if (activeDashboard()) {
            @if (dashboardEditing()) {
              <div class="ui-action-bar">
                <button
                  hlmBtn
                  type="button"
                  variant="ghost"
                  size="icon-sm"
                  aria-label="Edit dashboard properties"
                  title="Edit dashboard properties"
                  [disabled]="updateDashboardMutation.isPending()"
                  (click)="openDashboardPropertiesDialog()"
                >
                  <ng-icon name="lucideSettings2"></ng-icon>
                </button>
                <button
                  hlmBtn
                  type="button"
                  variant="ghost"
                  size="icon-sm"
                  class="text-destructive"
                  aria-label="Delete dashboard"
                  title="Delete dashboard"
                  [disabled]="updateDashboardMutation.isPending()"
                  (click)="requestDashboardDelete()"
                >
                  <ng-icon name="lucideTrash2"></ng-icon>
                </button>
                <button
                  hlmBtn
                  type="button"
                  variant="ghost"
                  size="icon-sm"
                  aria-label="Cancel dashboard changes"
                  title="Cancel dashboard changes"
                  [disabled]="updateDashboardMutation.isPending()"
                  (click)="cancelDashboardEditing()"
                >
                  <ng-icon name="lucideX"></ng-icon>
                </button>
                <button
                  hlmBtn
                  type="button"
                  variant="outline"
                  size="icon-sm"
                  class="text-success"
                  aria-label="Save dashboard changes"
                  title="Save dashboard changes"
                  [disabled]="updateDashboardMutation.isPending()"
                  (click)="confirmDashboardEditing()"
                >
                  <ng-icon name="lucideCheck"></ng-icon>
                </button>
              </div>
            } @else {
              <button
                hlmBtn
                type="button"
                variant="ghost"
                size="icon-sm"
                aria-label="Enter dashboard edit mode"
                title="Edit dashboard"
                (click)="startDashboardEditing()"
              >
                <ng-icon name="lucidePencil"></ng-icon>
              </button>
            }
          }

          @if (updateDashboardMutation.isError()) {
            <span class="hidden text-xs text-destructive xl:inline" role="alert">Save failed</span>
          }
        </div>
      </header>

      <div class="flex min-h-0 flex-1">
        <aside
          class="hidden w-60 shrink-0 flex-col border-r bg-sidebar md:flex"
          aria-label="Editor tools"
        >
          <section class="flex min-h-0 flex-1 flex-col border-b" aria-label="Query">
            <div class="flex items-center justify-between gap-2 border-b px-3 py-2">
              <div class="text-meta font-semibold uppercase tracking-label text-muted-foreground">
                Query
              </div>

              <hlm-dialog>
                <button
                  hlmBtn
                  hlmDialogTrigger
                  type="button"
                  variant="ghost"
                  size="icon-sm"
                  class="text-muted-foreground"
                  aria-label="Preview query"
                  title="Preview query"
                  [disabled]="!selectedWorkspaceId()"
                  (click)="loadQueryPreview()"
                >
                  <ng-icon name="lucideEye"></ng-icon>
                </button>

                <hlm-dialog-content
                  class="editor-preview-dialog flex h-[calc(100dvh-2rem)] max-h-[calc(100dvh-2rem)] flex-col gap-0 overflow-hidden p-0 sm:h-[min(44rem,calc(100dvh-4rem))]"
                  *brnDialogContent
                >
                  <div hlmDialogHeader class="shrink-0 border-b p-5 pr-12">
                    <h2 hlmDialogTitle>Query preview</h2>
                    <p hlmDialogDescription>
                      {{ previewSummary() || 'Contexts returned for the selected query.' }}
                    </p>
                  </div>
                  <div class="min-h-0 flex-1 overflow-y-auto p-5">
                    @if (previewShowError()) {
                      <ctx-query-error-state
                        [error]="queryPreview.error()"
                        [paused]="queryPreview.isPaused()"
                        resourceName="query preview"
                        [retrying]="queryPreview.isFetching()"
                        (retry)="loadQueryPreview()"
                      ></ctx-query-error-state>
                    } @else if (queryPreview.isLoading()) {
                      <div class="flex flex-col gap-2" role="status" aria-label="Loading contexts">
                        <span class="sr-only">Loading contexts</span>
                        @for (item of previewSkeletonItems; track item) {
                          <div class="rounded-lg border bg-card p-3">
                            <div class="mb-3 flex items-center gap-2">
                              <hlm-skeleton class="size-2 shrink-0"></hlm-skeleton>
                              <hlm-skeleton class="h-3.5 w-2/5"></hlm-skeleton>
                              <hlm-skeleton class="ml-auto h-3 w-12"></hlm-skeleton>
                            </div>
                            <hlm-skeleton class="mb-2 h-1.5 w-full"></hlm-skeleton>
                            <hlm-skeleton class="h-2.5 w-24"></hlm-skeleton>
                          </div>
                        }
                      </div>
                    } @else {
                      <ctx-context-list
                        [items]="previewContexts()"
                        emptyMessage="No contexts match this query."
                      ></ctx-context-list>
                    }
                  </div>
                </hlm-dialog-content>
              </hlm-dialog>
            </div>

            <div class="flex flex-1 items-start gap-3 overflow-y-auto bg-background/40 p-3">
              <div class="font-mono text-xs text-muted-foreground/50" aria-hidden="true">1</div>
              <div
                class="min-w-0 whitespace-pre-wrap break-words font-mono text-xs text-foreground"
              >
                @if (dashboardEditing()) {
                  <textarea
                    hlmTextarea
                    class="dashboard-field-control min-h-32 w-full resize-y font-mono text-xs"
                    aria-label="Dashboard query"
                    [value]="dashboardQueryText()"
                    (input)="updateDashboardQuery($event)"
                  ></textarea>
                } @else {
                  {{ queryText() }}
                }
              </div>
            </div>
          </section>

          <section class="flex min-h-0 flex-1 flex-col" aria-labelledby="saved-queries-title">
            <div class="border-b px-3 py-3">
              <div id="saved-queries-title" class="ui-section-label">Saved queries</div>
            </div>

            <ctx-search-select
              class="min-h-0"
              inputId="editor-saved-query"
              ariaLabel="Saved queries"
              searchPlaceholder="Search saved queries…"
              [emptyText]="savedQueryEmptyText()"
              [disabled]="savedQueriesQuery.isLoading()"
              [options]="savedQueryOptions()"
              [value]="selectedSavedQueryId()"
              [embedded]="true"
              (selectionChange)="selectSavedQuery($event)"
            ></ctx-search-select>
          </section>
        </aside>

        <main
          class="flex min-w-0 flex-1 flex-col items-center overflow-auto bg-muted p-3 dark:bg-background sm:p-6"
        >
          @if (dashboardDefinitionError()) {
            <div class="ui-notice ui-notice-destructive mb-3" role="alert">
              {{ dashboardDefinitionError() }}
            </div>
          }
          <ctx-dashboard-canvas
            #dashboardCanvas
            [widgets]="dashboardWidgets()"
            [editing]="dashboardEditing() && !updateDashboardMutation.isPending()"
            [selectedWidgetId]="selectedWidgetId()"
            (widgetSelected)="selectWidget($event)"
            (widgetAdded)="addWidget($event)"
            (widgetMoved)="moveWidgetTo($event.id, $event.layout)"
            (widgetResized)="applyWidgetLayout($event.id, $event.layout)"
            (widgetKeydown)="handleWidgetKeydown($event.event, $event.widget)"
          ></ctx-dashboard-canvas>
          @if (dashboardEditing()) {
            <section
              class="mt-4 w-full max-w-[1000px] rounded-xl border border-border/65 bg-sidebar lg:hidden"
              aria-label="Text widget tools"
            >
              <ctx-widget-palette
                [disabled]="updateDashboardMutation.isPending()"
                [draggable]="false"
                (widgetAdded)="addWidget()"
              ></ctx-widget-palette>
              @if (selectedWidget(); as widget) {
                <ctx-text-widget-properties
                  [widget]="widget"
                  idSuffix="-mobile"
                  [disabled]="updateDashboardMutation.isPending()"
                  (queryChange)="updateWidgetQuery(widget.id, $event)"
                  (propertiesChange)="updateWidgetProperties(widget.id, $event)"
                  (deleteRequested)="deleteWidget(widget.id)"
                ></ctx-text-widget-properties>
              }
            </section>
          }
        </main>

        <aside
          class="hidden w-72 shrink-0 flex-col border-l bg-sidebar lg:flex"
          aria-label="Widget tools"
        >
          <section class="flex shrink-0 flex-col border-b" aria-label="Widgets">
            <div class="border-b px-3 py-3">
              <div class="text-meta font-semibold uppercase tracking-label text-muted-foreground">
                Widgets
              </div>
            </div>

            @if (dashboardEditing()) {
              <ctx-widget-palette
                [overCanvas]="dashboardCanvas.palettePreviewVisible()"
                [disabled]="updateDashboardMutation.isPending()"
                (widgetAdded)="addWidget()"
                (dragStarted)="dashboardCanvas.startPaletteDrag()"
                (dragMoved)="dashboardCanvas.movePaletteDrag($event)"
                (dragEnded)="dashboardCanvas.finishPaletteDrag($event)"
              ></ctx-widget-palette>
            } @else {
              <div class="flex flex-1 items-center justify-center p-5">
                <div class="text-center text-xs text-muted-foreground">
                  {{
                    dashboardWidgets().length
                      ? 'Widget catalog is available in edit mode.'
                      : 'No widgets yet'
                  }}
                </div>
              </div>
            }
          </section>

          <section class="flex min-h-0 flex-1 flex-col" aria-label="Properties">
            <div class="border-b px-3 py-3">
              <div class="text-meta font-semibold uppercase tracking-label text-muted-foreground">
                Properties
              </div>
            </div>

            @if (selectedWidget(); as widget) {
              <ctx-text-widget-properties
                [widget]="widget"
                [disabled]="updateDashboardMutation.isPending()"
                (queryChange)="updateWidgetQuery(widget.id, $event)"
                (propertiesChange)="updateWidgetProperties(widget.id, $event)"
                (deleteRequested)="deleteWidget(widget.id)"
              ></ctx-text-widget-properties>
            } @else {
              <div class="flex flex-1 items-center justify-center p-5">
                <div class="text-center text-xs text-muted-foreground">Nothing selected</div>
              </div>
            }
          </section>
        </aside>
      </div>

      <hlm-alert-dialog
        [state]="pendingWidgetDelete() ? 'open' : 'closed'"
        (closed)="pendingWidgetDelete.set(null)"
      >
        <hlm-alert-dialog-content *brnAlertDialogContent>
          <hlm-alert-dialog-header>
            <h3 hlmAlertDialogTitle>Delete this text widget?</h3>
            <p hlmAlertDialogDescription>
              The widget will be removed from this draft. Save the dashboard to keep the change.
            </p>
          </hlm-alert-dialog-header>
          <hlm-alert-dialog-footer>
            <button hlmBtn type="button" variant="outline" (click)="pendingWidgetDelete.set(null)">
              Cancel
            </button>
            <button hlmBtn type="button" variant="destructive" (click)="confirmWidgetDelete()">
              Delete widget
            </button>
          </hlm-alert-dialog-footer>
        </hlm-alert-dialog-content>
      </hlm-alert-dialog>

      <hlm-dialog
        [state]="dashboardDialogMode() ? 'open' : 'closed'"
        [disableClose]="dashboardMutationPending()"
        (closed)="closeDashboardDialog()"
      >
        <hlm-dialog-content *brnDialogContent>
          <div hlmDialogHeader>
            <h2 hlmDialogTitle>{{ dashboardDialogTitle() }}</h2>
            <p hlmDialogDescription>
              @if (dashboardDialogContentMode() === 'properties') {
                Change the dashboard name and where it is used. Confirm all editor changes with the
                checkmark in the toolbar.
              } @else {
                Choose where this dashboard is used. Custom dashboards are stored for this
                workspace.
              }
            </p>
          </div>

          <div class="grid gap-4 py-2">
            <div class="ui-field">
              <label hlmLabel for="dashboard-name">Name</label>
              <input
                hlmInput
                class="dashboard-field-control h-9"
                id="dashboard-name"
                type="text"
                autocomplete="off"
                placeholder="Dashboard name"
                [value]="dashboardFormName()"
                [disabled]="dashboardMutationPending()"
                (input)="setDashboardFormName($event)"
              />
            </div>

            <div class="ui-field">
              <label hlmLabel for="dashboard-type">Type</label>
              <ctx-search-dropdown-select
                inputId="dashboard-type"
                ariaLabel="Dashboard type"
                panelWidth="100%"
                [searchable]="false"
                [disabled]="dashboardMutationPending()"
                [options]="dashboardTypeOptions"
                [value]="dashboardFormType()"
                (selectionChange)="setDashboardFormType($event)"
              ></ctx-search-dropdown-select>
            </div>

            @if (dashboardFormRequiresTarget()) {
              <div class="ui-field">
                <label hlmLabel for="dashboard-target">Target</label>
                <ctx-search-dropdown-select
                  inputId="dashboard-target"
                  [ariaLabel]="dashboardTargetAriaLabel()"
                  [placeholder]="dashboardTargetPlaceholder()"
                  [searchPlaceholder]="dashboardTargetSearchPlaceholder()"
                  [emptyText]="dashboardTargetEmptyText()"
                  [disabled]="dashboardMutationPending()"
                  [options]="dashboardTargetOptions()"
                  [value]="dashboardFormTargetId()"
                  (selectionChange)="dashboardFormTargetId.set($event)"
                ></ctx-search-dropdown-select>
              </div>
            } @else {
              <div class="rounded-lg border bg-muted/25 px-3 py-2.5 text-xs text-muted-foreground">
                {{ dashboardTypeDescription() }}
              </div>
            }

            @if (dashboardMutationError()) {
              <p class="text-xs text-destructive" role="alert">
                Could not save the dashboard. Check the selected type and target, then try again.
              </p>
            }
          </div>

          <div hlmDialogFooter>
            <button
              hlmBtn
              type="button"
              variant="ghost"
              [disabled]="dashboardMutationPending()"
              (click)="closeDashboardDialog()"
            >
              Cancel
            </button>
            <button
              hlmBtn
              type="button"
              variant="outline"
              [disabled]="!dashboardFormValid() || dashboardMutationPending()"
              (click)="saveDashboard()"
            >
              @if (dashboardMutationPending()) {
                Saving…
              } @else if (dashboardDialogContentMode() === 'properties') {
                Apply properties
              } @else {
                Create dashboard
              }
            </button>
          </div>
        </hlm-dialog-content>
      </hlm-dialog>

      <hlm-alert-dialog
        [state]="pendingDashboardDelete() ? 'open' : 'closed'"
        (closed)="pendingDashboardDelete.set(null)"
      >
        <hlm-alert-dialog-content *brnAlertDialogContent>
          <hlm-alert-dialog-header>
            <h3 hlmAlertDialogTitle>Delete this dashboard?</h3>
            <p hlmAlertDialogDescription>
              The dashboard definition and its saved layout will be permanently removed.
            </p>
          </hlm-alert-dialog-header>
          <hlm-alert-dialog-footer>
            <button
              hlmAlertDialogCancel
              [disabled]="deleteDashboardMutation.isPending()"
              (click)="pendingDashboardDelete.set(null)"
            >
              Cancel
            </button>
            <button
              hlmAlertDialogAction
              variant="destructive"
              [disabled]="deleteDashboardMutation.isPending()"
              (click)="confirmDeleteDashboard()"
            >
              {{ deleteDashboardMutation.isPending() ? 'Deleting…' : 'Delete' }}
            </button>
          </hlm-alert-dialog-footer>
        </hlm-alert-dialog-content>
      </hlm-alert-dialog>
    </div>
  `,
  styles: `
    :host {
      display: block;
    }

    .editor-preview-dialog {
      width: calc(100vw - 2rem);
      max-width: 72rem;
    }
  `,
})
export class EditorComponent {
  readonly previewSkeletonItems = [0, 1, 2];
  readonly selectedWidgetId = signal<string | null>(null);
  readonly dashboardWidgets = computed<readonly TextWidgetDefinition[]>(() => {
    const dashboard = this.dashboardDraft() ?? this.activeDashboard();
    if (!dashboard) return [];
    try {
      return normalizeDashboardDefinition(dashboard.definition).widgets;
    } catch {
      return [];
    }
  });
  private readonly projectQueries = inject(ProjectQueries);
  private readonly contextQueries = inject(ContextQueries);
  private readonly dashboardQueries = inject(DashboardQueries);
  private readonly dashboardMutations = inject(DashboardMutations);
  private readonly queryQueries = inject(QueryQueries);
  private readonly savedQueryQueries = inject(SavedQueryQueries);
  private readonly store = inject(Store);

  readonly dashboardTypeOptions = DASHBOARD_TYPE_OPTIONS;
  readonly activeDashboardId = signal('');
  readonly dashboardDraft = signal<Dashboard | null>(null);
  readonly dashboardDialogMode = signal<DashboardDialogMode | null>(null);
  readonly dashboardDialogContentMode = signal<DashboardDialogMode>('create');
  readonly dashboardFormName = signal('');
  readonly dashboardFormType = signal<DashboardType>('custom');
  readonly dashboardFormTargetId = signal('');
  readonly pendingDashboardDelete = signal<Dashboard | null>(null);
  readonly pendingWidgetDelete = signal<string | null>(null);
  readonly selectedSavedQueryId = signal('');
  readonly selectedWorkspaceId = this.store.selectSignal(WorkspaceState.selectedWorkspaceId);

  readonly projectsQuery = injectQuery(() =>
    this.projectQueries.all(this.selectedWorkspaceId() ?? ''),
  );
  readonly contextsQuery = injectQuery(() => this.contextQueries.list(this.selectedWorkspaceId()));
  readonly savedQueriesQuery = injectQuery(() =>
    this.savedQueryQueries.list(this.selectedWorkspaceId()),
  );
  readonly dashboardsQuery = injectQuery(() =>
    this.dashboardQueries.list(this.selectedWorkspaceId()),
  );
  readonly createDashboardMutation = injectMutation(() => this.dashboardMutations.create());
  readonly updateDashboardMutation = injectMutation(() => this.dashboardMutations.update());
  readonly deleteDashboardMutation = injectMutation(() => this.dashboardMutations.delete());

  readonly dashboards = computed(() => this.dashboardsQuery.data() ?? []);
  readonly activeDashboard = computed(
    () => this.dashboards().find((dashboard) => dashboard.id === this.activeDashboardId()) ?? null,
  );
  readonly dashboardOptions = computed<SearchSelectOption[]>(() =>
    dashboardsAsOptions(this.dashboards()),
  );
  readonly dashboardEditing = computed(() => this.dashboardDraft() !== null);
  readonly dashboardDefinitionError = computed(() => {
    const dashboard = this.activeDashboard();
    if (!dashboard) return '';
    try {
      normalizeDashboardDefinition(dashboard.definition);
      return '';
    } catch (error) {
      return error instanceof Error ? error.message : 'Dashboard definition is invalid.';
    }
  });
  readonly selectedWidget = computed(
    () => this.dashboardWidgets().find((widget) => widget.id === this.selectedWidgetId()) ?? null,
  );
  readonly dashboardQueryText = computed(() => {
    const dashboard = this.dashboardDraft();
    if (!dashboard) return this.queryText();
    try {
      return normalizeDashboardDefinition(dashboard.definition).query;
    } catch {
      return '';
    }
  });

  private readonly resetEditorOnWorkspaceChange = effect(() => {
    this.selectedWorkspaceId();
    this.activeDashboardId.set('');
    this.dashboardDraft.set(null);
    this.selectedSavedQueryId.set('');
    this.dashboardDialogMode.set(null);
    this.pendingDashboardDelete.set(null);
    this.pendingWidgetDelete.set(null);
  });

  private readonly selectAvailableDashboard = effect(() => {
    const dashboards = this.dashboards();
    const activeDashboardId = this.activeDashboardId();
    if (!dashboards.some((dashboard) => dashboard.id === activeDashboardId)) {
      this.activeDashboardId.set(dashboards[0]?.id ?? '');
    }
  });

  readonly projectOptions = computed<SearchSelectOption[]>(() =>
    (this.projectsQuery.data() ?? []).map((project) => ({
      value: project.id,
      label: project.name,
      color: colorHash(`project:${project.id}`),
      description: 'Project',
    })),
  );
  readonly contextOptions = computed<SearchSelectOption[]>(() =>
    (this.contextsQuery.data() ?? []).map((context) => ({
      value: context.id,
      label: context.name,
      color: colorHash(context.id),
      badge: context.project?.name,
      keywords: context.project?.name ? [context.project.name] : [],
    })),
  );
  readonly savedQueryOptions = computed<SearchSelectOption[]>(() =>
    savedQueriesAsOptions(this.savedQueriesQuery.data() ?? []),
  );
  readonly queryText = computed(() =>
    selectedSavedQueryText(this.savedQueriesQuery.data() ?? [], this.selectedSavedQueryId()),
  );
  readonly queryPreview = injectQuery(() =>
    this.queryQueries.result(this.selectedWorkspaceId() ?? '', this.queryText(), false),
  );
  readonly previewSummary = computed(() => queryPreviewSummary(this.queryPreview.data()));
  readonly previewContexts = computed(() =>
    contextQueryResultAsListItems(this.queryPreview.data()),
  );
  readonly previewShowError = computed(
    () =>
      this.queryPreview.data() === undefined &&
      (this.queryPreview.isError() || this.queryPreview.isPaused()),
  );
  readonly savedQueryEmptyText = computed(() => {
    if (this.savedQueriesQuery.isLoading()) {
      return 'Loading saved queries…';
    }
    if (this.savedQueriesQuery.isError()) {
      return 'Unable to load saved queries';
    }
    return 'No saved queries';
  });
  readonly dashboardListEmptyText = computed(() => {
    if (this.dashboardsQuery.isLoading()) {
      return 'Loading dashboards…';
    }
    if (this.dashboardsQuery.isError()) {
      return 'Unable to load dashboards';
    }
    return 'No dashboards';
  });
  readonly dashboardFormRequiresTarget = computed(() =>
    dashboardTypeRequiresTarget(this.dashboardFormType()),
  );
  readonly dashboardTargetOptions = computed<SearchSelectOption[]>(() => {
    switch (this.dashboardFormType()) {
      case 'project':
        return this.projectOptions();
      case 'context':
        return this.contextOptions();
      default:
        return [];
    }
  });
  readonly dashboardTargetLoading = computed(() => {
    switch (this.dashboardFormType()) {
      case 'project':
        return this.projectsQuery.isLoading();
      case 'context':
        return this.contextsQuery.isLoading();
      default:
        return false;
    }
  });
  readonly dashboardTargetName = computed(() => this.dashboardFormType());
  readonly dashboardTargetPluralName = computed(() => `${this.dashboardTargetName()}s`);
  readonly dashboardTargetAriaLabel = computed(() => `Select ${this.dashboardTargetName()}`);
  readonly dashboardTargetPlaceholder = computed(() => `Select ${this.dashboardTargetName()}…`);
  readonly dashboardTargetSearchPlaceholder = computed(
    () => `Search ${this.dashboardTargetPluralName()}…`,
  );
  readonly dashboardTargetEmptyText = computed(() =>
    this.dashboardTargetLoading()
      ? `Loading ${this.dashboardTargetPluralName()}…`
      : `No matching ${this.dashboardTargetPluralName()}`,
  );
  readonly dashboardTypeDescription = computed(() => {
    switch (this.dashboardFormType()) {
      case 'workspace':
        return 'This dashboard is used for the selected workspace insight.';
      case 'daily':
        return 'This dashboard is used for daily insights in the selected workspace.';
      case 'custom':
        return 'This dashboard is available as a custom dashboard in the selected workspace.';
      default:
        return '';
    }
  });
  readonly dashboardDialogTitle = computed(() =>
    this.dashboardDialogContentMode() === 'properties'
      ? 'Dashboard properties'
      : 'Create dashboard',
  );
  readonly dashboardMutationPending = computed(() => this.createDashboardMutation.isPending());
  readonly dashboardMutationError = computed(() => this.createDashboardMutation.isError());
  readonly dashboardFormValid = computed(() => {
    const workspaceId = this.selectedWorkspaceId();
    const name = this.dashboardFormName().trim();
    const targetValid =
      !this.dashboardFormRequiresTarget() || this.dashboardFormTargetId().length > 0;
    return Boolean(workspaceId && name && targetValid);
  });

  openCreateDashboardDialog(): void {
    if (!this.selectedWorkspaceId() || this.dashboardEditing()) {
      return;
    }
    this.resetDashboardMutations();
    this.dashboardFormName.set('');
    this.dashboardFormType.set('custom');
    this.dashboardFormTargetId.set('');
    this.dashboardDialogContentMode.set('create');
    this.dashboardDialogMode.set('create');
  }

  selectDashboard(dashboardId: string): void {
    if (
      this.dashboardEditing() ||
      !this.dashboards().some((dashboard) => dashboard.id === dashboardId)
    ) {
      return;
    }
    this.activeDashboardId.set(dashboardId);
  }

  startDashboardEditing(): void {
    const dashboard = this.activeDashboard();
    if (!dashboard) {
      return;
    }
    this.resetDashboardMutations();
    try {
      normalizeDashboardDefinition(dashboard.definition);
    } catch {
      return;
    }
    this.dashboardDraft.set({
      ...dashboard,
      definition: normalizeDashboardDefinition(dashboard.definition),
    });
    this.selectedWidgetId.set(null);
  }

  openDashboardPropertiesDialog(): void {
    const dashboard = this.dashboardDraft();
    if (!dashboard || this.updateDashboardMutation.isPending()) {
      return;
    }
    this.updateDashboardMutation.reset();
    this.dashboardFormName.set(dashboard.name);
    this.dashboardFormType.set(dashboard.type);
    this.dashboardFormTargetId.set(dashboard.targetId ?? '');
    this.dashboardDialogContentMode.set('properties');
    this.dashboardDialogMode.set('properties');
  }

  cancelDashboardEditing(): void {
    if (this.updateDashboardMutation.isPending()) {
      return;
    }
    this.dashboardDialogMode.set(null);
    this.dashboardDraft.set(null);
    this.selectedWidgetId.set(null);
    this.pendingWidgetDelete.set(null);
    this.updateDashboardMutation.reset();
  }

  confirmDashboardEditing(): void {
    const dashboard = this.dashboardDraft();
    if (!dashboard || this.updateDashboardMutation.isPending()) {
      return;
    }
    this.updateDashboardMutation.mutate(dashboard, {
      onSuccess: (savedDashboard) => {
        this.activeDashboardId.set(savedDashboard.id);
        this.dashboardDialogMode.set(null);
        this.dashboardDraft.set(null);
      },
    });
  }

  requestDashboardDelete(): void {
    const dashboard = this.activeDashboard();
    if (dashboard && this.dashboardEditing()) {
      this.pendingDashboardDelete.set(dashboard);
    }
  }

  closeDashboardDialog(): void {
    if (this.dashboardMutationPending()) {
      return;
    }
    this.dashboardDialogMode.set(null);
  }

  setDashboardFormName(event: Event): void {
    this.dashboardFormName.set((event.target as HTMLInputElement).value);
  }

  setDashboardFormType(value: string): void {
    if (!isDashboardType(value)) {
      return;
    }
    if (this.dashboardFormType() === value) {
      return;
    }
    this.dashboardFormType.set(value);
    this.dashboardFormTargetId.set('');
  }

  saveDashboard(): void {
    const workspaceId = this.selectedWorkspaceId();
    const mode = this.dashboardDialogMode();
    if (!workspaceId || !this.dashboardFormValid() || this.dashboardMutationPending()) {
      return;
    }

    if (!mode) {
      return;
    }
    const type = this.dashboardFormType();
    const targetId = dashboardTargetId(type, workspaceId, this.dashboardFormTargetId());
    const name = this.dashboardFormName().trim();

    if (mode === 'properties') {
      const dashboard = this.dashboardDraft();
      if (!dashboard) {
        return;
      }
      this.dashboardDraft.set({ ...dashboard, type, targetId, name });
      this.dashboardDialogMode.set(null);
      return;
    }

    const input: CreateDashboardInput = {
      workspaceId,
      type,
      targetId,
      name,
      definition: { version: 1, query: '', widgets: [] },
    };

    this.createDashboardMutation.mutate(input, {
      onSuccess: (dashboard) => {
        this.activeDashboardId.set(dashboard.id);
        this.dashboardDialogMode.set(null);
      },
    });
  }

  confirmDeleteDashboard(): void {
    const dashboard = this.pendingDashboardDelete();
    if (!dashboard || this.deleteDashboardMutation.isPending()) {
      return;
    }
    this.deleteDashboardMutation.mutate(dashboard, {
      onSuccess: () => {
        if (this.activeDashboardId() === dashboard.id) {
          this.activeDashboardId.set('');
        }
        this.dashboardDraft.set(null);
        this.pendingDashboardDelete.set(null);
      },
    });
  }

  selectSavedQuery(queryId: string): void {
    this.selectedSavedQueryId.set(queryId);
  }

  loadQueryPreview(): void {
    if (this.selectedWorkspaceId()) {
      void this.queryPreview.refetch();
    }
  }

  private resetDashboardMutations(): void {
    this.createDashboardMutation.reset();
    this.updateDashboardMutation.reset();
  }

  addWidget(position?: GridPoint): void {
    const dashboard = this.dashboardDraft();
    if (!dashboard || this.updateDashboardMutation.isPending()) return;
    const definition = normalizeDashboardDefinition(dashboard.definition);
    const next = addTextWidget(definition, crypto.randomUUID(), position);
    if (next === definition) return;
    const widget = next.widgets.at(-1);
    if (widget) {
      this.dashboardDraft.set({ ...dashboard, definition: next });
      this.selectedWidgetId.set(widget.id);
    }
  }

  selectWidget(widgetId: string | null): void {
    this.selectedWidgetId.set(widgetId);
  }

  updateDashboardQuery(event: Event): void {
    const dashboard = this.dashboardDraft();
    if (!dashboard) return;
    const definition = normalizeDashboardDefinition(dashboard.definition);
    this.dashboardDraft.set({
      ...dashboard,
      definition: { ...definition, query: (event.target as HTMLTextAreaElement).value },
    });
  }

  moveWidgetTo(widgetId: string, target: GridPoint): void {
    const dashboard = this.dashboardDraft();
    if (!dashboard || this.updateDashboardMutation.isPending()) return;
    const definition = normalizeDashboardDefinition(dashboard.definition);
    this.dashboardDraft.set({
      ...dashboard,
      definition: moveTextWidget(definition, widgetId, target.x, target.y),
    });
  }

  applyWidgetLayout(widgetId: string, layout: DashboardWidgetLayout): void {
    const dashboard = this.dashboardDraft();
    if (!dashboard || this.updateDashboardMutation.isPending()) return;
    const definition = normalizeDashboardDefinition(dashboard.definition);
    this.dashboardDraft.set({
      ...dashboard,
      definition: updateTextWidgetLayout(definition, widgetId, layout),
    });
  }

  updateWidgetQuery(widgetId: string, query: string): void {
    this.updateWidget(widgetId, { query });
  }

  updateWidgetProperties(widgetId: string, properties: TextWidgetProperties): void {
    this.updateWidget(widgetId, { properties });
  }

  deleteWidget(widgetId: string): void {
    if (this.dashboardEditing() && !this.updateDashboardMutation.isPending())
      this.pendingWidgetDelete.set(widgetId);
  }

  confirmWidgetDelete(): void {
    const widgetId = this.pendingWidgetDelete();
    if (!widgetId || this.updateDashboardMutation.isPending()) return;
    const dashboard = this.dashboardDraft();
    if (!dashboard) return;
    this.dashboardDraft.set({
      ...dashboard,
      definition: deleteTextWidget(normalizeDashboardDefinition(dashboard.definition), widgetId),
    });
    this.selectedWidgetId.set(null);
    this.pendingWidgetDelete.set(null);
  }

  handleWidgetKeydown(event: KeyboardEvent, widget: TextWidgetDefinition): void {
    if (event.key === 'Escape') {
      event.stopPropagation();
      this.selectedWidgetId.set(null);
      return;
    }
    if (
      !this.dashboardEditing() ||
      this.updateDashboardMutation.isPending() ||
      !['ArrowLeft', 'ArrowRight', 'ArrowUp', 'ArrowDown'].includes(event.key)
    )
      return;
    event.preventDefault();
    const selected = this.selectedWidgetId() === widget.id;
    if (!selected) {
      this.selectedWidgetId.set(widget.id);
      return;
    }
    const dashboard = this.dashboardDraft();
    if (!dashboard) return;
    const definition = normalizeDashboardDefinition(dashboard.definition);
    const dx = event.key === 'ArrowLeft' ? -1 : event.key === 'ArrowRight' ? 1 : 0;
    const dy = event.key === 'ArrowUp' ? -1 : event.key === 'ArrowDown' ? 1 : 0;
    const next = event.shiftKey
      ? resizeTextWidget(definition, widget.id, widget.layout.width + dx, widget.layout.height + dy)
      : moveTextWidget(definition, widget.id, widget.layout.x + dx, widget.layout.y + dy);
    this.dashboardDraft.set({ ...dashboard, definition: next });
  }

  private updateWidget(
    widgetId: string,
    update: Partial<Pick<TextWidgetDefinition, 'query' | 'properties'>>,
  ): void {
    const dashboard = this.dashboardDraft();
    if (!dashboard || this.updateDashboardMutation.isPending()) return;
    this.dashboardDraft.set({
      ...dashboard,
      definition: updateTextWidget(
        normalizeDashboardDefinition(dashboard.definition),
        widgetId,
        update,
      ),
    });
  }
}

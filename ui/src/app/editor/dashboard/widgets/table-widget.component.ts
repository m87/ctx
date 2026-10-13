import { Component, computed, input, linkedSignal } from '@angular/core';
import { RouterLink } from '@angular/router';
import {
  PaginationState,
  SortingState,
  Updater,
  functionalUpdate,
  injectTable,
} from '@tanstack/angular-table';
import { HlmButton } from '@spartan-ng/helm/button';
import { HlmTableImports } from '@spartan-ng/helm/table';
import { DashboardMetricResult } from '../../../../api/dashboard/dashboard.service';
import { TableWidgetDefinition } from '../dashboard-definition';
import { formatResultValue } from './format-result-value';
import {
  initialSorting,
  tableColumns,
  tableRowsByKey,
  tableWidgetFeatures,
} from './table-widget-columns';

type TableResult = Extract<DashboardMetricResult, { kind: 'table' }>;
type TableColumn = TableWidgetDefinition['properties']['columns'][number];
type Unit = 'ms' | 'count' | 'percent';

@Component({
  selector: 'ctx-table-widget',
  imports: [HlmButton, HlmTableImports, RouterLink],
  host: { class: 'flex min-h-0 flex-1 flex-col gap-3' },
  template: `
    <div hlmTableContainer class="min-h-0 flex-1 overflow-auto" [attr.aria-busy]="false">
      @if (result().table.totalRows === 0) {
        <p class="p-4 text-sm text-muted-foreground">No matching groups.</p>
      } @else {
        <table hlmTable class="text-xs">
          <thead hlmTHead class="sticky top-0 bg-card">
            <tr hlmTr>
              @for (header of table.getFlatHeaders(); track header.id) {
                <th
                  hlmTh
                  scope="col"
                  class="h-8"
                  [attr.aria-sort]="ariaSort(header.column.getIsSorted())"
                >
                  <button type="button" class="font-medium" (click)="header.column.toggleSorting()">
                    {{ header.column.columnDef.header
                    }}<span aria-hidden="true">{{ sortArrow(header.column.getIsSorted()) }}</span>
                  </button>
                </th>
              }
            </tr>
          </thead>
          <tbody hlmTBody>
            @for (row of table.getRowModel().rows; track row.id) {
              <tr hlmTr class="border-border/60">
                <th hlmTh scope="row" class="h-auto py-1 pr-2" [title]="row.original.label">
                  <a
                    [routerLink]="
                      row.original.entity
                        ? [
                            row.original.entity.type === 'project' ? '/project' : '/context',
                            row.original.entity.id,
                          ]
                        : null
                    "
                    >{{ row.original.label }}</a
                  >
                </th>
                @for (column of visibleColumns(); track column.id) {
                  <td hlmTd class="py-1 text-right tabular-nums">
                    {{ format(row.original.values[column.id] ?? null, column) }}
                    @if (column.showBar && row.original.values[column.id] !== null) {
                      <span
                        class="ml-1 inline-block h-1 rounded-sm bg-primary align-middle"
                        [style.width.%]="barWidth(row.original.values[column.id] ?? null, column)"
                      ></span>
                    }
                  </td>
                }
              </tr>
            }
          </tbody>
          @if (widget().properties.showSummary) {
            <tfoot hlmTFoot>
              <tr hlmTr class="border-border font-semibold">
                <th hlmTh scope="row" class="h-auto py-1 pr-2">Total</th>
                @for (column of visibleColumns(); track column.id) {
                  <td hlmTd class="py-1 text-right tabular-nums">
                    {{ format(result().table.summary[column.id] ?? null, column) }}
                  </td>
                }
              </tr>
            </tfoot>
          }
        </table>
      }
    </div>
    <footer class="flex items-center justify-between gap-2 text-xs text-muted-foreground">
      <span aria-live="polite">
        @if (result().table.totalRows) {
          {{ pageStart() }}–{{ pageEnd() }} of {{ result().table.totalRows }}
        }
      </span>
      <nav class="flex items-center gap-1" aria-label="Table pages">
        <button
          hlmBtn
          type="button"
          variant="outline"
          size="sm"
          [disabled]="!table.getCanPreviousPage()"
          (click)="table.previousPage()"
        >
          Previous
        </button>
        @if (table.getPageCount()) {
          <span class="px-1">{{ pageIndex() + 1 }} / {{ table.getPageCount() }}</span>
        }
        <button
          hlmBtn
          type="button"
          variant="outline"
          size="sm"
          [disabled]="!table.getCanNextPage()"
          (click)="table.nextPage()"
        >
          Next
        </button>
      </nav>
    </footer>
  `,
})
export class TableWidgetComponent {
  readonly widget = input.required<TableWidgetDefinition>();
  readonly result = input.required<TableResult>();
  readonly visibleColumns = computed(() =>
    this.widget().properties.columns.filter((column) => column.visible),
  );
  private readonly columns = computed(() => tableColumns(this.widget().properties));
  private readonly rows = computed(() => tableRowsByKey(this.result().table.rows));
  private readonly sorting = linkedSignal<SortingState>(() =>
    initialSorting(this.widget().properties),
  );
  private readonly requestedPage = linkedSignal<PaginationState>(() => ({
    pageIndex: 0,
    pageSize: this.widget().properties.pageSize,
  }));
  private readonly pagination = computed<PaginationState>(() => {
    const { pageIndex, pageSize } = this.requestedPage();
    const lastPage = Math.max(0, Math.ceil(this.rows().length / pageSize) - 1);
    return { pageIndex: Math.min(pageIndex, lastPage), pageSize };
  });
  readonly table = injectTable(() => ({
    features: tableWidgetFeatures,
    columns: this.columns(),
    data: this.rows(),
    state: { sorting: this.sorting(), pagination: this.pagination() },
    onSortingChange: (updater: Updater<SortingState>) => {
      this.sorting.set(functionalUpdate(updater, this.sorting()));
      this.requestedPage.update((page) => ({ ...page, pageIndex: 0 }));
    },
    onPaginationChange: (updater: Updater<PaginationState>) =>
      this.requestedPage.set(functionalUpdate(updater, this.pagination())),
    autoResetPageIndex: false,
    enableSortingRemoval: false,
    enableMultiSort: false,
  }));
  readonly pageIndex = computed(() => this.pagination().pageIndex);
  readonly pageStart = computed(() => this.pagination().pageIndex * this.pagination().pageSize + 1);
  readonly pageEnd = computed(() =>
    Math.min((this.pagination().pageIndex + 1) * this.pagination().pageSize, this.rows().length),
  );

  ariaSort(direction: false | 'asc' | 'desc'): 'ascending' | 'descending' | 'none' {
    return direction === false ? 'none' : direction === 'asc' ? 'ascending' : 'descending';
  }
  sortArrow(direction: false | 'asc' | 'desc'): string {
    return direction === false ? '' : direction === 'asc' ? ' ↑' : ' ↓';
  }
  format(value: number | null, column: TableColumn): string {
    return formatResultValue(value, this.metricUnit(column.metric), column);
  }
  barWidth(value: number | null, column: TableColumn): number {
    if (value === null) return 0;
    const maxValue =
      column.metric === 'time_share.ratio'
        ? 100
        : (this.result().table.columnMaxima[column.id] ?? 0);
    return maxValue <= 0 ? 0 : Math.min(100, (value / maxValue) * 100);
  }

  private metricUnit(metric: string): Unit {
    return metric.startsWith('duration.')
      ? 'ms'
      : metric === 'time_share.ratio'
        ? 'percent'
        : 'count';
  }
}

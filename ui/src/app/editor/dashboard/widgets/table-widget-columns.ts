import {
  ColumnDef,
  SortFn,
  SortingState,
  createPaginatedRowModel,
  createSortedRowModel,
  rowPaginationFeature,
  rowSortingFeature,
  tableFeatures,
} from '@tanstack/angular-table';
import { DashboardTablePayload } from '../../../../api/dashboard/dashboard.service';
import { TableWidgetProperties } from '../metric-widget-definition';

export const GROUP_COLUMN_ID = 'group';

export type TableRowData = DashboardTablePayload['rows'][number];

export const tableWidgetFeatures = tableFeatures({
  rowSortingFeature,
  sortedRowModel: createSortedRowModel(),
  rowPaginationFeature,
  paginatedRowModel: createPaginatedRowModel(),
});

type TableFeatures = typeof tableWidgetFeatures;
type TableColumnDef = ColumnDef<TableFeatures, TableRowData>;

const groupSort: SortFn<TableFeatures, TableRowData> = (a, b) => {
  const left = a.original.label.toLowerCase();
  const right = b.original.label.toLowerCase();
  return left < right ? -1 : left > right ? 1 : 0;
};

const numberSort: SortFn<TableFeatures, TableRowData> = (a, b, columnId) => {
  const left = a.getValue<number>(columnId);
  const right = b.getValue<number>(columnId);
  return left < right ? -1 : left > right ? 1 : 0;
};

export function metricLabel(metric: string): string {
  return metric.replaceAll('.', ' ');
}

export function tableColumns(properties: TableWidgetProperties): TableColumnDef[] {
  return [
    {
      id: GROUP_COLUMN_ID,
      header: properties.groupBy === 'project' ? 'Project' : 'Context',
      accessorFn: (row) => row.label,
      sortFn: groupSort,
      sortDescFirst: true,
    },
    ...properties.columns
      .filter((column) => column.visible)
      .map((column): TableColumnDef => ({
        id: column.id,
        header: column.label || metricLabel(column.metric),
        accessorFn: (row) => row.values[column.id] ?? undefined,
        sortFn: numberSort,
        sortUndefined: 'last',
        sortDescFirst: true,
      })),
  ];
}

export function tableRowsByKey(rows: readonly TableRowData[]): TableRowData[] {
  return [...rows].sort((a, b) => (a.key < b.key ? -1 : a.key > b.key ? 1 : 0));
}

export function initialSorting(properties: TableWidgetProperties): SortingState {
  const { columnId, direction } = properties.sort;
  const known =
    columnId === GROUP_COLUMN_ID ||
    properties.columns.some((column) => column.id === columnId && column.visible);
  return known ? [{ id: columnId, desc: direction === 'desc' }] : [];
}

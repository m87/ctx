import { defaultTableProperties } from '../metric-widget-definition';
import {
  GROUP_COLUMN_ID,
  TableRowData,
  initialSorting,
  tableColumns,
  tableRowsByKey,
} from './table-widget-columns';

describe('tableColumns', () => {
  it('starts with the group column named after the grouping and lists visible metrics only', () => {
    const properties = defaultTableProperties();
    const hidden = {
      ...properties,
      groupBy: 'context' as const,
      columns: properties.columns.map((column) =>
        column.id === 'share' ? { ...column, visible: false } : column,
      ),
    };
    const columns = tableColumns(hidden);
    expect(columns.map((column) => column.id)).toEqual([
      GROUP_COLUMN_ID,
      'time',
      'sessions',
      'average',
    ]);
    expect(columns[0].header).toBe('Context');
    expect(columns[1].header).toBe('Time');
  });

  it('uses the metric name as the header when a column has no label', () => {
    const properties = defaultTableProperties();
    const columns = tableColumns({
      ...properties,
      columns: properties.columns.map((column) => ({ ...column, label: '' })),
    });
    expect(columns[2].header).toBe('sessions count');
  });
});

describe('tableRowsByKey', () => {
  it('orders rows by key without touching the input', () => {
    const rows = [{ key: 'b' }, { key: 'a' }, { key: 'c' }] as TableRowData[];
    expect(tableRowsByKey(rows).map((row) => row.key)).toEqual(['a', 'b', 'c']);
    expect(rows.map((row) => row.key)).toEqual(['b', 'a', 'c']);
  });
});

describe('initialSorting', () => {
  it('turns the saved sort into table state', () => {
    const properties = defaultTableProperties();
    expect(
      initialSorting({ ...properties, sort: { columnId: 'time', direction: 'desc' } }),
    ).toEqual([{ id: 'time', desc: true }]);
    expect(
      initialSorting({ ...properties, sort: { columnId: GROUP_COLUMN_ID, direction: 'asc' } }),
    ).toEqual([{ id: GROUP_COLUMN_ID, desc: false }]);
  });

  it('ignores a sort on a hidden or unknown column', () => {
    const properties = defaultTableProperties();
    const hidden = {
      ...properties,
      sort: { columnId: 'share', direction: 'desc' as const },
      columns: properties.columns.map((column) =>
        column.id === 'share' ? { ...column, visible: false } : column,
      ),
    };
    expect(initialSorting(hidden)).toEqual([]);
    expect(initialSorting({ ...properties, sort: { columnId: 'nope', direction: 'asc' } })).toEqual(
      [],
    );
  });
});

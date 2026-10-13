import { TestBed } from '@angular/core/testing';
import { provideRouter } from '@angular/router';
import { DashboardMetricResult } from '../../../../api/dashboard/dashboard.service';
import { TableWidgetDefinition } from '../dashboard-definition';
import { defaultTableProperties } from '../metric-widget-definition';
import { TableWidgetComponent } from './table-widget.component';

type TableResult = Extract<DashboardMetricResult, { kind: 'table' }>;
type Row = TableResult['table']['rows'][number];

function widget(
  overrides: Partial<TableWidgetDefinition['properties']> = {},
): TableWidgetDefinition {
  const properties = defaultTableProperties();
  return {
    id: 'table-1',
    type: 'table',
    layout: { x: 0, y: 0, width: 10, height: 6 },
    query: '',
    properties: {
      ...properties,
      sort: { columnId: 'time', direction: 'desc' },
      columns: properties.columns.map((column) =>
        column.id === 'time' ? { ...column, showBar: true } : column,
      ),
      ...overrides,
    },
  };
}

function row(key: string, label: string, time: number | null, sessions: number | null): Row {
  return {
    key,
    label,
    entity: key === 'p1' ? { type: 'project', id: 'p1' } : null,
    values: { time, sessions, average: null, share: null },
  };
}

const HOUR = 3_600_000;

function result(rows: Row[]): TableResult {
  return {
    kind: 'table',
    value: null,
    timeZone: 'UTC',
    hasActivity: true,
    table: {
      columns: [],
      rows,
      summary: { time: 10 * HOUR, sessions: 4 },
      columnMaxima: { time: 4 * HOUR, sessions: 9 },
      totalRows: rows.length,
    },
  } as unknown as TableResult;
}

const sample = [
  row('p3', 'charlie', 2 * HOUR, 9),
  row('p1', 'Alpha', 4 * HOUR, 4),
  row('p4', 'delta', null, null),
  row('p2', 'bravo', 2 * HOUR, 4),
];

function manyRows(count: number): Row[] {
  return Array.from({ length: count }, (_, index) =>
    row(`k${String(index).padStart(2, '0')}`, `Row ${index}`, (index + 1) * HOUR, index),
  );
}

function create(
  properties: Partial<TableWidgetDefinition['properties']> = {},
  rows: Row[] = sample,
) {
  const fixture = TestBed.createComponent(TableWidgetComponent);
  fixture.componentRef.setInput('widget', widget(properties));
  fixture.componentRef.setInput('result', result(rows));
  fixture.detectChanges();
  const element = fixture.nativeElement as HTMLElement;
  return {
    fixture,
    element,
    labels: () =>
      Array.from(element.querySelectorAll('tbody th')).map((th) => th.textContent?.trim()),
    header: (name: string) =>
      Array.from(element.querySelectorAll<HTMLButtonElement>('thead button')).find((button) =>
        button.textContent?.trim().startsWith(name),
      ) as HTMLButtonElement,
    sorted: () =>
      Array.from(element.querySelectorAll('thead th'))
        .filter((th) => th.getAttribute('aria-sort') !== 'none')
        .map(
          (th) => `${th.textContent?.replace(/[↑↓]/g, '').trim()}:${th.getAttribute('aria-sort')}`,
        ),
    footerButton: (index: number) =>
      element.querySelectorAll<HTMLButtonElement>('footer button')[index],
    range: () => element.querySelector('[aria-live]')?.textContent?.replace(/\s+/g, ' ').trim(),
    update: (next: Row[]) => {
      fixture.componentRef.setInput('result', result(next));
      fixture.detectChanges();
    },
    click: (button: HTMLButtonElement) => {
      button.click();
      fixture.detectChanges();
    },
  };
}

describe('TableWidgetComponent', () => {
  beforeEach(() => TestBed.configureTestingModule({ providers: [provideRouter([])] }));
  afterEach(() => TestBed.resetTestingModule());

  it('renders the group and visible metric headers', () => {
    const { element } = create();
    const headers = Array.from(element.querySelectorAll('thead th')).map((header) =>
      header.textContent?.replace(/[↑↓]/g, '').trim(),
    );
    expect(headers).toEqual(['Project', 'Time', 'Sessions', 'Average session', 'Share']);
  });

  it('uses the context heading and hides hidden columns', () => {
    const properties = widget().properties;
    const { element } = create({
      groupBy: 'context',
      columns: properties.columns.map((column) =>
        column.id === 'share' ? { ...column, visible: false } : column,
      ),
    });
    const labels = Array.from(element.querySelectorAll('thead th')).map((th) =>
      th.textContent?.trim(),
    );
    expect(labels[0]).toBe('Context');
    expect(labels).not.toContain('Share');
  });

  it('starts with the saved sort and keeps ties ordered by key', () => {
    const table = create();
    expect(table.sorted()).toEqual(['Time:descending']);
    expect(table.labels()).toEqual(['Alpha', 'bravo', 'charlie', 'delta']);
    const asc = create({ sort: { columnId: 'time', direction: 'asc' } });
    expect(asc.sorted()).toEqual(['Time:ascending']);
    expect(asc.labels()).toEqual(['bravo', 'charlie', 'Alpha', 'delta']);
  });

  it('puts rows without a value last in both directions', () => {
    const table = create({ sort: { columnId: 'sessions', direction: 'desc' } });
    expect(table.labels()).toEqual(['charlie', 'Alpha', 'bravo', 'delta']);
    table.click(table.header('Sessions'));
    expect(table.sorted()).toEqual(['Sessions:ascending']);
    expect(table.labels()).toEqual(['Alpha', 'bravo', 'charlie', 'delta']);
  });

  it('cycles a clicked column through descending, ascending and descending without clearing', () => {
    const table = create();
    table.click(table.header('Sessions'));
    expect(table.sorted()).toEqual(['Sessions:descending']);
    table.click(table.header('Sessions'));
    expect(table.sorted()).toEqual(['Sessions:ascending']);
    table.click(table.header('Sessions'));
    expect(table.sorted()).toEqual(['Sessions:descending']);
  });

  it('sorts the group column without regard to case, starting descending', () => {
    const table = create();
    table.click(table.header('Project'));
    expect(table.sorted()).toEqual(['Project:descending']);
    expect(table.labels()).toEqual(['delta', 'charlie', 'bravo', 'Alpha']);
    table.click(table.header('Project'));
    expect(table.labels()).toEqual(['Alpha', 'bravo', 'charlie', 'delta']);
  });

  it('sorts without emitting anything or asking for new data', () => {
    const table = create();
    const result = table.fixture.componentInstance.result();
    table.click(table.header('Sessions'));
    expect(table.fixture.componentInstance.result()).toBe(result);
    expect('viewStateChange' in table.fixture.componentInstance).toBe(false);
  });

  it('formats values, links entities and draws a bar for columns with showBar', () => {
    const { element } = create();
    const rows = element.querySelectorAll('tbody tr');
    expect(rows[0].querySelector('a')?.getAttribute('href')).toBe('/project/p1');
    expect(rows[1].querySelector('a')?.getAttribute('href')).toBeNull();
    expect(rows[0].textContent).toContain('4 h');
    expect(rows[3].textContent).toContain('—');
    const bars = element.querySelectorAll<HTMLElement>('tbody span.bg-primary');
    expect(Array.from(bars).map((bar) => bar.style.width)).toEqual(['100%', '50%', '50%']);
  });

  it('shows the summary row only when requested', () => {
    expect(create({ showSummary: true }).element.querySelector('tfoot')?.textContent).toContain(
      '10 h',
    );
    expect(create({ showSummary: false }).element.querySelector('tfoot')).toBeNull();
  });

  it('pages in the browser and reports the visible range', () => {
    const table = create({}, manyRows(12));
    expect(table.labels()).toHaveLength(10);
    expect(table.range()).toBe('1–10 of 12');
    expect(table.element.textContent).toContain('1 / 2');
    expect(table.footerButton(0).disabled).toBe(true);
    expect(table.footerButton(1).disabled).toBe(false);
    table.click(table.footerButton(1));
    expect(table.labels()).toEqual(['Row 1', 'Row 0']);
    expect(table.range()).toBe('11–12 of 12');
    expect(table.footerButton(1).disabled).toBe(true);
    table.click(table.footerButton(0));
    expect(table.range()).toBe('1–10 of 12');
  });

  it('honours the saved page size', () => {
    const table = create({ pageSize: 25 }, manyRows(30));
    expect(table.labels()).toHaveLength(25);
    expect(table.range()).toBe('1–25 of 30');
  });

  it('goes back to the first page when the sorting changes', () => {
    const table = create({}, manyRows(12));
    table.click(table.footerButton(1));
    expect(table.range()).toBe('11–12 of 12');
    table.click(table.header('Sessions'));
    expect(table.range()).toBe('1–10 of 12');
  });

  it('keeps the sort and the page when the data refreshes with a new result', () => {
    const table = create({}, manyRows(12));
    table.click(table.header('Sessions'));
    table.click(table.header('Sessions'));
    table.click(table.footerButton(1));
    table.update(manyRows(12));
    expect(table.sorted()).toEqual(['Sessions:ascending']);
    expect(table.range()).toBe('11–12 of 12');
  });

  it('moves to the last page when a refresh leaves fewer rows', () => {
    const table = create({}, manyRows(12));
    table.click(table.footerButton(1));
    table.update(manyRows(5));
    expect(table.range()).toBe('1–5 of 5');
    expect(table.footerButton(0).disabled).toBe(true);
    expect(table.footerButton(1).disabled).toBe(true);
  });

  it('adopts a changed saved sort from the widget properties', () => {
    const table = create();
    table.click(table.header('Sessions'));
    table.fixture.componentRef.setInput(
      'widget',
      widget({ sort: { columnId: 'group', direction: 'asc' } }),
    );
    table.fixture.detectChanges();
    expect(table.sorted()).toEqual(['Project:ascending']);
    expect(table.labels()).toEqual(['Alpha', 'bravo', 'charlie', 'delta']);
  });

  it('shows the empty state and disables both page buttons without rows', () => {
    const { element, footerButton } = create({}, []);
    expect(element.textContent).toContain('No matching groups.');
    expect(footerButton(0).disabled).toBe(true);
    expect(footerButton(1).disabled).toBe(true);
  });
});

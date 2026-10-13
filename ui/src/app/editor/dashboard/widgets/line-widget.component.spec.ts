import { ChartPoint } from '@tanstack/charts';
import { TestBed } from '@angular/core/testing';
import { DashboardMetricResult } from '../../../../api/dashboard/dashboard.service';
import { LineWidgetDefinition } from '../dashboard-definition';
import { defaultLineProperties } from '../metric-widget-definition';
import { LineWidgetComponent } from './line-widget.component';

type SeriesResult = Extract<DashboardMetricResult, { kind: 'series' }>;

function widget(overrides: Partial<LineWidgetDefinition['properties']> = {}): LineWidgetDefinition {
  return {
    id: 'line-1',
    type: 'line',
    layout: { x: 0, y: 0, width: 4, height: 3 },
    query: '',
    properties: { ...defaultLineProperties(), ...overrides },
  };
}

function result(): SeriesResult {
  return {
    kind: 'series',
    value: null,
    unit: 'count',
    timeZone: 'UTC',
    hasActivity: true,
    series: {
      unit: 'count',
      buckets: [
        {
          key: 'a',
          start: '2026-10-10T00:00:00Z',
          end: '2026-10-11T00:00:00Z',
          status: 'complete',
        },
        {
          key: 'b',
          start: '2026-10-11T00:00:00Z',
          end: '2026-10-12T00:00:00Z',
          status: 'complete',
        },
        { key: 'c', start: '2026-10-12T00:00:00Z', end: '2026-10-13T00:00:00Z', status: 'future' },
      ],
      series: [
        { key: 's1', label: 'Alpha', values: [1, 3, null] },
        { key: 's2', label: 'Beta', values: [2, null, null] },
      ],
      totalValue: 6,
      totalSeries: 2,
      omittedSeries: 0,
    },
  } as unknown as SeriesResult;
}

function create(properties: Partial<LineWidgetDefinition['properties']> = {}) {
  const fixture = TestBed.createComponent(LineWidgetComponent);
  fixture.componentRef.setInput('widget', widget(properties));
  fixture.componentRef.setInput('result', result());
  fixture.detectChanges();
  return fixture;
}

describe('LineWidgetComponent', () => {
  beforeEach(() => {
    vi.stubGlobal(
      'ResizeObserver',
      class {
        observe(): void {}
        unobserve(): void {}
        disconnect(): void {}
      },
    );
  });
  afterEach(() => {
    TestBed.resetTestingModule();
    vi.unstubAllGlobals();
  });

  it('mounts a chart with an accessible name and lists every series in the legend', () => {
    const fixture = create();
    expect(fixture.nativeElement.querySelector('ctx-chart-frame')).not.toBeNull();
    expect(fixture.componentInstance.options().ariaLabel).toBe(
      'smooth chart with 2 series and 3 time buckets',
    );
    expect(fixture.nativeElement.textContent).toContain('Alpha');
    expect(fixture.nativeElement.textContent).toContain('Beta');
  });

  it('keeps the chart out of the tab order while the dashboard is being edited', () => {
    const fixture = create();
    expect(fixture.componentInstance.options().tabIndex).toBe(0);
    fixture.componentRef.setInput('editing', true);
    fixture.detectChanges();
    expect(fixture.componentInstance.options().tabIndex).toBe(-1);
  });

  it('announces the focused bucket with the value of every series', () => {
    const fixture = create({ format: 'number', precision: 0 });
    fixture.componentInstance.options().onFocusChange?.({
      xValue: '2026-10-11T00:00:00Z',
    } as unknown as ChartPoint);
    fixture.detectChanges();
    const status = fixture.nativeElement.querySelector('[aria-live="polite"]');
    expect(status.textContent).toContain('Oct 11, 2026');
    expect(status.textContent).toContain('Alpha: 3 · Beta: —');
    fixture.componentInstance.options().onFocusChange?.(null);
    fixture.detectChanges();
    expect(status.textContent.trim()).toBe('');
  });

  it('shows the total and omitted series in the summary', () => {
    const fixture = create({ showSummary: true, format: 'number', precision: 0 });
    expect(fixture.nativeElement.textContent).toContain('6');
    expect(fixture.nativeElement.textContent).toContain('total');
  });
});

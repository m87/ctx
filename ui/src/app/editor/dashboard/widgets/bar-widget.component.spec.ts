import { TestBed } from '@angular/core/testing';
import { Router, provideRouter } from '@angular/router';
import { ChartBounds, ChartPoint } from '@tanstack/charts';
import { DashboardMetricResult } from '../../../../api/dashboard/dashboard.service';
import { BarWidgetDefinition } from '../dashboard-definition';
import { defaultBarProperties } from '../metric-widget-definition';
import { BAR_LABEL_HEIGHT, BAR_LABEL_WIDTH } from './bar-chart-definition';
import { BarWidgetComponent } from './bar-widget.component';

type CategoriesResult = Extract<DashboardMetricResult, { kind: 'categories' }>;

function widget(overrides: Partial<BarWidgetDefinition['properties']> = {}): BarWidgetDefinition {
  return {
    id: 'bar-1',
    type: 'bar',
    layout: { x: 0, y: 0, width: 6, height: 4 },
    query: '',
    properties: { ...defaultBarProperties(), format: 'number', precision: 0, ...overrides },
  };
}

function result(count = 3): CategoriesResult {
  const categories = Array.from({ length: count }, (_, index) => ({
    key: `k${index}`,
    label: `Category ${index}`,
    value: (index + 1) * 2,
    entity: index === 1 ? null : { type: index === 0 ? 'project' : 'context', id: `id${index}` },
    synthetic: false,
    memberCount: 1,
    start: null,
    end: null,
    status: 'complete',
  }));
  return {
    kind: 'categories',
    value: null,
    unit: 'count',
    timeZone: 'UTC',
    hasActivity: count > 0,
    categories: {
      unit: 'count',
      categories,
      totalValue: 12,
      totalCategories: count,
      omittedCategories: 1,
    },
  } as unknown as CategoriesResult;
}

function create(
  properties: Partial<BarWidgetDefinition['properties']> = {},
  data: CategoriesResult = result(),
) {
  const fixture = TestBed.createComponent(BarWidgetComponent);
  fixture.componentRef.setInput('widget', widget(properties));
  fixture.componentRef.setInput('result', data);
  fixture.detectChanges();
  return fixture;
}

describe('BarWidgetComponent', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({ providers: [provideRouter([])] });
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

  it('shows an empty state without categories', () => {
    const fixture = create({}, result(0));
    expect(fixture.nativeElement.textContent).toContain('No activity in the selected period.');
    expect(fixture.nativeElement.querySelector('ctx-chart-frame')).toBeNull();
  });

  it('summarises the categories, the total and the omitted ones', () => {
    const fixture = create({ showSummary: true });
    expect(fixture.nativeElement.textContent.replace(/\s+/g, ' ')).toContain(
      '3 categories · 12 total · 1 omitted',
    );
  });

  it('places label links beside horizontal bars and plain labels where there is no entity', () => {
    const fixture = create();
    fixture.componentInstance.area.set({
      x: BAR_LABEL_WIDTH,
      y: 4,
      width: 300,
      height: 90,
    } as ChartBounds);
    fixture.detectChanges();
    const links: HTMLAnchorElement[] = Array.from(fixture.nativeElement.querySelectorAll('a'));
    expect(links.map((link) => link.getAttribute('href'))).toEqual([
      '/project/id0',
      '/context/id2',
    ]);
    expect(links.map((link) => link.textContent?.trim())).toEqual(['Category 0', 'Category 2']);
    expect(links[0].style.top).toBe('4px');
    expect(links[1].style.top).toBe('64px');
    expect(fixture.nativeElement.textContent).toContain('Category 1');
  });

  it('places the labels under vertical bars', () => {
    const fixture = create({ orientation: 'vertical' });
    fixture.componentInstance.area.set({ x: 40, y: 10, width: 300, height: 120 } as ChartBounds);
    fixture.detectChanges();
    const first: HTMLElement = fixture.nativeElement.querySelector('a');
    expect(first.style.left).toBe('40px');
    expect(first.style.width).toBe('100px');
    expect(first.style.top).toBe(`${10 + 120 + 8}px`);
    expect(first.style.height).toBe(`${BAR_LABEL_HEIGHT - 8}px`);
  });

  it('navigates to the project or context when a bar is selected', () => {
    const fixture = create();
    const navigate = vi.spyOn(TestBed.inject(Router), 'navigate').mockResolvedValue(true);
    const select = fixture.componentInstance.options().onSelect;
    select?.({ datum: { entity: { type: 'context', id: 'c1' } } } as unknown as ChartPoint);
    expect(navigate).toHaveBeenCalledWith(['/context', 'c1']);
    select?.({ datum: { entity: null } } as unknown as ChartPoint);
    select?.(null);
    expect(navigate).toHaveBeenCalledTimes(1);
  });

  it('keeps the chart out of the tab order while editing and sizes scroll areas', () => {
    const fixture = create({ orientation: 'vertical' });
    expect(fixture.componentInstance.options().tabIndex).toBe(0);
    expect(fixture.componentInstance.minimumWidth()).toBeGreaterThan(0);
    fixture.componentRef.setInput('editing', true);
    fixture.detectChanges();
    expect(fixture.componentInstance.options().tabIndex).toBe(-1);
  });
});

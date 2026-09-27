import { TestBed } from '@angular/core/testing';
import { DashboardMetricResult } from '../../../../api/dashboard/dashboard.service';
import { defaultMetricProperties } from '../metric-widget-definition';
import { KpiWidgetComponent } from './kpi-widget.component';

describe('KpiWidgetComponent', () => {
  afterEach(() => TestBed.resetTestingModule());

  it('shows loading, value, zero, no data and errors as separate states', () => {
    const fixture = TestBed.createComponent(KpiWidgetComponent);
    fixture.componentRef.setInput('properties', defaultMetricProperties());
    fixture.detectChanges();
    expect(fixture.nativeElement.textContent).toContain('Loading');
    const result: DashboardMetricResult = {
      value: 3_600_000,
      unit: 'ms',
      timeZone: 'UTC',
      start: '2026-09-27T00:00:00Z',
      end: '2026-09-28T00:00:00Z',
      ignoredIntervals: 0,
    };
    fixture.componentRef.setInput('result', result);
    fixture.componentRef.setInput('period', 'Sep 27');
    fixture.detectChanges();
    expect(fixture.nativeElement.textContent).toContain('1 h');
    expect(fixture.nativeElement.textContent).toContain('Sep 27');
    fixture.componentRef.setInput('result', { ...result, value: 0 });
    fixture.detectChanges();
    expect(fixture.nativeElement.textContent).toContain('0 min');
    fixture.componentRef.setInput('result', { ...result, value: null });
    fixture.detectChanges();
    expect(fixture.nativeElement.textContent).toContain('No data');
    fixture.componentRef.setInput('error', 'Unable to load metric data.');
    fixture.detectChanges();
    expect(fixture.nativeElement.textContent).toContain('Unable to load metric data.');
  });
});

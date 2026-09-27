import { TestBed } from '@angular/core/testing';
import { DashboardMetricResult } from '../../../../api/dashboard/dashboard.service';
import { defaultGaugeProperties } from '../metric-widget-definition';
import { GaugeWidgetComponent } from './gauge-widget.component';

describe('GaugeWidgetComponent', () => {
  afterEach(() => TestBed.resetTestingModule());

  it('clamps the arc while keeping over-target values and accessible limit wording', () => {
    const fixture = TestBed.createComponent(GaugeWidgetComponent);
    fixture.componentRef.setInput('properties', {
      ...defaultGaugeProperties(),
      maximum: 40 * 3_600_000,
      target: 40 * 3_600_000,
      goalMode: 'at_most',
    });
    const result: DashboardMetricResult = {
      value: 45 * 3_600_000,
      unit: 'ms',
      timeZone: 'UTC',
      start: '2026-09-21T00:00:00Z',
      end: '2026-09-28T00:00:00Z',
      ignoredIntervals: 0,
    };
    fixture.componentRef.setInput('result', result);
    fixture.componentRef.setInput('period', 'Sep 21 – Sep 27');
    fixture.detectChanges();
    expect(fixture.componentInstance.progress()).toBe(100);
    expect(fixture.componentInstance.percentage()).toBe('112.5');
    expect(fixture.nativeElement.textContent).toContain('45 h');
    expect(fixture.nativeElement.textContent).toContain('Limit exceeded');
    expect(fixture.nativeElement.textContent).toContain('Sep 21 – Sep 27');
    expect(
      fixture.nativeElement.querySelector('[role="img"]').getAttribute('aria-label'),
    ).toContain('limit');
  });

  it('shows no data and period even when the ordinary value is hidden', () => {
    const fixture = TestBed.createComponent(GaugeWidgetComponent);
    fixture.componentRef.setInput('properties', {
      ...defaultGaugeProperties(),
      showValue: false,
      description: 'Optional description',
    });
    fixture.componentRef.setInput('result', {
      value: null,
      unit: 'ms',
      timeZone: 'UTC',
      start: '2026-09-27T00:00:00Z',
      end: '2026-09-28T00:00:00Z',
      ignoredIntervals: 0,
    });
    fixture.componentRef.setInput('period', 'Sep 27');
    fixture.detectChanges();
    expect(fixture.nativeElement.textContent).toContain('No data');
    expect(fixture.nativeElement.textContent).toContain('Optional description');
    expect(fixture.nativeElement.textContent).toContain('Sep 27');
  });
});

import { TestBed } from '@angular/core/testing';
import { createWidget, GaugeWidgetDefinition } from '../dashboard-definition';
import { MetricWidgetPropertiesComponent } from './metric-widget-properties.component';

function gauge(id: string): GaugeWidgetDefinition {
  const widget = createWidget('gauge', id, { x: 0, y: 0, width: 6, height: 5 });
  if (widget.type !== 'gauge') throw new Error('Expected gauge widget');
  return widget;
}

describe('MetricWidgetPropertiesComponent', () => {
  afterEach(() => TestBed.resetTestingModule());

  it('resets gauge scale when units change and rejects dependent scale changes', () => {
    const fixture = TestBed.createComponent(MetricWidgetPropertiesComponent);
    fixture.componentRef.setInput('widget', gauge('first'));
    fixture.detectChanges();
    const changed = vi.fn();
    fixture.componentInstance.widgetChange.subscribe(changed);
    const metric = fixture.nativeElement.querySelector('#metric-first-metric') as HTMLSelectElement;
    metric.value = 'sessions.count';
    metric.dispatchEvent(new Event('change'));
    fixture.detectChanges();
    expect(changed).toHaveBeenLastCalledWith(
      expect.objectContaining({
        properties: expect.objectContaining({ maximum: 100, target: 100, thresholds: [] }),
      }),
    );
    expect(fixture.nativeElement.textContent).toContain('reset for the new unit');

    const nextWidget = changed.mock.lastCall?.[0] as GaugeWidgetDefinition;
    fixture.componentRef.setInput('widget', nextWidget);
    fixture.detectChanges();
    const maximum = fixture.nativeElement.querySelector(
      '#metric-first-maximum',
    ) as HTMLInputElement;
    maximum.value = '20';
    maximum.dispatchEvent(new Event('change'));
    fixture.detectChanges();
    expect(changed).toHaveBeenCalledTimes(1);
    expect(maximum.getAttribute('aria-invalid')).toBe('true');
    expect(fixture.nativeElement.textContent).toContain('within the scale');
  });

  it('clears a validation error when a different widget is selected', () => {
    const fixture = TestBed.createComponent(MetricWidgetPropertiesComponent);
    fixture.componentRef.setInput('widget', gauge('first'));
    fixture.detectChanges();
    fixture.componentInstance.error.set('The target must be positive and within the scale.');
    fixture.detectChanges();
    expect(fixture.nativeElement.textContent).toContain('within the scale');
    fixture.componentRef.setInput('widget', gauge('second'));
    fixture.detectChanges();
    expect(fixture.componentInstance.error()).toBe('');
    expect(fixture.nativeElement.textContent).not.toContain('within the scale');
  });
});

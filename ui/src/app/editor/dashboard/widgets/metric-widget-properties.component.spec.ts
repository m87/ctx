import { TestBed } from '@angular/core/testing';
import {
  createWidget,
  GaugeWidgetDefinition,
  RadarWidgetDefinition,
} from '../dashboard-definition';
import { MetricWidgetPropertiesComponent } from './metric-widget-properties.component';

function gauge(id: string): GaugeWidgetDefinition {
  const widget = createWidget('gauge', id, { x: 0, y: 0, width: 6, height: 5 });
  if (widget.type !== 'gauge') throw new Error('Expected gauge widget');
  return widget;
}

function radar(properties: Partial<RadarWidgetDefinition['properties']> = {}) {
  const widget = createWidget('radar', 'radar', { x: 0, y: 0, width: 6, height: 5 });
  if (widget.type !== 'radar') throw new Error('Expected radar widget');
  widget.properties = { ...widget.properties, ...properties };
  const fixture = TestBed.createComponent(MetricWidgetPropertiesComponent);
  fixture.componentRef.setInput('widget', widget);
  fixture.detectChanges();
  const changed = vi.fn();
  fixture.componentInstance.widgetChange.subscribe(changed);
  const element: HTMLElement = fixture.nativeElement;
  function control<T extends HTMLElement>(selector: string): T {
    const result = element.querySelector<T>(selector);
    if (!result) throw new Error(`Missing control: ${selector}`);
    return result;
  }
  function change(field: string, value: string) {
    const input = control<HTMLInputElement | HTMLSelectElement>(`#metric-radar-${field}`);
    input.value = value;
    input.dispatchEvent(new Event('change'));
    fixture.detectChanges();
  }
  function acceptChange() {
    expect(changed).toHaveBeenCalled();
    fixture.componentRef.setInput('widget', changed.mock.lastCall?.[0]);
    fixture.detectChanges();
  }
  return { fixture, element, control, changed, change, acceptChange };
}

describe('MetricWidgetPropertiesComponent', () => {
  afterEach(() => TestBed.resetTestingModule());

  it.each([
    ['kpi', 'KPI'],
    ['gauge', 'Gauge'],
    ['radar', 'Radar'],
  ] as const)('names the %s settings and renders only its own sections', (type, label) => {
    const fixture = TestBed.createComponent(MetricWidgetPropertiesComponent);
    fixture.componentRef.setInput(
      'widget',
      createWidget(type, 'widget', { x: 0, y: 0, width: 6, height: 5 }),
    );
    fixture.detectChanges();
    const element: HTMLElement = fixture.nativeElement;
    expect(element.querySelector('legend')?.textContent?.trim()).toBe(`${label} settings`);
    expect(!!element.querySelector('#metric-widget-axes-by')).toBe(type === 'radar');
    expect(!!element.querySelector('#metric-widget-goal-mode')).toBe(type === 'gauge');
  });

  it('switches radar grouping and clears axes selected for the previous group', () => {
    const { element, change, changed, acceptChange } = radar();
    expect(element.textContent).toContain('Each axis sums all observed occurrences');
    expect(element.querySelector('[role="group"]')).toBeNull();
    change('axes-by', 'project');
    expect(changed).toHaveBeenLastCalledWith(
      expect.objectContaining({
        properties: expect.objectContaining({ axesBy: 'project', axisKeys: [] }),
      }),
    );
    acceptChange();
    expect(element.querySelector('[aria-label="Radar project axes"]')?.textContent).toContain(
      'No project',
    );
    change('axes-by', 'context');
    acceptChange();
    expect(element.querySelector('[aria-label="Radar context axes"]')?.textContent).toContain(
      'Selected: 0',
    );
    expect(element.textContent).not.toContain('No project');
  });

  it('keeps unavailable axes visible and supports selecting, reordering, and removing axes', () => {
    const { fixture, element, control, changed, acceptChange, change } = radar({
      axesBy: 'project',
      axisKeys: ['project:a', 'project:missing'],
    });
    fixture.componentRef.setInput('projects', [{ value: 'a', label: 'Alpha' }]);
    fixture.detectChanges();
    expect(element.textContent).toContain('Unavailable category (project:missing)');
    expect(control<HTMLButtonElement>('[aria-label="Move Alpha axis up"]').disabled).toBe(true);
    control<HTMLButtonElement>('[aria-label="Move Alpha axis down"]').click();
    expect(changed).toHaveBeenLastCalledWith(
      expect.objectContaining({
        properties: expect.objectContaining({ axisKeys: ['project:missing', 'project:a'] }),
      }),
    );
    acceptChange();
    const missing = control<HTMLInputElement>(
      '[aria-label="Radar project axes"] > div:nth-child(3) input',
    );
    missing.checked = false;
    missing.dispatchEvent(new Event('change'));
    acceptChange();
    const noProject = control<HTMLInputElement>('[aria-label="Radar project axes"] input');
    noProject.checked = true;
    noProject.dispatchEvent(new Event('change'));
    expect(changed).toHaveBeenLastCalledWith(
      expect.objectContaining({
        properties: expect.objectContaining({ axisKeys: ['project:a', 'project:none'] }),
      }),
    );
    acceptChange();
    change('axes-by', 'weekday');
    expect(changed).toHaveBeenLastCalledWith(
      expect.objectContaining({
        properties: expect.objectContaining({ axesBy: 'weekday', axisKeys: [] }),
      }),
    );
  });

  it('limits radar selection to eight axes while allowing selected axes to be removed', () => {
    const contexts = Array.from({ length: 9 }, (_, i) => ({
      value: String(i),
      label: `Context ${i}`,
    }));
    const { fixture, element } = radar({
      axesBy: 'context',
      axisKeys: contexts.slice(0, 8).map(({ value }) => `context:${value}`),
    });
    fixture.componentRef.setInput('contexts', contexts);
    fixture.detectChanges();
    const inputs = Array.from(element.querySelectorAll<HTMLInputElement>('[role="group"] input'));
    expect(inputs).toHaveLength(9);
    expect(inputs.slice(0, 8).every((input) => input.checked && !input.disabled)).toBe(true);
    expect(inputs[8].disabled).toBe(true);
  });

  it('converts radar hours to milliseconds and allows clearing scale and target', () => {
    const { control, change, changed, acceptChange } = radar({
      scaleMaximum: 7200000,
      targetValue: 3600000,
    });
    expect(control<HTMLInputElement>('#metric-radar-scale-maximum').value).toBe('2');
    expect(control<HTMLInputElement>('#metric-radar-target-value').value).toBe('1');
    change('target-value', '1.5');
    expect(changed).toHaveBeenLastCalledWith(
      expect.objectContaining({
        properties: expect.objectContaining({ targetValue: 5400000 }),
      }),
    );
    acceptChange();
    for (const field of ['scale-maximum', 'target-value']) {
      change(field, '');
      acceptChange();
    }
    expect(changed).toHaveBeenLastCalledWith(
      expect.objectContaining({
        properties: expect.objectContaining({ scaleMaximum: null, targetValue: null }),
      }),
    );
  });

  it('keeps radar validation and unit resets in the parent form', () => {
    const { element, change, changed, acceptChange } = radar({
      scaleMaximum: 7200000,
      targetValue: 3600000,
    });
    change('metric', 'sessions.count');
    expect(changed).toHaveBeenLastCalledWith(
      expect.objectContaining({
        properties: expect.objectContaining({
          scaleMaximum: null,
          targetValue: null,
          format: 'number',
        }),
      }),
    );
    expect(element.textContent).toContain('Scale and target reset for the new unit.');
    acceptChange();
    changed.mockClear();
    change('scale-maximum', '1.5');
    expect(changed).not.toHaveBeenCalled();
    expect(element.querySelector('[role="alert"]')?.textContent).toContain(
      'Invalid radar scaleMaximum',
    );
    change('scale-maximum', '10');
    expect(changed).toHaveBeenCalledTimes(1);
    expect(element.querySelector('[role="alert"]')).toBeNull();
  });

  it.each([
    ['showValues', 'Show values by axis'],
    ['showGrid', 'Show grid'],
    ['showFill', 'Fill complete profile'],
    ['showSummary', 'Show period summary'],
  ] as const)('updates radar %s through the parent form', (key, label) => {
    const { element, changed } = radar();
    const input = Array.from(element.querySelectorAll('label'))
      .find((element) => element.textContent?.includes(label))
      ?.querySelector('input');
    if (!input) throw new Error(`Missing checkbox: ${label}`);
    input.checked = !input.checked;
    input.dispatchEvent(new Event('change'));
    expect(changed).toHaveBeenLastCalledWith(
      expect.objectContaining({
        properties: expect.objectContaining({ [key]: input.checked }),
      }),
    );
  });

  it('preserves radar field identifiers and prevents updates when disabled', () => {
    const { fixture, control, changed } = radar();
    fixture.componentRef.setInput('idSuffix', '-sidebar');
    fixture.componentRef.setInput('disabled', true);
    fixture.detectChanges();
    const input = control<HTMLInputElement>('#metric-radar-target-value-sidebar');
    expect(
      control<HTMLLabelElement>('label[for="metric-radar-target-value-sidebar"]').htmlFor,
    ).toBe(input.id);
    expect(input.matches(':disabled')).toBe(true);
    input.value = '2';
    input.dispatchEvent(new Event('change'));
    expect(changed).not.toHaveBeenCalled();
  });

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

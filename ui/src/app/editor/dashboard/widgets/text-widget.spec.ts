import { TestBed } from '@angular/core/testing';
import { TextWidgetComponent } from './text-widget.component';
import { TextWidgetPropertiesComponent } from './text-widget-properties.component';
import { TextWidgetDefinition } from '../dashboard-definition';

const widget: TextWidgetDefinition = {
  id: 'note',
  type: 'text',
  query: 'preserve query',
  layout: { x: 0, y: 0, width: 4, height: 3 },
  properties: {
    text: 'Focus\nOne thing at a time',
    horizontalAlign: 'left',
    verticalAlign: 'top',
    fontSize: 14,
  },
};

describe('text widget presentation', () => {
  afterEach(() => TestBed.resetTestingModule());

  it('renders stored alignment, size and plain multiline text, with legacy defaults', () => {
    const fixture = TestBed.createComponent(TextWidgetComponent);
    fixture.componentRef.setInput('properties', { text: '<b>Note</b>\nSecond line' });
    fixture.detectChanges();
    const content = () => fixture.nativeElement.querySelector('span') as HTMLSpanElement;
    expect(content().style.fontSize).toBe('14px');
    expect(content().classList.contains('text-left')).toBe(true);
    expect(content().textContent).toContain('<b>Note</b>\nSecond line');
    expect(content().querySelector('b')).toBeNull();
    fixture.componentRef.setInput('properties', {
      text: 'Centered',
      horizontalAlign: 'center',
      verticalAlign: 'center',
      fontSize: 32,
    });
    fixture.detectChanges();
    expect(content().classList.contains('text-center')).toBe(true);
    expect(content().style.fontSize).toBe('32px');
    expect(fixture.nativeElement.querySelector('.items-center')).not.toBeNull();
    fixture.componentRef.setInput('properties', {
      text: 'Corner',
      horizontalAlign: 'right',
      verticalAlign: 'bottom',
      fontSize: 24,
    });
    fixture.detectChanges();
    expect(content().classList.contains('text-right')).toBe(true);
    expect(fixture.nativeElement.querySelector('.items-end')).not.toBeNull();
  });

  it('updates the draft through controls without changing other properties and rejects invalid sizes', () => {
    const fixture = TestBed.createComponent(TextWidgetPropertiesComponent);
    fixture.componentRef.setInput('widget', structuredClone(widget));
    fixture.detectChanges();
    const changed = vi.fn();
    fixture.componentInstance.propertiesChange.subscribe(changed);
    const horizontal = fixture.nativeElement.querySelector(
      'select[id*="horizontal"]',
    ) as HTMLSelectElement;
    horizontal.value = 'center';
    horizontal.dispatchEvent(new Event('change'));
    expect(changed).toHaveBeenLastCalledWith({ ...widget.properties, horizontalAlign: 'center' });
    expect(widget.properties.horizontalAlign).toBe('left');
    const size = fixture.nativeElement.querySelector('input[type="number"]') as HTMLInputElement;
    size.value = '40';
    size.dispatchEvent(new Event('input'));
    expect(changed).toHaveBeenLastCalledWith({ ...widget.properties, fontSize: 40 });
    size.value = '0';
    size.dispatchEvent(new Event('input'));
    fixture.detectChanges();
    expect(size.getAttribute('aria-invalid')).toBe('true');
    expect(changed).toHaveBeenCalledTimes(2);
    size.dispatchEvent(new Event('blur'));
    fixture.detectChanges();
    expect(size.value).toBe('14');
    fixture.componentRef.setInput('disabled', true);
    fixture.detectChanges();
    expect(size.disabled).toBe(true);
    expect(horizontal.disabled).toBe(true);
  });
});

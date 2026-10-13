import { Component } from '@angular/core';
import { TestBed } from '@angular/core/testing';
import { ChartOptions } from '@tanstack/charts/angular';
import { ChartFrameComponent } from './chart-frame.component';

type ObserverCallback = (entries: { contentRect: { height: number } }[]) => void;

@Component({
  imports: [ChartFrameComponent],
  template: `<ctx-chart-frame [chart]="chart" />`,
})
class HostComponent {
  chart = { definition: {}, ariaLabel: 'Test chart' } as unknown as ChartOptions;
}

describe('ChartFrameComponent', () => {
  let callback: ObserverCallback | null = null;
  let disconnected = false;

  beforeEach(() => {
    callback = null;
    disconnected = false;
    vi.stubGlobal(
      'ResizeObserver',
      class {
        constructor(cb: ObserverCallback) {
          callback = cb;
        }
        observe(): void {}
        disconnect(): void {
          disconnected = true;
        }
      },
    );
  });
  afterEach(() => {
    TestBed.resetTestingModule();
    vi.unstubAllGlobals();
  });

  it('waits for a measured height before building the chart options', async () => {
    const fixture = TestBed.createComponent(HostComponent);
    fixture.detectChanges();
    await fixture.whenStable();
    const frame = fixture.debugElement.children[0].componentInstance as ChartFrameComponent;
    expect(frame.options()).toBeNull();
    callback?.([{ contentRect: { height: 187.6 } }]);
    expect(frame.options()?.height).toBe(187);
    expect(frame.options()?.ariaLabel).toBe('Test chart');
    callback?.([{ contentRect: { height: 0 } }]);
    expect(frame.options()).toBeNull();
  });

  it('stops observing when destroyed', async () => {
    const fixture = TestBed.createComponent(HostComponent);
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.destroy();
    expect(disconnected).toBe(true);
  });
});

import { TestBed } from '@angular/core/testing';
import { DashboardCanvasComponent } from './dashboard-canvas.component';
import {
  addTextWidget,
  normalizeDashboardDefinition,
  TextWidgetDefinition,
} from './dashboard-definition';

describe('DashboardCanvasComponent dragging', () => {
  beforeEach(() => {
    vi.stubGlobal(
      'ResizeObserver',
      class {
        observe() {}
        disconnect() {}
      },
    );
  });

  afterEach(() => {
    TestBed.resetTestingModule();
    vi.unstubAllGlobals();
  });

  it('renders plain text without editor controls outside edit mode', async () => {
    const definition = addTextWidget(normalizeDashboardDefinition({}), 'note');
    const widget = {
      ...definition.widgets[0],
      properties: { text: 'Dashboard heading', fontSize: 32, horizontalAlign: 'center' as const },
    };
    const fixture = await createFixture([widget]);
    const canvas = fixture.componentInstance;
    fixture.componentRef.setInput('selectedWidgetId', widget.id);
    fixture.detectChanges();
    expect(fixture.nativeElement.querySelector('.dashboard-grid-editing')).not.toBeNull();
    expect(fixture.nativeElement.querySelectorAll('.dashboard-resize-handle')).toHaveLength(8);
    const selected = vi.fn();
    canvas.widgetSelected.subscribe(selected);
    fixture.componentRef.setInput('editing', false);
    fixture.detectChanges();
    expect(fixture.nativeElement.querySelector('.dashboard-grid-editing')).toBeNull();
    expect(fixture.nativeElement.querySelector('button')).toBeNull();
    expect(fixture.nativeElement.querySelector('.dashboard-resize-handle')).toBeNull();
    const content = fixture.nativeElement.querySelector('ctx-text-widget span') as HTMLSpanElement;
    expect(content.textContent).toContain('Dashboard heading');
    expect(content.style.fontSize).toBe('32px');
    expect(canvas.rowCount()).toBe(widget.layout.y + widget.layout.height);
    canvas.grid().nativeElement.dispatchEvent(new PointerEvent('pointerdown', { bubbles: true }));
    expect(selected).not.toHaveBeenCalled();
    expect(canvas.widgets()[0]).toEqual(widget);
  });

  it('previews a palette drop, adds at the dropped cell, and rejects collisions and outside drops', async () => {
    const fixture = await createFixture();
    const canvas = fixture.componentInstance;
    const added = vi.fn();
    canvas.widgetAdded.subscribe(added);
    canvas.startPaletteDrag();
    canvas.movePaletteDrag({ x: 325, y: 285 });
    expect(canvas.preview()?.layout).toEqual({ x: 4, y: 4, width: 4, height: 3 });
    expect(added).not.toHaveBeenCalled();
    canvas.finishPaletteDrag({ x: 325, y: 285 });
    expect(added).toHaveBeenCalledExactlyOnceWith({ x: 4, y: 4 });
    expect(canvas.preview()).toBeNull();

    canvas.startPaletteDrag();
    canvas.finishPaletteDrag({ x: 120, y: 80 });
    canvas.startPaletteDrag();
    canvas.finishPaletteDrag({ x: 950, y: 80 });
    expect(added).toHaveBeenCalledTimes(1);
  });

  it('snaps consecutive real CDK gestures and resets pixel transforms before committing', async () => {
    const fixture = await createFixture();
    const canvas = fixture.componentInstance;
    const button = fixture.nativeElement.querySelector('button') as HTMLButtonElement;
    const moved = vi.fn();
    canvas.widgetMoved.subscribe((change) => {
      expect(button.style.transform).toBe('');
      moved(change);
    });
    let rect = new DOMRect(110, 70, 180, 130);
    vi.spyOn(button, 'getBoundingClientRect').mockImplementation(() => rect);

    mouse(button, 'mousedown', 120, 80);
    mouse(document, 'mousemove', 130, 80);
    mouse(document, 'mousemove', 270, 180);
    fixture.detectChanges();
    expect(canvas.preview()?.layout).toMatchObject({ x: 3, y: 2 });
    expect(canvas.widgets()[0].layout).toMatchObject({ x: 0, y: 0 });
    mouse(document, 'mouseup', 270, 180);
    expect(moved).toHaveBeenLastCalledWith({
      id: 'one',
      layout: { x: 3, y: 2, width: 4, height: 3 },
    });

    fixture.componentRef.setInput('widgets', [
      { ...canvas.widgets()[0], layout: { x: 3, y: 2, width: 4, height: 3 } },
    ]);
    fixture.detectChanges();
    rect = new DOMRect(260, 170, 180, 130);
    mouse(button, 'mousedown', 270, 180);
    mouse(document, 'mousemove', 280, 180);
    mouse(document, 'mousemove', 320, 180);
    mouse(document, 'mouseup', 320, 180);
    expect(moved).toHaveBeenLastCalledWith({
      id: 'one',
      layout: { x: 4, y: 2, width: 4, height: 3 },
    });
    expect(moved).toHaveBeenCalledTimes(2);
  });

  it('rejects a colliding drag and clears its pixel transform', async () => {
    const definition = addTextWidget(addTextWidget(normalizeDashboardDefinition({}), 'one'), 'two');
    const fixture = await createFixture(definition.widgets);
    const canvas = fixture.componentInstance;
    const button = fixture.nativeElement.querySelector('button') as HTMLButtonElement;
    vi.spyOn(button, 'getBoundingClientRect').mockReturnValue(new DOMRect(110, 70, 180, 130));
    const moved = vi.fn();
    canvas.widgetMoved.subscribe(moved);
    mouse(button, 'mousedown', 120, 80);
    mouse(document, 'mousemove', 130, 80);
    mouse(document, 'mousemove', 320, 80);
    fixture.detectChanges();
    expect(canvas.preview()?.valid).toBe(false);
    expect(fixture.nativeElement.querySelector('.dashboard-widget-invalid')).not.toBeNull();
    mouse(document, 'mouseup', 320, 80);
    expect(moved).not.toHaveBeenCalled();
    expect(button.style.transform).toBe('');
    expect(canvas.preview()).toBeNull();
  });

  it('discards cancelled palette and widget gestures and prevents adding in view mode', async () => {
    const fixture = await createFixture();
    const canvas = fixture.componentInstance;
    const added = vi.fn();
    canvas.widgetAdded.subscribe(added);
    canvas.startPaletteDrag();
    canvas.movePaletteDrag({ x: 325, y: 285 });
    canvas.cancelDrag();
    canvas.finishPaletteDrag({ x: 325, y: 285 });
    expect(added).not.toHaveBeenCalled();
    expect(canvas.preview()).toBeNull();
    fixture.componentRef.setInput('editing', false);
    fixture.detectChanges();
    canvas.startPaletteDrag();
    canvas.finishPaletteDrag({ x: 325, y: 285 });
    expect(added).not.toHaveBeenCalled();
    expect(fixture.nativeElement.querySelector('button')).toBeNull();
  });

  it('renders eight handles only for a selected widget and commits a pointer resize once', async () => {
    const fixture = await createFixture();
    const canvas = fixture.componentInstance;
    expect(fixture.nativeElement.querySelectorAll('.dashboard-resize-handle')).toHaveLength(0);
    fixture.componentRef.setInput('selectedWidgetId', 'one');
    canvas.cellSize.set(50);
    fixture.detectChanges();
    expect(fixture.nativeElement.querySelectorAll('.dashboard-resize-handle')).toHaveLength(8);
    const handle = fixture.nativeElement.querySelector('[data-edge="se"]') as HTMLButtonElement;
    const capture = mockPointerCapture(handle);
    const resized = vi.fn();
    canvas.widgetResized.subscribe(resized);
    pointer(handle, 'pointerdown', 300, 210);
    expect(capture.set).toHaveBeenCalledWith(1);
    pointer(handle, 'pointermove', 374, 284);
    expect(canvas.preview()?.layout).toEqual({ x: 0, y: 0, width: 5, height: 4 });
    expect(canvas.widgets()[0].layout).toEqual({ x: 0, y: 0, width: 4, height: 3 });
    expect(resized).not.toHaveBeenCalled();
    pointer(handle, 'pointerup', 374, 284);
    expect(resized).toHaveBeenCalledExactlyOnceWith({
      id: 'one',
      layout: { x: 0, y: 0, width: 5, height: 4 },
    });
    expect(capture.release).toHaveBeenCalledWith(1);
    expect(canvas.preview()).toBeNull();
  });

  it('rejects resizing into another widget and restores the original preview after cancellation', async () => {
    const definition = addTextWidget(addTextWidget(normalizeDashboardDefinition({}), 'one'), 'two');
    const fixture = await createFixture(definition.widgets);
    const canvas = fixture.componentInstance;
    fixture.componentRef.setInput('selectedWidgetId', 'one');
    canvas.cellSize.set(50);
    fixture.detectChanges();
    const right = fixture.nativeElement.querySelector('[data-edge="e"]') as HTMLButtonElement;
    mockPointerCapture(right);
    const resized = vi.fn();
    canvas.widgetResized.subscribe(resized);
    pointer(right, 'pointerdown', 300, 120);
    pointer(right, 'pointermove', 350, 120);
    expect(canvas.preview()?.valid).toBe(false);
    pointer(right, 'pointerup', 350, 120);
    expect(resized).not.toHaveBeenCalled();
    const bottom = fixture.nativeElement.querySelector('[data-edge="s"]') as HTMLButtonElement;
    mockPointerCapture(bottom);
    pointer(bottom, 'pointerdown', 200, 210);
    pointer(bottom, 'pointermove', 200, 310);
    expect(canvas.preview()?.layout.height).toBe(5);
    pointer(bottom, 'pointercancel', 200, 310);
    expect(resized).not.toHaveBeenCalled();
    expect(canvas.preview()).toBeNull();
    expect(canvas.resizingWidgetId()).toBeNull();
  });

  it('resizes with arrows on a focused handle and discards a pointer resize on Escape', async () => {
    const fixture = await createFixture();
    const canvas = fixture.componentInstance;
    fixture.componentRef.setInput('selectedWidgetId', 'one');
    canvas.cellSize.set(50);
    fixture.detectChanges();
    const handle = fixture.nativeElement.querySelector('[data-edge="e"]') as HTMLButtonElement;
    const resized = vi.fn();
    canvas.widgetResized.subscribe(resized);
    handle.dispatchEvent(new KeyboardEvent('keydown', { key: 'ArrowRight', bubbles: true }));
    expect(resized).toHaveBeenCalledExactlyOnceWith({
      id: 'one',
      layout: { x: 0, y: 0, width: 5, height: 3 },
    });
    mockPointerCapture(handle);
    pointer(handle, 'pointerdown', 300, 120);
    expect(document.activeElement).toBe(handle);
    pointer(handle, 'pointermove', 350, 120);
    handle.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape', bubbles: true }));
    pointer(handle, 'pointerup', 350, 120);
    expect(resized).toHaveBeenCalledTimes(1);
    expect(canvas.preview()).toBeNull();
  });

  it('clears a resize preview when selection changes and removes its handles', async () => {
    const fixture = await createFixture();
    const canvas = fixture.componentInstance;
    fixture.componentRef.setInput('selectedWidgetId', 'one');
    canvas.cellSize.set(50);
    fixture.detectChanges();
    const handle = fixture.nativeElement.querySelector('[data-edge="s"]') as HTMLButtonElement;
    const capture = mockPointerCapture(handle);
    pointer(handle, 'pointerdown', 200, 210);
    pointer(handle, 'pointermove', 200, 310);
    expect(canvas.preview()).not.toBeNull();
    fixture.componentRef.setInput('selectedWidgetId', null);
    fixture.detectChanges();
    expect(canvas.preview()).toBeNull();
    expect(canvas.resizingWidgetId()).toBeNull();
    expect(capture.release).toHaveBeenCalledWith(1);
    expect(fixture.nativeElement.querySelector('.dashboard-widget-dragging')).toBeNull();
  });
});

async function createFixture(widgets?: TextWidgetDefinition[]) {
  TestBed.configureTestingModule({ imports: [DashboardCanvasComponent] });
  const fixture = TestBed.createComponent(DashboardCanvasComponent);
  fixture.componentRef.setInput(
    'widgets',
    widgets ?? addTextWidget(normalizeDashboardDefinition({}), 'one').widgets,
  );
  fixture.componentRef.setInput('editing', true);
  fixture.detectChanges();
  await fixture.whenStable();
  vi.spyOn(fixture.componentInstance.grid().nativeElement, 'getBoundingClientRect').mockReturnValue(
    new DOMRect(100, 60, 800, 1600),
  );
  return fixture;
}

function mouse(target: EventTarget, type: string, clientX: number, clientY: number): void {
  target.dispatchEvent(
    new MouseEvent(type, {
      bubbles: true,
      cancelable: true,
      button: 0,
      detail: 1,
      buttons: type === 'mouseup' ? 0 : 1,
      clientX,
      clientY,
    }),
  );
}

function mockPointerCapture(element: HTMLElement) {
  const set = vi.fn();
  const release = vi.fn();
  Object.defineProperties(element, {
    setPointerCapture: { value: set, configurable: true },
    hasPointerCapture: { value: () => true, configurable: true },
    releasePointerCapture: { value: release, configurable: true },
  });
  return { set, release };
}

function pointer(target: EventTarget, type: string, clientX: number, clientY: number): void {
  target.dispatchEvent(
    new PointerEvent(type, {
      bubbles: true,
      cancelable: true,
      pointerId: 1,
      button: 0,
      clientX,
      clientY,
    }),
  );
}

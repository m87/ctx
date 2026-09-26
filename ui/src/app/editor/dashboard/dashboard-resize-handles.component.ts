import { Component, DestroyRef, inject, input, output } from '@angular/core';
import { DashboardWidgetLayout } from './dashboard-definition';
import { GridPoint, WidgetResizeEdge, widgetResizeLayout } from './dashboard-layout';

interface ResizeGesture {
  pointerId: number;
  target: HTMLElement;
  edge: WidgetResizeEdge;
  start: GridPoint;
  layout: DashboardWidgetLayout;
  cellSize: number;
}

@Component({
  selector: 'ctx-dashboard-resize-handles',
  host: {
    class: 'dashboard-resize-handles',
    '(document:keydown.escape)': 'cancelResize()',
  },
  template: `
    @for (handle of handles; track handle.edge) {
      <button
        type="button"
        class="dashboard-resize-handle"
        [attr.data-edge]="handle.edge"
        [attr.aria-label]="'Resize ' + handle.label"
        [disabled]="disabled()"
        (pointerdown)="startResize($event, handle.edge)"
        (pointermove)="moveResize($event)"
        (pointerup)="finishResize($event)"
        (pointercancel)="cancelResize($event)"
        (lostpointercapture)="cancelResize($event)"
        (keydown)="resizeByKeyboard($event, handle.edge)"
      >
        <span aria-hidden="true"></span>
      </button>
    }
  `,
})
export class DashboardResizeHandlesComponent {
  readonly layout = input.required<DashboardWidgetLayout>();
  readonly cellSize = input.required<number>();
  readonly disabled = input(false);
  readonly resizeStarted = output<void>();
  readonly resizePreview = output<DashboardWidgetLayout>();
  readonly resizeEnded = output<DashboardWidgetLayout | null>();
  readonly handles: readonly { edge: WidgetResizeEdge; label: string }[] = [
    { edge: 'n', label: 'top edge' },
    { edge: 'ne', label: 'top right corner' },
    { edge: 'e', label: 'right edge' },
    { edge: 'se', label: 'bottom right corner' },
    { edge: 's', label: 'bottom edge' },
    { edge: 'sw', label: 'bottom left corner' },
    { edge: 'w', label: 'left edge' },
    { edge: 'nw', label: 'top left corner' },
  ];
  private gesture: ResizeGesture | null = null;

  constructor() {
    inject(DestroyRef).onDestroy(() => this.releasePointer());
  }

  startResize(event: PointerEvent, edge: WidgetResizeEdge): void {
    if (this.disabled() || this.gesture || event.button !== 0 || this.cellSize() <= 0) return;
    const target = event.currentTarget;
    if (!(target instanceof HTMLElement)) return;
    event.preventDefault();
    event.stopPropagation();
    target.setPointerCapture(event.pointerId);
    target.focus({ preventScroll: true });
    this.gesture = {
      pointerId: event.pointerId,
      target,
      edge,
      start: { x: event.clientX, y: event.clientY },
      layout: { ...this.layout() },
      cellSize: this.cellSize(),
    };
    this.resizeStarted.emit();
  }

  moveResize(event: PointerEvent): void {
    const layout = this.layoutAtPointer(event);
    if (!layout) return;
    event.preventDefault();
    this.resizePreview.emit(layout);
  }

  finishResize(event: PointerEvent): void {
    const layout = this.layoutAtPointer(event);
    if (!layout) return;
    event.preventDefault();
    event.stopPropagation();
    this.releasePointer();
    this.resizeEnded.emit(layout);
  }

  cancelResize(event?: PointerEvent): void {
    if (!this.gesture || (event && event.pointerId !== this.gesture.pointerId)) return;
    this.releasePointer();
    this.resizeEnded.emit(null);
  }

  resizeByKeyboard(event: KeyboardEvent, edge: WidgetResizeEdge): void {
    if (this.disabled() || this.gesture) return;
    const dx = event.key === 'ArrowLeft' ? -1 : event.key === 'ArrowRight' ? 1 : 0;
    const dy = event.key === 'ArrowUp' ? -1 : event.key === 'ArrowDown' ? 1 : 0;
    if (!dx && !dy) return;
    event.preventDefault();
    event.stopPropagation();
    const layout = widgetResizeLayout(this.layout(), edge, { x: dx, y: dy }, 1);
    this.resizeStarted.emit();
    this.resizeEnded.emit(layout);
  }

  private layoutAtPointer(event: PointerEvent): DashboardWidgetLayout | null {
    const gesture = this.gesture;
    if (!gesture || gesture.pointerId !== event.pointerId) return null;
    return widgetResizeLayout(
      gesture.layout,
      gesture.edge,
      { x: event.clientX - gesture.start.x, y: event.clientY - gesture.start.y },
      gesture.cellSize,
    );
  }

  private releasePointer(): void {
    const gesture = this.gesture;
    this.gesture = null;
    if (gesture?.target.hasPointerCapture(gesture.pointerId)) {
      gesture.target.releasePointerCapture(gesture.pointerId);
    }
  }
}

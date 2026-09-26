import {
  afterNextRender,
  Component,
  computed,
  DestroyRef,
  effect,
  ElementRef,
  inject,
  input,
  output,
  signal,
  viewChild,
} from '@angular/core';
import { CdkDrag, CdkDragEnd, CdkDragMove, DragConstrainPosition } from '@angular/cdk/drag-drop';
import {
  DashboardWidgetLayout,
  dashboardGridRowCount,
  TextWidgetDefinition,
  TEXT_WIDGET_CATALOG_ENTRY,
} from './dashboard-definition';
import {
  canPlaceWidget,
  DASHBOARD_COLUMNS,
  DASHBOARD_WIDGET_INSET,
  GridPoint,
  isPointerOnGrid,
  widgetDropLayout,
} from './dashboard-layout';
import { DashboardResizeHandlesComponent } from './dashboard-resize-handles.component';

interface WidgetDragPreview {
  widgetId: string | null;
  layout: DashboardWidgetLayout;
  text: string;
  valid: boolean;
}

@Component({
  selector: 'ctx-dashboard-canvas',
  imports: [CdkDrag, DashboardResizeHandlesComponent],
  host: {
    class: 'block w-full max-w-[1000px] shrink-0 rounded-lg border bg-card shadow-xs',
    '(document:keydown.escape)': 'cancelDrag()',
  },
  template: `
    <div
      #grid
      class="dashboard-grid"
      [style.--dashboard-widget-inset.px]="widgetInset"
      [style.grid-auto-rows.px]="cellSize()"
      [style.min-height.px]="rowCount() * cellSize()"
      [style.background-size]="cellSize() + 'px ' + cellSize() + 'px'"
      [attr.aria-label]="'Dashboard grid, 16 columns by ' + rowCount() + ' rows'"
      (pointerdown)="clearSelection($event)"
    >
      @for (widget of widgets(); track widget.id) {
        <div
          class="dashboard-widget-slot"
          [style.grid-column]="widget.layout.x + 1 + ' / span ' + widget.layout.width"
          [style.grid-row]="widget.layout.y + 1 + ' / span ' + widget.layout.height"
        >
          @if (editing()) {
            <button
              cdkDrag
              type="button"
              class="dashboard-widget-surface dashboard-widget-editable"
              [class.dashboard-widget-selected]="selectedWidgetId() === widget.id"
              [class.dashboard-widget-dragging]="
                (draggedWidget()?.id === widget.id && !cancelled) ||
                resizingWidgetId() === widget.id
              "
              [cdkDragConstrainPosition]="keepWidgetAtOrigin"
              [attr.aria-label]="'Text widget: ' + (widget.properties.text || 'Empty text widget')"
              [attr.aria-pressed]="selectedWidgetId() === widget.id"
              (click)="widgetSelected.emit(widget.id)"
              (keydown)="widgetKeydown.emit({ event: $event, widget })"
              (cdkDragStarted)="startWidgetDrag(widget)"
              (cdkDragMoved)="moveWidgetDrag($event)"
              (cdkDragEnded)="finishWidgetDrag($event)"
            >
              {{ widget.properties.text || 'Empty text widget' }}
            </button>
            @if (selectedWidgetId() === widget.id) {
              <ctx-dashboard-resize-handles
                [layout]="widget.layout"
                [cellSize]="cellSize()"
                [disabled]="draggedWidget() !== null"
                [class.opacity-0]="resizingWidgetId() === widget.id"
                (resizeStarted)="startWidgetResize(widget)"
                (resizePreview)="previewWidgetResize(widget, $event)"
                (resizeEnded)="finishWidgetResize(widget, $event)"
              ></ctx-dashboard-resize-handles>
            }
          } @else {
            <div class="dashboard-widget-surface">{{ widget.properties.text }}</div>
          }
        </div>
      }
      @if (preview(); as target) {
        <div
          class="dashboard-widget-slot dashboard-widget-preview"
          [style.grid-column]="target.layout.x + 1 + ' / span ' + target.layout.width"
          [style.grid-row]="target.layout.y + 1 + ' / span ' + target.layout.height"
          aria-hidden="true"
        >
          <div
            class="dashboard-widget-surface dashboard-widget-selected"
            [class.dashboard-widget-invalid]="!target.valid"
          >
            {{ target.text || 'Empty text widget' }}
          </div>
        </div>
      }
    </div>
  `,
})
export class DashboardCanvasComponent {
  readonly widgets = input.required<readonly TextWidgetDefinition[]>();
  readonly editing = input(false);
  readonly selectedWidgetId = input<string | null>(null);
  readonly widgetSelected = output<string | null>();
  readonly widgetAdded = output<GridPoint>();
  readonly widgetMoved = output<{ id: string; layout: DashboardWidgetLayout }>();
  readonly widgetResized = output<{ id: string; layout: DashboardWidgetLayout }>();
  readonly widgetKeydown = output<{ event: KeyboardEvent; widget: TextWidgetDefinition }>();
  readonly grid = viewChild.required<ElementRef<HTMLElement>>('grid');
  readonly widgetInset = DASHBOARD_WIDGET_INSET;
  readonly cellSize = signal(1);
  readonly draggedWidget = signal<TextWidgetDefinition | null>(null);
  readonly resizingWidgetId = signal<string | null>(null);
  readonly preview = signal<WidgetDragPreview | null>(null);
  readonly palettePreviewVisible = computed(() => this.preview()?.widgetId === null);
  readonly rowCount = computed(() => {
    const layout = this.preview()?.layout;
    return Math.max(dashboardGridRowCount(this.widgets()), layout ? layout.y + layout.height : 0);
  });
  cancelled = false;
  private pickup: GridPoint | undefined;
  private readonly destroyRef = inject(DestroyRef);
  private readonly resetResizeOnSelectionChange = effect(() => {
    const widgetId = this.resizingWidgetId();
    if (
      widgetId &&
      (!this.editing() ||
        this.selectedWidgetId() !== widgetId ||
        !this.widgets().some((widget) => widget.id === widgetId))
    ) {
      this.cancelDrag();
    }
  });

  readonly keepWidgetAtOrigin: DragConstrainPosition = (_point, _ref, dimensions, pickup) => {
    this.pickup = pickup;
    return { x: dimensions.x, y: dimensions.y };
  };

  constructor() {
    afterNextRender(() => {
      const grid = this.grid().nativeElement;
      this.cellSize.set(grid.getBoundingClientRect().width / DASHBOARD_COLUMNS);
      const observer = new ResizeObserver(([entry]) => {
        if (entry) this.cellSize.set(entry.contentRect.width / DASHBOARD_COLUMNS);
      });
      observer.observe(grid);
      this.destroyRef.onDestroy(() => observer.disconnect());
    });
  }

  startWidgetDrag(widget: TextWidgetDefinition): void {
    this.cancelled = false;
    this.pickup = undefined;
    this.draggedWidget.set(widget);
    this.widgetSelected.emit(widget.id);
    this.preview.set({
      widgetId: widget.id,
      layout: widget.layout,
      text: widget.properties.text,
      valid: true,
    });
  }

  moveWidgetDrag(event: CdkDragMove): void {
    const pointer = 'touches' in event.event ? event.event.touches[0] : event.event;
    if (pointer) this.updateWidgetPreview({ x: pointer.clientX, y: pointer.clientY });
  }

  finishWidgetDrag(event: CdkDragEnd): void {
    this.updateWidgetPreview(event.dropPoint);
    const target = this.preview();
    event.source.reset();
    this.preview.set(null);
    this.draggedWidget.set(null);
    if (!this.cancelled && event.event.type !== 'touchcancel' && target?.valid && target.widgetId) {
      this.widgetMoved.emit({ id: target.widgetId, layout: target.layout });
    }
    this.pickup = undefined;
  }

  startPaletteDrag(): void {
    this.cancelled = false;
    this.preview.set(null);
  }

  movePaletteDrag(pointer: GridPoint): void {
    if (!this.editing() || this.cancelled) return;
    const geometry = this.grid().nativeElement.getBoundingClientRect();
    if (!isPointerOnGrid(pointer, geometry)) {
      this.preview.set(null);
      return;
    }
    const layout = widgetDropLayout(pointer, geometry, {
      x: 0,
      y: 0,
      ...TEXT_WIDGET_CATALOG_ENTRY.defaultLayout,
    });
    this.preview.set({
      widgetId: null,
      layout,
      text: TEXT_WIDGET_CATALOG_ENTRY.createProperties().text,
      valid: this.widgets().length < 200 && canPlaceWidget(this.widgets(), layout),
    });
  }

  finishPaletteDrag(pointer: GridPoint | null): void {
    if (pointer) this.movePaletteDrag(pointer);
    const target = this.preview();
    this.preview.set(null);
    if (pointer && !this.cancelled && this.editing() && target?.valid) {
      this.widgetAdded.emit({ x: target.layout.x, y: target.layout.y });
    }
  }

  cancelDrag(): void {
    this.cancelled = true;
    this.preview.set(null);
    this.resizingWidgetId.set(null);
  }

  startWidgetResize(widget: TextWidgetDefinition): void {
    this.cancelled = false;
    this.resizingWidgetId.set(widget.id);
    this.previewWidgetResize(widget, widget.layout);
  }

  previewWidgetResize(widget: TextWidgetDefinition, layout: DashboardWidgetLayout): void {
    if (!this.editing() || this.cancelled) return;
    this.preview.set({
      widgetId: widget.id,
      layout,
      text: widget.properties.text,
      valid: canPlaceWidget(this.widgets(), layout, widget.id),
    });
  }

  finishWidgetResize(widget: TextWidgetDefinition, layout: DashboardWidgetLayout | null): void {
    this.preview.set(null);
    this.resizingWidgetId.set(null);
    if (
      layout &&
      !this.cancelled &&
      this.editing() &&
      canPlaceWidget(this.widgets(), layout, widget.id)
    ) {
      this.widgetResized.emit({ id: widget.id, layout });
    }
  }

  clearSelection(event: PointerEvent): void {
    if (event.target === this.grid().nativeElement) this.widgetSelected.emit(null);
  }

  private updateWidgetPreview(pointer: GridPoint): void {
    const widget = this.draggedWidget();
    if (!widget || !this.pickup || this.cancelled || !this.editing()) return;
    const layout = widgetDropLayout(
      pointer,
      this.grid().nativeElement.getBoundingClientRect(),
      widget.layout,
      this.pickup,
    );
    this.preview.set({
      widgetId: widget.id,
      layout,
      text: widget.properties.text,
      valid: canPlaceWidget(this.widgets(), layout, widget.id),
    });
  }
}

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
  DashboardWidgetDefinition,
  DashboardWidgetType,
  widgetCatalogEntry,
  createWidget,
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
import { DashboardWidgetComponent } from './widgets/dashboard-widget.component';
import { DashboardMetricResult } from '../../../api/dashboard/dashboard.service';

interface WidgetDragPreview {
  widgetId: string | null;
  layout: DashboardWidgetLayout;
  widget: DashboardWidgetDefinition;
  valid: boolean;
}

@Component({
  selector: 'ctx-dashboard-canvas',
  imports: [CdkDrag, DashboardResizeHandlesComponent, DashboardWidgetComponent],
  host: {
    class: 'block w-full max-w-[1000px] shrink-0',
    '(document:keydown.escape)': 'cancelDrag()',
  },
  template: `
    <div
      #grid
      class="dashboard-grid"
      [class.dashboard-grid-editing]="editing()"
      [class.dashboard-grid-stacked]="!editing() && narrow()"
      [style.--dashboard-widget-inset.px]="widgetInset"
      [style.grid-auto-rows.px]="!narrow() ? cellSize() : null"
      [style.min-height.px]="narrow() ? null : rowCount() * cellSize()"
      [style.min-width.px]="editing() ? 720 : null"
      [style.background-size]="cellSize() + 'px ' + cellSize() + 'px'"
      [attr.aria-label]="
        editing() ? 'Dashboard grid, 16 columns by ' + rowCount() + ' rows' : 'Dashboard content'
      "
      (pointerdown)="clearSelection($event)"
    >
      @for (widget of presentationWidgets(); track widget.id) {
        <div
          class="dashboard-widget-slot"
          [style.grid-column]="widget.layout.x + 1 + ' / span ' + widget.layout.width"
          [style.grid-row]="widget.layout.y + 1 + ' / span ' + widget.layout.height"
          [style.height.px]="!editing() && narrow() ? presentationHeight(widget) : null"
        >
          @if (editing()) {
            <div
              cdkDrag
              role="button"
              tabindex="0"
              class="dashboard-widget-surface dashboard-widget-editable"
              [class.dashboard-text-surface]="widget.type === 'text'"
              [class.dashboard-widget-selected]="selectedWidgetId() === widget.id"
              [class.dashboard-widget-dragging]="
                (draggedWidget()?.id === widget.id && !cancelled) ||
                resizingWidgetId() === widget.id
              "
              [cdkDragConstrainPosition]="keepWidgetAtOrigin"
              [attr.aria-label]="widgetLabel(widget)"
              [attr.aria-pressed]="selectedWidgetId() === widget.id"
              (click)="widgetSelected.emit(widget.id)"
              (keydown)="onWidgetKeydown($event, widget)"
              (cdkDragStarted)="startWidgetDrag(widget)"
              (cdkDragMoved)="moveWidgetDrag($event)"
              (cdkDragEnded)="finishWidgetDrag($event)"
            >
              <ctx-dashboard-widget
                [widget]="widget"
                [editing]="true"
                [result]="results()[widget.id]"
                [error]="error()"
                [firstDayOfWeek]="firstDayOfWeek()"
              />
            </div>
            @if (selectedWidgetId() === widget.id) {
              <ctx-dashboard-resize-handles
                [layout]="widget.layout"
                [minimum]="catalogEntry(widget.type).minimumLayout"
                [cellSize]="cellSize()"
                [disabled]="draggedWidget() !== null"
                [class.opacity-0]="resizingWidgetId() === widget.id"
                (resizeStarted)="startWidgetResize(widget)"
                (resizePreview)="previewWidgetResize(widget, $event)"
                (resizeEnded)="finishWidgetResize(widget, $event)"
              ></ctx-dashboard-resize-handles>
            }
          } @else {
            <div
              class="dashboard-widget-surface"
              [class.dashboard-text-surface]="widget.type === 'text'"
            >
              <ctx-dashboard-widget
                [widget]="widget"
                [result]="results()[widget.id]"
                [error]="error()"
                [firstDayOfWeek]="firstDayOfWeek()"
              />
            </div>
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
            [class.dashboard-text-surface]="target.widget.type === 'text'"
          >
            <ctx-dashboard-widget
              [widget]="target.widget"
              [editing]="true"
              [result]="results()[target.widget.id]"
              [error]="error()"
              [firstDayOfWeek]="firstDayOfWeek()"
            />
          </div>
        </div>
      }
    </div>
  `,
})
export class DashboardCanvasComponent {
  readonly widgets = input.required<readonly DashboardWidgetDefinition[]>();
  readonly editing = input(false);
  readonly results = input<Record<string, DashboardMetricResult>>({});
  readonly firstDayOfWeek = input<0 | 1>(1);
  readonly error = input('');
  readonly catalogEntry = widgetCatalogEntry;
  private paletteType: DashboardWidgetType = 'text';
  readonly selectedWidgetId = input<string | null>(null);
  readonly widgetSelected = output<string | null>();
  readonly widgetAdded = output<{ type: DashboardWidgetType; position: GridPoint }>();
  readonly widgetMoved = output<{ id: string; layout: DashboardWidgetLayout }>();
  readonly widgetResized = output<{ id: string; layout: DashboardWidgetLayout }>();
  readonly widgetKeydown = output<{ event: KeyboardEvent; widget: DashboardWidgetDefinition }>();
  readonly grid = viewChild.required<ElementRef<HTMLElement>>('grid');
  readonly widgetInset = DASHBOARD_WIDGET_INSET;
  readonly cellSize = signal(1);
  readonly narrow = signal(false);
  readonly draggedWidget = signal<DashboardWidgetDefinition | null>(null);
  readonly resizingWidgetId = signal<string | null>(null);
  readonly preview = signal<WidgetDragPreview | null>(null);
  readonly palettePreviewVisible = computed(() => this.preview()?.widgetId === null);
  readonly rowCount = computed(() => {
    if (!this.editing()) {
      return this.widgets().reduce(
        (rows, widget) => Math.max(rows, widget.layout.y + widget.layout.height),
        0,
      );
    }
    const layout = this.preview()?.layout;
    return Math.max(dashboardGridRowCount(this.widgets()), layout ? layout.y + layout.height : 0);
  });
  readonly presentationWidgets = computed(() => {
    const widgets = this.widgets();
    if (this.editing() || !this.narrow()) return widgets;
    return [...widgets].sort(
      (left, right) =>
        left.layout.y - right.layout.y ||
        left.layout.x - right.layout.x ||
        left.id.localeCompare(right.id),
    );
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
      const updateSize = (width: number) => {
        this.cellSize.set(width / DASHBOARD_COLUMNS);
        this.narrow.set(width < 600);
      };
      updateSize(grid.getBoundingClientRect().width);
      const observer = new ResizeObserver(([entry]) => {
        if (entry) updateSize(entry.contentRect.width);
      });
      observer.observe(grid);
      this.destroyRef.onDestroy(() => observer.disconnect());
    });
  }

  presentationHeight(widget: DashboardWidgetDefinition): number {
    return widget.type === 'gauge' || widget.type === 'radar'
      ? 240
      : widget.type === 'kpi'
        ? 150
        : 96;
  }

  startWidgetDrag(widget: DashboardWidgetDefinition): void {
    this.cancelled = false;
    this.pickup = undefined;
    this.draggedWidget.set(widget);
    this.widgetSelected.emit(widget.id);
    this.preview.set({
      widgetId: widget.id,
      layout: widget.layout,
      widget,
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

  startPaletteDrag(type: DashboardWidgetType): void {
    this.paletteType = type;
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
      ...widgetCatalogEntry(this.paletteType).defaultLayout,
    });
    this.preview.set({
      widgetId: null,
      layout,
      widget: createWidget(this.paletteType, 'palette-preview', layout),
      valid: this.widgets().length < 200 && canPlaceWidget(this.widgets(), layout),
    });
  }

  finishPaletteDrag(pointer: GridPoint | null): void {
    if (pointer) this.movePaletteDrag(pointer);
    const target = this.preview();
    this.preview.set(null);
    if (pointer && !this.cancelled && this.editing() && target?.valid) {
      this.widgetAdded.emit({
        type: this.paletteType,
        position: { x: target.layout.x, y: target.layout.y },
      });
    }
  }

  widgetLabel(widget: DashboardWidgetDefinition): string {
    return widget.type === 'text'
      ? 'Text widget: ' + widget.properties.text
      : widgetCatalogEntry(widget.type).label +
          ' widget: ' +
          (widget.properties.title ||
            ('metric' in widget.properties ? widget.properties.metric : ''));
  }

  cancelDrag(): void {
    this.cancelled = true;
    this.preview.set(null);
    this.resizingWidgetId.set(null);
  }

  startWidgetResize(widget: DashboardWidgetDefinition): void {
    this.cancelled = false;
    this.resizingWidgetId.set(widget.id);
    this.previewWidgetResize(widget, widget.layout);
  }

  previewWidgetResize(widget: DashboardWidgetDefinition, layout: DashboardWidgetLayout): void {
    if (!this.editing() || this.cancelled) return;
    this.preview.set({
      widgetId: widget.id,
      layout,
      widget,
      valid: canPlaceWidget(this.widgets(), layout, widget.id),
    });
  }

  finishWidgetResize(
    widget: DashboardWidgetDefinition,
    layout: DashboardWidgetLayout | null,
  ): void {
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
    if (this.editing() && event.target === this.grid().nativeElement)
      this.widgetSelected.emit(null);
  }

  onWidgetKeydown(event: KeyboardEvent, widget: DashboardWidgetDefinition): void {
    const target = event.target;
    if (
      target instanceof Element &&
      target !== event.currentTarget &&
      target.closest('button, a, input, select, textarea, [role="gridcell"], [tabindex="0"]')
    )
      return;
    this.widgetKeydown.emit({ event, widget });
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
      widget,
      valid: canPlaceWidget(this.widgets(), layout, widget.id),
    });
  }
}

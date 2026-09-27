import { Component, input, output } from '@angular/core';
import { CdkDrag, CdkDragEnd, CdkDragPreview, CdkDropList } from '@angular/cdk/drag-drop';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideType, lucideHash, lucideGauge } from '@ng-icons/lucide';
import { WIDGET_CATALOG, DashboardWidgetType } from './dashboard-definition';
import { GridPoint } from './dashboard-layout';
import { HlmButtonImports } from '@spartan-ng/helm/button';

@Component({
  selector: 'ctx-widget-palette',
  imports: [CdkDrag, CdkDragPreview, CdkDropList, NgIcon, HlmButtonImports],
  providers: [provideIcons({ lucideType, lucideHash, lucideGauge })],
  template: `
    <div cdkDropList [cdkDropListSortingDisabled]="true" class="grid grid-cols-2 gap-2 p-3">
      @for (entry of entries; track entry.type) {
        <button
          cdkDrag
          hlmBtn
          variant="ghost"
          [disabled]="disabled()"
          [cdkDragDisabled]="disabled() || !draggable()"
          [cdkDragData]="entry.type"
          type="button"
          class="dashboard-palette-item h-auto min-h-16 w-full flex-col items-start justify-center gap-1 text-foreground"
          [attr.aria-label]="'Add ' + entry.label + ' widget'"
          title="Drag onto the dashboard or click to add"
          (cdkDragStarted)="startDrag(entry.type)"
          (cdkDragMoved)="dragMoved.emit($event.pointerPosition)"
          (cdkDragEnded)="finishDrag($event)"
          (click)="addFromButton(entry.type)"
        >
          <ng-icon [name]="entry.icon"></ng-icon>
          <span class="text-xs font-medium">{{ entry.label }}</span>
          <span class="text-caption text-muted-foreground">{{ purpose(entry.type) }}</span>
          <ng-template cdkDragPreview>
            <div
              class="dashboard-palette-preview flex items-center gap-2 rounded-lg border bg-card p-3 text-card-foreground shadow-xs"
              [class.invisible]="overCanvas()"
            >
              <ng-icon [name]="entry.icon"></ng-icon>
              {{ entry.label }}
            </div>
          </ng-template>
        </button>
      }
    </div>
  `,
})
export class WidgetPaletteComponent {
  readonly entries = WIDGET_CATALOG;
  readonly overCanvas = input(false);
  readonly disabled = input(false);
  readonly draggable = input(true);
  readonly widgetAdded = output<DashboardWidgetType>();
  readonly dragStarted = output<DashboardWidgetType>();
  readonly dragMoved = output<GridPoint>();
  readonly dragEnded = output<GridPoint | null>();
  private dragging = false;
  private suppressClickUntil = 0;

  purpose(type: DashboardWidgetType): string {
    switch (type) {
      case 'text':
        return 'Add a note';
      case 'kpi':
        return 'Show one value';
      case 'gauge':
        return 'Track a goal';
    }
  }

  startDrag(type: DashboardWidgetType): void {
    this.dragging = true;
    this.dragStarted.emit(type);
  }

  finishDrag(event: CdkDragEnd): void {
    this.dragging = false;
    this.suppressClickUntil = Date.now() + 250;
    event.source.reset();
    this.dragEnded.emit(event.event.type === 'touchcancel' ? null : event.dropPoint);
  }

  addFromButton(type: DashboardWidgetType): void {
    if (!this.disabled() && !this.dragging && Date.now() >= this.suppressClickUntil)
      this.widgetAdded.emit(type);
  }
}

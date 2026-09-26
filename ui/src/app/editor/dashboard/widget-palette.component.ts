import { Component, input, output } from '@angular/core';
import { CdkDrag, CdkDragEnd, CdkDragPreview, CdkDropList } from '@angular/cdk/drag-drop';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideType } from '@ng-icons/lucide';
import { TEXT_WIDGET_CATALOG_ENTRY } from './dashboard-definition';
import { GridPoint } from './dashboard-layout';

@Component({
  selector: 'ctx-widget-palette',
  imports: [CdkDrag, CdkDragPreview, CdkDropList, NgIcon],
  providers: [provideIcons({ lucideType })],
  template: `
    <div cdkDropList [cdkDropListSortingDisabled]="true" class="flex flex-col gap-2 p-3">
      <button
        cdkDrag
        [cdkDragData]="entry.type"
        type="button"
        class="dashboard-palette-item flex items-center gap-3 rounded-lg border border-border/70 bg-card p-3 text-left text-card-foreground hover:border-border focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        aria-label="Add Text widget"
        aria-describedby="widget-drag-hint"
        (cdkDragStarted)="startDrag()"
        (cdkDragMoved)="dragMoved.emit($event.pointerPosition)"
        (cdkDragEnded)="finishDrag($event)"
        (click)="addFromButton()"
      >
        <span
          class="flex size-10 items-center justify-center rounded-md border bg-muted/30 text-muted-foreground"
        >
          <ng-icon [name]="entry.icon"></ng-icon>
        </span>
        <span class="min-w-0">
          <span class="block text-sm font-medium">{{ entry.label }}</span>
          <span class="block text-xs text-muted-foreground">{{ entry.description }}</span>
        </span>
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
      <p id="widget-drag-hint" class="text-xs text-muted-foreground">
        Drag onto the dashboard or press Enter to add.
      </p>
    </div>
  `,
})
export class WidgetPaletteComponent {
  readonly entry = TEXT_WIDGET_CATALOG_ENTRY;
  readonly overCanvas = input(false);
  readonly widgetAdded = output<void>();
  readonly dragStarted = output<void>();
  readonly dragMoved = output<GridPoint>();
  readonly dragEnded = output<GridPoint | null>();
  private dragging = false;
  private suppressClickUntil = 0;

  startDrag(): void {
    this.dragging = true;
    this.dragStarted.emit();
  }

  finishDrag(event: CdkDragEnd): void {
    this.dragging = false;
    this.suppressClickUntil = Date.now() + 250;
    event.source.reset();
    this.dragEnded.emit(event.event.type === 'touchcancel' ? null : event.dropPoint);
  }

  addFromButton(): void {
    if (!this.dragging && Date.now() >= this.suppressClickUntil) this.widgetAdded.emit();
  }
}

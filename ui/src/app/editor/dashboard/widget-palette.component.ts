import { Component, input, output } from '@angular/core';
import { CdkDrag, CdkDragEnd, CdkDragPreview, CdkDropList } from '@angular/cdk/drag-drop';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideType } from '@ng-icons/lucide';
import { TEXT_WIDGET_CATALOG_ENTRY } from './dashboard-definition';
import { GridPoint } from './dashboard-layout';
import { HlmButtonImports } from '@spartan-ng/helm/button';

@Component({
  selector: 'ctx-widget-palette',
  imports: [CdkDrag, CdkDragPreview, CdkDropList, NgIcon, HlmButtonImports],
  providers: [provideIcons({ lucideType })],
  template: `
    <div cdkDropList [cdkDropListSortingDisabled]="true" class="flex flex-col gap-2 p-3">
      <button
        cdkDrag
        hlmBtn
        variant="ghost"
        [disabled]="disabled()"
        [cdkDragDisabled]="disabled() || !draggable()"
        [cdkDragData]="entry.type"
        type="button"
        class="dashboard-palette-item w-full justify-start text-foreground"
        aria-label="Add Text widget"
        title="Drag onto the dashboard or click to add"
        (cdkDragStarted)="startDrag()"
        (cdkDragMoved)="dragMoved.emit($event.pointerPosition)"
        (cdkDragEnded)="finishDrag($event)"
        (click)="addFromButton()"
      >
        <ng-icon [name]="entry.icon"></ng-icon>
        Add {{ entry.label }}
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
    </div>
  `,
})
export class WidgetPaletteComponent {
  readonly entry = TEXT_WIDGET_CATALOG_ENTRY;
  readonly overCanvas = input(false);
  readonly disabled = input(false);
  readonly draggable = input(true);
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
    if (!this.disabled() && !this.dragging && Date.now() >= this.suppressClickUntil)
      this.widgetAdded.emit();
  }
}

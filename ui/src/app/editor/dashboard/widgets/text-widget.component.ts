import { Component, input } from '@angular/core';
import { TextWidgetProperties, TEXT_WIDGET_FONT_SIZE } from '../dashboard-definition';

@Component({
  selector: 'ctx-text-widget',
  host: { class: 'block h-full min-h-0 min-w-0' },
  template: `
    <div
      class="flex h-full min-h-0 min-w-0 overflow-hidden"
      [class.items-start]="!properties().verticalAlign || properties().verticalAlign === 'top'"
      [class.items-center]="properties().verticalAlign === 'center'"
      [class.items-end]="properties().verticalAlign === 'bottom'"
    >
      <span
        class="block max-h-full w-full min-w-0 overflow-hidden whitespace-pre-wrap break-words leading-relaxed"
        [class.text-left]="!properties().horizontalAlign || properties().horizontalAlign === 'left'"
        [class.text-center]="properties().horizontalAlign === 'center'"
        [class.text-right]="properties().horizontalAlign === 'right'"
        [style.font-size.px]="properties().fontSize ?? defaultFontSize"
      >
        @if (properties().text) {
          {{ properties().text }}
        } @else if (editing()) {
          <span class="text-muted-foreground">Empty text widget</span>
        }
      </span>
    </div>
  `,
})
export class TextWidgetComponent {
  readonly properties = input.required<TextWidgetProperties>();
  readonly editing = input(false);
  readonly defaultFontSize = TEXT_WIDGET_FONT_SIZE.default;
}

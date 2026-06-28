import { Component, computed, input } from '@angular/core';
import { ContextLinkRule, linkifyContextText } from './context-link-text';

@Component({
  selector: 'app-context-link-text',
  template: `
    @for (part of parts(); track $index) {
      @if (part.href) {
        <a
          class="text-primary underline underline-offset-2 hover:text-primary/80"
          [href]="part.href"
          target="_blank"
          rel="noopener noreferrer"
          (click)="stopParentClick($event)"
        >
          {{ part.text }}
        </a>
      } @else {
        <span>{{ part.text }}</span>
      }
    }
  `,
  styles: `
    :host {
      display: inline;
    }
  `,
})
export class ContextLinkTextComponent {
  readonly text = input<string | null | undefined>('');
  readonly rules = input<readonly ContextLinkRule[] | null | undefined>([]);

  readonly parts = computed(() => linkifyContextText(this.text(), this.rules()));

  stopParentClick(event: MouseEvent): void {
    event.stopPropagation();
  }
}

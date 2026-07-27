import { Component, computed, inject, input } from '@angular/core';
import { LinkRulesService } from './link-rules.service';

@Component({
  selector: 'ctx-linkified-text',
  template: `
    @for (part of parts(); track $index) {
      @if (part.href; as href) {
        <a
          class="text-primary underline decoration-primary/40 underline-offset-2 hover:decoration-primary"
          [href]="href"
          target="_blank"
          rel="noopener noreferrer"
          (click)="$event.stopPropagation()"
          >{{ part.text }}</a
        >
      } @else {
        {{ part.text }}
      }
    }
  `,
  styles: `
    :host {
      display: contents;
    }
  `,
})
export class LinkifiedTextComponent {
  private linkRules = inject(LinkRulesService);

  readonly text = input('');
  readonly parts = computed(() => this.linkRules.linkify(this.text()));
}

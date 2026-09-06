import { Component, input, output } from '@angular/core';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideX } from '@ng-icons/lucide';
import { HlmButtonImports } from '@spartan-ng/helm/button';

@Component({
  selector: 'ctx-query-form',
  imports: [HlmButtonImports, NgIcon],
  providers: [provideIcons({ lucideX })],
  template: `
    <form
      class="overflow-hidden rounded-xl border bg-card shadow-sm transition-shadow focus-within:border-ring focus-within:ring-2 focus-within:ring-ring/20"
      (submit)="submitQuery($event)"
    >
      <label for="context-query" class="sr-only">Context query</label>
      <div class="flex min-h-28 items-start gap-3 px-4 py-3.5">
        <textarea
          id="context-query"
          rows="3"
          class="min-h-20 min-w-0 flex-1 resize-y bg-transparent py-1 font-mono text-sm leading-6 text-foreground outline-none placeholder:text-muted-foreground/65"
          placeholder="Write a context query…"
          [value]="query()"
          (input)="updateQuery($event)"
          (keydown)="onQueryKeydown($event)"
          autocomplete="off"
          autocapitalize="off"
          aria-describedby="context-query-help"
          spellcheck="false"
        ></textarea>
        @if (query()) {
          <button
            type="button"
            class="mt-0.5 flex size-8 shrink-0 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
            aria-label="Clear query"
            title="Clear query"
            (click)="clear.emit()"
          >
            <ng-icon name="lucideX" class="text-sm"></ng-icon>
          </button>
        }
      </div>
      <div
        class="flex flex-col gap-3 border-t bg-muted/25 px-4 py-3 sm:flex-row sm:items-center sm:justify-between"
      >
        <div id="context-query-help" class="flex min-w-0 items-center gap-2 text-xs">
          <span
            class="shrink-0 rounded-full border bg-background px-2 py-0.5 font-medium text-muted-foreground"
          >
            Preview
          </span>
          <span class="truncate text-muted-foreground">
            Interpreter not connected — showing all contexts
          </span>
        </div>
        <div class="flex shrink-0 items-center gap-3">
          <span class="hidden text-[11px] text-muted-foreground sm:inline">
            Ctrl&nbsp;+&nbsp;Enter
          </span>
          <button hlmBtn type="submit" variant="outline" class="h-8 px-3 text-xs">Run</button>
        </div>
      </div>
    </form>
  `,
})
export class QueryFormComponent {
  readonly query = input('');
  readonly queryChange = output<string>();
  readonly run = output<void>();
  readonly clear = output<void>();

  updateQuery(event: Event): void {
    this.queryChange.emit((event.target as HTMLTextAreaElement).value);
  }

  submitQuery(event: Event): void {
    event.preventDefault();
    this.run.emit();
  }

  onQueryKeydown(event: KeyboardEvent): void {
    if (event.key !== 'Enter' || (!event.ctrlKey && !event.metaKey)) {
      return;
    }
    event.preventDefault();
    this.run.emit();
  }
}

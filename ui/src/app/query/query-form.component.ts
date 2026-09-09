import { Component, ElementRef, effect, input, output, signal, viewChild } from '@angular/core';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideSave, lucideX } from '@ng-icons/lucide';
import { HlmButtonImports } from '@spartan-ng/helm/button';

export const QUERY_SUGGESTIONS = [
  'and',
  'or',
  'in',
  '!=',
  '==',
  '~',
  '<',
  '>',
  '<=',
  '>=',
  'text',
  'name',
  'label',
  'project',
  'start',
  'end',
  'duration',
  'sessions',
] as const;

interface SuggestionContext {
  start: number;
  end: number;
  suggestions: readonly string[];
}

export function getQuerySuggestionContext(query: string, cursor: number): SuggestionContext {
  const safeCursor = Math.max(0, Math.min(cursor, query.length));
  if (isInsideQuotedText(query, safeCursor)) {
    return { start: safeCursor, end: safeCursor, suggestions: [] };
  }

  const beforeCursor = query.slice(0, safeCursor);
  const fragmentMatch = beforeCursor.match(/[A-Za-z]+$|[!<>=~]+$/);
  const fragment = fragmentMatch?.[0] ?? '';
  const start = safeCursor - fragment.length;
  const fragmentPattern = /^[A-Za-z]/.test(fragment) ? /[A-Za-z]/ : /[!<>=~]/;
  let end = safeCursor;

  while (end < query.length && fragment && fragmentPattern.test(query[end])) {
    end++;
  }

  return {
    start,
    end,
    suggestions: QUERY_SUGGESTIONS.filter((suggestion) =>
      suggestion.toLowerCase().startsWith(fragment.toLowerCase()),
    ),
  };
}

function isInsideQuotedText(query: string, cursor: number): boolean {
  let quote = '';
  let escaped = false;

  for (const character of query.slice(0, cursor)) {
    if (escaped) {
      escaped = false;
    } else if (character === '\\') {
      escaped = true;
    } else if (quote) {
      if (character === quote) {
        quote = '';
      }
    } else if (character === '"' || character === "'") {
      quote = character;
    }
  }

  return quote.length > 0;
}

@Component({
  selector: 'ctx-query-form',
  imports: [HlmButtonImports, NgIcon],
  providers: [provideIcons({ lucideSave, lucideX })],
  template: `
    <form
      class="rounded-xl border bg-card shadow-sm transition-shadow focus-within:border-ring focus-within:ring-2 focus-within:ring-ring/20"
      (submit)="submitQuery($event)"
    >
      <label for="context-query" class="sr-only">Context query</label>
      <div class="flex min-h-28 items-start gap-3 px-4 py-3.5">
        <div class="relative min-w-0 flex-1">
          <textarea
            #queryInput
            id="context-query"
            rows="3"
            class="block min-h-20 w-full resize-y bg-transparent py-1 font-mono text-sm leading-6 text-foreground outline-none placeholder:text-muted-foreground/65"
            placeholder="Write a context query…"
            [value]="query()"
            (input)="updateQuery($event)"
            (focus)="refreshSuggestions($event)"
            (click)="refreshSuggestions($event)"
            (keyup)="onQueryKeyup($event)"
            (keydown)="onQueryKeydown($event)"
            (blur)="closeSuggestions()"
            autocomplete="off"
            autocapitalize="off"
            aria-autocomplete="list"
            aria-controls="context-query-suggestions"
            [attr.aria-activedescendant]="activeSuggestionId()"
            [attr.aria-expanded]="suggestionsOpen()"
            aria-describedby="context-query-help"
            spellcheck="false"
          ></textarea>
          @if (suggestionsOpen()) {
            <div
              id="context-query-suggestions"
              role="listbox"
              aria-label="Query suggestions"
              class="absolute top-[calc(100%+0.375rem)] left-0 z-50 max-h-56 min-w-44 overflow-auto rounded-lg border border-border/70 bg-popover/95 p-1 text-popover-foreground shadow-md backdrop-blur-sm"
            >
              @for (suggestion of suggestions(); track suggestion; let suggestionIndex = $index) {
                <button
                  type="button"
                  role="option"
                  class="flex w-full rounded-md px-2.5 py-1.5 text-left font-mono text-sm outline-none hover:bg-muted"
                  [id]="suggestionId(suggestionIndex)"
                  [attr.aria-selected]="activeSuggestionIndex() === suggestionIndex"
                  [class.bg-muted]="activeSuggestionIndex() === suggestionIndex"
                  (mouseenter)="activeSuggestionIndex.set(suggestionIndex)"
                  (mousedown)="$event.preventDefault()"
                  (click)="selectSuggestion(suggestion)"
                >
                  {{ suggestion }}
                </button>
              }
            </div>
          }
        </div>
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
      <div class="rounded-b-xl border-t bg-muted/25 px-4 py-3">
        @if (saveExpanded()) {
          <div class="flex flex-col gap-2 sm:flex-row sm:items-center">
            <label for="saved-query-name" class="sr-only">Saved query name</label>
            <input
              #saveNameInput
              id="saved-query-name"
              type="text"
              class="h-8 min-w-0 flex-1 rounded-md border bg-background px-2.5 text-xs outline-none focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/30"
              placeholder="Query name"
              [value]="saveName()"
              [disabled]="savePending()"
              (input)="updateSaveName($event)"
              (keydown.enter)="confirmSave($event)"
              (keydown.escape)="saveCancel.emit()"
            />
            <div class="flex shrink-0 items-center gap-2">
              <button
                hlmBtn
                type="button"
                class="h-8 px-3 text-xs"
                [disabled]="!saveName().trim() || savePending()"
                (click)="saveConfirm.emit()"
              >
                Save
              </button>
              <button
                hlmBtn
                type="button"
                variant="ghost"
                class="size-8 px-0 text-muted-foreground"
                [disabled]="savePending()"
                aria-label="Cancel saving query"
                title="Cancel"
                (click)="saveCancel.emit()"
              >
                <ng-icon name="lucideX" class="text-sm"></ng-icon>
              </button>
            </div>
          </div>
        } @else {
          <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <div id="context-query-help" class="flex min-w-0 items-center gap-2 text-xs">
              <span
                class="shrink-0 rounded-full border bg-background px-2 py-0.5 font-medium text-muted-foreground"
              >
                Preview
              </span>
              <span class="truncate text-muted-foreground">
                Interpreter not connected — showing all workspace contexts
              </span>
            </div>
            <div class="flex shrink-0 items-center gap-2">
              <button
                hlmBtn
                type="button"
                variant="ghost"
                class="h-8 px-2.5 text-xs"
                [disabled]="!canSave()"
                (click)="saveStart.emit()"
              >
                <ng-icon name="lucideSave" class="text-xs"></ng-icon>
                Save query
              </button>
              <span class="hidden text-[11px] text-muted-foreground sm:inline">
                Ctrl&nbsp;+&nbsp;Enter
              </span>
              <button hlmBtn type="submit" variant="outline" class="h-8 px-3 text-xs">Run</button>
            </div>
          </div>
        }
      </div>
    </form>
  `,
})
export class QueryFormComponent {
  readonly query = input('');
  readonly canSave = input(false);
  readonly saveExpanded = input(false);
  readonly saveName = input('');
  readonly savePending = input(false);
  readonly queryChange = output<string>();
  readonly run = output<void>();
  readonly clear = output<void>();
  readonly saveStart = output<void>();
  readonly saveNameChange = output<string>();
  readonly saveConfirm = output<void>();
  readonly saveCancel = output<void>();
  readonly suggestions = signal<readonly string[]>([]);
  readonly suggestionsOpen = signal(false);
  readonly activeSuggestionIndex = signal(0);
  private readonly queryInput = viewChild<ElementRef<HTMLTextAreaElement>>('queryInput');
  private readonly saveNameInput = viewChild<ElementRef<HTMLInputElement>>('saveNameInput');
  private suggestionStart = 0;
  private suggestionEnd = 0;
  private readonly focusSaveNameEffect = effect(() => {
    if (this.saveExpanded()) {
      this.saveNameInput()?.nativeElement.focus();
    }
  });

  updateQuery(event: Event): void {
    const textarea = event.target as HTMLTextAreaElement;
    this.queryChange.emit(textarea.value);
    this.setSuggestions(textarea);
  }

  updateSaveName(event: Event): void {
    this.saveNameChange.emit((event.target as HTMLInputElement).value);
  }

  confirmSave(event: Event): void {
    event.preventDefault();
    if (this.saveName().trim() && !this.savePending()) {
      this.saveConfirm.emit();
    }
  }

  submitQuery(event: Event): void {
    event.preventDefault();
    this.run.emit();
  }

  onQueryKeydown(event: KeyboardEvent): void {
    if (event.key === 'Enter' && (event.ctrlKey || event.metaKey)) {
      event.preventDefault();
      this.closeSuggestions();
      this.run.emit();
      return;
    }

    if (!this.suggestionsOpen()) {
      return;
    }

    if (event.key === 'ArrowDown') {
      event.preventDefault();
      this.moveActiveSuggestion(1);
    } else if (event.key === 'ArrowUp') {
      event.preventDefault();
      this.moveActiveSuggestion(-1);
    } else if (event.key === 'Enter' || event.key === 'Tab') {
      event.preventDefault();
      this.selectSuggestion(this.suggestions()[this.activeSuggestionIndex()]);
    } else if (event.key === 'Escape') {
      event.preventDefault();
      this.closeSuggestions();
    }
  }

  onQueryKeyup(event: KeyboardEvent): void {
    if (['ArrowDown', 'ArrowUp', 'Enter', 'Tab', 'Escape'].includes(event.key)) {
      return;
    }
    this.setSuggestions(event.target as HTMLTextAreaElement);
  }

  refreshSuggestions(event: Event): void {
    this.setSuggestions(event.target as HTMLTextAreaElement);
  }

  closeSuggestions(): void {
    this.suggestionsOpen.set(false);
  }

  selectSuggestion(suggestion: string | undefined): void {
    const textarea = this.queryInput()?.nativeElement;
    if (!textarea || !suggestion) {
      return;
    }

    const nextQuery =
      textarea.value.slice(0, this.suggestionStart) +
      suggestion +
      textarea.value.slice(this.suggestionEnd);
    const nextCursor = this.suggestionStart + suggestion.length;
    textarea.value = nextQuery;
    textarea.focus();
    textarea.setSelectionRange(nextCursor, nextCursor);
    this.queryChange.emit(nextQuery);
    this.closeSuggestions();
  }

  suggestionId(index: number): string {
    return `context-query-suggestion-${index}`;
  }

  activeSuggestionId(): string | null {
    return this.suggestionsOpen() ? this.suggestionId(this.activeSuggestionIndex()) : null;
  }

  private setSuggestions(textarea: HTMLTextAreaElement): void {
    const context = getQuerySuggestionContext(
      textarea.value,
      textarea.selectionStart ?? textarea.value.length,
    );
    this.suggestionStart = context.start;
    this.suggestionEnd = context.end;
    this.suggestions.set(context.suggestions);
    this.activeSuggestionIndex.set(0);
    this.suggestionsOpen.set(context.suggestions.length > 0);
  }

  private moveActiveSuggestion(offset: number): void {
    const suggestionCount = this.suggestions().length;
    this.activeSuggestionIndex.update(
      (index) => (index + offset + suggestionCount) % suggestionCount,
    );
  }
}

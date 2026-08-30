import { Component, computed, effect, input, output, signal } from '@angular/core';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideCheck, lucidePencil, lucideX } from '@ng-icons/lucide';
import type { Tag } from '../../api/context/context.service';
import { LinkifiedTextComponent } from './linkified-text.component';

export interface NameSaveValue {
  name: string;
  description: string;
  tags?: Tag[];
}

export function parseTagNames(value: string): string[] {
  return [
    ...new Set(
      value
        .split(',')
        .map((tag) => tag.trim())
        .filter((tag) => tag.length > 0),
    ),
  ];
}

export function resolveTags(value: string, tags: readonly Tag[]): Tag[] {
  const existingTags = new Map(tags.map((tag) => [tag.name, tag]));
  return parseTagNames(value).map((name) => existingTags.get(name) ?? { id: '', name });
}

@Component({
  selector: 'ctx-name',
  imports: [NgIcon, LinkifiedTextComponent],
  providers: [provideIcons({ lucideCheck, lucidePencil, lucideX })],
  template: `
    <div class="w-full min-w-0">
      <div
        class="text-[11px] font-semibold text-muted-foreground flex items-center gap-2 uppercase"
      >
        @if (accentColor()) {
          <span class="w-2 h-2 rounded-full" [style.background-color]="accentColor()"></span>
        }
        <span>{{ label() }}</span>
      </div>

      @if (isEditing()) {
        <div class="mt-2 rounded-lg border bg-card p-3 md:p-4">
          <div class="grid gap-3">
            <label class="flex flex-col gap-1">
              <span
                class="text-[11px] uppercase tracking-[0.08em] text-muted-foreground font-semibold"
                >Name</span
              >
              <input
                class="w-full h-9 rounded-md border border-border bg-background px-3 text-sm transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring/50"
                [value]="editName()"
                (input)="editName.set(getInputValue($event))"
                (keydown.escape)="cancelEdit()"
                [placeholder]="namePlaceholder()"
              />
            </label>

            @if (showDescription()) {
              <label class="flex flex-col gap-1">
                <span
                  class="text-[11px] uppercase tracking-[0.08em] text-muted-foreground font-semibold"
                  >Description</span
                >
                <textarea
                  class="w-full min-h-24 rounded-md border border-border bg-background px-3 py-2 text-sm transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring/50"
                  [value]="editDescription()"
                  (input)="editDescription.set(getInputValue($event))"
                  (keydown.escape)="cancelEdit()"
                  [placeholder]="descriptionPlaceholder()"
                ></textarea>
              </label>
            }

            @if (showTags()) {
              <label class="flex flex-col gap-1">
                <span
                  class="text-[11px] uppercase tracking-[0.08em] text-muted-foreground font-semibold"
                  >Tags</span
                >
                <input
                  class="w-full h-9 rounded-md border border-border bg-background px-3 text-sm transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring/50"
                  [value]="editTagsInput()"
                  (input)="editTagsInput.set(getInputValue($event))"
                  (keydown.escape)="cancelEdit()"
                  [placeholder]="tagsPlaceholder()"
                />
              </label>

              <div class="flex flex-wrap items-center gap-1.5 min-h-5">
                @if (editTagsPreview().length > 0) {
                  @for (tag of editTagsPreview(); track tag) {
                    <span
                      class="text-[10px] md:text-[11px] font-medium text-primary bg-primary/10 px-2 py-1 rounded-md"
                    >
                      #{{ tag }}
                    </span>
                  }
                } @else {
                  <span class="text-xs text-muted-foreground">No tags yet</span>
                }
              </div>
            }
          </div>

          <div class="mt-3 flex justify-end gap-2">
            <button
              type="button"
              class="h-8 px-3 rounded-md border text-xs hover:bg-muted/60 flex items-center gap-1.5"
              (click)="cancelEdit()"
            >
              <ng-icon name="lucideX"></ng-icon>
              <span>Cancel</span>
            </button>
            <button
              type="button"
              class="h-8 px-3 rounded-md bg-primary text-primary-foreground text-xs hover:bg-primary/90 disabled:opacity-60 flex items-center gap-1.5"
              [disabled]="savePending() || editName().trim().length === 0"
              (click)="saveEdit()"
            >
              <ng-icon name="lucideCheck"></ng-icon>
              <span>Save</span>
            </button>
          </div>
        </div>
      } @else {
        <div class="mt-1 flex items-start gap-3">
          <div class="min-w-0 flex-1">
            @if (compact()) {
              <div class="text-sm font-medium truncate">
                <ctx-linkified-text [text]="name()" />
              </div>
            } @else {
              <h1 class="text-2xl font-semibold tracking-tight truncate">
                <ctx-linkified-text [text]="name()" />
              </h1>
            }
            @if (showDescription()) {
              @if (description()) {
                <p class="mt-1 whitespace-pre-wrap text-sm text-muted-foreground/90">
                  <ctx-linkified-text [text]="description() ?? ''" />
                </p>
              } @else {
                <p class="mt-1 text-sm text-muted-foreground">{{ emptyDescription() }}</p>
              }
            }

            @if (showTags() && tags().length > 0) {
              <div class="flex flex-wrap items-center gap-1.5 md:gap-2 mt-2">
                @for (tag of tags(); track tag.id || tag.name) {
                  <span
                    class="text-[10px] md:text-[11px] font-medium text-blue-600 bg-blue-50/80 px-2 py-1 rounded-md"
                  >
                    #{{ tag.name }}
                  </span>
                }
              </div>
            }
          </div>

          @if (!readonly()) {
            <button
              type="button"
              class="size-9 rounded-md border text-muted-foreground hover:text-foreground hover:bg-muted/60 flex items-center justify-center shrink-0"
              aria-label="Edit"
              title="Edit"
              (click)="startEdit()"
            >
              <ng-icon name="lucidePencil"></ng-icon>
            </button>
          }
        </div>
      }
    </div>
  `,
})
export class NameComponent {
  readonly label = input('Name');
  readonly name = input('');
  readonly description = input<string | null>('');
  readonly tags = input<readonly Tag[]>([]);
  readonly showDescription = input(true);
  readonly showTags = input(false);
  readonly savePending = input(false);
  readonly readonly = input(false);
  readonly compact = input(false);
  readonly accentColor = input<string | null>(null);
  readonly emptyDescription = input('No description');
  readonly namePlaceholder = input('Name');
  readonly descriptionPlaceholder = input('Description');
  readonly tagsPlaceholder = input('Comma separated');
  readonly save = output<NameSaveValue>();

  readonly isEditing = signal(false);
  readonly editName = signal('');
  readonly editDescription = signal('');
  readonly editTagsInput = signal('');

  readonly editTagsPreview = computed(() => parseTagNames(this.editTagsInput()));

  constructor() {
    effect(() => {
      if (this.isEditing()) {
        return;
      }

      this.editName.set(this.name());
      this.editDescription.set(this.description() ?? '');
      this.editTagsInput.set(
        this.tags()
          .map((tag) => tag.name)
          .join(', '),
      );
    });
  }

  startEdit(): void {
    if (this.readonly()) {
      return;
    }

    this.editName.set(this.name());
    this.editDescription.set(this.description() ?? '');
    this.editTagsInput.set(
      this.tags()
        .map((tag) => tag.name)
        .join(', '),
    );
    this.isEditing.set(true);
  }

  cancelEdit(): void {
    this.isEditing.set(false);
  }

  saveEdit(): void {
    if (this.readonly()) {
      return;
    }

    const name = this.editName().trim();
    if (!name) {
      return;
    }

    this.save.emit({
      name,
      description: this.editDescription().trim(),
      tags: this.showTags() ? resolveTags(this.editTagsInput(), this.tags()) : undefined,
    });
    this.isEditing.set(false);
  }

  getInputValue(event: Event): string {
    return (event.target as HTMLInputElement | HTMLTextAreaElement).value;
  }
}

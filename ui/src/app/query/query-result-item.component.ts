import { Component, input } from '@angular/core';
import { RouterLink } from '@angular/router';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideChevronRight, lucideFolder } from '@ng-icons/lucide';
import { Context } from '../../api/context/context.service';
import { colorHash } from '../utils';

@Component({
  selector: 'ctx-query-result-item',
  imports: [NgIcon, RouterLink],
  providers: [provideIcons({ lucideChevronRight, lucideFolder })],
  host: {
    class: 'block',
    role: 'listitem',
  },
  template: `
    <a
      class="group flex cursor-pointer items-center gap-3 rounded-lg border bg-card p-4 transition-colors hover:bg-muted/30"
      [routerLink]="['/context', context().id]"
    >
      <span class="size-2.5 shrink-0 rounded-sm" [style.background-color]="contextColor()"></span>
      <div class="min-w-0 flex-1">
        <div class="flex min-w-0 flex-wrap items-center gap-2">
          <span class="min-w-0 truncate text-sm font-medium">{{ context().name }}</span>
          @if (context().archived) {
            <span
              class="rounded border px-1.5 py-0.5 text-[10px] font-medium text-muted-foreground"
            >
              Archived
            </span>
          }
        </div>
        @if (context().description) {
          <p class="mt-1 line-clamp-2 text-xs text-muted-foreground">
            {{ context().description }}
          </p>
        }
        @if (context().project || (context().tags ?? []).length > 0) {
          <div class="mt-2 flex flex-wrap items-center gap-1.5">
            @if (context().project; as project) {
              <span
                class="inline-flex max-w-full items-center gap-1 rounded-md bg-primary/10 px-1.5 py-0.5 text-[10px] font-medium text-primary"
                [title]="project.name"
              >
                <ng-icon name="lucideFolder" class="shrink-0"></ng-icon>
                <span class="truncate">{{ project.name }}</span>
              </span>
            }
            @for (tag of context().tags ?? []; track tag.id || tag.name) {
              <span class="rounded-md bg-muted px-1.5 py-0.5 text-[10px] text-muted-foreground">
                {{ tag.name }}
              </span>
            }
          </div>
        }
      </div>
      <ng-icon
        name="lucideChevronRight"
        class="shrink-0 text-sm text-muted-foreground/70 transition-transform group-hover:translate-x-0.5"
      ></ng-icon>
    </a>
  `,
})
export class QueryResultItemComponent {
  readonly context = input.required<Context>();

  contextColor(): string {
    return colorHash(this.context().id);
  }
}

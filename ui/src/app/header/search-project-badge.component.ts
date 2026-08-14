import { Component, input } from '@angular/core';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideFolder } from '@ng-icons/lucide';
import { ProjectMetadata } from '../../api/context/context.service';

@Component({
  selector: 'ctx-search-project-badge',
  imports: [NgIcon],
  providers: [provideIcons({ lucideFolder })],
  host: {
    class: 'min-w-0 max-w-32',
  },
  template: `
    <span
      class="inline-flex max-w-full items-center gap-1 rounded-md bg-primary/10 px-1.5 py-0.5 text-primary"
      [title]="project().name"
    >
      <ng-icon name="lucideFolder" class="shrink-0 text-[10px]"></ng-icon>
      <span class="truncate">{{ project().name }}</span>
    </span>
  `,
})
export class SearchProjectBadgeComponent {
  readonly project = input.required<ProjectMetadata>();
}

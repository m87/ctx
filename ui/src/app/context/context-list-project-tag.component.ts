import { Component, inject, input } from '@angular/core';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideFolder } from '@ng-icons/lucide';
import { RouterLink } from '@angular/router';
import { Store } from '@ngxs/store';
import { ProjectMetadata } from '../../api/context/context.service';
import { LinkifiedTextComponent } from '../shared/linkified-text.component';
import { SelectProject } from '../sidebar/workspace.state';

@Component({
  selector: 'ctx-context-list-project-tag',
  imports: [LinkifiedTextComponent, NgIcon, RouterLink],
  providers: [provideIcons({ lucideFolder })],
  host: {
    class: 'ml-auto min-w-0 max-w-[50%] shrink-0',
  },
  template: `
    <span
      class="inline-flex max-w-full items-center cursor-pointer rounded-md bg-primary/10 px-2 py-0.5 font-medium text-primary transition-colors hover:bg-primary/15"
      [routerLink]="['/project', project().id]"
      [title]="project().name"
      (click)="selectProject($event)"
    >
      <ng-icon
        name="lucideFolder"
        class="mr-1.5 shrink-0 inline-flex items-center justify-center text-[11px]"
      ></ng-icon>
      <span class="truncate">
        <ctx-linkified-text [text]="project().name" />
      </span>
    </span>
  `,
})
export class ContextListProjectTagComponent {
  private readonly store = inject(Store);

  readonly project = input.required<ProjectMetadata>();

  selectProject(event: MouseEvent): void {
    event.stopPropagation();
    this.store.dispatch(new SelectProject(this.project().id));
  }
}

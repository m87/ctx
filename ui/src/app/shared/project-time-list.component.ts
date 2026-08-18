import { Component, inject, input } from '@angular/core';
import { RouterLink } from '@angular/router';
import { Store } from '@ngxs/store';
import type { ProjectMetadata } from '../../api/context/context.service';
import { SelectProject } from '../sidebar/workspace.state';
import { LinkifiedTextComponent } from './linkified-text.component';

export interface ProjectTimeListItem {
  id: string;
  name: string;
  duration: string;
  percentage: number;
  color: string;
  contextCount: number;
  project?: ProjectMetadata;
}

@Component({
  selector: 'ctx-project-time-list',
  imports: [LinkifiedTextComponent, RouterLink],
  template: `
    @if (items().length > 0) {
      <div class="flex flex-col gap-2">
        @for (item of items(); track item.id) {
          <div
            class="block rounded-lg border bg-card p-3 transition-colors hover:bg-muted/30"
            [class.cursor-pointer]="item.project"
            [routerLink]="item.project ? ['/project', item.project.id] : null"
            [attr.role]="item.project ? 'link' : null"
            [attr.tabindex]="item.project ? 0 : null"
            (click)="selectProject(item)"
          >
            <div class="flex items-center gap-2 mb-2">
              <span class="size-2 rounded-sm shrink-0" [style.background-color]="item.color"></span>
              <span class="text-sm font-medium flex-1 truncate">
                <ctx-linkified-text [text]="item.name" />
              </span>
              <span class="text-xs text-muted-foreground tabular-nums">{{ item.duration }}</span>
            </div>
            <div class="h-1.5 rounded bg-muted/40 overflow-hidden">
              <div
                class="h-full rounded"
                [style.width.%]="boundedPercentage(item.percentage)"
                [style.background-color]="item.color"
              ></div>
            </div>
            <div class="mt-2 text-[10px] text-muted-foreground">
              {{ item.contextCount }}
              {{ item.contextCount === 1 ? 'context' : 'contexts' }} ·
              {{ boundedPercentage(item.percentage).toFixed(1) }}%
            </div>
          </div>
        }
      </div>
    } @else if (emptyMessage()) {
      <p class="text-xs text-muted-foreground">{{ emptyMessage() }}</p>
    }
  `,
})
export class ProjectTimeListComponent {
  private readonly store = inject(Store);

  readonly items = input<readonly ProjectTimeListItem[]>([]);
  readonly emptyMessage = input('');

  selectProject(item: ProjectTimeListItem): void {
    if (item.project) {
      this.store.dispatch(new SelectProject(item.project.id));
    }
  }

  boundedPercentage(value: number): number {
    if (!Number.isFinite(value)) {
      return 0;
    }
    return Math.min(100, Math.max(0, value));
  }
}

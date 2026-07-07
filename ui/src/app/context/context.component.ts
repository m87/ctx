import { Component, computed, inject, signal } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucidePlay, lucideTrash2 } from '@ng-icons/lucide';
import { HlmButtonImports } from '@spartan-ng/helm/button';
import { HlmCardImports } from '@spartan-ng/helm/card';
import { map } from 'rxjs';
import { toSignal } from '@angular/core/rxjs-interop';
import { ContextQueries } from '../../api/context.quries';
import { injectMutation, injectQuery } from '@tanstack/angular-query-experimental';
import { ContextMutations } from '../../api/context.mutations';
import { durationAsH, durationAsM } from '../utils';
import { DateTime } from 'luxon';
import { Store } from '@ngxs/store';
import { WorkspaceState } from '../sidebar/workspace.state';
import { NameComponent, NameSaveValue } from '../shared/name.component';
import { ContextIntervalListComponent } from './context-interval-list.component';

@Component({
  imports: [NameComponent, ContextIntervalListComponent, NgIcon, HlmButtonImports, HlmCardImports],
  providers: [
    provideIcons({
      lucidePlay,
      lucideTrash2,
    }),
  ],
  selector: 'ctx-context',
  template: `
    <div
      class="w-full h-full overflow-hidden flex flex-col items-start justify-start p-4 md:p-6 gap-5 relative"
    >
      <div class="w-full flex flex-col md:flex-row justify-between items-start gap-4">
        <ctx-name
          class="w-full min-w-0"
          label="Context"
          accentColor="#d97706"
          [name]="context().name"
          [description]="context().description"
          [tags]="context().tags ?? []"
          [showTags]="true"
          namePlaceholder="Context name"
          descriptionPlaceholder="What this context is for"
          tagsPlaceholder="Comma separated"
          [savePending]="updateContextMutation.isPending()"
          (save)="saveContextName($event)"
        ></ctx-name>

        <div class="flex items-center gap-2 w-full md:w-auto flex-nowrap md:pt-5">
          <button
            hlmBtn
            variant="outline"
            class="h-9 px-3 text-xs bg-red-100/70 text-red-700"
            [disabled]="deleteContextMutation.isPending()"
            (click)="deleteContext()"
          >
            <ng-icon name="lucideTrash2"></ng-icon>
          </button>
          <button
            hlmBtn
            variant="outline"
            class="h-9 px-3 text-xs bg-blue-200/70 text-blue-600"
            (click)="startContext()"
          >
            <ng-icon name="lucidePlay"></ng-icon>
            <span class="font-semibold text-blue-600">Start</span>
          </button>
        </div>
      </div>

      <div class="flex w-full">
        <div class="w-full flex items-center justify-center gap-4">
          <div hlmCard class="w-full p-3 rounded-lg border">
            <h3
              class="text-[11px] uppercase tracking-[0.08em] text-muted-foreground font-semibold"
              hlmCardTitle
            >
              Total time
            </h3>
            <div class="text-lg font-semibold" hlmCardContet>
              {{ parseDuration(contextStats()?.totalDuration) }}
            </div>
          </div>
          <div hlmCard class="w-full p-3 rounded-lg border">
            <h3
              class="text-[11px] uppercase tracking-[0.08em] text-muted-foreground font-semibold"
              hlmCardTitle
            >
              Today
            </h3>
            <div class="text-lg font-semibold" hlmCardContet>
              {{ parseDuration(contextStats()?.duration) }}
            </div>
          </div>
          <div hlmCard class="w-full p-3 rounded-lg border">
            <h3
              class="text-[11px] uppercase tracking-[0.08em] text-muted-foreground font-semibold"
              hlmCardTitle
            >
              Sessions
            </h3>
            <div class="text-lg font-semibold" hlmCardContet>
              {{ contextStats()?.totalSessions }}
            </div>
          </div>
          <div hlmCard class="w-full p-3 rounded-lg border">
            <h3
              class="text-[11px] uppercase tracking-[0.08em] text-muted-foreground font-semibold"
              hlmCardTitle
            >
              Today sessions
            </h3>
            <div class="text-lg font-semibold" hlmCardContet>{{ contextStats()?.sessions }}</div>
          </div>
        </div>
      </div>
      <ctx-context-interval-list
        [contextId]="contextId()"
        [activeWorkspaceId]="activeWorkspaceId()"
        [contexts]="contexts()"
      ></ctx-context-interval-list>
    </div>
  `,
  styles: `
    :host {
      display: block;
      width: 100%;
      max-width: 1000px;
      height: 100%;
      min-height: 0;
    }
  `,
})
export class ContextComponent {
  private contextQueries = inject(ContextQueries);
  private contextMutations = inject(ContextMutations);
  private router = inject(Router);
  private store = inject(Store);
  readonly activeWorkspaceId = this.store.selectSignal(WorkspaceState.selectedWorkspaceId);

  switchContextMutation = injectMutation(() => this.contextMutations.switch());
  updateContextMutation = injectMutation(() => this.contextMutations.update());
  deleteContextMutation = injectMutation(() => this.contextMutations.delete());
  contextQuery = injectQuery(() => this.contextQueries.get(this.contextId()));
  contextsQuery = injectQuery(() => this.contextQueries.list(this.activeWorkspaceId()));
  context = computed(() => this.contextQuery.data()!);
  contextStatsQuery = injectQuery(() => this.contextQueries.stats(this.contextId(), this.today()));
  contextStats = computed(() => this.contextStatsQuery.data());
  today = signal(DateTime.local().toFormat('yyyy-MM-dd'));
  contexts = computed(() => this.contextsQuery.data() ?? []);

  route = inject(ActivatedRoute);
  readonly contextId = toSignal(this.route.paramMap.pipe(map((pm) => pm.get('id') ?? '')), {
    initialValue: '',
  });

  startContext() {
    this.switchContextMutation.mutate(this.context()!);
  }

  deleteContext() {
    const context = this.context();

    if (!context.id) {
      return;
    }

    if (!window.confirm(`Delete context "${context.name}"?`)) {
      return;
    }

    this.deleteContextMutation.mutate(context.id, {
      onSuccess: () => {
        this.router.navigate(['/day', this.today()]);
      },
    });
  }

  saveContextName(value: NameSaveValue): void {
    const context = this.context();
    this.updateContextMutation.mutate({
      id: context.id,
      context: {
        ...context,
        name: value.name,
        description: value.description,
        tags: value.tags ?? [],
      },
    });
  }

  parseDuration(duration: number | undefined): string {
    if (duration === undefined) {
      return '0h 0m';
    }
    return `${durationAsH(duration)}h ${durationAsM(duration)}m`;
  }
}

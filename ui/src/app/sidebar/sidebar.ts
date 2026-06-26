import { Component, computed, inject, signal } from '@angular/core';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideSettings } from '@ng-icons/lucide';
import { RouterLink, RouterLinkActive } from '@angular/router';
import { SidebarContextListComponent } from './sidebar-context-list.component';
import { SidebarSettingsModalComponent } from './sidebar-settings-modal.component';
import { SidebarStore } from './sidebar.store';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { VersionQueries } from '../../api/version.queries';
import { SidebarWorkspaceSelectComponent } from './sidebar-workspace-select.component';
import { Store } from '@ngxs/store';
import { WorkspaceState } from './workspace.state';

@Component({
  selector: 'app-sidebar',
  imports: [SidebarContextListComponent, SidebarWorkspaceSelectComponent, RouterLink, RouterLinkActive, NgIcon, SidebarSettingsModalComponent],
  providers: [provideIcons({ lucideSettings })],
  template: ` <div class="h-full w-full min-h-0 flex flex-col">
    <div class="flex-1 min-h-0 flex flex-col border-b bg-sidebar">
      <div class="flex flex-col gap-2.5 p-2.5 border-b">
        <app-sidebar-workspace-select></app-sidebar-workspace-select>
      </div>
      <div class="flex flex-col gap-2.5 p-2.5 border-b">
        <div class="flex-1 min-h-0 flex flex-col gap-1.5 p-1">
          <a
            routerLink="/day"
            routerLinkActive="bg-muted text-foreground"
            class="uppercase flex justify-between items-center text-[11px] tracking-[0.08em] text-muted-foreground px-2 py-1 font-semibold rounded-md hover:bg-muted/50 cursor-pointer"
            (click)="sidebar.closeMobile()"
          >
            daily summary
          </a>
          <a
            [routerLink]="workspaceLink()"
            routerLinkActive="bg-muted text-foreground"
            class="uppercase flex justify-between items-center text-[11px] tracking-[0.08em] text-muted-foreground px-2 py-1 font-semibold rounded-md hover:bg-muted/50 cursor-pointer"
            (click)="sidebar.closeMobile()"
          >
            workspace
          </a>
          <!-- <div
            class="uppercase flex justify-between items-center text-sm text-muted-foreground p-1 font-semibold"
          >
            stats
          </div>
          <div
            class="uppercase flex justify-between items-center text-sm p-1 font-semibold text-blue-500 bg-blue-50 rounded-sm  "
          >
            overview
          </div> -->
        </div>
      </div>
      <div class="min-h-0 flex-1 overflow-auto">
        <app-sidebar-context-list></app-sidebar-context-list>
      </div>
      <div class="border-t mt-auto shrink-0 px-3 py-2 flex items-center justify-between">
        <span class="text-[11px] text-muted-foreground/70 tracking-[0.06em]">v{{ appVersion() }}</span>
        <button
          type="button"
          class="h-7 w-7 rounded-md text-muted-foreground hover:text-foreground hover:bg-muted/60 flex items-center justify-center"
          aria-label="Open settings"
          (click)="openSettings()"
        >
          <ng-icon name="lucideSettings" class="text-[15px]"></ng-icon>
        </button>
      </div>
    </div>

    <app-sidebar-settings-modal [open]="isSettingsOpen()" (openChange)="isSettingsOpen.set($event)">
    </app-sidebar-settings-modal>
  </div>`,
})
export class SidebarComponent {
  private versionQueries = inject(VersionQueries);
  private store = inject(Store);

  versionQuery = injectQuery(() => this.versionQueries.version());
  readonly appVersion = computed(() => this.versionQuery.data()?.version ?? 'dev');
  readonly isSettingsOpen = signal<boolean>(false);
  readonly selectedWorkspaceId = this.store.selectSignal(WorkspaceState.selectedWorkspaceId);
  readonly workspaceLink = computed(() => {
    const workspaceId = this.selectedWorkspaceId();
    return workspaceId ? ['/workspace', workspaceId] : ['/workspace'];
  });

  constructor(public sidebar: SidebarStore) {}

  openSettings(): void {
    this.isSettingsOpen.set(true);
  }
}

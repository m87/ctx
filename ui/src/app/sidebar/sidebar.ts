import { Component, computed, inject, signal } from '@angular/core';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideSettings } from '@ng-icons/lucide';
import { RouterLink, RouterLinkActive } from '@angular/router';
import { SidebarContextListComponent } from './sidebar-context-list.component';
import { SidebarSettingsModalComponent } from './sidebar-settings-modal.component';
import { SidebarStore } from './sidebar.store';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { VersionQueries } from '../../api/version.queries';

@Component({
  selector: 'app-sidebar',
  imports: [SidebarContextListComponent, RouterLink, RouterLinkActive, NgIcon, SidebarSettingsModalComponent],
  providers: [provideIcons({ lucideSettings })],
  template: ` <div class="h-full w-full min-h-0 flex flex-col">
    <div class="flex-1 min-h-0 flex flex-col border-b">
      <div class="flex flex-col gap-2 p-2.5 border-b">
        <!-- <div class="uppercase flex justify-between items-center text-sm text-muted-foreground">
          <span>Workspace</span><ng-icon name="lucidePlus" class="cursor-pointer"></ng-icon>
        </div>
        <div class="font-bold text-blue-500 bg-blue-50 rounded-sm p-1">personal</div>
        <div class="p-1">work</div>
      </div> -->
        <div class="flex-1 min-h-0 flex flex-col gap-1.5 p-1">
          <a
            routerLink="/day"
            routerLinkActive="bg-muted text-foreground"
            class="uppercase flex justify-between items-center text-[11px] tracking-[0.08em] text-muted-foreground px-2 py-1 font-semibold rounded-md hover:bg-muted/50 cursor-pointer"
            (click)="sidebar.closeMobile()"
          >
            day
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

  versionQuery = injectQuery(() => this.versionQueries.version());
  readonly appVersion = computed(() => this.versionQuery.data()?.version ?? 'dev');
  readonly isSettingsOpen = signal<boolean>(false);

  constructor(public sidebar: SidebarStore) {}

  openSettings(): void {
    this.isSettingsOpen.set(true);
  }
}

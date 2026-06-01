import { Component } from '@angular/core';
import { RouterLink, RouterLinkActive } from '@angular/router';
import { SidebarContextListComponent } from './sidebar-context-list.component';
import { SidebarStore } from './sidebar.store';
import { SidebarWorkspaceSelectComponent } from './sidebar-workspace-select.component';

@Component({
  selector: 'app-sidebar',
  imports: [SidebarContextListComponent, SidebarWorkspaceSelectComponent, RouterLink, RouterLinkActive],
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
      <!-- <div
        class="border-t mt-auto shrink-0 px-3 py-2 text-xs text-muted-foreground uppercase tracking-[0.08em] hover:text-foreground hover:bg-muted/40 cursor-pointer"
      >
        settings
      </div> -->
    </div>
  </div>`,
})
export class SidebarComponent {
  constructor(public sidebar: SidebarStore) { }
}

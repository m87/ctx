import { Component, inject } from '@angular/core';
import { HeaderComponent } from './header/header';
import { SidebarComponent } from './sidebar/sidebar';
import { MainComponent } from './main/main';
import { SidebarStore } from './sidebar/sidebar.store';

@Component({
  selector: 'app-root',
  imports: [HeaderComponent, SidebarComponent, MainComponent],
  template: `
    <div class="flex justify-center w-full h-screen bg-background text-foreground">
      <div class="w-full h-full flex flex-col">
        <app-header></app-header>
        <div class="relative flex flex-1 min-h-0" [class.sidebar-collapsed]="sidebar.collapsed()">
          <div
            class="md:hidden absolute inset-0 bg-black/40 z-30 transition-opacity"
            [class.opacity-100]="sidebar.mobileOpen()"
            [class.opacity-0]="!sidebar.mobileOpen()"
            [class.pointer-events-none]="!sidebar.mobileOpen()"
            (click)="sidebar.closeMobile()"
          ></div>

          <div
            class="md:hidden absolute left-0 top-0 bottom-0 z-40 w-64 border-r bg-card transition-transform duration-200"
            [class.translate-x-0]="sidebar.mobileOpen()"
            [class.-translate-x-full]="!sidebar.mobileOpen()"
          >
            <app-sidebar></app-sidebar>
          </div>

          <div
            class="hidden md:block w-60 h-full border-r bg-card/40 sidebar-collapsed:w-16 transition-all duration-200"
          >
            <app-sidebar></app-sidebar>
          </div>
          <div class="flex-1 flex flex-col h-full min-h-0">
            <main class="w-full flex-1 min-h-0">
              <app-main class="w-full"></app-main>
            </main>
          </div>
        </div>
      </div>
    </div>
  `,
})
export class App {
  sidebar = inject(SidebarStore);
}

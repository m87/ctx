import { Component } from '@angular/core';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucidePackage } from '@ng-icons/lucide';

@Component({
  selector: 'app-title',
  imports: [NgIcon],
  providers: [provideIcons({ lucidePackage })],
  template: /* html */ `
    <div class="flex items-center justify-between p-4 sidebar-collapsed:justify-center sidebar-collapsed:p-2">
      <div class="flex items-center gap-2 font-bold text-lg justify-start h-full sidebar-collapsed:flex-col sidebar-collapsed:gap-0">
        <span class="flex items-center sidebar-collapsed:hidden"><ng-icon name="lucidePackage" size="24px" /></span>
        <span class="hidden sidebar-collapsed:block"><ng-icon name="lucidePackage" size="18px" /></span>
        <span class="sidebar-collapsed:text-xs">inv</span>
      </div>
      <div class="text-sm text-muted-foreground sidebar-collapsed:hidden">
        <span>0.0.1</span>
      </div>
    </div>
  `,
})
export class TitleComponent {}

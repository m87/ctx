import { Component, EventEmitter, Input, Output, signal } from '@angular/core';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideX } from '@ng-icons/lucide';
import { SidebarSettingsDataIntegritySectionComponent } from './sidebar-settings-data-integrity-section.component';
import { SidebarSettingsGeneralSectionComponent } from './sidebar-settings-general-section.component';
import { SidebarSettingsLinkRulesSectionComponent } from './sidebar-settings-link-rules-section.component';

const settingsSections = ['General', 'Link rules', 'Data integrity'] as const;

type SettingsSection = (typeof settingsSections)[number];

@Component({
  selector: 'ctx-sidebar-settings-modal',
  imports: [
    NgIcon,
    SidebarSettingsGeneralSectionComponent,
    SidebarSettingsLinkRulesSectionComponent,
    SidebarSettingsDataIntegritySectionComponent,
  ],
  providers: [provideIcons({ lucideX })],
  template: `
    @if (open) {
      <div
        class="fixed inset-0 z-50 w-screen flex items-end sm:items-center sm:justify-center"
        (click)="requestClose()"
      >
        <div class="absolute inset-0 bg-background/60 backdrop-blur-[1px]"></div>
        <div
          class="relative w-full sm:w-[min(94vw,1040px)] h-[92vh] sm:h-auto max-h-[92vh] sm:max-h-[720px] bg-popover text-popover-foreground border rounded-t-2xl sm:rounded-xl shadow-xl flex flex-col sm:flex-row min-h-0"
          (click)="$event.stopPropagation()"
        >
          <div class="sm:w-60 border-b sm:border-b-0 sm:border-r p-3 sm:p-4 shrink-0">
            <div
              class="text-[11px] uppercase tracking-[0.08em] text-muted-foreground px-2 py-1 hidden sm:block"
            >
              Settings
            </div>
            <div class="flex items-center gap-2 sm:block">
              <div class="flex-1 flex sm:flex-col gap-1 overflow-x-auto sm:overflow-visible">
                @for (section of settingsSections; track section) {
                  <button
                    type="button"
                    class="text-left px-3 sm:px-2 py-2 sm:py-1.5 text-[14px] sm:text-[13px] rounded-md hover:bg-muted/60 whitespace-nowrap min-w-fit"
                    [class.bg-muted]="activeSettingsSection() === section"
                    [class.text-foreground]="activeSettingsSection() === section"
                    [class.text-muted-foreground]="activeSettingsSection() !== section"
                    (click)="activeSettingsSection.set(section)"
                  >
                    {{ section }}
                  </button>
                }
              </div>
              <button
                type="button"
                class="h-9 w-9 rounded-md text-muted-foreground hover:text-foreground hover:bg-muted/60 flex items-center justify-center sm:hidden"
                aria-label="Close settings"
                (click)="requestClose()"
              >
                <ng-icon name="lucideX" class="text-[14px]"></ng-icon>
              </button>
            </div>
          </div>
          <div class="flex-1 min-h-0 flex flex-col">
            <div
              class="h-14 sm:h-14 border-b px-5 sm:px-7 flex items-center justify-between hidden sm:flex"
            >
              <div class="flex items-center gap-2">
                <div class="font-semibold text-[15px]">{{ activeSettingsSection() }}</div>
                @if (activeSettingsSection() === 'Link rules') {
                  <span
                    class="rounded-full border border-border bg-muted/50 px-2 py-0.5 text-[10px] font-medium uppercase tracking-[0.06em] text-muted-foreground"
                  >
                    Per workspace
                  </span>
                }
              </div>
              <button
                type="button"
                class="h-9 w-9 rounded-md text-muted-foreground hover:text-foreground hover:bg-muted/60 flex items-center justify-center"
                aria-label="Close settings"
                (click)="requestClose()"
              >
                <ng-icon name="lucideX" class="text-[14px]"></ng-icon>
              </button>
            </div>
            <div class="p-5 sm:p-7 overflow-auto text-[14px] text-muted-foreground">
              <div class="max-w-[760px] pb-5">
                @switch (activeSettingsSection()) {
                  @case ('General') {
                    <ctx-sidebar-settings-general-section />
                  }
                  @case ('Link rules') {
                    <ctx-sidebar-settings-link-rules-section />
                  }
                  @case ('Data integrity') {
                    <ctx-sidebar-settings-data-integrity-section />
                  }
                }
              </div>
            </div>
          </div>
        </div>
      </div>
    }
  `,
})
export class SidebarSettingsModalComponent {
  @Input() open = false;
  @Output() openChange = new EventEmitter<boolean>();

  readonly settingsSections = settingsSections;
  readonly activeSettingsSection = signal<SettingsSection>('General');

  requestClose(): void {
    this.openChange.emit(false);
  }
}

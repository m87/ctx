import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { TimelineComponent } from './timeline.component';

@Component({
  selector: 'app-main',
  imports: [RouterOutlet, TimelineComponent],
  template: `
    <div class="w-full h-full min-h-0 flex flex-col overflow-hidden">
      <div
        class="w-full flex-1 min-h-0 overflow-hidden flex items-stretch justify-center pb-24 md:pb-0"
      >
        <router-outlet></router-outlet>
      </div>
      <div
        class="w-full mt-auto shrink-0 z-20 bg-background pb-[env(safe-area-inset-bottom)] fixed bottom-0 left-0 right-0 md:static"
      >
        <app-timeline class="block w-full"></app-timeline>
      </div>
    </div>
  `,
  styles: `
    :host {
      display: block;
      width: 100%;
      height: 100%;
      min-height: 0;
    }
  `,
})
export class MainComponent {}

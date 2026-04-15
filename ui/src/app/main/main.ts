import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { TimelineComponent } from './timeline.component';

@Component({
  selector: 'app-main',
  imports: [RouterOutlet, TimelineComponent],
  template: `
    <div class="w-full h-full min-h-0 flex flex-col overflow-hidden">
      <div class="w-full flex-1 min-h-0 overflow-hidden flex items-stretch justify-center">
        <router-outlet></router-outlet>
      </div>
      <div class="w-full mt-auto shrink-0 sticky bottom-0 z-20 bg-background">
        <app-timeline class="block w-full"></app-timeline>
      </div>
    </div>
  `,
})
export class MainComponent {}

import {
  Component,
  DestroyRef,
  ElementRef,
  afterNextRender,
  computed,
  inject,
  input,
  signal,
} from '@angular/core';
import { Chart, ChartOptions } from '@tanstack/charts/angular';

@Component({
  selector: 'ctx-chart-frame',
  imports: [Chart],
  host: { class: 'dashboard-chart relative block min-h-24 min-w-0 flex-1' },
  template: `
    @if (options(); as chartOptions) {
      <tanstack-chart class="absolute inset-0" [options]="chartOptions" />
    }
    <ng-content />
  `,
})
export class ChartFrameComponent {
  private readonly element = inject<ElementRef<HTMLElement>>(ElementRef);
  private readonly destroyRef = inject(DestroyRef);
  private readonly height = signal(0);
  readonly chart = input.required<ChartOptions>();
  readonly options = computed(() => {
    const height = this.height();
    return height > 0 ? { ...this.chart(), height } : null;
  });

  constructor() {
    afterNextRender(() => {
      const observer = new ResizeObserver((entries) =>
        this.height.set(Math.floor(entries[0]?.contentRect.height ?? 0)),
      );
      observer.observe(this.element.nativeElement);
      this.destroyRef.onDestroy(() => observer.disconnect());
    });
  }
}

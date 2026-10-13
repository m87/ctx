import { Component, computed, inject, input, signal } from '@angular/core';
import { Router, RouterLink } from '@angular/router';
import { ChartBounds } from '@tanstack/charts';
import { ChartOptions } from '@tanstack/charts/angular';
import { DashboardMetricResult } from '../../../../api/dashboard/dashboard.service';
import { BarWidgetDefinition } from '../dashboard-definition';
import {
  BAR_COLUMN_WIDTH,
  BAR_LABEL_HEIGHT,
  BAR_LABEL_WIDTH,
  BAR_ROW_HEIGHT,
  BarChartRow,
  barChartDefinition,
  barChartRows,
} from './bar-chart-definition';
import { ChartFrameComponent } from './chart-frame.component';
import { formatResultValue } from './format-result-value';

type CategoriesResult = Extract<DashboardMetricResult, { kind: 'categories' }>;

interface LabelBox {
  row: BarChartRow;
  left: number;
  top: number;
  width: number;
  height: number;
}

const LABEL_GAP = 8;
const AXIS_ALLOWANCE = 32;

@Component({
  selector: 'ctx-bar-widget',
  imports: [ChartFrameComponent, RouterLink],
  host: { class: 'flex min-h-0 flex-1 flex-col gap-3' },
  template: `
    @if (rows().length) {
      <div class="flex min-h-0 flex-1 flex-col overflow-auto">
        <ctx-chart-frame
          [chart]="options()"
          [style.min-height.px]="minimumHeight()"
          [style.min-width.px]="minimumWidth()"
        >
          @for (box of labelBoxes(); track box.row.key) {
            @if (box.row.entity; as entity) {
              <a
                class="absolute flex items-center truncate rounded-sm text-xs focus-visible:outline focus-visible:outline-2 focus-visible:outline-ring"
                [class.justify-center]="vertical()"
                [style.left.px]="box.left"
                [style.top.px]="box.top"
                [style.width.px]="box.width"
                [style.height.px]="box.height"
                [title]="box.row.label"
                [routerLink]="[entity.type === 'project' ? '/project' : '/context', entity.id]"
                ><span class="truncate">{{ box.row.label }}</span></a
              >
            } @else {
              <span
                class="absolute flex items-center truncate text-xs"
                [class.justify-center]="vertical()"
                [style.left.px]="box.left"
                [style.top.px]="box.top"
                [style.width.px]="box.width"
                [style.height.px]="box.height"
                [title]="box.row.label"
                ><span class="truncate">{{ box.row.label }}</span></span
              >
            }
          }
        </ctx-chart-frame>
      </div>
    } @else {
      <p role="status" class="text-sm text-muted-foreground">No activity in the selected period.</p>
    }
    @if (widget().properties.showSummary) {
      <p class="text-xs text-muted-foreground">
        {{ result().categories.totalCategories }} categories ·
        {{ format(result().categories.totalValue) }}
        total
        @if (result().categories.omittedCategories) {
          · {{ result().categories.omittedCategories }} omitted
        }
      </p>
    }
  `,
})
export class BarWidgetComponent {
  private readonly router = inject(Router);
  readonly widget = input.required<BarWidgetDefinition>();
  readonly result = input.required<CategoriesResult>();
  readonly editing = input(false);
  readonly area = signal<ChartBounds | null>(null);
  readonly vertical = computed(() => this.widget().properties.orientation === 'vertical');
  readonly rows = computed(() => barChartRows(this.result().categories.categories));
  readonly minimumHeight = computed(() =>
    this.vertical() ? 0 : this.rows().length * BAR_ROW_HEIGHT + AXIS_ALLOWANCE,
  );
  readonly minimumWidth = computed(() =>
    this.vertical() ? this.rows().length * BAR_COLUMN_WIDTH + BAR_LABEL_WIDTH : 0,
  );
  readonly options = computed<ChartOptions>(() => {
    const rows = this.rows();
    const properties = this.widget().properties;
    return {
      definition: barChartDefinition(rows, properties, (value) => this.format(value)),
      ariaLabel: `${properties.orientation} bar chart with ${rows.length} categories`,
      tabIndex: this.editing() ? -1 : 0,
      onRender: ({ scene }) => this.area.set(scene.chart),
      onSelect: (point) => {
        const entity = (point?.datum as BarChartRow | undefined)?.entity;
        if (entity) {
          void this.router.navigate([
            entity.type === 'project' ? '/project' : '/context',
            entity.id,
          ]);
        }
      },
    };
  });
  readonly labelBoxes = computed<LabelBox[]>(() => {
    const area = this.area();
    const rows = this.rows();
    if (!area || !rows.length) return [];
    if (this.vertical()) {
      const step = area.width / rows.length;
      return rows.map((row, index) => ({
        row,
        left: area.x + index * step,
        top: area.y + area.height + LABEL_GAP,
        width: step,
        height: BAR_LABEL_HEIGHT - LABEL_GAP,
      }));
    }
    const step = area.height / rows.length;
    return rows.map((row, index) => ({
      row,
      left: 0,
      top: area.y + index * step,
      width: BAR_LABEL_WIDTH - LABEL_GAP,
      height: step,
    }));
  });

  format(value: number | null): string {
    return formatResultValue(value, this.result().unit, this.widget().properties);
  }
}

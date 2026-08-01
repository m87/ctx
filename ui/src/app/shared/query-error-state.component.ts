import { Component, computed, input, output } from '@angular/core';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideRefreshCcw, lucideServerOff, lucideTriangleAlert } from '@ng-icons/lucide';
import { HlmButtonImports } from '@spartan-ng/helm/button';
import { HlmEmptyImports } from '@spartan-ng/helm/empty';
import { isConnectionError, normalizeError } from '../../api/error';

export interface QueryErrorContent {
  icon: 'lucideServerOff' | 'lucideTriangleAlert';
  title: string;
  description: string;
  retryable: boolean;
}

export function getQueryErrorContent(
  error: unknown,
  resourceName: string,
  paused = false,
): QueryErrorContent {
  const normalizedError = normalizeError(error);
  const normalizedResourceName = resourceName.trim() || 'data';

  if (paused || isConnectionError(error)) {
    return {
      icon: 'lucideServerOff',
      title: "Can't reach the server",
      description: `We couldn't load ${normalizedResourceName}. Check that ctx is running and try again.`,
      retryable: true,
    };
  }

  if (normalizedError.status === 404) {
    return {
      icon: 'lucideTriangleAlert',
      title: `${capitalize(normalizedResourceName)} not found`,
      description: `The requested ${normalizedResourceName} does not exist or may have been removed.`,
      retryable: false,
    };
  }

  return {
    icon: 'lucideTriangleAlert',
    title: `Couldn't load ${normalizedResourceName}`,
    description: normalizedError.message,
    retryable: true,
  };
}

@Component({
  selector: 'ctx-query-error-state',
  imports: [NgIcon, HlmButtonImports, HlmEmptyImports],
  providers: [provideIcons({ lucideRefreshCcw, lucideServerOff, lucideTriangleAlert })],
  template: `
    <div
      hlmEmpty
      class="h-full w-full min-h-56 border border-dashed bg-gradient-to-b from-muted/40 to-background"
      role="alert"
      aria-live="polite"
    >
      <div hlmEmptyHeader>
        <div hlmEmptyMedia variant="icon">
          <ng-icon [name]="content().icon"></ng-icon>
        </div>
        <div hlmEmptyTitle>{{ content().title }}</div>
        <div hlmEmptyDescription>{{ content().description }}</div>
      </div>

      @if (content().retryable) {
        <div hlmEmptyContent>
          <button
            hlmBtn
            type="button"
            variant="outline"
            [disabled]="retrying()"
            (click)="retry.emit()"
          >
            <ng-icon name="lucideRefreshCcw"></ng-icon>
            {{ retrying() ? 'Trying again...' : 'Try again' }}
          </button>
        </div>
      }
    </div>
  `,
  styles: `
    :host {
      display: flex;
      width: 100%;
      flex: 1 1 auto;
      min-height: 0;
    }
  `,
})
export class QueryErrorStateComponent {
  readonly error = input<unknown>(null);
  readonly paused = input(false);
  readonly resourceName = input('data');
  readonly retrying = input(false);
  readonly retry = output<void>();

  readonly content = computed(() =>
    getQueryErrorContent(this.error(), this.resourceName(), this.paused()),
  );
}

function capitalize(value: string): string {
  return value.charAt(0).toUpperCase() + value.slice(1);
}

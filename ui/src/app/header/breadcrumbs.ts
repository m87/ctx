import { Injectable, signal } from '@angular/core';

export interface Breadcrumb {
  label: string;
  id: string;
  icon?: string;
}

@Injectable({ providedIn: 'root' })
export class BreadcrumbService {
  private readonly _breadcrumbs = signal<Breadcrumb[]>([]);

  readonly breadcrumbs = this._breadcrumbs.asReadonly();

  setBreadcrumbs(breadcrumbs: Breadcrumb[]) {
    this._breadcrumbs.set(breadcrumbs);
  }

  appendBreadcrumb(breadcrumb: Breadcrumb) {
    this._breadcrumbs.update((current) => {
      const idx = current.findIndex((b) => b.id === breadcrumb.id);
      if (idx >= 0) {
        return current.slice(0, idx + 1);
      }
      return [...current, breadcrumb];
    });
  }

  resetBreadcrumbsTo(id: string) {
    this._breadcrumbs.update((current) => {
      const idx = current.findIndex((b) => b.id === id);
      return idx === -1 ? current : current.slice(0, idx + 1);
    });
  }
}

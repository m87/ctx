import { Injectable, signal } from '@angular/core';

@Injectable({ providedIn: 'root' })
export class SidebarStore {
  private readonly _collapsed = signal<boolean>(false);
  private readonly _mobileOpen = signal<boolean>(false);

  readonly collapsed = this._collapsed.asReadonly();
  readonly mobileOpen = this._mobileOpen.asReadonly();

  toggle(): void {
    this._collapsed.update((collapsed) => !collapsed);
  }

  collapse(): void {
    this._collapsed.set(true);
  }

  expand(): void {
    this._collapsed.set(false);
  }

  toggleMobile(): void {
    this._mobileOpen.update((open) => !open);
  }

  openMobile(): void {
    this._mobileOpen.set(true);
  }

  closeMobile(): void {
    this._mobileOpen.set(false);
  }
}

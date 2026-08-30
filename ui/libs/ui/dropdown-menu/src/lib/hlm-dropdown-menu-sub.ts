import { CdkMenu } from '@angular/cdk/menu';
import { Directive, inject, signal } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { classes } from '@spartan-ng/helm/utils';

@Directive({
	selector: '[hlmDropdownMenuSub],hlm-dropdown-menu-sub',
	hostDirectives: [CdkMenu],
	host: {
		'data-slot': 'dropdown-menu-sub',
		'[attr.data-state]': '_state()',
		'[attr.data-side]': '_side()',
	},
})
export class HlmDropdownMenuSub {
	private readonly _host = inject(CdkMenu);

	protected readonly _state = signal('open');
	protected readonly _side = signal('top');

	constructor() {
		this.setSideWithDarkMagic();
		this._host.closed.pipe(takeUntilDestroyed()).subscribe(() => this._state.set('closed'));

		classes(
			() =>
				'bg-popover text-popover-foreground data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 z-50 min-w-[8rem] origin-top overflow-hidden rounded-md border p-1 shadow-lg',
		);
	}

	private setSideWithDarkMagic() {
		const isRoot = this._host.menuStack.peek() === undefined;
		setTimeout(() => {
			const ps = (this._host as any)._parentTrigger._spartanLastPosition;
			if (!ps) {
				this._side.set(isRoot ? 'top' : 'left');
				return;
			}
			const side = isRoot ? ps.originY : ps.originX === 'end' ? 'right' : 'left';
			this._side.set(side);
		});
	}
}

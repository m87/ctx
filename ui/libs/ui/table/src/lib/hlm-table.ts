import { computed, Directive, inject, InjectionToken, input, type ValueProvider } from '@angular/core';
import { classes } from '@spartan-ng/helm/utils';

export const HlmTableConfigToken = new InjectionToken<HlmTableVariant>('HlmTableConfig');
export interface HlmTableVariant {
	tableContainer: string;
	table: string;
	thead: string;
	tbody: string;
	tfoot: string;
	tr: string;
	th: string;
	td: string;
	caption: string;
}

export const HlmTableVariantDefault: HlmTableVariant = {
	tableContainer: 'relative w-full overflow-x-auto',
	table: 'w-full caption-bottom text-sm',
	thead: '[&_tr]:border-b',
	tbody: '[&_tr:last-child]:border-0',
	tfoot: 'bg-muted/50 border-t font-medium [&>tr]:last:border-b-0',
	tr: 'hover:bg-muted/50 data-[state=selected]:bg-muted border-b transition-colors',
	th: 'text-foreground h-10 px-2 text-start align-middle font-medium whitespace-nowrap [&:has([role=checkbox])]:pe-0',
	td: 'p-2 align-middle whitespace-nowrap [&:has([role=checkbox])]:pe-0',
	caption: 'text-muted-foreground mt-4 text-sm',
};

export function provideHlmTableConfig(config: Partial<HlmTableVariant>): ValueProvider {
	return {
		provide: HlmTableConfigToken,
		useValue: { ...HlmTableVariantDefault, ...config },
	};
}

export function injectHlmTableConfig(): HlmTableVariant {
	return inject(HlmTableConfigToken, { optional: true }) ?? HlmTableVariantDefault;
}

@Directive({
	selector: 'div[hlmTableContainer]',
	host: { 'data-slot': 'table-container' },
})
export class HlmTableContainer {
	private readonly _globalOrDefaultConfig = injectHlmTableConfig();

	constructor() {
		classes(() => (this._globalOrDefaultConfig ? this._globalOrDefaultConfig.tableContainer.trim() : ''));
	}
}

@Directive({
	selector: 'table[hlmTable]',
	host: { 'data-slot': 'table' },
})
export class HlmTable {
	public readonly userVariant = input<Partial<HlmTableVariant> | string>({}, { alias: 'hlmTable' });

	private readonly _globalOrDefaultConfig = injectHlmTableConfig();

	protected readonly _variant = computed<HlmTableVariant>(() => {
		const globalOrDefaultConfig = this._globalOrDefaultConfig;
		const localInputConfig = this.userVariant();

		if (typeof localInputConfig === 'object' && localInputConfig !== null && Object.keys(localInputConfig).length > 0) {
			return { ...globalOrDefaultConfig, ...localInputConfig };
		}
		return globalOrDefaultConfig;
	});

	constructor() {
		classes(() => this._variant().table);
	}
}


@Directive({
	selector: 'thead[hlmTHead]',
	host: { 'data-slot': 'table-header' },
})
export class HlmTHead {
	private readonly _globalOrDefaultConfig = injectHlmTableConfig();

	constructor() {
		classes(() => (this._globalOrDefaultConfig ? this._globalOrDefaultConfig.thead.trim() : ''));
	}
}

@Directive({
	selector: 'tbody[hlmTBody]',
	host: { 'data-slot': 'table-body' },
})
export class HlmTBody {
	private readonly _globalOrDefaultConfig = injectHlmTableConfig();
	constructor() {
		classes(() => (this._globalOrDefaultConfig ? this._globalOrDefaultConfig.tbody.trim() : ''));
	}
}

@Directive({
	selector: 'tfoot[hlmTFoot]',
	host: { 'data-slot': 'table-footer' },
})
export class HlmTFoot {
	private readonly _globalOrDefaultConfig = injectHlmTableConfig();
	constructor() {
		classes(() => (this._globalOrDefaultConfig ? this._globalOrDefaultConfig.tfoot.trim() : ''));
	}
}

@Directive({
	selector: 'tr[hlmTr]',
	host: { 'data-slot': 'table-row' },
})
export class HlmTr {
	private readonly _globalOrDefaultConfig = injectHlmTableConfig();
	constructor() {
		classes(() => (this._globalOrDefaultConfig ? this._globalOrDefaultConfig.tr.trim() : ''));
	}
}

@Directive({
	selector: 'th[hlmTh]',
	host: { 'data-slot': 'table-head' },
})
export class HlmTh {
	private readonly _globalOrDefaultConfig = injectHlmTableConfig();
	constructor() {
		classes(() => (this._globalOrDefaultConfig ? this._globalOrDefaultConfig.th.trim() : ''));
	}
}

@Directive({
	selector: 'td[hlmTd]',
	host: { 'data-slot': 'table-cell' },
})
export class HlmTd {
	private readonly _globalOrDefaultConfig = injectHlmTableConfig();
	constructor() {
		classes(() => (this._globalOrDefaultConfig ? this._globalOrDefaultConfig.td.trim() : ''));
	}
}

@Directive({
	selector: 'caption[hlmCaption]',
	host: { 'data-slot': 'table-caption' },
})
export class HlmCaption {
	private readonly _globalOrDefaultConfig = injectHlmTableConfig();
	constructor() {
		classes(() => (this._globalOrDefaultConfig ? this._globalOrDefaultConfig.caption.trim() : ''));
	}
}

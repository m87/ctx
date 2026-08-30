import { Directive } from '@angular/core';
import { classes } from '@spartan-ng/helm/utils';

@Directive({
	selector: 'hlm-hint',
})
export class HlmHint {
	constructor() {
		classes(() => 'text-muted-foreground block text-sm');
	}
}

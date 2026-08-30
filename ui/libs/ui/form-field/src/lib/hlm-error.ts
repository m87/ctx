import { Directive } from '@angular/core';
import { classes } from '@spartan-ng/helm/utils';

@Directive({
	selector: 'hlm-error',
})
export class HlmError {
	constructor() {
		classes(() => 'text-destructive block text-sm font-medium');
	}
}

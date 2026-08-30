import type { BooleanInput } from '@angular/cdk/coercion';
import { booleanAttribute, ChangeDetectionStrategy, Component, computed, input } from '@angular/core';
import type { RouterLink } from '@angular/router';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideChevronLeft } from '@ng-icons/lucide';
import type { ButtonVariants } from '@spartan-ng/helm/button';
import { HlmIcon } from '@spartan-ng/helm/icon';
import { classes } from '@spartan-ng/helm/utils';
import { HlmPaginationLink } from './hlm-pagination-link';

@Component({
	selector: 'hlm-pagination-previous',
	imports: [HlmPaginationLink, NgIcon, HlmIcon],
	providers: [provideIcons({ lucideChevronLeft })],
	changeDetection: ChangeDetectionStrategy.OnPush,
	template: `
		<a
			hlmPaginationLink
			[link]="link()"
			[queryParams]="queryParams()"
			[queryParamsHandling]="queryParamsHandling()"
			[size]="_size()"
			[attr.aria-label]="ariaLabel()"
		>
			<ng-icon hlm size="sm" name="lucideChevronLeft" />
			<span [class]="_labelClass()">{{ text() }}</span>
		</a>
	`,
})
export class HlmPaginationPrevious {
	public readonly link = input<RouterLink['routerLink']>();
	public readonly queryParams = input<RouterLink['queryParams']>();
	public readonly queryParamsHandling = input<RouterLink['queryParamsHandling']>();

	public readonly ariaLabel = input<string>('Go to previous page', { alias: 'aria-label' });
	public readonly text = input<string>('Previous');
	public readonly iconOnly = input<boolean, BooleanInput>(false, {
		transform: booleanAttribute,
	});
	protected readonly _labelClass = computed(() => (this.iconOnly() ? 'sr-only' : 'hidden sm:block'));
	protected readonly _size = computed<ButtonVariants['size']>(() => (this.iconOnly() ? 'icon' : 'default'));

	constructor() {
		classes(() => ['gap-1 px-2.5', !this.iconOnly() ? 'sm:pl-2.5' : '']);
	}
}

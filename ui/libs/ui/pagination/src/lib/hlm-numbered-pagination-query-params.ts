import type { BooleanInput, NumberInput } from '@angular/cdk/coercion';
import {
	ChangeDetectionStrategy,
	Component,
	booleanAttribute,
	computed,
	input,
	model,
	numberAttribute,
	untracked,
} from '@angular/core';
import { FormsModule } from '@angular/forms';
import { BrnSelectImports } from '@spartan-ng/brain/select';
import { HlmSelectImports } from '@spartan-ng/helm/select';
import { createPageArray, outOfBoundCorrection } from './hlm-numbered-pagination';
import { HlmPagination } from './hlm-pagination';
import { HlmPaginationContent } from './hlm-pagination-content';
import { HlmPaginationEllipsis } from './hlm-pagination-ellipsis';
import { HlmPaginationItem } from './hlm-pagination-item';
import { HlmPaginationLink } from './hlm-pagination-link';
import { HlmPaginationNext } from './hlm-pagination-next';
import { HlmPaginationPrevious } from './hlm-pagination-previous';

@Component({
	selector: 'hlm-numbered-pagination-query-params',
	imports: [
		FormsModule,
		HlmPagination,
		HlmPaginationContent,
		HlmPaginationItem,
		HlmPaginationPrevious,
		HlmPaginationNext,
		HlmPaginationLink,
		HlmPaginationEllipsis,

		BrnSelectImports,
		HlmSelectImports,
	],
	changeDetection: ChangeDetectionStrategy.OnPush,
	template: `
		<div class="flex items-center justify-between gap-2 px-4 py-2">
			<div class="flex items-center gap-1 text-sm text-nowrap text-gray-600">
				<b>{{ totalItems() }}</b>
				total items |
				<b>{{ _lastPageNumber() }}</b>
				pages
			</div>

			<nav hlmPagination>
				<ul hlmPaginationContent>
					@if (showEdges() && !_isFirstPageActive()) {
						<li hlmPaginationItem>
							<hlm-pagination-previous
								[link]="link()"
								[queryParams]="{ page: currentPage() - 1 }"
								queryParamsHandling="merge"
							/>
						</li>
					}

					@for (page of _pages(); track page) {
						<li hlmPaginationItem>
							@if (page === '...') {
								<hlm-pagination-ellipsis />
							} @else {
								<a
									hlmPaginationLink
									[link]="currentPage() !== page ? link() : undefined"
									[queryParams]="{ page }"
									queryParamsHandling="merge"
									[isActive]="currentPage() === page"
								>
									{{ page }}
								</a>
							}
						</li>
					}

					@if (showEdges() && !_isLastPageActive()) {
						<li hlmPaginationItem>
							<hlm-pagination-next
								[link]="link()"
								[queryParams]="{ page: currentPage() + 1 }"
								queryParamsHandling="merge"
							/>
						</li>
					}
				</ul>
			</nav>

			<brn-select [(ngModel)]="itemsPerPage" class="ml-auto" placeholder="Page size">
				<hlm-select-trigger class="w-fit">
					<hlm-select-value />
				</hlm-select-trigger>
				<hlm-select-content>
					@for (pageSize of _pageSizesWithCurrent(); track pageSize) {
						<hlm-option [value]="pageSize">{{ pageSize }} / page</hlm-option>
					}
				</hlm-select-content>
			</brn-select>
		</div>
	`,
})
export class HlmNumberedPaginationQueryParams {
	public readonly currentPage = model.required<number>();

	public readonly itemsPerPage = model.required<number>();

	public readonly totalItems = input.required<number, NumberInput>({
		transform: numberAttribute,
	});

	public readonly link = input<string>('.');

	public readonly maxSize = input<number, NumberInput>(7, {
		transform: numberAttribute,
	});

	public readonly showEdges = input<boolean, BooleanInput>(true, {
		transform: booleanAttribute,
	});

	public readonly pageSizes = input<number[]>([10, 20, 50, 100]);

	protected readonly _pageSizesWithCurrent = computed(() => {
		const pageSizes = this.pageSizes();
		return pageSizes.includes(this.itemsPerPage())
			? pageSizes
			: [...pageSizes, this.itemsPerPage()].sort((a, b) => a - b);
	});

	protected readonly _isFirstPageActive = computed(() => this.currentPage() === 1);
	protected readonly _isLastPageActive = computed(() => this.currentPage() === this._lastPageNumber());

	protected readonly _lastPageNumber = computed(() => {
		if (this.totalItems() < 1) {
			return 1;
		}
		return Math.ceil(this.totalItems() / this.itemsPerPage());
	});

	protected readonly _pages = computed(() => {
		const correctedCurrentPage = outOfBoundCorrection(this.totalItems(), this.itemsPerPage(), this.currentPage());

		if (correctedCurrentPage !== this.currentPage()) {
			untracked(() => this.currentPage.set(correctedCurrentPage));
		}

		return createPageArray(correctedCurrentPage, this.itemsPerPage(), this.totalItems(), this.maxSize());
	});
}

<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	export let page: number;
	export let maxPage: number;
	export let resultsCount: number;
	export let showPageInput: boolean = false;
	export let zeroBased: boolean = false;

	const dispatch = createEventDispatcher();

	// For zero-based pagination: page 0 is first, maxPage is last valid index
	// For one-based pagination: page 1 is first, maxPage is total pages
	$: minPage = zeroBased ? 0 : 1;
	$: displayPage = zeroBased ? page + 1 : page;
	$: displayMaxPage = zeroBased ? maxPage + 1 : maxPage;
	$: isFirstPage = page <= minPage;
	$: isLastPage = page >= maxPage;

	function prevPage() {
		if (!isFirstPage) {
			dispatch('prevPage');
		}
	}

	function nextPage() {
		if (!isLastPage) {
			dispatch('nextPage');
		}
	}

	function handlePageInput() {
		dispatch('pageChange');
	}
</script>

<div class="pagination-footer">
	<div class="pagination-results">
		<span class="font-semibold">{resultsCount}</span> résultats
	</div>
	<div class="pagination-controls">
		<p class="pagination-info">
			Page
			{#if showPageInput}
				<input
					type="number"
					min="1"
					max={displayMaxPage}
					class="page-input"
					bind:value={page}
					on:change={handlePageInput}
				/>
			{:else}
				<span class="font-bold">{displayPage}</span>
			{/if}
			sur <span class="font-bold">{displayMaxPage}</span>
		</p>
		<div class="pagination-buttons">
			<button class="pagination-button" on:click={prevPage} disabled={isFirstPage}>
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M15 19l-7-7 7-7"
					/>
				</svg>
				<span>Précédent</span>
			</button>
			<button class="pagination-button" on:click={nextPage} disabled={isLastPage}>
				<span>Suivant</span>
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M9 5l7 7-7 7"
					/>
				</svg>
			</button>
		</div>
	</div>
</div>

<style>
	.pagination-footer {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 16px;
		background-color: white;
		border-top: 1px solid #e5e7eb;
	}

	.pagination-results {
		font-size: 14px;
		color: #6b7280;
	}

	.pagination-controls {
		display: flex;
		align-items: center;
		gap: 16px;
	}

	.pagination-info {
		font-size: 14px;
		font-weight: 500;
		color: #374151;
	}

	.page-input {
		width: 50px;
		padding: 4px 6px;
		border: 1px solid #d1d5db;
		border-radius: 4px;
		font-size: 14px;
		text-align: center;
		font-weight: 600;
		appearance: textfield;
		-moz-appearance: textfield;
	}

	.page-input::-webkit-outer-spin-button,
	.page-input::-webkit-inner-spin-button {
		-webkit-appearance: none;
		margin: 0;
	}

	.page-input:focus {
		outline: none;
		border-color: #3b82f6;
	}

	.pagination-buttons {
		display: flex;
		gap: 8px;
	}

	.pagination-button {
		display: flex;
		align-items: center;
		gap: 8px;
		padding: 8px 16px;
		background-color: white;
		color: #374151;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		font-size: 14px;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.pagination-button:hover:not(:disabled) {
		background-color: #f9fafb;
		border-color: #9ca3af;
	}

	.pagination-button:disabled {
		background-color: #f9fafb;
		color: #9ca3af;
		cursor: not-allowed;
	}

	@media (prefers-color-scheme: dark) {
		.pagination-footer {
			background-color: #1f2937;
			border-color: #374151;
		}

		.pagination-results {
			color: #9ca3af;
		}

		.pagination-info {
			color: #e5e7eb;
		}

		.page-input {
			background-color: #374151;
			color: #e5e7eb;
			border-color: #4b5563;
		}

		.pagination-button {
			background-color: #374151;
			color: #e5e7eb;
			border-color: #4b5563;
		}

		.pagination-button:hover:not(:disabled) {
			background-color: #4b5563;
			border-color: #6b7280;
		}

		.pagination-button:disabled {
			background-color: #374151;
			color: #6b7280;
		}
	}

	/* Mobile: compact two-line layout */
	@media (max-width: 480px) {
		.pagination-footer {
			flex-wrap: wrap;
			padding: 10px 12px;
			gap: 8px;
		}

		.pagination-results {
			font-size: 12px;
		}

		.pagination-info {
			font-size: 12px;
		}

		.pagination-controls {
			width: 100%;
			justify-content: space-between;
			gap: 8px;
		}

		.pagination-buttons {
			gap: 6px;
		}

		.pagination-button {
			padding: 6px 10px;
			font-size: 12px;
			gap: 4px;
		}

		.pagination-button svg {
			width: 12px;
			height: 12px;
		}

		.page-input {
			width: 40px;
			padding: 2px 4px;
			font-size: 12px;
		}
	}
</style>

<script lang="ts">
	import { goto } from '$app/navigation';
	import { searchName } from '$lib/store/store';

	export let activeTab: 'transactions' | 'resume' | 'refills' = 'transactions';
	export let showSearch: boolean = true;

	let searchNameValue: string;

	searchName.subscribe((value) => {
		searchNameValue = value;
	});

	function handleSearchInput(event: Event) {
		const target = event.currentTarget as HTMLInputElement;
		searchNameValue = target.value.toLowerCase();
	}
</script>

<div class="header-controls">
	<div class="title-section">
		<div class="view-tabs">
			<button
				class="tab"
				class:active={activeTab === 'transactions'}
				on:click={() => goto('/comptoir/c/transactions')}
			>
				<iconify-icon icon="mdi:format-list-bulleted" width="18" height="18" />
				Transactions
			</button>
			<button
				class="tab"
				class:active={activeTab === 'resume'}
				on:click={() => goto('/comptoir/c/resume')}
			>
				<iconify-icon icon="mdi:chart-box" width="18" height="18" />
				Résumé
			</button>
			<button
				class="tab"
				class:active={activeTab === 'refills'}
				on:click={() => goto('/comptoir/c/refills')}
			>
				<iconify-icon icon="mdi:cash-clock" width="18" height="18" />
				Historique recharges
			</button>
		</div>
		{#if showSearch}
			<input
				class="search-input"
				placeholder="Rechercher une personne"
				bind:value={searchNameValue}
				on:input={handleSearchInput}
			/>
		{/if}
	</div>
	<slot name="filters" />
</div>

<style>
	.header-controls {
		display: flex;
		flex-direction: row;
		align-items: center;
		gap: 20px;
		padding: 1rem;
		background-color: var(--bg-secondary);
		border-radius: 0.75rem;
		border: 1px solid var(--border-color);
		margin-bottom: 10px;
	}

	.title-section {
		display: flex;
		flex-direction: row;
		align-items: center;
		gap: 40px;
		flex-grow: 1;
	}

	.view-tabs {
		display: flex;
		background-color: #e5e7eb;
		border-radius: 8px;
		padding: 4px;
		gap: 0;
	}

	.tab {
		display: flex;
		align-items: center;
		gap: 6px;
		padding: 8px 16px;
		background-color: transparent;
		color: #374151;
		border: none;
		cursor: pointer;
		font-size: 14px;
		font-weight: 500;
		transition: all 0.2s ease;
		border-radius: 6px;
	}

	.tab:hover {
		background-color: rgba(0, 0, 0, 0.05);
	}

	.tab.active {
		background-color: #2563eb;
		color: white;
	}

	@media (prefers-color-scheme: dark) {
		.view-tabs {
			background-color: #374151;
		}

		.tab {
			color: #e5e7eb;
		}

		.tab:hover {
			background-color: rgba(255, 255, 255, 0.05);
		}

		.tab.active {
			background-color: #2563eb;
			color: white;
		}
	}

	.search-input {
		border-radius: 8px;
		padding: 8px;
		background-color: white;
		color: black;
		border: 1px solid #ccc;
	}

	@media (prefers-color-scheme: dark) {
		.search-input {
			background-color: #374151;
			color: white;
			border: 1px solid #6b7280;
		}
	}

	:global(.filters-section) {
		display: flex;
		flex-direction: row;
		align-items: center;
		gap: 20px;
		font-size: 18px;
		font-weight: 600;
	}

	:global(.filter-select) {
		background-color: #f3f4f6;
		padding: 8px;
		border-radius: 8px;
		color: #111827;
		border: 1px solid #d1d5db;
		cursor: pointer;
	}

	:global(.filter-select:hover) {
		background-color: #e5e7eb;
	}

	@media (prefers-color-scheme: dark) {
		:global(.filter-select) {
			background-color: #374151;
			color: white;
			border: none;
		}

		:global(.filter-select:hover) {
			background-color: #4b5563;
		}
	}

	:global(.checkbox-label) {
		display: flex;
		align-items: center;
		gap: 4px;
	}

	:global(.checkbox-input) {
		height: 24px;
		width: 24px;
		vertical-align: middle;
		margin-left: 4px;
	}

	@media (max-width: 768px) {
		.header-controls {
			flex-direction: column;
			gap: 12px;
		}

		.title-section {
			flex-direction: column;
			gap: 12px;
			align-items: stretch;
		}

		:global(.filters-section) {
			flex-direction: column;
			gap: 12px;
			align-items: stretch;
		}
	}

	@media (max-width: 480px) {
		.view-tabs button {
			font-size: 9px;
			padding: 3px 8px;
		}
	}
</style>
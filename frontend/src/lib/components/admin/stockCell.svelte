<script lang="ts">
	export let item: {
		id: string;
		amount_left: number;
		buy_limit?: number;
		optimal_amount: number;
		category_id: string;
	};
	export let editItem: (id: string, item: any, category_id: string) => void;

	let expanded = false;

	function handleStockChange(e: Event) {
		const target = e.target as HTMLInputElement;
		let stock = parseInt(target.value);
		editItem(item.id, { amount_left: stock }, item.category_id);
	}

	function handleLimitChange(e: Event) {
		const target = e.target as HTMLInputElement;
		if (target.value === '') {
			editItem(item.id, { buy_limit: -1 }, item.category_id);
			return;
		}
		let buy_limit = parseInt(target.value);
		editItem(item.id, { buy_limit: buy_limit }, item.category_id);
	}

	function handleOptimalChange(e: Event) {
		const target = e.target as HTMLInputElement;
		let optimal_amount = parseInt(target.value);
		editItem(item.id, { optimal_amount: optimal_amount }, item.category_id);
	}
</script>

<div class="relative">
	{#if !expanded}
		<div
			class="flex items-center gap-1 flex-wrap text-xs"
			role="button"
			tabindex="0"
			on:click={() => (expanded = true)}
			on:keydown={(e) => e.key === 'Enter' && (expanded = true)}
		>
			<span class="bg-gray-100 dark:bg-gray-700 px-1 rounded">Stock:{item.amount_left}</span>
			<span class="bg-gray-100 dark:bg-gray-700 px-1 rounded">Lim:{item.buy_limit ?? '-'}</span>
			<span class="bg-gray-100 dark:bg-gray-700 px-1 rounded">Opt:{item.optimal_amount}</span>
			<button
				class="ml-1 p-1 hover:bg-gray-200 dark:hover:bg-gray-600 rounded transition-colors"
				on:click|stopPropagation={() => (expanded = true)}
			>
				<iconify-icon icon="mdi:pencil" width="14" height="14"></iconify-icon>
			</button>
		</div>
	{:else}
		<div class="min-w-[140px] space-y-1 mt-8">
			<button
				class="absolute -top-8 right-2 px-2 py-1 bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 rounded z-10 flex items-center gap-1 text-xs"
				on:click={() => (expanded = false)}
			>
				<iconify-icon icon="mdi:close" width="14" height="14"></iconify-icon>
				<span>Fermer</span>
			</button>
			<div class="flex items-center gap-2">
				<span class="text-xs w-14 shrink-0 dark:text-gray-300">Stock:</span>
				<input
					type="number"
					value={item.amount_left}
					class="w-full text-xs px-2 py-1 border border-gray-200 dark:border-gray-600 rounded dark:bg-gray-700 dark:text-white focus:border-blue-500 focus:outline-none"
					on:input={handleStockChange}
				/>
			</div>
			<div class="flex items-center gap-2">
				<span class="text-xs w-14 shrink-0 dark:text-gray-300">Limite:</span>
				<input
					type="number"
					value={item.buy_limit}
					class="w-full text-xs px-2 py-1 border border-gray-200 dark:border-gray-600 rounded dark:bg-gray-700 dark:text-white focus:border-blue-500 focus:outline-none"
					on:input={handleLimitChange}
				/>
			</div>
			<div class="flex items-center gap-2">
				<span class="text-xs w-14 shrink-0 dark:text-gray-300">Optimal:</span>
				<input
					type="number"
					value={item.optimal_amount}
					class="w-full text-xs px-2 py-1 border border-gray-200 dark:border-gray-600 rounded dark:bg-gray-700 dark:text-white focus:border-blue-500 focus:outline-none"
					on:input={handleOptimalChange}
				/>
			</div>
		</div>
	{/if}
</div>

<style>
	input[type='number'] {
		-moz-appearance: textfield;
	}
	input[type='number']::-webkit-outer-spin-button,
	input[type='number']::-webkit-inner-spin-button {
		-webkit-appearance: none;
		margin: 0;
	}
</style>

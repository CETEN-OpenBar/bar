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
		<button
			class="w-full text-left p-2 rounded-md border border-gray-200 dark:border-gray-600 hover:border-blue-400 dark:hover:border-blue-500 hover:bg-blue-50 dark:hover:bg-blue-950/30 transition-all group"
			on:click={() => (expanded = true)}
		>
			<div class="flex items-center justify-between gap-2">
				<div class="flex flex-col gap-1 text-xs flex-1 min-w-0">
					<div class="flex items-center gap-2">
						<span class="text-gray-500 dark:text-gray-400 font-medium w-14">Stock:</span>
						<span class="font-semibold dark:text-gray-200">{item.amount_left}</span>
					</div>
					<div class="flex items-center gap-2">
						<span class="text-gray-500 dark:text-gray-400 font-medium w-14">Limite:</span>
						<span class="font-semibold dark:text-gray-200">{item.buy_limit ?? '-'}</span>
					</div>
					<div class="flex items-center gap-2">
						<span class="text-gray-500 dark:text-gray-400 font-medium w-14">Optimal:</span>
						<span class="font-semibold dark:text-gray-200">{item.optimal_amount}</span>
					</div>
				</div>
				<iconify-icon icon="mdi:chevron-down" width="16" height="16" class="text-gray-400 group-hover:text-blue-500 transition-colors flex-shrink-0"></iconify-icon>
			</div>
		</button>
	{:else}
		<div class="border border-blue-500 dark:border-blue-400 rounded-md p-3 bg-blue-50/50 dark:bg-blue-950/20">
			<div class="flex items-center justify-between mb-3">
				<span class="text-xs font-semibold text-blue-700 dark:text-blue-300">Modifier Stocks</span>
				<button
					class="px-2 py-1 text-xs bg-white dark:bg-gray-700 hover:bg-gray-100 dark:hover:bg-gray-600 border border-gray-300 dark:border-gray-500 rounded flex items-center gap-1 transition-colors"
					on:click={() => (expanded = false)}
				>
					<iconify-icon icon="mdi:check" width="14" height="14"></iconify-icon>
					<span>Fermer</span>
				</button>
			</div>
			<div class="space-y-2">
				<div class="flex items-center gap-2">
					<label class="text-xs w-16 shrink-0 font-medium dark:text-gray-300">Stock:</label>
					<input
						type="number"
						value={item.amount_left}
						class="flex-1 text-xs px-2 py-1.5 border border-gray-300 dark:border-gray-500 rounded bg-white dark:bg-gray-700 dark:text-white focus:border-blue-500 focus:ring-2 focus:ring-blue-200 dark:focus:ring-blue-800 outline-none transition-all"
						on:input={handleStockChange}
					/>
				</div>
				<div class="flex items-center gap-2">
					<label class="text-xs w-16 shrink-0 font-medium dark:text-gray-300">Limite:</label>
					<input
						type="number"
						value={item.buy_limit}
						class="flex-1 text-xs px-2 py-1.5 border border-gray-300 dark:border-gray-500 rounded bg-white dark:bg-gray-700 dark:text-white focus:border-blue-500 focus:ring-2 focus:ring-blue-200 dark:focus:ring-blue-800 outline-none transition-all"
						on:input={handleLimitChange}
					/>
				</div>
				<div class="flex items-center gap-2">
					<label class="text-xs w-16 shrink-0 font-medium dark:text-gray-300">Optimal:</label>
					<input
						type="number"
						value={item.optimal_amount}
						class="flex-1 text-xs px-2 py-1.5 border border-gray-300 dark:border-gray-500 rounded bg-white dark:bg-gray-700 dark:text-white focus:border-blue-500 focus:ring-2 focus:ring-blue-200 dark:focus:ring-blue-800 outline-none transition-all"
						on:input={handleOptimalChange}
					/>
				</div>
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

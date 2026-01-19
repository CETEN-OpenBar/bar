<script lang="ts">
	import type {Item, TransactionItem} from '$lib/api';
	import { api } from '$lib/config/config';
	import { transactionsApi, itemsApi } from '$lib/requests/requests';
	import { onMount } from 'svelte';

	let items: Array<TransactionItem> = [];
	let stockItems: Array<Item> = [];
	let name: string = "";
	
	$: filtered_items = name !== "" ?
		items.filter((item) => {
			return item.item_name.toLocaleLowerCase().includes(name.toLocaleLowerCase())
		}) : items
	
	onMount(() => {
		reloadTransactionsItems();
		reloadItems();
	});

	function reloadTransactionsItems() {
		transactionsApi()
			.getTransactionsItems(undefined, { withCredentials: true })
			.then((res) => {
				if (!(res.data instanceof Array)) return;
				items = res.data;
			});
	}

	function reloadItems() {
		itemsApi()
			.getAllItems(1, 1000, undefined, undefined, undefined, undefined, undefined, { withCredentials: true })
			.then((res) => {
				if (!(res.data.items instanceof Array)) return;
				stockItems = res.data.items;
			});
	}

	function getStockAmount(itemId: string): number {
		const stockItem = stockItems.find(i => i.id === itemId);
		return stockItem?.amount_left ?? 0;
	}
</script>

<div class="w-full p-4 md:p-0">
	<div class="hidden md:block bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 overflow-hidden">
		<div class="overflow-x-auto">
			<table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
				<thead class="bg-gray-50 dark:bg-gray-700">
					<tr class="grid grid-cols-5 divide-x divide-gray-200 dark:divide-gray-700">
						<th class="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider flex items-center gap-2 sticky left-0 bg-gray-50 dark:bg-gray-700 z-10">
							<span>Nom</span>
							<input
								type="text"
								class="ml-[10px] rounded-md p-1.5 text-sm text-black dark:text-white dark:bg-gray-600 placeholder-gray-400 w-24"
								placeholder="rechercher"
								bind:value={name}
							/>
						</th>
						<th class="px-4 py-3 text-center text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
							Commandés
						</th>
						<th class="px-4 py-3 text-center text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
							Terminés
						</th>
						<th class="px-4 py-3 text-center text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                            Restants
						</th>
                        <th class="px-4 py-3 text-center text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                            Stock
                        </th>
					</tr>
				</thead>
				<tbody class="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700 max-h-[calc(100vh-240px)] overflow-y-auto">
					{#each filtered_items as item}
						<tr class="grid grid-cols-5 divide-x divide-gray-200 dark:divide-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
							<td class="px-4 py-3 flex items-center gap-3 sticky left-0 bg-white dark:bg-gray-800 z-10">
								<img
									src={api() + item.picture_uri}
									alt="item"
									class="w-10 h-10 rounded-2xl flex-shrink-0"
								/>
								<span class="text-sm text-gray-900 dark:text-gray-300 truncate">{item.item_name}</span>
							</td>
							<td class="px-4 py-3 text-center text-sm text-gray-900 dark:text-gray-300">{item.item_amount}</td>
							<td class="px-4 py-3 text-center text-sm text-gray-900 dark:text-gray-300">{item.item_already_done}</td>
							<td class="px-4 py-3 text-center text-sm font-medium text-gray-900 dark:text-gray-300">
                                {item.item_amount - item.item_already_done}
							</td>
                            <td class="px-4 py-3 text-center text-sm font-medium text-gray-900 dark:text-gray-300">
                                {getStockAmount(item.item_id)}
                            </td>
						</tr>
					{/each}
					{#if filtered_items.length === 0}
						<tr>
							<td colspan="5" class="px-4 py-8 text-center text-sm text-gray-500 dark:text-gray-400">
								Aucun article trouvé
							</td>
						</tr>
					{/if}
				</tbody>
			</table>
		</div>
	</div>

	<div class="md:hidden space-y-4 max-h-[calc(100vh-240px)] overflow-y-auto">
		<div class="flex items-center gap-2 mb-2">
			<span class="text-sm font-medium text-gray-500 dark:text-gray-300">Nom</span>
			<input
				type="text"
				class="flex-1 rounded-md p-2 text-sm text-black dark:text-white dark:bg-gray-600 placeholder-gray-400"
				placeholder="rechercher"
				bind:value={name}
			/>
		</div>
		{#each filtered_items as item}
			<div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-4">
				<div class="flex items-center gap-3 mb-3">
					<img
						src={api() + item.picture_uri}
						alt="item"
						class="w-12 h-12 rounded-2xl"
					/>
					<span class="text-base font-medium text-gray-900 dark:text-gray-300">{item.item_name}</span>
				</div>
				<div class="grid grid-cols-4 gap-2 text-center">
					<div class="bg-gray-100 dark:bg-gray-700 rounded p-2">
						<div class="text-xs text-gray-500 dark:text-gray-400">Commandés</div>
						<div class="text-sm font-medium text-gray-900 dark:text-gray-300">{item.item_amount}</div>
					</div>
					<div class="bg-gray-100 dark:bg-gray-700 rounded p-2">
						<div class="text-xs text-gray-500 dark:text-gray-400">Terminés</div>
						<div class="text-sm font-medium text-gray-900 dark:text-gray-300">{item.item_already_done}</div>
					</div>
					<div class="bg-gray-100 dark:bg-gray-700 rounded p-2">
                        <div class="text-xs text-gray-500 dark:text-gray-400">Restants</div>
						<div class="text-sm font-medium text-gray-900 dark:text-gray-300">{item.item_amount - item.item_already_done}</div>
					</div>
                    <div class="bg-gray-100 dark:bg-gray-700 rounded p-2">
                        <div class="text-xs text-gray-500 dark:text-gray-400">Stock</div>
                        <div class="text-sm font-medium text-gray-900 dark:text-gray-300">{getStockAmount(item.item_id)}</div>
                    </div>
				</div>
			</div>
		{/each}
		{#if filtered_items.length === 0}
			<div class="text-center py-8 text-sm text-gray-500 dark:text-gray-400">
				Aucun article trouvé
			</div>
		{/if}
	</div>
</div>

<script lang="ts">
	import type {TransactionItem} from '$lib/api';
	import { api } from '$lib/config/config';
	import { transactionsApi } from '$lib/requests/requests';
	import { onMount } from 'svelte';

	let items: Array<TransactionItem> = [];
	let name: string = "";
	
	$: filtered_items = name !== "" ?
		items.filter((item) => {
			return item.item_name.toLocaleLowerCase().includes(name.toLocaleLowerCase())
		}) : items
	
	onMount(() => {
		reloadTransactionsItems();
	});

	function reloadTransactionsItems() {
		transactionsApi()
			.getTransactionsItems(undefined, { withCredentials: true })
			.then((res) => {
				if (!(res.data instanceof Array)) return;
				items = res.data;
			});
	}
</script>

<div class="w-full flex justify-center">
	<div class="overflow-x-auto overflow-y-visible w-full max-w-7xl">
		<div class="min-w-full bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700">
			<div class="grid grid-cols-[55%_15%_15%_15%] divide-x divide-gray-200 dark:divide-gray-700 bg-gray-50 dark:bg-gray-700 sticky top-0">
				<th class="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider flex items-center gap-2">
					<span>Nom</span>
					<input
						type="text"
                        class="ml-[10px] rounded-md p-1.5 text-sm text-black dark:text-white dark:bg-gray-600 placeholder-gray-400"
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
			</div>
			<div class="divide-y divide-gray-200 dark:divide-gray-700 max-h-[calc(100vh-240px)] overflow-y-auto">
				{#each filtered_items as item}
					<div class="grid grid-cols-[55%_15%_15%_15%] divide-x divide-gray-200 dark:divide-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
						<td class="px-4 py-3 flex items-center gap-3">
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
					</div>
				{/each}
				{#if filtered_items.length === 0}
					<div class="px-4 py-8 text-center text-sm text-gray-500 dark:text-gray-400">
						Aucun article trouvé
					</div>
				{/if}
			</div>
		</div>
	</div>
</div>

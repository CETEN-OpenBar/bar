<script lang="ts">
	import type {
		UpdateItem,
		Item,
		AccountPriceRole,
		ItemState,
		NewRestock,
		NewRestockItem,
		RestockItem,
		Restock,
		Category
	} from '$lib/api';
	import Items from '$lib/components/borne/items.svelte';
	import { api } from '$lib/config/config';
	import { categoriesApi, itemsApi, restocksApi } from '$lib/requests/requests';
	import { formatPrice, parsePrice } from '$lib/utils';
	import { onMount } from 'svelte';

	let sure: boolean = false;
	let items: Item[] = [];

	let restoks : Restock[] = [];
	let newRestock: NewRestock = {
		total_cost_ht: 0,
		total_cost_ttc: 0,
		driver_id: '',
		type: 'promocash',
		items: []
	};

	let page = 0;
	let max_page = 0;
	let itemsPerPage = 10;

	let itemName: string = '';
	let nameList: string[] = [];
	let newItem: NewRestockItem = {
		item_id: '',
		amount_of_bundle: 0,
		amount_per_bundle: 0,
		bundle_cost_ht: 0,
		tva: 0
	};
	let searchName: string = '';

	onMount(() => {
		reloadItems();
		restocksApi()
			.getRestocks(page, itemsPerPage, undefined, undefined, undefined)
	});

	function reloadItems() {
		itemsApi()
			.getAllItems(page, itemsPerPage, undefined, undefined, searchName, {
				withCredentials: true
			})
			.then((res) => {
				max_page = res.data.max_page ?? 0;
				page = res.data.page ?? 0;
				itemsPerPage = res.data.limit ?? 0;
				items = res.data.items ?? [];
			});
	}

	function updateTotalHTandTTC() {
		newRestock.total_cost_ht = 0.0;
		newRestock.items.forEach((item) => {
			newRestock.total_cost_ht +=
				item.amount_of_bundle * item.bundle_cost_ht;
		});
		newRestock.total_cost_ttc = 0.0;
		newRestock.items.forEach((item) => {
			newRestock.total_cost_ttc += Math.floor((item.amount_of_bundle * item.bundle_cost_ht * (1 + item.tva / 10000)));
		});

	}
    
function applyRestock() {
	if (!sure) return;
	restocksApi()
		.createRestock(newRestock, { withCredentials: true })
		.then((res) => {
			restoks = [... restoks, res.data]
		})

}

</script>

<div class="max-w-[95%] px-4 py-10 sm:px-6 lg:px-8 lg:py-14 mx-auto">
	<div class="py-3 px-2 w-1.0 flex m-auto">
		<select
			class="rounded-lg border-transparent appearance-none border border-gray-300 w-96 py-2 px-4 bg-white text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
			placeholder="Type"
			bind:value={newRestock.type}
		>
			<option value="promocash">Promocash</option>
			<option value="auchan_drive">Auchan drive</option>
			<option value="auchan">Auchan</option>
			<option value="viennoiserie">Boulangerie Benoist</option>
		</select>
		<input
			type="text"
			class="rounded-lg border-transparent ml-3 appearance-none border border-gray-300 w-100 py-2 px-4 bg-white text-gray-700 mr-auto placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
			placeholder="Conducteur"
			bind:value={newRestock.driver_id}
		/>
		<div>
			<p class="text-white text-2xl ml-5">Total HT : {newRestock.total_cost_ht / 100} €</p>
			<p class="text-white text-2xl ml-5">Total TTC : {newRestock.total_cost_ttc / 100} €</p>
		</div>
	</div>
	<div class="flex flex-col">
		<table class="mb-10 min-w-full divide-y divide-gray-200 dark:divide-gray-700 bg-blue-950">
			<thead class="bg-gray-50 dark:bg-blue-600">
				<tr>
					<th scope="col" class="px-12 py-3">
						<span
							class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
						>
							Nom
						</span>
					</th>
					<th scope="col" class="px-3 py-3">
						<span
							class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
						>
							Nombre de lots
						</span>
					</th>
					<th scope="col" class="px-3 py-3">
						<span
							class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
						>
							Nbr produits par lots
						</span>
					</th>
					<th scope="col" class="px-6 py-3">
						<span
							class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
						>
							Prix d'un lot HT
						</span>
					</th>
					<th scope="col" class="px-6 py-3">
						<span
							class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
						>
							TVA
						</span>
					</th>
					<th scope="col" class="px-6 py-3">
						<span
							class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
						>
							Prix d'un lot TTC
						</span>
					</th>
					<th scope="col" class="bg-blue-800 px-6 py-3">
						<span
							class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
						>
							Ajouter / Supprimer
						</span>
					</th>
				</tr>
			</thead>
			<tr>
				<td class="px-12 py-3">
					<div class="flex flex-col">
						<input
							type="text"
							class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-white text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							placeholder="Nom du produit"
							bind:value={itemName}
							on:input={(e) => {
								// @ts-ignore
								searchName = e.target.value.toLowerCase();
								reloadItems();
							}}
						/>
					</div>
					{#if searchName.length > 1}
						{#each items as item}
							<div
								class="bg-white p-2"
								on:click={() => {
									itemName = item.name;
									newItem.item_id = item.id;
									searchName = '';
									nameList.push(item.name);
								}}
							>
								{item.name}
							</div>
						{/each}
					{/if}
				</td>
				<td class="px-3 py-3">
					<div class="flex flex-col">
						<input
							type="number"
							class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-white text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							placeholder="Nombre de lots"
							min="0"
							max="1000"
							bind:value={newItem.amount_of_bundle}
						/>
					</div>
				</td>
				<td class="px-3 py-3">
					<div class="flex flex-col">
						<input
							type="number"
							class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-white text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							placeholder="Nombre de produits par lots"
							min="0"
							max="1000"
							bind:value={newItem.amount_per_bundle}
						/>
					</div>
				</td>
				<td class="px-6 py-3">
					<div class="flex flex-col">
						<input
							min="0"
							max="1000"
							step="0.01"
							type="number"
							class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-white text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							placeholder="Prix d'un lot HT"

							on:change={(e) => {
								// @ts-ignore
								newItem.bundle_cost_ht = parsePrice(e.target?.value);
							}}
						/>
					</div>
				</td>
				<td class="px-6 py-3">
					<div class="flex flex-col">
						<select
							class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-white text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							on:change={(e) => {
								// @ts-ignore
								newItem.tva = parseInt(e.target?.value);
							}}
						>
							<option value="0">0%</option>
							<option value="550">5.5%</option>
							<option value="1000">10%</option>
							<option value="2000">20%</option>
						</select>
					</div>
				</td>
				<td class="px-6 py-3">
					<div class="flex flex-col">
						<input
							min="0"
							max="1000"
							step="0.01"
							type="number"
							class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-white text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							value={formatPrice(newItem.bundle_cost_ht * (1 + newItem.tva / 10000))}
						/>
					</div>
				</td>
				<td class="px-6 py-3">
					<div class="flex flex-col">
						<button
							class="bg-green-600 hover:bg-green-700 text-white font-bold py-2 px-4 rounded"
							on:click={() => {
								newRestock.items.push(newItem);
								itemName = '';
								updateTotalHTandTTC();
							}}
						>
							Ajouter
						</button>
					</div>
				</td></tr
			>
			{#each newRestock.items as item}
				<tr>
					<td class="px-12 py-3">
						<div class="flex flex-col">
							<div
								class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-gray-300 text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							>
								<p>{nameList[newRestock.items.indexOf(item)]}</p>
							</div>
						</div>
					</td>
					<td class="px-3 py-3">
						<div class="flex flex-col">
							<div
								class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-gray-300 text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							>
								<p>{item.amount_of_bundle}</p>
							</div>
						</div>
					</td>
					<td class="px-3 py-3">
						<div class="flex flex-col">
							<div
								class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-gray-300 text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							>
								<p>{item.amount_per_bundle}</p>
							</div>
						</div>
					</td>
					<td class="px-6 py-3">
						<div class="flex flex-col">
							<div
								class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-gray-300 text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							>
								<p>{formatPrice(item.bundle_cost_ht)}</p>
							</div>
						</div>
					</td>
					<td class="px-6 py-3">
						<div class="flex flex-col">
							<div
								class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-gray-300 text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							>
								<p>{item.tva / 100}%</p>
							</div>
						</div>
					</td>
					<td class="px-6 py-3">
						<div class="flex flex-col">
							<div
								class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-gray-300 text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							>
								<p>{formatPrice(item.bundle_cost_ht * (1 + item.tva / 10000))}</p>
							</div>
						</div>
					</td>
					<td class="px-6 py-3">
						<div class="flex flex-col">
							<button
								class="bg-red-600 hover:bg-red-700 text-white font-bold py-2 px-4 rounded"
								on:click={() => {
									newRestock.items.splice(newRestock.items.indexOf(item), 1);
									updateTotalHTandTTC();
								}}
							>
								Supprimer
							</button>
						</div>
					</td>
				</tr>
			{/each}
		</table>

		<div class="flex p-2 m-8 bg-slate-600 items-center">
			<p class="font-bold text-white text-2xl">
				Ma réappro est irréprochable, et j'en suis responsable :
			</p>
			<input class="m-2 mr-auto max-w-lg w-6 h-6" type="checkbox" bind:checked={sure} />

			{#if sure}
				<button on:click={() => applyRestock()} class="bg-green-600 hover:bg-green-700 text-white font-bold py-2 px-4 rounded">
					<p class="font-bold text-white text-2xl">Terminer la réappro</p>
				</button>
			{/if}
		</div>
	</div>

	<table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
		<thead class="bg-gray-50 dark:bg-slate-800">
			<tr>
				<th scope="col" class="px-6 py-3">
					<span
						class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
					>
						Date
					</span>
				</th>
				<th scope="col" class="px-6 py-3">
					<span
						class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
					>
						Fournisseur
					</span>
				</th>
				<th scope="col" class="px-6 py-3">
					<span
						class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
					>
						Conducteur
					</span>
				</th>
				<th scope="col" class="px-2 py-3">
					<p
						class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
					>
						Prix total TTC
					</p>
				</th>
			</tr>
        </thead>
       

    </table>
</div>

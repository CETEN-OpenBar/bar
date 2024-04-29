<script lang="ts">
	import type { Item, NewRestock, NewRestockItem, Restock } from '$lib/api';
	import ConfirmationPopup from '$lib/components/confirmationPopup.svelte';
	import { itemsApi, restocksApi } from '$lib/requests/requests';
	import { formatPrice, parsePrice } from '$lib/utils';
	import { onMount } from 'svelte';

	let sure: boolean = false;
	let items: Item[] = [];

	let restocks: Restock[] = [];
	let newRestock: NewRestock = {
		total_cost_ht: 0,
		total_cost_ttc: 0,
		driver_id: '',
		type: 'promocash',
		items: []
	};

	let page = 0;
	let max_page = 0;
	let itemsPerPage = 4;

	let nameList: string[] = [];
	let newItem: NewRestockItem = {
		item_id: '',
		amount_of_bundle: 0,
		amount_per_bundle: 0,
		bundle_cost_ht: 0,
		bundle_cost_ttc: 0,
		bundle_cost_float_ttc: 0,
		tva: 0
	};
	let searchName: string = '';

	type dV = {
		name: string;
		item_price_calc: number;
		item_price: string;
		item_price_ht: string;
		amount_of_bundle: string;
		amount_per_bundle: string;
		bundle_cost_ht: string;
		tva: string;
		bundle_cost_ttc: string;
	};

	let displayedValues: dV = {
		name: 'Nom du produit',
		item_price_calc: 0,
		item_price: 'Prix coûtant TTC',
		item_price_ht: 'Prix coûtant HT',
		amount_of_bundle: 'Nombre de lots',
		amount_per_bundle: 'Nombre de produits par lots',
		bundle_cost_ht: "Prix d'un lot HT",
		tva: '0',
		bundle_cost_ttc: "Prix d'un lot TTC"
	};

	let deletingRestock = false;
	let deleteRestockCallback: VoidFunction = () => {};
	let confirmationMessage: string | undefined = undefined;

	let selectedRestock: Restock | undefined = undefined;

	onMount(() => {
		reloadItems();
		restocksApi()
			.getRestocks(page, itemsPerPage, undefined, undefined, {
				withCredentials: true
			})
			.then((res) => {
				restocks = res.data.restocks ?? [];
			});
	});

	function reloadItems() {
		itemsApi()
			.getAllItems(page, itemsPerPage, undefined, undefined, searchName, undefined, {
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
		newRestock.total_cost_ht = 0;
		newRestock.items.forEach((item) => {
			newRestock.total_cost_ht += item.amount_of_bundle * item.bundle_cost_ht;
		});
		newRestock.total_cost_ttc = 0.0;
		newRestock.items.forEach((item) => {
			if (item.bundle_cost_float_ttc === 0.0) {
				newRestock.total_cost_ttc += item.amount_of_bundle * item.bundle_cost_ttc;
			} else {
				newRestock.total_cost_ttc +=
					item.amount_of_bundle * (item.bundle_cost_float_ttc ?? item.bundle_cost_ttc);
			}
		});
	}

	async function applyRestock() {
		if (!sure) return;
		newRestock.driver_id = undefined;
		newRestock.total_cost_ttc = Math.round(newRestock.total_cost_ttc);
		newRestock.total_cost_ht = Math.round(newRestock.total_cost_ht);
		restocksApi()
			.createRestock(newRestock, { withCredentials: true })
			.then((res) => {
				restocks = [res.data, ...restocks];
				newRestock = {
					total_cost_ht: 0,
					total_cost_ttc: 0,
					driver_id: '',
					type: newRestock.type,
					items: []
				};
				displayedValues = {
					name: 'Nom du produit',
					item_price_calc: 0,
					item_price: 'Prix coûtant TTC',
					item_price_ht: 'Prix coûtant HT',
					amount_of_bundle: 'Nombre de lots',
					amount_per_bundle: 'Nombre de produits par lots',
					bundle_cost_ht: "Prix d'un lot HT",
					tva: '0',
					bundle_cost_ttc: "Prix d'un lot TTC"
				};
				newItem = {
					item_id: '',
					amount_of_bundle: 0,
					amount_per_bundle: 0,
					bundle_cost_ht: 0,
					bundle_cost_ttc: 0,
					bundle_cost_float_ttc: 0.0,
					tva: 0
				};
				sure = false;
			});
	}

	function updatePrices() {
		// Calculate from displayedValues.item_price_calc, displayedValues.amount_of_bundle and TVA
		if (newItem.amount_of_bundle === 0 || newItem.amount_per_bundle === 0) return;

		if (displayedValues.bundle_cost_ht === "Prix d'un lot HT") {
			newItem.bundle_cost_ht = Math.round(
				(displayedValues.item_price_calc * newItem.amount_per_bundle) / (1 + newItem.tva / 10000)
			);
			displayedValues.bundle_cost_ht = formatPrice(newItem.bundle_cost_ht);
		}
		if (displayedValues.bundle_cost_ttc === "Prix d'un lot TTC") {
			newItem.bundle_cost_ttc = Math.round(
				displayedValues.item_price_calc * newItem.amount_per_bundle
			);
			displayedValues.bundle_cost_ttc = formatPrice(newItem.bundle_cost_ttc);
		} else {
			newItem.bundle_cost_ttc = Math.round(newItem.bundle_cost_ht * (1 + newItem.tva / 10000));
			newItem.bundle_cost_float_ttc = newItem.bundle_cost_ht * (1 + newItem.tva / 10000);
			displayedValues.bundle_cost_ttc = formatPrice(newItem.bundle_cost_ttc);
		}
		displayedValues.item_price_ht = formatPrice(newItem.bundle_cost_ht * newItem.amount_of_bundle);
		displayedValues.item_price = formatPrice(newItem.bundle_cost_ttc * newItem.amount_of_bundle);
	}

	function deleteRestock(restockId: string) {
		restocksApi()
			.deleteRestock(restockId, { withCredentials: true })
			.then(() => {
				restocks = restocks.filter((ct) => ct.id !== restockId);
			});
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
		<div>
			<p class="dark:text-white text-2xl ml-5">
				Total HT : {formatPrice(newRestock.total_cost_ht)}
			</p>
			<p class="dark:text-white text-2xl ml-5">
				Total TTC : {formatPrice(newRestock.total_cost_ttc)}
			</p>
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
					<th scope="col" class="px-3 py-3 w-48">
						<span
							class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
						>
							Prix coûtant HT
						</span>
					</th>
					<th scope="col" class="px-3 py-3 w-48">
						<span
							class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
						>
							Prix coûtant TTC
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
					<th scope="col" class="px-6 py-3 w-48">
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
					<th scope="col" class="px-6 py-3 w-48">
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
				<td class="relative px-12 py-3">
					<div class="flex flex-col">
						<input
							type="text"
							class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-white text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent placeholder:text-slate-500"
							placeholder={displayedValues.name}
							on:input={(e) => {
								// @ts-ignore
								searchName = e.target.value.toLowerCase();
								reloadItems();
							}}
							on:change={(e) => {
								// @ts-ignore
								e.target.value = '';
							}}
							on:focusout={() => {
								setTimeout(() => {
									searchName = '';
								}, 200);
							}}
							on:focusin={(e) => {
								// @ts-ignore
								searchName = e.target.value.toLowerCase();
							}}
						/>
					</div>
					<div class="absolute rounded-b-lg bg-slate-100 -translate-y-2 flex flex-col">
						{#if searchName.length > 0}
							{#each items as item}
								<button
									class="p-2"
									on:click={() => {
										displayedValues.name = item.name;
										displayedValues.item_price = formatPrice(item.prices.coutant);
										displayedValues.item_price_ht = formatPrice(
											item.prices.coutant / (1 + (item.last_tva ?? 0) / 10000)
										);
										displayedValues.item_price_calc = item.prices.coutant;
										displayedValues.tva = (item.last_tva ?? 0).toString();
										newItem.tva = item.last_tva ?? 0;
										newItem.item_id = item.id;
										searchName = '';
									}}
								>
									{item.name}
								</button>
							{/each}
						{/if}
					</div>
				</td>
				<td class="px-3 py-3">
					<div class="flex flex-col">
						<input
							class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-white text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							disabled
							placeholder={displayedValues.item_price_ht}
						/>
					</div>
				</td>
				<td class="px-3 py-3">
					<div class="flex flex-col">
						<input
							class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-white text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							disabled
							placeholder={displayedValues.item_price}
						/>
					</div>
				</td>
				<td class="px-3 py-3">
					<div class="flex flex-col">
						<input
							type="number"
							class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-white text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							placeholder="Nombre de lots"
							min="1"
							max="1000"
							bind:value={newItem.amount_of_bundle}
							on:change={() => {
								updatePrices();
							}}
						/>
					</div>
				</td>
				<td class="px-3 py-3">
					<div class="flex flex-col">
						<input
							type="number"
							class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-white text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							placeholder="Nombre de produits par lots"
							min="1"
							max="1000"
							bind:value={newItem.amount_per_bundle}
							on:change={() => {
								updatePrices();
							}}
						/>
					</div>
				</td>
				<td class="px-6 py-3">
					<div class="flex flex-col">
						<input
							min="0"
							max="100000"
							type="number"
							class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-white text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							placeholder={displayedValues.bundle_cost_ht}
							on:change={(e) => {
								// @ts-ignore
								newItem.bundle_cost_ht = parsePrice(e.target?.value);
								let r = formatPrice(newItem.bundle_cost_ht);
								displayedValues.bundle_cost_ht = r;
								// @ts-ignore
								e.target.value = r;
								newItem.bundle_cost_ttc = Math.round(
									newItem.bundle_cost_ht * (1 + newItem.tva / 10000)
								);
								newItem.bundle_cost_float_ttc = newItem.bundle_cost_ht * (1 + newItem.tva / 10000);
								displayedValues.bundle_cost_ttc = formatPrice(
									Number((newItem.bundle_cost_ht * (1 + newItem.tva / 10000)).toFixed(0))
								);
							}}
						/>
					</div>
				</td>
				<td class="px-6 py-3">
					<div class="flex flex-col">
						<select
							class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-white text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							id="tva"
							on:change={(e) => {
								// @ts-ignore
								newItem.tva = parseInt(e.target?.value);
								updatePrices();
							}}
							bind:value={displayedValues.tva}
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
							max="100000"
							type="number"
							class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-white text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							placeholder={displayedValues.bundle_cost_ttc}
							on:change={(e) => {
								newItem.bundle_cost_ht = Math.round(
									// @ts-ignore
									parsePrice(e.target?.value) / (1 + (newItem.tva ?? 0) / 10000)
								);
								// @ts-ignore
								newItem.bundle_cost_ttc = parsePrice(e.target?.value);
								let r = formatPrice(newItem.bundle_cost_ttc);
								displayedValues.bundle_cost_ttc = r;
								displayedValues.bundle_cost_ht = formatPrice(newItem.bundle_cost_ht);
								// @ts-ignore
								e.target.value = r;
							}}
						/>
					</div>
				</td>
				<td class="px-6 py-3">
					<div class="flex flex-col">
						<button
							class="bg-green-600 hover:bg-green-700 text-white font-bold py-2 px-4 rounded"
							on:click={() => {
								let t = newRestock.items;
								let n = nameList;
								t.unshift(newItem);
								n.unshift(displayedValues.name);
								newRestock.items = t;
								nameList = n;
								displayedValues = {
									name: 'Nom du produit',
									item_price_calc: 0,
									item_price: 'Prix coûtant TTC',
									item_price_ht: 'Prix coûtant HT',
									amount_of_bundle: 'Nombre de lots',
									amount_per_bundle: 'Nombre de produits par lots',
									bundle_cost_ht: "Prix d'un lot HT",
									tva: '0',
									bundle_cost_ttc: "Prix d'un lot TTC"
								};
								newItem = {
									item_id: '',
									amount_of_bundle: 0,
									amount_per_bundle: 0,
									bundle_cost_ht: 0,
									bundle_cost_ttc: 0,
									bundle_cost_float_ttc: 0.0,
									tva: 0
								};
								updateTotalHTandTTC();
							}}
						>
							Ajouter
						</button>
					</div>
				</td></tr
			>
			{#each newRestock.items as item, i}
				<tr>
					<td class="px-12 py-3">
						<div class="flex flex-col">
							<div
								class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-gray-300 text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							>
								<p>{nameList[i]}</p>
							</div>
						</div>
					</td>
					<td class="px-3 py-3">
						<div class="flex flex-col">
							<div
								class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-gray-300 text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							>
								<p>{formatPrice(item.bundle_cost_ht * item.amount_of_bundle)}</p>
							</div>
						</div>
					</td>
					<td class="px-3 py-3">
						<div class="flex flex-col">
							<div
								class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-gray-300 text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							>
								<p>
									{formatPrice(item.bundle_cost_ttc * item.amount_of_bundle)}
								</p>
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
								<p>{formatPrice(item.bundle_cost_ttc)}</p>
							</div>
						</div>
					</td>
					<td class="px-6 py-3">
						<div class="flex flex-col">
							<button
								class="bg-red-600 hover:bg-red-700 text-white font-bold py-2 px-4 rounded"
								on:click={() => {
									newRestock.items = newRestock.items.filter((_, index) => index !== i);
									nameList = nameList.filter((_, index) => index !== i);
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
			<input
				id="CHECKBOX"
				class="m-2 mr-auto max-w-lg w-6 h-6"
				type="checkbox"
				bind:checked={sure}
			/>

			{#if sure}
				<button
					on:click={() => applyRestock()}
					class="bg-green-600 hover:bg-green-700 text-white font-bold py-2 px-4 rounded"
				>
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
				<th scope="col" class="px-2 py-3">
					<p
						class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
					>
						Actions
					</p>
				</th>
			</tr>

			{#if deletingRestock}
				<ConfirmationPopup
					message={confirmationMessage}
					confirm_text="Supprimer"
					cancel_callback={() => {
						deletingRestock = false;
					}}
					confirm_callback={deleteRestockCallback}
				/>
			{/if}
			{#each restocks as restock}
				<tr>
					<td class="px-6 py-3">
						<div class="flex flex-col">
							<div
								class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-gray-300 text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							>
								<p>{restock.created_at}</p>
							</div>
						</div>
					</td>
					<td class="px-6 py-3">
						<div class="flex flex-col">
							<div
								class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-gray-300 text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							>
								<p>{restock.type}</p>
							</div>
						</div>
					</td>
					<td class="px-6 py-3">
						<div class="flex flex-col">
							<div
								class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-gray-300 text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							>
								<p>{restock.created_by_name}</p>
							</div>
						</div>
					</td>
					<td class="px-2 py-3">
						<p
							class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
						>
							{formatPrice(restock.total_cost_ttc)}
						</p>
					</td>
					<td class="px-2 py-3">
						<div
							class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
						>
							<button
								class="px-2 inline-flex items-center text-sm text-blue-600 decoration-2 hover:underline font-medium"
								on:click={() => {
									deleteRestockCallback = () => {
										deletingRestock = false;
										deleteRestock(restock.id);
									};
									confirmationMessage =
										'Supprimer la réappro de ' +
										restock.created_by_name +
										' à ' +
										restock.type +
										' ?';
									deletingRestock = true;
								}}
							>
								Supprimer
							</button>
							<button
								class="px-2 inline-flex items-center text-sm text-blue-600 decoration-2 hover:underline font-medium"
								on:click={() => {
									if (selectedRestock == restock) {
										selectedRestock = undefined;
									} else {
										selectedRestock = restock;
									}
								}}
							>
								Voir
							</button>
						</div>
					</td>
				</tr>
			{/each}
		</thead>
	</table>
	{#if selectedRestock != undefined}
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
					<th scope="col" class="px-3 py-3 w-48">
						<span
							class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
						>
							Prix coûtant HT
						</span>
					</th>
					<th scope="col" class="px-3 py-3 w-48">
						<span
							class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
						>
							Prix coûtant TTC
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
					<th scope="col" class="px-6 py-3 w-48">
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
					<th scope="col" class="px-6 py-3 w-48">
						<span
							class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
						>
							Prix d'un lot TTC
						</span>
					</th>
				</tr>
			</thead>
			{#each selectedRestock.items as item}
				<tr>
					<td class="px-12 py-3">
						<div class="flex flex-col">
							<div
								class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-gray-300 text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							>
								<p>{item.item_name}</p>
							</div>
						</div>
					</td>
					<td class="px-3 py-3">
						<div class="flex flex-col">
							<div
								class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-gray-300 text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							>
								<p>{formatPrice(item.bundle_cost_ht * item.amount_of_bundle)}</p>
							</div>
						</div>
					</td>
					<td class="px-3 py-3">
						<div class="flex flex-col">
							<div
								class="rounded-lg border-transparent flex-1 appearance-none border border-gray-300 w-full py-2 px-4 bg-gray-300 text-gray-700 placeholder-gray-400 shadow-sm text-base focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
							>
								<p>
									{formatPrice(item.bundle_cost_ttc * item.amount_of_bundle)}
								</p>
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
								<p>{formatPrice(item.bundle_cost_ttc)}</p>
							</div>
						</div>
					</td>
				</tr>
			{/each}
		</table>
	{/if}
</div>

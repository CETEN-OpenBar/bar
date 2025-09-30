<script lang="ts">
	import { CourseApi, restocksApi, categoriesApi } from '$lib/requests/requests';
	import {
		type CourseItem,
		type NewRestock,
		type NewRestockItem,
		type Category,
		RestockState,
		RestockType
	} from '$lib/api';
	import { formatPrice, fournisseurIterator } from '$lib/utils';
	import { goto } from '$app/navigation';
	import { api } from '$lib/config/config';
	import { onMount } from 'svelte';
	let items: CourseItem[] = [];
	let categories: Category[] = [];
	let groupedItems: { [categoryId: string]: CourseItem[] } = {};
	let fournisseur = RestockType.RestockPromocash;
	let newRestock: NewRestock = {
		items: [],
		total_cost_ht: 0,
		total_cost_ttc: 0,
		type: fournisseur,
		state: RestockState.RestockPending
	};
	let newRestockItem: NewRestockItem = {
		item_id: '',
		item_name: '',
		amount_of_bundle: 0,
		amount_per_bundle: 0,
		bundle_cost_ht: 0,
		bundle_cost_ttc: 0,
		bundle_cost_float_ttc: 0,
		tva: 0
	};
	function reloadCourse() {
		CourseApi()
			.getCourse(fournisseur, {
				withCredentials: true
			})
			.then((res) => {
				if (res.data.items != null) {
					items = res.data.items.sort((a, b) => a.item.name.localeCompare(b.item.name));
					groupItemsByCategory();
				} else {
					items = [];
					groupedItems = {};
				}
			});
	}
	function loadCategories() {
		categoriesApi()
			.getCategories(true, { withCredentials: true })
			.then((res) => {
				categories = res.data ?? [];
				groupItemsByCategory();
			})
			.catch((error) => {
				console.error('Failed to load categories:', error);
			});
	}
	function groupItemsByCategory() {
		if (items.length === 0 || categories.length === 0) return;
		
		groupedItems = {};
		
		// Group items by category
		items.forEach(item => {
			const categoryId = item.item.category_id;
			if (!groupedItems[categoryId]) {
				groupedItems[categoryId] = [];
			}
			groupedItems[categoryId].push(item);
		});
		
		// Sort items within each category alphabetically
		Object.keys(groupedItems).forEach(categoryId => {
			groupedItems[categoryId].sort((a, b) => a.item.name.localeCompare(b.item.name));
		});
	}
	onMount(() => {
		reloadCourse();
		loadCategories();
	});
	function addNewRestockItem(courseItem: CourseItem) : CourseItem {
		if (courseItem.amountToBuy <= 0) {
			courseItem.amountToBuy = 0;
		}

		newRestockItem = {
			item_id: courseItem.item.id,
			item_name: courseItem.item.name,
			amount_of_bundle: courseItem.amountToBuy,
			amount_per_bundle: courseItem.item.amount_per_bundle ? courseItem.item.amount_per_bundle : 1,
			tva: courseItem.item.last_tva ? courseItem.item.last_tva : 0,
			bundle_cost_ht: 0,
			bundle_cost_ttc: 0
		};
		newRestockItem.bundle_cost_ttc = Math.round(
			courseItem.item.prices.coutant * newRestockItem.amount_per_bundle * newRestockItem.amount_of_bundle
		);
		newRestockItem.bundle_cost_ht = Math.round(
			(newRestockItem.bundle_cost_ttc) /
				(1 + newRestockItem.tva / 10000)
		);
		newRestock.items.push(newRestockItem);
		updateHTandTTC();
		return courseItem
	}

	function removeNewRestockItem(courseItem: CourseItem) {
		newRestock.items = newRestock.items.filter((item) => item.item_id !== courseItem.item.id);
		updateHTandTTC();
	}

	function updateHTandTTC() {
		newRestock.total_cost_ht = newRestock.items.reduce((acc, item) => acc + item.bundle_cost_ht, 0);
		newRestock.total_cost_ttc = newRestock.items.reduce(
			(acc, item) => acc + item.bundle_cost_ttc,
			0
		);
	}

	function generateRestock() {
		restocksApi()
			.createRestock(newRestock, {
				withCredentials: true
			})
			.then(() => {
				goto('/panel/products/restocks');
			});
	}
</script>

<div class="fixed inset-0 flex flex-col bg-white dark:bg-gray-900">
	<div class="flex-none bg-white dark:bg-gray-900 border-b border-gray-200 dark:border-gray-700 px-6 py-4">
		<h1 class="text-2xl font-semibold text-gray-800 dark:text-white">Réapprovisionnement</h1>
	</div>

	<div class="flex-none px-6 py-3 bg-white dark:bg-gray-900 border-b border-gray-200 dark:border-gray-700">
		<div class="flex items-center gap-3">
			<label for="fournisseur" class="text-gray-800 dark:text-white text-sm font-medium">Fournisseur :</label>
			<select
				id="fournisseur"
				name="fournisseur"
				class="py-3 px-4 w-64 border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
				required
				aria-describedby="text-error"
							on:change={(e) => {
								// @ts-ignore
								let val = e.target?.value;
								fournisseur = val;
								reloadCourse();
								loadCategories();
							}}
			>
			{#each fournisseurIterator as [val, name]}
				<option value="{val}" selected={val === RestockType.RestockPromocash}>{name}</option>
			{/each}
			</select>
		</div>
	</div>

	<div class="flex-1 overflow-y-auto px-6 py-4">
		{#each Object.entries(groupedItems) as [categoryId, categoryItems]}
			{#if categories.find(c => c.id === categoryId)}
				<div class="mb-6">
					<h2 class="text-xl font-semibold text-gray-800 dark:text-white mb-4">{categories.find(c => c.id === categoryId).name}</h2>
					<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
						{#each categoryItems as item}
							<div class="bg-white dark:bg-gray-800 rounded-lg shadow-md border border-gray-200 dark:border-gray-700 p-4">
								<div class="flex items-center mb-3">
									{#if item.item.picture_uri}
										<img 
											src="{api() + item.item.picture_uri}" 
											alt="{item.item.name}"
											class="w-16 h-16 object-cover rounded-md mr-3"
										/>
									{:else}
										<div class="w-16 h-16 bg-gray-200 dark:bg-gray-700 rounded-md mr-3 flex items-center justify-center">
											<svg class="w-8 h-8 text-gray-400" fill="currentColor" viewBox="0 0 20 20">
												<path fill-rule="evenodd" d="M4 3a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V5a2 2 0 00-2-2H4zm12 12H4l4-8 3 6 2-4 3 6z" clip-rule="evenodd"/>
											</svg>
										</div>
									{/if}
									<div class="flex-1">
										<h3 class="font-semibold text-gray-800 dark:text-gray-200 text-sm truncate">{item.item.name}</h3>
										<div class="text-xs text-gray-600 dark:text-gray-400 mt-1">
											<div>Stock: {item.item.amount_left}</div>
											<div>Optimal: {item.item.optimal_amount}</div>
										</div>
									</div>
								</div>
								
								<div class="flex items-center justify-between mt-3">
									<div class="flex items-center">
										<input
											class="w-20 py-2 px-2 border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
											type="number"
											min="0"
											bind:value={item.amountToBuy}
											on:input={() => {
												if (newRestock.items.some((restockItem) => restockItem.item_id === item.item.id)) {
													if (item.amountToBuy !== null) {
														removeNewRestockItem(item);
														item = addNewRestockItem(item);
													}
												}
											}}
											placeholder="0"
										/>
										<span class="ml-2 text-sm text-gray-600 dark:text-gray-400">packs</span>
									</div>
									<input
										class="w-6 h-6 border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
										type="checkbox"
										on:change={(event) => {
											// @ts-ignore
											if (event.target?.checked) {
												item = addNewRestockItem(item);
											} else {
												removeNewRestockItem(item);
											}
										}}
									/>
								</div>
							</div>
						{/each}
					</div>
				</div>
			{/if}
		{/each}
	</div>

	<div class="flex-none bg-white dark:bg-gray-900 border-t border-gray-200 dark:border-gray-700 px-6 py-4">
		<div class="flex flex-col md:flex-row justify-between items-center gap-3 md:gap-4">
			<div class="flex gap-4 md:gap-8">
				<p class="text-base md:text-xl text-gray-800 dark:text-white">Total HT: {formatPrice(newRestock.total_cost_ht)}</p>
				<p class="text-base md:text-xl text-gray-800 dark:text-white">Total TTC: {formatPrice(newRestock.total_cost_ttc)}</p>
			</div>
			<button
				class="text-base md:text-xl bg-blue-700 px-6 py-2 rounded-xl hover:bg-blue-900 transition-all text-white w-full md:w-auto"
				on:click={() => {
					generateRestock();
				}}
			>
				Générer la réappro
			</button>
		</div>
	</div>
</div>

<script lang="ts">
	import { CourseApi, restocksApi } from '$lib/requests/requests';
	import {
		type CourseItem,
		type NewRestock,
		type NewRestockItem,
		RestockState,
		RestockType
	} from '$lib/api';
	import { formatPrice } from '$lib/utils';
	import { goto } from '$app/navigation';
	let items: CourseItem[] = [];
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
				} else {
					items = [];
				}
			});
	}
	reloadCourse();
	function addNewRestockItem(courseItem: CourseItem) {
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

<div class="relative mt-4 w-96 md:mt-0">
	<!-- filter by state -->
	<select
		id="fournisseur"
		name="fournisseur"
		class="py-3 px-4 block w-full border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
		required
		aria-describedby="text-error"
		on:change={(e) => {
			// @ts-ignore
			let val = e.target?.value;
			fournisseur = val;
			reloadCourse();
		}}
	>
		<option value="promocash">Promocash</option>
		<option value="auchan_drive">Auchan drive</option>
		<option value="auchan">Auchan</option>
		<option value="viennoiserie">Boulangerie Benoist</option>
	</select>
</div>

<table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
	<thead class="bg-gray-50 dark:bg-slate-800">
		<tr>
			<th scope="col" class="px-6 py-3">
				<span
					class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
				>
					Nom
				</span>
			</th>
			<th scope="col" class="px-6 py-3">
				<span
					class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
				>
					Nombre à acheter
				</span>
			</th>
			<th scope="col" class="px-3">
				<span
					class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
				>
					Checkbox
				</span>
			</th>
		</tr>
	</thead>
	<tbody class="divide-y divide-gray-200 dark:divide-gray-700">
		{#each items as item}
			<tr>
				<td class="px-6 py-4">
					<p
						class="py-3 px-2 block border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
					>
						{item.item.name}
					</p>
				</td>
				<td class="px-6 py-4">
					<input
						class="w-full py-3 px-2 block border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
						type="number"
						min="0"
						bind:value={item.amountToBuy}
						on:input={() => {
							if (newRestock.items.some((restockItem) => restockItem.item_id === item.item.id)) {
								if (item.amountToBuy !== null) {
									removeNewRestockItem(item);
									addNewRestockItem(item);
								}
						}}}
					>
				</td>
				<td class="text-center px-3">
					<input
						class="w-12 h-12 block border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 mx-auto"
						type="checkbox"
						on:change={(event) => {
							// @ts-ignore
							if (event.target?.checked) {
								addNewRestockItem(item);
							} else {
								removeNewRestockItem(item);
							}
						}}
					/>
				</td>
			</tr>
		{/each}
	</tbody>
</table>
<div class="flex justify-center items-center w-full mt-4 px-6 space-x-8">
	<p class="text-2xl text-white">Total HT: {formatPrice(newRestock.total_cost_ht)}</p>
	<p class="text-2xl text-white">Total TTC: {formatPrice(newRestock.total_cost_ttc)}</p>
	<button
		class="text-xl bg-blue-700 p-2 rounded-xl hover:bg-blue-900 transition-all"
		on:click={() => {
			generateRestock();
		}}
	>
		Générer la réappro
	</button>
</div>

<script lang="ts">
	import { formatDateTime, formatPrice, time2Utc } from '$lib/utils';
	import type { Item, Transaction } from '$lib/api';
	import { itemsApi, transactionsApi } from '$lib/requests/requests';
	import { onMount } from 'svelte';
	import { stringify } from 'postcss';

	let items: Item[] = [];
	let itemsPerPage = 5;
	let searchItem: string = '';
	let selectedItem: Item | undefined = undefined;
	let hide_remotes = true;
	async function reloadItems() {
		let resp = await itemsApi().getAllItems(1, itemsPerPage, undefined, undefined, searchItem, undefined, undefined, {
			withCredentials: true
		});
		items = resp.data.items ?? [];
	}

	let transactions: Transaction[] = [];
	let reloading = false;

	async function reloadTransactions() {
		if (reloading) return;
		reloading = true;
		try {
			if (selectedItem != undefined){
				itemID = selectedItem.id;
			}
			else{
				var itemID = undefined;
			}
		
			let resp = await transactionsApi().getTransactions(1, undefined, "finished", true, !hide_remotes, undefined, startDate, endDate, itemID, {
				withCredentials: true
			});
		
			transactions = resp.data.transactions ?? [];
			

		} finally {
			reloading = false;
		}
	}

	let todayMorning = new Date();
	todayMorning.setHours(0, 0, 0, 0);
	let startDate = time2Utc(todayMorning.getTime() / 1000)
	let endDate = time2Utc(todayMorning.getTime() / 1000 + 24 * 60 * 60);
</script>

<div class="w-full flex flex-col items-center">
	<div class="flex flex-col p-5 gap-3">
		<div class="flex flex-row">
			<input
				class="rounded-l-md bg-slate-200 p-4"
				placeholder="Choisir un produit"
				type="text"
				list="items"
				bind:value={searchItem}
				on:keydown={reloadItems}
				on:keyup={reloadItems}
				on:change={() => {
					selectedItem = items.find((item) => item.name == searchItem); 
				}}
			/>
			
			<button class="rounded-r-lg bg-slate-200 p-4"> &#x1F50D; </button>
		</div>
	
	
		<datalist id="items">
			{#each items as item}
				<option value={item.name} />
			{/each}
		</datalist>
		<div class="flex flex-row gap-2 red-tex">
		<h1 class="text-l text-white">En ligne </h1>
		<input 
				type="checkbox" 
				class="h-6 w-6 align-middle ml-1"
				bind:checked={hide_remotes}
			/>
		</div>
		<div class="flex flex-row gap-2">
			<div class="flex flex-col">
				<h1 class="text-md font-semibold self-center text-white">Début:</h1>
				<input
					class="rounded-md bg-slate-200 p-2"
					type="date"
					value={todayMorning.toLocaleString('default', { year: 'numeric' }) +
						'-' +
						todayMorning.toLocaleString('default', { month: '2-digit' }) +
						'-' +
						todayMorning.toLocaleString('default', { day: '2-digit' })}
					on:change={(e) => {
						// @ts-ignore
						startDate = time2Utc(new Date(e.target.value).getTime() / 1000);
					}}
				/>
			</div>
			<div class="flex flex-col">
				<h1 class="text-md font-semibold self-center text-white">Fin:</h1>
				<input
					class="rounded-md bg-slate-200 p-2"
					type="date"
					value={todayMorning.toLocaleString('default', { year: 'numeric' }) +
						'-' +
						todayMorning.toLocaleString('default', { month: '2-digit' }) +
						'-' +
						todayMorning.toLocaleString('default', { day: '2-digit' })}
					on:change={(e) => {
						// @ts-ignore
						endDate = time2Utc(new Date(e.target.value).getTime() / 1000 + 24 * 60 * 60);
			
					}}
				/>
			</div>
		</div>

		<button
			class="rounded-md bg-slate-200 p-2 disabled:bg-slate-400"
			disabled={reloading}
			on:click={() => {
				reloadTransactions();
			}}>Rechercher</button
		>
	</div>
	<div
		class="flex flex-col items-center bg-blue-200 rounded-md p-5 overflow-auto"
		style="max-height: 50vh"
	>
		<div class="flex flex-col items-center">
			<h1 class="text-lg font-semibold">Stats (pour: {selectedItem?.name ?? ''})</h1>

			<table class="mt-4">
				<thead>
					<tr>
						<th class="px-4">Prix coûtant</th>
						<th class="px-4">Total sur la période</th>
						<th class="px-4">Nombre de produit</th>
					</tr>
				</thead>
				<tbody>
					<tr>
						{#if selectedItem}
							<td class="px-4">{formatPrice(selectedItem?.prices['coutant'] ?? -1)}</td>
						{:else}
							<td class="px-4" />
						{/if}
						<td class="px-4">
							{formatPrice(
								transactions
									.flatMap((t) => t.items)
									.filter((item) => item.item_id == selectedItem?.id)
									.reduce((acc, item) => acc + item.unit_cost * item.item_amount, 0)
							)}
						</td>
						<td class="px-4">
							<!-- Count amount of items in transaction.items[x].item_amount and transaction.items[x].menu_items[y].item_amount -->
							{transactions
								.flatMap((t) => t.items)
								.filter((item) => item.item_id == selectedItem?.id)
								.reduce((acc, item) => acc + item.item_amount, 0) +
								transactions
									.flatMap((t) => t.items)
									.flatMap((item) => item.menu_items ?? [])
									.filter((item) => item.id == selectedItem?.id)
									.reduce((acc, item) => acc + item.amount, 0) +
								transactions
									.flatMap((t) => t.items)
									.flatMap((item) => item.picked_categories_items ?? [])
									.filter((item) => item.item_id == selectedItem?.id)
									.reduce((acc, item) => acc + item.item_amount, 0)}
						</td>
					</tr>
				</tbody>
			</table>
		</div>

		<div class="w-full p-1 bg-blue-300 mt-4 mb-4" />

		<table>
			<thead>
				<tr>
					<th class="px-4">Client</th>
					<th class="px-4">Date</th>
					<th class="px-4">Prix d'achat</th>
					<th class="px-4">Quantité</th>
				</tr>
			</thead>
			<tbody>
				{#each transactions as t, index (t.id)}
						<tr>
							<td class="px-4">{t.account_name}</td>
							<td class="px-4">{formatDateTime(t.created_at)}</td>
							<td class="px-4">
								{formatPrice(
									t.items
										.filter((item) => item.item_id == selectedItem?.id)
										.reduce((acc, item) => acc + item.unit_cost * item.item_amount, 0)
								)}
							</td>
							<td class="px-4">
								{t.items
									.filter((item) => item.item_id == selectedItem?.id)
									.reduce((acc, item) => acc + item.item_amount, 0) +
									t.items
										.flatMap((item) => item.menu_items ?? [])
										.filter((item) => item.id == selectedItem?.id)
										.reduce((acc, item) => acc + item.amount, 0) +
									t.items
										.flatMap((item) => item.picked_categories_items ?? [])
										.filter((item) => item.item_id == selectedItem?.id)
										.reduce((acc, item) => acc + item.item_amount, 0)}
							</td>
						</tr>
			
				{/each}
			</tbody>
		</table>
	</div>
</div>

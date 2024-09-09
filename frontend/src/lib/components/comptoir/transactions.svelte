<script lang="ts">
	import type { Transaction, TransactionItem, TransactionState } from '$lib/api';
	import { api } from '$lib/config/config';
	import { transactionsApi } from '$lib/requests/requests';
	import { formatPrice } from '$lib/utils';
	import { onDestroy, onMount } from 'svelte';
	import TransactionPopup from './transactionPopup.svelte';
	import { dragscroll } from '@svelte-put/dragscroll';

	export let amount: number = 4;

	let transactions: Array<Transaction> = [];
	let maxItemPerTransaction: number = 6;
	let interval: number;
	let searchName: string | undefined;

	let page: number = 0;
	let maxPage: number = 0;
	let nextPage = () => {
		if (page < maxPage) {
			page++;
			reloadTransactions();
		}
	};
	let prevPage = () => {
		if (page > 0) {
			page--;
			reloadTransactions();
		}
	};

	onMount(() => {
		reloadTransactions();
		interval = setInterval(() => {
			reloadTransactions();
		}, 2000);
	});

	onDestroy(() => {
		clearInterval(interval);
	});

	let st: TransactionState | undefined = 'started';
	let showRemoteTransactions: boolean = true;

	type TransactionItemWithFakeAmount = TransactionItem & { item_fake_amount?: number };

	function reloadTransactions() {
		transactionsApi()
			.getTransactions(page, amount, st, !showRemoteTransactions, searchName, { withCredentials: true })
			.then((res) => {
				page = res.data.page ?? 0;
				maxPage = res.data.max_page ?? 0;
				if (!(res.data.transactions instanceof Array)) return;
				let countedItems = 0;
				let newTransactions = [];
				for (let transaction of res.data.transactions) {
					let items: Array<TransactionItemWithFakeAmount> = [];

					for (let r of transaction.items ?? []) {
						let item = r as TransactionItemWithFakeAmount;
						if (countedItems >= maxItemPerTransaction) break;
						if (countedItems + item.item_amount > maxItemPerTransaction) {
							item.item_fake_amount = maxItemPerTransaction - countedItems;
						}
						items.push(item);
					}
					transaction.items = items;
					newTransactions.push(transaction);
				}
				transactions = newTransactions;
			});
	}

	let displayTransaction: Transaction | null = null;
</script>

{#if displayTransaction}
	<TransactionPopup
		transaction={displayTransaction}
		close={() => {
			displayTransaction = null;
		}}
	/>
{/if}

<!-- Good looking dropdown for transaction -->
<div class="w-full h-full flex flex-col">
	<div class="flex flex-col flex-grow flex-1">
		<div class="flex flex-row items-center mt-2">
			<div class="flex flex-row items-center space-x-10 grow">
				<div class="text-lg font-semibold">Transactions</div>
				<input 
					class="rounded-lg p-2 text-black"
					placeholder="Rechercher une personne"
					on:input={(e) => {
						// @ts-ignore
						searchName = e.target.value.toLowerCase();
						page = 1;
						reloadTransactions();
					}}
				/>
			</div>
			<div class="text-lg font-semibold grow flex flex-row items-center space-x-5">
				<div>
					Filtre :
					<button
						class="bg-gray-700 p-2 rounded-lg"
						on:click={() => {
							if (st == 'started') {
								st = undefined;
							} else {
								st = 'started';
							}
							reloadTransactions();
						}}
					>
						{#if st == 'started'}
							En cours
						{:else}
							Tout
						{/if}
					</button>
				</div>
				<div>
					<label>
						Commandes en ligne
						<input 
							type="checkbox" 
							class="h-6 w-6 align-middle ml-1"
							bind:checked={showRemoteTransactions}
							on:change={reloadTransactions}
						/>
					</label>
				</div>
			</div>
			<div class="text-lg font-semibold grow text-end">Montant</div>
		</div>
		<!-- show clearly if there's a scrollbar -->
		<div use:dragscroll class="flex flex-col overflow-auto">
			{#each transactions as transaction}
				<button
					on:click={() => (displayTransaction = transaction)}
					class="flex flex-col mt-5 border-4 border-white rounded-xl {transaction.state == 'started'
						? 'animate-pulse bg-green-100'
						: ''} {transaction.state == 'canceled' ? 'bg-red-200' : ''} {transaction.state ==
					'finished'
						? 'bg-green-200'
						: ''}"
				>
					<div class="text-black">
						{#if transaction.is_remote}
							<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512" class="h-6 m-1 absolute"><!--!Font Awesome Free 6.6.0 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license/free Copyright 2024 Fonticons, Inc.--><path d="M483 159.7c10.9-24.6 21.4-60.4 21.4-87.9 0-72.7-79.6-98.4-209.7-38.6-107.6-7.2-211.2 73.7-237.1 186.5 30.9-34.9 78.3-82.3 122-101.2C125.4 166.9 79.1 228 44 291.7 23.2 329.7 0 390.9 0 436.7c0 98.6 92.9 86.5 180.3 42 31.4 15.4 66.6 15.6 101.7 15.6 97.1 0 184.2-54.3 216.8-146H377.9c-52.5 88.6-196.8 53-196.8-47.4H509.9c6.4-43.6-1.7-95.7-26.9-141.2zM64.6 346.9c17.7 51.2 53.7 95.9 100.3 123.3-88.7 48.9-173.3 29.1-100.3-123.3zm116-108.9c2-55.2 50.3-94.9 104-94.9 53.4 0 102 39.7 104 94.9H180.5zm184.5-187.6c21.4-10.3 48.6-22 72.6-22 31.4 0 54.3 21.7 54.3 53.7 0 20-7.4 49-14.6 67.9-26.3-42.3-66-81.6-112.3-99.6z"/></svg>
						{/if}
						Commande de : <b>{transaction.account_nick_name || transaction.account_name}</b>
					</div>
					<!-- first col is 90% second is 1% third is 9% -->
					<div class="grid grid-cols-11 grid-flow-row-dense gap-2">
						<div class="p-1 h-full pr-1 w-full col-span-10 border-r border-gray-400">
							<div class="grid grid-cols-3 gap-2">
								{#each transaction.items as item}
									<!-- One for each item.amount -->
									<div class="flex flex-col justify-center">
										<img
											src={api() + item.picture_uri}
											alt="ca charge"
											class="w-6 h-6 rounded-2xl self-center"
										/>
										<div class="text-xs text-center text-black">{item.item_name}</div>
										<div class="text-xs text-center text-black">x {item.item_amount}</div>
									</div>
								{/each}
							</div>
						</div>
						<div class="p-1 pl-1 w-full text-lg self-center text-center text-black col-span-1">
							{formatPrice(transaction.total_cost)}
						</div>
					</div>
				</button>
			{/each}
		</div>
	</div>

	<!-- Pagination -->
	<div class="flex flex-row justify-center mt-5">
		<button class="bg-blue-700 p-4 rounded-xl hover:bg-blue-900 transition-all text-2xl" on:click={prevPage}
			>&lt;</button
		>
		<div class="text-2xl font-semibold self-center mx-2">{page}/{maxPage}</div>
		<button class="bg-blue-700 p-4 rounded-xl hover:bg-blue-900 transition-all text-2xl" on:click={nextPage}
			>&gt;</button
		>
	</div>
</div>

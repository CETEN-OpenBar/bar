<script lang="ts">
	import type { Transaction, TransactionItem, TransactionState } from '$lib/api';
	import { api } from '$lib/config/config';
	import { transactionsApi } from '$lib/requests/requests';
	import { formatPrice } from '$lib/utils';
	import { onDestroy, onMount } from 'svelte';
	import TransactionPopup from './transactionPopup.svelte';
	import { dragscroll } from '@svelte-put/dragscroll';
	import { searchName } from '$lib/store/store';

	let searchNameValue:string;

	searchName.subscribe((value)=>{searchNameValue = value});	
	export let amount: number = 4;

	let transactions: Array<Transaction> = [];
	let maxItemPerTransaction: number = 6;
	let interval: number;
	// let searchName: string | undefined;

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
			.getTransactions(page, amount, st, !showRemoteTransactions, searchNameValue, undefined, undefined, undefined, { withCredentials: true })
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
			reloadTransactions();
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
					bind:value={$searchName}
					on:input={(e) => {
						// @ts-ignore
						searchNameValue = e.target.value.toLowerCase();
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
							<svg xmlns="http://www.w3.org/2000/svg" class="h-5 m-1 absolute" viewBox="0 0 640 512"><!--!Font Awesome Free 6.6.0 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license/free Copyright 2024 Fonticons, Inc.--><path d="M54.2 202.9C123.2 136.7 216.8 96 320 96s196.8 40.7 265.8 106.9c12.8 12.2 33 11.8 45.2-.9s11.8-33-.9-45.2C549.7 79.5 440.4 32 320 32S90.3 79.5 9.8 156.7C-2.9 169-3.3 189.2 8.9 202s32.5 13.2 45.2 .9zM320 256c56.8 0 108.6 21.1 148.2 56c13.3 11.7 33.5 10.4 45.2-2.8s10.4-33.5-2.8-45.2C459.8 219.2 393 192 320 192s-139.8 27.2-190.5 72c-13.3 11.7-14.5 31.9-2.8 45.2s31.9 14.5 45.2 2.8c39.5-34.9 91.3-56 148.2-56zm64 160a64 64 0 1 0 -128 0 64 64 0 1 0 128 0z"/></svg>
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

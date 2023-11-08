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

	let page: number = 0;
	let maxPage: number = 0;

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

	function reloadTransactions() {
		transactionsApi()
			.getTransactions(page, amount, st, { withCredentials: true })
			.then((res) => {
				page = res.data.page ?? 0;
				maxPage = res.data.max_page ?? 0;
				if (!(res.data.transactions instanceof Array)) return;
				let countedItems = 0;
				let newTransactions = [];
				for (let transaction of res.data.transactions) {
					let items: Array<TransactionItem> = [];

					for (let item of transaction.items ?? []) {
						if (countedItems >= maxItemPerTransaction) break;
						if (countedItems + item.item_amount > maxItemPerTransaction) {
							item.item_amount = maxItemPerTransaction - countedItems;
						}
						if (item.item_amount > 0) {
							items.push(item);
						}
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
<div class="w-full">
	<div class="flex flex-col">
		<div class="flex flex-row justify-between">
			<div class="text-lg font-semibold">Transactions</div>
			<div class="text-lg font-semibold">
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
			<div class="text-lg font-semibold">Montant</div>
		</div>
		<div class="flex flex-col max-h-[80vh] overflow-auto" use:dragscroll>
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
						Commande de : <b>{transaction.account_name}</b>
					</div>
					<div class="flex flex-row justify-between">
						<div class="p-5 h-full pr-4 w-full">
							<div class="grid grid-cols-3 gap-2">
								{#each transaction.items as item}
									<!-- One for each item.amount -->
									{#each Array(item.item_amount) as _}
										<div class="flex flex-col justify-center">
											<img
												src={api() + item.picture_uri}
												alt="ca charge"
												class="w-10 h-10 rounded-2xl self-center"
											/>
											<div class="text-center text-black">{item.item_name}</div>
										</div>
									{/each}
								{/each}
							</div>
						</div>
						<div class="border-r border-l border-gray-400" />
						<div class="p-5 pl-4 w-full text-lg self-center text-center text-black">
							{formatPrice(transaction.total_cost)}
						</div>
					</div>
				</button>
			{/each}
		</div>
	</div>

	<!-- Pagination -->
	<div class="flex flex-row justify-center mt-5">
		<button
			class="bg-blue-700 p-2 rounded-xl hover:bg-blue-900 transition-all"
			on:click={() => {
				if (page > 0) {
					page--;
					reloadTransactions();
				}
			}}>&lt;</button
		>
		<div class="text-lg font-semibold self-center mx-2">{page}/{maxPage + 1}</div>
		<button
			class="bg-blue-700 p-2 rounded-xl hover:bg-blue-900 transition-all"
			on:click={() => {
				if (page <= maxPage) {
					page++;
					reloadTransactions();
				}
			}}>&gt;</button
		>
	</div>
</div>

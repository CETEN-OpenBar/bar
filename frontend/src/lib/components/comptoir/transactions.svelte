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

	type TransactionItemWithFakeAmount = TransactionItem & { item_fake_amount?: number };

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
		<!-- show clearly if there's a scrollbar -->
		<div use:dragscroll class="flex flex-col max-h-[80vh] overflow-auto">
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
		<button class="bg-blue-700 p-2 rounded-xl hover:bg-blue-900 transition-all" on:click={prevPage}
			>&lt;</button
		>
		<div class="text-lg font-semibold self-center mx-2">{page}/{maxPage}</div>
		<button class="bg-blue-700 p-2 rounded-xl hover:bg-blue-900 transition-all" on:click={nextPage}
			>&gt;</button
		>
	</div>
</div>

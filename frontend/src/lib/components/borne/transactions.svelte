<script lang="ts">
	import type { Transaction, TransactionItem } from '$lib/api';
	import { api } from '$lib/config/config';
	import { transactionsApi } from '$lib/requests/requests';
	import { formatPrice } from '$lib/utils';
	import { onMount } from 'svelte';

	export let amount: number = 3;

	let transactions: Array<Transaction> = [];
	let maxItemPerTransaction: number = 6;

	onMount(() => {
		transactionsApi()
			.getCurrentAccountTransactions(0, amount, undefined, { withCredentials: true })
			.then((res) => {
				if (!(res.data.transactions instanceof Array)) return;
				let countedItems = 0;
				let newTransactions = [];
				for (let transaction of res.data.transactions) {
					let items: Array<TransactionItem> = [];
					for (let item of transaction.items) {
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
	});
</script>

<!-- Good looking dropdown for transaction -->
<div class="w-full">
	<div class="flex flex-col">
		<div class="flex flex-row justify-between">
			<div class="text-lg font-semibold">Transactions</div>
			<div class="text-lg font-semibold">Montant</div>
		</div>
		<div class="flex flex-col">
			{#each transactions as transaction}
				<div class="flex flex-row justify-between mt-5 border-4 border-white rounded-xl {transaction.state=="started" ? "animate-pulse bg-green-100":""} {transaction.state=="canceled" ? "bg-gray-200":""} {transaction.state=="finished" ? "bg-green-200":""}">
					<div class="p-5 h-full pr-4 w-full">
						{#each transaction.items as item}
							<div class="grid grid-cols-3 gap-2">
								<!-- One for each item.amount -->
								{#each Array(item.item_amount) as _}
									<img
										src={api() + item.picture_uri}
										alt="ca charge"
										class="w-10 h-10 rounded-2xl"
									/>
								{/each}
							</div>
						{/each}
					</div>
					<div class="border-r border-l border-gray-400" />
					<div class="p-5 pl-4 w-full text-lg self-center text-center text-red-600">
						-{formatPrice(transaction.total_cost)}
					</div>
				</div>
			{/each}
		</div>
	</div>
</div>

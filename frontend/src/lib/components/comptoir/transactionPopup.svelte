<script lang="ts">
	import type { Transaction } from '$lib/api';
	import { api } from '$lib/config/config';
	import { transactionsApi } from '$lib/requests/requests';
	import { formatPrice } from '$lib/utils';

	export let transaction: Transaction;
	export let close: () => void;

	let newTransaction: Transaction = transaction;

</script>

<!-- Popup overlay -->
<button
	id="overlay"
	class="absolute w-full h-full top-0 left-0 bg-black bg-opacity-50 flex justify-center items-center z-10 hover:cursor-default"
	on:click={() => {
		close();
	}}
/>

<div id="popup" class="absolute w-full h-full top-0 left-0 flex justify-center items-center">
	<!-- 
        We can modify the transaction here
        We can :
            - Check an item and mark it as completed
            - Check the transaction to mark it as completed
            - Cancel (item or transaction)
            - Undo a cancel (item or transaction)
            - Lower the amount of an item
            - Close the popup
    -->
	<div class="w-2/3 bg-white rounded-xl z-20 text-black">
		<div class="p-5 h-full pr-4 w-full">
			<div class="grid grid-cols-8 gap-2">
				{#each transaction.items as item, i}
					<!-- One for each item.amount -->
					<div class="flex flex-col justify-center text-center">
						{item.item_name?item.item_name:"Test"}
						<img src={api() + item.picture_uri} alt="ca charge" class="self-center w-10 h-10 rounded-2xl" />
						{newTransaction.items[i].item_amount}
						<div class="flex flex-row justify-center">
							<button
								class="bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded"
								on:click={() => {
									if (newTransaction.items[i].item_amount > 0)
										newTransaction.items[i].item_amount--;
								}}
							>
								-
							</button>
							{#if newTransaction.items[i].item_amount < item.item_amount}
							<button
								class="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded"
								on:click={() => {
									if (newTransaction.items[i].item_amount < item.item_amount)
										newTransaction.items[i].item_amount++;
								}}
							>
								+
							</button>
							{/if}
						</div>
					</div>
				{/each}
			</div>
		</div>
		<div class="border-r border-l border-gray-400" />
		<div class="p-5 pl-4 w-full text-lg self-center text-center">
			Prix: {formatPrice(transaction.total_cost)}
		</div>
	</div>
</div>

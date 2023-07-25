<script lang="ts">
	import Categories from '$lib/components/borne/categories.svelte';
	import Items from '$lib/components/borne/items.svelte';
	import { onMount, onDestroy } from 'svelte';
	import { formatPrice } from '$lib/utils';
	import { store } from '$lib/store/store';
	import { fly } from 'svelte/transition';
	import type { Account, Item, NewTransactionItem, TransactionItem } from '$lib/api';
	import Transactions from '$lib/components/borne/transactions.svelte';
	import { api } from '$lib/config/config';

	let account: Account | undefined = undefined;
	let unsub: () => void;

	type NewTransactionItemWithItem = NewTransactionItem & { item: Item };

	let order: NewTransactionItemWithItem[] = [];
	let orderPrice: number = 0;

	onMount(() => {
		unsub = store.subscribe((state) => {
			account = state.account;
		});
	});

	onDestroy(() => {
		unsub();
	});

	let currentCatgory: string = '';

	let changeCategory: (category: string) => void = (category: string) => {
		currentCatgory = '';
		setTimeout(() => {
			currentCatgory = category;
		}, 10);
	};

	let clickItem: (item: Item) => void = (item: Item) => {
		let newOrder = order;

		if (newOrder.find((i) => i.item_id == item.id)) {
			let found = newOrder.find((i) => i.item_id == item.id)!;
			if (found.amount >= found.item.buy_limit) {
				return;
			}
			found.amount++;
			order = newOrder;
			orderPrice += item.price;
			return;
		}

		let newTItem: NewTransactionItemWithItem = {
			item_id: item.id,
			amount: 1,
			item: item
		};
		newOrder.push(newTItem);
		order = newOrder;
		orderPrice += item.price;
	};

	function removeItem(item: NewTransactionItemWithItem) {
		return () => {
			let newOrder = order;

			if (newOrder.find((i) => i.item_id == item.item.id)) {
				newOrder.find((i) => i.item_id == item.item.id)!.amount--;

				if (newOrder.find((i) => i.item_id == item.item.id)!.amount == 0) {
					newOrder.splice(newOrder.indexOf(item), 1);
				}

				order = newOrder;
				orderPrice -= item.item.price;
				return;
			}
		};
	}

	let sidebar = true;
</script>

<div
	id="main"
	class="absolute w-screen h-screen top-0 left-0 overflow-y-hidden"
	style="background-color:#393E46"
>
	<div class="{sidebar ? 'w-4/5' : 'w-full'} h-full relative transition-all ease-in-out">
		<div class="p-4 flex justify-between" style="background-color:#222831">
			<Categories {changeCategory} />
			<button
				class="flex items-center space-x-2 px-4 py-2 rounded-lg bg-green-500 hover:bg-green-600 transition-colors duration-300 animate-pulse"
				on:click={() => {
					sidebar = !sidebar;
				}}
			>
				{#if sidebar}
					<iconify-icon class="text-white align-middle text-2xl" icon="akar-icons:chevron-right" />
				{:else}
					<iconify-icon class="text-white align-middle text-2xl" icon="akar-icons:chevron-left" />
				{/if}
			</button>
		</div>
		{#if currentCatgory != ''}
			<Items category={currentCatgory} click={clickItem} />
		{/if}
	</div>
	{#if sidebar}
		<div
			class="absolute top-0 right-0 w-1/5 h-screen"
			style="background-color:#222831"
			in:fly={{ x: 300, duration: 200 }}
			out:fly={{ x: 300, duration: 200 }}
		>
			<div class="p-4 flex justify-between h-1/6">
				<div
					class="flex flex-col gap-5 justify-center items-center w-full h-full overflow-x-auto overflow-y-hidden"
				>
					<!-- Commande en cours title -->
					<h1 class="text-white text-md md:text-md lg:text-2xl">Commande actuelle</h1>
					<!-- Subtitle with current balance -->
					<h2 class="text-white text-xs md:text-xs lg:text-xl">
						Disponible: {formatPrice(account?.balance ?? 0)}
					</h2>

					<!-- Spacer -->
				</div>
			</div>
			<hr class="w-full border-white" />

			<!-- Items in current commande with buttons for + and - with how much there is and the cost -->
			<div
				class="flex flex-col gap-5 justify-center items-center overflow-x-auto overflow-y-visible h-4/6 p-4"
			>
				<div class="grid grid-cols-2 gap-10">
					{#each order as item}
						<div class="flex flex-col justify-center gap-5 items-center w-full">
							<img class="w-16" src={api() + item.item.picture_uri} alt={item.item.name} />
							<div class="flex flex-row justify-center items-center">
								<button
									class="w-10 h-10 border-2 border-gray-300 rounded-full"
									on:click={removeItem(item)}
								>
									<iconify-icon class="text-white align-middle text-2xl" icon="akar-icons:minus" />
								</button>
								<span class="text-lg text-white mx-4">{item.amount}</span>
								<button
									class="w-10 h-10 border-2 border-gray-300 rounded-full"
									on:click={() => clickItem(item.item)}
								>
									<iconify-icon class="text-white align-middle text-2xl" icon="akar-icons:plus" />
								</button>
							</div>
							<span class="text-lg text-white">{formatPrice(item.item.price * item.amount)}</span>
						</div>
					{/each}
				</div>
			</div>
			<hr class="w-full border-white" />
			<div class="p-4 flex justify-between bottom-0 h-1/6">
				<div
					class="flex flex-col gap-1 justify-center items-center w-full h-full overflow-x-auto overflow-y-hidden"
				>
					<!-- Commande en cours title -->
					<h1 class="text-md md:text-md lg:text-2xl text-white">Total</h1>
					<!-- Subtitle with current balance -->
					<h2 class="text-xs md:text-xs lg:text-xl text-white">Coût: {formatPrice(orderPrice)}</h2>
					<h2 class="text-xs md:text-xs lg:text-xl text-white">
						Restant: {formatPrice((account?.balance ?? 0) - orderPrice)}
					</h2>

					<button class="w-full h-10 bg-green-500 rounded-lg text-white text-lg font-bold">
						Valider la commande
					</button>
				</div>
			</div>
		</div>
	{/if}
</div>

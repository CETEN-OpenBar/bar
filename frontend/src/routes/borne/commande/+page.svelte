<script lang="ts">
	import Categories from '$lib/components/borne/categories.svelte';
	import Items from '$lib/components/borne/items.svelte';
	import { onMount, onDestroy } from 'svelte';
	import { formatPrice } from '$lib/utils';
	import { store } from '$lib/store/store';
	import { fly } from 'svelte/transition';
	import type {
		Account,
		Item,
		NewTransaction,
		NewTransactionItem,
		TransactionItem
	} from '$lib/api';
	import Transactions from '$lib/components/borne/transactions.svelte';
	import { api } from '$lib/config/config';
	import Confirm from '$lib/components/borne/confirm.svelte';
	import { accountsApi, authApi, transactionsApi } from '$lib/requests/requests';
	import Pin from '$lib/components/borne/pin.svelte';
	import Error from '$lib/components/error.svelte';
	import Success from '$lib/components/success.svelte';
	import { goto } from '$app/navigation';

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
			if (found.amount >= (found.item.buy_limit ?? 9999)) {
				return;
			}
			found.amount++;
			order = newOrder;
			orderPrice += item.display_price ?? 999;
			return;
		}

		let newTItem: NewTransactionItemWithItem = {
			item_id: item.id,
			amount: 1,
			item: item
		};
		newOrder.push(newTItem);
		order = newOrder;
		orderPrice += item.display_price ?? 999;
	};

	function removeItem(item: NewTransactionItemWithItem, amount: number = 1) {
		return () => {
			let newOrder = order;
			let found = newOrder.find((i) => i.item_id == item.item.id)!;

			if (found) {
				found!.amount -= amount;

				if (found!.amount < 0) {
					amount += found!.amount;
				}

				if (found!.amount == 0) {
					newOrder.splice(newOrder.indexOf(item), 1);
				}

				order = newOrder;
				orderPrice -= amount * (item.item.display_price ?? 999);
				return;
			}
		};
	}

	function confirmOrder(response: Boolean) {
		confirm = false;
		if (!response) {
			return;
		}
		pin = true;
	}

	function finalizeTransaction(card_pin: string) {
		if (card_pin == '') {
			pin = false;
			error = "J'ai besoin de votre code pin pour valider la transaction";
			setTimeout(() => {
				error = '';
			}, 3000);
			return;
		}
		let transaction: NewTransaction = {
			items: order.map((item) => {
				return {
					item_id: item.item_id,
					amount: item.amount
				};
			}),
			card_pin: card_pin
		};

		transactionsApi()
			.postTransactions(transaction, { withCredentials: true })
			.then((res) => {
				success = 'Transaction effectuée avec succès';
				setTimeout(() => {
					success = '';
					authApi().logout({ withCredentials: true });
					goto('/borne');
				}, 3000);
				order = [];
				orderPrice = 0;
			})
			.catch((err) => {
				error = 'Erreur lors de la transaction';
				setTimeout(() => {
					error = '';
				}, 3000);
				pin = false;
			});
		confirm = false;
	}

	let error = '';
	let success = '';
	let pin = false;
	let confirm = false;
	let sidebar = true;
</script>

{#if confirm}
	<Confirm custom_text="Envoyer la commande ?" callback={confirmOrder} />
{/if}

{#if pin}
	<Pin callback={finalizeTransaction} />
{/if}

{#if error}
	<Error {error} />
{/if}

{#if success}
	<Success message={success} />
{/if}

<div
	id="main"
	class="absolute w-screen h-screen top-0 left-0 overflow-y-hidden"
	style="background-color:#393E46"
>
	<div class="{sidebar ? 'w-4/5' : 'w-full'} h-full relative transition-all ease-in-out">
		<div class="p-4 flex justify-between" style="background-color:#222831">
			<button
				class="flex items-center h-1/2 space-x-2 px-4 py-2 mr-2 rounded-lg bg-green-500 hover:bg-green-600 transition-colors duration-300"
				on:click={() => {
					goto('/borne/index');
				}}
			>
				<iconify-icon class="text-white align-middle text-2xl" icon="akar-icons:chevron-left" />
			</button>
			<Categories {changeCategory} />
			<button
				class="flex items-center space-x-2 px-4 py-2 ml-2 rounded-lg bg-green-500 hover:bg-green-600 transition-colors duration-300 animate-pulse"
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
			<div class="relative flex flex-col gap-5 justify-center items-center h-4/6 p-4">
				{#if order.length == 0}
					<h1 class="text-white text-md md:text-md lg:text-2xl">Aucun article</h1>
				{:else}
					<button
						class="w-10 h-10 rounded-full absolute top-2 bg-red-500"
						on:click={() => {
							order = [];
							orderPrice = 0;
						}}
					>
						<iconify-icon class="text-white align-middle text-2xl" icon="icomoon-free:bin" />
					</button>
				{/if}
				<div class="grid grid-cols-2 gap-10 overflow-x-auto overflow-y-scroll">
					{#each order as item}
						<div class="flex flex-col justify-center gap-5 items-center w-full">
							<button
								class="-mr-20 -my-6 w-6 h-6 rounded-full z-10"
								on:click={removeItem(item, item.amount)}
							>
								<iconify-icon class="text-white align-middle text-2xl" icon="ic:outline-cancel" />
							</button>
							<img
								draggable="false"
								class="w-16 h-16 object-contain"
								src={api() + item.item.picture_uri}
								alt={item.item.name}
							/>
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
							<span class="text-lg text-white"
								>{formatPrice((item.item.display_price ?? 999) * item.amount)}</span
							>
						</div>
					{/each}
				</div>
			</div>
			<hr class="w-full border-white" />
			<div class="p-1 flex justify-between bottom-0 h-1/6">
				<div
					class="flex flex-col gap-1 justify-center items-center w-full h-full overflow-x-auto overflow-y-hidden"
				>
					<h1 class="text-md md:text-md lg:text-2xl text-white">Total</h1>
					<h2 class="text-md text-white">Coût: {formatPrice(orderPrice)}</h2>
					<h2 class="text-md text-white">
						Reste: {formatPrice((account?.balance ?? 0) - orderPrice)}
					</h2>

					<button
						class="w-full h-16 bg-green-500 rounded-lg text-white text-lg font-bold"
						on:click={() => (confirm = true)}
					>
						Valider la commande
					</button>
				</div>
			</div>
		</div>
	{/if}
</div>

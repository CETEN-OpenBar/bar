<script lang="ts">
	import Items from '$lib/components/client/items.svelte';
	import { onMount, onDestroy } from 'svelte';
	import { formatPrice } from '$lib/utils';
	import { store } from '$lib/store/store';
	import { slide } from 'svelte/transition';
	import type {
		Account,
		Item,
		MenuCategory,
		NewTransaction,
		NewTransactionItem,
		Category
	} from '$lib/api';
	import { api } from '$lib/config/config';
	import Confirm from '$lib/components/client/confirm.svelte';
	import { transactionsApi, categoriesApi } from '$lib/requests/requests';
	import Pin from '$lib/components/client/pin.svelte';
	import Error from '$lib/components/error.svelte';
	import Success from '$lib/components/success.svelte';
	import { goto } from '$app/navigation';
	import Stars from '$lib/components/random/stars.svelte';
	import Price from '$lib/components/random/price.svelte';
	import Hamburger from '$lib/components/client/hamburger.svelte';

	let account: Account | undefined = undefined;
	let unsub: () => void;

	type NewTransactionItemWithItem = NewTransactionItem & {
		category: string;
		item: Item;
		pickedItems: NewTransactionItem[] | undefined;
	};

	let order: NewTransactionItemWithItem[] = [];
	let orderPrice: number = 0;
	let orderForLater: boolean = true;

	onMount(() => {
		unsub = store.subscribe((state) => {
			account = state.account;
		});

		/* Load categories */
		categoriesApi()
			.getCategories(undefined, { withCredentials: true })
			.then((res) => {
				categories = res.data??[];
				changeCategory(categories[0].id, categories[0].name);
			});

	});

	onDestroy(() => {
		unsub();
	});

	let categories: Category[] = [];
	let currentCatgoryId: string = '';
	let currentCategoryName: string = '';

	let changeCategory: (categoryId: string, categoryName: string) => void = (categoryId: string, categoryName: string) => {
		currentCatgoryId = ''
		currentCategoryName = '';
		setTimeout(() => {
			currentCategoryName = categoryName;
			currentCatgoryId = categoryId;
		}, 10);
	};

	type MenuPopup = {
		categories: MenuCategory[] | undefined;
		pickedItems: NewTransactionItemWithItem[];
		tItem: NewTransactionItemWithItem;
		step: number;
	};
	let menuPicks: MenuPopup | undefined = undefined;

	let clickItemMenu: (item: Item) => void = (item: Item) => {
		let newPicks = menuPicks?.pickedItems ?? [];

		if (item.amount_left == 0) {
			return;
		}

		if (newPicks.find((i) => i.item_id == item.id)) {
			let found = newPicks.find((i) => i.item_id == item.id)!;
			if (found.amount >= (found.item.buy_limit ?? 9999)) {
				return;
			}
			if (found.amount >= (found.item.amount_left ?? 9999)) {
				return;
			}
			found.amount++;
		} else {
			let newTItem: NewTransactionItemWithItem = {
				category: currentCatgoryId,
				item_id: item.id,
				amount: 1,
				item: item,
				pickedItems: undefined
			};

			newPicks.push(newTItem);
		}

		if (menuPicks) {
			let amt = 0;
			for (let i = 0; i < menuPicks.pickedItems.length; i++) {
				if (menuPicks.pickedItems[i].category == currentCatgoryId) {
					amt += menuPicks.pickedItems[i].amount;
				}
			}

			menuPicks.pickedItems = newPicks;

			if (amt >= (menuPicks.categories ?? [])[menuPicks.step].amount) {
				menuPicks.step++;
				if (menuPicks.step >= (menuPicks.categories ?? []).length) {
					menuPicks.step = 0;

					menuPicks.tItem.pickedItems = menuPicks.pickedItems;

					let newOrder = order;
					newOrder.push(menuPicks.tItem);
					order = newOrder;
					orderPrice += menuPicks.tItem.item.display_price ?? 999;
					menuPicks = undefined;
					return;
				}
				changeCategory((menuPicks.tItem.item.menu_categories ?? [])[menuPicks.step].id, (menuPicks.tItem.item.menu_categories ?? [])[menuPicks.step].name);
			}
		}
	};

	let clickItem: (item: Item) => void = (item: Item) => {
		let newOrder = order;

		if (newOrder.find((i) => i.item_id == item.id)) {
			let found = newOrder.find((i) => i.item_id == item.id)!;
			found.item = item;
			if (found.amount >= (found.item.buy_limit ?? 9999)) {
				return;
			}
			if (found.amount >= (found.item.amount_left ?? 9999)) {
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
			item: item,
			pickedItems: undefined,
			category: ''
		};

		if (item.is_menu && item.menu_categories) {
			menuPicks = {
				categories: item.menu_categories,
				pickedItems: [],
				tItem: newTItem,
				step: 0
			};
			changeCategory(item.menu_categories[0].id, item.menu_categories[0].name);
		} else {
			newOrder.push(newTItem);
			order = newOrder;
			orderPrice += item.display_price ?? 999;
		}
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
					amount: item.amount,
					picked_categories_items: item.pickedItems?.map((i) => {
						return {
							item_id: i.item_id,
							amount: i.amount
						};
					})
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
				}, 3000);
				order = [];
				orderPrice = 0;
				pin = false;
				order_menu = false;
				categories_menu = false;
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

	let categories_menu = false;
	let order_menu = false;

	function toggleCategoriesMenu() {
		categories_menu = !categories_menu;
	}

	function toggleOrderMenu() {
		order_menu = !order_menu;
	}

</script>

<style>
	/* Prevents reload when scrolling down on mobile */
	:global(body) {
		overflow-y: hidden;
	}
</style>

{#if confirm}
	<Confirm custom_text="Envoyer la commande ?" callback={confirmOrder}/>
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

<!-- Categories menu -->

<aside class="flex flex-col w-4/5 fixed h-full text-white transition-transform ease-in-out duration-150 bg-[#393E46] {categories_menu ? 'z-20' : '-translate-x-full shadow-2xl'}">
	<button class="w-full flex flex-row items-center bg-[#222831] p-3" on:click={() => {categories_menu = false;}}>
		<h1 class="flex-grow font-bold text-lg text-center h-fit">Catégories</h1>
		<Hamburger activated={categories_menu} toggle={() => {categories_menu = false}} />
	</button>
	<menu class="flex-grow w-full overflow-y-scroll">
		{#each categories as category}
			<button
				class="flex flex-row w-full items-center space-x-5 p-1 {category.id == currentCatgoryId ? 'bg-gray-600' : ''}"
				on:click={() => {changeCategory(category.id, category.name); categories_menu = false}}
			>
				<img
					draggable="false"
					class="h-20 w-20 object-contain"
					src={api() + category.picture_uri}
					alt={category.name}
				/>
				<span class="text-lg">{category.name}</span>
			</button>
		{/each}
	</menu>
</aside>

<!-- Order menu -->
<div 
	class="bg-[#222831] fixed bottom-0 w-full text-white flex flex-col items-center p-2 pb-4 max-h-[90vh] {order_menu ? 'z-20 space-y-5' : 'space-y-3'}"
>
	<!-- Handle -->
	<button class="w-[30px] h-[6px] rounded-md bg-slate-500" on:click={toggleOrderMenu}></button>
	
	<!-- Overlay to open -->
	{#if !order_menu}
		<button class="absolute w-full h-full" on:click={toggleOrderMenu}></button>
	{/if}

	{#if order_menu}
	<div transition:slide class="flex flex-col items-center space-y-5 w-full">
		<!-- Current order-->
		<h1 class="font-semibold text-xl">Commande actuelle</h1>
		<!-- Items -->
		{#if order.length == 0}
			<div class="h-16 flex items-center">Aucun articles</div>
		{:else}
			<ul class="overflow-y-scroll flex flex-col w-full max-w-full max-h-[50vh]">
				{#each order as item}
					<li class="max-h-16 flex flex-row items-center w-full max-w-full">
						<img
							draggable="false"
							class="w-16 h-16 object-contain mr-3"
							src={api() + item.item.picture_uri}
							alt={item.item.name}
						/>
						<div class="flex-grow">{item.item.name}</div>
						<div class="min-w-fit ml-2">
							{formatPrice((item.item.display_price ?? 999) * item.amount)}</div
						>
						<div class="min-w-fit flex flex-row items-center">
							<button
								class="w-7 h-7 ml-2 border-2 border-gray-300 rounded-full flex flex-row items-center justify-center"
								on:click={removeItem(item)}
							>
								<iconify-icon class="text-white" icon="akar-icons:minus" />
							</button>
							<span class=" text-white mx-2">{item.amount}</span>
							<button
								class="w-7 h-7 border-2 border-gray-300 rounded-full flex flex-row items-center justify-center"
								on:click={() => clickItem(item.item)}
							>
								<iconify-icon class="text-white" icon="akar-icons:plus" />
							</button>
						</div>
					</li>
				{/each}
			</ul>
		{/if}
		<div class="bg-gray-500 h-[1px] w-11/12"></div>
		<label class="text-lg flex flew-row space-x-3 items-center">
			<input 
				type="checkbox"
				bind:checked={orderForLater}
				class="h-5 w-5"
			/>
			<span>Je commande pour plus tard</span>
		</label>
	</div>
	{/if}

	<div class="flex flex-col items-center space-y-3">
		<div class="justify-self-start grow font-semibold">
			Total : 
			<Price amount={orderPrice} class="inline-block"/>
		</div>
		<div>
			Reste :
			<Price 
				amount={(account?.balance ?? 0) - orderPrice + (account?.points ?? 0)}
				class="inline-block px-2"
			/> 
			{#if (account?.points ?? 0) > 0}
				<Stars 
					stars={(account?.points ?? 0) - orderPrice}
					icon_size={5}
					class="inline-block"
				/>
			{/if}
		</div>
	</div>

	{#if order_menu}
		<button
			class="bg-green-500 rounded-md font-semibold w-11/12 p-3 block disabled:bg-gray-400"
			transition:slide
			on:click={() => {confirm = true}}
			disabled={order.length == 0}
		>
			Valider la commande
		</button>
	{/if}

</div>

<div
	id="main"
	class="w-screen h-screen top-0 left-0 overflow-y-hidden bg-[#393E46] text-white flex flex-col {categories_menu || order_menu ? 'blur-sm' : ''}"
>

	<!-- Menu close overlay -->
	{#if categories_menu || order_menu}
		<button
			class="fixed h-full w-full bg-black bg-opacity-60 z-10"
			on:click={() => {categories_menu = false; order_menu = false}}
		></button>
	{/if}

	<!-- Header -->
	<div class="bg-[#222831] flex flex-row w-full p-3 items-center">
		<Hamburger activated={categories_menu} toggle={toggleCategoriesMenu} />
		<button class="flex-grow font-semibold ml-3 text-left" on:click={toggleCategoriesMenu}>{currentCategoryName}</button>

		<div class="flex flex-col items-end mr-3 text-xs space-y-2">
			<Price amount={account?.balance ?? 0} />
			{#if (account?.points ?? 0) > 0}
				<Stars stars={account?.points ?? 0} icon_size={5}/>
			{/if}
		</div>

		<button 
		  on:click={() => {
		    goto('/client/index');
		  }}
		>
			<img src={account?.google_picture} alt="Profile" class="w-12 h-12 rounded-full border-2 border-gray-200"/>
		</button>
	</div>

	<!-- Items -->
	<div class="flex flex-col flex-grow items-center overflow-y-scroll">
		{#if currentCatgoryId != ''}
			<Items category={currentCatgoryId} click={clickItem} limit={24} />
		{/if}
	</div>


	<!-- Order menu duplicate to put the correct margin below the items-->
	<div class="bg-[#222831] flex flex-col items-center p-2 pb-4 space-y-3">
		<!-- Handle -->
		<div class="w-[30px] h-[6px] rounded-md bg-slate-500"></div>
		<div class="justify-self-start grow font-semibold">
			Total : 
			<Price amount={orderPrice} class="inline-block"/>
		</div>
		<div>
			Reste :
			<Price 
				amount={(account?.balance ?? 0) - orderPrice + (account?.points ?? 0)}
				class="inline-block px-2"
			/> 
			{#if (account?.points ?? 0) > 0}
				<Stars 
					stars={(account?.points ?? 0) - orderPrice}
					icon_size={5}
					class="inline-block"
				/>
			{/if}
		</div>
	</div>
</div>

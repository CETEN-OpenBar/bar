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

	let innerWidth = 0;
	$: screenXl = innerWidth >= 1280

</script>

<svelte:window bind:innerWidth />

<style>
	/* Prevents reload when scrolling down on mobile */
	:global(body) {
		overflow-y: hidden;
	}

	/* Prevents double-tap zooming */
	* {
		touch-action: manipulation;
	}
</style>

<svelte:head>
    <title>OpenBar - Commander</title>
</svelte:head>


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

<aside class="flex flex-col w-4/5 xl:w-1/3 fixed h-full text-white transition-transform ease-in-out duration-150 bg-[#393E46] {categories_menu ? 'z-20' : '-translate-x-full shadow-2xl'}">
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

<!-- Main content -->
<div
	id="main"
	class="w-screen h-screen top-0 left-0 overflow-y-hidden bg-[#393E46] text-white flex flex-col"
>

	<!-- Menu close overlay -->
	{#if categories_menu || order_menu}
		<button
			class="fixed h-full w-full bg-black bg-opacity-60 z-10"
			on:click={() => {categories_menu = false; order_menu = false}}
		></button>
	{/if}

	<!-- Header -->
	<div class="bg-[#222831] flex flex-row w-full p-3 items-center {categories_menu || order_menu ? 'blur-sm' : ''}">
		<Hamburger activated={categories_menu} toggle={toggleCategoriesMenu} />
		<button class="flex-grow font-semibold ml-3 text-left" on:click={toggleCategoriesMenu}>{currentCategoryName}</button>

		<div class="flex flex-col items-end mr-3 text-xs space-y-2">
			<Price amount={account?.balance ?? 0} />
			{#if (account?.points ?? 0) > 0}
				<Stars stars={account?.points ?? 0} icon_size={5} class="text-right"/>
			{/if}
		</div>

		<button 
		  on:click={() => {
		    goto('/client/index');
		  }}
		>
			<img src={account?.google_picture ?? 'https://www.gravatar.com/avatar/' + account?.email_address + '?d=mp'} alt="Profile" class="w-12 h-12 rounded-full border-2 border-gray-200"/>
		</button>
	</div>

	<div class="flex-grow flex flex-col xl:flex-row overflow-y-scroll">
		<!-- Items -->
		<div class="flex flex-col flex-grow items-center overflow-y-scroll {categories_menu || order_menu ? 'blur-sm' : ''}">
			{#if currentCatgoryId != ''}
				<Items category={currentCatgoryId} click={clickItem} limit={24} />
			{/if}
		</div>
	
		<!-- Order menu -->
		<div 
			class="bg-[#222831] bottom-0 w-full text-white flex flex-col items-center p-2 pb-4 max-h-[90vh] 
			{order_menu ? 'z-20 space-y-5 fixed' : 'relative space-y-3'}
			{categories_menu ? 'blur-sm' : ''}
			{categories_menu && !order_menu ? 'z-[1]' : ''}
			xl:max-w-md xl:h-full xl:max-h-full"
		>
			<!-- Handle -->
			<button 
				class="{order_menu ? 'w-[60px] h-6 p-1' : 'w-[30px] h-[6px]'} rounded-md bg-slate-500 
					flex flex-col items-center transition-[width,height] duration-200 xl:hidden" 
				on:click={toggleOrderMenu}>
				{#if order_menu}
					<svg xmlns="http://www.w3.org/2000/svg"
						class="h-5 w-5"
						viewBox="0 0 512 512"
					>
						<!--!Font Awesome Free 6.6.0 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license/free Copyright 2024 Fonticons, Inc.-->
						<path d="M233.4 406.6c12.5 12.5 32.8 12.5 45.3 0l192-192c12.5-12.5 12.5-32.8 0-45.3s-32.8-12.5-45.3 0L256 338.7 86.6 169.4c-12.5-12.5-32.8-12.5-45.3 0s-12.5 32.8 0 45.3l192 192z"/>
					</svg>
				{/if}
			</button>
	
			<!-- Overlay to open -->
			{#if !order_menu}
				<button class="absolute w-full h-full xl:hidden" on:click={toggleOrderMenu}></button>
			{/if}
	
			{#if order_menu || screenXl}
			<div transition:slide class="flex flex-col items-center space-y-5 w-full xl:grow">
				<!-- Current order-->
				<h1 class="font-semibold text-xl">Commande actuelle</h1>
				<!-- Items -->
				{#if order.length == 0}
					<div class="h-16 flex items-center xl:grow">Aucun articles</div>
				{:else}
					<ul class="overflow-y-scroll flex flex-col w-full max-w-full max-h-[50vh] xl:max-h-[60vh] xl:grow">
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
					<div class="hidden xl:block xl:grow"></div>
				{/if}
				<div class="bg-gray-500 h-[1px] w-11/12"></div>
			</div>
			{/if}
	
			<div class="flex flex-row px-5 w-full items-center">
				<div class="flex flex-col items-center space-y-3 grow">
					<div class="justify-self-start grow font-semibold">
						Total : 
						<Price amount={orderPrice} class="inline-block"/>
					</div>
					<div class="grid grid-cols-3 w-full">
						<div class="{order_menu ? 'col-span-3' : 'col-span-2 xl:col-span-3'} flex flex-col items-center transition-[width] duration-1000 ease-in-out">
							<div class="flex flex-row items-center space-x-2 {order_menu ? 'col-span-3' : 'col-span-2'}">
								<div>
									Reste :
								</div>
								<div class="inline-flex flex-col items-center">
									{#if orderPrice < (account?.points ?? 0)}
										<Price amount={account?.balance ?? 0} class="inline-block px-2"/>
										{#if (account?.points ?? 0) > 0}
											<Stars stars={(account?.points ?? 0) - orderPrice} icon_size={5} class="inline-block"/>
										{/if}
									{:else}
										<Price amount={(account?.balance ?? 0) - orderPrice + (account?.points ?? 0)} class="inline-block px-2" />
									{/if}
								</div>
							</div>
						</div>

						<!-- Cart logo -->
						{#if !order_menu}
						<div class="flex flex-col items-center col-span-1 xl:hidden">
							<button 
								class="rounded-md bg-green-500 p-3 flex flex-row items-center space-x-2 h-fit relative"
								on:click={toggleOrderMenu}
							>
								<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 512 512"><!--!Font Awesome Free 6.6.0 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license/free Copyright 2024 Fonticons, Inc.--><path d="M233.4 105.4c12.5-12.5 32.8-12.5 45.3 0l192 192c12.5 12.5 12.5 32.8 0 45.3s-32.8 12.5-45.3 0L256 173.3 86.6 342.6c-12.5 12.5-32.8 12.5-45.3 0s-12.5-32.8 0-45.3l192-192z"/></svg>
								<svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8" viewBox="0 0 576 512">
									<!--!Font Awesome Free 6.6.0 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license/free Copyright 2024 Fonticons, Inc.-->
									<path d="M0 24C0 10.7 10.7 0 24 0L69.5 0c22 0 41.5 12.8 50.6 32l411 0c26.3 0 45.5 25 38.6 50.4l-41 152.3c-8.5 31.4-37 53.3-69.5 53.3l-288.5 0 5.4 28.5c2.2 11.3 12.1 19.5 23.6 19.5L488 336c13.3 0 24 10.7 24 24s-10.7 24-24 24l-288.3 0c-34.6 0-64.3-24.6-70.7-58.5L77.4 54.5c-.7-3.8-4-6.5-7.9-6.5L24 48C10.7 48 0 37.3 0 24zM128 464a48 48 0 1 1 96 0 48 48 0 1 1 -96 0zm336-48a48 48 0 1 1 0 96 48 48 0 1 1 0-96z"/>
								</svg>
								{#if order.length > 0}
									<div class="absolute top-1 right-1 bg-red-600 rounded-full h-5 w-5 text-sm">
										{order.reduce((total, item) => total + item.amount, 0)}
									</div>
								{/if}
							</button>
						</div>
						{/if}
					</div>
				</div>
			</div>
	
			{#if order_menu || screenXl}
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
	</div>

</div>


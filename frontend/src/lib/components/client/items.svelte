<script lang="ts">
	import type { Item, ItemState, MenuCategory, MenuItem } from '$lib/api';
	import { api } from '$lib/config/config';
	import { itemsApi } from '$lib/requests/requests';
	import { formatPrice } from '$lib/utils';
	import { dragscroll } from '@svelte-put/dragscroll';
	import { onMount } from 'svelte';
	import { fade, fly } from 'svelte/transition';

	export let state: ItemState = 'buyable';
	export let category: string = '';

	export let click: (item: Item) => void;

	// Update item when clicking on it
	function clickWrapper(item: Item) {

		if (item.amount_left == 0) {
			return;
		}

		reloadItems(false);
		// Update current item
		for (let i = 0; i < items.length; i++) {
			if (items[i].id === item.id) {
				item = items[i];
				break;
			}
		}
		click(item);
	}

	let items: Item[] = [];

	let page: number = 0;
	let maxPage: number = 0;
	let nextPage = () => {
		if (page < maxPage) {
			page++;
			reloadItems();
		}
	};
	let prevPage = () => {
		if (page > 0) {
			page--;
			reloadItems();
		}
	};
	export let limit: number = 12;

	type MenuPopup = {
		items: MenuItem[] | undefined;
		categories: MenuCategory[] | undefined;
	};
	let menuPopup: MenuPopup | undefined;

	onMount(() => {
		reloadItems();
	});

	function reloadItems(anim: boolean = true) {
		itemsApi()
			.getCategoryItems(category, page, limit, state, { withCredentials: true })
			.then((res) => {
				maxPage = res.data.max_page ?? 0;
				page = res.data.page ?? 0;
				limit = res.data.limit ?? 15;

				if (anim) {
					let newItems = res.data.items ?? [];
					items = [];
					setTimeout(() => {
						items = newItems;
					}, 10);
				} else {
					items = res.data.items ?? [];
				}
			})
			.catch((err) => {
				console.log(err);
			});
	}

	let direction = 1;
</script>

{#if menuPopup}
	<!-- Show overlay -->
	<button
		class="absolute w-screen h-screen top-0 left-0 bg-black bg-opacity-50 flex flex-col items-center justify-center"
		in:fade={{ duration: 200 }}
		out:fade={{ duration: 200 }}
		style="z-index: 100;"
		on:click={() => {
			menuPopup = undefined;
		}}
	/>

	<div class="absolute w-full h-full top-0 left-0 flex justify-center items-center">
		<div
			class="flex flex-col items-center bg-white rounded-lg shadow-lg p-4"
			style="z-index: 101;"
			in:fade={{ duration: 200 }}
		>
			<div class="flex flex-col w-full justify-center">
				<div class="text-xl font-bold self-center">Contenu :</div>
				<div class="grid grid-cols-4 justify-between items-center w-full px-4 mt-4">
					{#each menuPopup.categories ?? [] as cat}
						<div class="flex flex-col justify-center">
							<img
								draggable="false"
								class="w-32 h-32 object-contain self-center bg-gray-200"
								src={api() + cat.picture_uri}
								alt={cat.name}
							/>
							<span class="w-full text-lg font-bold text-center">{cat.amount} {cat.name}</span>
						</div>
					{/each}
					{#each menuPopup.items ?? [] as item}
						<div class="flex flex-col justify-center">
							<img
								draggable="false"
								class="w-32 h-32 object-contain bg-gray-200"
								src={api() + item.picture_uri}
								alt={item.name}
							/>
							<span class="w-full text-lg font-bold text-center">{item.amount} {item.name}</span>
						</div>
					{/each}
				</div>
			</div>
		</div>
	</div>
{/if}

<!-- horizontal & overflows -->
{#if items.length === 0}
	<div class="flex flex-col flex-grow justify-items-center">
		<span class="text-3xl text-white" in:fade={{ duration: 200, delay: 100 }}>Aucun article</span>
	</div>
{:else}
	<div
		class="flex flex-wrap flex-grow max-w-full md:p-12 p-4"
		in:fly={{ x: -direction * 300, duration: 500 }}
		out:fly={{ x: direction * 300, duration: 500 }}
		use:dragscroll
	>
		{#each items as item}
			<!-- image wil be in a button box -->
			<div
				class="flex-1 flex flex-col basis-[min-content] my-1 min-w-40 max-w-40 min-h-50 h-50 max-h-50
				items-center rounded-lg text-white transition-colors duration-300 overflow-x-clip"
			>
				<!-- add info svg on the top right -->
				{#if item.is_menu}
					<button
						class="relative top-0 right-0 w-10 h-10"
						on:click={() => {
							menuPopup = {
								items: item.menu_items,
								categories: item.menu_categories
							};
						}}
					>
						<iconify-icon class="text-white align-middle text-2xl" icon="akar-icons:info" />
					</button>
				{/if}
				<button
					class="relative"
					on:click={() => {
						// check we are not clicking on the info button
						clickWrapper(item);
					}}
				>
					<img
						draggable="false"
						class="w-full h-28 md:h-32 object-contain"
						src={api() + item.picture_uri}
						alt={item.name}
					/>
					<div class="flex flex-col">
						<span class="text-lg font-bold w-40 m-[0_auto]">{item.name}</span>
						<span class="text-sm">Prix: {formatPrice(item.display_price ?? 999)}</span>
					</div>
					{#if item.amount_left <= 0}
						<!-- Stock épuisé icon -->
						<img
							class="absolute top-[50%] left-[50%] -translate-x-[50%] -translate-y-[70%] w-20 h-20 md:h-24 md:w-24 drop-shadow-2xl"
							alt="oof"
							src="/epuise.webp"
						/>
					{/if}
				</button>
			</div>
		{/each}
	</div>
{/if}

<!-- Navigation -->
<div class="flex flex-col justify-center">
	<div class="text-3xl text-white text-center">
		{page}/{maxPage}
	</div>
	<div class="flex flex-row gap-4 justify-center items-center w-full h-16">
		<button
			class="w-10 h-10 border-2 border-gray-300 rounded-full"
			on:click={() => {
				prevPage();
			}}
		>
			<iconify-icon class="text-white align-middle text-2xl" icon="akar-icons:chevron-left" />
		</button>
		<button
			class="w-10 h-10 text-center border-2 border-gray-300 rounded-full"
			on:click={() => {
				nextPage();
			}}
		>
			<iconify-icon class="text-white align-middle text-2xl" icon="akar-icons:chevron-right" />
		</button>
	</div>
</div>

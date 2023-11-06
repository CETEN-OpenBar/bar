<script lang="ts">
	import type { Item, ItemState, MenuItem } from '$lib/api';
	import { api } from '$lib/config/config';
	import { itemsApi } from '$lib/requests/requests';
	import { formatPrice } from '$lib/utils';
	import { onMount } from 'svelte';
	import { fade, fly } from 'svelte/transition';

	export let state: ItemState = 'buyable';
	export let category: string = '';
	export let click: (item: Item) => void;

	let items: Item[] = [];

	let page: number = 0;
	let maxPage: number = 0;
	let limit: number = 12;

	let menuPopup: MenuItem[] | undefined = undefined;

	onMount(() => {
		loadItems();
	});

	function loadItems() {
		itemsApi()
			.getCategoryItems(category, page, limit, state, { withCredentials: true })
			.then((res) => {
				maxPage = res.data.max_page ?? 0;
				page = res.data.page ?? 0;
				limit = res.data.limit ?? 15;

				let newItems = res.data.items ?? [];
				items = [];
				setTimeout(() => {
					items = newItems;
				}, 1);
			})
			.catch((err) => {
				console.log(err);
			});
	}

	let direction = 1;

	function nextPage() {
		if (page < maxPage) {
			page++;
			direction = -1;
			loadItems();
		}
	}

	function prevPage() {
		if (page > 1) {
			page--;
			direction = 1;
			loadItems();
		}
	}
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
					{#each menuPopup as item}
						<div class="flex flex-col justify-center">
							<img
								draggable="false"
								class="w-32 h-32 object-contain"
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
	<div class="col-span-7 flex flex-col items-center justify-center">
		<span class="text-3xl text-white">Aucun article</span>
	</div>
{:else}
	<div
		class="grid grid-cols-4 gap-8 w-full p-16"
		in:fly={{ x: -direction * 300, duration: 500 }}
		out:fly={{ x: direction * 300, duration: 500 }}
	>
		{#each items as item}
			{#if item.amount_left > 0}
				<!-- image wil be in a button box -->
				<button
					class="w-50 h-50 flex-shrink-0 flex flex-col items-center justify-between rounded-lg text-white transition-colors duration-300"
				>
					<!-- add info svg on the top right -->
					{#if item.is_menu}
						<button
							class="relative top-0 right-0 w-10 h-10"
							on:click={() => {
								menuPopup = item.items;
							}}
						>
							<iconify-icon class="text-white align-middle text-2xl" icon="akar-icons:info" />
						</button>
					{/if}
					<button
						on:click={() => {
							// check we are not clicking on the info button
							click(item);
						}}
					>
						<img
							draggable="false"
							class="w-32 h-32 object-contain"
							src={api() + item.picture_uri}
							alt={item.name}
						/>
						<div class="flex flex-col">
							<span class="text-lg font-bold">{item.name}</span>
							<span class="text-sm">Prix: {formatPrice(item.display_price ?? 999)}</span>
						</div>
					</button>
				</button>
			{/if}
		{/each}
	</div>
{/if}

<!-- Navigation -->
<div class="absolute bottom-5 left-[50%] -translate-x-[50%] flex flex-col justify-center">
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

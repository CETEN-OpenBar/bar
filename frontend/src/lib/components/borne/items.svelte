<script lang="ts">
	import type { Item, ItemState } from '$lib/api';
	import { api } from '$lib/config/config';
	import { itemsApi } from '$lib/requests/requests';
	import { formatPrice } from '$lib/utils';
	import { onMount } from 'svelte';
	import { fly } from 'svelte/transition';

	export let state: ItemState = 'buyable';
	export let category: string = '';
	export let click: (item: Item) => void;

	let items: Item[] = [];

	let page: number = 0;
	let maxPage: number = 0;
	let limit: number = 12;

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
					on:click={() => {
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

<script lang="ts">
	import type { Item, ItemState } from '$lib/api';
	import { api } from '$lib/config/config';
	import { itemsApi } from '$lib/requests/requests';
	import { onMount } from 'svelte';
	import { fly } from 'svelte/transition';

	export let state: ItemState = 'buyable';
	export let category: string = '';
	export let click: (item: string) => void;

	let items: Item[] = [];

	let page: number = 0;
	let maxPage: number = 0;
	let limit: number = 15;

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
	<div class="grid grid-cols-5 gap-3 w-full p-16" in:fly={{ x: -direction*300, duration: 500 }} out:fly={{ x: direction*300, duration: 500 }}>
		{#each items as item}
			<button
				class="w-32 flex-shrink-0 flex flex-col items-center justify-center rounded-lg text-white transition-colors duration-300"
				on:click={() => {
					click(item.id);
				}}
			>
				<img class="w-full" src={api() + item.picture_uri} alt={item.name} />
				<span class="text-lg font-bold">{item.name}</span>
			</button>
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

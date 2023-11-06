<script lang="ts">
	import type { Category } from '$lib/api';
	import { api } from '$lib/config/config';
	import { categoriesApi } from '$lib/requests/requests';
	import { onMount } from 'svelte';
	import { dragscroll } from '@svelte-put/dragscroll';

	export let changeCategory: (category: string) => void;

	let categories: Category[] = [];

	onMount(() => {
		categoriesApi()
			.getCategories({ withCredentials: true })
			.then((res) => {
				categories = res.data;
				changeCategory(categories[0].id);
			});
	});
</script>

<!-- horizontal & overflows -->
<div class="flex flex-row gap-5 items-center w-full h-full overflow-x-auto" use:dragscroll>
	{#each categories as category}
		<button
			class="w-32 flex-shrink-0 flex flex-col items-center justify-center m-2 rounded-lg text-white transition-colors duration-300"
			on:click={() => {
				changeCategory(category.id);
			}}
		>
			<img
				draggable="false"
				class="w-full"
				src={api() + category.picture_uri}
				alt={category.name}
			/>
			<span class="text-lg font-bold">{category.name}</span>
		</button>
	{/each}
</div>

<script lang="ts">
	import Categories from '$lib/components/borne/categories.svelte';
	import Items from '$lib/components/borne/items.svelte';
	import { fly } from 'svelte/transition';

	let currentCatgory: string = '';

	let changeCategory: (category: string) => void = (category: string) => {
		currentCatgory = '';
		setTimeout(() => {
			currentCatgory = category;
		}, 10);
	};

	let clickItem: (item: string) => void = (item: string) => {
		console.log(item);
	};

	let sidebar = true;
</script>

<div id="main" class="absolute w-screen h-screen top-0 left-0" style="background-color:#393E46">
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
			class="absolute top-0 right-0 w-1/5 h-full background-color:#393E4"
			in:fly={{ x: 300, duration: 500 }}
			out:fly={{ x: -300, duration: 500 }}
		>
		cc
		</div>
	{/if}
</div>

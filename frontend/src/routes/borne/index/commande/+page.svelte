<script lang="ts">
	import Categories from '$lib/components/borne/categories.svelte';
	import Items from '$lib/components/borne/items.svelte';

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
</script>

<div
	id="main"
	class="absolute w-screen h-screen top-0 left-0 bg-blue-50"
	style="background-color:#393E46"
>
	<div class="w-4/5 p-4 flex justify-between" style="background-color:#222831">
		<Categories {changeCategory} />
	</div>
	<div class="w-4/5 h-full">
		{#if currentCatgory != ''}
			<Items category={currentCatgory} click={clickItem} />
		{/if}
	</div>
</div>

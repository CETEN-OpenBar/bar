<script lang="ts">
	import type { Category } from '$lib/api';
	import { categoriesApi } from '$lib/requests/requests';
	import { onMount } from 'svelte';

	interface CategoryWithOldPosition extends Category {
		oldPosition: number;
	}
	let categories: CategoryWithOldPosition[] = [];

	onMount(async () => {
		let resp = await categoriesApi().getCategories(undefined, { withCredentials: true });
		if (Array.isArray(resp.data)) {
			categories = resp.data.map((category, index) => {
				return {
					...category,
					oldPosition: index
				};
			});
		}
	});

	let mouseYCoordinate: number | null = null; // pointer y coordinate within client
	let distanceTopGrabbedVsPointer: number = 0;

	let draggingItem: Category | null = null;
	let draggingItemIndex: number | null = null;

	let hoveredItemIndex: number | null = null;

	$: {
		// prevents the ghost flickering at the top
		if (mouseYCoordinate == null || mouseYCoordinate == 0) {
			// showGhost = false;
		}
	}

	$: {
		if (
			draggingItemIndex != null &&
			hoveredItemIndex != null &&
			draggingItemIndex != hoveredItemIndex
		) {
			// swap items
			[categories[draggingItemIndex], categories[hoveredItemIndex]] = [
				categories[hoveredItemIndex],
				categories[draggingItemIndex]
			];

			// balance
			draggingItemIndex = hoveredItemIndex;

			// Save all categories positions
			categories.forEach(async (category, index) => {
				if (category.position != index) {
					await categoriesApi().patchCategory(
						category.id,
						{ position: index },
						{ withCredentials: true }
					);
					category.position = index;
				}
			});
		}
	}

	let container = null;
</script>

<h1 class="w-full text-3xl font-bold text-center mt-16">Repositionner les catégories</h1>

<div
	class="flex flex-row bg-slate-200 p-5 gap-6 m-16 rounded-md overflow-x-auto"
	bind:this={container}
>
	{#if mouseYCoordinate}
		<div class="item ghost" style="top: {mouseYCoordinate + distanceTopGrabbedVsPointer}px;">
			{draggingItem?.name}
		</div>
	{/if}

	{#each categories as item, index (item)}
		<button
			class="bg-blue-200 p-4 rounded-lg {draggingItem == item ? 'invisible' : ''}"
			draggable="true"
			on:dragstart={(e) => {
				mouseYCoordinate = e.clientY;
				draggingItem = item;
				draggingItemIndex = index;

				// @ts-ignore
				distanceTopGrabbedVsPointer = e.target.getBoundingClientRect().y - e.clientY;
			}}
			on:drag={(e) => {
				mouseYCoordinate = e.clientY;
			}}
			on:dragover={(e) => {
				hoveredItemIndex = index;
			}}
			on:dragend={(e) => {
				draggingItem = null;
				hoveredItemIndex = null; // prevents instant swap
			}}
		>
			{item.name}
		</button>
	{/each}
</div>

<style>
	.item {
		width: 300px;
		background: white;
		padding: 10px;
		margin-bottom: 10px;
		cursor: grab;
	}

	.ghost {
		margin-bottom: 10px;
		pointer-events: none;
		z-index: 99;
		position: absolute;
		top: 0;
		left: 10;
	}

	.invisible {
		opacity: 0;
	}
</style>

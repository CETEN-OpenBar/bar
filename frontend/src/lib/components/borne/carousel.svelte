<script lang="ts">
	import { api } from '$lib/config/config';
	import type { CarouselImage, CarouselText } from '$lib/api';
	import { fly } from 'svelte/transition';
	import Carousel from 'svelte-carousel';

	export let images: Array<CarouselImage>;
	export let texts: Array<CarouselText>;
</script>

<div class="absolute bottom-0 w-full" style="z-index:1">
	<Carousel
		autoplay
		autoplayDuration={100}
		duration={15000}
		swiping={false}
		dots={false}
		arrows={false}
	>
		{#each texts as text, i}
			<div
				class="w-full h-full flex justify-center items-center p-2"
				style="text-color: {text.color}; background-color: #22283170"
				transition:fly={{ y: 100, duration: 1000, delay: 1000 }}
			>
				<p class="text-lg text-white">{text.text}</p>
			</div>
		{/each}
	</Carousel>
</div>

<Carousel autoplay autoplayDuration={3000} arrows={false} dots={false}>
	{#each images as image, i}
		<img
			src={api() + image.image_url}
			alt="dommage"
			draggable="false"
			class="w-full object-cover overflow-hidden"
		/>
	{/each}
</Carousel>

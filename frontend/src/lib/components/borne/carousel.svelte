<script lang="ts">
	import { api } from '$lib/config/config';
	import type { CarouselImage, CarouselText } from '$lib/api';
	import { onMount } from 'svelte';
	import { fly } from 'svelte/transition';
	import { quadIn, quadOut } from 'svelte/easing';

	export let images: Array<CarouselImage>;
	export let texts: Array<CarouselText>;

	let currentSlideIndex = 0;

	onMount(() => {
		startAutoScroll();
		// Change carousel-texts-roll animation duration to match the length of the text
		const carouselTextRoll = document.querySelector('.carousel-texts-roll');
		if (!carouselTextRoll) return;
		const carouselTextRollLength = carouselTextRoll.clientWidth;
		const carouselTextRollDuration = carouselTextRollLength / 200;
		// @ts-ignore
		carouselTextRoll.style.animationDuration = `${carouselTextRollDuration}s`;
	});

	function startAutoScroll() {
		nextSlide();
	}

	function nextSlide() {
		currentSlideIndex = (currentSlideIndex + 1) % images.length;

		setTimeout(() => {
			nextSlide();
		}, 5000);
	}

	function customEasingIn(t: number): number {
		// do nothing for 200ms
		if (t < 0.2) return 0;
		// ease In for 800ms (normalize t)
		return quadIn((t - 0.2) / 0.8);
	}

	function customEasingOut(t: number): number {
		// do nothing for 200ms
		if (t < 0.2) return 0;
		// ease Out for 800ms (normalize t)
		return quadIn((t - 0.2) / 0.8);
	}
</script>

<div class="carousel-container">
	<div class="carousel-images">
		{#each images as image, i}
			{#if i === currentSlideIndex}
					<img
						in:fly={{ x: -window.innerWidth, easing: customEasingIn, duration: 1000, opacity: 1 }}
						out:fly={{ x: window.innerWidth, easing: customEasingOut, duration: 1000, opacity: 1 }}
						src={image.image_url}
						alt="dommage"
						class="w-full h-full object-cover"
					/>
			{/if}
		{/each}
	</div>

	<div class="carousel-texts">
		<div class="carousel-texts-roll">
			<p class="text-center text-white text-4xl font-bold">
				{#each texts as text, i}
					<span style="color: {text.color}; margin-left: {window.innerWidth}px;">{text.text}</span>
				{/each}
			</p>
			<p class="text-center text-white text-4xl font-bold">
				{#each texts as text, i}
					<span style="color: {text.color}; margin-left: {window.innerWidth}px;">{text.text}</span>
				{/each}
			</p>
		</div>
	</div>
</div>

<style>
	.carousel-container {
		/* Set the container to take the full height and width of the viewport */
		width: 100vw;
		height: 100vh;
		position: relative;
	}

	.carousel-images {
		/* Flex container for horizontal scrolling */
		display: flex;
		overflow-x: hidden;
	}

	.carousel-images img {
		position: absolute;
		/* Make each image take the full width and height of the container */
		width: 100vw;
		height: 100vh;
	}

	.carousel-texts {
		position: absolute;
		bottom: 30px;
		max-width: 100%;
		overflow: hidden;
		background-color: #000a;
		padding: 5px;
	}

	.carousel-texts-roll {
		white-space: nowrap;
		overflow: hidden;
		display: inline-block;
		animation: carousel-texts-roll 20s linear infinite;
	}

	.carousel-texts-roll p {
		display: inline-block;
		user-select: none;
	}

	@keyframes carousel-texts-roll {
		0% {
			transform: translate3d(0%, 0, 0);
		}
		100% {
			transform: translate3d(-50%, 0, 0);
		}
	}
</style>

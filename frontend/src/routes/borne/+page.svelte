<script lang="ts">
	import Carousel from '$lib/components/carousel.svelte';
    import Pin from '$lib/components/pin.svelte';
	import type { CarouselImage, CarouselText } from '$lib/api';
	import { onMount, onDestroy } from 'svelte';
	import { carouselApi } from '$lib/requests/requests';

	let fakeImages: Array<CarouselImage> = [
		{
			id: '1',
			image_url: 'https://picsum.photos/seed/1/1920/1080'
		},
		{
			id: '2',
			image_url: 'https://picsum.photos/seed/2/1920/1080'
		}
	];

	let fakeTexts: Array<CarouselText> = [
		{
			id: '1',
			text: "je ne suis pas encore mis en place mais j'ai envie d'en apprendre plus sur la vie notemment celle de shrodinger",
			color: 'red'
		},
		{
			id: '2',
			text: 'aidez-moi :D',
			color: 'blue'
		}
	];

	let images: Array<CarouselImage> = fakeImages;
	let texts: Array<CarouselText> = fakeTexts;

	onMount(() => {
		fetchCarousel();
	});

	function fetchCarousel() {
		carouselApi()
			.getCarouselImages()
			.then((res) => {
				images = res.data;
				if (images.length === 0) {
					images = fakeImages;
				}
			});

		carouselApi()
			.getCarouselTexts()
			.then((res) => {
				texts = res.data;
				if (texts.length === 0) {
					texts = fakeTexts;
				}
			});

		setTimeout(fetchCarousel, 60000);
	}

    let card = {
        id: '',
        pin: ''
    }

	let buffer = '';
	function onType(e: KeyboardEvent) {
		if (e.key !== 'Enter') {
			buffer += e.key;
		} else {
            card.id = buffer;
        }
	}

	function pinCallback(pin:string) {
		card.pin = pin;
		card = {
			id: '',
			pin: ''
		};
	}
</script>

<svelte:window on:keydown={onType} />


{#if card.id !== ''}
    <Pin callback={pinCallback} />
{/if}

<Carousel {images} {texts} />
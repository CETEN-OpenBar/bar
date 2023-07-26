<script lang="ts">
	import Carousel from '$lib/components/borne/carousel.svelte';
	import Pin from '$lib/components/borne/pin.svelte';
	import Error from '$lib/components/error.svelte';

	import { AccountRole, type CarouselImage, type CarouselText } from '$lib/api';
	import { onMount } from 'svelte';
	import { authApi, carouselApi } from '$lib/requests/requests';
	import { goto } from '$app/navigation';

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
				if (res.data != null) images = res.data;
				if (images.length === 0) images = fakeImages;
			});

		carouselApi()
			.getCarouselTexts()
			.then((res) => {
				if (res.data != null)  texts = res.data;
				if (texts.length === 0) texts = fakeTexts;
			});

		setTimeout(fetchCarousel, 60000);
	}

	let card = {
		id: '',
		pin: ''
	};

	let buffer = '';
	function onType(e: KeyboardEvent) {
		if (e.key !== 'Enter') {
			buffer += e.key;
		} else {
			card.id = buffer;
			buffer = '';
		}
	}

	let incorrectPin = '';

	function pinCallback(pin: string) {
		card.pin = pin;

		authApi()
			.connectCard(
				{
					card_id: card.id,
					card_pin: card.pin
				},
				{
					withCredentials: true
				}
			)
			.then((res) => {
				if (res.data.account?.role === AccountRole.AccountAdmin || AccountRole.AccountSuperAdmin || AccountRole.AccountMember)
					goto('/comptoir');
				goto('/comptoir/index');
			})
			.catch(() => {
				incorrectPin = 'Mauvais code pin';
				setTimeout(() => {
					incorrectPin = '';
				}, 3000);
			});

		card = {
			id: '',
			pin: ''
		};
	}
</script>

<svelte:window on:keydown={onType} />

{#if incorrectPin != ''}
	<Error error={incorrectPin} />
{/if}

{#if card.id !== ''}
	<Pin callback={pinCallback} />
{/if}

<Carousel {images} {texts} />

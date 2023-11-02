<script lang="ts">
	// import Carousel from '$lib/components/borne/carousel.svelte';
	import Pin from '$lib/components/borne/pin.svelte';
	import Error from '$lib/components/error.svelte';
	import Carousel from 'svelte-carousel';


	import {
		AccountState,
		type CarouselImage,
		type CarouselText,
		type ConnectCardRequest
	} from '$lib/api';
	import { onMount } from 'svelte';
	import { authApi, carouselApi } from '$lib/requests/requests';
	import { goto } from '$app/navigation';
	import FsLoading from '$lib/components/borne/fs_loading.svelte';
	import { api } from '$lib/config/config';
	import { fly } from 'svelte/transition';
	import ReadCard from '$lib/components/readCard.svelte';

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

	let display: Boolean = false;
	let images: Array<CarouselImage> = fakeImages;
	let texts: Array<CarouselText> = fakeTexts;

	onMount(() => {
		fetchCarousel();
	});

	const preloadImage = (src: string) =>
		new Promise((resolve, reject) => {
			const image = new Image();
			image.onload = resolve;
			image.onerror = reject;
			image.src = src;
		});

	function fetchCarousel() {
		carouselApi()
			.getCarouselImages()
			.then((res) => {
				if (res.data != null) images = res.data;
				if (images.length === 0) images = fakeImages;
				Promise.all(images.map((x) => preloadImage(api() + x.image_url))).finally(() => {
					setTimeout(() => {
						display = true;
					}, 1500);
				});
			});

		carouselApi()
			.getCarouselTexts()
			.then((res) => {
				if (res.data != null) texts = res.data;
				if (texts.length === 0) texts = fakeTexts;
			});

		setTimeout(fetchCarousel, 60000);
	}

	let card = {
		id: '',
		pin: ''
	};

	function cardCallback(id: string) {
		card.id = id;
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
				if (res.data.account?.state == AccountState.AccountOK) goto('/borne/commande');
				if (res.data.account?.state == AccountState.AccountNotOnBoarded) goto('/borne/onboarding');
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

<ReadCard callback={cardCallback} />


{#if incorrectPin != ''}
	<Error error={incorrectPin} />
{/if}

{#if card.id !== ''}
	<Pin callback={pinCallback} />
{/if}

{#if display}
	<Carousel
	autoplay
	autoplayDuration={3000}

	arrows={false}
	dots={false}
	>
	{#each images as image, i}
			<img src={api() + image.image_url}
				alt="dommage"
				draggable="false"
				class="w-full h-full object-cover"
			/>
	{/each}
	</Carousel>
{:else}
	<FsLoading />
{/if}

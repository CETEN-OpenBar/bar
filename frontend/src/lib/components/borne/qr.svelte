<script lang="ts">
	import { apiws } from '$lib/config/config';
	import { authApi } from '$lib/requests/requests';
	import Error from '../error.svelte';
	import Pin from './pin.svelte';
	import { onMount, onDestroy } from 'svelte';

	let qr: any = undefined;
	let ask = true;
	let error = '';
	let ws: WebSocket;

	let spinner = false;
	let connected = false;

	onMount(() => {
		console.log(apiws());
		ws = new WebSocket(apiws() + '/account/qr');
		ws.onmessage = (event) => {
			if (event.data == 'scanned') {
				spinner = true;
			} else if (event.data == 'connected') {
				spinner = false;
				connected = true;
				// Send refresh command in 2 seconds
				setTimeout(() => {
					window.location.reload();
				}, 2000);
			}
		};
	});

	onDestroy(() => {
		ws.close();
	});

	function pinCallback(pin: string) {
		if (pin == '') {
			ask = false;
			return;
		}
		authApi()
			.getAccountQR({ card_pin: pin }, { withCredentials: true })
			.then((res) => {
				qr = res.data;
			})
			.catch(() => {
				error = 'Impossible de générer le QR code';
				setTimeout(() => {
					error = '';
				}, 3000);
			});
	}
</script>

{#if qr !== undefined || ask == false}
	<div class="grid">
		{#if spinner}
			<div id="overlay" class="w-full" style="background-color: #0009; grid-area: 1 / 1 / 2 / 2;">
				<div class="spinner" />
			</div>
		{/if}
		{#if connected}
			<div id="overlay" class="w-full" style="background-color: #0009; grid-area: 1 / 1 / 2 / 2;">
				<svg class="checkmark" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 52 52"
					><circle class="checkmark__circle" cx="26" cy="26" r="25" fill="none" /><path
						class="checkmark__check"
						fill="none"
						d="M14.1 27.2l7.1 7.2 16.7-16.8"
					/>
				</svg>
			</div>
		{/if}
		<img
			class="top-0 left-0 w-full"
			style="grid-area: 1 / 1 / 2 / 2;"
			src="data:image/png;base64,{qr}"
			alt="Entrez-votre pin pour voir le code"
		/>
	</div>
{:else}
	<Pin callback={pinCallback} />
	<!-- {pinCallback('1234')} -->
{/if}

{#if error !== ''}
	<Error {error} />
{/if}

<style>
	.checkmark {
		width: 60px;
		height: 60px;
		border-radius: 50%;
		display: block;
		stroke-width: 2;
		stroke: #4bb71b;
		stroke-miterlimit: 10;
		box-shadow: inset 0px 0px 0px #4bb71b;
		animation: fill 0.4s ease-in-out 0.4s forwards, scale 0.3s ease-in-out 0.9s both;
		position: relative;
		top: 50%;
		left: 50%;
	}
	.checkmark__circle {
		stroke-dasharray: 166;
		stroke-dashoffset: 166;
		stroke-width: 2;
		stroke-miterlimit: 10;
		stroke: #4bb71b;
		fill: #000a;
		animation: stroke 0.6s cubic-bezier(0.65, 0, 0.45, 1) forwards;
	}

	.checkmark__check {
		transform-origin: 50% 50%;
		stroke-dasharray: 48;
		stroke-dashoffset: 48;
		animation: stroke 0.3s cubic-bezier(0.65, 0, 0.45, 1) 0.8s forwards;
	}

	@keyframes stroke {
		100% {
			stroke-dashoffset: 0;
		}
	}

	@keyframes scale {
		0%,
		100% {
			transform: translate(-50%, -50%);
		}

		50% {
			transform: translate(-50%, -50%) scale3d(1.1, 1.1, 1);
		}
	}

	.spinner {
		position: relative;
		left: 50%;
		top: 50%;
		height: 60px;
		width: 60px;
		-webkit-animation: rotation 0.6s infinite linear;
		-moz-animation: rotation 0.6s infinite linear;
		-o-animation: rotation 0.6s infinite linear;
		animation: rotation 0.6s infinite linear;
		border-left: 6px solid rgba(0, 174, 239, 0.15);
		border-right: 6px solid rgba(0, 174, 239, 0.15);
		border-bottom: 6px solid rgba(0, 174, 239, 0.15);
		border-top: 6px solid rgba(0, 174, 239, 0.8);
		border-radius: 100%;
	}

	@-webkit-keyframes rotation {
		from {
			-webkit-transform: translate(-50%, -50%) rotate(0deg);
		}
		to {
			-webkit-transform: translate(-50%, -50%) rotate(359deg);
		}
	}
	@-moz-keyframes rotation {
		from {
			-moz-transform: translate(-50%, -50%) rotate(0deg);
		}
		to {
			-moz-transform: translate(-50%, -50%) rotate(359deg);
		}
	}
	@-o-keyframes rotation {
		from {
			-o-transform: translate(-50%, -50%) rotate(0deg);
		}
		to {
			-o-transform: translate(-50%, -50%) rotate(359deg);
		}
	}
	@keyframes rotation {
		from {
			transform: translate(-50%, -50%) rotate(0deg);
		}
		to {
			transform: translate(-50%, -50%) rotate(359deg);
		}
	}
</style>

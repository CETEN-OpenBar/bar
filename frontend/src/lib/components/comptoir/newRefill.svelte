<script lang="ts">
	import { refillsApi } from '$lib/requests/requests';
	import Error from '../error.svelte';
	import Success from '../success.svelte';

	export let close: () => void;

	let success = '';
	let error = '';

	let card = {
		id: '',
		amount_euros: 0,
		amount_cents: 0,
		amount: 0
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
</script>

{#if success != ''}
	<Success message={success} />
{/if}

{#if error != ''}
	<Error {error} />
{/if}

<svelte:window on:keydown={onType} />

<!-- Popup overlay -->
<button
	id="overlay"
	class="absolute w-full h-full top-0 left-0 bg-black bg-opacity-50 flex justify-center items-center z-10 hover:cursor-default"
	on:click={() => {
		close();
	}}
/>

<div id="popup" class="absolute w-full h-full top-0 left-0 flex justify-center items-center">
	<div
		class="relative text-black flex flex-col justify-center items-center gap-4 p-10 h-64 bg-white rounded-xl shadow-xl z-20"
	>
		<!-- button to close the popup -->
		<button
			class="absolute top-0 right-0 p-2 text-xl font-bold m-2 rounded-full transition-all text-black"
			on:click={() => {
				close();
			}}
		>
			<iconify-icon icon="mdi:close" />
		</button>
		<!-- prompt to scan the card -->
		{#if card.id == ''}
			<h1 class="text-3xl">Veuillez scanner la carte.</h1>
		{:else}
			<h1 class="text-3xl">Veuillez entrer le montant de la recharge.</h1>

			<div class="flex flex-row">
				<input
					type="number"
					class="text-sm bg-gray-200 rounded-md w-24 text-center"
					bind:value={card.amount_euros}
				/>
				<div class="self-end">.</div>
				<input
					type="number"
					class="text-sm bg-gray-200 rounded-md w-24 text-center"
					bind:value={card.amount_cents}
				/>
				<div class="self-end">€</div>
			</div>

			<button
				class="text-3xl bg-green-500 p-4 rounded-xl hover:bg-green-700 transition-all text-white"
				on:click={() => {
					if (card.amount_euros < 0 || card.amount_cents < 0 || card.amount_cents > 99) {
						error = 'Montant invalide.';
						return;
					}
					card.amount = card.amount_euros * 100 + card.amount_cents;
					refillsApi()
						.postRefill(card.id, card.amount, { withCredentials: true })
						.then(() => {
							success = 'Recharge effectuée avec succès.';
							close();
						})
						.catch(() => {
							error = 'Une erreur est survenue.';
						});
				}}
			>
				Valider
			</button>
		{/if}
	</div>
</div>

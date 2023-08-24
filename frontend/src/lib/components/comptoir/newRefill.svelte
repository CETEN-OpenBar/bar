<script lang="ts">
	import { RefillType } from '$lib/api';
	import { refillsApi } from '$lib/requests/requests';
	import { formatPrice, parsePrice } from '$lib/utils';
	import Error from '../error.svelte';
	import Success from '../success.svelte';

	export let close: () => void;

	let success = '';
	let error = '';

	let card = {
		id: '',
		amount: 0,
		type: RefillType.RefillOther
	};

	let rebounce = 0;

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
		class="relative text-black flex flex-col justify-center items-center gap-4 p-10 h-96 bg-white rounded-xl shadow-xl z-20"
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

			<div class="flex flex-col gap-8">
				<div class="flex flex-col">
					<label for="price-new" class="block text-xl mb-2 align-middle">Montant :</label>
					<input
						type="number"
						id="price-new"
						name="price"
						placeholder="Montant de la recharge"
						class="text-sm bg-gray-200 rounded-md p-2 text-center"
						required
						aria-describedby="text-error"
						on:input={(e) => {
							// @ts-ignore
							card.amount = parsePrice(e.target?.value);
							console.log(card.amount);

							// rebounce update
							clearTimeout(rebounce);
							rebounce = setTimeout(() => {
								let elem = document.getElementById('price-new');
								if (elem) {
									// @ts-ignore
									elem.value = '';
									// @ts-ignore
									elem.placeholder = formatPrice(card.amount);
								}
							}, 1000);
						}}
					/>
				</div>

				<div class="flex flex-col">
					<label for="refill-type" class="block text-xl mb-2 align-middle">Type :</label>
					<select
						id="refill-type"
						name="refill-type"
						class="text-sm bg-gray-200 rounded-md p-2 text-center"
						on:change={(e) => {
							// @ts-ignore
							card.type = e.target?.value;
						}}
					>
						<option value={RefillType.RefillOther}>Autre</option>
						<option value={RefillType.RefillCash}>Liquide</option>
						<option value={RefillType.RefillCard}>Carte</option>
					</select>
				</div>
			</div>

			<button
				class="text-3xl bg-green-500 p-4 rounded-xl hover:bg-green-700 transition-all text-white"
				on:click={() => {
					refillsApi()
						.postRefill(card.id, card.amount, card.type, { withCredentials: true })
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

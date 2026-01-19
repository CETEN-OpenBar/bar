<script lang="ts">
	import { RefillType } from '$lib/api';
	import { refillsApi } from '$lib/requests/requests';
	import { formatPrice, parsePrice } from '$lib/utils';
	import Error from '../error.svelte';
	import ReadCard from '../readCard.svelte';
	import Success from '../success.svelte';

	export let close: () => void;
	export let cardId: string | undefined = undefined;

	let success = '';
	let error = '';

	function cardCallback(id: string) {
		card.id = id;
	}
	let rebounce = 0;
	let card: {
		id: string;
		amount: number;
		type: RefillType | undefined;
	} = {
		id: cardId ?? '',
		amount: 0,
		type: undefined
	};
</script>

{#if success != ''}
	<Success message={success} />
{/if}

{#if error != ''}
	<Error {error} />
{/if}

{#if !cardId}
	<ReadCard callback={cardCallback} />
{/if}

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
		class="relative text-black dark:text-white flex flex-col justify-center items-center gap-4 p-10 h-96 bg-white dark:bg-gray-800 rounded-xl shadow-xl z-20"
	>
		<!-- button to close the popup -->
		<button
			class="absolute top-0 right-0 p-2 text-xl font-bold m-2 rounded-full transition-all text-black dark:text-white"
			on:click={() => {
				close();
			}}
		>
			<iconify-icon icon="mdi:close" />
		</button>
		<!-- prompt to scan the card -->
		{#if card.id == ''}
			<h1 class="text-3xl">Veuillez scanner la carte.</h1>
		{:else if card.type === undefined}
			<h1>Veuillez selectionner le moyen de paiement.</h1>
			<div class="grid grid-cols-2 w-full gap-8 p-4">
				<button
					class="bg-blue-500 hover:bg-blue-700 text-white font-bold text-xl py-12 px-4 rounded "
					on:click={() => (card.type = RefillType.RefillCard)}
				>
					Carte
				</button>
				<button
					class="bg-yellow-500 hover:bg-yellow-600 text-white font-bold text-xl py-12 px-4 rounded"
					on:click={() => (card.type = RefillType.RefillCash)}
				>
					Liquide
				</button>
			</div>
		{:else}
			<h1 class="text-3xl">Veuillez entrer le montant de la recharge.</h1>

			<!-- svelte-ignore a11y-no-static-element-interactions -->
			<div
				class="flex flex-col gap-8"
				on:keypress={(e) => {
					if (e.key == 'Enter' && card.type)
						refillsApi()
							.postRefill(card.id, card.amount, card.type, { withCredentials: true })
					.then(() => {
						success = 'Recharge effectuée avec succès.';
						setTimeout(() => close(), 2000);
					})
					.catch(() => {
						error = 'Une erreur est survenue.';
						setTimeout(() => close(), 2000);
					});
				}}
			>
				<div class="flex flex-col">
					<label for="price-new" class="block text-xl mb-2 align-middle dark:text-white">Montant :</label>
					<input
						type="number"
						id="price-new"
						name="price"
						placeholder="Montant de la recharge"
						class="text-sm bg-gray-200 dark:bg-gray-700 dark:text-white rounded-md p-2 text-center"
						required
						aria-describedby="text-error"
						on:input={(e) => {
							// @ts-ignore
							card.amount = parsePrice(e.target?.value);

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
			</div>

			<button
				class="text-3xl bg-green-500 p-4 rounded-xl hover:bg-green-700 transition-all text-white"
				on:click={() => {
					if (card.type) {
						refillsApi()
							.postRefill(card.id, card.amount, card.type, { withCredentials: true })
							.then(() => {
								success = 'Recharge effectuée avec succès.';
								setTimeout(() => close(), 2000);
							})
							.catch(() => {
								error = 'Une erreur est survenue.';
								setTimeout(() => close(), 2000);
							});
					}
				}}
			>
				Valider
			</button>
		{/if}
	</div>
</div>

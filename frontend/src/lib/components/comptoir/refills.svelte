<script lang="ts">
	import type { Refill } from '$lib/api';
	import { api } from '$lib/config/config';
	import { refillsApi } from '$lib/requests/requests';
	import { formatPrice } from '$lib/utils';
	import { onDestroy, onMount } from 'svelte';
	// import RefillPopup from './refillPopup.svelte';

	export let amount: number = 3;

	let refills: Array<Refill> = [];
	let interval: number;

	onMount(() => {
		reloadRefills();
		interval = setInterval(() => {
			reloadRefills();
		}, 2000);
	});

	onDestroy(() => {
		clearInterval(interval);
	});

	function reloadRefills() {
		refillsApi()
			.getRefills(0, amount, undefined, undefined, { withCredentials: true })
			.then((res) => {
				refills = res.data;
			});
	}
</script>

<!-- {#if displayRefill}
	<RefillPopup refill={displayRefill} close={()=>{
		displayRefill = null;

	}} />
{/if} -->

<!-- Good looking dropdown for refill -->
<div class="w-full">
	<div class="flex flex-col">
		<div class="flex flex-col p-20">
			{#each refills as refill}
				<div
					class="flex flex-row p-4 justify-between text-black mt-5 border-4 border-white rounded-xl {refill.state ==
					'canceled'
						? 'bg-gray-200'
						: ''} {refill.state == 'valid' ? 'bg-green-200' : ''}"
				>
					<!-- Display localTimeString and name of the account -->
					<div class="flex flex-col">
						<div class="text-sm">{new Date(refill.issued_at * 1000).toLocaleString()}</div>
						<div class="text-lg font-semibold">à : {refill.account_name}</div>
						{#if refill.canceled_by}
							<div class="text-sm">Annulé par : {refill.canceled_by_name}</div>
						{/if}
					</div>
					<div class="flex flex-col">
						<div class="text-lg self-center">{formatPrice(refill.amount)}</div>
						{#if refill.state == 'valid'}
							<button
								class="text-sm text-white self-center bg-red-500 rounded-xl p-2"
								on:click={() => {
									refillsApi()
										.patchRefillId(refill.account_id, refill.id, 'canceled', {
											withCredentials: true
										})
										.then(() => {
											reloadRefills();
										});
								}}>Annuler</button
							>
						{/if}
						{#if refill.state == 'canceled'}
							<button
								class="text-sm text-white  self-center bg-red-500 rounded-xl p-2"
								on:click={() => {
									refillsApi()
										.patchRefillId(refill.account_id, refill.id, 'valid', { withCredentials: true })
										.then(() => {
											reloadRefills();
										});
								}}>Re-valider</button
							>
						{/if}
					</div>
				</div>
			{/each}
		</div>
	</div>
</div>

<script lang="ts">
	import { goto } from '$app/navigation';
	import NewRefill from '$lib/components/comptoir/newRefill.svelte';
	import Refills from '$lib/components/comptoir/refills.svelte';
	import Pin from '$lib/components/borne/pin.svelte';
	import ReadCard from '$lib/components/readCard.svelte';
	import { open_caisse, open_door, open_ventilo } from '$lib/local/local';
	import ChangePassword from '$lib/components/comptoir/changePassword.svelte';
	import Password from '$lib/components/password.svelte';

	let newRefill = false;

	function close() {
		newRefill = false;
	}

	function reset(){
		askForCard = false;
		askForPassword = false;
		to_call = () => {};
	}

	let to_call = open_door;
	let infos = {
		card_id: '',
		card_pin: ''
	};
	let askForCard = false;
	let askForPassword = false;
</script>

{#if askForCard}
	<!-- Popup overlay -->
	<button
		id="overlay"
		class="absolute w-full h-full top-0 left-0 bg-black bg-opacity-50 flex justify-center items-center z-10 hover:cursor-default"
		on:click={() => {
			reset();
		}}
	/>

	<div id="popup" class="absolute w-full h-full top-0 left-0 flex justify-center items-center">
		<div
			class="relative text-black flex flex-col justify-center items-center gap-4 p-10 h-96 bg-white rounded-xl shadow-xl z-20"
		>
			<button
				class="absolute top-0 right-0 p-2 text-xl font-bold m-2 rounded-full transition-all text-black"
				on:click={() => {
					reset();
				}}
			>
				<iconify-icon icon="mdi:close" />
			</button>
			<h1 class="text-3xl">Veuillez scanner la carte.</h1>
		</div>
	</div>

	<ReadCard
		callback={(id) => {
			infos.card_id = id;
			askForPassword = true;
			askForCard = false;
		}}
	/>
{/if}

{#if askForPassword}
	<Password
		callback={(password) => {
			infos.card_pin = password;
			to_call(infos.card_id, infos.card_pin);
			reset();
		}}
	/>
{/if}

<div class="gap-16 p-5 w-full h-full text-white">
	<div class="flex flex-row justify-between gap-16 p-2 w-full text-white">
		<div class="flex flex-row">
			<h1 class="text-3xl">Panneau de contrôle</h1>
			<button
				class="text-3xl bg-blue-700 p-2 rounded-xl hover:bg-blue-900 transition-all ml-2"
				on:click={() => {
					to_call = open_door;
					askForCard = true;
				}}>porte</button
			>
			<button
				class="text-3xl bg-blue-700 p-2 rounded-xl hover:bg-blue-900 transition-all ml-2"
				on:click={() => {
					to_call = open_ventilo;
					askForCard = true;
				}}>ventilo</button
			>
			<button
				class="text-3xl bg-blue-700 p-2 rounded-xl hover:bg-blue-900 transition-all ml-2"
				on:click={() => {
					to_call = open_caisse;
					askForCard = true;
				}}>caisse</button
			>
		</div>
		
		<div class="flex flex-row">
			<button
				class="text-3xl bg-blue-700 p-2 rounded-xl hover:bg-blue-900 transition-all mr-2"
				on:click={() => (newRefill = true)}>Nouvelle Recharge</button
			>
			<button
				class="text-3xl bg-blue-700 p-2 rounded-xl hover:bg-blue-900 transition-all"
				on:click={() => goto('/comptoir/c/transactions')}>Transactions</button
			>
		</div>
	</div>
	<hr class="col-span-3" />

	{#if newRefill}
		<NewRefill {close} />
	{/if}

	<Refills amount={6} />
</div>

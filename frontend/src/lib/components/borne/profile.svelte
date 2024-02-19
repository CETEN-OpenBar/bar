<script lang="ts">
	import { goto } from '$app/navigation';
	import type { Account } from '$lib/api';
	import { accountsApi, authApi } from '$lib/requests/requests';
	import { formatPrice } from '$lib/utils';
	import Error from '../error.svelte';
	import Pin from './pin.svelte';
	import Success from '../success.svelte';
	import Stars from '../random/stars.svelte';
	import Price from '../random/price.svelte';

	export let account: Account;
	export let logout: () => void;

	let pin_step = 0;
	let card = {
		old_pin: '',
		new_pin: ''
	};

	let popup = false;
	let errorTimeout: number;

	if (account.google_picture == '') {
		account.google_picture = `https://www.gravatar.com/avatar/${account.email_address}?d=retro`;
	}

	function setNewPin(pin: string) {
		if (pin == '') {
			pin_step = 0;
			return;
		}
		card.new_pin = pin;

		accountsApi()
			.patchAccountPin(
				{
					old_card_pin: card.old_pin,
					new_card_pin: card.new_pin
				},
				{
					withCredentials: true
				}
			)
			.then(() => {
				pin_step = 3;
			})
			.catch(() => {
				pin_step = 4;
			})
			.finally(() => {
				card = {
					old_pin: '',
					new_pin: ''
				};
				errorTimeout = setTimeout(() => (pin_step = 0), 3000);
			});
	}

	function nextStep(pin: string) {
		if (pin == '') {
			pin_step = 0;
			return;
		}
		card.old_pin = pin;
		pin_step = 2;
	}

	function logoutAccount() {
		authApi()
			.logout({
				withCredentials: true
			})
			.then(() => {
				logout();
			})
			.catch(() => {
				logout();
			});
	}
</script>

<!-- Icon with a dropdown menu for the actions possible -->
<div class="relative">
	<div class="flex flex-row">
		<button
			on:click={() => {
				popup = !popup;
			}}
			class="{popup ? 'border-gray-200' : 'border-transparent'} rounded-full border-4"
		>
			<img
				src={account.google_picture}
				alt="Y'a pas"
				class={`w-14 h-14 rounded-full cursor-pointer text-white`}
			/>
		</button>
		<button
			on:click={() => {
				popup = !popup;
			}}
			class="ml-2 self-center flex flex-col justify-start"
		>
			<h1 class="text-lg text-white font-bold">{account.nickname || account.first_name +" "+ account.last_name}</h1>
			<div class="flex flex-col justify-center">
				<Price amount={account.balance} class="text-md text-white self-center" />
				<Stars stars={account.points} class="text-md text-white self-center" />
			</div>
		</button>
	</div>

	{#if popup}
		<!-- the button will be on the left of the screen, this pops up underneath -->
		<div class="absolute bg-gray-200 rounded-xl shadow-lg p-4 z-20" style="display:ruby">
			<div class="bg-white rounded-xl">
				<div class="flex flex-col p-4">
					<div class="flex flex-row align-middle">
						<img
							src={account.google_picture}
							alt="User Avatar"
							class={`w-11 rounded-full cursor-pointer align-middle`}
						/>
						<div class="flex flex-col ml-5">
							<h4 class="text-md font-bold">{account.first_name} {account.last_name}</h4>
							{#if account.email_address != ''}
								<h5 class="text-xs">{account.email_address}</h5>
							{/if}
						</div>
					</div>
					<hr class="my-2" />
					<button
						class=" bg-blue-500 text-white rounded-lg p-2 mb-2"
						on:click={() => {
							popup = false;
							goto('/borne/index');
						}}>Mon compte</button
					>
					<button
						class=" bg-blue-500 text-white rounded-lg p-2 mb-2"
						on:click={() => {
							pin_step = 1;
							popup = false;
							clearTimeout(errorTimeout);
						}}>Changer de PIN</button
					>
					<hr class="my-2" />
					<button class="w-full bg-red-500 text-white rounded-lg p-2 mb-2" on:click={logoutAccount}
						>Déconnexion</button
					>
				</div>
			</div>
		</div>
	{/if}
</div>

{#if pin_step == 1}
	<Pin custom_text="Entrez votre ancien code pin" callback={nextStep} />
{/if}
{#if pin_step == 2}
	<Pin custom_text="Entrez votre nouveau code pin" callback={setNewPin} />
{/if}
{#if pin_step == 3}
	<Success message="Votre code pin a bien été changé" />
{/if}
{#if pin_step == 4}
	<Error error="Une erreur est survenue" />
{/if}

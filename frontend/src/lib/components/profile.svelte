<script lang="ts">
	import type { Account } from '$lib/api';

	export let account: Account;
	export let size: string = '10';
	export let logout: () => void;

	let popup = false;

	if (account.google_picture == '') {
		account.google_picture = `https://www.gravatar.com/avatar/${account.email_address}?d=retro`;
	}
</script>

<!-- Icon with a dropdown menu for the actions possible -->
<div class="relative">
	<button
		on:click={() => {
			popup = !popup;
		}}
		class="{popup ? 'border-gray-200' : 'border-white'} rounded-full border-4"
	>
		<img
			src={account.google_picture}
			alt="User Avatar"
			class={`w-${size} h-${size} rounded-full cursor-pointer`}
		/>
	</button>

	{#if popup}
		<!-- the button will be on the left of the screen, this pops up underneath -->
		<div class="absolute bg-gray-200 rounded-xl shadow-lg p-4 z-20">
			<div class="bg-white rounded-xl">
				<div class="flex flex-col p-4">
					<div class="flex flex-row align-middle">
						<img
							src={account.google_picture}
							alt="User Avatar"
							class={`w-${size} h-${size} rounded-full cursor-pointer align-middle`}
						/>
						<div class="flex flex-col ml-5">
							<h4 class="text-md font-bold">{account.first_name} {account.last_name}</h4>
							{#if account.email_address != ''}
								<h5 class="text-xs">{account.email_address}</h5>
							{/if}
						</div>
					</div>
                    <hr class="my-2" />
                    <button class=" bg-blue-500 text-white rounded-lg p-2 mb-2">Mon compte</button>
                    <button class=" bg-blue-500 text-white rounded-lg p-2 mb-2">Changer de PIN</button>
                    <hr class="my-2" />
					<button class="w-full bg-red-500 text-white rounded-lg p-2 mb-2" on:click={logout}
						>Déconnexion</button
					>
				</div>
			</div>
		</div>
	{/if}
</div>

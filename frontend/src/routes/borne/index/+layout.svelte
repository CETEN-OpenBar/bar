<script lang="ts">
	import { goto } from '$app/navigation';
	import { accountsApi } from '$lib/requests/requests';
	import type { Account } from '$lib/api';
	import { onMount } from 'svelte';
	import { store } from '$lib/store/store';
	import Profile from '$lib/components/borne/profile.svelte';
	import 'iconify-icon';

	let account: Account | undefined = undefined;

	onMount(() => {
		accountsApi()
			.getAccount({
				withCredentials: true
			})
			.then((res) => {
				account = res.data.account;
				store.set({ account });
			})
			.catch(() => {
				goto('/borne');
			});
	});
</script>

{#if account !== undefined}
	<div id="main" class="absolute top-0 left-0 w-screen h-screen" style="background-color:#393E46">
		<div class="w-full p-4 flex justify-between" style="background-color:#222831">
			<!-- left div -->
			<Profile
				{account}
				logout={() => {
					goto('/borne');
				}}
			/>
			<!-- Good looking "order" button -->
			<button
				class="flex items-center space-x-2 px-4 py-2 rounded-lg bg-green-500 hover:bg-green-600 transition-colors duration-300 animate-pulse"
				on:click={() => {
					goto('/borne/index/commande');
				}}
			>
				<iconify-icon class="text-white text-2xl" icon="mingcute:fork-fill" />
				<span class="text-white text-lg font-bold">Commander</span>
			</button>
		</div>

		<slot />
	</div>
{/if}

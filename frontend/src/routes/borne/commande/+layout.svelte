<script lang="ts">
	import { goto } from '$app/navigation';
	import { accountsApi, authApi } from '$lib/requests/requests';
	import type { Account } from '$lib/api';
	import { onMount } from 'svelte';
	import { store } from '$lib/store/store';
	import Profile from '$lib/components/borne/profile.svelte';
	import 'iconify-icon';
	import { fly } from 'svelte/transition';

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
			
		disconnectInterval = setInterval(logout, 15000);
	});

	let disconnectInterval: number | undefined = undefined;

	function logout() {
		authApi()
			.logout({ withCredentials: true })
			.then(() => {
				clearInterval(disconnectInterval);
				goto('/comptoir');
			});
	}

	function onAction() {
		clearInterval(disconnectInterval);
		disconnectInterval = setInterval(logout, 15000);
	}
</script>


<svelte:window on:mousemove={onAction} on:click={onAction} on:scroll={onAction} on:keypress={onAction} />

{#if account !== undefined}
	<div
		id="main"
		class="absolute top-0 left-0 w-screen h-screen"
		style="background-color:#393E46"
		out:fly={{ y: window.innerHeight, delay: 500 }}
	>
		<slot />
	</div>
{/if}

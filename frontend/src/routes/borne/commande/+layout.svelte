<script lang="ts">
	import { goto } from '$app/navigation';
	import { accountsApi } from '$lib/requests/requests';
	import type { Account } from '$lib/api';
	import { onMount } from 'svelte';
	import { store } from '$lib/store/store';
	import 'iconify-icon';
	import { fly } from 'svelte/transition';
	import Autodisconnect from '$lib/components/random/autodisconnect.svelte';

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

<Autodisconnect delay={60000} location="/borne" />

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

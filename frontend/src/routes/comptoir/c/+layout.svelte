<script lang="ts">
	import { goto } from '$app/navigation';
	import { accountsApi } from '$lib/requests/requests';
	import type { Account } from '$lib/api';
	import { onMount } from 'svelte';
	import { store } from '$lib/store/store';
	import 'iconify-icon';
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
				goto('/comptoir');
			});
	});
</script>

<Autodisconnect delay={900000} location="/comptoir" />

{#if account !== undefined}
	<div id="main" class="absolute top-0 left-0 w-screen h-screen" style="background-color:black;">
		<slot />
	</div>
{/if}

<script lang="ts">
	import type { Account } from '$lib/api';
	import { onMount, onDestroy } from 'svelte';
	import { store } from '$lib/store/store';
    import { goto } from '$app/navigation';
	import Profile from '$lib/components/profile.svelte';

	let account: Account | undefined = undefined;
	let unsub: () => void;

	onMount(() => {
		unsub = store.subscribe((state) => {
			account = state.account;
		});
	});

	onDestroy(() => {
		unsub();
	});
</script>

{#if account !== undefined}
	<div id="main" class="w-full h-full">
		<Profile account={account} logout={() => {goto('/borne')}} />
	</div>
{/if}

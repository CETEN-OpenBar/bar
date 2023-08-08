<script lang="ts">
	import { fly } from 'svelte/transition';
	import QrOnboarding from '$lib/components/borne/qr_onboarding.svelte';
	import { onMount } from 'svelte';
	import { elasticIn } from 'svelte/easing';

	let shown = false;

	onMount(() => {
		setTimeout(() => {
			shown = true;
		}, 100);
	});
</script>

{#if shown}
	<div
		class="grid lg:grid-cols-3 grid-cols-1 gap-16 p-16 w-full"
		in:fly={{ x: -window.innerWidth, easing: elasticIn, duration: 1000, opacity: 1 }}
	>
		<QrOnboarding />
	</div>
{/if}

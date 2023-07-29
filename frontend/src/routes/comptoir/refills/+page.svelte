<script lang="ts">
	import { goto } from '$app/navigation';
	import Refills from '$lib/components/comptoir/refills.svelte';
</script>

<div class="gap-16 p-5 w-full h-full text-white">
	<div class="flex flex-row justify-between gap-16 p-2 w-full text-white">
		<h1 class="text-3xl">Panneau de contrôle</h1>
		<button
			class="text-3xl bg-blue-700 p-2 rounded-xl hover:bg-blue-900 transition-all"
			on:click={() => goto('/comptoir/transactions')}>Transactions</button
		>
	</div>
	<hr class="col-span-3" />

	<Refills amount={6} />
</div>

<script lang="ts">
	import { goto } from '$app/navigation';
	import NewRefill from '$lib/components/comptoir/newRefill.svelte';
	import Refills from '$lib/components/comptoir/refills.svelte';

  let newRefill = false;

  function close() {
    newRefill = false;
  }
</script>

<div class="gap-16 p-5 w-full h-full text-white">
	<div class="flex flex-row justify-between gap-16 p-2 w-full text-white">
		<h1 class="text-3xl">Panneau de contrôle</h1>
		<div class="flex flex-row">
			<button
				class="text-3xl bg-blue-700 p-2 rounded-xl hover:bg-blue-900 transition-all mr-2"
				on:click={() => newRefill = true}>Nouvelle Recharge</button
			>
			<button
				class="text-3xl bg-blue-700 p-2 rounded-xl hover:bg-blue-900 transition-all"
				on:click={() => goto('/comptoir/transactions')}>Transactions</button
			>
		</div>
	</div>
	<hr class="col-span-3" />

  {#if newRefill}
    <NewRefill close={close} />
  {/if}

	<Refills amount={6} />
</div>

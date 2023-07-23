<script lang="ts">
	import type { Refill } from '$lib/api';
	import { refillsApi } from '$lib/requests/requests';
	import { formatPrice } from '$lib/utils';
	import { onMount } from 'svelte';

	export let amount: number = 3;

	let refills: Array<Refill> = [];

	onMount(() => {
		refillsApi()
			.getSelfRefills(0, amount, undefined, undefined, { withCredentials: true })
			.then((res) => {
				if (!(res.data.refills instanceof Array)) return
				refills = res.data.refills;
			});
	});
</script>

<!-- Good looking dropdown for transaction -->
<div class="w-full">
	<div class="flex flex-col">
		<div class="flex flex-row justify-between">
			<div class="text-lg font-semibold">Transactions</div>
			<div class="text-lg font-semibold">Montant</div>
		</div>
		<div class="flex flex-col">
			{#each refills as refill}
				<div class="flex flex-row justify-between mt-5 border-4 border-white rounded-xl">
					<div class="p-5 h-full pr-4 w-full">
						<!-- refill.issued_at is unix seconds -->
						{new Date(refill.issued_at*1000).toLocaleString()}
					</div>
					<div class="border-r border-l border-gray-400">

					</div>
					<div class="p-5 pl-4 w-full text-lg text-center self-center text-green-600">
						+{formatPrice(refill.amount)}
					</div>
				</div>
			{/each}
		</div>
	</div>
</div>

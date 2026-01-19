<script lang="ts">
	import { formatDateTime, formatPrice, time2Utc } from '$lib/utils';
	import { RefillType } from '$lib/api';
    import type { Refill } from '$lib/api';
	import { refillsApi } from '$lib/requests/requests';
	import { onMount } from 'svelte';

	let refills: Refill[] = [];

	let page: number = 0;
	let maxPage: number = 0;
	let nextPage = () => {
		if (page <= maxPage) {
			page++;
			reloadItems();
		}
	};
	let prevPage = () => {
		if (page > 0) {
			page--;
			reloadItems();
		}
	};
	let itemsPerPage = 100;

	let today = new Date();
	// Set it precisely to midnight LOCAL time today
	today.setHours(0, 0, 0, 0);
	let startTimestampSeconds = today.getTime() / 1000;
	let endTimestampSeconds = startTimestampSeconds + (24 * 60 * 60);
	let startDate = Math.floor(startTimestampSeconds).toString();
	let endDate = Math.floor(endTimestampSeconds).toString();

	async function reloadItems() {
		let resp = await refillsApi().getRefills(page, itemsPerPage, startDate, endDate, {
			withCredentials: true
		});
		refills = resp.data.refills ?? [];

		maxPage = resp.data.max_page;
		// Load all pages
		for (let p = 0; p <= maxPage; p++) {
			if (p != page) {
				let resp = await refillsApi().getRefills(p, itemsPerPage, startDate, endDate, {
					withCredentials: true
				});
				refills.push(...(resp.data.refills ?? []));
			}
		}
	}

	onMount(() => {
		reloadItems();
	});

	let types = [RefillType.RefillCash, RefillType.RefillCard, RefillType.RefillHelloAsso, RefillType.RefillOther];
</script>

<div class="w-full p-5 mt-4 justify-center">
	<div class="flex flex-row justify-between p-5">
		<div class="flex flex-row items-center bg-blue-200 p-4 rounded-lg gap-3">
			<h1 class="text-md font-semibold">Jour :</h1>
			<button
				class="p-2 rounded-full bg-blue-300 hover:bg-blue-400 transition shadow"
				on:click={() => {
				let d = new Date(today);
				d.setDate(d.getDate() - 1);
				today = d;
				let s = time2Utc(d.getTime() / 1000);
				startDate = s.toString();
				endDate = (s + 24 * 60 * 60).toString();
				reloadItems();
				}}
			>
				◀
			</button>
			<input
				type="date"
				value={today.toLocaleString('default', { year: 'numeric' }) +
					'-' +
					today.toLocaleString('default', { month: '2-digit' }) +
					'-' +
					today.toLocaleString('default', { day: '2-digit' })}

				on:change={(e) => {
				// @ts-ignore
				let s = time2Utc(new Date(e.target.value).getTime() / 1000);
				today = new Date(e.target.value);
				startDate = s.toString();
				endDate = (s + 24 * 60 * 60).toString();
				reloadItems();
				}}
				class="rounded-md bg-white px-3 py-2 text-blue-900 font-medium shadow-sm 
					focus:outline-none focus:ring-2 focus:ring-blue-400 focus:border-blue-500"
			/>

			<button
				class="p-2 rounded-full bg-blue-300 hover:bg-blue-400 transition shadow"
				on:click={() => {
				let d = new Date(today);
				d.setDate(d.getDate() + 1);
				today = d;
				let s = time2Utc(d.getTime() / 1000);
				startDate = s.toString();
				endDate = (s + 24 * 60 * 60).toString();
				reloadItems();
				}}
			>
				▶
			</button>
		</div>

	</div>
	<div class="grid grid-cols-1 lg:grid-cols-2 w-full min-w-full mt-5 gap-4 justify-center">
		{#each types as t}
			<div class="w-full flex flex-col bg-blue-200 items-center rounded-lg">
				<h1 class="text-3xl font-semibold p-5">
					Recharges
					{#if t == RefillType.RefillCash}
						en espèces
					{:else if t == RefillType.RefillCard}
						par carte
                    {:else if t == RefillType.RefillHelloAsso}
						par HelloAsso
					{:else}
						autres
					{/if}
				</h1>

				<div class="w-full p-1 bg-blue-300" />

				<div class="flex flex-col items-center bg-blue-200 rounded-t-md p-5">
					<table>
						<thead>
							<tr>
								<th class="px-4">Nombre de recharges total pour ce type</th>
								<th class="px-4">Nombre de recharges réussies pour ce type</th>
								<th class="px-4">Total</th>
							</tr>
						</thead>
						<tbody>
							<tr>
								<td class="px-4">
									{refills.filter((refill) => refill.type == t).length}
									</td>
								<td class="px-4">
									{refills.filter((refill) => refill.type == t)
									.filter((refill) => refill.canceled_by == undefined).length}
									</td>
								<td class="px-4">
									{formatPrice(
										refills
											.filter((refill) => refill.type == t)
											.filter((refill) => refill.canceled_by == undefined)
											.reduce((acc, refill) => acc + refill.amount, 0)
									)}
								</td>
							</tr>
						</tbody>
					</table>
				</div>

				<div class="flex flex-col items-center bg-blue-200 rounded-md p-5">
					<table>
						<thead>
							<tr>
								<th class="px-4">Compte</th>
                                <!-- HelloAsso refills are validated by the user and cannot be canceled -->
                                {#if t != RefillType.RefillHelloAsso}
                                    <th class="px-4">Opérateur</th>
                                    <th class="px-4">Annulé par</th>
                                {/if}
								<th class="px-4">Montant</th>
								<th class="px-4">Date</th>
							</tr>
						</thead>
						<tbody>
							{#each refills as refill, index (refill.id)}
								{#if refill.type == t}
									<tr>
										<td class="px-4">{refill.account_name}</td>
                                        {#if t != RefillType.RefillHelloAsso}
                                            <td class="px-4">{refill.issued_by_name}</td>
                                            <td class="px-4">{refill.canceled_by_name ?? ''}</td>
                                        {/if}
										<td class="px-4">{formatPrice(refill.amount)}</td>
										<td class="px-4">{formatDateTime(refill.issued_at)}</td>
									</tr>
								{/if}
							{/each}
						</tbody>
					</table>
				</div>
			</div>
		{/each}
	</div>
</div>

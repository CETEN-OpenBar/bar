<script lang="ts">
	import { onMount } from 'svelte';
	import type { Starring } from '$lib/api';
	import { starsApi } from '$lib/requests/requests';

	let stars: Starring[] = [];
	let displayed: Starring[] = [];
	let start_date = '';
	let end_date = '';
	let nameFilter = '';
	let page = 0;
	let maxPage = 0;
	let stars_per_page = 15;
	const nextPage = () => {
		if (page < maxPage) {
			page++;
			fetchStars();
		}
	};
	const prevPage = () => {
		if (page > 0) {
			page--;
			fetchStars();
		}
	};

	async function fetchStars() {
		const response = await starsApi().getStarrings(page, stars_per_page, start_date, end_date, {
			withCredentials: true
		});
		const data = response.data;
		stars = Array.isArray(data.stars) ? data.stars : [];
		page = data.page;
		maxPage = data.max_page;
		stars_per_page = data.limit;
		applyFilter();
	}

	function applyFilter() {
		const term = nameFilter.trim().toLowerCase();
		displayed = term ? stars.filter((s) => s.account_name.toLowerCase().includes(term)) : stars;
	}

	onMount(fetchStars);
</script>

<div class="mb-4 flex space-x-4">
	<label>
		Date début:
		<input
			type="date"
			bind:value={start_date}
			on:change={fetchStars}
			class="border px-2 py-1 rounded"
		/>
	</label>
	<label>
		Date fin:
		<input
			type="date"
			bind:value={end_date}
			on:change={fetchStars}
			class="border px-2 py-1 rounded"
		/>
	</label>
	<label>
		Compte:
		<input
			type="text"
			placeholder="Nom"
			bind:value={nameFilter}
			on:input={applyFilter}
			class="border px-2 py-1 rounded"
		/>
	</label>
</div>

<table class="min-w-full border-collapse">
	<thead class="bg-gray-100">
		<tr>
			<th class="border p-2">Date</th>
			<th class="border p-2">Compte</th>
			<th class="border p-2">Émetteur</th>
			<th class="border p-2">Montant</th>
			<th class="border p-2">Type</th>
			<th class="border p-2">État</th>
		</tr>
	</thead>
	<tbody>
		{#each displayed as star}
			<tr class="hover:bg-gray-50">
				<td class="border p-2">{new Date(star.issued_at).toLocaleString()}</td>
				<td class="border p-2">{star.account_name}</td>
				<td class="border p-2">{star.issued_by_name}</td>
				<td class="border p-2">{star.amount}</td>
				<td class="border p-2">{star.type}</td>
				<td class="border p-2">{star.state}</td>
			</tr>
		{/each}
	</tbody>
</table>

<!-- Pagination -->
<div
	class="px-6 py-4 grid gap-3 md:flex md:justify-between md:items-center border-t border-gray-200 dark:border-gray-700"
>
	<div>
		<p class="text-sm text-gray-600 dark:text-gray-400">
			<span class="font-semibold text-gray-800 dark:text-gray-200">{stars.length}</span> résultats
		</p>
	</div>
	<div class="inline-flex gap-x-2">
		<button
			type="button"
			class="py-2 px-3 inline-flex justify-center items-center gap-2 rounded-md border font-medium bg-white text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-600 transition-all text-sm dark:bg-slate-900 dark:hover:bg-slate-800 dark:border-gray-700 dark:text-gray-400 dark:hover:text-white dark:focus:ring-offset-gray-800"
			on:click={prevPage}
			disabled={page === 0}
		>
			Précédent
		</button>
		<p class="text-sm self-center text-gray-600 dark:text-gray-400">Page {page} / {maxPage}</p>
		<button
			type="button"
			class="py-2 px-3 inline-flex justify-center items-center gap-2 rounded-md border font-medium bg-white text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-600 transition-all text-sm dark:bg-slate-900 dark:hover:bg-slate-800 dark:border-gray-700 dark:text-gray-400 dark:hover:text-white dark:focus:ring-offset-gray-800"
			on:click={nextPage}
			disabled={page === maxPage}
		>
			Suivant
		</button>
	</div>
</div>

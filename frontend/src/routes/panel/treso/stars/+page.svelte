<script lang="ts">
	import { onMount } from 'svelte';
	import type { Starring } from '$lib/api';
	import { starsApi } from '$lib/requests/requests';
	import 'iconify-icon';

	let stars: Starring[] = [];
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
		const response = await starsApi().getStarrings(
			page,
			stars_per_page,
			nameFilter.trim(),
			start_date,
			end_date,
			{ withCredentials: true }
		);
		const data = response.data;
		stars = Array.isArray(data.stars) ? data.stars : [];
		page = data.page;
		maxPage = data.max_page;
		stars_per_page = data.limit;
	}

	onMount(fetchStars);
</script>

<div class="flex-grow grid grid-cols-1 grid-rows-[auto_1fr_auto] bg-gray-50 dark:bg-gray-900">
	<div class="m-3 p-2">
		<div class="flex flex-wrap items-center gap-6">
			<div class="flex items-center gap-3">
				<span class="text-sm font-medium text-gray-700 dark:text-gray-300">Du:</span>
				<input
					type="date"
					bind:value={start_date}
					on:change={fetchStars}
					class="px-3 py-1.5 text-sm bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:text-white"
				/>
			</div>
			<div class="flex items-center gap-3">
				<span class="text-sm font-medium text-gray-700 dark:text-gray-300">Au:</span>
				<input
					type="date"
					bind:value={end_date}
					on:change={fetchStars}
					class="px-3 py-1.5 text-sm bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:text-white"
				/>
			</div>
			<div class="flex items-center gap-3">
				<span class="text-sm font-medium text-gray-700 dark:text-gray-300">Compte:</span>
				<input
					type="text"
					placeholder="Rechercher par nom..."
					bind:value={nameFilter}
					on:input={() => {
						page = 0;
						fetchStars();
					}}
					class="px-3 py-1.5 text-sm bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:text-white placeholder-gray-400 dark:placeholder-gray-500"
				/>
			</div>
		</div>
	</div>

	<div class="flex-grow w-full overflow-x-auto">
		<div class="min-w-full bg-white dark:bg-gray-800 rounded-lg shadow-sm">
			<div class="grid grid-cols-[1fr_1.5fr_1fr_0.75fr_0.5fr_0.5fr] bg-gray-50 dark:bg-gray-700 divide-x divide-gray-200 dark:divide-gray-700" >
				<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                    Date
                </th>
				<th	class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                    Compte
                </th>
				<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"> 
                    Émetteur
                </th>
				<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"> 
                    Montant
                </th>
				<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider opacity-60"> 
                    Type
                </th>
				<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider opacity-60"> 
                    État
                </th>
			</div>
			<div class="divide-y divide-gray-200 dark:divide-gray-700">
				{#each stars as star}
					<div class="grid grid-cols-[1fr_1.5fr_1fr_0.75fr_0.5fr_0.5fr] divide-x divide-gray-200 dark:divide-gray-700">
						<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
                            {new Date(star.issued_at * 1000).toLocaleDateString('fr-FR', {
								day: 'numeric',
								month: 'long',
								year: 'numeric',
								hour: '2-digit',
								minute: '2-digit'
							})}
                        </td>
						<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
                            {star.account_name}
                        </td>
						<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
                            {star.issued_by_name}
                        </td>
						<td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-gray-300">
                            {star.amount}
							<iconify-icon icon="mdi:star" class="text-yellow-500" />
						</td>
						<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300 opacity-60">
                            {star.type}
                        </td>
						<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300 opacity-60">
                            {star.state}
                        </td>
					</div>
				{/each}
			</div>
		</div>
	</div>

	<!-- Pagination -->
	<div class="sticky bottom-0 left-0 right-0 mt-4 px-6 py-4 bg-white dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700 flex flex-col sm:flex-row justify-between items-center gap-4">
		<div>
			<p class="text-sm text-gray-600 dark:text-gray-400">
				<span class="font-semibold text-gray-800 dark:text-gray-200">{stars.length}</span> résultats
			</p>
		</div>
		<div class="flex items-center gap-x-4">
			<button
				type="button"
				class="py-2 px-4 inline-flex justify-center items-center gap-2 rounded-md border font-medium bg-white text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-600 transition-all text-sm disabled:opacity-50 disabled:cursor-not-allowed dark:bg-gray-700 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-600 dark:focus:ring-offset-gray-800"
				on:click={prevPage}
				disabled={page === 0}
			>
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M15 19l-7-7 7-7"
					/>
				</svg>
				Précédent
			</button>
			<p class="text-sm font-medium text-gray-700 dark:text-gray-300">
				Page <span class="font-bold">{page}</span> sur <span class="font-bold">{maxPage}</span>
			</p>
			<button
				type="button"
				class="py-2 px-4 inline-flex justify-center items-center gap-2 rounded-md border font-medium bg-white text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-600 transition-all text-sm disabled:opacity-50 disabled:cursor-not-allowed dark:bg-gray-700 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-600 dark:focus:ring-offset-gray-800"
				on:click={nextPage}
				disabled={page === maxPage}
			>
				Suivant
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
				</svg>
			</button>
		</div>
	</div>
</div>

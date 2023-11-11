<script lang="ts">
	import { formatDateTime, formatPrice, time2Utc } from '$lib/utils';
	import type { CashMovement } from '$lib/api';
	import { cashMovementsApi } from '$lib/requests/requests';
	import { onMount } from 'svelte';
	import { stringify } from 'postcss';

	let cashMovements: CashMovement[] = [];
	let page: number = 0;
	let maxPage: number = 0;
	let nextPage = () => {
		if (page < maxPage) {
			page++;
			reloadCashMovements();
		}
	};
	let prevPage = () => {
		if (page > 0) {
			page--;
			reloadCashMovements();
		}
	};
	let limit = 30;
	let search: string | undefined = undefined;

	async function reloadCashMovements() {
		let resp = await cashMovementsApi().getCashMovements(page, limit, search, {
			withCredentials: true
		});
		cashMovements = resp.data.cash_movements ?? [];
		page = resp.data.page ?? 0;
		maxPage = resp.data.max_page ?? 0;
	}

	onMount(async () => {
		await reloadCashMovements();
	});
</script>

<div class="w-full flex flex-col items-center">
	<div class="flex flex-col p-5 gap-3">
		<div class="flex flex-row">
			<input
				class="rounded-l-md bg-slate-200 p-4"
				placeholder="Filtrer par raisons"
				type="text"
				list="items"
				on:keyup={async (e) => {
					// @ts-ignore
					search = e.target.value;
					await reloadCashMovements();
				}}
			/>
			<button class="rounded-r-lg bg-slate-200 p-4"> &#x1F50D; </button>
		</div>
	</div>
	<div
		class="flex flex-col items-center bg-blue-200 rounded-md p-5 overflow-auto"
		style="max-height: 50vh"
	>
		<div class="flex flex-col items-center">
			<h1 class="text-lg font-semibold">Etat de la caisse</h1>
		</div>
		<div class="w-full p-1 bg-blue-300 mt-4 mb-4" />

		<table>
			<thead>
				<tr>
					<th class="px-4">Responsable</th>
					<th class="px-4">Date</th>
					<th class="px-4">Ancienne valeur</th>
					<th class="px-4">Nouvelle valeur</th>
					<th class="px-4">Ecart</th>
					<th class="px-4">Raison</th>
				</tr>
			</thead>
			<tbody>
				{#each cashMovements as t}
					<tr>
						<td class="px-4">{t.created_by_name}</td>
						<td class="px-4">{formatDateTime(time2Utc(t.created_at))}</td>
						<td class="px-4">{formatPrice(t.old_amount)}</td>
						<td class="px-4">{formatPrice(t.amount)}</td>
						<td class="px-4">{formatPrice(t.amount - t.old_amount)}</td>
						<td class="px-4 max-w-3xl">{t.reason}</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
	<div class="flex flex-row mt-5 gap-3">
		<button
			type="button"
			class="py-2 px-3 inline-flex justify-center items-center gap-2 rounded-md border font-medium bg-white text-gray-700 shadow-sm align-middle hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-white focus:ring-blue-600 transition-all text-sm dark:bg-slate-900 dark:hover:bg-slate-800 dark:border-gray-700 dark:text-gray-400 dark:hover:text-white dark:focus:ring-offset-gray-800"
			on:click={prevPage}
		>
			<svg
				class="w-3 h-3"
				xmlns="http://www.w3.org/2000/svg"
				width="16"
				height="16"
				fill="currentColor"
				viewBox="0 0 16 16"
			>
				<path
					fill-rule="evenodd"
					d="M11.354 1.646a.5.5 0 0 1 0 .708L5.707 8l5.647 5.646a.5.5 0 0 1-.708.708l-6-6a.5.5 0 0 1 0-.708l6-6a.5.5 0 0 1 .708 0z"
				/>
			</svg>
			Précédent
		</button>

		<p class="text-sm self-center text-gray-600 dark:text-gray-400">
			Page {page} / {maxPage}
		</p>

		<button
			type="button"
			class="py-2 px-3 inline-flex justify-center items-center gap-2 rounded-md border font-medium bg-white text-gray-700 shadow-sm align-middle hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-white focus:ring-blue-600 transition-all text-sm dark:bg-slate-900 dark:hover:bg-slate-800 dark:border-gray-700 dark:text-gray-400 dark:hover:text-white dark:focus:ring-offset-gray-800"
			on:click={nextPage}
		>
			Suivant
			<svg
				class="w-3 h-3"
				xmlns="http://www.w3.org/2000/svg"
				width="16"
				height="16"
				fill="currentColor"
				viewBox="0 0 16 16"
			>
				<path
					fill-rule="evenodd"
					d="M4.646 1.646a.5.5 0 0 1 .708 0l6 6a.5.5 0 0 1 0 .708l-6 6a.5.5 0 0 1-.708-.708L10.293 8 4.646 2.354a.5.5 0 0 1 0-.708z"
				/>
			</svg>
		</button>
	</div>
</div>

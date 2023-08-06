<script lang="ts">
	import type { Category, Item } from '$lib/api';
	import { api } from '$lib/config/config';
	import { categoriesApi, deletedApi } from '$lib/requests/requests';
	import { onMount } from 'svelte';

	// item with category type
	interface ItemWithCategoryName extends Item {
		category_name: string;
	}

	let items: ItemWithCategoryName[] = [];
	let categories: Category[] = [];

	let page = 0;
	let max_page = 0;
	let itemPerPage = 10;

	onMount(() => {
		categoriesApi()
			.getCategories({ withCredentials: true })
			.then((res) => {
				categories = res.data ?? [];
				reloadItems();
			});
	});

	function reloadItems() {
		deletedApi()
			.getDeletedItems(page, itemPerPage, { withCredentials: true })
			.then((res) => {
				items = (res.data.items ?? []).map((item) => {
					const category = categories.find((ct) => ct.id === item.category_id);
					return {
						...item,
						category_name: category?.name ?? 'Catégorie supprimée'
					};
				});

				page = res.data.page;
				itemPerPage = res.data.limit;
				max_page = res.data.max_page;
			});
	}

	function deleteItem(id: string) {
		deletedApi()
			.deleteItem(id, { withCredentials: true })
			.then(() => {
				items = items.filter((ct) => ct.id !== id);
			});
	}

	function restoreItem(id: string) {
		deletedApi()
			.restoreDeletedItem(id, { withCredentials: true })
			.then(() => {
				items = items.filter((ct) => ct.id !== id);
			});
	}
</script>

<!-- Table Section -->
<div class="max-w-[85rem] px-4 py-10 sm:px-6 lg:px-8 lg:py-14 mx-auto">
	<!-- Card -->
	<div class="flex flex-col">
		<div class="-m-1.5 overflow-x-auto">
			<div class="p-1.5 min-w-full inline-block align-middle">
				<div
					class="bg-white border border-gray-200 rounded-xl shadow-sm overflow-hidden dark:bg-slate-900 dark:border-gray-700"
				>
					<!-- Header -->
					<div
						class="px-6 py-4 grid gap-3 md:flex md:justify-between md:items-center border-b border-gray-200 dark:border-gray-700"
					>
						<div>
							<h2 class="text-xl font-semibold text-gray-800 dark:text-gray-200">Produits</h2>
							<p class="text-sm text-gray-600 dark:text-gray-400">Restorer des produits</p>
						</div>
					</div>
					<!-- End Header -->

					<!-- Table -->
					<table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
						<thead class="bg-gray-50 dark:bg-slate-800">
							<tr>
								<th scope="col" class="px-6 py-3 text-left">
									<div class="flex items-center gap-x-2">
										<span
											class="text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
										>
											Nom
										</span>
									</div>
								</th>
								<th scope="col" class="px-6 py-3 text-left">
									<div class="flex items-center gap-x-2">
										<span
											class="text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
										>
											Catégorie
										</span>
									</div>
								</th>
								<th scope="col" class="px-6 py-3 text-left">
									<div class="flex items-center gap-x-2">
										<span
											class="text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
										>
											Image
										</span>
									</div>
								</th>

								<th scope="col" class="px-6 py-3 text-right" />
							</tr>
						</thead>

						<tbody class="divide-y divide-gray-200 dark:divide-gray-700">
							{#each items as item}
								<tr>
									<td class="h-px w-72">
										<div class="px-6 py-3">
											<p class="block text-sm text-gray-500 break-words">{item.name}</p>
										</div>
									</td>
									<td class="h-px w-72">
										<div class="px-6 py-3">
											<p class="block text-sm text-gray-500 break-words">{item.category_name}</p>
										</div>
									</td>
									<td class="h-px w-72">
										<!-- Display a miniature of the image -->
										<div class="px-6 py-3 w-24 relative">
											{#if item.picture_uri != ''}
												<img
													src={api() + item.picture_uri}
													alt="indisponible"
													class="w-full h-full rounded-md object-cover"
												/>
											{/if}
										</div>
									</td>
									<td class="h-px w-px whitespace-nowrap">
										<div class="px-6 py-1.5">
											<button
												class="inline-flex items-center gap-x-1.5 text-sm text-blue-600 decoration-2 hover:underline font-medium"
												on:click={() => restoreItem(item.id)}
											>
												Restorer
											</button>
											<button
												class="inline-flex items-center gap-x-1.5 text-sm text-blue-600 decoration-2 hover:underline font-medium"
												on:click={() => deleteItem(item.id)}
											>
												Supprimer
											</button>
										</div>
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
					<!-- End Table -->

					<!-- Footer -->
					<div
						class="px-6 py-4 grid gap-3 md:flex md:justify-between md:items-center border-t border-gray-200 dark:border-gray-700"
					>
						<div>
							<p class="text-sm text-gray-600 dark:text-gray-400">
								<span class="font-semibold text-gray-800 dark:text-gray-200">{items.length}</span> résultats
							</p>
						</div>

						<div>
							<div class="inline-flex gap-x-2">
								<button
									type="button"
									class="py-2 px-3 inline-flex justify-center items-center gap-2 rounded-md border font-medium bg-white text-gray-700 shadow-sm align-middle hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-white focus:ring-blue-600 transition-all text-sm dark:bg-slate-900 dark:hover:bg-slate-800 dark:border-gray-700 dark:text-gray-400 dark:hover:text-white dark:focus:ring-offset-gray-800"
									on:click={() => {
										if (page > 0) {
											page--;
											reloadItems();
										}
									}}
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
									Page {page+1} / {max_page+1}
								</p>

								<button
									type="button"
									class="py-2 px-3 inline-flex justify-center items-center gap-2 rounded-md border font-medium bg-white text-gray-700 shadow-sm align-middle hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-white focus:ring-blue-600 transition-all text-sm dark:bg-slate-900 dark:hover:bg-slate-800 dark:border-gray-700 dark:text-gray-400 dark:hover:text-white dark:focus:ring-offset-gray-800"
									on:click={() => {
										if (page < max_page) {
											page++;
											reloadItems();
										}
									}}
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
					</div>
					<!-- End Footer -->
				</div>
			</div>
		</div>
	</div>
	<!-- End Card -->
</div>
<!-- End Table Section -->

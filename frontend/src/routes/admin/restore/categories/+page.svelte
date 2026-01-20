<script lang="ts">
	import type { Category, NewCategory } from '$lib/api';
	import { api } from '$lib/config/config';
	import { categoriesApi, deletedApi } from '$lib/requests/requests';
	import { onMount } from 'svelte';
	import PaginationFooter from '$lib/components/PaginationFooter.svelte';

	let categories: Category[] = [];

	let page: number = 0;
	let maxPage: number = 0;
	let nextPage = () => {
		if (page < maxPage) {
			page++;
			reloadCategories();
		}
	};
	let prevPage = () => {
		if (page > 0) {
			page--;
			reloadCategories();
		}
	};
	let categoriesPerPage = 10;

	onMount(() => {
		reloadCategories();
	});

	function reloadCategories() {
		deletedApi()
			.getDeletedCategories(page, categoriesPerPage, { withCredentials: true })
			.then((res) => {
				categories = res.data.categories ?? [];
				page = res.data.page;
				categoriesPerPage = res.data.limit;
				maxPage = res.data.max_page;
			});
	}

	function deleteCategory(id: string) {
		deletedApi()
			.deleteCategory(id, { withCredentials: true })
			.then(() => {
				categories = categories.filter((ct) => ct.id !== id);
			});
	}

	function restoreCategory(id: string) {
		deletedApi()
			.restoreDeletedCategory(id, { withCredentials: true })
			.then(() => {
				categories = categories.filter((ct) => ct.id !== id);
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
							<h2 class="text-xl font-semibold text-gray-800 dark:text-gray-200">Catégories</h2>
							<p class="text-sm text-gray-600 dark:text-gray-400">Restorer des catégories</p>
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
											Image
										</span>
									</div>
								</th>

								<th scope="col" class="px-6 py-3 text-right" />
							</tr>
						</thead>

						<tbody class="divide-y divide-gray-200 dark:divide-gray-700">
							{#each categories as category}
								<tr>
									<td class="h-px w-72">
										<div class="px-6 py-3">
											<p class="block text-sm text-gray-500 break-words">{category.name}</p>
										</div>
									</td>
									<td class="h-px w-72">
										<!-- Display a miniature of the image -->
										<div class="px-6 py-3 w-24 relative">
											{#if category.picture_uri != ''}
												<img
													src={api() + category.picture_uri}
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
												on:click={() => restoreCategory(category.id)}
											>
												Restorer
											</button>
											<button
												class="inline-flex items-center gap-x-1.5 text-sm text-blue-600 decoration-2 hover:underline font-medium"
												on:click={() => deleteCategory(category.id)}
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
					<PaginationFooter
						{page}
						{maxPage}
						resultsCount={categories.length}
						zeroBased={true}
						on:prevPage={prevPage}
						on:nextPage={nextPage}
					/>
					<!-- End Footer -->
				</div>
			</div>
		</div>
	</div>
	<!-- End Card -->
</div>
<!-- End Table Section -->

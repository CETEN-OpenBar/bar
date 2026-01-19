<script lang="ts">
	import type { Category, NewCategory } from '$lib/api';
	import { api } from '$lib/config/config';
	import { categoriesApi } from '$lib/requests/requests';
	import ConfirmationPopup from '$lib/components/confirmationPopup.svelte';
	import { onMount } from 'svelte';
	import PaginationFooter from '$lib/components/PaginationFooter.svelte';

	let categories: Category[] = [];
	let newCategory: NewCategory = {
		name: '',
		picture: '',
		position: 0
	};

	let page: number = 1;
	let maxPage: number = 1;
	let categories_per_page = 10;

	let nextPage = () => {
		if (page < maxPage) {
			page++;
			reloadCategories();
		}
	};
	let prevPage = () => {
		if (page > 1) {
			page--;
			reloadCategories();
		}
	};
	let handlePageInput = () => {
		if (page < 1) {
			page = 1;
		} else if (page > maxPage) {
			page = maxPage;
		}
		reloadCategories();
	};

	let deletingCategory: boolean = false;
	let confirmationMessage: string | undefined = undefined;
	let deleteCategoryCallback: VoidFunction = () => {};

	onMount(() => {
		reloadCategories();
	});

	function reloadCategories() {
		categoriesApi()
			.getCategories(true, { withCredentials: true })
			.then((res) => {
				const allCategories = res.data ?? [];
				maxPage = Math.max(1, Math.ceil(allCategories.length / categories_per_page));
				if (page > maxPage) page = maxPage;
				const start = (page - 1) * categories_per_page;
				const end = start + categories_per_page;
				categories = allCategories.slice(start, end);
			});
	}

	function createNewCategory() {
		if (!newCategory) return;
		categoriesApi()
			.postCategory(newCategory, { withCredentials: true })
			.then((res) => {
				reloadCategories();
				newCategory = {
					name: '',
					picture: '',
					position: 0
				};
			});
	}

	function renameCategory(id: string, newName: string) {
		categoriesApi()
			.patchCategory(id, { name: newName }, { withCredentials: true })
			.then((res) => {
				categories = categories.map((ct) => {
					if (ct.id === id) {
						ct.name = newName;
					}
					return ct;
				});
			});
	}

	function toggleHidden(id: string, state: boolean) {
		categoriesApi()
			.patchCategory(id, { hidden: state }, { withCredentials: true })
			.then((res) => {
				categories = categories.map((ct) => {
					if (ct.id === id) {
						ct.hidden = state;
					}
					return ct;
				});
			});
	}

	function toggleSpecialPrice(id: string, state: boolean) {
		categoriesApi()
			.patchCategory(id, { special_price: state }, { withCredentials: true })
			.then((res) => {
				categories = categories.map((ct) => {
					if (ct.id === id) {
						ct.special_price = state;
					}
					return ct;
				});
			});
	}

	function reuploadCategoryPicture(id: string, file: File) {
		file2Base64(file).then((base64) => {
			base64 = base64.replace('data:', '').replace(/^.+,/, '');
			categoriesApi()
				.patchCategory(id, { picture: base64 }, { withCredentials: true })
				.then((res) => {
					categories = categories.map((ct) => {
						if (ct.id === id) {
							ct.picture_uri = res.data.picture_uri + '?' + Math.random();
						}
						return ct;
					});
				});
		});
	}

	function deleteCategory(id: string) {
		categoriesApi()
			.markDeleteCategory(id, { withCredentials: true })
			.then(() => {
				reloadCategories();
			});
	}

	const file2Base64 = (file: File): Promise<string> => {
		return new Promise<string>((resolve, reject) => {
			const reader = new FileReader();
			reader.readAsDataURL(file);
			reader.onload = () => resolve(reader.result?.toString() || '');
			reader.onerror = (error) => reject(error);
		});
	};
</script>

<!-- Popup -->
<div
	id="hs-modal-new-image"
	class="hs-overlay hidden w-full h-full fixed top-0 left-0 z-[60] overflow-x-hidden overflow-y-auto"
>
	<div
		class="hs-overlay-open:mt-7 hs-overlay-open:opacity-100 hs-overlay-open:duration-500 mt-0 opacity-0 ease-out transition-all sm:max-w-lg sm:w-full m-3 sm:mx-auto"
	>
		<div
			class="bg-white border border-gray-200 rounded-xl shadow-sm dark:bg-gray-800 dark:border-gray-700"
		>
			<div class="p-4 sm:p-7">
				<div class="text-center">
					<h2 class="block text-2xl font-bold text-gray-800 dark:text-gray-200">
						Ajouter une categorie
					</h2>
				</div>

				<div class="mt-5">
					<!-- Form -->
					<div class="grid gap-y-4">
						<!-- Form Group -->
						<div>
							<!-- name -->
							<label for="name" class="block text-sm mb-2 dark:text-white">Nom</label>
							<div class="relative">
								<input
									type="text"
									id="name"
									name="name"
									placeholder="Nom de la categorie"
									class="py-3 px-4 block w-full border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
									required
									aria-describedby="text-error"
									bind:value={newCategory.name}
								/>
							</div>

							<label for="image" class="block text-sm mb-2 dark:text-white">Image</label>
							<div class="relative">
								<input
									type="file"
									id="image"
									name="image"
									accept=".jpg, .jpeg, .png"
									class="py-3 px-4 block w-full border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
									required
									aria-describedby="text-error"
									on:change={(e) => {
										// @ts-ignore
										let file = e.target?.files[0];
										file2Base64(file).then((res) => {
											res = res.replace('data:', '').replace(/^.+,/, '');
											newCategory.picture = res;
										});
									}}
								/>
							</div>

							<button
								type="submit"
								class="mt-4 py-3 px-4 inline-flex justify-center items-center gap-2 rounded-md border border-transparent font-semibold bg-blue-500 text-white hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-all text-sm dark:focus:ring-offset-gray-800"
								on:click={() => createNewCategory()}
								data-hs-overlay="#hs-modal-new-image">Creer</button
							>
						</div>
					</div>
					<!-- End Form -->
				</div>
			</div>
		</div>
	</div>
</div>

{#if deletingCategory}
	<ConfirmationPopup
		message={confirmationMessage}
		confirm_text="Supprimer"
		cancel_callback={() => {
			deletingCategory = false;
		}} 
		confirm_callback={deleteCategoryCallback}
	/>
{/if}

<div class="h-[calc(100vh-69px)] grid grid-cols-1 grid-rows-[auto_1fr_100px] sm:grid-rows-[auto_1fr_80px] bg-gray-50 dark:bg-gray-900">
	<div class="m-3 p-2">
		<div class="flex flex-col sm:flex-row sm:flex-wrap sm:items-center gap-4 sm:gap-6">
			<div class="flex flex-col sm:flex-row sm:items-center gap-2 sm:gap-3">
				<button
					class="py-2 px-3 inline-flex justify-center items-center gap-2 rounded-md border border-transparent font-semibold bg-blue-500 text-white hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-all text-sm dark:focus:ring-offset-gray-800 w-full sm:w-auto"
					data-hs-overlay="#hs-modal-new-image"
				>
					<svg
						class="w-3 h-3"
						xmlns="http://www.w3.org/2000/svg"
						width="16"
						height="16"
						viewBox="0 0 16 16"
						fill="none"
					>
						<path
							d="M2.63452 7.50001L13.6345 7.5M8.13452 13V2"
							stroke="currentColor"
							stroke-width="2"
							stroke-linecap="round"
						/>
					</svg>
					<span class="sm:hidden">Ajouter</span>
					<span class="hidden sm:inline">Ajouter une categorie</span>
				</button>
			</div>
		</div>
	</div>

	<div class="flex-grow w-full overflow-x-auto overflow-y-visible">
		<!-- Desktop Table View -->
		<div class="hidden min-[800px]:block min-w-full bg-white dark:bg-gray-800 rounded-lg shadow-sm overflow-visible">
			<table class="min-w-full table-fixed divide-y divide-gray-200 dark:divide-gray-700">
				<colgroup>
					<col class="w-[25%]" />
					<col class="w-[25%]" />
					<col class="w-[15%]" />
					<col class="w-[15%]" />
					<col class="w-[20%]" />
				</colgroup>
				<thead class="bg-gray-50 dark:bg-gray-700">
					<tr>
						<th class="px-6 py-4 text-left text-sm font-semibold text-gray-900 dark:text-gray-300 border-r border-gray-200 dark:border-gray-600">
							Nom
						</th>
						<th class="px-6 py-4 text-left text-sm font-semibold text-gray-900 dark:text-gray-300 border-r border-gray-200 dark:border-gray-600">
							Image
						</th>
						<th class="px-6 py-4 text-left text-sm font-semibold text-gray-900 dark:text-gray-300 border-r border-gray-200 dark:border-gray-600">
							Cachee
						</th>
						<th class="px-6 py-4 text-left text-sm font-semibold text-gray-900 dark:text-gray-300 border-r border-gray-200 dark:border-gray-600">
							Prix Special
						</th>
						<th class="px-6 py-4 text-left text-sm font-semibold text-gray-900 dark:text-gray-300">
							Actions
						</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-gray-200 dark:divide-gray-700">
					{#each categories as category}
						<tr>
							<td class="px-6 py-4 text-sm text-gray-900 dark:text-gray-300 border-r border-gray-200 dark:border-gray-600 group relative">
								<input
									type="text"
									class="block w-full text-sm dark:text-white/[.8] p-2 bg-transparent border-none outline-none"
									value={category.name}
									on:input={(e) => {
										// @ts-ignore
										let name = e.target?.value;
										renameCategory(category.id, name);
									}}
								/>
								<iconify-icon 
									icon="mdi:pencil" 
									class="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400 opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none"
								/>
							</td>
							<td class="px-6 py-4 text-sm text-gray-900 dark:text-gray-300 border-r border-gray-200 dark:border-gray-600">
								<div class="w-16 relative">
									<input
										type="file"
										class="absolute w-full h-full opacity-0 cursor-pointer"
										accept=".jpg, .jpeg, .png"
										on:change={(e) => {
											// @ts-ignore
											let file = e.target?.files[0];
											reuploadCategoryPicture(category.id, file);
										}}
									/>
									{#if category.picture_uri != ''}
										<img
											src={api() + category.picture_uri}
											alt="indisponible"
											class="w-full h-full rounded-md object-cover"
										/>
									{:else}
										<div class="w-full h-12 bg-gray-200 dark:bg-gray-600 rounded-md flex items-center justify-center">
											<iconify-icon icon="mdi:image-plus" class="text-gray-400" />
										</div>
									{/if}
								</div>
							</td>
							<td class="px-6 py-4 text-sm text-gray-900 dark:text-gray-300 border-r border-gray-200 dark:border-gray-600">
								<input
									type="checkbox"
									checked={category.hidden}
									class="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
									on:change={() => toggleHidden(category.id, !category.hidden)}
								/>
							</td>
							<td class="px-6 py-4 text-sm text-gray-900 dark:text-gray-300 border-r border-gray-200 dark:border-gray-600">
								<input
									type="checkbox"
									checked={category.special_price}
									class="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
									on:change={() => toggleSpecialPrice(category.id, !category.special_price)}
								/>
							</td>
							<td class="px-6 py-4 text-sm text-gray-900 dark:text-gray-300">
								<button
									class="inline-flex items-center gap-x-1.5 text-sm text-red-600 dark:text-red-400 decoration-2 hover:underline font-medium"
									on:click={() => {
										deleteCategoryCallback = () => {
											deletingCategory = false;
											deleteCategory(category.id);
										};
										confirmationMessage = "Supprimer '" + category.name + "' ?";
										deletingCategory = true;
									}}
								>
									<iconify-icon icon="mdi:delete" width="20" height="20" />
									Supprimer
								</button>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>

		<!-- Mobile Card View -->
		<div class="block min-[800px]:hidden space-y-4 px-2 pb-[110px] sm:pb-0">
			{#if categories.length === 0}
				<div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-8 text-center">
					<p class="text-gray-500 dark:text-gray-400">Aucune categorie trouvee</p>
				</div>
			{:else}
				{#each categories as category}
					<div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-4">
						<div class="flex justify-between items-start mb-3">
							<div class="flex-1">
								<div class="mb-2">
									<label class="text-xs font-medium text-gray-500 dark:text-gray-400 block">Nom</label>
									<input
										type="text"
										class="w-full text-sm dark:text-white/[.8] bg-transparent border border-gray-200 dark:border-gray-600 rounded p-2 focus:border-blue-500 focus:outline-none"
										value={category.name}
										on:input={(e) => {
											// @ts-ignore
											let name = e.target?.value;
											renameCategory(category.id, name);
										}}
									/>
								</div>
								<div class="mb-2">
									<label class="text-xs font-medium text-gray-500 dark:text-gray-400 block">Image</label>
									<div class="w-20 relative mt-1">
										<input
											type="file"
											class="absolute w-full h-full opacity-0 cursor-pointer"
											accept=".jpg, .jpeg, .png"
											on:change={(e) => {
												// @ts-ignore
												let file = e.target?.files[0];
												reuploadCategoryPicture(category.id, file);
											}}
										/>
										{#if category.picture_uri != ''}
											<img
												src={api() + category.picture_uri}
												alt="indisponible"
												class="w-full h-full rounded-md object-cover"
											/>
										{:else}
											<div class="w-full h-12 bg-gray-200 dark:bg-gray-600 rounded-md flex items-center justify-center">
												<iconify-icon icon="mdi:image-plus" class="text-gray-400" />
											</div>
										{/if}
									</div>
								</div>
							</div>
						</div>
						<div class="grid grid-cols-2 gap-2 mb-3">
							<div>
								<label class="text-xs font-medium text-gray-500 dark:text-gray-400 block">Cachee</label>
								<div class="py-2">
									<input
										type="checkbox"
										checked={category.hidden}
										class="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
										on:change={() => toggleHidden(category.id, !category.hidden)}
									/>
								</div>
							</div>
							<div>
								<label class="text-xs font-medium text-gray-500 dark:text-gray-400 block">Prix Special</label>
								<div class="py-2">
									<input
										type="checkbox"
										checked={category.special_price}
										class="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
										on:change={() => toggleSpecialPrice(category.id, !category.special_price)}
									/>
								</div>
							</div>
						</div>
						<div class="flex flex-wrap gap-2">
							<button
								class="flex-1 min-w-0 px-3 py-2 text-xs bg-red-500 text-white rounded hover:bg-red-600 transition-colors flex items-center justify-center gap-1"
								on:click={() => {
									deleteCategoryCallback = () => {
										deletingCategory = false;
										deleteCategory(category.id);
									};
									confirmationMessage = "Supprimer '" + category.name + "' ?";
									deletingCategory = true;
								}}
							>
								<iconify-icon icon="mdi:delete" width="16" height="16" />
								Supprimer
							</button>
						</div>
					</div>
				{/each}
			{/if}
		</div>
	</div>

	<!-- Pagination -->
	<PaginationFooter
		bind:page
		{maxPage}
		resultsCount={categories.length}
		showPageInput={true}
		on:prevPage={prevPage}
		on:nextPage={nextPage}
		on:pageChange={handlePageInput}
	/>
</div>

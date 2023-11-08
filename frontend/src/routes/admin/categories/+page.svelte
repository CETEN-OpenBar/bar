<script lang="ts">
	import type { Category, NewCategory } from '$lib/api';
	import { api } from '$lib/config/config';
	import { categoriesApi } from '$lib/requests/requests';
	import { onMount } from 'svelte';

	let categories: Category[] = [];
	let selectedCategories: Category[] = [];
	let newCategory: NewCategory = {
		name: '',
		picture: '',
		position: 0
	};

	let page = 0;
	let categoriesPerPage = 10;

	onMount(() => {
		categoriesApi()
			.getCategories(true, { withCredentials: true })
			.then((res) => {
				categories = res.data ?? [];
			});
	});

	function createNewCategory() {
		if (!newCategory) return;
		categoriesApi()
			.postCategory(newCategory, { withCredentials: true })
			.then((res) => {
				categories = [...categories, res.data];
			});
	}

	function renameCategory(id: string, newName: string) {
		if (!newCategory) return;
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
		if (!newCategory) return;
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

	function reuploadCategoryPicture(id: string, file: File) {
		if (!newCategory) return;
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
				categories = categories.filter((ct) => ct.id !== id);
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
						Ajouter une catégorie
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
									placeholder="Nom de la catégorie"
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
								data-hs-overlay="#hs-modal-new-image">Créer</button
							>
						</div>
					</div>
					<!-- End Form -->
				</div>
			</div>
		</div>
	</div>
</div>

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
							<p class="text-sm text-gray-600 dark:text-gray-400">Ajouter des catégories</p>
						</div>

						<div>
							<div class="inline-flex gap-x-2">
								<button
									class="py-2 px-3 inline-flex justify-center items-center gap-2 rounded-md border border-transparent font-semibold bg-blue-500 text-white hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-all text-sm dark:focus:ring-offset-gray-800"
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
									Ajouter une catégorie
								</button>
							</div>
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
								<th scope="col" class="px-6 py-3 text-left">
									<div class="flex items-center gap-x-2">
										<span
											class="text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
										>
											Cachée
										</span>
									</div>
								</th>

								<th scope="col" class="px-6 py-3 text-right" />
							</tr>
						</thead>

						<tbody class="divide-y divide-gray-200 dark:divide-gray-700">
							{#each categories.slice(page * categoriesPerPage, (page + 1) * categoriesPerPage) as category}
								<tr>
									<td class="h-px w-72">
										<div class="px-6 py-3">
											<!-- <p class="block text-sm text-gray-500 break-words">{category.name}</p> -->

											<!-- editable p -->

											<input
												type="text"
												class="block text-sm dark:text-white/[.8] break-words p-2 bg-transparent"
												value={category.name}
												on:input={(e) => {
													// @ts-ignore
													let name = e.target?.value;
													renameCategory(category.id, name);
												}}
											/>
										</div>
									</td>
									<td class="h-px w-72">
										<!-- Display a miniature of the image -->
										<div class="px-6 py-3 w-24 relative">
											<!-- <img
												src={api() + category.picture_uri}
												alt="indisponible"
												class="w-full h-full rounded-md object-cover"
											/> -->

											<!-- input in front of the image to click & reupload -->
											<input
												type="file"
												class="absolute w-[50%] h-[70%] opacity-0 cursor-pointer"
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
											{/if}
										</div>
									</td>
									<td class="h-px w-px whitespace-nowrap">
										<div class="px-6 py-1.5">
											<input
												type="checkbox"
												checked={category.hidden}
												class="inline-flex items-center gap-x-1.5 text-sm text-blue-600 decoration-2 hover:underline font-medium"
												on:change={() => toggleHidden(category.id, !category.hidden)}
											/>
										</div>
									</td>
									<td class="h-px w-px whitespace-nowrap">
										<div class="px-6 py-1.5">
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
					<div
						class="px-6 py-4 grid gap-3 md:flex md:justify-between md:items-center border-t border-gray-200 dark:border-gray-700"
					>
						<div>
							<p class="text-sm text-gray-600 dark:text-gray-400">
								<span class="font-semibold text-gray-800 dark:text-gray-200"
									>{categories.length}</span
								> résultats
							</p>
						</div>

						<div>
							<div class="inline-flex gap-x-2">
								<button
									type="button"
									class="py-2 px-3 inline-flex justify-center items-center gap-2 rounded-md border font-medium bg-white text-gray-700 shadow-sm align-middle hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-white focus:ring-blue-600 transition-all text-sm dark:bg-slate-900 dark:hover:bg-slate-800 dark:border-gray-700 dark:text-gray-400 dark:hover:text-white dark:focus:ring-offset-gray-800"
									on:click={() => {
										if (page > 0) page--;
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
									Page {page + 1} / {Math.ceil(categories.length / categoriesPerPage)}
								</p>

								<button
									type="button"
									class="py-2 px-3 inline-flex justify-center items-center gap-2 rounded-md border font-medium bg-white text-gray-700 shadow-sm align-middle hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-white focus:ring-blue-600 transition-all text-sm dark:bg-slate-900 dark:hover:bg-slate-800 dark:border-gray-700 dark:text-gray-400 dark:hover:text-white dark:focus:ring-offset-gray-800"
									on:click={() => {
										if (page < Math.ceil(categories.length / categoriesPerPage) - 1) page++;
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

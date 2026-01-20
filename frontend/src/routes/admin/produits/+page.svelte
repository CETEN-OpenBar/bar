<script lang="ts">
	import type {
		Category,
		Item,
		NewItem,
		ItemPrices,
		UpdateItem,
		AccountPriceRole,
		ItemState,
		Fournisseur
	} from '$lib/api';
	import ConfirmationPopup from '$lib/components/confirmationPopup.svelte';
	import { api } from '$lib/config/config';
	import { categoriesApi, itemsApi } from '$lib/requests/requests';
	import { formatPrice, parsePrice, fournisseurIterator } from '$lib/utils';
	import { onMount } from 'svelte';
	import PaginationFooter from '$lib/components/PaginationFooter.svelte';
	import PriceCell from '$lib/components/admin/priceCell.svelte';
	import FournisseurCell from '$lib/components/admin/fournisseurCell.svelte';
	import StockCell from '$lib/components/admin/stockCell.svelte';

	interface NewItemWithCategory extends NewItem {
		category_id: string;
	}

	let categories: Category[] = [];
	let selectedItem: Item | null = null;
	let items: Item[] = [];
	let newItem: NewItemWithCategory = {
		name: '',
		picture: '',
		prices: {
			externe: 0,
			ceten: 0,
			staff_bar: 0,
			coutant: 0,
			privilegies: 0,
			menu: 0,
		} as ItemPrices,
		amount_left: 0,
		buy_limit: undefined,
		optimal_amount: 0,
		state: 'buyable',
		category_id: ''
	};
	let newItemPriceRole: AccountPriceRole = 'ceten';
	let promoItemPriceRole: AccountPriceRole = 'ceten';

	let page: number = 1;
	let maxPage: number = 0;
	let handlePageInput = () => {
		if (page < 1) {
			page = 1;
		} else if (page > maxPage) {
			page = maxPage;
		}
		reloadItems();
	};
	let nextPage = () => {
		if (page < maxPage) {
			page++;
			reloadItems();
		}
	};
	let prevPage = () => {
		if (page > 1) {
			page--;
			reloadItems();
		}
	};
	let itemsPerPage = 7;

	let searchState: ItemState | undefined = undefined;
	let searchCategory: string | undefined = undefined;
	let searchName: string | undefined = undefined;
	let searchFournisseur: Fournisseur | undefined = undefined;

	let deletingItem: boolean = false;
	let confirmationMessage: string | undefined = undefined;
	let deleteItemCallback: VoidFunction = () => {};

	let filtersOpen: boolean = false;

	onMount(() => {
		categoriesApi()
			.getCategories(true, { withCredentials: true })
			.then((res) => {
				categories = res.data ?? [];
			});

		reloadItems();
	});

	function reloadItems() {
		itemsApi()
			.getAllItems(page, itemsPerPage, searchState, searchCategory, searchName, searchFournisseur, undefined, {
				withCredentials: true
			})
			.then((res) => {
				maxPage = res.data.max_page ?? 1;
				page = res.data.page ?? 1;
				itemsPerPage = res.data.limit ?? 7;
				items = res.data.items ?? [];
			});
	}

	function createNewItem() {
		if (!newItem) return;
		itemsApi()
			.postItem(newItem.category_id, newItem, { withCredentials: true })
			.then((res) => {
				items = [...items, res.data];
			});
	}

	function editItem(id: string, item: UpdateItem, category_id: string) {
		if (!newItem) return;
		itemsApi()
			.patchItem(category_id, id, item, { withCredentials: true })
			.then((res) => {
				items = items.map((it) => {
					if (it.id === id) {
						it = res.data;
					}
					return it;
				});
				selectedItem = res.data;
			});
	}

	function reuploadItemPicture(id: string, file: File, category_id: string) {
		if (!newItem) return;
		file2Base64(file).then((base64) => {
			base64 = base64.replace('data:', '').replace(/^.+,/, '');
			itemsApi()
				.patchItem(category_id, id, { picture: base64 }, { withCredentials: true })
				.then((res) => {
					items = items.map((ct) => {
						if (ct.id === id) {
							ct.picture_uri = res.data.picture_uri + '?' + Math.random();
						}
						return ct;
					});
				});
		});
	}

	function deleteItem(id: string, category_id: string) {
		itemsApi()
			.markDeleteItem(category_id, id, { withCredentials: true })
			.then(() => {
				items = items.filter((ct) => ct.id !== id);
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
	id="hs-modal-new-item"
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
						Ajouter un produit
					</h2>
				</div>

				<div class="mt-5">
					<!-- Form -->
					<div class="grid gap-y-4">
						<!-- Form Group -->
						<div class="flex flex-col gap-2">
							<!-- name -->
							<label for="name" class="block text-sm mb-2 dark:text-white">Nom</label>
							<div class="relative">
								<input
									type="text"
									id="name"
									name="name"
									placeholder="Nom du produit"
									class="py-3 px-4 block w-full border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
									required
									aria-describedby="text-error"
									bind:value={newItem.name}
								/>
							</div>

							<label for="category" class="block text-sm mb-2 dark:text-white">Catégorie</label>
							<div class="relative">
								<select
									id="category"
									name="category"
									class="py-3 px-4 block w-full border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
									required
									aria-describedby="text-error"
									bind:value={newItem.category_id}
								>
									{#each categories as category}
										<option value={category.id}>{category.name}</option>
									{/each}
								</select>
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
											newItem.picture = res;
										});
									}}
								/>
							</div>

							<!-- Price role selector -->
							<label for="role" class="block text-sm mb-2 dark:text-white">
								<div class="relative">
									<select
										id="role"
										name="role"
										class="py-3 px-4 block w-full border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
										required
										aria-describedby="text-error"
										bind:value={newItemPriceRole}
										on:change={() => {
											// resets the value of the price input
											// Remove the value of elems with id "price"
											let elems = document.querySelectorAll('[id=price-new]');
											elems.forEach((elem) => {
												// @ts-ignore
												elem.value = '';
												// @ts-ignore
												elem.placeholder = formatPrice(newItem.prices[newItemPriceRole]);
											});
										}}
									>
										<option value="coutant">Prix réel coutant</option>
										<option value="externe">Prix réel externe</option>
										<option value="ceten">Prix réel ceten</option>
										<option value="staff_bar">Prix réel staff</option>
										<option value="privilegies">Prix réel privilégiés</option>
										<option value="menu">Prix réel menu</option>
									</select>
								</div>

								<label for="price" class="block text-sm mb-2 dark:text-white"
									>Prix ({newItemPriceRole})</label
								>
								<div class="relative">
									<input
										type="number"
										id="price-new"
										name="price"
										placeholder="Prix du produit"
										class="py-3 px-4 block w-[95%] border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
										required
										aria-describedby="text-error"
										on:input={(e) => {
											// @ts-ignore
											newItem.prices[newItemPriceRole] = parsePrice(e.target?.value);
										}}
									/>

									<span class="absolute top-3 right-4 text-sm text-gray-400"> € </span>
								</div>

								<label for="description" class="block text-sm mb-2 dark:text-white"
									>Stock disponible</label
								>
								<div class="relative">
									<input
										type="number"
										id="stock"
										name="stock"
										placeholder="Stock disponible"
										class="py-3 px-4 block w-full border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
										required
										aria-describedby="text-error"
										bind:value={newItem.amount_left}
									/>
								</div>

								<label for="description" class="block text-sm mb-2 dark:text-white"
									>Max par commande</label
								>
								<div class="relative">
									<input
										type="number"
										id="max"
										name="max"
										placeholder="Pas de max par commande"
										class="py-3 px-4 block w-full border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
										required
										aria-describedby="text-error"
										bind:value={newItem.buy_limit}
										on:input={(e) => {
											// @ts-ignore
											if (e.target?.value === '') {
												newItem.buy_limit = undefined;
												return;
											}
											// @ts-ignore
											let buy_limit = parseInt(e.target?.value);
											newItem.buy_limit = buy_limit;
										}}
									/>
								</div>

								<label for="description" class="block text-sm mb-2 dark:text-white"
									>Valeur optimale du stock</label
								>
								<div class="relative">
									<input
										type="number"
										id="max"
										name="max"
										placeholder="Valeur optimale du stock"
										class="py-3 px-4 block w-full border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
										required
										aria-describedby="text-error"
										bind:value={newItem.optimal_amount}
										on:input={(e) => {
											// @ts-ignore
											let optimal_amount = parseInt(e.target?.value);
											newItem.optimal_amount = optimal_amount;
										}}
									/>
								</div>

								<label for="description" class="block text-sm mb-2 dark:text-white">Status</label>
								<div class="relative">
									<select
										id="status"
										name="status"
										class="py-3 px-4 block w-full border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
										required
										aria-describedby="text-error"
										bind:value={newItem.state}
									>
										<option value="buyable">Achetable</option>
										<option value="not_buyable">Non achetable</option>
									</select>
								</div>

								<button
									type="submit"
									class="mt-4 py-3 px-4 inline-flex justify-center items-center gap-2 rounded-md border border-transparent font-semibold bg-blue-500 text-white hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-all text-sm dark:focus:ring-offset-gray-800"
									on:click={() => createNewItem()}
									data-hs-overlay="#hs-modal-new-item">Créer</button
								>
							</label>
						</div>
					</div>
					<!-- End Form -->
				</div>
			</div>
		</div>
	</div>
</div>

<div
	id="hs-modal-edit-item"
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
						Appliquer une promotion
					</h2>
				</div>

				<!-- One side to change its category, one side to modify the current promotion -->

				<div class="flex flex-col gap-4 mt-4">
					<label for="promo" class="block text-sm dark:text-white">Promotion en cours</label>
					<div class="relative flex flex-row">
						<input
							type="number"
							class="block text-sm dark:text-white/[.8] w-16 break-words p-2 bg-transparent"
							value={selectedItem?.promotion ?? 0}
							on:input={(e) => {
								// @ts-ignore
								let promotion = parseInt(e.target?.value);
								editItem(
									selectedItem?.id ?? '',
									{ promotion: promotion },
									selectedItem?.category_id ?? ''
								);
							}}
						/>
						<span class="self-center text-sm text-gray-400"> % </span>
					</div>

					<label for="end-promo" class="block text-sm dark:text-white">Fin de la promotion</label>
					<div class="relative">
						<!-- promotion ends at is unix seconds, we need to parse it -->
						<input
							type="date"
							id="end-promo"
							name="end-promo"
							value={new Date((selectedItem?.promotion_ends_at ?? 0) * 1000)
								.toISOString()
								.split('T')[0]}
							class="py-3 px-4 block w-full border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
							required
							aria-describedby="text-error"
							on:input={(e) => {
								// @ts-ignore
								let date = e.target?.value;
								let timestampSeconds = Date.parse(date) / 1000;
								editItem(
									selectedItem?.id ?? '',
									{ promotion_ends_at: timestampSeconds },
									selectedItem?.category_id ?? ''
								);
							}}
						/>
					</div>

					<select
						id="role"
						name="role"
						class="py-3 px-4 block w-full border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
						required
						aria-describedby="text-error"
						bind:value={promoItemPriceRole}
					>
						<option value="coutant">Prix réel coutant</option>
						<option value="externe">Prix réel externe</option>
						<option value="ceten">Prix réel ceten</option>
						<option value="staff_bar">Prix réel staff</option>
						<option value="privilegies">Prix réel privilégiés</option>
						<option value="menu">Prix réel menu</option>
					</select>
					<div class="relative">
						<span class="self-center text-sm text-gray-400"
							>{formatPrice((selectedItem?.display_prices ?? {})[promoItemPriceRole] ?? 0)}</span
						>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>

{#if deletingItem}
	<ConfirmationPopup 
		message={confirmationMessage}
		confirm_text="Supprimer"
		cancel_callback={() => {
			deletingItem = false;
		}} 
		confirm_callback={deleteItemCallback}
	/>
{/if}

<div class="h-full flex flex-col bg-gray-50 dark:bg-gray-900">
	<div class="m-3 p-2">
		<!-- Mobile: Compact filters with toggle -->
		<div class="lg:hidden">
			<div class="flex items-center gap-2 mb-2">
				<button
					class="flex-1 py-2 px-3 inline-flex justify-center items-center gap-2 rounded-md border border-gray-300 dark:border-gray-600 font-medium bg-white dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-600 transition-all text-sm"
					on:click={() => filtersOpen = !filtersOpen}
				>
					<iconify-icon icon="mdi:filter-variant" width="18" height="18"></iconify-icon>
					Filtres
					{#if searchCategory || searchState || searchName}
						<span class="bg-blue-500 text-white text-xs rounded-full px-1.5 py-0.5 min-w-[18px]">
							{(searchCategory ? 1 : 0) + (searchState ? 1 : 0) + (searchName ? 1 : 0)}
						</span>
					{/if}
					<iconify-icon icon={filtersOpen ? "mdi:chevron-up" : "mdi:chevron-down"} width="18" height="18"></iconify-icon>
				</button>
				<button
					class="py-2 px-3 inline-flex justify-center items-center gap-2 rounded-md border border-transparent font-semibold bg-blue-500 text-white hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-all text-sm dark:focus:ring-offset-gray-800"
					data-hs-overlay="#hs-modal-new-item"
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
					Ajouter
				</button>
			</div>
			{#if filtersOpen}
				<div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3 mb-2 space-y-3">
					<div class="grid grid-cols-2 gap-3">
						<div>
							<label class="text-xs font-medium text-gray-500 dark:text-gray-400 block mb-1">Catégorie</label>
							<select
								id="category-mobile"
								name="category"
								class="w-full px-2 py-1.5 text-sm bg-gray-50 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:text-white"
								on:change={(e) => {
									// @ts-ignore
									searchCategory = e.target?.value;
									if (searchCategory === '') searchCategory = undefined;
									page = 1;
									reloadItems();
								}}
							>
								<option value="">Toutes</option>
								{#each categories as category}
									<option value={category.id}>{category.name}</option>
								{/each}
							</select>
						</div>
						<div>
							<label class="text-xs font-medium text-gray-500 dark:text-gray-400 block mb-1">Etat</label>
							<select
								id="state-mobile"
								name="state"
								class="w-full px-2 py-1.5 text-sm bg-gray-50 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:text-white"
								on:change={(e) => {
									// @ts-ignore
									let val = e.target?.value;
									if (val == '') searchState = undefined;
									else searchState = val;
									page = 1;
									reloadItems();
								}}
							>
								<option value="">Tous</option>
								<option value="buyable">Achetable</option>
								<option value="not_buyable">Non achetable</option>
							</select>
						</div>
					</div>
					<div>
						<label class="text-xs font-medium text-gray-500 dark:text-gray-400 block mb-1">Fournisseur</label>
						<select
							id="fournisseur-mobile"
							name="fournisseur"
							class="w-full px-2 py-1.5 text-sm bg-gray-50 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:text-white"
							on:change={(e) => {
								// @ts-ignore
								let val = e.target?.value;
								if (val == '') searchFournisseur = undefined;
								else searchFournisseur = val;
								page = 1;
								reloadItems();
							}}
						>
							<option value="">Tous</option>
							{#each fournisseurIterator as [val, name]}
								<option value={val}>{name}</option>
							{/each}
						</select>
					</div>
					<div>
						<label class="text-xs font-medium text-gray-500 dark:text-gray-400 block mb-1">Rechercher</label>
						<input
							type="text"
							placeholder="Nom du produit..."
							bind:value={searchName}
							on:input={() => {
								page = 1;
								reloadItems();
							}}
							class="w-full px-2 py-1.5 text-sm bg-gray-50 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:text-white placeholder-gray-400 dark:placeholder-gray-500"
						/>
					</div>
				</div>
			{/if}
		</div>

		<!-- Desktop: Original layout -->
		<div class="hidden lg:flex flex-row flex-wrap items-center gap-6">
			<div class="flex flex-row items-center gap-3">
				<span class="text-sm font-medium text-gray-700 dark:text-gray-300">Par catégorie:</span>
				<select
					id="category"
					name="category"
					class="px-3 py-1.5 text-sm bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:text-white dark:placeholder-gray-500 w-auto"
					on:change={(e) => {
						// @ts-ignore
						searchCategory = e.target?.value;
						if (searchCategory === '') searchCategory = undefined;
						page = 1;
						reloadItems();
					}}
				>
					<option value="">Pas de filtre</option>
					{#each categories as category}
						<option value={category.id}>{category.name}</option>
					{/each}
				</select>
			</div>
			<div class="flex flex-row items-center gap-3">
				<span class="text-sm font-medium text-gray-700 dark:text-gray-300">Par état:</span>
				<select
					id="state"
					name="state"
					class="px-3 py-1.5 text-sm bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:text-white dark:placeholder-gray-500 w-auto"
					on:change={(e) => {
						// @ts-ignore
						let val = e.target?.value;
						if (val == '') searchState = undefined;
						else searchState = val;
						page = 1;
						reloadItems();
					}}
				>
					<option value="">Pas de filtre</option>
					<option value="buyable">Achetable</option>
					<option value="not_buyable">Non achetable</option>
				</select>
			</div>
			<div class="flex flex-row items-center gap-3">
				<span class="text-sm font-medium text-gray-700 dark:text-gray-300">Par nom:</span>
				<input
					type="text"
					placeholder="Rechercher..."
					bind:value={searchName}
					on:input={() => {
						page = 1;
						reloadItems();
					}}
					class="px-3 py-1.5 text-sm bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:text-white placeholder-gray-400 dark:placeholder-gray-500 w-auto"
				/>
			</div>
			<div class="flex flex-row items-center gap-3">
				<span class="text-sm font-medium text-gray-700 dark:text-gray-300">Par fourn.:</span>
				<select
					id="fournisseur"
					name="fournisseur"
					class="px-3 py-1.5 text-sm bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:text-white dark:placeholder-gray-500 w-auto"
					on:change={(e) => {
						// @ts-ignore
						let val = e.target?.value;
						if (val == '') searchFournisseur = undefined;
						else searchFournisseur = val;
						page = 1;
						reloadItems();
					}}
				>
					<option value="">Pas de filtre</option>
					{#each fournisseurIterator as [val, name]}
						<option value={val}>{name}</option>
					{/each}
				</select>
			</div>
			<div class="flex flex-row items-center gap-3">
				<button
					class="py-2 px-3 inline-flex justify-center items-center gap-2 rounded-md border border-transparent font-semibold bg-blue-500 text-white hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-all text-sm dark:focus:ring-offset-gray-800 w-auto"
					data-hs-overlay="#hs-modal-new-item"
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
					Ajouter un produit
				</button>
			</div>
		</div>
	</div>

	<div class="flex-1 min-h-0 w-full overflow-x-auto overflow-y-auto">
		<!-- Desktop Table View -->
		<div class="hidden min-[1300px]:block min-w-full bg-white dark:bg-gray-800 rounded-lg shadow-sm overflow-visible">
			<table class="min-w-full table-fixed divide-y divide-gray-200 dark:divide-gray-700">
				<colgroup>
					<col class="w-[14%]" />
					<col class="w-[12%]" />
					<col class="w-[8%]" />
					<col class="w-[8%]" />
					<col class="w-[14%]" />
					<col class="w-[14%]" />
					<col class="w-[16%]" />
					<col class="w-[14%]" />
				</colgroup>
				<thead class="bg-gray-50 dark:bg-gray-700">
					<tr>
						<th class="px-6 py-4 text-left text-sm font-semibold text-gray-900 dark:text-gray-300 border-r border-gray-200 dark:border-gray-600">
							Nom
						</th>
						<th class="px-6 py-4 text-left text-sm font-semibold text-gray-900 dark:text-gray-300 border-r border-gray-200 dark:border-gray-600">
							Catégorie
						</th>
						<th class="px-6 py-4 text-left text-sm font-semibold text-gray-900 dark:text-gray-300 border-r border-gray-200 dark:border-gray-600">
							Image
						</th>
						<th class="px-6 py-4 text-left text-sm font-semibold text-gray-900 dark:text-gray-300 border-r border-gray-200 dark:border-gray-600">
							Achetable
						</th>
						<th class="px-6 py-4 text-left text-sm font-semibold text-gray-900 dark:text-gray-300 border-r border-gray-200 dark:border-gray-600">
							Stock / Limite / Optimal
						</th>
						<th class="px-6 py-4 text-left text-sm font-semibold text-gray-900 dark:text-gray-300 border-r border-gray-200 dark:border-gray-600">
							Fournisseur
						</th>
						<th class="px-6 py-4 text-left text-sm font-semibold text-gray-900 dark:text-gray-300 border-r border-gray-200 dark:border-gray-600">
							Prix
						</th>
						<th class="px-6 py-4 text-left text-sm font-semibold text-gray-900 dark:text-gray-300">
							Actions
						</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-gray-200 dark:divide-gray-700">
					{#each items as item}
						<tr>
							<td class="px-6 py-4 text-sm text-gray-900 dark:text-gray-300 border-r border-gray-200 dark:border-gray-600">
								<div class="relative group">
									<input
										type="text"
										class="block w-full text-sm dark:text-white/[.8] p-2 bg-gray-50 dark:bg-slate-900 border border-gray-200 dark:border-gray-600 rounded focus:border-blue-500 focus:ring-2 focus:ring-blue-200 dark:focus:ring-blue-800 outline-none transition-all"
										value={item.name}
										on:input={(e) => {
											// @ts-ignore
											let name = e.target?.value;
											editItem(item.id, { name: name }, item.category_id);
										}}
									/>
									<iconify-icon icon="mdi:pencil" class="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400 opacity-0 group-hover:opacity-100 transition-opacity" width="14" height="14"></iconify-icon>
								</div>
							</td>
							<td class="px-6 py-4 text-sm text-gray-900 dark:text-gray-300 border-r border-gray-200 dark:border-gray-600">
								<select
									class="block w-full text-sm dark:text-white/[.8] dark:bg-slate-900 p-2 bg-gray-50 dark:bg-slate-900 border border-gray-200 dark:border-gray-600 rounded focus:border-blue-500 focus:ring-2 focus:ring-blue-200 dark:focus:ring-blue-800 outline-none transition-all"
									value={item.category_id}
									on:change={(e) => {
										// @ts-ignore
										let category_id = e.target?.value;
										editItem(item.id, { category_id: category_id }, item.category_id);
									}}
								>
									{#each categories as category}
										<option value={category.id}>{category.name}</option>
									{/each}
								</select>
							</td>
							<td class="px-6 py-4 text-sm text-gray-900 dark:text-gray-300 border-r border-gray-200 dark:border-gray-600 relative">
								<input
									type="file"
									class="absolute w-12 h-12 opacity-0 cursor-pointer"
									on:change={(e) => {
										// @ts-ignore
										let file = e.target?.files[0];
										reuploadItemPicture(item.id, file, item.category_id);
									}}
								/>
								{#if item.picture_uri != ''}
									<img
										src={api() + item.picture_uri}
										alt="indisponible"
										class="w-12 h-12 rounded-md object-cover"
									/>
								{/if}
							</td>
							<td class="px-6 py-4 text-sm text-gray-900 dark:text-gray-300 border-r border-gray-200 dark:border-gray-600">
								<select
									class="block w-full text-sm dark:text-white/[.8] dark:bg-slate-900 p-2 bg-gray-50 dark:bg-slate-900 border border-gray-200 dark:border-gray-600 rounded focus:border-blue-500 focus:ring-2 focus:ring-blue-200 dark:focus:ring-blue-800 outline-none transition-all"
									value={item.state}
									on:change={(e) => {
										// @ts-ignore
										let state = e.target?.value;
										editItem(item.id, { state: state }, item.category_id);
									}}
								>
									<option value="buyable">✅</option>
									<option value="not_buyable">❌</option>
								</select>
							</td>
							<td class="px-6 py-4 text-sm text-gray-900 dark:text-gray-300 border-r border-gray-200 dark:border-gray-600">
								<StockCell {item} {editItem} />
							</td>
							<td class="px-6 py-4 text-sm text-gray-900 dark:text-gray-300 border-r border-gray-200 dark:border-gray-600">
								<FournisseurCell {item} {editItem} />
							</td>
							<td class="px-6 py-4 text-sm text-gray-900 dark:text-gray-300 border-r border-gray-200 dark:border-gray-600">
								<PriceCell {item} {editItem} />
							</td>
							<td class="px-6 py-4 text-sm text-gray-900 dark:text-gray-300 relative flex flex-col gap-1">
								<button
									class="{item.promotion_ends_at ?? 0 > new Date().getTime() / 1000
										? 'animate-pulse'
										: ''} flex items-center gap-2 text-sm text-blue-600 dark:text-blue-400 font-medium hover:bg-blue-50 dark:hover:bg-blue-900/30 p-2 rounded-md transition-colors"
									data-hs-overlay="#hs-modal-edit-item"
									on:click={() => (selectedItem = item)}
								>
									<iconify-icon icon="mdi:percent" width="16" height="16"></iconify-icon>
									Promotions
								</button>
								<button
									class="flex items-center gap-2 w-full text-left px-2 py-1.5 text-sm text-red-700 dark:text-red-300 hover:bg-red-50 dark:hover:bg-red-900/30 rounded-md transition-colors"
									on:click={() => {
										deleteItemCallback = () => {
											deletingItem = false;
											deleteItem(item.id, item.category_id);
										};
										confirmationMessage = "Supprimer '" + item.name + "' ?";
										deletingItem = true;
									}}
								>
									<iconify-icon icon="mdi:delete" width="16" height="16"></iconify-icon>
									Supprimer
								</button>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>

		<!-- Mobile Card View -->
		<div class="block min-[1300px]:hidden space-y-4 px-2 pb-4">
			{#if items.length === 0}
				<div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-8 text-center">
					<p class="text-gray-500 dark:text-gray-400">Aucun produit trouvé</p>
				</div>
			{:else}
				{#each items as item}
					<div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-4">
						<div class="flex justify-between items-start mb-3">
							<div class="flex-1">
								<div class="grid grid-cols-2 gap-2 mb-2">
									<div>
										<label class="text-xs font-medium text-gray-500 dark:text-gray-400 block">Nom</label>
										<input
											type="text"
											class="w-full text-sm dark:text-white/[.8] bg-transparent border border-gray-200 dark:border-gray-600 rounded p-2 focus:border-blue-500 focus:outline-none"
											value={item.name}
											on:input={(e) => {
												// @ts-ignore
												let name = e.target?.value;
												editItem(item.id, { name: name }, item.category_id);
											}}
										/>
									</div>
									<div>
										<label class="text-xs font-medium text-gray-500 dark:text-gray-400 block">Catégorie</label>
										<select
											class="w-full text-sm dark:text-white/[.8] dark:bg-slate-900 bg-transparent border border-gray-200 dark:border-gray-600 rounded p-2 focus:border-blue-500 focus:outline-none"
											value={item.category_id}
											on:change={(e) => {
												// @ts-ignore
												let category_id = e.target?.value;
												editItem(item.id, { category_id: category_id }, item.category_id);
											}}
										>
											{#each categories as category}
												<option value={category.id}>{category.name}</option>
											{/each}
										</select>
									</div>
								</div>
								<div class="mb-2">
									<label class="text-xs font-medium text-gray-500 dark:text-gray-400 block">Image</label>
									<div class="relative inline-block">
										<input
											type="file"
											class="absolute w-12 h-12 opacity-0 cursor-pointer"
											on:change={(e) => {
												// @ts-ignore
												let file = e.target?.files[0];
												reuploadItemPicture(item.id, file, item.category_id);
											}}
										/>
										{#if item.picture_uri != ''}
											<img
												src={api() + item.picture_uri}
												alt="indisponible"
												class="w-12 h-12 rounded-md object-cover"
											/>
										{/if}
									</div>
								</div>
							</div>
						</div>
						<div class="grid grid-cols-2 sm:grid-cols-4 gap-2 mb-3">
							<div>
								<label class="text-xs font-medium text-gray-500 dark:text-gray-400 flex items-center gap-1">
									Achetable
								</label>
								<select
									class="w-full text-sm dark:text-white/[.8] dark:bg-slate-900 bg-gray-50 dark:bg-slate-900 border border-gray-200 dark:border-gray-600 rounded p-2 focus:border-blue-500 focus:ring-2 focus:ring-blue-200 dark:focus:ring-blue-800 outline-none transition-all"
									value={item.state}
									on:change={(e) => {
										// @ts-ignore
										let state = e.target?.value;
										editItem(item.id, { state: state }, item.category_id);
									}}
								>
									<option value="buyable">✅</option>
									<option value="not_buyable">❌</option>
								</select>
							</div>
							<div>
								<label class="text-xs font-medium text-gray-500 dark:text-gray-400 flex items-center gap-1">
									<iconify-icon icon="mdi:pencil" width="12" height="12"></iconify-icon>
									Stock
								</label>
								<input
									type="number"
									class="w-full text-sm dark:text-white/[.8] bg-gray-50 dark:bg-slate-900 border border-gray-200 dark:border-gray-600 rounded p-2 focus:border-blue-500 focus:ring-2 focus:ring-blue-200 dark:focus:ring-blue-800 outline-none transition-all"
									value={item.amount_left}
									on:input={(e) => {
										// @ts-ignore
										let stock = parseInt(e.target?.value);
										editItem(item.id, { amount_left: stock }, item.category_id);
									}}
								/>
							</div>
							<div>
								<label class="text-xs font-medium text-gray-500 dark:text-gray-400 flex items-center gap-1">
									<iconify-icon icon="mdi:pencil" width="12" height="12"></iconify-icon>
									Limite
								</label>
								<input
									type="number"
									class="w-full text-sm dark:text-white/[.8] bg-gray-50 dark:bg-slate-900 border border-gray-200 dark:border-gray-600 rounded p-2 focus:border-blue-500 focus:ring-2 focus:ring-blue-200 dark:focus:ring-blue-800 outline-none transition-all"
									value={item.buy_limit}
									on:input={(e) => {
										// @ts-ignore
										let buy_limit = parseInt(e.target?.value);
										// @ts-ignore
										if (e.target?.value === '') {
											editItem(item.id, { buy_limit: -1 }, item.category_id);
											return;
										}
										editItem(item.id, { buy_limit: buy_limit }, item.category_id);
									}}
								/>
							</div>
							<div>
								<label class="text-xs font-medium text-gray-500 dark:text-gray-400 flex items-center gap-1">
									<iconify-icon icon="mdi:pencil" width="12" height="12"></iconify-icon>
									Optimal
								</label>
								<input
									type="number"
									class="w-full text-sm dark:text-white/[.8] bg-gray-50 dark:bg-slate-900 border border-gray-200 dark:border-gray-600 rounded p-2 focus:border-blue-500 focus:ring-2 focus:ring-blue-200 dark:focus:ring-blue-800 outline-none transition-all"
									value={item.optimal_amount}
									on:input={(e) => {
										// @ts-ignore
										let optimal_amount = parseInt(e.target?.value);
										editItem(item.id, { optimal_amount: optimal_amount }, item.category_id);
									}}
								/>
							</div>
						</div>
						<div class="grid grid-cols-1 sm:grid-cols-2 gap-3 mb-3">
							<div>
								<label class="text-xs font-medium text-gray-500 dark:text-gray-400 block mb-1">Fournisseur</label>
								<FournisseurCell {item} {editItem} />
							</div>
							<div>
								<label class="text-xs font-medium text-gray-500 dark:text-gray-400 block mb-1">Prix</label>
								<PriceCell {item} {editItem} />
							</div>
						</div>
						<div class="flex flex-wrap gap-2">
							<button
								class="flex-1 min-w-0 px-3 py-2 text-xs bg-yellow-500 text-white rounded hover:bg-yellow-600 transition-colors flex items-center justify-center gap-1"
								data-hs-overlay="#hs-modal-edit-item"
								on:click={() => (selectedItem = item)}
							>
								<iconify-icon icon="mdi:tag" width="16" height="16" />
								Promotions
							</button>
							<button
								class="flex-1 min-w-0 px-3 py-2 text-xs bg-red-500 text-white rounded hover:bg-red-600 transition-colors flex items-center justify-center gap-1"
								on:click={() => {
									deleteItemCallback = () => {
										deletingItem = false;
										deleteItem(item.id, item.category_id);
									};
									confirmationMessage = "Supprimer '" + item.name + "' ?";
									deletingItem = true;
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
		resultsCount={items.length}
		showPageInput={true}
		on:prevPage={prevPage}
		on:nextPage={nextPage}
		on:pageChange={handlePageInput}
	/>
</div>

<style>
	input[type='number'] {
		-moz-appearance: textfield;
	}

	input[type='number']::-webkit-outer-spin-button,
	input[type='number']::-webkit-inner-spin-button {
		-webkit-appearance: none;
		margin: 0;
	}
</style>
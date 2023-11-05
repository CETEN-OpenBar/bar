<script lang="ts">
	import type {
		Category,
		Item,
		NewItem,
		ItemPrices,
		UpdateItem,
		AccountPriceRole,
		ItemState
	} from '$lib/api';
	import { api } from '$lib/config/config';
	import { categoriesApi, itemsApi } from '$lib/requests/requests';
	import { formatPrice, parsePrice } from '$lib/utils';
	import { onMount } from 'svelte';

	// Type for NewItem with categoryId
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
			exte: 0,
			ceten: 0,
			staff: 0,
			vip: 0,
			interne: 0,
			membre_bureau: 0,
			membre_privilegie: 0,
		} as ItemPrices,
		amount_left: 0,
		buy_limit: undefined,
		optimal_amount: 0,
		state: 'buyable',
		category_id: ''
	};
	let newItemPriceRole: AccountPriceRole = 'ceten';
	let editItemPriceRole: AccountPriceRole = 'ceten';
	let promoItemPriceRole: AccountPriceRole = 'ceten';

	let page = 0;
	let max_page = 0;
	let itemsPerPage = 10;

	let rebounceTimeout: number | null = null;

	let searchState: ItemState | undefined = undefined;
	let searchCategory: string | undefined = undefined;
	let searchName: string | undefined = undefined;

	onMount(() => {
		categoriesApi()
			.getCategories({ withCredentials: true })
			.then((res) => {
				categories = res.data ?? [];
			});

		reloadItems();
	});

	function reloadItems() {
		itemsApi()
			.getAllItems(page, itemsPerPage, searchState, searchCategory, searchName, {
				withCredentials: true
			})
			.then((res) => {
				max_page = res.data.max_page ?? 0;
				page = res.data.page ?? 0;
				itemsPerPage = res.data.limit ?? 0;
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
										on:change={// resets the value of the price input
										() => {
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
										<option value="interne">Prix interne</option>
										<option value="exte">Prix externe</option>
										<option value="ceten">Prix ceten</option>
										<option value="staff">Prix staff</option>
										<option value="vip">Prix VIP</option>
										<option value="membre">Prix Membre</option>
										<option value="membre_privilegie">Prix Membre privilégié</option>
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
										class="py-3 px-4 block w-[90%] border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
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
						<option value="interne">Prix réel interne</option>
						<option value="exte">Prix réel externe</option>
						<option value="ceten">Prix réel ceten</option>
						<option value="staff">Prix réel staff</option>
						<option value="vip">Prix réel VIP</option>
						<option value="membre">Prix réel Membre</option>
						<option value="membre_privilegie">Prix réel Membre privilégié</option>
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

<!-- Table Section -->
<div class="max-w-[95%] px-4 py-10 sm:px-6 lg:px-8 lg:py-14 mx-auto">
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
							<p class="text-sm text-gray-600 dark:text-gray-400">Ajouter des produits</p>
						</div>

						<div class="px-4 mx-4 grid grid-cols-3 gap-4">
							<!-- Titles -->
							<div class="text-lg dark:text-white">
								Par catégorie
							</div>
							<div class="text-lg dark:text-white">
								Par état
							</div>
							<div class="text-lg dark:text-white">
								Par nom
							</div>
							<div class="relative mt-4 w-96 md:mt-0">
								<!-- filter by category -->
								<select
									id="category"
									name="category"
									class="py-3 px-4 block w-full border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
									required
									aria-describedby="text-error"
									on:change={(e) => {
										// @ts-ignore
										searchCategory = e.target?.value;
										if (searchCategory === '') searchCategory = undefined;
										reloadItems();
									}}
								>
									<option value="">Pas de filtre</option>
									{#each categories as category}
										<option value={category.id}>{category.name}</option>
									{/each}
								</select>
							</div>
							<div class="relative mt-4 w-96 md:mt-0">
								<!-- filter by state -->
								<select
									id="state"
									name="state"
									class="py-3 px-4 block w-full border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
									required
									aria-describedby="text-error"
									on:change={(e) => {
										// @ts-ignore
										let val = e.target?.value;
										if (val == "") searchState = undefined;
										else searchState = val;
										reloadItems();
									}}
								>
									<option value="">Pas de filtre</option>
									<option value="buyable">Achetable</option>
									<option value="not_buyable">Non achetable</option>
								</select>
							</div>
							<div class="relative mt-4 w-96 md:mt-0">
								<input
									type="text"
									class="py-3 px-4 w-full border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
									placeholder="Rechercher"
									aria-label="Rechercher"
									on:input={(e) => {
										// @ts-ignore
										searchName = e.target.value.toLowerCase();
										reloadItems();
									}}
								/>
								<svg
									class="absolute w-4 h-4 right-3 top-3 text-gray-400 dark:text-gray-300 pointer-events-none"
									xmlns="http://www.w3.org/2000/svg"
									width="16"
									height="16"
									viewBox="0 0 16 16"
									fill="none"
								>
									<path
										d="M11.6667 11.6667L15.3333 15.3333"
										stroke="currentColor"
										stroke-width="1.5"
										stroke-linecap="round"
										stroke-linejoin="round"
									/>
									<path
										d="M6.66663 12.6667C9.53763 12.6667 12 10.2037 12 7.33337C12 4.46337 9.53763 2.00004 6.66663 2.00004C3.79563 2.00004 1.33329 4.46337 1.33329 7.33337C1.33329 10.2037 3.79563 12.6667 6.66663 12.6667Z"
										stroke="currentColor"
										stroke-width="1.5"
										stroke-linecap="round"
										stroke-linejoin="round"
									/>
								</svg>
							</div>
						</div>

						<div>
							<div class="inline-flex gap-x-2">
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
									Ajouter un produit
								</button>
							</div>
						</div>
					</div>
					<!-- End Header -->

					<!-- Table -->
					<table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
						<thead class="bg-gray-50 dark:bg-slate-800">
							<tr>
								<th scope="col" class="px-6 py-3">
									<span
										class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
									>
										Nom
									</span>
								</th>
								<th scope="col" class="px-6 py-3">
									<span
										class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
									>
										Catégorie
									</span>
								</th>
								<th scope="col" class="px-6 py-3">
									<span
										class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
									>
										Image
									</span>
								</th>
								<th scope="col" class="px-2 py-3">
									<p
										class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
									>
										Etat
									</p>
								</th>
								<th scope="col" class="px-2 py-3 w-2">
									<span
										class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
									>
										En stock
									</span>
								</th>
								<th scope="col" class="px-2 py-3">
									<span
										class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
									>
										Limite d'achat
									</span>
								</th>
								<th scope="col" class="px-2 py-3">
									<span
										class="text-center text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
									>
										Montant optimal en stock
									</span>
								</th>
								<th scope="col" class="px-6 py-3">
									<select
										id="role"
										name="role"
										class="py-3 px-4 block w-full border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
										required
										aria-describedby="text-error"
										bind:value={editItemPriceRole}
										on:change={() => {
											// Remove the value of elems with id "price"
											let elems = document.querySelectorAll('[id=price]');
											elems.forEach((elem) => {
												// @ts-ignore
												elem.value = '';
											});
										}}
									>
										<option value="interne">Prix interne</option>
										<option value="exte">Prix externe</option>
										<option value="ceten">Prix ceten</option>
										<option value="staff">Prix staff</option>
										<option value="vip">Prix VIP</option>
										<option value="membre">Prix Membre</option>
										<option value="membre_privilegie">Prix Membre privilégié</option>
									</select>
								</th>
								<th scope="col" class="px-6 py-3 text-right" />
							</tr>
						</thead>

						<tbody class="divide-y divide-gray-200 dark:divide-gray-700">
							{#each items as item}
								<tr>
									<td class="h-px w-72">
										<div class="px-6 py-3 grid justify-center">
											<input
												type="text"
												class="py-3 px-2 block w-[90%] border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
												value={item.name}
												on:input={(e) => {
													// @ts-ignore
													let name = e.target?.value;
													editItem(item.id, { name: name }, item.category_id);
												}}
											/>
										</div>
									</td>
									<td class="h-px w-52">
										<div class="px-6 py-3 grid justify-center">
											<select
												class="block text-sm dark:text-white/[.8] dark:bg-slate-900 break-words p-2 bg-transparent"
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
									</td>
									<td class="h-px w-36 grid justify-center">
										<!-- Display a miniature of the image -->
										<div class="px-6 py-3 relative">
											<!-- <img
												src={api() + category.picture_uri}
												alt="indisponible"
												class="w-full h-full rounded-md object-cover"
											/> -->

											<!-- input in front of the image to click & reupload -->
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
									</td>
									<td class="h-px w-52">
										<div class="px-2 py-3 grid justify-center">
											<select
												class="block text-sm dark:text-white/[.8] dark:bg-slate-900 break-words p-2 bg-transparent"
												value={item.state}
												on:change={(e) => {
													// @ts-ignore
													let state = e.target?.value;
													editItem(item.id, { state: state }, item.category_id);
												}}
											>
												<option value="buyable">Achetable</option>
												<option value="not_buyable">Pas achetable</option>
											</select>
										</div>
									</td>
									<td class="h-px w-52">
										<div class="px-2 py-3 grid justify-center">
											<input
												type="number"
												class="py-3 px-2 block w-[90%] border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
												value={item.amount_left}
												on:input={(e) => {
													// @ts-ignore
													let stock = parseInt(e.target?.value);
													editItem(item.id, { amount_left: stock }, item.category_id);
												}}
											/>
										</div>
									</td>
									<td class="h-px w-52">
										<div class="px-6 py-3 grid justify-center">
											<input
												type="number"
												class="py-3 px-2 block w-[90%] border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
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
									</td>
									<td class="h-px w-52">
										<div class="px-6 py-3 grid justify-center">
											<input
												type="number"
												class="py-3 px-2 block w-[90%] border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
												value={item.optimal_amount}
												on:input={(e) => {
													// @ts-ignore
													let optimal_amount = parseInt(e.target?.value);
													editItem(item.id, { optimal_amount: optimal_amount }, item.category_id);
												}}
											/>
										</div>
									</td>
									<td class="h-px w-52">
										<div class="px-6 py-3 grid justify-center">
											<input
												type="number"
												id="price"
												name="price"
												placeholder={formatPrice(item.prices[editItemPriceRole])}
												class="py-3 px-2 block w-[90%] border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
												required
												aria-describedby="text-error"
												on:input={(e) => {
													let prices = item.prices;
													// @ts-ignore
													prices[editItemPriceRole] = parsePrice(e.target?.value);

													editItem(item.id, { prices: prices }, item.category_id);

													// Remove the value of elems with id "price" with rebounce timeout
													if (rebounceTimeout) clearTimeout(rebounceTimeout);
													rebounceTimeout = setTimeout(() => {
														let elems = document.querySelectorAll('[id=price]');
														elems.forEach((elem) => {
															// @ts-ignore
															elem.value = '';
														});
													}, 1000);
												}}
											/>
										</div>
									</td>
									<td class="h-px w-px whitespace-nowrap">
										<div class="px-6 py-1.5 grid justify-center">
											<button
												class="{item.promotion_ends_at ?? 0 > new Date().getTime() / 1000
													? 'animate-pulse'
													: ''} inline-flex items-center gap-x-1.5 text-sm text-blue-600 decoration-2 hover:underline font-medium"
												data-hs-overlay="#hs-modal-edit-item"
												on:click={() => (selectedItem = item)}
											>
												Promotions
											</button>
											<button
												class="inline-flex items-center gap-x-1.5 text-sm text-blue-600 decoration-2 hover:underline font-medium"
												on:click={() => deleteItem(item.id, item.category_id)}
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
								<span class="font-semibold text-gray-800 dark:text-gray-200">{items.length}</span>
								résultats
							</p>
						</div>

						<div>
							<div class="inline-flex gap-x-2">
								<button
									type="button"
									class="py-2 px-3 inline-flex justify-center items-center gap-2 rounded-md border font-medium bg-white text-gray-700 shadow-sm align-middle hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-white focus:ring-blue-600 transition-all text-sm dark:bg-slate-900 dark:hover:bg-slate-800 dark:border-gray-700 dark:text-gray-400 dark:hover:text-white dark:focus:ring-offset-gray-800"
									on:click={() => {
										if (page > 1) page--;
										reloadItems();
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
									Page {page} / {max_page}
								</p>

								<button
									type="button"
									class="py-2 px-3 inline-flex justify-center items-center gap-2 rounded-md border font-medium bg-white text-gray-700 shadow-sm align-middle hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-white focus:ring-blue-600 transition-all text-sm dark:bg-slate-900 dark:hover:bg-slate-800 dark:border-gray-700 dark:text-gray-400 dark:hover:text-white dark:focus:ring-offset-gray-800"
									on:click={() => {
										if (page < max_page) page++;
										reloadItems();
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

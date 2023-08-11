<script lang="ts">
	import type { Category, Item, NewItem, ItemPrices, UpdateItem, AccountPriceRole } from '$lib/api';
	import { api } from '$lib/config/config';
	import { categoriesApi, itemsApi } from '$lib/requests/requests';
	import { formatPrice } from '$lib/utils';
	import { onMount } from 'svelte';

	let categories: Category[] = [];
	let selectedCategory: string = '';
	let selectedItem: Item | null = null;
	let items: Item[] = [];
	let newItem: NewItem = {
		name: '',
		picture: '',
		prices: {
			exte: 0,
			ceten: 0,
			normal: 0,
			staff: 0,
			vip: 0
		} as ItemPrices,
		amount_left: 0,
		buy_limit: 0,
		state: 'buyable'
	};
	let newItemPriceRole: AccountPriceRole = 'normal';
	let editItemPriceRole: AccountPriceRole = 'normal';
	let promoItemPriceRole: AccountPriceRole = 'normal';

	let page = 0;
	let maxPage = 0;
	let itemsPerPage = 10;

	onMount(() => {
		categoriesApi()
			.getCategories({ withCredentials: true })
			.then((res) => {
				categories = res.data ?? [];
			});
	});

	function changeCategory(category: string) {
		selectedCategory = category;
		page = 0;
		maxPage = 0;
		itemsApi()
			.getCategoryItems(selectedCategory, page, itemsPerPage, undefined, { withCredentials: true })
			.then((res) => {
				maxPage = res.data.max_page ?? 0;
				page = res.data.page ?? 0;
				itemsPerPage = res.data.limit ?? 0;
				items = res.data.items ?? [];
			});
	}

	function createNewItem() {
		if (!newItem) return;
		itemsApi()
			.postItem(selectedCategory, newItem, { withCredentials: true })
			.then((res) => {
				items = [...items, res.data];
			});
	}

	function editItem(id: string, item: UpdateItem) {
		if (!newItem) return;
		itemsApi()
			.patchItem(selectedCategory, id, item, { withCredentials: true })
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

	function reuploadItemPicture(id: string, file: File) {
		if (!newItem) return;
		file2Base64(file).then((base64) => {
			base64 = base64.replace('data:', '').replace(/^.+,/, '');
			itemsApi()
				.patchItem(selectedCategory, id, { picture: base64 }, { withCredentials: true })
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

	function deleteItem(id: string) {
		itemsApi()
			.markDeleteItem(selectedCategory, id, { withCredentials: true })
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
									>
										<option value="normal">Prix de base</option>
										<option value="exte">Prix externe</option>
										<option value="ceten">Prix ceten</option>
										<option value="staff">Prix staff</option>
										<option value="vip">Prix VIP</option>
									</select>
								</div>

								<label for="price" class="block text-sm mb-2 dark:text-white"
									>Prix ({newItemPriceRole})</label
								>
								<div class="relative">
									<input
										type="number"
										id="price"
										name="price"
										placeholder="Prix du produit"
										class="py-3 px-4 block w-[90%] border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
										required
										aria-describedby="text-error"
										on:input={(e) => {
											// @ts-ignore
											let price = e.target?.value;
											// split and parse price[0]
											if (price.includes(',')) price = price.split(',');
											else if (price.includes('.')) price = price.split('.');

											if (price.length > 1) price = parseInt(price[0]) * 100 + parseInt(price[1]);
											else price = parseInt(price[0]) * 100;

											newItem.prices[newItemPriceRole] = price;
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
										placeholder="Max par commande"
										class="py-3 px-4 block w-full border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
										required
										aria-describedby="text-error"
										bind:value={newItem.buy_limit}
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
							class="block text-sm dark:text-white/[.8] w-16  break-words p-2 bg-transparent"
							value={selectedItem?.promotion ?? 0}
							on:input={(e) => {
								// @ts-ignore
								let promotion = parseInt(e.target?.value);
								editItem(selectedItem?.id ?? '', { promotion: promotion });
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
								editItem(selectedItem?.id ?? '', { promotion_ends_at: timestampSeconds });
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
						<option value="normal">Prix réel de base</option>
						<option value="exte">Prix réel externe</option>
						<option value="ceten">Prix réel ceten</option>
						<option value="staff">Prix réel staff</option>
						<option value="vip">Prix réel VIP</option>
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

<!-- Categories section -->
<div
	class="flex flex-row gap-5 items-center overflow-x-auto overflow-y-hidden m-5 p-5 rounded-xl bg-slate-200 dark:bg-slate-900"
>
	{#each categories as category}
		<button
			class="w-16 flex-shrink-0 flex flex-col items-center justify-center m-2 rounded-lg dark:text-white transition-colors duration-300"
			on:click={() => {
				changeCategory(category.id);
			}}
		>
			<img class="w-full" src={api() + category.picture_uri} alt={category.name} />
			<span class="text-md font-bold">{category.name}</span>
		</button>
	{/each}
</div>

{#if selectedCategory != ''}
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
												Etat
											</span>
										</div>
									</th>
									<th scope="col" class="px-2 py-3 text-left">
										<div class="flex items-center gap-x-2">
											<span
												class="text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
											>
												En stock
											</span>
										</div>
									</th>
									<th scope="col" class="px-2 py-3 text-left">
										<div class="flex items-center gap-x-2">
											<span
												class="text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
											>
												Limite d'achat
											</span>
										</div>
									</th>
									<th scope="col" class="px-6 py-3 text-left">
										<div class="flex items-center gap-x-2">
											<select
												id="role"
												name="role"
												class="py-3 px-4 block w-full border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
												required
												aria-describedby="text-error"
												bind:value={editItemPriceRole}
												on:change={()=>{
													// Remove the value of elems with id "price"
													let elems = document.querySelectorAll('[id=price]');
													elems.forEach((elem) => {
														// @ts-ignore
														elem.value = "";
													});

												}}
											>
												<option value="normal">Prix de base</option>
												<option value="exte">Prix externe</option>
												<option value="ceten">Prix ceten</option>
												<option value="staff">Prix staff</option>
												<option value="vip">Prix VIP</option>
											</select>
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
												<input
													type="text"
													class="block text-sm dark:text-white/[.8]  break-words p-2 bg-transparent"
													value={item.name}
													on:input={(e) => {
														// @ts-ignore
														let name = e.target?.value;
														editItem(item.id, { name: name });
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
														reuploadItemPicture(item.id, file);
													}}
												/>
												{#if item.picture_uri != ''}
													<img
														src={api() + item.picture_uri}
														alt="indisponible"
														class="w-full h-full rounded-md object-cover"
													/>
												{/if}
											</div>
										</td>
										<td class="h-px w-72">
											<div class="px-6 py-3">
												<select
													class="block text-sm dark:text-white/[.8] dark:bg-slate-900 break-words p-2 bg-transparent"
													value={item.state}
													on:change={(e) => {
														// @ts-ignore
														let state = e.target?.value;
														editItem(item.id, { state: state });
													}}
												>
													<option value="buyable">Achetable</option>
													<option value="not_buyable">Pas achetable</option>
												</select>
											</div>
										</td>
										<td class="h-px w-72">
											<div class="px-6 py-3">
												<input
													type="number"
													class="block text-sm dark:text-white/[.8]  break-words p-2 bg-transparent"
													value={item.amount_left}
													on:input={(e) => {
														// @ts-ignore
														let stock = parseInt(e.target?.value);
														editItem(item.id, { amount_left: stock });
													}}
												/>
											</div>
										</td>
										<td class="h-px w-72">
											<div class="px-6 py-3">
												<input
													type="number"
													class="block text-sm dark:text-white/[.8]  break-words p-2 bg-transparent"
													value={item.buy_limit}
													on:input={(e) => {
														// @ts-ignore
														let buy_limit = parseInt(e.target?.value);
														editItem(item.id, { buy_limit: buy_limit });
													}}
												/>
											</div>
										</td>
										<td class="h-px w-72">
											<div class="px-6 py-3">
												<input
													type="number"
													id="price"
													name="price"
													placeholder={formatPrice(item.prices[editItemPriceRole])}
													class="py-3 px-4 block w-[90%] border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
													required
													aria-describedby="text-error"
													on:input={(e) => {
														// @ts-ignore
														let price = e.target?.value;
														// split and parse price[0]
														if (price.includes(',')) price = price.split(',');
														else if (price.includes('.')) price = price.split('.');

														if (price.length > 1)
															price = parseInt(price[0]) * 100 + parseInt(price[1]);
														else price = parseInt(price[0]) * 100;

														let prices = item.prices;
														prices[editItemPriceRole] = price;

														editItem(item.id, { prices: prices });
													}}
												/>
											</div>
										</td>
										<td class="h-px w-px whitespace-nowrap">
											<div class="px-6 py-1.5">
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
										Page {page} / {Math.ceil(categories.length / itemsPerPage)}
									</p>

									<button
										type="button"
										class="py-2 px-3 inline-flex justify-center items-center gap-2 rounded-md border font-medium bg-white text-gray-700 shadow-sm align-middle hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-white focus:ring-blue-600 transition-all text-sm dark:bg-slate-900 dark:hover:bg-slate-800 dark:border-gray-700 dark:text-gray-400 dark:hover:text-white dark:focus:ring-offset-gray-800"
										on:click={() => {
											if (page < Math.ceil(categories.length / itemsPerPage) - 1) page++;
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
{/if}

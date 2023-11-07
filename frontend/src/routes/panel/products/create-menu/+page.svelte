<script lang="ts">
	import type { Category, Item, ItemPrices, MenuCategory, MenuItem, NewItem } from '$lib/api';
	import { api } from '$lib/config/config';
	import { itemsApi, categoriesApi } from '$lib/requests/requests';
	import { file2Base64, formatPrice, parsePrice } from '$lib/utils';
	import { onMount } from 'svelte';

	onMount(async () => {
		const c = await categoriesApi().getCategories({ withCredentials: true });
		if (!Array.isArray(c.data)) {
			return;
		}

		categories = c.data;
	});

	let currentStep = 0;
	let rebounceTimeout: Map<string, number> = new Map();

	let searchCategory = '';
	let categories: Category[] = [];
	let items: Item[] = [];

	// type for steps
	type Steps = {
		name: string;
		category: Category | null;
		items: Item[] | null;
		categories: MenuCategory[] | undefined;
		prices: ItemPrices;
		available_from?: number;
		available_to?: number;
		amount_left: number;
		optimal_amount: number;
		buy_limit: number;
		picture: string;
	};

	let steps: Steps = {
		name: '',
		category: null,
		items: null,
		categories: undefined,
		prices: {
			interne: 0,
			exte: 0,
			membre_bureau: 0,
			membre_privilegie: 0,
			staff: 0,
			vip: 0,
			ceten: 0
		},
		available_from: undefined,
		available_to: undefined,
		amount_left: 0,
		optimal_amount: 0,
		buy_limit: 0,
		picture: ''
	};

	function createMenu() {
		if (steps.category === null || steps.items === null) {
			return;
		}

		let subItems: MenuItem[] = [];
		for (let i = 0; i < steps.items.length; i++) {
			subItems.push({
				id: steps.items[i].id,
				amount: steps.items[i].amount_left,
				picture_uri: '',
				name: ''
			});
		}

		let newItem: NewItem = {
			name: steps.name,
			prices: steps.prices,
			available_from: steps.available_from,
			available_until: steps.available_to,
			amount_left: steps.amount_left,
			optimal_amount: steps.optimal_amount,
			buy_limit: steps.buy_limit,
			picture: steps.picture,
			menu_items: subItems,
			menu_categories: steps.categories,
			is_menu: true,
			state: 'buyable'
		};

		itemsApi()
			.postItem(steps.category.id, newItem, { withCredentials: true })
			.then((res) => {
				if (res.status === 201) {
					window.location.href = '/panel/products';
				}
			});
	}
</script>

<!-- popup in center too choose category -->
<div class="w-full flex flex-row justify-center mt-5">
	<div class="w-[50%] flex flex-col gap-5">
		{#if currentStep >= 0}
			<div class="bg-slate-100 flex flex-col rounded-xl border-2 p-4">
				<div class="text-3xl text-center font-bold p-5">Nom du menu</div>

				<div class="flex flex-row justify-center">
					<input
						type="text"
						class="w-1/2 rounded-lg p-2 border-2 border-slate-200"
						placeholder="Super menu"
						bind:value={steps.name}
					/>
				</div>

				<div class="flex flex-row justify-center mt-5">
					<button
						class="w-32 flex flex-col flex-shrink-0 bg-blue-100 hover:bg-blue-200 rounded-sm p-5 text-center text-lg"
						on:click={() => {
							currentStep = 1;
						}}
					>
						Suivant
					</button>
				</div>
			</div>
		{/if}
		{#if currentStep >= 1}
			<div class="bg-slate-100 flex flex-col rounded-xl border-2 p-4">
				<div class="text-3xl text-center font-bold p-5">Choix de catégorie</div>

				<div class="flex flex-row justify-center">
					<input
						type="text"
						class="w-1/2 rounded-lg p-2 border-2 border-slate-200"
						placeholder="Rechercher une catégorie"
						bind:value={searchCategory}
					/>
				</div>

				<div class="grid grid-cols-4 h-32 items-center gap-3 overflow-x-auto mt-2 p-5">
					{#each categories as c}
						{#if c.name.toLowerCase().includes(searchCategory.toLowerCase())}
							<button
								class="w-32 flex flex-col flex-shrink-0 bg-blue-100 hover:bg-blue-200 rounded-sm p-5 text-center text-lg"
								on:click={async () => {
									steps.category = c;
									currentStep = 2;
									const res = await itemsApi().getAllItems(0, 16, undefined, undefined, '', {
										withCredentials: true
									});
									if (!Array.isArray(res.data.items)) {
										return;
									}

									items = res.data.items;

									// scroll to product
									const product = document.getElementById('product');
									if (product !== null) {
										product.scrollIntoView({ behavior: 'smooth' });
									}
								}}
							>
								{c.name}
							</button>
						{/if}
					{/each}
				</div>
			</div>
		{/if}
		{#if currentStep >= 2}
			<div id="product" class="bg-slate-100 flex flex-col rounded-xl border-2 p-4">
				<div class="text-3xl text-center font-bold p-5">Mettre des produits (catégories dans l'étape suivante)</div>

				<div class="flex flex-row justify-center">
					<input
						type="text"
						class="w-1/2 rounded-lg p-2 border-2 border-slate-200"
						placeholder="Rechercher un objet"
						on:keyup={async (e) => {
							// @ts-ignore
							let val = e.target?.value;
							const res = await itemsApi().getAllItems(0, 16, undefined, undefined, val, {
								withCredentials: true
							});
							if (!Array.isArray(res.data.items)) {
								return;
							}

							items = res.data.items;
						}}
					/>
				</div>

				<div class="grid grid-cols-4 h-32 items-center gap-3 overflow-x-auto mt-2">
					{#each items as i}
						{#if i.name.toLowerCase().includes(searchCategory.toLowerCase())}
							<button
								class="w-32 flex flex-col flex-shrink-0 bg-blue-100 hover:bg-blue-200 rounded-sm p-5 text-center text-lg"
								on:click={() => {
									currentStep = 3;
									if (steps.items === null) {
										steps.items = [];
									}

									for (let j = 0; j < steps.items.length; j++) {
										if (steps.items[j].id === i.id) {
											return;
										}
									}

									let temp = steps.items;
									i.amount_left = 1;
									temp.push(i);
									steps.items = temp;

									setTimeout(() => {
										// scroll to id
										const item = document.getElementById(i.id);
										if (item !== null) {
											item.scrollIntoView({ behavior: 'smooth' });
										}
									}, 100);
								}}
							>
								{i.name}
							</button>
						{/if}
					{/each}
				</div>

				<div class="grid grid-cols-4 max-h-96 justify-center mt-5 overflow-x-auto">
					{#each steps.items ?? [] as i}
						<div class="flex flex-col">
							<div
								class="flex flex-col items-center justify-center bg-slate-200 rounded-lg p-5 m-2 h-48"
							>
								<img
									draggable="false"
									class="w-32 h-32 object-contain"
									src={api() + i.picture_uri}
									alt={i.name}
								/>
								<span id={i.id} class="text-lg font-bold">{i.name}</span>

								<!-- - , amount, + -->
								<div class="grid grid-cols-3 bg-slate-300 w-full">
									<button
										class="flex flex-col items-center justify-center bg-red-100 rounded-lg p-2"
										on:click={() => {
											// Remove 1 from amount
											let temp = steps.items ?? [];
											for (let j = 0; j < temp.length; j++) {
												if (temp[j].id === i.id) {
													if (temp[j].amount_left === 1) {
														temp = temp.filter((item) => item.id !== i.id);
													} else {
														temp[j].amount_left--;
													}
													break;
												}
											}
											steps.items = temp;
										}}
									>
										-
									</button>

									<div
										class="flex flex-col items-center justify-center bg-slate-200 rounded-lg p-2"
									>
										{#each steps.items ?? [] as j}
											{#if j.id === i.id}
												{j.amount_left}
											{/if}
										{/each}
									</div>
									<button
										class="flex flex-col items-center justify-center bg-green-100 rounded-lg p-2"
										on:click={() => {
											// Add 1 to amount
											let temp = steps.items ?? [];
											for (let j = 0; j < temp.length; j++) {
												if (temp[j].id === i.id) {
													temp[j].amount_left++;
													break;
												}
											}
											steps.items = temp;
										}}
										>+
									</button>
								</div>
							</div>
						</div>
					{/each}
				</div>

				<div class="flex flex-row justify-center mt-5">
					<button
						class="w-32 flex flex-col flex-shrink-0 bg-blue-100 hover:bg-blue-200 rounded-sm p-5 text-center text-lg"
						on:click={() => {
							currentStep = 4;
							setTimeout(() => {
								document.getElementById('categ')?.scrollIntoView({ behavior: 'smooth' });
							}, 100);
						}}
					>
						Suivant
					</button>
				</div>
			</div>
		{/if}
		{#if currentStep >= 3}
			<div id="categ" class="bg-slate-100 flex flex-col rounded-xl border-2 p-4">
				<div class="text-3xl text-center font-bold p-5">Mettre des catégories</div>

				<div class="flex flex-row justify-center">
					<input
						type="text"
						class="w-1/2 rounded-lg p-2 border-2 border-slate-200"
						placeholder="Rechercher une categorie"
						bind:value={searchCategory}
					/>
				</div>

				<div class="grid grid-cols-4 h-32 items-center gap-3 overflow-x-auto mt-2">
					{#each categories as i}
						{#if i.name.toLowerCase().includes(searchCategory.toLowerCase())}
							<button
								class="w-32 flex flex-col flex-shrink-0 bg-blue-100 hover:bg-blue-200 rounded-sm p-5 text-center text-lg"
								on:click={() => {
									currentStep = 3;
									if (!steps.categories) steps.categories = [];

									for (let j = 0; j < steps.categories.length; j++) {
										if (steps.categories[j].id === i.id) {
											return;
										}
									}

									let temp = steps.categories;
									temp.push({
										id: i.id,
										amount: 1,
										name: i.name,
										picture_uri: i.picture_uri
									});
									steps.categories = temp;

									setTimeout(() => {
										// scroll to id
										const item = document.getElementById(i.id);
										if (item !== null) {
											item.scrollIntoView({ behavior: 'smooth' });
										}
									}, 100);
								}}
							>
								{i.name}
							</button>
						{/if}
					{/each}
				</div>

				<div class="grid grid-cols-4 max-h-96 justify-center mt-5 overflow-x-auto">
					{#each steps.categories ?? [] as i}
						<div class="flex flex-col">
							<div
								class="flex flex-col items-center justify-center bg-slate-200 rounded-lg p-5 m-2 h-48"
							>
								<img
									draggable="false"
									class="w-32 h-32 object-contain"
									src={api() + i.picture_uri}
									alt={i.name}
								/>
								<span id={i.id} class="text-lg font-bold">{i.name}</span>

								<!-- - , amount, + -->
								<div class="grid grid-cols-3 bg-slate-300 w-full">
									<button
										class="flex flex-col items-center justify-center bg-red-100 rounded-lg p-2"
										on:click={() => {
											// Remove 1 from amount
											let temp = steps.categories ?? [];
											for (let j = 0; j < temp.length; j++) {
												if (temp[j].id === i.id) {
													if (temp[j].amount === 1) {
														temp = temp.filter((item) => item.id !== i.id);
													} else {
														temp[j].amount--;
													}
													break;
												}
											}
											steps.categories = temp;
										}}
									>
										-
									</button>

									<div
										class="flex flex-col items-center justify-center bg-slate-200 rounded-lg p-2"
									>
										{#each steps.categories ?? [] as j}
											{#if j.id === i.id}
												{j.amount}
											{/if}
										{/each}
									</div>
									<button
										class="flex flex-col items-center justify-center bg-green-100 rounded-lg p-2"
										on:click={() => {
											// Add 1 to amount
											let temp = steps.categories ?? [];
											for (let j = 0; j < temp.length; j++) {
												if (temp[j].id === i.id) {
													temp[j].amount++;
													break;
												}
											}
											steps.categories = temp;
										}}
										>+
									</button>
								</div>
							</div>
						</div>
					{/each}
				</div>

				<div class="flex flex-row justify-center mt-5">
					<button
						class="w-32 flex flex-col flex-shrink-0 bg-blue-100 hover:bg-blue-200 rounded-sm p-5 text-center text-lg"
						on:click={() => {
							currentStep = 4;
						}}
					>
						Suivant
					</button>
				</div>
			</div>
		{/if}
		{#if currentStep >= 4}
			<div class="bg-slate-100 flex flex-col rounded-xl border-2 p-4">
				<div class="text-3xl text-center font-bold p-5">Informations</div>

				<!-- Ask for all the prices -->
				<div class="grid grid-cols-3 justify-center">
					{#each Object.keys(steps.prices) as c}
						<div class="block">
							<label for="price" class="block text-sm mb-2 dark:text-white">Prix ({c})</label>
							<div class="relative">
								<input
									type="number"
									name="price"
									placeholder="Prix du menu"
									class="py-3 px-4 block w-[90%] border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
									required
									aria-describedby="text-error"
									on:keyup={(e) => {
										// @ts-ignore
										steps.prices[c] = parsePrice(e.target?.value);
										let rT = rebounceTimeout.get(c);
										clearTimeout(rT ?? 0);
										rT = setTimeout(() => {
											// @ts-ignore
											let r = formatPrice(parsePrice(e.target?.value));
											// @ts-ignore
											e.target.placeholder = r;
											// @ts-ignore
											e.target.value = r;
										}, 500);
									}}
								/>
							</div>
						</div>
					{/each}
				</div>

				<div class="grid grid-cols-3 justify-center">
					<div class="block">
						<label for="price" class="block text-sm mb-2 dark:text-white">En stock</label>
						<div class="relative">
							<input
								type="number"
								name="price"
								placeholder="En stock"
								class="py-3 px-4 block w-[90%] border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
								required
								aria-describedby="text-error"
								on:keyup={(e) => {
									// @ts-ignore
									steps.amount_left = parseInt(e.target?.value);
								}}
							/>
						</div>
					</div>
					<div class="block">
						<label for="price" class="block text-sm mb-2 dark:text-white">Limite d'achat</label>
						<div class="relative">
							<input
								type="number"
								name="price"
								placeholder="Limite"
								class="py-3 px-4 block w-[90%] border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
								required
								aria-describedby="text-error"
								on:keyup={(e) => {
									// @ts-ignore
									steps.buy_limit = parseInt(e.target?.value);
								}}
							/>
						</div>
					</div>
					<div class="block">
						<label for="price" class="block text-sm mb-2 dark:text-white">Stock optimal</label>
						<div class="relative">
							<input
								type="number"
								name="price"
								placeholder="Stock opti"
								class="py-3 px-4 block w-[90%] border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
								required
								aria-describedby="text-error"
								on:keyup={(e) => {
									// @ts-ignore
									steps.optimal_amount = parseInt(e.target?.value);
								}}
							/>
						</div>
					</div>
				</div>

				<div class="block">
					<label for="price" class="block text-sm mb-2 dark:text-white">Image</label>
					<div class="relative">
						<input
							type="file"
							id="image"
							name="image"
							accept=".jpg, .jpeg, .png, .webp"
							class="py-3 px-4 block w-full border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
							required
							aria-describedby="text-error"
							on:change={(e) => {
								// @ts-ignore
								let file = e.target?.files[0];
								file2Base64(file).then((res) => {
									res = res.replace('data:', '').replace(/^.+,/, '');
									steps.picture = res;
								});
							}}
						/>
					</div>
				</div>

				<div class="text-xl text-center font-bold p-5">Optionnel</div>

				<div class="grid grid-cols-2 justify-center">
					<!-- Optional available_from and available_to (being seconds since day start) -->
					<div class="block">
						<label for="available_from" class="block text-sm mb-2 dark:text-white"
							>Disponible à partir de</label
						>
						<div class="relative">
							<!-- hour and minute input -->
							<input
								type="time"
								name="available_from"
								placeholder="Disponible à partir de"
								class="py-3 px-4 block w-[90%] border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
								required
								aria-describedby="text-error"
								on:change={(e) => {
									// @ts-ignore
									steps.available_from = e.target?.valueAsNumber / 1000;
								}}
							/>
						</div>
					</div>
					<div class="block">
						<label for="available_from" class="block text-sm mb-2 dark:text-white"
							>Disponible jusqu'à</label
						>
						<div class="relative">
							<!-- hour and minute input -->
							<input
								type="time"
								name="available_from"
								placeholder="Disponible à partir de"
								class="py-3 px-4 block w-[90%] border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
								required
								aria-describedby="text-error"
								on:change={(e) => {
									// @ts-ignore
									steps.available_to = e.target?.valueAsNumber / 1000;
								}}
							/>
						</div>
					</div>
				</div>

				<div class="flex flex-row justify-center mt-5">
					<button
						class="w-32 flex flex-col flex-shrink-0 bg-blue-100 hover:bg-blue-200 rounded-sm p-5 text-center text-lg"
						on:click={() => {
							currentStep = 5;
							setTimeout(() => {
								document.getElementById('end')?.scrollIntoView({ behavior: 'smooth' });
							}, 100);
						}}
					>
						Suivant
					</button>
				</div>
			</div>
		{/if}
		{#if currentStep >= 5}
			<div id="end" class="bg-slate-100 flex flex-col rounded-xl border-2 p-4 mb-20">
				<div class="text-3xl text-center font-bold p-5">Récapitulatif</div>

				<div class="flex flex-row justify-center">
					<div class="w-1/2 rounded-lg p-2 border-2 border-slate-200">
						<div class="flex flex-row justify-between">
							<div class="text-lg font-bold">Nom du menu</div>
							<div class="text-lg">{steps.name}</div>
						</div>
						<div class="flex flex-row justify-between">
							<div class="text-lg font-bold">Catégorie</div>
							<div class="text-lg">{steps.category?.name}</div>
						</div>
						<div class="flex flex-row justify-between">
							<div class="text-lg font-bold">Produits</div>
							<div class="text-lg">
								{#each steps.items ?? [] as i}
									<div>{i.name}</div>
								{/each}
							</div>
						</div>
					</div>
				</div>

				<div class="flex flex-row justify-center mt-5">
					<button
						class="w-32 flex flex-col flex-shrink-0 bg-blue-100 hover:bg-blue-200 rounded-sm p-5 text-center text-lg"
						on:click={() => {
							createMenu();
						}}
					>
						Créer
					</button>
				</div>
			</div>
		{/if}
	</div>
</div>

 <script lang="ts">
	import type { Account, NewAccount, NewCategory } from '$lib/api';
	import Refills from '$lib/components/admin/refills.svelte';
	import Stars from '$lib/components/admin/stars.svelte';
	import NewRefill from '$lib/components/comptoir/newRefill.svelte';
	import ConfirmationPopup from '$lib/components/confirmationPopup.svelte';
	import { accountsApi } from '$lib/requests/requests';
	import { formatPrice } from '$lib/utils';
	import { onMount } from 'svelte';
	import { browser } from '$app/environment';
	import ReadCard from '$lib/components/readCard.svelte';
	import PaginationFooter from '$lib/components/PaginationFooter.svelte';

	let accounts: Account[] = [];
	let newAccount: NewAccount = {
		first_name: '',
		last_name: '',
		email_address: '',
		card_id: '',
		balance: 0,
		role: 'student',
		price_role: 'externe'
	};

	let searchPriceRole: string | undefined = undefined;
	let searchRole: string | undefined = undefined;


	let searchQuery = '';
	let filtersOpen: boolean = false;
	let page: number = 1;
	let maxPage: number = 0;
	let nextPage = () => {
		if (page < maxPage) {
			page++;
			reloadAccounts();
		}
	};
	let prevPage = () => {
		if (page > 1) {
			page--;
			reloadAccounts();
		}
	};
	let handlePageInput = () => {
		if (page < 1) {
			page = 1;
		} else if (page > maxPage) {
			page = maxPage;
		}
		reloadAccounts();
	};
	let accounts_per_page = 10;
	let shown_refill: Account | undefined = undefined;
	let shown_stars: Account | undefined = undefined;
	let recharging_account: Account | undefined = undefined;

	let deletingAccount: boolean = false;
	let confirmationMessage: string | undefined = undefined;
	let deleteAccountCallback: VoidFunction = () => {};
	let selectedAccount: Account | undefined = undefined;
	let askForCard = false;
	onMount(() => {
		reloadAccounts();

		if (browser) {
			// Close dropdown menus when clicking outside
			document.addEventListener('click', () => {
				const dropdownMenus = document.querySelectorAll('.dropdown-menu');
				dropdownMenus.forEach((menu) => {
					if (!menu.classList.contains('hidden')) {
						menu.classList.add('hidden');
						menu.classList.remove('flex');
						menu.classList.remove('flex-col');
					}
				});
			});
		}
	});

	function reloadAccounts() {
		accountsApi()
			.getAccounts(page, accounts_per_page, searchQuery, searchPriceRole, searchRole, { withCredentials: true })
			.then((res) => {
				accounts = res.data.accounts ?? [];
				page = res.data.page;
				maxPage = res.data.max_page;
				accounts_per_page = res.data.limit;
			});
	}

	function createNewAccount() {
		if (!newAccount) return;
		accountsApi()
			.postAccounts(newAccount, { withCredentials: true })
			.then((res) => {
				accounts = [...accounts, res.data];
				newAccount = {
					first_name: '',
					last_name: '',
					email_address: '',
					card_id: '',
					balance: 0,
					role: 'student',
					price_role: 'ceten'
				};
			});
	}

	function deleteAccount(id: string) {
		accountsApi()
			.markDeleteAccountId(id, { withCredentials: true })
			.then(() => {
				accounts = accounts.filter((ct) => ct.id !== id);
			});
	}

	function importAccounts(e: any) {
		const csv = e.target.files[0];

		accountsApi()
			.importAccounts(csv, { withCredentials: true })
			.then(() => {
				reloadAccounts();
			});
	}

	function changeCardId(account: Account | undefined, card_id: string) {
		if (account != undefined) {
			accountsApi()
				.patchAccountId(
					account.id,
					{
						card_id: card_id
					},
					{ withCredentials: true }
				)
				.then((res) => {
					account = res.data ?? account;
				})
				.catch((err) => {
					if (account != undefined) {
						account.card_id = account.card_id ?? '';
					}
				});
		}
	}

	function reset() {
		askForCard = false;
		selectedAccount = undefined;
	}
</script>

{#if shown_refill}
	<Refills
		account={shown_refill}
		onClose={() => {
			shown_refill = undefined;
		}}
	/>
{/if}

{#if shown_stars}
	<Stars
		account={shown_stars}
		close={() => {
			shown_stars = undefined;
		}}
		onStarsAdded={reloadAccounts}
	/>
{/if}

{#if recharging_account}
	<NewRefill
		close={() => {
			recharging_account = undefined;
			reloadAccounts();
		}}
		cardId={recharging_account.card_id}
	/>
{/if}

{#if askForCard}
	<!-- Popup overlay -->
	<button
		id="overlay"
		class="fixed w-full h-full top-0 left-0 bg-black bg-opacity-50 flex justify-center items-center z-50 hover:cursor-default"
		on:click={() => {
			reset();
		}}
	/>

	<div id="popup" class="fixed w-full h-full top-0 left-0 flex justify-center items-center z-50">
		<div
			class="relative flex flex-col justify-center items-center gap-4 p-10 h-96 bg-white dark:bg-slate-900 text-black dark:text-white rounded-xl shadow-xl z-60"
		>
			<button
				class="absolute top-0 right-0 p-2 text-xl font-bold m-2 rounded-full transition-all text-black dark:text-white"
				on:click={() => {
					reset();
				}}
			>
				<iconify-icon icon="mdi:close" />
			</button>
			<h1 class="text-3xl">Veuillez scanner la carte.</h1>
		</div>
	</div>

	<ReadCard
		callback={(id) => {
			changeCardId(selectedAccount, id);
			reset();
		}}
	/>
{/if}
<!-- Popup -->
<div
	id="hs-modal-new-account"
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
						Ajouter un compte
					</h2>
				</div>

			

				<div class="mt-5">
					<!-- Form -->
					<div class="grid gap-y-4">
						<!-- Form Group -->
						<div>
							<!-- name -->
							<label for="first_name" class="block text-sm mb-2 dark:text-white">Prénom</label>
							<div class="relative">
								<input
									type="text"
									id="first_name"
									name="first_name"
									placeholder="Prénom"
									class="py-3 px-4 block w-full border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
									required
									aria-describedby="text-error"
									bind:value={newAccount.first_name}
								/>
							</div>

							<label for="last_name" class="block text-sm mb-2 dark:text-white">Nom</label>
							<div class="relative mt-3">
								<input
									type="text"
									id="last_name"
									name="last_name"
									placeholder="Nom"
									class="py-3 px-4 block w-full border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
									required
									aria-describedby="text-error"
									bind:value={newAccount.last_name}
								/>
							</div>

							<label for="email_address" class="block text-sm mb-2 dark:text-white"
								>Adresse E-Mail</label
							>
							<div class="relative mt-3">
								<input
									type="text"
									id="email_address"
									name="email_address"
									placeholder="Adresse email"
									class="py-3 px-4 block w-full border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
									required
									aria-describedby="text-error"
									bind:value={newAccount.email_address}
								/>
							</div>

							<label for="card_id" class="block text-sm mb-2 dark:text-white"
								>Identifiant de la carte</label
							>
							<div class="relative mt-3">
								<input
									type="text"
									id="card_id"
									name="card_id"
									placeholder="ID de la carte"
									class="py-3 px-4 block w-full border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
									required
									aria-describedby="text-error"
									bind:value={newAccount.card_id}
								/>
							</div>

							<button
								type="submit"
								class="mt-4 py-3 px-4 inline-flex justify-center items-center gap-2 rounded-md border border-transparent font-semibold bg-blue-500 text-white hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-all text-sm dark:focus:ring-offset-gray-800"
								on:click={() => createNewAccount()}
								data-hs-overlay="#hs-modal-new-account">Créer</button
							>
						</div>
					</div>
					<!-- End Form -->
				</div>
			</div>
		</div>
	</div>
</div>

{#if deletingAccount}
	<ConfirmationPopup
		message={confirmationMessage}
		confirm_text="Supprimer"
		cancel_callback={() => {
			deletingAccount = false;
		}}
		confirm_callback={deleteAccountCallback}
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
					{#if searchPriceRole || searchRole || searchQuery}
						<span class="bg-blue-500 text-white text-xs rounded-full px-1.5 py-0.5 min-w-[18px]">
							{(searchPriceRole ? 1 : 0) + (searchRole ? 1 : 0) + (searchQuery ? 1 : 0)}
						</span>
					{/if}
					<iconify-icon icon={filtersOpen ? "mdi:chevron-up" : "mdi:chevron-down"} width="18" height="18"></iconify-icon>
				</button>
				<button
					class="py-2 px-3 inline-flex justify-center items-center gap-2 rounded-md border border-transparent font-semibold bg-blue-500 text-white hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-all text-sm dark:focus:ring-offset-gray-800"
					data-hs-overlay="#hs-modal-new-account"
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
					<div>
						<label class="text-xs font-medium text-gray-500 dark:text-gray-400 block mb-1">Rechercher</label>
						<input
							type="text"
							placeholder="Rechercher par nom..."
							bind:value={searchQuery}
							on:input={() => {
								page = 1;
								reloadAccounts();
							}}
							class="w-full px-2 py-1.5 text-sm bg-gray-50 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:text-white placeholder-gray-400 dark:placeholder-gray-500"
						/>
					</div>
					<div class="grid grid-cols-2 gap-3">
						<div>
							<label class="text-xs font-medium text-gray-500 dark:text-gray-400 block mb-1">Prix</label>
							<select
								id="category-mobile"
								name="category"
								class="w-full px-2 py-1.5 text-sm bg-gray-50 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:text-white"
								on:change={(e) => {
									// @ts-ignore
									searchPriceRole = e.target?.value;
									if (searchPriceRole === '') searchPriceRole = undefined;
									page = 1;
									reloadAccounts();
								}}
							>
								<option value="">Tous</option>
								<option value="ceten">CETEN</option>
								<option value="coutant">Coutant</option>
								<option value="staff_bar">Staff</option>
								<option value="externe">Externe</option>
								<option value="privilegies">Privilégié</option>
							</select>
						</div>
						<div>
							<label class="text-xs font-medium text-gray-500 dark:text-gray-400 block mb-1">Rôle</label>
							<select
								id="category-mobile"
								name="category"
								class="w-full px-2 py-1.5 text-sm bg-gray-50 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:text-white"
								on:change={(e) => {
									// @ts-ignore
									searchRole = e.target?.value;
									if (searchRole === '') searchRole = undefined;
									page = 1;
									reloadAccounts();
								}}
							>
								<option value="">Tous</option>
								<option value="student">Étudiant</option>
								<option value="student_with_benefits">Avec avantages</option>
								<option value="member">Membre</option>
								<option value="admin">Admin</option>
								<option value="ghost">Fantôme</option>
								<option value="superadmin">SuperAdmin</option>
							</select>
						</div>
					</div>
				</div>
			{/if}
		</div>

		<!-- Desktop: Original layout -->
		<div class="hidden lg:flex flex-row flex-wrap items-center gap-6">
			<div class="flex flex-row items-center gap-3">
				<span class="text-sm font-medium text-gray-700 dark:text-gray-300">Rechercher:</span>
				<input
					type="text"
					placeholder="Rechercher par nom..."
					bind:value={searchQuery}
					on:input={() => {
						page = 1;
						reloadAccounts();
					}}
					class="px-3 py-1.5 text-sm bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:text-white placeholder-gray-400 dark:placeholder-gray-500 w-auto"
				/>
			</div>
			<div class="flex flex-row items-center gap-3">
				<!-- Import -->
				<input
					type="file"
					id="import"
					name="import"
					class="hidden"
					accept=".csv"
					on:change={(e) => importAccounts(e)}
				/>
				<button
					class="py-2 px-3 inline-flex justify-center items-center gap-2 rounded-md border border-transparent font-semibold bg-blue-500 text-white hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-all text-sm dark:focus:ring-offset-gray-800"
					on:click={() => {
						// @ts-ignore
						document.getElementById('import').click();
					}}
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
					Importer des Comptes
				</button>
			</div>
			<div class="flex flex-row items-center gap-3">
				<button
					class="py-2 px-3 inline-flex justify-center items-center gap-2 rounded-md border border-transparent font-semibold bg-blue-500 text-white hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-all text-sm dark:focus:ring-offset-gray-800"
					data-hs-overlay="#hs-modal-new-account"
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
					Ajouter un compte
				</button>
			</div>
			<div class="flex flex-row items-center gap-3">
				<span class="text-sm font-medium text-gray-700 dark:text-gray-300">Prix:</span>
				<select
					id="category"
					name="category"
					class="px-3 py-1.5 text-sm bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:text-white dark:placeholder-gray-500 w-auto"
					on:change={(e) => {
						// @ts-ignore
						searchPriceRole = e.target?.value;
						if (searchPriceRole === '') searchPriceRole = undefined;
						page = 1;
						reloadAccounts();
					}}
				>
					<option value="">Pas de filtre</option>
					<option value="ceten">Prix CETEN</option>
					<option value="coutant">Prix coutant</option>
					<option value="staff_bar">Prix staff</option>
					<option value="externe">Prix externe</option>
					<option value="privilegies">Prix privilégié</option>
				</select>
			</div>
			<div class="flex flex-row items-center gap-3">
				<span class="text-sm font-medium text-gray-700 dark:text-gray-300">Rôle:</span>
				<select
					id="category"
					name="category"
					class="px-3 py-1.5 text-sm bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:text-white dark:placeholder-gray-500 w-auto"
					on:change={(e) => {
						// @ts-ignore
						searchRole = e.target?.value;
						if (searchRole === '') searchRole = undefined;
						page = 1;
						reloadAccounts();
					}}
				>
					<option value="">Pas de filtre</option>
					<option value="student">Étudiant</option>
					<option value="student_with_benefits">Étudiant avec avantages</option>
					<option value="member">Membre</option>
					<option value="admin">Admin</option>
					<option value="ghost">Fantôme</option>
					<option value="superadmin">SuperAdmin</option>
				</select>
			</div>
		</div>
	</div>

	<div class="flex-1 min-h-0 w-full overflow-x-auto overflow-y-auto">
		<!-- Desktop Table View -->
		<div class="hidden min-[1300px]:block min-w-full bg-white dark:bg-gray-800 rounded-lg shadow-sm overflow-visible">
			<table class="min-w-full table-fixed divide-y divide-gray-200 dark:divide-gray-700">
				<colgroup>
					<col class="w-[12%]" />
					<col class="w-[12%]" />
					<col class="w-[18%]" />
					<col class="w-[10%]" />
					<col class="w-[10%]" />
					<col class="w-[12%]" />
					<col class="w-[12%]" />
					<col class="w-[14%]" />
				</colgroup>
				<thead class="bg-gray-50 dark:bg-gray-700">
					<tr>
						<th class="px-6 py-4 text-left text-sm font-semibold text-gray-900 dark:text-gray-300 border-r border-gray-200 dark:border-gray-600">
							Nom
						</th>
						<th class="px-6 py-4 text-left text-sm font-semibold text-gray-900 dark:text-gray-300 border-r border-gray-200 dark:border-gray-600">
							Prénom
						</th>
						<th class="px-6 py-4 text-left text-sm font-semibold text-gray-900 dark:text-gray-300 border-r border-gray-200 dark:border-gray-600">
							Adresse E-Mail
						</th>
						<th class="px-6 py-4 text-left text-sm font-semibold text-gray-900 dark:text-gray-300 border-r border-gray-200 dark:border-gray-600">
							Solde
						</th>
						<th class="px-6 py-4 text-left text-sm font-semibold text-gray-900 dark:text-gray-300 border-r border-gray-200 dark:border-gray-600">
							Étoiles
						</th>
						<th class="px-6 py-4 text-left text-sm font-semibold text-gray-900 dark:text-gray-300 border-r border-gray-200 dark:border-gray-600">
							Rôle
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
					{#each accounts as account}
						<tr>
							<td class="px-6 py-4 text-sm text-gray-900 dark:text-gray-300 border-r border-gray-200 dark:border-gray-600 group relative">
								<input
									type="text"
									class="block w-full text-sm dark:text-white/[.8] p-2 bg-transparent border-none outline-none"
									value={account.last_name}
									on:input={(e) => {
										// @ts-ignore
										let name = e.target?.value;
										accountsApi()
											.patchAccountId(
												account.id,
												{
													last_name: name
												},
												{ withCredentials: true }
											)
											.then((res) => {
												account = res.data ?? account;
											})
											.catch((err) => {
												account.last_name = account.last_name ?? '';
											});
									}}
								/>
								<iconify-icon 
									icon="mdi:pencil" 
									class="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400 opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none"
								/>
							</td>
							<td class="px-6 py-4 text-sm text-gray-900 dark:text-gray-300 border-r border-gray-200 dark:border-gray-600 group relative">
								<input
									type="text"
									class="block w-full text-sm dark:text-white/[.8] p-2 bg-transparent border-none outline-none"
									value={account.first_name}
									on:input={(e) => {
										// @ts-ignore
										let name = e.target?.value;
										accountsApi()
											.patchAccountId(
												account.id,
												{
													first_name: name
												},
												{ withCredentials: true }
											)
											.then((res) => {
												account = res.data ?? account;
											})
											.catch((err) => {
												account.first_name = account.first_name ?? '';
											});
									}}
								/>
								<iconify-icon 
									icon="mdi:pencil" 
									class="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400 opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none"
								/>
							</td>
							<td class="px-6 py-4 text-sm text-gray-900 dark:text-gray-300 border-r border-gray-200 dark:border-gray-600 group relative">
								<input
									type="text"
									class="block w-full text-sm dark:text-white/[.8] p-2 bg-transparent border-none outline-none"
									value={account.email_address}
									on:input={(e) => {
										// @ts-ignore
										let name = e.target?.value;
										accountsApi()
											.patchAccountId(
												account.id,
												{
													email_address: name
												},
												{ withCredentials: true }
											)
											.then((res) => {
												account = res.data ?? account;
											})
											.catch((err) => {
												account.email_address = account.email_address ?? '';
											});
									}}
								/>
								<iconify-icon 
									icon="mdi:pencil" 
									class="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400 opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none"
								/>
							</td>
							<td class="px-6 py-4 text-sm font-medium text-gray-900 dark:text-gray-300 border-r border-gray-200 dark:border-gray-600">
								{formatPrice(account.balance)}
							</td>
							<td class="px-6 py-4 text-sm font-medium text-gray-900 dark:text-gray-300 border-r border-gray-200 dark:border-gray-600">
								{#if account.points >= 1000000000}
									{(account.points / 1000000000).toFixed(account.points % 1000000000 === 0 ? 0 : 1)}G
								{:else if account.points >= 1000000}
									{(account.points / 1000000).toFixed(account.points % 1000000 === 0 ? 0 : 1)}M
								{:else if account.points >= 1000}
									{(account.points / 1000).toFixed(account.points % 1000 === 0 ? 0 : 1)}k
								{:else}
									{account.points}
								{/if}
								<iconify-icon icon="mdi:star" class="ml-1 text-yellow-500" />
							</td>
							<td class="px-6 py-4 text-sm text-gray-900 dark:text-gray-300 border-r border-gray-200 dark:border-gray-600">
								<select
									class="block w-full text-sm dark:text-white/[.8] dark:bg-slate-900 p-2 bg-transparent border-none outline-none"
									value={account.role}
									on:change={(e) => {
										// @ts-ignore
										let role = e.target?.value;
										accountsApi()
											.patchAccountId(
												account.id,
												{
													role: role
												},
												{ withCredentials: true }
											)
											.then((res) => {
												account = res.data ?? account;
											})
											.catch((err) => {
												account.role = account.role ?? '';
											});
									}}
								>
									<option value="student">Étudiant</option>
									<option value="student_with_benefits">Étudiant avec avantages</option>
									<option value="member">Membre</option>
									<option value="admin">Admin</option>
									<option value="ghost">Fantôme</option>
									<option value="superadmin">Superadmin</option>
								</select>
							</td>
							<td class="px-6 py-4 text-sm text-gray-900 dark:text-gray-300 border-r border-gray-200 dark:border-gray-600">
								<select
									class="block w-full text-sm dark:text-white/[.8] dark:bg-slate-900 p-2 bg-transparent border-none outline-none"
									value={account.price_role}
									on:change={(e) => {
										// @ts-ignore
										let role = e.target?.value;
										accountsApi()
											.patchAccountId(
												account.id,
												{
													price_role: role
												},
												{ withCredentials: true }
											)
											.then((res) => {
												account = res.data ?? account;
											})
											.catch((err) => {
												account.price_role = account.price_role ?? '';
											});
									}}
								>
									<option value="externe">Externe</option>
									<option value="ceten">CETEN</option>
									<option value="staff_bar">Staff</option>
									<option value="coutant">Coutant</option>
									<option value="privilegies">Membre privilégié</option>
								</select>
							</td>
							<td class="px-6 py-4 text-sm text-gray-900 dark:text-gray-300 relative">
								<button
									class="text-sm text-blue-600 dark:text-blue-400 font-medium hover:bg-gray-100 dark:hover:bg-slate-800 p-2 rounded-md flex items-center gap-2"
									on:click={(e) => {
										const menu = e.currentTarget.nextElementSibling;
										if (menu) {
											const rect = e.currentTarget.getBoundingClientRect();
											const menuHeight = 200; // Approximate dropdown height
											const viewportHeight = window.innerHeight;
											const spaceBelow = viewportHeight - rect.bottom;
											const spaceAbove = rect.top;
											
											// @ts-ignore
											menu.style.left = `${rect.left}px`;
											
											// Position above if not enough space below and enough space above
											if (spaceBelow < menuHeight && spaceAbove > menuHeight) {
												// @ts-ignore
												menu.style.top = `${rect.top - menuHeight - 8}px`;
											} else {
												// @ts-ignore
												menu.style.top = `${rect.bottom + 8}px`;
											}
											
											menu.classList.toggle('hidden');
											menu.classList.toggle('flex');
											menu.classList.toggle('flex-col');
										}
										e.stopPropagation();
									}}
								>
									Actions
									<svg
										xmlns="http://www.w3.org/2000/svg"
										width="16"
										height="16"
										viewBox="0 0 24 24"
										fill="none"
										stroke="currentColor"
										stroke-width="2"
										stroke-linecap="round"
										stroke-linejoin="round"
									>
										<path d="m6 9 6 6 6-6" />
									</svg>
								</button>

								<!-- Dropdown menu -->
								<div
									class="dropdown-menu hidden fixed mt-2 py-2 w-48 bg-white dark:bg-slate-900 rounded-md shadow-lg z-[100] border border-gray-200 dark:border-gray-700"
								>
									{#if askForCard == false}
										<button
											class="w-full text-left px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-slate-800 flex items-center gap-2"
											on:click={() => {
												selectedAccount = account;
												askForCard = true;
											}}
										>
											<iconify-icon icon="mdi:card-account-details" width="20" height="20" />
											Nouvelle Carte
										</button>
									{/if}
									<button
										class="w-full text-left px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-slate-800 flex items-center gap-2"
										on:click={() => (shown_refill = account)}
									>
										<iconify-icon icon="mdi:history" width="20" height="20" />
										Transactions
									</button>
									<button
										class="w-full text-left px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-slate-800 flex items-center gap-2"
										on:click={() => (recharging_account = account)}
									>
										<iconify-icon icon="mdi:wallet-plus" width="20" height="20" />
										Recharger
									</button>
									<button
										class="w-full text-left px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-slate-800 flex items-center gap-2"
										on:click={() => {
											shown_stars = account;
										}}
									>
										<iconify-icon icon="mdi:star-plus" width="20" height="20" />
										Ajouter des étoiles
									</button>
									<button
										class="w-full text-left px-4 py-2 text-sm text-red-700 dark:text-red-300 hover:bg-red-100 dark:hover:bg-red-800 flex items-center gap-2"
										on:click={() => {
											deleteAccountCallback = () => {
												deletingAccount = false;
												deleteAccount(account.id);
											};
											confirmationMessage =
												'Supprimer le compte de ' +
												account.first_name +
												' ' +
												account.last_name +
												' ?';
											deletingAccount = true;
										}}
									>
										<iconify-icon icon="mdi:delete" width="20" height="20" />
										Supprimer
									</button>
								</div>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>

		<!-- Mobile Card View -->
		<div class="block min-[1300px]:hidden space-y-4 px-2 pb-4">
			{#if accounts.length === 0}
				<div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-8 text-center">
					<p class="text-gray-500 dark:text-gray-400">Aucun compte trouvé</p>
				</div>
			{:else}
				{#each accounts as account}
					<div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-4">
						<div class="flex justify-between items-start mb-3">
							<div class="flex-1">
								<div class="grid grid-cols-2 gap-2 mb-2">
									<div>
										<label class="text-xs font-medium text-gray-500 dark:text-gray-400 block">Nom</label>
										<input
											type="text"
											class="w-full text-sm dark:text-white/[.8] bg-transparent border border-gray-200 dark:border-gray-600 rounded p-2 focus:border-blue-500 focus:outline-none"
											value={account.last_name}
											on:input={(e) => {
												// @ts-ignore
												let name = e.target?.value;
												accountsApi()
													.patchAccountId(
														account.id,
														{
															last_name: name
														},
														{ withCredentials: true }
													)
													.then((res) => {
														account = res.data ?? account;
													})
													.catch((err) => {
														account.last_name = account.last_name ?? '';
													});
											}}
										/>
									</div>
									<div>
										<label class="text-xs font-medium text-gray-500 dark:text-gray-400 block">Prénom</label>
										<input
											type="text"
											class="w-full text-sm dark:text-white/[.8] bg-transparent border border-gray-200 dark:border-gray-600 rounded p-2 focus:border-blue-500 focus:outline-none"
											value={account.first_name}
											on:input={(e) => {
												// @ts-ignore
												let name = e.target?.value;
												accountsApi()
													.patchAccountId(
														account.id,
														{
															first_name: name
														},
														{ withCredentials: true }
													)
													.then((res) => {
														account = res.data ?? account;
													})
													.catch((err) => {
														account.first_name = account.first_name ?? '';
													});
											}}
										/>
									</div>
								</div>
								<div class="mb-2">
									<label class="text-xs font-medium text-gray-500 dark:text-gray-400 block">E-mail</label>
									<input
										type="text"
										class="w-full text-sm dark:text-white/[.8] bg-transparent border border-gray-200 dark:border-gray-600 rounded p-2 focus:border-blue-500 focus:outline-none"
										value={account.email_address}
										on:input={(e) => {
											// @ts-ignore
											let name = e.target?.value;
											accountsApi()
												.patchAccountId(
													account.id,
													{
														email_address: name
													},
													{ withCredentials: true }
												)
												.then((res) => {
													account = res.data ?? account;
												})
												.catch((err) => {
													account.email_address = account.email_address ?? '';
												});
										}}
									/>
								</div>
							</div>
						</div>
						<div class="grid grid-cols-2 sm:grid-cols-4 gap-2 mb-3">
						<div>
							<label class="text-xs font-medium text-gray-500 dark:text-gray-400 block">Solde</label>
							<div class="text-sm font-medium text-gray-900 dark:text-gray-300 py-2">
								{formatPrice(account.balance)}
							</div>
						</div>
						<div>
							<label class="text-xs font-medium text-gray-500 dark:text-gray-400 block">Étoiles</label>
							<div class="text-sm font-medium text-gray-900 dark:text-gray-300 py-2 flex items-center">
								{#if account.points >= 1000000000}
									{(account.points / 1000000000).toFixed(account.points % 1000000000 === 0 ? 0 : 1)}G
								{:else if account.points >= 1000000}
									{(account.points / 1000000).toFixed(account.points % 1000000 === 0 ? 0 : 1)}M
								{:else if account.points >= 1000}
									{(account.points / 1000).toFixed(account.points % 1000 === 0 ? 0 : 1)}k
								{:else}
									{account.points}
								{/if}
								<iconify-icon icon="mdi:star" class="ml-1 text-yellow-500" />
							</div>
						</div>
						<div>
							<label class="text-xs font-medium text-gray-500 dark:text-gray-400 block">Rôle</label>
							<select
									class="w-full text-sm dark:text-white/[.8] dark:bg-slate-900 bg-transparent border border-gray-200 dark:border-gray-600 rounded p-2 focus:border-blue-500 focus:outline-none"
									value={account.role}
									on:change={(e) => {
										// @ts-ignore
										let role = e.target?.value;
										accountsApi()
											.patchAccountId(
												account.id,
												{
													role: role
												},
												{ withCredentials: true }
											)
											.then((res) => {
												account = res.data ?? account;
											})
											.catch((err) => {
												account.role = account.role ?? '';
											});
									}}
								>
									<option value="student">Étudiant</option>
									<option value="student_with_benefits">Avec avantages</option>
									<option value="member">Membre</option>
									<option value="admin">Admin</option>
									<option value="ghost">Fantôme</option>
									<option value="superadmin">Superadmin</option>
								</select>
							</div>
							<div>
								<label class="text-xs font-medium text-gray-500 dark:text-gray-400 block">Prix</label>
								<select
									class="w-full text-sm dark:text-white/[.8] dark:bg-slate-900 bg-transparent border border-gray-200 dark:border-gray-600 rounded p-2 focus:border-blue-500 focus:outline-none"
									value={account.price_role}
									on:change={(e) => {
										// @ts-ignore
										let role = e.target?.value;
										accountsApi()
											.patchAccountId(
												account.id,
												{
													price_role: role
												},
												{ withCredentials: true }
											)
											.then((res) => {
												account = res.data ?? account;
											})
											.catch((err) => {
												account.price_role = account.price_role ?? '';
											});
									}}
								>
									<option value="externe">Externe</option>
									<option value="ceten">CETEN</option>
									<option value="staff_bar">Staff</option>
									<option value="coutant">Coutant</option>
									<option value="privilegies">Privilégié</option>
								</select>
							</div>
						</div>
						<div class="flex flex-wrap gap-2">
							{#if askForCard == false}
								<button
									class="flex-1 min-w-0 px-3 py-2 text-xs bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors flex items-center justify-center gap-1"
									on:click={() => {
										selectedAccount = account;
										askForCard = true;
									}}
								>
									<iconify-icon icon="mdi:card-account-details" width="16" height="16" />
									Carte
								</button>
							{/if}
							<button
								class="flex-1 min-w-0 px-3 py-2 text-xs bg-gray-500 text-white rounded hover:bg-gray-600 transition-colors flex items-center justify-center gap-1"
								on:click={() => (shown_refill = account)}
							>
								<iconify-icon icon="mdi:history" width="16" height="16" />
								Transactions
							</button>
							<button
								class="flex-1 min-w-0 px-3 py-2 text-xs bg-green-500 text-white rounded hover:bg-green-600 transition-colors flex items-center justify-center gap-1"
								on:click={() => (recharging_account = account)}
							>
								<iconify-icon icon="mdi:wallet-plus" width="16" height="16" />
								Recharger
							</button>
							<button
								class="flex-1 min-w-0 px-3 py-2 text-xs bg-yellow-500 text-white rounded hover:bg-yellow-600 transition-colors flex items-center justify-center gap-1"
								on:click={() => {
									shown_stars = account;
								}}
							>
								<iconify-icon icon="mdi:star-plus" width="16" height="16" />
								Étoiles
							</button>
							<button
								class="flex-1 min-w-0 px-3 py-2 text-xs bg-red-500 text-white rounded hover:bg-red-600 transition-colors flex items-center justify-center gap-1"
								on:click={() => {
									deleteAccountCallback = () => {
										deletingAccount = false;
										deleteAccount(account.id);
									};
									confirmationMessage =
										'Supprimer le compte de ' +
										account.first_name +
										' ' +
										account.last_name +
										' ?';
									deletingAccount = true;
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
		resultsCount={accounts.length}
		showPageInput={true}
		on:prevPage={prevPage}
		on:nextPage={nextPage}
		on:pageChange={handlePageInput}
	/>
</div>

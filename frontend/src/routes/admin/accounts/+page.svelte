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

	let accounts: Account[] = [];
	let newAccount: NewAccount = {
		first_name: '',
		last_name: '',
		email_address: '',
		card_id: '',
		balance: 0,
		role: 'student',
		price_role: 'ceten'
	};

	let searchQuery = '';
	let page: number = 0;
	let maxPage: number = 0;
	let nextPage = () => {
		if (page < maxPage) {
			page++;
			reloadAccounts();
		}
	};
	let prevPage = () => {
		if (page > 0) {
			page--;
			reloadAccounts();
		}
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
				const dropdownMenus = document.querySelectorAll('.fixed.mt-2');
				dropdownMenus.forEach((menu) => {
					if (!menu.classList.contains('hidden')) {
						menu.classList.add('hidden');
					}
				});
			});
		}
	});

	function reloadAccounts() {
		accountsApi()
			.getAccounts(page, accounts_per_page, searchQuery, { withCredentials: true })
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

<div class="h-[calc(100vh-80px)] grid grid-cols-1 grid-rows-[auto_1fr_80px] bg-gray-50 dark:bg-gray-900">
	<div class="m-3 p-2">
		<div class="flex flex-wrap items-center gap-6">
			<div class="flex items-center gap-3">
				<span class="text-sm font-medium text-gray-700 dark:text-gray-300">Rechercher:</span>
				<input
					type="text"
					placeholder="Rechercher par nom..."
					bind:value={searchQuery}
					on:input={() => {
						page = 0;
						reloadAccounts();
					}}
					class="px-3 py-1.5 text-sm bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:text-white placeholder-gray-400 dark:placeholder-gray-500"
				/>
			</div>
			<div class="flex items-center gap-3">
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
					Importer des Comptes
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
					Ajouter un compte
				</button>
			</div>
		</div>
	</div>

	<div class="flex-grow w-full overflow-x-auto overflow-y-visible">
		<div class="min-w-full bg-white dark:bg-gray-800 rounded-lg shadow-sm overflow-visible">
			<div class="grid grid-cols-[1fr_1fr_1.5fr_0.8fr_0.8fr_0.8fr_1fr] bg-gray-50 dark:bg-gray-700 divide-x divide-gray-200 dark:divide-gray-700">
				<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
					Nom
				</th>
				<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
					Prénom
				</th>
				<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
					Adresse E-Mail
				</th>
				<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
					Solde
				</th>
				<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
					Rôle
				</th>
				<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
					Prix
				</th>
				<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
					Actions
				</th>
			</div>
			<div class="divide-y divide-gray-200 dark:divide-gray-700">
				{#each accounts as account}
					<div class="grid grid-cols-[1fr_1fr_1.5fr_0.8fr_0.8fr_0.8fr_1fr] divide-x divide-gray-200 dark:divide-gray-700">
						<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
							<input
								type="text"
								class="block text-sm dark:text-white/[.8] break-words p-2 bg-transparent border-none outline-none w-full"
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
						</td>
						<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
							<input
								type="text"
								class="block text-sm dark:text-white/[.8] break-words p-2 bg-transparent border-none outline-none w-full"
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
						</td>
						<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
							<input
								type="text"
								class="block text-sm dark:text-white/[.8] break-words p-2 bg-transparent border-none outline-none w-full"
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
						</td>
						<td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-gray-300">
							{formatPrice(account.balance)}
						</td>
						<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
							<select
								class="block text-sm dark:text-white/[.8] dark:bg-slate-900 break-words p-2 bg-transparent border-none outline-none w-full"
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
						<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
							<select
								class="block text-sm dark:text-white/[.8] dark:bg-slate-900 break-words p-2 bg-transparent border-none outline-none w-full"
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
						<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300 relative">
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
								class="hidden fixed mt-2 py-2 w-48 bg-white dark:bg-slate-900 rounded-md shadow-lg z-[100] border border-gray-200 dark:border-gray-700 flex flex-col"
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
					</div>
				{/each}
			</div>
		</div>
	</div>

	<!-- Pagination -->
	<div class="h-20 px-6 py-4 bg-white dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700 flex flex-col sm:flex-row justify-between items-center gap-4">
		<div>
			<p class="text-sm text-gray-600 dark:text-gray-400">
				<span class="font-semibold text-gray-800 dark:text-gray-200">{accounts.length}</span> résultats
			</p>
		</div>
		<div class="flex items-center gap-x-4">
			<button
				type="button"
				class="py-2 px-4 inline-flex justify-center items-center gap-2 rounded-md border font-medium bg-white text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-600 transition-all text-sm disabled:opacity-50 disabled:cursor-not-allowed dark:bg-gray-700 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-600 dark:focus:ring-offset-gray-800"
				on:click={prevPage}
				disabled={page === 0}
			>
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M15 19l-7-7 7-7"
					/>
				</svg>
				Précédent
			</button>
			<p class="text-sm font-medium text-gray-700 dark:text-gray-300">
				Page <span class="font-bold">{page}</span> sur <span class="font-bold">{maxPage}</span>
			</p>
			<button
				type="button"
				class="py-2 px-4 inline-flex justify-center items-center gap-2 rounded-md border font-medium bg-white text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-600 transition-all text-sm disabled:opacity-50 disabled:cursor-not-allowed dark:bg-gray-700 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-600 dark:focus:ring-offset-gray-800"
				on:click={nextPage}
				disabled={page === maxPage}
			>
				Suivant
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
				</svg>
			</button>
		</div>
	</div>
</div>

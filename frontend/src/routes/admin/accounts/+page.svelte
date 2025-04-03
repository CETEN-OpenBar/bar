<script lang="ts">
	import type { Account, NewAccount, NewCategory } from '$lib/api';
	import Refills from '$lib/components/admin/refills.svelte';
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

{#if askForCard}
	<!-- Popup overlay -->
	<button
		id="overlay"
		class="absolute w-full h-full top-0 left-0 bg-black bg-opacity-50 flex justify-center items-center z-10 hover:cursor-default"
		on:click={() => {
			reset();
		}}
	/>

	<div id="popup" class="absolute w-full h-full top-0 left-0 flex justify-center items-center">
		<div
			class="relative text-black flex flex-col justify-center items-center gap-4 p-10 h-96 bg-white rounded-xl shadow-xl z-20"
		>
			<button
				class="absolute top-0 right-0 p-2 text-xl font-bold m-2 rounded-full transition-all text-black"
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
							<h2 class="text-xl font-semibold text-gray-800 dark:text-gray-200">Comptes</h2>
							<p class="text-sm text-gray-600 dark:text-gray-400">Ajouter des comptes</p>
						</div>

						<!-- search bar -->
						<div class="relative mt-4 w-96 md:mt-0">
							<input
								type="text"
								class="py-3 px-4 w-full border-gray-200 border-2 rounded-md text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
								placeholder="Rechercher"
								aria-label="Rechercher"
								on:input={(e) => {
									// @ts-ignore
									searchQuery = e.target.value.toLowerCase();

									reloadAccounts();
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

						<div>
							<div class="inline-flex gap-x-2">
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
											Prénom
										</span>
									</div>
								</th>
								<th scope="col" class="px-6 py-3 text-left">
									<div class="flex items-center gap-x-2">
										<span
											class="text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
										>
											Adresse E-Mail
										</span>
									</div>
								</th>
								<th scope="col" class="px-6 py-3 text-left">
									<div class="flex items-center gap-x-2">
										<span
											class="text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
										>
											Solde
										</span>
									</div>
								</th>
								<th scope="col" class="px-6 py-3 text-left">
									<div class="flex items-center gap-x-2">
										<span
											class="text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
										>
											Rôle
										</span>
									</div>
								</th>
								<th scope="col" class="px-6 py-3 text-left">
									<div class="flex items-center gap-x-2">
										<span
											class="text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
										>
											Prix
										</span>
									</div>
								</th>
								<th scope="col" class="px-6 py-3 text-left">
									<div class="flex items-center gap-x-2">
										<span
											class="text-xs font-semibold uppercase tracking-wide text-gray-800 dark:text-gray-200"
										>
											Actions
										</span>
									</div>
								</th>
								<th scope="col" class="px-6 py-3 text-right" />
							</tr>
						</thead>

						<tbody class="divide-y divide-gray-200 dark:divide-gray-700">
							{#each accounts as account}
								<tr>
									<td class="h-px w-72">
										<div class="px-6 py-3">
											<input
												type="text"
												class="block text-sm dark:text-white/[.8] break-words p-2 bg-transparent"
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
									</td>
									<td class="h-px w-72">
										<div class="px-6 py-3">
											<input
												type="text"
												class="block text-sm dark:text-white/[.8] break-words p-2 bg-transparent"
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
									</td>
									<td class="h-px w-96">
										<div class="px-6 py-3">
											<input
												type="text"
												class="w-72 block text-sm dark:text-white/[.8] break-words p-2 bg-transparent"
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
									</td>
									<td class="h-px w-72">
										<div class="px-6 py-3">
											<p class="text-sm dark:text-white/[.8] break-words p-2 bg-transparent">
												{formatPrice(account.balance)}
											</p>
										</div>
									</td>
									<td class="h-px w-72">
										<div class="px-6 py-3">
											<select
												class="block text-sm dark:text-white/[.8] dark:bg-slate-900 break-words p-2 bg-transparent"
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
										</div>
									</td>
									<td class="h-px w-72">
										<div class="px-6 py-3">
											<select
												class="block text-sm dark:text-white/[.8] dark:bg-slate-900 break-words p-2"
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
										</div>
									</td>
									<td class="h-px w-px whitespace-nowrap relative">
										<div class="px-6 py-1.5">
											<button
												class="text-sm text-blue-600 font-medium hover:bg-gray-100 p-2 rounded-md flex items-center gap-2"
												on:click={(e) => {
													const menu = e.currentTarget.nextElementSibling;
													menu?.classList.toggle('hidden');
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
												class="hidden fixed mt-2 py-2 w-48 bg-white rounded-md shadow-lg z-50 border border-gray-200 flex flex-col"
												style="position: fixed; right: 5%;"
											>
												{#if askForCard == false}
													<button
														class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
														on:click={() => {
															selectedAccount = account;
															askForCard = true;
														}}
													>
														Nouvelle Carte
													</button>
												{/if}
												<button
													class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
													on:click={() => (shown_refill = account)}
												>
													Transactions
												</button>
												<button
													class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
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
													Supprimer
												</button>
												<button
													class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
													on:click={() => {
														console.log('stars');
													}}
												>
													Ajouter des étoiles
												</button>
											</div>
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
								<span class="font-semibold text-gray-800 dark:text-gray-200">{accounts.length}</span
								> résultats
							</p>
						</div>

						<div>
							<div class="inline-flex gap-x-2">
								<button
									type="button"
									class="py-2 px-3 inline-flex justify-center items-center gap-2 rounded-md border font-medium bg-white text-gray-700 shadow-sm align-middle hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-white focus:ring-blue-600 transition-all text-sm dark:bg-slate-900 dark:hover:bg-slate-800 dark:border-gray-700 dark:text-gray-400 dark:hover:text-white dark:focus:ring-offset-gray-800"
									on:click={prevPage}
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
									Page {page} / {maxPage}
								</p>

								<button
									type="button"
									class="py-2 px-3 inline-flex justify-center items-center gap-2 rounded-md border font-medium bg-white text-gray-700 shadow-sm align-middle hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-white focus:ring-blue-600 transition-all text-sm dark:bg-slate-900 dark:hover:bg-slate-800 dark:border-gray-700 dark:text-gray-400 dark:hover:text-white dark:focus:ring-offset-gray-800"
									on:click={nextPage}
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

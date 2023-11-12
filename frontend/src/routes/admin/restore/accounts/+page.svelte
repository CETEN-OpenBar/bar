<script lang="ts">
	import type { Account, NewAccount, NewCategory } from '$lib/api';
	import { api } from '$lib/config/config';
	import { accountsApi, deletedApi } from '$lib/requests/requests';
	import { formatPrice } from '$lib/utils';
	import { onMount } from 'svelte';

	let accounts: Account[] = [];

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

	onMount(() => {
		reloadAccounts();
	});

	function reloadAccounts() {
		deletedApi()
			.getDeletedAccounts(page, accounts_per_page, searchQuery, { withCredentials: true })
			.then((res) => {
				accounts = res.data.accounts ?? [];
				page = res.data.page;
				maxPage = res.data.max_page;
				accounts_per_page = res.data.limit;
			});
	}

	function restoreAccount(id: string) {
		deletedApi()
			.restoreDeletedAccount(id, { withCredentials: true })
			.then(() => {
				accounts = accounts.filter((ct) => ct.id !== id);
			});
	}

	function deleteAccount(id: string) {
		deletedApi()
			.deleteAccount(id, { withCredentials: true })
			.then(() => {
				accounts = accounts.filter((ct) => ct.id !== id);
			});
	}
</script>

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
							<p class="text-sm text-gray-600 dark:text-gray-400">Restorer des comptes</p>
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
								<th scope="col" class="px-6 py-3 text-right" />
							</tr>
						</thead>

						<tbody class="divide-y divide-gray-200 dark:divide-gray-700">
							{#each accounts as account}
								<tr>
									<td class="h-px w-72">
										<div class="px-6 py-3">
											<p
												class="w-72 block text-sm dark:text-white/[.8]  break-words p-2 bg-transparent"
											>
												{account.last_name}
											<p/>
										</div>
									</td>
									<td class="h-px w-72">
										<div class="px-6 py-3">
											<p
												class="w-72 block text-sm dark:text-white/[.8]  break-words p-2 bg-transparent"
											>
												{account.first_name}
											<p/>
										</div>
									</td>
									<td class="h-px w-96">
										<div class="px-6 py-3">
											<p
												class="w-72 block text-sm dark:text-white/[.8]  break-words p-2 bg-transparent"
											>
												{account.email_address}
											<p/>
										</div>
									</td>
									<td class="h-px w-72">
										<div class="px-6 py-3">
											<p
												class="text-sm dark:text-white/[.8]  break-words p-2 bg-transparent"
											>
												{formatPrice(account.balance)}
											</p>
										</div>
									</td>
									<td class="h-px w-72">
										<div class="px-6 py-3">
											<select
												class="block text-sm dark:text-white/[.8] dark:bg-slate-900 break-words p-2 bg-transparent"
												value={account.role}
												disabled
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
												class="block text-sm dark:text-white/[.8] dark:bg-slate-900 break-words p-2 bg-transparent"
												value={account.price_role}
												disabled
											>
												<option value="normal">Normal</option>
												<option value="exte">Exte</option>
												<option value="ceten">CETEN</option>
												<option value="vip">VIP</option>
												<option value="staff">Staff</option>
											</select>
										</div>
									</td>
									<td class="h-px w-px whitespace-nowrap">
										<div class="px-6 py-1.5">
											<button
												class="inline-flex items-center gap-x-1.5 text-sm text-blue-600 decoration-2 hover:underline font-medium"
												on:click={() => restoreAccount(account.id)}
											>
												Restorer
											</button>
											<button
												class="inline-flex items-center gap-x-1.5 text-sm text-blue-600 decoration-2 hover:underline font-medium"
												on:click={() => deleteAccount(account.id)}
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

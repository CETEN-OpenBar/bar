<script lang="ts">
	import { goto } from '$app/navigation';
	import { accountsApi } from '$lib/requests/requests';

	let can_restore = false;

	accountsApi()
		.getAccountAdmin({ withCredentials: true })
		.then((res) => {
			can_restore = res.data.can_restore;
			if (!res.data.is_allowed) {
				goto('/auth');
			}
		})
		.catch(() => {
			goto('/auth');
		});
</script>

<!-- svelte-ignore a11y-invalid-attribute -->

<div class="min-h-screen dark:bg-slate-950 flex flex-col">
	<header
		class="flex flex-wrap sm:justify-start sm:flex-nowrap z-50 w-full bg-blue-600 dark:bg-blue-900 text-sm py-3 sm:py-0"
	>
		<nav
			class="relative max-w-[85rem] w-full mx-auto px-4 sm:flex sm:items-center sm:justify-between sm:px-6 lg:px-8"
			aria-label="Global"
		>
			<div class="flex items-center justify-between">
				<a class="flex-none text-xl font-semibold text-white" href="#" aria-label="Panel admin"
					>Panel admin</a
				>
				<div class="sm:hidden">
					<button
						type="button"
						class="hs-collapse-toggle p-2 inline-flex justify-center items-center gap-2 rounded-md border font-medium text-white/[.5] shadow-sm align-middle hover:bg-white/[.1] hover:text-white focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-white focus:ring-blue-600 transition-all text-sm"
						data-hs-collapse="#navbar-collapse-with-animation"
						aria-controls="navbar-collapse-with-animation"
						aria-label="Toggle navigation"
					>
						<svg
							class="hs-collapse-open:hidden w-4 h-4"
							width="16"
							height="16"
							fill="currentColor"
							viewBox="0 0 16 16"
						>
							<path
								fill-rule="evenodd"
								d="M2.5 12a.5.5 0 0 1 .5-.5h10a.5.5 0 0 1 0 1H3a.5.5 0 0 1-.5-.5zm0-4a.5.5 0 0 1 .5-.5h10a.5.5 0 0 1 0 1H3a.5.5 0 0 1-.5-.5zm0-4a.5.5 0 0 1 .5-.5h10a.5.5 0 0 1 0 1H3a.5.5 0 0 1-.5-.5z"
							/>
						</svg>
						<svg
							class="hs-collapse-open:block hidden w-4 h-4"
							width="16"
							height="16"
							fill="currentColor"
							viewBox="0 0 16 16"
						>
							<path
								d="M4.646 4.646a.5.5 0 0 1 .708 0L8 7.293l2.646-2.647a.5.5 0 0 1 .708.708L8.707 8l2.647 2.646a.5.5 0 0 1-.708.708L8 8.707l-2.646 2.647a.5.5 0 0 1-.708-.708L7.293 8 4.646 5.354a.5.5 0 0 1 0-.708z"
							/>
						</svg>
					</button>
				</div>
			</div>
			<div
				id="navbar-collapse-with-animation"
				class="hs-collapse hidden overflow-hidden transition-all duration-300 basis-full grow sm:block"
			>
				<div
					class="flex flex-col gap-y-4 gap-x-0 mt-5 sm:flex-row sm:items-center sm:justify-end sm:gap-y-0 sm:gap-x-7 sm:mt-0 sm:pl-7"
				>
					{#if can_restore}
						<div
							class="hs-dropdown [--strategy:static] sm:[--strategy:fixed] [--adaptive:none] sm:[--trigger:hover] sm:py-4"
						>
							<button
								type="button"
								class="flex items-center w-full text-white/[.8] hover:text-white font-medium"
							>
								Restauration
								<svg
									class="ml-2 w-2.5 h-2.5"
									width="16"
									height="16"
									viewBox="0 0 16 16"
									fill="none"
									xmlns="http://www.w3.org/2000/svg"
								>
									<path
										d="M2 5L8.16086 10.6869C8.35239 10.8637 8.64761 10.8637 8.83914 10.6869L15 5"
										stroke="currentColor"
										stroke-width="2"
										stroke-linecap="round"
									/>
								</svg>
							</button>

							<div
								class="hs-dropdown-menu transition-[opacity,margin] duration-[0.1ms] sm:duration-[150ms] hs-dropdown-open:opacity-100 opacity-0 sm:w-48 hidden z-10 bg-white sm:shadow-md rounded-lg p-2 dark:bg-gray-800 dark:divide-gray-700 before:absolute top-full sm:border before:-top-5 before:left-0 before:w-full before:h-5"
							>
								<a
									class="flex items-center gap-x-3.5 py-2 px-3 rounded-md text-sm text-gray-800 hover:bg-gray-100 focus:ring-2 focus:ring-blue-500 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-gray-300"
									href="/admin/restore/accounts"
								>
									Comptes
								</a>
								<a
									class="flex items-center gap-x-3.5 py-2 px-3 rounded-md text-sm text-gray-800 hover:bg-gray-100 focus:ring-2 focus:ring-blue-500 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-gray-300"
									href="/admin/restore/categories"
								>
									Catégories
								</a>
								<a
									class="flex items-center gap-x-3.5 py-2 px-3 rounded-md text-sm text-gray-800 hover:bg-gray-100 focus:ring-2 focus:ring-blue-500 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-gray-300"
									href="/admin/restore/produits"
								>
									Produits
								</a>
								<a
									class="flex items-center gap-x-3.5 py-2 px-3 rounded-md text-sm text-gray-800 hover:bg-gray-100 focus:ring-2 focus:ring-blue-500 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-gray-300"
									href="/admin/restore/carousel/textes"
								>
									Textes du carousel
								</a>
								<a
									class="flex items-center gap-x-3.5 py-2 px-3 rounded-md text-sm text-gray-800 hover:bg-gray-100 focus:ring-2 focus:ring-blue-500 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-gray-300"
									href="/admin/restore/carousel/images"
								>
									Images du carousel
								</a>
							</div>
						</div>
					{/if}
					<a
						class="font-medium text-white/[.8] hover:text-white sm:py-6"
						href="/admin/accounts"
						aria-current="page">Comptes</a
					>
					<a
						class="font-medium text-white/[.8] hover:text-white sm:py-6"
						href="/admin/categories"
						aria-current="page">Catégories</a
					>
					<a
						class="font-medium text-white/[.8] hover:text-white sm:py-6"
						href="/admin/produits"
						aria-current="page">Produits</a
					>
					<a
						class="font-medium text-white/[.8] hover:text-white sm:py-6"
						href="/panel"
						aria-current="page">Panel</a
					>

					<div
						class="hs-dropdown [--strategy:static] sm:[--strategy:fixed] [--adaptive:none] sm:[--trigger:hover] sm:py-4"
					>
						<button
							type="button"
							class="flex items-center w-full text-white/[.8] hover:text-white font-medium"
						>
							Carousel
							<svg
								class="ml-2 w-2.5 h-2.5"
								width="16"
								height="16"
								viewBox="0 0 16 16"
								fill="none"
								xmlns="http://www.w3.org/2000/svg"
							>
								<path
									d="M2 5L8.16086 10.6869C8.35239 10.8637 8.64761 10.8637 8.83914 10.6869L15 5"
									stroke="currentColor"
									stroke-width="2"
									stroke-linecap="round"
								/>
							</svg>
						</button>

						<div
							class="hs-dropdown-menu transition-[opacity,margin] duration-[0.1ms] sm:duration-[150ms] hs-dropdown-open:opacity-100 opacity-0 sm:w-48 hidden z-10 bg-white sm:shadow-md rounded-lg p-2 dark:bg-gray-800 dark:divide-gray-700 before:absolute top-full sm:border before:-top-5 before:left-0 before:w-full before:h-5"
						>
							<a
								class="flex items-center gap-x-3.5 py-2 px-3 rounded-md text-sm text-gray-800 hover:bg-gray-100 focus:ring-2 focus:ring-blue-500 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-gray-300"
								href="/admin/carousel/textes"
							>
								Textes
							</a>
							<a
								class="flex items-center gap-x-3.5 py-2 px-3 rounded-md text-sm text-gray-800 hover:bg-gray-100 focus:ring-2 focus:ring-blue-500 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-gray-300"
								href="/admin/carousel/images"
							>
								Images
							</a>
						</div>
					</div>
				</div>
			</div>
		</nav>
	</header>
	<slot />
</div>

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

<div class="min-h-screen dark:bg-slate-950">
	<header
		class="flex flex-wrap sm:justify-start sm:flex-nowrap z-50 w-full bg-blue-600 dark:bg-blue-900 text-sm py-3 sm:py-0"
	>
		<nav
			class="relative max-w-[85rem] w-full mx-auto px-4 sm:flex sm:items-center sm:justify-between sm:px-6 lg:px-8 p-5"
			aria-label="Global"
		>
			<div class="flex items-center justify-between">
				<a class="flex-none text-xl font-semibold text-white" href="#" aria-label="Panel admin"
					>(Nouveau) Panel admin</a
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
		</nav>
	</header>
	<slot />
</div>

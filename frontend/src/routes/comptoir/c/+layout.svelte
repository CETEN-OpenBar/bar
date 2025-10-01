<script lang="ts">
	import { goto } from '$app/navigation';
	import { accountsApi } from '$lib/requests/requests';
	import type { Account } from '$lib/api';
	import { onMount } from 'svelte';
	import { store } from '$lib/store/store';
	import 'iconify-icon';
	import Autodisconnect from '$lib/components/random/autodisconnect.svelte';

	let account: Account | undefined = undefined;

	onMount(() => {
		accountsApi()
			.getAccount({
				withCredentials: true
			})
			.then((res) => {
				account = res.data.account;
				store.set({ account });
			})
			.catch(() => {
				goto('/comptoir');
			});
	});
</script>

<style>
	:root {
		--bg-primary: #000000;
		--bg-secondary: #1a1a1a;
		--bg-tertiary: #2d2d2d;
		--text-primary: #ffffff;
		--text-secondary: #b0b0b0;
		--accent-primary: #2563eb;
		--accent-hover: #1d4ed8;
		--border-color: #404040;
		--shadow-color: rgba(0, 0, 0, 0.3);
	}

	@media (prefers-color-scheme: light) {
		:root {
			--bg-primary: #ffffff;
			--bg-secondary: #f8fafc;
			--bg-tertiary: #e2e8f0;
			--text-primary: #1a1a1a;
			--text-secondary: #64748b;
			--accent-primary: #2563eb;
			--accent-hover: #1d4ed8;
			--border-color: #e2e8f0;
			--shadow-color: rgba(0, 0, 0, 0.1);
		}
	}

	.main-container {
		background-color: var(--bg-primary);
		color: var(--text-primary);
		min-height: 100vh;
		width: 100vw;
		overflow: hidden;
	}

	.content-wrapper {
		display: flex;
		flex-direction: column;
		height: 100vh;
		padding: 1rem;
		gap: 1rem;
	}
</style>

<Autodisconnect delay={900000} location="/comptoir" />

{#if account !== undefined}
	<main class="main-container">
		<div class="content-wrapper">
			<slot />
		</div>
	</main>
{/if}

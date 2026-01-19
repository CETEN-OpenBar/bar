<script lang="ts">
	import { goto } from '$app/navigation';
	import { accountsApi } from '$lib/requests/requests';
	import type { Account } from '$lib/api';
	import { onMount } from 'svelte';
	import { store } from '$lib/store/store';
	import 'iconify-icon';
	import Autodisconnect from '$lib/components/random/autodisconnect.svelte';
	import ReadCard from '$lib/components/readCard.svelte';
	import { open_caisse, open_door, open_ventilo } from '$lib/local/local';
	import NewRefill from '$lib/components/comptoir/newRefill.svelte';
	import ChangePassword from '$lib/components/comptoir/changePassword.svelte';
	import Password from '$lib/components/password.svelte';
	import type { ConnectPasswordRequest } from '$lib/api';
	import { authApi } from '$lib/requests/requests';

	let account: Account | undefined = undefined;

	function reset() {
		askForCard = false;
		askForPassword = false;
		to_call = () => {};
	}

	let to_call = open_door;
	let infos = {
		card_id: '',
		card_pin: ''
	};
	let askForCard = false;
	let askForPassword = false;
	let newRefill = false;
	let changePassword = false;
	let info: ConnectPasswordRequest;
	let mobileMenuOpen = false;

	function close() {
		newRefill = false;
	}

	function logoutAccount() {
		authApi()
			.logout({
				withCredentials: true
			})
			.then(() => {
				goto('/comptoir');
			})
			.catch(() => {
				goto('/comptoir');
			});
	}

	const go_admin_panel = (card_id: string, password: string) => {
		info = {
			card_id: card_id,
			password: password,
		}
		authApi()
			.connectPassword(info, { withCredentials: true })
			.then(() => {
				goto('/panel');
			});
	};

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
		--overlay-bg: rgba(0, 0, 0, 0.5);
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
			--overlay-bg: rgba(0, 0, 0, 0.3);
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
	}

	.slot-container {
		flex: 1;
		min-height: 0;
		display: flex;
		flex-direction: column;
	}

	.main-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1rem;
		background-color: var(--bg-secondary);
		border-bottom: 1px solid var(--border-color);
	}

	.actions-group {
		display: flex;
		gap: 0.5rem;
		flex-wrap: wrap;
		align-items: center;
	}

	.nav-group {
		display: flex;
		gap: 0.5rem;
		align-items: center;
	}

	.right-section {
		display: flex;
		gap: 0.5rem;
		align-items: center;
	}

	.mobile-menu-btn {
		display: none;
		background: none;
		border: none;
		color: var(--text-primary);
		cursor: pointer;
		padding: 0.5rem;
		border-radius: 0.5rem;
		transition: all 0.2s ease;
	}

	.mobile-menu-btn:hover {
		background-color: var(--bg-tertiary);
	}

	.mobile-menu {
		display: none;
		position: absolute;
		top: 100%;
		left: 0;
		right: 0;
		background-color: var(--bg-tertiary);
		border-radius: 0.75rem;
		border: 1px solid var(--border-color);
		padding: 1rem;
		margin-top: 0.5rem;
		z-index: 100;
		box-shadow: 0 10px 15px -3px var(--shadow-color);
	}

	.mobile-menu.open {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.btn-small {
		padding: 0.5rem 0.75rem;
		font-size: 0.875rem;
	}

	@media (max-width: 768px) {
		.mobile-menu-btn {
			display: block;
		}

		.nav-group {
			display: none;
		}

		.actions-group {
			display: none;
		}

		.main-header {
			position: relative;
		}

		.right-section {
			margin-left: auto;
		}

		.mobile-menu {
			display: none;
			position: absolute;
			top: 100%;
			right: 0;
			background-color: var(--bg-tertiary);
			border-radius: 0.75rem;
			border: 1px solid var(--border-color);
			padding: 1rem;
			margin-top: 0.5rem;
			z-index: 100;
			box-shadow: 0 10px 15px -3px var(--shadow-color);
			min-width: 200px;
		}

		.mobile-menu.open {
			display: flex;
			flex-direction: column;
			gap: 0.5rem;
		}

		.mobile-menu .btn {
			width: 100%;
			justify-content: flex-start;
		}
	}

	.btn {
		padding: 0.75rem 1rem;
		border-radius: 0.5rem;
		background-color: var(--accent-primary);
		color: white;
		border: none;
		cursor: pointer;
		font-size: 1rem;
		font-weight: 500;
		transition: all 0.2s ease;
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.btn:hover {
		background-color: var(--accent-hover);
		transform: translateY(-1px);
	}

	.btn-danger {
		background-color: #dc2626;
	}

	.btn-danger:hover {
		background-color: #b91c1c;
	}

	.modal-overlay {
		position: fixed;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		background-color: var(--overlay-bg);
		display: flex;
		justify-content: center;
		align-items: center;
		z-index: 1000;
		cursor: pointer;
	}

	.modal-content {
		background-color: var(--bg-tertiary);
		padding: 2rem;
		border-radius: 1rem;
		box-shadow: 0 20px 25px -5px var(--shadow-color);
		position: relative;
		min-width: 400px;
		min-height: 350px;
		text-align: center;
		cursor: default;
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		gap: 1rem;
	}

	.modal-close {
		position: absolute;
		top: 1rem;
		right: 1rem;
		background: none;
		border: none;
		color: var(--text-secondary);
		cursor: pointer;
		padding: 0.5rem;
		border-radius: 50%;
		transition: all 0.2s ease;
	}


	.modal-title {
		font-size: 1.5rem;
		font-weight: 600;
		color: var(--text-primary);
		margin-bottom: 1rem;
	}
</style>

<Autodisconnect delay={900000} location="/comptoir" />

{#if account !== undefined}
	<main class="main-container">
		<div class="content-wrapper">
			<header class="main-header">
				<div class="actions-group">
					<button class="btn btn-small" on:click={() => newRefill = true}>
						<iconify-icon icon="mdi:cash-plus" width="16" height="16" />
						Recharge
					</button>

					<button class="btn" on:click={() => changePassword = true}>
						<iconify-icon icon="mdi:key-change" width="20" height="20" />
						Changer MDP
					</button>
					<button class="btn" on:click={() => { to_call = open_door; askForCard = true; }}>
						<iconify-icon icon="mdi:door-open" width="20" height="20" />
						Porte
					</button>
					<button class="btn" on:click={() => { to_call = open_ventilo; askForCard = true; }}>
						<iconify-icon icon="mdi:fan" width="20" height="20" />
						Ventilo
					</button>
					<button class="btn" on:click={() => { to_call = open_caisse; askForCard = true; }}>
						<iconify-icon icon="mdi:cash-register" width="20" height="20" />
						Caisse
					</button>
					<button class="btn" on:click={() => { to_call = go_admin_panel; askForCard = true; }}>
						<iconify-icon icon="mdi:shield-crown" width="20" height="20" />
						Admin Panel
					</button>
				</div>

				<div class="right-section">
					<div class="nav-group">
						<button class="btn btn-danger" on:click={logoutAccount}>
							<iconify-icon icon="mdi:logout" width="24" height="24" />
							Déconnexion
						</button>
					</div>

					<button class="mobile-menu-btn" on:click={() => mobileMenuOpen = !mobileMenuOpen}>
						<iconify-icon icon="mdi:menu" width="24" height="24" />
					</button>
				</div>

				<div class="mobile-menu" class:open={mobileMenuOpen}>
					<button class="btn btn-small" on:click={() => { newRefill = true; mobileMenuOpen = false; }}>
						<iconify-icon icon="mdi:cash-plus" width="16" height="16" />
						Recharge
					</button>

					<button class="btn" on:click={() => { changePassword = true; mobileMenuOpen = false; }}>
						<iconify-icon icon="mdi:key-change" width="20" height="20" />
						Changer MDP
					</button>
					<button class="btn" on:click={() => { to_call = open_door; askForCard = true; mobileMenuOpen = false; }}>
						<iconify-icon icon="mdi:door-open" width="20" height="20" />
						Porte
					</button>
					<button class="btn" on:click={() => { to_call = open_ventilo; askForCard = true; mobileMenuOpen = false; }}>
						<iconify-icon icon="mdi:fan" width="20" height="20" />
						Ventilo
					</button>
					<button class="btn" on:click={() => { to_call = open_caisse; askForCard = true; mobileMenuOpen = false; }}>
						<iconify-icon icon="mdi:cash-register" width="20" height="20" />
						Caisse
					</button>
					<button class="btn" on:click={() => { to_call = go_admin_panel; askForCard = true; mobileMenuOpen = false; }}>
						<iconify-icon icon="mdi:shield-crown" width="20" height="20" />
						Admin Panel
					</button>

					<button class="btn btn-danger" on:click={logoutAccount}>
						<iconify-icon icon="mdi:logout" width="24" height="24" />
						Déconnexion
					</button>
				</div>
			</header>

			<div class="slot-container">
				<slot />
			</div>

			{#if askForCard}
				<div class="modal-overlay" on:click={reset}>
					<div class="modal-content" on:click|stopPropagation>
						<button class="modal-close" on:click={reset}>
							<iconify-icon icon="mdi:close" width="24" height="24" />
						</button>
						<h2 class="modal-title">Veuillez scanner la carte</h2>
					</div>
				</div>
				<ReadCard
					callback={(id) => {
						infos.card_id = id;
						askForPassword = true;
						askForCard = false;
					}}
				/>
			{/if}

			{#if askForPassword}
				<Password
					callback={(password) => {
						infos.card_pin = password;
						to_call(infos.card_id, infos.card_pin);
						reset();
					}}
				/>
			{/if}

			{#if changePassword}
				<ChangePassword onEnd={() => (changePassword = false)} />
			{/if}

			{#if newRefill}
				<NewRefill {close} />
			{/if}
		</div>
	</main>
{/if}
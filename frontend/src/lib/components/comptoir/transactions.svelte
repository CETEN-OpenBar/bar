<script lang="ts">
	import { goto } from '$app/navigation';
	import {
		TransactionItemState,
		type Transaction,
		type TransactionItem,
		type TransactionState
	} from '$lib/api';
	import { api } from '$lib/config/config';
	import { transactionsApi } from '$lib/requests/requests';
	import { formatPrice } from '$lib/utils';
	import { onDestroy, onMount } from 'svelte';
	import TransactionPopup from './transactionPopup.svelte';
	import { dragscroll } from '@svelte-put/dragscroll';
	import { searchName } from '$lib/store/store';
	import ComptoirHeaderControls from './headerControls.svelte';
	import PaginationFooter from '$lib/components/PaginationFooter.svelte';

	function handleAvatarError(e: Event) {
		const target = e.currentTarget as HTMLImageElement;
		target.style.display = 'none';
		target.nextElementSibling?.classList.remove('hidden');
	}

	let searchNameValue: string;

	searchName.subscribe((value) => {
		searchNameValue = value;
	});

	let transactions: Array<Transaction> = [];
	let maxItemPerTransaction: number = 6;
	let interval: number;
	let transactionAmount: number = 4;

	let page: number = 1;
	let maxPage: number = 1;
	let nextPage = () => {
		if (page < maxPage) {
			page++;
			reloadTransactions();
		}
	};
	let prevPage = () => {
		if (page > 1) {
			page--;
			reloadTransactions();
		}
	};

	onMount(() => {
		reloadTransactions();
		interval = setInterval(() => {
			reloadTransactions();
		}, 2000);

		window.addEventListener('keydown', handleKeyDown);

		const transactionsList = document.querySelector('.transactions-list');
		if (transactionsList) {
			let touchStartX = 0;
			let touchEndX = 0;
			let isTouch = false;

			transactionsList.addEventListener(
				'touchstart',
				(e: Event) => {
					const touchEvent = e as TouchEvent;
					touchStartX = touchEvent.changedTouches[0].screenX;
					isTouch = true;
				},
				{ passive: true }
			);

			transactionsList.addEventListener(
				'touchend',
				(e: Event) => {
					if (isTouch) {
						const touchEvent = e as TouchEvent;
						touchEndX = touchEvent.changedTouches[0].screenX;
						handleSwipe();
						isTouch = false;
					}
				},
				{ passive: true }
			);

			transactionsList.addEventListener('mousedown', (e: Event) => {
				const mouseEvent = e as MouseEvent;
				if (!isTouch) {
					touchStartX = mouseEvent.screenX;
					isTouch = false;
				}
			});

			transactionsList.addEventListener('mouseup', (e: Event) => {
				const mouseEvent = e as MouseEvent;
				if (!isTouch) {
					touchEndX = mouseEvent.screenX;
					handleMouseSwipe();
				}
			});

			const handleSwipe = () => {
				const swipeThreshold = 50;
				const diff = touchStartX - touchEndX;

				if (Math.abs(diff) > swipeThreshold) {
					if (diff > 0) {
						if (page < maxPage) nextPage();
					} else {
						if (page > 1) prevPage();
					}
				}
			};

			const handleMouseSwipe = () => {
				const swipeThreshold = 100;
				const diff = touchStartX - touchEndX;

				if (Math.abs(diff) > swipeThreshold) {
					if (diff > 0) {
						if (page < maxPage) nextPage();
					} else {
						if (page > 1) prevPage();
					}
				}
			};
		}
	});

	onDestroy(() => {
		clearInterval(interval);
		if (typeof window !== 'undefined') {
			window.removeEventListener('keydown', handleKeyDown);
		}
	});

	let st: TransactionState | undefined = 'started';

	let showRemoteTransactions: boolean = true;

	type TransactionItemWithFakeAmount = TransactionItem & { item_fake_amount?: number };

	function reloadTransactions() {
		transactionsApi()
			.getTransactions(
				page,
				transactionAmount,
				st,
				false,
				!showRemoteTransactions,
				searchNameValue,
				undefined,
				undefined,
				undefined,
				{ withCredentials: true }
			)
			.then((res) => {
				page = res.data.page ?? 1;
				maxPage = res.data.max_page ?? 1;
				if (!(res.data.transactions instanceof Array)) return;
				let countedItems = 0;
				let newTransactions = [];
				for (let transaction of res.data.transactions) {
					let items: Array<TransactionItemWithFakeAmount> = [];

					for (let r of transaction.items ?? []) {
						let item = r as TransactionItemWithFakeAmount;
						if (countedItems >= maxItemPerTransaction) break;
						if (countedItems + item.item_amount > maxItemPerTransaction) {
							item.item_fake_amount = maxItemPerTransaction - countedItems;
						}
						items.push(item);
					}
					transaction.items = items;
					newTransactions.push(transaction);
				}
				transactions = newTransactions;
			});
	}

	let displayTransaction: Transaction | null = null;

	function handleSearchInput(event: Event) {
		const target = event.currentTarget as HTMLInputElement;
		searchNameValue = target.value.toLowerCase();
		page = 1;
		reloadTransactions();
	}

	function handleKeyDown(event: KeyboardEvent) {
		// Only handle navigation keys if not typing in an input
		if (
			event.target instanceof HTMLInputElement ||
			event.target instanceof HTMLTextAreaElement ||
			event.target instanceof HTMLSelectElement
		) {
			return;
		}

		switch (event.key) {
			case 'ArrowLeft':
				event.preventDefault();
				if (page > 1) prevPage();
				break;
			case 'ArrowRight':
				event.preventDefault();
				if (page < maxPage) nextPage();
				break;
			case 'PageUp':
				event.preventDefault();
				if (page > 1) prevPage();
				break;
			case 'PageDown':
				event.preventDefault();
				if (page < maxPage) nextPage();
				break;
			case 'Home':
				event.preventDefault();
				if (page > 1) {
					page = 1;
					reloadTransactions();
				}
				break;
			case 'End':
				event.preventDefault();
				if (page < maxPage) {
					page = maxPage;
					reloadTransactions();
				}
				break;
		}
	}
</script>

{#if displayTransaction}
	<TransactionPopup
		transaction={displayTransaction}
		close={() => {
			reloadTransactions();
			displayTransaction = null;
		}}
	/>
{/if}

<div class="transactions-wrapper">
	<div class="transactions-content">
		<ComptoirHeaderControls activeTab="transactions">
			<div class="filters-section" slot="filters">
				<div>
					Filtre :
					<select class="filter-select" bind:value={st} on:change={reloadTransactions}>
						<option value={undefined}>Tout</option>
						<option value="started">En cours</option>
						<option value="finished">Terminées</option>
						<option value="canceled">Annulées</option>
					</select>
				</div>
				<div>
					<label class="checkbox-label">
						Commandes en ligne
						<input
							type="checkbox"
							class="checkbox-input"
							bind:checked={showRemoteTransactions}
							on:change={reloadTransactions}
						/>
					</label>
				</div>
			</div>
		</ComptoirHeaderControls>

		<div use:dragscroll class="transactions-list">
			{#each transactions as transaction}
				<button
					class="transaction-card {transaction.state}"
					on:click={() => (displayTransaction = transaction)}
				>
					<div class="transaction-header">
						{#if transaction.is_remote}
							<iconify-icon icon="mdi:wifi" class="remote-icon" />
						{:else}
							<iconify-icon icon="mdi:monitor" class="remote-icon" />
						{/if}
						{#if transaction.account_google_picture}
							<img
								src={transaction.account_google_picture}
								alt="Avatar"
								class="transaction-avatar"
								on:error={handleAvatarError}
							/>
						{/if}
						<iconify-icon
							icon="mdi:account-circle"
							class="transaction-avatar placeholder-icon {transaction.account_google_picture
								? 'hidden'
								: ''}"
						/>
						<b>{transaction.account_nick_name || transaction.account_name}</b>
						{#if transaction.account_nick_name && transaction.account_name && transaction.account_nick_name !== transaction.account_name}
							({transaction.account_name})
						{/if}
					</div>

					<div class="transaction-content">
						<div class="items-section">
							<div class="items-grid">
								{#each transaction.items as item}
									<div
										class="item-card {item.state == TransactionItemState.TransactionItemCanceled
											? 'canceled'
											: ''} {item.state == TransactionItemState.TransactionItemFinished
											? 'finished'
											: ''}"
									>
										{#if item.state == TransactionItemState.TransactionItemCanceled || item.state == TransactionItemState.TransactionItemFinished}
											<img src={api() + item.picture_uri} alt="item" class="item-image completed" />
											<div class="item-name completed">{item.item_name}</div>
											<div class="item-amount completed">x {item.item_amount}</div>
										{:else}
											<img src={api() + item.picture_uri} alt="item" class="item-image" />
											<div class="item-name">{item.item_name}</div>
											<div class="item-amount">x {item.item_amount}</div>
										{/if}
									</div>
								{/each}
							</div>
						</div>
						<div class="price-section">
							{formatPrice(transaction.total_cost)}
						</div>
					</div>

					<div
						class="status-led {transaction.state == 'finished'
							? 'validated'
							: transaction.state == 'canceled'
							? 'cancelled'
							: 'waiting'}"
					/>
				</button>
			{/each}
		</div>
	</div>

	<PaginationFooter
		bind:page
		{maxPage}
		resultsCount={transactions.length}
		showPageInput={true}
		on:prevPage={prevPage}
		on:nextPage={nextPage}
		on:pageChange={reloadTransactions}
	/>
</div>

<style>
	.transactions-wrapper {
		width: 100%;
		height: 100%;
		display: flex;
		flex-direction: column;
	}

	.transactions-content {
		display: flex;
		flex-direction: column;
		flex: 1;
		min-height: 0;
		padding: 15px;
		padding-bottom: 0;
	}

	.header-controls {
		display: flex;
		flex-direction: row;
		align-items: center;
		gap: 20px;
	}

	.title-section {
		display: flex;
		flex-direction: row;
		align-items: center;
		gap: 40px;
		flex-grow: 1;
	}

	.title {
		font-size: 18px;
		font-weight: 600;
	}

	.view-tabs {
		display: flex;
		background-color: #e5e7eb;
		border-radius: 8px;
		padding: 4px;
		gap: 0;
	}

	.tab {
		display: flex;
		align-items: center;
		gap: 6px;
		padding: 8px 16px;
		background-color: transparent;
		color: #374151;
		border: none;
		cursor: pointer;
		font-size: 14px;
		font-weight: 500;
		transition: all 0.2s ease;
		border-radius: 6px;
	}

	.tab:hover {
		background-color: rgba(0, 0, 0, 0.05);
	}

	.tab.active {
		background-color: #2563eb;
		color: white;
	}

	@media (prefers-color-scheme: dark) {
		.view-tabs {
			background-color: #374151;
		}

		.tab {
			color: #e5e7eb;
		}

		.tab:hover {
			background-color: rgba(255, 255, 255, 0.05);
		}

		.tab.active {
			background-color: #2563eb;
			color: white;
		}
	}

	.search-input {
		border-radius: 8px;
		padding: 8px;
		background-color: white;
		color: black;
		border: 1px solid #ccc;
	}

	@media (prefers-color-scheme: dark) {
		.search-input {
			background-color: #374151;
			color: white;
			border: 1px solid #6b7280;
		}
	}

	.filters-section {
		display: flex;
		flex-direction: row;
		align-items: center;
		gap: 20px;
		font-size: 18px;
		font-weight: 600;
	}

	.filter-select {
		background-color: #f3f4f6;
		padding: 8px;
		border-radius: 8px;
		color: #111827;
		border: 1px solid #d1d5db;
		cursor: pointer;
	}

	.filter-select:hover {
		background-color: #e5e7eb;
	}

	@media (prefers-color-scheme: dark) {
		.filter-select {
			background-color: #374151;
			color: white;
			border: none;
		}

		.filter-select:hover {
			background-color: #4b5563;
		}
	}

	.checkbox-label {
		display: flex;
		align-items: center;
		gap: 4px;
	}

	.checkbox-input {
		height: 24px;
		width: 24px;
		vertical-align: middle;
		margin-left: 4px;
	}

	.transactions-list {
		display: flex;
		flex-direction: column;
		overflow-y: auto;
		flex: 1;
		min-height: 0;
		gap: 20px;
	}

	.transaction-card {
		display: flex;
		flex-direction: column;
		border: 2px solid #e5e7eb;
		border-radius: 12px;
		cursor: pointer;
		transition: all 0.2s ease;
		position: relative;
		background-color: #f9fafb;
		flex-shrink: 0;
	}

	@media (prefers-color-scheme: dark) {
		.transaction-card {
			background-color: var(--bg-secondary);
			border: 2px solid #3c4554;
		}
	}

	.transaction-card:hover {
		border: 2px solid #9ca3af;
	}

	.transaction-card.canceled {
		border-color: #ef4444;
	}

	.transaction-card.finished {
		border-color: #22c55e;
	}

	.transaction-header {
		color: black;
		padding: 12px;
		padding-right: 35px;
		display: flex;
		align-items: center;
		gap: 8px;
		position: relative;
		border-bottom: 1px solid #9ca3af;
	}

	@media (prefers-color-scheme: dark) {
		.transaction-header {
			color: white;
			border-bottom-color: #4b5563;
		}
	}

	.transaction-avatar {
		width: 40px;
		height: 40px;
		border-radius: 50%;
		object-fit: cover;
		flex-shrink: 0;
	}

	.placeholder-icon {
		color: #9ca3af;
		font-size: 40px;
		flex-shrink: 0;
	}

	@media (prefers-color-scheme: dark) {
		.placeholder-icon {
			color: #6b7280;
		}
	}

	.remote-icon {
		height: 20px;
		width: 20px;
		margin: 4px;
		position: absolute;
		right: 12px;
		top: 17px;
		color: #374151;
	}

	@media (prefers-color-scheme: dark) {
		.remote-icon {
			color: #e5e7eb;
		}
	}

	.transaction-content {
		display: grid;
		grid-template-columns: 10fr 1fr;
		gap: 8px;
	}

	.items-section {
		padding: 4px;
		padding-right: 4px;
		height: 100%;
		width: 100%;
		border-right: 1px solid #9ca3af;
	}

	@media (prefers-color-scheme: dark) {
		.items-section {
			border-right-color: #4b5563;
		}
	}

	.items-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(80px, 1fr));
		gap: 8px;
	}

	.item-card {
		display: flex;
		flex-direction: column;
		height: auto;
		justify-content: center;
		padding: 8px;
		border-radius: 4px;
	}

	.item-card.canceled {
		background-color: #fb9393;
		opacity: 0.85;
	}

	.item-card.finished {
		background-color: #86efac;
		opacity: 0.85;
	}

	@media (prefers-color-scheme: dark) {
		.item-card.canceled {
			background-color: #843535;
		}

		.item-card.finished {
			background-color: #145e31;
		}
	}

	.item-image {
		width: 36px;
		height: 36px;
		border-radius: 16px;
		align-self: center;
	}

	.item-image.completed {
		opacity: 0.75;
	}

	.item-name {
		font-size: 14px;
		text-align: center;
		color: black;
	}

	.item-name.completed {
		opacity: 0.75;
	}

	.item-amount {
		font-size: 12px;
		text-align: center;
		color: black;
	}

	.item-amount.completed {
		opacity: 0.75;
	}

	@media (prefers-color-scheme: dark) {
		.item-name,
		.item-amount {
			color: white;
		}
	}

	.price-section {
		padding: 4px;
		padding-left: 4px;
		width: 100%;
		font-size: 18px;
		align-self: center;
		text-align: center;
		color: black;
	}

	@media (prefers-color-scheme: dark) {
		.price-section {
			color: white;
		}
	}

	.status-led {
		width: 12px;
		height: 12px;
		border-radius: 50%;
		position: absolute;
		left: 12px;
		top: 70%;
		transform: translateY(-50%);
		box-shadow: 0 0 4px rgba(0, 0, 0, 0.3);
	}

	.status-led.validated {
		background-color: #22c55e;
		box-shadow: 0 0 8px #22c55e;
	}

	.status-led.cancelled {
		background-color: #ef4444;
		box-shadow: 0 0 8px #ef4444;
	}

	.status-led.waiting {
		background-color: #f97316;
		box-shadow: 0 0 8px #f97316;
		animation: pulse-orange 2s infinite;
	}

	@keyframes pulse-orange {
		0% {
			opacity: 1;
		}
		50% {
			opacity: 0.6;
		}
		100% {
			opacity: 1;
		}
	}

	@media (max-width: 1024px) {
		.items-grid {
			grid-template-columns: repeat(auto-fit, minmax(70px, 1fr));
		}
	}

	@media (max-width: 768px) {
		.header-controls {
			flex-direction: column;
			gap: 12px;
		}

		.title-section {
			flex-direction: column;
			gap: 12px;
			align-items: stretch;
		}

		.filters-section {
			flex-direction: column;
			gap: 12px;
			align-items: stretch;
		}

		.items-grid {
			grid-template-columns: repeat(auto-fit, minmax(60px, 1fr));
			gap: 6px;
		}

		.transaction-content {
			grid-template-columns: 8fr 1fr;
		}

		.item-card {
			padding: 6px;
		}

		.item-image {
			width: 30px;
			height: 30px;
		}

		.item-name,
		.item-amount {
			font-size: 10px;
		}
	}

	@media (max-width: 480px) {
        .view-tabs button {
            font-size: 9px;
            padding: 3px 8px;
        }

		.items-grid {
			grid-template-columns: repeat(auto-fit, minmax(50px, 1fr));
		}

		.transaction-content {
			grid-template-columns: 1fr;
		}

		.items-section {
			border-right: none;
			border-bottom: 1px solid #9ca3af;
		}

		@media (prefers-color-scheme: dark) {
			.items-section {
				border-bottom-color: #4b5563;
			}
		}

		.price-section {
			text-align: center;
			padding-top: 8px;
		}

		.transaction-header {
			font-size: 14px;
		}

		.remote-icon {
			width: 16px;
			height: 16px;
		}

		.transaction-avatar {
			width: 32px;
			height: 32px;
			font-size: 32px;
		}

		.placeholder-icon {
			font-size: 32px;
		}

		.status-led {
			top: 53%;
		}


	}
</style>

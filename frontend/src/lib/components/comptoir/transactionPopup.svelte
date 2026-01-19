<script lang="ts">
	import {
		ItemState,
		TransactionItemState,
		type Transaction,
		type MenuItem,
		type MenuCategory,
		type TransactionItem
	} from '$lib/api';
	import { api } from '$lib/config/config';
	import { transactionsApi } from '$lib/requests/requests';
	import { formatPrice } from '$lib/utils';
	import Error from '../error.svelte';
	import Success from '../success.svelte';
	import Transactions from './transactions.svelte';
	import { createEventDispatcher } from 'svelte';
	import { searchName } from '$lib/store/store';

	export let transaction: Transaction;
	export let close: () => void;

	let newTransaction: Transaction = structuredClone(transaction);
	let success = '';
	let error = '';

	type MenuPopup = {
		items: MenuItem[] | undefined;
		categories: MenuCategory[] | undefined;
		pickedItems: TransactionItem[] | undefined;
	};
	let menuPopup: MenuPopup | undefined;

	const eventDispatcher = createEventDispatcher();

	async function cancelTransaction() {
		let res = await transactionsApi().patchTransactionId(
			newTransaction.account_id,
			newTransaction.id,
			'canceled',
			{
				withCredentials: true
			}
		);

		if (res.status != 200) {
			error = "Une erreur s'est produite";
			setTimeout(() => {
				error = '';
			}, 1500);
			return;
		}

		transaction = newTransaction;
		success = 'Commande annulée';
		setTimeout(() => {
			success = '';
			close();
		}, 1500);
	}

	async function putBackTransaction() {
		let res = await transactionsApi().patchTransactionId(
			newTransaction.account_id,
			newTransaction.id,
			'started',
			{
				withCredentials: true
			}
		);

		if (res.status != 200) {
			error = "Une erreur s'est produite";
			setTimeout(() => {
				error = '';
			}, 1500);
			return;
		}

		transaction = newTransaction;
		success = 'Commande remise en attente';
		setTimeout(() => {
			success = '';
			close();
		}, 1500);
	}

	async function finishTransaction() {
		for (let i = 0; i < newTransaction.items.length; i++) {
			let item = newTransaction.items[i];

			// @ts-ignore
			if (item.state == transaction.items[i].state) item.state = undefined;

			let res = await transactionsApi().patchTransactionItemId(
				newTransaction.account_id,
				newTransaction.id,
				item.item_id,
				item.state,
				item.item_amount,
				item.item_already_done,
				{
					withCredentials: true
				}
			);

			if (res.status != 200) {
				error = "Une erreur s'est produite";
				setTimeout(() => {
					error = '';
				}, 1500);
				return;
			}
		}

		if (!error) {
			let res = await transactionsApi().patchTransactionId(
				newTransaction.account_id,
				newTransaction.id,
				'finished',
				{
					withCredentials: true
				}
			);

			if (res.status != 200) {
				error = "Une erreur s'est produite";
				setTimeout(() => {
					error = '';
				}, 1500);
				return;
			}

			transaction = newTransaction;
			success = 'Commande terminée';
			setTimeout(() => {
				success = '';
				searchName.set('');
				close();
			}, 1500);
		}
	}

	async function saveTransaction() {
		for (let i = 0; i < newTransaction.items.length; i++) {
			let item = newTransaction.items[i];

			// @ts-ignore
			if (item.state == transaction.items[i].state) item.state = undefined;

			let res = await transactionsApi().patchTransactionItemId(
				newTransaction.account_id,
				newTransaction.id,
				item.item_id,
				item.state,
				item.item_amount,
				item.item_already_done,
				{
					withCredentials: true
				}
			);

			if (res.status != 200) {
				error = "Une erreur s'est produite";
				setTimeout(() => {
					error = '';
				}, 1500);
				return;
			}
		}

		if (!error) {
			transaction = newTransaction;
			success = 'Changements enregistrée';
			setTimeout(() => {
				success = '';
			}, 1500);
			reloadTransaction();
		}
	}

	function reloadTransaction() {
		transactionsApi()
			.getTransactionId(transaction.account_id, transaction.id, { withCredentials: true })
			.then((res) => {
				transaction = res.data;
				newTransaction = structuredClone(transaction);
			});
	}

	function formatTimestampToReadableDate(timestamp: number): string {
        return new Date(timestamp).toLocaleTimeString('fr-FR', {
			day: '2-digit',
			month: "2-digit",
            hour: '2-digit',
            minute: '2-digit',
			second: '2-digit'
        });
    }

	let readableDate: string = formatTimestampToReadableDate(transaction.created_at*1000);
</script>

{#if success != ''}
	<Success message={success} />
{/if}

{#if error != ''}
	<Error {error} />
{/if}

<button
	id="overlay"
	class="fixed inset-0 bg-black bg-opacity-50 flex justify-center items-center z-50 cursor-default"
	on:click={() => {
		close();
	}}
/>

<div class="fixed inset-0 z-50 flex justify-center items-center pointer-events-none">
	<div class="popup-container pointer-events-auto">
		<div class="popup-header">
			<div class="header-left">
				{#if transaction.account_google_picture}
					<img
						src={transaction.account_google_picture}
						alt="Avatar"
						class="header-avatar"
						on:error={(e) => {
							const target = e.currentTarget;
							if (target instanceof HTMLImageElement) {
								target.style.display = 'none';
								target.nextElementSibling?.classList.remove('hidden');
							}
						}}
					/>
				{/if}
				<iconify-icon
					icon="mdi:account-circle"
					class="header-icon {transaction.account_google_picture ? 'hidden' : ''}"
				/>
				<span class="account-name">
					{transaction.account_nick_name || transaction.account_name}
					{#if transaction.account_nick_name && transaction.account_name && transaction.account_nick_name !== transaction.account_name}
						<span class="account-name-paren">({transaction.account_name})</span>
					{/if}
				</span>
			</div>
			<div class="header-right">
				<span class="transaction-date">{readableDate}</span>
				<button class="close-btn" on:click={close}>
					<iconify-icon icon="mdi:close" />
				</button>
			</div>
		</div>

		<div class="popup-body">
			<div class="items-grid">
				{#each newTransaction.items as item, i}
					<div
						class="item-card {item.state ==
						TransactionItemState.TransactionItemCanceled
							? 'canceled'
							: ''} {item.state == TransactionItemState.TransactionItemFinished
							? 'finished'
							: ''}"
					>
						{#if item.state == TransactionItemState.TransactionItemCanceled}
							<div class="item-status-badge canceled">
								<iconify-icon icon="mdi:cancel" />
							</div>
						{:else if item.state == TransactionItemState.TransactionItemFinished}
							<div class="item-status-badge finished">
								<iconify-icon icon="mdi:check-circle" />
							</div>
						{/if}

						<div class="item-main">
							{#if item.is_menu}
								<button
									class="item-image-container"
									on:click={() => {
										menuPopup = {
											items: item.menu_items,
											categories: item.menu_categories,
											pickedItems: item.picked_categories_items
										};
									}}
								>
									{#if item.picture_uri}
										<img
											src={api() + item.picture_uri}
											alt={item.item_name}
											class="item-image"
										/>
									{:else}
										<div class="item-image-placeholder">
											<iconify-icon icon="mdi:food-off" class="placeholder-icon" />
										</div>
									{/if}
									<iconify-icon icon="mdi:expand" class="menu-expand-icon" />
								</button>
							{:else}
								{#if item.picture_uri}
									<img
										src={api() + item.picture_uri}
										alt={item.item_name}
										class="item-image"
									/>
								{:else}
									<div class="item-image-placeholder">
										<iconify-icon icon="mdi:food-off" class="placeholder-icon" />
									</div>
								{/if}
							{/if}
							<div class="item-name">{item.item_name ? item.item_name : 'Test'}</div>
						</div>

						<div class="item-footer">
							<div class="item-amounts">
								<span class="amount-badge">
									Quantité: {item.item_already_done}/{item.item_amount}
								</span>
							</div>
						</div>

						{#if item.state == TransactionItemState.TransactionItemStarted}
							<div class="item-actions">
								<div class="action-buttons">
									{#if item.item_amount > 1 && item.item_already_done == 0}
										<button
											class="action-btn minus"
											on:click={() => {
												if (item.item_amount > 1) item.item_amount--;
											}}
										>
											<iconify-icon icon="mdi:minus" />
										</button>
									{/if}
									<button
										class="action-btn cancel"
										on:click={() => {
											item.state = TransactionItemState.TransactionItemCanceled;
										}}
									>
										<iconify-icon icon="mdi:close-circle" />
									</button>
									<button
										class="action-btn complete"
										on:click={() => {
											if (item.item_already_done < item.item_amount) item.item_already_done += 1;
											if (item.item_already_done == item.item_amount)
												item.state = TransactionItemState.TransactionItemFinished;
										}}
									>
										<iconify-icon icon="mdi:check-circle" />
									</button>
									{#if item.item_amount < transaction.items[i].item_amount}
										<button
											class="action-btn add"
											on:click={() => {
												if (item.item_amount < transaction.items[i].item_amount) item.item_amount++;
											}}
										>
											<iconify-icon icon="mdi:plus" />
										</button>
									{/if}
								</div>
							</div>
						{/if}
					</div>
				{/each}
			</div>
		</div>

		{#if menuPopup != undefined}
			<div class="menu-divider" />
			<div class="menu-header">
				<iconify-icon icon="mdi:food" class="menu-header-icon" />
				<span>Ce menu contient:</span>
			</div>
			<div class="menu-content">
				{#each menuPopup.categories ?? [] as cat}
					<div class="menu-item">
						{#if cat.picture_uri}
							<img
								src={api() + cat.picture_uri}
								alt={cat.name}
								class="menu-item-image"
							/>
						{:else}
							<div class="menu-item-placeholder">
								<iconify-icon icon="mdi:folder-outline" class="menu-placeholder-icon" />
							</div>
						{/if}
						<span class="menu-item-name">{cat.name ? cat.name : 'Test'}</span>
					</div>
				{/each}
				{#each menuPopup.items ?? [] as item}
					<div class="menu-item">
						{#if item.picture_uri}
							<img
								src={api() + item.picture_uri}
								alt={item.name}
								class="menu-item-image"
							/>
						{:else}
							<div class="menu-item-placeholder">
								<iconify-icon icon="mdi:food-off" class="menu-placeholder-icon" />
							</div>
						{/if}
						<span class="menu-item-name">{item.name ? item.name : 'Test'}</span>
					</div>
				{/each}
				{#each menuPopup.pickedItems ?? [] as pickedItem}
					<div class="menu-item selected">
						{#if pickedItem.picture_uri}
							<img
								src={api() + pickedItem.picture_uri}
								alt={pickedItem.item_name}
								class="menu-item-image"
							/>
						{:else}
							<div class="menu-item-placeholder">
								<iconify-icon icon="mdi:food-off" class="menu-placeholder-icon" />
							</div>
						{/if}
						<span class="menu-item-name">{pickedItem.item_name ? pickedItem.item_name : 'Test'}</span>
					</div>
				{/each}
			</div>
			<div class="menu-divider" />
		{/if}

		<div class="popup-footer">
			<div class="footer-left">
				<div class="transaction-type-badge">
					{#if transaction.is_remote}
						<iconify-icon icon="mdi:wifi" class="type-icon" />
						<span>En ligne</span>
					{:else}
						<iconify-icon icon="mdi:monitor" class="type-icon" />
						<span>Borne</span>
					{/if}
				</div>
			</div>
			<div class="price-display">
				<iconify-icon icon="mdi:cash" class="price-icon" />
				<span class="price-label">Total:</span>
				<span class="price-value">{formatPrice(transaction.total_cost)}</span>
			</div>
		</div>

		<div class="popup-actions">
			<div class="action-column primary-actions">
				<button
					class="action-main-btn save"
					disabled={transaction.state !== 'started'}
					on:click={saveTransaction}
				>
					<iconify-icon icon="mdi:content-save" />
					<span>Enregistrer</span>
				</button>
				<button
					class="action-main-btn finish"
					disabled={transaction.state !== 'started'}
					on:click={finishTransaction}
				>
					<iconify-icon icon="mdi:check-all" />
					<span>Terminer (paiement)</span>
				</button>
			</div>
			<div class="action-column secondary-actions">
				<button
					class="action-main-btn undo"
					disabled={transaction.state === 'started'}
					on:click={putBackTransaction}
				>
					<iconify-icon icon="mdi:arrow-u-left-top" />
					<span>Remettre en attente</span>
				</button>
				<button
					class="action-main-btn cancel"
					disabled={transaction.state !== 'started'}
					on:click={cancelTransaction}
				>
					<iconify-icon icon="mdi:cash-refund" />
					<span>Annuler (remboursement)</span>
				</button>
			</div>
		</div>
	</div>
</div>

<style>
	.fixed {
		position: fixed;
	}

	.inset-0 {
		top: 0;
		right: 0;
		bottom: 0;
		left: 0;
	}

	.z-50 {
		z-index: 50;
	}

	.flex {
		display: flex;
	}

	.justify-center {
		justify-content: center;
	}

	.items-center {
		align-items: center;
	}

	.cursor-default {
		cursor: default;
	}

	.bg-black {
		background-color: rgb(0 0 0);
	}

	.bg-opacity-50 {
		background-color: rgba(0, 0, 0, 0.5);
	}

	.pointer-events-none {
		pointer-events: none;
	}

	.pointer-events-auto {
		pointer-events: auto;
	}

	.popup-container {
		width: 80%;
		max-width: 1200px;
		max-height: 90vh;
		background-color: white;
		border-radius: 16px;
		box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
		overflow: hidden;
		display: flex;
		flex-direction: column;
	}

	@media (prefers-color-scheme: dark) {
		.popup-container {
			background-color: #2f2f41;
		}
	}

	.popup-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 20px 32px;
		border-bottom: 1px solid #e5e7eb;
		background-color: #f9fafb;
	}

	@media (prefers-color-scheme: dark) {
		.popup-header {
			background-color: #1f2937;
			border-bottom-color: #374151;
		}
	}

	.header-left {
		display: flex;
		align-items: center;
		gap: 12px;
	}

	.header-avatar {
		width: 48px;
		height: 48px;
		border-radius: 50%;
		object-fit: cover;
	}

	.header-icon {
		font-size: 40px;
		color: #374151;
	}

	.header-icon.hidden {
		display: none;
	}

	@media (prefers-color-scheme: dark) {
		.header-icon {
			color: #9ca3af;
		}
	}

	.account-name {
		font-size: 24px;
		font-weight: 600;
		color: #111827;
	}

	.account-name-paren {
		font-weight: 400;
		color: #6b7280;
	}

	@media (prefers-color-scheme: dark) {
		.account-name {
			color: white;
		}

		.account-name-paren {
			color: #9ca3af;
		}
	}

	.header-right {
		display: flex;
		align-items: center;
		gap: 16px;
	}

	.transaction-date {
		font-size: 16px;
		color: #6b7280;
		font-weight: 500;
	}

	.close-btn {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 44px;
		height: 44px;
		border-radius: 10px;
		border: none;
		background-color: #e5e7eb;
		color: #374151;
		cursor: pointer;
		transition: all 0.2s ease;
		font-size: 20px;
	}

	.close-btn:hover {
		background-color: #d1d5db;
	}

	@media (prefers-color-scheme: dark) {
		.close-btn {
			background-color: #374151;
			color: #d1d5db;
		}

		.close-btn:hover {
			background-color: #4b5563;
		}
	}

	.popup-body {
		padding: 24px;
		overflow-y: auto;
		flex: 1;
	}

	.items-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
		gap: 20px;
	}

	.item-card {
		display: flex;
		flex-direction: column;
		background-color: #f3f4f6;
		border-radius: 16px;
		padding: 20px;
		position: relative;
		transition: all 0.2s ease;
		min-height: 200px;
	}

	@media (prefers-color-scheme: dark) {
		.item-card {
			background-color: #374151;
		}
	}

	.item-card.canceled {
		background-color: #fee2e2;
	}

	@media (prefers-color-scheme: dark) {
		.item-card.canceled {
			background-color: #7f1d1d;
		}
	}

	.item-card.finished {
		background-color: #dcfce7;
	}

	@media (prefers-color-scheme: dark) {
		.item-card.finished {
			background-color: #14532d;
		}
	}

	.item-status-badge {
		position: absolute;
		top: 12px;
		right: 12px;
		width: 28px;
		height: 28px;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 16px;
	}

	.item-status-badge.canceled {
		background-color: #ef4444;
		color: white;
	}

	.item-status-badge.finished {
		background-color: #22c55e;
		color: white;
	}

	.item-main {
		display: flex;
		flex-direction: column;
		align-items: center;
		text-align: center;
	}

	.item-image-container {
		position: relative;
		background: none;
		border: none;
		padding: 0;
		cursor: pointer;
	}

	.item-image {
		width: 64px;
		height: 64px;
		border-radius: 32px;
		object-fit: cover;
		align-self: center;
	}

	.item-image-placeholder {
		width: 64px;
		height: 64px;
		border-radius: 32px;
		background-color: #e5e7eb;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	@media (prefers-color-scheme: dark) {
		.item-image-placeholder {
			background-color: #4b5563;
		}
	}

	.item-image-placeholder .placeholder-icon {
		font-size: 32px;
		color: #6b7280;
	}

	@media (prefers-color-scheme: dark) {
		.item-image-placeholder .placeholder-icon {
			color: #9ca3af;
		}
	}

	.menu-expand-icon {
		position: absolute;
		bottom: -4px;
		right: -4px;
		background-color: #2563eb;
		color: white;
		border-radius: 50%;
		width: 20px;
		height: 20px;
		font-size: 12px;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.item-name {
		font-size: 18px;
		font-weight: 600;
		color: #111827;
		margin-top: 12px;
		text-align: center;
		line-height: 1.2;
	}

	@media (prefers-color-scheme: dark) {
		.item-name {
			color: white;
		}
	}

	.item-footer {
		margin-top: 12px;
	}

	.item-amounts {
		display: flex;
		flex-direction: column;
		gap: 6px;
	}

	.amount-badge {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 6px;
		font-size: 14px;
		color: #374151;
		background-color: rgba(255, 255, 255, 0.7);
		padding: 6px 12px;
		border-radius: 12px;
	}

	@media (prefers-color-scheme: dark) {
		.amount-badge {
			background-color: rgba(0, 0, 0, 0.2);
			color: #d1d5db;
		}
	}

	.amount-badge.pending {
		color: #2563eb;
		font-weight: 600;
	}

	.item-actions {
		margin-top: 16px;
		padding-top: 16px;
		border-top: 1px solid rgba(0, 0, 0, 0.1);
	}

	@media (prefers-color-scheme: dark) {
		.item-actions {
			border-top-color: rgba(255, 255, 255, 0.1);
		}
	}

	.action-buttons {
		display: flex;
		justify-content: center;
		gap: 12px;
	}

	.action-btn {
		width: 44px;
		height: 44px;
		border-radius: 10px;
		border: none;
		display: flex;
		align-items: center;
		justify-content: center;
		cursor: pointer;
		transition: all 0.2s ease;
		font-size: 22px;
	}

	.action-btn.minus {
		background-color: #f97316;
		color: white;
	}

	.action-btn.minus:hover {
		background-color: #ea580c;
	}

	.action-btn.cancel {
		background-color: #6b7280;
		color: white;
	}

	.action-btn.cancel:hover {
		background-color: #4b5563;
	}

	.action-btn.complete {
		background-color: #22c55e;
		color: white;
	}

	.action-btn.complete:hover {
		background-color: #16a34a;
	}

	.action-btn.add {
		background-color: #2563eb;
		color: white;
	}

	.action-btn.add:hover {
		background-color: #1d4ed8;
	}

	.menu-divider {
		height: 1px;
		background-color: #e5e7eb;
	}

	@media (prefers-color-scheme: dark) {
		.menu-divider {
			background-color: #374151;
		}
	}

	.menu-header {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 8px;
		padding: 12px;
		font-weight: 600;
		color: #374151;
		font-size: 16px;
	}

	@media (prefers-color-scheme: dark) {
		.menu-header {
			color: #d1d5db;
		}
	}

	.menu-header-icon {
		font-size: 20px;
		color: #f97316;
	}

	.menu-content {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
		gap: 16px;
		padding: 16px 24px;
		max-height: 240px;
		overflow-y: auto;
	}

	.menu-item {
		display: flex;
		flex-direction: column;
		align-items: center;
		text-align: center;
		padding: 12px;
		background-color: #f3f4f6;
		border-radius: 10px;
	}

	@media (prefers-color-scheme: dark) {
		.menu-item {
			background-color: #374151;
		}
	}

	.menu-item.selected {
		background-color: #dbeafe;
	}

	@media (prefers-color-scheme: dark) {
		.menu-item.selected {
			background-color: #1e3a8a;
		}
	}

	.menu-item-image {
		width: 40px;
		height: 40px;
		border-radius: 20px;
		object-fit: cover;
	}

	.menu-item-placeholder {
		width: 40px;
		height: 40px;
		border-radius: 20px;
		background-color: #e5e7eb;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	@media (prefers-color-scheme: dark) {
		.menu-item-placeholder {
			background-color: #4b5563;
		}
	}

	.menu-placeholder-icon {
		font-size: 22px;
		color: #6b7280;
	}

	@media (prefers-color-scheme: dark) {
		.menu-placeholder-icon {
			color: #9ca3af;
		}
	}

	.menu-item-name {
		font-size: 14px;
		color: #374151;
		margin-top: 6px;
	}

	@media (prefers-color-scheme: dark) {
		.menu-item-name {
			color: #d1d5db;
		}
	}

	.popup-footer {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 20px 32px;
		border-top: 1px solid #e5e7eb;
		background-color: #f9fafb;
	}

	@media (prefers-color-scheme: dark) {
		.popup-footer {
			background-color: #1f2937;
			border-top-color: #374151;
		}
	}

	.footer-left {
		display: flex;
		align-items: center;
	}

	.transaction-type-badge {
		display: flex;
		align-items: center;
		gap: 8px;
		padding: 10px 20px;
		background-color: #e5e7eb;
		border-radius: 10px;
		font-size: 16px;
		font-weight: 600;
		color: #374151;
	}

	.transaction-type-badge .type-icon {
		font-size: 22px;
	}

	@media (prefers-color-scheme: dark) {
		.transaction-type-badge {
			background-color: #374151;
			color: #d1d5db;
		}
	}

	.price-display {
		display: flex;
		align-items: center;
		gap: 12px;
	}

	.price-icon {
		font-size: 32px;
		color: #22c55e;
	}

	.price-label {
		font-size: 20px;
		font-weight: 500;
		color: #6b7280;
	}

	.price-value {
		font-size: 28px;
		font-weight: 700;
		color: #111827;
	}

	@media (prefers-color-scheme: dark) {
		.price-value {
			color: white;
		}
	}

	.action-group {
		display: flex;
		gap: 12px;
	}

	.footer-btn {
		display: flex;
		align-items: center;
		gap: 6px;
		padding: 10px 16px;
		border-radius: 8px;
		border: none;
		font-size: 14px;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.footer-btn.secondary {
		background-color: #e5e7eb;
		color: #374151;
	}

	.footer-btn.secondary:hover {
		background-color: #d1d5db;
	}

	@media (prefers-color-scheme: dark) {
		.footer-btn.secondary {
			background-color: #374151;
			color: #d1d5db;
		}

		.footer-btn.secondary:hover {
			background-color: #4b5563;
		}
	}

	.popup-actions {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 16px;
		padding: 16px 24px;
		border-top: 1px solid #e5e7eb;
		background-color: #f9fafb;
	}

	@media (prefers-color-scheme: dark) {
		.popup-actions {
			background-color: #1f2937;
			border-top-color: #374151;
		}
	}

	.action-column {
		display: flex;
		flex-direction: column;
		gap: 12px;
	}

	.action-main-btn {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 8px;
		padding: 14px 20px;
		border-radius: 10px;
		border: none;
		font-size: 15px;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s ease;
		color: white;
	}

	.action-main-btn:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.action-main-btn:disabled:hover {
		transform: none;
	}

	.action-main-btn.save {
		background-color: #2563eb;
	}

	.action-main-btn.save:hover:not(:disabled) {
		background-color: #1d4ed8;
	}

	.action-main-btn.finish {
		background-color: #22c55e;
	}

	.action-main-btn.finish:hover:not(:disabled) {
		background-color: #16a34a;
	}

	.action-main-btn.undo {
		background-color: #f97316;
	}

	.action-main-btn.undo:hover:not(:disabled) {
		background-color: #ea580c;
	}

	.action-main-btn.cancel {
		background-color: #ef4444;
	}

	.action-main-btn.cancel:hover:not(:disabled) {
		background-color: #dc2626;
	}

	@media (max-width: 640px) {
		.popup-container {
			width: 95%;
			max-height: 95vh;
		}

		.popup-header {
			flex-direction: column;
			gap: 8px;
			text-align: center;
		}

		.header-left {
			justify-content: center;
		}

		.header-right {
			width: 100%;
			justify-content: center;
		}

		.items-grid {
			grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
			gap: 12px;
		}

		.popup-footer {
			flex-direction: column;
			gap: 12px;
		}

		.footer-left {
			width: 100%;
			justify-content: center;
		}

		.popup-actions {
			grid-template-columns: 1fr;
		}

		.action-column {
			flex-direction: row;
		}

		.action-main-btn {
			flex: 1;
			padding: 12px 16px;
			font-size: 13px;
		}
	}
</style>